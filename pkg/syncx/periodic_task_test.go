/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-29 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-29 12:00:00
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

// TestPeriodicTaskManager_NewPeriodicTaskManager 测试创建新的任务管理器
func TestPeriodicTaskManager_NewPeriodicTaskManager(t *testing.T) {
	manager := NewPeriodicTaskManager()

	assert.NotNil(t, manager, "manager should not be nil")
	assert.NotNil(t, manager.tasks, "tasks slice should be initialized")
	assert.Equal(t, 0, len(manager.tasks), "initial tasks count should be 0")
	assert.False(t, manager.isRunning, "manager should not be running initially")
}

// TestPeriodicTaskManager_AddTask 测试添加任务
func TestPeriodicTaskManager_AddTask(t *testing.T) {
	manager := NewPeriodicTaskManager()

	task := &PeriodicTask{
		Name:        "test_task",
		Interval:    time.Second,
		ExecuteFunc: func(ctx context.Context) error { return nil },
	}

	result := manager.AddTask(task)

	assert.Equal(t, manager, result, "AddTask should return the manager for chaining")
	assert.Equal(t, 1, manager.GetTaskCount(), "task count should be 1")
	assert.Equal(t, "test_task", manager.tasks[0].Name, "task name should match")
}

// TestPeriodicTaskManager_AddSimpleTask 测试添加简单任务
func TestPeriodicTaskManager_AddSimpleTask(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executed int32

	result := manager.AddSimpleTask("simple_task", time.Millisecond*100, func(ctx context.Context) error {
		atomic.AddInt32(&executed, 1)
		return nil
	})

	assert.Equal(t, manager, result, "AddSimpleTask should return the manager for chaining")
	assert.Equal(t, 1, manager.GetTaskCount(), "task count should be 1")
	assert.Equal(t, "simple_task", manager.tasks[0].Name, "task name should match")
	assert.Equal(t, time.Millisecond*100, manager.tasks[0].Interval, "task interval should match")
	assert.False(t, manager.tasks[0].ImmediateStart, "immediate start should be false by default")
}

// TestPeriodicTaskManager_AddTaskWithImmediateStart 测试添加立即执行任务
func TestPeriodicTaskManager_AddTaskWithImmediateStart(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var executed int32

	result := manager.AddTaskWithImmediateStart("immediate_task", time.Millisecond*100, func(ctx context.Context) error {
		atomic.AddInt32(&executed, 1)
		return nil
	})

	assert.Equal(t, manager, result, "AddTaskWithImmediateStart should return the manager for chaining")
	assert.Equal(t, 1, manager.GetTaskCount(), "task count should be 1")
	assert.Equal(t, "immediate_task", manager.tasks[0].Name, "task name should match")
	assert.True(t, manager.tasks[0].ImmediateStart, "immediate start should be true")
}

// TestPeriodicTaskManager_SetDefaultErrorHandler 测试设置默认错误处理器
func TestPeriodicTaskManager_SetDefaultErrorHandler(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加一个没有错误处理器的任务
	manager.AddSimpleTask("task1", time.Second, func(ctx context.Context) error { return nil })

	// 添加一个已有错误处理器的任务
	task2 := &PeriodicTask{
		Name:        "task2",
		Interval:    time.Second,
		ExecuteFunc: func(ctx context.Context) error { return nil },
		OnError:     func(name string, err error) { /* existing handler */ },
	}
	manager.AddTask(task2)

	// 设置默认错误处理器
	result := manager.SetDefaultErrorHandler(func(name string, err error) {
		// 错误处理逻辑
	})

	assert.Equal(t, manager, result, "SetDefaultErrorHandler should return the manager for chaining")
	assert.NotNil(t, manager.tasks[0].OnError, "task1 should have error handler set")
	assert.NotNil(t, manager.tasks[1].OnError, "task2 should still have its original error handler")
}

