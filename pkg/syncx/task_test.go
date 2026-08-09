/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-05 16:31:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-11 11:15:15
 * @FilePath: \go-toolbox\pkg\syncx\task_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 任务管理器创建与基本功能

// 测试创建任务管理器
func TestNewTaskManager(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)
	assert.NotNil(t, tm, "确保任务管理器不为 nil") // 确保任务管理器不为 nil
}

// 测试添加任务
func TestAddTask(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")

	tm.AddTask(task) // 添加任务

	// 验证任务是否成功添加
	assert.Equal(t, 1, len(tm.GetTasks()), "任务数量应为 1")
	assert.NotNil(t, tm.GetTasks()["task1"], "任务 task1 应存在")
}

// 任务执行

// 测试运行任务
func TestRunTask(t *testing.T) {
	tm := NewTaskManager[string, string, string](10)
	index := 10
	for i := 0; i < 50; i++ {
		taskName := fmt.Sprintf("task%d", i)
		task := NewTask[string, string, string](taskName, func(ctx context.Context, input string) (string, error) {
			return input, nil
		}, taskName).SetTimeout(10).SetMaxHistorySize(index).SetMaxHistorySize(index).SetMaxRetries(int32(index))

		for j := 0; j < 50; j++ {
			dependName := fmt.Sprintf("depend%d", j) // 使用不同的名称
			depend := NewTask[string, string, string](dependName, func(ctx context.Context, input string) (string, error) {
				return input, nil
			}, dependName) // 使用 dependName 作为输入
			task.AddDependency(depend)
		}

		// 设置依赖执行模式为并发
		task.SetDependExecutionMode(Concurrent)
		tm.AddTask(task)
	}

	err := tm.Run()                    // 执行所有任务
	assert.NoError(t, err, "任务执行应无错误") // 添加错误检查

	for _, task := range tm.GetTasks() {
		assert.NoError(t, task.GetError(), "没有错误")
		assert.Equal(t, task.GetCallbackDuration(), time.Duration(0))
		assert.Equal(t, task.GetTaskType(), MainTask)
		assert.Equal(t, Completed, task.GetState(), "状态应为已完成")
		assert.Equal(t, task.GetInput(), task.GetName(), "任务结果应为 'task XXX'")
		assert.Equal(t, task.GetMaxRetries(), int32(index))
		assert.Equal(t, task.GetCallbackResult(), "")
	}
}

// 测试任务执行错误
func TestTaskWithError(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return "", errors.New("任务执行失败") // 模拟错误
	}, "hello").SetFailureCallback(func(result string, err error) (string, error) {
		return "callback failure", nil
	})

	tm.AddTask(task) // 添加任务
	tm.Run()         // 执行任务

	// 验证任务状态和错误
	assert.Equal(t, Failed, task.GetState(), "任务状态应为失败")
	assert.NotNil(t, task.GetError(), "任务应返回错误")
	assert.Equal(t, "任务执行失败", task.GetError().Error(), "错误信息应为 '任务执行失败'")
	assert.Equal(t, task.GetCallbackResult(), "callback failure")
}

// 测试任务成功后的回调
func TestTaskWithCallback(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	var callbackResult string
	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello").SetSuccessCallback(func(result string, err error) (string, error) {
		callbackResult = result // 保存回调结果
		return "callback success", nil
	}).SetFailureCallback(func(result string, err error) (string, error) {
		return "callback failure", nil
	})

	tm.AddTask(task) // 添加任务
	tm.Run()         // 执行任务

	// 验证回调结果
	assert.Equal(t, "hello", callbackResult, "回调结果应为 'hello'")
	assert.Equal(t, task.GetCallbackResult(), "callback success")
}

// 测试任务取消
func TestTaskCancellation(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")

	tm.AddTask(task) // 添加任务
	tm.Cancel("task1")
	tm.Run() // 执行任务

	// 验证任务状态
	assert.Equal(t, Cancelled, task.GetState(), "任务状态应为已取消")
}

