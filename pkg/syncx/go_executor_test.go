/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 00:00:00 09:00:00
 * @FilePath: \go-toolbox\pkg\syncx\go_executor_test.go
 * @Description: Goroutine 执行器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGo_Exec(t *testing.T) {
	t.Run("basic execution", func(t *testing.T) {
		executed := int32(0)

		Go().Exec(func() {
			atomic.AddInt32(&executed, 1)
		})

		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
	})

	t.Run("with context", func(t *testing.T) {
		executed := int32(0)
		ctx := context.Background()

		Go(ctx).Exec(func() {
			atomic.AddInt32(&executed, 1)
		})

		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
	})
}

func TestGo_ExecWithPanic(t *testing.T) {
	panicCaught := int32(0)

	Go().
		OnPanic(func(r interface{}) {
			atomic.AddInt32(&panicCaught, 1)
		}).
		Exec(func() {
			panic("test panic")
		})

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&panicCaught))
}

func TestGo_ExecWithDelay(t *testing.T) {
	start := time.Now()
	executed := int32(0)

	Go().
		WithDelay(100 * time.Millisecond).
		Exec(func() {
			atomic.AddInt32(&executed, 1)
		})

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(0), atomic.LoadInt32(&executed)) // 还未执行

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
	assert.GreaterOrEqual(t, time.Since(start), 100*time.Millisecond)
}

func TestGo_ExecWithContext(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		executed := int32(0)

		Go().ExecWithContext(func(ctx context.Context) error {
			atomic.AddInt32(&executed, 1)
			return nil
		})

		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
	})

	t.Run("with error", func(t *testing.T) {
		errorCaught := int32(0)
		testErr := errors.New("test error")

		Go().
			OnError(func(err error) {
				if err == testErr {
					atomic.AddInt32(&errorCaught, 1)
				}
			}).
			ExecWithContext(func(ctx context.Context) error {
				return testErr
			})

		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&errorCaught))
	})
}

func TestGo_ExecWithTimeout(t *testing.T) {
	t.Run("complete before timeout", func(t *testing.T) {
		executed := int32(0)

		Go().
			WithTimeout(200 * time.Millisecond).
			ExecWithContext(func(ctx context.Context) error {
				atomic.AddInt32(&executed, 1)
				return nil
			})

		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
	})

	t.Run("timeout", func(t *testing.T) {
		errorCaught := int32(0)

		Go().
			WithTimeout(50 * time.Millisecond).
			OnError(func(err error) {
				if err == context.DeadlineExceeded {
					atomic.AddInt32(&errorCaught, 1)
				}
			}).
			ExecWithContext(func(ctx context.Context) error {
				<-ctx.Done()
				return ctx.Err()
			})

		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&errorCaught))
	})
}

func TestGo_ExecWithDelayAndTimeout(t *testing.T) {
	executed := int32(0)

	Go().
		WithDelay(50 * time.Millisecond).
		WithTimeout(200 * time.Millisecond).
		ExecWithContext(func(ctx context.Context) error {
			atomic.AddInt32(&executed, 1)
			return nil
		})

	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, int32(0), atomic.LoadInt32(&executed)) // 延迟中

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&executed)) // 已执行
}

func TestGo_ExecWithCancel(t *testing.T) {
	cancelCalled := int32(0)
	executed := int32(0)

	ctx, cancel := context.WithCancel(context.Background())

	Go(ctx).
		WithDelay(100 * time.Millisecond).
		OnCancel(func() {
			atomic.AddInt32(&cancelCalled, 1)
		}).
		ExecWithContext(func(ctx context.Context) error {
			atomic.AddInt32(&executed, 1)
			return nil
		})

	// 在延迟期间取消
	time.Sleep(30 * time.Millisecond)
	cancel()

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&cancelCalled))
	assert.Equal(t, int32(0), atomic.LoadInt32(&executed)) // 未执行
}

func TestGo_ExecWithResult(t *testing.T) {
	resultChan := Go().ExecWithResult(func() (interface{}, error) {
		return 42, nil
	})

	result := <-resultChan
	assert.NoError(t, result.Err)
	assert.Equal(t, 42, result.Value)
}

