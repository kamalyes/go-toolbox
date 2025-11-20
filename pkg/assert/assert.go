/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-20 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 12:00:00
 * @FilePath: \go-toolbox\pkg\assert\assert.go
 * @Description: 业务断言库，提供运行时断言功能，支持自定义错误处理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package assert

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// AssertionError 断言错误类型
type AssertionError struct {
	Message string
	File    string
	Line    int
}

func (e *AssertionError) Error() string {
	if e.File != "" && e.Line > 0 {
		return fmt.Sprintf("Assertion failed at %s:%d: %s", e.File, e.Line, e.Message)
	}
	return fmt.Sprintf("Assertion failed: %s", e.Message)
}

// ErrorHandler 错误处理函数类型
type ErrorHandler func(err *AssertionError)

var (
	// DefaultHandler 默认错误处理器 - 直接 panic
	DefaultHandler ErrorHandler = func(err *AssertionError) {
		panic(err)
	}

	// GlobalHandler 全局错误处理器，可以自定义
	GlobalHandler ErrorHandler = DefaultHandler
)

// SetGlobalHandler 设置全局错误处理器
func SetGlobalHandler(handler ErrorHandler) {
	GlobalHandler = handler
}

// newError 创建断言错误
func newError(message string) *AssertionError {
	return &AssertionError{
		Message: message,
	}
}

// handleError 处理断言错误
func handleError(err *AssertionError) {
	if GlobalHandler != nil {
		GlobalHandler(err)
	}
}

