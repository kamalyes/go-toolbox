/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-21 11:57:55
 * @FilePath: \go-toolbox\pkg\contextx\context.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package contextx

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// Context 是一个自定义的上下文，支持多个值的存储
type Context struct {
	mu     sync.RWMutex
	values map[interface{}]interface{}
	context.Context
	cancelFunc context.CancelFunc // 添加取消函数
	pool       *syncx.LimitedPool // 引入字节切片池
}

// NewContext 创建一个新的 Context，允许用户传入自定义的字节切片池
func NewContext(parent context.Context, pool *syncx.LimitedPool) *Context {
	if pool == nil {
		pool = syncx.NewLimitedPool(32, 1024)
	}
	return &Context{
		values:  make(map[interface{}]interface{}),
		Context: parent,
		pool:    pool,
	}
}

// NewContextWithTimeout 创建一个带有超时的 Context
func NewContextWithTimeout(parent context.Context, timeout time.Duration, pool *syncx.LimitedPool) *Context {
	ctx, cancel := context.WithTimeout(parent, timeout)
	return &Context{
		values:     make(map[interface{}]interface{}),
		Context:    ctx,
		cancelFunc: cancel,
		pool:       pool,
	}
}

// NewContextWithCancel 创建一个可以手动取消的 Context
func NewContextWithCancel(parent context.Context, pool *syncx.LimitedPool) *Context {
	ctx, cancel := context.WithCancel(parent)
	return &Context{
		values:     make(map[interface{}]interface{}),
		Context:    ctx,
		cancelFunc: cancel,
		pool:       pool,
	}
}

// NewContextWithValue 在父上下文中设置值并返回新的 Context
func NewContextWithValue(parent context.Context, key, val interface{}, pool *syncx.LimitedPool) (*Context, error) {
	ctx := NewContext(parent, pool)
	if err := ctx.WithValue(key, val); err != nil {
		return nil, err
	}
	return ctx, nil
}

// NewLocalContextWithValue 在当前 Context 中设置局部值
func NewLocalContextWithValue(ctx *Context, key, val interface{}) (*Context, error) {
	return ctx, ctx.WithValue(key, val)
}

// Values 返回当前上下文中所有的键值对
func (c *Context) Values() map[interface{}]interface{} {
	return syncx.WithRLockReturnValue(&c.mu, func() map[interface{}]interface{} {
		// 创建一个新的映射以返回
		valuesCopy := make(map[interface{}]interface{}, len(c.values))
		for k, v := range c.values {
			valuesCopy[k] = v
		}
		return valuesCopy
	})
}

// validateKey 检查键是否有效
func validateKey(key interface{}) error {
	if key == nil {
		return fmt.Errorf("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		return fmt.Errorf("key is not comparable")
	}
	return nil
}

// Value 获取指定键的值
func (c *Context) Value(key interface{}) interface{} {
	return syncx.WithRLockReturnValue(&c.mu, func() any {

		if v, ok := c.values[key]; ok {
			return v
		}
		return c.Context.Value(key)
	})
}

// WithByteSlice 处理字节切片的存储
func (c *Context) WithByteSlice(key interface{}, value []byte) error {
	if buf := c.pool.Get(len(value)); buf != nil {
		copy(*buf, value)
		c.values[key] = buf
		return nil
	}
	c.values[key] = value
	return nil
}

// WithValue 设置指定键的值并返回错误
func (c *Context) WithValue(key, value interface{}) error {
	if err := validateKey(key); err != nil {
		return err
	}

	if byteSlice, ok := value.([]byte); ok {
		return c.WithByteSlice(key, byteSlice)
	}

	return syncx.WithLockReturnValue(&c.mu, func() error {
		c.values[key] = value
		return nil
	})
}

// Remove 删除指定键的键值对
func (c *Context) Remove(key interface{}) {
	syncx.WithLock(&c.mu, func() {
		delete(c.values, key)
	})
}

// Cancel 取消上下文
func (c *Context) Cancel() {
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
}

// Deadline 返回上下文的截止时间
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Deadline()
}

// String 返回上下文的字符串表示
func (c *Context) String() string {
	return fmt.Sprintf("%v.WithValue(%v)", c.Context, c.Values())
}

// IsContext 检查上下文是否是 Context
func IsContext(ctx context.Context) bool {
	_, ok := ctx.(*Context)
	return ok
}

// MergeContext 合并多个上下文为一个 Context
func MergeContext(ctxs ...context.Context) *Context {
	if len(ctxs) == 0 {
		return NewContext(context.Background(), nil) // 如果没有上下文，返回背景上下文
	}

	merged := NewContext(ctxs[0], nil) // 使用第一个上下文作为基础

	for _, ctx := range ctxs {
		if customCtx, ok := ctx.(*Context); ok { // 确保 ctx 是 Context 类型
			for key, value := range customCtx.values {
				merged.WithValue(key, value) // 设置值时，如果键已经存在，则会覆盖
			}
		}
	}

	return merged
}

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