func TestGo_ExecWithResultError(t *testing.T) {
	testErr := errors.New("test error")
	errorCaught := int32(0)

	resultChan := Go().
		OnError(func(err error) {
			atomic.AddInt32(&errorCaught, 1)
		}).
		ExecWithResult(func() (interface{}, error) {
			return nil, testErr
		})

	result := <-resultChan
	assert.Error(t, result.Err)
	assert.Equal(t, testErr, result.Err)

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&errorCaught))
}

func TestGo_Wait(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		err := Go().Wait(func() error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("with error", func(t *testing.T) {
		testErr := errors.New("test error")
		errorCaught := int32(0)

		err := Go().
			OnError(func(err error) {
				atomic.AddInt32(&errorCaught, 1)
			}).
			Wait(func() error {
				return testErr
			})

		assert.Error(t, err)
		assert.Equal(t, testErr, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&errorCaught))
	})

	t.Run("with panic", func(t *testing.T) {
		panicCaught := int32(0)

		err := Go().
			OnPanic(func(r interface{}) {
				atomic.AddInt32(&panicCaught, 1)
			}).
			Wait(func() error {
				panic("test panic")
			})

		assert.NoError(t, err) // panic 被捕获了
		assert.Equal(t, int32(1), atomic.LoadInt32(&panicCaught))
	})
}

func TestGo_ChainedCallbacks(t *testing.T) {
	errorCaught := int32(0)
	panicCaught := int32(0)
	cancelCalled := int32(0)

	Go().
		WithTimeout(100 * time.Millisecond).
		WithDelay(50 * time.Millisecond).
		OnError(func(err error) {
			atomic.AddInt32(&errorCaught, 1)
		}).
		OnPanic(func(r interface{}) {
			atomic.AddInt32(&panicCaught, 1)
		}).
		OnCancel(func() {
			atomic.AddInt32(&cancelCalled, 1)
		}).
		ExecWithContext(func(ctx context.Context) error {
			return nil
		})

	time.Sleep(150 * time.Millisecond)
	// 应该正常完成,不触发任何错误回调
	assert.Equal(t, int32(0), atomic.LoadInt32(&errorCaught))
	assert.Equal(t, int32(0), atomic.LoadInt32(&panicCaught))
	assert.Equal(t, int32(0), atomic.LoadInt32(&cancelCalled))
}

// 基准测试
func BenchmarkGo_Exec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		done := make(chan struct{})
		Go().Exec(func() {
			close(done)
		})
		<-done
	}
}

func BenchmarkGo_ExecWithContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		done := make(chan struct{})
		Go().ExecWithContext(func(ctx context.Context) error {
			close(done)
			return nil
		})
		<-done
	}
}

func BenchmarkGo_ExecWithTimeout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		done := make(chan struct{})
		Go().
			WithTimeout(1 * time.Second).
			ExecWithContext(func(ctx context.Context) error {
				close(done)
				return nil
			})
		<-done
	}
}

// 示例测试
func ExampleGo() {
	// 基础执行
	Go().Exec(func() {
		// 执行任务
	})

	// 带错误处理
	Go().
		OnError(func(err error) {
			// log.Error(err)
		}).
		OnPanic(func(r interface{}) {
			// log.Error("panic", r)
		}).
		ExecWithContext(func(ctx context.Context) error {
			// 执行可能失败的任务
			return nil
		})

	// 带超时和延迟
	Go().
		WithTimeout(5 * time.Second).
		WithDelay(1 * time.Second).
		ExecWithContext(func(ctx context.Context) error {
			// 延迟1秒后执行,最多等待5秒
			return nil
		})

	// Output:
}

// TestGo_WithNilContext 测试传入 nil context
func TestGo_WithNilContext(t *testing.T) {
	executed := int32(0)

	Go(nil).Exec(func() {
		atomic.AddInt32(&executed, 1)
	})

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
}

// TestGo_WithParentContext 测试传入父 context
func TestGo_WithParentContext(t *testing.T) {
	parentCtx := context.Background()
	executed := int32(0)

	Go(parentCtx).
		WithTimeout(100 * time.Millisecond).
		ExecWithContext(func(ctx context.Context) error {
			atomic.AddInt32(&executed, 1)
			return nil
		})

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
}

// TestGo_DeprecatedGoWithContext 测试废弃的 GoWithContext 方法仍然可用
func TestGo_DeprecatedGoWithContext(t *testing.T) {
	ctx := context.Background()
	executed := int32(0)

	GoWithContext(ctx).Exec(func() {
		atomic.AddInt32(&executed, 1)
	})

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
}
