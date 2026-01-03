/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 00:00:00 10:15:15
 * @FilePath: \go-toolbox\pkg\syncx\event_loop.go
 * @Description: 事件循环执行器 - 处理多路复用的事件分发
 *
 * 使用说明:
 *
 * 1. 基础事件循环:
 *    loop := NewEventLoop(ctx)
 *    loop.OnChannel(registerChan, handleRegister).
 *         OnChannel(unregisterChan, handleUnregister).
 *         OnTicker(5*time.Second, checkHeartbeat).
 *         Run()
 *
 * 2. 完整示例:
 *    loop := NewEventLoop(ctx).
 *        OnChannel(h.register, h.handleRegister).
 *        OnChannel(h.unregister, h.handleUnregister).
 *        OnChannel(h.broadcast, h.handleBroadcast).
 *        OnTicker(5*time.Second, h.checkHeartbeat).
 *        OnTicker(10*time.Second, h.reportMetrics).
 *        OnShutdown(h.cleanup).
 *        Run()
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"reflect"
	"time"
)

// EventLoop 事件循环执行器
// 用于统一管理多个通道和定时器的事件分发，替代复杂的 select 循环
//
// 主要功能:
//   - 通道事件处理: 监听多个通道，自动分发到对应的处理函数
//   - 定时器管理: 支持多个定时器，自动管理生命周期
//   - Panic 恢复: 自动捕获处理函数中的 panic
//   - 优雅关闭: 支持 context 取消和资源清理
//
// 使用场景:
//   - WebSocket/SSE 服务端的事件循环
//   - 消息队列消费者
//   - 任何需要监听多个事件源的场景
type EventLoop struct {
	ctx          context.Context       // 上下文，用于控制事件循环的生命周期
	channels     []channelHandler      // 通道处理器列表
	tickers      []tickerHandler       // 定时器处理器列表
	onShutdown   func()                // 关闭时的回调函数
	onPanic      func(interface{})     // panic 处理函数
	selectCases  []reflect.SelectCase  // reflect.Select 使用的 case 列表
	caseHandlers []func(reflect.Value) // 每个 case 对应的处理函数
}

// channelHandler 通道处理器
// 封装了通道及其对应的处理函数
type channelHandler struct {
	ch      interface{} // 通道（必须是 chan 类型）
	handler interface{} // 处理函数（必须接受通道元素类型的参数）
}

// tickerHandler 定时器处理器
// 封装了定时器的间隔、处理函数和 ticker 实例
type tickerHandler struct {
	interval time.Duration // 定时间隔
	handler  func()        // 定时触发的处理函数
	ticker   *time.Ticker  // time.Ticker 实例（在 Run 时创建）
}

// NewEventLoop 创建新的事件循环
// 参数 ctx 用于控制事件循环的生命周期，当 ctx 被取消时，事件循环会自动退出
//
// 示例:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	loop := NewEventLoop(ctx)
//	// ... 配置事件处理器 ...
//	go loop.Run()
//	// 需要停止时
//	cancel()
func NewEventLoop(ctx context.Context) *EventLoop {
	if ctx == nil {
		ctx = context.Background()
	}
	return &EventLoop{
		ctx:          ctx,
		channels:     make([]channelHandler, 0),
		tickers:      make([]tickerHandler, 0),
		selectCases:  make([]reflect.SelectCase, 0),
		caseHandlers: make([]func(reflect.Value), 0),
	}
}

// OnChannel 注册通道事件处理器
// 当通道接收到数据时，会调用对应的处理函数
//
// 参数:
//   - ch: 必须是一个通道类型（如 chan int, chan string, chan *Client 等）
//   - handler: 必须是一个函数，接受通道元素类型的参数
//
// 类型检查:
//   - 会在运行时通过反射验证 ch 是否为通道类型
//   - 会验证 handler 的参数类型是否匹配通道元素类型
//
// 示例:
//
//	registerChan := make(chan *Client, 100)
//	loop.OnChannel(registerChan, func(client *Client) {
//	    handleRegister(client)
//	})
//
//	messageChan := make(chan Message, 100)
//	loop.OnChannel(messageChan, func(msg Message) {
//	    handleMessage(msg)
//	})
func (el *EventLoop) OnChannel(ch interface{}, handler interface{}) *EventLoop {
	el.channels = append(el.channels, channelHandler{
		ch:      ch,
		handler: handler,
	})
	return el
}

