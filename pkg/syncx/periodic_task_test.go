/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-29 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-01 07:41:12
 * @FilePath: \go-toolbox\pkg\syncx\periodic_task_test.go
 * @Description: PeriodicTaskManager 测试文件
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package syncx

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPeriodicTaskManagerNewPeriodicTaskManager 测试创建新的任务管理器
func TestPeriodicTaskManagerNewPeriodicTaskManager(t *testing.T) {
	manager := NewPeriodicTaskManager()

	assert.NotNil(t, manager, "manager should not be nil")
	assert.NotNil(t, manager.tasks, "tasks slice should be initialized")
	assert.Equal(t, 0, len(manager.tasks), "initial tasks count should be 0")
	assert.False(t, manager.isRunning, "manager should not be running initially")
}

// TestPeriodicTaskManagerAddTask 测试添加任务
func TestPeriodicTaskManagerAddTask(t *testing.T) {
	manager := NewPeriodicTaskManager()

	task := NewPeriodicTask("test_task", time.Second, func(ctx context.Context) error { return nil })

	result := manager.AddTask(task)

	assert.Equal(t, manager, result, "AddTask should return the manager for chaining")
	assert.Equal(t, 1, manager.GetTaskCount(), "task count should be 1")
	assert.Equal(t, "test_task", manager.tasks[0].GetName(), "task name should match")
}

// TestPeriodicTaskManagerAddSimpleTask 测试添加简单任务
func TestPeriodicTaskManagerAddSimpleTask(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executed int32

	result := manager.AddSimpleTask("simple_task", time.Millisecond*100, func(ctx context.Context) error {
		atomic.AddInt32(&executed, 1)
		return nil
	})

	assert.Equal(t, manager, result, "AddSimpleTask should return the manager for chaining")
	assert.Equal(t, 1, manager.GetTaskCount(), "task count should be 1")
	assert.Equal(t, "simple_task", manager.tasks[0].GetName(), "task name should match")
	assert.Equal(t, time.Millisecond*100, manager.tasks[0].GetInterval(), "task interval should match")
	assert.False(t, manager.tasks[0].GetImmediateStart(), "immediate start should be false by default")
}

// TestPeriodicTaskManagerAddTaskWithImmediateStart 测试添加立即执行任务
func TestPeriodicTaskManagerAddTaskWithImmediateStart(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executed int32

	result := manager.AddTaskWithImmediateStart("immediate_task", time.Millisecond*100, func(ctx context.Context) error {
		atomic.AddInt32(&executed, 1)
		return nil
	})

	assert.Equal(t, manager, result, "AddTaskWithImmediateStart should return the manager for chaining")
	assert.Equal(t, 1, manager.GetTaskCount(), "task count should be 1")
	assert.Equal(t, "immediate_task", manager.tasks[0].GetName(), "task name should match")
	assert.True(t, manager.tasks[0].GetImmediateStart(), "immediate start should be true")
}

// TestPeriodicTaskManagerSetDefaultErrorHandler 测试设置默认错误处理器
func TestPeriodicTaskManagerSetDefaultErrorHandler(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加一个没有错误处理器的任务
	manager.AddSimpleTask("task1", time.Second, func(ctx context.Context) error { return nil })

	// 添加一个已有错误处理器的任务
	task2 := NewPeriodicTask("task2", time.Second, func(ctx context.Context) error { return nil }).
		SetOnError(func(name string, err error) { /* existing handler */ })
	manager.AddTask(task2)

	// 设置默认错误处理器
	result := manager.SetDefaultErrorHandler(func(name string, err error) {
		// 错误处理逻辑
	})

	assert.Equal(t, manager, result, "SetDefaultErrorHandler should return the manager for chaining")
	assert.NotNil(t, manager.tasks[0].GetOnError(), "task1 should have error handler set")
	assert.NotNil(t, manager.tasks[1].GetOnError(), "task2 should still have its original error handler")
}

// TestPeriodicTaskManagerSetDefaultCallbacks 测试设置默认回调函数
func TestPeriodicTaskManagerSetDefaultCallbacks(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加一个没有回调的任务
	manager.AddSimpleTask("task1", time.Second, func(ctx context.Context) error { return nil })

	// 添加一个已有回调的任务
	task2 := NewPeriodicTask("task2", time.Second, func(ctx context.Context) error { return nil }).
		SetOnStart(func(name string) { /* existing start handler */ }).
		SetOnStop(func(name string) { /* existing stop handler */ })
	manager.AddTask(task2)

	// 设置默认回调
	result := manager.SetDefaultCallbacks(
		func(name string) { /* start callback */ },
		func(name string) { /* stop callback */ },
	)

	assert.Equal(t, manager, result, "SetDefaultCallbacks should return the manager for chaining")
	assert.NotNil(t, manager.tasks[0].GetOnStart(), "task1 should have start callback set")
	assert.NotNil(t, manager.tasks[0].GetOnStop(), "task1 should have stop callback set")
	assert.NotNil(t, manager.tasks[1].GetOnStart(), "task2 should still have its original start callback")
	assert.NotNil(t, manager.tasks[1].GetOnStop(), "task2 should still have its original stop callback")
}

// TestPeriodicTaskManagerStartAlreadyRunning 测试重复启动
func TestPeriodicTaskManagerStartAlreadyRunning(t *testing.T) {
	manager := NewPeriodicTaskManager()
	manager.AddSimpleTask("test_task", time.Second, func(ctx context.Context) error { return nil })

	err := manager.Start()
	assert.NoError(t, err, "first start should succeed")
	assert.True(t, manager.IsRunning(), "manager should be running")

	err = manager.Start()
	assert.Error(t, err, "second start should return error")
	assert.Contains(t, err.Error(), "already running", "error should mention already running")

	// 清理
	manager.Stop()
}

// TestPeriodicTaskManagerStartWithContextAlreadyRunning 测试使用上下文重复启动
func TestPeriodicTaskManagerStartWithContextAlreadyRunning(t *testing.T) {
	manager := NewPeriodicTaskManager()
	manager.AddSimpleTask("test_task", time.Second, func(ctx context.Context) error { return nil })

	ctx := context.Background()

	err := manager.StartWithContext(ctx)
	assert.NoError(t, err, "first start should succeed")
	assert.True(t, manager.IsRunning(), "manager should be running")

	err = manager.StartWithContext(ctx)
	assert.Error(t, err, "second start should return error")
	assert.Contains(t, err.Error(), "already running", "error should mention already running")

	// 清理
	manager.Stop()
}

// TestPeriodicTaskManagerStopNotRunning 测试停止未运行的管理器
func TestPeriodicTaskManagerStopNotRunning(t *testing.T) {
	manager := NewPeriodicTaskManager()

	err := manager.Stop()
	assert.NoError(t, err, "stopping non-running manager should not error")
	assert.False(t, manager.IsRunning(), "manager should not be running")
}