// 测试任务管理器的取消所有任务
func TestCancelAllTasks(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	for i := 0; i < 5; i++ {
		task := NewTask[string, string, string](fmt.Sprintf("task%d", i), func(ctx context.Context, input string) (string, error) {
			return input, nil
		}, fmt.Sprintf("hello %d", i))
		tm.AddTask(task)
	}
	tm.CancelAll() // 取消所有任务
	tm.Run()       // 执行所有任务

	for _, task := range tm.GetTasks() {
		assert.Equal(t, Cancelled, task.GetState(), "所有任务状态应为已取消")
	}
}

// 任务依赖

// TestCircularDependency 测试循环依赖和自我依赖
func TestCircularDependency(t *testing.T) {
	// 创建任务管理器
	tm := NewTaskManager[string, string, string](2)

	// 创建任务 A
	taskA := NewTask[string, string, string]("taskA", func(ctx context.Context, input string) (string, error) {
		return "taskA executed", nil
	}, "inputA")

	// 创建任务 B
	taskB := NewTask[string, string, string]("taskB", func(ctx context.Context, input string) (string, error) {
		return "taskB executed", nil
	}, "inputB")

	// 1. 测试循环依赖
	// 添加 B 作为 A 的依赖
	taskA.AddDependency(taskB)

	// 尝试将 A 添加为 B 的依赖，应该会引发循环依赖错误
	assert.Panics(t, func() {
		taskB.AddDependency(taskA)
	}, "应该引发循环依赖错误")

	// 2. 测试自我依赖
	assert.Panics(t, func() {
		taskA.AddDependency(taskA) // 尝试将自己添加为依赖
	}, "应该引发自我依赖错误")

	tm.AddTask(taskA)
	tm.AddTask(taskB)
	assert.Equal(t, 2, len(tm.GetTasks()), fmt.Sprintf("任务数量应为 %d", 2))
}

// 测试依赖任务成功时主任务的执行
func TestTaskDependenciesSuccess(t *testing.T) {
	// 创建任务管理器
	tm := NewTaskManager[string, int, string](2)

	// 创建依赖任务1，成功返回结果
	depTask1 := NewTask[string, int, string]("dep1", func(ctx context.Context, input string) (int, error) {
		return 1, nil // 成功
	}, "input1")

	// 创建依赖任务2，成功返回结果
	depTask2 := NewTask[string, int, string]("dep2", func(ctx context.Context, input string) (int, error) {
		return 2, nil // 成功
	}, "input2").SetPriority(999)

	// 创建主任务，依赖于depTask1和depTask2
	mainTask := NewTask[string, int, string]("main", func(ctx context.Context, input string) (int, error) {
		return 3, nil // 成功
	}, "inputMain")

	// 添加依赖
	mainTask.AddDependency(depTask1)
	mainTask.AddDependency(depTask2)

	// 添加任务到任务管理器
	tm.AddTask(depTask1)
	tm.AddTask(depTask2)
	tm.AddTask(mainTask)

	// 运行任务管理器
	tm.Run()

	// 验证主任务的状态
	assert.Equal(t, Completed, mainTask.GetState(), "主任务状态应为已完成")
	assert.Equal(t, 3, mainTask.GetResult(), "主任务结果应为 3")
}

// 测试依赖任务失败时主任务的行为
func TestTaskDependenciesFailure(t *testing.T) {
	// 创建任务管理器
	tm := NewTaskManager[string, int, string](2)

	// 创建依赖任务1，成功返回结果
	depTask1 := NewTask[string, int, string]("dep1", func(ctx context.Context, input string) (int, error) {
		return 1, nil // 成功
	}, "input1")

	// 创建依赖任务2，失败返回错误
	depTask2 := NewTask[string, int, string]("dep2", func(ctx context.Context, input string) (int, error) {
		return 0, errors.New("test err") // 失败
	}, "input2")

	// 创建主任务，依赖于depTask1和depTask2
	mainTask := NewTask[string, int, string]("main", func(ctx context.Context, input string) (int, error) {
		return 3, nil // 成功
	}, "inputMain")

	// 添加依赖
	mainTask.AddDependency(depTask1)
	mainTask.AddDependency(depTask2)

	// 添加任务到任务管理器
	tm.AddTask(mainTask)

	// 运行任务管理器
	tm.Run()

	// 验证主任务的状态
	assert.Equal(t, Failed, mainTask.GetState(), "主任务状态应为失败")
	assert.Error(t, mainTask.GetError(), "主任务应返回错误")
	assert.Equal(t, "dependency 'dep2' failed: test err", mainTask.GetError().Error(), "错误信息应为 'dependency 'dep2' failed: test err'")
}

