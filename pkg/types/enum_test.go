/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 00:00:00
 * @FilePath: \go-toolbox\pkg\types\enum_test.go
 * @Description: 枚举验证器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义测试用的枚举类型
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
	RoleGuest UserRole = "guest"
)

type Status int

const (
	StatusPending  Status = 1
	StatusApproved Status = 2
	StatusRejected Status = 3
)

func TestNewEnumValidator(t *testing.T) {
	t.Run("string enum", func(t *testing.T) {
		validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
		assert.NotNil(t, validator)
		assert.Equal(t, 3, validator.Count())
	})

	t.Run("int enum", func(t *testing.T) {
		validator := NewEnumValidator(StatusPending, StatusApproved, StatusRejected)
		assert.NotNil(t, validator)
		assert.Equal(t, 3, validator.Count())
	})

	t.Run("empty validator", func(t *testing.T) {
		validator := NewEnumValidator[string]()
		assert.NotNil(t, validator)
		assert.Equal(t, 0, validator.Count())
	})
}

func TestEnumValidator_IsValid(t *testing.T) {
	validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)

	tests := []struct {
		name     string
		value    UserRole
		expected bool
	}{
		{"valid admin", RoleAdmin, true},
		{"valid user", RoleUser, true},
		{"valid guest", RoleGuest, true},
		{"invalid role", UserRole("superadmin"), false},
		{"empty string", UserRole(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.IsValid(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnumValidator_MustBeValid(t *testing.T) {
	validator := NewEnumValidator(StatusPending, StatusApproved, StatusRejected)

	t.Run("valid value", func(t *testing.T) {
		err := validator.MustBeValid(StatusPending)
		assert.NoError(t, err)
	})

	t.Run("invalid value", func(t *testing.T) {
		err := validator.MustBeValid(Status(999))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid enum value")
	})
}

func TestEnumValidator_GetValidValues(t *testing.T) {
	validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
	values := validator.GetValidValues()

	assert.Equal(t, 3, len(values))
	// 验证所有值都存在
	valueMap := make(map[UserRole]bool)
	for _, v := range values {
		valueMap[v] = true
	}
	assert.True(t, valueMap[RoleAdmin])
	assert.True(t, valueMap[RoleUser])
	assert.True(t, valueMap[RoleGuest])
}

func TestEnumValidator_GetValidValuesString(t *testing.T) {
	validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
	values := validator.GetValidValuesString()

	assert.Equal(t, 3, len(values))
	// 验证是否排序
	assert.Equal(t, []string{"admin", "guest", "user"}, values)
}

func TestEnumValidator_Contains(t *testing.T) {
	validator := NewEnumValidator(RoleAdmin, RoleUser)

	assert.True(t, validator.Contains(RoleAdmin))
	assert.True(t, validator.Contains(RoleUser))
	assert.False(t, validator.Contains(RoleGuest))
}

func TestEnumValidator_Add(t *testing.T) {
	validator := NewEnumValidator(RoleAdmin, RoleUser)
	assert.Equal(t, 2, validator.Count())

	// 添加新值
	validator.Add(RoleGuest)
	assert.Equal(t, 3, validator.Count())
	assert.True(t, validator.IsValid(RoleGuest))

	// 重复添加不会增加数量
	validator.Add(RoleGuest)
	assert.Equal(t, 3, validator.Count())

	// 批量添加
	validator.Add(UserRole("moderator"), UserRole("vip"))
	assert.Equal(t, 5, validator.Count())
}

func TestEnumValidator_Remove(t *testing.T) {
	validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
	assert.Equal(t, 3, validator.Count())

	// 移除一个值
	validator.Remove(RoleGuest)
	assert.Equal(t, 2, validator.Count())
	assert.False(t, validator.IsValid(RoleGuest))

	// 移除不存在的值不会报错
	validator.Remove(UserRole("nonexistent"))
	assert.Equal(t, 2, validator.Count())

	// 批量移除
	validator.Remove(RoleAdmin, RoleUser)
	assert.Equal(t, 0, validator.Count())
}

func TestEnumValidator_Clear(t *testing.T) {
	validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
	assert.Equal(t, 3, validator.Count())

	validator.Clear()
	assert.Equal(t, 0, validator.Count())
	assert.False(t, validator.IsValid(RoleAdmin))
}

func TestEnumValidator_Clone(t *testing.T) {
	original := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
	clone := original.Clone()

	// 验证克隆的内容相同
	assert.Equal(t, original.Count(), clone.Count())
	assert.True(t, clone.IsValid(RoleAdmin))
	assert.True(t, clone.IsValid(RoleUser))
	assert.True(t, clone.IsValid(RoleGuest))

	// 修改克隆不影响原对象
	clone.Add(UserRole("moderator"))
	assert.Equal(t, 3, original.Count())
	assert.Equal(t, 4, clone.Count())

	// 修改原对象不影响克隆
	original.Remove(RoleAdmin)
	assert.False(t, original.IsValid(RoleAdmin))
	assert.True(t, clone.IsValid(RoleAdmin))
}

// 基准测试
func BenchmarkEnumValidator_IsValid(b *testing.B) {
	validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.IsValid(RoleUser)
	}
}

func BenchmarkEnumValidator_MustBeValid(b *testing.B) {
	validator := NewEnumValidator(StatusPending, StatusApproved, StatusRejected)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.MustBeValid(StatusApproved)
	}
}

func BenchmarkEnumValidator_GetValidValues(b *testing.B) {
	validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.GetValidValues()
	}
}

// 示例测试
func ExampleEnumValidator() {
	// 创建用户角色验证器
	roleValidator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)

	// 验证角色
	if roleValidator.IsValid(RoleAdmin) {
		// 角色有效
	}

	// 强制验证
	if err := roleValidator.MustBeValid(RoleUser); err != nil {
		// 处理错误
	}

	// 获取所有有效值
	validRoles := roleValidator.GetValidValues()
	_ = validRoles

	// Output:
}

func ExampleEnumValidator_Add() {
	validator := NewEnumValidator(RoleAdmin, RoleUser)

	// 添加新的角色
	validator.Add(RoleGuest)
	validator.Add(UserRole("moderator"), UserRole("vip"))

	// Output:
}
