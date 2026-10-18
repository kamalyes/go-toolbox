/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-04 09:59:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 10:15:07
 * @FilePath: \go-toolbox\pkg\safe\protobuf_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package safe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	x := 123
	res := Ptr(&x, func(v int) int { return v * 2 })
	assert.Equal(t, 246, *res)
	assert.Nil(t, Ptr(nil, func(v int) int { return v }))
}

func TestStringPtr(t *testing.T) {
	s := "abc"
	assert.Equal(t, "abc", StringPtr(&s))
	assert.Equal(t, "", StringPtr(nil))
}

func TestIntPtr(t *testing.T) {
	i := 42
	assert.Equal(t, 42, IntPtr(&i))
	assert.Equal(t, 0, IntPtr(nil))
}

func TestBoolPtr(t *testing.T) {
	b := true
	assert.Equal(t, true, BoolPtr(&b))
	assert.Equal(t, false, BoolPtr(nil))
}

func TestSlicePtr(t *testing.T) {
	s := []int{1, 2, 3}
	assert.Equal(t, []int{1, 2, 3}, SlicePtr(&s))
	assert.Equal(t, []int{}, SlicePtr[int](nil))
}

func TestFloat32Ptr(t *testing.T) {
	f := float32(1.5)
	assert.Equal(t, float32(1.5), Float32Ptr(&f))
	assert.Equal(t, float32(0), Float32Ptr(nil))
}

func TestFloat64Ptr(t *testing.T) {
	f := 2.5
	assert.Equal(t, 2.5, Float64Ptr(&f))
	assert.Equal(t, 0.0, Float64Ptr(nil))
}

func TestUintPtr(t *testing.T) {
	u := uint(7)
	assert.Equal(t, uint(7), UintPtr(&u))
	assert.Equal(t, uint(0), UintPtr(nil))
}

func TestInt32Ptr(t *testing.T) {
	i := int32(8)
	assert.Equal(t, int32(8), Int32Ptr(&i))
	assert.Equal(t, int32(0), Int32Ptr(nil))
}

func TestInt64Ptr(t *testing.T) {
	i := int64(9)
	assert.Equal(t, int64(9), Int64Ptr(&i))
	assert.Equal(t, int64(0), Int64Ptr(nil))
}

func TestDurationPtr(t *testing.T) {
	d := time.Second
	assert.Equal(t, time.Second, DurationPtr(&d))
	assert.Equal(t, time.Duration(0), DurationPtr(nil))
}

func TestBytesPtr(t *testing.T) {
	b := []byte{1, 2}
	assert.Equal(t, []byte{1, 2}, BytesPtr(&b))
	assert.Equal(t, []byte{}, BytesPtr(nil))
}

func TestPtrKV(t *testing.T) {
	kv := KV[string, int]{"key1": 1, "key2": 2}
	assert.Equal(t, kv, PtrKV(&kv))
	assert.Equal(t, KV[string, int]{}, PtrKV[string, int](nil))
}

// ——— 恶意输入审计：本库的承诺是「不崩」，任何 panic 都是缺陷 ———

func TestAudit_NoPanicOnPtrHelpers(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic: ptr helpers -> %v", r)
			}
		}()
		_ = StringPtr(nil)
		_ = IntPtr(nil)
		_ = BoolPtr(nil)
		_ = Float32Ptr(nil)
		_ = Float64Ptr(nil)
		_ = UintPtr(nil)
		_ = Int32Ptr(nil)
		_ = Int64Ptr(nil)
		_ = DurationPtr(nil)
		_ = BytesPtr(nil)
		_ = SlicePtr[int](nil)
		_ = PtrKV[string, int](nil)
		_ = Ptr[int, string](nil, func(i int) string { return "x" })
	}()
}
