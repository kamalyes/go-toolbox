/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-20 12:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 12:30:00
 * @FilePath: \go-toolbox\tests\assert_test.go
 * @Description: 业务断言库测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"github.com/kamalyes/go-toolbox/pkg/assert"
	"testing"
)

func TestAssertTrue(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.True should not panic when condition is true")
		}
	}()
	assert.True(true, "This should pass")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.True should panic when condition is false")
		}
	}()
	assert.True(false, "This should fail")
}

func TestAssertFalse(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.False should not panic when condition is false")
		}
	}()
	assert.False(false, "This should pass")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.False should panic when condition is true")
		}
	}()
	assert.False(true, "This should fail")
}

func TestAssertEqual(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Equal should not panic when values are equal")
		}
	}()
	assert.Equal(42, 42, "Numbers should be equal")
	assert.Equal("hello", "hello", "Strings should be equal")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Equal should panic when values are not equal")
		}
	}()
	assert.Equal(42, 24, "This should fail")
}

func TestAssertNotEqual(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.NotEqual should not panic when values are different")
		}
	}()
	assert.NotEqual(42, 24, "Numbers should be different")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.NotEqual should panic when values are equal")
		}
	}()
	assert.NotEqual(42, 42, "This should fail")
}

func TestAssertNil(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Nil should not panic when value is nil")
		}
	}()
	var ptr *int
	assert.Nil(ptr, "Pointer should be nil")
	assert.Nil(nil, "Nil should be nil")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Nil should panic when value is not nil")
		}
	}()
	value := 42
	assert.Nil(&value, "This should fail")
}

func TestAssertNotNil(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.NotNil should not panic when value is not nil")
		}
	}()
	value := 42
	assert.NotNil(&value, "Pointer should not be nil")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.NotNil should panic when value is nil")
		}
	}()
	var ptr *int
	assert.NotNil(ptr, "This should fail")
}

func TestAssertEmpty(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Empty should not panic when value is empty")
		}
	}()
	assert.Empty("", "Empty string should be empty")
	assert.Empty([]int{}, "Empty slice should be empty")
	assert.Empty(map[string]int{}, "Empty map should be empty")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Empty should panic when value is not empty")
		}
	}()
	assert.Empty("hello", "This should fail")
}

func TestAssertNotEmpty(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.NotEmpty should not panic when value is not empty")
		}
	}()
	assert.NotEmpty("hello", "Non-empty string should not be empty")
	assert.NotEmpty([]int{1, 2, 3}, "Non-empty slice should not be empty")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.NotEmpty should panic when value is empty")
		}
	}()
	assert.NotEmpty("", "This should fail")
}

func TestAssertZero(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Zero should not panic when value is zero")
		}
	}()
	assert.Zero(0, "Zero should be zero")
	assert.Zero("", "Empty string should be zero")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Zero should panic when value is not zero")
		}
	}()
	assert.Zero(42, "This should fail")
}

func TestAssertNotZero(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.NotZero should not panic when value is not zero")
		}
	}()
	assert.NotZero(42, "Non-zero should not be zero")
	assert.NotZero("hello", "Non-empty string should not be zero")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.NotZero should panic when value is zero")
		}
	}()
	assert.NotZero(0, "This should fail")
}

func TestAssertGreater(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Greater should not panic when a > b")
		}
	}()
	assert.Greater(5, 3, "5 should be greater than 3")
	assert.Greater(10.5, 10.1, "10.5 should be greater than 10.1")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Greater should panic when a <= b")
		}
	}()
	assert.Greater(3, 5, "This should fail")
}

func TestAssertLess(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Less should not panic when a < b")
		}
	}()
	assert.Less(3, 5, "3 should be less than 5")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Less should panic when a >= b")
		}
	}()
	assert.Less(5, 3, "This should fail")
}

func TestAssertContains(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Contains should not panic when string contains substring")
		}
	}()
	assert.Contains("hello world", "world", "String should contain substring")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Contains should panic when string does not contain substring")
		}
	}()
	assert.Contains("hello", "world", "This should fail")
}

