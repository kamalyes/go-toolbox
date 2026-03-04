/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 10:15:15
 * @FilePath: \go-toolbox\pkg\contextx\helpers.go
 * @Description: Context 全局辅助函数
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package contextx

import (
	"context"
	"fmt"
	"time"
)

// WithTimeout 创建带超时的 context 并执行函数
// 自动处理 cancel 调用,并监听 context 超时
//
// 使用示例:
//
//	err := contextx.WithTimeout(2*time.Second, func(ctx context.Context) error {
//	    return repo.SaveData(ctx, data)
//	})
func WithTimeout(timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// WithTimeoutValue 创建带超时的 context 并执行函数,返回结果和错误
// 自动处理 cancel 调用,并监听 context 超时
//
// 使用示例:
//
//	result, err := contextx.WithTimeoutValue(2*time.Second, func(ctx context.Context) (int, error) {
//	    return repo.GetCount(ctx)
//	})
func WithTimeoutValue[T any](timeout time.Duration, fn func(context.Context) (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type result struct {
		value T
		err   error
	}

	done := make(chan result, 1)
	go func() {
		v, e := fn(ctx)
		done <- result{value: v, err: e}
	}()

	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case res := <-done:
		return res.value, res.err
	}
}

// OrBackground 如果ctx已取消则返回context.Background()，否则返回ctx本身
// 用于在关闭流程中确保异步任务仍能完成
func OrBackground(ctx context.Context) context.Context {
	select {
	case <-ctx.Done():
		return context.Background()
	default:
		return ctx
	}
}

// WithTimeoutFrom 在指定的父context基础上创建带超时的context并执行函数
// 自动处理 cancel 调用,并监听 context 超时
//
// 使用示例:
//
//	err := contextx.WithTimeoutFrom(parentCtx, 2*time.Second, func(ctx context.Context) error {
//	    return repo.SaveData(ctx, data)
//	})
func WithTimeoutFrom(parent context.Context, timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// WithTimeoutOrBackground 在父context基础上创建带超时的context并执行函数
// 如果父context已取消，则使用Background作为父context，确保异步任务仍能完成
//
// 使用示例:
//
//	err := contextx.WithTimeoutOrBackground(h.ctx, 2*time.Second, func(ctx context.Context) error {
//	    return repo.SaveData(ctx, data)
//	})
func WithTimeoutOrBackground(parent context.Context, timeout time.Duration, fn func(context.Context) error) error {
	return WithTimeoutFrom(OrBackground(parent), timeout, fn)
}

// WithTimeoutDecorators 创建带超时的标准 context 并应用装饰器
// 支持可选的装饰器函数来增强 context
//
// 使用示例:
//
//	ctx, cancel := contextx.WithTimeoutDecorators(5*time.Second)
//	defer cancel()
//
//	或者带装饰器
//	ctx, cancel := contextx.WithTimeoutDecorators(5*time.Second, WithUser(user), WithTraceID(id))
//	defer cancel()
func WithTimeoutDecorators(timeout time.Duration, decorators ...func(context.Context) context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	for _, decorator := range decorators {
		ctx = decorator(ctx)
	}
	return ctx, cancel
}

// WithDeadlineDecorators 创建带截止时间的 context 并应用装饰器
// 支持可选的装饰器函数来增强 context
//
// 使用示例:
//
//	deadline := time.Now().Add(10 * time.Second)
//	ctx, cancel := contextx.WithDeadlineDecorators(deadline)
//	defer cancel()
//
//	或者带装饰器
//	ctx, cancel := contextx.WithDeadlineDecorators(deadline, WithUser(user))
//	defer cancel()
func WithDeadlineDecorators(deadline time.Time, decorators ...func(context.Context) context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	for _, decorator := range decorators {
		ctx = decorator(ctx)
	}
	return ctx, cancel
}

// MustGet 从上下文中获取值，如果不存在则 panic
// 这是一个通用的 Must 模式实现
//
// 使用示例:
//
//	config := contextx.MustGet[*Config](ctx, "config")
//	userID := contextx.MustGet[string](ctx, "user_id")
func MustGet[T any](ctx context.Context, key any) T {
	value := ctx.Value(key)
	if value == nil {
		panic(fmt.Sprintf("value for key %v not found in context", key))
	}

	result, ok := value.(T)
	if !ok {
		panic(fmt.Sprintf("value for key %v has wrong type: expected %T, got %T", key, *new(T), value))
	}

	return result
}

// MustGetWithMessage 从上下文中获取值，如果不存在则 panic 并显示自定义消息
//
// 使用示例:
//
//	config := contextx.MustGetWithMessage[*Config](ctx, "config", "配置信息不存在于上下文中")
func MustGetWithMessage[T any](ctx context.Context, key any, message string) T {
	value := ctx.Value(key)
	if value == nil {
		panic(message)
	}

	result, ok := value.(T)
	if !ok {
		panic(fmt.Sprintf("%s (type mismatch: expected %T, got %T)", message, *new(T), value))
	}

	return result
}

// GetOrDefault 从上下文中获取值，如果不存在则返回默认值
//
// 使用示例:
//
//	timeout := contextx.GetOrDefault(ctx, "timeout", 30*time.Second)
//	userID := contextx.GetOrDefault(ctx, "user_id", "anonymous")
func GetOrDefault[T any](ctx context.Context, key any, defaultValue T) T {
	value := ctx.Value(key)
	if value == nil {
		return defaultValue
	}

	result, ok := value.(T)
	if !ok {
		return defaultValue
	}

	return result
}