// 测试添加多个依赖
func TestAddMultipleDependencies(t *testing.T) {
	tm := NewTaskManager[string, int, string](2)

	depTask1 := NewTask[string, int, string]("dep1", func(ctx context.Context, input string) (int, error) {
		return 1, nil
	}, "input1")

	depTask2 := NewTask[string, int, string]("dep2", func(ctx context.Context, input string) (int, error) {
		return 2, nil
	}, "input2")

	mainTask := NewTask[string, int, string]("main", func(ctx context.Context, input string) (int, error) {
		return 3, nil
	}, "inputMain")

	mainTask.AddDependency(depTask1)
	mainTask.AddDependency(depTask2)

	tm.AddTask(depTask1)
	tm.AddTask(depTask2)
	tm.AddTask(mainTask)

	tm.Run()

	assert.Equal(t, Completed, mainTask.GetState(), "主任务状态应为已完成")
	assert.Equal(t, Completed, depTask1.GetState(), "依赖任务1状态应为已完成")
	assert.Equal(t, Completed, depTask2.GetState(), "依赖任务2状态应为已完成")
}

// 测试取消主任务及其依赖
func TestCancelMainTaskWithDependencies(t *testing.T) {
	tm := NewTaskManager[string, int, string](2)

	depTask := NewTask[string, int, string]("dep", func(ctx context.Context, input string) (int, error) {
		return 1, nil
	}, "input")

	mainTask := NewTask[string, int, string]("main", func(ctx context.Context, input string) (int, error) {
		return 2, nil
	}, "inputMain")

	mainTask.AddDependency(depTask)

	tm.AddTask(depTask)
	tm.AddTask(mainTask)

	tm.Cancel("main") // 取消主任务
	tm.Run()

	assert.Equal(t, Cancelled, mainTask.GetState(), "主任务状态应为已取消")
	assert.Equal(t, Cancelled, depTask.GetState(), "依赖任务应为已取消")
}

// 任务历史记录

// 测试任务历史记录
func TestTaskHistory(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")

	tm.AddTask(task) // 添加任务
	tm.Run()         // 执行任务

	history := task.GetTaskHistory("task1")
	assert.NotNil(t, history, "任务历史记录应存在")
	assert.Equal(t, 1, len(history), "任务历史记录应包含一条记录")
	assert.Equal(t, Completed, history[0].GetState(), "历史记录状态应为已完成")
	assert.Equal(t, "hello", history[0].GetResult(), "历史记录结果应为 'hello'")
}

// 测试任务多次执行
func TestTaskMultipleExecutions(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	// 定义任务函数
	taskFunc := func(ctx context.Context, input string) (string, error) {
		return input, nil
	}

	// 第一次执行
	task1 := NewTask[string, string, string]("task1", taskFunc, "hello")
	tm.AddTask(task1) // 添加任务
	tm.Run()          // 执行任务
	assert.Equal(t, Completed, task1.GetState(), "第一次执行后任务状态应为已完成")

	// 第二次执行
	task2 := NewTask[string, string, string]("task1", taskFunc, "hello")
	tm.AddTask(task2) // 添加新的任务
	tm.Run()          // 执行任务
	assert.Equal(t, Completed, task2.GetState(), "第二次执行后任务状态应为已完成")
}

// 测试任务回调失败
func TestTaskCallbackFailure(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello").SetSuccessCallback(func(result string, err error) (string, error) {
		return "", errors.New("callback failed") // 模拟回调失败
	})

	tm.AddTask(task)
	tm.Run()

	assert.Equal(t, Failed, task.GetState(), "任务状态应为失败")
	assert.Equal(t, Failed, task.GetCallbackState(), "回调状态应为失败")
	assert.Nil(t, task.GetError(), "主任务没有错误")
	assert.NotNil(t, task.GetCallbackError(), "回调任务应返回错误")
	assert.Equal(t, "callback failed", task.GetCallbackError().Error(), "错误信息应为 'callback failed'")
}