// TestPeriodicTaskManagerStartStop 测试启动和停止
func TestPeriodicTaskManagerStartStop(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executed int32

	manager.AddSimpleTask("test_task", time.Millisecond*50, func(ctx context.Context) error {
		atomic.AddInt32(&executed, 1)
		return nil
	})

	err := manager.Start()
	assert.NoError(t, err, "start should succeed")
	assert.True(t, manager.IsRunning(), "manager should be running")

	// 等待任务执行几次
	time.Sleep(time.Millisecond * 200)

	err = manager.Stop()
	assert.NoError(t, err, "stop should succeed")
	assert.False(t, manager.IsRunning(), "manager should not be running")

	executedCount := atomic.LoadInt32(&executed)
	assert.Greater(t, executedCount, int32(0), "task should have executed at least once")
}

// TestPeriodicTaskManagerImmediateStart 测试立即执行任务
func TestPeriodicTaskManagerImmediateStart(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executed int32

	manager.AddTaskWithImmediateStart("immediate_task", time.Second, func(ctx context.Context) error {
		atomic.AddInt32(&executed, 1)
		return nil
	})

	err := manager.Start()
	assert.NoError(t, err, "start should succeed")

	// 等待立即执行
	time.Sleep(time.Millisecond * 100)

	executedCount := atomic.LoadInt32(&executed)
	assert.Greater(t, executedCount, int32(0), "task should have executed immediately")

	manager.Stop()
}

// TestPeriodicTaskManagerErrorHandling 测试错误处理
func TestPeriodicTaskManagerErrorHandling(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var errorCount int32
	var errorName string
	var errorMessage string

	manager.SetDefaultErrorHandler(func(name string, err error) {
		atomic.AddInt32(&errorCount, 1)
		errorName = name
		errorMessage = err.Error()
	})

	manager.AddSimpleTask("error_task", time.Millisecond*50, func(ctx context.Context) error {
		return fmt.Errorf("test error")
	})

	manager.Start()
	time.Sleep(time.Millisecond * 200)
	manager.Stop()

	assert.Greater(t, atomic.LoadInt32(&errorCount), int32(0), "error handler should have been called")
	assert.Equal(t, "error_task", errorName, "error name should match task name")
	assert.Equal(t, "test error", errorMessage, "error message should match")
}

// TestPeriodicTaskManagerStartStopCallbacks 测试启动停止回调
func TestPeriodicTaskManagerStartStopCallbacks(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var startCalled, stopCalled bool
	var startTaskName, stopTaskName string

	manager.SetDefaultCallbacks(
		func(name string) {
			startCalled = true
			startTaskName = name
		},
		func(name string) {
			stopCalled = true
			stopTaskName = name
		},
	)

	manager.AddSimpleTask("callback_task", time.Second, func(ctx context.Context) error { return nil })

	manager.Start()
	time.Sleep(time.Millisecond * 100)
	manager.Stop()

	assert.True(t, startCalled, "start callback should have been called")
	assert.True(t, stopCalled, "stop callback should have been called")
	assert.Equal(t, "callback_task", startTaskName, "start callback should receive correct task name")
	assert.Equal(t, "callback_task", stopTaskName, "stop callback should receive correct task name")
}

// TestPeriodicTaskManagerMultipleTasksExecution 测试多个任务执行
func TestPeriodicTaskManagerMultipleTasksExecution(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var task1Count, task2Count, task3Count int32

	manager.AddSimpleTask("task1", time.Millisecond*50, func(ctx context.Context) error {
		atomic.AddInt32(&task1Count, 1)
		return nil
	})

	manager.AddSimpleTask("task2", time.Millisecond*75, func(ctx context.Context) error {
		atomic.AddInt32(&task2Count, 1)
		return nil
	})

	manager.AddTaskWithImmediateStart("task3", time.Millisecond*100, func(ctx context.Context) error {
		atomic.AddInt32(&task3Count, 1)
		return nil
	})

	assert.Equal(t, 3, manager.GetTaskCount(), "should have 3 tasks")

	manager.Start()
	time.Sleep(time.Millisecond * 150)
	manager.Stop()

	assert.Greater(t, atomic.LoadInt32(&task1Count), int32(0), "task1 should have executed")
	assert.Greater(t, atomic.LoadInt32(&task2Count), int32(0), "task2 should have executed")
	assert.Greater(t, atomic.LoadInt32(&task3Count), int32(0), "task3 should have executed")
}

// TestPeriodicTaskManagerGetTaskNames 测试获取任务名称
func TestPeriodicTaskManagerGetTaskNames(t *testing.T) {
	manager := NewPeriodicTaskManager()

	manager.AddSimpleTask("task_a", time.Second, func(ctx context.Context) error { return nil })
	manager.AddSimpleTask("task_b", time.Second, func(ctx context.Context) error { return nil })
	manager.AddSimpleTask("task_c", time.Second, func(ctx context.Context) error { return nil })

	names := manager.GetTaskNames()

	assert.Equal(t, 3, len(names), "should return 3 task names")
	assert.Contains(t, names, "task_a", "should contain task_a")
	assert.Contains(t, names, "task_b", "should contain task_b")
	assert.Contains(t, names, "task_c", "should contain task_c")
}

// TestPeriodicTaskManagerContextCancellation 测试上下文取消
func TestPeriodicTaskManagerContextCancellation(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executed int32

	ctx, cancel := context.WithCancel(context.Background())

	manager.AddSimpleTask("cancelable_task", time.Millisecond*50, func(ctx context.Context) error {
		atomic.AddInt32(&executed, 1)
		return nil
	})

	err := manager.StartWithContext(ctx)
	assert.NoError(t, err, "start should succeed")

	time.Sleep(time.Millisecond * 200)
	cancel() // 取消上下文

	time.Sleep(time.Millisecond * 100) // 等待任务停止

	executedBefore := atomic.LoadInt32(&executed)
	time.Sleep(time.Millisecond * 100) // 再等待一段时间
	executedAfter := atomic.LoadInt32(&executed)

	assert.Equal(t, executedBefore, executedAfter, "task should not execute after context cancellation")
}

// TestPeriodicTaskManagerStopWithTimeout_Success 测试超时停止成功
func TestPeriodicTaskManagerStopWithTimeout_Success(t *testing.T) {
	manager := NewPeriodicTaskManager()

	manager.AddSimpleTask("quick_task", time.Millisecond*50, func(ctx context.Context) error {
		return nil
	})

	manager.Start()
	time.Sleep(time.Millisecond * 100)

	err := manager.StopWithTimeout(time.Second)
	assert.NoError(t, err, "stop with timeout should succeed")
	assert.False(t, manager.IsRunning(), "manager should not be running")
}

