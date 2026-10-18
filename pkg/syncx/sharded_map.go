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
//
// 分组用计数排序而不是 map 嵌 map
//
// 原先为每个 shard 建一个 map[K]V 来装它的那批键，一千个元素分到 64 个 shard
// 要分配 451 次；而 shard 下标本就是 [0, shardCount) 的小整数，
// 计数排序两趟即可把元素按 shard 排好，总共三次分配
//
// 为什么保留分组而不是干脆逐项 Store：量过两种写法，结论随 shard 数翻转——
// 64 个 shard 时逐项更快（争用小，省下的是建表开销），
// 但 8 个 shard 并发写时逐项反而慢一倍（十个 goroutine 抢八把锁，争用占主导，一千次加锁远贵于一次）
// 分组的价值在争用大的那一侧，去掉分配之后两侧都不亏
//
// [EN] Grouping is done with a counting sort rather than a map of maps. The
// previous version built one map[K]V per shard, which cost 451 allocations for a
// thousand elements across 64 shards, even though a shard index is just a small
// integer in [0, shardCount) that two passes can order by. Grouping itself is
// kept rather than falling back to per-item Store because the two were measured
// and the winner flips with the shard count: per-item is faster with 64 shards,
// where contention is low and the saving is the bookkeeping, but half as fast
// with 8 shards under concurrent writers, where ten goroutines contend for eight
// locks and a thousand acquisitions cost far more than one. Grouping earns its
// keep on the contended side, and without the allocations neither side loses
// StoreBatch 批量写入
//
// 元素较多时按 shard 分组，每个 shard 只加锁一次，把加锁次数从 len(items) 降到
// 最多 shardCount 次分组用计数排序，两趟线性扫描，不做比较排序
//
// 分组并非总是更快，阈值取 shardCount
//
// 分组省下来的是「多个元素落进同一个 shard」时的重复加锁，所以只有
// len(items) 明显大于 shardCount 时才有的省元素更少时多数 shard 只摊到一个
// 元素，省不下任何一次加锁，却要先付三次切片分配（entries/counts/sorted）
// 实测（Apple M4，10 并发，ns/op，越小越好）：
//
//	元素数   8 分片 逐项/分组     64 分片 逐项/分组
//	8        426 / 375           290 / 396
//	16       917 / 549           539 / 632
//	32       1829 / 848          913 / 1115
//	64       4029 / 1498         1897 / 1747   ← 64 分片的交点
//	512      34825 / 11060       15759 / 10993
//
// 两种分片数的交点都落在 len(items) ≈ shardCount 上，与上面的机理一致
// 无争用的单协程场景加锁只是几纳秒的原子操作，分组要到约 512 项才扳平，
// 因此这个阈值偏向争用场景——分片 map 本来就是为争用而存在的
//
// [EN] Bulk store. With enough items it groups them by shard with a counting
// sort so each shard is locked once, cutting lock acquisitions from len(items)
// to at most shardCount.
//
// Grouping is not always faster, and the threshold is shardCount. What grouping
// saves is the repeated locking when several items land in the same shard, so it
// only pays off once len(items) clearly exceeds shardCount; below that most
// shards hold a single item, no lock acquisition is saved, and three slice
// allocations are spent for nothing. Measured on an Apple M4 at 10-way
// parallelism, both shard counts cross over at len(items) ~ shardCount. With no
// contention at all grouping needs about 512 items to break even, so this
// threshold favours the contended case — which is what a sharded map is for
func (m *ShardedMap[K, V]) StoreBatch(items map[K]V) {
	if len(items) == 0 {
		return
	}
	if len(items) <= m.shardCount {
		m.storeBatchSmall(items)
		return
	}
	m.storeBatchGrouped(items)
}

func (m *ShardedMap[K, V]) storeBatchSmall(items map[K]V) {
	for k, v := range items {
		shard := m.getShard(k)
		shard.mu.Lock()
		_, exists := shard.data[k]
		shard.data[k] = v
		shard.mu.Unlock()
		if !exists {
			m.count.Add(1)
		}
	}
}

