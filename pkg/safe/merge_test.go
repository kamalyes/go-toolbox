/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-05 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-05 00:00:00
 * @FilePath: \go-toolbox\pkg\safe\merge_test.go
 * @Description: 泛型配置合并测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用的结构体
type TestConfig struct {
	Name        string
	Port        int
	Enabled     bool
	Tags        []string
	Metadata    map[string]string
	Nested      *NestedConfig
	SliceNested []*NestedConfig
}

type NestedConfig struct {
	Host    string
	Timeout int
	Debug   bool
	Options *Options
}

type Options struct {
	MaxRetry int
	UseCache bool
}

// TestMergeWithDefaultsNilConfig 测试 nil config 使用默认配置
func TestMergeWithDefaultsNilConfig(t *testing.T) {
	var nilConfig *TestConfig
	defaultCfg := &TestConfig{
		Name:    "Default",
		Port:    8080,
		Enabled: true,
	}

	result := MergeWithDefaults(nilConfig, defaultCfg)

	assert.NotNil(t, result, "Expected non-nil result")
	assert.Equal(t, defaultCfg.Name, result.Name)
	assert.Equal(t, defaultCfg.Port, result.Port)
	assert.Equal(t, defaultCfg.Enabled, result.Enabled)
}

// TestMergeWithDefaultsNoDefault 测试没有默认配置
func TestMergeWithDefaultsNoDefault(t *testing.T) {
	config := &TestConfig{
		Name: "Custom",
		Port: 9000,
	}

	result := MergeWithDefaults(config)

	assert.NotNil(t, result, "Expected non-nil result")
	assert.Equal(t, config.Name, result.Name)
	assert.Equal(t, config.Port, result.Port)
}

// TestMergeWithDefaultsPartialConfig 测试部分配置合并
func TestMergeWithDefaultsPartialConfig(t *testing.T) {
	partialCfg := &TestConfig{
		Name: "Partial",
	}
	defaultCfg := &TestConfig{
		Name:    "Default",
		Port:    8080,
		Enabled: true,
		Tags:    []string{"tag1", "tag2"},
	}

	result := MergeWithDefaults(partialCfg, defaultCfg)

	assert.Equal(t, partialCfg.Name, result.Name)
	assert.Equal(t, defaultCfg.Port, result.Port)
	assert.Equal(t, defaultCfg.Enabled, result.Enabled)
	assert.Len(t, result.Tags, 2)
}

// TestMergeWithDefaultsNestedStruct 测试嵌套结构体合并
func TestMergeWithDefaultsNestedStruct(t *testing.T) {
	partialCfg := &TestConfig{
		Name: "Test",
		Nested: &NestedConfig{
			Host: "127.0.0.1",
		},
	}
	defaultCfg := &TestConfig{
		Name: "Default",
		Port: 8080,
		Nested: &NestedConfig{
			Host:    "localhost",
			Timeout: 30,
			Debug:   true,
		},
	}

	result := MergeWithDefaults(partialCfg, defaultCfg)

	assert.NotNil(t, result.Nested, "Expected non-nil Nested")
	assert.Equal(t, partialCfg.Nested.Host, result.Nested.Host)
	assert.Equal(t, defaultCfg.Nested.Timeout, result.Nested.Timeout)
	assert.Equal(t, defaultCfg.Nested.Debug, result.Nested.Debug)
}

// TestMergeWithDefaultsNilNestedStruct 测试 nil 嵌套结构体
func TestMergeWithDefaultsNilNestedStruct(t *testing.T) {
	partialCfg := &TestConfig{
		Name: "Test",
	}
	defaultCfg := &TestConfig{
		Name: "Default",
		Nested: &NestedConfig{
			Host:    "localhost",
			Timeout: 30,
		},
	}

	result := MergeWithDefaults(partialCfg, defaultCfg)

	assert.NotNil(t, result.Nested, "Expected non-nil Nested from default")
	assert.Equal(t, defaultCfg.Nested.Host, result.Nested.Host, "Expected Host='localhost'")
	assert.Equal(t, defaultCfg.Nested.Timeout, result.Nested.Timeout, "Expected Timeout=30")
}

// TestMergeWithDefaultsDeepNested 测试深度嵌套结构体
func TestMergeWithDefaultsDeepNested(t *testing.T) {
	partialCfg := &TestConfig{
		Name: "Test",
		Nested: &NestedConfig{
			Host: "127.0.0.1",
			Options: &Options{
				MaxRetry: 5,
			},
		},
	}
	defaultCfg := &TestConfig{
		Name: "Default",
		Nested: &NestedConfig{
			Host:    "localhost",
			Timeout: 30,
			Options: &Options{
				MaxRetry: 3,
				UseCache: true,
			},
		},
	}

	result := MergeWithDefaults(partialCfg, defaultCfg)

	assert.NotNil(t, result.Nested, "Expected non-nil nested Options")
	assert.NotNil(t, result.Nested.Options, "Expected non-nil nested Options")
	assert.Equal(t, 5, result.Nested.Options.MaxRetry, "Expected MaxRetry=5")
	assert.Equal(t, defaultCfg.Nested.Options.UseCache, result.Nested.Options.UseCache)
	assert.Equal(t, defaultCfg.Nested.Timeout, result.Nested.Timeout)
}