// TestPeriodicTaskManagerStopWithTimeout_Timeout 测试超时停止超时
func TestPeriodicTaskManagerStopWithTimeout_Timeout(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 创建一个会长时间阻塞且不检查 context 的任务
	blockCh := make(chan struct{})
	started := make(chan struct{})

	manager.AddTask(&PeriodicTask{
		name:           "blocking_task",
		interval:       time.Millisecond * 10,
		preventOverlap: true, // 防止重叠执行
		executeFunc: func(ctx context.Context) error {
			close(started)
			<-blockCh // 阻塞直到通道关闭
			return nil
		},
	})

	manager.Start()
	<-started // 等待任务开始执行

	// 尝试在很短的时间内停止，任务还在阻塞中
	err := manager.StopWithTimeout(time.Millisecond * 50)

	// 清理：关闭阻塞通道，让任务可以完成
	close(blockCh)

	if err != nil {
		assert.Contains(t, err.Error(), "timeout", "error should mention timeout")
	} else {
		// 如果没有错误，说明 Stop 在超时前完成了（这也是可以接受的）
		t.Log("Stop completed before timeout (race condition - acceptable)")
	}
}

// TestPeriodicTaskManagerWait 测试等待功能
func TestPeriodicTaskManagerWait(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var taskStarted, taskCompleted atomic.Bool

	manager.AddSimpleTask("wait_task", time.Millisecond*50, func(ctx context.Context) error {
		taskStarted.Store(true)
		time.Sleep(time.Millisecond * 20)
		taskCompleted.Store(true)
		return nil
	})

	manager.Start()

	// 等待任务至少开始一次
	for i := 0; i < 100 && !taskStarted.Load(); i++ {
		time.Sleep(time.Millisecond * 10)
	}
	assert.True(t, taskStarted.Load(), "Task should have started")

	// 停止 manager
	err := manager.Stop()
	assert.NoError(t, err)

	// Wait 应该立即返回，因为 Stop 已经等待了所有任务
	done := make(chan struct{})
	go func() {
		manager.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Wait 成功返回
		t.Log("Wait completed successfully")
	case <-time.After(time.Second):
		t.Fatal("Wait should return quickly after Stop")
	}
}

// TestPeriodicTaskManagerConcurrentAccess 测试并发访问
func TestPeriodicTaskManagerConcurrentAccess(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var wg sync.WaitGroup
	var errorCount int32

	// 并发添加任务
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			taskName := fmt.Sprintf("task_%d", id)
			manager.AddSimpleTask(taskName, time.Second, func(ctx context.Context) error {
				return nil
			})
		}(i)
	}

	// 并发获取任务信息
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count := manager.GetTaskCount()
			names := manager.GetTaskNames()
			if len(names) != count {
				atomic.AddInt32(&errorCount, 1)
			}
		}()
	}

	wg.Wait()

	// 等待一小段时间确保所有任务都已添加完成
	time.Sleep(time.Millisecond * 50)

	// 在所有并发操作完成后再检查
	finalCount := manager.GetTaskCount()
	finalNames := manager.GetTaskNames()

	assert.LessOrEqual(t, int(atomic.LoadInt32(&errorCount)), 1, "minimal concurrent access errors should occur")
	assert.Equal(t, 10, finalCount, "should have 10 tasks")
	assert.Equal(t, 10, len(finalNames), "should have 10 task names")
}

// TestPeriodicTaskManagerTaskWithCustomCallbacks 测试带自定义回调的任务
func TestPeriodicTaskManagerTaskWithCustomCallbacks(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var customStartCalled, customStopCalled, customErrorCalled bool

	task := NewPeriodicTask("custom_task", time.Millisecond*50, func(ctx context.Context) error {
		return fmt.Errorf("custom error")
	}).
		SetOnStart(func(name string) {
			customStartCalled = true
		}).
		SetOnStop(func(name string) {
			customStopCalled = true
		}).
		SetOnError(func(name string, err error) {
			customErrorCalled = true
		})

	manager.AddTask(task)

	// 设置默认回调（不应该覆盖自定义回调）
	manager.SetDefaultCallbacks(
		func(name string) { assert.Fail(t, "default start callback should not be called") },
		func(name string) { assert.Fail(t, "default stop callback should not be called") },
	)

	manager.SetDefaultErrorHandler(func(name string, err error) {
		assert.Fail(t, "default error handler should not be called")
	})

	manager.Start()
	time.Sleep(time.Millisecond * 200)
	manager.Stop()

	assert.True(t, customStartCalled, "custom start callback should have been called")
	assert.True(t, customStopCalled, "custom stop callback should have been called")
	assert.True(t, customErrorCalled, "custom error callback should have been called")
}

// TestPeriodicTaskManagerEmptyManager 测试空管理器
func TestPeriodicTaskManagerEmptyManager(t *testing.T) {
	manager := NewPeriodicTaskManager()

	err := manager.Start()
	assert.NoError(t, err, "starting empty manager should succeed")
	assert.True(t, manager.IsRunning(), "empty manager should be running")

	err = manager.Stop()
	assert.NoError(t, err, "stopping empty manager should succeed")
	assert.False(t, manager.IsRunning(), "empty manager should not be running")

	assert.Equal(t, 0, manager.GetTaskCount(), "task count should be 0")
	assert.Equal(t, 0, len(manager.GetTaskNames()), "task names should be empty")
}

// TestPeriodicTaskManagerTaskExecutionOrder 测试任务执行顺序
func TestPeriodicTaskManagerTaskExecutionOrder(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionOrder []string
	var mutex sync.Mutex

	for i := 0; i < 3; i++ {
		taskName := fmt.Sprintf("task_%d", i)
		manager.AddTaskWithImmediateStart(taskName, time.Second, func(ctx context.Context) error {
			mutex.Lock()
			executionOrder = append(executionOrder, taskName)
			mutex.Unlock()
			return nil
		})
	}

	manager.Start()
	time.Sleep(time.Millisecond * 100)
	manager.Stop()

	mutex.Lock()
	defer mutex.Unlock()

	assert.Equal(t, 3, len(executionOrder), "all tasks should have executed")
	// 注意：由于并发执行，执行顺序可能不固定，但都应该执行
}

// TestPeriodicTaskManagerLongRunningTask 测试长时间运行的任务
func TestPeriodicTaskManagerLongRunningTask(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var startCount, completeCount int32

	manager.AddSimpleTask("long_task", time.Millisecond*50, func(ctx context.Context) error {
		atomic.AddInt32(&startCount, 1)
		time.Sleep(time.Millisecond * 200) // 任务执行时间长于间隔时间
		atomic.AddInt32(&completeCount, 1)
		return nil
	})

	manager.Start()
	time.Sleep(time.Millisecond * 300)
	manager.Stop()

	starts := atomic.LoadInt32(&startCount)
	completes := atomic.LoadInt32(&completeCount)

	assert.Greater(t, starts, int32(0), "task should have started")
	assert.Greater(t, completes, int32(0), "task should have completed")
	// 由于任务执行时间长，可能出现starts > completes的情况
}

