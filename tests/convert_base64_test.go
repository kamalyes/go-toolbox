/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 10:56:16
 * @FilePath: \go-toolbox\tests\convert_base64_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/convert"
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
		result, err := convert.B64Encode(test.input)
		if err != nil {
			t.Fatal(err)
		}
		if result != test.expected {
			t.Errorf("B64Encode(%s) = %s; want %s", test.input, result, test.expected)
		}
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
		result, err := convert.B64Decode(test.input)
		if err != nil {
			t.Errorf("B64Decode(%s) returned an error: %v", test.input, err)
			continue
		}
		if !equalBytes(result, test.expected) {
			t.Errorf("B64Decode(%s) = %v; want %v", test.input, result, test.expected)
		}
	}
}
