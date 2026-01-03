/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 00:00:00 10:00:00
 * @FilePath: \go-toolbox\pkg\syncx\event_loop_test.go
 * @Description: 事件循环执行器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventLoop_BasicChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int, 10)
	received := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ch, func(val int) {
			atomic.AddInt32(&received, int32(val))
		})

	go loop.Run()

	ch <- 1
	ch <- 2
	ch <- 3

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(6), atomic.LoadInt32(&received))

	cancel()
	time.Sleep(50 * time.Millisecond)
}

func TestEventLoop_MultipleChannels(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch1 := make(chan string, 10)
	ch2 := make(chan int, 10)

	count1 := int32(0)
	count2 := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ch1, func(val string) {
			atomic.AddInt32(&count1, 1)
		}).
		OnChannel(ch2, func(val int) {
			atomic.AddInt32(&count2, 1)
		})

	go loop.Run()

	ch1 <- "a"
	ch1 <- "b"
	ch2 <- 1
	ch2 <- 2
	ch2 <- 3

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(2), atomic.LoadInt32(&count1))
	assert.Equal(t, int32(3), atomic.LoadInt32(&count2))

	cancel()
	time.Sleep(50 * time.Millisecond)
}

func TestEventLoop_Ticker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tickCount := int32(0)

	loop := NewEventLoop(ctx).
		OnTicker(50*time.Millisecond, func() {
			atomic.AddInt32(&tickCount, 1)
		})

	go loop.Run()

	time.Sleep(160 * time.Millisecond)
	count := atomic.LoadInt32(&tickCount)
	assert.GreaterOrEqual(t, count, int32(2))
	assert.LessOrEqual(t, count, int32(4))

	cancel()
	time.Sleep(50 * time.Millisecond)
}

func TestEventLoop_MultipleTickers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tick1Count := int32(0)
	tick2Count := int32(0)

	loop := NewEventLoop(ctx).
		OnTicker(50*time.Millisecond, func() {
			atomic.AddInt32(&tick1Count, 1)
		}).
		OnTicker(100*time.Millisecond, func() {
			atomic.AddInt32(&tick2Count, 1)
		})

	go loop.Run()

	time.Sleep(260 * time.Millisecond)

	count1 := atomic.LoadInt32(&tick1Count)
	count2 := atomic.LoadInt32(&tick2Count)

	assert.GreaterOrEqual(t, count1, int32(4))
	assert.GreaterOrEqual(t, count2, int32(2))

	cancel()
	time.Sleep(50 * time.Millisecond)
}

func TestEventLoop_ChannelAndTicker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int, 10)
	channelCount := int32(0)
	tickCount := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ch, func(val int) {
			atomic.AddInt32(&channelCount, 1)
		}).
		OnTicker(50*time.Millisecond, func() {
			atomic.AddInt32(&tickCount, 1)
		})

	go loop.Run()

	ch <- 1
	ch <- 2

	time.Sleep(160 * time.Millisecond)

	assert.Equal(t, int32(2), atomic.LoadInt32(&channelCount))
	assert.GreaterOrEqual(t, atomic.LoadInt32(&tickCount), int32(2))

	cancel()
	time.Sleep(50 * time.Millisecond)
}

func TestEventLoop_OnShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	shutdownCalled := int32(0)

	loop := NewEventLoop(ctx).
		OnShutdown(func() {
			atomic.AddInt32(&shutdownCalled, 1)
		})

	go loop.Run()

	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, int32(1), atomic.LoadInt32(&shutdownCalled))
}

func TestEventLoop_OnPanic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int, 10)
	panicCaught := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ch, func(val int) {
			panic("test panic")
		}).
		OnPanic(func(r interface{}) {
			atomic.AddInt32(&panicCaught, 1)
		})

	go loop.Run()

	ch <- 1
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, int32(1), atomic.LoadInt32(&panicCaught))

	cancel()
	time.Sleep(50 * time.Millisecond)
}

func TestEventLoop_RunAsync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int, 10)
	received := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(ch, func(val int) {
			atomic.AddInt32(&received, 1)
		})

	loop.RunAsync()

	ch <- 1
	ch <- 2

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(2), atomic.LoadInt32(&received))

	cancel()
	time.Sleep(50 * time.Millisecond)
}