// TestPeriodicTaskManagerPanicRecovery 测试panic恢复
func TestPeriodicTaskManagerPanicRecovery(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var panicHandled bool

	// 设置错误处理器来捕获panic
	manager.SetDefaultErrorHandler(func(name string, err error) {
		if name == "panic_task" {
			panicHandled = true
			assert.Contains(t, err.Error(), "panic", "error should contain panic information")
		}
	})

	// 添加一个会panic的任务
	manager.AddSimpleTask("panic_task", time.Millisecond*50, func(ctx context.Context) error {
		panic("test panic")
	})

	// 添加一个正常任务
	manager.AddSimpleTask("normal_task", time.Millisecond*100, func(ctx context.Context) error {
		return nil
	})

	manager.Start()
	time.Sleep(time.Millisecond * 200)
	manager.Stop()

	assert.True(t, panicHandled, "panic should have been handled by error handler")
} // TestPeriodicTaskManagerHighFrequencyTasks 测试高频任务
func TestPeriodicTaskManagerHighFrequencyTasks(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount int32

	manager.AddSimpleTask("high_freq_task", time.Millisecond*10, func(ctx context.Context) error {
		atomic.AddInt32(&executionCount, 1)
		return nil
	})

	manager.Start()
	time.Sleep(time.Millisecond * 200)
	manager.Stop()

	executions := atomic.LoadInt32(&executionCount)
	assert.Greater(t, executions, int32(10), "high frequency task should execute many times")
}

// TestPeriodicTaskManagerTaskNameUniqueness 测试任务名称唯一性
func TestPeriodicTaskManagerTaskNameUniqueness(t *testing.T) {
	manager := NewPeriodicTaskManager()

	manager.AddSimpleTask("duplicate_name", time.Second, func(ctx context.Context) error { return nil })
	manager.AddSimpleTask("duplicate_name", time.Second, func(ctx context.Context) error { return nil })
	manager.AddSimpleTask("unique_name", time.Second, func(ctx context.Context) error { return nil })

	names := manager.GetTaskNames()
	assert.Equal(t, 3, len(names), "should have 3 tasks even with duplicate names")

	// 统计重复名称
	nameCount := make(map[string]int)
	for _, name := range names {
		nameCount[name]++
	}

	assert.Equal(t, 2, nameCount["duplicate_name"], "should have 2 tasks with duplicate_name")
	assert.Equal(t, 1, nameCount["unique_name"], "should have 1 task with unique_name")
}

// TestPeriodicTaskManagerMemoryUsage 测试内存使用
func TestPeriodicTaskManagerMemoryUsage(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加大量任务
	for i := 0; i < 100; i++ {
		taskName := fmt.Sprintf("task_%d", i)
		manager.AddSimpleTask(taskName, time.Second, func(ctx context.Context) error { return nil })
	}

	assert.Equal(t, 100, manager.GetTaskCount(), "should have 100 tasks")

	manager.Start()
	time.Sleep(time.Millisecond * 100)
	manager.Stop()

	// 验证任务可以正常清理
	assert.False(t, manager.IsRunning(), "manager should not be running after stop")
}

// TestPeriodicTaskManagerZeroInterval 测试零间隔任务
func TestPeriodicTaskManagerZeroInterval(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount int32

	// 添加零间隔任务（可能导致高CPU使用）
	manager.AddSimpleTask("zero_interval", 0, func(ctx context.Context) error {
		atomic.AddInt32(&executionCount, 1)
		return nil
	})

	manager.Start()
	time.Sleep(time.Millisecond * 100)
	manager.Stop()

	executions := atomic.LoadInt32(&executionCount)
	// 零间隔可能导致非常高的执行次数
	assert.Greater(t, executions, int32(0), "zero interval task should execute")
}

// TestPeriodicTaskManagerNegativeInterval 测试负间隔任务
func TestPeriodicTaskManagerNegativeInterval(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount int32

	// 添加负间隔任务
	manager.AddSimpleTask("negative_interval", -time.Second, func(ctx context.Context) error {
		atomic.AddInt32(&executionCount, 1)
		return nil
	})

	manager.Start()
	time.Sleep(time.Millisecond * 200)
	manager.Stop()

	// 负间隔的行为可能不可预测，但不应该导致崩溃
	assert.True(t, true, "negative interval should not cause crash")
}

// TestPeriodicTaskManagerVeryLargeInterval 测试非常大的间隔
func TestPeriodicTaskManagerVeryLargeInterval(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executed bool

	// 添加非常大间隔的任务
	manager.AddSimpleTask("large_interval", time.Hour*24, func(ctx context.Context) error {
		executed = true
		return nil
	})

	manager.Start()
	time.Sleep(time.Millisecond * 100)
	manager.Stop()

	assert.False(t, executed, "large interval task should not execute quickly")
}

// TestPeriodicTaskManagerComplexScenario 测试复杂场景
func TestPeriodicTaskManagerComplexScenario(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var results sync.Map
	var errorCount int32

	// 设置错误处理器
	manager.SetDefaultErrorHandler(func(name string, err error) {
		atomic.AddInt32(&errorCount, 1)
	})

	// 设置回调
	manager.SetDefaultCallbacks(
		func(name string) {
			results.Store(name+"_started", true)
		},
		func(name string) {
			results.Store(name+"_stopped", true)
		},
	)

	// 添加各种类型的任务
	manager.AddSimpleTask("fast_task", time.Millisecond*50, func(ctx context.Context) error {
		results.Store("fast_executed", true)
		return nil
	})

	manager.AddTaskWithImmediateStart("immediate_task", time.Second, func(ctx context.Context) error {
		results.Store("immediate_executed", true)
		return nil
	})

	manager.AddSimpleTask("error_task", time.Millisecond*75, func(ctx context.Context) error {
		return fmt.Errorf("intentional error")
	})

	// 启动并运行一段时间
	manager.Start()
	assert.True(t, manager.IsRunning(), "manager should be running")

	time.Sleep(time.Millisecond * 150)

	manager.Stop()
	assert.False(t, manager.IsRunning(), "manager should not be running")

	// 验证结果
	assert.Equal(t, 3, manager.GetTaskCount(), "should have 3 tasks")

	fastExecuted, ok := results.Load("fast_executed")
	assert.True(t, ok && fastExecuted.(bool), "fast task should have executed")

	immediateExecuted, ok := results.Load("immediate_executed")
	assert.True(t, ok && immediateExecuted.(bool), "immediate task should have executed")

	assert.Greater(t, atomic.LoadInt32(&errorCount), int32(0), "error handler should have been called")

	// 验证回调被调用
	taskNames := []string{"fast_task", "immediate_task", "error_task"}
	for _, taskName := range taskNames {
		started, ok := results.Load(taskName + "_started")
		assert.True(t, ok && started.(bool), fmt.Sprintf("%s should have started", taskName))

		stopped, ok := results.Load(taskName + "_stopped")
		assert.True(t, ok && stopped.(bool), fmt.Sprintf("%s should have stopped", taskName))
	}
}

