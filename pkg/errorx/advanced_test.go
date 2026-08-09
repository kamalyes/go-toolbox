/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-10 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-10 00:00:00
 * @FilePath: \go-toolbox\pkg\errorx\advanced_test.go
 * @Description: advanced.go 的单元测试，覆盖错误链、重试、验证等高级错误处理功能
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package errorx

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewErrorChain 测试创建新的错误链
func TestNewErrorChain(t *testing.T) {
	chain := NewErrorChain()
	assert.NotNil(t, chain)
	assert.False(t, chain.HasErrors())
	assert.Empty(t, chain.GetErrors())
	assert.Nil(t, chain.GetLastError())
}

// TestErrorChainAddError 测试向错误链中添加错误
func TestErrorChainAddError(t *testing.T) {
	t.Run("添加普通错误", func(t *testing.T) {
		chain := NewErrorChain()
		chain.AddError(errors.New("first error"))
		assert.True(t, chain.HasErrors())
		errs := chain.GetErrors()
		require.Len(t, errs, 1)
		assert.Equal(t, "first error", errs[0].Error.Error())
		assert.NotEmpty(t, errs[0].Location)
		assert.False(t, errs[0].Timestamp.IsZero())
		assert.NotNil(t, errs[0].Context)
	})

	t.Run("添加 nil 错误应被忽略", func(t *testing.T) {
		chain := NewErrorChain()
		result := chain.AddError(nil)
		assert.Same(t, chain, result)
		assert.False(t, chain.HasErrors())
	})

	t.Run("链式添加多个错误", func(t *testing.T) {
		chain := NewErrorChain()
		chain.AddError(errors.New("first")).AddError(errors.New("second"))
		assert.Len(t, chain.GetErrors(), 2)
	})
}

// TestErrorChainAddErrorWithContext 测试带上下文添加错误
func TestErrorChainAddErrorWithContext(t *testing.T) {
	t.Run("添加带上下文的错误", func(t *testing.T) {
		chain := NewErrorChain()
		ctx := map[string]interface{}{"user": "kamal", "code": 500}
		chain.AddErrorWithContext(errors.New("ctx error"), ctx)
		errs := chain.GetErrors()
		require.Len(t, errs, 1)
		assert.Equal(t, ctx, errs[0].Context)
		assert.NotEmpty(t, errs[0].Location)
	})

	t.Run("nil 错误带上下文应被忽略", func(t *testing.T) {
		chain := NewErrorChain()
		result := chain.AddErrorWithContext(nil, map[string]interface{}{"k": "v"})
		assert.Same(t, chain, result)
		assert.False(t, chain.HasErrors())
	})
}

// TestErrorChainGetLastError 测试获取最后一个错误
func TestErrorChainGetLastError(t *testing.T) {
	t.Run("空链返回 nil", func(t *testing.T) {
		chain := NewErrorChain()
		assert.Nil(t, chain.GetLastError())
	})

	t.Run("返回最后添加的错误", func(t *testing.T) {
		chain := NewErrorChain()
		chain.AddError(errors.New("first"))
		chain.AddError(errors.New("second"))
		last := chain.GetLastError()
		require.NotNil(t, last)
		assert.Equal(t, "second", last.Error.Error())
	})
}

// TestErrorChainError 测试 ErrorChain 的 Error 方法
func TestErrorChainError(t *testing.T) {
	t.Run("空链返回空字符串", func(t *testing.T) {
		chain := NewErrorChain()
		assert.Empty(t, chain.Error())
	})

	t.Run("单个错误格式", func(t *testing.T) {
		chain := NewErrorChain()
		chain.AddError(errors.New("single error"))
		errStr := chain.Error()
		assert.Contains(t, errStr, "single error")
		assert.Contains(t, errStr, "error at")
	})

	t.Run("多个错误格式", func(t *testing.T) {
		chain := NewErrorChain()
		chain.AddError(errors.New("first"))
		chain.AddError(errors.New("second"))
		errStr := chain.Error()
		assert.Contains(t, errStr, "error chain:")
		assert.Contains(t, errStr, "first")
		assert.Contains(t, errStr, "second")
		assert.Contains(t, errStr, " -> ")
	})
}

// TestErrorChainString 测试 ErrorChain 的 String 方法
func TestErrorChainString(t *testing.T) {
	t.Run("空链返回 no errors", func(t *testing.T) {
		chain := NewErrorChain()
		assert.Equal(t, "no errors", chain.String())
	})

	t.Run("带错误和上下文的详细字符串", func(t *testing.T) {
		chain := NewErrorChain()
		chain.AddErrorWithContext(errors.New("detail error"), map[string]interface{}{"key": "value"})
		s := chain.String()
		assert.Contains(t, s, "detail error")
		assert.Contains(t, s, "Message:")
		assert.Contains(t, s, "Time:")
		assert.Contains(t, s, "Location:")
		assert.Contains(t, s, "Context:")
	})

	t.Run("无上下文时不输出 Context", func(t *testing.T) {
		chain := NewErrorChain()
		chain.AddError(errors.New("no ctx error"))
		s := chain.String()
		assert.Contains(t, s, "no ctx error")
		assert.NotContains(t, s, "Context:")
	})
}

