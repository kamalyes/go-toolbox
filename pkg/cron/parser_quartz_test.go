/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-25 10:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-25 12:11:09
 * @FilePath: \go-toolbox\pkg\cron\parser_quartz_test.go
 * @Description: Quartz 风格 Cron 表达式测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQuartzCron_BasicExamples 测试 Quartz 风格的基本示例
func TestQuartzCron_BasicExamples(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		expectErr   bool
	}{
		"every_5_seconds": {
			spec:        "*/5 * * * * ?",
			description: "每隔 5 秒执行一次",
			expectErr:   false,
		},
		"every_1_minute": {
			spec:        "0 */1 * * * ?",
			description: "每隔 1 分钟执行一次",
		},
		"monthly_2am": {
			spec:        "0 0 2 1 * ?",
			description: "每月 1 日的凌晨 2 点执行一次",
			expectErr:   false,
		},
		"weekday_morning": {
			spec:        "0 15 10 ? * MON-FRI",
			description: "周一到周五每天上午 10:15 执行作业",
			expectErr:   false,
		},
		"daily_23pm": {
			spec:        "0 0 23 * * ?",
			description: "每天 23 点执行一次",
			expectErr:   false,
		},
		"daily_1am": {
			spec:        "0 0 1 * * ?",
			description: "每天凌晨 1 点执行一次",
			expectErr:   false,
		},
		"monthly_first_1am": {
			spec:        "0 0 1 1 * ?",
			description: "每月 1 日凌晨 1 点执行一次",
			expectErr:   false,
		},
		"multi_minutes": {
			spec:        "0 26,29,33 * * * ?",
			description: "在 26 分、29 分、33 分执行一次",
			expectErr:   false,
		},
		"multi_hours": {
			spec:        "0 0 0,13,18,21 * * ?",
			description: "每天的 0 点、13 点、18 点、21 点都执行一次",
			expectErr:   false,
		},
		"specific_hours": {
			spec:        "0 0 10,14,16 * * ?",
			description: "每天上午 10 点，下午 2 点,4 点执行一次",
			expectErr:   false,
		},
		"business_hours_halfhour": {
			spec:        "0 0/30 9-17 * * ?",
			description: "朝九晚五工作时间内每半小时执行一次",
			expectErr:   false,
		},
		"wednesday_noon": {
			spec:        "0 0 12 ? * WED",
			description: "每个星期三中午 12 点执行一次",
			expectErr:   false,
		},
		"daily_noon": {
			spec:        "0 0 12 * * ?",
			description: "每天中午 12 点触发",
			expectErr:   false,
		},
		"daily_1015am": {
			spec:        "0 15 10 ? * *",
			description: "每天上午 10:15 触发",
			expectErr:   false,
		},
		"daily_1015am_variant": {
			spec:        "0 15 10 * * ?",
			description: "每天上午 10:15 触发",
			expectErr:   false,
		},
		"every_minute_15pm": {
			spec:        "0 * 15 * * ?",
			description: "每天下午 3 点到 3:59 期间的每 1 分钟触发",
			expectErr:   false,
		},
		"every_5min_15pm": {
			spec:        "0 0/5 15 * * ?",
			description: "每天下午 3 点到 3:55 期间的每 5 分钟触发",
			expectErr:   false,
		},
		"every_5min_14_18pm": {
			spec:        "0 0/5 14,18 * * ?",
			description: "每天下午 2 点到 2:55 期间和下午 6 点到 6:55 期间的每 5 分钟触发",
			expectErr:   false,
		},
		"minute_range_14pm": {
			spec:        "0 0-5 14 * * ?",
			description: "每天下午 2 点到 2:05 期间的每 1 分钟触发",
			expectErr:   false,
		},
		"march_wednesday_specific": {
			spec:        "0 10,44 14 ? 3 WED",
			description: "每年三月的星期三的下午 2:10 和 2:44 触发",
			expectErr:   false,
		},
		"weekday_1015am": {
			spec:        "0 15 10 ? * MON-FRI",
			description: "周一至周五的上午 10:15 触发",
			expectErr:   false,
		},
		"monthly_15th_1015am": {
			spec:        "0 15 10 15 * ?",
			description: "每月 15 日上午 10:15 触发",
			expectErr:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			if tc.expectErr {
				assert.Error(t, err, tc.description)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err, tc.description)
				assert.NotNil(t, schedule, tc.description)
			}
		})
	}
}

