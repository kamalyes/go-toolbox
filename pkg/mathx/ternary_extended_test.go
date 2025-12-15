/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-20 11:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 21:27:25
 * @FilePath: \go-toolbox\pkg\mathx\ternary_extended_test.go
 * @Description: 扩展三元运算符功能测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIfNotNil(t *testing.T) {
	// 测试指针不为空的情况
	value := 42
	result := IfNotNil(&value, 0)
	assert.Equal(t, 42, result)

	// 测试指针为空的情况
	var nilPtr *int
	result = IfNotNil(nilPtr, 100)
	assert.Equal(t, 100, result)
}

func TestIfNotEmpty(t *testing.T) {
	// 测试非空字符串
	result := IfNotEmpty("hello", "default")
	assert.Equal(t, "hello", result)

	// 测试空字符串
	result = IfNotEmpty("", "default")
	assert.Equal(t, "default", result)
}

func TestIfNotZero(t *testing.T) {
	// 测试非零值
	result := IfNotZero(42, 100)
	assert.Equal(t, 42, result)

	// 测试零值
	result = IfNotZero(0, 100)
	assert.Equal(t, 100, result)

	// 测试字符串零值
	strResult := IfNotZero("", "default")
	assert.Equal(t, "default", strResult)
}

func TestIfContains(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	// 测试包含的情况
	result := IfContains(slice, "banana", "found", "not found")
	assert.Equal(t, "found", result)

	// 测试不包含的情况
	result = IfContains(slice, "grape", "found", "not found")
	assert.Equal(t, "not found", result)
}

func TestIfAny(t *testing.T) {
	// 测试有条件为真的情况
	conditions := []bool{false, false, true, false}
	result := IfAny(conditions, "success", "failure")
	assert.Equal(t, "success", result)

	// 测试所有条件为假的情况
	conditions = []bool{false, false, false}
	result = IfAny(conditions, "success", "failure")
	assert.Equal(t, "failure", result)
}

func TestIfAll(t *testing.T) {
	// 测试所有条件为真的情况
	conditions := []bool{true, true, true}
	result := IfAll(conditions, "success", "failure")
	assert.Equal(t, "success", result)

	// 测试有条件为假的情况
	conditions = []bool{true, false, true}
	result = IfAll(conditions, "success", "failure")
	assert.Equal(t, "failure", result)

	// 测试空条件数组
	conditions = []bool{}
	result = IfAll(conditions, "success", "failure")
	assert.Equal(t, "failure", result)
}

func TestIfCount(t *testing.T) {
	conditions := []bool{true, false, true, true, false}

	// 测试达到阈值的情况
	result := IfCount(conditions, 3, "enough", "not enough")
	assert.Equal(t, "enough", result)

	// 测试未达到阈值的情况
	result = IfCount(conditions, 4, "enough", "not enough")
	assert.Equal(t, "not enough", result)
}

func TestIfMap(t *testing.T) {
	// 测试条件为真时的映射
	result := IfMap(true, "hello", strings.ToUpper, "default")
	assert.Equal(t, "HELLO", result)

	// 测试条件为假时返回默认值
	result = IfMap(false, "hello", strings.ToUpper, "default")
	assert.Equal(t, "default", result)
}

func TestIfMapElse(t *testing.T) {
	toUpper := func(s string) string { return strings.ToUpper(s) }
	toLower := func(s string) string { return strings.ToLower(s) }

	// 测试条件为真
	result := IfMapElse(true, "Hello", toUpper, toLower)
	assert.Equal(t, "HELLO", result)

	// 测试条件为假
	result = IfMapElse(false, "Hello", toUpper, toLower)
	assert.Equal(t, "hello", result)
}

func TestIfFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}
	isEven := func(n int) bool { return n%2 == 0 }

	// 测试使用过滤
	result := IfFilter(true, numbers, isEven)
	expected := []int{2, 4, 6}
	assert.Equal(t, expected, result)

	// 测试不使用过滤
	result = IfFilter(false, numbers, isEven)
	assert.Equal(t, numbers, result)
}