// 测试任务管理器的并发执行
func TestTaskManagerConcurrency(t *testing.T) {
	const numTasks = 1000
	const concurrency = 300
	const resultX = 2

	// 创建 TaskManager 实例
	tm := NewTaskManager[int, int, string](concurrency)

	// 添加任务
	for i := 1; i <= numTasks; i++ {
		task := NewTask[int, int, string](fmt.Sprintf("task-%d", i), func(ctx context.Context, input int) (int, error) {
			return input, nil
		}, i*resultX)
		tm.AddTask(task)
	}

	// 开始性能测试
	startTime := time.Now()
	tm.Run() // 执行所有任务
	duration := time.Since(startTime)

	// 验证任务执行时间
	assert.Less(t, duration, 2*time.Second, "Tasks should complete in less than 2 seconds")

	// 验证任务结果
	for i := 1; i <= numTasks; i++ {
		taskName := fmt.Sprintf("task-%d", i)
		task := tm.GetTasks()[taskName]
		assert.Equal(t, task.GetResult(), i*resultX, "Task result should be double the input")
		assert.Equal(t, task.GetState(), Completed, "Task should be completed")
	}
}

// 测试任务管理器的最大并发限制
func TestTaskManagerMaxConcurrency(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	for i := 0; i < 5; i++ {
		task := NewTask[string, string, string](fmt.Sprintf("task%d", i), func(ctx context.Context, input string) (string, error) {
			return input, nil
		}, fmt.Sprintf("hello %d", i))
		tm.AddTask(task) // 添加任务
	}

	tm.Run() // 执行所有任务

	// 验证当前正在执行的任务数量
	activeTasks := 0
	for _, task := range tm.GetTasks() {
		if task.GetState() == Running {
			activeTasks++
		}
	}
	assert.LessOrEqual(t, activeTasks, 2, "当前正在执行的任务数量应小于或等于最大并发数")
}

// 测试任务管理器取消正在运行的任务
func TestTaskManagerCancelRunningTasks(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task1 := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")

	task2 := NewTask[string, string, string]("task2", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "world")

	tm.AddTask(task1) // 添加任务1
	tm.AddTask(task2) // 添加任务2

	tm.Cancel("task1") // 在执行前取消任务1
	tm.Run()           // 执行任务
	tm.Cancel("task2") // 在此期间取消任务2

	// 验证任务状态
	assert.Equal(t, Cancelled, task1.GetState(), "任务1状态应为已取消")
	assert.Equal(t, Completed, task2.GetState(), "任务2状态应为已执行")
}

// 测试任务在超时情况下的行为
func TestTaskTimeout(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		select {
		case <-time.After(100 * time.Millisecond): // 模拟长时间运行的任务
			return input, nil
		case <-ctx.Done(): // 处理超时
			return "", ctx.Err()
		}
	}, "hello")

	tm.AddTask(task)
	tm.Run()

	assert.Equal(t, Completed, task.GetState(), "任务状态应为已取消")
}

// 测试任务状态转换
func TestTaskStateTransitions(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")

	assert.Equal(t, Pending, task.GetState(), "任务初始状态应为 Pending")

	tm.AddTask(task)
	assert.Equal(t, Pending, task.GetState(), "任务状态应仍为 Pending")

	tm.Run()
	assert.Equal(t, Completed, task.GetState(), "任务状态应为 Completed")
}

// 测试任务状态转换与取消
func TestTaskStateTransitionsWithCancel(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")

	assert.Equal(t, Pending, task.GetState(), "任务初始状态应为 Pending")

	tm.AddTask(task)
	assert.Equal(t, Pending, task.GetState(), "任务状态应仍为 Pending")

	tm.Cancel("task1")
	assert.Equal(t, Cancelled, task.GetState(), "任务状态应为 Cancelled")
}