// TestNewErrorWithStack 测试创建带调用栈的错误
func TestNewErrorWithStack(t *testing.T) {
	err := NewErrorWithStack("stack error")
	assert.NotNil(t, err)
	assert.NotEmpty(t, err.Stack)
	assert.Equal(t, "stack error", err.Error())
}

// TestErrorWithStackGetStackTrace 测试获取调用栈信息
func TestErrorWithStackGetStackTrace(t *testing.T) {
	err := NewErrorWithStack("stack error")
	trace := err.GetStackTrace()
	assert.NotEmpty(t, trace)
	assert.Contains(t, trace, "advanced_test.go")
}

// TestErrorWithStackString 测试带调用栈错误的字符串表示
func TestErrorWithStackString(t *testing.T) {
	err := NewErrorWithStack("stack error")
	s := err.String()
	assert.Contains(t, s, "Error: stack error")
	assert.Contains(t, s, "Stack trace:")
}

// TestNewRetryableError 测试创建可重试错误
func TestNewRetryableError(t *testing.T) {
	err := NewRetryableError("retry me", 3, time.Second)
	assert.NotNil(t, err)
	assert.Equal(t, 3, err.MaxRetries)
	assert.Equal(t, 0, err.CurrentRetry)
	assert.Equal(t, time.Second, err.RetryAfter)
	assert.True(t, err.Retryable)
}

// TestRetryableErrorError 测试可重试错误的 Error 方法
func TestRetryableErrorError(t *testing.T) {
	err := NewRetryableError("retry me", 3, time.Second)
	err.CurrentRetry = 1
	assert.Equal(t, "retry me (retry 1/3)", err.Error())
}

// TestRetryableErrorShouldRetry 测试 ShouldRetry 逻辑
func TestRetryableErrorShouldRetry(t *testing.T) {
	t.Run("可重试且未达上限", func(t *testing.T) {
		err := NewRetryableError("retry", 3, time.Second)
		assert.True(t, err.ShouldRetry())
	})

	t.Run("达到重试上限", func(t *testing.T) {
		err := NewRetryableError("retry", 3, time.Second)
		err.CurrentRetry = 3
		assert.False(t, err.ShouldRetry())
	})

	t.Run("禁用重试后不再重试", func(t *testing.T) {
		err := NewRetryableError("retry", 3, time.Second)
		err.DisableRetry()
		assert.False(t, err.ShouldRetry())
	})
}

// TestRetryableErrorIncrementRetry 测试增加重试次数
func TestRetryableErrorIncrementRetry(t *testing.T) {
	err := NewRetryableError("retry", 3, time.Second)
	assert.Equal(t, 0, err.CurrentRetry)
	err.IncrementRetry()
	assert.Equal(t, 1, err.CurrentRetry)
	err.IncrementRetry()
	assert.Equal(t, 2, err.CurrentRetry)
}

// TestRetryableErrorGetRetryAfter 测试获取重试间隔（指数退避）
func TestRetryableErrorGetRetryAfter(t *testing.T) {
	err := NewRetryableError("retry", 3, time.Second)
	// CurrentRetry=0 时，1<<0 = 1
	assert.Equal(t, time.Second, err.GetRetryAfter())
	err.IncrementRetry()
	// CurrentRetry=1 时，1<<1 = 2
	assert.Equal(t, 2*time.Second, err.GetRetryAfter())
	err.IncrementRetry()
	// CurrentRetry=2 时，1<<2 = 4
	assert.Equal(t, 4*time.Second, err.GetRetryAfter())
}

// TestNewValidationError 测试创建验证错误
func TestNewValidationError(t *testing.T) {
	ve := NewValidationError("username", "required", "", "must not be empty")
	assert.NotNil(t, ve)
	assert.Equal(t, "username", ve.Field)
	assert.Equal(t, "required", ve.Rule)
	assert.NotNil(t, ve.Details)
}

// TestValidationErrorError 测试验证错误的 Error 方法
func TestValidationErrorError(t *testing.T) {
	ve := NewValidationError("username", "required", "", "must not be empty")
	errStr := ve.Error()
	assert.Contains(t, errStr, "username")
	assert.Contains(t, errStr, "required")
	assert.Contains(t, errStr, "must not be empty")
}

// TestValidationErrorAddDetail 测试添加验证错误详情
func TestValidationErrorAddDetail(t *testing.T) {
	ve := NewValidationError("field", "rule", "", "msg")
	ve.AddDetail("min", 1)
	ve.AddDetail("max", 100)
	assert.Equal(t, 1, ve.Details["min"])
	assert.Equal(t, 100, ve.Details["max"])
}

