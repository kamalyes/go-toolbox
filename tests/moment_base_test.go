/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-13 15:28:36
 * @FilePath: \go-toolbox\tests\moment_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/moment"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTimeDifference(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected moment.TimeDifference
	}{
		{time.Duration(31536000 * time.Second), moment.TimeDifference{Years: 1, Days: 0, Hours: 0, Minutes: 0, Seconds: 0}}, // 1 year
		{time.Duration(86400 * time.Second), moment.TimeDifference{Years: 0, Days: 1, Hours: 0, Minutes: 0, Seconds: 0}},    // 1 day
		{time.Duration(366 * 24 * time.Hour), moment.TimeDifference{Years: 1, Days: 1, Hours: 0, Minutes: 0, Seconds: 0}},   // 1 year and 1 day
		{time.Duration(3600 * time.Second), moment.TimeDifference{Years: 0, Days: 0, Hours: 1, Minutes: 0, Seconds: 0}},     // 1 hour
		{time.Duration(61 * time.Second), moment.TimeDifference{Years: 0, Days: 0, Hours: 0, Minutes: 1, Seconds: 1}},       // 1 minute and 1 second
	}

	for _, tt := range tests {
		t.Run(tt.duration.String(), func(t *testing.T) {
			result := moment.CalculateTimeDifference(tt.duration)
			if result != tt.expected {
				t.Errorf("期望: %+v, 实际: %+v", tt.expected, result)
			}
		})
	}
}

func TestSafeTimeToUnixNano(t *testing.T) {
	tests := []struct {
		timeStr  string
		expected int64
	}{
		{"2023-10-01T12:00:00Z", 1696161600000}, // 以毫秒表示的 Unix 时间戳
		{"2023-10-01T12:00:00+00:00", 1696161600000},
		{"2023-10-01T12:00:00+08:00", 1696132800000}, // 注意时区差异
		{"invalid-time", 0},                          // 无效时间
	}

	for _, tt := range tests {
		t.Run(tt.timeStr, func(t *testing.T) {
			result := moment.SafeParseTimeToUnixNano(tt.timeStr)
			if result != tt.expected {
				t.Errorf("期望: %d, 实际: %d", tt.expected, result)
			}
		})
	}
}

func TestGetCurrentTimeInfo(t *testing.T) {
	date, hour, currentTime := moment.GetCurrentTimeInfo()
	if date == "" {
		t.Error("当前日期应该不为空")
	}

	if hour < 0 || hour > 23 {
		t.Errorf("小时应在 0 到 23 之间，但实际: %d", hour)
	}

	if currentTime.IsZero() {
		t.Error("当前时间应该不为零")
	}
}

const (
	defaultDateFormat = "2006-01-02 15:04:05"
)

// TestGetServerTimezone 测试获取服务器时区
func TestGetServerTimezone(t *testing.T) {
	timezone := moment.GetServerTimezone()
	if timezone == "" {
		t.Errorf("Expected non-empty timezone, got: %s", timezone)
	}
}

// TestGetTimeOffset 测试获取时间戳偏移
func TestGetTimeOffset(t *testing.T) {
	ts := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC).Unix()
	offset := moment.GetTimeOffset("Asia/Shanghai", ts)
	expectedOffset := 8 * 3600 // 上海时间相对于 UTC 的偏移，单位为秒

	if offset != expectedOffset {
		t.Errorf("Expected offset %d, got: %d", expectedOffset, offset)
	}
}

// TestFormatWithLocation 测试时间戳格式化
func TestFormatWithLocation(t *testing.T) {
	ts := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC).Unix()
	formatted := moment.FormatWithLocation("Asia/Shanghai", ts, defaultDateFormat)

	expected := "2023-10-01 20:00:00" // 上海时间对应的格式化字符串
	if formatted != expected {
		t.Errorf("Expected formatted time %s, got: %s", expected, formatted)
	}
}

// TestParseWithLocation 测试时间字符串解析
func TestParseWithLocation(t *testing.T) {
	timeStr := "2023-10-01 20:00:00"
	ts := moment.ParseWithLocation("Asia/Shanghai", timeStr, defaultDateFormat)

	expectedTS := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC).Unix() // 2023-10-01 12:00:00 UTC
	if ts != expectedTS {
		t.Errorf("Expected timestamp %d, got: %d", expectedTS, ts)
	}
}

func TestConvertStringToTimestamp(t *testing.T) {
	expectedTimestamp := int64(1628424042)
	dateString := "2021-08-08 12:03:42"
	layout := "2006-01-02 15:05:05"
	timeZone := "UTC"

	timestamp, err := moment.ConvertStringToTimestamp(dateString, layout, timeZone)

	assert.NoError(t, err)
	assert.Equal(t, expectedTimestamp, timestamp, "Timestamps should match")
}