func TestIfValidate(t *testing.T) {
	isPositive := func(n int) bool { return n > 0 }

	// 测试验证通过
	result := IfValidate(5, isPositive, "valid", "invalid")
	assert.Equal(t, "valid", result)

	// 测试验证失败
	result = IfValidate(-5, isPositive, "valid", "invalid")
	assert.Equal(t, "invalid", result)
}

func TestIfCast(t *testing.T) {
	var value interface{} = "hello"

	// 测试成功转换
	result := IfCast[string](value, "default")
	assert.Equal(t, "hello", result)

	// 测试转换失败
	result = IfCast[string](123, "default")
	assert.Equal(t, "default", result)
}

func TestIfBetween(t *testing.T) {
	// 测试在范围内
	result := IfBetween(5, 1, 10, 100, 200)
	assert.Equal(t, 100, result)

	// 测试超出范围
	result = IfBetween(15, 1, 10, 100, 200)
	assert.Equal(t, 200, result)

	// 测试边界值
	result = IfBetween(1, 1, 10, 100, 200)
	assert.Equal(t, 100, result)
}

func TestIfSwitch(t *testing.T) {
	cases := map[string]string{
		"red":    "stop",
		"green":  "go",
		"yellow": "caution",
	}

	// 测试匹配的情况
	result := IfSwitch("red", cases, "unknown")
	assert.Equal(t, "stop", result)

	// 测试不匹配的情况
	result = IfSwitch("blue", cases, "unknown")
	assert.Equal(t, "unknown", result)
}

func TestIfTryParse(t *testing.T) {
	parser := func(s string) (int, error) {
		return strconv.Atoi(s)
	}

	// 测试解析成功
	result := IfTryParse("123", parser, 0)
	assert.Equal(t, 123, result)

	// 测试解析失败
	result = IfTryParse("abc", parser, 0)
	assert.Equal(t, 0, result)
}

func TestIfSafeIndex(t *testing.T) {
	slice := []string{"a", "b", "c"}

	// 测试有效索引
	result := IfSafeIndex(slice, 1, "default")
	assert.Equal(t, "b", result)

	// 测试无效索引
	result = IfSafeIndex(slice, 5, "default")
	assert.Equal(t, "default", result)

	// 测试负索引
	result = IfSafeIndex(slice, -1, "default")
	assert.Equal(t, "default", result)
}

func TestIfSafeKey(t *testing.T) {
	m := map[string]int{
		"apple":  1,
		"banana": 2,
		"orange": 3,
	}

	// 测试存在的键
	result := IfSafeKey(m, "banana", 0)
	assert.Equal(t, 2, result)

	// 测试不存在的键
	result = IfSafeKey(m, "grape", 0)
	assert.Equal(t, 0, result)
}

func TestIfMulti(t *testing.T) {
	values := []string{"red", "green", "blue"}

	// 测试匹配的情况
	result := IfMulti("green", values, "primary", "not primary")
	assert.Equal(t, "primary", result)

	// 测试不匹配的情况
	result = IfMulti("purple", values, "primary", "not primary")
	assert.Equal(t, "not primary", result)
}

func TestIfPipeline(t *testing.T) {
	funcs := []func(string) string{
		strings.ToUpper,
		func(s string) string { return s + "!" },
		func(s string) string { return ">>> " + s },
	}

	// 测试条件为真时的管道处理
	result := IfPipeline(true, "hello", funcs, "default")
	expected := ">>> HELLO!"
	assert.Equal(t, expected, result)

	// 测试条件为假时返回默认值
	result = IfPipeline(false, "hello", funcs, "default")
	assert.Equal(t, "default", result)
}

func TestIfMemoized(t *testing.T) {
	cache := make(map[string]string)
	callCount := 0

	computeFn := func() string {
		callCount++
		return "computed result"
	}

	// 第一次调用
	result := IfMemoized(true, "key1", cache, computeFn, "default")
	assert.Equal(t, "computed result", result)
	assert.Equal(t, 1, callCount)

	// 第二次调用相同key，应该使用缓存
	result = IfMemoized(true, "key1", cache, computeFn, "default")
	assert.Equal(t, "computed result", result)
	assert.Equal(t, 1, callCount)

	// 条件为假时返回默认值
	result = IfMemoized(false, "key2", cache, computeFn, "default")
	assert.Equal(t, "default", result)
	assert.Equal(t, 1, callCount)
}