// OnTicker 注册定时器事件处理器
// 按指定的时间间隔周期性地执行处理函数
//
// 参数:
//   - interval: 定时间隔（如 5*time.Second, 1*time.Minute）
//   - handler: 定时触发时执行的函数，无参数无返回值
//
// 特性:
//   - 定时器会在 Run() 方法中自动创建
//   - 定时器会在事件循环结束时自动停止
//   - 支持注册多个定时器，它们互不干扰
//
// 示例:
//
//	// 每5秒检查一次心跳
//	loop.OnTicker(5*time.Second, func() {
//	    checkHeartbeat()
//	})
//
//	// 每1分钟清理一次过期数据
//	loop.OnTicker(1*time.Minute, func() {
//	    cleanupExpiredData()
//	})
func (el *EventLoop) OnTicker(interval time.Duration, handler func()) *EventLoop {
	el.tickers = append(el.tickers, tickerHandler{
		interval: interval,
		handler:  handler,
	})
	return el
}

// IfTicker 条件注册定时器事件处理器
// 只有当条件为 true 时才注册定时器
//
// 参数:
//   - condition: 条件判断，为 true 时才注册
//   - interval: 定时间隔（如 5*time.Second, 1*time.Minute）
//   - handler: 定时触发时执行的函数，无参数无返回值
//
// 示例:
//
//	// 只有当启用清理时才注册定时器
//	loop.IfTicker(config.EnableCleanup, 30*time.Minute, func() {
//	    cleanupExpiredData()
//	})
func (el *EventLoop) IfTicker(condition bool, interval time.Duration, handler func()) *EventLoop {
	if condition {
		return el.OnTicker(interval, handler)
	}
	return el
}

// IfChannel 条件注册通道事件处理器
// 只有当条件为 true 时才注册通道监听
//
// 参数:
//   - condition: 条件判断，为 true 时才注册
//   - ch: 必须是一个通道类型
//   - handler: 必须是一个函数，接受通道元素类型的参数
//
// 示例:
//
//	// 只有当启用某功能时才监听该通道
//	loop.IfChannel(config.EnableFeature, featureChan, func(msg Message) {
//	    handleFeature(msg)
//	})
func (el *EventLoop) IfChannel(condition bool, ch interface{}, handler interface{}) *EventLoop {
	if condition {
		return el.OnChannel(ch, handler)
	}
	return el
}

// OnShutdown 设置关闭时的回调
// 当事件循环退出时（context 取消或发生 panic），会调用此函数进行清理
//
// 用途:
//   - 关闭数据库连接
//   - 清理临时文件
//   - 保存状态数据
//   - 记录关闭日志
//
// 注意:
//   - 此函数会在 defer 中执行，保证一定会被调用
//   - 如果设置了多次，只有最后一次设置的函数会生效
//
// 示例:
//
//	loop.OnShutdown(func() {
//	    log.Info("事件循环已停止")
//	    db.Close()
//	    saveState()
//	})
func (el *EventLoop) OnShutdown(fn func()) *EventLoop {
	el.onShutdown = fn
	return el
}

// OnPanic 设置 panic 处理器
// 当事件处理函数中发生 panic 时，会调用此函数
//
// 特性:
//   - 捕获所有处理函数中的 panic，防止整个事件循环崩溃
//   - panic 发生后，事件循环会继续运行，处理其他事件
//   - 如果设置了多次，只有最后一次设置的函数会生效
//
// 用途:
//   - 记录 panic 日志
//   - 上报错误到监控系统
//   - 发送告警通知
//
// 示例:
//
//	loop.OnPanic(func(r interface{}) {
//	    log.Error("事件处理panic", "panic", r, "stack", string(debug.Stack()))
//	    sentry.CaptureException(fmt.Errorf("panic: %v", r))
//	})
func (el *EventLoop) OnPanic(fn func(interface{})) *EventLoop {
	el.onPanic = fn
	return el
}

