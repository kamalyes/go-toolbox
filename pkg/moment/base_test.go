/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-03 15:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:00:57
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

func TestCalculateTimeDifference(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected TimeDifference
	}{
		{time.Duration(366 * 24 * time.Hour), TimeDifference{Years: 1, Days: 0, Hours: 0, Minutes: 0, Seconds: 0}}, // 1 year
		{time.Duration(86400 * time.Second), TimeDifference{Years: 0, Days: 1, Hours: 0, Minutes: 0, Seconds: 0}},  // 1 day
		{time.Duration(367 * 24 * time.Hour), TimeDifference{Years: 1, Days: 1, Hours: 0, Minutes: 0, Seconds: 0}}, // 1 year and 1 day
		{time.Duration(3600 * time.Second), TimeDifference{Years: 0, Days: 0, Hours: 1, Minutes: 0, Seconds: 0}},   // 1 hour
		{time.Duration(61 * time.Second), TimeDifference{Years: 0, Days: 0, Hours: 0, Minutes: 1, Seconds: 1}},     // 1 minute and 1 second
	}

	for _, tt := range tests {
		t.Run(tt.duration.String(), func(t *testing.T) {
			result := CalculateTimeDifference(tt.duration)
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
		result := CalculateBirthDate(test.age)
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
			result := SafeParseTimeToUnixNano(tt.timeStr)
			if result != tt.expected {
				t.Errorf("期望: %d, 实际: %d", tt.expected, result)
			}
		})
	}
}

func TestGetCurrentTimeInfo(t *testing.T) {
	date, hour, currentTime := GetCurrentTimeInfo()
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
	timezone := GetServerTimezone()
	if timezone == "" {
		t.Errorf("Expected non-empty timezone, got: %s", timezone)
	}
}

// TestGetTimeOffset 测试获取时间戳偏移
func TestGetTimeOffset(t *testing.T) {
	ts := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC).Unix()
	offset := GetTimeOffset("Asia/Shanghai", ts)
	expectedOffset := 8 * 3600 // 上海时间相对于 UTC 的偏移，单位为秒

	if offset != expectedOffset {
		t.Errorf("Expected offset %d, got: %d", expectedOffset, offset)
	}
}

// TestFormatWithLocation 测试时间戳格式化
func TestFormatWithLocation(t *testing.T) {
	ts := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC).Unix()
	formatted := FormatWithLocation("Asia/Shanghai", ts, defaultDateFormat)

	expected := "2023-10-01 20:00:00" // 上海时间对应的格式化字符串
	if formatted != expected {
		t.Errorf("Expected formatted time %s, got: %s", expected, formatted)
	}
}

// TestParseWithLocation 测试时间字符串解析
func TestParseWithLocation(t *testing.T) {
	timeStr := "2023-10-01 20:00:00"
	ts := ParseWithLocation("Asia/Shanghai", timeStr, defaultDateFormat)

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

	timestamp, err := ConvertStringToTimestamp(dateString, layout, timeZone)

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
		age, err := CalculateAge(test.birthday, test.currentTime)
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
		age, err := CalculateAge(test.birthday, test.currentTime)
		assert.Error(t, err, "对于生日 %s,期望返回错误", test.birthday)
		assert.Equal(t, test.expected, age, "对于生日 %s,期望年龄 %d,但得到 %d", test.birthday, test.expected, age)
	}
}

func TestHour(t *testing.T) {
	now := time.Now()
	hour := Hour(now)
	if hour < 0 || hour > 23 {
		t.Errorf("Hour() returned an invalid hour value: %d", hour)
	}
}

func TestMinute(t *testing.T) {
	now := time.Now()
	minute := Minute(now)
	if minute < 0 || minute > 59 {
		t.Errorf("Minute() returned an invalid minute value: %d", minute)
	}
}

