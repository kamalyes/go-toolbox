/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-10 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-10 00:00:00
 * @FilePath: \go-toolbox\pkg\types\kind_test.go
 * @Description: kind.go 的单元测试，覆盖 ToFloat64OK 类型转换函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestToFloat64OK 测试将各种数值类型转换为 float64
func TestToFloat64OK(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		v, ok := ToFloat64OK(float64(3.14))
		assert.True(t, ok)
		assert.Equal(t, 3.14, v)
	})

	t.Run("float32", func(t *testing.T) {
		v, ok := ToFloat64OK(float32(2.5))
		assert.True(t, ok)
		assert.Equal(t, float64(2.5), v)
	})

	t.Run("int", func(t *testing.T) {
		v, ok := ToFloat64OK(int(42))
		assert.True(t, ok)
		assert.Equal(t, float64(42), v)
	})

	t.Run("int8", func(t *testing.T) {
		v, ok := ToFloat64OK(int8(-128))
		assert.True(t, ok)
		assert.Equal(t, float64(-128), v)
	})

	t.Run("int16", func(t *testing.T) {
		v, ok := ToFloat64OK(int16(1000))
		assert.True(t, ok)
		assert.Equal(t, float64(1000), v)
	})

	t.Run("int32", func(t *testing.T) {
		v, ok := ToFloat64OK(int32(100000))
		assert.True(t, ok)
		assert.Equal(t, float64(100000), v)
	})

	t.Run("int64", func(t *testing.T) {
		v, ok := ToFloat64OK(int64(9999999999))
		assert.True(t, ok)
		assert.Equal(t, float64(9999999999), v)
	})

	t.Run("uint", func(t *testing.T) {
		v, ok := ToFloat64OK(uint(42))
		assert.True(t, ok)
		assert.Equal(t, float64(42), v)
	})

	t.Run("uint8", func(t *testing.T) {
		v, ok := ToFloat64OK(uint8(255))
		assert.True(t, ok)
		assert.Equal(t, float64(255), v)
	})

	t.Run("uint16", func(t *testing.T) {
		v, ok := ToFloat64OK(uint16(65535))
		assert.True(t, ok)
		assert.Equal(t, float64(65535), v)
	})

	t.Run("uint32", func(t *testing.T) {
		v, ok := ToFloat64OK(uint32(4294967295))
		assert.True(t, ok)
		assert.Equal(t, float64(4294967295), v)
	})

	t.Run("uint64", func(t *testing.T) {
		v, ok := ToFloat64OK(uint64(18446744073709551615))
		assert.True(t, ok)
		assert.Equal(t, float64(18446744073709551615), v)
	})

	t.Run("不支持类型返回 false", func(t *testing.T) {
		// 字符串类型不支持
		v, ok := ToFloat64OK("hello")
		assert.False(t, ok)
		assert.Equal(t, float64(0), v)

		// bool 类型不支持
		v, ok = ToFloat64OK(true)
		assert.False(t, ok)
		assert.Equal(t, float64(0), v)

		// nil 不支持
		v, ok = ToFloat64OK(nil)
		assert.False(t, ok)
		assert.Equal(t, float64(0), v)

		// 切片不支持
		v, ok = ToFloat64OK([]int{1, 2, 3})
		assert.False(t, ok)
		assert.Equal(t, float64(0), v)

		// map 不支持
		v, ok = ToFloat64OK(map[string]int{"a": 1})
		assert.False(t, ok)
		assert.Equal(t, float64(0), v)
	})

	t.Run("特殊浮点数值", func(t *testing.T) {
		// 正无穷
		v, ok := ToFloat64OK(float64(math.Inf(1)))
		assert.True(t, ok)
		assert.True(t, math.IsInf(v, 1))

		// 负无穷
		v, ok = ToFloat64OK(float64(math.Inf(-1)))
		assert.True(t, ok)
		assert.True(t, math.IsInf(v, -1))

		// NaN
		v, ok = ToFloat64OK(float64(math.NaN()))
		assert.True(t, ok)
		assert.True(t, math.IsNaN(v))
	})

	t.Run("零值", func(t *testing.T) {
		v, ok := ToFloat64OK(int(0))
		assert.True(t, ok)
		assert.Equal(t, float64(0), v)

		v, ok = ToFloat64OK(float64(0))
		assert.True(t, ok)
		assert.Equal(t, float64(0), v)
	})
}
