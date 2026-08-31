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
	"runtime"
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
	shards       []*shardEntry[K, V]
	shardCount   int
	mask         int            // shardCount-1，用于位运算取模（shardCount 必须是 2 的幂）
	hasher       func(K) uint32 // key 的 hash 函数
	count        atomic.Int64   // 元素总数（原子计数，零锁开销）
	perShardHint int            // 每 shard 预分配容量提示（Clear 时复用）
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
		shards:       shards,
		shardCount:   shardCount,
		mask:         shardCount - 1,
		hasher:       KvHasher[K](),
		perShardHint: cfg.perShardHint,
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

// Swap 替换 key 的值并返回之前的值
// 如果 key 不存在，存储新值并返回零值和 false
func (m *ShardedMap[K, V]) Swap(key K, value V) (V, bool) {
	shard := m.getShard(key)
	shard.mu.Lock()
	old, exists := shard.data[key]
	shard.data[key] = value
	shard.mu.Unlock()
	if !exists {
		m.count.Add(1)
	}
	return old, exists
}

// CompareAndSwap 原子比较并交换
// 仅当 key 存在且当前值等于 old 时，才替换为 new
// 注意：V 为非可比较类型（如 slice/map）时会 panic，与 sync.Map 行为一致
func (m *ShardedMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	current, exists := shard.data[key]
	if !exists {
		return false
	}
	if any(current) == any(old) {
		shard.data[key] = new
		return true
	}
	return false
}

// CompareAndDelete 原子比较并删除
// 仅当 key 存在且当前值等于 value 时，才删除
// 注意：V 为非可比较类型（如 slice/map）时会 panic，与 sync.Map 行为一致
func (m *ShardedMap[K, V]) CompareAndDelete(key K, value V) bool {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	current, exists := shard.data[key]
	if !exists {
		return false
	}
	if any(current) == any(value) {
		delete(shard.data, key)
		m.count.Add(-1)
		return true
	}
	return false
}

// ============================================================================
// 批量操作
// ============================================================================

// StoreBatch 批量存储，按 shard 分组减少锁获取次数
// 相比逐个 Store，N 个元素分散到 S 个 shard 时锁获取次数从 N 降到 S
func (m *ShardedMap[K, V]) StoreBatch(items map[K]V) {
	if len(items) == 0 {
		return
	}

	// 按 shard 索引分组
	groups := make(map[int]map[K]V, m.shardCount)
	for k, v := range items {
		idx := int(m.hasher(k)) & m.mask
		if groups[idx] == nil {
			groups[idx] = make(map[K]V)
		}
		groups[idx][k] = v
	}

	// 每个 shard 只获取一次写锁
	for idx, group := range groups {
		shard := m.shards[idx]
		shard.mu.Lock()
		added := 0
		for k, v := range group {
			if _, exists := shard.data[k]; !exists {
				added++
			}
			shard.data[k] = v
		}
		shard.mu.Unlock()
		if added > 0 {
			m.count.Add(int64(added))
		}
	}
}

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

// RangeParallel 并行遍历所有键值对（百万级数据优化）
//
// 将 shards 分配到 workerNum 个 goroutine 并行遍历，每个 goroutine 独立持有所分配 shard 的读锁。
// 不同 goroutine 处理不同 shard，读锁无竞争；同一 shard 仅被一个 goroutine 读取，零锁争用。
//
// 对比 Range（串行遍历）：
//   - N 核 CPU 上百万级数据遍历提速约 N 倍
//   - 零额外分配（直接遍历 shards 数组，不需要 Keys() 中间切片）
//
// 注意：
//   - fn 无返回值（始终遍历全部），不支持提前终止——并行场景下"全局停止"语义不明确
//   - fn 在持有 shard 读锁时执行，应为非阻塞操作（TrySend 等非阻塞操作可安全调用）
//   - workerNum <= 0 时使用 GOMAXPROCS；shard 数量少于 worker 时退化为按 shard 数并行
func (m *ShardedMap[K, V]) RangeParallel(workerNum int, fn func(key K, value V)) {
	if fn == nil {
		return
	}
	shardCount := len(m.shards)
	if shardCount == 0 {
		return
	}
	if workerNum <= 0 {
		workerNum = runtime.GOMAXPROCS(0)
	}
	if workerNum < 1 {
		workerNum = 1
	}
	// shard 数量少于 worker 时，每个 shard 一个 goroutine 即可
	if workerNum > shardCount {
		workerNum = shardCount
	}

	// 每个 worker 负责一段连续的 shard（分区方式，避免 channel 同步开销）
	chunkSize := (shardCount + workerNum - 1) / workerNum
	var wg sync.WaitGroup
	for i := 0; i < workerNum; i++ {
		start := i * chunkSize
		if start >= shardCount {
			break
		}
		end := start + chunkSize
		if end > shardCount {
			end = shardCount
		}
		wg.Add(1)
		go func(shards []*shardEntry[K, V]) {
			defer wg.Done()
			for _, shard := range shards {
				shard.mu.RLock()
				for k, v := range shard.data {
					fn(k, v)
				}
				shard.mu.RUnlock()
			}
		}(m.shards[start:end])
	}
	wg.Wait()
}

// Len 返回元素总数（原子读取，零锁开销）
func (m *ShardedMap[K, V]) Len() int {
	return int(m.count.Load())
}

// Clear 清空所有元素，保留预分配容量提示
func (m *ShardedMap[K, V]) Clear() {
	for _, shard := range m.shards {
		shard.mu.Lock()
		if m.perShardHint > 0 {
			shard.data = make(map[K]V, m.perShardHint)
		} else {
			shard.data = make(map[K]V)
		}
		shard.mu.Unlock()
	}
	m.count.Store(0)
}

// ClearHalf 随机清空每个 shard 内一半的元素
//
// 用于大容量缓存的超限驱逐：相比整表 Clear 保留一半数据，避免缓存命中率雪崩
// （重建窗口内的查询全部回源），驱逐成本相同（无 LRU 链表维护开销）
// 每 shard 独立驱逐一半（map 遍历序随机），无轮换状态、多次调用持续有效
func (m *ShardedMap[K, V]) ClearHalf() {
	var removed int64
	for _, shard := range m.shards {
		shard.mu.Lock()
		n := len(shard.data)
		if n > 0 {
			evict := n / 2
			if evict < 1 {
				evict = 1
			}
			i := 0
			for k := range shard.data {
				delete(shard.data, k)
				i++
				if i >= evict {
					break
				}
			}
			removed += int64(i)
		}
		shard.mu.Unlock()
	}
	m.count.Add(-removed)
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
// 内联实现，零分配（对比 fnv.New32a()+Write 消除 heap alloc 和 []byte 转换）
func FNVHashString32(s string) uint32 {
	const (
		offsetBasis uint32 = 2166136261
		prime       uint32 = 16777619
	)
	h := offsetBasis
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= prime
	}
	return h
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
