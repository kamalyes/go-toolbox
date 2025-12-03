/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-03 15:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-03 15:01:01
 * @FilePath: \go-toolbox\pkg\moment\base_test.go
 * @Description: moment 包的测试文件
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package moment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试用常量，避免字符串重复
const (
	testDateStr           = "2025-12-03"
	testTimeStr           = "14:30:45"
	testFullDateTime      = testDateStr + " " + testTimeStr
	testRFC3339Str        = testDateStr + "T" + testTimeStr
	testRFC3339UTCStr     = testDateStr + "T" + testTimeStr + "Z"
	testRFC3339MillisStr  = testDateStr + "T" + testTimeStr + ".000Z"
	testRFC3339PlusTZStr  = testDateStr + "T22:30:45+08:00"
	testRFC3339MinusTZStr = testDateStr + "T06:30:45-08:00"
	testDateTimeStr       = testDateStr + " " + testTimeStr
	testDateTimeMillisStr = testDateStr + " " + testTimeStr + ".000"
	testWhitespaceStr     = "  " + testDateStr + "  "
	testChineseDateStr    = "2025年12月3日"
	testSlashDateStr      = "2025/12/03"
	testDotDateStr        = "2025.12.03"
	testCompactDateStr    = "20251203"
)

func TestParseFlexibleDate(t *testing.T) {
	// 期望的基准时间：2025年12月3日 14:30:45 UTC+8
	expectedDate := time.Date(2025, 12, 3, 14, 30, 45, 0, time.UTC)
	expectedDateOnly := time.Date(2025, 12, 3, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		input       string
		expected    time.Time
		shouldError bool
		description string
	}{
		// 标准日期格式
		{
			name:        "DateOnly",
			input:       testDateStr,
			expected:    expectedDateOnly,
			shouldError: false,
			description: "标准日期格式 YYYY-MM-DD",
		},
		{
			name:        "RFC3339",
			input:       testRFC3339Str,
			expected:    expectedDate,
			shouldError: false,
			description: "RFC3339 格式",
		},
		{
			name:        "DateTime",
			input:       testFullDateTime,
			expected:    expectedDate,
			shouldError: false,
			description: "日期时间格式",
		},

		// ISO 格式变体
		{
			name:        "ISO_NoTimezone",
			input:       testRFC3339Str,
			expected:    expectedDate,
			shouldError: false,
			description: "ISO 格式（无时区）",
		},
		{
			name:        "ISO_UTC_Z",
			input:       testRFC3339UTCStr,
			expected:    expectedDate,
			shouldError: false,
			description: "ISO UTC 格式",
		},
		{
			name:        "ISO_Milliseconds_UTC",
			input:       testRFC3339MillisStr,
			expected:    expectedDate,
			shouldError: false,
			description: "ISO 毫秒 UTC 格式",
		},

		// 带时区偏移的格式
		{
			name:        "ISO_PositiveOffset",
			input:       testRFC3339PlusTZStr, // UTC+8 时区，所以是 22:30
			expected:    expectedDate,
			shouldError: false,
			description: "ISO 正时区偏移格式",
		},
		{
			name:        "ISO_NegativeOffset",
			input:       testRFC3339MinusTZStr, // UTC-8 时区，所以是 06:30
			expected:    expectedDate,
			shouldError: false,
			description: "ISO 负时区偏移格式",
		},

		// 中国常用格式
		{
			name:        "Chinese_Date_WithZero",
			input:       "2025年12月03日",
			expected:    expectedDateOnly,
			shouldError: false,
			description: "中文日期格式（带前导零）",
		},
		{
			name:        "Chinese_Date_NoZero",
			input:       testChineseDateStr,
			expected:    expectedDateOnly,
			shouldError: false,
			description: "中文日期格式（无前导零）",
		},

		// 斜杠分隔格式
		{
			name:        "Slash_WithZero",
			input:       testSlashDateStr,
			expected:    expectedDateOnly,
			shouldError: false,
			description: "斜杠日期格式（带前导零）",
		},
		{
			name:        "Slash_NoZero",
			input:       "2025/12/3",
			expected:    expectedDateOnly,
			shouldError: false,
			description: "斜杠日期格式（无前导零）",
		},
		{
			name:        "US_Format",
			input:       "12/03/2025",
			expected:    expectedDateOnly,
			shouldError: false,
			description: "美式日期格式 MM/DD/YYYY",
		},

		// 时间格式（今日日期）
		{
			name:        "TimeOnly_24Hour",
			input:       testTimeStr,
			expected:    time.Date(0, 1, 1, 14, 30, 45, 0, time.UTC), // time.Parse 会使用年份 0
			shouldError: false,
			description: "24小时制时间格式",
		},

		// 完整日期时间格式
		{
			name:        "DateTime_Space",
			input:       testDateTimeStr,
			expected:    expectedDate,
			shouldError: false,
			description: "空格分隔的日期时间",
		},
		{
			name:        "DateTime_Slash_Space",
			input:       testSlashDateStr + " 14:30:45",
			expected:    expectedDate,
			shouldError: false,
			description: "斜杠日期 + 空格 + 时间",
		},

		// 数据库常见格式
		{
			name:        "DB_Milliseconds",
			input:       testDateTimeMillisStr,
			expected:    expectedDate,
			shouldError: false,
			description: "数据库毫秒格式",
		},

		// 紧凑格式
		{
			name:        "Compact_Date",
			input:       testCompactDateStr,
			expected:    expectedDateOnly,
			shouldError: false,
			description: "紧凑日期格式 YYYYMMDD",
		},
		{
			name:        "Compact_DateTime",
			input:       "20251203143045",
			expected:    expectedDate,
			shouldError: false,
			description: "紧凑日期时间格式 YYYYMMDDHHMMSS",
		},

		// 错误情况
		{
			name:        "Invalid_Format",
			input:       "invalid-date",
			expected:    time.Time{},
			shouldError: true,
			description: "无效的日期格式",
		},
		{
			name:        "Empty_String",
			input:       "",
			expected:    time.Time{},
			shouldError: true,
			description: "空字符串",
		},
		{
			name:        "Invalid_Date_Values",
			input:       "2025-13-45",
			expected:    time.Time{},
			shouldError: true,
			description: "无效的日期数值",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFlexibleDate(tt.input)

			if tt.shouldError {
				assert.Error(t, err, "应该返回错误: %s", tt.description)
				assert.True(t, result.IsZero(), "错误时应返回零值时间")
			} else {
				assert.NoError(t, err, "不应该返回错误: %s", tt.description)

				// 对于只有日期的格式，只比较日期部分
				if tt.input == testTimeStr {
					// 时间格式特殊处理，只比较时分秒
					assert.Equal(t, tt.expected.Hour(), result.Hour(), "小时应该匹配")
					assert.Equal(t, tt.expected.Minute(), result.Minute(), "分钟应该匹配")
					assert.Equal(t, tt.expected.Second(), result.Second(), "秒应该匹配")
				} else if containsOnlyDate(tt.input) {
					// 只包含日期的格式，只比较日期部分
					assert.Equal(t, tt.expected.Year(), result.Year(), "年应该匹配")
					assert.Equal(t, tt.expected.Month(), result.Month(), "月应该匹配")
					assert.Equal(t, tt.expected.Day(), result.Day(), "日应该匹配")
				} else {
					// 包含日期和时间的格式，比较完整时间（允许时区差异）
					expectedUTC := tt.expected.UTC()
					resultUTC := result.UTC()
					assert.True(t,
						expectedUTC.Equal(resultUTC) ||
							expectedUTC.Sub(resultUTC).Abs() < time.Minute,
						"时间应该匹配（允许1分钟内差异）: expected %v, got %v", expectedUTC, resultUTC)
				}
			}
		})
	}
}

