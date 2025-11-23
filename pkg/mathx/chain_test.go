/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 09:00:15
 * @FilePath: \go-toolbox\pkg\mathx\chain_test.go
 * @Description: 链式条件构建器测试 - 全面覆盖所有功能
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package mathx

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewIFChain 测试创建新的链式构建器
func TestNewIFChain(t *testing.T) {
	chain := NewIFChain[int]()
	assert.NotNil(t, chain)
	assert.False(t, chain.executed)
	assert.False(t, chain.hasReturn)
	assert.Zero(t, chain.returnValue)
}

// TestChain 测试全局链式构建器入口
func TestChain(t *testing.T) {
	chain := IFChain()
	assert.NotNil(t, chain)
	assert.IsType(t, &IFChainBuilder[any]{}, chain)
}

// TestChainFor 测试为特定类型创建链式构建器
func TestChainFor(t *testing.T) {
	intChain := IFChainFor[int]()
	assert.NotNil(t, intChain)
	assert.IsType(t, &IFChainBuilder[int]{}, intChain)

	stringChain := IFChainFor[string]()
	assert.NotNil(t, stringChain)
	assert.IsType(t, &IFChainBuilder[string]{}, stringChain)

	errorChain := IFChainFor[error]()
	assert.NotNil(t, errorChain)
	assert.IsType(t, &IFChainBuilder[error]{}, errorChain)
}

// TestIFChainBuilderWhen 测试 When 方法
func TestIFChainBuilderWhen(t *testing.T) {
	t.Run("condition true", func(t *testing.T) {
		chain := NewIFChain[int]()
		condition := chain.When(true)
		assert.NotNil(t, condition)
		assert.True(t, condition.condition)
		assert.Equal(t, chain, condition.chain)
	})

	t.Run("condition false", func(t *testing.T) {
		chain := NewIFChain[int]()
		condition := chain.When(false)
		assert.NotNil(t, condition)
		assert.False(t, condition.condition)
		assert.Equal(t, chain, condition.chain)
	})

	t.Run("already executed", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.executed = true
		condition := chain.When(true)
		assert.NotNil(t, condition)
		assert.False(t, condition.condition) // 应该跳过
		assert.Equal(t, chain, condition.chain)
	})
}

// TestIFChainBuilderConditionThen 测试 Then 方法
func TestIFChainBuilderConditionThen(t *testing.T) {
	chain := NewIFChain[int]()
	condition := chain.When(true)
	action := condition.Then(func() {})

	assert.NotNil(t, action)
	assert.Equal(t, chain, action.chain)
	assert.True(t, action.condition)
	assert.NotNil(t, action.action)
}

// TestIFChainBuilderConditionThenReturn 测试 ThenReturn 方法
func TestIFChainBuilderConditionThenReturn(t *testing.T) {
	t.Run("condition true, not executed", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		result := chain.When(true).ThenReturn(42, func() {
			executed = true
		})

		assert.Equal(t, chain, result)
		assert.True(t, executed)
		assert.True(t, chain.executed)
		assert.True(t, chain.hasReturn)
		assert.Equal(t, 42, chain.returnValue)
	})

	t.Run("condition true, no action", func(t *testing.T) {
		chain := NewIFChain[int]()
		result := chain.When(true).ThenReturn(42)

		assert.Equal(t, chain, result)
		assert.True(t, chain.executed)
		assert.True(t, chain.hasReturn)
		assert.Equal(t, 42, chain.returnValue)
	})

	t.Run("condition false", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		result := chain.When(false).ThenReturn(42, func() {
			executed = true
		})

		assert.Equal(t, chain, result)
		assert.False(t, executed)
		assert.False(t, chain.executed)
		assert.False(t, chain.hasReturn)
		assert.Zero(t, chain.returnValue)
	})

	t.Run("already executed", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		chain.executed = true
		result := chain.When(true).ThenReturn(42, func() {
			executed = true
		})

		assert.Equal(t, chain, result)
		assert.False(t, executed) // 不应该执行
	})

	t.Run("nil action", func(t *testing.T) {
		chain := NewIFChain[int]()
		result := chain.When(true).ThenReturn(42, nil)

		assert.Equal(t, chain, result)
		assert.True(t, chain.executed)
		assert.Equal(t, 42, chain.returnValue)
	})
}

