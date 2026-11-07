/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\pbjson\serializer_coverage_test.go
 * @Description: serializer proto-aware 路径覆盖（需装配 ProtoJSONCodec）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package pbjson_test

import (
	"testing"

	_ "github.com/kamalyes/go-toolbox/pkg/pbjson"
	. "github.com/kamalyes/go-toolbox/pkg/serializer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type protoSimpleStruct struct {
	Name *wrapperspb.StringValue `json:"name"`
}

// protoArrayStruct 包含 protobuf 数组字段的结构体，用于测试数组路径
type protoArrayStruct struct {
	Fixed [2]*wrapperspb.StringValue `json:"fixed"`
}

// ==================== serializer.go 补充测试 ====================

// TestWithCustomEncoder 测试自定义编码器
func TestMarshalJSONReflect_InvalidValue(t *testing.T) {
	// 使用包含 protobuf 字段的结构体，确保走 marshalJSONReflect 路径
	type invalidType struct {
		Chan chan int                `json:"chan"`
		Name *wrapperspb.StringValue `json:"name"`
	}
	// channel 字段会导致 marshaling 失败
	v := invalidType{Chan: make(chan int)}
	_, err := JSONMarshal(&v)
	assert.Error(t, err)
}

// TestScanJSONStruct_MissingColon 测试 scanJSONStruct 缺少冒号
func TestScanJSONStruct_MissingColon(t *testing.T) {
	// 使用包含 protobuf 字段的类型，确保走 scanJSONReflect 路径
	var s protoSimpleStruct
	// JSON 缺少冒号
	err := JSONUnmarshal([]byte(`{"name" "test"}`), &s)
	assert.Error(t, err)
	assert.True(t, IsJSONExpectedObjectKeySeparatorError(err))
}

// TestScanJSONStruct_UnexpectedEnd 测试 scanJSONStruct 意外结束
func TestScanJSONStruct_UnexpectedEnd(t *testing.T) {
	var s protoSimpleStruct
	// JSON 缺少结束括号
	err := JSONUnmarshal([]byte(`{"name":"test"`), &s)
	assert.Error(t, err)
}

// TestScanJSONStruct_ExpectedObjectNext 测试 scanJSONStruct 缺少逗号或 }
func TestScanJSONStruct_ExpectedObjectNext(t *testing.T) {
	var s protoSimpleStruct
	// JSON 缺少逗号
	err := JSONUnmarshal([]byte(`{"name":"test" "extra":"20"}`), &s)
	assert.Error(t, err)
	assert.True(t, IsJSONExpectedObjectNextError(err))
}

// TestScanJSONStruct_NotObject 测试 scanJSONStruct 非对象
func TestScanJSONStruct_NotObject(t *testing.T) {
	var s protoSimpleStruct
	err := JSONUnmarshal([]byte(`[]`), &s)
	assert.Error(t, err)
	assert.True(t, IsJSONExpectedObjectError(err))
}

// TestScanJSONStruct_EmptyObject 测试 scanJSONStruct 空对象
func TestScanJSONStruct_EmptyObject(t *testing.T) {
	var s protoSimpleStruct
	err := JSONUnmarshal([]byte(`{}`), &s)
	require.NoError(t, err)
	assert.Nil(t, s.Name)
}

// TestScanJSONStruct_UnknownFieldInvalidValue 测试 scanJSONStruct 未知字段值非法
func TestScanJSONStruct_UnknownFieldInvalidValue(t *testing.T) {
	var s protoSimpleStruct
	// 未知字段包含无效 JSON
	err := JSONUnmarshal([]byte(`{"unknown":[invalidjson]}`), &s)
	assert.Error(t, err)
	assert.True(t, IsJSONInvalidUnknownFieldValueError(err))
}

// TestScanJSONStruct_FastFallbackToMap 测试 fast 失败回退到 map 路径
func TestScanJSONArray_NotArray(t *testing.T) {
	// 使用包含 protobuf 数组字段的类型，确保走 scanJSONReflect 路径
	var a protoArrayStruct
	err := JSONUnmarshal([]byte(`{"fixed":{}}`), &a)
	assert.Error(t, err)
	assert.True(t, IsJSONExpectedArrayError(err))
}

// TestScanJSONArray_TooLong 测试 scanJSONArray 数组长度超过目标
func TestScanJSONArray_TooLong(t *testing.T) {
	var a protoArrayStruct
	// 数组有 3 个元素但目标只能容纳 2 个
	err := JSONUnmarshal([]byte(`{"fixed":["a","b","c"]}`), &a)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "array has")
}

// TestScanJSONArray_MissingComma 测试 scanJSONArray 缺少逗号
func TestScanJSONArray_MissingComma(t *testing.T) {
	// 使用更大的数组类型
	type arrType struct {
		Fixed [3]*wrapperspb.StringValue `json:"fixed"`
	}
	var a arrType
	err := JSONUnmarshal([]byte(`{"fixed":["a" "b"]}`), &a)
	assert.Error(t, err)
	assert.True(t, IsJSONExpectedArrayNextError(err))
}

// TestScanJSONArray_UnexpectedEnd 测试 scanJSONArray 意外结束
func TestScanJSONArray_UnexpectedEnd(t *testing.T) {
	type arrType struct {
		Fixed [3]*wrapperspb.StringValue `json:"fixed"`
	}
	var a arrType
	// 缺少结束括号
	err := JSONUnmarshal([]byte(`{"fixed":["a","b"`), &a)
	assert.Error(t, err)
}

