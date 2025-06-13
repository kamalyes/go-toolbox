/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-09 10:19:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-12 11:15:55
 * @FilePath: \go-toolbox\pkg\schedule\errors.go
 * @Description: 公共错误信息常量
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

const (
	ErrWeekdayHashFormat       = "星期字段中 '#' 语法格式错误: %s"
	ErrWeekdayHashValue        = "星期字段中 '#' 语法值无效: %s"
	ErrWeekdayValueRangeZero   = "星期值超出范围 [0-6]: %d"
	ErrWeekdayValueRangeOne    = "星期值超出范围 [1-7]: %d"
	ErrWeekdayUnknownNumbering = "未知的星期数字体系"
	ErrWeekdayValueInvalid     = "星期字段值无效: %s"
	ErrWeekdayQuestionMutual   = "日字段和星期字段不能同时为 '?'"
	ErrDayWFormat              = "日字段 W 格式无效: %s"
	ErrDayWOutOfRange          = "日字段 W 的日期超出范围: %d"
	ErrDayLComplexFormat       = "暂不支持日字段中 L 的复杂格式: %s"
	ErrStepFormat              = "步进表达式格式错误: %s"
	ErrStepValueInvalid        = "步进值无效: %s"
	ErrRangeFormat             = "范围格式错误: %s"
	ErrRangeFormatRightMissing = "范围格式错误，缺少右边界: %s"
	ErrRangeStartInvalid       = "范围起始值无效: %s"
	ErrRangeEndInvalid         = "范围结束值无效: %s"
	ErrRangeStartGreater       = "范围起始值大于结束值: %s"
	ErrValueInvalid            = "值无效: %s"
	ErrValueOutOfRange         = "值超出范围 [%d-%d]: %v"
	ErrFieldCount              = "Cron表达式必须包含7个字段（秒 分 时 日 月 星期 年）"
	ErrParseSecond             = "秒字段解析失败"
	ErrParseMinute             = "分字段解析失败"
	ErrParseHour               = "时字段解析失败"
	ErrParseDay                = "日字段解析失败"
	ErrParseMonth              = "月字段解析失败"
	ErrParseWeekday            = "星期字段解析失败"
	ErrParseYear               = "年字段解析失败"
)

const (
	ErrJobIDEmpty                  = "任务ID不能为空"
	ErrJobNameEmpty                = "任务名称不能为空"
	ErrJobCallbackNil              = "任务回调函数不能为空"
	ErrJobExpressionEmpty          = "任务调度表达式不能为空"
	ErrJobIDAlreadyExists          = "任务ID %s 已存在"
	ErrAtLeastOneRuleMustBeDefined = "至少必须定义一条规则"
)
