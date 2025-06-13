/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-24 15:00:55
 * @FilePath: \go-toolbox\pkg\schedule\worker.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"fmt"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// Start 启动 Schedule 定时器
func (c *Schedule) Start() {
	syncx.WithLock(&c.mu, func() {
		if c.isBroken {
			fmt.Println("Schedule is broken, cannot start.")
			return
		}
		go c.run()
	})
}

// Stop 停止 Schedule 定时器
func (c *Schedule) Stop() {
	syncx.WithLock(&c.mu, func() {
		c.isBroken = true
		close(c.taskChan)
		c.jobWaiter.Wait() // 等待所有任务完成
	})
}

// run 运行定时器，调度任务
func (c *Schedule) run() {
	ticker := time.NewTicker(time.Second) // 每秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.scheduleTasks() // 调度任务
		case task := <-c.taskChan:
			c.executeTask(task) // 执行任务
		}
	}
}

// scheduleTasks 调度任务
func (c *Schedule) scheduleTasks() {
	syncx.WithRLock(&c.mu, func() {
		for _, rule := range c.jobs {
			if time.Now().After(rule.GetNextTime()) {
				c.taskChan <- rule
				rule.SetPrevTime(time.Now())
				rule.SetNextTime(c.calculateNextRunTime()) // 计算下次运行时间
			}
		}
	})
}

// calculateNextRunTime 计算下次运行时间
func (c *Schedule) calculateNextRunTime() time.Time {
	// 这里可以实现根据 Schedule 表达式计算下次运行时间的逻辑
	// 目前仅返回当前时间加上1分钟作为示例
	return time.Now().Add(time.Minute)
}

// executeTask 执行任务
func (c *Schedule) executeTask(rule *JobRule) {
	defer c.jobWaiter.Done() // 确保任务完成时调用 Done

	// 生成一个新的 traceId
	traceId := fmt.Sprintf("%s-%s", rule.GetId(), osx.HashUnixMicroCipherText())

	// 执行前的函数
	if rule.GetBeforeFunc() != nil {
		rule.GetBeforeFunc()
	}

	// 执行任务回调
	err := rule.GetCallback()()
	if err != nil {
		c.recordTaskFailure(rule, traceId, err)
	} else {
		c.recordTaskSuccess(rule, traceId)
	}

	// 执行后的函数 细化执行后回调，区分成功和失败，安全调用
	mathx.IfChain([]mathx.ConditionValue[func()]{
		{Cond: err != nil, Value: mathx.IF(rule.GetAfterFailureFunc() != nil, rule.GetAfterFailureFunc(), func() {})},
		{Cond: err == nil, Value: mathx.IF(rule.GetAfterSuccessFunc() != nil, rule.GetAfterSuccessFunc(), func() {})},
	}, func() {})()
}

// recordTaskStatus 记录任务状态信息
func (c *Schedule) recordTaskStatus(rule *JobRule, traceId string, status execStatus, err error) {
	syncx.WithLock(&c.mu, func() {
		// 查找是否已存在对应的快照
		currentSnapshot, exists := rule.exceedTaskSnapshots[traceId]

		// 如果该 traceId 不存在，初始化快照
		if !exists {
			currentSnapshot = NewExceedTaskSnapshot()
			currentSnapshot.SetTraceId(traceId) // 设置 traceId
			currentSnapshot.SetExecStatus(status)
			rule.exceedTaskSnapshots[traceId] = currentSnapshot // 添加新快照
		}

		// 更新执行频率
		currentSnapshot.execFrequency++

		// 根据状态更新快照
		if messageTemplate, ok := execStatusLogMessages[status]; ok {
			switch status {
			case Failure:
				currentSnapshot.failureFrequency++
				currentSnapshot.execLogRecord = append(currentSnapshot.execLogRecord, fmt.Sprintf(messageTemplate, err))
			default:
				currentSnapshot.execLogRecord = append(currentSnapshot.execLogRecord, messageTemplate)
			}
		}

		// 更新执行状态
		currentSnapshot.execStatus = status
	})
}

// recordTaskFailure 记录任务失败信息
func (c *Schedule) recordTaskFailure(rule *JobRule, traceId string, err error) {
	c.recordTaskStatus(rule, traceId, Failure, err)
}

// recordTaskSuccess 记录任务成功信息
func (c *Schedule) recordTaskSuccess(rule *JobRule, traceId string) {
	c.recordTaskStatus(rule, traceId, Success, nil)
}

// recordTaskPending 记录任务等待中信息
func (c *Schedule) recordTaskPending(rule *JobRule, traceId string) {
	c.recordTaskStatus(rule, traceId, Pending, nil)
}

// recordTaskRunning 记录任务运行中信息
func (c *Schedule) recordTaskRunning(rule *JobRule, traceId string) {
	c.recordTaskStatus(rule, traceId, Running, nil)
}

// recordTaskSysTermination 记录任务系统终止信息
func (c *Schedule) recordTaskSysTermination(rule *JobRule, traceId string) {
	c.recordTaskStatus(rule, traceId, SysTermination, nil)
}

// recordTaskUserTermination 记录任务用户终止信息
func (c *Schedule) recordTaskUserTermination(rule *JobRule, traceId string) {
	c.recordTaskStatus(rule, traceId, UserTermination, nil)
}
