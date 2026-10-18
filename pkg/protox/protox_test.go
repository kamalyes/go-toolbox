/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 15:36:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 15:52:36
 * @FilePath: \go-toolbox\pkg\protox\protox_test.go
 * @Description: protobuf 反射辅助测试（随函数从 types 迁入）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package protox

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUnwrapPBValue(t *testing.T) {
	tests := []struct{ input, expected interface{} }{
		{wrapperspb.String("hello"), "hello"},
		{wrapperspb.Bool(true), true},
		{wrapperspb.Int32(5), int32(5)},
		{wrapperspb.Int64(99), int64(99)},
		{wrapperspb.UInt32(7), uint32(7)},
		{wrapperspb.UInt64(42), uint64(42)},
		{wrapperspb.Float(3.14), float32(3.14)},
		{wrapperspb.Double(99.5), 99.5},
		{wrapperspb.Bytes([]byte("hi")), []byte("hi")},
		{"plain", "plain"},
		{(*wrapperspb.StringValue)(nil), nil},
		{(*wrapperspb.BoolValue)(nil), nil},
		{(*wrapperspb.Int32Value)(nil), nil},
		{(*wrapperspb.Int64Value)(nil), nil},
		{(*wrapperspb.UInt32Value)(nil), nil},
		{(*wrapperspb.UInt64Value)(nil), nil},
		{(*wrapperspb.FloatValue)(nil), nil},
		{(*wrapperspb.DoubleValue)(nil), nil},
		{(*wrapperspb.BytesValue)(nil), nil},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, UnwrapPBValue(tt.input))
	}
}

// TestNewProtoMessage 测试 NewProtoMessage 泛型函数
func TestNewProtoMessage(t *testing.T) {
	t.Run("创建 StringValue", func(t *testing.T) {
		msg := NewProtoMessage[*wrapperspb.StringValue]()
		assert.NotNil(t, msg)
		assert.IsType(t, &wrapperspb.StringValue{}, msg)
	})

	t.Run("创建 Int32Value", func(t *testing.T) {
		msg := NewProtoMessage[*wrapperspb.Int32Value]()
		assert.NotNil(t, msg)
		assert.IsType(t, &wrapperspb.Int32Value{}, msg)
	})

	t.Run("创建 BoolValue", func(t *testing.T) {
		msg := NewProtoMessage[*wrapperspb.BoolValue]()
		assert.NotNil(t, msg)
		assert.IsType(t, &wrapperspb.BoolValue{}, msg)
	})
}

// TestIsProtoMessageType 测试 proto.Message 类型判定（原在 types.TestReflectJSONHelpers 内）
// [EN] Moved out of types' test together with the function itself.
func TestIsProtoMessageType(t *testing.T) {
	type plainStruct struct{ Name string }

	assert.True(t, IsProtoMessageType(reflect.TypeOf(wrapperspb.String("x"))))
	assert.False(t, IsProtoMessageType(reflect.TypeOf(plainStruct{})))
	assert.False(t, IsProtoMessageType(nil))
}
