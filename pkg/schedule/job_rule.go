/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-24 13:26:15
 * @FilePath: \go-toolbox\pkg\schedule\job_rule.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package schedule

import (
	"time"
)

type CustomFunc func() error // 自定义函数类型

// JobRule 结构体表示一个Job规则
type JobRule struct {
	id                  string                         // 作业Id
	name                string                         // 作业名称
	nextTime            time.Time                      // 此作业下一个运行的时间，如果 Cron 尚未启动或此条目的计划不可满足，则为零时间
	prevTime            time.Time                      // 此作业最后一次运行的时间，如果从未运行则为零时间
	cooldownDuration    time.Duration                  // 冷却时间
	sleepDuration       time.Duration                  // 睡眠时间
	maxFailureCount     int                            // 最大的失败次数
	expression          string                         // Cron 表达式
	callback            func() error                   // 任务回调函数
	beforeFunc          func()                         // 执行前的函数
	afterFunc           func()                         // 执行后的函数
	skipFunc            func() bool                    // 跳过执行的条件函数
	exceedTaskSnapshots map[string]*ExceedTaskSnapshot // 对应的任务快照信息
	abort               bool                           // 终止、后续不会再执行
}

// GetId 获取Job名称
func (t *JobRule) GetId() string {
	return t.id
}

// GetName 获取名称
func (t *JobRule) GetName() string {
	return t.name
}

// GetNextTime 获取下次运行时间
func (t *JobRule) GetNextTime() time.Time {
	return t.nextTime
}

// GetPrevTime 获取上次运行时间
func (t *JobRule) GetPrevTime() time.Time {
	return t.prevTime
}

// SetId 设置Id
func (t *JobRule) SetId(id string) *JobRule {
	t.id = id
	return t
}

// SetName 设置名称
func (t *JobRule) SetName(name string) *JobRule {
	t.name = name
	return t
}

// SetNextTime 设置下次运行时间
func (t *JobRule) SetNextTime(nextTime time.Time) *JobRule {
	t.nextTime = nextTime
	return t
}

// SetPrevTime 设置上次运行时间
func (t *JobRule) SetPrevTime(prevTime time.Time) *JobRule {
	t.prevTime = prevTime
	return t
}

// GetCooldownDuration 获取冷却时间
func (t *JobRule) GetCooldownDuration() time.Duration {
	return t.cooldownDuration
}

// GetSleepDuration 获取睡眠时间
func (t *JobRule) GetSleepDuration() time.Duration {
	return t.sleepDuration
}

// SetCooldownDuration 设置冷却时间
func (t *JobRule) SetCooldownDuration(duration time.Duration) *JobRule {
	t.cooldownDuration = duration
	return t
}

// SetSleepDuration 设置睡眠时间
func (t *JobRule) SetSleepDuration(duration time.Duration) *JobRule {
	t.sleepDuration = duration
	return t
}

// GetExpression 获取 Cron 表达式
func (t *JobRule) GetExpression() string {
	return t.expression
}

// SetExpression 设置 Cron 表达式
func (t *JobRule) SetExpression(expression string) *JobRule {
	t.expression = expression
	return t
}

// GetMaxFailureCount 获取最大失败次数
func (t *JobRule) GetMaxFailureCount() int {
	return t.maxFailureCount
}

// SetMaxFailureCount 设置最大失败次数
func (t *JobRule) SetMaxFailureCount(count int) *JobRule {
	t.maxFailureCount = count
	return t
}

// GetCallback 获取任务回调函数
func (t *JobRule) GetCallback() func() error {
	return t.callback
}

// SetCallback 设置任务回调函数
func (t *JobRule) SetCallback(callback func() error) *JobRule {
	t.callback = callback
	return t
}

// GetBeforeFunc 获取执行前的函数
func (t *JobRule) GetBeforeFunc() func() {
	return t.beforeFunc
}

// SetBeforeFunc 设置执行前的函数
func (t *JobRule) SetBeforeFunc(beforeFunc func()) *JobRule {
	t.beforeFunc = beforeFunc
	return t
}

// GetAfterFunc 获取执行后的函数
func (t *JobRule) GetAfterFunc() func() {
	return t.afterFunc
}

// SetAfterFunc 设置执行后的函数
func (t *JobRule) SetAfterFunc(afterFunc func()) *JobRule {
	t.afterFunc = afterFunc
	return t
}

// GetSkipFunc 获取跳过执行的条件函数
func (t *JobRule) GetSkipFunc() func() bool {
	return t.skipFunc
}

// SetSkipFunc 设置跳过执行的条件函数
func (t *JobRule) SetSkipFunc(skipFunc func() bool) *JobRule {
	t.skipFunc = skipFunc
	return t
}

func (t *JobRule) Abort() *JobRule {
	t.abort = true
	return t
}

func (t *JobRule) Aborted() bool {
	return t.abort
}
