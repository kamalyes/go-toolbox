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