func TestSecond(t *testing.T) {
	now := time.Now()
	second := Second(now)
	if second < 0 || second > 59 {
		t.Errorf("Second() returned an invalid second value: %d", second)
	}
}

func TestCurrentMillisecond(t *testing.T) {
	ms := CurrentMillisecond()
	if ms <= 0 {
		t.Errorf("CurrentMillisecond() returned an invalid millisecond value: %d", ms)
	}
}

func TestCurrentMicrosecond(t *testing.T) {
	micros := CurrentMicrosecond()
	if micros <= 0 {
		t.Errorf("CurrentMicrosecond() returned an invalid microsecond value: %d", micros)
	}
}

func TestCurrentNanosecond(t *testing.T) {
	nanos := CurrentNanosecond()
	if nanos <= 0 {
		t.Errorf("CurrentNanosecond() returned an invalid nanosecond value: %d", nanos)
	}
}

func TestStrtoTime(t *testing.T) {
	testString := "2024-12-25 23:59:59"
	result, err := StrtoTime(testString)
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
	result := CharToCode(layout)
	if result != expectedLayout {
		t.Errorf("CharToCode() returned %s, expected %s", result, expectedLayout)
	}
}

func TestYear(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	year := Year(testTime)
	expectedYear := 2024
	if year != expectedYear {
		t.Errorf("Year() returned %d, expected %d", year, expectedYear)
	}
}

func TestMonth(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	month := Month(testTime)
	expectedMonth := 4
	if month != expectedMonth {
		t.Errorf("Month() returned %d, expected %d", month, expectedMonth)
	}
}

func TestDay(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	day := Day(testTime)
	expectedDay := 2
	if day != expectedDay {
		t.Errorf("Day() returned %d, expected %d", day, expectedDay)
	}
}

func TestYearDay(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	yearDay := YearDay(testTime)
	expectedYearDay := 93
	if yearDay != expectedYearDay {
		t.Errorf("YearDay() returned %d, expected %d", yearDay, expectedYearDay)
	}
}

func TestYearDefault(t *testing.T) {
	year := Year()
	currentYear := time.Now().Year()
	if year != currentYear {
		t.Errorf("YearDefault() returned %d, expected %d", year, currentYear)
	}
}

func TestMonthDefault(t *testing.T) {
	month := Month()
	currentMonth := int(time.Now().Month())
	if month != currentMonth {
		t.Errorf("MonthDefault() returned %d, expected %d", month, currentMonth)
	}
}

func TestDayDefault(t *testing.T) {
	day := Day()
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
			got := DayOfYear(tc.year, tc.n)
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
		got := LastDayOfMonth(tc.year, tc.month)
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
		got := LastWeekdayOfMonth(tc.year, tc.month, tc.target)
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
		got := NextWorkDay(tt.year, time.Month(tt.month), tt.day)
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
		got := NextWeekday(tt.year, time.Month(tt.month), tt.day, tt.target)
		assert.Equal(t, tt.expYear, got.Year(), "Year mismatch for %s", tt.desc)
		assert.Equal(t, tt.expMonth, got.Month(), "Month mismatch for %s", tt.desc)
		assert.Equal(t, tt.expDay, got.Day(), "Day mismatch for %s", tt.desc)
	}
}