// True 断言条件为真
// 如果 condition 为 false，触发断言错误
func True(condition bool, message ...string) {
	if !condition {
		msg := "Expected true, but got false"
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// False 断言条件为假
// 如果 condition 为 true，触发断言错误
func False(condition bool, message ...string) {
	if condition {
		msg := "Expected false, but got true"
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Equal 断言两个值相等
// 使用 reflect.DeepEqual 进行深度比较
func Equal[T comparable](expected, actual T, message ...string) {
	if expected != actual {
		msg := fmt.Sprintf("Expected %v, but got %v", expected, actual)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// NotEqual 断言两个值不相等
func NotEqual[T comparable](expected, actual T, message ...string) {
	if expected == actual {
		msg := fmt.Sprintf("Expected values to be different, but both are %v", expected)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// DeepEqual 断言两个值深度相等
// 适用于 slice、map、struct 等复杂类型
func DeepEqual(expected, actual interface{}, message ...string) {
	if !reflect.DeepEqual(expected, actual) {
		msg := fmt.Sprintf("Expected %+v, but got %+v", expected, actual)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// NotDeepEqual 断言两个值深度不相等
func NotDeepEqual(expected, actual interface{}, message ...string) {
	if reflect.DeepEqual(expected, actual) {
		msg := fmt.Sprintf("Expected values to be different, but both are %+v", expected)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Nil 断言值为 nil
func Nil(value interface{}, message ...string) {
	if value != nil {
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface ||
			rv.Kind() == reflect.Slice || rv.Kind() == reflect.Map ||
			rv.Kind() == reflect.Chan || rv.Kind() == reflect.Func {
			if !rv.IsNil() {
				msg := fmt.Sprintf("Expected nil, but got %v", value)
				if len(message) > 0 {
					msg = message[0]
				}
				handleError(newError(msg))
			}
		} else {
			msg := fmt.Sprintf("Expected nil, but got %v", value)
			if len(message) > 0 {
				msg = message[0]
			}
			handleError(newError(msg))
		}
	}
}

// NotNil 断言值不为 nil
func NotNil(value interface{}, message ...string) {
	if value == nil {
		msg := "Expected non-nil value, but got nil"
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
		return
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface ||
		rv.Kind() == reflect.Slice || rv.Kind() == reflect.Map ||
		rv.Kind() == reflect.Chan || rv.Kind() == reflect.Func {
		if rv.IsNil() {
			msg := "Expected non-nil value, but got nil"
			if len(message) > 0 {
				msg = message[0]
			}
			handleError(newError(msg))
		}
	}
}

// Empty 断言容器为空（string、slice、map、array、channel）
func Empty(value interface{}, message ...string) {
	if !isEmpty(value) {
		msg := fmt.Sprintf("Expected empty, but got %v", value)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// NotEmpty 断言容器不为空
func NotEmpty(value interface{}, message ...string) {
	if isEmpty(value) {
		msg := "Expected non-empty value, but got empty"
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// isEmpty 检查值是否为空
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String:
		return rv.String() == ""
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return true
		}
		return isEmpty(rv.Elem().Interface())
	default:
		return false
	}
}

// Zero 断言值为零值
func Zero[T comparable](value T, message ...string) {
	var zero T
	if value != zero {
		msg := fmt.Sprintf("Expected zero value, but got %v", value)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// NotZero 断言值不为零值
func NotZero[T comparable](value T, message ...string) {
	var zero T
	if value == zero {
		msg := fmt.Sprintf("Expected non-zero value, but got %v", value)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Greater 断言 a > b（支持数值类型）
func Greater[T int | int64 | float32 | float64](a, b T, message ...string) {
	if a <= b {
		msg := fmt.Sprintf("Expected %v > %v", a, b)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// GreaterOrEqual 断言 a >= b
func GreaterOrEqual[T int | int64 | float32 | float64](a, b T, message ...string) {
	if a < b {
		msg := fmt.Sprintf("Expected %v >= %v", a, b)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Less 断言 a < b
func Less[T int | int64 | float32 | float64](a, b T, message ...string) {
	if a >= b {
		msg := fmt.Sprintf("Expected %v < %v", a, b)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// LessOrEqual 断言 a <= b
func LessOrEqual[T int | int64 | float32 | float64](a, b T, message ...string) {
	if a > b {
		msg := fmt.Sprintf("Expected %v <= %v", a, b)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Contains 断言字符串包含子串
func Contains(str, substr string, message ...string) {
	if !strings.Contains(str, substr) {
		msg := fmt.Sprintf("Expected string '%s' to contain '%s'", str, substr)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// NotContains 断言字符串不包含子串
func NotContains(str, substr string, message ...string) {
	if strings.Contains(str, substr) {
		msg := fmt.Sprintf("Expected string '%s' to not contain '%s'", str, substr)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// HasPrefix 断言字符串有指定前缀
func HasPrefix(str, prefix string, message ...string) {
	if !strings.HasPrefix(str, prefix) {
		msg := fmt.Sprintf("Expected string '%s' to have prefix '%s'", str, prefix)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// HasSuffix 断言字符串有指定后缀
func HasSuffix(str, suffix string, message ...string) {
	if !strings.HasSuffix(str, suffix) {
		msg := fmt.Sprintf("Expected string '%s' to have suffix '%s'", str, suffix)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// InSlice 断言值在切片中
func InSlice[T comparable](value T, slice []T, message ...string) {
	for _, item := range slice {
		if item == value {
			return
		}
	}
	msg := fmt.Sprintf("Expected value %v to be in slice %v", value, slice)
	if len(message) > 0 {
		msg = message[0]
	}
	handleError(newError(msg))
}

// NotInSlice 断言值不在切片中
func NotInSlice[T comparable](value T, slice []T, message ...string) {
	for _, item := range slice {
		if item == value {
			msg := fmt.Sprintf("Expected value %v to not be in slice %v", value, slice)
			if len(message) > 0 {
				msg = message[0]
			}
			handleError(newError(msg))
			return
		}
	}
}

// Error 断言错误不为 nil
func Error(err error, message ...string) {
	if err == nil {
		msg := "Expected error, but got nil"
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// NoError 断言错误为 nil
func NoError(err error, message ...string) {
	if err != nil {
		msg := fmt.Sprintf("Expected no error, but got: %v", err)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// ErrorIs 断言错误是指定类型
func ErrorIs(err, target error, message ...string) {
	if !errors.Is(err, target) {
		msg := fmt.Sprintf("Expected error to be %v, but got %v", target, err)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Panic 断言函数会 panic
func Panic(fn func(), message ...string) {
	defer func() {
		if r := recover(); r == nil {
			msg := "Expected function to panic, but it didn't"
			if len(message) > 0 {
				msg = message[0]
			}
			handleError(newError(msg))
		}
	}()
	fn()
}

// NotPanic 断言函数不会 panic
func NotPanic(fn func(), message ...string) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Expected function to not panic, but it panicked with: %v", r)
			if len(message) > 0 {
				msg = message[0]
			}
			handleError(newError(msg))
		}
	}()
	fn()
}

// InRange 断言数值在范围内 [min, max]
func InRange[T int | int64 | float32 | float64](value, min, max T, message ...string) {
	if value < min || value > max {
		msg := fmt.Sprintf("Expected %v to be in range [%v, %v]", value, min, max)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// NotInRange 断言数值不在范围内
func NotInRange[T int | int64 | float32 | float64](value, min, max T, message ...string) {
	if value >= min && value <= max {
		msg := fmt.Sprintf("Expected %v to not be in range [%v, %v]", value, min, max)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Length 断言容器长度
func Length(value interface{}, expectedLen int, message ...string) {
	rv := reflect.ValueOf(value)
	var actualLen int

	switch rv.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		actualLen = rv.Len()
	default:
		msg := fmt.Sprintf("Value of type %T does not have a length", value)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
		return
	}

	if actualLen != expectedLen {
		msg := fmt.Sprintf("Expected length %d, but got %d", expectedLen, actualLen)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Type 断言值的类型
func Type(value interface{}, expectedType interface{}, message ...string) {
	expectedT := reflect.TypeOf(expectedType)
	actualT := reflect.TypeOf(value)

	if actualT != expectedT {
		msg := fmt.Sprintf("Expected type %v, but got %v", expectedT, actualT)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}

// Implements 断言类型实现了指定接口
func Implements(value interface{}, interfaceType interface{}, message ...string) {
	valueType := reflect.TypeOf(value)
	interfaceT := reflect.TypeOf(interfaceType).Elem()

	if !valueType.Implements(interfaceT) {
		msg := fmt.Sprintf("Expected type %v to implement %v", valueType, interfaceT)
		if len(message) > 0 {
			msg = message[0]
		}
		handleError(newError(msg))
	}
}
