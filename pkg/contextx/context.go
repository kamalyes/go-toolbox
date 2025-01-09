/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-09 15:55:51
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

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// Context 是一个自定义的上下文，支持多个值的存储
type Context struct {
	mu     sync.RWMutex
	values map[interface{}]interface{}
	context.Context
	cancelFunc context.CancelFunc // 添加取消函数
	pool       *osx.LimitedPool   // 引入字节切片池
}

// NewContext 创建一个新的 Context，允许用户传入自定义的字节切片池
func NewContext(parent context.Context, pool *osx.LimitedPool) *Context {
	if pool == nil {
		pool = osx.NewLimitedPool(32, 1024)
	}
	return &Context{
		values:  make(map[interface{}]interface{}),
		Context: parent,
		pool:    pool,
	}
}

// NewContextWithTimeout 创建一个带有超时的 Context
func NewContextWithTimeout(parent context.Context, timeout time.Duration, pool *osx.LimitedPool) *Context {
	ctx, cancel := context.WithTimeout(parent, timeout)
	return &Context{
		values:     make(map[interface{}]interface{}),
		Context:    ctx,
		cancelFunc: cancel,
		pool:       pool,
	}
}

// NewContextWithCancel 创建一个可以手动取消的 Context
func NewContextWithCancel(parent context.Context, pool *osx.LimitedPool) *Context {
	ctx, cancel := context.WithCancel(parent)
	return &Context{
		values:     make(map[interface{}]interface{}),
		Context:    ctx,
		cancelFunc: cancel,
		pool:       pool,
	}
}

// NewContextWithValue 在父上下文中设置值并返回新的 Context
func NewContextWithValue(parent context.Context, key, val interface{}, pool *osx.LimitedPool) (*Context, error) {
	customCtx := NewContext(parent, pool)
	if err := customCtx.Set(key, val); err != nil {
		return nil, err
	}
	return customCtx, nil
}

// NewLocalContextWithValue 在当前 Context 中设置局部值
func NewLocalContextWithValue(ctx *Context, key, val interface{}) (*Context, error) {
	if err := ctx.Set(key, val); err != nil {
		return nil, err
	}
	return ctx, nil
}

// Values 返回当前上下文中所有的键值对
func (c *Context) Values() map[interface{}]interface{} {
	return syncx.WithRLockReturnValue(&c.mu, func() map[interface{}]interface{} {
		// 创建一个新的映射以返回
		valuesCopy := make(map[interface{}]interface{})
		for k, v := range c.values {
			valuesCopy[k] = v
		}
		return valuesCopy
	})
}

// validateKey 检查键是否有效
func validateKey(key interface{}) error {
	if key == nil || !reflect.TypeOf(key).Comparable() {
		if key == nil {
			return fmt.Errorf("nil key")
		}
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

// setByteSlice 处理字节切片的存储
func (c *Context) setByteSlice(key interface{}, value []byte) error {
	if buf := c.pool.Get(len(value)); buf != nil {
		copy(*buf, value)
		c.values[key] = buf
		return nil
	}
	c.values[key] = value
	return nil
}

// Set 设置指定键的值并返回错误
func (c *Context) Set(key, value interface{}) error {
	if err := validateKey(key); err != nil {
		return err
	}

	if byteSlice, ok := value.([]byte); ok {
		return c.setByteSlice(key, byteSlice)
	}

	syncx.WithLock(&c.mu, func() {
		c.values[key] = value
	})
	return nil
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
	return fmt.Sprintf("%v.WithValue(%v)", c.Context, c.values)
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
			syncx.WithLock(&customCtx.mu, func() {
				for key, value := range customCtx.values {
					// 设置值时，如果键已经存在，则会覆盖
					if err := merged.Set(key, value); err != nil {
						// 处理错误（可选）
						fmt.Printf("Error setting value: %v\n", err)
					}
				}
			})
		}
	}

	return merged
}
