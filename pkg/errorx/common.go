/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 11:35:00
 * @FilePath: \go-toolbox\pkg\errorx\common.go
 * @Description: 常用错误类型定义和工厂函数
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package errorx

import "fmt"

// 预定义的错误类型
const (
	// 参数错误
	ErrTypeInvalidParam ErrorType = iota + 1000
	ErrTypeMissingParam
	ErrTypeInvalidFormat

	// 业务错误
	ErrTypeNotFound
	ErrTypeAlreadyExists
	ErrTypeConflict
	ErrTypeUnauthorized
	ErrTypeForbidden

	// 系统错误
	ErrTypeInternal
	ErrTypeTimeout
	ErrTypeResourceExhausted
	ErrTypeUnavailable
	ErrTypeNotImplemented

	// 网络错误
	ErrTypeNetworkError
	ErrTypeConnectionLost
	ErrTypeConnectionTimeout

	// 数据错误
	ErrTypeDataCorrupted
	ErrTypeDataNotFound
	ErrTypeDuplicateData

	// 配置错误
	ErrTypeConfigError
	ErrTypeConfigMissing
	ErrTypeConfigInvalid

	// 状态错误
	ErrTypeInvalidState
	ErrTypeConcurrentOperation

	// 事件错误
	ErrTypeHandlerPanic
	ErrTypeHandlerNotFound
	ErrTypeQueueFull
	ErrTypeHandlerTimeout
	ErrTypeInvalidHandler
	ErrTypeInvalidFilter
	ErrTypeInvalidMiddleware
	ErrTypeEventProcessingFailed
)

// 初始化默认错误消息
func init() {
	// 参数错误
	RegisterError(ErrTypeInvalidParam, "invalid parameter: %s")
	RegisterError(ErrTypeMissingParam, "missing required parameter: %s")
	RegisterError(ErrTypeInvalidFormat, "invalid format: %s")

	// 业务错误
	RegisterError(ErrTypeNotFound, "resource not found: %s")
	RegisterError(ErrTypeAlreadyExists, "resource already exists: %s")
	RegisterError(ErrTypeConflict, "resource conflict: %s")
	RegisterError(ErrTypeUnauthorized, "unauthorized: %s")
	RegisterError(ErrTypeForbidden, "forbidden: %s")

	// 系统错误
	RegisterError(ErrTypeInternal, "internal error: %s")
	RegisterError(ErrTypeTimeout, "operation timeout: %s")
	RegisterError(ErrTypeResourceExhausted, "resource exhausted: %s")
	RegisterError(ErrTypeUnavailable, "service unavailable: %s")
	RegisterError(ErrTypeNotImplemented, "not implemented: %s")

	// 网络错误
	RegisterError(ErrTypeNetworkError, "network error: %s")
	RegisterError(ErrTypeConnectionLost, "connection lost: %s")
	RegisterError(ErrTypeConnectionTimeout, "connection timeout: %s")

	// 数据错误
	RegisterError(ErrTypeDataCorrupted, "data corrupted: %s")
	RegisterError(ErrTypeDataNotFound, "data not found: %s")
	RegisterError(ErrTypeDuplicateData, "duplicate data: %s")

	// 配置错误
	RegisterError(ErrTypeConfigError, "configuration error: %s")
	RegisterError(ErrTypeConfigMissing, "missing configuration: %s")
	RegisterError(ErrTypeConfigInvalid, "invalid configuration: %s")

	// 状态错误
	RegisterError(ErrTypeInvalidState, "invalid state: %s")
	RegisterError(ErrTypeConcurrentOperation, "concurrent operation error: %s")

	// 事件错误
	RegisterError(ErrTypeHandlerPanic, "handler panic: %s")
	RegisterError(ErrTypeHandlerNotFound, "handler not found: %s")
	RegisterError(ErrTypeQueueFull, "queue is full: %s")
	RegisterError(ErrTypeHandlerTimeout, "handler execution timeout: %s")
	RegisterError(ErrTypeInvalidHandler, "invalid handler: %s")
	RegisterError(ErrTypeInvalidFilter, "invalid filter: %s")
	RegisterError(ErrTypeInvalidMiddleware, "invalid middleware: %s")
	RegisterError(ErrTypeEventProcessingFailed, "processing failed: %s")
}

// CustomError 自定义错误结构，包含错误码和详细信息
type CustomError struct {
	BaseError
	Code    ErrorType
	Details map[string]interface{}
}

// NewCustomError 创建自定义错误
func NewCustomError(code ErrorType, message string, details map[string]interface{}) *CustomError {
	return &CustomError{
		BaseError: NewBaseError(message),
		Code:      code,
		Details:   details,
	}
}

