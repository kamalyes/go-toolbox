/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:05:15
 * @FilePath: \go-toolbox\pkg\contextx\contextx_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package contextx

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext(parentCtx, nil)

	assert.Equal(t, parentCtx, customCtx.Context, "Expected parent context to be equal")

	// 创建一个基础上下文
	parentWithCtx := context.Background()
	customPool := syncx.NewLimitedPool(32, 1024)

	// 测试 NewContextWithValue
	customCtx, err := NewContextWithValue(parentWithCtx, "key1", "value1", customPool)
	assert.NoError(t, err, "Expected no error when creating Context with NewContextWithValue")

	// 测试 Value
	assert.Equal(t, "value1", customCtx.Value("key1"), "Expected value1 for key1")

	// 测试 NewLocalContextWithValue
	customCtx, err = NewLocalContextWithValue(customCtx, "key2", "value2")
	assert.NoError(t, err, "Expected no error when setting local value with NewLocalContextWithValue")

	// 测试 Value
	assert.Equal(t, "value2", customCtx.Value("key2"), "Expected value2 for key2")

	// 测试父上下文中的值
	assert.Equal(t, "value1", customCtx.Value("key1"), "Expected value1 from parent context")

	// 测试 DeleteKey
	customCtx.Remove("key1")
	assert.Nil(t, customCtx.Value("key1"), "Expected nil for key1 after deletion")

	// 测试 IsContext
	assert.True(t, IsContext(customCtx), "Expected customCtx to be a Context")
	assert.False(t, IsContext(parentCtx), "Expected parentCtx not to be a Context")
}

func TestSetAndGetValue(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext(parentCtx, nil)

	key := "testKey"
	value := "testValue"
	customCtx.WithValue(key, value)

	got := customCtx.Value(key)
	assert.Equal(t, value, got, "Expected value to be equal")
}

func TestValueFromParentContext(t *testing.T) {
	parentCtx := context.WithValue(context.Background(), "parentKey", "parentValue")
	customCtx := NewContext(parentCtx, nil)
	got := customCtx.Value("parentKey")
	assert.Equal(t, "parentValue", got, "Expected value from parent context to be 'parentValue'")
}

func TestDeleteKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext(parentCtx, nil)

	key := "testKey"
	value := "testValue"
	customCtx.WithValue(key, value)

	customCtx.Remove(key)

	got := customCtx.Value(key)
	assert.Nil(t, got, "Expected value to be nil after deletion")
}

// 测试 Context 的 String 方法
func TestContext_String(t *testing.T) {
	// 创建一个背景上下文
	ctx := context.Background()

	// 创建一个 Context 实例
	customCtx := &Context{Context: ctx}
	// Create a map with interface{} as key and value
	myMap := make(map[interface{}]interface{})

	// Adding different types of keys and values
	myMap["stringKey"] = "stringValue"
	myMap[42] = "integerValue"
	myMap[3.14] = "floatValue"
	myMap[true] = "booleanValue"

	// 预期的字符串输出
	expected := fmt.Sprintf("%v.WithValue(%v)", ctx, customCtx.Values())

	// 调用 String 方法
	result := customCtx.String()

	// 验证结果
	assert.Equal(t, expected, result, "Expected String output to match")
}

// TestMergeContext 测试合并多个上下文
func TestMergeContext(t *testing.T) {
	// 创建上下文并设置一些值
	ctx1 := NewContext(context.Background(), nil)
	_ = ctx1.WithValue("key1", "value1")
	_ = ctx1.WithValue("key2", "value2")

	ctx2 := NewContext(context.Background(), nil)
	_ = ctx2.WithValue("key2", "newValue2") // 这个值会覆盖 ctx1 中的值
	_ = ctx2.WithValue("key3", "value3")

	ctx3 := NewContext(context.Background(), nil)
	_ = ctx3.WithValue("key4", "value4")

	// 合并上下文
	merged := MergeContext(ctx1, ctx2, ctx3)

	// 断言合并后的值
	assert.Equal(t, "value1", merged.Value("key1"), "期望值为 'value1'")
	assert.Equal(t, "newValue2", merged.Value("key2"), "期望值为 'newValue2'，应覆盖之前的值")
	assert.Equal(t, "value3", merged.Value("key3"), "期望值为 'value3'")
	assert.Equal(t, "value4", merged.Value("key4"), "期望值为 'value4'")
	assert.Nil(t, merged.Value("key5"), "期望值为 nil，因为 key5 不存在")
}

