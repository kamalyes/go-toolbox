/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-05-27 18:51:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\convert\reflect_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"testing"
	"time"

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

func TestTransformFieldsAllTypes(t *testing.T) {
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

	err := TransformFields(&dst, src, &TransformFieldsOptions{
		StrictTypeCheck: true,
	})
	assert.NoError(t, err)

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

func TestTransformValueNilSourceCases(t *testing.T) {
	type Inner struct {
		Field int
	}

	type Src struct {
		Ptr       *Inner
		Slice     []int
		Map       map[string]int
		Interface interface{}
		Func      func()
		Chan      chan int

		StrPtr  *string
		IntPtr  *int
		BoolPtr *bool
		TimePtr *time.Time

		// 非指针基本类型，无法为 nil，测试零值赋值
		Str  string
		Int  int
		Bool bool
		Time time.Time
	}

	type Dst struct {
		Ptr       *Inner
		Slice     []int
		Map       map[string]int
		Interface interface{}
		Func      func()
		Chan      chan int

		StrPtr  string
		IntPtr  int
		BoolPtr bool
		TimePtr string

		Str  string
		Int  int
		Bool bool
		Time string
	}

	var src Src // 所有指针、slice、map、interface、func、chan 字段默认 nil，非指针基本类型为零值
	var dst Dst

	err := TransformFields(&dst, &src, nil)
	assert.NoError(t, err)

	// 指针字段 nil -> 目标为 nil 或零值
	assert.Nil(t, dst.Ptr)
	assert.Nil(t, dst.Slice)
	assert.Nil(t, dst.Map)
	assert.Nil(t, dst.Interface)
	assert.Nil(t, dst.Func)
	assert.Nil(t, dst.Chan)

	// 指向基本类型的指针为 nil，目标赋零值
	assert.Equal(t, "", dst.StrPtr)
	assert.Equal(t, 0, dst.IntPtr)
	assert.Equal(t, false, dst.BoolPtr)
	assert.Equal(t, "", dst.TimePtr) // time.Time 转 string，nil 源应赋空字符串

	// 非指针基本类型，源为零值，目标赋零值
	assert.Equal(t, "", dst.Str)
	assert.Equal(t, 0, dst.Int)
	assert.Equal(t, false, dst.Bool)
	assert.Equal(t, "0001-01-01 00:00:00", dst.Time) // time.Time 转 string，零值时间格式化为空字符串或默认格式

}

func TestTransformFieldsStrictTypeCheckError(t *testing.T) {
	src := struct {
		Field int
	}{
		Field: 10,
	}
	dst := struct {
		Field string
	}{}

	err := TransformFields(&dst, src, &TransformFieldsOptions{
		StrictTypeCheck: true,
	})
	assert.Error(t, err)
}

func TestTransformFieldsNonStrictNoError(t *testing.T) {
	src := struct {
		Field int
	}{
		Field: 10,
	}
	dst := struct {
		Field string
	}{}

	err := TransformFields(&dst, src, &TransformFieldsOptions{
		StrictTypeCheck: false,
	})
	assert.NoError(t, err)
	assert.Equal(t, "10", dst.Field)
}

func TestTransformerSettersChaining(t *testing.T) {
	t1 := NewTransformer()

	dst := &struct{ A int }{A: 1}
	src := map[string]interface{}{"A": 2}
	opts := &TransformFieldsOptions{StrictTypeCheck: true}

	t2 := t1.SetDst(dst).SetSrc(src).SetOptions(opts)

	// 断言返回的是同一个实例，保证链式调用正确
	assert.Equal(t, t1, t2)

	// 断言字段设置正确
	assert.Equal(t, dst, t1.GetDst())
	assert.Equal(t, src, t1.GetSrc())
	assert.Equal(t, opts, t1.GetOptions())
}

func TestTransformFieldsNilPointerSrcToNonPtrDst(t *testing.T) {
	type Src struct {
		Ptr *int
	}
	type Dst struct {
		Ptr int
	}
	src := Src{Ptr: nil}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, dst.Ptr) // nil指针源转目标值类型，目标字段应为零值
}

func TestTransformFieldsSrcPointerToDstPointer(t *testing.T) {
	type Src struct {
		Val *int
	}
	type Dst struct {
		Val *int
	}
	v := 123
	src := Src{Val: &v}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dst.Val)
	assert.Equal(t, 123, *dst.Val)
}

func TestTransformFieldsSliceNilSrc(t *testing.T) {
	type Src struct {
		S []int
	}
	type Dst struct {
		S []int
	}
	src := Src{S: nil}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.Nil(t, dst.S)
}

func TestTransformFieldsSliceEmptySrc(t *testing.T) {
	type Src struct {
		S []int
	}
	type Dst struct {
		S []int
	}
	src := Src{S: []int{}}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dst.S)
	assert.Len(t, dst.S, 0)
}

func TestTransformFieldsMapNilSrc(t *testing.T) {
	type Src struct {
		M map[string]int
	}
	type Dst struct {
		M map[string]int
	}
	src := Src{M: nil}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.Nil(t, dst.M)
}

func TestTransformFieldsMapEmptySrc(t *testing.T) {
	type Src struct {
		M map[string]int
	}
	type Dst struct {
		M map[string]int
	}
	src := Src{M: map[string]int{}}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dst.M)
	assert.Len(t, dst.M, 0)
}

func TestTransformFieldsUnsupportedType(t *testing.T) {
	// 结构体字段类型不匹配且不支持转换
	type Src struct {
		Field chan int
	}
	type Dst struct {
		Field int
	}
	src := Src{Field: make(chan int)}
	var dst Dst

	err := TransformFields(&dst, src, &TransformFieldsOptions{StrictTypeCheck: true})
	assert.Error(t, err)
}

func TestTransformFieldsPtrToNonPtrWithNilSrc(t *testing.T) {
	type Src struct {
		Val *int
	}
	type Dst struct {
		Val int
	}
	src := Src{Val: nil}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, dst.Val)
}

func TestTransformFieldsNonPtrToPtrDst(t *testing.T) {
	type Src struct {
		Val int
	}
	type Dst struct {
		Val *int
	}
	src := Src{Val: 42}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.NotNil(t, dst.Val)
	assert.Equal(t, 42, *dst.Val)
}

func TestTransformFieldsEmptyStructs(t *testing.T) {
	type Src struct{}
	type Dst struct{}
	src := Src{}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
}

func TestTransformFieldsIgnorePrivateFields(t *testing.T) {
	type Src struct {
		Public  int
		private int
	}
	type Dst struct {
		Public  int
		private int
	}
	src := Src{Public: 1, private: 2}
	var dst Dst

	err := TransformFields(&dst, src, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, dst.Public)
	assert.Equal(t, 0, dst.private) // 私有字段不赋值，保持默认零值
}

func TestTransformFieldsTagNameMapping(t *testing.T) {
	type Src struct {
		Foo int
	}
	type Dst struct {
		Bar int `convert:"Foo"`
	}
	src := Src{Foo: 100}
	var dst Dst

	err := TransformFields(&dst, src, &TransformFieldsOptions{TransTagName: "convert"})
	assert.NoError(t, err)
	assert.Equal(t, 100, dst.Bar)
}

func TestTransformFieldsTimeFormatCustom(t *testing.T) {
	type Src struct {
		T time.Time
	}
	type Dst struct {
		T string
	}
	now := time.Now()
	src := Src{T: now}
	var dst Dst

	err := TransformFields(&dst, src, &TransformFieldsOptions{TimeFormat: time.RFC1123})
	assert.NoError(t, err)
	assert.Equal(t, now.Format(time.RFC1123), dst.T)
}
