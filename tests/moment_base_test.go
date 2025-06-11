/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-11 18:15:26
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

// TestCalculateBirthDate 测试 CalculateBirthDate 函数
func TestCalculateBirthDate(t *testing.T) {
	tests := []struct {
		age          int
		expectedYear int
	}{
		{0, time.Now().Year()},       // 0岁，应该是今年
		{1, time.Now().Year() - 1},   // 1岁，应该是去年
		{25, time.Now().Year() - 25}, // 25岁
	}

	for _, test := range tests {
		result := moment.CalculateBirthDate(test.age)
		assert.Equal(t, test.expectedYear, result.Year(), "年龄为 %d 的出生年份应该是 %d", test.age, test.expectedYear)
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

func TestDayOfYear(t *testing.T) {
	tests := map[string]struct {
		year int
		n    int
		want time.Time
	}{
		"2025-1":   {2025, 1, time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)},
		"2025-200": {2025, 162, time.Date(2025, 6, 11, 0, 0, 0, 0, time.Local)},
		"2028-60":  {2028, 60, time.Date(2028, 2, 29, 0, 0, 0, 0, time.Local)},
		"2028-366": {2028, 366, time.Date(2028, 12, 31, 0, 0, 0, 0, time.Local)},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := moment.DayOfYear(tc.year, tc.n)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestLastDayOfMonth(t *testing.T) {
	// 测试用例：year, month, 期望的最后一天日期字符串
	testCases := []struct {
		year     int
		month    time.Month
		wantDate string
	}{
		{2020, time.February, "2020-02-29"},  // 2020年是闰年，2月有29天
		{2023, time.December, "2023-12-31"},  // 2023年12月最后一天
		{2026, time.February, "2026-02-28"},  // 2026年不是闰年，2月28天
		{2025, time.June, "2025-06-30"},      // 6月最后一天
		{2025, time.July, "2025-07-31"},      // 7月最后一天
		{2025, time.August, "2025-08-31"},    // 8月最后一天
		{2025, time.September, "2025-09-30"}, // 9月最后一天
		{2025, time.October, "2025-10-31"},   // 10月最后一天
		{2025, time.December, "2025-12-31"},  // 12月最后一天
	}

	for _, tc := range testCases {
		got := moment.LastDayOfMonth(tc.year, tc.month)
		assert.Equal(t, tc.wantDate, got.Format("2006-01-02"), "LastDayOfMonth(%d, %d)", tc.year, tc.month)
	}
}

func TestLastWeekdayOfMonth(t *testing.T) {
	// 测试用例：year, month, 目标weekday, 期望返回日期字符串
	testCases := []struct {
		year     int
		month    time.Month
		target   time.Weekday
		wantDate string
	}{
		{2023, time.October, time.Friday, "2023-10-27"},     // 2023年10月最后一个周五是27号
		{2023, time.October, time.Monday, "2023-10-30"},     // 2023年10月最后一个周一是30号
		{2023, time.January, time.Sunday, "2023-01-29"},     // 2023年1月最后一个周日是29号
		{2020, time.February, time.Saturday, "2020-02-29"},  // 2020年闰年2月最后一天是周六29号
		{2025, time.June, time.Thursday, "2025-06-26"},      // 6月最后一个周四是26号
		{2025, time.June, time.Monday, "2025-06-30"},        // 6月最后一个周一是30号
		{2025, time.July, time.Friday, "2025-07-25"},        // 7月最后一个周五是25号
		{2025, time.July, time.Wednesday, "2025-07-30"},     // 7月最后一个周三是30号
		{2025, time.August, time.Sunday, "2025-08-31"},      // 8月最后一个周日是31号
		{2025, time.August, time.Tuesday, "2025-08-26"},     // 8月最后一个周二是26号
		{2025, time.September, time.Monday, "2025-09-29"},   // 9月最后一个周一是29号
		{2025, time.September, time.Saturday, "2025-09-27"}, // 9月最后一个周六是27号
		{2025, time.October, time.Friday, "2025-10-31"},     // 10月最后一个周五是31号
		{2025, time.October, time.Monday, "2025-10-27"},     // 10月最后一个周一是27号
		{2025, time.December, time.Thursday, "2025-12-25"},  // 12月最后一个周四是25号
		{2025, time.December, time.Tuesday, "2025-12-30"},   // 12月最后一个周二是30号
	}

	for _, tc := range testCases {
		got := moment.LastWeekdayOfMonth(tc.year, tc.month, tc.target)
		assert.Equal(t, tc.wantDate, got.Format("2006-01-02"),
			"LastWeekdayOfMonth(%d, %d, %s)", tc.year, tc.month, tc.target.String())
	}
}

func TestNextWorkDay(t *testing.T) {
	tests := []struct {
		year, month, day int
		want             int    // 期望返回的“最近工作日”的日（1~31）
		desc             string // 用例说明
	}{
		// 2023-06-04 是周日，返回下周一 2023-06-05
		{2023, int(time.June), 4, 5, "Sunday -> next Monday (+1 day)"},

		// 2023-06-05 是周一，返回周二 2023-06-06
		{2023, int(time.June), 5, 6, "Monday -> next day (Tuesday)"},

		// 2023-06-08 是周四，返回周五 2023-06-09
		{2023, int(time.June), 8, 9, "Thursday -> next day (Friday)"},

		// 2023-06-09 是周五，返回下周一 2023-06-12 (+3 days)
		{2023, int(time.June), 9, 12, "Friday -> next Monday (+3 days)"},

		// 2023-06-10 是周六，返回下周一 2023-06-12 (+2 days)
		{2023, int(time.June), 10, 12, "Saturday -> next Monday (+2 days)"},

		// 跨月测试：2023-12-29 是周五，返回下周一 2024-01-01
		{2023, int(time.December), 29, 1, "Friday (end of year) -> next Monday next year (Jan 1)"},

		// 跨年测试：2023-12-31 是周日，返回下周一 2024-01-01
		{2023, int(time.December), 31, 1, "Sunday (end of year) -> next Monday next year (Jan 1)"},
	}

	for _, tt := range tests {
		got := moment.NextWorkDay(tt.year, time.Month(tt.month), tt.day)
		assert.Equalf(t, tt.want, got.Day(), "NextWorkDay(%d, %d, %d) failed: %s", tt.year, tt.month, tt.day, tt.desc)
	}
}

func TestNextWeekday(t *testing.T) {
	// 基准日期：2025-06-11 周三（实际2025年6月11日是周三）
	// 这里保持和之前一致，覆盖周一到周日，跨月跨年

	type testCase struct {
		year, month, day int
		target           time.Weekday
		expYear          int
		expMonth         time.Month
		expDay           int
		desc             string
	}

	tests := []testCase{
		// 基准日期场景（当天不是目标）
		{2025, 6, 11, time.Monday, 2025, 6, 16, "Mon after Wed 2025-06-11"},
		{2025, 6, 11, time.Tuesday, 2025, 6, 17, "Tue after Wed 2025-06-11"},
		{2025, 6, 11, time.Wednesday, 2025, 6, 18, "Wed after Wed 2025-06-11 (skip same day)"},
		{2025, 6, 11, time.Thursday, 2025, 6, 12, "Thu after Wed 2025-06-11"},
		{2025, 6, 11, time.Friday, 2025, 6, 13, "Fri after Wed 2025-06-11"},
		{2025, 6, 11, time.Saturday, 2025, 6, 14, "Sat after Wed 2025-06-11"},
		{2025, 6, 11, time.Sunday, 2025, 6, 15, "Sun after Wed 2025-06-11"},

		// 当天是目标，跳到下一周
		{2025, 6, 16, time.Monday, 2025, 6, 23, "Mon on Mon 2025-06-16 (next week)"},
		{2025, 6, 17, time.Tuesday, 2025, 6, 24, "Tue on Tue 2025-06-17 (next week)"},
		{2025, 6, 18, time.Wednesday, 2025, 6, 25, "Wed on Wed 2025-06-18 (next week)"},
		{2025, 6, 19, time.Thursday, 2025, 6, 26, "Thu on Thu 2025-06-19 (next week)"},
		{2025, 6, 20, time.Friday, 2025, 6, 27, "Fri on Fri 2025-06-20 (next week)"},
		{2025, 6, 21, time.Saturday, 2025, 6, 28, "Sat on Sat 2025-06-21 (next week)"},
		{2025, 6, 22, time.Sunday, 2025, 6, 29, "Sun on Sun 2025-06-22 (next week)"},

		// 跨月场景：2025-06-30是周一，找下一个周一、周五、周日等
		{2025, 6, 30, time.Monday, 2025, 7, 7, "Mon on Mon 2025-06-30 (next week)"},
		{2025, 6, 30, time.Friday, 2025, 7, 4, "Fri after Mon 2025-06-30"},
		{2025, 6, 30, time.Sunday, 2025, 7, 6, "Sun after Mon 2025-06-30"},

		// 跨年场景：2025-12-31是周三，找下一个周三、周四
		{2025, 12, 31, time.Wednesday, 2026, 1, 7, "Wed on Wed 2025-12-31 (next week)"},
		{2025, 12, 31, time.Thursday, 2026, 1, 1, "Thu after Wed 2025-12-31"},
	}

	for _, tt := range tests {
		got := moment.NextWeekday(tt.year, time.Month(tt.month), tt.day, tt.target)
		assert.Equal(t, tt.expYear, got.Year(), "Year mismatch for %s", tt.desc)
		assert.Equal(t, tt.expMonth, got.Month(), "Month mismatch for %s", tt.desc)
		assert.Equal(t, tt.expDay, got.Day(), "Day mismatch for %s", tt.desc)
	}
}
