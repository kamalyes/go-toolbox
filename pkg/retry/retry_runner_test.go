/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-11 15:57:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 23:07:29
 * @FilePath: \go-toolbox\pkg\retry\retry_runner_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package retry

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunnerRunSuccess(t *testing.T) {
	r := NewRunner[int]()

	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 42, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestRunnerRunFnIsNil(t *testing.T) {
	r := NewRunner[int]()
	result, err := r.Run(nil)

	assert.Error(t, err)
	assert.Equal(t, ErrFunIsNil, err)
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunnerRunPanicRecovered(t *testing.T) {
	r := NewRunner[int]()

	result, err := r.Run(func(ctx context.Context) (int, error) {
		panic("something went wrong")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic recovered")
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunnerRunTimeout(t *testing.T) {
	r := NewRunner[int]().Timeout(50 * time.Millisecond)

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
	assert.True(t, errors.Is(err, ErrTimeout))
	assert.True(t, timeoutCalled)
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunnerRunCustomTimeoutErr(t *testing.T) {
	customErr := errors.New("custom timeout error")
	r := NewRunner[int]().Timeout(50 * time.Millisecond).CustomTimeoutErr(customErr)

	result, err := r.Run(func(ctx context.Context) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 1, nil
	})

	assert.Error(t, err)
	assert.Equal(t, customErr, err)
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunnerRunOnSuccessCalled(t *testing.T) {
	r := NewRunner[int]()

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

func TestRunnerRunOnErrorCalled(t *testing.T) {
	r := NewRunner[int]()

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

func TestRunnerRunWithLockSuccess(t *testing.T) {
	r := NewRunner[int]()
	mu := &sync.Mutex{}

	result, err := r.RunWithLock(mu, func(ctx context.Context) (int, error) {
		return 99, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 99, result)
}

func TestRunnerRunWithLockLockIsNil(t *testing.T) {
	r := NewRunner[int]()

	result, err := r.RunWithLock(nil, func(ctx context.Context) (int, error) {
		return 1, nil
	})

	assert.Error(t, err)
	assert.Equal(t, ErrLockIsNil, err)
	var zero int
	assert.Equal(t, zero, result)
}

func TestRunnerRunWithLockConcurrent(t *testing.T) {
	r := NewRunner[int]()
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

// 测试 GetTimeout 方法
func TestRunnerGetTimeout(t *testing.T) {
	r := NewRunner[int]()

	// 默认值应该是0
	assert.Equal(t, time.Duration(0), r.GetTimeout())

	// 设置超时后验证
	r.Timeout(5 * time.Second)
	assert.Equal(t, 5*time.Second, r.GetTimeout())
}

// 测试链式调用
func TestRunnerChainedCalls(t *testing.T) {
	successCalled := false
	errorCalled := false
	timeoutCalled := false

	r := NewRunner[string]().
		Timeout(time.Second).
		OnSuccess(func(result string, err error) { successCalled = true }).
		OnError(func(result string, err error) { errorCalled = true }).
		OnTimeout(func() { timeoutCalled = true }).
		CustomTimeoutErr(errors.New("custom"))

	assert.Equal(t, time.Second, r.GetTimeout())

	// 测试成功路径
	result, err := r.Run(func(ctx context.Context) (string, error) {
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.True(t, successCalled)
	assert.False(t, errorCalled)
	assert.False(t, timeoutCalled)
}

// 测试无超时情况下的 panic
func TestRunnerPanicWithoutTimeout(t *testing.T) {
	r := NewRunner[int]()
	// 不设置超时

	result, err := r.Run(func(ctx context.Context) (int, error) {
		panic("panic without timeout")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic recovered")
	assert.Equal(t, 0, result)
}

// 测试无超时情况下的错误
func TestRunnerErrorWithoutTimeout(t *testing.T) {
	r := NewRunner[string]()
	testErr := errors.New("test error without timeout")

	result, err := r.Run(func(ctx context.Context) (string, error) {
		return "", testErr
	})

	assert.Equal(t, testErr, err)
	assert.Equal(t, "", result)
}

// 测试无超时情况下的成功
func TestRunnerSuccessWithoutTimeout(t *testing.T) {
	r := NewRunner[float64]()

	result, err := r.Run(func(ctx context.Context) (float64, error) {
		return 3.14, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 3.14, result)
}

// 测试带超时的 panic
func TestRunnerPanicWithTimeout(t *testing.T) {
	r := NewRunner[int]().
		Timeout(time.Second)

	result, err := r.Run(func(ctx context.Context) (int, error) {
		panic("panic with timeout")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic recovered")
	assert.Equal(t, 0, result)
}

// 测试带超时的错误
func TestRunnerErrorWithTimeout(t *testing.T) {
	errorCalled := false
	r := NewRunner[int]().
		Timeout(time.Second).
		OnError(func(result int, err error) {
			errorCalled = true
		})

	testErr := errors.New("test error with timeout")
	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 0, testErr
	})

	assert.Equal(t, testErr, err)
	assert.Equal(t, 0, result)
	assert.True(t, errorCalled)
}

// 测试带超时的成功
func TestRunnerSuccessWithTimeout(t *testing.T) {
	successCalled := false
	r := NewRunner[int]().
		Timeout(time.Second).
		OnSuccess(func(result int, err error) {
			successCalled = true
		})

	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 42, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
	assert.True(t, successCalled)
}

// 测试超时回调在超时时被调用
func TestRunnerTimeoutCallbackCalled(t *testing.T) {
	timeoutCalled := false
	r := NewRunner[int]().
		Timeout(50 * time.Millisecond).
		OnTimeout(func() {
			timeoutCalled = true
		})

	_, err := r.Run(func(ctx context.Context) (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 1, nil
	})

	assert.Error(t, err)
	assert.True(t, timeoutCalled)
}

// 测试超时时不设置超时回调
func TestRunnerTimeoutWithoutCallback(t *testing.T) {
	r := NewRunner[int]().
		Timeout(50 * time.Millisecond)
	// 不设置 OnTimeout

	_, err := r.Run(func(ctx context.Context) (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 1, nil
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrTimeout))
}

// 测试不设置任何回调的情况
func TestRunnerNoCallbacks(t *testing.T) {
	r := NewRunner[int]().
		Timeout(time.Second)

	// 成功情况
	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 123, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 123, result)

	// 失败情况
	result, err = r.Run(func(ctx context.Context) (int, error) {
		return 0, errors.New("error")
	})
	assert.Error(t, err)
	assert.Equal(t, 0, result)
}

// 测试不同泛型类型
func TestRunnerDifferentGenericTypes(t *testing.T) {
	t.Run("string type", func(t *testing.T) {
		r := NewRunner[string]()
		result, err := r.Run(func(ctx context.Context) (string, error) {
			return "hello", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("struct type", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		r := NewRunner[Person]()
		result, err := r.Run(func(ctx context.Context) (Person, error) {
			return Person{Name: "Alice", Age: 30}, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "Alice", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("slice type", func(t *testing.T) {
		r := NewRunner[[]int]()
		result, err := r.Run(func(ctx context.Context) ([]int, error) {
			return []int{1, 2, 3}, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("map type", func(t *testing.T) {
		r := NewRunner[map[string]int]()
		result, err := r.Run(func(ctx context.Context) (map[string]int, error) {
			return map[string]int{"a": 1, "b": 2}, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, result["a"])
		assert.Equal(t, 2, result["b"])
	})

	t.Run("pointer type", func(t *testing.T) {
		r := NewRunner[*int]()
		val := 42
		result, err := r.Run(func(ctx context.Context) (*int, error) {
			return &val, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 42, *result)
	})

	t.Run("interface type", func(t *testing.T) {
		r := NewRunner[interface{}]()
		result, err := r.Run(func(ctx context.Context) (interface{}, error) {
			return "interface value", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "interface value", result)
	})
}

// 测试零值返回
func TestRunnerZeroValueReturn(t *testing.T) {
	t.Run("int zero value on error", func(t *testing.T) {
		r := NewRunner[int]()
		result, err := r.Run(func(ctx context.Context) (int, error) {
			return 0, errors.New("error")
		})
		assert.Error(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("string zero value on panic", func(t *testing.T) {
		r := NewRunner[string]()
		result, err := r.Run(func(ctx context.Context) (string, error) {
			panic("oops")
		})
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("bool zero value on timeout", func(t *testing.T) {
		r := NewRunner[bool]().Timeout(10 * time.Millisecond)
		result, err := r.Run(func(ctx context.Context) (bool, error) {
			time.Sleep(100 * time.Millisecond)
			return true, nil
		})
		assert.Error(t, err)
		assert.False(t, result)
	})
}

// 测试并发执行多个 Runner
func TestRunnerConcurrentRunners(t *testing.T) {
	var wg sync.WaitGroup
	const goroutines = 10

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			r := NewRunner[int]().Timeout(time.Second)
			result, err := r.Run(func(ctx context.Context) (int, error) {
				return idx * 2, nil
			})
			assert.NoError(t, err)
			assert.Equal(t, idx*2, result)
		}(i)
	}

	wg.Wait()
}

// 测试 RunWithLock 带 panic
func TestRunnerRunWithLockPanic(t *testing.T) {
	r := NewRunner[int]()
	mu := &sync.Mutex{}

	result, err := r.RunWithLock(mu, func(ctx context.Context) (int, error) {
		panic("lock panic")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic recovered")
	assert.Equal(t, 0, result)
}

// 测试 RunWithLock 带超时
func TestRunnerRunWithLockTimeout(t *testing.T) {
	r := NewRunner[int]().Timeout(50 * time.Millisecond)
	mu := &sync.Mutex{}

	result, err := r.RunWithLock(mu, func(ctx context.Context) (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 1, nil
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrTimeout))
	assert.Equal(t, 0, result)
}

// 测试 RunWithLock 带自定义超时错误
func TestRunnerRunWithLockCustomTimeout(t *testing.T) {
	customErr := errors.New("custom lock timeout")
	r := NewRunner[int]().
		Timeout(50 * time.Millisecond).
		CustomTimeoutErr(customErr)
	mu := &sync.Mutex{}

	result, err := r.RunWithLock(mu, func(ctx context.Context) (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 1, nil
	})

	assert.Equal(t, customErr, err)
	assert.Equal(t, 0, result)
}

// 测试 RunWithLock fn 为 nil
func TestRunnerRunWithLockFnIsNil(t *testing.T) {
	r := NewRunner[int]()
	mu := &sync.Mutex{}

	result, err := r.RunWithLock(mu, nil)

	assert.Error(t, err)
	assert.Equal(t, ErrFunIsNil, err)
	assert.Equal(t, 0, result)
}

// 测试上下文检测
func TestRunnerContextRespect(t *testing.T) {
	r := NewRunner[int]().Timeout(time.Second)

	// 任务应该能够检测到上下文
	result, err := r.Run(func(ctx context.Context) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return 42, nil
		}
	})

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

// 测试快速成功（不应该触发超时）
func TestRunnerFastSuccess(t *testing.T) {
	timeoutCalled := false
	r := NewRunner[int]().
		Timeout(time.Second).
		OnTimeout(func() {
			timeoutCalled = true
		})

	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 1, nil // 立即返回
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, result)
	assert.False(t, timeoutCalled)
}

// 测试快速错误（不应该触发超时）
func TestRunnerFastError(t *testing.T) {
	timeoutCalled := false
	errorCalled := false
	r := NewRunner[int]().
		Timeout(time.Second).
		OnTimeout(func() {
			timeoutCalled = true
		}).
		OnError(func(result int, err error) {
			errorCalled = true
		})

	result, err := r.Run(func(ctx context.Context) (int, error) {
		return 0, errors.New("fast error")
	})

	assert.Error(t, err)
	assert.Equal(t, 0, result)
	assert.False(t, timeoutCalled)
	assert.True(t, errorCalled)
}

// 测试回调接收正确的参数
func TestRunnerCallbackParams(t *testing.T) {
	t.Run("success callback params", func(t *testing.T) {
		var cbResult int
		var cbErr error

		r := NewRunner[int]().
			OnSuccess(func(result int, err error) {
				cbResult = result
				cbErr = err
			})

		result, err := r.Run(func(ctx context.Context) (int, error) {
			return 100, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, 100, result)
		assert.Equal(t, 100, cbResult)
		assert.NoError(t, cbErr)
	})

	t.Run("error callback params", func(t *testing.T) {
		var cbResult int
		var cbErr error
		testErr := errors.New("test error")

		r := NewRunner[int]().
			OnError(func(result int, err error) {
				cbResult = result
				cbErr = err
			})

		result, err := r.Run(func(ctx context.Context) (int, error) {
			return 50, testErr
		})

		assert.Equal(t, testErr, err)
		assert.Equal(t, 50, result)
		assert.Equal(t, 50, cbResult)
		assert.Equal(t, testErr, cbErr)
	})
}

// 测试多次设置超时
func TestRunnerMultipleTimeoutSets(t *testing.T) {
	r := NewRunner[int]()

	r.Timeout(100 * time.Millisecond)
	assert.Equal(t, 100*time.Millisecond, r.GetTimeout())

	r.Timeout(200 * time.Millisecond)
	assert.Equal(t, 200*time.Millisecond, r.GetTimeout())

	r.Timeout(50 * time.Millisecond)
	assert.Equal(t, 50*time.Millisecond, r.GetTimeout())
}

// 测试多次设置回调
func TestRunnerMultipleCallbackSets(t *testing.T) {
	call1 := false
	call2 := false

	r := NewRunner[int]().
		OnSuccess(func(result int, err error) {
			call1 = true
		}).
		OnSuccess(func(result int, err error) {
			call2 = true
		})

	_, err := r.Run(func(ctx context.Context) (int, error) {
		return 1, nil
	})

	assert.NoError(t, err)
	// 最后设置的回调应该生效
	assert.False(t, call1)
	assert.True(t, call2)
}

// 测试超时边界值
func TestRunnerTimeoutBoundary(t *testing.T) {
	r := NewRunner[int]().Timeout(100 * time.Millisecond)

	// 任务执行时间刚好在边界
	result, err := r.Run(func(ctx context.Context) (int, error) {
		time.Sleep(50 * time.Millisecond) // 在超时之前完成
		return 42, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

// 测试非常短的超时
func TestRunnerVeryShortTimeout(t *testing.T) {
	r := NewRunner[int]().Timeout(1 * time.Nanosecond)

	result, err := r.Run(func(ctx context.Context) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 1, nil
	})

	assert.Error(t, err)
	assert.Equal(t, 0, result)
}

// 测试错误变量
func TestRunnerErrorVariables(t *testing.T) {
	assert.Equal(t, "function execution timeout", ErrTimeout.Error())
	assert.Equal(t, "fn cannot be nil", ErrFunIsNil.Error())
	assert.Equal(t, "lock cannot be nil", ErrLockIsNil.Error())
	assert.Equal(t, "panic recovered", ErrPanic)
}