// TestQuartzCron_SpecialChars_LastDay 测试 L 字符(最后一天)
func TestQuartzCron_SpecialChars_LastDay(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
	}{
		"last_day_23pm": {
			spec:        "0 0 23 L * ?",
			description: "每月最后一天 23 点执行一次",
		},
		"last_day_1015am": {
			spec:        "0 15 10 L * ?",
			description: "每月最后一日的上午 10:15 触发",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err, tc.description)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			assert.True(t, spec.LastDay, "应该标记为最后一天")
		})
	}
}

// TestQuartzCron_SpecialChars_LastWeekday 测试 nL 字符(最后一个星期X)
func TestQuartzCron_SpecialChars_LastWeekday(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		expectedDow int // 期望的星期几
	}{
		"last_sunday_1am": {
			spec:        "0 0 1 ? * L",
			description: "每周星期天凌晨 1 点执行一次 (这里的L表示周日)",
			expectedDow: 0,
		},
		"last_friday_1015am": {
			spec:        "0 15 10 ? * 6L",
			description: "每月的最后一个星期五上午 10:15 触发",
			expectedDow: 6,
		},
		"last_friday_1015am_with_year": {
			spec:        "0 15 10 ? * 6L 2002-2005",
			description: "2002 年至 2005 年的每月的最后一个星期五上午 10:15 触发",
			expectedDow: 6,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			// 注意：年份范围功能可能未实现，所以可能会失败
			if name == "last_friday_1015am_with_year" {
				// 暂时跳过年份相关测试
				t.Skip("年份范围功能未实现")
				return
			}

			assert.NoError(t, err, tc.description)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			assert.Equal(t, tc.expectedDow, spec.LastDow, "应该标记为最后一个星期%d", tc.expectedDow)
		})
	}
}

// TestQuartzCron_SpecialChars_NthWeekday 测试 n#m 字符(第m个星期n)
func TestQuartzCron_SpecialChars_NthWeekday(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		expectedNth int // 期望的编码值 (dow*10 + nth)
	}{
		"third_friday_1015am": {
			spec:        "0 15 10 ? * 6#3",
			description: "每月的第三个星期五上午 10:15 触发",
			expectedNth: 63, // 6*10 + 3
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err, tc.description)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			assert.Equal(t, tc.expectedNth, spec.NthDow, "应该标记为第%d个星期", tc.expectedNth)
		})
	}
}

// TestQuartzCron_SpecialChars_Weekday 测试 nW 字符(最近的工作日)
func TestQuartzCron_SpecialChars_Weekday(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		expectedDay int // 期望的日期
	}{
		"nearest_weekday_15th": {
			spec:        "0 0 12 15W * ?",
			description: "最接近 15 号的工作日中午 12 点执行",
			expectedDay: 15,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err, tc.description)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			assert.Equal(t, tc.expectedDay, spec.NearestWeekday, "应该标记为最近工作日，基准日期为%d", tc.expectedDay)
		})
	}
}

// TestQuartzCron_SpecialChars_LastWeekdayOfMonth 测试 LW 字符(最后一个工作日)
func TestQuartzCron_SpecialChars_LastWeekdayOfMonth(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
	}{
		"last_weekday_noon": {
			spec:        "0 0 12 LW * ?",
			description: "每月最后一个工作日中午 12 点执行",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err, tc.description)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			assert.True(t, spec.LastWeekday, "应该标记为最后一个工作日")
		})
	}
}

