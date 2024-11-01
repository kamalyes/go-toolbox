/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-17 15:59:36
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 01:37:37
 * @FilePath: \go-toolbox\tests\moment_constants_test.go
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
		{"Nanosecond", nanosecond, moment.Nanosecond},
		{"Microsecond", microsecond, moment.Microsecond},
		{"Millisecond", millisecond, moment.Millisecond},
		{"Second", second, moment.Second},
		{"Minute", minute, moment.Minute},
		{"Hour", hour, moment.Hour},
		{"Day", day, moment.Day},
		{"Week", week, moment.Week},
		{"Year", year, moment.Year},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("%s expected: %d, got: %d", test.name, test.expected, test.actual)
		}
	}
}
