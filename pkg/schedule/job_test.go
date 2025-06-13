/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 17:55:19
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
	c := &Schedule{jobs: make(map[string]*JobRule)}

	tests := []struct {
		job         *JobRule
		expectPanic bool
		expectedErr string
	}{
		{&JobRule{id: "job1", name: "Test Job 1", expression: "*/5 * * * *", callback: func() error { return nil }}, false, ""},
		{&JobRule{id: "", name: "Test Job 2", expression: "*/5 * * * *", callback: func() error { return nil }}, true, ErrJobIDEmpty},
		{&JobRule{id: "job2", name: "", expression: "*/5 * * * *", callback: func() error { return nil }}, true, ErrJobNameEmpty},
		{&JobRule{id: "job3", name: "Test Job 3", expression: "", callback: func() error { return nil }}, true, ErrJobExpressionEmpty},
		{&JobRule{id: "job4", name: "Test Job 4", expression: "*/5 * * * *", callback: nil}, true, ErrJobCallbackNil},
		{&JobRule{id: "job1", name: "Test Job 5", expression: "*/5 * * * *", callback: func() error { return nil }}, true, fmt.Sprintf(ErrJobIDAlreadyExists, "job1")},
	}

	for _, tt := range tests {
		if tt.expectPanic {
			assert.PanicsWithValue(t, tt.expectedErr, func() {
				c.AddJob(tt.job)
			}, fmt.Sprintf("expected panic for input %+v, but did not", tt.job))
		} else {
			c.AddJob(tt.job)
			assert.NotNil(t, c.jobs[tt.job.id], "Job should be added")
		}
	}
}

func TestConcurrentAddJobs(t *testing.T) {
	schedule := NewSchedule()
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

			var beforeExecuted, afterExecuted bool
			beforeFunc := func() { beforeExecuted = true }
			afterSuccessFunc := func() { afterExecuted = true }
			timezone, _ := time.LoadLocation("America/New_York")
			jobConfig := &JobRule{
				id:               fmt.Sprintf("job%d", id),
				name:             fmt.Sprintf("Complex Job %d%s", id, osx.HashUnixMicroCipherText()),
				expression:       fmt.Sprintf("* * * * %d", id),
				cooldownDuration: 5 * time.Second,
				sleepDuration:    2 * time.Second,
				maxFailureCount:  3,
				callback:         jobFunc,
				beforeFunc:       beforeFunc,
				afterSuccessFunc: afterSuccessFunc,
				timezone:         timezone,
			}

			schedule.AddJob(jobConfig)

			job := schedule.GetJob(jobConfig.id)
			// 断言任务的各个属性是否被正确设置

			assert.NotNil(t, job, "Job should not be nil after addition")
			assert.Equal(t, jobConfig.id, job.GetId(), "Job ID should match")
			assert.Equal(t, jobConfig.name, job.GetName(), "Job Name should match")
			assert.Equal(t, jobConfig.expression, job.GetExpression(), "Job Expression should match")
			assert.Equal(t, jobConfig.cooldownDuration, job.GetCooldownDuration(), "Cooldown Duration should match")
			assert.Equal(t, jobConfig.sleepDuration, job.GetSleepDuration(), "Sleep Duration should match")
			assert.Equal(t, jobConfig.maxFailureCount, job.GetMaxFailureCount(), "Max Failure Count should match")
			assert.Equal(t, jobConfig.timezone, job.GetTimezone(), "Timezone Count should match")

			// 执行前后的函数应该被调用
			job.GetBeforeFunc()()
			assert.True(t, beforeExecuted, "Before function should be executed")

			job.GetAfterSuccessFunc()()
			assert.True(t, afterExecuted, "AfterSuccessFunc function should be executed")
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 检查作业数量
	assert.Equal(t, jobCount, len(schedule.jobs), "The number of jobs in the schedule should match the number of added jobs")

	// 可选：检查每个作业是否存在
	for i := 0; i < jobCount; i++ {
		jobID := fmt.Sprintf("job%d", i)
		job, exists := schedule.jobs[jobID]
		assert.True(t, exists, "Job should exist in the schedule jobs")
		assert.Equal(t, jobID, job.GetId(), "Job ID should match")
	}
}

func TestDelJob(t *testing.T) {
	schedule := NewSchedule()
	job := &JobRule{
		id:         "job1",
		name:       "Test Job",
		expression: "* * * * *",
		callback:   func() error { return nil },
	}
	schedule.AddJob(job)

	// 确认任务已添加
	assert.NotNil(t, schedule.GetJob(job.id))

	// 删除任务
	schedule.DelJob(job.id)

	// 删除后任务应不存在
	assert.Nil(t, schedule.GetJob(job.id))
}

func TestAbortJob(t *testing.T) {
	schedule := NewSchedule()
	job := &JobRule{
		id:         "job2",
		name:       "Test Job 2",
		expression: "* * * * *",
		callback:   func() error { return nil },
	}
	schedule.AddJob(job)

	// 任务初始状态未被终止
	assert.False(t, job.Aborted())

	// 终止任务
	schedule.AbortJob(job.id)

	// 任务状态应被标记为终止
	abortedJob := schedule.GetJob(job.id)
	assert.NotNil(t, abortedJob)
	assert.True(t, abortedJob.Aborted())
}

func TestGetJob(t *testing.T) {
	schedule := NewSchedule()
	job := &JobRule{
		id:         "job3",
		name:       "Test Job 3",
		expression: "* * * * *",
		callback:   func() error { return nil },
	}
	schedule.AddJob(job)

	// 获取已存在任务状态
	retrieved := schedule.GetJob(job.id)
	assert.NotNil(t, retrieved)
	assert.Equal(t, job.id, retrieved.GetId())

	// 获取不存在任务状态应返回 nil
	nonExistent := schedule.GetJob("nonexistent")
	assert.Nil(t, nonExistent)
}