// TestMergeContextEmpty 测试合并空上下文
func TestMergeContextEmpty(t *testing.T) {
	merged := MergeContext()

	assert.NotNil(t, merged, "期望合并后的上下文不为 nil")
	assert.Equal(t, context.Background(), merged.Context, "期望合并后的上下文为背景上下文")
}

func TestNewContextWithTimeout(t *testing.T) {
	parentCtx := context.Background()
	timeout := 1 * time.Second
	customCtx := NewContextWithTimeout(parentCtx, timeout, nil)

	// 等待超时
	time.Sleep(timeout + 100*time.Millisecond)
	select {
	case <-customCtx.Done():
	default:
		t.Error("Expected context to be done after timeout")
	}
}

func TestNewContextWithCancel(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContextWithCancel(parentCtx, nil)

	// 取消上下文
	customCtx.Cancel()
	select {
	case <-customCtx.Done():
	default:
		t.Error("Expected context to be done after cancellation")
	}
}

func TestSetNilKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext(parentCtx, nil)

	err := customCtx.WithValue(nil, "value")
	assert.Error(t, err, "Expected error when setting nil key")
}

func TestRemoveNonExistentKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext(parentCtx, nil)

	customCtx.Remove("nonExistentKey") // should not panic or error
	assert.Nil(t, customCtx.Value("nonExistentKey"), "Expected nil for non-existent key")
}

func TestValues(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext(parentCtx, nil)

	customCtx.WithValue("key1", "value1")
	customCtx.WithValue("key2", "value2")

	values := customCtx.Values()
	assert.Equal(t, 2, len(values), "Expected 2 values in context")
	assert.Equal(t, "value1", values["key1"], "Expected value1 for key1")
	assert.Equal(t, "value2", values["key2"], "Expected value2 for key2")
}

func TestCancel(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContextWithCancel(parentCtx, nil)
	customCtx.Cancel()

	select {
	case <-customCtx.Done():
		// Expected behavior
	default:
		t.Error("Expected context to be done after cancellation")
	}
}

func TestDeadline(t *testing.T) {
	parentCtx := context.Background()
	timeout := 1 * time.Second
	customCtx := NewContextWithTimeout(parentCtx, timeout, nil)

	deadline, ok := customCtx.Deadline()
	assert.True(t, ok, "Expected deadline to be set")
	assert.WithinDuration(t, time.Now().Add(timeout), deadline, time.Second, "Expected deadline to be within duration of timeout")
}

func TestSetByteSlice(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext(parentCtx, nil)

	byteSlice := []byte("test")
	err := customCtx.WithValue("byteKey", byteSlice)
	assert.NoError(t, err, "Expected no error when setting byte slice")

	got := customCtx.Value("byteKey")
	assert.Equal(t, byteSlice, got, "Expected byte slice to be equal")
}

