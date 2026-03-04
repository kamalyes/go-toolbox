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

// TestWithTimeoutDecorators_Success 测试成功创建带超时的context
func TestWithTimeoutDecorators_Success(t *testing.T) {
	ctx, cancel := WithTimeoutDecorators(TestTimeout1s)
	defer cancel()

	assert.NotNil(t, ctx)
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, time.Until(deadline) > 0)
}

// TestWithTimeoutDecorators_Decorator 测试带装饰器的context创建
func TestWithTimeoutDecorators_Decorator(t *testing.T) {
	type contextKey string
	const userKey contextKey = "user"

	decorator := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, userKey, "test_user")
	}

	ctx, cancel := WithTimeoutDecorators(TestTimeout1s, decorator)
	defer cancel()

	assert.NotNil(t, ctx)
	assert.Equal(t, "test_user", ctx.Value(userKey))

	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, time.Until(deadline) > 0)
}

// TestWithTimeoutDecorators_MultipleDecorators 测试多个装饰器
func TestWithTimeoutDecorators_MultipleDecorators(t *testing.T) {
	type contextKey string
	const (
		userKey  contextKey = "user"
		traceKey contextKey = "trace_id"
	)

	decorator1 := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, userKey, "test_user")
	}
	decorator2 := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, traceKey, "trace_123")
	}

	ctx, cancel := WithTimeoutDecorators(TestTimeout1s, decorator1, decorator2)
	defer cancel()

	assert.NotNil(t, ctx)
	assert.Equal(t, "test_user", ctx.Value(userKey))
	assert.Equal(t, "trace_123", ctx.Value(traceKey))
}

// TestWithTimeoutDecorators_Timeout 测试超时触发
func TestWithTimeoutDecorators_Timeout(t *testing.T) {
	ctx, cancel := WithTimeoutDecorators(TestTimeout100ms)
	defer cancel()

	select {
	case <-time.After(TestTimeout200ms):
		t.Fatal("context should have timed out")
	case <-ctx.Done():
		assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	}
}

// TestWithDeadlineDecorators_Success 测试成功创建带截止时间的context
func TestWithDeadlineDecorators_Success(t *testing.T) {
	deadline := time.Now().Add(TestTimeout1s)
	ctx, cancel := WithDeadlineDecorators(deadline)
	defer cancel()

	assert.NotNil(t, ctx)
	ctxDeadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, ctxDeadline.Equal(deadline) || ctxDeadline.Before(deadline.Add(time.Millisecond)))
}

// TestWithDeadlineDecorators_Decorator 测试带装饰器的deadline context
func TestWithDeadlineDecorators_Decorator(t *testing.T) {
	type contextKey string
	const requestIDKey contextKey = "request_id"

	decorator := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestIDKey, "req_456")
	}

	deadline := time.Now().Add(TestTimeout1s)
	ctx, cancel := WithDeadlineDecorators(deadline, decorator)
	defer cancel()

	assert.NotNil(t, ctx)
	assert.Equal(t, "req_456", ctx.Value(requestIDKey))

	ctxDeadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, time.Until(ctxDeadline) > 0)
}

// TestWithDeadlineDecorators_Timeout 测试截止时间触发
func TestWithDeadlineDecorators_Timeout(t *testing.T) {
	deadline := time.Now().Add(TestTimeout100ms)
	ctx, cancel := WithDeadlineDecorators(deadline)
	defer cancel()

	select {
	case <-time.After(TestTimeout200ms):
		t.Fatal("context should have reached deadline")
	case <-ctx.Done():
		assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	}
}

// TestMustGet_Success 测试成功获取值
func TestMustGet_Success(t *testing.T) {
	type contextKey string
	const configKey contextKey = "config"

	type Config struct {
		Value string
	}

	ctx := context.WithValue(context.Background(), configKey, &Config{Value: "test"})
	config := MustGet[*Config](ctx, configKey)

	assert.NotNil(t, config)
	assert.Equal(t, "test", config.Value)
}

// TestMustGet_Panic_NotFound 测试值不存在时panic
func TestMustGet_Panic_NotFound(t *testing.T) {
	type contextKey string
	const missingKey contextKey = "missing"

	ctx := context.Background()

	assert.Panics(t, func() {
		MustGet[string](ctx, missingKey)
	})
}

// TestMustGet_Panic_TypeMismatch 测试类型不匹配时panic
func TestMustGet_Panic_TypeMismatch(t *testing.T) {
	type contextKey string
	const valueKey contextKey = "value"

	ctx := context.WithValue(context.Background(), valueKey, "string_value")

	assert.Panics(t, func() {
		MustGet[int](ctx, valueKey)
	})
}