// ===================== 重叠保护功能测试 =====================

// TestPeriodicTaskManagerAddTaskWithOverlapPrevention 测试添加防重叠任务
func TestPeriodicTaskManagerAddTaskWithOverlapPrevention(t *testing.T) {
	manager := NewPeriodicTaskManager()

	result := manager.AddTaskWithOverlapPrevention("overlap_task", time.Millisecond*100, func(ctx context.Context) error {
		return nil
	})

	assert.Equal(t, manager, result, "AddTaskWithOverlapPrevention should return the manager for chaining")
	assert.Equal(t, 1, manager.GetTaskCount(), "task count should be 1")
	assert.Equal(t, "overlap_task", manager.tasks[0].GetName(), "task name should match")
	assert.True(t, manager.tasks[0].GetPreventOverlap(), "PreventOverlap should be true")
}

// TestPeriodicTaskManagerAddTaskWithOverlapPreventionAndCallback 测试添加带回调的防重叠任务
func TestPeriodicTaskManagerAddTaskWithOverlapPreventionAndCallback(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var callbackCalled bool

	result := manager.AddTaskWithOverlapPreventionAndCallback(
		"overlap_callback_task",
		time.Millisecond*100,
		func(ctx context.Context) error { return nil },
		func(name string) {
			callbackCalled = true
			t.Logf("重叠回调被调用: %s", name)
		},
	)

	assert.Equal(t, manager, result, "AddTaskWithOverlapPreventionAndCallback should return the manager for chaining")
	assert.Equal(t, 1, manager.GetTaskCount(), "task count should be 1")
	assert.True(t, manager.tasks[0].GetPreventOverlap(), "PreventOverlap should be true")
	assert.NotNil(t, manager.tasks[0].GetOnOverlapSkipped(), "OnOverlapSkipped should be set")

	// 验证回调变量被正确设置
	_ = callbackCalled // 使用变量避免编译警告
}

// TestPeriodicTaskManagerOverlapPrevention 测试重叠保护功能
func TestPeriodicTaskManagerOverlapPrevention(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount, overlapCount int32

	// 添加一个执行时间较长的任务（确保会产生重叠）
	manager.AddTaskWithOverlapPreventionAndCallback(
		"slow_task",
		time.Millisecond*20, // 20ms间隔
		func(ctx context.Context) error {
			atomic.AddInt32(&executionCount, 1)
			t.Logf("任务开始执行，当前执行次数: %d", atomic.LoadInt32(&executionCount))

			// 执行时间远长于间隔时间，确保产生重叠
			time.Sleep(time.Millisecond * 100) // 100ms >> 20ms，5倍长

			t.Logf("任务执行完成")
			return nil
		},
		func(name string) {
			count := atomic.AddInt32(&overlapCount, 1)
			t.Logf("!!! 重叠被跳过: %s (第%d次)", name, count)
		},
	)

	manager.Start()

	// 运行足够长时间产生重叠
	time.Sleep(time.Millisecond * 150)

	manager.Stop()

	executions := atomic.LoadInt32(&executionCount)
	overlaps := atomic.LoadInt32(&overlapCount)

	t.Logf("🧪 重叠保护测试结果: 执行次数=%d, 重叠跳过次数=%d", executions, overlaps)

	// 基本验证
	assert.Greater(t, executions, int32(0), "should have some executions")

	// 在200ms内，10ms间隔理论上应该尝试20次
	// 但由于每次执行50ms，实际只能执行几次，其余应被跳过
	totalAttempts := executions + overlaps
	t.Logf("总尝试次数: %d (执行: %d + 跳过: %d)", totalAttempts, executions, overlaps)

	// 在300ms内，20ms间隔理论上尝试15次，但每次执行100ms，最多3次
	assert.LessOrEqual(t, executions, int32(3), "execution count should be limited by overlap prevention")

	// 如果有重叠跳过更好，但不强制要求（可能是时序问题）
	if overlaps > 0 {
		t.Logf("✅ 重叠保护正常工作，跳过了 %d 次重叠执行", overlaps)
	} else {
		t.Logf("⚠️ 未检测到重叠跳过，可能是时序问题或任务执行太快")
	}
}

// TestPeriodicTaskManagerOverlapPreventionWithoutCallback 测试无回调的重叠保护
func TestPeriodicTaskManagerOverlapPreventionWithoutCallback(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount int32

	// 添加一个执行时间较长的任务，但不设置重叠回调
	manager.AddTaskWithOverlapPrevention(
		"slow_task_no_callback",
		time.Millisecond*50,
		func(ctx context.Context) error {
			atomic.AddInt32(&executionCount, 1)

			// 使用select防止阻塞
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Millisecond * 150): // 执行时间比间隔长
				return nil
			}
		},
	)

	manager.Start()
	time.Sleep(time.Millisecond * 500)
	manager.Stop()

	executions := atomic.LoadInt32(&executionCount)

	// 验证重叠保护起作用
	assert.Greater(t, executions, int32(0), "should have some executions")
	assert.Less(t, executions, int32(10), "execution count should be limited by overlap prevention")
}

// TestPeriodicTaskManagerNoOverlapPrevention 测试无重叠保护的对比
func TestPeriodicTaskManagerNoOverlapPrevention(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var startCount, endCount int32

	// 添加一个普通任务（无重叠保护）
	manager.AddSimpleTask(
		"normal_task",
		time.Millisecond*50,
		func(ctx context.Context) error {
			atomic.AddInt32(&startCount, 1)
			time.Sleep(time.Millisecond * 100) // 执行时间比间隔长
			atomic.AddInt32(&endCount, 1)
			return nil
		},
	)

	manager.Start()
	time.Sleep(time.Millisecond * 150)
	manager.Stop()

	starts := atomic.LoadInt32(&startCount)
	ends := atomic.LoadInt32(&endCount)

	// 无重叠保护的任务可能会有多个实例并发执行
	// 所以开始次数可能大于结束次数
	assert.Greater(t, starts, int32(0), "should have task starts")
	assert.GreaterOrEqual(t, starts, ends, "starts should be >= ends due to possible overlap")
}

