/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 15:36:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 15:52:36
 * @FilePath: \go-toolbox\pkg\protox\safe_test.go
 * @Description: proto 空安全转换测试（随函数从 safe 迁入）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package protox

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
func TestPtrToTime(t *testing.T) {
	now := time.Now().In(time.UTC) // 确保使用 UTC 时区
	timestamp := timestamppb.New(now)
	result := PtrToTime(timestamp)

	// 将预期值转换为 UTC
	assert.NotNil(t, result)
	assert.Equal(t, now.UTC(), *result) // 使用 UTC 进行比较

	assert.Nil(t, PtrToTime(nil))
}
func TestPtrToString(t *testing.T) {
	strValue := wrapperspb.String("test")
	result := PtrToString(strValue)
	assert.NotNil(t, result)
	assert.Equal(t, "test", *result)

	assert.Nil(t, PtrToString(nil))
}
func TestPtrToBool(t *testing.T) {
	boolValue := wrapperspb.Bool(true)
	result := PtrToBool(boolValue)
	assert.NotNil(t, result)
	assert.True(t, *result)

	assert.Nil(t, PtrToBool(nil))
}
func TestPtrToInt32(t *testing.T) {
	int32Value := wrapperspb.Int32(42)
	result := PtrToInt32(int32Value)
	assert.NotNil(t, result)
	assert.Equal(t, int32(42), *result)

	assert.Nil(t, PtrToInt32(nil))
}
func TestPtrToInt64(t *testing.T) {
	int64Value := wrapperspb.Int64(42)
	result := PtrToInt64(int64Value)
	assert.NotNil(t, result)
	assert.Equal(t, int64(42), *result)

	assert.Nil(t, PtrToInt64(nil))
}
func TestPtrToDouble(t *testing.T) {
	doubleValue := wrapperspb.Double(3.14)
	result := PtrToDouble(doubleValue)
	assert.NotNil(t, result)
	assert.Equal(t, 3.14, *result)

	assert.Nil(t, PtrToDouble(nil))
}
