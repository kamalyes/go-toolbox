/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-17 15:59:36
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:00:20
 * @FilePath: \go-toolbox\pkg\moment\constants_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import (
	"testing"
	"time"
)

func TestTimeUnits(t *testing.T) {
	// Define expected durations
	nanosecond := time.Duration(1)
	microsecond := 1000 * nanosecond
	millisecond := 1000 * microsecond
	second := 1000 * millisecond
	minute := 60 * second
	hour := 60 * minute
	day := 24 * hour
	week := 7 * day
	year := 366 * day

	// Testing the predefined constants
	tests := []struct {
		name     string
		expected time.Duration
		actual   time.Duration
	}{
		{"Nanosecond", nanosecond, NanosecondDuration},
		{"Microsecond", microsecond, MicrosecondDuration},
		{"Millisecond", millisecond, MillisecondDuration},
		{"Second", second, SecondDuration},
		{"Minute", minute, MinuteDuration},
		{"Hour", hour, HourDuration},
		{"Day", day, DayDuration},
		{"Week", week, WeekDuration},
		{"Year", year, Year366Duration},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("%s expected: %d, got: %d", test.name, test.expected, test.actual)
		}
	}
}