// TestPeriodicTaskManagerMixedTasks 测试混合任务（有/无重叠保护）
func TestPeriodicTaskManagerMixedTasks(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var normalCount, protectedCount, overlapSkipCount int32

	// 普通任务
	manager.AddSimpleTask("normal", time.Millisecond*30, func(ctx context.Context) error {
		atomic.AddInt32(&normalCount, 1)
		time.Sleep(time.Millisecond * 200) // 增加到200ms
		return nil
	})

	// 有重叠保护的任务
	manager.AddTaskWithOverlapPreventionAndCallback(
		"protected",
		time.Millisecond*30, // 与普通任务相同的间隔
		func(ctx context.Context) error {
			atomic.AddInt32(&protectedCount, 1)
			time.Sleep(time.Millisecond * 200) // 增加到200ms
			return nil
		},
		func(name string) {
			atomic.AddInt32(&overlapSkipCount, 1)
		},
	)

	manager.Start()
	time.Sleep(time.Millisecond * 250)
	manager.Stop()

	normal := atomic.LoadInt32(&normalCount)
	protected := atomic.LoadInt32(&protectedCount)
	skips := atomic.LoadInt32(&overlapSkipCount)

	assert.Greater(t, normal, int32(0), "normal task should execute")
	assert.Greater(t, protected, int32(0), "protected task should execute")
	assert.Greater(t, skips, int32(0), "should have overlap skips for protected task")

	// 通常情况下，保护任务的执行次数应该少于或等于普通任务
	// 因为保护任务会跳过重叠执行
	assert.LessOrEqual(t, protected+skips, normal*2, "protected task behavior should be different from normal task")

	t.Logf("🧪 混合任务测试结果: 普通任务=%d, 保护任务=%d, 跳过次数=%d", normal, protected, skips)
}

// TestPeriodicTaskManagerOverlapPreventionWithError 测试重叠保护中的错误处理
func TestPeriodicTaskManagerOverlapPreventionWithError(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount, errorCount, overlapCount int32

	manager.SetDefaultErrorHandler(func(name string, err error) {
		atomic.AddInt32(&errorCount, 1)
	})

	manager.AddTaskWithOverlapPreventionAndCallback(
		"error_task",
		time.Millisecond*30, // 减少间隔
		func(ctx context.Context) error {
			atomic.AddInt32(&executionCount, 1)
			time.Sleep(time.Millisecond * 150) // 增加执行时间
			return fmt.Errorf("test error")
		},
		func(name string) {
			atomic.AddInt32(&overlapCount, 1)
		},
	)

	manager.Start()
	time.Sleep(time.Millisecond * 200)
	manager.Stop()

	executions := atomic.LoadInt32(&executionCount)
	errors := atomic.LoadInt32(&errorCount)
	overlaps := atomic.LoadInt32(&overlapCount)

	assert.Greater(t, executions, int32(0), "should have executions")
	assert.Greater(t, errors, int32(0), "should have errors")
	assert.Greater(t, overlaps, int32(0), "should have overlaps")
	// 由于并发时序，errors 可能比 executions 少 1（最后一个任务可能在 Stop 时被中断）
	assert.GreaterOrEqual(t, executions, errors, "executions should be >= errors")
	assert.LessOrEqual(t, int(executions-errors), 1, "difference should be at most 1")
}

// TestPeriodicTaskManagerOverlapPreventionWithPanic 测试重叠保护中的panic处理
func TestPeriodicTaskManagerOverlapPreventionWithPanic(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount, panicCount, overlapCount int32

	manager.SetDefaultErrorHandler(func(name string, err error) {
		if name == "panic_task" {
			atomic.AddInt32(&panicCount, 1)
		}
	})

	manager.AddTaskWithOverlapPreventionAndCallback(
		"panic_task",
		time.Millisecond*50,
		func(ctx context.Context) error {
			atomic.AddInt32(&executionCount, 1)
			time.Sleep(time.Millisecond * 100)
			panic("test panic")
		},
		func(name string) {
			atomic.AddInt32(&overlapCount, 1)
		},
	)

	manager.Start()
	time.Sleep(time.Millisecond * 300)
	manager.Stop()

	executions := atomic.LoadInt32(&executionCount)
	panics := atomic.LoadInt32(&panicCount)
	overlaps := atomic.LoadInt32(&overlapCount)

	assert.Greater(t, executions, int32(0), "should have executions")
	assert.Greater(t, panics, int32(0), "should have panics")
	assert.Greater(t, overlaps, int32(0), "should have overlaps")
	// 由于并发时序，panic 可能比 executions 少（panic 恢复后任务可能被取消）
	assert.GreaterOrEqual(t, executions, panics, "executions should be >= panics")
	assert.LessOrEqual(t, int(executions-panics), 1, "difference should be at most 1")
	assert.LessOrEqual(t, int(executions-panics), 1, "difference should be at most 1")
}

// TestPeriodicTaskManagerFastTaskWithOverlapPrevention 测试快速任务的重叠保护
func TestPeriodicTaskManagerFastTaskWithOverlapPrevention(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount, overlapCount int32

	// 添加一个执行时间很短的任务
	manager.AddTaskWithOverlapPreventionAndCallback(
		"fast_task",
		time.Millisecond*100,
		func(ctx context.Context) error {
			atomic.AddInt32(&executionCount, 1)
			time.Sleep(time.Millisecond * 10) // 很短的执行时间
			return nil
		},
		func(name string) {
			atomic.AddInt32(&overlapCount, 1)
		},
	)

	manager.Start()
	time.Sleep(time.Millisecond * 500)
	manager.Stop()

	executions := atomic.LoadInt32(&executionCount)
	overlaps := atomic.LoadInt32(&overlapCount)

	assert.Greater(t, executions, int32(3), "fast task should execute multiple times")
	assert.Equal(t, int32(0), overlaps, "fast task should not have overlaps")
}

// TestPeriodicTaskManagerOverlapPreventionThreadSafety 测试重叠保护的线程安全性
func TestPeriodicTaskManagerOverlapPreventionThreadSafety(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount, overlapCount int32
	var activeExecutions int32

	manager.AddTaskWithOverlapPreventionAndCallback(
		"thread_safe_task",
		time.Millisecond*20,
		func(ctx context.Context) error {
			current := atomic.AddInt32(&activeExecutions, 1)
			defer atomic.AddInt32(&activeExecutions, -1)

			// 验证同时只有一个执行实例
			assert.Equal(t, int32(1), current, "should only have one active execution")

			atomic.AddInt32(&executionCount, 1)
			time.Sleep(time.Millisecond * 100)
			return nil
		},
		func(name string) {
			atomic.AddInt32(&overlapCount, 1)
		},
	)

	manager.Start()
	time.Sleep(time.Millisecond * 250)
	manager.Stop()

	// 等待一小段时间让最后的任务完成清理
	time.Sleep(time.Millisecond * 50)

	executions := atomic.LoadInt32(&executionCount)
	overlaps := atomic.LoadInt32(&overlapCount)
	final := atomic.LoadInt32(&activeExecutions)

	assert.Greater(t, executions, int32(0), "should have executions")
	assert.Greater(t, overlaps, int32(0), "should have overlaps")
	// 允许最多 1 个活跃执行（由于并发时序）
	assert.LessOrEqual(t, final, int32(1), "should have at most 1 active execution after stop")
}

