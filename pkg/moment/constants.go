/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-18 11:20:07
 * @FilePath: \go-toolbox\pkg\moment\constants.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import "time"

// 通用的时间单位持续时长，这里泛指国际基础单位制（民用日）所理解的时间，不考虑夏时制，不用作科学与天文学
const (
	// 默认时区
	DefaultTimezone = "Asia/Shanghai"

	// 默认时间格式
	DefaultTimeFormat = "2006-01-02 15:04:05"

	// 紧凑格式 - 用于文件名、日志、缓存键等场景
	CompactDateFormat          = "20060102"           // 紧凑日期格式 YYYYMMDD
	CompactDateHourFormat      = "2006010215"         // 紧凑日期+小时格式 YYYYMMDDHH
	CompactDateTimeFormat      = "200601021504"       // 紧凑日期时间格式 YYYYMMDDHHMM
	CompactDateTimeSecFormat   = "20060102150405"     // 紧凑日期时间秒格式 YYYYMMDDHHMMSS
	CompactDateTimeMilliFormat = "20060102150405.000" // 紧凑日期时间毫秒格式 YYYYMMDDHHMMSS.sss

	// 纳秒，作为最基础的时间单位
	NanosecondDuration time.Duration = 1

	// 微秒，表示1微秒持续的纳秒时长
	MicrosecondDuration = 1000 * NanosecondDuration

	// 毫秒，表示1毫秒持续的纳秒时长
	MillisecondDuration = 1000 * MicrosecondDuration

	// 秒，表示1秒持续的纳秒时长
	SecondDuration = 1000 * MillisecondDuration

	// 分钟，表示1分钟持续的纳秒时长
	MinuteDuration = 60 * SecondDuration

	// 小时，表示1小时持续的纳秒时长
	HourDuration = 60 * MinuteDuration

	// 天，表示1天持续的纳秒时长
	// 这里不考虑夏时制问题，泛指国际基础单位制（民用日）所理解的时间
	DayDuration = 24 * HourDuration

	// 周，表示1周持续的纳秒时长
	// 这里不考虑夏时制问题，泛指国际基础单位制（民用日）所理解的时间
	WeekDuration = 7 * DayDuration

	// 月，表示30天
	Month28Duration = 28 * DayDuration

	// 月，表示29天，等于28天加1天
	Month29Duration = Month28Duration + 1*DayDuration

	// 月，表示30天，等于28天加2天
	Month30Duration = Month28Duration + 2*DayDuration

	// 年，表示365天（平年）
	Year365Duration = 365 * DayDuration

	// 年，表示366天（闰年），等于365天加1天
	Year366Duration = Year365Duration + 1*DayDuration
)
