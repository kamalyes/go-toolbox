/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-17 15:07:01
 * @FilePath: \go-toolbox\tests\contextx_test.go
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

	"github.com/kamalyes/go-toolbox/pkg/contextx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	assert.Equal(t, parentCtx, customCtx.Context, "Expected parent context to be equal")

	// 创建一个基础上下文
	parentWithCtx := context.Background()
	customPool := syncx.NewLimitedPool(32, 1024)

	// 测试 NewContextWithValue
	customCtx, err := contextx.NewContextWithValue(parentWithCtx, "key1", "value1", customPool)
	assert.NoError(t, err, "Expected no error when creating Context with NewContextWithValue")

	// 测试 Value
	assert.Equal(t, "value1", customCtx.Value("key1"), "Expected value1 for key1")

	// 测试 NewLocalContextWithValue
	customCtx, err = contextx.NewLocalContextWithValue(customCtx, "key2", "value2")
	assert.NoError(t, err, "Expected no error when setting local value with NewLocalContextWithValue")

	// 测试 Value
	assert.Equal(t, "value2", customCtx.Value("key2"), "Expected value2 for key2")

	// 测试父上下文中的值
	assert.Equal(t, "value1", customCtx.Value("key1"), "Expected value1 from parent context")

	// 测试 DeleteKey
	customCtx.Remove("key1")
	assert.Nil(t, customCtx.Value("key1"), "Expected nil for key1 after deletion")

	// 测试 IsContext
	assert.True(t, contextx.IsContext(customCtx), "Expected customCtx to be a Context")
	assert.False(t, contextx.IsContext(parentCtx), "Expected parentCtx not to be a Context")
}

func TestSetAndGetValue(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	key := "testKey"
	value := "testValue"
	customCtx.Set(key, value)

	got := customCtx.Value(key)
	assert.Equal(t, value, got, "Expected value to be equal")
}

func TestValueFromParentContext(t *testing.T) {
	parentCtx := context.WithValue(context.Background(), "parentKey", "parentValue")
	customCtx := contextx.NewContext(parentCtx, nil)
	got := customCtx.Value("parentKey")
	assert.Equal(t, "parentValue", got, "Expected value from parent context to be 'parentValue'")
}

func TestDeleteKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	key := "testKey"
	value := "testValue"
	customCtx.Set(key, value)

	customCtx.Remove(key)

	got := customCtx.Value(key)
	assert.Nil(t, got, "Expected value to be nil after deletion")
}

// 测试 Context 的 String 方法
func TestContext_String(t *testing.T) {
	// 创建一个背景上下文
	ctx := context.Background()

	// 创建一个 Context 实例
	customCtx := &contextx.Context{Context: ctx}
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
	ctx1 := contextx.NewContext(context.Background(), nil)
	_ = ctx1.Set("key1", "value1")
	_ = ctx1.Set("key2", "value2")

	ctx2 := contextx.NewContext(context.Background(), nil)
	_ = ctx2.Set("key2", "newValue2") // 这个值会覆盖 ctx1 中的值
	_ = ctx2.Set("key3", "value3")

	ctx3 := contextx.NewContext(context.Background(), nil)
	_ = ctx3.Set("key4", "value4")

	// 合并上下文
	merged := contextx.MergeContext(ctx1, ctx2, ctx3)

	// 断言合并后的值
	assert.Equal(t, "value1", merged.Value("key1"), "期望值为 'value1'")
	assert.Equal(t, "newValue2", merged.Value("key2"), "期望值为 'newValue2'，应覆盖之前的值")
	assert.Equal(t, "value3", merged.Value("key3"), "期望值为 'value3'")
	assert.Equal(t, "value4", merged.Value("key4"), "期望值为 'value4'")
	assert.Nil(t, merged.Value("key5"), "期望值为 nil，因为 key5 不存在")
}

// TestMergeContextEmpty 测试合并空上下文
func TestMergeContextEmpty(t *testing.T) {
	merged := contextx.MergeContext()

	assert.NotNil(t, merged, "期望合并后的上下文不为 nil")
	assert.Equal(t, context.Background(), merged.Context, "期望合并后的上下文为背景上下文")
}

func TestNewContextWithTimeout(t *testing.T) {
	parentCtx := context.Background()
	timeout := 1 * time.Second
	customCtx := contextx.NewContextWithTimeout(parentCtx, timeout, nil)

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
	customCtx := contextx.NewContextWithCancel(parentCtx, nil)

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
	customCtx := contextx.NewContext(parentCtx, nil)

	err := customCtx.Set(nil, "value")
	assert.Error(t, err, "Expected error when setting nil key")
}

func TestRemoveNonExistentKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	customCtx.Remove("nonExistentKey") // should not panic or error
	assert.Nil(t, customCtx.Value("nonExistentKey"), "Expected nil for non-existent key")
}

func TestValues(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	customCtx.Set("key1", "value1")
	customCtx.Set("key2", "value2")

	values := customCtx.Values()
	assert.Equal(t, 2, len(values), "Expected 2 values in context")
	assert.Equal(t, "value1", values["key1"], "Expected value1 for key1")
	assert.Equal(t, "value2", values["key2"], "Expected value2 for key2")
}

func TestCancel(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewContextWithCancel(parentCtx, nil)
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
	customCtx := contextx.NewContextWithTimeout(parentCtx, timeout, nil)

	deadline, ok := customCtx.Deadline()
	assert.True(t, ok, "Expected deadline to be set")
	assert.WithinDuration(t, time.Now().Add(timeout), deadline, time.Second, "Expected deadline to be within duration of timeout")
}

func TestSetByteSlice(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	byteSlice := []byte("test")
	err := customCtx.Set("byteKey", byteSlice)
	assert.NoError(t, err, "Expected no error when setting byte slice")

	got := customCtx.Value("byteKey")
	assert.Equal(t, byteSlice, got, "Expected byte slice to be equal")
}