// 测试计算年龄的函数
func TestCalculateAge(t *testing.T) {
	tests := []struct {
		birthday    string
		currentTime time.Time
		expected    int
	}{
		{"1990-05-15", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 33}, // 生日当天
		{"2000-01-01", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 23}, // 生日已过
		{"1985-12-31", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 37}, // 生日未到
		{"2020-01-01", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 3},  // 生日未到
		{"2000-02-29", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 23}, // 闰年出生
	}

	for _, test := range tests {
		age, err := moment.CalculateAge(test.birthday, test.currentTime)
		assert.NoError(t, err, "计算 %s 的年龄时出错", test.birthday)
		assert.Equal(t, test.expected, age, "对于生日 %s,期望年龄 %d,但得到 %d", test.birthday, test.expected, age)
	}
}

// 测试异常用例
func TestCalculateAgeErrors(t *testing.T) {
	// 异常用例
	invalidTests := []struct {
		birthday    string
		currentTime time.Time
		expected    int
	}{
		{"invalid-date", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 0}, // 无效的日期格式
		{"2023-02-30", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 0},   // 不存在的日期
		{"2023-13-01", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 0},   // 无效的月份
		{"2023-00-01", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 0},   // 无效的月份
		{"2023-01-32", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 0},   // 不存在的日期
	}

	for _, test := range invalidTests {
		age, err := moment.CalculateAge(test.birthday, test.currentTime)
		assert.Error(t, err, "对于生日 %s,期望返回错误", test.birthday)
		assert.Equal(t, test.expected, age, "对于生日 %s,期望年龄 %d,但得到 %d", test.birthday, test.expected, age)
	}
}

func TestHour(t *testing.T) {
	now := time.Now()
	hour := moment.Hour(now)
	if hour < 0 || hour > 23 {
		t.Errorf("Hour() returned an invalid hour value: %d", hour)
	}
}

func TestMinute(t *testing.T) {
	now := time.Now()
	minute := moment.Minute(now)
	if minute < 0 || minute > 59 {
		t.Errorf("Minute() returned an invalid minute value: %d", minute)
	}
}

func TestSecond(t *testing.T) {
	now := time.Now()
	second := moment.Second(now)
	if second < 0 || second > 59 {
		t.Errorf("Second() returned an invalid second value: %d", second)
	}
}

func TestCurrentMillisecond(t *testing.T) {
	ms := moment.CurrentMillisecond()
	if ms <= 0 {
		t.Errorf("CurrentMillisecond() returned an invalid millisecond value: %d", ms)
	}
}

func TestCurrentMicrosecond(t *testing.T) {
	micros := moment.CurrentMicrosecond()
	if micros <= 0 {
		t.Errorf("CurrentMicrosecond() returned an invalid microsecond value: %d", micros)
	}
}

func TestCurrentNanosecond(t *testing.T) {
	nanos := moment.CurrentNanosecond()
	if nanos <= 0 {
		t.Errorf("CurrentNanosecond() returned an invalid nanosecond value: %d", nanos)
	}
}

func TestStrtoTime(t *testing.T) {
	testString := "2024-12-25 23:59:59"
	result, err := moment.StrtoTime(testString)
	if err != nil {
		t.Errorf("StrtoTime() error: %v", err)
	}

	expectedDate := time.Date(2024, time.December, 25, 23, 59, 59, 0, time.Local)
	if !result.Equal(expectedDate) {
		t.Errorf("Strtotime() did not parse the string correctly")
	}
}

func TestCharToCode(t *testing.T) {
	layout := "Y-m-d H:i:s"
	expectedLayout := "2006-1-2 15:4:5"
	result := moment.CharToCode(layout)
	if result != expectedLayout {
		t.Errorf("CharToCode() returned %s, expected %s", result, expectedLayout)
	}
}

func TestYear(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	year := moment.Year(testTime)
	expectedYear := 2024
	if year != expectedYear {
		t.Errorf("Year() returned %d, expected %d", year, expectedYear)
	}
}

func TestMonth(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	month := moment.Month(testTime)
	expectedMonth := 4
	if month != expectedMonth {
		t.Errorf("Month() returned %d, expected %d", month, expectedMonth)
	}
}

func TestDay(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	day := moment.Day(testTime)
	expectedDay := 2
	if day != expectedDay {
		t.Errorf("Day() returned %d, expected %d", day, expectedDay)
	}
}

func TestYearDay(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	yearDay := moment.YearDay(testTime)
	expectedYearDay := 93
	if yearDay != expectedYearDay {
		t.Errorf("YearDay() returned %d, expected %d", yearDay, expectedYearDay)
	}
}

func TestYearDefault(t *testing.T) {
	year := moment.Year()
	currentYear := time.Now().Year()
	if year != currentYear {
		t.Errorf("YearDefault() returned %d, expected %d", year, currentYear)
	}
}

func TestMonthDefault(t *testing.T) {
	month := moment.Month()
	currentMonth := int(time.Now().Month())
	if month != currentMonth {
		t.Errorf("MonthDefault() returned %d, expected %d", month, currentMonth)
	}
}

func TestDayDefault(t *testing.T) {
	day := moment.Day()
	currentDay := int(time.Now().Day())
	if day != currentDay {
		t.Errorf("DayDefault() returned %d, expected %d", day, currentDay)
	}
}
