/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-05 15:27:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-05 15:27:15
 * @FilePath: \go-toolbox\pkg\syncx\clone_test.go
 * @Description: syncx 克隆单元测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"testing"

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
	Pointer  *NestedStruct
	private  string // 未导出字段，不应该被复制
}

type ComplexMapStruct struct {
	Data map[string]map[string]int
	Tags map[int][]string
}

type SliceOfStructs struct {
	Items []NestedStruct
}

type ArrayStruct struct {
	Numbers [5]int
	Texts   [3]string
}

// 测试基本类型的深拷贝
func TestDeepCopyBasicType(t *testing.T) {
	var intSrc = new(int)
	*intSrc = 42
	var intDst int
	err := DeepCopy(&intDst, intSrc)
	assert.NoError(t, err)
	assert.Equal(t, *intSrc, intDst)

	// 修改源不影响目标
	*intSrc = 100
	assert.NotEqual(t, *intSrc, intDst)
}

// 测试字符串的深拷贝
func TestDeepCopyString(t *testing.T) {
	src := "Hello World"
	var dst string
	err := DeepCopy(&dst, &src)
	assert.NoError(t, err)
	assert.Equal(t, src, dst)
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
			"count": 42,
		},
		Pointer: &NestedStruct{"Pointer", 200},
		private: "private_value",
	}
	var dst TestCloneStruct
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)

	// 断言源和目标相等
	assert.Equal(t, src.Name, dst.Name)
	assert.Equal(t, src.Age, dst.Age)
	assert.Equal(t, src.Nested, dst.Nested)
	assert.Equal(t, src.Friends, dst.Friends)
	assert.Equal(t, src.Settings, dst.Settings)
	assert.Equal(t, *src.Pointer, *dst.Pointer)

	// 修改源数据，确保目标数据不受影响
	src.Name = "Bob"
	src.Age = 40
	src.Friends[0] = "Dave"
	src.Settings["theme"] = "light"
	src.Pointer.Field1 = "Modified"

	assert.NotEqual(t, src.Name, dst.Name)
	assert.NotEqual(t, src.Age, dst.Age)
	assert.NotEqual(t, src.Friends[0], dst.Friends[0])
	assert.NotEqual(t, src.Settings["theme"], dst.Settings["theme"])
	assert.NotEqual(t, src.Pointer.Field1, dst.Pointer.Field1)
}

// 测试 Map 的深拷贝
func TestDeepCopyMap(t *testing.T) {
	src := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	var dst map[string]int
	err := DeepCopy(&dst, &src)
	assert.NoError(t, err)
	assert.Equal(t, src, dst)

	// 修改源 map
	src["one"] = 100
	src["four"] = 4
	delete(src, "two")

	// 验证目标 map 未受影响
	assert.Equal(t, 1, dst["one"])
	assert.Equal(t, 2, dst["two"])
	assert.NotContains(t, dst, "four")
}

// 测试复杂嵌套 Map 的深拷贝
func TestDeepCopyComplexMap(t *testing.T) {
	src := &ComplexMapStruct{
		Data: map[string]map[string]int{
			"group1": {"a": 1, "b": 2},
			"group2": {"c": 3, "d": 4},
		},
		Tags: map[int][]string{
			1: {"tag1", "tag2"},
			2: {"tag3", "tag4"},
		},
	}
	var dst ComplexMapStruct
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)

	// 验证相等
	assert.Equal(t, src.Data, dst.Data)
	assert.Equal(t, src.Tags, dst.Tags)

	// 修改源数据
	src.Data["group1"]["a"] = 100
	src.Data["group3"] = map[string]int{"e": 5}
	src.Tags[1][0] = "modified"
	src.Tags[3] = []string{"tag5"}

	// 验证目标未受影响
	assert.Equal(t, 1, dst.Data["group1"]["a"])
	assert.NotContains(t, dst.Data, "group3")
	assert.Equal(t, "tag1", dst.Tags[1][0])
	assert.NotContains(t, dst.Tags, 3)
}

