/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-19 10:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-19 10:00:00
 * @FilePath: \go-toolbox\pkg\desensitize\masker_test.go
 * @Description: 数据脱敏器测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataMasker_MaskJSON(t *testing.T) {
	masker := NewMasker()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "脱敏密码字段",
			input:    `{"username":"admin","password":"secret123"}`,
			expected: `{"username":"admin","password":"***"}`,
		},
		{
			name:     "脱敏 token 字段",
			input:    `{"user_id":123,"access_token":"abc123xyz"}`,
			expected: `{"user_id":123,"access_token":"***"}`,
		},
		{
			name:     "嵌套 JSON 脱敏",
			input:    `{"user":{"name":"test","password":"pass123"}}`,
			expected: `{"user":{"name":"test","password":"***"}}`,
		},
		{
			name:     "数组中的敏感数据",
			input:    `{"users":[{"name":"user1","token":"token1"},{"name":"user2","token":"token2"}]}`,
			expected: `{"users":[{"name":"user1","token":"***"},{"name":"user2","token":"***"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := masker.Mask([]byte(tt.input))
			assertJSONEqual(t, tt.expected, result)
		})
	}
}

// assertJSONEqual 比较两个 JSON 字符串是否相等（忽略字段顺序）
func assertJSONEqual(t *testing.T, expected, actual string) {
	var expectedObj, actualObj interface{}

	err := json.Unmarshal([]byte(expected), &expectedObj)
	require.NoError(t, err, "期望的 JSON 格式错误")

	err = json.Unmarshal([]byte(actual), &actualObj)
	require.NoError(t, err, "实际的 JSON 格式错误")

	assert.Equal(t, expectedObj, actualObj)
}

func TestDataMasker_MaskText(t *testing.T) {
	masker := NewMasker()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL 参数脱敏",
			input:    "password=secret123&username=admin",
			expected: "password=***&username=admin",
		},
		{
			name:     "表单数据脱敏",
			input:    "token: abc123xyz",
			expected: "token=***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := masker.Mask([]byte(tt.input))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataMasker_CustomConfig(t *testing.T) {
	config := &MaskerConfig{
		SensitiveKeys: []string{"custom_field"},
		SensitiveMask: "[REDACTED]",
		MaxBodySize:   1024,
	}
	masker := NewMasker(config)

	input := `{"custom_field":"sensitive_data","normal_field":"normal_data"}`
	expected := `{"custom_field":"[REDACTED]","normal_field":"normal_data"}`
	result := masker.Mask([]byte(input))

	assertJSONEqual(t, expected, result)
}

func TestDataMasker_MaxBodySize(t *testing.T) {
	config := &MaskerConfig{
		SensitiveKeys: []string{"password"},
		SensitiveMask: "***",
		MaxBodySize:   20,
	}
	masker := NewMasker(config)

	input := `{"password":"secret123","other":"data that exceeds max size"}`
	result := masker.Mask([]byte(input))

	assert.LessOrEqual(t, len(result), 20)
}

func TestDataMasker_EmptyData(t *testing.T) {
	masker := NewMasker()

	t.Run("空字节数组", func(t *testing.T) {
		result := masker.Mask([]byte{})
		assert.Equal(t, "", result)
	})

	t.Run("nil 数据", func(t *testing.T) {
		result := masker.Mask(nil)
		assert.Equal(t, "", result)
	})
}

func TestDataMasker_NonJSONData(t *testing.T) {
	masker := NewMasker()

	t.Run("纯文本数据", func(t *testing.T) {
		input := "password=secret123"
		expected := "password=***"
		result := masker.Mask([]byte(input))
		assert.Equal(t, expected, result)
	})

	t.Run("无效 JSON", func(t *testing.T) {
		input := `{"invalid json`
		result := masker.Mask([]byte(input))
		assert.NotEmpty(t, result)
	})
}

func BenchmarkDataMasker_MaskJSON(b *testing.B) {
	masker := NewMasker()
	data := []byte(`{"username":"admin","password":"secret123","token":"abc123xyz"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		masker.Mask(data)
	}
}

func BenchmarkDataMasker_MaskText(b *testing.B) {
	masker := NewMasker()
	data := []byte("password=secret123&token=abc123xyz&username=admin")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		masker.Mask(data)
	}
}

func TestDataMasker_ChainCalls(t *testing.T) {
	t.Run("链式调用设置配置", func(t *testing.T) {
		input := `{"api_key":"12345","secret":"abcde","normal":"data"}`
		expected := `{"api_key":"[HIDDEN]","secret":"[HIDDEN]","normal":"data"}`
		result := NewMasker().
			WithSensitiveKeys("api_key", "secret").
			WithMask("[HIDDEN]").
			WithMaxBodySize(2048).
			MaskString(input)

		assertJSONEqual(t, expected, result)
	})

	t.Run("链式调用添加敏感字段", func(t *testing.T) {
		input := `{"custom_field":"sensitive","password":"pass123"}`
		expected := `{"custom_field":"***","password":"***"}`
		result := NewMasker().
			AddSensitiveKeys("custom_field").
			MaskString(input)

		assertJSONEqual(t, expected, result)
	})

	t.Run("链式调用处理字节数组", func(t *testing.T) {
		data := []byte(`{"token":"abc123"}`)
		expected := `{"token":"***"}`
		result := NewMasker().
			WithMask("***").
			MaskBytes(data)

		assertJSONEqual(t, expected, result)
	})

	t.Run("链式调用组合使用", func(t *testing.T) {
		masker := NewMasker().
			WithSensitiveKeys("key1", "key2").
			AddSensitiveKeys("key3").
			WithMask("[REDACTED]")

		input1 := `{"key1":"val1","key2":"val2"}`
		expected1 := `{"key1":"[REDACTED]","key2":"[REDACTED]"}`
		result1 := masker.MaskString(input1)

		input2 := `{"key3":"val3"}`
		expected2 := `{"key3":"[REDACTED]"}`
		result2 := masker.MaskString(input2)

		assertJSONEqual(t, expected1, result1)
		assertJSONEqual(t, expected2, result2)
	})
}