// TestMergeWithDefaultsSlice 测试切片合并
func TestMergeWithDefaultsSlice(t *testing.T) {
	t.Run("nil slice uses default", func(t *testing.T) {
		partialCfg := &TestConfig{
			Name: "Test",
		}
		defaultCfg := &TestConfig{
			Tags: []string{"tag1", "tag2"},
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Len(t, result.Tags, 2, "Expected 2 tags")
	})

	t.Run("empty slice uses default", func(t *testing.T) {
		partialCfg := &TestConfig{
			Name: "Test",
			Tags: []string{},
		}
		defaultCfg := &TestConfig{
			Tags: []string{"tag1", "tag2"},
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Len(t, result.Tags, 2, "Expected 2 tags")
	})

	t.Run("non-empty slice keeps original", func(t *testing.T) {
		partialCfg := &TestConfig{
			Name: "Test",
			Tags: []string{"custom"},
		}
		defaultCfg := &TestConfig{
			Tags: []string{"tag1", "tag2"},
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Len(t, result.Tags, 1, "Expected 1 tag")
		assert.Equal(t, "custom", result.Tags[0], "Expected 'custom'")
	})
}

// TestMergeWithDefaultsMap 测试 Map 合并
func TestMergeWithDefaultsMap(t *testing.T) {
	t.Run("nil map uses default", func(t *testing.T) {
		partialCfg := &TestConfig{
			Name: "Test",
		}
		defaultCfg := &TestConfig{
			Metadata: map[string]string{
				"env":     "dev",
				"version": "1.0",
			},
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Len(t, result.Metadata, 2, "Expected 2 metadata entries")
		assert.Equal(t, "dev", result.Metadata["env"], "Expected env='dev'")
	})

	t.Run("existing map merges with default", func(t *testing.T) {
		partialCfg := &TestConfig{
			Name: "Test",
			Metadata: map[string]string{
				"custom": "value",
			},
		}
		defaultCfg := &TestConfig{
			Metadata: map[string]string{
				"env":     "dev",
				"version": "1.0",
			},
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Equal(t, "value", result.Metadata["custom"], "Expected custom='value'")
		assert.Equal(t, "dev", result.Metadata["env"], "Expected env='dev'")
	})
}

// TestMergeWithDefaultsString 测试字符串合并
func TestMergeWithDefaultsString(t *testing.T) {
	t.Run("empty string uses default", func(t *testing.T) {
		partialCfg := &TestConfig{
			Name: "",
		}
		defaultCfg := &TestConfig{
			Name: "Default",
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Equal(t, "Default", result.Name, "Expected Name='Default'")
	})

	t.Run("non-empty string keeps original", func(t *testing.T) {
		partialCfg := &TestConfig{
			Name: "Custom",
		}
		defaultCfg := &TestConfig{
			Name: "Default",
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Equal(t, "Custom", result.Name, "Expected Name='Custom'")
	})
}

// TestMergeWithDefaultsInt 测试整数合并
func TestMergeWithDefaultsInt(t *testing.T) {
	t.Run("zero int uses default", func(t *testing.T) {
		partialCfg := &TestConfig{
			Port: 0,
		}
		defaultCfg := &TestConfig{
			Port: 8080,
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Equal(t, 8080, result.Port, "Expected Port=8080")
	})

	t.Run("non-zero int keeps original", func(t *testing.T) {
		partialCfg := &TestConfig{
			Port: 9000,
		}
		defaultCfg := &TestConfig{
			Port: 8080,
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.Equal(t, 9000, result.Port, "Expected Port=9000")
	})
}

// TestMergeWithDefaultsBool 测试布尔值合并
func TestMergeWithDefaultsBool(t *testing.T) {
	t.Run("false bool uses default true", func(t *testing.T) {
		partialCfg := &TestConfig{
			Enabled: false,
		}
		defaultCfg := &TestConfig{
			Enabled: true,
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.True(t, result.Enabled, "Expected Enabled=true from default")
	})

	t.Run("true bool keeps original", func(t *testing.T) {
		partialCfg := &TestConfig{
			Enabled: true,
		}
		defaultCfg := &TestConfig{
			Enabled: false,
		}

		result := MergeWithDefaults(partialCfg, defaultCfg)

		assert.True(t, result.Enabled, "Expected Enabled=true")
	})
}

// TestMergeWithDefaultsMultipleDefaults 测试多个默认配置
func TestMergeWithDefaultsMultipleDefaults(t *testing.T) {
	partialCfg := &TestConfig{
		Name: "Test",
	}
	default1 := &TestConfig{
		Port: 8080,
	}
	default2 := &TestConfig{
		Port:    9090,
		Enabled: true,
		Tags:    []string{"tag1"},
	}

	result := MergeWithDefaults(partialCfg, default1, default2)

	assert.Equal(t, "Test", result.Name, "Expected Name='Test'")
	assert.Equal(t, 8080, result.Port, "Expected Port=8080")
	assert.True(t, result.Enabled, "Expected Enabled=true")
	assert.Len(t, result.Tags, 1, "Expected 1 tag")
}