// 测试 Slice 的深拷贝
func TestDeepCopySlice(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	var dst []int
	err := DeepCopy(&dst, &src)
	assert.NoError(t, err)
	assert.Equal(t, src, dst)

	// 修改源 slice
	src[0] = 100
	src = append(src, 6)

	// 验证目标 slice 未受影响
	assert.Equal(t, 1, dst[0])
	assert.Equal(t, 5, len(dst))
}

// 测试结构体切片的深拷贝
func TestDeepCopySliceOfStructs(t *testing.T) {
	src := &SliceOfStructs{
		Items: []NestedStruct{
			{"First", 1},
			{"Second", 2},
			{"Third", 3},
		},
	}
	var dst SliceOfStructs
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	assert.Equal(t, src.Items, dst.Items)

	// 修改源
	src.Items[0].Field1 = "Modified"
	src.Items[0].Field2 = 100
	src.Items = append(src.Items, NestedStruct{"Fourth", 4})

	// 验证目标未受影响
	assert.Equal(t, "First", dst.Items[0].Field1)
	assert.Equal(t, 1, dst.Items[0].Field2)
	assert.Equal(t, 3, len(dst.Items))
}

// 测试数组的深拷贝
func TestDeepCopyArray(t *testing.T) {
	src := &ArrayStruct{
		Numbers: [5]int{1, 2, 3, 4, 5},
		Texts:   [3]string{"a", "b", "c"},
	}
	var dst ArrayStruct
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	assert.Equal(t, src.Numbers, dst.Numbers)
	assert.Equal(t, src.Texts, dst.Texts)

	// 修改源
	src.Numbers[0] = 100
	src.Texts[0] = "modified"

	// 验证目标未受影响
	assert.Equal(t, 1, dst.Numbers[0])
	assert.Equal(t, "a", dst.Texts[0])
}

// 测试指针的深拷贝
func TestDeepCopyPointer(t *testing.T) {
	nested := &NestedStruct{"Original", 100}
	src := &nested
	var dst *NestedStruct
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	assert.Equal(t, **src, *dst)

	// 修改源指针指向的值
	(*src).Field1 = "Modified"
	(*src).Field2 = 200

	// 验证目标未受影响
	assert.Equal(t, "Original", dst.Field1)
	assert.Equal(t, 100, dst.Field2)
}

// 测试 nil 指针的深拷贝
func TestDeepCopyNilPointer(t *testing.T) {
	var nilSrc *TestCloneStruct
	var nilDst TestCloneStruct
	err := DeepCopy(&nilDst, nilSrc)
	assert.Error(t, err)
	assert.Equal(t, nilDst, TestCloneStruct{})
}

// 测试结构体中的 nil 指针字段
func TestDeepCopyStructWithNilPointer(t *testing.T) {
	src := &TestCloneStruct{
		Name:    "Test",
		Pointer: nil,
	}
	var dst TestCloneStruct
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	assert.Nil(t, dst.Pointer)
}

// 测试空切片的深拷贝
func TestDeepCopyEmptySlice(t *testing.T) {
	srcSlice := &[]string{}
	var dstSlice []string
	err := DeepCopy(&dstSlice, srcSlice)
	assert.NoError(t, err)
	assert.Equal(t, dstSlice, []string{})
	assert.NotNil(t, dstSlice)
}

// 测试 nil 切片的深拷贝
func TestDeepCopyNilSlice(t *testing.T) {
	var srcSlice []string
	var dstSlice []string
	err := DeepCopy(&dstSlice, &srcSlice)
	assert.NoError(t, err)
	assert.Nil(t, dstSlice)
}

// 测试空 Map 的深拷贝
func TestDeepCopyEmptyMap(t *testing.T) {
	srcMap := make(map[string]int)
	var dstMap map[string]int
	err := DeepCopy(&dstMap, &srcMap)
	assert.NoError(t, err)
	assert.NotNil(t, dstMap)
	assert.Equal(t, 0, len(dstMap))
}

// 测试 nil Map 的深拷贝
func TestDeepCopyNilMap(t *testing.T) {
	var srcMap map[string]int
	var dstMap map[string]int
	err := DeepCopy(&dstMap, &srcMap)
	assert.NoError(t, err)
	assert.Nil(t, dstMap)
}

