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
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type NestedStruct struct {
	Field1 string
	Field2 int
}

type TestCloneStruct struct {
	Name      string
	Age       int
	Nested    NestedStruct
	Friends   []string
	Settings  map[string]interface{}
	Pointer   *NestedStruct
	CreatedAt time.Time  // time.Time 字段
	UpdatedAt *time.Time // time.Time 指针字段
	private   string     // 未导出字段，不应该被复制
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
	now := time.Now()
	updatedAt := time.Now().Add(1 * time.Hour)
	src := &TestCloneStruct{
		Name:    "Alice",
		Age:     30,
		Nested:  NestedStruct{"Inner", 100},
		Friends: []string{"Bob", "Charlie"},
		Settings: map[string]interface{}{
			"theme": "dark",
			"count": 42,
		},
		Pointer:   &NestedStruct{"Pointer", 200},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
		private:   "private_value",
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
	assert.True(t, src.CreatedAt.Equal(dst.CreatedAt), "CreatedAt should be equal")
	assert.True(t, src.UpdatedAt.Equal(*dst.UpdatedAt), "UpdatedAt should be equal")

	// 修改源数据，确保目标数据不受影响
	src.Name = "Bob"
	src.Age = 40
	src.Friends[0] = "Dave"
	src.Settings["theme"] = "light"
	src.Pointer.Field1 = "Modified"
	src.CreatedAt = time.Now().Add(24 * time.Hour)
	newUpdatedAt := time.Now().Add(48 * time.Hour)
	src.UpdatedAt = &newUpdatedAt

	assert.NotEqual(t, src.Name, dst.Name)
	assert.NotEqual(t, src.Age, dst.Age)
	assert.NotEqual(t, src.Friends[0], dst.Friends[0])
	assert.NotEqual(t, src.Settings["theme"], dst.Settings["theme"])
	assert.NotEqual(t, src.Pointer.Field1, dst.Pointer.Field1)
	assert.False(t, src.CreatedAt.Equal(dst.CreatedAt), "CreatedAt should not be equal after modification")
	assert.False(t, src.UpdatedAt.Equal(*dst.UpdatedAt), "UpdatedAt should not be equal after modification")
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

// 测试 time.Time 类型的深拷贝
func TestDeepCopyTimeTime(t *testing.T) {
	now := time.Now()
	src := &now
	var dst time.Time
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	assert.True(t, src.Equal(dst), "time.Time should be copied correctly")
	assert.False(t, dst.IsZero(), "copied time.Time should not be zero")

	// 修改源时间
	newTime := time.Now().Add(24 * time.Hour)
	src = &newTime

	// 验证目标未受影响
	assert.False(t, src.Equal(dst), "dst should not change when src changes")
}

// 测试零值 time.Time 的深拷贝
func TestDeepCopyZeroTimeTime(t *testing.T) {
	src := time.Time{}
	var dst time.Time
	err := DeepCopy(&dst, &src)
	assert.NoError(t, err)
	assert.True(t, dst.IsZero(), "zero time.Time should remain zero after copy")
	assert.True(t, src.Equal(dst), "zero time.Time should be equal")
}

// 测试包含 time.Time 的结构体深拷贝
func TestDeepCopyStructWithTimeTime(t *testing.T) {
	type Message struct {
		ID        string
		Content   string
		CreatedAt time.Time
		UpdatedAt *time.Time
		DeletedAt *time.Time
	}

	now := time.Now()
	updatedAt := now.Add(1 * time.Hour)
	src := &Message{
		ID:        "msg-123",
		Content:   "Hello World",
		CreatedAt: now,
		UpdatedAt: &updatedAt,
		DeletedAt: nil,
	}

	var dst Message
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)

	// 验证所有字段
	assert.Equal(t, src.ID, dst.ID)
	assert.Equal(t, src.Content, dst.Content)
	assert.True(t, src.CreatedAt.Equal(dst.CreatedAt), "CreatedAt should be equal")
	assert.False(t, dst.CreatedAt.IsZero(), "CreatedAt should not be zero")
	assert.NotNil(t, dst.UpdatedAt)
	assert.True(t, src.UpdatedAt.Equal(*dst.UpdatedAt), "UpdatedAt should be equal")
	assert.Nil(t, dst.DeletedAt)

	// 修改源
	src.ID = "msg-456"
	src.Content = "Modified"
	src.CreatedAt = time.Now().Add(24 * time.Hour)
	newUpdatedAt := time.Now().Add(48 * time.Hour)
	src.UpdatedAt = &newUpdatedAt
	deletedAt := time.Now()
	src.DeletedAt = &deletedAt

	// 验证目标未受影响
	assert.Equal(t, "msg-123", dst.ID)
	assert.Equal(t, "Hello World", dst.Content)
	assert.False(t, src.CreatedAt.Equal(dst.CreatedAt))
	assert.False(t, src.UpdatedAt.Equal(*dst.UpdatedAt))
	assert.Nil(t, dst.DeletedAt)
}

// 测试包含 time.Time 的 Map 深拷贝
func TestDeepCopyMapWithTimeTime(t *testing.T) {
	now := time.Now()
	src := map[string]time.Time{
		"created": now,
		"updated": now.Add(1 * time.Hour),
	}
	var dst map[string]time.Time
	err := DeepCopy(&dst, &src)
	assert.NoError(t, err)

	// 验证相等
	assert.True(t, src["created"].Equal(dst["created"]))
	assert.True(t, src["updated"].Equal(dst["updated"]))
	assert.False(t, dst["created"].IsZero())
	assert.False(t, dst["updated"].IsZero())

	// 修改源
	src["created"] = time.Now().Add(24 * time.Hour)
	src["deleted"] = time.Now().Add(48 * time.Hour)

	// 验证目标未受影响
	assert.False(t, src["created"].Equal(dst["created"]))
	assert.NotContains(t, dst, "deleted")
}

// 测试包含 time.Time 的 Slice 深拷贝
func TestDeepCopySliceWithTimeTime(t *testing.T) {
	now := time.Now()
	src := []time.Time{
		now,
		now.Add(1 * time.Hour),
		now.Add(2 * time.Hour),
	}
	var dst []time.Time
	err := DeepCopy(&dst, &src)
	assert.NoError(t, err)

	// 验证相等
	assert.Equal(t, len(src), len(dst))
	for i := range src {
		assert.True(t, src[i].Equal(dst[i]), "time.Time at index %d should be equal", i)
		assert.False(t, dst[i].IsZero(), "time.Time at index %d should not be zero", i)
	}

	// 修改源
	src[0] = time.Now().Add(24 * time.Hour)
	src = append(src, time.Now().Add(48*time.Hour))

	// 验证目标未受影响
	assert.False(t, src[0].Equal(dst[0]))
	assert.Equal(t, 3, len(dst))
}

// ============================================================================
// 测试类型定义
// ============================================================================

// clonerTestType 实现 Cloner 接口的测试类型
type clonerTestType struct {
	Name     string
	Age      int
	Tags     []string
	Settings map[string]interface{}
	inner    string // 未导出字段
}

// CloneDeep 实现 Cloner 接口 — 手动深拷贝，零反射
func (c *clonerTestType) CloneDeep() any {
	cp := &clonerTestType{
		Name: c.Name,
		Age:  c.Age,
	}
	if c.Tags != nil {
		cp.Tags = make([]string, len(c.Tags))
		copy(cp.Tags, c.Tags)
	}
	if c.Settings != nil {
		cp.Settings = make(map[string]interface{}, len(c.Settings))
		for k, v := range c.Settings {
			cp.Settings[k] = v
		}
	}
	return cp
}

// tagSkipType 测试 deepcopy:"-" 标签跳过（仅对引用类型字段有效，值类型被 dst.Set 拷贝）
type tagSkipType struct {
	Name     string
	KeepName string
	Data     map[string]int `deepcopy:"-"`
}

// noExportType 仅含未导出字段
type noExportType struct {
	x int
	y string
}

// largeStruct 模拟真实业务结构体（多字段、混合类型）
type largeStruct struct {
	ID        string
	UserID    string
	Title     string
	Content   string
	Status    int
	Priority  int
	Tags      []string
	Metadata  map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
	Inner     *nestedData
}

type nestedData struct {
	A string
	B int
	C map[string]string
}

// ============================================================================
// 1. Cloner 接口测试
// ============================================================================

// 测试 Cloner 接口快速路径正确性
func TestClonerFastPath(t *testing.T) {
	src := &clonerTestType{
		Name:     "Alice",
		Age:      30,
		Tags:     []string{"go", "rust"},
		Settings: map[string]interface{}{"theme": "dark", "count": 42},
		inner:    "private",
	}
	var dst clonerTestType
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)

	assert.Equal(t, src.Name, dst.Name)
	assert.Equal(t, src.Age, dst.Age)
	assert.Equal(t, src.Tags, dst.Tags)
	assert.Equal(t, src.Settings, dst.Settings)

	// 修改源，验证独立性
	src.Name = "Bob"
	src.Tags[0] = "python"
	src.Settings["theme"] = "light"
	assert.Equal(t, "Alice", dst.Name)
	assert.Equal(t, "go", dst.Tags[0])
	assert.Equal(t, "dark", dst.Settings["theme"])
}

