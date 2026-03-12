/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-03-12 16:30:21
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-03-12 16:35:55
 * @FilePath: \go-toolbox\pkg\convert\format_stats_test.go
 * @Description: 统计数据格式化工具测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		// 基础场景
		{"秒数", 45, "45s"},
		{"分钟+秒", 90, "1m 30s"},
		{"小时+分钟+秒", 3665, "1h 1m 5s"},
		{"整小时", 3600, "1h"},
		{"整天", 86400, "1d"},

		// 复杂场景
		{"天+小时+分钟", 90061, "1d 1h 1m"},
		{"月+天+分钟", 2678461, "1mo 1d 1m"}, // 修正：2678461秒 = 1个月(2592000) + 1天(86400) + 61秒(1m 1s)
		{"年+月+天", 34214400, "1y 1mo 1d"},
		{"多年", 63072000, "2y"},

		// 边界情况
		{"零秒", 0, "N/A"},
		{"负数", -10, "N/A"},
		{"nil值", nil, "N/A"},

		// 类型支持
		{"float64", 65.8, "1m 5s"},
		{"int64", int64(200), "3m 20s"},
		{"int32", int32(150), "2m 30s"},
		{"float32", float32(120.5), "2m"},

		// 特殊情况
		{"不支持的类型", "invalid", "N/A"},
		{"大数值", 31536000, "1y"},
		{"超大数值", 94608000, "3y"},
	}

	for _, tt := range tests {
		got := FormatDuration(tt.value)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestFormatCount(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"正常数字", 123, "123"},
		{"零值", 0, "0"},
		{"nil值", nil, "0"},
		{"float64", 45.67, "45.67"},
		{"字符串", "test", "test"},
		{"负数", -10, "-10"},
		{"大数", 1000000, "1000000"},
	}

	for _, tt := range tests {
		got := FormatCount(tt.value)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestFormatPercentage(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		precision int
		want      string
	}{
		{"正常百分比", 85.567, 1, "85.6%"},
		{"整数百分比", 100, 0, "100%"},
		{"nil值", nil, 1, "0%"},
		{"零值", 0.0, 2, "0.00%"},
		{"小数", 12.3456, 2, "12.35%"},
		{"int类型", 75, 1, "75.0%"},
		{"不支持的类型", "invalid", 1, "0%"},
		{"高精度", 99.9999, 3, "100.000%"},
	}

	for _, tt := range tests {
		got := FormatPercentage(tt.value, tt.precision)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

// TestFormatDurationEdgeCases 测试边界情况
func TestFormatDurationEdgeCases(t *testing.T) {
	// 1秒
	assert.Equal(t, "1s", FormatDuration(1))

	// 59秒
	assert.Equal(t, "59s", FormatDuration(59))

	// 1分钟
	assert.Equal(t, "1m", FormatDuration(60))

	// 1小时
	assert.Equal(t, "1h", FormatDuration(3600))

	// 23小时59分59秒
	assert.Equal(t, "23h 59m 59s", FormatDuration(86399))

	// 1天
	assert.Equal(t, "1d", FormatDuration(86400))

	// 29天23小时59分
	assert.Equal(t, "29d 23h 59m", FormatDuration(2591999))

	// 1个月
	assert.Equal(t, "1mo", FormatDuration(2592000))

	// 12个月4天 (实际计算：31449600秒 = 12个月(31104000) + 4天(345600))
	assert.Equal(t, "12mo 4d", FormatDuration(31449600))

	// 1年
	assert.Equal(t, "1y", FormatDuration(31536000))
}