// Error 实现error接口
func (e *CustomError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("%s (details: %+v)", e.BaseError.Error(), e.Details)
	}
	return e.BaseError.Error()
}

// GetCode 获取错误码
func (e *CustomError) GetCode() ErrorType {
	return e.Code
}

// GetDetails 获取错误详情
func (e *CustomError) GetDetails() map[string]interface{} {
	return e.Details
}

// IsType 检查错误是否为指定类型
func IsType(err error, errType ErrorType) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Code == errType
	}
	return false
}

// 参数错误工厂函数
func NewInvalidParamError(param string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeInvalidParam, param),
		Code:      ErrTypeInvalidParam,
		Details:   map[string]interface{}{"parameter": param},
	}
}

func NewMissingParamError(param string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeMissingParam, param),
		Code:      ErrTypeMissingParam,
		Details:   map[string]interface{}{"parameter": param},
	}
}

func NewInvalidFormatError(format string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeInvalidFormat, format),
		Code:      ErrTypeInvalidFormat,
		Details:   map[string]interface{}{"format": format},
	}
}

// 业务错误工厂函数
func NewNotFoundError(resource string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeNotFound, resource),
		Code:      ErrTypeNotFound,
		Details:   map[string]interface{}{"resource": resource},
	}
}

func NewAlreadyExistsError(resource string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeAlreadyExists, resource),
		Code:      ErrTypeAlreadyExists,
		Details:   map[string]interface{}{"resource": resource},
	}
}

func NewConflictError(resource string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeConflict, resource),
		Code:      ErrTypeConflict,
		Details:   map[string]interface{}{"resource": resource},
	}
}

func NewUnauthorizedError(reason string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeUnauthorized, reason),
		Code:      ErrTypeUnauthorized,
		Details:   map[string]interface{}{"reason": reason},
	}
}

func NewForbiddenError(reason string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeForbidden, reason),
		Code:      ErrTypeForbidden,
		Details:   map[string]interface{}{"reason": reason},
	}
}

// 系统错误工厂函数
func NewInternalError(message string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeInternal, message),
		Code:      ErrTypeInternal,
		Details:   map[string]interface{}{"message": message},
	}
}

func NewTimeoutError(operation string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeTimeout, operation),
		Code:      ErrTypeTimeout,
		Details:   map[string]interface{}{"operation": operation},
	}
}

func NewResourceExhaustedError(resource string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeResourceExhausted, resource),
		Code:      ErrTypeResourceExhausted,
		Details:   map[string]interface{}{"resource": resource},
	}
}

func NewUnavailableError(service string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeUnavailable, service),
		Code:      ErrTypeUnavailable,
		Details:   map[string]interface{}{"service": service},
	}
}

func NewNotImplementedError(feature string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeNotImplemented, feature),
		Code:      ErrTypeNotImplemented,
		Details:   map[string]interface{}{"feature": feature},
	}
}

// 网络错误工厂函数
func NewNetworkError(message string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeNetworkError, message),
		Code:      ErrTypeNetworkError,
		Details:   map[string]interface{}{"message": message},
	}
}

func NewConnectionLostError(target string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeConnectionLost, target),
		Code:      ErrTypeConnectionLost,
		Details:   map[string]interface{}{"target": target},
	}
}

func NewConnectionTimeoutError(target string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeConnectionTimeout, target),
		Code:      ErrTypeConnectionTimeout,
		Details:   map[string]interface{}{"target": target},
	}
}

// 数据错误工厂函数
func NewDataCorruptedError(data string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeDataCorrupted, data),
		Code:      ErrTypeDataCorrupted,
		Details:   map[string]interface{}{"data": data},
	}
}

func NewDataNotFoundError(data string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeDataNotFound, data),
		Code:      ErrTypeDataNotFound,
		Details:   map[string]interface{}{"data": data},
	}
}

func NewDuplicateDataError(data string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeDuplicateData, data),
		Code:      ErrTypeDuplicateData,
		Details:   map[string]interface{}{"data": data},
	}
}

// 配置错误工厂函数
func NewConfigError(config string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeConfigError, config),
		Code:      ErrTypeConfigError,
		Details:   map[string]interface{}{"config": config},
	}
}

func NewConfigMissingError(config string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeConfigMissing, config),
		Code:      ErrTypeConfigMissing,
		Details:   map[string]interface{}{"config": config},
	}
}

func NewConfigInvalidError(config string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeConfigInvalid, config),
		Code:      ErrTypeConfigInvalid,
		Details:   map[string]interface{}{"config": config},
	}
}

// 状态错误工厂函数
func NewInvalidStateError(state string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeInvalidState, state),
		Code:      ErrTypeInvalidState,
		Details:   map[string]interface{}{"state": state},
	}
}

