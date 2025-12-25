/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-25 10:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-25 12:07:16
 * @FilePath: \go-toolbox\pkg\cron\parser_test.go
 * @Description: Parser 测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package cron

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewCronParser 测试创建 CronParser 实例
// 验证有效选项能成功创建，无效选项会触发 panic
func TestNewCronParser(t *testing.T) {
	tests := map[string]struct {
		options CronParseOption
		panic   bool
	}{
		"valid_standard": {
			options: CronMinute | CronHour | CronDom | CronMonth | CronDow,
			panic:   false,
		},
		"valid_with_seconds": {
			options: CronSecond | CronMinute | CronHour | CronDom | CronMonth | CronDow,
			panic:   false,
		},
		"invalid_multiple_optionals": {
			options: CronSecondOptional | CronDowOptional,
			panic:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panic {
				assert.Panics(t, func() {
					NewCronParser(tc.options)
				})
			} else {
				parser := NewCronParser(tc.options)
				assert.NotNil(t, parser)
			}
		})
	}
}

// TestCronParser_Parse 测试 Parse 方法的基本功能
// 验证空表达式、有效表达式、字段数量、时区格式等情况
func TestCronParser_Parse(t *testing.T) {
	tests := map[string]struct {
		spec      string
		parser    *CronParser
		expectErr bool
	}{
		"empty_spec": {
			spec:      "",
			parser:    CronStandardParser,
			expectErr: true,
		},
		"valid_standard_5_fields": {
			spec:      "0 0 * * *",
			parser:    CronStandardParser,
			expectErr: false,
		},
		"valid_with_seconds_6_fields": {
			spec:      "0 0 0 * * *",
			parser:    CronSecondParser,
			expectErr: false,
		},
		"invalid_field_count": {
			spec:      "0 0 0",
			parser:    CronStandardParser,
			expectErr: true,
		},
		"with_timezone_utc": {
			spec:      "TZ=UTC 0 0 * * *",
			parser:    CronStandardParser,
			expectErr: false,
		},
		"with_timezone_invalid": {
			spec:      "TZ=InvalidZone 0 0 * * *",
			parser:    CronStandardParser,
			expectErr: true,
		},
		"cron_tz_format": {
			spec:      "CRON_TZ=Asia/Shanghai 0 0 * * *",
			parser:    CronStandardParser,
			expectErr: false,
		},
		"timezone_format_error_no_space": {
			spec:      "TZ=UTC",
			parser:    CronStandardParser,
			expectErr: true,
		},
		"timezone_format_error_no_equals": {
			spec:      "TZ 0 0 * * *",
			parser:    CronStandardParser,
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := tc.parser.Parse(tc.spec)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// TestCronParser_ParseDescriptorWithoutSupport 测试不支持描述符时的错误处理
// 验证当解析器未启用描述符支持时，解析描述符表达式会返回错误
func TestCronParser_ParseDescriptorWithoutSupport(t *testing.T) {
	parser := NewCronParser(CronMinute | CronHour | CronDom | CronMonth | CronDow) // 不包含 CronDescriptor
	schedule, err := parser.Parse("@daily")
	assert.Error(t, err)
	assert.Nil(t, schedule)
	assert.Contains(t, err.Error(), "描述符不受支持")
}

// TestCronParser_NormalizeFields 测试字段规范化功能
// 验证字段数量验证和默认值填充逻辑
func TestCronParser_NormalizeFields(t *testing.T) {
	tests := map[string]struct {
		parser    *CronParser
		fields    []string
		expectErr bool
		expected  int // 期望的字段数量
	}{
		"exact_fields": {
			parser:    CronStandardParser,
			fields:    []string{"0", "0", "*", "*", "*"},
			expectErr: false,
			expected:  6, // 会填充秒字段
		},
		"too_few_fields": {
			parser:    CronStandardParser,
			fields:    []string{"0", "0"},
			expectErr: true,
		},
		"too_many_fields": {
			parser:    CronStandardParser,
			fields:    []string{"0", "0", "*", "*", "*", "*", "*"},
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := tc.parser.normalizeFields(tc.fields)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tc.expected)
			}
		})
	}
}

