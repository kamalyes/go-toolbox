/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-25 10:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-25 11:05:59
 * @FilePath: \go-toolbox\pkg\cron\doc.go
 * @Description: Cron 包 - 主入口文件，提供 cron 表达式解析和调度功能
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cron

// 此包提供完整的 Cron 表达式解析和调度功能
//
// 包的实现已被拆分到以下文件中：
// - types.go: 基本类型定义(CronParseOption, CronSchedule, CronSpecSchedule 等)
// - constants.go: 常量定义(字段范围、位掩码等)
// - parser.go: 解析器实现(NewCronParser, Parse 等)
// - schedule.go: 调度逻辑实现(Next 方法等)
// - descriptor.go: 描述符解析(@yearly, @monthly 等)
// - expression.go: 表达式解析辅助函数(ParseFieldToBits 等)
//
// 使用示例：
//
//	使用标准解析器(5个字段：分 时 日 月 周)
//	schedule, err := ParseCronStandard("0 2 * * *")
//	nextTime := schedule.Next(time.Now())
//
//	使用带秒的解析器(6个字段：秒 分 时 日 月 周)
//	schedule, err := ParseCronWithSeconds("0 0 2 * * *")
//	nextTime := schedule.Next(time.Now())
//
//	使用描述符
//	schedule, err := ParseCronStandard("@daily")
//	nextTime := schedule.Next(time.Now())
//
//	使用间隔调度
//	schedule, err := ParseCronStandard("@every 1h30m")
//	nextTime := schedule.Next(time.Now())
//
//	自定义解析器
//	parser := NewCronParser(CronSecond | CronMinute | CronHour | CronDom | CronMonth | CronDow)
//	schedule, err := parser.Parse("*/5 * * * * *")
//	nextTime := schedule.Next(time.Now())
