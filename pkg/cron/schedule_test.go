/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-25 15:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-25 11:15:55
 * @FilePath: \go-toolbox\pkg\cron\schedule_test.go
 * @Description: Schedule 测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCronSpecSchedule_Next_Daily(t *testing.T) {
	// 每天午夜执行
	schedule, err := ParseCronWithSeconds("0 0 0 * * *")
	assert.NoError(t, err)

	now := time.Date(2025, 12, 25, 10, 30, 0, 0, time.UTC)
	next := schedule.Next(now)

	expected := time.Date(2025, 12, 26, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_Weekly(t *testing.T) {
	// 每周日午夜执行
	schedule, err := ParseCronWithSeconds("0 0 0 * * 0")
	assert.NoError(t, err)

	// 从周四开始
	now := time.Date(2025, 12, 25, 10, 0, 0, 0, time.UTC) // Thursday
	next := schedule.Next(now)

	// 应该是下个周日
	expected := time.Date(2025, 12, 28, 0, 0, 0, 0, time.UTC) // Sunday
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_Monthly(t *testing.T) {
	// 每月1号午夜执行
	schedule, err := ParseCronWithSeconds("0 0 0 1 * *")
	assert.NoError(t, err)

	now := time.Date(2025, 12, 25, 10, 0, 0, 0, time.UTC)
	next := schedule.Next(now)

	expected := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_Yearly(t *testing.T) {
	// 每年1月1日午夜执行
	schedule, err := ParseCronWithSeconds("0 0 0 1 1 *")
	assert.NoError(t, err)

	now := time.Date(2025, 12, 25, 10, 0, 0, 0, time.UTC)
	next := schedule.Next(now)

	expected := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_EveryMinute(t *testing.T) {
	// 每分钟执行
	schedule, err := ParseCronWithSeconds("0 * * * * *")
	assert.NoError(t, err)

	now := time.Date(2025, 12, 25, 10, 30, 30, 0, time.UTC)
	next := schedule.Next(now)

	expected := time.Date(2025, 12, 25, 10, 31, 0, 0, time.UTC)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_EveryHour(t *testing.T) {
	// 每小时执行
	schedule, err := ParseCronWithSeconds("0 0 * * * *")
	assert.NoError(t, err)

	now := time.Date(2025, 12, 25, 10, 30, 0, 0, time.UTC)
	next := schedule.Next(now)

	expected := time.Date(2025, 12, 25, 11, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_SpecificTime(t *testing.T) {
	// 每天9:30执行
	schedule, err := ParseCronWithSeconds("0 30 9 * * *")
	assert.NoError(t, err)

	now := time.Date(2025, 12, 25, 8, 0, 0, 0, time.UTC)
	next := schedule.Next(now)

	expected := time.Date(2025, 12, 25, 9, 30, 0, 0, time.UTC)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_Workdays(t *testing.T) {
	// 工作日执行 (周一到周五)
	schedule, err := ParseCronWithSeconds("0 0 9 * * 1-5")
	assert.NoError(t, err)

	// 从周六开始
	now := time.Date(2025, 12, 27, 10, 0, 0, 0, time.UTC) // Saturday
	next := schedule.Next(now)

	// 应该跳到周一
	expected := time.Date(2025, 12, 29, 9, 0, 0, 0, time.UTC) // Monday
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_WithTimezone(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	assert.NoError(t, err)

	schedule, err := ParseCronWithSeconds("TZ=Asia/Shanghai 0 0 0 * * *")
	assert.NoError(t, err)

	now := time.Date(2025, 12, 25, 10, 0, 0, 0, loc)
	next := schedule.Next(now)

	expected := time.Date(2025, 12, 26, 0, 0, 0, 0, loc)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_NanosecondPrecision(t *testing.T) {
	schedule, err := ParseCronWithSeconds("0 0 0 * * *")
	assert.NoError(t, err)

	now := time.Date(2025, 12, 25, 10, 30, 45, 123456789, time.UTC)
	next := schedule.Next(now)

	// 应该从下一秒开始，纳秒为0
	assert.Equal(t, 0, next.Nanosecond())
}

func TestCronSpecSchedule_Next_YearLimit(t *testing.T) {
	// 一个永远不会匹配的表达式(2月30日)
	schedule := &CronSpecSchedule{
		Second:   1 << 0,
		Minute:   1 << 0,
		Hour:     1 << 0,
		Dom:      1 << 30,
		Month:    1 << 2, // 2月
		Dow:      cronStarBit | cronAllWeekdays,
		Location: time.UTC,
	}

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	next := schedule.Next(now)

	// 超过5年限制，返回零值
	assert.True(t, next.IsZero())
}

func TestCronSpecSchedule_MatchBit(t *testing.T) {
	schedule := &CronSpecSchedule{}

	// 测试位匹配
	bits := uint64(1<<5 | 1<<10 | 1<<15)
	assert.True(t, schedule.matchBit(bits, 5))
	assert.True(t, schedule.matchBit(bits, 10))
	assert.True(t, schedule.matchBit(bits, 15))
	assert.False(t, schedule.matchBit(bits, 7))
}

func TestCronSpecSchedule_MatchDayOfMonthAndWeek(t *testing.T) {
	tests := map[string]struct {
		dom      uint64
		dow      uint64
		date     time.Time
		expected bool
	}{
		"both_match": {
			dom:      1 << 15,
			dow:      1 << 1,                                        // Monday
			date:     time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC), // Monday 15th
			expected: true,
		},
		"dom_star_dow_match": {
			dom:      cronStarBit | cronAllDaysOfMon,
			dow:      1 << 1,                                        // Monday
			date:     time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC), // Monday
			expected: true,
		},
		"dow_star_dom_match": {
			dom:      1 << 15,
			dow:      cronStarBit | cronAllWeekdays,
			date:     time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		"both_star": {
			dom:      cronStarBit | cronAllDaysOfMon,
			dow:      cronStarBit | cronAllWeekdays,
			date:     time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		"neither_match": {
			dom:      1 << 1,
			dow:      1 << 0,                                        // Sunday
			date:     time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC), // Thursday 25th
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule := &CronSpecSchedule{
				Dom: tc.dom,
				Dow: tc.dow,
			}
			result := schedule.matchDayOfMonthAndWeek(tc.date)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCronEverySchedule_Next(t *testing.T) {
	tests := map[string]struct {
		duration time.Duration
		check    func(time.Time, time.Time)
	}{
		"1_hour": {
			duration: 1 * time.Hour,
			check: func(now, next time.Time) {
				assert.Equal(t, now.Add(1*time.Hour), next)
			},
		},
		"30_minutes": {
			duration: 30 * time.Minute,
			check: func(now, next time.Time) {
				assert.Equal(t, now.Add(30*time.Minute), next)
			},
		},
		"10_seconds": {
			duration: 10 * time.Second,
			check: func(now, next time.Time) {
				assert.Equal(t, now.Add(10*time.Second), next)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule := &CronEverySchedule{Duration: tc.duration}
			now := time.Now()
			next := schedule.Next(now)
			tc.check(now, next)
		})
	}
}

func TestCronSpecSchedule_Next_ComplexExpression(t *testing.T) {
	// 每周一、三、五的9:00, 14:00, 18:00执行
	schedule, err := ParseCronWithSeconds("0 0 9,14,18 * * 1,3,5")
	assert.NoError(t, err)

	// 从周一早上8点开始
	now := time.Date(2025, 12, 29, 8, 0, 0, 0, time.UTC) // Monday
	next := schedule.Next(now)

	// 应该是同一天9点
	expected := time.Date(2025, 12, 29, 9, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, next)

	// 继续下一次
	next = schedule.Next(next)
	expected = time.Date(2025, 12, 29, 14, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_Next_LastDayOfMonth(t *testing.T) {
	// 每月最后几天
	schedule, err := ParseCronWithSeconds("0 0 0 28-31 * *")
	assert.NoError(t, err)

	now := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	next := schedule.Next(now)

	// 应该是1月28日
	expected := time.Date(2025, 1, 28, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, next)
}

func TestCronSpecSchedule_WithMethods(t *testing.T) {
	schedule := NewZeroCronSpecSchedule(time.UTC)

	// 测试链式调用
	schedule.WithSecond(1 << 30).
		WithMinute(1 << 15).
		WithHour(1 << 9).
		WithDom(1 << 1).
		WithMonth(1 << 6).
		WithDow(1 << 1)

	assert.Equal(t, uint64(1<<30), schedule.Second)
	assert.Equal(t, uint64(1<<15), schedule.Minute)
	assert.Equal(t, uint64(1<<9), schedule.Hour)
	assert.Equal(t, uint64(1<<1), schedule.Dom)
	assert.Equal(t, uint64(1<<6), schedule.Month)
	assert.Equal(t, uint64(1<<1), schedule.Dow)
}