func NewConcurrentOperationError(operation string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeConcurrentOperation, operation),
		Code:      ErrTypeConcurrentOperation,
		Details:   map[string]interface{}{"operation": operation},
	}
}

// 事件错误工厂函数
func NewHandlerPanicError(handler string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeHandlerPanic, handler),
		Code:      ErrTypeHandlerPanic,
		Details:   map[string]interface{}{"handler": handler},
	}
}

func NewHandlerNotFoundError(handler string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeHandlerNotFound, handler),
		Code:      ErrTypeHandlerNotFound,
		Details:   map[string]interface{}{"handler": handler},
	}
}

func NewQueueFullError(queue string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeQueueFull, queue),
		Code:      ErrTypeQueueFull,
		Details:   map[string]interface{}{"queue": queue},
	}
}

func NewHandlerTimeoutError(handler string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeHandlerTimeout, handler),
		Code:      ErrTypeHandlerTimeout,
		Details:   map[string]interface{}{"handler": handler},
	}
}

func NewInvalidHandlerError(handler string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeInvalidHandler, handler),
		Code:      ErrTypeInvalidHandler,
		Details:   map[string]interface{}{"handler": handler},
	}
}

func NewInvalidFilterError(filter string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeInvalidFilter, filter),
		Code:      ErrTypeInvalidFilter,
		Details:   map[string]interface{}{"filter": filter},
	}
}

func NewInvalidMiddlewareError(middleware string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeInvalidMiddleware, middleware),
		Code:      ErrTypeInvalidMiddleware,
		Details:   map[string]interface{}{"middleware": middleware},
	}
}

func NewEventProcessingFailedError(event string) error {
	return &CustomError{
		BaseError: NewError(ErrTypeEventProcessingFailed, event),
		Code:      ErrTypeEventProcessingFailed,
		Details:   map[string]interface{}{"event": event},
	}
}

// 错误检查函数
func IsInvalidParamError(err error) bool {
	return IsType(err, ErrTypeInvalidParam)
}

func IsNotFoundError(err error) bool {
	return IsType(err, ErrTypeNotFound)
}

func IsTimeoutError(err error) bool {
	return IsType(err, ErrTypeTimeout)
}

func IsResourceExhaustedError(err error) bool {
	return IsType(err, ErrTypeResourceExhausted)
}

func IsNetworkError(err error) bool {
	return IsType(err, ErrTypeNetworkError) ||
		IsType(err, ErrTypeConnectionLost) ||
		IsType(err, ErrTypeConnectionTimeout)
}

// 状态错误检查函数
func IsInvalidStateError(err error) bool {
	return IsType(err, ErrTypeInvalidState)
}

func IsConcurrentOperationError(err error) bool {
	return IsType(err, ErrTypeConcurrentOperation)
}

// 事件错误检查函数
func IsHandlerPanicError(err error) bool {
	return IsType(err, ErrTypeHandlerPanic)
}

func IsHandlerNotFoundError(err error) bool {
	return IsType(err, ErrTypeHandlerNotFound)
}

func IsQueueFullError(err error) bool {
	return IsType(err, ErrTypeQueueFull)
}

func IsHandlerTimeoutError(err error) bool {
	return IsType(err, ErrTypeHandlerTimeout)
}

func IsInvalidHandlerError(err error) bool {
	return IsType(err, ErrTypeInvalidHandler)
}

func IsInvalidFilterError(err error) bool {
	return IsType(err, ErrTypeInvalidFilter)
}

func IsInvalidMiddlewareError(err error) bool {
	return IsType(err, ErrTypeInvalidMiddleware)
}

func IsEventProcessingFailedError(err error) bool {
	return IsType(err, ErrTypeEventProcessingFailed)
}

// 错误转换函数
func ToCustomError(err error, fallbackType ErrorType) *CustomError {
	if customErr, ok := err.(*CustomError); ok {
		return customErr
	}

	return &CustomError{
		BaseError: NewBaseError(err.Error()),
		Code:      fallbackType,
		Details:   map[string]interface{}{"original_error": err.Error()},
	}
}

// 批量错误处理
type ErrorCollector struct {
	errors []error
}

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]error, 0),
	}
}

func (c *ErrorCollector) Add(err error) {
	if err != nil {
		c.errors = append(c.errors, err)
	}
}

func (c *ErrorCollector) HasErrors() bool {
	return len(c.errors) > 0
}

func (c *ErrorCollector) GetErrors() []error {
	return c.errors
}

func (c *ErrorCollector) Error() string {
	if len(c.errors) == 0 {
		return ""
	}

	if len(c.errors) == 1 {
		return c.errors[0].Error()
	}

	var result string
	for i, err := range c.errors {
		if i > 0 {
			result += "; "
		}
		result += err.Error()
	}

	return fmt.Sprintf("multiple errors: %s", result)
}