// TestIFChainBuilderConditionThenReturnNil 测试 ThenReturnNil 方法
func TestIFChainBuilderConditionThenReturnNil(t *testing.T) {
	t.Run("with action", func(t *testing.T) {
		executed := false
		chain := NewIFChain[*string]()
		result := chain.When(true).ThenReturnNil(func() {
			executed = true
		})

		assert.Equal(t, chain, result)
		assert.True(t, executed)
		assert.True(t, chain.executed)
		assert.Nil(t, chain.returnValue)
	})

	t.Run("without action", func(t *testing.T) {
		chain := NewIFChain[*string]()
		result := chain.When(true).ThenReturnNil()

		assert.Equal(t, chain, result)
		assert.True(t, chain.executed)
		assert.Nil(t, chain.returnValue)
	})

	t.Run("condition false", func(t *testing.T) {
		chain := NewIFChain[*string]()
		result := chain.When(false).ThenReturnNil(func() {})

		assert.Equal(t, chain, result)
		assert.False(t, chain.executed)
	})
}

// TestIFChainBuilderActionReturn 测试 Return 方法
func TestIFChainBuilderActionReturn(t *testing.T) {
	t.Run("condition true, not executed", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		result := chain.When(true).Then(func() {
			executed = true
		}).Return(42)

		assert.Equal(t, chain, result)
		assert.True(t, executed)
		assert.True(t, chain.executed)
		assert.True(t, chain.hasReturn)
		assert.Equal(t, 42, chain.returnValue)
	})

	t.Run("condition false", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		result := chain.When(false).Then(func() {
			executed = true
		}).Return(42)

		assert.Equal(t, chain, result)
		assert.False(t, executed)
		assert.False(t, chain.executed)
		assert.False(t, chain.hasReturn)
		assert.Zero(t, chain.returnValue)
	})

	t.Run("already executed", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		chain.executed = true
		result := chain.When(true).Then(func() {
			executed = true
		}).Return(42)

		assert.Equal(t, chain, result)
		assert.False(t, executed)
	})

	t.Run("nil action", func(t *testing.T) {
		chain := NewIFChain[int]()
		result := chain.When(true).Then(nil).Return(42)

		assert.Equal(t, chain, result)
		assert.True(t, chain.executed)
		assert.Equal(t, 42, chain.returnValue)
	})
}

// TestIFChainBuilderActionReturnNil 测试 ReturnNil 方法
func TestIFChainBuilderActionReturnNil(t *testing.T) {
	chain := NewIFChain[*string]()
	result := chain.When(true).Then(func() {}).ReturnNil()

	assert.Equal(t, chain, result)
	assert.True(t, chain.executed)
	assert.Nil(t, chain.returnValue)
}

// TestIFChainBuilderActionContinueChain 测试 ContinueChain 方法
func TestIFChainBuilderActionContinueChain(t *testing.T) {
	t.Run("condition true, not executed", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		result := chain.When(true).Then(func() {
			executed = true
		}).ContinueChain()

		assert.Equal(t, chain, result)
		assert.True(t, executed)
		assert.False(t, chain.executed) // ContinueChain 不设置 executed
	})

	t.Run("condition false", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		result := chain.When(false).Then(func() {
			executed = true
		}).ContinueChain()

		assert.Equal(t, chain, result)
		assert.False(t, executed)
	})

	t.Run("already executed", func(t *testing.T) {
		executed := false
		chain := NewIFChain[int]()
		chain.executed = true
		result := chain.When(true).Then(func() {
			executed = true
		}).ContinueChain()

		assert.Equal(t, chain, result)
		assert.False(t, executed)
	})

	t.Run("nil action", func(t *testing.T) {
		chain := NewIFChain[int]()
		result := chain.When(true).Then(nil).ContinueChain()

		assert.Equal(t, chain, result)
		assert.False(t, chain.executed)
	})
}

// TestIFChainBuilderExecute 测试 Execute 方法
func TestIFChainBuilderExecute(t *testing.T) {
	t.Run("has return value", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.When(true).ThenReturn(42)

		value, hasReturn := chain.Execute()
		assert.True(t, hasReturn)
		assert.Equal(t, 42, value)
	})

	t.Run("no return value", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.When(false).ThenReturn(42)

		value, hasReturn := chain.Execute()
		assert.False(t, hasReturn)
		assert.Zero(t, value)
	})
}

