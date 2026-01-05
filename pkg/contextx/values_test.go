/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 10:15:55
 * @FilePath: \go-toolbox\pkg\contextx\values_test.go
 * @Description: Context 值操作测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package contextx

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetValue(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext().WithParent(parentCtx)

	key := "testKey"
	value := "testValue"
	ctx := customCtx.WithValue(key, value)

	got := ctx.Value(key)
	assert.Equal(t, value, got, "Expected value to be equal")
}

func TestValueFromParentContext(t *testing.T) {
	parentCtx := context.WithValue(context.Background(), TestParentKey, TestParentValue)
	customCtx := NewContext().WithParent(parentCtx)
	got := customCtx.Value(TestParentKey)
	assert.Equal(t, TestParentValue, got, "Expected value from parent context to be 'parentValue'")
}

func TestDeleteKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext().WithParent(parentCtx)

	key := "testKey"
	value := "testValue"
	ctx := customCtx.WithValue(key, value)

	ctx = ctx.Remove(key)

	got := ctx.Value(key)
	assert.Nil(t, got, "Expected value to be nil after deletion")
}

func TestSetNilKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext().WithParent(parentCtx)

	// WithValue 现在不返回 error，但会打印错误日志
	ctx := customCtx.WithValue(nil, "value")
	// 验证 nil key 没有被设置
	assert.Nil(t, ctx.Value(nil), "Expected nil key not to be set")
}

func TestRemoveNonExistentKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext().WithParent(parentCtx)

	ctx := customCtx.Remove(TestNonExistentKey) // should not panic
	assert.Nil(t, ctx.Value(TestNonExistentKey), "Expected nil for non-existent key")
}

func TestValues(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext().WithParent(parentCtx)

	ctx := customCtx.WithValue(TestKey1, TestValue1).WithValue(TestKey2, TestValue2)

	values := ctx.Values()
	assert.Equal(t, 2, len(values), "Expected 2 values in context")
	assert.Equal(t, TestValue1, values[TestKey1], "Expected value1 for key1")
	assert.Equal(t, TestValue2, values[TestKey2], "Expected value2 for key2")
}

func TestSetByteSlice(t *testing.T) {
	parentCtx := context.Background()
	customCtx := NewContext().WithParent(parentCtx)

	byteSlice := []byte(TestByteValue)
	ctx := customCtx.WithValue(TestByteKey, byteSlice)

	got := ctx.Value(TestByteKey)
	assert.Equal(t, byteSlice, got, "Expected byte slice to be equal")
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
	myMap[TestInt] = "integerValue"
	myMap[TestFloat64314] = "floatValue"
	myMap[true] = "booleanValue"

	// 预期的字符串输出
	expected := fmt.Sprintf("%v.WithValue(%v)", ctx, customCtx.Values())

	// 调用 String 方法
	result := customCtx.String()

	// 验证结果
	assert.Equal(t, expected, result, "Expected String output to match")
}
