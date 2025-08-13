/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 11:51:55
 * @FilePath: \go-toolbox\tests\json_base_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/json"
)

// BenchmarkAppendKeysToJSON 基准测试 AppendKeysToJSON 函数
func BenchmarkAppendKeysToJSON(b *testing.B) {
	originalJSON := `{"name": "Alice", "age": 30}`
	pairs := json.NewKeyValuePairs().
		Add("city", "New York").
		Add("country", "USA").
		Add("occupation", "Engineer").
		Add("hobby", "Photography")

	expected := map[string]interface{}{
		"name":       "Alice",
		"age":        float64(30),
		"city":       "New York",
		"country":    "USA",
		"occupation": "Engineer",
		"hobby":      "Photography",
	}

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		updatedJSON, err := json.AppendKeysToJSONMarshal(originalJSON, pairs)
		if err != nil {
			b.Fatalf("期望没有错误，实际错误为 %v", err)
		}

		result := parseJSON(b, updatedJSON)
		checkJSONEquality(b, result, expected)
	}
}