// containsOnlyDate 判断输入字符串是否只包含日期（不包含时间）
func containsOnlyDate(input string) bool {
	dateOnlyFormats := []string{
		testDateStr, testChineseDateStr, "2025年12月3日",
		"2025/12/3", testSlashDateStr, "12/03/2025", "12/3/2025",
		"03/12/2025", "3/12/2025", testCompactDateStr,
	}

	for _, format := range dateOnlyFormats {
		if input == format {
			return true
		}
	}
	return false
}

// TestParseFlexibleDate_EdgeCases 测试边界情况
func TestParseFlexibleDateEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name        string
		input       string
		shouldError bool
		description string
	}{
		{
			name:        "LeapYear_Feb29",
			input:       "2024-02-29",
			shouldError: false,
			description: "闰年2月29日",
		},
		{
			name:        "NonLeapYear_Feb29",
			input:       "2025-02-29",
			shouldError: true,
			description: "非闰年2月29日",
		},
		{
			name:        "Year_1900",
			input:       "1900-01-01",
			shouldError: false,
			description: "1900年（特殊年份）",
		},
		{
			name:        "Year_2000",
			input:       "2000-01-01",
			shouldError: false,
			description: "2000年（闰年）",
		},
		{
			name:        "Whitespace_Prefix",
			input:       testWhitespaceStr,
			shouldError: true, // 我们的函数不处理前后空格
			description: "带前后空格的日期",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseFlexibleDate(tc.input)

			if tc.shouldError {
				assert.Error(t, err, "应该返回错误: %s", tc.description)
			} else {
				assert.NoError(t, err, "不应该返回错误: %s", tc.description)
			}
		})
	}
}

// BenchmarkParseFlexibleDate 性能基准测试
func BenchmarkParseFlexibleDate(b *testing.B) {
	testCases := []string{
		testDateStr,        // 最常用格式，应该最快
		testRFC3339UTCStr,  // RFC3339
		testChineseDateStr, // 中文格式
		testCompactDateStr, // 紧凑格式
		"invalid-date",     // 无效格式，会尝试所有格式
	}

	for _, tc := range testCases {
		b.Run(tc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ParseFlexibleDate(tc)
			}
		})
	}
}

// TestParseFlexibleDate_Concurrent 并发安全测试
func TestParseFlexibleDateConcurrent(t *testing.T) {
	const numGoroutines = 100
	const numIterations = 1000

	testInputs := []string{
		testDateStr,
		testRFC3339UTCStr,
		testChineseDateStr,
		testSlashDateStr,
	}

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIterations; j++ {
				for _, input := range testInputs {
					_, err := ParseFlexibleDate(input)
					assert.NoError(t, err, "并发调用不应该出错")
				}
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
