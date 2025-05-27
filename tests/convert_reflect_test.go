/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-05-27 18:51:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-27 18:52:03
 * @FilePath: \go-toolbox\tests\convert_reflect_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/stretchr/testify/assert"
)

type InnerSrc struct {
	IntVal    int32
	FloatVal  float64
	TimeField time.Time
}

type InnerDst struct {
	IntVal    int32
	FloatVal  float64
	TimeField string
}

type AllTypesSrc struct {
	StringVal string
	BoolVal   bool

	IntVal   int
	Int8Val  int8
	Int16Val int16
	Int32Val int32
	Int64Val int64

	UintVal   uint
	Uint8Val  uint8
	Uint16Val uint16
	Uint32Val uint32
	Uint64Val uint64

	Float32Val float32
	Float64Val float64

	PtrField *InnerSrc
	Slice    []InnerSrc
	Map      map[string]int

	Inner InnerSrc

	TimeField time.Time
}

type AllTypesDst struct {
	StringVal string
	BoolVal   bool

	IntVal   int
	Int8Val  int8
	Int16Val int16
	Int32Val int32
	Int64Val int64

	UintVal   uint
	Uint8Val  uint8
	Uint16Val uint16
	Uint32Val uint32
	Uint64Val uint64

	Float32Val float32
	Float64Val float64

	PtrField *InnerDst
	Slice    []InnerDst
	Map      map[string]int

	Inner InnerDst

	TimeField string
}

func TestTransformFields_AllTypes(t *testing.T) {
	now := time.Date(2025, 5, 27, 18, 52, 0, 0, time.UTC)

	src := AllTypesSrc{
		StringVal: "hello",
		BoolVal:   true,

		IntVal:   -123,
		Int8Val:  -8,
		Int16Val: -16,
		Int32Val: -32,
		Int64Val: -64,

		UintVal:   123,
		Uint8Val:  8,
		Uint16Val: 16,
		Uint32Val: 32,
		Uint64Val: 64,

		Float32Val: 3.14,
		Float64Val: 6.28,

		PtrField: &InnerSrc{
			IntVal:    1000,
			FloatVal:  1.618,
			TimeField: now,
		},

		Slice: []InnerSrc{
			{IntVal: 1, FloatVal: 2.2, TimeField: now},
			{IntVal: 3, FloatVal: 4.4, TimeField: now},
		},

		Map: map[string]int{
			"a": 1,
			"b": 2,
		},

		Inner: InnerSrc{
			IntVal:    555,
			FloatVal:  9.9,
			TimeField: now,
		},

		TimeField: now,
	}

	var dst AllTypesDst

	convert.TransformFields(&dst, src, &convert.TransformFieldsOptions{
		StrictTypeCheck: true,
	})

	assert.Equal(t, "hello", dst.StringVal)
	assert.Equal(t, true, dst.BoolVal)

	assert.Equal(t, int(-123), dst.IntVal)
	assert.Equal(t, int8(-8), dst.Int8Val)
	assert.Equal(t, int16(-16), dst.Int16Val)
	assert.Equal(t, int32(-32), dst.Int32Val)
	assert.Equal(t, int64(-64), dst.Int64Val)

	assert.Equal(t, uint(123), dst.UintVal)
	assert.Equal(t, uint8(8), dst.Uint8Val)
	assert.Equal(t, uint16(16), dst.Uint16Val)
	assert.Equal(t, uint32(32), dst.Uint32Val)
	assert.Equal(t, uint64(64), dst.Uint64Val)

	assert.InDelta(t, float32(3.14), dst.Float32Val, 0.0001)
	assert.InDelta(t, float64(6.28), dst.Float64Val, 0.0001)

	// 指针字段
	assert.NotNil(t, dst.PtrField)
	assert.Equal(t, int32(1000), dst.PtrField.IntVal)
	assert.InDelta(t, 1.618, dst.PtrField.FloatVal, 0.0001)
	assert.Equal(t, now.Format(time.DateTime), dst.PtrField.TimeField)

	// 切片字段
	assert.Len(t, dst.Slice, 2)
	assert.Equal(t, int32(1), dst.Slice[0].IntVal)
	assert.InDelta(t, 2.2, dst.Slice[0].FloatVal, 0.0001)
	assert.Equal(t, now.Format(time.DateTime), dst.Slice[0].TimeField)

	assert.Equal(t, int32(3), dst.Slice[1].IntVal)
	assert.InDelta(t, 4.4, dst.Slice[1].FloatVal, 0.0001)
	assert.Equal(t, now.Format(time.DateTime), dst.Slice[1].TimeField)

	// map字段
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, dst.Map)

	// 内嵌结构体
	assert.Equal(t, int32(555), dst.Inner.IntVal)
	assert.InDelta(t, 9.9, dst.Inner.FloatVal, 0.0001)
	assert.Equal(t, now.Format(time.DateTime), dst.Inner.TimeField)

	// 顶层time.Time转string
	assert.Equal(t, now.Format(time.DateTime), dst.TimeField)
}

func TestTransformFields_StrictTypeCheck_Panic(t *testing.T) {
	src := struct {
		Field int
	}{
		Field: 10,
	}
	dst := struct {
		Field string
	}{}

	assert.Panics(t, func() {
		convert.TransformFields(&dst, src, &convert.TransformFieldsOptions{
			StrictTypeCheck: true,
		})
	})
}

func TestTransformFields_NonStrict_NoPanic(t *testing.T) {
	src := struct {
		Field int
	}{
		Field: 10,
	}
	dst := struct {
		Field string
	}{}

	convert.TransformFields(&dst, src, &convert.TransformFieldsOptions{
		StrictTypeCheck: false,
	})

	assert.Equal(t, "", dst.Field)
}
