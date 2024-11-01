/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-17 15:05:00
 * @FilePath: \go-toolbox\tests\date_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/moment"
)

func TestDateAllFunctions(t *testing.T) {
	t.Run("TestParseYear", TestParseYear)
	t.Run("TestParseMonth", TestParseMonth)
	t.Run("TestParseDay", TestParseDay)
	t.Run("TestParseYearDay", TestParseYearDay)
	t.Run("TestParseYearDefault", TestParseYearDefault)
	t.Run("TestParseMonthDefault", TestParseMonthDefault)
	t.Run("TestParseDayDefault", TestParseDayDefault)
}

func TestParseYear(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	year := moment.ParseYear(testTime)
	expectedYear := 2024
	if year != expectedYear {
		t.Errorf("ParseYear() returned %d, expected %d", year, expectedYear)
	}
}

func TestParseMonth(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	month := moment.ParseMonth(testTime)
	expectedMonth := 4
	if month != expectedMonth {
		t.Errorf("ParseMonth() returned %d, expected %d", month, expectedMonth)
	}
}

func TestParseDay(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	day := moment.ParseDay(testTime)
	expectedDay := 2
	if day != expectedDay {
		t.Errorf("ParseDay() returned %d, expected %d", day, expectedDay)
	}
}

func TestParseYearDay(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	yearDay := moment.ParseYearDay(testTime)
	expectedYearDay := 93
	if yearDay != expectedYearDay {
		t.Errorf("ParseYearDay() returned %d, expected %d", yearDay, expectedYearDay)
	}
}

func TestParseYearDefault(t *testing.T) {
	year := moment.ParseYear()
	currentYear := time.Now().Year()
	if year != currentYear {
		t.Errorf("YearDefault() returned %d, expected %d", year, currentYear)
	}
}

func TestParseMonthDefault(t *testing.T) {
	month := moment.ParseMonth()
	currentMonth := int(time.Now().Month())
	if month != currentMonth {
		t.Errorf("MonthDefault() returned %d, expected %d", month, currentMonth)
	}
}

func TestParseDayDefault(t *testing.T) {
	day := moment.ParseDay()
	currentDay := int(time.Now().Day())
	if day != currentDay {
		t.Errorf("DayDefault() returned %d, expected %d", day, currentDay)
	}
}