// 基准测试
func BenchmarkIfNotNil(b *testing.B) {
	value := 42
	for i := 0; i < b.N; i++ {
		IfNotNil(&value, 0)
	}
}

func BenchmarkIfContains(b *testing.B) {
	slice := []string{"apple", "banana", "orange", "grape", "watermelon"}
	for i := 0; i < b.N; i++ {
		IfContains(slice, "banana", "found", "not found")
	}
}

func BenchmarkIfDoAF(b *testing.B) {
	trueFn := func() string { return "true result" }
	falseFn := func() string { return "false result" }

	for i := 0; i < b.N; i++ {
		IfDoAF(i%2 == 0, trueFn, falseFn)
	}
}

// TestIfStrFmt 测试条件格式化字符串
func TestIfStrFmt(t *testing.T) {
	tests := []struct {
		name         string
		condition    bool
		trueFormat   string
		trueArgs     []any
		falseFormat  string
		falseArgs    []any
		expectFormat string
		expectArgs   []any
	}{
		{
			name:         "true condition",
			condition:    true,
			trueFormat:   "Success: %s",
			trueArgs:     []any{"operation completed"},
			falseFormat:  "Error: %v",
			falseArgs:    []any{errors.New("failed")},
			expectFormat: "Success: %s",
			expectArgs:   []any{"operation completed"},
		},
		{
			name:         "false condition",
			condition:    false,
			trueFormat:   "Success: %s",
			trueArgs:     []any{"operation completed"},
			falseFormat:  "Error: %v",
			falseArgs:    []any{errors.New("failed")},
			expectFormat: "Error: %v",
			expectArgs:   []any{errors.New("failed")},
		},
		{
			name:         "empty args",
			condition:    true,
			trueFormat:   "No args",
			trueArgs:     []any{},
			falseFormat:  "Also no args",
			falseArgs:    []any{},
			expectFormat: "No args",
			expectArgs:   []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, args := IfStrFmt(tt.condition, tt.trueFormat, tt.trueArgs, tt.falseFormat, tt.falseArgs)
			assert.Equal(t, tt.expectFormat, format)
			assert.Equal(t, tt.expectArgs, args)

			// 验证格式化后的字符串
			if len(args) > 0 {
				formatted := fmt.Sprintf(format, args...)
				expectedFormatted := fmt.Sprintf(tt.expectFormat, tt.expectArgs...)
				assert.Equal(t, expectedFormatted, formatted)
			}
		})
	}
}

