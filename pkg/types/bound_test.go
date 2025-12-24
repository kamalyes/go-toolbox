/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-24 19:20:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-24 19:01:15
 * @FilePath: \go-scheduler\go-toolbox\pkg\types\bound_test.go
 * @Description: 边界和范围相关类型的单元测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package types

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/assert"
)

// TestBounds 测试 Bounds 结构体
func TestBounds(t *testing.T) {
	t.Run("基本整数范围", func(t *testing.T) {
		bounds := Bounds[int]{
			Min: 0,
			Max: 59,
		}
		assert.Equal(0, bounds.Min)
		assert.Equal(59, bounds.Max)
		assert.Nil(bounds.Names)
	})

	t.Run("带名称映射的范围", func(t *testing.T) {
		bounds := Bounds[int]{
			Min: 1,
			Max: 12,
			Names: map[string]int{
				"jan": 1,
				"feb": 2,
				"dec": 12,
			},
		}
		assert.Equal(1, bounds.Min)
		assert.Equal(12, bounds.Max)
		assert.NotNil(bounds.Names)
		assert.Equal(1, bounds.Names["jan"])
		assert.Equal(12, bounds.Names["dec"])
	})

	t.Run("uint 类型范围", func(t *testing.T) {
		bounds := Bounds[uint]{
			Min: 0,
			Max: 100,
		}
		assert.Equal(uint(0), bounds.Min)
		assert.Equal(uint(100), bounds.Max)
	})

	t.Run("int8 类型范围", func(t *testing.T) {
		bounds := Bounds[int8]{
			Min: -128,
			Max: 127,
		}
		assert.Equal(int8(-128), bounds.Min)
		assert.Equal(int8(127), bounds.Max)
	})
}

// TestBoundType 测试边界类型枚举
func TestBoundType(t *testing.T) {
	tests := []struct {
		name      string
		boundType BoundType
		expected  int
	}{
		{"闭区间", BoundClosed, 0},
		{"开区间", BoundOpen, 1},
		{"左开右闭", BoundLeftOpen, 2},
		{"左闭右开", BoundRightOpen, 3},
		{"左无界", BoundLeftUnbounded, 4},
		{"右无界", BoundRightUnbounded, 5},
		{"无界", BoundUnbounded, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.expected, int(tt.boundType))
		})
	}
}

// TestRangeMode 测试范围解析模式
func TestRangeMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     RangeMode
		expected int
	}{
		{"普通模式", RangeModeNormal, 0},
		{"通配符模式", RangeModeWildcard, 1},
		{"步长模式", RangeModeStep, 2},
		{"列表模式", RangeModeList, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.expected, int(tt.mode))
		})
	}
}

// TestBoundError 测试边界错误类型
func TestBoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      BoundError
		expected int
	}{
		{"无错误", BoundErrorNone, 0},
		{"低于最小值", BoundErrorBelowMin, 1},
		{"超过最大值", BoundErrorAboveMax, 2},
		{"无效范围", BoundErrorInvalidRange, 3},
		{"步长为零", BoundErrorZeroStep, 4},
		{"负数错误", BoundErrorNegative, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.expected, int(tt.err))
		})
	}
}

// TestRangeValidator 测试范围验证器函数类型
func TestRangeValidator(t *testing.T) {
	t.Run("验证器函数签名", func(t *testing.T) {
		bounds := Bounds[int]{Min: 0, Max: 100}

		// 正常值验证器
		normalValidator := func(value int, b Bounds[int]) BoundError {
			if value < b.Min {
				return BoundErrorBelowMin
			}
			if value > b.Max {
				return BoundErrorAboveMax
			}
			return BoundErrorNone
		}

		// 测试正常值
		err := normalValidator(50, bounds)
		assert.Equal(BoundErrorNone, err)

		// 测试低于最小值
		err = normalValidator(-1, bounds)
		assert.Equal(BoundErrorBelowMin, err)

		// 测试超过最大值
		err = normalValidator(101, bounds)
		assert.Equal(BoundErrorAboveMax, err)
	})
}

// TestRangeParser 测试范围解析器函数类型
func TestRangeParser(t *testing.T) {
	t.Run("解析器函数签名", func(t *testing.T) {
		bounds := Bounds[int]{
			Min: 1,
			Max: 12,
			Names: map[string]int{
				"jan": 1,
				"feb": 2,
			},
		}

		// 简单解析器
		simpleParser := func(expr string, b Bounds[int]) (int, error) {
			if val, ok := b.Names[expr]; ok {
				return val, nil
			}
			return 0, nil
		}

		// 测试名称解析
		result, err := simpleParser("jan", bounds)
		assert.Nil(err)
		assert.Equal(1, result)

		result, err = simpleParser("feb", bounds)
		assert.Nil(err)
		assert.Equal(2, result)
	})
}

// TestRangeTransformer 测试范围转换器函数类型
func TestRangeTransformer(t *testing.T) {
	t.Run("转换器函数签名", func(t *testing.T) {
		bounds := Bounds[int]{Min: 0, Max: 100}

		// 百分比转换器
		percentTransformer := func(value int, b Bounds[int]) (string, error) {
			percent := float64(value-b.Min) / float64(b.Max-b.Min) * 100
			return string(rune(int(percent))) + "%", nil
		}

		// 测试转换
		result, err := percentTransformer(50, bounds)
		assert.Nil(err)
		assert.NotEmpty(result)
	})

	t.Run("位掩码转换器", func(t *testing.T) {
		bounds := Bounds[uint]{Min: 0, Max: 63}

		// 位掩码转换器
		bitTransformer := func(value uint, b Bounds[uint]) (uint64, error) {
			return 1 << value, nil
		}

		// 测试转换
		result, err := bitTransformer(5, bounds)
		assert.Nil(err)
		assert.Equal(uint64(1<<5), result)

		result, err = bitTransformer(10, bounds)
		assert.Nil(err)
		assert.Equal(uint64(1<<10), result)
	})
}

// TestBoundsWithDifferentTypes 测试不同整数类型的 Bounds
func TestBoundsWithDifferentTypes(t *testing.T) {
	t.Run("int16 类型", func(t *testing.T) {
		bounds := Bounds[int16]{
			Min: -1000,
			Max: 1000,
		}
		assert.Equal(int16(-1000), bounds.Min)
		assert.Equal(int16(1000), bounds.Max)
	})

	t.Run("uint32 类型", func(t *testing.T) {
		bounds := Bounds[uint32]{
			Min: 0,
			Max: 4294967295,
		}
		assert.Equal(uint32(0), bounds.Min)
		assert.Equal(uint32(4294967295), bounds.Max)
	})

	t.Run("int64 类型", func(t *testing.T) {
		bounds := Bounds[int64]{
			Min: -9223372036854775808,
			Max: 9223372036854775807,
		}
		assert.Equal(int64(-9223372036854775808), bounds.Min)
		assert.Equal(int64(9223372036854775807), bounds.Max)
	})
}

// BenchmarkBoundsCreation 基准测试：创建 Bounds 结构体
func BenchmarkBoundsCreation(b *testing.B) {
	b.Run("无名称映射", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Bounds[int]{
				Min: 0,
				Max: 100,
			}
		}
	})

	b.Run("带名称映射", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Bounds[int]{
				Min: 1,
				Max: 12,
				Names: map[string]int{
					"jan": 1,
					"feb": 2,
					"mar": 3,
				},
			}
		}
	})
}
