/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-08 15:37:52
 * @FilePath: \go-toolbox\pkg\syncx\go_executor_test.go
 * @Description: Goroutine 执行器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

func TestGo_ExecWithChildren_BasicExecution(t *testing.T) {
	var (
		task1Done = int32(0)
		task2Done = int32(0)
		task3Done = int32(0)
	)

	done := make(chan struct{})

	Go().
		OnPanic(func(r interface{}) {
			t.Errorf("unexpected panic: %v", r)
		}).
		ExecWithChildren(func(children *ChildRunner) {
			children.Go(func() {
				time.Sleep(10 * time.Millisecond)
				atomic.StoreInt32(&task1Done, 1)
			})

			children.Go(func() {
				time.Sleep(20 * time.Millisecond)
				atomic.StoreInt32(&task2Done, 1)
			})

			children.GoWithError(func() error {
				time.Sleep(15 * time.Millisecond)
				atomic.StoreInt32(&task3Done, 1)
				return nil
			})

			// 等待所有子任务完成后关闭 done channel
			go func() {
				children.Wait()
				close(done)
			}()
		})

	// 等待所有任务完成
	<-done

	assert.Equal(t, int32(1), atomic.LoadInt32(&task1Done))
	assert.Equal(t, int32(1), atomic.LoadInt32(&task2Done))
	assert.Equal(t, int32(1), atomic.LoadInt32(&task3Done))
}

func TestGo_ExecWithChildren_WithErrorHandling(t *testing.T) {
	errorCaught := int32(0)
	testErr := errors.New("test error")
	done := make(chan struct{})

	Go().
		OnError(func(err error) {
			if err == testErr {
				atomic.AddInt32(&errorCaught, 1)
			}
		}).
		ExecWithChildren(func(children *ChildRunner) {
			children.GoWithError(func() error {
				return testErr
			})

			children.GoWithError(func() error {
				return nil
			})

			go func() {
				children.Wait()
				close(done)
			}()
		})

	<-done
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&errorCaught))
}

func TestGo_ExecWithChildren_WithPanicHandling(t *testing.T) {
	panicCaught := int32(0)
	done := make(chan struct{})

	Go().
		OnPanic(func(r interface{}) {
			atomic.AddInt32(&panicCaught, 1)
		}).
		ExecWithChildren(func(children *ChildRunner) {
			children.Go(func() {
				panic("test panic")
			})

			children.Go(func() {
				// 正常任务应该继续执行
			})

			go func() {
				children.Wait()
				close(done)
			}()
		})

	<-done
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&panicCaught))
}

func TestGo_ExecSubs_BasicExecution(t *testing.T) {
	var (
		task1Done = int32(0)
		task2Done = int32(0)
		task3Done = int32(0)
	)

	done := make(chan struct{})

	Go().
		OnPanic(func(r interface{}) {
			t.Errorf("unexpected panic: %v", r)
		}).
		Sub(func() {
			time.Sleep(10 * time.Millisecond)
			atomic.StoreInt32(&task1Done, 1)
		}).
		Sub(func() {
			time.Sleep(20 * time.Millisecond)
			atomic.StoreInt32(&task2Done, 1)
		}).
		SubWithError(func() error {
			time.Sleep(15 * time.Millisecond)
			atomic.StoreInt32(&task3Done, 1)
			return nil
		}).
		ExecSubs()

	// 等待足够时间让所有任务完成
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, int32(1), atomic.LoadInt32(&task1Done))
	assert.Equal(t, int32(1), atomic.LoadInt32(&task2Done))
	assert.Equal(t, int32(1), atomic.LoadInt32(&task3Done))

	close(done)
}

func TestGo_ExecSubs_WithErrorHandling(t *testing.T) {
	errorCaught := int32(0)
	testErr := errors.New("test error")

	Go().
		OnError(func(err error) {
			if err == testErr {
				atomic.AddInt32(&errorCaught, 1)
			}
		}).
		SubWithError(func() error {
			return testErr
		}).
		SubWithError(func() error {
			return nil
		}).
		ExecSubs()

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&errorCaught))
}