// TestIfEmptySlice 测试空切片判断
func TestIfEmptySlice(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		trueVal  string
		falseVal string
		expected string
	}{
		{
			name:     "empty slice",
			slice:    []int{},
			trueVal:  "empty",
			falseVal: "not empty",
			expected: "empty",
		},
		{
			name:     "nil slice",
			slice:    nil,
			trueVal:  "empty",
			falseVal: "not empty",
			expected: "empty",
		},
		{
			name:     "non-empty slice",
			slice:    []int{1, 2, 3},
			trueVal:  "empty",
			falseVal: "not empty",
			expected: "not empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfEmptySlice(tt.slice, tt.trueVal, tt.falseVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIfEmptySlice_DifferentTypes 测试不同切片类型
func TestIfEmptySliceDifferentTypes(t *testing.T) {
	t.Run("string slice", func(t *testing.T) {
		result := IfEmptySlice([]string{}, "empty", "not empty")
		assert.Equal(t, "empty", result)

		result = IfEmptySlice([]string{"a", "b"}, "empty", "not empty")
		assert.Equal(t, "not empty", result)
	})

	t.Run("struct slice", func(t *testing.T) {
		type TestStruct struct {
			Value string
		}

		result := IfEmptySlice([]TestStruct{}, 0, 1)
		assert.Equal(t, 0, result)

		result = IfEmptySlice([]TestStruct{{Value: "test"}}, 0, 1)
		assert.Equal(t, 1, result)
	})
}

// TestIfLenGt 测试长度大于比较
func TestIfLenGt(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		trueVal  string
		falseVal string
		expected string
	}{
		{
			name:     "length greater than n",
			slice:    []int{1, 2, 3, 4},
			n:        3,
			trueVal:  "greater",
			falseVal: "not greater",
			expected: "greater",
		},
		{
			name:     "length equal to n",
			slice:    []int{1, 2, 3},
			n:        3,
			trueVal:  "greater",
			falseVal: "not greater",
			expected: "not greater",
		},
		{
			name:     "length less than n",
			slice:    []int{1, 2},
			n:        3,
			trueVal:  "greater",
			falseVal: "not greater",
			expected: "not greater",
		},
		{
			name:     "empty slice",
			slice:    []int{},
			n:        0,
			trueVal:  "greater",
			falseVal: "not greater",
			expected: "not greater",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfLenGt(tt.slice, tt.n, tt.trueVal, tt.falseVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIfLenEq 测试长度等于比较
func TestIfLenEq(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		n        int
		trueVal  int
		falseVal int
		expected int
	}{
		{
			name:     "length equals n",
			slice:    []string{"a", "b", "c"},
			n:        3,
			trueVal:  1,
			falseVal: 0,
			expected: 1,
		},
		{
			name:     "length not equals n",
			slice:    []string{"a", "b"},
			n:        3,
			trueVal:  1,
			falseVal: 0,
			expected: 0,
		},
		{
			name:     "empty slice with n=0",
			slice:    []string{},
			n:        0,
			trueVal:  1,
			falseVal: 0,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfLenEq(tt.slice, tt.n, tt.trueVal, tt.falseVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIfErrOrNil 测试错误或 nil 返回
func TestIfErrOrNil(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		val          int
		trueVal      string
		falseVal     string
		expectResult string
	}{
		{
			name:         "error case - return trueVal",
			err:          errors.New("test error"),
			val:          100,
			trueVal:      "error",
			falseVal:     "ok",
			expectResult: "error",
		},
		{
			name:         "success case with non-zero value - return falseVal",
			err:          nil,
			val:          100,
			trueVal:      "error",
			falseVal:     "ok",
			expectResult: "ok",
		},
		{
			name:         "success case but zero value - return trueVal",
			err:          nil,
			val:          0,
			trueVal:      "error",
			falseVal:     "ok",
			expectResult: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfErrOrNil(tt.val, tt.err, tt.trueVal, tt.falseVal)
			assert.Equal(t, tt.expectResult, result)
		})
	}
}

// TestIfCountGt 测试计数大于判断
func TestIfCountGt(t *testing.T) {
	tests := []struct {
		name     string
		count    int64
		n        int64
		trueVal  bool
		falseVal bool
		expected bool
	}{
		{
			name:     "count greater than n",
			count:    10,
			n:        5,
			trueVal:  true,
			falseVal: false,
			expected: true,
		},
		{
			name:     "count equal to n",
			count:    5,
			n:        5,
			trueVal:  true,
			falseVal: false,
			expected: false,
		},
		{
			name:     "count less than n",
			count:    3,
			n:        5,
			trueVal:  true,
			falseVal: false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfCountGt(tt.count, tt.n, tt.trueVal, tt.falseVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIfExecElseVariadicParams 测试 IfExecElse 的可变参数功能
func TestIfExecElseVariadicParams(t *testing.T) {
	// 测试完整两个参数版本
	t.Run("两个参数都提供", func(t *testing.T) {
		trueExecuted := false
		falseExecuted := false

		IfExecElse(true,
			func() { trueExecuted = true },
			func() { falseExecuted = true },
		)

		assert.True(t, trueExecuted)
		assert.False(t, falseExecuted)
	})

	t.Run("条件为false时执行第二个回调", func(t *testing.T) {
		trueExecuted := false
		falseExecuted := false

		IfExecElse(false,
			func() { trueExecuted = true },
			func() { falseExecuted = true },
		)

		assert.False(t, trueExecuted)
		assert.True(t, falseExecuted)
	})

	// 测试可变参数：省略第二个参数
	t.Run("省略第二个参数", func(t *testing.T) {
		trueExecuted := false

		IfExecElse(true, func() { trueExecuted = true })

		assert.True(t, trueExecuted)
	})

	t.Run("省略第二个参数且条件为false", func(t *testing.T) {
		trueExecuted := false

		// 不应该执行任何操作
		IfExecElse(false, func() { trueExecuted = true })

		assert.False(t, trueExecuted)
	})

	// 测试 nil 处理
	t.Run("第一个参数为nil", func(t *testing.T) {
		falseExecuted := false

		IfExecElse(false, nil, func() { falseExecuted = true })

		assert.True(t, falseExecuted)
	})

	t.Run("第二个参数为nil", func(t *testing.T) {
		trueExecuted := false

		IfExecElse(true,
			func() { trueExecuted = true },
			nil,
		)

		assert.True(t, trueExecuted)
	})
}

// TestIfCallVariadicParams 测试 IfCall 的可变参数功能
func TestIfCallVariadicParams(t *testing.T) {
	// 测试完整版本（两个回调）
	t.Run("两个回调都提供", func(t *testing.T) {
		onTrueCalled := false
		onFalseCalled := false

		IfCall(true, 42, nil,
			func(r int, e error) { onTrueCalled = true },
			func(r int, e error) { onFalseCalled = true },
		)

		assert.True(t, onTrueCalled)
		assert.False(t, onFalseCalled)
	})

	t.Run("条件为false时调用第二个回调", func(t *testing.T) {
		onTrueCalled := false
		onFalseCalled := false

		IfCall(false, 42, nil,
			func(r int, e error) { onTrueCalled = true },
			func(r int, e error) { onFalseCalled = true },
		)

		assert.False(t, onTrueCalled)
		assert.True(t, onFalseCalled)
	})

	// 测试只提供一个回调（true 分支）
	t.Run("只提供true分支回调", func(t *testing.T) {
		result := ""

		IfCall(true, "success", nil,
			func(r string, e error) { result = r },
		)

		assert.Equal(t, "success", result)
	})

	t.Run("只提供true分支但条件为false", func(t *testing.T) {
		called := false

		IfCall(false, "data", nil,
			func(r string, e error) { called = true },
		)

		assert.False(t, called)
	})

	// 测试不提供回调
	t.Run("不提供任何回调", func(t *testing.T) {
		// 不应该 panic
		IfCall(true, 123, nil)
		IfCall(false, 123, nil)
	})

	// 测试 nil 回调
	t.Run("第一个回调为nil", func(t *testing.T) {
		onFalseCalled := false

		IfCall(false, 0, nil,
			nil,
			func(r int, e error) { onFalseCalled = true },
		)

		assert.True(t, onFalseCalled)
	})

	// 测试错误处理
	t.Run("处理错误情况", func(t *testing.T) {
		var capturedErr error

		IfCall(true, "", assert.AnError,
			func(r string, e error) { capturedErr = e },
		)

		assert.Equal(t, assert.AnError, capturedErr)
	})
}

// TestIfMapElseVariadicParams 测试 IfMapElse 的可变参数功能
func TestIfMapElseVariadicParams(t *testing.T) {
	// 测试完整版本
	t.Run("两个mapper都提供", func(t *testing.T) {
		result := IfMapElse(true, 10,
			func(x int) string { return "positive" },
			func(x int) string { return "negative" },
		)
		assert.Equal(t, "positive", result)
	})

	t.Run("条件为false时使用第二个mapper", func(t *testing.T) {
		result := IfMapElse(false, 10,
			func(x int) string { return "positive" },
			func(x int) string { return "negative" },
		)
		assert.Equal(t, "negative", result)
	})

	// 测试省略第二个mapper（返回零值）
	t.Run("省略第二个mapper返回零值", func(t *testing.T) {
		result := IfMapElse(false, 10,
			func(x int) string { return "mapped" },
		)
		assert.Equal(t, "", result) // string 零值
	})

	t.Run("省略第二个mapper但条件为true", func(t *testing.T) {
		result := IfMapElse(true, 100,
			func(x int) int { return x * 2 },
		)
		assert.Equal(t, 200, result)
	})

	// 测试复杂类型转换
	t.Run("复杂类型转换", func(t *testing.T) {
		type Data struct{ Value int }

		result := IfMapElse(true, Data{Value: 42},
			func(d Data) string { return "JSON: " + string(rune(d.Value)) },
			func(d Data) string { return "XML: " + string(rune(d.Value)) },
		)

		assert.NotEmpty(t, result)
	})

	// 测试 nil mapper
	t.Run("第二个mapper为nil", func(t *testing.T) {
		result := IfMapElse(false, 5,
			func(x int) string { return "value" },
			nil,
		)
		assert.Equal(t, "", result)
	})
}

// TestIfDoAsyncVariadicParams 测试 IfDoAsync 的可变参数功能
func TestIfDoAsyncVariadicParams(t *testing.T) {
	// 测试提供默认值
	t.Run("提供默认值且条件为true", func(t *testing.T) {
		ch := IfDoAsync(true,
			func() int { return 100 },
			999,
		)
		result := <-ch
		assert.Equal(t, 100, result)
	})

	t.Run("提供默认值且条件为false", func(t *testing.T) {
		ch := IfDoAsync(false,
			func() int { return 100 },
			999,
		)
		result := <-ch
		assert.Equal(t, 999, result)
	})

	// 测试不提供默认值（返回零值）
	t.Run("不提供默认值且条件为false", func(t *testing.T) {
		ch := IfDoAsync(false,
			func() int { return 100 },
		)
		result := <-ch
		assert.Equal(t, 0, result) // int 零值
	})

	t.Run("不提供默认值且条件为true", func(t *testing.T) {
		ch := IfDoAsync(true,
			func() string { return "executed" },
		)
		result := <-ch
		assert.Equal(t, "executed", result)
	})

	// 测试复杂类型的零值
	t.Run("复杂类型的零值", func(t *testing.T) {
		type CustomType struct {
			Name  string
			Value int
		}

		ch := IfDoAsync(false,
			func() CustomType { return CustomType{"test", 42} },
		)
		result := <-ch
		assert.Equal(t, CustomType{}, result) // 零值
	})
}

// TestIfDoAsyncWithTimeoutVariadicParams 测试 IfDoAsyncWithTimeout 的可变参数功能
func TestIfDoAsyncWithTimeoutVariadicParams(t *testing.T) {
	// 测试提供默认值
	t.Run("提供默认值且条件为true", func(t *testing.T) {
		ch := IfDoAsyncWithTimeout(true,
			func() int { return 200 },
			100,
			999,
		)
		result := <-ch
		assert.Equal(t, 200, result)
	})

	t.Run("提供默认值且条件为false", func(t *testing.T) {
		ch := IfDoAsyncWithTimeout(false,
			func() int { return 200 },
			100,
			999,
		)
		result := <-ch
		assert.Equal(t, 999, result)
	})

	// 测试不提供默认值
	t.Run("不提供默认值且条件为false", func(t *testing.T) {
		ch := IfDoAsyncWithTimeout(false,
			func() string { return "data" },
			100,
		)
		result := <-ch
		assert.Equal(t, "", result) // string 零值
	})

	// 测试超时场景
	t.Run("超时返回零值", func(t *testing.T) {
		ch := IfDoAsyncWithTimeout(true,
			func() int {
				time.Sleep(200 * time.Millisecond)
				return 100
			},
			50, // 50ms 超时
		)
		result := <-ch
		assert.Equal(t, 0, result) // 超时返回零值
	})

	t.Run("超时返回零值_提供默认值也被忽略", func(t *testing.T) {
		// 注意：超时时返回零值，而不是defaultVal
		ch := IfDoAsyncWithTimeout(true,
			func() int {
				time.Sleep(200 * time.Millisecond)
				return 100
			},
			50,  // 50ms 超时
			999, // 这个默认值只在 condition=false 时使用
		)
		result := <-ch
		assert.Equal(t, 0, result) // 超时时返回零值
	})

	// 测试正常完成（不超时）
	t.Run("正常完成不超时", func(t *testing.T) {
		ch := IfDoAsyncWithTimeout(true,
			func() int {
				time.Sleep(10 * time.Millisecond)
				return 42
			},
			100,
		)
		result := <-ch
		assert.Equal(t, 42, result)
	})
}
