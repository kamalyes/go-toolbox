/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-11 15:57:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-16 10:15:03
 * @FilePath: \go-toolbox\tests\retry_runner_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/retry"
	"github.com/stretchr/testify/assert"
)

func TestRunner_Run_Success(t *testing.T) {
	r := retry.NewRunner[int]()

	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 42, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestRunner_Run_FnIsNil(t *testing.T) {
	r := retry.NewRunner[int]()
	result, err := r.Run(nil)

	assert.Error(t, err)
	assert.Equal(t, retry.ErrFunIsNil, err)
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunner_Run_PanicRecovered(t *testing.T) {
	r := retry.NewRunner[int]()

	result, err := r.Run(func(ctx context.Context) (int, error) {
		panic("something went wrong")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic recovered")
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunner_Run_Timeout(t *testing.T) {
	r := retry.NewRunner[int]().Timeout(50 * time.Millisecond)

	timeoutCalled := false
	r.OnTimeout(func() {
		timeoutCalled = true
	})

	result, err := r.Run(func(ctx context.Context) (int, error) {
		// 模拟长时间阻塞，超过超时限制
		time.Sleep(100 * time.Millisecond)
		return 1, nil
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, retry.ErrTimeout))
	assert.True(t, timeoutCalled)
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunner_Run_CustomTimeoutErr(t *testing.T) {
	customErr := errors.New("custom timeout error")
	r := retry.NewRunner[int]().Timeout(50 * time.Millisecond).CustomTimeoutErr(customErr)

	result, err := r.Run(func(ctx context.Context) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 1, nil
	})

	assert.Error(t, err)
	assert.Equal(t, customErr, err)
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunner_Run_OnSuccessCalled(t *testing.T) {
	r := retry.NewRunner[int]()

	doneCalled := false
	var doneResult int
	var doneErr error

	// 注册成功回调，接收泛型 int 和 error（成功时 error 应该为 nil）
	r.OnSuccess(func(result int, err error) {
		doneCalled = true
		doneResult = result
		doneErr = err
	})

	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 123, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 123, result)

	// 成功回调必须被调用，且参数正确
	assert.True(t, doneCalled)
	assert.Equal(t, 123, doneResult)
	assert.NoError(t, doneErr)
}

func TestRunner_Run_OnErrorCalled(t *testing.T) {
	r := retry.NewRunner[int]()

	errorCalled := false
	var errorResult int
	var errorErr error

	// 注册失败回调，接收泛型 int 和 error（失败时 error 不为 nil）
	r.OnError(func(result int, err error) {
		errorCalled = true
		errorResult = result
		errorErr = err
	})

	testErr := errors.New("test error")

	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 0, testErr
	})

	assert.Error(t, err)
	assert.Equal(t, 0, result)

	// 失败回调必须被调用，且参数正确
	assert.True(t, errorCalled)
	assert.Equal(t, 0, errorResult)
	assert.Equal(t, testErr, errorErr)
}

func TestRunner_RunWithLock_Success(t *testing.T) {
	r := retry.NewRunner[int]()
	mu := &sync.Mutex{}

	result, err := r.RunWithLock(mu, func(ctx context.Context) (int, error) {
		return 99, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 99, result)
}

func TestRunner_RunWithLock_LockIsNil(t *testing.T) {
	r := retry.NewRunner[int]()

	result, err := r.RunWithLock(nil, func(ctx context.Context) (int, error) {
		return 1, nil
	})

	assert.Error(t, err)
	assert.Equal(t, retry.ErrLockIsNil, err)
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunner_RunWithLock_Concurrent(t *testing.T) {
	r := retry.NewRunner[int]()
	mu := &sync.Mutex{}

	counter := 0
	const goroutines = 10
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := r.RunWithLock(mu, func(ctx context.Context) (int, error) {
				time.Sleep(10 * time.Millisecond)
				counter++ // 有锁保护，安全
				return counter, nil
			})
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	assert.Equal(t, goroutines, counter)
}
