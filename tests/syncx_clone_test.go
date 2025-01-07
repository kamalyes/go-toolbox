/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-05 15:27:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-05 15:27:15
 * @FilePath: \go-toolbox\tests\syncx_clone_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

type NestedStruct struct {
	Field1 string
	Field2 int
}

type TestCloneStruct struct {
	Name     string
	Age      int
	Nested   NestedStruct
	Friends  []string
	Settings map[string]interface{}
}

// 测试基本类型的深拷贝
func TestDeepCopyBasicType(t *testing.T) {
	var intSrc = new(int)
	*intSrc = 42
	var intDst int
	syncx.DeepCopy(&intDst, intSrc)
	assert.Equal(t, *intSrc, intDst)
}

// 测试结构体的深拷贝
func TestDeepCopyStruct(t *testing.T) {
	src := &TestCloneStruct{
		Name:    "Alice",
		Age:     30,
		Nested:  NestedStruct{"Inner", 100},
		Friends: []string{"Bob", "Charlie"},
		Settings: map[string]interface{}{
			"theme": "dark",
		},
	}
	var dst TestCloneStruct
	syncx.DeepCopy(&dst, src)

	// 断言源和目标相等
	assert.Equal(t, src.Name, dst.Name)
	assert.Equal(t, src.Age, dst.Age)
	assert.Equal(t, src.Nested, dst.Nested)
	assert.Equal(t, src.Friends, dst.Friends)
	assert.Equal(t, src.Settings, dst.Settings)

	// 修改源数据，确保目标数据不受影响
	src.Name = "Bob"
	src.Friends[0] = "Dave"
	src.Settings["theme"] = "light"
	assert.NotEqual(t, src.Name, dst.Name)
	assert.NotEqual(t, src.Friends[0], dst.Friends[0])
	assert.NotEqual(t, src.Settings["theme"], dst.Settings["theme"])
}

// 测试空指针的深拷贝
func TestDeepCopyNilPointer(t *testing.T) {
	var nilSrc *TestCloneStruct
	var nilDst TestCloneStruct
	err := syncx.DeepCopy(&nilDst, nilSrc)
	assert.Error(t, err)
	assert.Equal(t, nilDst, TestCloneStruct{})
}

// 测试空切片的深拷贝
func TestDeepCopyEmptySlice(t *testing.T) {
	srcSlice := &[]string{}
	var dstSlice []string
	err := syncx.DeepCopy(&dstSlice, srcSlice)
	assert.NoError(t, err)                // 确保没有错误
	assert.Equal(t, dstSlice, []string{}) // 检查目标为空切片
}

// 测试嵌套指针的深拷贝
func TestDeepCopyNestedPointer(t *testing.T) {
	nestedSrc := &NestedStruct{"Inner", 100}
	testStructWithPointer := &TestCloneStruct{Nested: *nestedSrc}
	var testStructWithPointerDst TestCloneStruct
	err := syncx.DeepCopy(&testStructWithPointerDst, testStructWithPointer)
	assert.NoError(t, err) // 确保没有错误
	assert.Equal(t, testStructWithPointer.Nested, testStructWithPointerDst.Nested)
}
