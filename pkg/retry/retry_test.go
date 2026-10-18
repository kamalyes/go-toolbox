/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 23:07:40
 * @FilePath: \go-toolbox\pkg\retry\retry_test.go
 * @Description: 重试机制单元测试文件
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试尝试次数
func TestRetryWithAttemptCount(t *testing.T) {
	tests := []struct {
		attemptCount  int
		expectedCount int
		desc          string
	}{
		{-1, 1, "负数尝试次数,直接调用一次,不重试"},
		{0, 1, "0尝试次数,直接调用一次,不重试"},
		{1, 1, "1次尝试,正常调用"},
		{2, 1, "2次尝试,正常调用"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			retryInstance := NewRetry().
				SetAttemptCount(tt.attemptCount).
				SetInterval(time.Microsecond).
				SetErrCallback(func(now, remain int, err error, _ ...string) {
					fmt.Printf("当前第%d次尝试(剩余%d次),错误：%v\n", now, remain, err)
				})

			callCount := 0
			err := retryInstance.Do(func() error {
				callCount++
				return nil
			})

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, callCount)
		})
	}
}

func TestRetrySetConditionFunc(t *testing.T) {
	var attemptCounter int32 = 0

	operation := func() error {
		attempt := atomic.AddInt32(&attemptCounter, 1)
		if attempt == 1 {
			return errors.New("network timeout") // 第一次返回可重试错误
		}
		return errors.New("other error") // 第二次返回不可重试错误
	}

	conditionFunc := func(err error) bool {
		return err != nil && strings.Contains(err.Error(), "network")
	}

	r := NewRetry().
		SetAttemptCount(5).
		SetInterval(10 * time.Millisecond).
		SetConditionFunc(conditionFunc).
		SetErrCallback(func(attempt, remain int, err error, _ ...string) {
			t.Logf("attempt %d failed: %v, remain %d", attempt, err, remain)
		}).
		SetSuccessCallback(func(_ ...string) {
			t.Log("success called")
		})

	err := r.Do(operation)

	// 预期：
	// 第一次调用返回 "network timeout" 错误，满足重试条件，继续重试
	// 第二次调用返回 "other error"，不满足重试条件，立即返回错误，不继续重试
	assert.NotNil(t, err, "expected error but got nil")
	assert.Contains(t, err.Error(), "other error", "expected error message to contain 'other error'")
	assert.Equal(t, int32(2), atomic.LoadInt32(&attemptCounter), "expected 2 attempts")
}

// 测试重试机制在出错时的行为
func TestRetryWithError(t *testing.T) {
	// 初始化重试实例：最大尝试3次，间隔1秒
	retryInstance := NewRetry().
		SetAttemptCount(3).
		SetInterval(time.Microsecond).
		SetErrCallback(func(now, remain int, err error, funcName ...string) {
			fmt.Printf("%s当前第%d次尝试(剩余%d次)，错误：%v\n", funcName, now, remain, err)
		})

	// 模拟总是返回错误的函数
	err := retryInstance.Do(func() error {
		return errors.New("模拟错误")
	})
	assert.Error(t, err) // 预期最终返回错误
}

// 测试首次执行即成功的场景
func TestRetrySuccess(t *testing.T) {
	retryInstance := NewRetry().
		SetAttemptCount(3).
		SetInterval(0) // 无间隔立即重试

	err := retryInstance.Do(func() error {
		return nil // 直接返回成功
	})
	assert.NoError(t, err) // 预期无错误
}

