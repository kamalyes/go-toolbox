/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\yamlx\logger_level_test.go
 * @Description: logger.LogLevel 的 YAML 往返测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

// 这两个测试验证 logger.LogLevel 与 yaml.v3 的真实集成（v3 对 v2 风格
// UnmarshalYAML 的兼容调用）。go-toolbox 核心不携带 yaml 依赖，
// 故与 yaml 一同落在本模块，而非 pkg/logger
package yamlx_test

import (
	"encoding/json"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// TestLogLevelUnmarshalYAML 测试 YAML 反序列化
func TestLogLevelUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected logger.LogLevel
		hasError bool
	}{
		{"string_debug", "level: debug", logger.DEBUG, false},
		{"string_info", "level: info", logger.INFO, false},
		{"string_warn", "level: warn", logger.WARN, false},
		{"string_error", "level: error", logger.ERROR, false},
		{"string_uppercase", "level: DEBUG", logger.DEBUG, false},
		{"integer_value", "level: 1", logger.DEBUG, true},
		{"invalid_string", "level: invalid", logger.DEBUG, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config struct {
				Level logger.LogLevel `yaml:"level"`
			}
			err := yaml.Unmarshal([]byte(tt.yaml), &config)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, config.Level)
			}
		})
	}
}

// TestLogLevelRoundTrip 测试序列化和反序列化往返
func TestLogLevelRoundTrip(t *testing.T) {
	levels := []logger.LogLevel{logger.DEBUG, logger.INFO, logger.WARN, logger.ERROR, logger.FATAL}
	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			// JSON 往返
			jsonData, err := json.Marshal(struct {
				Level logger.LogLevel `json:"level"`
			}{Level: level})
			assert.NoError(t, err)
			var jsonResult struct {
				Level logger.LogLevel `json:"level"`
			}
			err = json.Unmarshal(jsonData, &jsonResult)
			assert.NoError(t, err)
			assert.Equal(t, level, jsonResult.Level)

			// YAML 往返
			yamlData, err := yaml.Marshal(struct {
				Level logger.LogLevel `yaml:"level"`
			}{Level: level})
			assert.NoError(t, err)
			var yamlResult struct {
				Level logger.LogLevel `yaml:"level"`
			}
			err = yaml.Unmarshal(yamlData, &yamlResult)
			assert.NoError(t, err)
			assert.Equal(t, level, yamlResult.Level)
		})
	}
}