func TestAssertNotContains(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.NotContains should not panic when string does not contain substring")
		}
	}()
	assert.NotContains("hello", "world", "String should not contain substring")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.NotContains should panic when string contains substring")
		}
	}()
	assert.NotContains("hello world", "world", "This should fail")
}

func TestAssertHasPrefix(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.HasPrefix should not panic when string has prefix")
		}
	}()
	assert.HasPrefix("hello world", "hello", "String should have prefix")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.HasPrefix should panic when string does not have prefix")
		}
	}()
	assert.HasPrefix("hello world", "world", "This should fail")
}

func TestAssertHasSuffix(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.HasSuffix should not panic when string has suffix")
		}
	}()
	assert.HasSuffix("hello world", "world", "String should have suffix")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.HasSuffix should panic when string does not have suffix")
		}
	}()
	assert.HasSuffix("hello world", "hello", "This should fail")
}

func TestAssertInSlice(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.InSlice should not panic when value is in slice")
		}
	}()
	assert.InSlice("banana", slice, "Value should be in slice")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.InSlice should panic when value is not in slice")
		}
	}()
	assert.InSlice("grape", slice, "This should fail")
}

func TestAssertNotInSlice(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.NotInSlice should not panic when value is not in slice")
		}
	}()
	assert.NotInSlice("grape", slice, "Value should not be in slice")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.NotInSlice should panic when value is in slice")
		}
	}()
	assert.NotInSlice("banana", slice, "This should fail")
}

func TestAssertError(t *testing.T) {
	err := errors.New("test error")

	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Error should not panic when error is not nil")
		}
	}()
	assert.Error(err, "Error should not be nil")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Error should panic when error is nil")
		}
	}()
	assert.Error(nil, "This should fail")
}

func TestAssertNoError(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.NoError should not panic when error is nil")
		}
	}()
	assert.NoError(nil, "Error should be nil")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.NoError should panic when error is not nil")
		}
	}()
	err := errors.New("test error")
	assert.NoError(err, "This should fail")
}

func TestAssertPanic(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Panic should not panic when function panics")
		}
	}()
	assert.Panic(func() {
		panic("test panic")
	}, "Function should panic")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Panic should panic when function does not panic")
		}
	}()
	assert.Panic(func() {
		// This function does not panic
	}, "This should fail")
}

func TestAssertNotPanic(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.NotPanic should not panic when function does not panic")
		}
	}()
	assert.NotPanic(func() {
		// This function does not panic
	}, "Function should not panic")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.NotPanic should panic when function panics")
		}
	}()
	assert.NotPanic(func() {
		panic("test panic")
	}, "This should fail")
}

func TestAssertInRange(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.InRange should not panic when value is in range")
		}
	}()
	assert.InRange(5, 1, 10, "Value should be in range")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.InRange should panic when value is not in range")
		}
	}()
	assert.InRange(15, 1, 10, "This should fail")
}

func TestAssertLength(t *testing.T) {
	// 测试成功情况
	defer func() {
		if r := recover(); r != nil {
			t.Error("Assert.Length should not panic when length is correct")
		}
	}()
	assert.Length("hello", 5, "String length should be 5")
	assert.Length([]int{1, 2, 3}, 3, "Slice length should be 3")

	// 测试失败情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert.Length should panic when length is incorrect")
		}
	}()
	assert.Length("hello", 3, "This should fail")
}

// TestCustomErrorHandler 测试自定义错误处理器
func TestCustomErrorHandler(t *testing.T) {
	// 保存原始处理器
	originalHandler := assert.GlobalHandler

	// 设置自定义处理器
	errorCaught := false
	assert.SetGlobalHandler(func(err *assert.AssertionError) {
		errorCaught = true
		if err.Message != "test message" {
			t.Errorf("Expected message 'test message', got '%s'", err.Message)
		}
	})

	// 触发断言错误
	assert.True(false, "test message")

	// 验证错误被捕获
	if !errorCaught {
		t.Error("Custom error handler should have caught the assertion error")
	}

	// 恢复原始处理器
	assert.SetGlobalHandler(originalHandler)
}
