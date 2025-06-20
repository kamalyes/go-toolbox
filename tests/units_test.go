/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-17 19:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-18 13:20:31
 * @FilePath: \go-toolbox\tests\units_test.go
 * @Description: 单位格式化与解析单元测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/units"
	"github.com/stretchr/testify/assert"
)

func TestFormatFunctions(t *testing.T) {
	tests := []struct {
		name   string
		fn     func() string
		expect string
	}{
		{"HumanSize decimal", func() string { return units.HumanSize(123456789) }, "123.5MB"},
		{"CustomSize decimal 2f", func() string { return units.CustomSize("%.2f%s", 1234, 1000, units.DecimalAbbrs) }, "1.23kB"},
		{"CustomSize decimal 1f", func() string { return units.CustomSize("%.1f%s", 1200000000, 1000, units.DecimalAbbrs) }, "1.2GB"},
		{"BytesSize binary", func() string { return units.BytesSize(123456789) }, "117.7MiB"},
		{"CustomSize binary 0f", func() string { return units.CustomSize("%.0f%s", 22528, 1024, units.BinaryAbbrs) }, "22KiB"},
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
		{"decimal 22kB", "22kB", 22000, false, units.ParseSizeDecimal},
		{"decimal 17MB", "17MB", 17000000, false, units.ParseSizeDecimal},
		{"decimal 100", "100", 100, false, units.ParseSizeDecimal},
		{"decimal 100b", "100b", 100, false, units.ParseSizeDecimal},

		// 二进制解析正确示例
		{"binary 22KiB", "22KiB", 22528, false, units.ParseSizeBinary},
		{"binary 17MiB", "17MiB", 17825792, false, units.ParseSizeBinary},
		{"binary 100", "100", 100, false, units.ParseSizeBinary},
		{"binary 100b", "100b", 100, false, units.ParseSizeBinary},

		// 错误格式测试
		{"decimal invalid unit", "22XB", 0, true, units.ParseSizeDecimal},
		{"binary invalid unit", "22KiX", 0, true, units.ParseSizeBinary},
		{"decimal invalid number", "abc", 0, true, units.ParseSizeDecimal},
		{"binary negative number", "-22KiB", 0, true, units.ParseSizeBinary},
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
			got, err := units.ParseSizeDecimal(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
