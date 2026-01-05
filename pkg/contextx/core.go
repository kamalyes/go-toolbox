/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 09:55:06
 * @FilePath: \go-toolbox\pkg\contextx\core.go
 * @Description: Context 核心定义和构造函数
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package contextx

import (
	"context"
	"sync"
	"sync/atomic"
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
	deadline   atomic.Int64       // 超时时间（UnixNano）
	metadata   sync.Map           // 元数据存储（并发安全）
}

// NewContext 创建一个新的 Context（使用默认配置）
func NewContext() *Context {
	return &Context{
		values:  make(map[interface{}]interface{}),
		Context: context.Background(),
		pool:    syncx.NewLimitedPool(32, 1024),
	}
}

// WithParent 设置父上下文
func (c *Context) WithParent(parent context.Context) *Context {
	if parent != nil {
		c.Context = parent
	}
	return c
}

// WithPool 设置对象池
func (c *Context) WithPool(pool *syncx.LimitedPool) *Context {
	if pool != nil {
		c.pool = pool
	}
	return c
}

// WithCancel 添加取消功能
func (c *Context) WithCancel() *Context {
	ctx, cancel := context.WithCancel(c.Context)
	c.Context = ctx
	c.cancelFunc = cancel
	return c
}

// WithTimeout 设置超时时间
func (c *Context) WithTimeout(timeout time.Duration) *Context {
	ctx, cancel := context.WithTimeout(c.Context, timeout)
	c.Context = ctx
	c.cancelFunc = cancel
	c.deadline.Store(time.Now().Add(timeout).UnixNano())
	return c
}

// WithDeadline 设置绝对截止时间
func (c *Context) WithDeadline(deadline time.Time) *Context {
	ctx, cancel := context.WithDeadline(c.Context, deadline)
	c.Context = ctx
	c.cancelFunc = cancel
	c.deadline.Store(deadline.UnixNano())
	return c
}

// NewContextWithTimeout 创建一个带有超时的 Context（便捷方法）
func NewContextWithTimeout(timeout time.Duration) *Context {
	return NewContext().WithTimeout(timeout)
}

// NewContextWithValue 在上下文中设置值并返回新的 Context（便捷方法）
func NewContextWithValue(key, val interface{}) *Context {
	return NewContext().WithValue(key, val)
}

// IsContext 检查上下文是否是 Context
func IsContext(ctx context.Context) bool {
	_, ok := ctx.(*Context)
	return ok
}