// 测试 Cloner 接口与 nil dst
func TestClonerFastPathNilDst(t *testing.T) {
	src := &clonerTestType{Name: "Test", Age: 1}
	var dst *clonerTestType
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	assert.NotNil(t, dst)
	assert.Equal(t, src.Name, dst.Name)
}

// 测试 Cloner 接口与空集合
func TestClonerFastPathEmptyCollections(t *testing.T) {
	src := &clonerTestType{
		Name:     "Empty",
		Tags:     []string{},
		Settings: map[string]interface{}{},
	}
	var dst clonerTestType
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	assert.NotNil(t, dst.Tags)
	assert.NotNil(t, dst.Settings)
	assert.Equal(t, 0, len(dst.Tags))
	assert.Equal(t, 0, len(dst.Settings))
}

// 测试 Cloner 接口与 nil 集合
func TestClonerFastPathNilCollections(t *testing.T) {
	src := &clonerTestType{
		Name:     "Nil",
		Tags:     nil,
		Settings: nil,
	}
	var dst clonerTestType
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	assert.Nil(t, dst.Tags)
	assert.Nil(t, dst.Settings)
}

// ============================================================================
// 2. struct 字段元数据缓存测试
// ============================================================================

// 测试预生成克隆闭包正确性
func TestStructCloneFnCache(t *testing.T) {
	// 有导出字段的结构体 — 验证 clone 函数可正常获取
	fn1 := getStructCloneFn(reflect.TypeOf(largeStruct{}))
	assert.NotNil(t, fn1)

	// 无导出字段的结构体 — 返回直接 Set 的函数
	fn2 := getStructCloneFn(reflect.TypeOf(noExportType{}))
	assert.NotNil(t, fn2)

	// time.Time 无导出字段
	fn3 := getStructCloneFn(reflect.TypeOf(time.Time{}))
	assert.NotNil(t, fn3)
}

