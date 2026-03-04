/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-04 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-04 12:00:00
 * @FilePath: \go-toolbox\pkg\syncx\dynamic_ticker_test.go
 * @Description: DynamicTicker 测试文件
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// 传统方式测试（直接使用 DynamicTicker）
// ============================================================================

// TestDynamicTicker_BasicUsage 测试基础用法
func TestDynamicTicker_BasicUsage(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	// 接收至少 3 个 tick
	count := 0
	timeout := time.After(500 * time.Millisecond)

	for count < 3 {
		select {
		case <-ticker.C:
			count++
		case <-timeout:
			t.Fatalf("超时：只收到 %d 个 tick，期望至少 3 个", count)
		}
	}

	assert.GreaterOrEqual(t, count, 3, "应该收到至少 3 个 tick")
}

// TestDynamicTicker_UpdateInterval 测试动态更新间隔
func TestDynamicTicker_UpdateInterval(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	// 等待第一个 tick
	select {
	case <-ticker.C:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("未收到第一个 tick")
	}

	// 更新间隔为 50ms
	ticker.UpdateInterval(50 * time.Millisecond)

	// 验证新间隔
	assert.Equal(t, 50*time.Millisecond, ticker.GetInterval(), "间隔应该更新为 50ms")

	// 验证新间隔生效（应该更快收到 tick）
	start := time.Now()
	select {
	case <-ticker.C:
		elapsed := time.Since(start)
		assert.Less(t, elapsed, 100*time.Millisecond, "新间隔应该生效，耗时应小于 100ms")
	case <-time.After(200 * time.Millisecond):
		t.Fatal("更新间隔后未收到 tick")
	}
}

// TestDynamicTicker_UpdateSameInterval 测试更新相同间隔
func TestDynamicTicker_UpdateSameInterval(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	// 更新为相同间隔（应该不做任何操作）
	ticker.UpdateInterval(100 * time.Millisecond)

	assert.Equal(t, 100*time.Millisecond, ticker.GetInterval(), "间隔应该保持不变")
}

// TestDynamicTicker_UpdateBeforeStart 测试启动前更新间隔
func TestDynamicTicker_UpdateBeforeStart(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)

	// 启动前更新间隔
	ticker.UpdateInterval(50 * time.Millisecond)

	assert.Equal(t, 50*time.Millisecond, ticker.GetInterval(), "启动前应该可以更新间隔")

	ticker.Start()
	defer ticker.Stop()

	// 验证启动后使用新间隔
	start := time.Now()
	select {
	case <-ticker.C:
		elapsed := time.Since(start)
		assert.Less(t, elapsed, 100*time.Millisecond, "启动后应该使用新间隔")
	case <-time.After(200 * time.Millisecond):
		t.Fatal("启动后未收到 tick")
	}
}

// TestDynamicTicker_Stop 测试停止定时器
func TestDynamicTicker_Stop(t *testing.T) {
	ticker := NewDynamicTicker(50 * time.Millisecond)
	ticker.Start()

	// 等待第一个 tick
	select {
	case <-ticker.C:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("未收到第一个 tick")
	}

	// 停止定时器
	ticker.Stop()

	// 等待一段时间，确保不再收到 tick
	time.Sleep(150 * time.Millisecond)

	select {
	case <-ticker.C:
		t.Error("停止后不应该再收到 tick")
	default:
		// 正确：没有收到 tick
	}
}

// TestDynamicTicker_MultipleStops 测试多次停止
func TestDynamicTicker_MultipleStops(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()

	// 多次停止（不应该 panic）
	assert.NotPanics(t, func() {
		ticker.Stop()
		ticker.Stop()
		ticker.Stop()
	}, "多次停止不应该 panic")
}

// TestDynamicTicker_MultipleStarts 测试多次启动
func TestDynamicTicker_MultipleStarts(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)

	// 多次启动（只有第一次生效）
	assert.NotPanics(t, func() {
		ticker.Start()
		ticker.Start()
		ticker.Start()
	}, "多次启动不应该 panic")

	defer ticker.Stop()

	// 验证定时器正常工作
	select {
	case <-ticker.C:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("多次启动后应该正常工作")
	}
}

