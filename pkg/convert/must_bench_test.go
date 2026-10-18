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
//
// 正确性只在计时前校验一次断言不能放进 b.N 循环：testify 的 assert.Equal 走
// reflect.DeepEqual 并格式化消息，实测它自己就要 196ns、2 次分配，原先占了
// BenchmarkStringSliceToIntSlice 总耗时的 44%，量到的其实是 testify
// 另外每个基准都要 ReportAllocs——这组转换的问题本来就在分配次数上，不报出来
// 等于把要看的东西藏起来
//
// [EN] Correctness is checked once before timing. Assertions must not sit inside
// the b.N loop: testify's assert.Equal goes through reflect.DeepEqual and formats
// its message, measuring 196ns and 2 allocations on its own — 44% of what
// BenchmarkStringSliceToIntSlice used to report, so the number described testify
// rather than the code. ReportAllocs is on everywhere because allocation count is
// exactly what matters for these conversions

func BenchmarkNumberSliceToStringSlice(b *testing.B) {
	numbers := []uint64{1, 2, 3, 4, 5}
	assert.Equal(b, []string{"1", "2", "3", "4", "5"}, NumberSliceToStringSlice(numbers))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NumberSliceToStringSlice(numbers)
	}
}

func BenchmarkStringSliceToIntSlice(b *testing.B) {
	input := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	got, err := StringSliceToNumberSlice[int](input, nil)
	assert.NoError(b, err, unexpectedErrorMsg, err)
	assert.Equal(b, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, got)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := StringSliceToNumberSlice[int](input, nil); err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
	}
}

func BenchmarkStringSliceToFloat64SliceRoundUp(b *testing.B) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	mode := RoundUp
	got, err := StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(b, err, unexpectedErrorMsg, err)
	assert.Equal(b, []float64{2.0, 3.0, 4.0, 4.0, 6.0}, got)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := StringSliceToNumberSlice[float64](input, &mode); err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
	}
}

func BenchmarkStringSliceToFloat64SliceRoundDown(b *testing.B) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	mode := RoundDown
	got, err := StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(b, err, unexpectedErrorMsg, err)
	assert.Equal(b, []float64{1.0, 2.0, 3.0, 4.0, 5.0}, got)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := StringSliceToNumberSlice[float64](input, &mode); err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
	}
}

func BenchmarkStringSliceToInt64Slice(b *testing.B) {
	input := []string{"1000", "2000", "3000"}
	got, err := StringSliceToNumberSlice[int64](input, nil)
	assert.NoError(b, err, unexpectedErrorMsg, err)
	assert.Equal(b, []int64{1000, 2000, 3000}, got)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := StringSliceToNumberSlice[int64](input, nil); err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
	}
}

// BenchmarkStringSliceToNumberSliceInvalidInput 失败路径：首个元素即报错后立刻返回
func BenchmarkStringSliceToNumberSliceInvalidInput(b *testing.B) {
	input := []string{"a", "b", "c"}
	mode := RoundDown
	_, err := StringSliceToNumberSlice[int](input, &mode)
	assert.Error(b, err, "expected error, but got none")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := StringSliceToNumberSlice[int](input, &mode); err == nil {
			b.Fatal("expected error, but got none")
		}
	}
}

// BenchmarkStringSliceToFloatSlice StringSliceToFloatSlice 原先没有基准覆盖
func BenchmarkStringSliceToFloatSlice(b *testing.B) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	got, err := StringSliceToFloatSlice[float64](input, RoundNone)
	assert.NoError(b, err, unexpectedErrorMsg, err)
	assert.Len(b, got, len(input))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := StringSliceToFloatSlice[float64](input, RoundNone); err != nil {
			b.Fatalf(unexpectedErrorMsg, err)
		}
	}
}
