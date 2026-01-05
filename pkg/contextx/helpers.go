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