// TestDynamicTicker_NonBlockingSend 测试非阻塞发送
func TestDynamicTicker_NonBlockingSend(t *testing.T) {
	ticker := NewDynamicTicker(10 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	// 不接收 tick，让通道填满
	time.Sleep(100 * time.Millisecond)

	// 定时器应该继续运行，不会阻塞
	// 清空通道
	drained := 0
	for {
		select {
		case <-ticker.C:
			drained++
		default:
			goto done
		}
	}

done:
	// 应该只有 1 个缓冲的 tick（其他被丢弃）
	assert.LessOrEqual(t, drained, 1, "通道缓冲应该最多 1 个")

	// 验证定时器仍然工作
	select {
	case <-ticker.C:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("定时器应该继续工作")
	}
}

// TestDynamicTicker_ConcurrentUpdate 测试并发更新间隔
func TestDynamicTicker_ConcurrentUpdate(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	// 并发更新间隔
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(interval time.Duration) {
			ticker.UpdateInterval(interval)
			done <- true
		}(time.Duration(i+1) * 10 * time.Millisecond)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证定时器仍然工作
	select {
	case <-ticker.C:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("并发更新后定时器应该继续工作")
	}
}

// TestDynamicTicker_ConcurrentGetInterval 测试并发获取间隔
func TestDynamicTicker_ConcurrentGetInterval(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	// 并发读取间隔
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_ = ticker.GetInterval()
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestDynamicTicker_RapidIntervalChanges 测试快速切换间隔
func TestDynamicTicker_RapidIntervalChanges(t *testing.T) {
	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	// 快速切换间隔
	intervals := []time.Duration{
		50 * time.Millisecond,
		100 * time.Millisecond,
		30 * time.Millisecond,
		80 * time.Millisecond,
	}

	for _, interval := range intervals {
		ticker.UpdateInterval(interval)
		time.Sleep(10 * time.Millisecond)
	}

	// 验证定时器仍然工作
	select {
	case <-ticker.C:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("快速切换间隔后定时器应该继续工作")
	}
}

// ============================================================================
// EventLoop 方式测试（与 EventLoop 配合使用）
// ============================================================================

// TestDynamicTicker_WithEventLoop_Basic 测试与 EventLoop 基础集成
func TestDynamicTicker_WithEventLoop_Basic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewDynamicTicker(50 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	tickCount := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ticker.C, func(tickTime time.Time) {
			atomic.AddInt32(&tickCount, 1)
		})

	go loop.Run()

	// 等待接收多个 tick
	time.Sleep(200 * time.Millisecond)

	count := atomic.LoadInt32(&tickCount)
	assert.GreaterOrEqual(t, count, int32(3), "应该收到至少 3 个 tick")
}

// TestDynamicTicker_WithEventLoop_DynamicAdjust 测试在 EventLoop 中动态调整频率
func TestDynamicTicker_WithEventLoop_DynamicAdjust(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	tickCount := int32(0)
	adjustedAt := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ticker.C, func(tickTime time.Time) {
			count := atomic.AddInt32(&tickCount, 1)
			// 收到第 3 个 tick 时，加快频率
			if count == 3 && atomic.CompareAndSwapInt32(&adjustedAt, 0, 1) {
				ticker.UpdateInterval(30 * time.Millisecond)
			}
		})

	go loop.Run()

	// 等待足够长的时间
	time.Sleep(500 * time.Millisecond)

	count := atomic.LoadInt32(&tickCount)
	// 前 3 个 tick: 100ms * 3 = 300ms
	// 后续 tick: 30ms 间隔，200ms 内应该有 6-7 个
	// 总共应该有 9+ 个 tick
	assert.GreaterOrEqual(t, count, int32(8), "动态调整后应该收到更多 tick")
	assert.Equal(t, int32(1), atomic.LoadInt32(&adjustedAt), "应该触发了频率调整")
}

// TestDynamicTicker_WithEventLoop_MultipleChannels 测试 EventLoop 处理多个通道和 DynamicTicker
func TestDynamicTicker_WithEventLoop_MultipleChannels(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewDynamicTicker(50 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	messageCh := make(chan string, 10)

	tickCount := int32(0)
	messageCount := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ticker.C, func(tickTime time.Time) {
			atomic.AddInt32(&tickCount, 1)
		}).
		OnChannel(messageCh, func(msg string) {
			atomic.AddInt32(&messageCount, 1)
		})

	go loop.Run()

	// 发送消息
	messageCh <- "msg1"
	messageCh <- "msg2"
	messageCh <- "msg3"

	// 等待 tick
	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, int32(3), atomic.LoadInt32(&messageCount), "应该收到 3 条消息")
	assert.GreaterOrEqual(t, atomic.LoadInt32(&tickCount), int32(3), "应该收到至少 3 个 tick")
}

// TestDynamicTicker_WithEventLoop_ConditionalAdjustment 测试根据条件动态调整频率
func TestDynamicTicker_WithEventLoop_ConditionalAdjustment(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	loadCh := make(chan int, 10)
	tickCount := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(loadCh, func(load int) {
			// 根据负载调整频率
			if load > 80 {
				// 高负载：降低频率
				ticker.UpdateInterval(200 * time.Millisecond)
			} else if load < 20 {
				// 低负载：提高频率
				ticker.UpdateInterval(50 * time.Millisecond)
			} else {
				// 正常负载
				ticker.UpdateInterval(100 * time.Millisecond)
			}
		}).
		OnChannel(ticker.C, func(tickTime time.Time) {
			atomic.AddInt32(&tickCount, 1)
		})

	go loop.Run()

	// 模拟负载变化
	loadCh <- 10 // 低负载，频率提高到 50ms
	time.Sleep(250 * time.Millisecond)

	count1 := atomic.LoadInt32(&tickCount)
	assert.GreaterOrEqual(t, count1, int32(4), "低负载时应该有更多 tick")

	loadCh <- 90 // 高负载，频率降低到 200ms
	atomic.StoreInt32(&tickCount, 0)
	time.Sleep(450 * time.Millisecond)

	count2 := atomic.LoadInt32(&tickCount)
	assert.LessOrEqual(t, count2, int32(3), "高负载时应该有更少 tick")
}