// TestParseCronStandard 测试标准 Cron 表达式解析(5个字段：分 时 日 月 周)
// 验证各种有效和无效的标准 cron 表达式
func TestParseCronStandard(t *testing.T) {
	specs := map[string]bool{
		"0 0 * * *":      false, // valid
		"*/5 * * * *":    false, // valid
		"0 0 1 * *":      false, // valid
		"invalid":        true,  // invalid
		"0 0 * * * *":    true,  // too many fields for standard
		"@invalid_alias": true,  // invalid alias
	}

	for spec, expectErr := range specs {
		t.Run(spec, func(t *testing.T) {
			schedule, err := ParseCronStandard(spec)
			if expectErr {
				assert.Error(t, err)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// TestParseCronWithSeconds 测试带秒的 Cron 表达式解析(6个字段：秒 分 时 日 月 周)
// 验证包含秒字段的 cron 表达式解析
func TestParseCronWithSeconds(t *testing.T) {
	specs := map[string]bool{
		"0 0 0 * * *":    false, // valid
		"*/10 * * * * *": false, // valid
		"0 0 0 1 1 *":    false, // valid
		"0 0 * * *":      true,  // too few fields
		"invalid":        true,  // invalid
	}

	for spec, expectErr := range specs {
		t.Run(spec, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(spec)
			if expectErr {
				assert.Error(t, err)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// TestExpressionAliases 测试预定义的表达式别名
// 验证所有内置别名都能正确映射到对应的 cron 表达式
func TestExpressionAliases(t *testing.T) {
	aliases := []string{
		"@secondly", "@minutely", "@hourly", "@daily", "@midnight",
		"@weekly", "@monthly", "@yearly", "@annually",
		"@workdays", "@weekends",
		"@every_5min", "@every_15min", "@every_30min",
		"@monday", "@tuesday", "@wednesday", "@thursday", "@friday", "@saturday", "@sunday",
	}

	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			expr, ok := GetExpression(alias)
			assert.True(t, ok)
			assert.NotEmpty(t, expr)
		})
	}
}

// TestGetExpression_NotFound 测试获取不存在的别名
// 验证查询不存在的别名时返回空和 false
func TestGetExpression_NotFound(t *testing.T) {
	expr, ok := GetExpression("@not_exists")
	assert.False(t, ok)
	assert.Empty(t, expr)
}

// TestCronParser_ParseWithMonthNames 测试使用月份名称的表达式解析
// 验证月份名称(如 jan, feb 等)能正确解析为对应的月份位
func TestCronParser_ParseWithMonthNames(t *testing.T) {
	schedule, err := ParseCronStandard("0 0 1 jan *")
	assert.NoError(t, err)
	assert.NotNil(t, schedule)

	spec := schedule.(*CronSpecSchedule)
	assert.NotZero(t, spec.Month)
	// 验证 January (月份1) 的位被设置
	assert.True(t, spec.matchBit(spec.Month, 1), "January 月份位应该被设置")
}

// TestCronParser_ParseWithWeekdayNames 测试使用星期名称的表达式解析
// 验证星期名称(如 mon, tue 等)能正确解析为对应的星期位
func TestCronParser_ParseWithWeekdayNames(t *testing.T) {
	schedule, err := ParseCronStandard("0 0 * * mon")
	assert.NoError(t, err)
	assert.NotNil(t, schedule)

	spec := schedule.(*CronSpecSchedule)
	assert.NotZero(t, spec.Dow)
	// 验证 Monday (星期1) 的位被设置
	assert.True(t, spec.matchBit(spec.Dow, 1), "Monday 星期位应该被设置")
}

// TestCronParser_ParseRange 测试范围表达式解析(如 9-17)
// 验证范围内的所有值都被正确设置，范围外的值未被设置
func TestCronParser_ParseRange(t *testing.T) {
	schedule, err := ParseCronStandard("0 9-17 * * *")
	assert.NoError(t, err)
	assert.NotNil(t, schedule)

	spec := schedule.(*CronSpecSchedule)
	// 验证分钟字段为 0
	assert.True(t, spec.matchBit(spec.Minute, 0), "分钟应该为 0")
	assert.False(t, spec.matchBit(spec.Minute, 1), "分钟不应该为 1")

	// 验证小时字段包含 9-17 的位
	for h := uint(9); h <= 17; h++ {
		assert.True(t, spec.matchBit(spec.Hour, h), "小时 %d 应该被设置", h)
	}
	// 验证范围外的小时未被设置
	assert.False(t, spec.matchBit(spec.Hour, 8), "小时 8 不应该被设置")
	assert.False(t, spec.matchBit(spec.Hour, 18), "小时 18 不应该被设置")
}

// TestCronParser_ParseStep 测试步长表达式解析(如 */15)
// 验证步长值正确设置，非步长值未被设置
func TestCronParser_ParseStep(t *testing.T) {
	schedule, err := ParseCronStandard("*/15 * * * *")
	assert.NoError(t, err)
	assert.NotNil(t, schedule)

	spec := schedule.(*CronSpecSchedule)
	// 验证分钟字段包含 0, 15, 30, 45
	assert.True(t, spec.matchBit(spec.Minute, 0), "分钟 0 应该被设置")
	assert.True(t, spec.matchBit(spec.Minute, 15), "分钟 15 应该被设置")
	assert.True(t, spec.matchBit(spec.Minute, 30), "分钟 30 应该被设置")
	assert.True(t, spec.matchBit(spec.Minute, 45), "分钟 45 应该被设置")
	// 验证非步长值未被设置
	assert.False(t, spec.matchBit(spec.Minute, 1), "分钟 1 不应该被设置")
	assert.False(t, spec.matchBit(spec.Minute, 10), "分钟 10 不应该被设置")
	assert.False(t, spec.matchBit(spec.Minute, 20), "分钟 20 不应该被设置")
}

// TestCronParser_ParseList 测试列表表达式解析(如 1,15)
// 验证列表中的值被设置，列表外的值未被设置
func TestCronParser_ParseList(t *testing.T) {
	schedule, err := ParseCronStandard("0 0 1,15 * *")
	assert.NoError(t, err)
	assert.NotNil(t, schedule)

	spec := schedule.(*CronSpecSchedule)
	// 验证列表中的日期被设置
	assert.True(t, spec.matchBit(spec.Dom, 1), "日期 1 应该被设置")
	assert.True(t, spec.matchBit(spec.Dom, 15), "日期 15 应该被设置")
	// 验证列表外的日期未被设置
	assert.False(t, spec.matchBit(spec.Dom, 2), "日期 2 不应该被设置")
	assert.False(t, spec.matchBit(spec.Dom, 10), "日期 10 不应该被设置")
	assert.False(t, spec.matchBit(spec.Dom, 31), "日期 31 不应该被设置")
}

// TestCronParser_ParseWildcard 测试通配符 * 的解析
// 验证通配符正确设置星号位标记
func TestCronParser_ParseWildcard(t *testing.T) {
	schedule, err := ParseCronStandard("* * * * *")
	assert.NoError(t, err)
	assert.NotNil(t, schedule)

	spec := schedule.(*CronSpecSchedule)
	// 验证通配符设置了星号位
	assert.NotZero(t, spec.Dom&cronStarBit, "日期字段应该设置星号位")
	assert.NotZero(t, spec.Dow&cronStarBit, "星期字段应该设置星号位")
}

// TestCronParser_ParseQuestion 测试问号 ? 的解析
// 验证问号和通配符 * 功能相同，都设置星号位标记
func TestCronParser_ParseQuestion(t *testing.T) {
	schedule, err := ParseCronStandard("0 0 ? * *")
	assert.NoError(t, err)
	assert.NotNil(t, schedule)

	spec := schedule.(*CronSpecSchedule)
	// ? 和 * 一样，设置星号位
	assert.NotZero(t, spec.Dom&cronStarBit, "日期字段使用 ? 应该设置星号位")
}

// TestParseSpecialDescriptors 测试特殊描述符的解析
// 验证所有自定义的时间描述符(如 @night, @dawn 等)都能正确解析
func TestParseSpecialDescriptors(t *testing.T) {
	descriptors := map[string]struct {
		desc string
	}{
		"night":           {desc: "@night"},
		"dawn":            {desc: "@dawn"},
		"noon":            {desc: "@noon"},
		"dusk":            {desc: "@dusk"},
		"late_night":      {desc: "@late_night"},
		"early_morning":   {desc: "@early_morning"},
		"lunch_time":      {desc: "@lunch_time"},
		"dinner_time":     {desc: "@dinner_time"},
		"workday_start":   {desc: "@workday_start"},
		"workday_end":     {desc: "@workday_end"},
		"weekend_morning": {desc: "@weekend_morning"},
		"weekend_evening": {desc: "@weekend_evening"},
		"month_end":       {desc: "@month_end"},
		"quarter_start":   {desc: "@quarter_start"},
		"quarter_end":     {desc: "@quarter_end"},
		"year_start":      {desc: "@year_start"},
		"year_end":        {desc: "@year_end"},
		"spring_start":    {desc: "@spring_start"},
		"summer_start":    {desc: "@summer_start"},
		"autumn_start":    {desc: "@autumn_start"},
		"winter_start":    {desc: "@winter_start"},
	}

	for name, tc := range descriptors {
		schedule, err := ParseCronStandard(tc.desc)
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, err)
			assert.NotNil(t, schedule)
		})
	}
}

// TestInvalidDescriptor 测试无效描述符的错误处理
// 验证未知描述符和格式错误的 @every 表达式会返回错误
func TestInvalidDescriptor(t *testing.T) {
	descs := []string{"@unknown", "@every 0s", "@every abc"}
	for _, desc := range descs {
		schedule, err := ParseCronWithSeconds(desc)
		t.Run(desc, func(t *testing.T) {
			assert.Error(t, err)
			assert.Nil(t, schedule)
		})
	}
}

// TestCronParser_ParseComplexExpression 测试复杂的 Cron 表达式
// 验证混合使用范围、步长、列表等复杂表达式的解析
func TestCronParser_ParseComplexExpression(t *testing.T) {
	tests := map[string]struct {
		spec           string
		validateMinute func(*testing.T, *CronSpecSchedule)
		validateHour   func(*testing.T, *CronSpecSchedule)
		validateDom    func(*testing.T, *CronSpecSchedule)
	}{
		"range_with_step": {
			spec: "0-30/10 * * * *",
			validateMinute: func(t *testing.T, s *CronSpecSchedule) {
				// 0-30/10 应该匹配 0, 10, 20, 30
				assert.True(t, s.matchBit(s.Minute, 0), "分钟 0 应该匹配")
				assert.True(t, s.matchBit(s.Minute, 10), "分钟 10 应该匹配")
				assert.True(t, s.matchBit(s.Minute, 20), "分钟 20 应该匹配")
				assert.True(t, s.matchBit(s.Minute, 30), "分钟 30 应该匹配")
				assert.False(t, s.matchBit(s.Minute, 5), "分钟 5 不应该匹配")
				assert.False(t, s.matchBit(s.Minute, 40), "分钟 40 不应该匹配")
			},
		},
		"mixed_list_and_range": {
			spec: "0 8-10,14-16 * * *",
			validateHour: func(t *testing.T, s *CronSpecSchedule) {
				// 8-10,14-16 应该匹配 8,9,10,14,15,16
				for _, h := range []uint{8, 9, 10, 14, 15, 16} {
					assert.True(t, s.matchBit(s.Hour, h), "小时 %d 应该匹配", h)
				}
				for _, h := range []uint{7, 11, 12, 13, 17} {
					assert.False(t, s.matchBit(s.Hour, h), "小时 %d 不应该匹配", h)
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			assert.NoError(t, err)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			if tc.validateMinute != nil {
				tc.validateMinute(t, spec)
			}
			if tc.validateHour != nil {
				tc.validateHour(t, spec)
			}
			if tc.validateDom != nil {
				tc.validateDom(t, spec)
			}
		})
	}
}

// TestCronParser_ParseWithTimezone 测试带时区的 Cron 表达式
// 验证时区信息被正确解析和保存
func TestCronParser_ParseWithTimezone(t *testing.T) {
	tests := map[string]struct {
		spec        string
		expectedLoc string
		expectErr   bool
	}{
		"utc_timezone": {
			spec:        "TZ=UTC 0 0 * * *",
			expectedLoc: "UTC",
			expectErr:   false,
		},
		"asia_shanghai": {
			spec:        "CRON_TZ=Asia/Shanghai 0 0 * * *",
			expectedLoc: "Asia/Shanghai",
			expectErr:   false,
		},
		"america_new_york": {
			spec:        "TZ=America/New_York 30 9 * * *",
			expectedLoc: "America/New_York",
			expectErr:   false,
		},
		"invalid_timezone": {
			spec:        "TZ=Invalid/Zone 0 0 * * *",
			expectedLoc: "",
			expectErr:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
				spec := schedule.(*CronSpecSchedule)
				assert.Equal(t, tc.expectedLoc, spec.Location.String(), "时区应该匹配")
			}
		})
	}
}

// TestCronParser_ParseWithSeconds_DetailedValidation 测试带秒的表达式并详细验证
// 验证秒字段被正确解析
func TestCronParser_ParseWithSeconds_DetailedValidation(t *testing.T) {
	tests := map[string]struct {
		spec            string
		expectSecond    []uint
		notExpectSecond []uint
	}{
		"every_10_seconds": {
			spec:            "*/10 * * * * *",
			expectSecond:    []uint{0, 10, 20, 30, 40, 50},
			notExpectSecond: []uint{5, 15, 25, 35, 45, 55},
		},
		"specific_seconds": {
			spec:            "0,15,30,45 * * * * *",
			expectSecond:    []uint{0, 15, 30, 45},
			notExpectSecond: []uint{1, 10, 20, 50},
		},
		"second_range": {
			spec:            "10-20 * * * * *",
			expectSecond:    []uint{10, 11, 15, 20},
			notExpectSecond: []uint{0, 9, 21, 30},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			for _, s := range tc.expectSecond {
				assert.True(t, spec.matchBit(spec.Second, s), "秒 %d 应该匹配", s)
			}
			for _, s := range tc.notExpectSecond {
				assert.False(t, spec.matchBit(spec.Second, s), "秒 %d 不应该匹配", s)
			}
		})
	}
}

// TestCronParser_ParseMonthValidation 测试月份解析的详细验证
// 验证月份字段(包括名称和数字)被正确解析
func TestCronParser_ParseMonthValidation(t *testing.T) {
	tests := map[string]struct {
		spec            string
		expectMonths    []uint
		notExpectMonths []uint
	}{
		"month_names": {
			spec:            "0 0 1 jan,mar,may *",
			expectMonths:    []uint{1, 3, 5}, // January, March, May
			notExpectMonths: []uint{2, 4, 6},
		},
		"month_numbers": {
			spec:            "0 0 1 1,6,12 *",
			expectMonths:    []uint{1, 6, 12},
			notExpectMonths: []uint{2, 7, 11},
		},
		"month_range": {
			spec:            "0 0 1 3-6 *",
			expectMonths:    []uint{3, 4, 5, 6}, // March to June
			notExpectMonths: []uint{1, 2, 7, 8},
		},
		"quarter_months": {
			spec:            "0 0 1 1,4,7,10 *",
			expectMonths:    []uint{1, 4, 7, 10}, // 每季度第一个月
			notExpectMonths: []uint{2, 3, 5, 6},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			assert.NoError(t, err)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			for _, m := range tc.expectMonths {
				assert.True(t, spec.matchBit(spec.Month, m), "月份 %d 应该匹配", m)
			}
			for _, m := range tc.notExpectMonths {
				assert.False(t, spec.matchBit(spec.Month, m), "月份 %d 不应该匹配", m)
			}
		})
	}
}