// 测试任务重试次数超过最大值
func TestTaskRetriesExceedMax(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return "", errors.New("task failed") // 模拟失败
	}, "input")

	maxRetries := int32(3)
	task.SetMaxRetries(maxRetries).SetRetryInterval(10 * time.Millisecond)

	tm.AddTask(task)
	tm.Run()

	assert.Equal(t, Failed, task.GetState(), "任务状态应为失败")
	assert.Equal(t, maxRetries, task.GetRetryCount(), "任务重试次数应为 3")
}

// 测试任务在取消期间的重试
func TestTaskRetriesDuringCancellation(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return "", errors.New("任务失败") // 模拟失败
	}, "input")

	maxRetries := int32(5)
	task.SetMaxRetries(maxRetries).SetRetryInterval(10 * time.Millisecond)

	tm.AddTask(task)

	// 启动任务管理器，开始执行任务
	tm.Run()
	tm.Cancel("task1") // 取消任务

	// 验证任务状态，确保任务失败
	assert.Equal(t, Failed, task.GetState(), "任务状态应为失败")
}

// 测试任务管理器启动
func TestTaskManager_TrunUp(t *testing.T) {
	// 创建一个新的 TaskManager
	taskManager := NewTaskManager[string, string, string](2)

	// 设置启动时的自定义函数
	taskManager.SetTrunUp(func() (string, error) {
		return "Task Manager Started", nil
	})

	// 启动任务管理器
	result, err := taskManager.TrunUp()
	assert.NoError(t, err)                          // 确保没有错误
	assert.Equal(t, "Task Manager Started", result) // 确保返回结果正确
}

// 测试任务管理器关闭
func TestTaskManager_TrunDown(t *testing.T) {
	// 创建一个新的 TaskManager
	taskManager := NewTaskManager[string, string, string](2)

	// 设置关闭时的自定义函数
	taskManager.SetTrunDown(func() (string, error) {
		return "Task Manager Stopped", nil
	})

	// 关闭任务管理器
	result, err := taskManager.TrunDown()
	assert.NoError(t, err)                          // 确保没有错误
	assert.Equal(t, "Task Manager Stopped", result) // 确保返回结果正确
}

// 测试 casState 原子比较并交换（成功与失败路径）
func TestTaskCasState(t *testing.T) {
	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "input")

	// 初始状态为 Pending
	assert.Equal(t, Pending, task.GetState(), "初始状态应为 Pending")

	// CAS 成功：Pending -> Running
	assert.True(t, task.casState(Pending, Running), "CAS 应成功")
	assert.Equal(t, Running, task.GetState(), "状态应为 Running")

	// CAS 失败：old 不匹配（当前是 Running，不是 Pending）
	assert.False(t, task.casState(Pending, Completed), "CAS 应失败")
	assert.Equal(t, Running, task.GetState(), "状态应仍为 Running")

	// CAS 成功：Running -> Completed
	assert.True(t, task.casState(Running, Completed), "CAS 应成功")
	assert.Equal(t, Completed, task.GetState(), "状态应为 Completed")
}

// 测试 TaskHistory 的所有 Get 方法（成功场景）
func TestTaskHistoryGetters(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		time.Sleep(time.Millisecond) // 确保FnDuration大于0
		return input, nil
	}, "hello").SetSuccessCallback(func(result string, err error) (string, error) {
		time.Sleep(time.Millisecond) // 确保CallbackDuration大于0
		return "callback success", nil
	})

	tm.AddTask(task)
	tm.Run()

	history := task.GetTaskHistory("task1")
	assert.NotNil(t, history, "任务历史记录应存在")
	assert.Equal(t, 1, len(history), "任务历史记录应包含一条记录")

	h := history[0]
	assert.True(t, h.GetTimestamp() > int64(0), "时间戳应大于 0")
	assert.Nil(t, h.GetError(), "错误应为 nil")
	assert.True(t, h.GetFnDuration() > time.Duration(0), "函数执行时间应大于 0")
	assert.True(t, h.GetCallbackDuration() > time.Duration(0), "回调执行时间应大于 0")
	assert.Equal(t, MainTask, h.GetTaskType(), "任务类型应为主任务")
}

