/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\pbjson\coverage_test.go
 * @Description: pbjson 内部覆盖测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package pbjson

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ==================== proto_json.go 补充测试 ====================

// TestClassifyParam_Nil 测试 classifyParam 处理 nil
func TestClassifyParam_Nil(t *testing.T) {
	msg, data := classifyParam(nil)
	assert.Nil(t, msg)
	assert.Nil(t, data)
}

// TestClassifyParam_Default 测试 classifyParam 处理未知类型
func TestClassifyParam_Default(t *testing.T) {
	msg, data := classifyParam(12345)
	assert.Nil(t, msg)
	assert.Nil(t, data)
}

// TestClassifyParam_ProtoMessage 测试 classifyParam 处理 proto.Message
func TestClassifyParam_ProtoMessage(t *testing.T) {
	// wrapperspb.StringValue 实现了 fmt.Stringer，但它也是 proto.Message
	msg := wrapperspb.String("test")
	// proto.Message 分支优先匹配
	m, d := classifyParam(msg)
	assert.NotNil(t, m)
	assert.Nil(t, d)
}

// TestExtractMessageAndJSON_SecondParamData 测试 extractMessageAndJSON 第二个参数提供数据
func TestExtractMessageAndJSON_SecondParamData(t *testing.T) {
	// 第一个参数提供 message，第二个参数提供 data
	msg := wrapperspb.String("test")
	jsonStr := `"hello"`
	m, d := extractMessageAndJSON(msg, jsonStr)
	assert.NotNil(t, m)
	assert.NotEmpty(t, d)
}

// TestExtractMessageAndJSON_FirstParamData 测试 extractMessageAndJSON 第一个参数提供 data
func TestExtractMessageAndJSON_FirstParamData(t *testing.T) {
	msg := wrapperspb.String("test")
	jsonStr := `"hello"`
	// 顺序反转
	m, d := extractMessageAndJSON(jsonStr, msg)
	assert.NotNil(t, m)
	assert.NotEmpty(t, d)
}

// TestProtoJSONUnmarshal_NilMessageAndData 测试 ProtoJSONUnmarshal 当 msg 和 data 都为 nil 时
func TestProtoJSONUnmarshal_StringerParam(t *testing.T) {
	// 使用实现了 fmt.Stringer 但不是 proto.Message 的类型
	cs := &testStringer{data: `"stringer-data"`}

	// 传入 testStringer 作为参数，应该走 classifyParam 的 fmt.Stringer 分支
	var restored wrapperspb.StringValue
	err := ProtoJSONUnmarshal(&restored, cs)
	require.NoError(t, err)
	assert.Equal(t, "stringer-data", restored.GetValue())
}

// TestClassifyParam_Stringer 测试 classifyParam 处理 fmt.Stringer
func TestClassifyParam_Stringer(t *testing.T) {
	cs := &testStringer{data: "stringer-value"}
	msg, data := classifyParam(cs)
	assert.Nil(t, msg)
	assert.NotNil(t, data)
	assert.Equal(t, "stringer-value", string(data))
}

// ==================== proto_json_lenient.go 补充测试 ====================

// TestIsNumericTypeMismatchErr_NilError 测试 isNumericTypeMismatchErr 处理 nil 错误
func TestLenientProtoJSONOptions_Unmarshal_NonNumericError(t *testing.T) {
	opts := LenientProtoJSONOptions{}
	msg := &wrapperspb.StringValue{}
	// 传入数字给 string 字段，不是数值类型不匹配
	err := opts.Unmarshal([]byte(`123`), msg)
	assert.Error(t, err)
	// 错误应该不是数值类型不匹配
}

// TestLenientProtoJSONUnmarshal_NonNumericStringForInt 测试非数字字符串给 int64 字段
func TestLenientProtoJSONUnmarshal_NonNumericStringForInt(t *testing.T) {
	msg := &wrapperspb.Int64Value{}
	// 字符串 "hello" 不是数字，转换不会成功
	err := LenientProtoJSONUnmarshal([]byte(`"hello"`), msg)
	assert.Error(t, err)
}

