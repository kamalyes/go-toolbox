/*
 * @Author: kamalyes 501893067@qq.com
 * @Date:2024-12-18 22:53:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-18 22:50:29
 * @FilePath: \go-toolbox\tests\desensitize_adapter_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/desensitize"
	"github.com/stretchr/testify/assert"
)

// 定义一个测试结构体
type TestJsonDesensitizationStruct struct {
	Email       string `desensitize:"email"`
	PhoneNumber string `desensitize:"phoneNumber"`
	Name        string `desensitize:"name"`
	Age         int
}

// 自定义脱敏器示例
type MyCustomDesensitizer struct{}

func (d *MyCustomDesensitizer) Desensitize(value string) string {
	// 自定义脱敏逻辑
	return "*****" // 示例：返回固定的脱敏值
}

func init() {
	// 注册自定义脱敏器
	desensitize.RegisterDesensitizer("myCustom", &MyCustomDesensitizer{})
}

func TestDesensitization(t *testing.T) {
	// 创建一个测试对象
	testObj := &TestJsonDesensitizationStruct{
		Email:       "test@example.com",
		PhoneNumber: "1234567890",
		Name:        "张三",
		Age:         30,
	}

	// 执行脱敏操作
	err := desensitize.Desensitization(testObj)

	// 断言没有错误
	assert.NoError(t, err)

	// 断言脱敏后的结果
	assert.Equal(t, "t***@example.com", testObj.Email)
	assert.Equal(t, "123****7890", testObj.PhoneNumber)
	assert.Equal(t, "张*", testObj.Name)
	assert.Equal(t, 30, testObj.Age) // 确保年龄没有被脱敏
}

func TestDesensitization_NonStruct(t *testing.T) {
	// 测试非结构体输入
	err := desensitize.Desensitization("not a struct")
	assert.Error(t, err)
	assert.Equal(t, "expected a struct or pointer to a struct", err.Error())
}

type TestEmptyDesensitizerStruct struct {
	A string `desensitize:"abc"`
}

func TestDesensitization_EmptyDesensitizer(t *testing.T) {
	// 测试未注册脱敏器的情况
	testObj := &TestEmptyDesensitizerStruct{
		A: "10",
	}

	err := desensitize.Desensitization(testObj)
	assert.NoError(t, err)           // 不会报错
	assert.Equal(t, "10", testObj.A) // 应保持原值
}

func TestDesensitization_CustomDesensitizer(t *testing.T) {
	// 定义一个测试结构体，使用自定义脱敏器
	type TestCustomDesensitizationStruct struct {
		SecretField string `desensitize:"myCustom"`
	}

	// 创建一个测试对象
	testObj := &TestCustomDesensitizationStruct{
		SecretField: "SensitiveData",
	}

	// 执行脱敏操作
	err := desensitize.Desensitization(testObj)

	// 断言没有错误
	assert.NoError(t, err)

	// 断言脱敏后的结果
	assert.Equal(t, "*****", testObj.SecretField) // 应该被自定义脱敏器处理
}