// 测试 TaskHistory GetError 方法（失败场景）
func TestTaskHistoryGetError(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	expectedErr := errors.New("任务执行失败")
	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return "", expectedErr
	}, "hello").SetMaxRetries(1).SetRetryInterval(time.Millisecond)

	tm.AddTask(task)
	tm.Run()

	history := task.GetTaskHistory("task1")
	assert.NotNil(t, history, "任务历史记录应存在")
	assert.Equal(t, 1, len(history), "任务历史记录应包含一条记录")

	h := history[0]
	assert.Equal(t, expectedErr, h.GetError(), "错误应匹配")
	assert.Equal(t, Failed, h.GetState(), "状态应为失败")
	assert.Equal(t, MainTask, h.GetTaskType(), "任务类型应为主任务")
	assert.True(t, h.GetTimestamp() > int64(0), "时间戳应大于 0")
}

// 测试 Task 的 Get 方法（执行前后对比）
func TestTaskGetters(t *testing.T) {
	tm := NewTaskManager[string, int, string](2)

	depTask := NewTask[string, int, string]("dep", func(ctx context.Context, input string) (int, error) {
		time.Sleep(time.Millisecond) // 确保FnDuration大于0
		return 1, nil
	}, "input")

	mainTask := NewTask[string, int, string]("main", func(ctx context.Context, input string) (int, error) {
		time.Sleep(time.Millisecond) // 确保FnDuration大于0
		return 2, nil
	}, "inputMain")

	mainTask.AddDependency(depTask)
	mainTask.SetDependExecutionMode(Sequential)

	// 执行前验证零值
	assert.Equal(t, time.Duration(0), mainTask.GetFnDuration(), "执行前 FnDuration 应为 0")
	assert.Equal(t, time.Duration(0), mainTask.GetTotalDuration(), "执行前 TotalDuration 应为 0")
	assert.Equal(t, int64(0), mainTask.GetTimestamp(), "执行前 Timestamp 应为 0")

	// 验证依赖和执行模式
	depends := mainTask.GetDepends()
	assert.Equal(t, 1, len(depends), "应有 1 个依赖任务")
	assert.Equal(t, "dep", depends[0].GetName(), "依赖任务名称应为 dep")
	assert.Equal(t, Sequential, mainTask.GetDependExecutionMode(), "执行模式应为顺序")

	// 验证依赖任务执行前的状态
	states := mainTask.GetDependencyStates()
	assert.Equal(t, Pending, states["dep"], "依赖任务状态应为 Pending")

	tm.AddTask(depTask)
	tm.AddTask(mainTask)
	tm.Run()

	// 执行后验证非零值
	assert.True(t, mainTask.GetFnDuration() > time.Duration(0), "执行后 FnDuration 应大于 0")
	assert.True(t, mainTask.GetTotalDuration() > time.Duration(0), "执行后 TotalDuration 应大于 0")
	assert.True(t, mainTask.GetTimestamp() > int64(0), "执行后 Timestamp 应大于 0")

	// 验证依赖任务执行后的状态
	states = mainTask.GetDependencyStates()
	assert.Equal(t, Completed, states["dep"], "依赖任务状态应为 Completed")
}

// 测试 TaskManager.GetTask 方法（存在与不存在）
func TestTaskManagerGetTask(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")
	tm.AddTask(task)

	// 测试获取存在的任务
	foundTask, exists := tm.GetTask("task1")
	assert.True(t, exists, "任务应存在")
	assert.NotNil(t, foundTask, "任务不应为 nil")
	assert.Equal(t, "task1", foundTask.GetName(), "任务名称应为 task1")

	// 测试获取不存在的任务
	notFoundTask, exists := tm.GetTask("nonexistent")
	assert.False(t, exists, "任务不应存在")
	assert.Nil(t, notFoundTask, "任务应为 nil")
}

// 测试 TrunUp 未设置函数时返回零值
func TestTaskManagerTrunUpNilFunc(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)
	// 不设置 trunUpFunc
	result, err := tm.TrunUp()
	assert.NoError(t, err, "未设置函数时不应返回错误")
	assert.Equal(t, "", result, "未设置函数时应返回零值")
}

