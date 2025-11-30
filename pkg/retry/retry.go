/*
* @Author: kamalyes 501893067@qq.com
* @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-23 09:50:58
 * @FilePath: \go-toolbox\pkg\retry\retry.go
* @Description: 重试机制
*
* 主要功能:
*   - 支持自定义最大尝试次数和重试间隔
*   - 支持指数退避策略（Exponential Backoff）
*   - 支持随机抖动（Jitter）避免惊群效应
*   - 支持上下文取消和超时控制
*   - 支持自定义重试条件函数
*   - 支持成功/失败回调函数
*
* 退避倍数（Backoff Multiplier）说明:
*   退避倍数用于实现指数退避策略，每次重试失败后，下一次重试的等待时间会乘以退避倍数。
*   公式: 下一次间隔 = 当前间隔 × 退避倍数
*
*   示例（初始间隔100ms，退避倍数2.0，最大间隔1s）:
*   | 重试次数 | 等待时间 |
*   |---------|---------|
*   | 第1次失败后 | 100ms |
*   | 第2次失败后 | 200ms (100 × 2) |
*   | 第3次失败后 | 400ms (200 × 2) |
*   | 第4次失败后 | 800ms (400 × 2) |
*   | 第5次失败后 | 1000ms (受 maxInterval 限制) |
*
* 使用场景:
*   - 避免服务雪崩：当服务暂时不可用时，逐渐增大重试间隔
*   - 网络抖动恢复：给网络足够的恢复时间
*   - 资源竞争：减少对有限资源的争抢
*
* 使用示例:
*   retry.NewRetry().
*       SetAttemptCount(5).
*       SetInterval(100 * time.Millisecond).
*       SetBackoffMultiplier(2.0).
*       SetMaxInterval(5 * time.Second).
*       SetJitter(true).
*       Do(func() error {
*           return someOperation()
*       })
*
* Copyright (c) 2024 by kamalyes, All Rights Reserved.
*/
package retry

import (
	"context"
	"fmt"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"sync"
	"time"
)

// Retry 结构体用于实现重试机制
type Retry struct {
	ctx               context.Context     // 上下文
	mu                sync.RWMutex        // 并发锁
	caller            string              // 调用者
	attemptCount      int                 // 最大尝试次数
	interval          time.Duration       // 重试间隔时间
	maxInterval       time.Duration       // 最大重试间隔时间
	backoffMultiplier float64             // 退避倍数
	jitter            bool                // 是否添加随机抖动
	errCallFun        ErrCallbackFunc     // 错误回调函数
	successCallFun    SuccessCallbackFunc // 成功回调函数
	conditionFunc     func(error) bool    // 重试条件函数
}

// DoFun 定义执行函数的类型
type DoFun func() error

// ErrCallbackFunc 是重试时的错误回调函数类型
// nowAttemptCount 表示当前尝试次数
// remainCount 表示剩余尝试次数
// err 是当前执行时的错误信息
type ErrCallbackFunc func(nowAttemptCount, remainCount int, err error, funcName ...string)

// SuccessCallbackFunc 是成功回调函数类型
type SuccessCallbackFunc func(funcName ...string)

// NewRetry 创建一个重试器，返回一个 Retry 实例
func NewRetry() *Retry {
	return NewRetryWithCtx(context.Background())
}

// NewRetryWithCtx 创建一个自定义上下文的重试器，返回一个 Retry实例
func NewRetryWithCtx(ctx context.Context) *Retry {
	return &Retry{
		ctx: ctx,
	}
}

// SetCaller 设置调用者信息，返回 Retry 实例以支持链式调用
func (r *Retry) SetCaller(caller string) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.caller = caller
		return r
	})
}

// SetAttemptCount 设置最大尝试次数，返回 Retry 实例以支持链式调用
func (r *Retry) SetAttemptCount(attemptCount int) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.attemptCount = attemptCount
		return r
	})
}

// SetInterval 设置重试间隔时间，返回 Retry 实例以支持链式调用
func (r *Retry) SetInterval(interval time.Duration) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.interval = interval
		return r
	})
}

// SetMaxInterval 设置最大重试间隔时间，返回 Retry 实例以支持链式调用
func (r *Retry) SetMaxInterval(maxInterval time.Duration) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.maxInterval = maxInterval
		return r
	})
}

// SetBackoffMultiplier 设置退避倍数，返回 Retry 实例以支持链式调用
func (r *Retry) SetBackoffMultiplier(multiplier float64) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.backoffMultiplier = multiplier
		return r
	})
}

// SetJitter 设置是否添加随机抖动，返回 Retry 实例以支持链式调用
func (r *Retry) SetJitter(jitter bool) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.jitter = jitter
		return r
	})
}

// SetErrCallback 设置错误回调函数，返回 Retry 实例以支持链式调用
func (r *Retry) SetErrCallback(fn ErrCallbackFunc) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.errCallFun = fn
		return r
	})
}

// SetSuccessCallback 设置成功回调函数，返回 Retry 实例以支持链式调用
func (r *Retry) SetSuccessCallback(fn SuccessCallbackFunc) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.successCallFun = fn
		return r
	})
}