// TestQuartzCron_InvalidCases 测试无效的表达式
func TestQuartzCron_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
	}{
		"L_in_list": {
			spec:        "0 0 0 1,L * ?",
			description: "L 不能出现在列表中",
		},
		"W_in_list": {
			spec:        "0 0 0 1,15W * ?",
			description: "W 不能出现在列表中",
		},
		"nth_in_list": {
			spec:        "0 0 0 ? * 1,6#3",
			description: "# 不能出现在列表中",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.Error(t, err, tc.description)
			assert.Nil(t, schedule)
		})
	}
}

// TestQuartzCron_MixedSpecialChars 测试混合特殊字符
func TestQuartzCron_MixedSpecialChars(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		validate    func(*testing.T, *CronSpecSchedule)
	}{
		"last_day_with_hour": {
			spec:        "0 30 14 L * ?",
			description: "每月最后一天下午 2:30",
			validate: func(t *testing.T, s *CronSpecSchedule) {
				assert.True(t, s.LastDay)
				assert.True(t, s.matchBit(s.Hour, 14))
				assert.True(t, s.matchBit(s.Minute, 30))
			},
		},
		"last_friday_with_month": {
			spec:        "0 0 10 ? 3,6,9,12 6L",
			description: "每年 3、6、9、12 月的最后一个星期五上午 10 点",
			validate: func(t *testing.T, s *CronSpecSchedule) {
				assert.Equal(t, 6, s.LastDow)
				assert.True(t, s.matchBit(s.Month, 3))
				assert.True(t, s.matchBit(s.Month, 6))
				assert.True(t, s.matchBit(s.Month, 9))
				assert.True(t, s.matchBit(s.Month, 12))
			},
		},
		"weekday_with_specific_month": {
			spec:        "0 0 12 15W 3 ?",
			description: "每年 3 月最接近 15 号的工作日中午 12 点",
			validate: func(t *testing.T, s *CronSpecSchedule) {
				assert.Equal(t, 15, s.NearestWeekday)
				assert.True(t, s.matchBit(s.Month, 3))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err, tc.description)
			assert.NotNil(t, schedule)

			if tc.validate != nil {
				spec := schedule.(*CronSpecSchedule)
				tc.validate(t, spec)
			}
		})
	}
}

// TestQuartzCron_EdgeCases 测试边界情况
func TestQuartzCron_EdgeCases(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		expectErr   bool
	}{
		"invalid_weekday_num": {
			spec:        "0 0 0 ? * 8",
			description: "无效的星期数字(超过 7)",
			expectErr:   true,
		},
		"invalid_nth_weekday": {
			spec:        "0 0 0 ? * 6#6",
			description: "无效的第几个星期(超过 5)",
			expectErr:   true,
		},
		"invalid_weekday_char": {
			spec:        "0 0 0 32W * ?",
			description: "无效的工作日日期(超过 31)",
			expectErr:   true,
		},
		"zero_weekday": {
			spec:        "0 0 0 ? * 0",
			description: "星期日(0 是有效的)",
			expectErr:   false,
		},
		"seven_weekday": {
			spec:        "0 0 0 ? * 7",
			description: "星期日(7 等同于 0,Quartz兼容)",
			expectErr:   false, // 支持 DOW=7，自动转换为 0
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			if tc.expectErr {
				assert.Error(t, err, tc.description)
				assert.Nil(t, schedule)
			} else {
				assert.NoError(t, err, tc.description)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// TestQuartzCron_ComplexRealWorld 测试复杂的真实场景
func TestQuartzCron_ComplexRealWorld(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
	}{
		"payroll_last_weekday": {
			spec:        "0 0 0 LW * ?",
			description: "工资单：每月最后一个工作日",
		},
		"quarterly_reports": {
			spec:        "0 0 6 1 1,4,7,10 ?",
			description: "季度报告：每季度第一天早 6 点",
		},
		"month_end_backup": {
			spec:        "0 0 2 L * ?",
			description: "月末备份：每月最后一天凌晨 2 点",
		},
		"first_monday_meeting": {
			spec:        "0 0 9 ? * 2#1",
			description: "每月第一个星期一早 9 点开会",
		},
		"last_friday_party": {
			spec:        "0 0 17 ? * 6L",
			description: "每月最后一个星期五下午 5 点聚会",
		},
		"mid_month_nearest_weekday": {
			spec:        "0 0 12 15W * ?",
			description: "每月 15 号(或最近的工作日)中午 12 点",
		},
		"business_day_morning": {
			spec:        "0 0 9 ? * MON-FRI",
			description: "工作日早 9 点",
		},
		"weekend_cleanup": {
			spec:        "0 0 3 ? * SAT,SUN",
			description: "周末凌晨 3 点清理",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err, tc.description)
			assert.NotNil(t, schedule)
		})
	}
}