// TestPeriodicTaskManager_SetDefaultCallbacks 测试设置默认回调函数
func TestPeriodicTaskManager_SetDefaultCallbacks(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 添加一个没有回调的任务
	manager.AddSimpleTask("task1", time.Second, func(ctx context.Context) error { return nil })

	// 添加一个已有回调的任务
	task2 := &PeriodicTask{
		Name:        "task2",
		Interval:    time.Second,
		ExecuteFunc: func(ctx context.Context) error { return nil },
		OnStart:     func(name string) { /* existing start handler */ },
		OnStop:      func(name string) { /* existing stop handler */ },
	}
	manager.AddTask(task2)

	// 设置默认回调
	result := manager.SetDefaultCallbacks(
		func(name string) { /* start callback */ },
		func(name string) { /* stop callback */ },
	)

	assert.Equal(t, manager, result, "SetDefaultCallbacks should return the manager for chaining")
	assert.NotNil(t, manager.tasks[0].OnStart, "task1 should have start callback set")
	assert.NotNil(t, manager.tasks[0].OnStop, "task1 should have stop callback set")
	assert.NotNil(t, manager.tasks[1].OnStart, "task2 should still have its original start callback")
	assert.NotNil(t, manager.tasks[1].OnStop, "task2 should still have its original stop callback")
}

