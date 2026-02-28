/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-28 00:00:00
 * @FilePath: \go-toolbox\pkg\convert\fast_format_test.go
 * @Description: 快速格式化函数测试和性能对比
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFastAppendInt 测试快速整数追加功能
func TestFastAppendInt(t *testing.T) {
	tests := []struct {
		name     string
		val      int
		expected string
	}{
		{"零", 0, "0"},
		{"单位数", 5, "5"},
		{"两位数", 42, "42"},
		{"三位数", 123, "123"},
		{"四位数", 1234, "1234"},
		{"大数", 123567, "123567"},
		{"负数单位", -5, "-5"},
		{"负数两位", -42, "-42"},
		{"负数三位", -123, "-123"},
		{"负数大数", -123567, "-123567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 0, 32)
			result := FastAppendInt(buf, tt.val)
			assert.Equal(t, tt.expected, string(result), "FastAppendInt(%d) should return %s", tt.val, tt.expected)
		})
	}
}

// TestFastFormatTime 测试快速时间格式化功能
func TestFastFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{"标准时间", time.Date(2026, 2, 28, 18, 32, 7, 0, time.UTC), "2026/2/28 18:32:07 "},
		{"单位数月日", time.Date(2026, 1, 5, 9, 5, 3, 0, time.UTC), "2026/1/5 9:05:03 "},
		{"双位数月日", time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), "2026/12/31 23:59:59 "},
		{"零点时刻", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), "2026/1/1 0:00:00 "},
		{"午夜前一秒", time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), "2026/12/31 23:59:59 "},
		{"正午时刻", time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC), "2026/6/15 12:00:00 "},
		{"早晨时刻", time.Date(2026, 3, 10, 6, 30, 56, 0, time.UTC), "2026/3/10 6:30:56 "},
		{"傍晚时刻", time.Date(2026, 9, 20, 18, 56, 30, 0, time.UTC), "2026/9/20 18:56:30 "},
		{"闰年2月29日", time.Date(2024, 2, 29, 15, 20, 10, 0, time.UTC), "2024/2/29 15:20:10 "},
		{"年初", time.Date(2026, 1, 1, 0, 0, 1, 0, time.UTC), "2026/1/1 0:00:01 "},
		{"年末", time.Date(2026, 12, 31, 23, 59, 58, 0, time.UTC), "2026/12/31 23:59:58 "},
		{"单位数时分秒", time.Date(2026, 5, 5, 5, 5, 5, 0, time.UTC), "2026/5/5 5:05:05 "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 0, 64)
			result := FastFormatTime(buf, tt.time)
			assert.Equal(t, tt.expected, string(result), "FastFormatTime should format time correctly")
		})
	}
}

// TestFastFormatTimeISO 测试 ISO 格式时间格式化
func TestFastFormatTimeISO(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{"标准时间", time.Date(2026, 2, 28, 18, 32, 7, 0, time.UTC), "2026-02-28 18:32:07"},
		{"单位数月日", time.Date(2026, 1, 5, 9, 5, 3, 0, time.UTC), "2026-01-05 09:05:03"},
		{"双位数月日", time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), "2026-12-31 23:59:59"},
		{"零点时刻", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), "2026-01-01 00:00:00"},
		{"正午时刻", time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC), "2026-06-15 12:00:00"},
		{"闰年2月29日", time.Date(2024, 2, 29, 15, 20, 10, 0, time.UTC), "2024-02-29 15:20:10"},
		{"年初", time.Date(2026, 1, 1, 0, 0, 1, 0, time.UTC), "2026-01-01 00:00:01"},
		{"年末", time.Date(2026, 12, 31, 23, 59, 58, 0, time.UTC), "2026-12-31 23:59:58"},
		{"早晨时刻", time.Date(2026, 3, 10, 6, 30, 56, 0, time.UTC), "2026-03-10 06:30:56"},
		{"傍晚时刻", time.Date(2026, 9, 20, 18, 56, 30, 0, time.UTC), "2026-09-20 18:56:30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 0, 64)
			result := FastFormatTimeISO(buf, tt.time)
			assert.Equal(t, tt.expected, string(result), "FastFormatTimeISO should format time in ISO format")
		})
	}
}

// TestFastFormatTimeCompact 测试紧凑格式时间格式化
func TestFastFormatTimeCompact(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{"标准时间", time.Date(2026, 2, 28, 18, 32, 7, 0, time.UTC), "20260228183207"},
		{"单位数月日", time.Date(2026, 1, 5, 9, 5, 3, 0, time.UTC), "20260105090503"},
		{"双位数月日", time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), "20261231235959"},
		{"零点时刻", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), "20260101000000"},
		{"正午时刻", time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC), "20260615120000"},
		{"闰年2月29日", time.Date(2024, 2, 29, 15, 20, 10, 0, time.UTC), "20240229152010"},
		{"年初", time.Date(2026, 1, 1, 0, 0, 1, 0, time.UTC), "20260101000001"},
		{"年末", time.Date(2026, 12, 31, 23, 59, 58, 0, time.UTC), "20261231235958"},
		{"早晨时刻", time.Date(2026, 3, 10, 6, 30, 56, 0, time.UTC), "20260310063056"},
		{"傍晚时刻", time.Date(2026, 9, 20, 18, 56, 30, 0, time.UTC), "20260920185630"},
		{"单位数时分秒", time.Date(2026, 5, 5, 5, 5, 5, 0, time.UTC), "20260505050505"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 0, 64)
			result := FastFormatTimeCompact(buf, tt.time)
			assert.Equal(t, tt.expected, string(result), "FastFormatTimeCompact should format time in compact format")
		})
	}
}

