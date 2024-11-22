/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-13 11:57:23
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 19:57:23
 * @FilePath: \go-toolbox\tests\stringx_repeat_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		count  int
		output string
	}{
		{
			name:   "Repeat hello 3 times",
			input:  "hello",
			count:  3,
			output: "hellohellohello",
		},
		{
			name:   "Repeat empty string 5 times",
			input:  "",
			count:  5,
			output: "",
		},
		{
			name:   "Repeat string with special chars",
			input:  "!@#",
			count:  2,
			output: "!@#!@#",
		},
		{
			name:   "Repeat single character",
			input:  "a",
			count:  4,
			output: "aaaa",
		},
		{
			name:   "Repeat with count zero",
			input:  "test",
			count:  0,
			output: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringx.Repeat(tt.input, tt.count)
			assert.Equal(t, tt.output, result)
		})
	}
}