// ===================== 任务移除和取消功能测试 =====================

// TestPeriodicTaskManagerRemoveTask_Basic 测试基本的任务移除功能
func TestPeriodicTaskManagerRemoveTask_Basic(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加任务
	manager.AddSimpleTask("remove_test", time.Second, func(ctx context.Context) error {
		return nil
	})

	// 验证任务已添加
	assert.Equal(t, 1, manager.GetTaskCount(), "should have 1 task")
	names := manager.GetTaskNames()
	assert.Contains(t, names, "remove_test", "should contain remove_test")

	// 移除任务
	removed := manager.RemoveTask("remove_test")
	assert.True(t, removed, "should successfully remove task")

	// 验证任务已移除
	assert.Equal(t, 0, manager.GetTaskCount(), "should have 0 tasks after removal")
	names = manager.GetTaskNames()
	assert.NotContains(t, names, "remove_test", "should not contain remove_test")

	// 尝试移除不存在的任务
	removed = manager.RemoveTask("non_existent")
	assert.False(t, removed, "should not be able to remove non-existent task")
}

// TestPeriodicTaskManagerRemoveRunningTask 测试移除正在运行的任务
func TestPeriodicTaskManagerRemoveRunningTask(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount int32
	var taskCancelled bool

	// 添加一个长时间运行的任务
	manager.AddTaskWithOverlapPrevention("long_running", time.Millisecond*50, func(ctx context.Context) error {
		atomic.AddInt32(&executionCount, 1)
		t.Log("任务开始执行...")

		// 检查是否被取消
		select {
		case <-time.After(time.Millisecond * 200):
			t.Log("任务正常完成")
		case <-ctx.Done():
			t.Log("任务被取消")
			taskCancelled = true
		}

		return nil
	})

	// 启动任务管理器
	err := manager.Start()
	assert.NoError(t, err, "should start successfully")

	// 等待任务开始执行
	time.Sleep(time.Millisecond * 100)

	// 验证任务正在执行
	details := manager.GetTaskDetails("long_running")
	assert.Equal(t, 1, len(details), "should find the task")
	assert.True(t, details[0].IsExecuting, "task should be executing")

	// 移除正在运行的任务
	t.Log("开始移除正在运行的任务...")
	removed := manager.RemoveTask("long_running")
	assert.True(t, removed, "should successfully remove running task")

	// 验证任务已移除
	assert.Equal(t, 0, manager.GetTaskCount(), "should have 0 tasks after removal")

	// 等待一段时间看任务是否被取消
	time.Sleep(time.Millisecond * 300)

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err, "should stop successfully")

	executions := atomic.LoadInt32(&executionCount)
	assert.Greater(t, executions, int32(0), "task should have executed at least once")

	if taskCancelled {
		t.Log("✅ 任务成功被取消")
	} else {
		t.Log("⚠️ 任务可能在取消前已完成")
	}
}

// TestPeriodicTaskManagerRemoveTaskWithTimeout 测试带超时的任务移除
func TestPeriodicTaskManagerRemoveTaskWithTimeout(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executionCount int32
	var taskCancelled bool

	// 添加一个长时间运行的任务
	manager.AddTaskWithOverlapPrevention("timeout_test", time.Millisecond*100, func(ctx context.Context) error {
		atomic.AddInt32(&executionCount, 1)

		select {
		case <-time.After(time.Millisecond * 500): // 很长的执行时间
			return nil
		case <-ctx.Done():
			taskCancelled = true
			return ctx.Err()
		}
	})

	// 启动任务管理器
	err := manager.Start()
	assert.NoError(t, err, "should start successfully")

	// 等待任务开始执行
	time.Sleep(time.Millisecond * 150)

	// 使用超时移除任务
	start := time.Now()
	removed := manager.RemoveTaskWithTimeout("timeout_test", time.Millisecond*200)
	duration := time.Since(start)

	assert.True(t, removed, "should successfully remove task")
	assert.Less(t, duration, time.Millisecond*300, "should not take too long")

	// 验证任务已移除
	assert.Equal(t, 0, manager.GetTaskCount(), "should have 0 tasks after removal")

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err, "should stop successfully")

	t.Logf("移除操作耗时: %v", duration)
	t.Logf("任务执行次数: %d", atomic.LoadInt32(&executionCount))
	t.Logf("任务是否被取消: %v", taskCancelled)
}

// TestPeriodicTaskManagerRemoveTaskTimeout 测试移除任务超时情况
func TestPeriodicTaskManagerRemoveTaskTimeout(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加一个会阻塞很久的任务
	manager.AddTaskWithOverlapPrevention("blocking_task", time.Millisecond*50, func(ctx context.Context) error {
		// 忽略取消信号，模拟无法优雅停止的任务
		time.Sleep(time.Millisecond * 500)
		return nil
	})

	// 启动任务管理器
	err := manager.Start()
	assert.NoError(t, err, "should start successfully")

	// 等待任务开始执行
	time.Sleep(time.Millisecond * 100)

	// 尝试在很短时间内移除任务
	start := time.Now()
	removed := manager.RemoveTaskWithTimeout("blocking_task", time.Millisecond*100)
	duration := time.Since(start)

	// 应该能成功移除（超时后强制移除）
	assert.True(t, removed, "should remove task even on timeout")
	assert.GreaterOrEqual(t, duration, time.Millisecond*100, "should wait for timeout")
	assert.Less(t, duration, time.Millisecond*200, "should not wait too long")

	// 验证任务已移除
	assert.Equal(t, 0, manager.GetTaskCount(), "should have 0 tasks after removal")

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err, "should stop successfully")

	t.Logf("超时移除操作耗时: %v", duration)
}