// TestFastItoa 测试快速整数转字符串功能
func TestFastItoa(t *testing.T) {
	tests := []struct {
		name     string
		val      int
		expected string
	}{
		{"零", 0, "0"},
		{"单位数", 5, "5"},
		{"两位数", 42, "42"},
		{"三位数_100", 100, "100"},
		{"三位数_123", 123, "123"},
		{"三位数_999", 999, "999"},
		{"四位数", 1234, "1234"},
		{"大数", 123567, "123567"},
		{"负数", -123, "-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FastItoa(tt.val)
			assert.Equal(t, tt.expected, result, "FastItoa(%d) should return %s", tt.val, tt.expected)
		})
	}
}

// TestFastFloat 测试快速浮点数转字符串功能
func TestFastFloat(t *testing.T) {
	tests := []struct {
		name     string
		val      float64
		prec     int
		expected string
	}{
		// 基本场景
		{"两位小数_正数", 3.14159, 2, "3.14"},
		{"两位小数_整数", 100.0, 2, "100.00"},
		{"一位小数", 123.456, 1, "123.5"},
		{"零位小数", 99.99, 0, "100"},
		{"三位小数", 1.23456, 3, "1.235"},
		{"四位小数", 0.123456, 4, "0.1235"},

		// 最少位数（-1）
		{"最少位数_整数", 100.0, -1, "100"},
		{"最少位数_小数", 3.14, -1, "3.14"},
		{"最少位数_零", 0.0, -1, "0"},

		// 负数
		{"负数_两位", -123.45, 2, "-123.45"},
		{"负数_零位", -99.99, 0, "-100"},
		{"负数_小数", -0.123, 3, "-0.123"},

		// 零值
		{"零值_两位", 0.0, 2, "0.00"},
		{"零值_零位", 0.0, 0, "0"},
		{"零值_最少", 0.0, -1, "0"},

		// 小数
		{"小数_科学计数", 0.00123, 5, "0.00123"},
		{"小数_极小", 0.0001, 4, "0.0001"},
		{"小数_四舍五入", 0.9999, 2, "1.00"},

		// 大数
		{"大数_两位", 123456.789, 2, "123456.79"},
		{"大数_零位", 123456.789, 0, "123457"},
		{"大数_整数", 999999.0, 2, "999999.00"},

		// 边界值
		{"极小正数", 0.000001, 6, "0.000001"},
		{"接近零", 0.00000001, 8, "0.00000001"},
		{"1.0", 1.0, 2, "1.00"},
		{"0.1", 0.1, 1, "0.1"},
		{"0.01", 0.01, 2, "0.01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FastFloat(tt.val, tt.prec)
			assert.Equal(t, tt.expected, result, "FastFloat(%f, %d) should return %s", tt.val, tt.prec, tt.expected)
		})
	}
}

// ============================================================================
// 性能对比基准测试
// ============================================================================

// BenchmarkFastAppendInt_vs_Itoa 对比 FastAppendInt 和 strconv.Itoa 的性能
func BenchmarkFastAppendInt_vs_Itoa(b *testing.B) {
	testCases := []struct {
		name string
		val  int
	}{
		{"单位数", 5},
		{"两位数", 42},
		{"三位数", 123},
		{"四位数", 1234},
		{"大数", 123567},
	}

	for _, tc := range testCases {
		b.Run("FastAppendInt_"+tc.name, func(b *testing.B) {
			buf := make([]byte, 0, 32)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf = FastAppendInt(buf[:0], tc.val)
			}
		})

		b.Run("strconv.Itoa_"+tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = strconv.Itoa(tc.val)
			}
		})

		b.Run("strconv.AppendInt_"+tc.name, func(b *testing.B) {
			buf := make([]byte, 0, 32)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf = strconv.AppendInt(buf[:0], int64(tc.val), 10)
			}
		})
	}
}