// storeBatchGrouped 按 shard 分组后批量写入，每个 shard 只加锁一次
//
// [EN] Groups by shard, then locks each shard once
func (m *ShardedMap[K, V]) storeBatchGrouped(items map[K]V) {
	type entry struct {
		key K
		val V
		idx int32
	}

	// 第一趟：算出每个元素的 shard 下标，并统计每个 shard 的元素个数
	entries := make([]entry, 0, len(items))
	counts := make([]int32, m.shardCount+1)
	for k, v := range items {
		idx := int32(int(m.hasher(k)) & m.mask)
		entries = append(entries, entry{key: k, val: v, idx: idx})
		counts[idx+1]++
	}

	// 前缀和：counts[i] 变为 shard i 在 sorted 中的起始位置
	for i := 1; i <= m.shardCount; i++ {
		counts[i] += counts[i-1]
	}

	// 第二趟：按 shard 下标把元素放到各自的区段里
	sorted := make([]entry, len(entries))
	offsets := counts[:m.shardCount]
	for _, e := range entries {
		sorted[offsets[e.idx]] = e
		offsets[e.idx]++
	}

	// offsets 已被推到各段末尾，正好可以用作区段的右边界
	begin := int32(0)
	for i := 0; i < m.shardCount; i++ {
		stop := offsets[i]
		if begin == stop {
			continue // 该 shard 没有元素
		}
		shard := m.shards[i]
		shard.mu.Lock()
		added := 0
		for _, e := range sorted[begin:stop] {
			if _, exists := shard.data[e.key]; !exists {
				added++
			}
			shard.data[e.key] = e.val
		}
		shard.mu.Unlock()
		if added > 0 {
			m.count.Add(int64(added))
		}
		begin = stop
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
// 对比 Range（串行遍历）：百万级数据提速数倍，不需要 Keys() 那样的中间切片
//
// 实测（Apple M4，10 核，一百万个 int 键值对，64 shard）：
//
//	回调内容                        Range        RangeParallel
//	本地累加（无争用）               5.91 ms      1.56 ms      快 3.8 倍
//	atomic.Int64.Add（共享计数器）    6.17 ms      44.2 ms      慢 7.2 倍
//
// **fn 里的争用会把并行的收益吃光并倒转**：十个 goroutine 抢同一个原子计数器时，
// 缓存行在核间来回弹跳的代价远超并行省下的时间，比串行还慢七倍
// 想累加就让每个 goroutine 写各自的槽位（例如按 key 分桶），遍历结束再合并；
// 往同一个 channel 发、抢同一把锁、累加同一个计数器都属于这一类
// 换句话说「提速数倍」的前提是 fn 之间互不干扰，而这一条比并行本身更决定结果
//
// [EN] Several times faster than Range on million-scale data, with no
// intermediate slice. Measured on an Apple M4 with 10 cores and one million int
// pairs across 64 shards: a contention-free callback runs 5.91 ms serially
// against 1.56 ms in parallel, a 3.8x gain, while a callback adding to one
// shared atomic.Int64 runs 6.17 ms serially against 44.2 ms in parallel — 7.2x
// *slower*. Contention inside fn does not merely erase the gain, it inverts it:
// ten goroutines bouncing one cache line between cores cost far more than
// parallelism saves. Accumulate into per-goroutine slots and merge afterwards;
// sending to one channel, taking one lock or incrementing one counter all fall
// into the same trap. The promised speedup is conditional on fn not
// interfering with itself, and that condition matters more than the parallelism
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

// getShard 根据 key 获取对应的 shard，分片号取键 hash 的低位
//
// 为什么不改成「按负载调度」：不把分片号定成键的纯函数，而是存一张
// key -> 分片号 的目录表、每次往最闲的分片放——这个方案行不通
//
// Load 必须能按键找回数据，而目录表覆盖全部键、每次操作都要先查它，它本身就是
// 分片要消除的那个单点（目录表再分片的话又得按键 hash，绕回原点）
// 实测一万个键、64 分片、读取：
//
//	           hash 分片   调度+目录表
//	串行       17.8 ns     20.6 ns     慢 1.16 倍
//	并发 10    13.4 ns     96.9 ns     慢 7.2 倍
//
// 串行几乎免费、并发差一个数量级，正是争用瓶颈的签名目录表那把读写锁把所有
// 读都串行化了此外每个键要存两份，写入还要再上一次写锁
// 而收益一侧：位混合后最大分片 159 对理想 156，不均只剩约 2%
//
// least-loaded 调度在「不需要按键找回」的场景才成立，比如任务派发
// （WorkerPool / BatchProcessor），不是键值存储
//
// [EN] Picks a shard from the key's hash. Not load-balanced placement: keeping a
// key -> shard directory and placing each key in the least loaded shard cannot
// work, because Load must find a key again, so the directory covers every key
// and is consulted on every operation — it is itself the single point sharding
// exists to remove (and sharding the directory means hashing the key again).
// Measured over ten thousand keys across 64 shards, reads cost 17.8 ns against
// 20.6 ns serially but 13.4 ns against 96.9 ns at 10-way parallelism: nearly
// free alone, an order of magnitude worse under contention — the signature of a
// contention bottleneck. It also stores every key twice. Against that, mixing
// already leaves only about 2% imbalance. Least-loaded dispatch belongs where
// nothing has to be found again by key, such as task dispatch, not in a map
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

// mixUint64 把 64 位整数的全部比特混进低位，供掩码取分片用
//
// 乘上 64 位黄金分割常数后取高 32 位（multiply-shift）
//
// 必须取高位，不能取低位乘法的进位只向左传，乘积的第 i 位只取决于输入的
// 第 0..i 位，所以乘积低位拿不到任何高位信息（常数为奇数时最低位就等于输入
// 最低位原样）实测一万个键、64 分片，取低位和完全不混合一样坏：
//
//	键步长    取高 32 位          取低 32 位
//	1         64/64 分片, 最大 159   64/64, 157
//	8         64/64, 159          8/64, 1250
//	64        64/64, 159          1/64, 10000
//	1000      64/64, 159          8/64, 1250
//
// 不按平台位宽分叉64 位乘法在 amd64 与 arm64 上编译结果与 32 位混合完全
// 等长（各 12 条指令）只有 386（33 对 28 条）和 arm32（21 对 14 条）要多付
// 五到七条，而这是在一次互斥加锁加一次 map 写入（上百纳秒）旁边，量不出来
// 一个混合函数只需要调对一次，按位宽分叉会让这件事变成两件
//
// [EN] Multiply-shift with the 64-bit golden-ratio constant, taking the HIGH 32
// bits. High, not low: a multiply's carries only travel left, so bit i of the
// product depends only on bits 0..i of the input and the product's low bits
// carry no information from the high ones — with an odd constant the lowest bit
// is the input's lowest bit unchanged. Measured, taking the low bits is exactly
// as bad as no mixing at all. Not split by platform word size either: the
// 64-bit multiply compiles to the same instruction count as a 32-bit mix on
// amd64 and arm64, and the five to seven extra instructions on 386 and arm32
// sit next to a mutex acquisition and a map write, so they cannot be measured
//
// [EN] Folds every bit of a 64-bit key down into the low bits for masking.
// Multiply-shift with the 64-bit golden-ratio constant: multiplication spreads
// the influence of each input bit upwards, so taking the high half and masking
// it distributes well. Costs one multiply and one shift
func mixUint64(v uint64) uint32 {
	return uint32((v * 0x9E3779B97F4A7C15) >> 32)
}

// KvHasher 为 ShardedMap 的 K 类型选择最优的 hash 函数
// 通过类型断言为常见 K 类型（string/整数族等）选择专用 hasher，避免反射开销
//
// 整数键必须经过位混合，不能直接把原值当 hash
//
// 分片下标取的是 hash 的低位（h & mask，mask = shardCount-1），而原值直接当
// hash 意味着分片**完全由键的低位决定**真实的键常常在低位带结构：按分片数
// 对齐的 ID、乘过 1000 的编号、以分为单位的金额、毫秒时间戳
// 实测一万个键、64 个分片（混合前）：
//
//	键步长    用到的分片   最大分片装入（理想 156）
//	1         64/64        157
//	8         8/64         1250
//	64        1/64         10000   ← 退化成一把全局锁
//	100       16/64        625
//	1000      8/64         1250
//
// 步长 64 时所有键挤进同一个分片，分片这件事等于没做，而且比普通 map 更糟——
// 多付了分片的开销却只用到一把锁这类退化不会报错，只会在压测时表现为
// 「分片 map 没有想象中快」，极难归因，所以必须在 hasher 里挡住
//
// [EN] Integer keys must be bit-mixed rather than used as their own hash. The
// shard index comes from the low bits of the hash (h & mask with
// mask = shardCount-1), so using the raw value means the shard is decided
// entirely by the key's low bits — and real keys often carry structure there:
// IDs aligned to the shard count, numbers multiplied by 1000, amounts in cents,
// millisecond timestamps. Measured over ten thousand keys across 64 shards
// before mixing, a stride of 64 put every key in one shard, degenerating to a
// single global lock while still paying for sharding. Such degradation raises no
// error; it merely shows up as "the sharded map is not as fast as expected",
// which is very hard to attribute, so the hasher has to prevent it
func KvHasher[K comparable]() func(K) uint32 {
	switch any(*new(K)).(type) {
	case string:
		return func(k K) uint32 { return FNVHashString32(any(k).(string)) }
	case int:
		return func(k K) uint32 { return mixUint64(uint64(any(k).(int))) }
	case int64:
		return func(k K) uint32 { return mixUint64(uint64(any(k).(int64))) }
	case int32:
		return func(k K) uint32 { return mixUint64(uint64(uint32(any(k).(int32)))) }
	case int16:
		return func(k K) uint32 { return mixUint64(uint64(uint16(any(k).(int16)))) }
	case int8:
		return func(k K) uint32 { return mixUint64(uint64(uint8(any(k).(int8)))) }
	case uint:
		return func(k K) uint32 { return mixUint64(uint64(any(k).(uint))) }
	case uint64:
		return func(k K) uint32 { return mixUint64(any(k).(uint64)) }
	case uint32:
		return func(k K) uint32 { return mixUint64(uint64(any(k).(uint32))) }
	case uint16:
		return func(k K) uint32 { return mixUint64(uint64(any(k).(uint16))) }
	case uint8:
		return func(k K) uint32 { return mixUint64(uint64(any(k).(uint8))) }
	default:
		// 其他类型用 fmt.Sprintf 转 string 再 FNV-1a hash
		//
		// 这条路每次操作都要格式化并分配一次字符串，是上面所有专用分支的几十倍
		// 键类型不在上面列表里（例如结构体键）时会落到这里，而它不会有任何提示
		// 这类键建议自己换成 string 或整数，或在 ShardedMap 外面做一层映射
		//
		// [EN] This path formats and allocates a string on every operation, tens
		// of times dearer than the specialised branches above. Any key type not
		// listed lands here silently, struct keys included; prefer converting
		// such keys to a string or integer, or map them outside ShardedMap
		return func(k K) uint32 { return FNVHashString32(fmt.Sprintf("%v", k)) }
	}
}
