/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: root 501893067@qq.com
 * @LastEditTime: 2025-02-13 15:20:08
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

	DefaultTimeFormat = "2006-01-02 15:04:05"

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

	// 年，表示1年持续的纳秒时长
	YearDuration = 365 * DayDuration
)