// 测试缓存命中（多次调用返回同一函数）
func TestStructCloneFnCacheHit(t *testing.T) {
	typ := reflect.TypeOf(largeStruct{})
	fn1 := getStructCloneFn(typ)
	fn2 := getStructCloneFn(typ)
	// 闭包函数指针不可直接比较，但通过 fmt.Sprintf 可验证
	assert.Equal(t, fmt.Sprintf("%p", fn1), fmt.Sprintf("%p", fn2), "缓存应返回同一函数")
}

// 测试 deepcopy:"-" 标签跳过引用类型字段的深拷贝
func TestDeepCopyTagSkip(t *testing.T) {
	src := &tagSkipType{
		Name:     "Alice",
		KeepName: "ShouldKeep",
		Data:     map[string]int{"a": 1},
	}
	var dst tagSkipType
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)

	// 值类型字段被 dst.Set 拷贝
	assert.Equal(t, "Alice", dst.Name)
	assert.Equal(t, "ShouldKeep", dst.KeepName)
	// Data 有 deepcopy:"-" 标签，不被深拷贝，共享引用（dst.Set 浅拷贝）
	assert.Equal(t, 1, dst.Data["a"])

	// 修改源 Data，目标也受影响（浅拷贝共享引用，证明跳过了深拷贝）
	src.Data["a"] = 100
	assert.Equal(t, 100, dst.Data["a"], "deepcopy:\"-\" 字段应共享引用")
}

// 测试缓存路径下复杂结构体深拷贝正确性
func TestCachedStructDeepCopy(t *testing.T) {
	now := time.Now()
	src := &largeStruct{
		ID:        "msg-1",
		UserID:    "user-1",
		Title:     "Hello",
		Content:   "World",
		Status:    1,
		Priority:  5,
		Tags:      []string{"urgent", "review"},
		Metadata:  map[string]interface{}{"ip": "127.0.0.1", "port": 8080},
		CreatedAt: now,
		UpdatedAt: now.Add(1 * time.Hour),
		Inner: &nestedData{
			A: "inner-a",
			B: 42,
			C: map[string]string{"key": "val"},
		},
	}

	var dst largeStruct
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)

	// 验证所有字段
	assert.Equal(t, src.ID, dst.ID)
	assert.Equal(t, src.Tags, dst.Tags)
	assert.Equal(t, src.Metadata, dst.Metadata)
	assert.True(t, src.CreatedAt.Equal(dst.CreatedAt))
	assert.NotNil(t, dst.Inner)
	assert.Equal(t, src.Inner.A, dst.Inner.A)
	assert.Equal(t, src.Inner.C["key"], dst.Inner.C["key"])

	// 修改源，验证独立性
	src.Tags[0] = "modified"
	src.Metadata["ip"] = "0.0.0.0"
	src.Inner.A = "modified-a"
	src.Inner.C["key"] = "modified-val"

	assert.Equal(t, "urgent", dst.Tags[0])
	assert.Equal(t, "127.0.0.1", dst.Metadata["ip"])
	assert.Equal(t, "inner-a", dst.Inner.A)
	assert.Equal(t, "val", dst.Inner.C["key"])
}

