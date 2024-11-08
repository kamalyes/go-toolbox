/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 21:08:05
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

	"github.com/kamalyes/go-toolbox/pkg/contextx"
	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/stretchr/testify/assert"
)

func TestNewCustomContext(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewCustomContext(parentCtx, nil)

	if customCtx.Context != parentCtx {
		t.Errorf("Expected parent context to be %v, got %v", parentCtx, customCtx.Context)
	}
	// 创建一个基础上下文
	parentWithCtx := context.Background()
	customPool := osx.NewLimitedPool(32, 1024)

	// 测试 NewContextWithValue
	customCtx, err := contextx.NewContextWithValue(parentWithCtx, "key1", "value1", customPool)
	assert.NoError(t, err, "Expected no error when creating CustomContext with NewContextWithValue")

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

	// 测试 IsCustomContext
	assert.True(t, contextx.IsCustomContext(customCtx), "Expected customCtx to be a CustomContext")
	assert.False(t, contextx.IsCustomContext(parentCtx), "Expected parentCtx not to be a CustomContext")
}

func TestSetAndGetValue(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewCustomContext(parentCtx, nil)

	key := "testKey"
	value := "testValue"
	customCtx.Set(key, value)

	got := customCtx.Value(key)
	if got != value {
		t.Errorf("Expected value to be %v, got %v", value, got)
	}
}

func TestValueFromParentContext(t *testing.T) {
	parentCtx := context.WithValue(context.Background(), "parentKey", "parentValue")
	customCtx := contextx.NewCustomContext(parentCtx, nil)
	got := customCtx.Value("parentKey")
	if got != "parentValue" {
		t.Errorf("Expected value from parent context to be 'parentValue', got %v", got)
	}
}

func TestDeleteKey(t *testing.T) {
	parentCtx := context.Background()
	customCtx := contextx.NewCustomContext(parentCtx, nil)

	key := "testKey"
	value := "testValue"
	customCtx.Set(key, value)

	customCtx.Remove(key)

	got := customCtx.Value(key)
	if got != nil {
		t.Errorf("Expected value to be nil after deletion, got %v", got)
	}
}

// 测试 CustomContext 的 String 方法
func TestCustomContext_String(t *testing.T) {
	// 创建一个背景上下文
	ctx := context.Background()

	// 创建一个 CustomContext 实例
	customCtx := &contextx.CustomContext{Context: ctx}
	// Create a map with interface{} as key and value
	myMap := make(map[interface{}]interface{})

	// Adding different types of keys and values
	myMap["stringKey"] = "stringValue"
	myMap[42] = "integerValue"
	myMap[3.14] = "floatValue"
	myMap[true] = "booleanValue"
	customCtx.Tags = myMap

	// 预期的字符串输出
	expected := fmt.Sprintf("%v.WithValue(%v)", ctx, customCtx.Tags)

	// 调用 String 方法
	result := customCtx.String()

	// 验证结果
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