// BenchmarkFastFormatTime_vs_Format 对比 FastFormatTime 和 time.Format 的性能
func BenchmarkFastFormatTime_vs_Format(b *testing.B) {
	testTime := time.Date(2026, 2, 28, 18, 32, 7, 0, time.UTC)

	b.Run("FastFormatTime", func(b *testing.B) {
		buf := make([]byte, 0, 64)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf = FastFormatTime(buf[:0], testTime)
		}
	})

	b.Run("time.Format", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = testTime.Format("2006/1/2 15:04:05 ")
		}
	})

	b.Run("FastFormatTimeISO", func(b *testing.B) {
		buf := make([]byte, 0, 64)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf = FastFormatTimeISO(buf[:0], testTime)
		}
	})

	b.Run("time.Format_ISO", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = testTime.Format("2006-01-02 15:04:05")
		}
	})

	b.Run("FastFormatTimeCompact", func(b *testing.B) {
		buf := make([]byte, 0, 64)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf = FastFormatTimeCompact(buf[:0], testTime)
		}
	})

	b.Run("time.Format_Compact", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = testTime.Format("20060102150405")
		}
	})
}

// BenchmarkMemoryAllocation 内存分配对比
func BenchmarkMemoryAllocation(b *testing.B) {
	testTime := time.Date(2026, 2, 28, 18, 32, 7, 0, time.UTC)

	b.Run("FastFormatTime_NoAlloc", func(b *testing.B) {
		buf := make([]byte, 0, 64)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf = FastFormatTime(buf[:0], testTime)
		}
	})

	b.Run("time.Format_WithAlloc", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = testTime.Format("2006/1/2 15:04:05 ")
		}
	})

	b.Run("FastAppendInt_NoAlloc", func(b *testing.B) {
		buf := make([]byte, 0, 32)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf = FastAppendInt(buf[:0], 123)
		}
	})

	b.Run("strconv.Itoa_WithAlloc", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = strconv.Itoa(123)
		}
	})
}

// BenchmarkConcurrent 并发性能测试
func BenchmarkConcurrent(b *testing.B) {
	testTime := time.Date(2026, 2, 28, 18, 32, 7, 0, time.UTC)

	b.Run("FastFormatTime_Concurrent", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			buf := make([]byte, 0, 64)
			for pb.Next() {
				buf = FastFormatTime(buf[:0], testTime)
			}
		})
	})

	b.Run("time.Format_Concurrent", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = testTime.Format("2006/1/2 15:04:05 ")
			}
		})
	})
}

// BenchmarkRealWorldScenario 真实场景性能测试（日志记录）
func BenchmarkRealWorldScenario(b *testing.B) {
	testTime := time.Date(2026, 2, 28, 18, 32, 7, 0, time.UTC)

	b.Run("FastFormat_LogLine", func(b *testing.B) {
		buf := make([]byte, 0, 128)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf = buf[:0]
			buf = FastFormatTime(buf, testTime)
			buf = append(buf, "[INFO] "...)
			buf = append(buf, "Log message with ID: "...)
			buf = FastAppendInt(buf, 12356)
			buf = append(buf, '\n')
		}
	})

	b.Run("StandardFormat_LogLine", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = testTime.Format("2006/1/2 15:04:05 ") + "[INFO] " + "Log message with ID: " + strconv.Itoa(12356) + "\n"
		}
	})
}

// BenchmarkFastItoa_vs_Itoa 对比 FastItoa 和 strconv.Itoa 的性能
func BenchmarkFastItoa_vs_Itoa(b *testing.B) {
	testCases := []struct {
		name string
		val  int
	}{
		// 单位数（0-9）
		{"单位数_0", 0},
		{"单位数_5", 5},
		{"单位数_9", 9},

		// 两位数（10-99）
		{"两位数_10", 10},
		{"两位数_42", 42},
		{"两位数_99", 99},

		// 三位数（100-999）
		{"三位数_100", 100},
		{"三位数_500", 500},
		{"三位数_999", 999},

		// 四位数（1000-9999）
		{"四位数_1000", 1000},
		{"四位数_1234", 1234},
		{"四位数_9999", 9999},

		// 五位数及以上
		{"五位数_10000", 10000},
		{"五位数_50000", 50000},
		{"六位数_123567", 123567},
		{"六位数_999999", 999999},

		// 负数
		{"负数_单位", -5},
		{"负数_两位", -42},
		{"负数_三位", -123},
		{"负数_四位", -1234},
		{"负数_大数", -123567},
	}

	for _, tc := range testCases {
		b.Run("FastItoa_"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FastItoa(tc.val)
			}
		})

		b.Run("strconv.Itoa_"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = strconv.Itoa(tc.val)
			}
		})
	}
}

// BenchmarkFastFloat_vs_FormatFloat 对比 FastFloat 和 strconv.FormatFloat 的性能
func BenchmarkFastFloat_vs_FormatFloat(b *testing.B) {
	testCases := []struct {
		name string
		val  float64
		prec int
	}{
		{"两位小数", 3.14159, 2},
		{"一位小数", 123.456, 1},
		{"零位小数", 99.99, 0},
		{"最少位数", 100.0, -1},
	}

	for _, tc := range testCases {
		b.Run("FastFloat_"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FastFloat(tc.val, tc.prec)
			}
		})

		b.Run("strconv.FormatFloat_"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = strconv.FormatFloat(tc.val, 'f', tc.prec, 64)
			}
		})
	}
}
