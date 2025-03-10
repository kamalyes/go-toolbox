/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-05 17:55:55
 * @FilePath: \go-toolbox\pkg\moment\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import (
	"strings"
	"time"
)

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
	years := totalSeconds / int(YearDuration.Seconds())
	days := (totalSeconds / int(DayDuration.Seconds())) % 365
	hours := (totalSeconds / int(HourDuration.Seconds())) % 24
	minutes := (totalSeconds / int(MinuteDuration.Seconds())) % 60
	seconds := totalSeconds % 60

	return TimeDifference{
		Years:   years,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}

// CalculateBirthDate 根据年龄计算出生日期
func CalculateBirthDate(age int) time.Time {
	// 获取当前时间
	currentTime := time.Now()

	// 计算出生年份
	birthYear := currentTime.Year() - age

	// 创建出生日期（假设出生日期为1月1日）
	birthDate := time.Date(birthYear, 1, 1, 0, 0, 0, 0, currentTime.Location())

	return birthDate
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

// 计算年龄的函数，currentTime 为计算年龄时的参考时间
func CalculateAge(birthday string, currentTime time.Time) (int, error) {
	// 解析生日字符串为 time.Time 对象
	birthDate, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return 0, err // 如果解析失败，返回错误
	}

	// 计算年龄
	age := currentTime.Year() - birthDate.Year()

	// 检查是否已经过了生日
	if currentTime.YearDay() < birthDate.YearDay() {
		age-- // 如果还没到生日，年龄减一
	}

	return age, nil
}

// DaysInMonth 获取指定年份和月份的天数
func DaysInMonth(month, year int) int {
	if month == 2 {
		if (year%4 == 0 && year%100 != 0) || (year%400 == 0) {
			return 29 // 闰年
		}
		return 28 // 平年
	}
	if month == 4 || month == 6 || month == 9 || month == 11 {
		return 30 // 30天的月份
	}
	return 31 // 31天的月份
}

// NowTime 获取时间的辅助函数
func NowTime(t []time.Time) time.Time {
	if len(t) > 0 {
		return t[0]
	}
	return time.Now()
}

// CalculateStartAndEndTime 根据给定的年、月、日和持续时间计算开始和结束时间
func CalculateStartAndEndTime(year int, month time.Month, day int, duration time.Duration) (int64, int64) {
	startTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local).UnixMilli()
	endTime := time.Now().Add(duration).UnixMilli()
	return startTime, endTime
}

// 获取年份
func Year(t ...time.Time) int {
	return NowTime(t).Year()
}

// 获取月份
func Month(t ...time.Time) int {
	return int(NowTime(t).Month())
}

// 获取日期
func Day(t ...time.Time) int {
	return NowTime(t).Day()
}

// 获取一年中的第几天
func YearDay(t ...time.Time) int {
	return NowTime(t).YearDay()
}

// 今天的开始和结束时间
func Today() (int64, int64) {
	now := time.Now()
	return CalculateStartAndEndTime(now.Year(), now.Month(), now.Day(), 24*time.Hour)
}

// 昨天的开始和结束时间
func Yesterday() (int64, int64) {
	now := time.Now().Add(-24 * time.Hour)
	return CalculateStartAndEndTime(now.Year(), now.Month(), now.Day(), 24*time.Hour)
}

// 最近N天的开始和结束时间
func LastNDays(num int) (int64, int64) {
	now := time.Now()
	start := now.Add(-24 * time.Hour * time.Duration(num-1))
	return CalculateStartAndEndTime(start.Year(), start.Month(), start.Day(), 24*time.Hour)
}

// 最近N个月的开始和结束时间
func LastNMonths(num int) (int64, int64) {
	now := time.Now()
	start := now.AddDate(0, -num, 0)
	return CalculateStartAndEndTime(start.Year(), start.Month(), start.Day(), 30*24*time.Hour)
}

// 最近N周的开始和结束时间
func LastNWeeks(num int) (int64, int64) {
	now := time.Now()
	start := now.AddDate(0, 0, -7*num)
	return CalculateStartAndEndTime(start.Year(), start.Month(), start.Day(), 7*24*time.Hour)
}

// 最近N年的开始和结束时间
func LastNYears(num int) (int64, int64) {
	now := time.Now()
	start := now.AddDate(-num, 0, 0)
	return CalculateStartAndEndTime(start.Year(), start.Month(), start.Day(), 365*24*time.Hour)
}

// 获取小时
func Hour(t ...time.Time) int {
	return NowTime(t).Hour()
}

// 获取分钟
func Minute(t ...time.Time) int {
	return NowTime(t).Minute()
}

// 获取秒
func Second(t ...time.Time) int {
	return NowTime(t).Second()
}

// 获取当前时间的毫秒数
func CurrentMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取当前时间的微秒数
func CurrentMicrosecond() int64 {
	return time.Now().UnixNano() / 1e3
}

// 获取当前时间的纳秒数
func CurrentNanosecond() int64 {
	return time.Now().UnixNano()
}

// 字符串转换为 time.Time
func StrtoTime(s string, format ...string) (time.Time, error) {
	if len(format) > 0 {
		return time.ParseInLocation(strings.TrimSpace(format[0]), s, time.Local)
	}
	return time.ParseInLocation(DefaultTimeFormat, s, time.Local)
}

// 将自定义布局字符替换为 Go 时间格式
func CharToCode(layout string) string {
	characters := []string{
		"y", "06",
		"m", "1",
		"d", "2",
		"Y", "2006",
		"M", "01",
		"D", "02",
		"h", "03",
		"H", "15",
		"i", "4",
		"s", "5",
		"I", "04",
		"S", "05",
		"t", "pm",
		"T", "PM",
	}
	replacer := strings.NewReplacer(characters...)
	return replacer.Replace(layout)
}
