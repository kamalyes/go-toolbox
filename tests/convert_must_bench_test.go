/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 10:55:57
 * @FilePath: \go-toolbox\tests\convert_must_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/stretchr/testify/assert"
)

// 基准测试
func BenchmarkNumberSliceToStringSlice(b *testing.B) {
	numbers := []uint64{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		convert.NumberSliceToStringSlice(numbers)
	}
}

func BenchmarkStringSliceToIntSlice(b *testing.B) {
	input := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		result, err := convert.StringSliceToNumberSlice[int](input, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		assert.Equal(b, expected, result, "expected: %v, got: %v", expected, result)
	}
}

func BenchmarkStringSliceToFloat64Slice_RoundUp(b *testing.B) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	expected := []float64{2.0, 3.0, 4.0, 4.0, 6.0}
	mode := convert.RoundUp
	for i := 0; i < b.N; i++ {
		result, err := convert.StringSliceToNumberSlice[float64](input, &mode)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		assert.Equal(b, expected, result, "expected: %v, got: %v", expected, result)
	}
}

func BenchmarkStringSliceToFloat64Slice_RoundDown(b *testing.B) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	mode := convert.RoundDown
	for i := 0; i < b.N; i++ {
		result, err := convert.StringSliceToNumberSlice[float64](input, &mode)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		assert.Equal(b, expected, result, "expected: %v, got: %v", expected, result)
	}
}

func BenchmarkStringSliceToInt64Slice(b *testing.B) {
	input := []string{"1000", "2000", "3000"}
	expected := []int64{1000, 2000, 3000}
	for i := 0; i < b.N; i++ {
		result, err := convert.StringSliceToNumberSlice[int64](input, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		assert.Equal(b, expected, result, "expected: %v, got: %v", expected, result)
	}
}

func BenchmarkStringSliceToNumberSlice_InvalidInput(b *testing.B) {
	input := []string{"a", "b", "c"}
	mode := convert.RoundDown
	for i := 0; i < b.N; i++ {
		_, err := convert.StringSliceToNumberSlice[int](input, &mode)
		assert.Error(b, err, "expected error, but got none")
	}
}
