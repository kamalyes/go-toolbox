/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-24 13:55:29
 * @FilePath: \go-toolbox\pkg\schedule\job.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"fmt"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// validateJobParameters 验证任务参数的有效性
func validateJobParameters(id, name, expression string, fn func() error) error {
	if id == "" {
		return fmt.Errorf("job id cannot be empty")
	}
	if name == "" {
		return fmt.Errorf("job name cannot be empty")
	}
	if fn == nil {
		return fmt.Errorf("job callback function cannot be nil")
	}
	if expression == "" {
		return fmt.Errorf("cron expression cannot be empty")
	}
	return nil
}

// AddJob 添加任务
func (c *Cron) AddJob(id, name, expression string, fn func() error) *JobRule {
	// 参数有效性检查
	if err := validateJobParameters(id, name, expression, fn); err != nil {
		panic(err.Error())
	}

	return syncx.WithLockReturnValue(&c.runningMux, func() *JobRule {
		if _, exists := c.jobs[id]; exists {
			panic(fmt.Sprintf("job with id %s already exists", id))
		}

		job := &JobRule{
			id:                  id,
			name:                name,
			expression:          expression,
			callback:            fn,
			exceedTaskSnapshots: make(map[string]*ExceedTaskSnapshot), // 初始化快照映射
		}

		c.jobs[id] = job
		return job
	})
}

// DelJob 删除任务
func (c *Cron) DelJob(id string) *Cron {
	return syncx.WithLockReturnValue(&c.runningMux, func() *Cron {
		delete(c.jobs, id)
		return c
	})
}

// AbortJob 终止任务
func (c *Cron) AbortJob(id string) *Cron {
	return syncx.WithLockReturnValue(&c.runningMux, func() *Cron {
		if job, exists := c.jobs[id]; exists {
			job.Abort()
		}
		return c
	})
}

// GetJobStatus 获取任务状态
func (c *Cron) GetJobStatus(id string) *JobRule {
	return syncx.WithRLockReturnValue(&c.runningMux, func() *JobRule {
		if job, exists := c.jobs[id]; exists {
			return job // 返回对应的任务规则
		}
		return nil // 或者返回一个错误指示任务不存在
	})
}
