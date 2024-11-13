/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 18:55:56
 * @FilePath: \go-toolbox\pkg\moment\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import "time"

// TimeDifference 结构体用于存储年、天、小时、分钟和秒
type TimeDifference struct {
	Years   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// CalculateTimeDifference 计算给定的 time.Duration，并返回 TimeDifference
func CalculateTimeDifference(duration time.Duration) TimeDifference {
	totalSeconds := int(duration.Seconds())

	// 使用常量计算年、天、小时、分钟和秒
	years := totalSeconds / int(Year.Seconds())
	days := (totalSeconds / int(Day.Seconds())) % 365
	hours := (totalSeconds / int(Hour.Seconds())) % 24
	minutes := (totalSeconds / int(Minute.Seconds())) % 60
	seconds := totalSeconds % 60

	return TimeDifference{
		Years:   years,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}

// SafeParseTimeToUnixNano
func SafeParseTimeToUnixNano(timeStr string) int64 {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0
	}
	return parsedTime.UnixNano() / int64(time.Millisecond)
}

// GetCurrentTimeInfo 获取当前日期、小时和时间的通用函数
func GetCurrentTimeInfo() (string, int, time.Time) {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02"), currentTime.Hour(), currentTime
}

// GetServerTimezone 获取服务器的本地时区信息
func GetServerTimezone() string {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		timezone := DefaultTimezone
		return timezone
	}
	return loc.String()
}

// GetTimeOffset 国际化时间戳偏移
func GetTimeOffset(timezone string, ts int64) (offset int) {
	var loc, _ = time.LoadLocation(timezone)
	_, offset = time.Unix(ts, 0).In(loc).Zone()
	return
}

// FormatWithLocation 国际化时间戳转换字符串
func FormatWithLocation(timezone string, ts int64, defaultDateFormat string) string {
	lt, _ := time.LoadLocation(timezone)
	str := time.Unix(ts, 0).In(lt).Format(defaultDateFormat)
	return str
}

// ParseWithLocation 国际化时间字符串转换时间戳
func ParseWithLocation(timezone string, timeStr string, defaultDateFormat string) int64 {
	l, _ := time.LoadLocation(timezone)
	lt, _ := time.ParseInLocation(defaultDateFormat, timeStr, l)
	return lt.Unix()
}

// ConvertStringToTimestamp String时间类型转换为时间戳
func ConvertStringToTimestamp(dateString, layout string, timeZone string) (int64, error) {
	// 加载时区
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return 0, err
	}

	// 将字符串解析为时间
	t, err := time.ParseInLocation(layout, dateString, loc)
	if err != nil {
		return 0, err
	}
	// 转换时间为时间戳
	timestamp := t.Unix()
	return timestamp, nil
}