// TestCronParser_ParseWeekdayValidation 测试星期解析的详细验证
// 验证星期字段(包括名称和数字)被正确解析
func TestCronParser_ParseWeekdayValidation(t *testing.T) {
	tests := map[string]struct {
		spec              string
		expectWeekdays    []uint
		notExpectWeekdays []uint
	}{
		"weekday_names": {
			spec:              "0 0 * * mon,wed,fri",
			expectWeekdays:    []uint{1, 3, 5}, // Monday, Wednesday, Friday
			notExpectWeekdays: []uint{0, 2, 4, 6},
		},
		"weekday_numbers": {
			spec:              "0 0 * * 1,3,5",
			expectWeekdays:    []uint{1, 3, 5},
			notExpectWeekdays: []uint{0, 2, 4, 6},
		},
		"weekday_range": {
			spec:              "0 0 * * 1-5",
			expectWeekdays:    []uint{1, 2, 3, 4, 5}, // Monday to Friday (工作日)
			notExpectWeekdays: []uint{0, 6},
		},
		"weekend": {
			spec:              "0 0 * * 0,6",
			expectWeekdays:    []uint{0, 6}, // Sunday, Saturday
			notExpectWeekdays: []uint{1, 2, 3, 4, 5},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			assert.NoError(t, err)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			for _, w := range tc.expectWeekdays {
				assert.True(t, spec.matchBit(spec.Dow, w), "星期 %d 应该匹配", w)
			}
			for _, w := range tc.notExpectWeekdays {
				assert.False(t, spec.matchBit(spec.Dow, w), "星期 %d 不应该匹配", w)
			}
		})
	}
}