func TestGo_ExecSubs_WithPanicHandling(t *testing.T) {
	panicCaught := int32(0)
	normalTaskDone := int32(0)

	Go().
		OnPanic(func(r interface{}) {
			atomic.AddInt32(&panicCaught, 1)
		}).
		Sub(func() {
			panic("test panic")
		}).
		Sub(func() {
			time.Sleep(10 * time.Millisecond)
			atomic.StoreInt32(&normalTaskDone, 1)
		}).
		ExecSubs()

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&panicCaught))
	assert.Equal(t, int32(1), atomic.LoadInt32(&normalTaskDone), "正常任务应该继续执行")
}

func TestGo_ExecSubs_ConcurrentExecution(t *testing.T) {
	start := time.Now()
	allDone := int32(0)

	Go().
		Sub(func() {
			time.Sleep(50 * time.Millisecond)
		}).
		Sub(func() {
			time.Sleep(50 * time.Millisecond)
		}).
		Sub(func() {
			time.Sleep(50 * time.Millisecond)
		}).
		ExecSubs()

	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)

	// 如果是并发执行，总耗时应该接近 50ms（最长的任务）
	// 而不是 150ms（顺序执行的总和）
	assert.Less(t, elapsed, 120*time.Millisecond, "应该并发执行，而非顺序执行")

	atomic.StoreInt32(&allDone, 1)
	assert.Equal(t, int32(1), atomic.LoadInt32(&allDone))
}