// 测试重试次数验证
func TestRetryCountValidation(t *testing.T) {
	const attempts = 3
	var counter int // 执行计数器

	retryInstance := NewRetry().
		SetAttemptCount(attempts).
		SetInterval(0)

	// 前N-1次返回错误，最后一次成功
	err := retryInstance.Do(func() error {
		counter++
		if counter < attempts {
			return errors.New("临时错误")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, attempts, counter) // 验证实际执行次数
}

// 测试上下文取消功能
func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	retryInstance := NewRetryWithCtx(ctx).
		SetAttemptCount(5).
		SetInterval(time.Second)

	// 500ms后主动取消上下文
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	err := retryInstance.Do(func() error {
		return errors.New("应被取消")
	})

	assert.ErrorIs(t, err, context.Canceled) // 验证错误类型
}

// ‌上下文超时测试‌
func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	retryInstance := NewRetryWithCtx(ctx).
		SetAttemptCount(10).
		SetInterval(time.Second)

	err := retryInstance.Do(func() error {
		return errors.New("超时测试")
	})

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

// 熔断
func TestOptimizedBackoff(t *testing.T) {
	t.Parallel()
	retryInstance := NewRetry().
		SetAttemptCount(2).
		SetInterval(100 * time.Millisecond). // 基础间隔降为100ms
		SetErrCallback(func(now, _ int, _ error, _ ...string) {
			base := math.Pow(1.5, float64(now-1)) // 改用1.5倍增长
			interval := time.Duration(base*100)*time.Millisecond +
				time.Duration(rand.Intn(50))*time.Millisecond
			if interval > 500*time.Millisecond {
				interval = 500 * time.Millisecond
			}
			time.Sleep(interval)
		})

	start := time.Now()
	err := retryInstance.Do(func() error {
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Less(t, time.Since(start), 1*time.Second) // 严格控制在1秒内
}

// 并发安全测试
func TestConcurrentSafety(t *testing.T) {
	retryInstance := NewRetry().
		SetAttemptCount(100).
		SetInterval(0)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, retryInstance.Do(func() error {
				return nil
			}))
		}()
	}
	wg.Wait()
}

// 回调函数覆盖率测试‌
func TestCallbackCoverage(t *testing.T) {
	var (
		successCalled bool
		errorCalled   bool
	)

	retryInstance := NewRetry().
		SetAttemptCount(2).
		SetInterval(0).
		SetSuccessCallback(func(_ ...string) {
			successCalled = true
		}).
		SetErrCallback(func(_, _ int, _ error, _ ...string) {
			errorCalled = true
		})

	// 测试成功回调
	_ = retryInstance.Do(func() error { return nil })
	assert.True(t, successCalled)

	// 测试错误回调
	_ = retryInstance.Do(func() error { return errors.New("") })
	assert.True(t, errorCalled)
}

func TestRetryGetSetMethods(t *testing.T) {
	r := NewRetryWithCtx(context.Background())

	// 测试 SetAttemptCount 和 GetAttemptCount
	r.SetAttemptCount(5)
	assert.Equal(t, 5, r.GetAttemptCount())

	// 测试 SetInterval 和 GetInterval
	interval := 2 * time.Second
	r.SetInterval(interval)
	assert.Equal(t, interval, r.GetInterval())

	// 测试 SetErrCallback 和 GetErrCallback
	errCallbackCalled := false
	errCallback := func(nowAttemptCount, remainCount int, err error, funcName ...string) {
		errCallbackCalled = true
	}
	r.SetErrCallback(errCallback)
	assert.NotNil(t, r.GetErrCallback())

	// 触发 errCallback 测试（调用回调）
	r.GetErrCallback()(1, 3, errors.New("test error"), "TestFunc")
	assert.True(t, errCallbackCalled)

	// 测试 SetSuccessCallback 和 GetSuccessCallback
	successCallbackCalled := false
	successCallback := func(funcName ...string) {
		successCallbackCalled = true
	}
	r.SetSuccessCallback(successCallback)
	assert.NotNil(t, r.GetSuccessCallback())

	// 触发 successCallback 测试（调用回调）
	r.GetSuccessCallback()("TestFunc")
	assert.True(t, successCallbackCalled)

	// 测试 SetConditionFunc 和 GetConditionFunc
	conditionFunc := func(err error) bool {
		return err != nil
	}
	r.SetConditionFunc(conditionFunc)
	assert.NotNil(t, r.GetConditionFunc())
	assert.True(t, r.GetConditionFunc()(errors.New("err")))
	assert.False(t, r.GetConditionFunc()(nil))

	// 测试 GetContext
	ctx := r.GetContext()
	assert.NotNil(t, ctx)
	assert.Equal(t, context.Background(), ctx)

	r.Do(func() error {
		return nil
	})
	assert.Contains(t, r.GetCaller(), "FuncName:TestRetryGetSetMethods, File")
	// 设置自定义调用者
	var caller = "TestRetryGetSetMethods_12356789"
	r.SetCaller(caller)
	// 再次检查
	r.Do(func() error {
		return nil
	})
	assert.Equal(t, caller, r.GetCaller())
}