// SetConditionFunc 设置重试条件函数
func (r *Retry) SetConditionFunc(fn func(error) bool) *Retry {
	return syncx.WithLockReturnValue(&r.mu, func() *Retry {
		r.conditionFunc = fn
		return r
	})
}

// GetCaller 获取调用者信息
func (r *Retry) GetCaller() string {
	return syncx.WithLockReturnValue(&r.mu, func() string {
		return r.caller
	})
}

// GetAttemptCount 获取最大尝试次数
func (r *Retry) GetAttemptCount() int {
	return syncx.WithRLockReturnValue(&r.mu, func() int {
		return r.attemptCount
	})
}

// GetInterval 获取重试间隔时间
func (r *Retry) GetInterval() time.Duration {
	return syncx.WithRLockReturnValue(&r.mu, func() time.Duration {
		return r.interval
	})
}

// GetMaxInterval 获取最大重试间隔时间
func (r *Retry) GetMaxInterval() time.Duration {
	return syncx.WithRLockReturnValue(&r.mu, func() time.Duration {
		return r.maxInterval
	})
}

// GetBackoffMultiplier 获取退避倍数
func (r *Retry) GetBackoffMultiplier() float64 {
	return syncx.WithRLockReturnValue(&r.mu, func() float64 {
		return r.backoffMultiplier
	})
}

// GetJitter 获取是否添加随机抖动
func (r *Retry) GetJitter() bool {
	return syncx.WithRLockReturnValue(&r.mu, func() bool {
		return r.jitter
	})
}

// GetErrCallback 获取错误回调函数
func (r *Retry) GetErrCallback() ErrCallbackFunc {
	return syncx.WithRLockReturnValue(&r.mu, func() ErrCallbackFunc {
		return r.errCallFun
	})
}

// GetSuccessCallback 获取成功回调函数
func (r *Retry) GetSuccessCallback() SuccessCallbackFunc {
	return syncx.WithRLockReturnValue(&r.mu, func() SuccessCallbackFunc {
		return r.successCallFun
	})
}

// GetConditionFunc 获取重试条件函数
func (r *Retry) GetConditionFunc() func(error) bool {
	return syncx.WithRLockReturnValue(&r.mu, func() func(error) bool {
		return r.conditionFunc
	})
}

// GetContext 获取上下文
func (r *Retry) GetContext() context.Context {
	return syncx.WithRLockReturnValue(&r.mu, func() context.Context {
		return r.ctx
	})
}

// Do 为 Retry 结构体定义执行函数，执行指定函数 f
func (r *Retry) Do(fn DoFun) (err error) {
	return syncx.WithLockReturnValue(&r.mu, func() error {
		r.caller = mathx.IfDo(r.caller == "", func() string {
			caller := osx.GetRuntimeCaller(7)
			defer caller.Release()
			return caller.String()
		}, r.caller)
		// 确保尝试次数为正数
		r.attemptCount = mathx.IF(r.attemptCount < 1, 1, r.attemptCount)
		return doRetryWithCondition(r.ctx, r.attemptCount, r.interval, r.maxInterval, r.backoffMultiplier, r.jitter, fn, r.errCallFun, r.successCallFun, r.conditionFunc, r.caller)
	})
}

// doRetryWithCondition 内部函数，定义了重试操作，执行指定次数的尝试
func doRetryWithCondition(ctx context.Context, attemptCount int, interval, maxInterval time.Duration, backoffMultiplier float64, jitter bool, fn DoFun, errCallFun ErrCallbackFunc, successCallFun SuccessCallbackFunc, conditionFunc func(error) bool, funcName ...string) (err error) {
	var (
		fName           = mathx.IF(len(funcName) > 0, funcName[0], "") // 获取函数名称
		nowAttemptCount int
		currentInterval = interval // 当前重试间隔
	)

	for nowAttemptCount < attemptCount {
		nowAttemptCount++

		select {
		case <-ctx.Done():
			return ctx.Err() // 如果上下文被取消，返回错误
		default:
			// 捕获 panic
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic occurred: %v", r)
				}
			}()

			if err = fn(); err == nil { // 执行传入的函数
				if successCallFun != nil {
					successCallFun(fName) // 调用成功回调函数
				}
				return // 如果没有错误，返回
			}

			// 判断是否满足重试条件，默认全部重试
			if conditionFunc != nil && !conditionFunc(err) {
				return err // 不满足重试条件，直接返回错误
			}

			if errCallFun != nil {
				errCallFun(nowAttemptCount, attemptCount-nowAttemptCount, err, fName) // 调用错误回调函数
			}

			// 等待指定的间隔时间
			if currentInterval > 0 {
				waitTime := currentInterval

				// 添加随机抖动
				if jitter {
					jitterRange := float64(currentInterval) * 0.2 // 20% 的抖动范围
					waitTime = currentInterval + time.Duration(random.RandFloat(0, jitterRange))
				}

				time.Sleep(waitTime)

				// 应用退避倍数
				if backoffMultiplier > 1.0 {
					currentInterval = time.Duration(float64(currentInterval) * backoffMultiplier)
					// 限制最大间隔
					if maxInterval > 0 && currentInterval > maxInterval {
						currentInterval = maxInterval
					}
				}
			}
		}
	}

	return err // 返回最后的错误
}
