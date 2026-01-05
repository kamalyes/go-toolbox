/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 10:08:00
 * @FilePath: \go-toolbox\pkg\contextx\helpers_test.go
 * @Description: Context 全局辅助函数测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package contextx

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestWithTimeout_Success 测试成功执行
func TestWithTimeout_Success(t *testing.T) {
	executed := false
	err := WithTimeout(TestTimeout1s, func(ctx context.Context) error {
		executed = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, executed)
}

// TestWithTimeout_WithError 测试函数返回错误
func TestWithTimeout_WithError(t *testing.T) {
	testErr := errors.New("test error")
	err := WithTimeout(TestTimeout1s, func(ctx context.Context) error {
		return testErr
	})
	assert.Equal(t, testErr, err)
}

// TestWithTimeout_TimeoutByDelay 测试通过延迟触发超时(函数不检查context)
func TestWithTimeout_TimeoutByDelay(t *testing.T) {
	start := time.Now()
	err := WithTimeout(TestTimeout100ms, func(ctx context.Context) error {
		time.Sleep(TestTimeout200ms)
		return nil
	})
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	// 应该在100ms左右超时,而不是等待200ms
	assert.Less(t, elapsed, TestTimeoutMargin)
}

// TestWithTimeout_TimeoutByContextCheck 测试函数主动检查context.Done()
func TestWithTimeout_TimeoutByContextCheck(t *testing.T) {
	err := WithTimeout(TestTimeout100ms, func(ctx context.Context) error {
		select {
		case <-time.After(TestTimeout200ms):
			return errors.New("should not reach here")
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

// TestWithTimeout_ContextAvailable 测试context是否正常传递
func TestWithTimeout_ContextAvailable(t *testing.T) {
	var receivedCtx context.Context
	err := WithTimeout(TestTimeout1s, func(ctx context.Context) error {
		receivedCtx = ctx
		return nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, receivedCtx)

	// 验证context有deadline
	_, ok := receivedCtx.Deadline()
	assert.True(t, ok)
}

// TestWithTimeout_CancelPropagation 测试cancel传播
func TestWithTimeout_CancelPropagation(t *testing.T) {
	err := WithTimeout(TestTimeout100ms, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	assert.Equal(t, context.DeadlineExceeded, err)
}

// TestWithTimeoutValue_Success 测试成功返回值
func TestWithTimeoutValue_Success(t *testing.T) {
	result, err := WithTimeoutValue(TestTimeout1s, func(ctx context.Context) (int, error) {
		return TestInt, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, TestInt, result)
}

// TestWithTimeoutValue_WithError 测试返回错误
func TestWithTimeoutValue_WithError(t *testing.T) {
	testErr := errors.New("test error")
	result, err := WithTimeoutValue(TestTimeout1s, func(ctx context.Context) (string, error) {
		return "", testErr
	})
	assert.Equal(t, testErr, err)
	assert.Equal(t, "", result)
}

// TestWithTimeoutValue_TimeoutByDelay 测试通过延迟触发超时
func TestWithTimeoutValue_TimeoutByDelay(t *testing.T) {
	start := time.Now()
	result, err := WithTimeoutValue(TestTimeout100ms, func(ctx context.Context) (int, error) {
		time.Sleep(TestTimeout200ms)
		return TestInt, nil
	})
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, 0, result)
	assert.Less(t, elapsed, TestTimeoutMargin)
}

// TestWithTimeoutValue_TimeoutByContextCheck 测试函数主动检查context.Done()
func TestWithTimeoutValue_TimeoutByContextCheck(t *testing.T) {
	result, err := WithTimeoutValue(TestTimeout100ms, func(ctx context.Context) (string, error) {
		select {
		case <-time.After(TestTimeout200ms):
			return "should not reach here", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	})
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, "", result)
}

// TestWithTimeoutValue_ZeroValue 测试超时时返回零值
func TestWithTimeoutValue_ZeroValue(t *testing.T) {
	type CustomStruct struct {
		Value int
	}

	result, err := WithTimeoutValue(TestTimeout100ms, func(ctx context.Context) (CustomStruct, error) {
		time.Sleep(TestTimeout200ms)
		return CustomStruct{Value: TestInt}, nil
	})

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, CustomStruct{}, result) // 应该返回零值
}

// TestWithTimeoutValue_ContextAvailable 测试context是否正常传递
func TestWithTimeoutValue_ContextAvailable(t *testing.T) {
	var receivedCtx context.Context
	result, err := WithTimeoutValue(TestTimeout1s, func(ctx context.Context) (bool, error) {
		receivedCtx = ctx
		return true, nil
	})

	assert.NoError(t, err)
	assert.True(t, result)
	assert.NotNil(t, receivedCtx)

	// 验证context有deadline
	_, ok := receivedCtx.Deadline()
	assert.True(t, ok)
}

// TestWithTimeoutValue_MultipleTypes 测试不同类型的返回值
func TestWithTimeoutValue_MultipleTypes(t *testing.T) {
	// string类型
	strResult, err := WithTimeoutValue(TestTimeout1s, func(ctx context.Context) (string, error) {
		return "hello", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "hello", strResult)

	// slice类型
	sliceResult, err := WithTimeoutValue(TestTimeout1s, func(ctx context.Context) ([]int, error) {
		return []int{1, 2, 3}, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, sliceResult)

	// pointer类型
	type TestStruct struct{ Value int }
	ptrResult, err := WithTimeoutValue(TestTimeout1s, func(ctx context.Context) (*TestStruct, error) {
		return &TestStruct{Value: TestInt99}, nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, ptrResult)
	assert.Equal(t, TestInt99, ptrResult.Value)
}

// TestWithTimeoutValue_CancelPropagation 测试cancel传播
func TestWithTimeoutValue_CancelPropagation(t *testing.T) {
	result, err := WithTimeoutValue(TestTimeout100ms, func(ctx context.Context) (int, error) {
		<-ctx.Done()
		return 0, ctx.Err()
	})
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, 0, result)
}

// TestOrBackground_WithCancelledContext 测试取消的context返回Background
func TestOrBackground_WithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := OrBackground(ctx)
	assert.Equal(t, context.Background(), result)
}

// TestOrBackground_WithValidContext 测试有效的context返回原context
func TestOrBackground_WithValidContext(t *testing.T) {
	ctx := context.Background()
	result := OrBackground(ctx)
	assert.Equal(t, ctx, result)
}

// TestWithTimeoutFrom_Success 测试从父context创建超时context
func TestWithTimeoutFrom_Success(t *testing.T) {
	parent := context.Background()
	executed := false
	err := WithTimeoutFrom(parent, TestTimeout1s, func(ctx context.Context) error {
		executed = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, executed)
}

// TestWithTimeoutFrom_ParentCancelled 测试父context取消时的行为
func TestWithTimeoutFrom_ParentCancelled(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()

	err := WithTimeoutFrom(parent, 1*time.Second, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

// TestWithTimeoutOrBackground_WithCancelledParent 测试父context取消时使用Background
func TestWithTimeoutOrBackground_WithCancelledParent(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()

	executed := false
	err := WithTimeoutOrBackground(parent, 1*time.Second, func(ctx context.Context) error {
		executed = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, executed)
}

// TestWithTimeoutOrBackground_WithValidParent 测试有效父context时的行为
func TestWithTimeoutOrBackground_WithValidParent(t *testing.T) {
	parent := context.Background()
	executed := false
	err := WithTimeoutOrBackground(parent, 1*time.Second, func(ctx context.Context) error {
		executed = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, executed)
}
