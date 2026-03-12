/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-03-12 16:30:21
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-03-12 16:35:55
 * @FilePath: \go-toolbox\pkg\convert\format_stats.go
 * @Description: 统计数据格式化工具 - 用于格式化时长、数量等统计指标
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"fmt"
	"strings"
)

// FormatDuration 格式化时长（秒 -> 人类可读格式）
// 根据时长自动选择合适的单位（年、月、日、时、分、秒）
//
// 参数：
//   - value: 时长值（any 类型，支持 float64、int、int64 等）
//
// 返回：
//   - 格式化后的字符串
//   - 如果值为 nil 或 <= 0，返回 "N/A"
//
// 时间单位换算：
//   - 1 分钟 = 60 秒
//   - 1 小时 = 3600 秒
//   - 1 天 = 86400 秒
//   - 1 月 = 2592000 秒（30天）
//   - 1 年 = 31536000 秒（365天）
//
// 示例：
//   - FormatDuration(45) -> "45s"
//   - FormatDuration(90) -> "1m 30s"
//   - FormatDuration(3665) -> "1h 1m 5s"
//   - FormatDuration(86400) -> "1d"
//   - FormatDuration(2678400) -> "1mo"
//   - FormatDuration(31536000) -> "1y"
//   - FormatDuration(nil) -> "N/A"
func FormatDuration(value any) string {
	if value == nil {
		return "N/A"
	}

	// 使用 ToFloat64 统一转换
	seconds, err := ToFloat64(value)
	if err != nil || seconds <= 0 {
		return "N/A"
	}

	return formatDurationFromSeconds(int64(seconds))
}

// formatDurationFromSeconds 从秒数格式化为人类可读格式
func formatDurationFromSeconds(totalSeconds int64) string {
	if totalSeconds <= 0 {
		return "N/A"
	}

	const (
		secondsPerMinute = 60
		secondsPerHour   = 3600
		secondsPerDay    = 86400
		secondsPerMonth  = 2592000  // 30 天
		secondsPerYear   = 31536000 // 365 天
	)

	var parts []string

	// 年
	if totalSeconds >= secondsPerYear {
		years := totalSeconds / secondsPerYear
		parts = append(parts, fmt.Sprintf("%dy", years))
		totalSeconds %= secondsPerYear
	}

	// 月
	if totalSeconds >= secondsPerMonth {
		months := totalSeconds / secondsPerMonth
		parts = append(parts, fmt.Sprintf("%dmo", months))
		totalSeconds %= secondsPerMonth
	}

	// 天
	if totalSeconds >= secondsPerDay {
		days := totalSeconds / secondsPerDay
		parts = append(parts, fmt.Sprintf("%dd", days))
		totalSeconds %= secondsPerDay
	}

	// 小时
	if totalSeconds >= secondsPerHour {
		hours := totalSeconds / secondsPerHour
		parts = append(parts, fmt.Sprintf("%dh", hours))
		totalSeconds %= secondsPerHour
	}

	// 分钟
	if totalSeconds >= secondsPerMinute {
		minutes := totalSeconds / secondsPerMinute
		parts = append(parts, fmt.Sprintf("%dm", minutes))
		totalSeconds %= secondsPerMinute
	}

	// 秒
	if totalSeconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", totalSeconds))
	}

	// 最多显示前 3 个单位，避免过长
	if len(parts) > 3 {
		parts = parts[:3]
	}

	return strings.Join(parts, " ")
}

// FormatCount 格式化数量（处理 nil 值）
// 将数量值格式化为字符串，nil 值返回 "0"
//
// 参数：
//   - value: 数量值（any 类型，支持各种数值类型）
//
// 返回：
//   - 格式化后的字符串
//   - 如果值为 nil，返回 "0"
//
// 示例：
//   - FormatCount(123) -> "123"
//   - FormatCount(nil) -> "0"
//   - FormatCount(0) -> "0"
func FormatCount(value any) string {
	if value == nil {
		return "0"
	}
	return MustString(value)
}

// FormatPercentage 格式化百分比
// 将浮点数格式化为百分比字符串
//
// 参数：
//   - value: 百分比值（any 类型，支持 float64、int 等）
//   - precision: 小数位数（默认 1）
//
// 返回：
//   - 格式化后的百分比字符串（如 "85.5%"）
//   - 如果值为 nil，返回 "0%"
//
// 示例：
//   - FormatPercentage(85.567, 1) -> "85.6%"
//   - FormatPercentage(nil, 1) -> "0%"
//   - FormatPercentage(100, 0) -> "100%"
func FormatPercentage(value any, precision int) string {
	if value == nil {
		return "0%"
	}

	rate, err := ToFloat64(value)
	if err != nil {
		return "0%"
	}

	format := fmt.Sprintf("%%.%df%%%%", precision)
	return fmt.Sprintf(format, rate)
}
