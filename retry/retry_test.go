/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 18:16:47
 * @FilePath: \go-toolbox\retry\retry_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package retry

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoRetry(t *testing.T) {
	// 测试Retry结构的Do方法
	t.Run("TestRetry", func(t *testing.T) {
		retry := NewRetry(
			WithAttemptCount(3),       // 设置尝试次数
			WithInterval(time.Second), // 设置重试间隔时间
			WithErrCallback(func(nowAttemptCount, remainCount int, err error, funcName ...string) {
				fmt.Printf("Current attempts: %v, Residue attempts: %v, Err: %v\n", nowAttemptCount, remainCount, err)
			}),
		)

		// 函数返回错误，预期触发重试
		err := retry.Do(func() error {
			fmt.Println(time.Now())
			return errors.New("error occurred")
		})

		assert.Error(t, err, "Expected error")
	})

	// 正常执行函数，预期无错误发生
	t.Run("NormalExecution", func(t *testing.T) {
		err := DoRetry(3, 0, func() error {
			return nil
		}, nil)

		assert.NoError(t, err, "Expected no error")
	})

	// 出错重试，验证重试此数是否符合预期
	t.Run("RetryOnError", func(t *testing.T) {
		attemptCount := 3
		var count int
		err := DoRetry(attemptCount, 0, func() error {
			count++
			if count < attemptCount {
				return errors.New("error occurred")
			}
			return nil
		}, nil)

		assert.NoError(t, err, "Expected no error")
		assert.Equal(t, attemptCount, count, "Unexpected number of attempts")
	})

	// 错误回调，验证错误回调是否被调用
	t.Run("ErrorCallback", func(t *testing.T) {
		var callbackCalled bool
		err := DoRetry(3, 0, func() error {
			return errors.New("error occurred")
		}, func(nowAttemptCount, remainCount int, err error, funcName ...string) {
			callbackCalled = true
		})

		assert.Error(t, err, "Expected error")
		assert.True(t, callbackCalled, "Expected error callback to be called")
	})

	// 带间隔时间，验证重试间隔时间是否生效
	t.Run("WithInterval", func(t *testing.T) {
		attemptCount := 3
		interval := time.Millisecond * 100
		start := time.Now()
		err := DoRetry(attemptCount, interval, func() error {
			return errors.New("error occurred")
		}, nil)

		elapsed := time.Since(start)
		expectedDuration := interval * time.Duration(attemptCount-1)

		assert.Error(t, err, "Expected error")
		assert.True(t, elapsed >= expectedDuration, "Expected total duration to be at least %v", expectedDuration)
	})
}
