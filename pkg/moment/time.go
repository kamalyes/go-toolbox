/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 02:01:27
 * @FilePath: \go-toolbox\pkg\moment\time.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// timestamp
func ParseString(f ...string) string {
	format := "2006-01-02 15:04:05"
	if len(f) > 0 {
		format = f[0]
	}
	return time.Now().Format(format)
}

// 获取小时
func ParseHour(t ...time.Time) int {
	tmp := time.Now()
	if len(t) > 0 {
		tmp = t[0]
	}
	return tmp.Hour()
}

// 获取分钟
func ParseMinute(t ...time.Time) int {
	tmp := time.Now()
	if len(t) > 0 {
		tmp = t[0]
	}
	return tmp.Minute()
}

// 获取秒
func ParseSecond(t ...time.Time) int {
	tmp := time.Now()
	if len(t) > 0 {
		tmp = t[0]
	}
	return tmp.Second()
}

// 字符串转时间戳
func ParseTimestamp(args ...string) int64 {
	var timestamp int64 = 0
	l := len(args)
	if l == 0 {
		timestamp = time.Now().Unix()
	}
	if l > 0 {
		if reflect.TypeOf(args[0]).String() == "string" {
			t, err := StrtoTime(string(args[0]), "2006-01-02 15:04:05")
			if err != nil {
				log.Println("datetime.Timestamp error:", err)
				return -1
			}
			timestamp = t.Unix()
		}
	}
	return timestamp
}

// 毫秒
func CurrentMillisecond() int64 {
	tmp := time.Now().UnixNano()
	return tmp / 1e6
}

// 微秒
func CurrentMicrosecond() int64 {
	tmp := time.Now().UnixNano()
	return tmp / 1e3
}

// 纳秒
func CurrentNanosecond() int64 {
	return time.Now().UnixNano()
}

// GMT TIME
func CurrentGmtTime() string {
	now := time.Now()
	year, mon, day := now.UTC().Date()
	hour, min, sec := now.UTC().Clock()
	zone, _ := now.UTC().Zone()
	return fmt.Sprintf("%d-%d-%d %02d:%02d:%02d %s", year, mon, day, hour, min, sec, zone)
}

// 本地时区（年-月-日 时:分:秒）
func LocalTime() string {
	now := time.Now().Local()
	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	zone, _ := now.Zone()
	return fmt.Sprintf("%d-%d-%d %02d:%02d:%02d %s", year, mon, day, hour, min, sec, zone)
}

// String To time.Time
func StrtoTime(s string, args ...string) (time.Time, error) {
	format := "2006-01-02 15:04:05"
	if len(args) > 0 {
		format = strings.Trim(args[0], " ")
	}
	if len(s) != len(format) {
		return time.Now(), errors.New("String to time: parameter format error")
	}
	return time.ParseInLocation(format, s, time.Local)
}

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
