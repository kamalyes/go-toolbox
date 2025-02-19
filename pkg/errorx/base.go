/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-15 09:19:15
 * @FilePath: \go-toolbox\pkg\errorx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package errorx

import (
	"fmt"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
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
	Msg string
}

// 错误类型常量
type ErrorType int

// 错误映射类型
type ErrorMapType map[ErrorType]string

// 错误消息映射
var (
	errorMessages   = make(map[ErrorType]string) // 初始化映射
	defaultErrorMap = make(ErrorMapType)         // 初始化默认错误映射
	mu              sync.Mutex                   // 互斥锁、保护并发访问
)

// NewBaseError 创建一个新的 BaseError 实例
func NewBaseError(msg string) BaseError {
	return BaseError{Msg: msg}
}

// Error 实现 error 接口，返回错误信息
func (e BaseError) Error() string {
	return e.Msg
}

// RegisterError 注册错误类型和消息
func RegisterError(errType ErrorType, msg string) {
	syncx.WithLock(&mu, func() {
		// 检查是否已经注册过该错误类型
		if _, exists := errorMessages[errType]; exists {
			fmt.Printf("ErrorType %d is already registered\n", errType)
			return
		}
		errorMessages[errType] = msg
		defaultErrorMap[errType] = msg // 将错误类型和消息添加到 defaultErrorMap
	})
}

// NewError 创建一个新的错误实例
func NewError(errType ErrorType, args ...interface{}) BaseError {
	var result BaseError
	syncx.WithLock(&mu, func() {
		if msg, ok := errorMessages[errType]; ok {
			result = NewBaseError(fmt.Sprintf(msg, args...))
		} else {
			result = NewBaseError("unknown error")
		}
	})
	return result
}

// 打印错误映射（调试用）
func PrintErrorMap() {
	syncx.WithLock(&mu, func() {
		for errType, msg := range defaultErrorMap {
			fmt.Printf("ErrorType: %d, Message: %s\n", errType, msg)
		}
	})
}

// GetErrorMap 返回当前错误映射
func GetErrorMap() ErrorMapType {
	result := make(ErrorMapType)
	syncx.WithLock(&mu, func() {
		// 深拷贝以避免数据竞争
		for k, v := range defaultErrorMap {
			result[k] = v
		}
	})
	return result // 返回深拷贝的错误映射
}

// ResetErrorMap 重置错误映射
func ResetErrorMap() {
	syncx.WithLock(&mu, func() {
		errorMessages = make(map[ErrorType]string) // 重置错误类型映射
		defaultErrorMap = make(ErrorMapType)       // 重置默认错误映射
	})
}