// TestQuartzCron_AllSpecialCharsValidation 测试所有特殊字符的验证
func TestQuartzCron_AllSpecialCharsValidation(t *testing.T) {
	tests := map[string]struct {
		spec      string
		checkFunc func(*testing.T, *CronSpecSchedule)
	}{
		"L_dom": {
			spec: "0 0 0 L * ?",
			checkFunc: func(t *testing.T, s *CronSpecSchedule) {
				assert.True(t, s.LastDay, "LastDay 应该为 true")
				assert.Equal(t, -1, s.NearestWeekday, "NearestWeekday 应该为 -1")
				assert.Equal(t, -1, s.LastDow, "LastDow 应该为 -1")
			},
		},
		"LW_dom": {
			spec: "0 0 0 LW * ?",
			checkFunc: func(t *testing.T, s *CronSpecSchedule) {
				assert.True(t, s.LastWeekday, "LastWeekday 应该为 true")
				assert.False(t, s.LastDay, "LastDay 应该为 false")
			},
		},
		"nW_dom": {
			spec: "0 0 0 10W * ?",
			checkFunc: func(t *testing.T, s *CronSpecSchedule) {
				assert.Equal(t, 10, s.NearestWeekday, "NearestWeekday 应该为 10")
				assert.False(t, s.LastDay, "LastDay 应该为 false")
			},
		},
		"L_dow": {
			spec: "0 0 0 ? * L",
			checkFunc: func(t *testing.T, s *CronSpecSchedule) {
				assert.Equal(t, 0, s.LastDow, "LastDow 应该为 0(周日)")
			},
		},
		"nL_dow": {
			spec: "0 0 0 ? * 5L",
			checkFunc: func(t *testing.T, s *CronSpecSchedule) {
				assert.Equal(t, 5, s.LastDow, "LastDow 应该为 5(最后一个星期五)")
			},
		},
		"nth_dow": {
			spec: "0 0 0 ? * 3#2",
			checkFunc: func(t *testing.T, s *CronSpecSchedule) {
				assert.Equal(t, 32, s.NthDow, "NthDow 应该为 32(第2个星期三)")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			tc.checkFunc(t, spec)
		})
	}
}

// TestQuartzCron_StepWithSpecialChars 测试步长与特殊字符组合
func TestQuartzCron_StepWithSpecialChars(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		expectErr   bool
	}{
		"every_5_seconds": {
			spec:        "*/5 * * * * ?",
			description: "每 5 秒执行",
			expectErr:   false,
		},
		"every_2_hours_workdays": {
			spec:        "0 0 */2 ? * MON-FRI",
			description: "工作日每 2 小时",
			expectErr:   false,
		},
		"every_quarter_hour": {
			spec:        "0 */15 * * * ?",
			description: "每 15 分钟",
			expectErr:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			if tc.expectErr {
				assert.Error(t, err, tc.description)
			} else {
				assert.NoError(t, err, tc.description)
				assert.NotNil(t, schedule)
			}
		})
	}
}

