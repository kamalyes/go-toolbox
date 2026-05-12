/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-16 13:59:55
 * @FilePath: \go-toolbox\pkg\validator\base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TestStruct struct {
	Name  string `validate:"notEmpty"`
	Age   int    `validate:"ge=0"`
	Email string `validate:"regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
}

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{"", true},                                // 空字符串
		{"Hello", false},                          // 非空字符串
		{"null", false},                           // "null" 字符串
		{"NULL", false},                           // "NULL" 字符串（大写）
		{"Null", false},                           // "Null" 字符串（混合）
		{" null ", false},                         // 带空格的 "null"
		{"undefined", false},                      // "undefined" 字符串
		{"UNDEFINED", false},                      // "UNDEFINED" 字符串（大写）
		{" undefined ", false},                    // 带空格的 "undefined"
		{nil, true},                               // nil 值
		{0, true},                                 // 整数 0
		{1, false},                                // 非零整数
		{[]int{}, true},                           // 空切片
		{[]int{1, 2}, false},                      // 非空切片
		{map[string]int{}, true},                  // 空映射
		{map[string]int{"key": 1}, false},         // 非空映射
		{struct{}{}, true},                        // 空结构体
		{TestStruct{}, true},                      // 自定义结构体，所有字段零值
		{TestStruct{Name: "Test"}, false},         // 自定义结构体，非零字段
		{TestStruct{Name: "", Age: 0}, true},      // 自定义结构体，所有字段零值
		{TestStruct{Name: "Test", Age: 1}, false}, // 自定义结构体，至少一个非零字段
		{struct{ A int }{1}, false},               // 非空结构体
		{struct{ A interface{} }{nil}, true},      // 包含 nil 的结构体
		{make(chan int), false},                   // 非空通道
	}

	for _, test := range tests {
		t.Run(func() string {
			if test.value == nil {
				return "nil"
			}
			return reflect.TypeOf(test.value).String()
		}(), func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestHasEmpty(t *testing.T) {
	tests := []struct {
		elems    []interface{}
		expected bool
		count    int
	}{
		{[]interface{}{"", "data", nil}, true, 2},
		{[]interface{}{"data1", "data2"}, false, 0},
		{[]interface{}{0, 1, 2}, true, 1},
		{[]interface{}{0, "", nil}, true, 3},
	}

	for _, test := range tests {
		t.Run("HasEmpty", func(t *testing.T) {
			result, count := HasEmpty(test.elems)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.count, count)
		})
	}
}

func TestIsAllEmpty(t *testing.T) {
	tests := []struct {
		elems    []interface{}
		expected bool
	}{
		{[]interface{}{"", nil}, true},
		{[]interface{}{"data", nil}, false},
		{[]interface{}{0, 0}, true},
		{[]interface{}{1, 0}, false},
	}

	for _, test := range tests {
		t.Run("IsAllEmpty", func(t *testing.T) {
			result := IsAllEmpty(test.elems)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsUndefined(t *testing.T) {
	tests := []struct {
		str      string
		expected bool
	}{
		{"undefined", true},
		{"Undefined", true},
		{"UNDEFINED", true},
		{" undefined ", true},
		{"defined", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			result := IsUndefined(test.str)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsNull(t *testing.T) {
	tests := []struct {
		str      string
		expected bool
	}{
		{"null", true},
		{"Null", true},
		{"NULL", true},
		{" null ", true},
		{"", false},
		{"nil", false},
		{"nothing", false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			result := IsNull(test.str)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIfNullOrUndefined(t *testing.T) {
	assert := require.New(t)
	assert.True(IfNullOrUndefined("null"))
	assert.True(IfNullOrUndefined("NULL"))
	assert.True(IfNullOrUndefined(" null "))
	assert.True(IfNullOrUndefined("undefined"))
	assert.True(IfNullOrUndefined("UNDEFINED"))
	assert.True(IfNullOrUndefined(" undefined "))
	assert.False(IfNullOrUndefined(""))
	assert.False(IfNullOrUndefined("hello"))
	assert.False(IfNullOrUndefined("nil"))
}

func TestContainsChinese(t *testing.T) {
	tests := []struct {
		str      string
		expected bool
	}{
		{"Hello 你好", true},
		{"Hello World", false},
		{"", false},
		{"123", false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			result := ContainsChinese(test.str)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestEmptyToDefault(t *testing.T) {
	tests := []struct {
		str        string
		defaultStr string
		expected   string
	}{
		{"", "default", "default"},
		{"value", "default", "value"},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			result := EmptyToDefault(test.str, test.defaultStr)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsNil(t *testing.T) {
	// 测试 nil interface
	var nilInterface interface{}
	assert.True(t, IsNil(nilInterface), "Expected nil interface to return true")

	// 测试 nil map
	var nilMap map[string]int
	assert.True(t, IsNil(nilMap), "Expected nil map to return true")

	// 测试空 map
	emptyMap := make(map[string]int)
	assert.False(t, IsNil(emptyMap), "Expected empty map to return false")

	// 测试非 nil map
	nonNilMap := map[string]int{"key": 1}
	assert.False(t, IsNil(nonNilMap), "Expected non-nil map to return false")

	// 测试指向 nil 的 map
	var ptrToNilMap *map[string]int
	assert.True(t, IsNil(ptrToNilMap), "Expected pointer to nil map to return true")

	// 测试指向空 map 的指针
	ptrToEmptyMap := &emptyMap
	assert.False(t, IsNil(ptrToEmptyMap), "Expected pointer to empty map to return false")

	// 测试非 nil 指针
	num := 42
	ptrToNum := &num
	assert.False(t, IsNil(ptrToNum), "Expected pointer to non-nil value to return false")

	// 测试 nil 切片
	var nilSlice []int
	assert.True(t, IsNil(nilSlice), "Expected nil slice to return true")

	// 测试空切片
	emptySlice := []int{}
	assert.False(t, IsNil(emptySlice), "Expected empty slice to return false")

	// 测试非 nil 切片
	nonNilSlice := []int{1, 2, 3}
	assert.False(t, IsNil(nonNilSlice), "Expected non-nil slice to return false")

	// 测试指向 nil 切片的指针
	var ptrToNilSlice *[]int
	assert.True(t, IsNil(ptrToNilSlice), "Expected pointer to nil slice to return true")

	// 测试指向空切片的指针
	ptrToEmptySlice := &emptySlice
	assert.False(t, IsNil(ptrToEmptySlice), "Expected pointer to empty slice to return false")

	// 测试 nil 通道
	var nilChan chan int
	assert.True(t, IsNil(nilChan), "Expected nil channel to return true")

	// 测试空通道
	emptyChan := make(chan int)
	assert.False(t, IsNil(emptyChan), "Expected empty channel to return false")

	// 测试指向 nil 通道的指针
	var ptrToNilChan *chan int
	assert.True(t, IsNil(ptrToNilChan), "Expected pointer to nil channel to return true")

	// 测试指向非 nil 通道的指针
	nonNilChan := make(chan int, 1)
	assert.False(t, IsNil(nonNilChan), "Expected non-nil channel to return false")

	// 测试 nil 接口
	var nilInterfaceValue interface{}
	assert.True(t, IsNil(nilInterfaceValue), "Expected nil interface value to return true")

	// 测试指向非 nil 接口的指针
	var nonNilInterfaceValue interface{} = 42
	ptrToNonNilInterface := &nonNilInterfaceValue
	assert.False(t, IsNil(ptrToNonNilInterface), "Expected pointer to non-nil interface to return false")
}

func TestIsFuncType(t *testing.T) {
	type FuncType func(int) int
	type MyStruct struct{ A int }

	tests := []struct {
		name     string
		typCheck func() bool
		want     bool
	}{
		{"int", func() bool { return IsFuncType[int]() }, false},
		{"string", func() bool { return IsFuncType[string]() }, false},
		{"struct", func() bool { return IsFuncType[MyStruct]() }, false},
		{"pointer", func() bool { return IsFuncType[*MyStruct]() }, false},
		{"slice", func() bool { return IsFuncType[[]int]() }, false},
		{"map", func() bool { return IsFuncType[map[string]int]() }, false},
		{"func type", func() bool { return IsFuncType[FuncType]() }, true},
		{"func literal type", func() bool { return IsFuncType[func(int) int]() }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.typCheck()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsEmptyPointer(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{(*int)(nil), true}, // nil 指针
		{new(int), true},    // 指向零值的指针
		{func() *int { i := 1; return &i }(), false},                // 指向非零值的指针
		{(*time.Time)(nil), true},                                   // nil time.Time 指针
		{new(time.Time), true},                                      // 零值 time.Time 指针
		{func() *time.Time { t := time.Now(); return &t }(), false}, // 非零值 time.Time 指针
	}

	for _, test := range tests {
		t.Run(func() string {
			if test.value == nil {
				return "nil"
			}
			return reflect.TypeOf(test.value).String()
		}(), func(t *testing.T) {
			result := IsEmptyPointer(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsEmptyStruct(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{struct{}{}, true},                        // 空结构体
		{TestStruct{}, true},                      // 自定义结构体，所有字段零值
		{TestStruct{Name: "Test"}, false},         // 自定义结构体，非零字段
		{TestStruct{Name: "", Age: 0}, true},      // 自定义结构体，所有字段零值
		{TestStruct{Name: "Test", Age: 1}, false}, // 自定义结构体，至少一个非零字段
		{struct{ A int }{1}, false},               // 非空结构体
		{struct{ A interface{} }{nil}, true},      // 包含 nil 的结构体
	}

	for _, test := range tests {
		t.Run(func() string {
			return reflect.TypeOf(test.value).String()
		}(), func(t *testing.T) {
			result := IsEmptyStruct(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCheckEmptyTimePointer(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{(*time.Time)(nil), true}, // nil time.Time 指针
		{new(time.Time), true},    // 零值 time.Time 指针
		{func() *time.Time { t := time.Now(); return &t }(), false}, // 非零值 time.Time 指针
		{(*timestamppb.Timestamp)(nil), true},                       // nil protobuf Timestamp 指针
		{&timestamppb.Timestamp{}, true},                            // 零值 protobuf Timestamp 指针
		{&timestamppb.Timestamp{Seconds: 0}, true},                  // Seconds=0 protobuf Timestamp 指针
		{&timestamppb.Timestamp{Seconds: 1}, false},                 // 非零值 protobuf Timestamp 指针
		{timestamppb.New(time.Now()), false},                        // 非零值 protobuf Timestamp 指针
	}

	for _, test := range tests {
		t.Run(func() string {
			return reflect.TypeOf(test.value).String()
		}(), func(t *testing.T) {
			result, _ := CheckEmptyTimePointer(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCheckEmptyTimeStruct(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"zero time.Time", time.Time{}, true},
		{"non-zero time.Time", time.Now(), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, _ := CheckEmptyTimeStruct(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, result)
		})
	}

	// 单独测试 protobuf Timestamp（避免 copylocks）
	t.Run("zero protobuf Timestamp", func(t *testing.T) {
		ts := timestamppb.Timestamp{}
		result, _ := CheckEmptyTimeStruct(reflect.ValueOf(ts))
		assert.Equal(t, true, result)
	})
}

func TestIsTimeEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    *time.Time
		expected bool
	}{
		{"nil pointer", nil, true},
		{"zero value", &time.Time{}, true},
		{"unix zero (1970-01-01)", func() *time.Time { t := time.Unix(0, 0); return &t }(), true},
		{"before unix zero", func() *time.Time { t := time.Unix(-1, 0); return &t }(), true},
		{"after unix zero", func() *time.Time { t := time.Unix(1, 0); return &t }(), false},
		{"now", func() *time.Time { t := time.Now(); return &t }(), false},
		{"specific date", func() *time.Time { t := time.Date(2025, 12, 16, 0, 0, 0, 0, time.UTC); return &t }(), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsTimeEmpty(test.value)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsProtobufTimestampEmpty(t *testing.T) {
	tests := []struct {
		name     string
		create   func() timestamppb.Timestamp
		expected bool
	}{
		{"zero value", func() timestamppb.Timestamp { return timestamppb.Timestamp{} }, true},
		{"seconds = 0", func() timestamppb.Timestamp { return timestamppb.Timestamp{Seconds: 0} }, true},
		{"negative seconds", func() timestamppb.Timestamp { return timestamppb.Timestamp{Seconds: -1} }, true},
		{"positive seconds", func() timestamppb.Timestamp { return timestamppb.Timestamp{Seconds: 1} }, false},
		{"from time.Now()", func() timestamppb.Timestamp { return *timestamppb.New(time.Now()) }, false},
		{"from specific date", func() timestamppb.Timestamp { return *timestamppb.New(time.Date(2025, 12, 16, 0, 0, 0, 0, time.UTC)) }, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := test.create()
			result := IsProtobufTimestampEmpty(reflect.ValueOf(ts))
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsCEmpty(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() bool
		expected bool
	}{
		// 基本类型
		{"int zero", func() bool { return IsCEmpty(0) }, true},
		{"int non-zero", func() bool { return IsCEmpty(1) }, false},
		{"string empty", func() bool { return IsCEmpty("") }, true},
		{"string non-empty", func() bool { return IsCEmpty("hello") }, false},
		{"bool false", func() bool { return IsCEmpty(false) }, true},
		{"bool true", func() bool { return IsCEmpty(true) }, false},
		{"float zero", func() bool { return IsCEmpty(0.0) }, true},
		{"float non-zero", func() bool { return IsCEmpty(1.5) }, false},

		// 指针类型
		{"nil pointer", func() bool { return IsCEmpty((*int)(nil)) }, true},
		{"non-nil pointer", func() bool {
			i := 42
			return IsCEmpty(&i)
		}, false},

		// 结构体类型
		{"empty struct", func() bool {
			type Empty struct{}
			return IsCEmpty(Empty{})
		}, true},
		{"struct with value", func() bool {
			type Point struct{ X, Y int }
			return IsCEmpty(Point{X: 1, Y: 0})
		}, false},
		{"zero struct", func() bool {
			type Point struct{ X, Y int }
			return IsCEmpty(Point{})
		}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.testFunc()
			assert.Equal(t, test.expected, result)
		})
	}
}

// 定义常量
const (
	invalidField = "invalid-field"
	invalidSpace = "invalid field"
)

// TestIsSafeFieldName 测试字段名安全检查函数
func TestIsSafeFieldName(t *testing.T) {
	testCases := []struct {
		name     string
		field    string
		expected bool
	}{
		{"空字符串", "", false},
		{"简单字段名", "id", true},
		{"下划线字段名", "user_id", true},
		{"数字结尾", "field123", true},
		{"大写字母", "UserId", true},
		{"点号表示法", "users.id", true},
		{"包含空格", "user id", false},
		{"包含单引号", "user'id", false},
		{"包含分号", "id;DROP", false},
		{"包含星号", "id*", false},
		{"包含减号", "user-id", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsSafeFieldName(tc.field)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestIsAllowedField 测试 IsAllowedField 函数
func TestIsAllowedField(t *testing.T) {
	allowedFields := []string{"name", "age", "email"}

	tests := []struct {
		field   string
		allowed bool
	}{
		{"name", true},           // 在白名单中
		{"age", true},            // 在白名单中
		{"email", true},          // 在白名单中
		{"invalid_field", false}, // 不在白名单中
		{"", false},              // 空字符串
		{invalidSpace, false},    // 包含空格
		{invalidField, false},    // 包含连字符
	}

	for _, test := range tests {
		result := IsAllowedField(test.field, allowedFields)
		assert.Equal(t, test.allowed, result, "Expected IsAllowedField(%q) to be %v", test.field, test.allowed)
	}

	// 测试没有白名单的情况
	testsNoWhitelist := []struct {
		field  string
		isSafe bool
	}{
		{"valid_field", true},
		{invalidSpace, false}, // 包含空格
		{invalidField, false}, // 包含连字符
		{"", false},           // 空字符串
	}

	for _, test := range testsNoWhitelist {
		result := IsAllowedField(test.field) // 不传白名单
		assert.Equal(t, test.isSafe, result, "Expected IsAllowedField(%q) to be %v", test.field, test.isSafe)
	}
}

// TestCheckEmptyTimePointer_UnexportedFields 测试修复后的代码能正确处理未导出字段
// 这个测试验证 CanInterface() 检查是否正确工作，避免反射 panic
func TestCheckEmptyTimePointer_UnexportedFields(t *testing.T) {
	// 定义一个包含未导出字段的结构体
	type ConfigWithUnexportedFields struct {
		PublicField  string
		privateField string // 未导出字段
	}

	tests := []struct {
		name     string
		value    any
		expected bool
		handled  bool
	}{
		{
			name:     "struct with unexported fields",
			value:    &ConfigWithUnexportedFields{PublicField: "test", privateField: "private"},
			expected: false,
			handled:  false, // 不是时间类型，handled 应该为 false
		},
		{
			name:     "nil pointer",
			value:    (*ConfigWithUnexportedFields)(nil),
			expected: false,
			handled:  false,
		},
		{
			name:     "normal time pointer",
			value:    func() *time.Time { t := time.Now(); return &t }(),
			expected: false,
			handled:  true,
		},
		{
			name:     "nil time pointer",
			value:    (*time.Time)(nil),
			expected: true,
			handled:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 这个测试的关键是：不应该 panic
			isEmpty, handled := CheckEmptyTimePointer(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, isEmpty, "isEmpty mismatch")
			assert.Equal(t, test.handled, handled, "handled mismatch")
		})
	}
}

// TestCheckEmptyTimeStruct_UnexportedFields 测试修复后的代码能正确处理未导出字段的结构体
func TestCheckEmptyTimeStruct_UnexportedFields(t *testing.T) {
	// 定义一个包含未导出字段的结构体
	type StructWithUnexportedFields struct {
		PublicField  string
		privateField string // 未导出字段
	}

	tests := []struct {
		name     string
		value    any
		expected bool
		handled  bool
	}{
		{
			name:     "struct with unexported fields",
			value:    StructWithUnexportedFields{PublicField: "test", privateField: "private"},
			expected: false,
			handled:  false, // 不是时间类型，handled 应该为 false
		},
		{
			name:     "zero time.Time",
			value:    time.Time{},
			expected: true,
			handled:  true,
		},
		{
			name:     "non-zero time.Time",
			value:    time.Now(),
			expected: false,
			handled:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 这个测试的关键是：不应该 panic
			isEmpty, handled := CheckEmptyTimeStruct(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, isEmpty, "isEmpty mismatch")
			assert.Equal(t, test.handled, handled, "handled mismatch")
		})
	}
}

// TestIsEmptyValue_WithUnexportedFields 测试 IsEmptyValue 能正确处理包含未导出字段的结构体
// 这是对整个修复的集成测试
func TestIsEmptyValue_WithUnexportedFields(t *testing.T) {
	// 模拟 logger.ILogger 接口的场景
	type mockLogger struct {
		level        string // 未导出字段
		output       string // 未导出字段
		PublicConfig string // 导出字段
	}

	type ConfigWithLogger struct {
		Name   string
		Logger *mockLogger // 包含未导出字段的结构体指针
	}

	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{
			name: "config with nil logger",
			value: ConfigWithLogger{
				Name:   "test",
				Logger: nil,
			},
			expected: false, // Name 非空，所以整体非空
		},
		{
			name: "config with logger",
			value: ConfigWithLogger{
				Name: "test",
				Logger: &mockLogger{
					level:        "info",
					output:       "stdout",
					PublicConfig: "public",
				},
			},
			expected: false, // 有非空字段
		},
		{
			name: "empty config",
			value: ConfigWithLogger{
				Name:   "",
				Logger: nil,
			},
			expected: true, // 所有字段都为空
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 这个测试的关键是：不应该 panic
			result := IsEmptyValue(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestIsEmptyValue_StringValue 测试 StringValue wrapper
func TestIsEmptyValue_StringValue(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.StringValue
		expected bool
	}{
		{"nil StringValue", nil, true},
		{"empty StringValue", wrapperspb.String(""), true},
		{"non-empty StringValue", wrapperspb.String("test"), false},
		{"whitespace StringValue", wrapperspb.String("  "), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(tt.value))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsEmptyValue_Int32Value 测试 Int32Value wrapper
func TestIsEmptyValue_Int32Value(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.Int32Value
		expected bool
	}{
		{"nil Int32Value", nil, true},
		{"zero Int32Value", wrapperspb.Int32(0), false}, // 0 是有效值
		{"positive Int32Value", wrapperspb.Int32(42), false},
		{"negative Int32Value", wrapperspb.Int32(-1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(tt.value))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsEmptyValue_Int64Value 测试 Int64Value wrapper
func TestIsEmptyValue_Int64Value(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.Int64Value
		expected bool
	}{
		{"nil Int64Value", nil, true},
		{"zero Int64Value", wrapperspb.Int64(0), false}, // 0 是有效值
		{"positive Int64Value", wrapperspb.Int64(12345), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(tt.value))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsEmptyValue_BoolValue 测试 BoolValue wrapper
func TestIsEmptyValue_BoolValue(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.BoolValue
		expected bool
	}{
		{"nil BoolValue", nil, true},
		{"false BoolValue", wrapperspb.Bool(false), false}, // false 是有效值
		{"true BoolValue", wrapperspb.Bool(true), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(tt.value))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsEmptyValue_UInt32Value 测试 UInt32Value wrapper
func TestIsEmptyValue_UInt32Value(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.UInt32Value
		expected bool
	}{
		{"nil UInt32Value", nil, true},
		{"zero UInt32Value", wrapperspb.UInt32(0), false}, // 0 是有效值
		{"positive UInt32Value", wrapperspb.UInt32(100), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(tt.value))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsEmptyValue_FloatValue 测试 FloatValue wrapper
func TestIsEmptyValue_FloatValue(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.FloatValue
		expected bool
	}{
		{"nil FloatValue", nil, true},
		{"zero FloatValue", wrapperspb.Float(0.0), false}, // 0.0 是有效值
		{"positive FloatValue", wrapperspb.Float(3.14), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(tt.value))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsEmptyValue_DoubleValue 测试 DoubleValue wrapper
func TestIsEmptyValue_DoubleValue(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.DoubleValue
		expected bool
	}{
		{"nil DoubleValue", nil, true},
		{"zero DoubleValue", wrapperspb.Double(0.0), false}, // 0.0 是有效值
		{"positive DoubleValue", wrapperspb.Double(1234.5678), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(tt.value))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsEmptyValue_BytesValue 测试 BytesValue wrapper
func TestIsEmptyValue_BytesValue(t *testing.T) {
	tests := []struct {
		name     string
		value    *wrapperspb.BytesValue
		expected bool
	}{
		{"nil BytesValue", nil, true},
		{"empty BytesValue", wrapperspb.Bytes([]byte{}), true},
		{"non-empty BytesValue", wrapperspb.Bytes([]byte{1, 2, 3}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmptyValue(reflect.ValueOf(tt.value))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCheckEmptyWrapperPointer 测试 CheckEmptyWrapperPointer 函数
func TestCheckEmptyWrapperPointer(t *testing.T) {
	t.Run("StringValue", func(t *testing.T) {
		sv := wrapperspb.String("test")

		isEmpty, handled := CheckEmptyWrapperPointer(reflect.ValueOf(sv))

		assert.True(t, handled)
		assert.False(t, isEmpty)
	})

	t.Run("nil StringValue", func(t *testing.T) {
		var nilSV *wrapperspb.StringValue

		isEmpty, handled := CheckEmptyWrapperPointer(reflect.ValueOf(nilSV))

		assert.True(t, handled)
		assert.True(t, isEmpty)
	})

	t.Run("Int32Value with 0", func(t *testing.T) {
		iv := wrapperspb.Int32(0)

		isEmpty, handled := CheckEmptyWrapperPointer(reflect.ValueOf(iv))

		assert.True(t, handled)
		assert.False(t, isEmpty)
	})

	t.Run("normal int pointer", func(t *testing.T) {
		normalInt := 42

		isEmpty, handled := CheckEmptyWrapperPointer(reflect.ValueOf(&normalInt))

		assert.False(t, handled)
		assert.False(t, isEmpty)
	})
}

// TestIsTimeValid 测试 IsTimeValid 函数
func TestIsTimeValid(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"nil", nil, false},
		{"zero time.Time", time.Time{}, false},
		{"unix zero (1970-01-01)", time.Unix(0, 0), false},
		{"before unix zero", time.Unix(-1, 0), false},
		{"after unix zero", time.Unix(1, 0), true},
		{"now", time.Now(), true},
		{"specific date", time.Date(2025, 12, 16, 0, 0, 0, 0, time.UTC), true},
		{"nil *time.Time", (*time.Time)(nil), false},
		{"zero *time.Time", func() *time.Time { t := time.Time{}; return &t }(), false},
		{"unix zero *time.Time", func() *time.Time { t := time.Unix(0, 0); return &t }(), false},
		{"before unix zero *time.Time", func() *time.Time { t := time.Unix(-1, 0); return &t }(), false},
		{"valid *time.Time", func() *time.Time { t := time.Now(); return &t }(), true},
		{"non-time type string", "2025-12-16", true},
		{"non-time type int", 1234567890, true},
		{"non-time type bool", true, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsTimeValid(test.value)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestDerefValue 测试 DerefValue 函数
func TestDerefValue(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		deref, ok := DerefValue(nil)
		assert.False(t, ok)
		assert.Nil(t, deref)
	})

	t.Run("nil string pointer", func(t *testing.T) {
		var p *string
		deref, ok := DerefValue(p)
		assert.False(t, ok)
		assert.Nil(t, deref)
	})

	t.Run("valid string pointer", func(t *testing.T) {
		s := "hello"
		deref, ok := DerefValue(&s)
		assert.True(t, ok)
		assert.Equal(t, "hello", deref)
	})

	t.Run("nil int pointer", func(t *testing.T) {
		var p *int
		deref, ok := DerefValue(p)
		assert.False(t, ok)
		assert.Nil(t, deref)
	})

	t.Run("valid int pointer (non-zero)", func(t *testing.T) {
		i := 42
		deref, ok := DerefValue(&i)
		assert.True(t, ok)
		assert.Equal(t, 42, deref)
	})

	t.Run("valid int pointer (zero)", func(t *testing.T) {
		i := 0
		deref, ok := DerefValue(&i)
		assert.True(t, ok)
		assert.Equal(t, 0, deref)
	})

	t.Run("nil bool pointer", func(t *testing.T) {
		var p *bool
		deref, ok := DerefValue(p)
		assert.False(t, ok)
		assert.Nil(t, deref)
	})

	t.Run("valid bool pointer (false)", func(t *testing.T) {
		b := false
		deref, ok := DerefValue(&b)
		assert.True(t, ok)
		assert.Equal(t, false, deref)
	})

	t.Run("valid bool pointer (true)", func(t *testing.T) {
		b := true
		deref, ok := DerefValue(&b)
		assert.True(t, ok)
		assert.Equal(t, true, deref)
	})

	t.Run("non-pointer value", func(t *testing.T) {
		deref, ok := DerefValue(42)
		assert.True(t, ok)
		assert.Equal(t, 42, deref)
	})

	t.Run("non-pointer string", func(t *testing.T) {
		deref, ok := DerefValue("hello")
		assert.True(t, ok)
		assert.Equal(t, "hello", deref)
	})

	t.Run("non-pointer nil interface", func(t *testing.T) {
		var x interface{}
		deref, ok := DerefValue(x)
		assert.False(t, ok)
		assert.Nil(t, deref)
	})
}

// TestUnwrapProtobufWrapper 测试 protobuf wrapper 解包
func TestUnwrapProtobufWrapper(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expected  interface{}
		unwrapped bool
	}{
		{"string wrapper", wrapperspb.String("hello"), "hello", true},
		{"string wrapper value", wrapperspb.StringValue{Value: "hello"}, "hello", true},
		{"int32 wrapper", wrapperspb.Int32(42), int32(42), true},
		{"int64 wrapper", wrapperspb.Int64(42), int64(42), true},
		{"uint32 wrapper", wrapperspb.UInt32(42), uint32(42), true},
		{"uint64 wrapper", wrapperspb.UInt64(42), uint64(42), true},
		{"bool wrapper false", wrapperspb.Bool(false), false, true},
		{"float wrapper", wrapperspb.Float(3.14), float32(3.14), true},
		{"double wrapper", wrapperspb.Double(3.14), 3.14, true},
		{"bytes wrapper", wrapperspb.Bytes([]byte{1, 2}), []byte{1, 2}, true},
		{"non wrapper", "hello", nil, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, ok := UnwrapProtobufWrapper(test.value)
			assert.Equal(t, test.unwrapped, ok)
			assert.Equal(t, test.expected, actual)
		})
	}
}

// TestIsEmptyAfterDeref 测试 IsEmptyAfterDeref 函数
func TestIsEmptyAfterDeref(t *testing.T) {
	var nilStringValue *wrapperspb.StringValue
	var nilBoolValue *wrapperspb.BoolValue

	tests := []struct {
		name          string
		value         interface{}
		expectedEmpty bool
		expectedDeref interface{}
	}{
		// nil
		{"nil", nil, true, nil},

		// string
		{"empty string", "", true, nil},
		{"whitespace string", "  ", true, nil},
		{"null string", "null", false, "null"},
		{"undefined string", "undefined", false, "undefined"},
		{"valid string", "hello", false, "hello"},
		{"string with spaces", " hello ", false, " hello "},

		// string pointer
		{"nil *string", (*string)(nil), true, nil},
		{"empty *string", func() *string { s := ""; return &s }(), true, nil},
		{"valid *string", func() *string { s := "hello"; return &s }(), false, "hello"},

		// int (零值视为空，非零值视为非空)
		{"zero int", 0, true, nil},
		{"positive int", 42, false, 42},
		{"negative int", -1, false, -1},

		// int pointer
		{"nil *int", (*int)(nil), true, nil},
		{"zero *int", func() *int { i := 0; return &i }(), true, nil},
		{"valid *int", func() *int { i := 42; return &i }(), false, 42},

		// bool (false 和 true 都是有效值)
		{"false bool", false, false, false},
		{"true bool", true, false, true},

		// bool pointer
		{"nil *bool", (*bool)(nil), true, nil},
		{"false *bool", func() *bool { b := false; return &b }(), false, false},
		{"true *bool", func() *bool { b := true; return &b }(), false, true},

		// protobuf wrapper
		{"nil StringValue", nilStringValue, true, nil},
		{"empty StringValue", wrapperspb.String(""), true, nil},
		{"whitespace StringValue", wrapperspb.String("  "), true, nil},
		{"null StringValue", wrapperspb.String("null"), false, "null"},
		{name: "undefined StringValue", value: wrapperspb.String("undefined"), expectedEmpty: false, expectedDeref: "undefined"},
		{"valid StringValue", wrapperspb.String("hello"), false, "hello"},
		{"false BoolValue", wrapperspb.Bool(false), false, false},
		{"true BoolValue", wrapperspb.Bool(true), false, true},
		{"nil BoolValue", nilBoolValue, true, nil},
		{"zero Int32Value", wrapperspb.Int32(0), false, int32(0)},
		{"valid Int64Value", wrapperspb.Int64(42), false, int64(42)},
		{"zero FloatValue", wrapperspb.Float(0), false, float32(0)},
		{"zero DoubleValue", wrapperspb.Double(0), false, float64(0)},
		{"empty BytesValue", wrapperspb.Bytes([]byte{}), true, nil},
		{"valid BytesValue", wrapperspb.Bytes([]byte{1, 2}), false, []byte{1, 2}},

		// float
		{"zero float", 0.0, true, nil},
		{"valid float", 3.14, false, 3.14},

		// float pointer
		{"nil *float64", (*float64)(nil), true, nil},
		{"zero *float64", func() *float64 { f := 0.0; return &f }(), true, nil},
		{"valid *float64", func() *float64 { f := 3.14; return &f }(), false, 3.14},

		// slice
		{"nil slice", ([]int)(nil), true, nil},
		{"empty slice", []int{}, true, nil},
		{"non-empty slice", []int{1, 2}, false, []int{1, 2}},

		// map
		{"nil map", (map[string]int)(nil), true, nil},
		{"empty map", map[string]int{}, true, nil},
		{"non-empty map", map[string]int{"a": 1}, false, map[string]int{"a": 1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deref, empty := IsEmptyAfterDeref(test.value)
			assert.Equal(t, test.expectedEmpty, empty)
			assert.Equal(t, test.expectedDeref, deref)
		})
	}
}
