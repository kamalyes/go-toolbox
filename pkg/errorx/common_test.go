/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-10 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-10 00:00:00
 * @FilePath: \go-toolbox\pkg\errorx\common_test.go
 * @Description: common.go 的单元测试，覆盖自定义错误、工厂函数、错误检查、转换和收集器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewCustomError 测试创建自定义错误
func TestNewCustomError(t *testing.T) {
	details := map[string]interface{}{"key": "value"}
	ce := NewCustomError(ErrTypeInvalidParam, "custom message", details)
	assert.NotNil(t, ce)
	assert.Equal(t, ErrorType(ErrTypeInvalidParam), ce.Code)
	assert.Equal(t, details, ce.Details)
	assert.Equal(t, "custom message", ce.BaseError.Error())
}

// TestCustomErrorError 测试 CustomError 的 Error 方法
func TestCustomErrorError(t *testing.T) {
	t.Run("有详情时输出详情", func(t *testing.T) {
		ce := NewCustomError(ErrTypeInternal, "msg", map[string]interface{}{"k": "v"})
		errStr := ce.Error()
		assert.Contains(t, errStr, "msg")
		assert.Contains(t, errStr, "details:")
		assert.Contains(t, errStr, "k")
	})

	t.Run("无详情时只输出消息", func(t *testing.T) {
		ce := NewCustomError(ErrTypeInternal, "msg", nil)
		assert.Equal(t, "msg", ce.Error())
	})
}

// TestCustomErrorGetCode 测试 GetCode 方法
func TestCustomErrorGetCode(t *testing.T) {
	ce := NewCustomError(ErrTypeTimeout, "msg", nil)
	assert.Equal(t, ErrorType(ErrTypeTimeout), ce.GetCode())
}

// TestCustomErrorGetDetails 测试 GetDetails 方法
func TestCustomErrorGetDetails(t *testing.T) {
	details := map[string]interface{}{"a": 1}
	ce := NewCustomError(ErrTypeInternal, "msg", details)
	assert.Equal(t, details, ce.GetDetails())
}

// TestIsType 测试 IsType 函数
func TestIsType(t *testing.T) {
	t.Run("类型匹配返回 true", func(t *testing.T) {
		err := NewInvalidParamError("param")
		assert.True(t, IsType(err, ErrTypeInvalidParam))
	})

	t.Run("类型不匹配返回 false", func(t *testing.T) {
		err := NewInvalidParamError("param")
		assert.False(t, IsType(err, ErrTypeNotFound))
	})

	t.Run("非 CustomError 类型返回 false", func(t *testing.T) {
		err := errors.New("plain error")
		assert.False(t, IsType(err, ErrTypeInvalidParam))
	})
}

