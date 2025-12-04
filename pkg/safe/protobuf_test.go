/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-04 09:59:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 10:13:10
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
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestSafeTimeToTimestampPB(t *testing.T) {
	now := time.Now()
	timestamp := TimeToTimestampPB(&now)
	assert.NotNil(t, timestamp)
	assert.Equal(t, timestamppb.New(now), timestamp)

	assert.Nil(t, TimeToTimestampPB(nil))
}

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

func TestStringToPB(t *testing.T) {
	s := "hello"
	assert.Equal(t, wrapperspb.String("hello"), StringToPB(&s))
	assert.Nil(t, StringToPB(nil))
}

func TestBoolToPB(t *testing.T) {
	b := true
	assert.Equal(t, wrapperspb.Bool(true), BoolToPB(&b))
	assert.Nil(t, BoolToPB(nil))
}

func TestInt32ToPB(t *testing.T) {
	i := int32(10)
	assert.Equal(t, wrapperspb.Int32(10), Int32ToPB(&i))
	assert.Nil(t, Int32ToPB(nil))
}

func TestInt64ToPB(t *testing.T) {
	i := int64(11)
	assert.Equal(t, wrapperspb.Int64(11), Int64ToPB(&i))
	assert.Nil(t, Int64ToPB(nil))
}

func TestDoubleToPB(t *testing.T) {
	d := 3.14
	assert.Equal(t, wrapperspb.Double(3.14), DoubleToPB(&d))
	assert.Nil(t, DoubleToPB(nil))
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