// TestLenientProtoJSONUnmarshal_ConvertFails 测试转换数字字符串失败时返回原始错误
func TestLenientProtoJSONUnmarshal_ConvertFails(t *testing.T) {
	opts := LenientProtoJSONOptions{}
	msg := &wrapperspb.Int64Value{}
	// 传入 "abc" 给 int64 字段，转换后仍失败
	err := opts.Unmarshal([]byte(`"abc"`), msg)
	assert.Error(t, err)
}

// TestLenientProtoJSONUnmarshal_AllowPartial 测试 AllowPartial 选项
func TestLenientProtoJSONUnmarshal_AllowPartial(t *testing.T) {
	opts := LenientProtoJSONOptions{
		AllowPartial: true,
	}
	msg := &wrapperspb.Int64Value{}
	err := opts.Unmarshal([]byte(`42`), msg)
	require.NoError(t, err)
	assert.Equal(t, int64(42), msg.GetValue())
}

// TestLenientProtoJSONUnmarshal_DiscardUnknownAndAllowPartial 测试同时设置 DiscardUnknown 和 AllowPartial
func TestLenientProtoJSONUnmarshal_DiscardUnknownAndAllowPartial(t *testing.T) {
	opts := LenientProtoJSONOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	msg := &wrapperspb.Int64Value{}
	err := opts.Unmarshal([]byte(`42`), msg)
	require.NoError(t, err)
	assert.Equal(t, int64(42), msg.GetValue())
}
func (s *testStringer) String() string {
	return s.data
}

// protoSimpleStruct 包含 protobuf 字段的简单结构体，用于触发 scanJSONReflect 路径

// testStringer 用于测试 classifyParam 的 fmt.Stringer 分支
type testStringer struct {
	data string
}

func TestProtoJSONUnmarshal_NilMessageAndData(t *testing.T) {
	err := ProtoJSONUnmarshal(nil, nil)
	assert.NoError(t, err)
}

func TestIsNumericTypeMismatchErr_NilError(t *testing.T) {
	assert.False(t, isNumericTypeMismatchErr(nil))
}

// TestIsNumericTypeMismatchErr_NoInvalidValue 测试 isNumericTypeMismatchErr 错误信息中不含 "invalid value for"
func TestIsNumericTypeMismatchErr_NoInvalidValue(t *testing.T) {
	err := errors.New("some other error")
	assert.False(t, isNumericTypeMismatchErr(err))
}

// TestIsNumericTypeMismatchErr_HasInvalidValueButNoNumericType 测试 isNumericTypeMismatchErr 含 "invalid value for" 但无数值类型
func TestIsNumericTypeMismatchErr_HasInvalidValueButNoNumericType(t *testing.T) {
	err := errors.New("proto: invalid value for string field")
	assert.False(t, isNumericTypeMismatchErr(err))
}

// TestIsNumericTypeMismatchErr_Int64Type 测试 isNumericTypeMismatchErr 匹配 int64 类型
func TestIsNumericTypeMismatchErr_Int64Type(t *testing.T) {
	err := errors.New("proto: invalid value for int64 field: bad value")
	assert.True(t, isNumericTypeMismatchErr(err))
}

// TestIsNumericTypeMismatchErr_FloatType 测试 isNumericTypeMismatchErr 匹配 float 类型
func TestIsNumericTypeMismatchErr_FloatType(t *testing.T) {
	err := errors.New("proto: invalid value for float field: bad value")
	assert.True(t, isNumericTypeMismatchErr(err))
}

// TestIsNumericTypeMismatchErr_DoubleType 测试 isNumericTypeMismatchErr 匹配 double 类型
func TestIsNumericTypeMismatchErr_DoubleType(t *testing.T) {
	err := errors.New("proto: invalid value for double field: bad value")
	assert.True(t, isNumericTypeMismatchErr(err))
}

// TestIsNumericTypeMismatchErr_Uint32Type 测试 isNumericTypeMismatchErr 匹配 uint32 类型
func TestIsNumericTypeMismatchErr_Uint32Type(t *testing.T) {
	err := errors.New("proto: invalid value for uint32 field: bad value")
	assert.True(t, isNumericTypeMismatchErr(err))
}

// TestLenientProtoJSONOptions_Unmarshal_NonNumericError 测试 Unmarshal 非数值类型错误直接返回