// TestErrorFactoryFunctions 测试所有错误工厂函数
func TestErrorFactoryFunctions(t *testing.T) {
	t.Run("参数错误工厂", func(t *testing.T) {
		err := NewInvalidParamError("userId")
		assert.NotNil(t, err)
		assert.True(t, IsInvalidParamError(err))
		assert.Equal(t, ErrorType(ErrTypeInvalidParam), err.(*CustomError).GetCode())

		err = NewMissingParamError("token")
		assert.NotNil(t, err)
		assert.Equal(t, ErrorType(ErrTypeMissingParam), err.(*CustomError).GetCode())

		err = NewInvalidFormatError("email")
		assert.NotNil(t, err)
		assert.Equal(t, ErrorType(ErrTypeInvalidFormat), err.(*CustomError).GetCode())
	})

	t.Run("业务错误工厂", func(t *testing.T) {
		err := NewNotFoundError("user")
		assert.NotNil(t, err)
		assert.True(t, IsNotFoundError(err))
		assert.Equal(t, ErrorType(ErrTypeNotFound), err.(*CustomError).GetCode())

		err = NewAlreadyExistsError("user")
		assert.Equal(t, ErrorType(ErrTypeAlreadyExists), err.(*CustomError).GetCode())

		err = NewConflictError("user")
		assert.Equal(t, ErrorType(ErrTypeConflict), err.(*CustomError).GetCode())

		err = NewUnauthorizedError("invalid token")
		assert.Equal(t, ErrorType(ErrTypeUnauthorized), err.(*CustomError).GetCode())

		err = NewForbiddenError("no permission")
		assert.Equal(t, ErrorType(ErrTypeForbidden), err.(*CustomError).GetCode())
	})

	t.Run("系统错误工厂", func(t *testing.T) {
		err := NewInternalError("boom")
		assert.NotNil(t, err)
		assert.Equal(t, ErrorType(ErrTypeInternal), err.(*CustomError).GetCode())

		err = NewTimeoutError("db query")
		assert.True(t, IsTimeoutError(err))
		assert.Equal(t, ErrorType(ErrTypeTimeout), err.(*CustomError).GetCode())

		err = NewResourceExhaustedError("memory")
		assert.True(t, IsResourceExhaustedError(err))
		assert.Equal(t, ErrorType(ErrTypeResourceExhausted), err.(*CustomError).GetCode())

		err = NewUnavailableError("redis")
		assert.Equal(t, ErrorType(ErrTypeUnavailable), err.(*CustomError).GetCode())

		err = NewNotImplementedError("feature")
		assert.Equal(t, ErrorType(ErrTypeNotImplemented), err.(*CustomError).GetCode())
	})

	t.Run("网络错误工厂", func(t *testing.T) {
		err := NewNetworkError("timeout")
		assert.True(t, IsNetworkError(err))
		assert.Equal(t, ErrorType(ErrTypeNetworkError), err.(*CustomError).GetCode())

		err = NewConnectionLostError("db")
		assert.True(t, IsNetworkError(err))
		assert.Equal(t, ErrorType(ErrTypeConnectionLost), err.(*CustomError).GetCode())

		err = NewConnectionTimeoutError("db")
		assert.True(t, IsNetworkError(err))
		assert.Equal(t, ErrorType(ErrTypeConnectionTimeout), err.(*CustomError).GetCode())
	})

	t.Run("数据错误工厂", func(t *testing.T) {
		err := NewDataCorruptedError("record")
		assert.Equal(t, ErrorType(ErrTypeDataCorrupted), err.(*CustomError).GetCode())

		err = NewDataNotFoundError("record")
		assert.Equal(t, ErrorType(ErrTypeDataNotFound), err.(*CustomError).GetCode())

		err = NewDuplicateDataError("record")
		assert.Equal(t, ErrorType(ErrTypeDuplicateData), err.(*CustomError).GetCode())
	})

	t.Run("配置错误工厂", func(t *testing.T) {
		err := NewConfigError("db.host")
		assert.Equal(t, ErrorType(ErrTypeConfigError), err.(*CustomError).GetCode())

		err = NewConfigMissingError("db.port")
		assert.Equal(t, ErrorType(ErrTypeConfigMissing), err.(*CustomError).GetCode())

		err = NewConfigInvalidError("db.url")
		assert.Equal(t, ErrorType(ErrTypeConfigInvalid), err.(*CustomError).GetCode())
	})

	t.Run("状态错误工厂", func(t *testing.T) {
		err := NewInvalidStateError("running")
		assert.True(t, IsInvalidStateError(err))
		assert.Equal(t, ErrorType(ErrTypeInvalidState), err.(*CustomError).GetCode())

		err = NewConcurrentOperationError("update")
		assert.True(t, IsConcurrentOperationError(err))
		assert.Equal(t, ErrorType(ErrTypeConcurrentOperation), err.(*CustomError).GetCode())
	})

	t.Run("事件错误工厂", func(t *testing.T) {
		err := NewHandlerPanicError("handler")
		assert.True(t, IsHandlerPanicError(err))
		assert.Equal(t, ErrorType(ErrTypeHandlerPanic), err.(*CustomError).GetCode())

		err = NewHandlerNotFoundError("handler")
		assert.True(t, IsHandlerNotFoundError(err))
		assert.Equal(t, ErrorType(ErrTypeHandlerNotFound), err.(*CustomError).GetCode())

		err = NewQueueFullError("queue")
		assert.True(t, IsQueueFullError(err))
		assert.Equal(t, ErrorType(ErrTypeQueueFull), err.(*CustomError).GetCode())

		err = NewHandlerTimeoutError("handler")
		assert.True(t, IsHandlerTimeoutError(err))
		assert.Equal(t, ErrorType(ErrTypeHandlerTimeout), err.(*CustomError).GetCode())

		err = NewInvalidHandlerError("handler")
		assert.True(t, IsInvalidHandlerError(err))
		assert.Equal(t, ErrorType(ErrTypeInvalidHandler), err.(*CustomError).GetCode())

		err = NewInvalidFilterError("filter")
		assert.True(t, IsInvalidFilterError(err))
		assert.Equal(t, ErrorType(ErrTypeInvalidFilter), err.(*CustomError).GetCode())

		err = NewInvalidMiddlewareError("middleware")
		assert.True(t, IsInvalidMiddlewareError(err))
		assert.Equal(t, ErrorType(ErrTypeInvalidMiddleware), err.(*CustomError).GetCode())

		err = NewEventProcessingFailedError("event")
		assert.True(t, IsEventProcessingFailedError(err))
		assert.Equal(t, ErrorType(ErrTypeEventProcessingFailed), err.(*CustomError).GetCode())
	})
}

