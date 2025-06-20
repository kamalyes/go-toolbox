/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-11 18:16:16
 * @FilePath: \go-toolbox\pkg\moment\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import (
	"fmt"
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
	years := totalSeconds / int(Year366Duration.Seconds())
	days := (totalSeconds / int(DayDuration.Seconds())) % 366
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

// DaysInMonth 计算指定年月的天数
func DaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
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

// DayOfYear 获取指定年份的第n天，时分秒归零
func DayOfYear(year int, n int) time.Time {
	// 先获取该年1月1日
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	// 加上n-1天
	return firstDay.AddDate(0, 0, n-1)
}

// LastDayOfMonth 获取当月最后一天，时分秒归零
func LastDayOfMonth(year int, month time.Month) time.Time {
	// 当月第一天
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	// 下个月第一天
	nextMonth := firstDay.AddDate(0, 1, 0)
	// 下个月第一天减一天即当月最后一天
	lastDay := nextMonth.AddDate(0, 0, -1)
	return lastDay
}

// LastWeekdayOfMonth 获取当月最后一个指定星期几（如最后一个周五）
// 返回时间，时分秒归零
func LastWeekdayOfMonth(year int, month time.Month, target time.Weekday) time.Time {
	// 获取当月最后一天
	lastDay := LastDayOfMonth(year, month)

	// lastDay的星期几
	wd := lastDay.Weekday()

	// 计算距离目标星期几需要往前推几天
	// 例如：lastDay是周日(0)，目标是周五(5)，往前推2天
	daysBack := (int(wd) - int(target) + 7) % 7

	// 往前推daysBack天
	lastTargetDay := lastDay.AddDate(0, 0, -daysBack)

	return lastTargetDay
}

// NextWorkDay 返回指定日期的下一个工作日（周一至周五）
// 规则说明：
// - 输入日期是周一到周四，返回下一天（一定是工作日）
// - 输入日期是周五，返回下周一（跳过周六、周日）
// - 输入日期是周六，返回下周一（加2天）
// - 输入日期是周日，返回下周一（加1天）
// 跨月跨年由 time.AddDate 自动处理，无需额外判断
func NextWorkDay(year int, month time.Month, day int) time.Time {
	// 构造时间对象，时分秒设为0，使用本地时区
	t := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	// 定义一个数组，索引为 Weekday 值（0=周日，1=周一，...，6=周六）
	// 数组值表示从该日期到下一个工作日需要加的天数
	// 解释：
	// 周日(0)  -> +1 天，变成周一
	// 周一(1)  -> +1 天，变成周二
	// 周二(2)  -> +1 天，变成周三
	// 周三(3)  -> +1 天，变成周四
	// 周四(4)  -> +1 天，变成周五
	// 周五(5)  -> +3 天，跳过周六、周日，变成下周一
	// 周六(6)  -> +2 天，跳过周日，变成下周一
	daysToAddMap := make(map[time.Weekday]int)

	// 先把周日、周一到周四全部赋值为1
	for _, day := range []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday} {
		daysToAddMap[day] = 1
	}

	// 单独赋值剩下的周五、周六
	daysToAddMap[time.Friday] = 3
	daysToAddMap[time.Saturday] = 2

	// 获取输入日期的星期几（0-6）
	wd := t.Weekday()

	// 根据星期几从数组中取对应的加天数
	daysToAdd := daysToAddMap[wd]

	// 返回加上对应天数后的日期，即下一个工作日
	return t.AddDate(0, 0, daysToAdd)
}

// NextWeekday 返回指定日期的下一个目标星期几
// 如果当天是目标星期几，则返回下一周的同一天（加7天）
func NextWeekday(year int, month time.Month, day int, target time.Weekday) time.Time {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	wd := t.Weekday()

	daysToAdd := (int(target) - int(wd) + 7) % 7
	if daysToAdd == 0 {
		daysToAdd = 7
	}

	return t.AddDate(0, 0, daysToAdd)
}

// HumanDuration 计算两个时间点之间的年月日时分秒差距，
// 并返回一个中文格式的字符串，比如“1年零2个月零3天零4小时零5分零6秒”
// start 和 end 可任意顺序，函数内部会自动调整顺序
func HumanDuration(start, end time.Time) string {
	// 如果 end 时间早于 start，交换两者，确保 start <= end
	if end.Before(start) {
		start, end = end, start
	}

	// 计算各个时间单位的初步差值
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	days := end.Day() - start.Day()
	hours := end.Hour() - start.Hour()
	minutes := end.Minute() - start.Minute()
	seconds := end.Second() - start.Second()

	// 处理秒的借位，比如秒为负数，向分钟借1分钟（60秒）
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	// 处理分钟的借位
	if minutes < 0 {
		minutes += 60
		hours--
	}
	// 处理小时的借位
	if hours < 0 {
		hours += 24
		days--
	}
	// 处理天的借位
	if days < 0 {
		// 计算上个月的天数，用于借位
		previousMonth := end.AddDate(0, -1, 0)
		days += DaysInMonth(previousMonth.Year(), previousMonth.Month())
		months--
	}
	// 处理月的借位
	if months < 0 {
		months += 12
		years--
	}

	// 定义一个结构体，便于统一处理单位和值
	type unit struct {
		value int
		name  string
	}
	units := []unit{
		{years, "年"},
		{months, "个月"},
		{days, "天"},
		{hours, "小时"},
		{minutes, "分"},
		{seconds, "秒"},
	}

	// 找第一个非零单位和最后一个非零单位索引
	firstIdx, lastIdx := -1, -1
	for i, u := range units {
		if u.value > 0 {
			if firstIdx == -1 {
				firstIdx = i
			}
			lastIdx = i
		}
	}
	// 全零返回"0秒"
	if firstIdx == -1 {
		return "0秒"
	}

	// 构造结果字符串
	result := ""
	prevNonZeroIdx := -1 // 记录上一个非零单位索引
	for i := firstIdx; i <= lastIdx; i++ {
		if units[i].value > 0 {
			// 如果当前非零单位和上一个非零单位不相邻，中间有零单位，插入“零”
			if prevNonZeroIdx != -1 && i-prevNonZeroIdx > 1 {
				result += "零"
			}
			result += fmt.Sprintf("%d%s", units[i].value, units[i].name)
			prevNonZeroIdx = i
		}
	}

	return result
}
