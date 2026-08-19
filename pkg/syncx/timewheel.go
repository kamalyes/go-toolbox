/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-20 09:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-20 09:00:00
 * @FilePath: \go-toolbox\pkg\syncx\timewheel.go
 * @Description:
 * 分片时间轮（HashedWheelTimer）——64 分片 + 双向链表 bucket + 惰性取消，
 * O(1) 调度/取消/刷新，替代 O(N) 全量扫描，适用于百万级连接心跳超时管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// 分片时间轮（Hashed Wheel Timer）
// ============================================================================
//
// 设计要点：
//   - 64 个 shard 分散锁竞争，每个 shard 独立 wheel + worker goroutine
//   - 每个 wheel 是环形 bucket 数组（默认 512），每个 bucket 持有任务双向链表
//   - FNV-1a hash 分配 shard，与 sharded_registry 模式对齐
//   - 惰性取消：atomic.Bool 标记，worker 遍历时统一清理（取消 O(1) 无锁开销）
//   - O(1) 调度 / O(1) 取消 / O(1) 心跳刷新（取消旧 + 调度新）
//   - 零额外分配：任务复用双向链表节点，无中间 slice
//
// 精度说明：
//   - 触发精度为 ±1 tick（默认 1ms），不保证精确到纳秒
//   - 任务可能因 worker 调度竞态而延迟至多 1 圈（bucketCount * tickInterval）触发
//   - 任务永远不会提前触发（ticks 向上取整 + rounds 保证）
//
// 适用场景：
//   - 百万级连接心跳超时管理（替代 O(N) 全量扫描）
//   - 跨节点投递 ACK 超时兜底
//   - 任意大量短/中周期定时任务管理

const (
	defaultTimerShardCount   = 64                   // 默认分片数（与 sharded_registry 对齐）
	defaultTimerBucketCount  = 512                  // 默认每分片 bucket 数
	defaultTimerTickInterval = 1 * time.Millisecond // 默认 tick 间隔（1ms 极致精度，64 分片共 64000 tick/s，CPU 可控）
)

// Timer 时间轮定时器接口
type Timer interface {
	// Schedule 调度一个无 key 的定时任务（轮询分配 shard）
	Schedule(delay time.Duration, task func()) *TimerTask
	// ScheduleWithKey 调度一个带 key 的定时任务（同 key 覆盖会惰性取消旧任务）
	ScheduleWithKey(key string, delay time.Duration, task func()) *TimerTask
	// Refresh 刷新任务（O(1) 心跳刷新：取消旧 + 调度新，等价于 ScheduleWithKey）
	Refresh(key string, delay time.Duration, task func()) *TimerTask
	// Cancel 取消指定任务（惰性取消，O(1)）
	Cancel(t *TimerTask)
	// CancelByKey 按 key 取消任务（惰性取消，O(1)）
	CancelByKey(key string)
	// Stop 停止时间轮（停止所有 worker，不再执行待触发任务）
	Stop()
	// Stats 返回时间轮统计信息
	Stats() TimerStats
}

// TimerTask 时间轮任务（双向链表节点）
type TimerTask struct {
	key       string      // 任务键（空表示无 key 任务）
	callback  func()      // 回调函数
	rounds    int64       // 剩余圈数（>0 表示需绕整圈后才到期）
	bucketIdx int         // 所在 bucket 索引
	shardIdx  int         // 所在 shard 索引
	cancelled atomic.Bool // 惰性取消标记
	prev      *TimerTask  // 前驱节点
	next      *TimerTask  // 后继节点
	wheel     *HashedWheelTimer
}

// TimerStats 时间轮统计信息
type TimerStats struct {
	ActiveTasks    int64 // 活跃任务数
	CompletedTasks int64 // 已完成任务数（已执行回调）
	CancelledTasks int64 // 已取消任务数（惰性取消）
}

// HashedWheelTimer 分片时间轮主实现
type HashedWheelTimer struct {
	shards         []*wheelShard
	shardCount     int
	mask           int          // shardCount-1，位运算取模
	keyIndex       sync.Map     // key → *TimerTask（仅 keyed 任务）
	rrCounter      atomic.Int64 // 无 key 任务的轮询分配计数器
	started        atomic.Bool
	stopped        atomic.Bool
	executor       func(func()) // 回调执行器（默认安全异步执行）
	completedCount atomic.Int64
	cancelledCount atomic.Int64
}

// wheelShard 单个分片时间轮
type wheelShard struct {
	buckets      []*bucketEntry
	bucketCount  int
	mask         int // bucketCount-1，位运算取模
	tickInterval time.Duration
	currentPos   atomic.Int64 // 下一个待处理 bucket 位置
	taskCount    atomic.Int64 // 本 shard 活跃任务数
	stopCh       chan struct{}
	wg           sync.WaitGroup
	parent       *HashedWheelTimer
}