// 测试嵌套指针的深拷贝
func TestDeepCopyNestedPointer(t *testing.T) {
	nestedSrc := &NestedStruct{"Inner", 100}
	testStructWithPointer := &TestCloneStruct{
		Name:    "Test",
		Nested:  *nestedSrc,
		Pointer: nestedSrc,
	}
	var testStructWithPointerDst TestCloneStruct
	err := DeepCopy(&testStructWithPointerDst, testStructWithPointer)
	assert.NoError(t, err)
	assert.Equal(t, testStructWithPointer.Nested, testStructWithPointerDst.Nested)
	assert.Equal(t, *testStructWithPointer.Pointer, *testStructWithPointerDst.Pointer)

	// 修改源
	nestedSrc.Field1 = "Modified"

	// 验证目标未受影响
	assert.Equal(t, "Inner", testStructWithPointerDst.Nested.Field1)
	assert.Equal(t, "Inner", testStructWithPointerDst.Pointer.Field1)
}

// 测试接口类型的深拷贝
func TestDeepCopyInterface(t *testing.T) {
	src := map[string]interface{}{
		"string": "hello",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"nested": map[string]interface{}{
			"key": "value",
		},
		"slice": []interface{}{1, 2, 3},
	}
	var dst map[string]interface{}
	err := DeepCopy(&dst, &src)
	assert.NoError(t, err)
	assert.Equal(t, src, dst)

	// 修改源
	src["string"] = "world"
	src["nested"].(map[string]interface{})["key"] = "modified"
	src["slice"].([]interface{})[0] = 100

	// 验证目标未受影响
	assert.Equal(t, "hello", dst["string"])
	assert.Equal(t, "value", dst["nested"].(map[string]interface{})["key"])
	assert.Equal(t, 1, dst["slice"].([]interface{})[0])
}

// 测试类型不匹配
func TestDeepCopyTypeMismatch(t *testing.T) {
	src := "string"
	var dst int
	assert.Panics(t, func() {
		DeepCopy(&dst, &src)
	})
}

// 测试非指针参数
func TestDeepCopyNonPointer(t *testing.T) {
	src := 42
	dst := 0
	assert.Panics(t, func() {
		DeepCopy(dst, src)
	})
}

// 测试复杂的多层嵌套结构
func TestDeepCopyComplexNested(t *testing.T) {
	type Level3 struct {
		Value string
	}
	type Level2 struct {
		Data  map[string]*Level3
		Items []Level3
	}
	type Level1 struct {
		Nested   *Level2
		MapData  map[string]Level2
		SlicePtr []*Level3
	}

	src := &Level1{
		Nested: &Level2{
			Data: map[string]*Level3{
				"a": {"value_a"},
				"b": {"value_b"},
			},
			Items: []Level3{
				{"item1"},
				{"item2"},
			},
		},
		MapData: map[string]Level2{
			"key1": {
				Data:  map[string]*Level3{"x": {"x_value"}},
				Items: []Level3{{"nested_item"}},
			},
		},
		SlicePtr: []*Level3{
			{"ptr1"},
			{"ptr2"},
		},
	}

	var dst Level1
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)

	// 验证深层嵌套数据
	assert.Equal(t, src.Nested.Data["a"].Value, dst.Nested.Data["a"].Value)
	assert.Equal(t, src.Nested.Items[0].Value, dst.Nested.Items[0].Value)
	assert.Equal(t, src.MapData["key1"].Data["x"].Value, dst.MapData["key1"].Data["x"].Value)
	assert.Equal(t, src.SlicePtr[0].Value, dst.SlicePtr[0].Value)

	// 修改源数据的深层值
	src.Nested.Data["a"].Value = "modified_a"
	src.Nested.Items[0].Value = "modified_item1"
	src.MapData["key1"].Data["x"].Value = "modified_x"
	src.SlicePtr[0].Value = "modified_ptr1"

	// 验证目标未受影响
	assert.Equal(t, "value_a", dst.Nested.Data["a"].Value)
	assert.Equal(t, "item1", dst.Nested.Items[0].Value)
	assert.Equal(t, "x_value", dst.MapData["key1"].Data["x"].Value)
	assert.Equal(t, "ptr1", dst.SlicePtr[0].Value)
}
