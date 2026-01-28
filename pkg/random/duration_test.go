/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-28 00:00:00
 * @FilePath: \go-toolbox\pkg\random\duration_test.go
 * @Description: 随机时间相关工具函数测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package random

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandDuration(t *testing.T) {
	tests := []struct {
		name string
		min  time.Duration
		max  time.Duration
	}{
		{
			name: "毫秒级别",
			min:  100 * time.Millisecond,
			max:  500 * time.Millisecond,
		},
		{
			name: "秒级别",
			min:  1 * time.Second,
			max:  5 * time.Second,
		},
		{
			name: "分钟级别",
			min:  1 * time.Minute,
			max:  5 * time.Minute,
		},
		{
			name: "相同值",
			min:  100 * time.Millisecond,
			max:  100 * time.Millisecond,
		},
		{
			name: "反向范围（自动交换）",
			min:  500 * time.Millisecond,
			max:  100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				result := RandDuration(tt.min, tt.max)

				min, max := tt.min, tt.max
				if max < min {
					min, max = max, min
				}

				if min != max {
					assert.GreaterOrEqual(t, result, min, "结果应该大于等于最小值")
					assert.Less(t, result, max, "结果应该小于最大值")
				} else {
					assert.Equal(t, min, result, "相同值时应该返回该值")
				}
			}
		})
	}
}

func TestRandTimeBetween(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	tests := []struct {
		name  string
		start time.Time
		end   time.Time
	}{
		{
			name:  "过去到现在",
			start: yesterday,
			end:   now,
		},
		{
			name:  "现在到未来",
			start: now,
			end:   tomorrow,
		},
		{
			name:  "相同时间",
			start: now,
			end:   now,
		},
		{
			name:  "反向范围（自动交换）",
			start: tomorrow,
			end:   yesterday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 50; i++ {
				result := RandTimeBetween(tt.start, tt.end)

				start, end := tt.start, tt.end
				if end.Before(start) {
					start, end = end, start
				}

				if !start.Equal(end) {
					assert.True(t, !result.Before(start), "结果应该不早于开始时间")
					assert.True(t, result.Before(end) || result.Equal(end), "结果应该不晚于结束时间")
				} else {
					assert.Equal(t, start, result, "相同时间应该返回该时间")
				}
			}
		})
	}
}

func TestRandTimeInPast(t *testing.T) {
	durations := []time.Duration{
		1 * time.Hour,
		24 * time.Hour,
		7 * 24 * time.Hour,
	}

	for _, duration := range durations {
		t.Run(duration.String(), func(t *testing.T) {
			now := time.Now()
			for i := 0; i < 50; i++ {
				result := RandTimeInPast(duration)
				assert.True(t, !result.After(now), "结果应该不晚于现在")
				assert.True(t, !result.Before(now.Add(-duration)), "结果应该在指定范围内")
			}
		})
	}
}

func TestRandTimeInFuture(t *testing.T) {
	durations := []time.Duration{
		1 * time.Hour,
		24 * time.Hour,
		7 * 24 * time.Hour,
	}

	for _, duration := range durations {
		t.Run(duration.String(), func(t *testing.T) {
			now := time.Now()
			for i := 0; i < 50; i++ {
				result := RandTimeInFuture(duration)
				assert.True(t, !result.Before(now), "结果应该不早于现在")
				assert.True(t, !result.After(now.Add(duration)), "结果应该在指定范围内")
			}
		})
	}
}

func TestRandDate(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 100; i++ {
		result := RandDate(start, end)

		// 检查时间部分是否为 00:00:00
		assert.Equal(t, 0, result.Hour(), "小时应该为0")
		assert.Equal(t, 0, result.Minute(), "分钟应该为0")
		assert.Equal(t, 0, result.Second(), "秒应该为0")

		// 检查日期范围
		assert.True(t, !result.Before(start), "结果应该不早于开始日期")
		assert.True(t, !result.After(end), "结果应该不晚于结束日期")
	}
}

func TestRandWeekday(t *testing.T) {
	weekdays := make(map[time.Weekday]bool)

	// 多次调用，确保能生成所有星期
	for i := 0; i < 1000; i++ {
		weekday := RandWeekday()
		assert.GreaterOrEqual(t, int(weekday), 0, "星期应该>=0")
		assert.LessOrEqual(t, int(weekday), 6, "星期应该<=6")
		weekdays[weekday] = true
	}

	// 检查是否生成了多种星期
	assert.Greater(t, len(weekdays), 1, "应该生成多种星期")
}

func TestRandMonth(t *testing.T) {
	months := make(map[time.Month]bool)

	// 多次调用，确保能生成所有月份
	for i := 0; i < 1000; i++ {
		month := RandMonth()
		assert.GreaterOrEqual(t, int(month), 1, "月份应该>=1")
		assert.LessOrEqual(t, int(month), 12, "月份应该<=12")
		months[month] = true
	}

	// 检查是否生成了多种月份
	assert.Greater(t, len(months), 1, "应该生成多种月份")
}