// buildSelectCases 构建 select cases
// 将注册的通道和定时器转换为 reflect.Select 所需的 case 列表
//
// 构建顺序:
//  1. context.Done() - 用于控制循环退出
//  2. 所有注册的通道 - 按注册顺序
//  3. 所有定时器的通道 - 按注册顺序
//
// 内部实现:
//   - 使用反射动态构建 select cases
//   - 为每个 case 创建对应的处理函数
//   - 定时器在此方法中被创建和启动
func (el *EventLoop) buildSelectCases() {
	// 添加 context.Done() case
	el.selectCases = append(el.selectCases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(el.ctx.Done()),
	})
	el.caseHandlers = append(el.caseHandlers, nil) // context done 不需要处理器

	// 添加通道 cases
	for _, ch := range el.channels {
		chanValue := reflect.ValueOf(ch.ch)
		handlerValue := reflect.ValueOf(ch.handler)

		el.selectCases = append(el.selectCases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: chanValue,
		})

		// 创建处理函数
		el.caseHandlers = append(el.caseHandlers, func(recvValue reflect.Value) {
			if recvValue.IsValid() {
				// 调用处理器
				handlerValue.Call([]reflect.Value{recvValue})
			}
		})
	}

	// 添加定时器 cases
	for i := range el.tickers {
		ticker := &el.tickers[i]
		ticker.ticker = time.NewTicker(ticker.interval)

		el.selectCases = append(el.selectCases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ticker.ticker.C),
		})

		handler := ticker.handler
		el.caseHandlers = append(el.caseHandlers, func(recvValue reflect.Value) {
			handler()
		})
	}
}

// stopTickers 停止所有定时器
// 在事件循环结束时调用，释放定时器资源
func (el *EventLoop) stopTickers() {
	for i := range el.tickers {
		if el.tickers[i].ticker != nil {
			el.tickers[i].ticker.Stop()
		}
	}
}

// Run 运行事件循环（阻塞）
// 开始监听所有注册的事件源，直到 context 被取消
//
// 执行流程:
//  1. 构建 select cases（包括通道、定时器）
//  2. 启动所有定时器
//  3. 进入 select 循环，等待事件
//  4. 接收到事件时，调用对应的处理函数
//  5. context 取消时退出循环
//  6. 停止所有定时器，调用 OnShutdown 回调
//
// 错误处理:
//   - 自动捕获所有处理函数中的 panic
//   - panic 不会导致事件循环崩溃
//   - 通过 OnPanic 回调通知上层
//
// 示例:
//
//	loop := NewEventLoop(ctx).
//	    OnChannel(ch1, handler1).
//	    OnChannel(ch2, handler2).
//	    OnTicker(5*time.Second, tickHandler)
//
//	// 阻塞运行，直到 context 被取消
//	loop.Run()
func (el *EventLoop) Run() {
	defer func() {
		if r := recover(); r != nil {
			if el.onPanic != nil {
				el.onPanic(r)
			}
		}
		el.stopTickers()
		if el.onShutdown != nil {
			el.onShutdown()
		}
	}()

	el.buildSelectCases()

	for {
		chosen, recv, ok := reflect.Select(el.selectCases)

		// context.Done() 的情况
		if chosen == 0 {
			return
		}

		// 处理接收到的事件
		if ok && el.caseHandlers[chosen] != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						if el.onPanic != nil {
							el.onPanic(r)
						}
					}
				}()
				el.caseHandlers[chosen](recv)
			}()
		}
	}
}

// RunAsync 异步运行事件循环（非阻塞）
// 在新的 goroutine 中运行事件循环，立即返回
//
// 与 Run 的区别:
//   - Run: 阻塞当前 goroutine，直到 context 被取消
//   - RunAsync: 在新 goroutine 中运行，立即返回
//
// 使用场景:
//   - 需要在主 goroutine 中继续执行其他任务
//   - 需要同时运行多个事件循环
//
// 示例:
//
//	loop := NewEventLoop(ctx).
//	    OnChannel(ch, handler).
//	    OnTicker(5*time.Second, tickHandler)
//
//	// 异步运行，不阻塞
//	loop.RunAsync()
//
//	// 可以继续执行其他代码
//	doOtherWork()
func (el *EventLoop) RunAsync() {
	Go(el.ctx).
		OnPanic(el.onPanic).
		Exec(func() {
			el.Run()
		})
}
