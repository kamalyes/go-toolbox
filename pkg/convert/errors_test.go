/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-09 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\convert\errors_test.go
 * @Description: 转换错误包装与泛型 JSON 转换
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonInput := []byte(`{"name":"Charlie","age":35}`)

	result, err := UnmarshalJSON[Person](jsonInput)
	assert.NoError(t, err, "JSON 反序列化失败")
	assert.Equal(t, "Charlie", result.Name)
	assert.Equal(t, 35, result.Age)
}

// 测试 MarshalJSON 泛型
func TestMarshalJSON(t *testing.T) {
	type Product struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	product := Product{
		ID:   1,
		Name: "Book",
	}

	jsonOutput, err := MarshalJSON(product)
	assert.NoError(t, err, "JSON 序列化失败")
	assert.Contains(t, string(jsonOutput), `"id":1`)
	assert.Contains(t, string(jsonOutput), `"name":"Book"`)
}