// TestDynamicTicker_WithEventLoop_ComplexScenario 测试复杂场景
func TestDynamicTicker_WithEventLoop_ComplexScenario(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建两个 DynamicTicker
	fastTicker := NewDynamicTicker(30 * time.Millisecond)
	slowTicker := NewDynamicTicker(100 * time.Millisecond)
	fastTicker.Start()
	slowTicker.Start()
	defer fastTicker.Stop()
	defer slowTicker.Stop()

	commandCh := make(chan string, 10)

	fastCount := int32(0)
	slowCount := int32(0)
	commandCount := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(fastTicker.C, func(tickTime time.Time) {
			atomic.AddInt32(&fastCount, 1)
		}).
		OnChannel(slowTicker.C, func(tickTime time.Time) {
			atomic.AddInt32(&slowCount, 1)
		}).
		OnChannel(commandCh, func(cmd string) {
			atomic.AddInt32(&commandCount, 1)
			switch cmd {
			case "speed_up":
				slowTicker.UpdateInterval(50 * time.Millisecond)
			case "slow_down":
				fastTicker.UpdateInterval(80 * time.Millisecond)
			case "reset":
				fastTicker.UpdateInterval(30 * time.Millisecond)
				slowTicker.UpdateInterval(100 * time.Millisecond)
			}
		}).
		OnShutdown(func() {
			// 清理资源
		})

	go loop.Run()

	// 发送命令
	commandCh <- "speed_up"
	time.Sleep(200 * time.Millisecond)

	commandCh <- "slow_down"
	time.Sleep(200 * time.Millisecond)

	commandCh <- "reset"
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, int32(3), atomic.LoadInt32(&commandCount), "应该处理 3 个命令")
	assert.Greater(t, atomic.LoadInt32(&fastCount), int32(0), "快速定时器应该触发")
	assert.Greater(t, atomic.LoadInt32(&slowCount), int32(0), "慢速定时器应该触发")
}

// TestDynamicTicker_WithEventLoop_PanicRecovery 测试 EventLoop 中的 panic 恢复
func TestDynamicTicker_WithEventLoop_PanicRecovery(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewDynamicTicker(50 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	tickCount := int32(0)
	panicCount := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ticker.C, func(tickTime time.Time) {
			count := atomic.AddInt32(&tickCount, 1)
			// 第 2 个 tick 时触发 panic
			if count == 2 {
				panic("test panic")
			}
		}).
		OnPanic(func(r any) {
			atomic.AddInt32(&panicCount, 1)
		})

	go loop.Run()

	// 等待足够长的时间
	time.Sleep(300 * time.Millisecond)

	assert.GreaterOrEqual(t, atomic.LoadInt32(&tickCount), int32(4), "panic 后应该继续接收 tick")
	assert.Equal(t, int32(1), atomic.LoadInt32(&panicCount), "应该捕获 1 次 panic")
}

// ============================================================================
// 基准测试
// ============================================================================

// BenchmarkDynamicTicker_Tick 基准测试：tick 性能
func BenchmarkDynamicTicker_Tick(b *testing.B) {
	ticker := NewDynamicTicker(1 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-ticker.C
	}
}

// BenchmarkDynamicTicker_UpdateInterval 基准测试：更新间隔性能
func BenchmarkDynamicTicker_UpdateInterval(b *testing.B) {
	ticker := NewDynamicTicker(100 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ticker.UpdateInterval(time.Duration(i%100+1) * time.Millisecond)
	}
}

// BenchmarkDynamicTicker_GetInterval 基准测试：获取间隔性能
func BenchmarkDynamicTicker_GetInterval(b *testing.B) {
	ticker := NewDynamicTicker(100 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ticker.GetInterval()
	}
}

// BenchmarkDynamicTicker_WithEventLoop 基准测试：EventLoop 集成性能
func BenchmarkDynamicTicker_WithEventLoop(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewDynamicTicker(1 * time.Millisecond)
	ticker.Start()
	defer ticker.Stop()

	count := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ticker.C, func(tickTime time.Time) {
			atomic.AddInt32(&count, 1)
		})

	go loop.Run()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for atomic.LoadInt32(&count) < int32(i+1) {
			time.Sleep(100 * time.Microsecond)
		}
	}
}