// 测试 TrunDown 未设置函数时返回零值
func TestTaskManagerTrunDownNilFunc(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)
	// 不设置 trunDownFunc
	result, err := tm.TrunDown()
	assert.NoError(t, err, "未设置函数时不应返回错误")
	assert.Equal(t, "", result, "未设置函数时应返回零值")
}

// 测试取消已完成任务（不可取消状态分支）
func TestTaskCancelNonCancellableState(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")
	tm.AddTask(task)
	tm.Run()

	// 任务已完成
	assert.Equal(t, Completed, task.GetState(), "任务状态应为已完成")

	// 取消已完成的任务，状态不应改变（canCancel 为 false）
	tm.Cancel("task1")
	assert.Equal(t, Completed, task.GetState(), "已完成任务的状态不应被取消")
}

// 测试取消已取消任务（提前返回分支）
func TestTaskCancelAlreadyCancelled(t *testing.T) {
	tm := NewTaskManager[string, string, string](2)

	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "hello")
	tm.AddTask(task)

	// 第一次取消
	tm.Cancel("task1")
	assert.Equal(t, Cancelled, task.GetState(), "任务状态应为已取消")

	// 第二次取消，触发 early return（current == Cancelled）
	tm.Cancel("task1")
	assert.Equal(t, Cancelled, task.GetState(), "任务状态应仍为已取消")
}

// 测试取消带 cancel 函数和依赖的任务（cancel 非空分支与递归取消）
func TestTaskCancelWithCancelFunc(t *testing.T) {
	tm := NewTaskManager[string, int, string](2)

	depTask := NewTask[string, int, string]("dep", func(ctx context.Context, input string) (int, error) {
		return 1, nil
	}, "input")

	mainTask := NewTask[string, int, string]("main", func(ctx context.Context, input string) (int, error) {
		return 2, nil
	}, "inputMain")

	mainTask.AddDependency(depTask)

	// 设置 cancel 函数（覆盖 task.cancel != nil 分支）
	cancelCalled := false
	mainTask.cancel = func() { cancelCalled = true }
	depCancelCalled := false
	depTask.cancel = func() { depCancelCalled = true }

	tm.AddTask(mainTask)

	// 取消主任务及其依赖（递归调用）
	tm.Cancel("main")

	assert.True(t, cancelCalled, "主任务的 cancel 函数应被调用")
	assert.True(t, depCancelCalled, "依赖任务的 cancel 函数应被调用")
	assert.Equal(t, Cancelled, mainTask.GetState(), "主任务状态应为已取消")
	assert.Equal(t, Cancelled, depTask.GetState(), "依赖任务状态应为已取消")
}

// 测试 runWithRetries 在已取消状态下的行为（getState == Cancelled 分支）
func TestTaskRunWithRetriesCancelledState(t *testing.T) {
	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return "", errors.New("fail")
	}, "input")

	// 直接设置为已取消状态
	task.setState(Cancelled)
	task.maxRetries = 3

	result, err := task.runWithRetries()
	assert.Equal(t, "", result, "结果应为零值")
	assert.NoError(t, err, "已取消任务不应返回错误")
	assert.Equal(t, Cancelled, task.GetState(), "状态应仍为已取消")
}

// 测试 runWithRetries 在上下文取消时的行为（ctx.Done 分支）
func TestTaskRunWithRetriesContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消上下文

	task := NewTaskWithOptions[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "input", ctx, 3, 1*time.Second)

	result, err := task.runWithRetries()
	assert.Equal(t, "", result, "结果应为零值")
	assert.NoError(t, err, "上下文取消不应返回错误")
	assert.Equal(t, Cancelled, task.GetState(), "状态应为已取消")
}

// 测试 logHistory 的历史记录裁剪逻辑（maxHistorySize > 0 分支）
func TestTaskLogHistoryTrimming(t *testing.T) {
	task := NewTask[string, string, string]("task1", func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, "input")

	task.SetMaxHistorySize(2)

	// 调用 logHistory 三次，触发裁剪逻辑
	task.logHistory()
	task.logHistory()
	task.logHistory()

	history := task.GetTaskHistory("task1")
	assert.Equal(t, 2, len(history), "历史记录应被裁剪为最大数量")
}
