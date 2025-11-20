/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 00:00:00
 * @FilePath: \go-toolbox\pkg\mathx\ternary_exec_test.go
 * @Description: IfExec 和 IfExecElse 函数测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package mathx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIfExec 测试 IfExec 函数
func TestIfExec(t *testing.T) {
	t.Run("条件为true时执行action", func(t *testing.T) {
		executed := false
		IfExec(true, func() {
			executed = true
		})
		assert.True(t, executed, "条件为true时应该执行action")
	})

	t.Run("条件为false时不执行action", func(t *testing.T) {
		executed := false
		IfExec(false, func() {
			executed = true
		})
		assert.False(t, executed, "条件为false时不应该执行action")
	})

	t.Run("action为nil时不会panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			IfExec(true, nil)
		}, "action为nil时不应该panic")
	})

	t.Run("副作用操作修改外部变量", func(t *testing.T) {
		counter := 0
		IfExec(true, func() {
			counter++
			counter++
		})
		assert.Equal(t, 2, counter, "应该执行副作用操作并修改外部变量")
	})

	t.Run("多次调用累积副作用", func(t *testing.T) {
		sum := 0
		for i := 1; i <= 5; i++ {
			IfExec(i%2 == 0, func() {
				sum += i
			})
		}
		assert.Equal(t, 6, sum, "应该累加偶数: 2+4=6")
	})
}

// TestIfExecElse 测试 IfExecElse 函数
func TestIfExecElse(t *testing.T) {
	t.Run("条件为true时执行onTrue", func(t *testing.T) {
		result := ""
		IfExecElse(true,
			func() { result = "true" },
			func() { result = "false" },
		)
		assert.Equal(t, "true", result, "条件为true时应该执行onTrue")
	})

	t.Run("条件为false时执行onFalse", func(t *testing.T) {
		result := ""
		IfExecElse(false,
			func() { result = "true" },
			func() { result = "false" },
		)
		assert.Equal(t, "false", result, "条件为false时应该执行onFalse")
	})

	t.Run("onTrue为nil时不会panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			IfExecElse(true, nil, func() {})
		}, "onTrue为nil时不应该panic")
	})

	t.Run("onFalse为nil时不会panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			IfExecElse(false, func() {}, nil)
		}, "onFalse为nil时不应该panic")
	})

	t.Run("两个回调都为nil时不会panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			IfExecElse(true, nil, nil)
		}, "两个回调都为nil时不应该panic")
	})

	t.Run("错误处理场景", func(t *testing.T) {
		err := errors.New("test error")
		logMessage := ""

		IfExecElse(err == nil,
			func() { logMessage = "Success" },
			func() { logMessage = "Error: " + err.Error() },
		)

		assert.Equal(t, "Error: test error", logMessage, "应该执行错误分支")
	})

	t.Run("成功处理场景", func(t *testing.T) {
		var err error = nil
		logMessage := ""

		IfExecElse(err == nil,
			func() { logMessage = "Success" },
			func() { logMessage = "Error: " + err.Error() },
		)

		assert.Equal(t, "Success", logMessage, "应该执行成功分支")
	})

	t.Run("复杂副作用操作", func(t *testing.T) {
		type User struct {
			Name  string
			Count int
		}

		user := &User{Name: "Alice", Count: 0}
		isValid := true

		IfExecElse(isValid,
			func() {
				user.Count++
				user.Name = user.Name + " (validated)"
			},
			func() {
				user.Count = -1
				user.Name = user.Name + " (invalid)"
			},
		)

		assert.Equal(t, "Alice (validated)", user.Name)
		assert.Equal(t, 1, user.Count)
	})

	t.Run("数值比较场景", func(t *testing.T) {
		score := 85
		grade := ""

		IfExecElse(score >= 60,
			func() { grade = "Pass" },
			func() { grade = "Fail" },
		)

		assert.Equal(t, "Pass", grade)
	})
}

// TestIfExecWithIfExecElse 测试 IfExec 和 IfExecElse 组合使用
func TestIfExecWithIfExecElse(t *testing.T) {
	t.Run("嵌套使用场景", func(t *testing.T) {
		user := struct {
			Name    string
			IsAdmin bool
			Logged  bool
		}{Name: "Bob", IsAdmin: true, Logged: false}

		IfExec(user.IsAdmin, func() {
			IfExecElse(user.Logged,
				func() { user.Name = user.Name + " (admin, online)" },
				func() { user.Name = user.Name + " (admin, offline)" },
			)
		})

		assert.Equal(t, "Bob (admin, offline)", user.Name)
	})

	t.Run("多条件判断", func(t *testing.T) {
		flags := []bool{false, true, false}
		results := []string{}

		for i, flag := range flags {
			IfExecElse(flag,
				func() { results = append(results, "on") },
				func() { results = append(results, "off") },
			)
			_ = i
		}

		assert.Equal(t, []string{"off", "on", "off"}, results)
	})
}

// BenchmarkIfExec 性能测试
func BenchmarkIfExec(b *testing.B) {
	counter := 0
	for i := 0; i < b.N; i++ {
		IfExec(true, func() {
			counter++
		})
	}
}

// BenchmarkIfExecElse 性能测试
func BenchmarkIfExecElse(b *testing.B) {
	counter := 0
	for i := 0; i < b.N; i++ {
		IfExecElse(i%2 == 0,
			func() { counter++ },
			func() { counter-- },
		)
	}
}

// BenchmarkTraditionalIf 传统if性能对比
func BenchmarkTraditionalIf(b *testing.B) {
	counter := 0
	for i := 0; i < b.N; i++ {
		if true {
			counter++
		}
	}
}

// BenchmarkTraditionalIfElse 传统if-else性能对比
func BenchmarkTraditionalIfElse(b *testing.B) {
	counter := 0
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			counter++
		} else {
			counter--
		}
	}
}
