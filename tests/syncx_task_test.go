/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-05 16:31:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-05 16:31:35
 * @FilePath: \go-toolbox\tests\syncx_task_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

// 模拟正常任务
func normalTask(name string, delayMs int) syncx.TaskFunc[string] {
	return func() (string, error) {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
		return "result of " + name, nil
	}
}

// 模拟返回错误任务
func errorTask(name string) syncx.TaskFunc[string] {
	return func() (string, error) {
		return "", errors.New("error in " + name)
	}
}

// 模拟panic任务
func panicTask(name string) syncx.TaskFunc[string] {
    return func() (string, error) {
        panic("task " + name + " panicked")
    }
}

func TestTaskRunner_Run_WithMap(t *testing.T) {
	assert := assert.New(t) // 创建断言对象
	// 定义测试用例集合，key是任务名，value是任务函数和期望的错误标志
	testCases := map[string]struct {
		task        syncx.TaskFunc[string]
		expectError bool
	}{
		"task1": {task: normalTask("task1", 50), expectError: false},
		"task2": {task: normalTask("task2", 100), expectError: false},
		"task3": {task: errorTask("task3"), expectError: true},
		"task4": {task: panicTask("task4"), expectError: true},
	}

	tr := syncx.NewTaskRunner[string]()

	// 添加任务
	for key, tc := range testCases {
		tr.Add(key, tc.task)
	}

	// 执行任务
	results := tr.Run()

	// 校验结果
	for key, tc := range testCases {
		result, ok := results[key]
		assert.True(ok, "missing result for task %s", key)

		if tc.expectError {
			assert.Error(result.Err, "expected error for task %s", key)
		} else {
			assert.NoError(result.Err, "unexpected error for task %s", key)
			assert.NotEmpty(result.Result, "empty result for task %s", key)
		}
	}
}