// TestErrorCheckFunctionsOnNonCustomError 测试错误检查函数对非 CustomError 的处理
func TestErrorCheckFunctionsOnNonCustomError(t *testing.T) {
	plain := errors.New("plain")
	assert.False(t, IsInvalidParamError(plain))
	assert.False(t, IsNotFoundError(plain))
	assert.False(t, IsTimeoutError(plain))
	assert.False(t, IsResourceExhaustedError(plain))
	assert.False(t, IsNetworkError(plain))
	assert.False(t, IsInvalidStateError(plain))
	assert.False(t, IsConcurrentOperationError(plain))
	assert.False(t, IsHandlerPanicError(plain))
	assert.False(t, IsHandlerNotFoundError(plain))
	assert.False(t, IsQueueFullError(plain))
	assert.False(t, IsHandlerTimeoutError(plain))
	assert.False(t, IsInvalidHandlerError(plain))
	assert.False(t, IsInvalidFilterError(plain))
	assert.False(t, IsInvalidMiddlewareError(plain))
	assert.False(t, IsEventProcessingFailedError(plain))
}

// TestIsNetworkErrorMultipleTypes 测试 IsNetworkError 匹配多种网络错误
func TestIsNetworkErrorMultipleTypes(t *testing.T) {
	// 测试不是网络错误的其他类型
	err := NewNotFoundError("x")
	assert.False(t, IsNetworkError(err))
}

// TestToCustomError 测试错误转换为 CustomError
func TestToCustomError(t *testing.T) {
	t.Run("已经是 CustomError 直接返回", func(t *testing.T) {
		original := NewInvalidParamError("param").(*CustomError)
		result := ToCustomError(original, ErrTypeInternal)
		assert.Same(t, original, result)
	})

	t.Run("普通错误转换为 CustomError", func(t *testing.T) {
		plain := errors.New("plain error")
		result := ToCustomError(plain, ErrTypeInternal)
		require.NotNil(t, result)
		assert.Equal(t, ErrorType(ErrTypeInternal), result.Code)
		assert.Equal(t, "plain error", result.BaseError.Error())
		assert.Equal(t, "plain error", result.Details["original_error"])
	})
}

// TestErrorCollector 测试错误收集器
func TestErrorCollector(t *testing.T) {
	t.Run("空收集器", func(t *testing.T) {
		c := NewErrorCollector()
		assert.NotNil(t, c)
		assert.False(t, c.HasErrors())
		assert.Empty(t, c.GetErrors())
		assert.Empty(t, c.Error())
	})

	t.Run("添加 nil 错误被忽略", func(t *testing.T) {
		c := NewErrorCollector()
		c.Add(nil)
		assert.False(t, c.HasErrors())
	})

	t.Run("单个错误", func(t *testing.T) {
		c := NewErrorCollector()
		c.Add(errors.New("single"))
		assert.True(t, c.HasErrors())
		assert.Len(t, c.GetErrors(), 1)
		assert.Equal(t, "single", c.Error())
	})

	t.Run("多个错误", func(t *testing.T) {
		c := NewErrorCollector()
		c.Add(errors.New("first"))
		c.Add(errors.New("second"))
		errStr := c.Error()
		assert.Contains(t, errStr, "multiple errors:")
		assert.Contains(t, errStr, "first")
		assert.Contains(t, errStr, "second")
		assert.Contains(t, errStr, "; ")
	})
}