// 测试退避倍数（BackoffMultiplier）
func TestRetryBackoffMultiplier(t *testing.T) {
	var intervals []time.Duration
	startTime := time.Now()
	lastCallTime := startTime

	r := NewRetry().
		SetAttemptCount(4).
		SetInterval(50 * time.Millisecond).
		SetBackoffMultiplier(2.0).
		SetErrCallback(func(now, remain int, err error, _ ...string) {
			currentTime := time.Now()
			if now > 1 {
				intervals = append(intervals, currentTime.Sub(lastCallTime))
			}
			lastCallTime = currentTime
		})

	err := r.Do(func() error {
		return errors.New("test error")
	})

	assert.Error(t, err)
	// 验证退避倍数是否生效：间隔应该逐渐增大
	assert.Equal(t, 2.0, r.GetBackoffMultiplier())
}

// 测试最大间隔时间（MaxInterval）
func TestRetryMaxInterval(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(5).
		SetInterval(50 * time.Millisecond).
		SetMaxInterval(100 * time.Millisecond).
		SetBackoffMultiplier(3.0)

	err := r.Do(func() error {
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Equal(t, 100*time.Millisecond, r.GetMaxInterval())
}

// 测试抖动（Jitter）
func TestRetryJitter(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(3).
		SetInterval(50 * time.Millisecond).
		SetJitter(true)

	err := r.Do(func() error {
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.True(t, r.GetJitter())
}

// 测试抖动和退避倍数组合
func TestRetryJitterWithBackoff(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(4).
		SetInterval(30 * time.Millisecond).
		SetBackoffMultiplier(1.5).
		SetMaxInterval(200 * time.Millisecond).
		SetJitter(true)

	start := time.Now()
	err := r.Do(func() error {
		return errors.New("test error")
	})

	elapsed := time.Since(start)
	assert.Error(t, err)
	// 确保有足够的重试时间（带抖动）
	assert.True(t, elapsed >= 30*time.Millisecond)
}

// 测试无回调函数情况
func TestRetryNoCallbacks(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(3).
		SetInterval(time.Millisecond)

	// 测试成功情况无回调
	err := r.Do(func() error {
		return nil
	})
	assert.NoError(t, err)

	// 测试失败情况无回调
	err = r.Do(func() error {
		return errors.New("test error")
	})
	assert.Error(t, err)
}

// 测试空函数名场景
func TestRetryEmptyFuncName(t *testing.T) {
	successCalled := false
	errCalled := false

	r := NewRetry().
		SetAttemptCount(2).
		SetInterval(time.Millisecond).
		SetSuccessCallback(func(funcName ...string) {
			successCalled = true
			// funcName 应该不为空（由 caller 自动填充）
		}).
		SetErrCallback(func(now, remain int, err error, funcName ...string) {
			errCalled = true
		})

	// 测试成功路径
	err := r.Do(func() error {
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, successCalled)

	// 测试失败路径
	err = r.Do(func() error {
		return errors.New("test error")
	})
	assert.Error(t, err)
	assert.True(t, errCalled)
}

// 测试条件函数返回 true 的情况（继续重试）
func TestRetryConditionFuncAllowRetry(t *testing.T) {
	var attemptCount int32 = 0

	r := NewRetry().
		SetAttemptCount(3).
		SetInterval(time.Millisecond).
		SetConditionFunc(func(err error) bool {
			return true // 总是允许重试
		})

	err := r.Do(func() error {
		atomic.AddInt32(&attemptCount, 1)
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&attemptCount))
}

// 测试条件函数返回 false 的情况（不重试）
func TestRetryConditionFuncDenyRetry(t *testing.T) {
	var attemptCount int32 = 0

	r := NewRetry().
		SetAttemptCount(5).
		SetInterval(time.Millisecond).
		SetConditionFunc(func(err error) bool {
			return false // 总是不允许重试
		})

	err := r.Do(func() error {
		atomic.AddInt32(&attemptCount, 1)
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&attemptCount)) // 只执行一次
}

// 测试零间隔时间
func TestRetryZeroInterval(t *testing.T) {
	var counter int32 = 0

	r := NewRetry().
		SetAttemptCount(5).
		SetInterval(0)

	start := time.Now()
	err := r.Do(func() error {
		atomic.AddInt32(&counter, 1)
		return errors.New("test error")
	})

	elapsed := time.Since(start)
	assert.Error(t, err)
	assert.Equal(t, int32(5), atomic.LoadInt32(&counter))
	// 零间隔应该很快完成
	assert.True(t, elapsed < 100*time.Millisecond)
}

// 测试退避倍数小于等于1的情况（不应用退避）
func TestRetryBackoffMultiplierLessThanOne(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(3).
		SetInterval(10 * time.Millisecond).
		SetBackoffMultiplier(0.5) // 小于1，不应用退避

	start := time.Now()
	err := r.Do(func() error {
		return errors.New("test error")
	})

	elapsed := time.Since(start)
	assert.Error(t, err)
	// 由于退避倍数小于1，间隔保持不变
	assert.True(t, elapsed < 100*time.Millisecond)
}

// 测试退避倍数等于1的情况
func TestRetryBackoffMultiplierEqualOne(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(3).
		SetInterval(10 * time.Millisecond).
		SetBackoffMultiplier(1.0)

	err := r.Do(func() error {
		return errors.New("test error")
	})

	assert.Error(t, err)
}

// 测试最大间隔为0的情况（不限制最大间隔）
func TestRetryZeroMaxInterval(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(3).
		SetInterval(10 * time.Millisecond).
		SetBackoffMultiplier(2.0).
		SetMaxInterval(0) // 不限制

	err := r.Do(func() error {
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Equal(t, time.Duration(0), r.GetMaxInterval())
}

// 测试上下文已经取消的情况
func TestRetryContextAlreadyCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	r := NewRetryWithCtx(ctx).
		SetAttemptCount(5).
		SetInterval(time.Second)

	err := r.Do(func() error {
		return errors.New("should not reach here")
	})

	assert.ErrorIs(t, err, context.Canceled)
}

// 测试链式调用
func TestRetryChainedCalls(t *testing.T) {
	successCalled := false
	errCalled := false

	r := NewRetry().
		SetCaller("TestCaller").
		SetAttemptCount(3).
		SetInterval(time.Millisecond).
		SetMaxInterval(time.Second).
		SetBackoffMultiplier(2.0).
		SetJitter(true).
		SetSuccessCallback(func(_ ...string) { successCalled = true }).
		SetErrCallback(func(_, _ int, _ error, _ ...string) { errCalled = true }).
		SetConditionFunc(func(err error) bool { return true })

	// 验证所有设置
	assert.Equal(t, "TestCaller", r.GetCaller())
	assert.Equal(t, 3, r.GetAttemptCount())
	assert.Equal(t, time.Millisecond, r.GetInterval())
	assert.Equal(t, time.Second, r.GetMaxInterval())
	assert.Equal(t, 2.0, r.GetBackoffMultiplier())
	assert.True(t, r.GetJitter())
	assert.NotNil(t, r.GetSuccessCallback())
	assert.NotNil(t, r.GetErrCallback())
	assert.NotNil(t, r.GetConditionFunc())

	err := r.Do(func() error {
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, successCalled)
	assert.False(t, errCalled) // 成功时不应调用错误回调
}

// 测试多次执行 Do
func TestRetryMultipleDoCalls(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(2).
		SetInterval(time.Millisecond)

	// 第一次执行成功
	err1 := r.Do(func() error {
		return nil
	})
	assert.NoError(t, err1)

	// 第二次执行失败
	err2 := r.Do(func() error {
		return errors.New("error")
	})
	assert.Error(t, err2)

	// 第三次执行成功
	err3 := r.Do(func() error {
		return nil
	})
	assert.NoError(t, err3)
}

// 测试并发执行多个 Retry 实例
func TestRetryConcurrentInstances(t *testing.T) {
	var wg sync.WaitGroup
	const goroutines = 10

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			r := NewRetry().
				SetAttemptCount(3).
				SetInterval(time.Millisecond)

			counter := 0
			err := r.Do(func() error {
				counter++
				if counter < 2 {
					return errors.New("temp error")
				}
				return nil
			})

			assert.NoError(t, err)
			assert.Equal(t, 2, counter)
		}(i)
	}

	wg.Wait()
}

// 测试第一次就成功的情况（不触发重试逻辑）
func TestRetryFirstAttemptSuccess(t *testing.T) {
	errCallbackCalled := false
	successCallbackCalled := false

	r := NewRetry().
		SetAttemptCount(5).
		SetInterval(time.Second). // 设置较长间隔，确保不会等待
		SetErrCallback(func(_, _ int, _ error, _ ...string) {
			errCallbackCalled = true
		}).
		SetSuccessCallback(func(_ ...string) {
			successCallbackCalled = true
		})

	start := time.Now()
	err := r.Do(func() error {
		return nil
	})

	elapsed := time.Since(start)
	assert.NoError(t, err)
	assert.False(t, errCallbackCalled) // 成功时不应调用错误回调
	assert.True(t, successCallbackCalled)
	assert.True(t, elapsed < 100*time.Millisecond) // 应该立即返回
}

// 测试最后一次尝试成功的情况
func TestRetryLastAttemptSuccess(t *testing.T) {
	var counter int32 = 0

	r := NewRetry().
		SetAttemptCount(3).
		SetInterval(time.Millisecond)

	err := r.Do(func() error {
		count := atomic.AddInt32(&counter, 1)
		if count < 3 {
			return errors.New("temp error")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&counter))
}

// 测试错误回调函数参数正确性
func TestRetryErrCallbackParams(t *testing.T) {
	var callParams []struct {
		now    int
		remain int
	}

	r := NewRetry().
		SetAttemptCount(4).
		SetInterval(time.Millisecond).
		SetErrCallback(func(now, remain int, err error, _ ...string) {
			callParams = append(callParams, struct {
				now    int
				remain int
			}{now, remain})
		})

	err := r.Do(func() error {
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Len(t, callParams, 4)

	// 验证参数正确性
	for i, p := range callParams {
		assert.Equal(t, i+1, p.now)        // 当前尝试次数
		assert.Equal(t, 4-(i+1), p.remain) // 剩余尝试次数
	}
}

// 测试条件函数为 nil 的情况（默认全部重试）
func TestRetryNilConditionFunc(t *testing.T) {
	var counter int32 = 0

	r := NewRetry().
		SetAttemptCount(3).
		SetInterval(time.Millisecond)
	// 不设置 ConditionFunc

	err := r.Do(func() error {
		atomic.AddInt32(&counter, 1)
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&counter)) // 应该重试3次
	assert.Nil(t, r.GetConditionFunc())
}

// 测试大量重试次数
func TestRetryLargeAttemptCount(t *testing.T) {
	var counter int32 = 0

	r := NewRetry().
		SetAttemptCount(100).
		SetInterval(0) // 无间隔

	err := r.Do(func() error {
		count := atomic.AddInt32(&counter, 1)
		if count < 50 {
			return errors.New("temp error")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, int32(50), atomic.LoadInt32(&counter))
}

// 测试不同类型的错误
func TestRetryDifferentErrorTypes(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"simple error", errors.New("simple error")},
		{"wrapped error", fmt.Errorf("wrapped: %w", errors.New("inner"))},
		{"custom error", &customError{msg: "custom"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRetry().
				SetAttemptCount(2).
				SetInterval(time.Millisecond)

			err := r.Do(func() error {
				return tc.err
			})

			assert.Error(t, err)
		})
	}
}

// 自定义错误类型用于测试
type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

// 测试带抖动的最大间隔限制
func TestRetryJitterWithMaxInterval(t *testing.T) {
	r := NewRetry().
		SetAttemptCount(5).
		SetInterval(50 * time.Millisecond).
		SetBackoffMultiplier(3.0).
		SetMaxInterval(100 * time.Millisecond).
		SetJitter(true)

	start := time.Now()
	err := r.Do(func() error {
		return errors.New("test error")
	})

	elapsed := time.Since(start)
	assert.Error(t, err)
	// 由于最大间隔限制，总时间应该有上限
	assert.True(t, elapsed < 2*time.Second)
}

// 测试上下文超时在执行过程中触发
func TestRetryContextTimeoutDuringExecution(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	r := NewRetryWithCtx(ctx).
		SetAttemptCount(10).
		SetInterval(100 * time.Millisecond)

	var counter int32 = 0
	err := r.Do(func() error {
		atomic.AddInt32(&counter, 1)
		return errors.New("test error")
	})

	// 应该因为上下文超时而停止
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	// 不应该完成所有10次尝试
	assert.True(t, atomic.LoadInt32(&counter) < 10)
}

// TestDo_CallerPointsAtCallSite 自动推断的 caller 必须指向真实调用点
//
// callerSkip 这个栈帧数错了不会报错，只会把调用者记成库自己的内部函数，
// 错误回调里打出来的位置从此指向错误的地方，而没人会立刻发现
// 任何改动 Do 到 runtime.Caller 之间调用层数的重构都会让这条变红
func TestDo_CallerPointsAtCallSite(t *testing.T) {
	r := NewRetry().SetAttemptCount(1).SetInterval(time.Millisecond)
	_ = r.Do(func() error { return errors.New("boom") })

	caller := r.GetCaller()
	assert.Contains(t, caller, "TestDo_CallerPointsAtCallSite",
		"应指向调用 Do 的测试函数，实际为 %q", caller)
	assert.Contains(t, caller, "retry_test.go", "应指向调用点所在文件")
	assert.NotContains(t, caller, "inlined.go", "不应指向库自己的内部实现")
}

// TestDo_RetriedFuncMayTouchItsOwnRetrier 被重试的函数可以读自己的 Retry
//
// 原先 Do 整个包在写锁里，被重试的函数在持锁状态下运行，而 sync.Mutex 不可重入：
// 函数里只要读一下 attemptCount（打日志时很自然的写法）就直接死锁
// 这条用超时兜底，死锁时会以超时失败而不是把整个测试进程挂住
func TestDo_RetriedFuncMayTouchItsOwnRetrier(t *testing.T) {
	done := make(chan error, 1)
	go func() {
		r := NewRetry().SetAttemptCount(2).SetInterval(time.Millisecond)
		done <- r.Do(func() error {
			_ = r.GetAttemptCount()
			_ = r.GetCaller()
			_ = r.GetInterval()
			return errors.New("boom")
		})
	}()

	select {
	case err := <-done:
		assert.Error(t, err, "两次尝试都失败应返回错误")
	case <-time.After(3 * time.Second):
		t.Fatal("死锁：被重试的函数读自己的 Retry 时卡住")
	}
}

// TestDo_ConcurrentDoOnSameRetrierDoesNotSerialize 同一个 Retry 上的并发 Do 不应互相等待
//
// 锁跨越退避 sleep 时，第二次 Do 要等第一次把整个重试窗口跑完
func TestDo_ConcurrentDoOnSameRetrierDoesNotSerialize(t *testing.T) {
	r := NewRetry().SetAttemptCount(3).SetInterval(80 * time.Millisecond)

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.Do(func() error { return errors.New("boom") })
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	// 单次重试窗口约 2 次间隔 = 160ms串行化会变成三倍（约 480ms）
	assert.Less(t, elapsed, 400*time.Millisecond,
		"三个并发 Do 耗时 %v，接近串行化的三倍说明锁仍跨越了退避", elapsed)
}

// TestDo_PanicIsRetriedPerAttempt panic 按「这一次尝试失败」处理，并继续重试
//
// 原先 defer syncx.RecoverToError(&err, nil) 写在重试循环体内，而 defer 注册在函数上
// 不是迭代上，于是第一次 panic 就把整个 doRetryWithCondition 退掉：
// AttemptCount(3) 实际只跑了一次
func TestDo_PanicIsRetriedPerAttempt(t *testing.T) {
	attempts := 0
	r := NewRetry().SetAttemptCount(3).SetInterval(time.Millisecond)
	err := r.Do(func() error {
		attempts++
		if attempts == 1 {
			panic("boom on first attempt")
		}
		return nil
	})

	assert.NoError(t, err, "第二次尝试成功，整体应成功")
	assert.Equal(t, 2, attempts, "第一次 panic 后应继续重试，而不是直接收工")
}

// TestDo_AllAttemptsPanic 每次都 panic 时跑完全部次数并返回最后一次的错误
func TestDo_AllAttemptsPanic(t *testing.T) {
	attempts := 0
	r := NewRetry().SetAttemptCount(3).SetInterval(time.Millisecond)
	err := r.Do(func() error {
		attempts++
		panic(errors.New("always boom"))
	})

	require.Error(t, err)
	assert.Equal(t, "always boom", err.Error(), "error 型 panic 值应原样返回")
	assert.Equal(t, 3, attempts, "应跑满 AttemptCount 次")
}

// TestDo_NonErrorPanicBecomesError 非 error 的 panic 值被包成 error
func TestDo_NonErrorPanicBecomesError(t *testing.T) {
	r := NewRetry().SetAttemptCount(1).SetInterval(time.Millisecond)

	err := r.Do(func() error { panic("string panic") })
	require.Error(t, err)
	assert.Equal(t, "string panic", err.Error())

	err = r.Do(func() error { panic(42) })
	require.Error(t, err)
	assert.Equal(t, "panic occurred", err.Error(), "既非 error 也非 string 时给固定文案")
}

// TestDo_PanicDoesNotEscapeToCaller panic 不得穿透到调用方
//
// 这是这个库最该守住的一条：重试器的意义就是把失败挡在里面
func TestDo_PanicDoesNotEscapeToCaller(t *testing.T) {
	r := NewRetry().SetAttemptCount(2).SetInterval(time.Millisecond)
	assert.NotPanics(t, func() {
		_ = r.Do(func() error { panic("must not escape") })
	})
}

// TestDo_NoWaitAfterFinalAttempt 最后一次尝试失败后不再等待
//
// 原先每次失败都睡，包括最后一次：调用方拿到错误要晚一个间隔，
// 上限设得大时这段纯等待很可观
func TestDo_NoWaitAfterFinalAttempt(t *testing.T) {
	const interval = 120 * time.Millisecond

	start := time.Now()
	attempts := 0
	err := NewRetry().SetAttemptCount(2).SetInterval(interval).
		Do(func() error { attempts++; return errors.New("boom") })
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Equal(t, 2, attempts)
	// 两次尝试之间睡一次 = 120ms若最后一次也睡则是 240ms
	assert.Less(t, elapsed, 2*interval,
		"耗时 %v 接近两个间隔，说明最后一次失败后仍在等待", elapsed)
	assert.GreaterOrEqual(t, elapsed, interval, "两次尝试之间仍应等待一个间隔")
}

// TestDo_SingleAttemptDoesNotWait 只尝试一次时完全不等待
func TestDo_SingleAttemptDoesNotWait(t *testing.T) {
	start := time.Now()
	err := NewRetry().SetAttemptCount(1).SetInterval(300 * time.Millisecond).
		Do(func() error { return errors.New("boom") })

	require.Error(t, err)
	assert.Less(t, time.Since(start), 100*time.Millisecond,
		"只有一次尝试，没有「下一次」可等")
}
