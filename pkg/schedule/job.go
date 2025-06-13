/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 15:51:15
 * @FilePath: \go-toolbox\pkg\schedule\job.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"errors"
	"fmt"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// validateJobParameters 验证任务参数的有效性
//
// Params:
//   - job: 需要验证的任务规则指针
//
// Returns:
//   - error: 如果参数无效返回对应错误，否则返回nil
//
// 说明:
//
//	该函数用于校验传入的任务规则的各个必要字段是否符合要求。
//	如果任务规则为空，直接panic，因为至少要有一条规则。
//	如果任务ID、名称、回调函数或表达式为空，则返回对应错误。
func validateJobParameters(job *JobRule) error {
	if job == nil {
		// 任务规则不能为空，直接panic
		panic(errors.New(ErrAtLeastOneRuleMustBeDefined))
	}
	if job.id == "" {
		return errors.New(ErrJobIDEmpty) // 任务ID不能为空
	}
	if job.name == "" {
		return errors.New(ErrJobNameEmpty) // 任务名称不能为空
	}
	if job.callback == nil {
		return errors.New(ErrJobCallbackNil) // 任务回调函数不能为空
	}
	if job.expression == "" {
		return errors.New(ErrJobExpressionEmpty) // 任务表达式不能为空
	}
	return nil // 所有参数合法，返回nil
}

// AddJob 添加任务到调度器
//
// Params:
//   - job: 任务规则指针，包含任务的所有信息
//
// Returns:
//   - *Schedule: 返回调度器本身，支持链式调用
//
// 说明:
//
//	该方法先调用 validateJobParameters 对任务参数进行校验，
//	如果校验失败会panic，确保传入的任务参数合法。
//	通过加锁保证并发安全，检查任务ID是否重复，避免覆盖已有任务。
//	初始化任务的 exceedTaskSnapshots（如果未初始化），
//	最后将任务加入调度器的任务列表中。
//	返回调度器本身，方便链式调用。
func (c *Schedule) AddJob(job *JobRule) *Schedule {
	// 参数校验，失败时panic
	if err := validateJobParameters(job); err != nil {
		panic(err.Error())
	}

	// 并发安全地向任务列表中添加任务
	return syncx.WithLockReturnValue(&c.mu, func() *Schedule {
		// 检查任务ID是否已存在，避免重复添加
		if _, exists := c.jobs[job.id]; exists {
			panic(fmt.Sprintf(ErrJobIDAlreadyExists, job.id))
		}
		// 初始化任务的exceedTaskSnapshots，防止后续空指针异常
		if job.exceedTaskSnapshots == nil {
			job.exceedTaskSnapshots = make(map[string]*ExceedTaskSnapshot)
		}
		// 添加任务到调度器任务列表
		c.jobs[job.id] = job

		// 返回调度器实例，支持链式调用
		return c
	})
}

// DelJob 删除任务
func (c *Schedule) DelJob(id string) *Schedule {
	return syncx.WithLockReturnValue(&c.mu, func() *Schedule {
		delete(c.jobs, id)
		return c
	})
}

// AbortJob 终止任务
func (c *Schedule) AbortJob(id string) *Schedule {
	return syncx.WithLockReturnValue(&c.mu, func() *Schedule {
		if job, exists := c.jobs[id]; exists {
			job.Abort()
		}
		return c
	})
}

// GetJob 获取任务信息
func (c *Schedule) GetJob(id string) *JobRule {
	return syncx.WithRLockReturnValue(&c.mu, func() *JobRule {
		if job, exists := c.jobs[id]; exists {
			return job // 返回对应的任务规则
		}
		return nil // 或者返回一个错误指示任务不存在
	})
}