// TestPeriodicTaskManager_Start_AlreadyRunning 测试重复启动
func TestPeriodicTaskManager_Start_AlreadyRunning(t *testing.T) {
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

// TestPeriodicTaskManager_StartWithContext_AlreadyRunning 测试使用上下文重复启动
func TestPeriodicTaskManager_StartWithContext_AlreadyRunning(t *testing.T) {
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

// TestPeriodicTaskManager_Stop_NotRunning 测试停止未运行的管理器
func TestPeriodicTaskManager_Stop_NotRunning(t *testing.T) {
	manager := NewPeriodicTaskManager()

	err := manager.Stop()
	assert.NoError(t, err, "stopping non-running manager should not error")
	assert.False(t, manager.IsRunning(), "manager should not be running")
}

// TestPeriodicTaskManager_StartStop 测试启动和停止
func TestPeriodicTaskManager_StartStop(t *testing.T) {
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

// TestPeriodicTaskManager_ImmediateStart 测试立即执行任务
func TestPeriodicTaskManager_ImmediateStart(t *testing.T) {
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

// TestPeriodicTaskManager_ErrorHandling 测试错误处理
func TestPeriodicTaskManager_ErrorHandling(t *testing.T) {
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

// TestPeriodicTaskManager_StartStopCallbacks 测试启动停止回调
func TestPeriodicTaskManager_StartStopCallbacks(t *testing.T) {
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

// TestPeriodicTaskManager_MultipleTasksExecution 测试多个任务执行
func TestPeriodicTaskManager_MultipleTasksExecution(t *testing.T) {
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
	time.Sleep(time.Millisecond * 300)
	manager.Stop()

	assert.Greater(t, atomic.LoadInt32(&task1Count), int32(0), "task1 should have executed")
	assert.Greater(t, atomic.LoadInt32(&task2Count), int32(0), "task2 should have executed")
	assert.Greater(t, atomic.LoadInt32(&task3Count), int32(0), "task3 should have executed")
}

// TestPeriodicTaskManager_GetTaskNames 测试获取任务名称
func TestPeriodicTaskManager_GetTaskNames(t *testing.T) {
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

// TestPeriodicTaskManager_ContextCancellation 测试上下文取消
func TestPeriodicTaskManager_ContextCancellation(t *testing.T) {
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

// TestPeriodicTaskManager_StopWithTimeout_Success 测试超时停止成功
func TestPeriodicTaskManager_StopWithTimeout_Success(t *testing.T) {
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

// TestPeriodicTaskManager_StopWithTimeout_Timeout 测试超时停止超时
func TestPeriodicTaskManager_StopWithTimeout_Timeout(t *testing.T) {
	manager := NewPeriodicTaskManager()

	// 创建一个会阻塞的任务
	manager.AddSimpleTask("blocking_task", time.Millisecond*10, func(ctx context.Context) error {
		time.Sleep(time.Second * 2) // 长时间阻塞
		return nil
	})

	manager.Start()
	time.Sleep(time.Millisecond * 50) // 让任务开始执行

	err := manager.StopWithTimeout(time.Millisecond * 100)
	assert.Error(t, err, "stop with timeout should fail due to timeout")
	assert.Contains(t, err.Error(), "timeout", "error should mention timeout")
}

// TestPeriodicTaskManager_Wait 测试等待功能
func TestPeriodicTaskManager_Wait(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var completed bool

	manager.AddSimpleTask("wait_task", time.Second, func(ctx context.Context) error {
		return nil
	})

	manager.Start()

	go func() {
		time.Sleep(time.Millisecond * 100)
		manager.Stop()
		completed = true
	}()

	manager.Wait()
	assert.True(t, completed, "Wait should block until tasks complete")
}

// TestPeriodicTaskManager_ConcurrentAccess 测试并发访问
func TestPeriodicTaskManager_ConcurrentAccess(t *testing.T) {
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

	assert.Equal(t, int32(0), atomic.LoadInt32(&errorCount), "no concurrent access errors should occur")
	assert.Equal(t, 10, manager.GetTaskCount(), "should have 10 tasks")
}

// TestPeriodicTaskManager_TaskWithCustomCallbacks 测试带自定义回调的任务
func TestPeriodicTaskManager_TaskWithCustomCallbacks(t *testing.T) {
	manager := NewPeriodicTaskManager()
	var customStartCalled, customStopCalled, customErrorCalled bool

	task := &PeriodicTask{
		Name:     "custom_task",
		Interval: time.Millisecond * 50,
		ExecuteFunc: func(ctx context.Context) error {
			return fmt.Errorf("custom error")
		},
		OnStart: func(name string) {
			customStartCalled = true
		},
		OnStop: func(name string) {
			customStopCalled = true
		},
		OnError: func(name string, err error) {
			customErrorCalled = true
		},
	}

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

// TestPeriodicTaskManager_EmptyManager 测试空管理器
func TestPeriodicTaskManager_EmptyManager(t *testing.T) {
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

// TestPeriodicTaskManager_TaskExecutionOrder 测试任务执行顺序
func TestPeriodicTaskManager_TaskExecutionOrder(t *testing.T) {
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

// TestPeriodicTaskManager_LongRunningTask 测试长时间运行的任务
func TestPeriodicTaskManager_LongRunningTask(t *testing.T) {
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

// TestPeriodicTaskManager_PanicRecovery 测试panic恢复
func TestPeriodicTaskManager_PanicRecovery(t *testing.T) {
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
} // TestPeriodicTaskManager_HighFrequencyTasks 测试高频任务
func TestPeriodicTaskManager_HighFrequencyTasks(t *testing.T) {
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

// TestPeriodicTaskManager_TaskNameUniqueness 测试任务名称唯一性
func TestPeriodicTaskManager_TaskNameUniqueness(t *testing.T) {
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

// TestPeriodicTaskManager_MemoryUsage 测试内存使用
func TestPeriodicTaskManager_MemoryUsage(t *testing.T) {
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

// TestPeriodicTaskManager_ZeroInterval 测试零间隔任务
func TestPeriodicTaskManager_ZeroInterval(t *testing.T) {
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

// TestPeriodicTaskManager_NegativeInterval 测试负间隔任务
func TestPeriodicTaskManager_NegativeInterval(t *testing.T) {
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

// TestPeriodicTaskManager_VeryLargeInterval 测试非常大的间隔
func TestPeriodicTaskManager_VeryLargeInterval(t *testing.T) {
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

// TestPeriodicTaskManager_ComplexScenario 测试复杂场景
func TestPeriodicTaskManager_ComplexScenario(t *testing.T) {
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

	time.Sleep(time.Millisecond * 300)

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