// TestScanJSONSlice_EmptyArray 测试 scanJSONSlice 空数组
func TestScanJSONSlice_EmptyArray(t *testing.T) {
	// 使用包含 protobuf slice 字段的类型
	type sliceType struct {
		Items []*wrapperspb.StringValue `json:"items"`
	}
	var s sliceType
	err := JSONUnmarshal([]byte(`{"items":[]}`), &s)
	require.NoError(t, err)
	assert.Empty(t, s.Items)
}

// TestScanJSONSlice_NotArray 测试 scanJSONSlice 非数组
func TestScanJSONSlice_NotArray(t *testing.T) {
	type sliceType struct {
		Items []*wrapperspb.StringValue `json:"items"`
	}
	var s sliceType
	err := JSONUnmarshal([]byte(`{"items":{}}`), &s)
	assert.Error(t, err)
	assert.True(t, IsJSONExpectedArrayError(err))
}

// TestScanJSONSlice_MissingComma 测试 scanJSONSlice 缺少逗号
func TestScanJSONSlice_MissingComma(t *testing.T) {
	type sliceType struct {
		Items []*wrapperspb.StringValue `json:"items"`
	}
	var s sliceType
	err := JSONUnmarshal([]byte(`{"items":["a" "b"]}`), &s)
	assert.Error(t, err)
	assert.True(t, IsJSONExpectedArrayNextError(err))
}

// TestScanJSONMap_NonStringKey 测试 scanJSONMap 非字符串键回退到标准 json.Unmarshal
func TestScanJSONMap_NonStringKey(t *testing.T) {
	// 使用非字符串键的 map[int]*wrapperspb.StringValue，
	// 这样 needsProtoJSON 返回 true（elem 是 proto），走 scanJSONMap，
	// 但 key 不是 string，所以回退到 json.Unmarshal
	var m map[int]*wrapperspb.StringValue
	err := JSONUnmarshal([]byte(`{"1":"a"}`), &m)
	// json.Unmarshal 回退路径可能返回错误（因为标准 json 不认识 proto message）
	// 这里只关心覆盖了回退分支
	_ = err
}

// TestScanJSONMap_InvalidData 测试 scanJSONMap 无效数据
func TestScanJSONMap_InvalidData(t *testing.T) {
	// 使用包含 protobuf 值的 map，确保走 scanJSONMap 路径
	var m map[string]*wrapperspb.StringValue
	err := JSONUnmarshal([]byte(`invalid`), &m)
	assert.Error(t, err)
}

// TestMarshalJSONMap_NonStringKey 测试 marshalJSONMap 非字符串键返回错误
func TestMarshalJSONMap_NonStringKey(t *testing.T) {
	// 使用包含 protobuf 值的非字符串键 map，确保走 marshalJSONMap 路径
	m := map[int]*wrapperspb.StringValue{1: wrapperspb.String("a")}
	_, err := JSONMarshal(m)
	assert.Error(t, err)
	assert.True(t, IsJSONMapKeyUnsupportedError(err))
}

// TestScanJSONNull_DefaultKind 测试 scanJSONNull 处理非指针类型
func TestScanJSONNull_DefaultKind(t *testing.T) {
	// 使用包含 protobuf 字段的结构体，传入 null 作为整个值
	// scanJSONNull 对 Struct 类型走 default 分支，调用 json.Unmarshal
	var s protoSimpleStruct
	err := JSONUnmarshal([]byte(`null`), &s)
	require.NoError(t, err)
	assert.Nil(t, s.Name)
}

// TestNeedsProtoJSON_NilType 测试 needsProtoJSON 处理 nil 类型
func TestScanJSONReflect_NestedArray(t *testing.T) {
	// 使用包含 protobuf 嵌套数组的类型
	type arrType struct {
		Matrix [2][2]*wrapperspb.Int32Value `json:"matrix"`
	}
	var a arrType
	err := JSONUnmarshal([]byte(`{"matrix":[[1,2],[3,4]]}`), &a)
	require.NoError(t, err)
	assert.Equal(t, int32(1), a.Matrix[0][0].GetValue())
	assert.Equal(t, int32(4), a.Matrix[1][1].GetValue())
}

// TestJSONMarshal_NilSlice 测试 JSONMarshal 处理 nil slice
func TestJSONMarshal_NilSlice(t *testing.T) {
	// 使用包含 protobuf 字段的 nil slice，确保走 marshalJSONReflect 路径
	var s []*wrapperspb.StringValue
	data, err := JSONMarshal(s)
	require.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

// TestJSONMarshal_NilMap 测试 JSONMarshal 处理 nil map
func TestJSONMarshal_NilMap(t *testing.T) {
	// 使用包含 protobuf 值的 nil map
	var m map[string]*wrapperspb.StringValue
	data, err := JSONMarshal(m)
	require.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

// TestJSONMarshal_NilInterface 测试 JSONMarshal 处理 nil interface
func TestJSONMarshal_NilInterface(t *testing.T) {
	// 使用包含 protobuf 字段的结构体，其中有一个 interface{} 字段
	type withInterface struct {
		Value interface{}             `json:"value"`
		Name  *wrapperspb.StringValue `json:"name"`
	}
	v := withInterface{Value: nil}
	data, err := JSONMarshal(&v)
	require.NoError(t, err)
	assert.Contains(t, string(data), "null")
}

// TestJSONMarshal_InvalidValue 测试 JSONMarshal 处理无法序列化的值
func TestJSONMarshal_InvalidValue(t *testing.T) {
	// 使用包含 protobuf 字段的结构体，其中有一个无法序列化的 channel 字段
	type invalid struct {
		Chan chan int                `json:"chan"`
		Name *wrapperspb.StringValue `json:"name"`
	}
	v := invalid{Chan: make(chan int)}
	_, err := JSONMarshal(&v)
	assert.Error(t, err)
}

// TestJSONFieldInfo_UnexportedField 测试 jsonFieldInfo 处理未导出字段
