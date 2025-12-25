/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-25 10:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-25 15:00:00
 * @FilePath: \go-toolbox\pkg\cron\schedule.go
 * @Description: Cron 调度逻辑实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cron

import "time"

// Next 计算下次执行时间(支持纳秒精度)
func (s *CronSpecSchedule) Next(t time.Time) time.Time {
	// 转换到调度器时区
	origLocation := t.Location()
	loc := s.Location
	if loc == time.Local {
		loc = t.Location()
	}
	if s.Location != time.Local {
		t = t.In(s.Location)
	}

	// 从下一秒开始(保留纳秒精度)
	t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

	// 标记是否已增加字段
	added := false

	// 最多查找 5 年
	yearLimit := t.Year() + 5

WRAP:
	if t.Year() > yearLimit {
		return time.Time{}
	}

	// 查找匹配的月份
	for !s.matchBit(s.Month, uint(t.Month())) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 1, 0)
		if t.Month() == time.January {
			goto WRAP
		}
	}

	// 查找匹配的日期
	for !s.matchDayOfMonthAndWeek(t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 0, 1)
		if t.Day() == 1 {
			goto WRAP
		}
	}

	// 查找匹配的小时
	for !s.matchBit(s.Hour, uint(t.Hour())) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
		}
		t = t.Add(1 * time.Hour)
		if t.Hour() == 0 {
			goto WRAP
		}
	}

	// 查找匹配的分钟
	for !s.matchBit(s.Minute, uint(t.Minute())) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, loc)
		}
		t = t.Add(1 * time.Minute)
		if t.Minute() == 0 {
			goto WRAP
		}
	}

	// 查找匹配的秒
	for !s.matchBit(s.Second, uint(t.Second())) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, loc)
		}
		t = t.Add(1 * time.Second)
		if t.Second() == 0 {
			goto WRAP
		}
	}

	return t.In(origLocation)
}

// matchBit 检查值是否匹配位集合
func (s *CronSpecSchedule) matchBit(bits uint64, value uint) bool {
	return bits&(1<<value) != 0
}

// matchDayOfMonthAndWeek 检查日期和星期是否匹配
func (s *CronSpecSchedule) matchDayOfMonthAndWeek(t time.Time) bool {
	domMatch := s.matchBit(s.Dom, uint(t.Day()))
	dowMatch := s.matchBit(s.Dow, uint(t.Weekday()))

	// 如果 Dom 或 Dow 设置了星号位，表示是通配符
	domStar := s.Dom&cronStarBit != 0
	dowStar := s.Dow&cronStarBit != 0

	if domStar && dowStar {
		return true
	}
	if domStar {
		return dowMatch
	}
	if dowStar {
		return domMatch
	}
	return domMatch && dowMatch
}

// Next 计算下次执行时间
func (s *CronEverySchedule) Next(t time.Time) time.Time {
	return t.Add(s.Duration)
}