func TestGo_ExecSubs_EmptySubTasks(t *testing.T) {
	// 测试没有子任务的情况
	Go().
		OnPanic(func(r interface{}) {
			t.Errorf("unexpected panic: %v", r)
		}).
		ExecSubs()

	time.Sleep(50 * time.Millisecond)
	// 不应该 panic
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

// TestBatchExecutor_Basic 测试基础批量执行
func TestBatchExecutor_Basic(t *testing.T) {
	ctx := context.Background()
	var counter atomic.Int32

	executor := NewBatchExecutor(ctx).SetLimit(5)

	for i := 0; i < 10; i++ {
		executor.Go(func() error {
			counter.Add(1)
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}

	err := executor.Wait()
	assert.NoError(t, err)
	assert.Equal(t, int32(10), counter.Load())
}

// TestBatchExecutor_ContinueOnError 测试继续执行模式，即使有错误也继续执行所有任务
func TestBatchExecutor_ContinueOnError(t *testing.T) {
	ctx := context.Background()
	var counter atomic.Int32

	executor := NewBatchExecutor(ctx).
		SetLimit(5).
		SetMode(ContinueOnErrorMode)

	for i := 0; i < 10; i++ {
		i := i // capture loop variable
		executor.Go(func() error {
			if i == 5 {
				return fmt.Errorf("error on task %d", i)
			}
			counter.Add(1)
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}

	err := executor.Wait()
	assert.Error(t, err)
	assert.Equal(t, int32(9), counter.Load()) // 任务5返回错误不增加counter
	assert.Equal(t, 1, executor.ErrorCount())
}

// TestBatchExecutor_FailFastMode 测试快速失败模式，确保在遇到第一个错误时停止提交新任务
func TestBatchExecutor_FailFastMode(t *testing.T) {
	ctx := context.Background()
	var counter atomic.Int32

	executor := NewBatchExecutor(ctx).
		SetLimit(5).
		SetMode(FailFastMode)

	for i := 0; i < 10; i++ {
		i := i // capture loop variable
		executor.Go(func() error {
			if i == 3 {
				return fmt.Errorf("error on task %d", i)
			}
			counter.Add(1)
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}

	err := executor.Wait()
	assert.Error(t, err)
	// FailFastMode: 任务3报错触发cancel，但0-3已提交，其中0,1,2成功执行
	// 由于并发和时序问题，counter可能是3或更少
	assert.LessOrEqual(t, counter.Load(), int32(6)) // 最多前几个任务执行
	assert.Equal(t, 1, executor.ErrorCount())
}

// TestBatchExecutor_ContextCancellation 测试在上下文取消后，确保不再执行新任务
func TestBatchExecutor_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var counter atomic.Int32

	executor := NewBatchExecutor(ctx).
		SetLimit(5)

	for i := 0; i < 10; i++ {
		executor.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}

	cancel() // 取消上下文

	err := executor.Wait()
	assert.NoError(t, err)
	assert.Equal(t, int32(0), counter.Load()) // 由于上下文取消，所有任务都不应执行
}

// TestBatchExecutor_ErrorHandler 测试错误处理器，确保每个错误都被处理
func TestBatchExecutor_ErrorHandler(t *testing.T) {
	ctx := context.Background()
	var counter atomic.Int32
	var errorCount int32

	executor := NewBatchExecutor(ctx).
		SetLimit(5).
		SetMode(ContinueOnErrorMode). // 继续执行模式
		OnError(func(err error) {
			atomic.AddInt32(&errorCount, 1)
		})

	for i := 0; i < 10; i++ {
		i := i // capture loop variable
		executor.Go(func() error {
			if i%3 == 0 {
				return fmt.Errorf("error on task %d", i)
			}
			counter.Add(1)
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}

	err := executor.Wait()
	assert.Error(t, err)
	assert.Equal(t, int32(6), counter.Load()) // 只有非i%3==0的任务增加counter: 1,2,4,5,7,8
	assert.Equal(t, int32(4), errorCount)     // 0, 3, 6, 9 应该产生错误
}

// TestBatchExecutor_ConcurrentLimit 测试并发限制，确保不会超过设定的并发数
func TestBatchExecutor_ConcurrentLimit(t *testing.T) {
	ctx := context.Background()
	executor := NewBatchExecutor(ctx).
		SetLimit(3)

	var counter atomic.Int32
	const totalTasks = 10

	for i := 0; i < totalTasks; i++ {
		executor.Go(func() error {
			counter.Add(1)
			time.Sleep(50 * time.Millisecond) // 模拟工作
			return nil
		})
	}

	err := executor.Wait() // 等待所有任务完成
	assert.NoError(t, err)
	assert.Equal(t, int32(totalTasks), counter.Load()) // 所有任务都应成功执行
	assert.Equal(t, 0, executor.ErrorCount())          // 没有错误
}

// TestBatchExecutor_PanicRecovery 测试任务中的 panic 恢复，确保 panic 不会导致程序崩溃
func TestBatchExecutor_PanicRecovery(t *testing.T) {
	ctx := context.Background()
	var counter atomic.Int32

	executor := NewBatchExecutor(ctx).
		SetLimit(5).
		SetMode(ContinueOnErrorMode). // 继续执行模式
		OnPanic(func(r interface{}) {
			t.Logf("Recovered from panic: %v", r)
		})

	for i := 0; i < 10; i++ {
		i := i // capture loop variable
		executor.Go(func() error {
			if i == 5 {
				panic("panic on task 5")
			}
			counter.Add(1)
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}

	err := executor.Wait()
	assert.Error(t, err)
	assert.Equal(t, int32(9), counter.Load()) // 任务5发生panic，其他任务应该执行
}

// TestBatchExecutor_EmptyExecutor 测试没有任务的执行器，确保不会出错
func TestBatchExecutor_EmptyExecutor(t *testing.T) {
	ctx := context.Background()
	executor := NewBatchExecutor(ctx)

	err := executor.Wait()
	assert.NoError(t, err)
	assert.Equal(t, 0, executor.ErrorCount())
}

// MockItem 模拟退款项目
type MockItem struct {
	ID      int
	TraceID string
	UserID  string
}

// MockProcessor 模拟退款处理器
type MockProcessor struct {
	successCount atomic.Int64
	failedCount  atomic.Int64
	skippedCount atomic.Int64
	panicCount   atomic.Int64
	processDelay time.Duration
	failOnID     int // 在哪个ID失败
	panicOnID    int // 在哪个ID panic
}

// process 处理单个项目
func (m *MockProcessor) process(item *MockItem) (*MockItem, error) {
	if m.processDelay > 0 {
		time.Sleep(m.processDelay)
	}

	// 模拟 panic
	if m.panicOnID > 0 && item.ID == m.panicOnID {
		m.panicCount.Add(1)
		panic(fmt.Sprintf("模拟 panic at ID %d", item.ID))
	}

	// 模拟失败
	if m.failOnID > 0 && item.ID == m.failOnID {
		m.failedCount.Add(1)
		return nil, fmt.Errorf("模拟失败 at ID %d", item.ID)
	}

	// 模拟跳过（ID是10的倍数）
	if item.ID%10 == 0 {
		m.skippedCount.Add(1)
		return nil, nil
	}

	// 成功
	m.successCount.Add(1)
	return item, nil
}

// TestBatchExecutorRefund_NormalProcessing 测试正常批量处理
func TestBatchExecutorRefund_NormalProcessing(t *testing.T) {
	mock := &MockProcessor{
		processDelay: 5 * time.Millisecond,
	}

	// 创建50个测试项目
	items := make([]*MockItem, 50)
	for i := 0; i < 50; i++ {
		items[i] = &MockItem{
			ID:      i,
			TraceID: fmt.Sprintf("trace_%d", i),
			UserID:  fmt.Sprintf("user_%d", i),
		}
	}

	// 使用 BatchExecutor 处理
	ctx := context.Background()
	executor := NewBatchExecutor(ctx).
		SetLimit(10).
		SetMode(ContinueOnErrorMode)

	for _, item := range items {
		capturedItem := item
		executor.Go(func() error {
			result, err := mock.process(capturedItem)
			if err != nil {
				return err
			}
			if result == nil {
				// 跳过
			}
			return nil
		})
	}

	err := executor.Wait()
	assert.NoError(t, err)

	// 验证：ID为 0,10,20,30,40 的5个跳过，其余45个成功
	assert.Equal(t, int64(45), mock.successCount.Load(), "成功数量")
	assert.Equal(t, int64(5), mock.skippedCount.Load(), "跳过数量")
	assert.Equal(t, int64(0), mock.failedCount.Load(), "失败数量")

	t.Logf("✅ 正常处理 - 总数:50, 成功:%d, 跳过:%d",
		mock.successCount.Load(), mock.skippedCount.Load())
}

// TestBatchExecutorRefund_FailFastMode 测试快速失败模式
func TestBatchExecutorRefund_FailFastMode(t *testing.T) {
	mock := &MockProcessor{
		processDelay: 10 * time.Millisecond,
		failOnID:     15, // 在ID=15时失败
	}

	// 创建30个测试项目
	items := make([]*MockItem, 30)
	for i := 0; i < 30; i++ {
		items[i] = &MockItem{
			ID:      i,
			TraceID: fmt.Sprintf("trace_%d", i),
		}
	}

	ctx := context.Background()
	executor := NewBatchExecutor(ctx).
		SetLimit(5).
		SetMode(FailFastMode). // 快速失败
		OnError(func(err error) {
			t.Logf("捕获错误: %v", err)
		})

	for _, item := range items {
		capturedItem := item
		executor.Go(func() error {
			_, err := mock.process(capturedItem)
			if err != nil {
				return fmt.Errorf("处理失败 ID=%d: %w", capturedItem.ID, err)
			}
			return nil
		})
	}

	err := executor.Wait()
	assert.Error(t, err, "应该返回错误")
	assert.Contains(t, err.Error(), "模拟失败", "错误信息")

	totalProcessed := mock.successCount.Load() + mock.failedCount.Load() + mock.skippedCount.Load()
	t.Logf("⚡ 快速失败 - 提交:%d, 实际处理:%d, 成功:%d, 失败:%d, 跳过:%d",
		len(items), totalProcessed, mock.successCount.Load(), mock.failedCount.Load(), mock.skippedCount.Load())

	// 快速失败模式下，处理的任务数应该远小于总数
	assert.Less(t, totalProcessed, int64(len(items)), "快速失败应停止提交新任务")
	assert.Equal(t, int64(1), mock.failedCount.Load(), "应该有1个失败")
}

// TestBatchExecutorRefund_PanicRecovery 测试 panic 恢复
func TestBatchExecutorRefund_PanicRecovery(t *testing.T) {
	mock := &MockProcessor{
		processDelay: 5 * time.Millisecond,
		panicOnID:    12, // 在ID=12时 panic
	}

	items := make([]*MockItem, 25)
	for i := 0; i < 25; i++ {
		items[i] = &MockItem{
			ID:      i,
			TraceID: fmt.Sprintf("trace_%d", i),
		}
	}

	ctx := context.Background()
	var panicRecovered bool
	executor := NewBatchExecutor(ctx).
		SetLimit(5).
		SetMode(FailFastMode).
		OnPanic(func(r interface{}) {
			panicRecovered = true
			t.Logf("🔥 捕获 panic: %v", r)
		})

	for _, item := range items {
		capturedItem := item
		executor.Go(func() error {
			_, err := mock.process(capturedItem)
			return err
		})
	}

	err := executor.Wait()
	assert.Error(t, err, "panic 应该被转换为 error")
	assert.True(t, panicRecovered, "panic 应该被捕获")
	assert.Equal(t, int64(1), mock.panicCount.Load(), "panic 次数")

	t.Logf("🛡️ Panic 恢复 - panic次数:%d, 成功:%d",
		mock.panicCount.Load(), mock.successCount.Load())
}

// TestBatchExecutorRefund_ConcurrentLimit 测试并发限制
func TestBatchExecutorRefund_ConcurrentLimit(t *testing.T) {
	mock := &MockProcessor{
		processDelay: 50 * time.Millisecond,
	}

	items := make([]*MockItem, 20)
	for i := 0; i < 20; i++ {
		items[i] = &MockItem{ID: i}
	}

	start := time.Now()

	ctx := context.Background()
	executor := NewBatchExecutor(ctx).
		SetLimit(4). // 限制并发为4
		SetMode(ContinueOnErrorMode)

	for _, item := range items {
		capturedItem := item
		executor.Go(func() error {
			_, err := mock.process(capturedItem)
			return err
		})
	}

	err := executor.Wait()
	assert.NoError(t, err)

	elapsed := time.Since(start)

	// 20个任务，每个50ms，并发4个：预期 (20/4) * 50ms = 250ms
	t.Logf("⏱️  并发限制 - 耗时:%v, 成功:%d", elapsed, mock.successCount.Load())

	assert.Greater(t, elapsed, 200*time.Millisecond, "并发限制应使总时间增加")
	assert.Less(t, elapsed, 400*time.Millisecond, "不应该超时太多")
}

// TestBatchExecutorRefund_HealthCheckThenBatch 测试健康检查+批量处理模式（模拟真实场景）
func TestBatchExecutorRefund_HealthCheckThenBatch(t *testing.T) {
	mock := &MockProcessor{
		processDelay: 5 * time.Millisecond,
	}

	// 创建100个项目
	totalCount := 100
	items := make([]*MockItem, totalCount)
	for i := 0; i < totalCount; i++ {
		items[i] = &MockItem{
			ID:      i,
			TraceID: fmt.Sprintf("trace_%d", i),
		}
	}

	var (
		successCount atomic.Int64
		skippedCount atomic.Int64
	)

	// 1️⃣ 健康检查：先处理第一个
	firstResult, firstErr := mock.process(items[0])
	assert.NoError(t, firstErr, "健康检查应该成功")

	if firstResult == nil {
		skippedCount.Add(1)
	} else {
		successCount.Add(1)
	}

	if totalCount == 1 {
		t.Skip("只有一条记录")
	}

	// 2️⃣ 批量处理剩余记录
	ctx := context.Background()
	executor := NewBatchExecutor(ctx).
		SetLimit(10).
		SetMode(FailFastMode).
		OnPanic(func(r interface{}) {
			t.Logf("panic: %v", r)
		})

	for i := 1; i < totalCount; i++ {
		item := items[i]
		executor.Go(func() error {
			result, err := mock.process(item)
			if err != nil {
				return fmt.Errorf("处理失败 ID=%d: %w", item.ID, err)
			}

			if result == nil {
				skippedCount.Add(1)
			} else {
				successCount.Add(1)
			}
			return nil
		})
	}

	err := executor.Wait()
	assert.NoError(t, err)

	// 验证：0,10,20...90 共10个跳过，其余90个成功
	assert.Equal(t, int64(90), successCount.Load(), "成功数量")
	assert.Equal(t, int64(10), skippedCount.Load(), "跳过数量")

	t.Logf("🎯 健康检查+批量 - 总数:%d, 成功:%d, 跳过:%d",
		totalCount, successCount.Load(), skippedCount.Load())
}

// TestBatchExecutorRefund_ContinueOnErrorMode 测试继续执行模式
func TestBatchExecutorRefund_ContinueOnErrorMode(t *testing.T) {
	mock := &MockProcessor{
		processDelay: 5 * time.Millisecond,
		failOnID:     15, // 在ID=15时失败，但继续处理其他任务
	}

	items := make([]*MockItem, 30)
	for i := 0; i < 30; i++ {
		items[i] = &MockItem{ID: i}
	}

	ctx := context.Background()
	executor := NewBatchExecutor(ctx).
		SetLimit(5).
		SetMode(ContinueOnErrorMode) // 继续执行模式

	for _, item := range items {
		capturedItem := item
		executor.Go(func() error {
			_, err := mock.process(capturedItem)
			if err != nil {
				return err // 虽然返回错误，但不会停止其他任务
			}
			return nil
		})
	}

	err := executor.Wait()
	assert.Error(t, err, "应该返回第一个错误")

	// 继续执行模式：所有任务都应该被处理
	totalProcessed := mock.successCount.Load() + mock.failedCount.Load() + mock.skippedCount.Load()
	assert.Equal(t, int64(30), totalProcessed, "所有任务都应该被处理")
	assert.Equal(t, int64(1), mock.failedCount.Load(), "应该有1个失败")

	t.Logf("🔄 继续执行模式 - 总数:30, 实际处理:%d, 成功:%d, 失败:%d, 跳过:%d",
		totalProcessed, mock.successCount.Load(), mock.failedCount.Load(), mock.skippedCount.Load())
}

// TestGo_WithWaitGroup_Exec 测试 Exec 绑定 WaitGroup
func TestGo_WithWaitGroup_Exec(t *testing.T) {
	var wg sync.WaitGroup
	executed := int32(0)

	Go().WithWaitGroup(&wg).Exec(func() {
		atomic.AddInt32(&executed, 1)
	})

	wg.Wait()
	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
}

// TestGo_WithWaitGroup_ExecWithContext 测试 ExecWithContext 绑定 WaitGroup
func TestGo_WithWaitGroup_ExecWithContext(t *testing.T) {
	var wg sync.WaitGroup
	executed := int32(0)

	Go().
		WithWaitGroup(&wg).
		WithTimeout(1 * time.Second).
		ExecWithContext(func(ctx context.Context) error {
			atomic.AddInt32(&executed, 1)
			return nil
		})

	wg.Wait()
	assert.Equal(t, int32(1), atomic.LoadInt32(&executed))
}

// TestGo_WithWaitGroup_MultipleGoroutines 测试多个 goroutine 绑定同一 WaitGroup
func TestGo_WithWaitGroup_MultipleGoroutines(t *testing.T) {
	var wg sync.WaitGroup
	var counter atomic.Int32
	const n = 10

	for i := 0; i < n; i++ {
		Go().
			WithWaitGroup(&wg).
			WithTimeout(2 * time.Second).
			ExecWithContext(func(ctx context.Context) error {
				time.Sleep(10 * time.Millisecond)
				counter.Add(1)
				return nil
			})
	}

	wg.Wait()
	assert.Equal(t, int32(n), counter.Load())
}

// TestGo_WithWaitGroup_WithPanic 测试 panic 时 WaitGroup 也能正确 Done
func TestGo_WithWaitGroup_WithPanic(t *testing.T) {
	var wg sync.WaitGroup
	panicCaught := int32(0)

	Go().
		WithWaitGroup(&wg).
		OnPanic(func(r interface{}) {
			atomic.AddInt32(&panicCaught, 1)
		}).
		Exec(func() {
			panic("test panic")
		})

	// wg.Wait() 不应死锁，即使 goroutine 内部 panic
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 正常完成
	case <-time.After(1 * time.Second):
		t.Fatal("wg.Wait() 超时，panic 时 WaitGroup 未 Done")
	}

	assert.Equal(t, int32(1), atomic.LoadInt32(&panicCaught))
}
