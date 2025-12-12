/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:25
 * @FilePath: \go-toolbox\pkg\convert\must_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	unexpectedErrorMsg = "unexpected error: %v"
	expectedResultMsg  = "expected: %v, got: %v"
)

// 基准测试
func BenchmarkNumberSliceToStringSlice(b *testing.B) {
	numbers := []uint64{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		NumberSliceToStringSlice(numbers)
	}
}

func BenchmarkStringSliceToIntSlice(b *testing.B) {
	input := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		result, err := StringSliceToNumberSlice[int](input, nil)
		if err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
		assert.Equal(b, expected, result, expectedResultMsg, expected, result)
	}
}

func BenchmarkStringSliceToFloat64SliceRoundUp(b *testing.B) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	expected := []float64{2.0, 3.0, 4.0, 4.0, 6.0}
	mode := RoundUp
	for i := 0; i < b.N; i++ {
		result, err := StringSliceToNumberSlice[float64](input, &mode)
		if err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
		assert.Equal(b, expected, result, expectedResultMsg, expected, result)
	}
}

func BenchmarkStringSliceToFloat64SliceRoundDown(b *testing.B) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	mode := RoundDown
	for i := 0; i < b.N; i++ {
		result, err := StringSliceToNumberSlice[float64](input, &mode)
		if err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
		assert.Equal(b, expected, result, expectedResultMsg, expected, result)
	}
}

func BenchmarkStringSliceToInt64Slice(b *testing.B) {
	input := []string{"1000", "2000", "3000"}
	expected := []int64{1000, 2000, 3000}
	for i := 0; i < b.N; i++ {
		result, err := StringSliceToNumberSlice[int64](input, nil)
		if err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
		assert.Equal(b, expected, result, expectedResultMsg, expected, result)
	}
}

func BenchmarkStringSliceToNumberSliceInvalidInput(b *testing.B) {
	input := []string{"a", "b", "c"}
	mode := RoundDown
	for i := 0; i < b.N; i++ {
		_, err := StringSliceToNumberSlice[int](input, &mode)
		assert.Error(b, err, "expected error, but got none")
	}
}
