/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-30 16:44:10
 * @FilePath: \go-toolbox\pkg\errorx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package errorx

import (
	"context"
	"fmt"
	"sync"
)

// WrapError 是一个通用的错误包装函数
func WrapError(message string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// 定义 BaseError 结构体
type BaseError struct {
	ctx context.Context
	msg string
}

// 错误类型常量
type ErrorType int

// 错误映射类型
type ErrorMapType map[ErrorType]string

// 错误消息映射
var (
	errorMessages   = make(map[ErrorType]string) // 初始化映射
	defaultErrorMap = make(ErrorMapType)         // 初始化默认错误映射
	mu              sync.Mutex                   // 保护并发访问
)

// NewBaseError 创建一个新的 BaseError 实例
func NewBaseError(msg string) BaseError {
	return BaseError{msg: msg}
}

// NewBaseErrorWithCtx 创建一个新的 BaseError Ctx 实例
func NewBaseErrorWithCtx(ctx context.Context, msg string) BaseError {
	return BaseError{ctx: ctx, msg: msg}
}

// Error 实现 error 接口，返回错误信息
func (e BaseError) Error() string {
	return e.msg
}

// RegisterError 注册错误类型和消息
func RegisterError(errType ErrorType, msg string) {
	mu.Lock()
	defer mu.Unlock()

	// 检查是否已经注册过该错误类型
	if _, exists := errorMessages[errType]; exists {
		fmt.Printf("ErrorType %d is already registered\n", errType)
		return
	}

	errorMessages[errType] = msg
	defaultErrorMap[errType] = msg // 将错误类型和消息添加到 defaultErrorMap
}

// NewError 创建一个新的错误实例
func NewError(errType ErrorType, args ...interface{}) BaseError {
	mu.Lock()
	defer mu.Unlock()

	if msg, ok := errorMessages[errType]; ok {
		return NewBaseError(fmt.Sprintf(msg, args...))
	}
	return NewBaseError("unknown error")
}

// 打印错误映射（调试用）
func PrintErrorMap() {
	mu.Lock()
	defer mu.Unlock()
	for errType, msg := range defaultErrorMap {
		fmt.Printf("ErrorType: %d, Message: %s\n", errType, msg)
	}
}

// GetErrorMap 返回当前错误映射
func GetErrorMap() ErrorMapType {
	mu.Lock()
	defer mu.Unlock()
	return defaultErrorMap // 返回错误映射
}

// ResetErrorMap 重置错误映射
func ResetErrorMap() {
	mu.Lock()
	defer mu.Unlock()
	defaultErrorMap = make(ErrorMapType) // 重置错误映射
}
