/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-05 17:36:27
 * @FilePath: \go-toolbox\tests\moment_timer_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/moment"
	"github.com/stretchr/testify/assert"
)

// 公共测试数据
var ruleExpression = "*/1 * * * * ?" // 每秒执行一次

// 辅助函数：启动定时器并等待一定时间
func startTimerAndWait(t *testing.T, timer *moment.Timer, duration time.Duration) {
	err := timer.Start()
	assert.NoError(t, err, "Expected no error when starting timer")
	time.Sleep(duration)
	timer.Stop()
}

// 公共断言检查
func assertTaskExecution(t *testing.T, timer *moment.Timer, rule string, expectedCount int) {
	taskInfo, exist := timer.GetTask(rule)
	assert.True(t, exist, "Expected task to exist")
	assert.Equal(t, expectedCount, taskInfo[0].GetExecCount(), "Expected execution count to match")
}

// TestAddRule 测试添加调度规则
func TestAddRule(t *testing.T) {
	timer := moment.NewTimerWithCtx(context.Background())
	timer.AddRule(ruleExpression, func() error { return nil })

	rules := timer.GetRules()
	assert.Len(t, rules, 1, "Expected one rule to be added.")
	assert.Equal(t, ruleExpression, rules[0].GetExpression(), "Expected rule expression to match.")
}

// TestTimer 测试定时器的不同场景
func TestTimer(t *testing.T) {
	tests := []struct {
		name          string
		taskFunc      func() error
		expectedCount int
		sleepDuration time.Duration
	}{
		{
			name:          "Normal Execution",
			taskFunc:      func() error { return nil },
			expectedCount: 3,
			sleepDuration: 3 * time.Second,
		},
		{
			name:          "Execution with Error",
			taskFunc:      func() error { return fmt.Errorf("intentional error") },
			expectedCount: 3,
			sleepDuration: 3 * time.Second,
		},
		{
			name:          "Cooldown Test",
			taskFunc:      func() error { return nil },
			expectedCount: 5,
			sleepDuration: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := moment.NewTimerWithCtx(context.Background())
			timer.AddRule(ruleExpression, tt.taskFunc)

			if tt.name == "Cooldown Test" {
				timer.SetCooldownDuration(2 * time.Second)
			}

			startTimerAndWait(t, timer, tt.sleepDuration)
			assertTaskExecution(t, timer, ruleExpression, tt.expectedCount)
		})
	}
}

// TestTimerWithSkip 测试跳过规则的功能
func TestTimerWithSkip(t *testing.T) {
	timer := moment.NewTimerWithCtx(context.Background())
	executed := false

	timer.AddRule(ruleExpression, func() error {
		executed = true
		return nil
	}).Skip(func() bool { return true }) // 始终跳过执行

	startTimerAndWait(t, timer, 3*time.Second)
	assert.False(t, executed, "Expected task to be skipped.")
}

// TestTimerSetDefaultScheduleRule
func TestTimerSetDefaultScheduleRule(t *testing.T) {
	timer := moment.NewTimerWithCtx(context.Background())
	timer.SetDefaultScheduleRule(moment.EverySecond, func() error {
		return nil
	})

	startTimerAndWait(t, timer, 3*time.Second)
	rule, _ := moment.EverySecond.String()
	taskInfo, exist := timer.GetTask(rule)
	assert.True(t, exist, "Expected task to exist")
	assert.Equal(t, 3, taskInfo[0].GetExecCount(), "Expected execution count to match")
}

// TestTimerWithCustomFunc 测试自定义函数的执行
func TestTimerWithCustomFunc(t *testing.T) {
	timer := moment.NewTimerWithCtx(context.Background())
	executed := false

	timer.AddRule(ruleExpression, func() error {
		executed = true
		return nil
	}).SetCustomFunc(func() error {
		executed = true // 自定义函数也执行
		return nil
	})

	startTimerAndWait(t, timer, 3*time.Second)
	assert.True(t, executed, "Expected task to be executed.")
}

// TestTimerConcurrency 测试并发启动和停止
func TestTimerConcurrency(t *testing.T) {
	timer := moment.NewTimerWithCtx(context.Background())
	timer.AddRule(ruleExpression, func() error { return nil })

	go func() {
		err := timer.Start()
		assert.NoError(t, err, "Expected no error when starting timer")
	}()

	time.Sleep(1 * time.Second)
	go timer.Stop()
	time.Sleep(2 * time.Second) // 确保定时器有时间执行
}

// TestTimerBreak 测试熔断功能
func TestTimerBreak(t *testing.T) {
	timer := moment.NewTimerWithCtx(context.Background())
	timer.AddRule(ruleExpression, func() error {
		return fmt.Errorf("error") // 让任务失败
	})

	startTimerAndWait(t, timer, 5*time.Second)
	taskInfo, exist := timer.GetTask(ruleExpression)
	assert.True(t, exist, "Expected task to exist")
	assert.Equal(t, 5, taskInfo[0].GetExecCount(), "Expected execution count to match")
}
