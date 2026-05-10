/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-10 10:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-10 10:00:00
 * @FilePath: \go-toolbox\pkg\types\pointer_test.go
 * @Description: 指针工具函数测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		p := Ptr(42)
		assert.NotNil(t, p)
		assert.Equal(t, 42, *p)
	})

	t.Run("string", func(t *testing.T) {
		p := Ptr("hello")
		assert.NotNil(t, p)
		assert.Equal(t, "hello", *p)
	})

	t.Run("bool", func(t *testing.T) {
		p := Ptr(true)
		assert.NotNil(t, p)
		assert.Equal(t, true, *p)
	})

	t.Run("float64", func(t *testing.T) {
		p := Ptr(3.14)
		assert.NotNil(t, p)
		assert.Equal(t, 3.14, *p)
	})

	t.Run("zero value int", func(t *testing.T) {
		p := Ptr(0)
		assert.NotNil(t, p)
		assert.Equal(t, 0, *p)
	})

	t.Run("empty string", func(t *testing.T) {
		p := Ptr("")
		assert.NotNil(t, p)
		assert.Equal(t, "", *p)
	})

	t.Run("time.Time", func(t *testing.T) {
		now := time.Now()
		p := Ptr(now)
		assert.NotNil(t, p)
		assert.Equal(t, now, *p)
	})
}

func TestDeref(t *testing.T) {
	t.Run("nil *int", func(t *testing.T) {
		var p *int
		assert.Equal(t, 0, Deref(p))
	})

	t.Run("valid *int", func(t *testing.T) {
		p := Ptr(42)
		assert.Equal(t, 42, Deref(p))
	})

	t.Run("nil *string", func(t *testing.T) {
		var p *string
		assert.Equal(t, "", Deref(p))
	})

	t.Run("valid *string", func(t *testing.T) {
		p := Ptr("hello")
		assert.Equal(t, "hello", Deref(p))
	})

	t.Run("nil *bool", func(t *testing.T) {
		var p *bool
		assert.Equal(t, false, Deref(p))
	})

	t.Run("valid *bool (true)", func(t *testing.T) {
		p := Ptr(true)
		assert.Equal(t, true, Deref(p))
	})

	t.Run("valid *bool (false)", func(t *testing.T) {
		p := Ptr(false)
		assert.Equal(t, false, Deref(p))
	})

	t.Run("nil *float64", func(t *testing.T) {
		var p *float64
		assert.Equal(t, 0.0, Deref(p))
	})

	t.Run("valid *float64", func(t *testing.T) {
		p := Ptr(3.14)
		assert.Equal(t, 3.14, Deref(p))
	})

	t.Run("nil *time.Time", func(t *testing.T) {
		var p *time.Time
		assert.Equal(t, time.Time{}, Deref(p))
	})

	t.Run("valid *time.Time", func(t *testing.T) {
		now := time.Now()
		p := Ptr(now)
		assert.Equal(t, now, Deref(p))
	})
}

func TestDerefOrDefault(t *testing.T) {
	t.Run("nil *int returns default", func(t *testing.T) {
		var p *int
		assert.Equal(t, 10, DerefOrDefault(p, 10))
	})

	t.Run("valid *int returns value", func(t *testing.T) {
		p := Ptr(42)
		assert.Equal(t, 42, DerefOrDefault(p, 10))
	})

	t.Run("nil *string returns default", func(t *testing.T) {
		var p *string
		assert.Equal(t, "default", DerefOrDefault(p, "default"))
	})

	t.Run("valid *string returns value", func(t *testing.T) {
		p := Ptr("hello")
		assert.Equal(t, "hello", DerefOrDefault(p, "default"))
	})

	t.Run("nil *bool returns default", func(t *testing.T) {
		var p *bool
		assert.Equal(t, true, DerefOrDefault(p, true))
	})

	t.Run("valid *bool returns value", func(t *testing.T) {
		p := Ptr(false)
		assert.Equal(t, false, DerefOrDefault(p, true))
	})

	t.Run("nil *float64 returns default", func(t *testing.T) {
		var p *float64
		assert.Equal(t, 99.9, DerefOrDefault(p, 99.9))
	})

	t.Run("valid *float64 returns value", func(t *testing.T) {
		p := Ptr(3.14)
		assert.Equal(t, 3.14, DerefOrDefault(p, 99.9))
	})

	t.Run("nil *time.Time returns default", func(t *testing.T) {
		var p *time.Time
		defaultTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, defaultTime, DerefOrDefault(p, defaultTime))
	})

	t.Run("valid *time.Time returns value", func(t *testing.T) {
		now := time.Now()
		p := Ptr(now)
		defaultTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, now, DerefOrDefault(p, defaultTime))
	})
}

func TestIsNilPtr(t *testing.T) {
	t.Run("nil *int", func(t *testing.T) {
		var p *int
		assert.True(t, IsNilPtr(p))
	})

	t.Run("valid *int", func(t *testing.T) {
		p := Ptr(42)
		assert.False(t, IsNilPtr(p))
	})

	t.Run("nil *string", func(t *testing.T) {
		var p *string
		assert.True(t, IsNilPtr(p))
	})

	t.Run("valid *string", func(t *testing.T) {
		p := Ptr("hello")
		assert.False(t, IsNilPtr(p))
	})

	t.Run("nil *bool", func(t *testing.T) {
		var p *bool
		assert.True(t, IsNilPtr(p))
	})

	t.Run("valid *bool", func(t *testing.T) {
		p := Ptr(true)
		assert.False(t, IsNilPtr(p))
	})
}

func TestIsNonNilPtr(t *testing.T) {
	t.Run("nil *int", func(t *testing.T) {
		var p *int
		assert.False(t, IsNonNilPtr(p))
	})

	t.Run("valid *int", func(t *testing.T) {
		p := Ptr(42)
		assert.True(t, IsNonNilPtr(p))
	})

	t.Run("nil *string", func(t *testing.T) {
		var p *string
		assert.False(t, IsNonNilPtr(p))
	})

	t.Run("valid *string", func(t *testing.T) {
		p := Ptr("hello")
		assert.True(t, IsNonNilPtr(p))
	})

	t.Run("nil *bool", func(t *testing.T) {
		var p *bool
		assert.False(t, IsNonNilPtr(p))
	})

	t.Run("valid *bool", func(t *testing.T) {
		p := Ptr(false)
		assert.True(t, IsNonNilPtr(p))
	})
}
