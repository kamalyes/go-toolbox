/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 10:00:00
 * @FilePath: \go-toolbox\pkg\contextx\metadata.go
 * @Description: Context 元数据操作方法
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package contextx

// WithMetadata 设置元数据（并发安全）
func (c *Context) WithMetadata(key, value string) *Context {
	c.metadata.Store(key, value)
	return c
}

// GetMetadata 获取元数据
func (c *Context) GetMetadata(key string) string {
	if val, ok := c.metadata.Load(key); ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// SetMetadataBatch 批量设置元数据
func (c *Context) SetMetadataBatch(kvs map[string]string) *Context {
	for k, v := range kvs {
		c.metadata.Store(k, v)
	}
	return c
}

// GetAllMetadata 获取所有元数据
func (c *Context) GetAllMetadata() map[string]string {
	result := make(map[string]string)
	c.metadata.Range(func(k, v interface{}) bool {
		if key, ok := k.(string); ok {
			if val, ok := v.(string); ok {
				result[key] = val
			}
		}
		return true
	})
	return result
}