// TestWithTimeout_Success 测试成功执行
func TestWithTimeout_Success(t *testing.T) {
	executed := false
	err := WithTimeout(1*time.Second, func(ctx context.Context) error {
		executed = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, executed)
}

// TestWithTimeout_WithError 测试函数返回错误
func TestWithTimeout_WithError(t *testing.T) {
	testErr := errors.New("test error")
	err := WithTimeout(1*time.Second, func(ctx context.Context) error {
		return testErr
	})
	assert.Equal(t, testErr, err)
}

// TestWithTimeout_TimeoutByDelay 测试通过延迟触发超时(函数不检查context)
func TestWithTimeout_TimeoutByDelay(t *testing.T) {
	start := time.Now()
	err := WithTimeout(100*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	// 应该在100ms左右超时,而不是等待200ms
	assert.Less(t, elapsed, 150*time.Millisecond)
}

// TestWithTimeout_TimeoutByContextCheck 测试函数主动检查context.Done()
func TestWithTimeout_TimeoutByContextCheck(t *testing.T) {
	err := WithTimeout(100*time.Millisecond, func(ctx context.Context) error {
		select {
		case <-time.After(200 * time.Millisecond):
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
	err := WithTimeout(1*time.Second, func(ctx context.Context) error {
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
	err := WithTimeout(100*time.Millisecond, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	assert.Equal(t, context.DeadlineExceeded, err)
}

// TestWithTimeoutValue_Success 测试成功返回值
func TestWithTimeoutValue_Success(t *testing.T) {
	result, err := WithTimeoutValue(1*time.Second, func(ctx context.Context) (int, error) {
		return 42, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

// TestWithTimeoutValue_WithError 测试返回错误
func TestWithTimeoutValue_WithError(t *testing.T) {
	testErr := errors.New("test error")
	result, err := WithTimeoutValue(1*time.Second, func(ctx context.Context) (string, error) {
		return "", testErr
	})
	assert.Equal(t, testErr, err)
	assert.Equal(t, "", result)
}

// TestWithTimeoutValue_TimeoutByDelay 测试通过延迟触发超时
func TestWithTimeoutValue_TimeoutByDelay(t *testing.T) {
	start := time.Now()
	result, err := WithTimeoutValue(100*time.Millisecond, func(ctx context.Context) (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 42, nil
	})
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, 0, result)
	assert.Less(t, elapsed, 150*time.Millisecond)
}

// TestWithTimeoutValue_TimeoutByContextCheck 测试函数主动检查context.Done()
func TestWithTimeoutValue_TimeoutByContextCheck(t *testing.T) {
	result, err := WithTimeoutValue(100*time.Millisecond, func(ctx context.Context) (string, error) {
		select {
		case <-time.After(200 * time.Millisecond):
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

	result, err := WithTimeoutValue(100*time.Millisecond, func(ctx context.Context) (CustomStruct, error) {
		time.Sleep(200 * time.Millisecond)
		return CustomStruct{Value: 42}, nil
	})

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, CustomStruct{}, result) // 应该返回零值
}

// TestWithTimeoutValue_ContextAvailable 测试context是否正常传递
func TestWithTimeoutValue_ContextAvailable(t *testing.T) {
	var receivedCtx context.Context
	result, err := WithTimeoutValue(1*time.Second, func(ctx context.Context) (bool, error) {
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
	strResult, err := WithTimeoutValue(1*time.Second, func(ctx context.Context) (string, error) {
		return "hello", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "hello", strResult)

	// slice类型
	sliceResult, err := WithTimeoutValue(1*time.Second, func(ctx context.Context) ([]int, error) {
		return []int{1, 2, 3}, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, sliceResult)

	// pointer类型
	type TestStruct struct{ Value int }
	ptrResult, err := WithTimeoutValue(1*time.Second, func(ctx context.Context) (*TestStruct, error) {
		return &TestStruct{Value: 99}, nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, ptrResult)
	assert.Equal(t, 99, ptrResult.Value)
}

// TestWithTimeoutValue_CancelPropagation 测试cancel传播
func TestWithTimeoutValue_CancelPropagation(t *testing.T) {
	result, err := WithTimeoutValue(100*time.Millisecond, func(ctx context.Context) (int, error) {
		<-ctx.Done()
		return 0, ctx.Err()
	})
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, 0, result)
}
