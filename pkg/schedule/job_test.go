/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-24 13:27:59
 * @FilePath: \go-toolbox\pkg\schedule\job_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/stretchr/testify/assert"
)

// TestAddJob 测试 AddJob 方法
func TestAddJob(t *testing.T) {
	c := &Cron{jobs: make(map[string]*JobRule)}

	tests := []struct {
		id          string
		name        string
		expression  string
		fn          func() error
		expectPanic bool
		expectedErr string
	}{
		{"job1", "Test Job 1", "*/5 * * * *", func() error { return nil }, false, ""},
		{"", "Test Job 2", "*/5 * * * *", func() error { return nil }, true, "job id cannot be empty"},
		{"job2", "", "*/5 * * * *", func() error { return nil }, true, "job name cannot be empty"},
		{"job3", "Test Job 3", "", func() error { return nil }, true, "cron expression cannot be empty"},
		{"job4", "Test Job 4", "*/5 * * * *", nil, true, "job callback function cannot be nil"},
		{"job1", "Test Job 5", "*/5 * * * *", func() error { return nil }, true, "job with id job1 already exists"},
	}

	for _, tt := range tests {
		if tt.expectPanic {
			assert.PanicsWithValue(t, tt.expectedErr, func() {
				c.AddJob(tt.id, tt.name, tt.expression, tt.fn)
			}, fmt.Sprintf("expected panic for input %v, but did not", tt))
		} else {
			c.AddJob(tt.id, tt.name, tt.expression, tt.fn)
			assert.NotNil(t, c.jobs[tt.id], "Job should be added")
		}
	}
}

func TestConcurrentAddJobs(t *testing.T) {
	cron := NewCron()
	var wg sync.WaitGroup

	// 定义回调函数
	jobFunc := func() error {
		return nil
	}

	// 并发添加多个任务
	jobCount := 10
	for i := 0; i < jobCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 定义任务配置
			jobConfig := JobRule{
				id:               fmt.Sprintf("job%d", id),
				name:             fmt.Sprintf("Complex Job %d%s", id, osx.HashUnixMicroCipherText()),
				expression:       fmt.Sprintf("* * * * %d", id),
				cooldownDuration: 5 * time.Second,
				sleepDuration:    2 * time.Second,
				maxFailureCount:  3,
			}

			// 定义执行前和执行后的函数
			var beforeExecuted, afterExecuted bool
			beforeFunc := func() {
				beforeExecuted = true
			}
			afterFunc := func() {
				afterExecuted = true
			}

			// 添加任务时设置更多属性
			job := cron.AddJob(jobConfig.id, jobConfig.name, jobConfig.expression, jobFunc).
				SetCooldownDuration(jobConfig.cooldownDuration).
				SetSleepDuration(jobConfig.sleepDuration).
				SetMaxFailureCount(jobConfig.maxFailureCount).
				SetBeforeFunc(beforeFunc).
				SetAfterFunc(afterFunc)

			// 断言任务的各个属性是否被正确设置
			assert.NotNil(t, job, "Job should not be nil after addition")
			assert.Equal(t, jobConfig.id, job.GetId(), "Job ID should match")
			assert.Equal(t, jobConfig.name, job.GetName(), "Job Name should match")
			assert.Equal(t, jobConfig.expression, job.GetExpression(), "Job Expression should match")
			assert.Equal(t, jobConfig.cooldownDuration, job.GetCooldownDuration(), "Cooldown Duration should match")
			assert.Equal(t, jobConfig.sleepDuration, job.GetSleepDuration(), "Sleep Duration should match")
			assert.Equal(t, jobConfig.maxFailureCount, job.GetMaxFailureCount(), "Max Failure Count should match")

			// 执行前后的函数应该被调用
			job.GetBeforeFunc()()
			assert.True(t, beforeExecuted, "Before function should be executed")

			job.GetAfterFunc()()
			assert.True(t, afterExecuted, "After function should be executed")
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 检查作业数量
	assert.Equal(t, jobCount, len(cron.jobs), "The number of jobs in the cron should match the number of added jobs")

	// 可选：检查每个作业是否存在
	for i := 0; i < jobCount; i++ {
		jobID := fmt.Sprintf("job%d", i)
		job, exists := cron.jobs[jobID]
		assert.True(t, exists, "Job should exist in the cron jobs")
		assert.Equal(t, jobID, job.GetId(), "Job ID should match")
	}
}

func TestGetJobStatus(t *testing.T) {
	cron := NewCron()
	jobFunc := func() error { return nil }
	cron.AddJob("job1", "Test Job", "* * * * *", jobFunc)

	// 测试获取任务状态
	retrievedJob := cron.GetJobStatus("job1")
	assert.NotNil(t, retrievedJob, "Retrieved job should not be nil")
	assert.Equal(t, "job1", retrievedJob.GetId(), "Retrieved job ID should match")
}

func TestDelJob(t *testing.T) {
	cron := NewCron()
	jobFunc := func() error { return nil }
	cron.AddJob("job1", "Test Job", "* * * * *", jobFunc)

	// 测试删除任务
	cron.DelJob("job1")
	assert.Nil(t, cron.GetJobStatus("job1"), "Job should be nil after deletion")
}

func TestAbortJob(t *testing.T) {
	cron := NewCron()
	jobFunc := func() error { return nil }
	cron.AddJob("job2", "Test Job 2", "* * * * *", jobFunc)

	// 测试终止任务
	cron.AbortJob("job2")
	abortedJob := cron.GetJobStatus("job2")
	assert.NotNil(t, abortedJob, "Aborted job should still exist")
	assert.True(t, abortedJob.Aborted(), "Job should be marked as aborted")
}
