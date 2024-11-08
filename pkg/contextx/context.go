/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 20:22:05
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

	"github.com/kamalyes/go-toolbox/pkg/osx"
)

// CustomContext 是一个自定义的上下文，支持多个值的存储
type CustomContext struct {
	mu   sync.RWMutex
	Tags map[interface{}]interface{}
	context.Context
	pool *osx.LimitedPool // 引入字节切片池
}

// NewCustomContext 创建一个新的 CustomContext，允许用户传入自定义的字节切片池
func NewCustomContext(parent context.Context, pool *osx.LimitedPool) *CustomContext {
	if pool == nil {
		pool = osx.NewLimitedPool(32, 1024)
	}
	return &CustomContext{
		Tags:    make(map[interface{}]interface{}),
		Context: parent,
		pool:    pool,
	}
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
func (c *CustomContext) Value(key interface{}) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if v, ok := c.Tags[key]; ok {
		return v
	}
	return c.Context.Value(key)
}

// setByteSlice 处理字节切片的存储
func (c *CustomContext) setByteSlice(key interface{}, value []byte) error {
	if buf := c.pool.Get(len(value)); buf != nil {
		copy(*buf, value)
		c.Tags[key] = buf
		return nil
	}
	c.Tags[key] = value
	return nil
}

// Set 设置指定键的值并返回错误
func (c *CustomContext) Set(key, value interface{}) error {
	if err := validateKey(key); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if byteSlice, ok := value.([]byte); ok {
		return c.setByteSlice(key, byteSlice)
	}

	c.Tags[key] = value
	return nil
}

// Remove 删除指定键的键值对
func (c *CustomContext) Remove(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Tags, key)
}

// String 返回上下文的字符串表示
func (c *CustomContext) String() string {
	return fmt.Sprintf("%v.WithValue(%v)", c.Context, c.Tags)
}

// NewContextWithValue 在父上下文中设置值并返回新的 CustomContext
func NewContextWithValue(parent context.Context, key, val interface{}, pool *osx.LimitedPool) (*CustomContext, error) {
	customCtx := NewCustomContext(parent, pool)
	if err := customCtx.Set(key, val); err != nil {
		return nil, err
	}
	return customCtx, nil
}

// NewLocalContextWithValue 在当前 CustomContext 中设置局部值
func NewLocalContextWithValue(ctx *CustomContext, key, val interface{}) (*CustomContext, error) {
	if err := ctx.Set(key, val); err != nil {
		return nil, err
	}
	return ctx, nil
}

// IsCustomContext 检查上下文是否是 CustomContext
func IsCustomContext(ctx context.Context) bool {
	_, ok := ctx.(*CustomContext)
	return ok
}
