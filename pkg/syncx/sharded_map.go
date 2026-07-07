/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-01 23:51:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-07 00:56:00
 * @FilePath: \go-toolbox\pkg\syncx\sharded_map.go
 * @Description: 泛型分片映射表
 *
 * 将 key 按 hash 分散到 N 个 shard，每个 shard 独立 RWMutex
 * 适用于高并发读写场景，相比 sync.Map 提供更好的写性能和 Range 性能
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package syncx

import (
	"fmt"
	"hash/fnv"
	"sync"
	"sync/atomic"
)

// ============================================================================
// 分片映射表
// ============================================================================

// ShardedMap 分片映射表
// 将 key 按 hash 分散到 N 个 shard，每个 shard 独立锁
// 适用于高并发读写场景，相比 sync.Map：
//   - 写性能更好（多 shard 分散锁竞争）
//   - Range 性能更好（可并行遍历不同 shard）
//   - Len 性能更好（原子计数，无需遍历）
type ShardedMap[K comparable, V any] struct {
	shards     []*shardEntry[K, V]
	shardCount int
	mask       int            // shardCount-1，用于位运算取模（shardCount 必须是 2 的幂）
	hasher     func(K) uint32 // key 的 hash 函数
	count      atomic.Int64   // 元素总数（原子计数，零锁开销）
}

