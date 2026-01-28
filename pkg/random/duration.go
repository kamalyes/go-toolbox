/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-28 00:00:00
 * @FilePath: \go-toolbox\pkg\random\duration.go
 * @Description: 随机时间相关工具函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package random

import "time"

// RandDuration 生成指定范围内的随机时间间隔
// 参数:
//   - min: 最小时间间隔
//   - max: 最大时间间隔
//
// 返回:
//   - time.Duration: 随机生成的时间间隔，范围在 [min, max) 之间
//
// 示例:
//
//	// 生成 100ms 到 500ms 之间的随机延迟
//	delay := random.RandDuration(100*time.Millisecond, 500*time.Millisecond)
//
//	// 生成 1s 到 5s 之间的随机延迟
//	delay := random.RandDuration(1*time.Second, 5*time.Second)
func RandDuration(min, max time.Duration) time.Duration {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	diff := max - min
	return min + time.Duration(NewRand().Int63n(int64(diff)))
}

// RandTimeBetween 生成指定时间范围内的随机时间点
// 参数:
//   - start: 开始时间
//   - end: 结束时间
//
// 返回:
//   - time.Time: 随机生成的时间点，范围在 [start, end) 之间
//
// 示例:
//
//	// 生成今天 00:00 到现在之间的随机时间
//	today := time.Now().Truncate(24 * time.Hour)
//	randTime := random.RandTimeBetween(today, time.Now())
//
//	// 生成过去7天内的随机时间
//	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
//	randTime := random.RandTimeBetween(weekAgo, time.Now())
func RandTimeBetween(start, end time.Time) time.Time {
	if end.Before(start) {
		start, end = end, start
	}
	if start.Equal(end) {
		return start
	}
	diff := end.Sub(start)
	randDuration := time.Duration(NewRand().Int63n(int64(diff)))
	return start.Add(randDuration)
}

// RandTimeInPast 生成过去指定时间范围内的随机时间点
// 参数:
//   - duration: 过去的时间范围
//
// 返回:
//   - time.Time: 随机生成的过去时间点
//
// 示例:
//
//	// 生成过去24小时内的随机时间
//	randTime := random.RandTimeInPast(24 * time.Hour)
//
//	// 生成过去30天内的随机时间
//	randTime := random.RandTimeInPast(30 * 24 * time.Hour)
func RandTimeInPast(duration time.Duration) time.Time {
	now := time.Now()
	past := now.Add(-duration)
	return RandTimeBetween(past, now)
}

// RandTimeInFuture 生成未来指定时间范围内的随机时间点
// 参数:
//   - duration: 未来的时间范围
//
// 返回:
//   - time.Time: 随机生成的未来时间点
//
// 示例:
//
//	// 生成未来24小时内的随机时间
//	randTime := random.RandTimeInFuture(24 * time.Hour)
//
//	// 生成未来7天内的随机时间
//	randTime := random.RandTimeInFuture(7 * 24 * time.Hour)
func RandTimeInFuture(duration time.Duration) time.Time {
	now := time.Now()
	future := now.Add(duration)
	return RandTimeBetween(now, future)
}

// RandDate 生成指定日期范围内的随机日期（时间部分为 00:00:00）
// 参数:
//   - startDate: 开始日期
//   - endDate: 结束日期
//
// 返回:
//   - time.Time: 随机生成的日期
//
// 示例:
//
//	// 生成2024年内的随机日期
//	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
//	randDate := random.RandDate(start, end)
func RandDate(startDate, endDate time.Time) time.Time {
	// 截断到日期（去掉时分秒）
	start := startDate.Truncate(24 * time.Hour)
	end := endDate.Truncate(24 * time.Hour)

	if end.Before(start) {
		start, end = end, start
	}
	if start.Equal(end) {
		return start
	}

	days := int(end.Sub(start).Hours() / 24)
	if days == 0 {
		return start
	}

	randDays := NewRand().Intn(days + 1)
	return start.Add(time.Duration(randDays) * 24 * time.Hour)
}

// RandWeekday 生成随机的星期几
// 返回:
//   - time.Weekday: 随机的星期（0=Sunday, 1=Monday, ..., 6=Saturday）
//
// 示例:
//
//	weekday := random.RandWeekday()
//	fmt.Println(weekday) // 输出: Monday, Tuesday, 等等
func RandWeekday() time.Weekday {
	return time.Weekday(NewRand().Intn(7))
}

// RandMonth 生成随机的月份
// 返回:
//   - time.Month: 随机的月份（1=January, 2=February, ..., 12=December）
//
// 示例:
//
//	month := random.RandMonth()
//	fmt.Println(month) // 输出: January, February, 等等
func RandMonth() time.Month {
	return time.Month(NewRand().Intn(12) + 1)
}

// RandHour 生成随机的小时数（0-23）
// 返回:
//   - int: 随机的小时数
//
// 示例:
//
//	hour := random.RandHour()
//	fmt.Println(hour) // 输出: 0-23 之间的数字
func RandHour() int {
	return NewRand().Intn(24)
}

// RandMinute 生成随机的分钟数（0-59）
// 返回:
//   - int: 随机的分钟数
//
// 示例:
//
//	minute := random.RandMinute()
//	fmt.Println(minute) // 输出: 0-59 之间的数字
func RandMinute() int {
	return NewRand().Intn(60)
}

// RandSecond 生成随机的秒数（0-59）
// 返回:
//   - int: 随机的秒数
//
// 示例:
//
//	second := random.RandSecond()
//	fmt.Println(second) // 输出: 0-59 之间的数字
func RandSecond() int {
	return NewRand().Intn(60)
}

// RandTimeOfDay 生成当天的随机时间点
// 返回:
//   - time.Time: 今天的随机时间点
//
// 示例:
//
//	randTime := random.RandTimeOfDay()
//	fmt.Println(randTime.Format("15:04:05")) // 输出: 今天的某个时间，如 14:23:45
func RandTimeOfDay() time.Time {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	return RandTimeBetween(startOfDay, endOfDay)
}

// RandBusinessHour 生成工作时间内的随机时间点（9:00-18:00）
// 返回:
//   - time.Time: 今天工作时间内的随机时间点
//
// 示例:
//
//	randTime := random.RandBusinessHour()
//	fmt.Println(randTime.Format("15:04:05")) // 输出: 09:00:00 到 18:00:00 之间的时间
func RandBusinessHour() time.Time {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	return RandTimeBetween(start, end)
}

// RandUnixTimestamp 生成指定范围内的随机 Unix 时间戳（秒）
// 参数:
//   - minTimestamp: 最小时间戳
//   - maxTimestamp: 最大时间戳
//
// 返回:
//   - int64: 随机的 Unix 时间戳
//
// 示例:
//
//	// 生成2024年内的随机时间戳
//	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
//	end := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
//	timestamp := random.RandUnixTimestamp(start, end)
func RandUnixTimestamp(minTimestamp, maxTimestamp int64) int64 {
	if maxTimestamp == minTimestamp {
		return minTimestamp
	}
	if maxTimestamp < minTimestamp {
		minTimestamp, maxTimestamp = maxTimestamp, minTimestamp
	}
	diff := maxTimestamp - minTimestamp
	return minTimestamp + NewRand().Int63n(diff+1)
}