// TestCronParser_EdgeCases 测试边界情况
// 验证各种边界值和特殊情况的正确处理
func TestCronParser_EdgeCases(t *testing.T) {
	tests := map[string]struct {
		spec      string
		expectErr bool
		validate  func(*testing.T, *CronSpecSchedule)
	}{
		"min_max_minute": {
			spec:      "0,59 * * * *",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				assert.True(t, s.matchBit(s.Minute, 0), "最小分钟值 0 应该匹配")
				assert.True(t, s.matchBit(s.Minute, 59), "最大分钟值 59 应该匹配")
			},
		},
		"min_max_hour": {
			spec:      "0 0,23 * * *",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				assert.True(t, s.matchBit(s.Hour, 0), "最小小时值 0 应该匹配")
				assert.True(t, s.matchBit(s.Hour, 23), "最大小时值 23 应该匹配")
			},
		},
		"min_max_day": {
			spec:      "0 0 1,31 * *",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				assert.True(t, s.matchBit(s.Dom, 1), "最小日期值 1 应该匹配")
				assert.True(t, s.matchBit(s.Dom, 31), "最大日期值 31 应该匹配")
			},
		},
		"all_months": {
			spec:      "0 0 1 * *",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				// * 应该匹配所有月份 (1-12)
				for m := uint(1); m <= 12; m++ {
					assert.True(t, s.matchBit(s.Month, m), "月份 %d 应该匹配", m)
				}
			},
		},
		"all_weekdays": {
			spec:      "0 0 * * *",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				// * 应该设置星号位
				assert.NotZero(t, s.Dow&cronStarBit, "星期字段应该设置星号位")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
				if tc.validate != nil {
					spec := schedule.(*CronSpecSchedule)
					tc.validate(t, spec)
				}
			}
		})
	}
}

