/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 15:32:55
 * @FilePath: \go-toolbox\pkg\mathx\ternary_chain_test.go
 * @Description: 链式调用三元运算符测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhen(t *testing.T) {
	t.Run("条件为true时执行Then", func(t *testing.T) {
		executed := false
		When(true).
			Then(func() { executed = true }).
			Else(func() { executed = false }).
			Do()

		assert.True(t, executed, "期望执行Then分支")
	})

	t.Run("条件为false时执行Else", func(t *testing.T) {
		executed := false
		When(false).
			Then(func() { executed = true }).
			Else(func() { executed = false }).
			Do()

		assert.False(t, executed, "期望执行Else分支")
	})

	t.Run("只有Then没有Else", func(t *testing.T) {
		count := 0
		When(true).
			Then(func() { count++ }).
			Do()

		assert.Equal(t, 1, count, "期望count=1")
	})

	t.Run("日志场景示例", func(t *testing.T) {
		err := fmt.Errorf("测试错误")
		logs := []string{}

		When(err != nil).
			Then(func() { logs = append(logs, "❌ 操作失败") }).
			Else(func() { logs = append(logs, "✅ 操作成功") }).
			Do()

		assert.Len(t, logs, 1, "期望日志长度为1")
		assert.Equal(t, "❌ 操作失败", logs[0], "期望日志='❌ 操作失败'")
	})
}

func TestWhenValue(t *testing.T) {
	t.Run("返回值-条件为true", func(t *testing.T) {
		result := WhenValue[int](true).
			ThenReturn(100).
			ElseReturn(0).
			Get()

		assert.Equal(t, 100, result, "期望result=100")
	})

	t.Run("返回值-条件为false", func(t *testing.T) {
		result := WhenValue[int](false).
			ThenReturn(100).
			ElseReturn(0).
			Get()

		assert.Equal(t, 0, result, "期望result=0")
	})

	t.Run("返回值-使用函数", func(t *testing.T) {
		x := 10
		result := WhenValue[string](x > 0).
			ThenDo(func() string { return "正数" }).
			ElseDo(func() string { return "负数或零" }).
			Get()

		assert.Equal(t, "正数", result, "期望result='正数'")
	})

	t.Run("返回值-混合使用", func(t *testing.T) {
		result := WhenValue[int](true).
			ThenDo(func() int { return 42 }).
			ElseReturn(0).
			Get()

		assert.Equal(t, 42, result, "期望result=42")
	})

	t.Run("字符串类型", func(t *testing.T) {
		isSuccess := true
		msg := WhenValue[string](isSuccess).
			ThenReturn("成功").
			ElseReturn("失败").
			Get()

		assert.Equal(t, "成功", msg, "期望msg='成功'")
	})
}

func ExampleWhen() {
	err := fmt.Errorf("操作失败")

	// 链式调用示例
	When(err != nil).
		Then(func() { fmt.Println("❌ 发生错误") }).
		Else(func() { fmt.Println("✅ 操作成功") }).
		Do()

	// Output:
	// ❌ 发生错误
}

func ExampleWhenValue() {
	x := 5

	// 返回值链式调用
	result := WhenValue[string](x > 0).
		ThenReturn("正数").
		ElseReturn("非正数").
		Get()

	fmt.Println(result)

	// Output:
	// 正数
}

// BenchmarkWhen 性能测试
func BenchmarkWhen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		When(i%2 == 0).
			Then(func() {}).
			Else(func() {}).
			Do()
	}
}

// BenchmarkWhenValue 性能测试
func BenchmarkWhenValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = WhenValue[int](i%2 == 0).
			ThenReturn(1).
			ElseReturn(0).
			Get()
	}
}
