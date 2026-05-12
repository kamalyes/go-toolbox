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