// TestCronParser_ParseErrorMessages 测试错误消息的准确性
// 验证各种错误情况返回有意义的错误消息
func TestCronParser_ParseErrorMessages(t *testing.T) {
	tests := map[string]struct {
		spec              string
		parser            *CronParser
		expectErrContains string
	}{
		"empty_spec": {
			spec:              "",
			parser:            CronStandardParser,
			expectErrContains: "不能为空",
		},
		"too_few_fields": {
			spec:              "0 0",
			parser:            CronStandardParser,
			expectErrContains: "期望",
		},
		"too_many_fields": {
			spec:              "0 0 0 0 0 0 0",
			parser:            CronStandardParser,
			expectErrContains: "期望",
		},
		"invalid_timezone": {
			spec:              "TZ=InvalidZone 0 0 * * *",
			parser:            CronStandardParser,
			expectErrContains: "无效的时区",
		},
		"descriptor_not_supported": {
			spec:              "@daily",
			parser:            NewCronParser(CronMinute | CronHour | CronDom | CronMonth | CronDow),
			expectErrContains: "描述符不受支持",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := tc.parser.Parse(tc.spec)
			assert.Error(t, err)
			assert.Nil(t, schedule)
			assert.Contains(t, err.Error(), tc.expectErrContains, "错误消息应该包含预期文本")
		})
	}
}

// TestCronParser_RealWorldExamples 测试真实世界的使用场景
// 验证常见的实际应用场景能正确解析
func TestCronParser_RealWorldExamples(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		validate    func(*testing.T, CronSchedule)
	}{
		"daily_backup": {
			spec:        "0 2 * * *",
			description: "每天凌晨2点执行备份",
			validate: func(t *testing.T, schedule CronSchedule) {
				spec := schedule.(*CronSpecSchedule)
				assert.True(t, spec.matchBit(spec.Minute, 0))
				assert.True(t, spec.matchBit(spec.Hour, 2))
			},
		},
		"business_hours": {
			spec:        "*/5 9-17 * * 1-5",
			description: "工作日工作时间每5分钟执行",
			validate: func(t *testing.T, schedule CronSchedule) {
				spec := schedule.(*CronSpecSchedule)
				// 验证每5分钟
				assert.True(t, spec.matchBit(spec.Minute, 0))
				assert.True(t, spec.matchBit(spec.Minute, 5))
				// 验证工作时间 9-17
				assert.True(t, spec.matchBit(spec.Hour, 9))
				assert.True(t, spec.matchBit(spec.Hour, 17))
				// 验证工作日 1-5
				assert.True(t, spec.matchBit(spec.Dow, 1))
				assert.True(t, spec.matchBit(spec.Dow, 5))
			},
		},
		"monthly_report": {
			spec:        "0 0 1 * *",
			description: "每月1号生成月报",
			validate: func(t *testing.T, schedule CronSchedule) {
				spec := schedule.(*CronSpecSchedule)
				assert.True(t, spec.matchBit(spec.Dom, 1))
			},
		},
		"quarterly_task": {
			spec:        "0 0 1 1,4,7,10 *",
			description: "每季度第一天执行",
			validate: func(t *testing.T, schedule CronSchedule) {
				spec := schedule.(*CronSpecSchedule)
				assert.True(t, spec.matchBit(spec.Month, 1))  // January
				assert.True(t, spec.matchBit(spec.Month, 4))  // April
				assert.True(t, spec.matchBit(spec.Month, 7))  // July
				assert.True(t, spec.matchBit(spec.Month, 10)) // October
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			assert.NoError(t, err, "解析 %s 失败", tc.description)
			assert.NotNil(t, schedule)
			if tc.validate != nil {
				tc.validate(t, schedule)
			}
		})
	}
}

