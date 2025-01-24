/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-24 13:55:15
 * @FilePath: \go-toolbox\pkg\schedule\constant.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

// 调度规则的常量
const (
	EverySecond     = "0 * * * * *"       // 每秒，表示每分钟的第0秒执行
	EveryMinute     = "0 */1 * * * *"     // 每分钟，表示每分钟的第0秒执行
	EveryHalfMinute = "*/30 * * * * *"    // 每半分钟，表示每30秒执行一次
	EveryHour       = "0 0 */1 * * *"     // 每小时，表示每小时的第0分钟和第0秒执行
	EveryHalfHour   = "0 */30 * * *"      // 每半小时，表示每半小时的第0分钟和第0秒执行
	EveryDay        = "0 0 1 * * *"       // 每天，表示每天的0点0分执行
	EveryHalfDay    = "0 0,12 * * *"      // 每半天，表示每天的0点和12点执行
	EveryWeek       = "0 0 * * 0"         // 每周，表示每周日的0点0分执行（可以根据需要调整）
	WeekdaysOnly    = "0 0 * * 1-5"       // 仅工作日，表示每个工作日（周一到周五）的0点0分执行
	WeekendsOnly    = "0 0 * * 6,0"       // 仅非工作日，表示每个非工作日（周六和周日）的0点0分执行
	PeakHours       = "0 8-9,17-18 * * *" // 高峰期，表示每天的8:00-9:00和17:00-18:00执行
	OffPeakHours    = "0 9-17 * * *"      // 低谷期，表示每天的9:00-17:00执行
	Monday          = "0 0 * * 1"         // 星期一，表示每周一的0点0分执行
	Tuesday         = "0 0 * * 2"         // 星期二，表示每周二的0点0分执行
	Wednesday       = "0 0 * * 3"         // 星期三，表示每周三的0点0分执行
	Thursday        = "0 0 * * 4"         // 星期四，表示每周四的0点0分执行
	Friday          = "0 0 * * 5"         // 星期五，表示每周五的0点0分执行
	Saturday        = "0 0 * * 6"         // 星期六，表示每周六的0点0分执行
	Sunday          = "0 0 * * 0"         // 星期日，表示每周日的0点0分执行
	FirstDayOfMonth = "0 0 1 * *"         // 每月的第一天，表示每月的第一天的0点0分执行
	LastDayOfMonth  = "0 0 L * *"         // 每月的最后一天，表示每月最后一天的0点0分执行
)
