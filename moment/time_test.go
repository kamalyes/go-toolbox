/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-17 15:05:55
 * @FilePath: \go-toolbox\moment\time_test.go
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
	t.Run("String", TestParseString)
	t.Run("Hour", TestParseHour)
	t.Run("Minute", TestParseMinute)
	t.Run("Second", TestParseSecond)
	t.Run("Timestamp", TestParseTimestamp)
	t.Run("Milliseconds", TestCurrentMillisecond)
	t.Run("Microsecond", TestCurrentMicrosecond)
	t.Run("Nanosecond", TestCurrentNanosecond)
	t.Run("GmtTime", TestCurrentGmtTime)
	t.Run("LocalTime", TestLocalTime)
	t.Run("Strtotime", TestStrtotime)
	t.Run("CharToCode", TestCharToCode)
}

func TestParseString(t *testing.T) {
	result := ParseString()
	if len(result) == 0 {
		t.Errorf("ParseString() returned an empty string")
	}
}

func TestParseHour(t *testing.T) {
	now := time.Now()
	hour := ParseHour(now)
	if hour < 0 || hour > 23 {
		t.Errorf("ParseHour() returned an invalid hour value: %d", hour)
	}
}

func TestParseMinute(t *testing.T) {
	now := time.Now()
	minute := ParseMinute(now)
	if minute < 0 || minute > 59 {
		t.Errorf("ParseMinute() returned an invalid minute value: %d", minute)
	}
}

func TestParseSecond(t *testing.T) {
	now := time.Now()
	second := ParseSecond(now)
	if second < 0 || second > 59 {
		t.Errorf("ParseSecond() returned an invalid second value: %d", second)
	}
}

func TestParseTimestamp(t *testing.T) {
	timestamp := ParseTimestamp("2024-01-01 00:00:00")
	if timestamp <= 0 {
		t.Errorf("ParseTimestamp() returned an invalid timestamp value: %d", timestamp)
	}
}

func TestCurrentMillisecond(t *testing.T) {
	ms := CurrentMillisecond()
	if ms <= 0 {
		t.Errorf("CurrentMillisecond() returned an invalid millisecond value: %d", ms)
	}
}

func TestCurrentMicrosecond(t *testing.T) {
	micros := CurrentMicrosecond()
	if micros <= 0 {
		t.Errorf("CurrentMicrosecond() returned an invalid microsecond value: %d", micros)
	}
}

func TestCurrentNanosecond(t *testing.T) {
	nanos := CurrentNanosecond()
	if nanos <= 0 {
		t.Errorf("CurrentNanosecond() returned an invalid nanosecond value: %d", nanos)
	}
}

func TestCurrentGmtTime(t *testing.T) {
	gmtTime := CurrentGmtTime()
	if len(gmtTime) == 0 {
		t.Errorf("CurrentGmtTime() returned an empty string")
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
