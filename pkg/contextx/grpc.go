/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-12 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-12 00:00:00
 * @FilePath: \go-toolbox\pkg\contextx\grpc.go
 * @Description: Context metadata 泛型操作接口（轻量级，无外部依赖）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package contextx

import (
	"context"
	"errors"
)

var (
	// ErrKeyNotFound 键不存在错误
	ErrKeyNotFound = errors.New("key not found in context")
	// ErrMarshalFailed 序列化失败错误
	ErrMarshalFailed = errors.New("failed to marshal value")
	// ErrUnmarshalFailed 反序列化失败错误
	ErrUnmarshalFailed = errors.New("failed to unmarshal value")
)

// MetadataAdapter 定义 metadata 操作的适配器接口
// 使用适配器模式，让不同的实现（gRPC、HTTP、自定义）可以灵活接入
type MetadataAdapter interface {
	// Set 将字符串值写入 context metadata
	Set(ctx context.Context, key, value string) context.Context

	// Get 从 context metadata 获取字符串值
	Get(ctx context.Context, key string) (string, bool)

	// Append 向现有 context 追加 metadata（不覆盖已有的）
	Append(ctx context.Context, key, value string) context.Context
}

// Marshaler 定义序列化接口
type Marshaler interface {
	Marshal(v any) (string, error)
	Unmarshal(data string, v any) error
}

// MetadataManager 封装 adapter 和 marshaler，提供简洁的 API
type MetadataManager struct {
	adapter   MetadataAdapter
	marshaler Marshaler
}

// NewMetadataManager 创建 metadata 管理器
func NewMetadataManager(adapter MetadataAdapter, marshaler Marshaler) *MetadataManager {
	return &MetadataManager{
		adapter:   adapter,
		marshaler: marshaler,
	}
}

// Set 将任意类型的数据写入 context metadata（泛型方法）
func (m *MetadataManager) Set(ctx context.Context, key string, value any) (context.Context, error) {
	data, err := m.marshaler.Marshal(value)
	if err != nil {
		return ctx, errors.Join(ErrMarshalFailed, err)
	}

	return m.adapter.Set(ctx, key, data), nil
}

// Get 从 context metadata 获取数据（泛型方法）
func (m *MetadataManager) Get(ctx context.Context, key string, result any) error {
	val, ok := m.adapter.Get(ctx, key)
	if !ok {
		return ErrKeyNotFound
	}

	if err := m.marshaler.Unmarshal(val, result); err != nil {
		return errors.Join(ErrUnmarshalFailed, err)
	}
	return nil
}

// GetOrDefault 获取数据，如果不存在则返回默认值（泛型方法）
func (m *MetadataManager) GetOrDefault(ctx context.Context, key string, defaultValue any) any {
	val, ok := m.adapter.Get(ctx, key)
	if !ok {
		return defaultValue
	}

	if err := m.marshaler.Unmarshal(val, defaultValue); err != nil {
		return defaultValue
	}
	return defaultValue
}

// Append 向现有 context 追加 metadata
func (m *MetadataManager) Append(ctx context.Context, key string, value any) (context.Context, error) {
	data, err := m.marshaler.Marshal(value)
	if err != nil {
		return ctx, errors.Join(ErrMarshalFailed, err)
	}

	return m.adapter.Append(ctx, key, data), nil
}