// 测试无导出字段结构体走快速 Set 路径
func TestNoExportStructDeepCopy(t *testing.T) {
	src := &noExportType{x: 1, y: "hello"}
	var dst noExportType
	err := DeepCopy(&dst, src)
	assert.NoError(t, err)
	// 未导出字段不会被拷贝（反射限制），但不报错
}

// ============================================================================
// 3. 并发安全测试
// ============================================================================

// 测试并发 DeepCopy + 缓存写入安全
func TestConcurrentDeepCopy(t *testing.T) {
	src := &largeStruct{
		ID:       "concurrent",
		Tags:     []string{"a", "b"},
		Metadata: map[string]interface{}{"k": "v"},
		Inner:    &nestedData{A: "x", C: map[string]string{"k": "v"}},
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var dst largeStruct
			err := DeepCopy(&dst, src)
			assert.NoError(t, err)
			assert.Equal(t, src.ID, dst.ID)
		}()
	}
	wg.Wait()
}

// 测试并发 Cloner 快速路径
func TestConcurrentClonerFastPath(t *testing.T) {
	src := &clonerTestType{
		Name:     "concurrent-cloner",
		Tags:     []string{"x"},
		Settings: map[string]interface{}{"k": 1},
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var dst clonerTestType
			err := DeepCopy(&dst, src)
			assert.NoError(t, err)
			assert.Equal(t, src.Name, dst.Name)
		}()
	}
	wg.Wait()
}

// 测试不同类型并发首次缓存写入（缓存竞争安全）
func TestConcurrentCachePopulation(t *testing.T) {
	types := []interface{}{
		&largeStruct{},
		&tagSkipType{},
		&clonerTestType{},
		&nestedData{},
		&TestCloneStruct{},
	}

	var wg sync.WaitGroup
	for _, typ := range types {
		wg.Add(1)
		go func(v interface{}) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				// 并发触发不同类型的缓存首次写入
				dst := v // shallow copy
				err := DeepCopy(&dst, &v)
				_ = err
			}
		}(typ)
	}
	wg.Wait()
}

// ============================================================================
// 4. 性能压测
// ============================================================================

// 反射缓存路径 vs Cloner 快速路径 vs 无缓存（首次）
// 对比三种路径的深拷贝性能

// BenchmarkDeepCopyLargeStruct 反射 + 缓存路径深拷贝大型结构体
func BenchmarkDeepCopyLargeStruct(b *testing.B) {
	now := time.Now()
	src := &largeStruct{
		ID:        "bench-1",
		UserID:    "user-1",
		Title:     "Benchmark",
		Content:   "Content",
		Status:    1,
		Priority:  5,
		Tags:      []string{"tag1", "tag2", "tag3"},
		Metadata:  map[string]interface{}{"ip": "127.0.0.1", "port": 8080, "ttl": 3600},
		CreatedAt: now,
		UpdatedAt: now,
		Inner: &nestedData{
			A: "a",
			B: 42,
			C: map[string]string{"key1": "val1", "key2": "val2"},
		},
	}
	// 预热缓存
	var warmup largeStruct
	_ = DeepCopy(&warmup, src)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var dst largeStruct
		_ = DeepCopy(&dst, src)
	}
}

// BenchmarkDeepCopyLargeStructParallel 反射 + 缓存路径并行深拷贝
func BenchmarkDeepCopyLargeStructParallel(b *testing.B) {
	src := &largeStruct{
		ID:       "bench-par",
		Tags:     []string{"a", "b"},
		Metadata: map[string]interface{}{"k": "v"},
		Inner:    &nestedData{A: "x", C: map[string]string{"k": "v"}},
	}
	// 预热缓存
	var warmup largeStruct
	_ = DeepCopy(&warmup, src)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var dst largeStruct
			_ = DeepCopy(&dst, src)
		}
	})
}

