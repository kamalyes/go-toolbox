/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-17 15:59:36
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-17 15:59:54
 * @FilePath: \go-toolbox\moment\constants_test.go
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
	year := 365 * day

	// Testing the predefined constants
	tests := []struct {
		name     string
		expected time.Duration
		actual   time.Duration
	}{
		{"Nanosecond", nanosecond, Nanosecond},
		{"Microsecond", microsecond, Microsecond},
		{"Millisecond", millisecond, Millisecond},
		{"Second", second, Second},
		{"Minute", minute, Minute},
		{"Hour", hour, Hour},
		{"Day", day, Day},
		{"Week", week, Week},
		{"Year", year, Year},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("%s expected: %d, got: %d", test.name, test.expected, test.actual)
		}
	}
}
