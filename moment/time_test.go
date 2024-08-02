/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 09:36:36
 * @FilePath: \go-middleware\moment\time_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import (
	"testing"
	"time"
)

func TestMomentFunctions(t *testing.T) {
	t.Run("String", TestString)
	t.Run("Hour", TestHour)
	t.Run("Minute", TestMinute)
	t.Run("Second", TestSecond)
	t.Run("Timestamp", TestTimestamp)
	t.Run("Milliseconds", TestMilliseconds)
	t.Run("Microsecond", TestMicrosecond)
	t.Run("Nanosecond", TestNanosecond)
	t.Run("GmtTime", TestGmtTime)
	t.Run("LocalTime", TestLocalTime)
	t.Run("Strtotime", TestStrtotime)
	t.Run("CharToCode", TestCharToCode)
}

func TestString(t *testing.T) {
	result := String()
	if len(result) == 0 {
		t.Errorf("String() returned an empty string")
	}
}

func TestHour(t *testing.T) {
	now := time.Now()
	hour := Hour(now)
	if hour < 0 || hour > 23 {
		t.Errorf("Hour() returned an invalid hour value: %d", hour)
	}
}

func TestMinute(t *testing.T) {
	now := time.Now()
	minute := Minute(now)
	if minute < 0 || minute > 59 {
		t.Errorf("Minute() returned an invalid minute value: %d", minute)
	}
}

func TestSecond(t *testing.T) {
	now := time.Now()
	second := Second(now)
	if second < 0 || second > 59 {
		t.Errorf("Second() returned an invalid second value: %d", second)
	}
}

func TestTimestamp(t *testing.T) {
	timestamp := Timestamp("2024-01-01 00:00:00")
	if timestamp <= 0 {
		t.Errorf("Timestamp() returned an invalid timestamp value: %d", timestamp)
	}
}

func TestMilliseconds(t *testing.T) {
	ms := Millisecond()
	if ms <= 0 {
		t.Errorf("Millisecond() returned an invalid millisecond value: %d", ms)
	}
}

func TestMicrosecond(t *testing.T) {
	micros := Microsecond()
	if micros <= 0 {
		t.Errorf("Microsecond() returned an invalid microsecond value: %d", micros)
	}
}

func TestNanosecond(t *testing.T) {
	nanos := Nanosecond()
	if nanos <= 0 {
		t.Errorf("Nanosecond() returned an invalid nanosecond value: %d", nanos)
	}
}

func TestGmtTime(t *testing.T) {
	gmtTime := GmtTime()
	if len(gmtTime) == 0 {
		t.Errorf("GmtTime() returned an empty string")
	}
}

func TestLocalTime(t *testing.T) {
	localTime := LocalTime()
	if len(localTime) == 0 {
		t.Errorf("LocalTime() returned an empty string")
	}
}

func TestStrtotime(t *testing.T) {
	testString := "2024-12-25 23:59:59"
	result, err := Strtotime(testString)
	if err != nil {
		t.Errorf("Strtotime() error: %v", err)
	}

	expectedDate := time.Date(2024, time.December, 25, 23, 59, 59, 0, time.Local)
	if !result.Equal(expectedDate) {
		t.Errorf("Strtotime() did not parse the string correctly")
	}
}

func TestCharToCode(t *testing.T) {
	layout := "Y-m-d H:i:s"
	expectedLayout := "2006-1-2 15:4:5"
	result := CharToCode(layout)
	if result != expectedLayout {
		t.Errorf("CharToCode() returned %s, expected %s", result, expectedLayout)
	}
}