// shardEntry 单个分片
type shardEntry[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewShardedMap 创建分片映射表（无预分配容量）
//
// 参数：
//   - shardCount: 分片数量，必须是 2 的幂（如 32/64/128）
//
// 返回：*ShardedMap[K, V]
//
// 注意：每个 shard 内部 map 不预分配容量，适用于容量未知或较小的场景
// 已知总容量的大数据场景请使用 NewShardedMapWithOptions 配合 WithPerShardHint
func NewShardedMap[K comparable, V any](shardCount int) *ShardedMap[K, V] {
	return NewShardedMapWithOptions[K, V](shardCount)
}

// ShardedMapOption ShardedMap 配置选项（修改内部 config）
type ShardedMapOption[K comparable, V any] func(*shardedMapConfig)

// shardedMapConfig ShardedMap 初始化配置（私有，避免外部直接修改）
type shardedMapConfig struct {
	// perShardHint 每个 shard 内部 map 的预分配容量提示
	// 用于已知总容量场景，减少 map 扩容次数，提升写入性能
	perShardHint int
}

// WithPerShardHint 设置每个 shard 内部 map 的预分配容量提示
//
// 用于已知总容量的场景，减少 shard 内部 map 扩容次数，提升写入性能
// 建议值：预估总元素数 / shardCount（向上取整）
//
// 示例：总容量 10000，分 64 个 shard，每 shard 提示 = 10000/64 ≈ 157
func WithPerShardHint[K comparable, V any](perShardHint int) ShardedMapOption[K, V] {
	if perShardHint < 0 {
		perShardHint = 0
	}
	return func(cfg *shardedMapConfig) {
		cfg.perShardHint = perShardHint
	}
}

// NewShardedMapWithOptions 创建分片映射表（支持配置选项）
//
// 参数：
//   - shardCount: 分片数量，必须是 2 的幂（如 32/64/128）
//   - opts: 配置选项（如 WithPerShardHint）
//
// 返回：*ShardedMap[K, V]
func NewShardedMapWithOptions[K comparable, V any](shardCount int, opts ...ShardedMapOption[K, V]) *ShardedMap[K, V] {
	// 解析选项
	cfg := shardedMapConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	// 确保 shardCount 是 2 的幂
	if shardCount <= 0 {
		shardCount = 64
	}
	if shardCount&(shardCount-1) != 0 {
		// 不是 2 的幂，向上取最近的 2 的幂
		shardCount = NextPowerOfTwo(shardCount)
	}

	shards := make([]*shardEntry[K, V], shardCount)
	for i := range shards {
		if cfg.perShardHint > 0 {
			shards[i] = &shardEntry[K, V]{
				data: make(map[K]V, cfg.perShardHint),
			}
		} else {
			shards[i] = &shardEntry[K, V]{
				data: make(map[K]V),
			}
		}
	}

	return &ShardedMap[K, V]{
		shards:     shards,
		shardCount: shardCount,
		mask:       shardCount - 1,
		hasher:     KvHasher[K](),
	}
}

// ============================================================================
// 基础操作
// ============================================================================

// Store 存储 key→value
func (m *ShardedMap[K, V]) Store(key K, value V) {
	shard := m.getShard(key)
	shard.mu.Lock()
	_, exists := shard.data[key]
	shard.data[key] = value
	shard.mu.Unlock()
	if !exists {
		m.count.Add(1)
	}
}

// Load 加载 key 的 value
// 返回：(value, exists)
func (m *ShardedMap[K, V]) Load(key K) (V, bool) {
	shard := m.getShard(key)
	shard.mu.RLock()
	value, exists := shard.data[key]
	shard.mu.RUnlock()
	return value, exists
}

// Delete 删除 key
func (m *ShardedMap[K, V]) Delete(key K) {
	shard := m.getShard(key)
	shard.mu.Lock()
	_, exists := shard.data[key]
	if exists {
		delete(shard.data, key)
	}
	shard.mu.Unlock()
	if exists {
		m.count.Add(-1)
	}
}

// LoadAndDelete 加载并删除
// 返回：(value, exists)
func (m *ShardedMap[K, V]) LoadAndDelete(key K) (V, bool) {
	shard := m.getShard(key)
	shard.mu.Lock()
	value, exists := shard.data[key]
	if exists {
		delete(shard.data, key)
	}
	shard.mu.Unlock()
	if exists {
		m.count.Add(-1)
	}
	return value, exists
}

// LoadOrStore 加载或存储
// 如果 key 存在返回 (existing, true)，否则存储 value 并返回 (value, false)
func (m *ShardedMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	shard := m.getShard(key)
	shard.mu.Lock()
	existing, exists := shard.data[key]
	if !exists {
		shard.data[key] = value
	}
	shard.mu.Unlock()
	if !exists {
		m.count.Add(1)
		return value, false
	}
	return existing, true
}

// Has 检查 key 是否存在
func (m *ShardedMap[K, V]) Has(key K) bool {
	_, exists := m.Load(key)
	return exists
}

// ============================================================================
// 批量操作
// ============================================================================

// Range 遍历所有键值对
// 遍历期间每个 shard 持有读锁（分片粒度，不影响其他 shard 写入）
// fn 返回 false 时停止遍历
func (m *ShardedMap[K, V]) Range(fn func(key K, value V) bool) {
	for _, shard := range m.shards {
		shard.mu.RLock()
		stop := false
		for k, v := range shard.data {
			if !fn(k, v) {
				stop = true
				break
			}
		}
		shard.mu.RUnlock()
		if stop {
			return
		}
	}
}

// Len 返回元素总数（原子读取，零锁开销）
func (m *ShardedMap[K, V]) Len() int {
	return int(m.count.Load())
}

// Clear 清空所有元素
func (m *ShardedMap[K, V]) Clear() {
	for _, shard := range m.shards {
		shard.mu.Lock()
		shard.data = make(map[K]V)
		shard.mu.Unlock()
	}
	m.count.Store(0)
}

// Count 返回满足条件的元素数量
// filter 为 nil 时等价于 Len
func (m *ShardedMap[K, V]) Count(filter func(K, V) bool) int {
	if filter == nil {
		return m.Len()
	}
	var count int
	m.Range(func(k K, v V) bool {
		if filter(k, v) {
			count++
		}
		return true
	})
	return count
}

// Keys 返回所有 key 的切片
func (m *ShardedMap[K, V]) Keys() []K {
	result := make([]K, 0, m.Len())
	m.Range(func(k K, _ V) bool {
		result = append(result, k)
		return true
	})
	return result
}

// Values 返回所有 value 的切片
func (m *ShardedMap[K, V]) Values() []V {
	result := make([]V, 0, m.Len())
	m.Range(func(_ K, v V) bool {
		result = append(result, v)
		return true
	})
	return result
}

// ============================================================================
// 分片级操作（高级 API，用于需要跨索引原子性的场景）
// ============================================================================

// WithShardLock 在 key 对应的 shard 锁内执行操作（写锁）
// 用于需要在同一 shard 内原子操作多个 map 的场景
func (m *ShardedMap[K, V]) WithShardLock(key K, fn func(shardData map[K]V)) {
	shard := m.getShard(key)
	shard.mu.Lock()
	fn(shard.data)
	shard.mu.Unlock()
}

// WithShardRLock 在 key 对应的 shard 读锁内执行操作（读锁）
func (m *ShardedMap[K, V]) WithShardRLock(key K, fn func(shardData map[K]V)) {
	shard := m.getShard(key)
	shard.mu.RLock()
	fn(shard.data)
	shard.mu.RUnlock()
}

// ============================================================================
// 内部方法
// ============================================================================

// getShard 根据 key 获取对应的 shard
func (m *ShardedMap[K, V]) getShard(key K) *shardEntry[K, V] {
	h := m.hasher(key)
	return m.shards[int(h)&m.mask]
}

// NextPowerOfTwo 返回不小于 n 的最小的 2 的幂
func NextPowerOfTwo(n int) int {
	if n <= 1 {
		return 2
	}
	// 如果 n 已经是 2 的幂，返回下一个
	if n&(n-1) == 0 {
		return n << 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n + 1
}

// FNVHashString32 使用 FNV-1a 算法计算 string 的 32 位 hash
func FNVHashString32(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// KvHasher 为 ShardedMap 的 K 类型选择最优的 hash 函数
// 通过类型断言为常见 K 类型（string/int/int64 等）选择专用 hasher，避免反射开销
func KvHasher[K comparable]() func(K) uint32 {
	switch any(*new(K)).(type) {
	case string:
		return func(k K) uint32 { return FNVHashString32(any(k).(string)) }
	case int:
		return func(k K) uint32 { return uint32(any(k).(int)) }
	case int64:
		return func(k K) uint32 { return uint32(any(k).(int64)) }
	case int32:
		return func(k K) uint32 { return uint32(any(k).(int32)) }
	case uint:
		return func(k K) uint32 { return uint32(any(k).(uint)) }
	case uint64:
		return func(k K) uint32 { return uint32(any(k).(uint64)) }
	default:
		// 其他类型用 fmt.Sprintf 转 string 再 FNV-1a hash
		return func(k K) uint32 { return FNVHashString32(fmt.Sprintf("%v", k)) }
	}
}