// BenchmarkClonerFastPath Cloner 接口快速路径深拷贝
func BenchmarkClonerFastPath(b *testing.B) {
	src := &clonerTestType{
		Name:     "bench-cloner",
		Age:      30,
		Tags:     []string{"go", "rust", "python"},
		Settings: map[string]interface{}{"theme": "dark", "count": 42, "verbose": true},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var dst clonerTestType
		_ = DeepCopy(&dst, src)
	}
}

// BenchmarkClonerFastPathParallel Cloner 接口快速路径并行深拷贝
func BenchmarkClonerFastPathParallel(b *testing.B) {
	src := &clonerTestType{
		Name:     "bench-cloner-par",
		Tags:     []string{"a", "b"},
		Settings: map[string]interface{}{"k": "v"},
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var dst clonerTestType
			_ = DeepCopy(&dst, src)
		}
	})
}

// BenchmarkDeepCopyMap 反射路径深拷贝 map
func BenchmarkDeepCopyMap(b *testing.B) {
	src := &map[string]interface{}{
		"string": "hello",
		"int":    42,
		"nested": map[string]interface{}{"key": "value"},
		"slice":  []interface{}{1, 2, 3},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var dst map[string]interface{}
		_ = DeepCopy(&dst, src)
	}
}

// BenchmarkDeepCopySlice 反射路径深拷贝 slice
func BenchmarkDeepCopySlice(b *testing.B) {
	src := &[]NestedStruct{
		{"item1", 1},
		{"item2", 2},
		{"item3", 3},
		{"item4", 4},
		{"item5", 5},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var dst []NestedStruct
		_ = DeepCopy(&dst, src)
	}
}

// BenchmarkGetStructCloneFn 缓存查找性能（命中 vs 首次）
func BenchmarkGetStructCloneFn(b *testing.B) {
	typ := reflect.TypeOf(largeStruct{})
	// 预热
	getStructCloneFn(typ)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = getStructCloneFn(typ)
	}
}

// BenchmarkDeepCopySimpleStruct 简单结构体（仅值类型字段）深拷贝
func BenchmarkDeepCopySimpleStruct(b *testing.B) {
	type simple struct {
		A string
		B int
		C bool
		D time.Time
	}
	now := time.Now()
	src := &simple{A: "hello", B: 42, C: true, D: now}
	// 预热
	var warmup simple
	_ = DeepCopy(&warmup, src)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var dst simple
		_ = DeepCopy(&dst, src)
	}
}

// BenchmarkCloneGeneric Clone[T] 泛型函数 + Cloner 零反射路径
func BenchmarkCloneGeneric(b *testing.B) {
	src := &clonerTestType{
		Name:     "bench-clone-generic",
		Age:      30,
		Tags:     []string{"go", "rust", "python"},
		Settings: map[string]interface{}{"theme": "dark", "count": 42, "verbose": true},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Clone(src)
	}
}

// BenchmarkCloneGenericParallel Clone[T] 泛型函数并行
func BenchmarkCloneGenericParallel(b *testing.B) {
	src := &clonerTestType{
		Name:     "bench-clone-par",
		Tags:     []string{"a", "b"},
		Settings: map[string]interface{}{"k": "v"},
	}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Clone(src)
		}
	})
}

// TestCloneGeneric Clone[T] 泛型函数正确性
func TestCloneGeneric(t *testing.T) {
	src := &clonerTestType{
		Name:     "clone-test",
		Age:      25,
		Tags:     []string{"a", "b"},
		Settings: map[string]interface{}{"k": "v"},
	}
	cloned := Clone(src)
	assert.NotNil(t, cloned)
	assert.Equal(t, src.Name, cloned.Name)
	assert.Equal(t, src.Tags, cloned.Tags)

	// 修改源，验证独立性
	src.Tags[0] = "modified"
	src.Settings["k"] = "modified"
	assert.Equal(t, "a", cloned.Tags[0])
	assert.Equal(t, "v", cloned.Settings["k"])
}

// TestCloneGenericNil Clone[T] 泛型函数 nil 处理
func TestCloneGenericNil(t *testing.T) {
	var src *clonerTestType
	cloned := Clone(src)
	assert.Nil(t, cloned)
}

// TestCloneGenericFallback Clone[T] 非 Cloner 类型走反射兜底
func TestCloneGenericFallback(t *testing.T) {
	type noCloner struct {
		Name string
		Tags []string
	}
	src := &noCloner{Name: "fallback", Tags: []string{"x"}}
	cloned := Clone(src)
	assert.NotNil(t, cloned)
	assert.Equal(t, src.Name, cloned.Name)
	assert.Equal(t, src.Tags, cloned.Tags)

	src.Tags[0] = "modified"
	assert.Equal(t, "x", cloned.Tags[0])
}
