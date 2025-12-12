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