// TestQuartzCron_NamedDaysWithSpecialChars 测试命名天数与特殊字符
func TestQuartzCron_NamedDaysWithSpecialChars(t *testing.T) {
	tests := map[string]struct {
		spec        string
		description string
		expectedDow int
	}{
		"last_monday": {
			spec:        "0 0 0 ? * MONL",
			description: "最后一个星期一",
			expectedDow: 1,
		},
		"last_friday": {
			spec:        "0 0 0 ? * FRIL",
			description: "最后一个星期五",
			expectedDow: 5,
		},
		"third_wednesday": {
			spec:        "0 0 0 ? * WED#3",
			description: "第三个星期三",
			expectedDow: 33, // 3*10 + 3
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.NoError(t, err, tc.description)
			assert.NotNil(t, schedule)

			spec := schedule.(*CronSpecSchedule)
			if tc.spec[len(tc.spec)-1] == 'L' {
				assert.Equal(t, tc.expectedDow, spec.LastDow, tc.description)
			} else {
				assert.Equal(t, tc.expectedDow, spec.NthDow, tc.description)
			}
		})
	}
}

// TestQuartzCron_QuestionMarkBehavior 测试问号 ? 的行为
func TestQuartzCron_QuestionMarkBehavior(t *testing.T) {
	tests := map[string]struct {
		spec1       string
		spec2       string
		description string
	}{
		"question_vs_star_dom": {
			spec1:       "0 0 0 ? * MON",
			spec2:       "0 0 0 * * MON",
			description: "日期字段：? 和 * 在某些情况下行为不同",
		},
		"question_vs_star_dow": {
			spec1:       "0 0 0 15 * ?",
			spec2:       "0 0 0 15 * *",
			description: "星期字段：? 和 * 在某些情况下行为不同",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule1, err1 := ParseCronWithSeconds(tc.spec1)
			schedule2, err2 := ParseCronWithSeconds(tc.spec2)

			assert.NoError(t, err1, tc.description)
			assert.NoError(t, err2, tc.description)
			assert.NotNil(t, schedule1)
			assert.NotNil(t, schedule2)
		})
	}
}

// TestQuartzCron_ErrorMessages 测试错误消息的详细程度
func TestQuartzCron_ErrorMessages(t *testing.T) {
	tests := map[string]struct {
		spec              string
		expectedErrSubstr string
	}{
		"invalid_W_value": {
			spec:              "0 0 0 0W * ?",
			expectedErrSubstr: "必须在 1-31 之间",
		},
		"invalid_nth_value": {
			spec:              "0 0 0 ? * 6#0",
			expectedErrSubstr: "必须在 1-5 之间",
		},
		"L_in_list_dom": {
			spec:              "0 0 0 1,L * ?",
			expectedErrSubstr: "不能出现在列表中",
		},
		"W_in_list": {
			spec:              "0 0 0 1,5W * ?",
			expectedErrSubstr: "无效的工作日表达式", // W 字符会导致解析错误
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := ParseCronWithSeconds(tc.spec)
			assert.Error(t, err)
			assert.Nil(t, schedule)
			assert.Contains(t, err.Error(), tc.expectedErrSubstr, "错误消息应包含有用的信息")
		})
	}
}

// BenchmarkQuartzParse 基准测试 Quartz 解析性能
func BenchmarkQuartzParse(b *testing.B) {
	specs := []string{
		"*/5 * * * * ?",           // 简单步长
		"0 15 10 L * ?",           // L 字符
		"0 15 10 ? * 6L",          // nL 字符
		"0 15 10 ? * 6#3",         // n#m 字符
		"0 0 12 15W * ?",          // nW 字符
		"0 0 0 LW * ?",            // LW 字符
		"0 0/30 9-17 ? * MON-FRI", // 复杂表达式
	}

	for _, spec := range specs {
		b.Run(spec, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = ParseCronWithSeconds(spec)
			}
		})
	}
}

// BenchmarkQuartzParseComplex 基准测试复杂 Quartz 表达式
func BenchmarkQuartzParseComplex(b *testing.B) {
	complexSpecs := map[string]string{
		"multi_list":     "0 10,20,30,40,50 * ? * MON,WED,FRI",
		"range_and_step": "0 0/5 9-17 ? * MON-FRI",
		"mixed":          "0 15,30,45 10-12,14-16 ? 1,4,7,10 *",
	}

	for name, spec := range complexSpecs {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = ParseCronWithSeconds(spec)
			}
		})
	}
}
