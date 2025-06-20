/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-20 19:11:55
 * @FilePath: \go-toolbox\tests\retry_test.go
 * @Description: 重试机制单元测试文件
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

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

	"github.com/kamalyes/go-toolbox/pkg/retry"
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
			retryInstance := retry.NewRetry().
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

func TestRetry_SetConditionFunc(t *testing.T) {
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

	r := retry.NewRetry().
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
	retryInstance := retry.NewRetry().
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
	retryInstance := retry.NewRetry().
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

	retryInstance := retry.NewRetry().
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
	retryInstance := retry.NewRetryWithCtx(ctx).
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

	retryInstance := retry.NewRetryWithCtx(ctx).
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
	retryInstance := retry.NewRetry().
		SetAttemptCount(3).
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
	retryInstance := retry.NewRetry().
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

	retryInstance := retry.NewRetry().
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

func TestRetry_GetSetMethods(t *testing.T) {
	r := retry.NewRetryWithCtx(context.Background())

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
	assert.Contains(t, r.GetCaller(), "FuncName:TestRetry_GetSetMethods, File")
	// 设置自定义调用者
	var caller = "TestRetry_GetSetMethods_12356789"
	r.SetCaller(caller)
	// 再次检查
	r.Do(func() error {
		return nil
	})
	assert.Equal(t, caller, r.GetCaller())
}