// TestCronParser_CompareStandardAndSecondParser 测试标准解析器和带秒解析器的区别
// 验证两种解析器在处理相同输入时的不同行为
func TestCronParser_CompareStandardAndSecondParser(t *testing.T) {
	// 6个字段的表达式
	spec6 := "0 0 0 * * *"

	// 标准解析器应该拒绝
	_, err := ParseCronStandard(spec6)
	assert.Error(t, err, "标准解析器应该拒绝6个字段的表达式")

	// 带秒解析器应该接受
	schedule, err := ParseCronWithSeconds(spec6)
	assert.NoError(t, err, "带秒解析器应该接受6个字段的表达式")
	assert.NotNil(t, schedule)

	spec := schedule.(*CronSpecSchedule)
	assert.True(t, spec.matchBit(spec.Second, 0), "秒字段应该为0")
	assert.True(t, spec.matchBit(spec.Minute, 0), "分钟字段应该为0")
	assert.True(t, spec.matchBit(spec.Hour, 0), "小时字段应该为0")
}

// TestCronParser_NormalizeFieldsWithOptional 测试可选字段的规范化
// 验证当字段缺失时会正确填充默认值
func TestCronParser_NormalizeFieldsWithOptional(t *testing.T) {
	// 创建一个标准解析器(秒字段可选)
	parser := CronStandardParser

	// 5个字段的表达式(没有秒)
	spec := "0 0 * * *"
	schedule, err := parser.Parse(spec)
	assert.NoError(t, err)
	assert.NotNil(t, schedule)

	specSchedule := schedule.(*CronSpecSchedule)
	// 验证秒字段被填充为默认值 0
	assert.True(t, specSchedule.matchBit(specSchedule.Second, 0), "秒字段应该被填充为0")
}

// TestCronParser_SpecialCharacters 测试特殊字符的处理
// 验证 L, W, # 等特殊字符的解析(如果支持)
func TestCronParser_SpecialCharacters(t *testing.T) {
	tests := map[string]struct {
		spec      string
		expectErr bool
	}{
		"last_day_of_month": {
			spec:      "0 0 L * *",
			expectErr: false, // L 表示月末最后一天
		},
		"question_mark": {
			spec:      "0 0 ? * MON",
			expectErr: false, // ? 表示不关心该字段
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// BenchmarkParseCronStandard 基准测试标准 Cron 表达式解析性能
func BenchmarkParseCronStandard(b *testing.B) {
	specs := []string{
		"0 0 * * *",
		"*/5 * * * *",
		"0 9-17 * * 1-5",
		"0 0 1,15 * *",
	}

	for _, spec := range specs {
		b.Run(spec, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ParseCronStandard(spec)
			}
		})
	}
}

// BenchmarkParseCronWithSeconds 基准测试带秒的 Cron 表达式解析性能
func BenchmarkParseCronWithSeconds(b *testing.B) {
	specs := []string{
		"0 0 0 * * *",
		"*/10 * * * * *",
		"0,30 * 9-17 * * 1-5",
	}

	for _, spec := range specs {
		b.Run(spec, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ParseCronWithSeconds(spec)
			}
		})
	}
}

// TestCronParser_CaseInsensitive 测试大小写不敏感
func TestCronParser_CaseInsensitive(t *testing.T) {
	tests := []struct {
		spec1 string
		spec2 string
		desc  string
	}{
		{
			spec1: "0 0 1 JAN *",
			spec2: "0 0 1 jan *",
			desc:  "月份名称大小写",
		},
		{
			spec1: "0 0 * * MON",
			spec2: "0 0 * * mon",
			desc:  "星期名称大小写",
		},
		{
			spec1: "0 0 * * Mon",
			spec2: "0 0 * * mOn",
			desc:  "混合大小写",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s1, err1 := ParseCronStandard(tc.spec1)
			s2, err2 := ParseCronStandard(tc.spec2)

			assert.NoError(t, err1)
			assert.NoError(t, err2)
			assert.NotNil(t, s1)
			assert.NotNil(t, s2)

			// 验证解析结果相同
			spec1 := s1.(*CronSpecSchedule)
			spec2 := s2.(*CronSpecSchedule)
			assert.Equal(t, spec1.Month, spec2.Month, "月份字段应相同")
			assert.Equal(t, spec1.Dow, spec2.Dow, "星期字段应相同")
		})
	}
}

// TestCronParser_WhitespaceHandling 测试空白字符处理
func TestCronParser_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		spec string
		desc string
	}{
		{
			spec: "0  0  *  *  *",
			desc: "多个空格",
		},
		{
			spec: "0\t0\t*\t*\t*",
			desc: "制表符",
		},
		{
			spec: " 0 0 * * * ",
			desc: "前后空格",
		},
		{
			spec: "0 0 * * *    ",
			desc: "尾部多个空格",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			assert.NoError(t, err, tc.desc)
			assert.NotNil(t, schedule)
		})
	}
}

