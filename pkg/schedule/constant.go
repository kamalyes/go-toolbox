/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 11:55:15
 * @FilePath: \go-toolbox\pkg\schedule\constant.go
 * @Description: 公共配置结构和常量
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import "time"

// ------------------- 配置结构和常量 -------------------

// 字段取值范围定义
type FieldRange struct {
	Min, Max int
}

// 各字段对应的合法范围
var (
	nowYear     = time.Now().Year()
	secondRange = FieldRange{0, 59}              // 秒 0-59
	minuteRange = FieldRange{0, 59}              // 分 0-59
	hourRange   = FieldRange{0, 23}              // 时 0-23
	dayRange    = FieldRange{1, 31}              // 日 1-31
	monthRange  = FieldRange{1, 12}              // 月 1-12
	yearRange   = FieldRange{1970, nowYear + 99} // 年 (1970)~(nowYear+99)
	cronRange   = FieldRange{5, 7}               // cron表达式 5~7
)

// 星期字段数字体系配置
type WeekdayNumbering int

const (
	WeekdayZeroBased WeekdayNumbering = iota // 0=周日，6=周六，Go标准
	WeekdayOneBased                          // 1=周日，7=周六，传统Cron
)

const (
	// 星期字段数字范围定义
	WeekdayZeroBasedMin = 0 // 0基表示，周日为0
	WeekdayZeroBasedMax = 6 // 0基最大值6（周六）
	WeekdayOneBasedMin  = 1 // 1基表示，周日为1
	WeekdayOneBasedMax  = 7 // 1基最大值7（周六）
)

// 调度规则的常量 秒、分、时、日、月、星期、年
const (
	EverySecond     = "*/1 * * * * *"       // 每秒执行一次
	EveryMinute     = "0 */1 * * * *"       // 每分钟的第0秒执行
	EveryHalfMinute = "0,30 * * * * *"      // 每分钟的第0秒和第30秒执行
	EveryHour       = "0 0 */1 * * *"       // 每小时的第0分0秒执行
	EveryHalfHour   = "0 0,30 * * * *"      // 每小时的第0分和第30分的第0秒执行
	EveryDay        = "0 0 0 * * *"         // 每天的0点0分0秒执行
	EveryHalfDay    = "0 0 0,12 * * *"      // 每天的0点和12点的0分0秒执行
	EveryWeek       = "0 0 0 * * 0"         // 每周日的0点0分0秒执行
	WeekdaysOnly    = "0 0 0 * * 1-5"       // 每周一到周五的0点0分0秒执行
	WeekendsOnly    = "0 0 0 * * 6,0"       // 每周六和周日的0点0分0秒执行
	PeakHours       = "0 0 8-9,17-18 * * *" // 每天8点到9点和17点到18点，每分钟的第0秒执行
	OffPeakHours    = "0 0 9-16 * * *"      // 每天9点到16点，每分钟的第0秒执行
	Monday          = "0 0 0 * * 1"         // 每周一的0点0分0秒执行
	Tuesday         = "0 0 0 * * 2"         // 每周二的0点0分0秒执行
	Wednesday       = "0 0 0 * * 3"         // 每周三的0点0分0秒执行
	Thursday        = "0 0 0 * * 4"         // 每周四的0点0分0秒执行
	Friday          = "0 0 0 * * 5"         // 每周五的0点0分0秒执行
	Saturday        = "0 0 0 * * 6"         // 每周六的0点0分0秒执行
	Sunday          = "0 0 0 * * 0"         // 每周日的0点0分0秒执行
	FirstDayOfMonth = "0 0 0 1 * *"         // 每月第一天的0点0分0秒执行
	LastDayOfMonth  = "0 0 0 L * *"         // 每月的最后一天，表示每月最后一天的0点0分执行 魔改下
)

// ------------------- 公共符号常量 -------------------

// 星期字段相关符号
const (
	StarSymbol     = "*" // 星号，表示全范围或任意值
	HashSymbol     = "#" // 星期字段的扩展符号，表示第几个星期几
	QuestionSymbol = "?" // 表示不指定该字段
	WSymbol        = "W" // 日字段中表示最近工作日的符号
	LSymbol        = "L" // 日字段中表示最后一天的符号
	CommaSymbol    = "," // 多值分隔符
	RangeSymbol    = "-" // 范围符号
	StepSymbol     = "/" // 步进符号
)

// 定义状态与日志消息的映射
var execStatusLogMessages = map[execStatus]string{
	Failure:         "Task failed: %v",
	Success:         "Task executed successfully",
	Pending:         "Task is pending",
	Running:         "Task is running",
	SysTermination:  "Task terminated by system",
	UserTermination: "Task terminated by user",
}

var DefaultTimeZone, _ = time.LoadLocation("Asia/Shanghai") // 默认时区 上海