// TestIFChainBuilderMustExecute 测试 MustExecute 方法
func TestIFChainBuilderMustExecute(t *testing.T) {
	t.Run("has return value", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.When(true).ThenReturn(42)

		value := chain.MustExecute()
		assert.Equal(t, 42, value)
	})

	t.Run("no return value - should panic", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.When(false).ThenReturn(42)

		assert.Panics(t, func() {
			chain.MustExecute()
		})
	})
}

// TestIFChainBuilderExecuteOr 测试 ExecuteOr 方法
func TestIFChainBuilderExecuteOr(t *testing.T) {
	t.Run("has return value", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.When(true).ThenReturn(42)

		value := chain.ExecuteOr(99)
		assert.Equal(t, 42, value)
	})

	t.Run("no return value - use default", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.When(false).ThenReturn(42)

		value := chain.ExecuteOr(99)
		assert.Equal(t, 99, value)
	})
}

// TestIFChainBuilderHasResult 测试 HasResult 方法
func TestIFChainBuilderHasResult(t *testing.T) {
	t.Run("has result", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.When(true).ThenReturn(42)

		assert.True(t, chain.HasResult())
	})

	t.Run("no result", func(t *testing.T) {
		chain := NewIFChain[int]()
		chain.When(false).ThenReturn(42)

		assert.False(t, chain.HasResult())
	})
}

// TestErrorChain 测试 ErrorChain 便利函数
func TestErrorChain(t *testing.T) {
	chain := IFErrorChain()
	assert.NotNil(t, chain)
	assert.IsType(t, &IFChainBuilder[error]{}, chain)

	testErr := errors.New("test error")
	err, hasErr := chain.When(true).ThenReturn(testErr).Execute()
	assert.True(t, hasErr)
	assert.Equal(t, testErr, err)
}

// TestNilChain 测试 NilChain 便利函数
func TestNilChain(t *testing.T) {
	chain := IFNilChain()
	assert.NotNil(t, chain)
	assert.IsType(t, &IFChainBuilder[any]{}, chain)

	value, hasValue := chain.When(true).ThenReturnNil().Execute()
	assert.True(t, hasValue)
	assert.Nil(t, value)
}

// TestComplexChaining 测试复杂的链式调用
func TestComplexChaining(t *testing.T) {
	t.Run("multiple conditions - first matches", func(t *testing.T) {
		executed1, executed2, executed3 := false, false, false

		value, hasValue := NewIFChain[string]().
			When(true).
			ThenReturn("first", func() { executed1 = true }).
			When(true).
			ThenReturn("second", func() { executed2 = true }).
			When(true).
			ThenReturn("third", func() { executed3 = true }).
			Execute()

		assert.True(t, hasValue)
		assert.Equal(t, "first", value)
		assert.True(t, executed1)
		assert.False(t, executed2) // 不应该执行
		assert.False(t, executed3) // 不应该执行
	})

	t.Run("multiple conditions - second matches", func(t *testing.T) {
		executed1, executed2, executed3 := false, false, false

		value, hasValue := NewIFChain[string]().
			When(false).
			ThenReturn("first", func() { executed1 = true }).
			When(true).
			ThenReturn("second", func() { executed2 = true }).
			When(true).
			ThenReturn("third", func() { executed3 = true }).
			Execute()

		assert.True(t, hasValue)
		assert.Equal(t, "second", value)
		assert.False(t, executed1)
		assert.True(t, executed2)
		assert.False(t, executed3) // 不应该执行
	})

	t.Run("no conditions match", func(t *testing.T) {
		executed1, executed2, executed3 := false, false, false

		value, hasValue := NewIFChain[string]().
			When(false).
			ThenReturn("first", func() { executed1 = true }).
			When(false).
			ThenReturn("second", func() { executed2 = true }).
			When(false).
			ThenReturn("third", func() { executed3 = true }).
			Execute()

		assert.False(t, hasValue)
		assert.Empty(t, value)
		assert.False(t, executed1)
		assert.False(t, executed2)
		assert.False(t, executed3)
	})
}

