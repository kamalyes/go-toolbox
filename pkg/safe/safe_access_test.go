/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-13 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 13:55:40
 * @FilePath: \go-toolbox\pkg\safe\safe_access_test.go
 * @Description: 安全访问装饰器 - 类似JavaScript的可选链操作符
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package safe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Address 结构体
type Address struct {
	City  string
	State string
}

// User 结构体
type User struct {
	Name    string
	Age     int
	Address *Address
	Alive   *bool
}

func TestSafeAccessWithNestedStructs(t *testing.T) {
	alive := true
	address := &Address{
		City:  "New York",
		State: "NY",
	}
	user := User{
		Name:    "Alice",
		Age:     30,
		Address: address,
		Alive:   &alive,
	}
	safe := Safe(user)

	// 测试正常情况
	assert.True(t, safe.Field("Name").IsValid())
	assert.Equal(t, "Alice", safe.Field("Name").String())

	assert.True(t, safe.Field("Age").IsValid())
	assert.Equal(t, 30, safe.Field("Age").Int())

	// 测试嵌套字段访问
	assert.True(t, safe.Field("Address").IsValid())
	assert.True(t, safe.Field("Address").Field("City").IsValid())
	assert.Equal(t, "New York", safe.Field("Address").Field("City").String())
	assert.True(t, safe.Field("Address").Field("State").IsValid())
	assert.Equal(t, "NY", safe.Field("Address").Field("State").String())

	// 测试布尔值
	assert.True(t, safe.Field("Alive").IsValid())
	assert.True(t, safe.Field("Alive").Bool())

	// 测试无效字段
	assert.False(t, safe.Field("NonExistent").IsValid())
	assert.False(t, safe.Field("Address").Field("NonExistent").IsValid())

	// 测试 nil 值
	safeNil := Safe(nil)
	assert.False(t, safeNil.IsValid())
	assert.Equal(t, "", safeNil.String(""))
	assert.Equal(t, 0, safeNil.Int(0))

	// 测试默认值
	assert.Equal(t, "default", safeNil.String("default"))
	assert.Equal(t, 42, safeNil.Int(42))

	// 测试 Duration
	duration := time.Duration(5 * time.Second)
	safeDuration := Safe(duration)
	assert.Equal(t, duration, safeDuration.Duration())

	// 测试 OrElse
	assert.Equal(t, "default", safeNil.OrElse("default").String())
	assert.Equal(t, 42, safeNil.OrElse(42).Int())

	// 测试 IfPresent
	var executed bool
	safe.IfPresent(func(value interface{}) {
		executed = true
	})
	assert.True(t, executed)

	// 测试 Map
	transformed := safe.Map(func(value interface{}) interface{} {
		return value.(User).Name + " is awesome"
	})
	assert.True(t, transformed.IsValid())
	assert.Equal(t, "Alice is awesome", transformed.String())

	// 测试 Filter
	isAdult := func(value interface{}) bool {
		return value.(User).Age >= 18
	}
	assert.True(t, safe.Filter(isAdult).IsValid())
	assert.False(t, safe.Filter(func(value interface{}) bool {
		return value.(User).Age < 18
	}).IsValid())
}