// TestEventLoop_ComplexScenario 复杂场景测试
func TestEventLoop_ComplexScenario(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type Message struct {
		ID   int
		Data string
	}

	registerCh := make(chan string, 10)
	unregisterCh := make(chan string, 10)
	messageCh := make(chan Message, 10)

	registerCount := int32(0)
	unregisterCount := int32(0)
	messageCount := int32(0)
	heartbeatCount := int32(0)
	cleanupCount := int32(0)

	loop := NewEventLoop(ctx).
		OnChannel(registerCh, func(id string) {
			atomic.AddInt32(&registerCount, 1)
		}).
		OnChannel(unregisterCh, func(id string) {
			atomic.AddInt32(&unregisterCount, 1)
		}).
		OnChannel(messageCh, func(msg Message) {
			atomic.AddInt32(&messageCount, 1)
		}).
		OnTicker(50*time.Millisecond, func() {
			atomic.AddInt32(&heartbeatCount, 1)
		}).
		OnTicker(100*time.Millisecond, func() {
			atomic.AddInt32(&cleanupCount, 1)
		}).
		OnShutdown(func() {
			// 清理资源
		}).
		OnPanic(func(r interface{}) {
			t.Errorf("Unexpected panic: %v", r)
		})

	go loop.Run()

	// 模拟事件
	registerCh <- "user1"
	registerCh <- "user2"
	messageCh <- Message{ID: 1, Data: "hello"}
	messageCh <- Message{ID: 2, Data: "world"}
	unregisterCh <- "user1"

	time.Sleep(160 * time.Millisecond)

	assert.Equal(t, int32(2), atomic.LoadInt32(&registerCount))
	assert.Equal(t, int32(1), atomic.LoadInt32(&unregisterCount))
	assert.Equal(t, int32(2), atomic.LoadInt32(&messageCount))
	assert.GreaterOrEqual(t, atomic.LoadInt32(&heartbeatCount), int32(2))
	assert.GreaterOrEqual(t, atomic.LoadInt32(&cleanupCount), int32(1))

	cancel()
	time.Sleep(50 * time.Millisecond)
}

// Example 示例
func ExampleEventLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	registerCh := make(chan string, 10)
	unregisterCh := make(chan string, 10)

	loop := NewEventLoop(ctx).
		OnChannel(registerCh, func(id string) {
			// 处理注册
		}).
		OnChannel(unregisterCh, func(id string) {
			// 处理注销
		}).
		OnTicker(5*time.Second, func() {
			// 定期心跳检查
		}).
		OnShutdown(func() {
			// 清理资源
		})

	// 运行事件循环
	go loop.Run()

	// 或者异步运行
	// loop.RunAsync()

	// Output:
}

// TestEventLoop_IfTicker 测试条件定时器
func TestEventLoop_IfTicker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	enabledCount := int32(0)
	disabledCount := int32(0)

	loop := NewEventLoop(ctx).
		// 条件为 true，应该执行
		IfTicker(true, 50*time.Millisecond, func() {
			atomic.AddInt32(&enabledCount, 1)
		}).
		// 条件为 false，不应该执行
		IfTicker(false, 50*time.Millisecond, func() {
			atomic.AddInt32(&disabledCount, 1)
		})

	loop.Run()

	// 300ms / 50ms = 6次触发（大约5-6次）
	enabled := atomic.LoadInt32(&enabledCount)
	disabled := atomic.LoadInt32(&disabledCount)

	assert.GreaterOrEqual(t, enabled, int32(4), "启用的定时器应该触发至少4次")
	assert.Equal(t, int32(0), disabled, "禁用的定时器不应该触发")
}

// TestEventLoop_IfChannel 测试条件通道
func TestEventLoop_IfChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	enabledCh := make(chan int, 10)
	disabledCh := make(chan int, 10)

	enabledReceived := int32(0)
	disabledReceived := int32(0)

	loop := NewEventLoop(ctx).
		// 条件为 true，应该监听
		IfChannel(true, enabledCh, func(val int) {
			atomic.AddInt32(&enabledReceived, int32(val))
		}).
		// 条件为 false，不应该监听
		IfChannel(false, disabledCh, func(val int) {
			atomic.AddInt32(&disabledReceived, int32(val))
		})

	go loop.Run()

	// 发送数据
	enabledCh <- 10
	enabledCh <- 20
	disabledCh <- 100

	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, int32(30), atomic.LoadInt32(&enabledReceived), "启用的通道应该接收到数据")
	assert.Equal(t, int32(0), atomic.LoadInt32(&disabledReceived), "禁用的通道不应该接收到数据")

	cancel()
	time.Sleep(50 * time.Millisecond)
}

// TestEventLoop_MixedConditions 测试混合条件
func TestEventLoop_MixedConditions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)

	ch1Received := int32(0)
	ch2Received := int32(0)
	ticker1Count := int32(0)
	ticker2Count := int32(0)

	enableFeature1 := true
	enableFeature2 := false

	loop := NewEventLoop(ctx).
		IfChannel(enableFeature1, ch1, func(val int) {
			atomic.AddInt32(&ch1Received, int32(val))
		}).
		IfChannel(enableFeature2, ch2, func(val int) {
			atomic.AddInt32(&ch2Received, int32(val))
		}).
		IfTicker(enableFeature1, 50*time.Millisecond, func() {
			atomic.AddInt32(&ticker1Count, 1)
		}).
		IfTicker(enableFeature2, 50*time.Millisecond, func() {
			atomic.AddInt32(&ticker2Count, 1)
		})

	go loop.Run()

	ch1 <- 5
	ch2 <- 10

	time.Sleep(250 * time.Millisecond)

	assert.Equal(t, int32(5), atomic.LoadInt32(&ch1Received), "feature1的通道应该工作")
	assert.Equal(t, int32(0), atomic.LoadInt32(&ch2Received), "feature2的通道不应该工作")
	assert.GreaterOrEqual(t, atomic.LoadInt32(&ticker1Count), int32(3), "feature1的定时器应该工作")
	assert.Equal(t, int32(0), atomic.LoadInt32(&ticker2Count), "feature2的定时器不应该工作")
}