// TestContinueChainFlow 测试 ContinueChain 流程
func TestContinueChainFlow(t *testing.T) {
	executed1, executed2, executed3 := false, false, false

	value, hasValue := NewIFChain[string]().
		When(true).
		Then(func() { executed1 = true }).
		ContinueChain().
		When(true).
		Then(func() { executed2 = true }).
		ContinueChain().
		When(true).
		ThenReturn("final", func() { executed3 = true }).
		Execute()

	assert.True(t, hasValue)
	assert.Equal(t, "final", value)
	assert.True(t, executed1)
	assert.True(t, executed2)
	assert.True(t, executed3)
}

// TestMixedChainFlow 测试混合链式流程
func TestMixedChainFlow(t *testing.T) {
	step1, step2, step3, step4 := false, false, false, false

	value, hasValue := NewIFChain[int]().
		When(true).
		Then(func() { step1 = true }).
		ContinueChain().
		When(false).
		ThenReturn(99, func() { step2 = true }).
		When(true).
		Then(func() { step3 = true }).
		Return(42).
		When(true).
		ThenReturn(66, func() { step4 = true }). // 不应该执行
		Execute()

	assert.True(t, hasValue)
	assert.Equal(t, 42, value)
	assert.True(t, step1)
	assert.False(t, step2)
	assert.True(t, step3)
	assert.False(t, step4) // 已经执行过了
}

// TestTypeSafety 测试类型安全
func TestTypeSafety(t *testing.T) {
	t.Run("int type", func(t *testing.T) {
		value, _ := IFChainFor[int]().When(true).ThenReturn(42).Execute()
		assert.Equal(t, 42, value)
	})

	t.Run("string type", func(t *testing.T) {
		value, _ := IFChainFor[string]().When(true).ThenReturn("hello").Execute()
		assert.Equal(t, "hello", value)
	})

	t.Run("custom struct type", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}

		expected := TestStruct{Name: "test", Age: 25}
		value, _ := IFChainFor[TestStruct]().When(true).ThenReturn(expected).Execute()
		assert.Equal(t, expected, value)
	})

	t.Run("pointer type", func(t *testing.T) {
		str := "test"
		value, _ := IFChainFor[*string]().When(true).ThenReturn(&str).Execute()
		assert.Equal(t, &str, value)
		assert.Equal(t, "test", *value)
	})
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Run("empty action function", func(t *testing.T) {
		value, hasValue := NewIFChain[int]().
			When(true).
			ThenReturn(42, func() {}).
			Execute()

		assert.True(t, hasValue)
		assert.Equal(t, 42, value)
	})

	t.Run("multiple ThenReturn calls", func(t *testing.T) {
		step1, step2 := false, false

		chain := NewIFChain[int]()
		chain.When(true).ThenReturn(42, func() { step1 = true })
		chain.When(true).ThenReturn(99, func() { step2 = true }) // 不应该执行

		value, hasValue := chain.Execute()
		assert.True(t, hasValue)
		assert.Equal(t, 42, value)
		assert.True(t, step1)
		assert.False(t, step2)
	})

	t.Run("zero value types", func(t *testing.T) {
		t.Run("int zero", func(t *testing.T) {
			value, hasValue := IFChainFor[int]().When(true).ThenReturn(0).Execute()
			assert.True(t, hasValue)
			assert.Equal(t, 0, value)
		})

		t.Run("string zero", func(t *testing.T) {
			value, hasValue := IFChainFor[string]().When(true).ThenReturn("").Execute()
			assert.True(t, hasValue)
			assert.Equal(t, "", value)
		})

		t.Run("bool zero", func(t *testing.T) {
			value, hasValue := IFChainFor[bool]().When(true).ThenReturn(false).Execute()
			assert.True(t, hasValue)
			assert.Equal(t, false, value)
		})
	})
}

// BenchmarkIFChainBuilder 性能基准测试
func BenchmarkIFChainBuilder(b *testing.B) {
	b.Run("simple chain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = NewIFChain[int]().When(true).ThenReturn(42).Execute()
		}
	})

	b.Run("complex chain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = NewIFChain[string]().
				When(false).ThenReturn("first").
				When(false).ThenReturn("second").
				When(true).ThenReturn("third").
				Execute()
		}
	})

	b.Run("with actions", func(b *testing.B) {
		counter := 0
		for i := 0; i < b.N; i++ {
			_, _ = NewIFChain[int]().
				When(true).
				ThenReturn(42, func() { counter++ }).
				Execute()
		}
	})
}