func TestHumanDuration(t *testing.T) {
	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		want  string
	}{
		{
			name:  "same time",
			start: time.Date(2023, 6, 18, 10, 0, 0, 0, time.UTC),
			end:   time.Date(2023, 6, 18, 10, 0, 0, 0, time.UTC),
			want:  "0秒",
		},
		{
			name:  "seconds difference",
			start: time.Date(2023, 6, 18, 10, 0, 0, 0, time.UTC),
			end:   time.Date(2023, 6, 18, 10, 0, 5, 0, time.UTC),
			want:  "5秒",
		},
		{
			name:  "minutes and seconds",
			start: time.Date(2023, 6, 18, 10, 0, 0, 0, time.UTC),
			end:   time.Date(2023, 6, 18, 10, 3, 10, 0, time.UTC),
			want:  "3分10秒",
		},
		{
			name:  "hours, minutes, seconds",
			start: time.Date(2023, 6, 18, 8, 25, 0, 0, time.UTC),
			end:   time.Date(2023, 6, 18, 10, 30, 15, 0, time.UTC),
			want:  "2小时5分15秒",
		},
		{
			name:  "days, hours, minutes",
			start: time.Date(2023, 6, 15, 5, 10, 0, 0, time.UTC),
			end:   time.Date(2023, 6, 18, 10, 30, 0, 0, time.UTC),
			want:  "3天5小时20分",
		},
		{
			name:  "months, days, hours",
			start: time.Date(2023, 4, 1, 7, 0, 0, 0, time.UTC),
			end:   time.Date(2023, 5, 16, 10, 0, 0, 0, time.UTC),
			want:  "1个月15天3小时",
		},
		{
			name:  "years, months, days, hours",
			start: time.Date(2022, 1, 10, 12, 0, 0, 0, time.UTC),
			end:   time.Date(2023, 8, 11, 10, 0, 0, 0, time.UTC),
			want:  "1年7个月零22小时",
		},
		{
			name:  "reverse order (end before start)",
			start: time.Date(2023, 6, 18, 10, 0, 0, 0, time.UTC),
			end:   time.Date(2023, 6, 15, 5, 10, 0, 0, time.UTC),
			want:  "3天4小时50分",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HumanDuration(tt.start, tt.end)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestFormatCompact 测试通用紧凑格式化函数
func TestFormatCompact(t *testing.T) {
	testTime := time.Date(2024, 2, 14, 15, 30, 45, 123456789, time.UTC)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{"Compact Date", CompactDateFormat, "20240214"},
		{"Compact Date Hour", CompactDateHourFormat, "2024021415"},
		{"Compact Date Time", CompactDateTimeFormat, "202402141530"},
		{"Compact Date Time Sec", CompactDateTimeSecFormat, "20240214153045"},
		{"Compact Date Time Milli", CompactDateTimeMilliFormat, "20240214153045.123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCompact(testTime, tt.format)
			assert.Equal(t, tt.expected, result, "FormatCompact with %s should match", tt.name)
		})
	}
}

// TestParseCompact 测试通用紧凑格式解析函数
func TestParseCompact(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		format      string
		expectYear  int
		expectMonth time.Month
		expectDay   int
		expectHour  int
		expectMin   int
		expectSec   int
		shouldError bool
	}{
		{
			name:        "Parse Compact Date",
			value:       "20240214",
			format:      CompactDateFormat,
			expectYear:  2024,
			expectMonth: time.February,
			expectDay:   14,
			expectHour:  0,
			expectMin:   0,
			expectSec:   0,
			shouldError: false,
		},
		{
			name:        "Parse Compact Date Hour",
			value:       "2024021415",
			format:      CompactDateHourFormat,
			expectYear:  2024,
			expectMonth: time.February,
			expectDay:   14,
			expectHour:  15,
			expectMin:   0,
			expectSec:   0,
			shouldError: false,
		},
		{
			name:        "Parse Compact Date Time",
			value:       "202402141530",
			format:      CompactDateTimeFormat,
			expectYear:  2024,
			expectMonth: time.February,
			expectDay:   14,
			expectHour:  15,
			expectMin:   30,
			expectSec:   0,
			shouldError: false,
		},
		{
			name:        "Parse Compact Date Time Sec",
			value:       "20240214153045",
			format:      CompactDateTimeSecFormat,
			expectYear:  2024,
			expectMonth: time.February,
			expectDay:   14,
			expectHour:  15,
			expectMin:   30,
			expectSec:   45,
			shouldError: false,
		},
		{
			name:        "Parse Invalid Format",
			value:       "invalid",
			format:      CompactDateFormat,
			shouldError: true,
		},
		{
			name:        "Parse Wrong Length",
			value:       "2024",
			format:      CompactDateFormat,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseCompact(tt.value, tt.format)

			if tt.shouldError {
				assert.Error(t, err, "Should return error for %s", tt.name)
			} else {
				assert.NoError(t, err, "Should not return error for %s", tt.name)
				assert.Equal(t, tt.expectYear, result.Year(), "Year should match")
				assert.Equal(t, tt.expectMonth, result.Month(), "Month should match")
				assert.Equal(t, tt.expectDay, result.Day(), "Day should match")
				assert.Equal(t, tt.expectHour, result.Hour(), "Hour should match")
				assert.Equal(t, tt.expectMin, result.Minute(), "Minute should match")
				assert.Equal(t, tt.expectSec, result.Second(), "Second should match")
			}
		})
	}
}

