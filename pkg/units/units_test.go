/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-17 19:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 23:19:55
 * @FilePath: \go-toolbox\pkg\units\units_test.go
 * @Description: 单位格式化与解析单元测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFunctions(t *testing.T) {
	tests := []struct {
		name   string
		fn     func() string
		expect string
	}{
		{"HumanSize decimal", func() string { return HumanSize(123456789) }, "123.5MB"},
		{"CustomSize decimal 2f", func() string { return CustomSize("%.2f%s", 1234, 1000, DecimalAbbrs) }, "1.23kB"},
		{"CustomSize decimal 1f", func() string { return CustomSize("%.1f%s", 1200000000, 1000, DecimalAbbrs) }, "1.2GB"},
		{"BytesSize binary", func() string { return BytesSize(123456789) }, "117.7MiB"},
		{"CustomSize binary 0f", func() string { return CustomSize("%.0f%s", 22528, 1024, BinaryAbbrs) }, "22KiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.fn())
		})
	}
}

func TestParseFunctions(t *testing.T) {
	type parseTest struct {
		name      string
		input     string
		want      int64
		wantError bool
		fn        func(string) (int64, error)
	}

	tests := []parseTest{
		// 十进制解析正确示例
		{"decimal 22kB", "22kB", 22000, false, ParseSizeDecimal},
		{"decimal 17MB", "17MB", 17000000, false, ParseSizeDecimal},
		{"decimal 100", "100", 100, false, ParseSizeDecimal},
		{"decimal 100b", "100b", 100, false, ParseSizeDecimal},

		// 二进制解析正确示例
		{"binary 22KiB", "22KiB", 22528, false, ParseSizeBinary},
		{"binary 17MiB", "17MiB", 17825792, false, ParseSizeBinary},
		{"binary 100", "100", 100, false, ParseSizeBinary},
		{"binary 100b", "100b", 100, false, ParseSizeBinary},

		// 错误格式测试
		{"decimal invalid unit", "22XB", 0, true, ParseSizeDecimal},
		{"binary invalid unit", "22KiX", 0, true, ParseSizeBinary},
		{"decimal invalid number", "abc", 0, true, ParseSizeDecimal},
		{"binary negative number", "-22KiB", 0, true, ParseSizeBinary},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fn(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		want      int64
	}{
		{"empty string", "", true, 0},
		{"too long suffix", "123kibb", true, 0},
		{"only number", "12345", false, 12345},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSizeDecimal(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseBytes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      uint64
		wantError bool
	}{
		// 二进制单位（优先匹配）
		{"1GB binary", "1GB", 1073741824, false},    // 1 GiB = 1024^3
		{"512MB binary", "512MB", 536870912, false}, // 512 MiB
		{"2048KB binary", "2048KB", 2097152, false}, // 2048 KiB
		{"1GiB explicit", "1GiB", 1073741824, false},
		{"100MiB explicit", "100MiB", 104857600, false},

		// 无单位（纯数字）
		{"plain number", "1024", 1024, false},
		{"plain with b", "2048b", 2048, false},

		// 大小写不敏感
		{"lowercase gb", "1gb", 1073741824, false},
		{"uppercase GB", "1GB", 1073741824, false},
		{"mixed case Mb", "512Mb", 536870912, false},

		// 边界情况
		{"zero", "0", 0, false},
		{"zero with unit", "0MB", 0, false},

		// 错误情况
		{"invalid unit", "100XB", 0, true},
		{"invalid format", "abc", 0, true},
		{"empty string", "", 0, true},
		{"negative", "-100MB", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBytes(tt.input)
			if tt.wantError {
				assert.Error(t, err, "Expected error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Unexpected error for input: %s", tt.input)
				assert.Equal(t, tt.want, got, "ParseBytes(%s) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
