/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\convert\base64_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestB64Encode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "aGVsbG8="},
		{[]byte("world"), "d29ybGQ="},
		{[]byte(""), ""},
	}

	for _, test := range tests {
		result, err := B64Encode(test.input)
		assert.NoError(t, err, "B64Encode(%s) returned an error", test.input)
		assert.Equal(t, test.expected, result, "B64Encode(%s) = %s; want %s", test.input, result, test.expected)
	}
}

func TestB64Decode(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"aGVsbG8=", []byte("hello")},
		{"d29ybGQ=", []byte("world")},
		{"", []byte{}},
	}

	for _, test := range tests {
		result, err := B64Decode(test.input)
		if err != nil {
			assert.Fail(t, "B64Decode(%s) returned an error: %v", test.input, err)
			continue
		}
		assert.Equal(t, test.expected, result, "B64Decode(%s) = %v; want %v", test.input, result, test.expected)
	}
}
