/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 13:15:06
 * @FilePath: \go-toolbox\pkg\types\reflect_test.go
 * @Description: 反射工具测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type jsonTagTestStruct struct {
	Name   string `json:"name,omitempty"`
	Age    int    `json:",omitzero"`
	Hidden string `json:"-"`
	Plain  string
	inner  string
}

func TestReflectJSONHelpers(t *testing.T) {
	typeOfStruct := reflect.TypeOf(jsonTagTestStruct{})
	nameField, _ := typeOfStruct.FieldByName("Name")
	ageField, _ := typeOfStruct.FieldByName("Age")
	hiddenField, _ := typeOfStruct.FieldByName("Hidden")
	plainField, _ := typeOfStruct.FieldByName("Plain")
	innerField, _ := typeOfStruct.FieldByName("inner")

	assert.True(t, IsProtoMessageType(reflect.TypeOf(wrapperspb.String("x"))))
	assert.False(t, IsProtoMessageType(reflect.TypeOf(jsonTagTestStruct{})))
	assert.True(t, IsExportedField(nameField))
	assert.False(t, IsExportedField(innerField))
	assert.Equal(t, "name", ExtractJSONKey(nameField))
	assert.Equal(t, "name", JSONFieldName(nameField))
	assert.Equal(t, "Age", ExtractJSONKey(ageField))
	assert.Equal(t, "Plain", JSONFieldName(plainField))
	assert.Empty(t, ExtractJSONKey(hiddenField))
	assert.True(t, HasJSONTagOption(nameField, "omitempty"))
	assert.True(t, HasJSONTagOption(ageField, "omitzero"))
	assert.False(t, HasJSONTagOption(plainField, "omitempty"))
}

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

func TestResolveModelKey(t *testing.T) {
	tests := []struct {
		name     string
		tag      reflect.StructField
		expected string
	}{
		{"gorm column", reflect.StructField{Name: "Name", Tag: reflect.StructTag(`gorm:"column:name;type:varchar(255)" json:"name"`)}, "name"},
		{"json fallback", reflect.StructField{Name: "Label", Tag: reflect.StructTag(`json:"label,omitempty"`)}, "label"},
		{"json dash", reflect.StructField{Name: "Secret", Tag: reflect.StructTag(`json:"-"`)}, "-"},
		{"no tag", reflect.StructField{Name: "Count"}, "Count"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ResolveModelKey(tt.tag))
		})
	}
}

func TestResolvePBKey(t *testing.T) {
	tests := []struct {
		name     string
		tag      reflect.StructField
		expected string
	}{
		{"protobuf name", reflect.StructField{Name: "HostStatus", Tag: reflect.StructTag(`protobuf:"varint,2,opt,name=host_status,json=hostStatus,proto3" json:"host_status,omitempty"`)}, "host_status"},
		{"json fallback", reflect.StructField{Name: "Name", Tag: reflect.StructTag(`json:"name,omitempty"`)}, "name"},
		{"json dash", reflect.StructField{Name: "Ignored", Tag: reflect.StructTag(`json:"-"`)}, "-"},
		{"no tag", reflect.StructField{Name: "Count"}, "Count"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ResolvePBKey(tt.tag))
		})
	}
}

func TestExtractGormColumn(t *testing.T) {
	assert.Equal(t, "name", ExtractGormColumn("column:name;type:varchar(255)"))
	assert.Equal(t, "status", ExtractGormColumn("column:status;type:int"))
	assert.Equal(t, "", ExtractGormColumn("type:varchar(255)"))
	assert.Equal(t, "", ExtractGormColumn(""))
}

func TestExtractPBTagValue(t *testing.T) {
	assert.Equal(t, "host_status", ExtractPBTagValue("varint,2,opt,name=host_status,json=hostStatus,proto3", "name"))
	assert.Equal(t, "hostStatus", ExtractPBTagValue("varint,2,opt,name=host_status,json=hostStatus,proto3", "json"))
	assert.Equal(t, "", ExtractPBTagValue("varint,1,opt,proto3", "name"))
}

func TestUnwrapModelValue(t *testing.T) {
	i := 42
	assert.Equal(t, 42, UnwrapModelValue(&i))
	assert.Equal(t, "hello", UnwrapModelValue("hello"))
	assert.Nil(t, UnwrapModelValue((*int)(nil)))
	assert.Equal(t, 3.14, UnwrapModelValue(3.14))
}

