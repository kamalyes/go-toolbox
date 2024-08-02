/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 09:36:36
 * @FilePath: \go-middleware\moment\date_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import (
	"testing"
	"time"
)

func TestDateAllFunctions(t *testing.T) {
	t.Run("TestYear", TestYear)
	t.Run("TestMonth", TestMonth)
	t.Run("TestDay", TestDay)
	t.Run("TestYearDay", TestYearDay)
	t.Run("TestYearDefault", TestYearDefault)
	t.Run("TestMonthDefault", TestMonthDefault)
	t.Run("TestDayDefault", TestDayDefault)
}

func TestYear(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	year := Year(testTime)
	expectedYear := 2024
	if year != expectedYear {
		t.Errorf("Year() returned %d, expected %d", year, expectedYear)
	}
}

func TestMonth(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	month := Month(testTime)
	expectedMonth := 4
	if month != expectedMonth {
		t.Errorf("Month() returned %d, expected %d", month, expectedMonth)
	}
}

func TestDay(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	day := Day(testTime)
	expectedDay := 2
	if day != expectedDay {
		t.Errorf("Day() returned %d, expected %d", day, expectedDay)
	}
}

func TestYearDay(t *testing.T) {
	testTime := time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC)
	yearDay := YearDay(testTime)
	expectedYearDay := 93
	if yearDay != expectedYearDay {
		t.Errorf("YearDay() returned %d, expected %d", yearDay, expectedYearDay)
	}
}

func TestYearDefault(t *testing.T) {
	year := Year()
	currentYear := time.Now().Year()
	if year != currentYear {
		t.Errorf("YearDefault() returned %d, expected %d", year, currentYear)
	}
}

func TestMonthDefault(t *testing.T) {
	month := Month()
	currentMonth := int(time.Now().Month())
	if month != currentMonth {
		t.Errorf("MonthDefault() returned %d, expected %d", month, currentMonth)
	}
}

func TestDayDefault(t *testing.T) {
	day := Day()
	currentDay := int(time.Now().Day())
	if day != currentDay {
		t.Errorf("DayDefault() returned %d, expected %d", day, currentDay)
	}
}
