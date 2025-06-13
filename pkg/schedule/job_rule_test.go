/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 17:17:15
 * @FilePath: \go-toolbox\pkg\schedule\job_rule_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobRuleJobRule(t *testing.T) {
	rule := &JobRule{}

	// 测试 SetNextTime 和 GetNextTime
	nextTime := time.Now().Add(1 * time.Hour)
	rule.SetNextTime(nextTime)
	assert.Equal(t, nextTime, rule.GetNextTime(), "GetNextTime should return the same time as SetNextTime")

	// 测试 SetPrevTime 和 GetPrevTime
	prevTime := time.Now().Add(-1 * time.Hour)
	rule.SetPrevTime(prevTime)
	assert.Equal(t, prevTime, rule.GetPrevTime(), "GetPrevTime should return the same time as SetPrevTime")

	// 测试 SetCooldownDuration 和 GetCooldownDuration
	cooldownDuration := 5 * time.Minute
	rule.SetCooldownDuration(cooldownDuration)
	assert.Equal(t, cooldownDuration, rule.GetCooldownDuration(), "GetCooldownDuration should return the same duration as SetCooldownDuration")

	// 测试 SetSleepDuration 和 GetSleepDuration
	sleepDuration := 10 * time.Minute
	rule.SetSleepDuration(sleepDuration)
	assert.Equal(t, sleepDuration, rule.GetSleepDuration(), "GetSleepDuration should return the same duration as SetSleepDuration")
}

// TestSetAndGetExpression 测试 SetExpression 和 GetExpression
func TestJobRuleSetAndGetExpression(t *testing.T) {
	rule := &JobRule{}

	expression := "0 0 * * *" // 每天午夜
	rule.SetExpression(expression)
	assert.Equal(t, expression, rule.GetExpression(), "GetExpression should return the same expression as SetExpression")
}

// TestSetAndGetMaxFailureCount 测试 SetMaxFailureCount 和 GetMaxFailureCount
func TestJobRuleSetAndGetMaxFailureCount(t *testing.T) {
	rule := &JobRule{}

	maxFailureCount := 5
	rule.SetMaxFailureCount(maxFailureCount)
	assert.Equal(t, maxFailureCount, rule.GetMaxFailureCount(), "GetMaxFailureCount should return the same maxFailureCount as SetMaxFailureCount")
}

// TestSetAndGetCallback 测试 SetCallback 和 GetCallback
func TestJobRuleSetAndGetCallback(t *testing.T) {
	rule := &JobRule{}

	expectedCalled := false
	callback := func() error {
		expectedCalled = true
		return nil
	}
	rule.SetCallback(callback)

	// 调用回调函数并验证
	err := rule.GetCallback()()
	assert.NoError(t, err, "Callback should not return an error")
	assert.True(t, expectedCalled, "Callback should have been called")
}

// TestSetAndGetBeforeFunc 测试 SetBeforeFunc 和 GetBeforeFunc
func TestJobRuleSetAndGetBeforeFunc(t *testing.T) {
	rule := &JobRule{}

	expectedCalled := false
	beforeFunc := func() {
		expectedCalled = true
	}
	rule.SetBeforeFunc(beforeFunc)

	// 调用前置函数并验证
	rule.GetBeforeFunc()()
	assert.True(t, expectedCalled, "BeforeFunc should have been called")
}

// 测试 SetAfterSuccessFunc 和 GetAfterSuccessFunc
func TestJobRuleSetAndGetAfterSuccessFunc(t *testing.T) {
	rule := &JobRule{}

	called := false
	afterSuccessFunc := func() {
		called = true
	}
	rule.SetAfterSuccessFunc(afterSuccessFunc)

	f := rule.GetAfterSuccessFunc()
	f()
	assert.True(t, called, "AfterSuccessFunc should have been called")
}

// 测试 SetAfterFailureFunc 和 GetAfterFailureFunc
func TestJobRuleSetAndGetAfterFailureFunc(t *testing.T) {
	rule := &JobRule{}

	called := false
	afterFailureFunc := func() {
		called = true
	}
	rule.SetAfterFailureFunc(afterFailureFunc)

	f := rule.GetAfterFailureFunc()
	f()
	assert.True(t, called, "AfterFailureFunc should have been called")
}

// TestSetAndGetSkipFunc 测试 SetSkipFunc 和 GetSkipFunc
func TestJobRuleSetAndGetSkipFunc(t *testing.T) {
	rule := &JobRule{}

	skipFunc := func() bool {
		return true
	}
	rule.SetSkipFunc(skipFunc)

	// 调用跳过函数并验证
	result := rule.GetSkipFunc()()
	assert.True(t, result, "SkipFunc should return true")
}

func TestJobRuleTimeout(t *testing.T) {
	job := &JobRule{}

	// 默认timeout为0
	assert.Equal(t, time.Duration(0), job.GetTimeout(), "默认timeout应为0")

	// 设置超时时间为5秒
	d := 5 * time.Second
	job.SetTimeout(d)
	assert.Equal(t, d, job.GetTimeout(), "设置timeout后获取值应一致")

	// 再设置超时时间为10毫秒
	d2 := 10 * time.Millisecond
	job.SetTimeout(d2)
	assert.Equal(t, d2, job.GetTimeout(), "更新timeout后获取值应一致")
}

func TestJobRule_SetTimezone(t *testing.T) {
	jr := &JobRule{}

	jr.SetTimezone("Asia/Tokyo")
	assert.Equal(t, "Asia/Tokyo", jr.timezone.String())

	jr.SetTimezone("Invalid/Zone")
	assert.Equal(t, DefaultTimeZone, jr.timezone)
	assert.Equal(t, DefaultTimeZone, jr.GetTimezone())
}