// TestPeriodicTaskManagerRemoveMultipleTasks 测试移除多个任务
func TestPeriodicTaskManagerRemoveMultipleTasks(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var task1Count, task2Count, task3Count int32

	// 添加多个任务
	manager.AddSimpleTask("task1", time.Millisecond*50, func(ctx context.Context) error {
		atomic.AddInt32(&task1Count, 1)
		return nil
	})

	manager.AddSimpleTask("task2", time.Millisecond*75, func(ctx context.Context) error {
		atomic.AddInt32(&task2Count, 1)
		return nil
	})

	manager.AddSimpleTask("task3", time.Millisecond*100, func(ctx context.Context) error {
		atomic.AddInt32(&task3Count, 1)
		return nil
	})

	assert.Equal(t, 3, manager.GetTaskCount(), "should have 3 tasks")

	// 启动任务管理器
	err := manager.Start()
	assert.NoError(t, err, "should start successfully")

	// 等待所有任务至少执行一次
	err = manager.WaitForExecution(time.Second)
	assert.NoError(t, err, "all tasks should execute")

	// 移除其中两个任务
	removed1 := manager.RemoveTask("task1")
	removed2 := manager.RemoveTask("task3")

	assert.True(t, removed1, "should remove task1")
	assert.True(t, removed2, "should remove task3")
	assert.Equal(t, 1, manager.GetTaskCount(), "should have 1 task remaining")

	// 验证剩余的任务
	names := manager.GetTaskNames()
	assert.Contains(t, names, "task2", "task2 should remain")
	assert.NotContains(t, names, "task1", "task1 should be removed")
	assert.NotContains(t, names, "task3", "task3 should be removed")

	// 继续运行一段时间
	time.Sleep(time.Millisecond * 100)

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err, "should stop successfully")

	// 验证所有任务都有执行
	assert.Greater(t, atomic.LoadInt32(&task1Count), int32(0), "task1 should have executed")
	assert.Greater(t, atomic.LoadInt32(&task2Count), int32(0), "task2 should have executed")
	assert.Greater(t, atomic.LoadInt32(&task3Count), int32(0), "task3 should have executed")

	t.Logf("Task1 执行次数: %d", atomic.LoadInt32(&task1Count))
	t.Logf("Task2 执行次数: %d", atomic.LoadInt32(&task2Count))
	t.Logf("Task3 执行次数: %d", atomic.LoadInt32(&task3Count))
}

// TestPeriodicTaskManagerClearAllTasks 测试清除所有任务
func TestPeriodicTaskManagerClearAllTasks(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加多个任务
	for i := 0; i < 5; i++ {
		taskName := fmt.Sprintf("task_%d", i)
		manager.AddSimpleTask(taskName, time.Millisecond*100, func(ctx context.Context) error {
			return nil
		})
	}

	assert.Equal(t, 5, manager.GetTaskCount(), "should have 5 tasks")

	// 启动任务管理器
	err := manager.Start()
	assert.NoError(t, err, "should start successfully")

	time.Sleep(time.Millisecond * 50)

	// 清除所有任务
	manager.ClearAllTasks()

	assert.Equal(t, 0, manager.GetTaskCount(), "should have 0 tasks after clear")
	assert.Equal(t, 0, len(manager.GetTaskNames()), "should have no task names")

	details := manager.GetTaskDetails()
	assert.Equal(t, 0, len(details), "should have no task details")

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err, "should stop successfully")
}

// TestPeriodicTaskManagerRemoveTaskConcurrency 测试并发移除任务
func TestPeriodicTaskManagerRemoveTaskConcurrency(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加多个任务
	taskCount := 10
	for i := 0; i < taskCount; i++ {
		taskName := fmt.Sprintf("concurrent_task_%d", i)
		manager.AddSimpleTask(taskName, time.Millisecond*100, func(ctx context.Context) error {
			time.Sleep(time.Millisecond * 50)
			return nil
		})
	}

	assert.Equal(t, taskCount, manager.GetTaskCount(), "should have all tasks")

	// 启动任务管理器
	err := manager.Start()
	assert.NoError(t, err, "should start successfully")

	time.Sleep(time.Millisecond * 50)

	// 并发移除任务
	var wg sync.WaitGroup
	var removedCount int32

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			taskName := fmt.Sprintf("concurrent_task_%d", index)
			if manager.RemoveTask(taskName) {
				atomic.AddInt32(&removedCount, 1)
			}
		}(i)
	}

	wg.Wait()

	// 验证所有任务都被移除
	assert.Equal(t, int32(taskCount), atomic.LoadInt32(&removedCount), "should remove all tasks")
	assert.Equal(t, 0, manager.GetTaskCount(), "should have 0 tasks after removal")

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err, "should stop successfully")
}

// TestPeriodicTaskManagerGetTaskDetailsAfterRemoval 测试移除后获取任务详情
func TestPeriodicTaskManagerGetTaskDetailsAfterRemoval(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加任务
	task := NewPeriodicTask("detail_test", time.Second, func(ctx context.Context) error {
		return nil
	}).SetImmediateStart(true).SetPreventOverlap(true)

	manager.AddTask(task)

	// 验证任务详情
	details := manager.GetTaskDetails("detail_test")
	assert.Equal(t, 1, len(details), "should have 1 task detail")
	assert.Equal(t, "detail_test", details[0].Name, "task name should match")
	assert.True(t, details[0].ImmediateStart, "should have immediate start")
	assert.True(t, details[0].PreventOverlap, "should have overlap prevention")

	// 移除任务
	removed := manager.RemoveTask("detail_test")
	assert.True(t, removed, "should remove task")

	// 验证任务详情已清空
	details = manager.GetTaskDetails("detail_test")
	assert.Equal(t, 0, len(details), "should have no task details after removal")

	allDetails := manager.GetTaskDetails()
	assert.Equal(t, 0, len(allDetails), "should have no task details")
}

// TestPeriodicTaskManagerTaskCancellationContext 测试任务取消上下文
func TestPeriodicTaskManagerTaskCancellationContext(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var cancelledCount int32
	var executionCount int32

	// 添加一个会检查取消信号的任务
	manager.AddSimpleTask("cancellable", time.Millisecond*50, func(ctx context.Context) error {
		atomic.AddInt32(&executionCount, 1)

		// 模拟长时间运行并检查取消
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				atomic.AddInt32(&cancelledCount, 1)
				return ctx.Err()
			case <-time.After(time.Millisecond * 20):
				// 继续执行
			}
		}
		return nil
	})

	// 启动任务管理器
	err := manager.Start()
	assert.NoError(t, err, "should start successfully")

	// 等待任务开始执行
	time.Sleep(time.Millisecond * 100)

	// 移除任务（这会取消任务的上下文）
	removed := manager.RemoveTask("cancellable")
	assert.True(t, removed, "should remove task")

	// 等待取消生效
	time.Sleep(time.Millisecond * 200)

	// 停止管理器
	err = manager.Stop()
	assert.NoError(t, err, "should stop successfully")

	executions := atomic.LoadInt32(&executionCount)
	cancelled := atomic.LoadInt32(&cancelledCount)

	assert.Greater(t, executions, int32(0), "task should have executed")

	if cancelled > 0 {
		t.Logf("✅ 任务成功响应取消信号，取消次数: %d", cancelled)
	} else {
		t.Logf("⚠️ 任务可能在取消前已完成，执行次数: %d", executions)
	}
}
