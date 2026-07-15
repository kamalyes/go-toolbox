/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-11 22:08:59
 * @FilePath: \go-toolbox\pkg\syncx\batch_processor.go
 * @Description: 泛型批量处理器
 *   持续收集 T 类型的请求，满 batchSize 或每 flushInterval 触发一次 flush
 *   适用于高并发写入场景（如批量 DB 更新、批量日志写入、批量指标上报等）
 *
 * 工作流程：
 *   1. 调用方调用 Submit（非阻塞，队列满时返回 false）
 *   2. 后台 worker 收集请求，满 batchSize 或每 flushInterval 触发一次 flush
 *   3. flush 时调用 flushFn 回调，由调用方处理批量逻辑
 *   4. Stop 时 drain channel 并 flush 剩余数据后退出
 *
 * 示例:
 *
 *	processor := NewBatchProcessor(4096, 100, 500*time.Millisecond, func(batch []string) {
 *	    db.BatchInsert(batch)
 *	})
 *	defer processor.Stop()
 *
 *	for _, item := range items {
// 	    processor.Submit(item)
 *	}
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
*/

package syncx

import (
	"context"
	"time"
)

// BatchProcessor 泛型批量处理器
// 持续收集 T 类型的请求，满 batchSize 或每 flushInterval 触发一次 flush
type BatchProcessor[T any] struct {
	queue         chan T          // 请求队列
	flushInterval time.Duration   // 最大 flush 间隔
	batchSize     int             // 每批最大数量
	flushFn       func(batch []T) // flush 回调
	stopChan      chan struct{}   // 停止信号通道
	done          chan struct{}   // 完成信号通道
}

// NewBatchProcessor 创建批量处理器并启动后台 worker
//
// 参数：
//   - queueSize: channel 缓冲大小（建议 4096）
//   - batchSize: 每批最大数量（建议 100）
//   - flushInterval: 最大 flush 间隔（建议 500ms）
//   - flushFn: flush 回调，接收一批数据
func NewBatchProcessor[T any](queueSize, batchSize int, flushInterval time.Duration, flushFn func(batch []T)) *BatchProcessor[T] {
	p := &BatchProcessor[T]{
		queue:         make(chan T, queueSize),
		flushInterval: flushInterval,
		batchSize:     batchSize,
		flushFn:       flushFn,
		stopChan:      make(chan struct{}),
		done:          make(chan struct{}),
	}
	go p.run()
	return p
}

// Submit 非阻塞提交一条数据
//
// 行为说明：
//   - 队列未满：item 写入 channel，返回 true
//   - 队列已满：不会阻塞等待，也不会重试，直接走 default 分支丢弃该 item，返回 false
//
// 也就是说队列满之后，新提交的数据进不来，会被静默丢弃（最终一致性语义）
// 调用方可根据返回值决定后续处理：
//   - 记录 metric（如丢弃计数）便于监控背压
//   - 降级处理（如改走同步更新、写本地缓冲等）
//
// 这种"非阻塞 + 丢弃"策略适用于可丢失的非核心路径（如消息状态更新），
// 避免生产者因队列满而阻塞，拖垮上游主流程
func (p *BatchProcessor[T]) Submit(item T) bool {
	select {
	case p.queue <- item:
		return true
	default:
		// 队列满，直接丢弃，不阻塞调用方
		return false
	}
}

// SubmitBlocking 阻塞式提交一条数据（不丢数据的兜底方案）
//
// 行为说明：
//   - 队列未满：item 立即写入 channel，返回 true
//   - 队列已满：阻塞等待，直到有空位可写入（返回 true）或 ctx 被取消/超时（返回 false）
//
// 与 Submit 的区别：
//   - Submit：队列满立即丢弃，适合可丢失的非核心路径（如消息状态更新）
//   - SubmitBlocking：队列满时背压传导给调用方，适合不能丢数据的核心路径（如订单、扣款）
//
// 注意事项：
//   - 阻塞期间会占用调用 goroutine，建议配合 ctx 超时使用，避免长时间卡住
//   - 若下游 flush 持续慢导致队列长期满，阻塞会向上游传导背压，上游需有相应限流/超时机制
//   - ctx.Done() 触发时该条数据被丢弃（由调用方决定是否重试/落盘补偿）
//
// 兜底策略选择：
//   - 不想丢 + 可接受短暂阻塞 → SubmitBlocking(ctx, item)
//   - 不想阻塞 + 可接受丢失 → Submit(item)
//   - 想要"先阻塞再降级" → SubmitBlocking 配合较短 ctx 超时，超时后走降级逻辑
func (p *BatchProcessor[T]) SubmitBlocking(ctx context.Context, item T) bool {
	select {
	case p.queue <- item:
		return true
	case <-ctx.Done():
		// ctx 超时或取消，该条数据未能写入（由调用方决定是否补偿）
		return false
	}
}

// run 后台 worker，收集并批量 flush
func (p *BatchProcessor[T]) run() {
	defer close(p.done)

	batch := make([]T, 0, p.batchSize)
	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case item := <-p.queue:
			batch = append(batch, item)
			if len(batch) >= p.batchSize {
				batch = p.flushAndReset(batch)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				batch = p.flushAndReset(batch)
			}
		case <-p.stopChan:
			// drain channel，flush 剩余数据
			for len(p.queue) > 0 {
				batch = append(batch, <-p.queue)
			}
			if len(batch) > 0 {
				p.flushFn(batch)
			}
			return
		}
	}
}

// flushAndReset 执行 flush 并重置 batch（复用底层数组，减少 GC）
func (p *BatchProcessor[T]) flushAndReset(batch []T) []T {
	p.flushFn(batch)
	return batch[:0]
}

// Stop 停止处理器，flush 剩余数据后退出
func (p *BatchProcessor[T]) Stop() {
	close(p.stopChan)
	<-p.done
}