// bucketEntry 单个桶（双向链表 + mutex）
type bucketEntry struct {
	mu   sync.Mutex
	head *TimerTask
}

// timerConfig 时间轮初始化配置（私有，避免外部直接修改）
type timerConfig struct {
	shardCount   int
	bucketCount  int
	tickInterval time.Duration
	executor     func(func())
}

// TimerOption 时间轮配置选项
type TimerOption func(*timerConfig)

// WithTimerShardCount 设置分片数量（必须是 2 的幂）
func WithTimerShardCount(count int) TimerOption {
	return func(c *timerConfig) {
		if count > 0 && count&(count-1) == 0 {
			c.shardCount = count
		}
	}
}

// WithTimerBucketCount 设置每分片 bucket 数量（必须是 2 的幂）
func WithTimerBucketCount(count int) TimerOption {
	return func(c *timerConfig) {
		if count > 0 && count&(count-1) == 0 {
			c.bucketCount = count
		}
	}
}

// WithTimerTickInterval 设置 tick 间隔
func WithTimerTickInterval(d time.Duration) TimerOption {
	return func(c *timerConfig) {
		if d > 0 {
			c.tickInterval = d
		}
	}
}

// WithTimerExecutor 设置回调执行器
func WithTimerExecutor(fn func(func())) TimerOption {
	return func(c *timerConfig) {
		if fn != nil {
			c.executor = fn
		}
	}
}

// NewHashedWheelTimer 创建分片时间轮
func NewHashedWheelTimer(opts ...TimerOption) *HashedWheelTimer {
	cfg := &timerConfig{
		shardCount:   defaultTimerShardCount,
		bucketCount:  defaultTimerBucketCount,
		tickInterval: defaultTimerTickInterval,
		executor:     safeGoExec,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	t := &HashedWheelTimer{
		shardCount: cfg.shardCount,
		mask:       cfg.shardCount - 1,
		executor:   cfg.executor,
	}

	t.shards = make([]*wheelShard, cfg.shardCount)
	for i := range t.shards {
		shard := &wheelShard{
			buckets:      make([]*bucketEntry, cfg.bucketCount),
			bucketCount:  cfg.bucketCount,
			mask:         cfg.bucketCount - 1,
			tickInterval: cfg.tickInterval,
			stopCh:       make(chan struct{}),
			parent:       t,
		}
		for j := range shard.buckets {
			shard.buckets[j] = &bucketEntry{}
		}
		t.shards[i] = shard
	}

	t.started.Store(true)
	for _, shard := range t.shards {
		shard.wg.Add(1)
		go shard.run()
	}

	return t
}

// safeGoExec 默认回调执行器：异步执行 + panic 恢复
// 防止单个任务 panic 导致 worker goroutine 崩溃
func safeGoExec(f func()) {
	go func() {
		defer func() { _ = recover() }()
		f()
	}()
}

// run shard worker 主循环：每 tick 推进一个 bucket
func (s *wheelShard) run() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.tickInterval)
	defer ticker.Stop()

	pos := int(s.currentPos.Load())
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.processBucket(pos)
			pos = (pos + 1) & s.mask
			s.currentPos.Store(int64(pos))
		}
	}
}

// processBucket 处理指定 bucket 的所有任务
// 三阶段：提取链表（持锁最短）→ 锁外分类处理 → 重新挂回未到期任务
func (s *wheelShard) processBucket(pos int) {
	b := s.buckets[pos]

	// Phase 1：提取整个 bucket 链表（持锁时间最短）
	b.mu.Lock()
	head := b.head
	b.head = nil
	b.mu.Unlock()

	if head == nil {
		return
	}

	// Phase 2：遍历提取的任务，分类处理（锁外执行，避免回调阻塞 worker）
	var keepHead, keepTail *TimerTask
	current := head
	for current != nil {
		next := current.next
		current.prev = nil
		current.next = nil

		if current.cancelled.Load() {
			// 惰性取消：从链表移除，不执行回调
			s.parent.cancelledCount.Add(1)
			s.taskCount.Add(-1)
			if current.key != "" {
				// 清理 key 索引（仅当索引仍指向此任务时）
				s.parent.keyIndex.CompareAndDelete(current.key, current)
			}
			current = next
			continue
		}

		if current.rounds > 0 {
			// 未到期，圈数减一，保留到下一轮
			current.rounds--
			if keepHead == nil {
				keepHead = current
				keepTail = current
			} else {
				keepTail.next = current
				current.prev = keepTail
				keepTail = current
			}
			current = next
			continue
		}

		// 到期，执行回调
		s.parent.completedCount.Add(1)
		s.taskCount.Add(-1)
		if current.key != "" {
			// 清理已执行任务的 key 索引
			s.parent.keyIndex.CompareAndDelete(current.key, current)
		}
		task := current
		s.parent.executor(task.callback)
		current = next
	}

	// Phase 3：重新挂回未到期任务（持锁，合并处理期间新 Schedule 的任务）
	if keepHead != nil {
		b.mu.Lock()
		keepTail.next = b.head
		if b.head != nil {
			b.head.prev = keepTail
		}
		b.head = keepHead
		b.mu.Unlock()
	}
}