// TestValidationErrors 测试多个验证错误的集合操作
func TestValidationErrors(t *testing.T) {
	t.Run("空集合", func(t *testing.T) {
		ves := NewValidationErrors()
		assert.False(t, ves.HasErrors())
		assert.Empty(t, ves.Error())
	})

	t.Run("单个错误", func(t *testing.T) {
		ves := NewValidationErrors()
		ve := NewValidationError("name", "required", "", "empty")
		ves.Add(ve)
		assert.True(t, ves.HasErrors())
		assert.Equal(t, ve.Error(), ves.Error())
	})

	t.Run("多个错误", func(t *testing.T) {
		ves := NewValidationErrors()
		ves.Add(NewValidationError("name", "required", "", "empty"))
		ves.Add(NewValidationError("age", "min", "", "too small"))
		errStr := ves.Error()
		assert.Contains(t, errStr, "multiple validation errors:")
		assert.Contains(t, errStr, "name")
		assert.Contains(t, errStr, "age")
		assert.Contains(t, errStr, "; ")
	})
}

// TestValidationErrorsGetFieldErrors 测试按字段获取验证错误
func TestValidationErrorsGetFieldErrors(t *testing.T) {
	ves := NewValidationErrors()
	ves.Add(NewValidationError("name", "required", "", "empty"))
	ves.Add(NewValidationError("name", "min_len", "", "too short"))
	ves.Add(NewValidationError("age", "min", "", "too small"))

	nameErrs := ves.GetFieldErrors("name")
	assert.Len(t, nameErrs, 2)

	ageErrs := ves.GetFieldErrors("age")
	assert.Len(t, ageErrs, 1)

	missingErrs := ves.GetFieldErrors("missing")
	assert.Empty(t, missingErrs)
}

// TestMust 测试 Must 函数
func TestMust(t *testing.T) {
	t.Run("nil 错误不 panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Must(nil)
		})
	})

	t.Run("非 nil 错误 panic", func(t *testing.T) {
		err := errors.New("must error")
		assert.PanicsWithValue(t, err, func() {
			Must(err)
		})
	})
}

// TestMustNot 测试 MustNot 函数
func TestMustNot(t *testing.T) {
	t.Run("条件为 false 不 panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MustNot(false, "should not panic")
		})
	})

	t.Run("条件为 true 触发 panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNot(true, "condition is true")
		})
	})
}

// TestAssert 测试 Assert 函数
func TestAssert(t *testing.T) {
	t.Run("条件为 true 返回 nil", func(t *testing.T) {
		assert.NoError(t, Assert(true, "should be nil"))
	})

	t.Run("条件为 false 返回错误", func(t *testing.T) {
		err := Assert(false, "assertion failed")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "assertion failed")
	})
}

// TestRecover 测试 Recover 函数
func TestRecover(t *testing.T) {
	t.Run("无 panic 返回 nil", func(t *testing.T) {
		err := Recover(func() {
			// 正常执行
		})
		assert.NoError(t, err)
	})

	t.Run("panic error 类型返回该错误", func(t *testing.T) {
		expected := errors.New("panic error")
		err := Recover(func() {
			panic(expected)
		})
		assert.Equal(t, expected, err)
	})

	t.Run("panic 非 error 类型返回 InternalError", func(t *testing.T) {
		err := Recover(func() {
			panic("string panic")
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic recovered: string panic")
	})

	t.Run("panic 整数类型", func(t *testing.T) {
		err := Recover(func() {
			panic(42)
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic recovered: 42")
	})
}

// TestWrap 测试 Wrap 函数
func TestWrap(t *testing.T) {
	t.Run("包装 nil 错误返回 nil", func(t *testing.T) {
		assert.Nil(t, Wrap(nil, "context"))
	})

	t.Run("包装普通错误", func(t *testing.T) {
		original := errors.New("original")
		wrapped := Wrap(original, "context")
		assert.Equal(t, "context: original", wrapped.Error())
	})
}

// TestUnwrap 测试 Unwrap 函数
func TestUnwrap(t *testing.T) {
	t.Run("可解包的错误", func(t *testing.T) {
		original := errors.New("original")
		wrapped := fmt.Errorf("context: %w", original)
		unwrapped := Unwrap(wrapped)
		assert.Equal(t, original, unwrapped)
	})

	t.Run("不可解包的错误返回 nil", func(t *testing.T) {
		plain := errors.New("plain")
		assert.Nil(t, Unwrap(plain))
	})
}

// TestIs 测试 Is 函数（递归匹配）
func TestIs(t *testing.T) {
	t.Run("target 为 nil", func(t *testing.T) {
		assert.True(t, Is(nil, nil))
		assert.False(t, Is(errors.New("err"), nil))
	})

	t.Run("err 与 target 相同", func(t *testing.T) {
		err := errors.New("same")
		assert.True(t, Is(err, err))
	})

	t.Run("嵌套包装匹配", func(t *testing.T) {
		target := errors.New("target")
		wrapped := fmt.Errorf("level2: %w", target)
		outer := fmt.Errorf("level1: %w", wrapped)
		assert.True(t, Is(outer, target))
	})

	t.Run("不匹配返回 false", func(t *testing.T) {
		err := errors.New("err")
		target := errors.New("target")
		assert.False(t, Is(err, target))
	})
}
