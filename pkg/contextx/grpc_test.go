/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-12 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-12 00:00:00
 * @FilePath: \go-toolbox\pkg\contextx\grpc_test.go
 * @Description: Context metadata 泛型操作测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package contextx

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockAdapter 模拟 metadata 适配器（用于测试）
type mockAdapter struct {
	data map[string]string
}

// newMockAdapter 创建模拟适配器
func newMockAdapter() *mockAdapter {
	return &mockAdapter{
		data: make(map[string]string),
	}
}

// Set 模拟写入 metadata
func (m *mockAdapter) Set(ctx context.Context, key, value string) context.Context {
	m.data[key] = value
	return ctx
}

// Get 模拟获取 metadata
func (m *mockAdapter) Get(ctx context.Context, key string) (string, bool) {
	val, ok := m.data[key]
	return val, ok
}

// Append 模拟追加 metadata
func (m *mockAdapter) Append(ctx context.Context, key, value string) context.Context {
	if existing, ok := m.data[key]; ok {
		m.data[key] = existing + "," + value
	} else {
		m.data[key] = value
	}
	return ctx
}

// mockMarshaler 模拟 JSON 序列化器（用于测试）
type mockMarshaler struct{}

// Marshal 模拟序列化
func (m *mockMarshaler) Marshal(v any) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Unmarshal 模拟反序列化
func (m *mockMarshaler) Unmarshal(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}

// TestUser 测试用户结构
type TestUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TestMetadataManager_Set 测试 Set 方法
func TestMetadataManager_Set(t *testing.T) {
	adapter := newMockAdapter()
	marshaler := &mockMarshaler{}
	manager := NewMetadataManager(adapter, marshaler)

	ctx := context.Background()
	user := TestUser{ID: "123", Name: "Alice", Age: 30}

	newCtx, err := manager.Set(ctx, "user", user)
	assert.NoError(t, err)
	assert.NotNil(t, newCtx)

	// 验证数据已写入
	val, ok := adapter.Get(ctx, "user")
	assert.True(t, ok, "Expected key 'user' to exist")

	var result TestUser
	err = marshaler.Unmarshal(val, &result)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Age, result.Age)
}

// TestMetadataManager_Get 测试 Get 方法
func TestMetadataManager_Get(t *testing.T) {
	adapter := newMockAdapter()
	marshaler := &mockMarshaler{}
	manager := NewMetadataManager(adapter, marshaler)

	ctx := context.Background()
	user := TestUser{ID: "456", Name: "Bob", Age: 25}

	// 先写入数据
	ctx, err := manager.Set(ctx, "user", user)
	assert.NoError(t, err)

	// 读取数据
	var result TestUser
	err = manager.Get(ctx, "user", &result)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Age, result.Age)
}

// TestMetadataManager_Get_KeyNotFound 测试获取不存在的 key
func TestMetadataManager_Get_KeyNotFound(t *testing.T) {
	adapter := newMockAdapter()
	marshaler := &mockMarshaler{}
	manager := NewMetadataManager(adapter, marshaler)

	ctx := context.Background()
	var result TestUser

	err := manager.Get(ctx, "nonexistent", &result)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

// TestMetadataManager_Append 测试 Append 方法
func TestMetadataManager_Append(t *testing.T) {
	adapter := newMockAdapter()
	marshaler := &mockMarshaler{}
	manager := NewMetadataManager(adapter, marshaler)

	ctx := context.Background()

	// 第一次写入
	ctx, err := manager.Set(ctx, "tags", "tag1")
	assert.NoError(t, err)

	// 追加数据
	ctx, err = manager.Append(ctx, "tags", "tag2")
	assert.NoError(t, err)

	// 验证追加结果
	val, ok := adapter.Get(ctx, "tags")
	assert.True(t, ok, "Expected key 'tags' to exist")

	// 模拟适配器会用逗号连接
	expected := `"tag1","tag2"`
	assert.Equal(t, expected, val)
}

// TestMetadataManager_ComplexStruct 测试复杂结构体
func TestMetadataManager_ComplexStruct(t *testing.T) {
	adapter := newMockAdapter()
	marshaler := &mockMarshaler{}
	manager := NewMetadataManager(adapter, marshaler)

	ctx := context.Background()

	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}

	type ComplexUser struct {
		ID      string            `json:"id"`
		Name    string            `json:"name"`
		Tags    []string          `json:"tags"`
		Meta    map[string]string `json:"meta"`
		Address Address           `json:"address"`
	}

	user := ComplexUser{
		ID:   "789",
		Name: "Charlie",
		Tags: []string{"admin", "developer"},
		Meta: map[string]string{"role": "admin", "level": "senior"},
		Address: Address{
			City:    "Beijing",
			Country: "China",
		},
	}

	// 写入复杂结构
	ctx, err := manager.Set(ctx, "complex_user", user)
	assert.NoError(t, err)

	// 读取复杂结构
	var result ComplexUser
	err = manager.Get(ctx, "complex_user", &result)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, len(user.Tags), len(result.Tags))
	assert.Equal(t, user.Address.City, result.Address.City)
}

// TestMetadataManager_MarshalError 测试序列化错误
func TestMetadataManager_MarshalError(t *testing.T) {
	adapter := newMockAdapter()
	marshaler := &mockMarshaler{}
	manager := NewMetadataManager(adapter, marshaler)

	ctx := context.Background()

	// channel 无法序列化
	ch := make(chan int)
	_, err := manager.Set(ctx, "channel", ch)
	assert.Error(t, err)
}

// TestMetadataManager_UnmarshalError 测试反序列化错误
func TestMetadataManager_UnmarshalError(t *testing.T) {
	adapter := newMockAdapter()
	marshaler := &mockMarshaler{}
	manager := NewMetadataManager(adapter, marshaler)

	ctx := context.Background()

	// 写入无效的 JSON
	adapter.data["invalid"] = "invalid json"

	var result TestUser
	err := manager.Get(ctx, "invalid", &result)
	assert.Error(t, err)
}