// schedule 在当前 shard 上调度任务（内部方法）
func (s *wheelShard) schedule(delay time.Duration, callback func(), key string, shardIdx int) *TimerTask {
	// 向上取整计算 tick 数（保证不提前触发）
	ticks := int((delay + s.tickInterval - 1) / s.tickInterval)
	if ticks < 1 {
		ticks = 1
	}

	pos := int(s.currentPos.Load())
	bucketIdx := (pos + ticks) & s.mask
	rounds := int64(ticks / s.bucketCount)

	task := &TimerTask{
		key:       key,
		callback:  callback,
		rounds:    rounds,
		bucketIdx: bucketIdx,
		shardIdx:  shardIdx,
		wheel:     s.parent,
	}

	// 头插法加入 bucket 链表（O(1)）
	b := s.buckets[bucketIdx]
	b.mu.Lock()
	task.next = b.head
	if b.head != nil {
		b.head.prev = task
	}
	b.head = task
	b.mu.Unlock()

	s.taskCount.Add(1)
	return task
}

// Schedule 调度一个无 key 的定时任务（轮询分配 shard）
func (t *HashedWheelTimer) Schedule(delay time.Duration, task func()) *TimerTask {
	if t.stopped.Load() || delay <= 0 || task == nil {
		return nil
	}
	idx := int(uint64(t.rrCounter.Add(1))) & t.mask
	return t.shards[idx].schedule(delay, task, "", idx)
}

// ScheduleWithKey 调度一个带 key 的定时任务（同 key 覆盖会惰性取消旧任务）
func (t *HashedWheelTimer) ScheduleWithKey(key string, delay time.Duration, task func()) *TimerTask {
	if t.stopped.Load() || delay <= 0 || task == nil || key == "" {
		return nil
	}
	idx := int(FNVHashString32(key)) & t.mask
	tt := t.shards[idx].schedule(delay, task, key, idx)
	if tt != nil {
		// 原子交换：新任务入索引，旧任务（若存在）惰性取消
		if old, loaded := t.keyIndex.Swap(key, tt); loaded {
			if oldTask, ok := old.(*TimerTask); ok {
				oldTask.cancelled.Store(true)
			}
		}
	}
	return tt
}

// Refresh 刷新任务（O(1) 心跳刷新：取消旧 + 调度新）
func (t *HashedWheelTimer) Refresh(key string, delay time.Duration, task func()) *TimerTask {
	return t.ScheduleWithKey(key, delay, task)
}

// Cancel 取消指定任务（惰性取消，O(1)）
func (t *HashedWheelTimer) Cancel(tt *TimerTask) {
	if tt == nil || tt.wheel != t {
		return
	}
	tt.cancelled.Store(true)
	if tt.key != "" {
		// CAS 删除：仅当索引仍指向此任务时删除（避免误删 Refresh 新建的任务）
		t.keyIndex.CompareAndDelete(tt.key, tt)
	}
}

// CancelByKey 按 key 取消任务（惰性取消，O(1)）
func (t *HashedWheelTimer) CancelByKey(key string) {
	if key == "" {
		return
	}
	if old, loaded := t.keyIndex.LoadAndDelete(key); loaded {
		if oldTask, ok := old.(*TimerTask); ok {
			oldTask.cancelled.Store(true)
		}
	}
}

// Stop 停止时间轮（停止所有 worker，不再执行待触发任务）
func (t *HashedWheelTimer) Stop() {
	if !t.stopped.CompareAndSwap(false, true) {
		return
	}
	t.started.Store(false)
	for _, shard := range t.shards {
		close(shard.stopCh)
	}
	for _, shard := range t.shards {
		shard.wg.Wait()
	}
}

// Stats 返回时间轮统计信息
func (t *HashedWheelTimer) Stats() TimerStats {
	var active int64
	for _, shard := range t.shards {
		active += shard.taskCount.Load()
	}
	return TimerStats{
		ActiveTasks:    active,
		CompletedTasks: t.completedCount.Load(),
		CancelledTasks: t.cancelledCount.Load(),
	}
}
