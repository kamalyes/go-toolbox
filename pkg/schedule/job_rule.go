/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 17:17:15
 * @FilePath: \go-toolbox\pkg\schedule\job_rule.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package schedule

import (
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

type CustomFunc func() error // 自定义函数类型

// JobRule 结构体表示一个Job规则
type JobRule struct {
	id                  string                         // 作业Id
	name                string                         // 作业名称
	timezone            *time.Location                 // 时区
	nextTime            time.Time                      // 此作业下一个运行的时间，如果 Cron 尚未启动或此条目的计划不可满足，则为零时间
	prevTime            time.Time                      // 此作业最后一次运行的时间，如果从未运行则为零时间
	timeout             time.Duration                  // 执行超时时间，0 表示不超时
	cooldownDuration    time.Duration                  // 冷却时间
	sleepDuration       time.Duration                  // 睡眠时间
	maxFailureCount     int                            // 最大的失败次数
	expression          string                         // Cron 表达式
	callback            func() error                   // 任务回调函数
	beforeFunc          func()                         // 执行前的函数
	afterSuccessFunc    func()                         // 执行成功后的回调
	afterFailureFunc    func()                         // 执行失败后的回调
	skipFunc            func() bool                    // 跳过执行的条件函数
	exceedTaskSnapshots map[string]*ExceedTaskSnapshot // 对应的任务快照信息
	abort               bool                           // 终止、后续不会再执行
	mu                  sync.RWMutex                   // 读写锁
}

// GetId 获取Job名称
func (t *JobRule) GetId() string {
	return syncx.WithRLockReturnValue(&t.mu, func() string {
		return t.id
	})
}

// GetName 获取名称
func (t *JobRule) GetName() string {
	return syncx.WithRLockReturnValue(&t.mu, func() string {
		return t.name
	})
}

// GetTimezone 获取时区
func (t *JobRule) GetTimezone() *time.Location {
	return syncx.WithRLockReturnValue(&t.mu, func() *time.Location {
		return t.timezone
	})
}

// GetNextTime 获取下次运行时间
func (t *JobRule) GetNextTime() time.Time {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Time {
		return t.nextTime
	})
}

// GetPrevTime 获取上次运行时间
func (t *JobRule) GetPrevTime() time.Time {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Time {
		return t.prevTime
	})
}

// GetTimeout 获取执行超时时间
func (t *JobRule) GetTimeout() time.Duration {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Duration {
		return t.timeout
	})
}

// SetTimeout 设置执行超时时间
func (t *JobRule) SetTimeout(timeout time.Duration) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.timeout = timeout
		return t
	})
}

// SetTimezone 设置时区
func (t *JobRule) SetTimezone(tz string) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.timezone = mathx.IfDoWithErrorDefault(true, func() (*time.Location, error) {
			return time.LoadLocation(tz)
		}, DefaultTimeZone)
		return t
	})
}

// SetId 设置Id
func (t *JobRule) SetId(id string) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.id = id
		return t
	})
}

// SetName 设置名称
func (t *JobRule) SetName(name string) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.name = name
		return t
	})
}

// SetNextTime 设置下次运行时间
func (t *JobRule) SetNextTime(nextTime time.Time) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.nextTime = nextTime
		return t
	})
}

// SetPrevTime 设置上次运行时间
func (t *JobRule) SetPrevTime(prevTime time.Time) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.prevTime = prevTime
		return t
	})
}

// GetCooldownDuration 获取冷却时间
func (t *JobRule) GetCooldownDuration() time.Duration {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Duration {
		return t.cooldownDuration
	})
}

// GetSleepDuration 获取睡眠时间
func (t *JobRule) GetSleepDuration() time.Duration {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Duration {
		return t.sleepDuration
	})
}

// SetCooldownDuration 设置冷却时间
func (t *JobRule) SetCooldownDuration(duration time.Duration) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.cooldownDuration = duration
		return t
	})
}

// SetSleepDuration 设置睡眠时间
func (t *JobRule) SetSleepDuration(duration time.Duration) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.sleepDuration = duration
		return t
	})
}

// GetExpression 获取 Cron 表达式
func (t *JobRule) GetExpression() string {
	return syncx.WithRLockReturnValue(&t.mu, func() string {
		return t.expression
	})
}

// SetExpression 设置 Cron 表达式
func (t *JobRule) SetExpression(expression string) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.expression = expression
		return t
	})
}

// GetMaxFailureCount 获取最大失败次数
func (t *JobRule) GetMaxFailureCount() int {
	return syncx.WithRLockReturnValue(&t.mu, func() int {
		return t.maxFailureCount
	})
}

// SetMaxFailureCount 设置最大失败次数
func (t *JobRule) SetMaxFailureCount(count int) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.maxFailureCount = count
		return t
	})
}

// GetCallback 获取任务回调函数
func (t *JobRule) GetCallback() func() error {
	return syncx.WithRLockReturnValue(&t.mu, func() func() error {
		return t.callback
	})
}

// SetCallback 设置任务回调函数
func (t *JobRule) SetCallback(callback func() error) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.callback = callback
		return t
	})
}

// GetBeforeFunc 获取执行前的函数
func (t *JobRule) GetBeforeFunc() func() {
	return syncx.WithRLockReturnValue(&t.mu, func() func() {
		return t.beforeFunc
	})
}

// SetBeforeFunc 设置执行前的函数
func (t *JobRule) SetBeforeFunc(beforeFunc func()) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.beforeFunc = beforeFunc
		return t
	})
}

// GetAfterSuccessFunc 获取执行成功后的函数
func (t *JobRule) GetAfterSuccessFunc() func() {
	return syncx.WithRLockReturnValue(&t.mu, func() func() {
		return t.afterSuccessFunc
	})
}

// SetAfterSuccessFunc 设置执行成功后的函数
func (t *JobRule) SetAfterSuccessFunc(afterSuccessFunc func()) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.afterSuccessFunc = afterSuccessFunc
		return t
	})
}

// GetAfterFailureFunc 获取执行失败后的函数
func (t *JobRule) GetAfterFailureFunc() func() {
	return syncx.WithRLockReturnValue(&t.mu, func() func() {
		return t.afterFailureFunc
	})
}

// SetAfterFailureFunc 设置执行失败后的函数
func (t *JobRule) SetAfterFailureFunc(afterFailureFunc func()) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.afterFailureFunc = afterFailureFunc
		return t
	})
}

// GetSkipFunc 获取跳过执行的条件函数
func (t *JobRule) GetSkipFunc() func() bool {
	return syncx.WithRLockReturnValue(&t.mu, func() func() bool {
		return t.skipFunc
	})
}

// SetSkipFunc 设置跳过执行的条件函数
func (t *JobRule) SetSkipFunc(skipFunc func() bool) *JobRule {
	return syncx.WithLockReturnValue(&t.mu, func() *JobRule {
		t.skipFunc = skipFunc
		return t
	})
}

func (t *JobRule) Abort() {
	syncx.WithLock(&t.mu, func() {
		t.abort = true
	})
}

func (t *JobRule) Aborted() bool {
	return syncx.WithRLockReturnValue(&t.mu, func() bool {
		return t.abort
	})
}