// TestHasJSONTagOptionEdgeCases 测试 HasJSONTagOption 的各种边界情况
func TestHasJSONTagOptionEdgeCases(t *testing.T) {
	t.Run("空 options 返回 false", func(t *testing.T) {
		f := reflect.StructField{Name: "F", Tag: reflect.StructTag(`json:"name,omitempty"`)}
		assert.False(t, HasJSONTagOption(f))
	})

	t.Run("json tag 为空返回 false", func(t *testing.T) {
		f := reflect.StructField{Name: "F", Tag: reflect.StructTag(``)}
		assert.False(t, HasJSONTagOption(f, "omitempty"))
	})

	t.Run("json tag 为 dash 返回 false", func(t *testing.T) {
		f := reflect.StructField{Name: "F", Tag: reflect.StructTag(`json:"-"`)}
		assert.False(t, HasJSONTagOption(f, "omitempty"))
	})

	t.Run("json tag 无逗号返回 false", func(t *testing.T) {
		// tag 为 `json:"name"`，没有逗号，没有选项
		f := reflect.StructField{Name: "F", Tag: reflect.StructTag(`json:"name"`)}
		assert.False(t, HasJSONTagOption(f, "omitempty"))
	})

	t.Run("逗号在末尾返回 false", func(t *testing.T) {
		// tag 为 `json:"name,"`，逗号后面没有选项
		f := reflect.StructField{Name: "F", Tag: reflect.StructTag(`json:"name,"`)}
		assert.False(t, HasJSONTagOption(f, "omitempty"))
	})

	t.Run("选项不匹配返回 false", func(t *testing.T) {
		f := reflect.StructField{Name: "F", Tag: reflect.StructTag(`json:"name,omitempty"`)}
		assert.False(t, HasJSONTagOption(f, "omitzero"))
	})

	t.Run("多个选项中有一个匹配", func(t *testing.T) {
		f := reflect.StructField{Name: "F", Tag: reflect.StructTag(`json:"name,omitempty,string"`)}
		assert.True(t, HasJSONTagOption(f, "string"))
		assert.True(t, HasJSONTagOption(f, "omitempty"))
		assert.False(t, HasJSONTagOption(f, "omitzero"))
	})

	t.Run("多个期望选项中匹配任意一个", func(t *testing.T) {
		f := reflect.StructField{Name: "F", Tag: reflect.StructTag(`json:"name,omitempty"`)}
		assert.True(t, HasJSONTagOption(f, "omitzero", "omitempty"))
	})
}

// TestEnsureStructDefaults 测试 EnsureStructDefaults 函数
func TestEnsureStructDefaults(t *testing.T) {
	t.Run("非结构体类型直接返回", func(t *testing.T) {
		// 传入一个 int 值
		v := reflect.ValueOf(42)
		// 不应 panic
		EnsureStructDefaults(v)
	})

	t.Run("初始化 nil proto 指针字段", func(t *testing.T) {
		type S struct {
			PB *wrapperspb.StringValue
		}
		s := &S{}
		v := reflect.ValueOf(s).Elem()
		EnsureStructDefaults(v)
		assert.NotNil(t, s.PB)
	})

	t.Run("初始化 nil 结构体指针字段", func(t *testing.T) {
		type Inner struct {
			Name string
		}
		type S struct {
			Inner *Inner
		}
		s := &S{}
		v := reflect.ValueOf(s).Elem()
		EnsureStructDefaults(v)
		assert.NotNil(t, s.Inner)
	})

	t.Run("nil 基本类型指针字段不被初始化", func(t *testing.T) {
		type S struct {
			Num *int
		}
		s := &S{}
		v := reflect.ValueOf(s).Elem()
		EnsureStructDefaults(v)
		assert.Nil(t, s.Num)
	})

	t.Run("已初始化的指针字段保持不变", func(t *testing.T) {
		type Inner struct {
			Name string
		}
		existing := &Inner{Name: "existing"}
		type S struct {
			Inner *Inner
		}
		s := &S{Inner: existing}
		v := reflect.ValueOf(s).Elem()
		EnsureStructDefaults(v)
		assert.Same(t, existing, s.Inner)
	})

	t.Run("非指针字段不受影响", func(t *testing.T) {
		type S struct {
			Name string
			Age  int
		}
		s := &S{Name: "test", Age: 20}
		v := reflect.ValueOf(s).Elem()
		EnsureStructDefaults(v)
		assert.Equal(t, "test", s.Name)
		assert.Equal(t, 20, s.Age)
	})
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