// TestFormatCompactDate 测试紧凑日期格式化
func TestFormatCompactDate(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Normal Date",
			time:     time.Date(2024, 2, 14, 15, 30, 45, 0, time.UTC),
			expected: "20240214",
		},
		{
			name:     "First Day of Year",
			time:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "20240101",
		},
		{
			name:     "Last Day of Year",
			time:     time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: "20241231",
		},
		{
			name:     "Leap Year Feb 29",
			time:     time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			expected: "20240229",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCompactDate(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatCompactDateHour 测试紧凑日期+小时格式化
func TestFormatCompactDateHour(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Afternoon Hour",
			time:     time.Date(2024, 2, 14, 15, 30, 45, 0, time.UTC),
			expected: "2024021415",
		},
		{
			name:     "Midnight",
			time:     time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC),
			expected: "2024021400",
		},
		{
			name:     "Last Hour of Day",
			time:     time.Date(2024, 2, 14, 23, 59, 59, 0, time.UTC),
			expected: "2024021423",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCompactDateHour(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatCompactDateTime 测试紧凑日期时间格式化
func TestFormatCompactDateTime(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Normal DateTime",
			time:     time.Date(2024, 2, 14, 15, 30, 45, 0, time.UTC),
			expected: "202402141530",
		},
		{
			name:     "Midnight",
			time:     time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC),
			expected: "202402140000",
		},
		{
			name:     "Last Minute of Day",
			time:     time.Date(2024, 2, 14, 23, 59, 0, 0, time.UTC),
			expected: "202402142359",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCompactDateTime(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatCompactDateTimeSec 测试紧凑日期时间秒格式化
func TestFormatCompactDateTimeSec(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Normal DateTime with Seconds",
			time:     time.Date(2024, 2, 14, 15, 30, 45, 0, time.UTC),
			expected: "20240214153045",
		},
		{
			name:     "Midnight",
			time:     time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC),
			expected: "20240214000000",
		},
		{
			name:     "Last Second of Day",
			time:     time.Date(2024, 2, 14, 23, 59, 59, 0, time.UTC),
			expected: "20240214235959",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCompactDateTimeSec(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatCompactDateTimeMilli 测试紧凑日期时间毫秒格式化
func TestFormatCompactDateTimeMilli(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "With Milliseconds",
			time:     time.Date(2024, 2, 14, 15, 30, 45, 123000000, time.UTC),
			expected: "20240214153045.123",
		},
		{
			name:     "Zero Milliseconds",
			time:     time.Date(2024, 2, 14, 15, 30, 45, 0, time.UTC),
			expected: "20240214153045.000",
		},
		{
			name:     "Max Milliseconds",
			time:     time.Date(2024, 2, 14, 15, 30, 45, 999000000, time.UTC),
			expected: "20240214153045.999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCompactDateTimeMilli(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNowCompactDate 测试获取当前紧凑日期
func TestNowCompactDate(t *testing.T) {
	result := NowCompactDate()
	assert.Len(t, result, 8, "Compact date should be 8 characters (YYYYMMDD)")
	assert.Regexp(t, `^\d{8}$`, result, "Should match YYYYMMDD pattern")

	// 验证可以解析回时间
	parsed, err := ParseCompact(result, CompactDateFormat)
	assert.NoError(t, err, "Should be able to parse back")
	assert.Equal(t, time.Now().Year(), parsed.Year(), "Year should match current year")
	assert.Equal(t, time.Now().Month(), parsed.Month(), "Month should match current month")
	assert.Equal(t, time.Now().Day(), parsed.Day(), "Day should match current day")
}

// TestNowCompactDateHour 测试获取当前紧凑日期+小时
func TestNowCompactDateHour(t *testing.T) {
	result := NowCompactDateHour()
	assert.Len(t, result, 10, "Compact date hour should be 10 characters (YYYYMMDDHH)")
	assert.Regexp(t, `^\d{10}$`, result, "Should match YYYYMMDDHH pattern")

	// 验证可以解析回时间
	parsed, err := ParseCompact(result, CompactDateHourFormat)
	assert.NoError(t, err, "Should be able to parse back")
	assert.Equal(t, time.Now().Hour(), parsed.Hour(), "Hour should match current hour")
}

// TestNowCompactDateTime 测试获取当前紧凑日期时间
func TestNowCompactDateTime(t *testing.T) {
	result := NowCompactDateTime()
	assert.Len(t, result, 12, "Compact date time should be 12 characters (YYYYMMDDHHMM)")
	assert.Regexp(t, `^\d{12}$`, result, "Should match YYYYMMDDHHMM pattern")

	// 验证可以解析回时间
	parsed, err := ParseCompact(result, CompactDateTimeFormat)
	assert.NoError(t, err, "Should be able to parse back")
	now := time.Now()
	assert.Equal(t, now.Year(), parsed.Year(), "Year should match")
	assert.Equal(t, now.Month(), parsed.Month(), "Month should match")
	assert.Equal(t, now.Day(), parsed.Day(), "Day should match")
	assert.Equal(t, now.Hour(), parsed.Hour(), "Hour should match")
	assert.Equal(t, now.Minute(), parsed.Minute(), "Minute should match")
}

// TestNowCompactDateTimeSec 测试获取当前紧凑日期时间秒
func TestNowCompactDateTimeSec(t *testing.T) {
	before := time.Now()
	result := NowCompactDateTimeSec()
	after := time.Now()

	assert.Len(t, result, 14, "Compact date time sec should be 14 characters (YYYYMMDDHHMMSS)")
	assert.Regexp(t, `^\d{14}$`, result, "Should match YYYYMMDDHHMMSS pattern")

	// 验证可以解析回时间 - 使用 ParseInLocation 因为格式化时使用的是本地时间
	parsedLocal, err := time.ParseInLocation(CompactDateTimeSecFormat, result, time.Local)
	assert.NoError(t, err, "Should parse in local timezone")

	// 验证解析的时间在 before 和 after 之间（允许1秒误差）
	assert.True(t,
		parsedLocal.Unix() >= before.Unix()-1 && parsedLocal.Unix() <= after.Unix()+1,
		"Parsed time should be between before (%v) and after (%v), got %v",
		before, after, parsedLocal)
}

// TestNowCompactDateTimeMilli 测试获取当前紧凑日期时间毫秒
func TestNowCompactDateTimeMilli(t *testing.T) {
	result := NowCompactDateTimeMilli()
	assert.Len(t, result, 18, "Compact date time milli should be 18 characters (YYYYMMDDHHMMSS.sss)")
	assert.Regexp(t, `^\d{14}\.\d{3}$`, result, "Should match YYYYMMDDHHMMSS.sss pattern")
}

// TestCompactFormatRoundTrip 测试格式化和解析的往返转换
func TestCompactFormatRoundTrip(t *testing.T) {
	originalTime := time.Date(2024, 2, 14, 15, 30, 45, 123000000, time.UTC)

	tests := []struct {
		name       string
		format     string
		formatFunc func(time.Time) string
	}{
		{"Date", CompactDateFormat, FormatCompactDate},
		{"Date Hour", CompactDateHourFormat, FormatCompactDateHour},
		{"Date Time", CompactDateTimeFormat, FormatCompactDateTime},
		{"Date Time Sec", CompactDateTimeSecFormat, FormatCompactDateTimeSec},
		{"Date Time Milli", CompactDateTimeMilliFormat, FormatCompactDateTimeMilli},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 格式化
			formatted := tt.formatFunc(originalTime)

			// 解析回来
			parsed, err := ParseCompact(formatted, tt.format)
			assert.NoError(t, err, "Should parse successfully")

			// 根据格式精度验证
			assert.Equal(t, originalTime.Year(), parsed.Year(), "Year should match")
			assert.Equal(t, originalTime.Month(), parsed.Month(), "Month should match")
			assert.Equal(t, originalTime.Day(), parsed.Day(), "Day should match")

			// 根据格式检查时分秒
			if tt.format != CompactDateFormat {
				assert.Equal(t, originalTime.Hour(), parsed.Hour(), "Hour should match")
			}
			if tt.format == CompactDateTimeFormat || tt.format == CompactDateTimeSecFormat || tt.format == CompactDateTimeMilliFormat {
				assert.Equal(t, originalTime.Minute(), parsed.Minute(), "Minute should match")
			}
			if tt.format == CompactDateTimeSecFormat || tt.format == CompactDateTimeMilliFormat {
				assert.Equal(t, originalTime.Second(), parsed.Second(), "Second should match")
			}
		})
	}
}

// BenchmarkFormatCompact 性能基准测试
func BenchmarkFormatCompact(b *testing.B) {
	testTime := time.Now()

	benchmarks := []struct {
		name   string
		format string
	}{
		{"Date", CompactDateFormat},
		{"DateHour", CompactDateHourFormat},
		{"DateTime", CompactDateTimeFormat},
		{"DateTimeSec", CompactDateTimeSecFormat},
		{"DateTimeMilli", CompactDateTimeMilliFormat},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = FormatCompact(testTime, bm.format)
			}
		})
	}
}

// BenchmarkParseCompact 解析性能基准测试
func BenchmarkParseCompact(b *testing.B) {
	benchmarks := []struct {
		name   string
		value  string
		format string
	}{
		{"Date", "20240214", CompactDateFormat},
		{"DateHour", "2024021415", CompactDateHourFormat},
		{"DateTime", "202402141530", CompactDateTimeFormat},
		{"DateTimeSec", "20240214153045", CompactDateTimeSecFormat},
		{"DateTimeMilli", "20240214153045.123", CompactDateTimeMilliFormat},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ParseCompact(bm.value, bm.format)
			}
		})
	}
}

// TestCompactFormatEdgeCases 测试边界情况
func TestCompactFormatEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Leap Year Feb 29",
			time:     time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: "20240229",
		},
		{
			name:     "Year 2000",
			time:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "20000101",
		},
		{
			name:     "Year 1900",
			time:     time.Date(1900, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: "19001231",
		},
		{
			name:     "Future Year 2100",
			time:     time.Date(2100, 6, 15, 12, 30, 45, 0, time.UTC),
			expected: "21000615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCompactDate(tt.time)
			assert.Equal(t, tt.expected, result)

			// 验证可以解析回来
			parsed, err := ParseCompact(result, CompactDateFormat)
			assert.NoError(t, err)
			assert.Equal(t, tt.time.Year(), parsed.Year())
			assert.Equal(t, tt.time.Month(), parsed.Month())
			assert.Equal(t, tt.time.Day(), parsed.Day())
		})
	}
}