// TestMustGet_DifferentTypes 测试不同类型的值获取
func TestMustGet_DifferentTypes(t *testing.T) {
	type contextKey string

	// string类型
	ctx := context.WithValue(context.Background(), contextKey("str"), "hello")
	strVal := MustGet[string](ctx, contextKey("str"))
	assert.Equal(t, "hello", strVal)

	// int类型
	ctx = context.WithValue(ctx, contextKey("num"), TestInt99)
	intVal := MustGet[int](ctx, contextKey("num"))
	assert.Equal(t, TestInt99, intVal)

	// slice类型
	ctx = context.WithValue(ctx, contextKey("slice"), []int{1, 2, 3})
	sliceVal := MustGet[[]int](ctx, contextKey("slice"))
	assert.Equal(t, []int{1, 2, 3}, sliceVal)
}

// TestMustGetWithMessage_Success 测试成功获取值
func TestMustGetWithMessage_Success(t *testing.T) {
	type contextKey string
	const userKey contextKey = "user"

	ctx := context.WithValue(context.Background(), userKey, "john")
	user := MustGetWithMessage[string](ctx, userKey, "用户信息不存在")

	assert.Equal(t, "john", user)
}

// TestMustGetWithMessage_Panic_CustomMessage 测试自定义panic消息
func TestMustGetWithMessage_Panic_CustomMessage(t *testing.T) {
	type contextKey string
	const missingKey contextKey = "missing"

	ctx := context.Background()
	customMsg := "配置信息不存在于上下文中"

	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, customMsg, r)
	}()

	MustGetWithMessage[string](ctx, missingKey, customMsg)
}

// TestMustGetWithMessage_Panic_TypeMismatch 测试类型不匹配时的自定义消息
func TestMustGetWithMessage_Panic_TypeMismatch(t *testing.T) {
	type contextKey string
	const valueKey contextKey = "value"

	ctx := context.WithValue(context.Background(), valueKey, "string_value")
	customMsg := "环境信息不存在于上下文中"

	defer func() {
		r := recover()
		assert.NotNil(t, r)
		panicMsg := r.(string)
		assert.Contains(t, panicMsg, customMsg)
		assert.Contains(t, panicMsg, "type mismatch")
	}()

	MustGetWithMessage[int](ctx, valueKey, customMsg)
}

// TestGetOrDefault_ValueExists 测试值存在时返回值
func TestGetOrDefault_ValueExists(t *testing.T) {
	type contextKey string
	const timeoutKey contextKey = "timeout"

	ctx := context.WithValue(context.Background(), timeoutKey, 30*time.Second)
	timeout := GetOrDefault(ctx, timeoutKey, 10*time.Second)

	assert.Equal(t, 30*time.Second, timeout)
}

// TestGetOrDefault_ValueNotExists 测试值不存在时返回默认值
func TestGetOrDefault_ValueNotExists(t *testing.T) {
	type contextKey string
	const missingKey contextKey = "missing"

	ctx := context.Background()
	defaultValue := "anonymous"
	userID := GetOrDefault(ctx, missingKey, defaultValue)

	assert.Equal(t, defaultValue, userID)
}

// TestGetOrDefault_TypeMismatch 测试类型不匹配时返回默认值
func TestGetOrDefault_TypeMismatch(t *testing.T) {
	type contextKey string
	const valueKey contextKey = "value"

	ctx := context.WithValue(context.Background(), valueKey, "string_value")
	defaultValue := TestInt99
	intVal := GetOrDefault(ctx, valueKey, defaultValue)

	assert.Equal(t, defaultValue, intVal)
}

// TestGetOrDefault_DifferentTypes 测试不同类型的默认值
func TestGetOrDefault_DifferentTypes(t *testing.T) {
	type contextKey string
	ctx := context.Background()

	// string类型
	strVal := GetOrDefault(ctx, contextKey("str"), "default")
	assert.Equal(t, "default", strVal)

	// int类型
	intVal := GetOrDefault(ctx, contextKey("num"), TestInt)
	assert.Equal(t, TestInt, intVal)

	// bool类型
	boolVal := GetOrDefault(ctx, contextKey("flag"), true)
	assert.True(t, boolVal)

	// duration类型
	durationVal := GetOrDefault(ctx, contextKey("timeout"), TestTimeout1s)
	assert.Equal(t, TestTimeout1s, durationVal)
}

// TestGetOrDefault_NilValue 测试nil值时返回默认值
func TestGetOrDefault_NilValue(t *testing.T) {
	type contextKey string
	const nilKey contextKey = "nil_value"

	ctx := context.WithValue(context.Background(), nilKey, nil)
	defaultValue := "default"
	result := GetOrDefault(ctx, nilKey, defaultValue)

	assert.Equal(t, defaultValue, result)
}