func TestRandHour(t *testing.T) {
	for i := 0; i < 100; i++ {
		hour := RandHour()
		assert.GreaterOrEqual(t, hour, 0, "小时应该>=0")
		assert.LessOrEqual(t, hour, 23, "小时应该<=23")
	}
}

func TestRandMinute(t *testing.T) {
	for i := 0; i < 100; i++ {
		minute := RandMinute()
		assert.GreaterOrEqual(t, minute, 0, "分钟应该>=0")
		assert.LessOrEqual(t, minute, 59, "分钟应该<=59")
	}
}

func TestRandSecond(t *testing.T) {
	for i := 0; i < 100; i++ {
		second := RandSecond()
		assert.GreaterOrEqual(t, second, 0, "秒应该>=0")
		assert.LessOrEqual(t, second, 59, "秒应该<=59")
	}
}

func TestRandTimeOfDay(t *testing.T) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	for i := 0; i < 100; i++ {
		result := RandTimeOfDay()

		// 检查是否在今天范围内
		assert.True(t, !result.Before(startOfDay), "应该不早于今天00:00")
		assert.True(t, result.Before(endOfDay), "应该早于明天00:00")

		// 检查日期是否是今天
		assert.Equal(t, now.Year(), result.Year(), "年份应该相同")
		assert.Equal(t, now.Month(), result.Month(), "月份应该相同")
		assert.Equal(t, now.Day(), result.Day(), "日期应该相同")
	}
}

func TestRandBusinessHour(t *testing.T) {
	now := time.Now()

	for i := 0; i < 100; i++ {
		result := RandBusinessHour()

		// 检查小时范围
		hour := result.Hour()
		assert.GreaterOrEqual(t, hour, 9, "小时应该>=9")
		assert.Less(t, hour, 18, "小时应该<18")

		// 检查日期是否是今天
		assert.Equal(t, now.Year(), result.Year(), "年份应该相同")
		assert.Equal(t, now.Month(), result.Month(), "月份应该相同")
		assert.Equal(t, now.Day(), result.Day(), "日期应该相同")
	}
}

func TestRandUnixTimestamp(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC).Unix()

	for i := 0; i < 100; i++ {
		result := RandUnixTimestamp(start, end)
		assert.GreaterOrEqual(t, result, start, "时间戳应该>=开始时间")
		assert.LessOrEqual(t, result, end, "时间戳应该<=结束时间")
	}

	// 测试相同值
	same := RandUnixTimestamp(start, start)
	assert.Equal(t, start, same, "相同值应该返回该值")

	// 测试反向范围
	reversed := RandUnixTimestamp(end, start)
	assert.GreaterOrEqual(t, reversed, start, "反向范围应该自动交换")
	assert.LessOrEqual(t, reversed, end, "反向范围应该自动交换")
}

// 基准测试
func BenchmarkRandDuration(b *testing.B) {
	min := 100 * time.Millisecond
	max := 500 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RandDuration(min, max)
	}
}

func BenchmarkRandTimeBetween(b *testing.B) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RandTimeBetween(yesterday, now)
	}
}

func BenchmarkRandTimeInPast(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RandTimeInPast(24 * time.Hour)
	}
}

func BenchmarkRandTimeOfDay(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RandTimeOfDay()
	}
}

func BenchmarkRandBusinessHour(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RandBusinessHour()
	}
}

// 测试随机性分布
func TestRandDurationDistribution(t *testing.T) {
	min := 100 * time.Millisecond
	max := 500 * time.Millisecond
	samples := 10000

	// 统计分布
	buckets := make([]int, 5) // 分成5个桶
	bucketSize := (max - min) / time.Duration(len(buckets))

	for i := 0; i < samples; i++ {
		result := RandDuration(min, max)
		bucketIndex := int((result - min) / bucketSize)
		if bucketIndex >= len(buckets) {
			bucketIndex = len(buckets) - 1
		}
		buckets[bucketIndex]++
	}

	// 检查每个桶的样本数是否在合理范围内（允许20%的偏差）
	expectedPerBucket := samples / len(buckets)
	tolerance := expectedPerBucket / 5 // 20%的容差

	for i, count := range buckets {
		assert.GreaterOrEqual(t, count, expectedPerBucket-tolerance,
			"桶 %d 的样本数应该在合理范围内", i)
		assert.LessOrEqual(t, count, expectedPerBucket+tolerance,
			"桶 %d 的样本数应该在合理范围内", i)
	}

	t.Logf("分布统计: %v (期望每桶约 %d 个样本)", buckets, expectedPerBucket)
}