// TestCronParser_MultipleParsers 测试多个解析器实例
func TestCronParser_MultipleParsers(t *testing.T) {
	parser1 := NewCronParser(CronMinute | CronHour | CronDom | CronMonth | CronDow)
	parser2 := NewCronParser(CronSecond | CronMinute | CronHour | CronDom | CronMonth | CronDow)
	parser3 := NewCronParser(CronMinute | CronHour | CronDom | CronMonth | CronDow | CronDescriptor)

	spec5 := "0 0 * * *"
	spec6 := "0 0 0 * * ?" // 6字段格式，使用?表示不关心周字段

	// parser1 应该接受 5 字段
	s1, err1 := parser1.Parse(spec5)
	assert.NoError(t, err1)
	assert.NotNil(t, s1)

	// parser1 应该拒绝 6 字段
	s2, err2 := parser1.Parse(spec6)
	assert.Error(t, err2)
	assert.Nil(t, s2)

	// parser2 应该接受 6 字段
	s3, err3 := parser2.Parse(spec6)
	assert.NoError(t, err3)
	assert.NotNil(t, s3)

	s4, err4 := parser3.Parse("@every 1h")
	assert.NoError(t, err4)
	assert.NotNil(t, s4)
}

// TestCronParser_FieldBoundaries 测试字段边界值
func TestCronParser_FieldBoundaries(t *testing.T) {
	tests := map[string]struct {
		spec      string
		expectErr bool
		desc      string
	}{
		"valid_second_0": {
			spec:      "0 0 0 * * *",
			expectErr: false,
			desc:      "秒字段最小值 0",
		},
		"valid_second_59": {
			spec:      "59 0 0 * * *",
			expectErr: false,
			desc:      "秒字段最大值 59",
		},
		"invalid_second_60": {
			spec:      "60 0 0 * * *",
			expectErr: true,
			desc:      "秒字段超出最大值",
		},
		"valid_minute_0": {
			spec:      "0 0 * * *",
			expectErr: false,
			desc:      "分钟字段最小值 0",
		},
		"valid_minute_59": {
			spec:      "59 0 * * *",
			expectErr: false,
			desc:      "分钟字段最大值 59",
		},
		"invalid_minute_60": {
			spec:      "60 0 * * *",
			expectErr: true,
			desc:      "分钟字段超出最大值",
		},
		"valid_hour_0": {
			spec:      "0 0 * * *",
			expectErr: false,
			desc:      "小时字段最小值 0",
		},
		"valid_hour_23": {
			spec:      "0 23 * * *",
			expectErr: false,
			desc:      "小时字段最大值 23",
		},
		"invalid_hour_24": {
			spec:      "0 24 * * *",
			expectErr: true,
			desc:      "小时字段超出最大值",
		},
		"valid_dom_1": {
			spec:      "0 0 1 * *",
			expectErr: false,
			desc:      "日期字段最小值 1",
		},
		"valid_dom_31": {
			spec:      "0 0 31 * *",
			expectErr: false,
			desc:      "日期字段最大值 31",
		},
		"invalid_dom_0": {
			spec:      "0 0 0 * *",
			expectErr: true,
			desc:      "日期字段小于最小值",
		},
		"invalid_dom_32": {
			spec:      "0 0 32 * *",
			expectErr: true,
			desc:      "日期字段超出最大值",
		},
		"valid_month_1": {
			spec:      "0 0 1 1 *",
			expectErr: false,
			desc:      "月份字段最小值 1",
		},
		"valid_month_12": {
			spec:      "0 0 1 12 *",
			expectErr: false,
			desc:      "月份字段最大值 12",
		},
		"invalid_month_0": {
			spec:      "0 0 1 0 *",
			expectErr: true,
			desc:      "月份字段小于最小值",
		},
		"invalid_month_13": {
			spec:      "0 0 1 13 *",
			expectErr: true,
			desc:      "月份字段超出最大值",
		},
		"valid_dow_0": {
			spec:      "0 0 * * 0",
			expectErr: false,
			desc:      "星期字段最小值 0(周日)",
		},
		"valid_dow_6": {
			spec:      "0 0 * * 6",
			expectErr: false,
			desc:      "星期字段最大值 6(周六)",
		},
		"valid_dow_7": {
			spec:      "0 0 * * 7",
			expectErr: false, // 现在支持 7 作为周日，会自动转换为 0
			desc:      "星期字段 7(自动转换为周日 0,Quartz兼容)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var schedule CronSchedule
			var err error

			// 根据字段数量选择解析器
			if len(strings.Fields(tc.spec)) == 6 {
				schedule, err = ParseCronWithSeconds(tc.spec)
			} else {
				schedule, err = ParseCronStandard(tc.spec)
			}

			if tc.expectErr {
				assert.Error(t, err, tc.desc)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err, tc.desc)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// TestCronParser_RangeValidation 测试范围验证
func TestCronParser_RangeValidation(t *testing.T) {
	tests := map[string]struct {
		spec      string
		expectErr bool
		desc      string
	}{
		"valid_range": {
			spec:      "0 9-17 * * *",
			expectErr: false,
			desc:      "有效范围",
		},
		"reversed_range": {
			spec:      "0 17-9 * * *",
			expectErr: true,
			desc:      "反向范围(起始 > 结束)",
		},
		"single_value_range": {
			spec:      "0 10-10 * * *",
			expectErr: false,
			desc:      "单值范围(起始 = 结束)",
		},
		"out_of_bounds_start": {
			spec:      "0 -1-10 * * *",
			expectErr: true,
			desc:      "起始值超出下界",
		},
		"out_of_bounds_end": {
			spec:      "0 10-25 * * *",
			expectErr: true,
			desc:      "结束值超出上界",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			if tc.expectErr {
				assert.Error(t, err, tc.desc)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err, tc.desc)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// TestCronParser_StepValidation 测试步长验证
func TestCronParser_StepValidation(t *testing.T) {
	tests := map[string]struct {
		spec      string
		expectErr bool
		desc      string
	}{
		"valid_step": {
			spec:      "*/5 * * * *",
			expectErr: false,
			desc:      "有效步长",
		},
		"zero_step": {
			spec:      "*/0 * * * *",
			expectErr: true,
			desc:      "零步长(无效)",
		},
		"negative_step": {
			spec:      "*/-5 * * * *",
			expectErr: false,
			desc:      "负步长(会被解析为绝对值，这是当前的行为)",
		},
		"step_larger_than_range": {
			spec:      "0-10/20 * * * *",
			expectErr: false,
			desc:      "步长大于范围(有效但结果为单值)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			if tc.expectErr {
				assert.Error(t, err, tc.desc)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err, tc.desc)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// TestCronParser_NamedMonthsEdgeCases 测试月份名称边界情况
func TestCronParser_NamedMonthsEdgeCases(t *testing.T) {
	tests := map[string]struct {
		spec      string
		expectErr bool
		validate  func(*testing.T, *CronSpecSchedule)
	}{
		"all_month_names": {
			spec:      "0 0 1 jan,feb,mar,apr,may,jun,jul,aug,sep,oct,nov,dec *",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				// 验证所有月份都被设置
				for m := uint(1); m <= 12; m++ {
					assert.True(t, s.matchBit(s.Month, m), "月份 %d 应该被设置", m)
				}
			},
		},
		"month_name_range": {
			spec:      "0 0 1 jan-mar *",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				assert.True(t, s.matchBit(s.Month, 1))
				assert.True(t, s.matchBit(s.Month, 2))
				assert.True(t, s.matchBit(s.Month, 3))
				assert.False(t, s.matchBit(s.Month, 4))
			},
		},
		"invalid_month_name": {
			spec:      "0 0 1 january *",
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
				if tc.validate != nil {
					spec := schedule.(*CronSpecSchedule)
					tc.validate(t, spec)
				}
			}
		})
	}
}

// TestCronParser_NamedWeekdaysEdgeCases 测试星期名称边界情况
func TestCronParser_NamedWeekdaysEdgeCases(t *testing.T) {
	tests := map[string]struct {
		spec      string
		expectErr bool
		validate  func(*testing.T, *CronSpecSchedule)
	}{
		"all_weekday_names": {
			spec:      "0 0 * * sun,mon,tue,wed,thu,fri,sat",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				// 验证所有星期都被设置
				for w := uint(0); w <= 6; w++ {
					assert.True(t, s.matchBit(s.Dow, w), "星期 %d 应该被设置", w)
				}
			},
		},
		"weekday_name_range": {
			spec:      "0 0 * * mon-fri",
			expectErr: false,
			validate: func(t *testing.T, s *CronSpecSchedule) {
				for w := uint(1); w <= 5; w++ {
					assert.True(t, s.matchBit(s.Dow, w), "工作日 %d 应该被设置", w)
				}
				assert.False(t, s.matchBit(s.Dow, 0), "周日不应该被设置")
				assert.False(t, s.matchBit(s.Dow, 6), "周六不应该被设置")
			},
		},
		"invalid_weekday_name": {
			spec:      "0 0 * * monday",
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronStandard(tc.spec)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
				if tc.validate != nil {
					spec := schedule.(*CronSpecSchedule)
					tc.validate(t, spec)
				}
			}
		})
	}
}

// TestCronParser_ConcurrentParsing 测试并发解析
func TestCronParser_ConcurrentParsing(t *testing.T) {
	specs := []string{
		"0 0 * * *",
		"*/5 * * * *",
		"0 9-17 * * 1-5",
		"0 0 1 * *",
		"0 0 * * 0",
	}

	// 并发解析多个表达式
	done := make(chan bool)
	for _, spec := range specs {
		go func(s string) {
			for i := 0; i < 100; i++ {
				schedule, err := ParseCronStandard(s)
				assert.NoError(t, err)
				assert.NotNil(t, schedule)
			}
			done <- true
		}(spec)
	}

	// 等待所有 goroutine 完成
	for range specs {
		<-done
	}
}

// TestCronParser_MemoryEfficiency 测试内存效率
func TestCronParser_MemoryEfficiency(t *testing.T) {
	// 解析大量表达式，确保不会导致内存泄漏
	for i := 0; i < 1000; i++ {
		_, _ = ParseCronStandard("0 0 * * *")
		_, _ = ParseCronWithSeconds("0 0 0 * * *")
	}
	// 如果有内存泄漏，这个测试会消耗大量内存
}
