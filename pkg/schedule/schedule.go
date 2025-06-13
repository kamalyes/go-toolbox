/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-24 13:26:15
 * @FilePath: \go-toolbox\pkg\schedule\schedule.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"errors"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// Schedule 结构体表示一个定时器 负责跟踪任意数量的条目，并根据指定的计划调用相关的函数。
// 它可以被启动、停止，并且在运行时可以检查条目。
type Schedule struct {
	jobs        map[string]*JobRule // 任务映射，使用任务ID作为键
	isBroken    bool                // 标记定时器是否处于故障状态
	mu          sync.RWMutex        // 读写锁
	jobWaiter   sync.WaitGroup      // 用于等待任务完成
	taskChan    chan *JobRule       // 任务通道
	workerCount int                 // 工作池大小
}

// NewSchedule 初始化创建一个新的 Schedule 实例
func NewSchedule() *Schedule {
	return &Schedule{
		jobs:        make(map[string]*JobRule), // 初始化为一个空的映射
		taskChan:    make(chan *JobRule, 100),  // 创建带缓冲的通道
		workerCount: 20,                        // 默认工作池大小
	}
}

// GetJobRulesCopy 获取调度规则的副本
func (c *Schedule) GetJobRulesCopy() []*JobRule {
	return syncx.WithRLockReturnValue(&c.mu, func() []*JobRule {
		jobsCopy := make([]*JobRule, 0, len(c.jobs)) // 创建规则副本切片
		for _, rule := range c.jobs {                // 遍历规则映射
			jobsCopy = append(jobsCopy, rule) // 添加每个规则到副本中
		}
		return jobsCopy
	})
}

// GetJobRules 获取调度规则
func (c *Schedule) GetJobRules() []*JobRule {
	return syncx.WithRLockReturnValue(&c.mu, func() []*JobRule {
		jobs := make([]*JobRule, 0, len(c.jobs)) // 创建规则切片
		for _, rule := range c.jobs {            // 遍历规则映射
			jobs = append(jobs, rule) // 添加每个规则到切片中
		}
		return jobs
	})
}

// IsBroken 获取定时器是否处于熔断状态
func (c *Schedule) IsBroken() bool {
	return syncx.WithRLockReturnValue(&c.mu, func() bool {
		return c.isBroken
	})
}

// GetWorkerCount 获取工作池大小
func (c *Schedule) GetWorkerCount() int {
	return syncx.WithRLockReturnValue(&c.mu, func() int {
		return c.workerCount
	})
}

// SetTaskChanCapacity 设置任务通道的容量
func (c *Schedule) SetTaskChanCapacity(capacity int) *Schedule {
	return syncx.WithLockReturnValue(&c.mu, func() *Schedule {
		// 关闭旧通道，创建新通道
		close(c.taskChan)
		c.taskChan = make(chan *JobRule, capacity)
		return c
	})
}

// SetWorkerCount 设置工作池的大小
func (c *Schedule) SetWorkerCount(count int) *Schedule {
	return syncx.WithLockReturnValue(&c.mu, func() *Schedule {
		c.workerCount = count
		return c
	})
}

// Validate 检查 Schedule 的配置是否有效
func (c *Schedule) Validate() error {
	return syncx.WithRLockReturnValue(&c.mu, func() error {
		if len(c.jobs) == 0 {
			return errors.New(ErrAtLeastOneRuleMustBeDefined)
		}
		return nil
	})
}
