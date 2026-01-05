/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 10:00:00
 * @FilePath: \go-toolbox\pkg\contextx\values.go
 * @Description: Context 值操作相关方法
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package contextx

import (
	"fmt"
	"reflect"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

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

// Range 遍历所有键值对（类似 sync.Map.Range）
func (c *Context) Range(f func(key, value interface{}) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for k, v := range c.values {
		if !f(k, v) {
			break
		}
	}
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

// WithByteSlice 处理字节切片的存储，支持链式调用
func (c *Context) WithByteSlice(key interface{}, value []byte) *Context {
	syncx.WithLock(&c.mu, func() {
		if buf := c.pool.Get(len(value)); buf != nil {
			copy(*buf, value)
			c.values[key] = buf
		} else {
			c.values[key] = value
		}
	})
	return c
}

// WithValue 设置指定键的值，支持链式调用
func (c *Context) WithValue(key, value interface{}) *Context {
	if err := validateKey(key); err != nil {
		fmt.Printf("contextx.WithValue error: %v\n", err)
		return c
	}

	if byteSlice, ok := value.([]byte); ok {
		return c.WithByteSlice(key, byteSlice)
	}

	syncx.WithLock(&c.mu, func() {
		c.values[key] = value
	})
	return c
}

// Remove 删除指定键的键值对，支持链式调用
func (c *Context) Remove(key interface{}) *Context {
	syncx.WithLock(&c.mu, func() {
		delete(c.values, key)
	})
	return c
}

// SetBatch 批量设置键值对，支持链式调用
func (c *Context) SetBatch(kvs map[interface{}]interface{}) *Context {
	for k, v := range kvs {
		c.WithValue(k, v)
	}
	return c
}

// MustValue 获取值，不存在则 panic
func (c *Context) MustValue(key interface{}) interface{} {
	if val := c.Value(key); val != nil {
		return val
	}
	panic(fmt.Sprintf("key %v not found in context", key))
}
