/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 02:11:22
 * @FilePath: \go-toolbox\tests\retry_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/retry"
	"github.com/stretchr/testify/assert"
)

func TestDoRetry(t *testing.T) {
	// 测试Retry结构的Do方法
	t.Run("TestRetry", func(t *testing.T) {
		retryInstance := retry.NewRetry().
			SetAttemptCount(3).
			SetInterval(time.Second).
			SetErrCallback(func(nowAttemptCount, remainCount int, err error, funcName ...string) {
				fmt.Printf("当前尝试次数: %v, 剩余尝试次数: %v, 错误: %v\n", nowAttemptCount, remainCount, err)
			})

		// 函数返回错误，预期触发重试
		err := retryInstance.Do(func() error {
			fmt.Println(time.Now())
			return errors.New("发生错误")
		})

		assert.Error(t, err, "预期发生错误")
	})

	// 正常执行函数，预期无错误发生
	t.Run("NormalExecution", func(t *testing.T) {
		retryInstance := retry.NewRetry().
			SetAttemptCount(3).
			SetInterval(0) // 无间隔

		err := retryInstance.Do(func() error {
			return nil
		})

		assert.NoError(t, err, "预期无错误")
	})

	// 出错重试，验证重试次数是否符合预期
	t.Run("RetryOnError", func(t *testing.T) {
		attemptCount := 3
		var count int
		retryInstance := retry.NewRetry().
			SetAttemptCount(attemptCount).
			SetInterval(0) // 无间隔

		err := retryInstance.Do(func() error {
			count++
			if count < attemptCount {
				return errors.New("发生错误")
			}
			return nil
		})

		assert.NoError(t, err, "预期无错误")
		assert.Equal(t, attemptCount, count, "尝试次数不符合预期")
	})

	// 错误回调，验证错误回调是否被调用
	t.Run("ErrorCallback", func(t *testing.T) {
		var callbackCalled bool
		retryInstance := retry.NewRetry().
			SetAttemptCount(3).
			SetInterval(0).
			SetErrCallback(func(nowAttemptCount, remainCount int, err error, funcName ...string) {
				callbackCalled = true
			})

		err := retryInstance.Do(func() error {
			return errors.New("发生错误")
		})

		assert.Error(t, err, "预期发生错误")
		assert.True(t, callbackCalled, "预期错误回调被调用")
	})
}
