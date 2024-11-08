/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 10:50:50
 * @FilePath: \go-toolbox\tests\convert_radix_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/convert"
)

func TestHexToBytes(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"AABB", []byte{0xAA, 0xBB}},
		{"", []byte{}}, // empty input should return empty byte array
		{"GG", nil},    // invalid hex should return error
	}

	for _, test := range tests {
		result, err := convert.HexToBytes(test.input)
		if test.expected == nil && err == nil {
			t.Errorf("HexToBytes(%s) should return an error", test.input)
			continue
		} else if test.expected != nil && err != nil {
			t.Errorf("HexToBytes(%s) returned an error: %v", test.input, err)
			continue
		}
		if !equalBytes(result, test.expected) {
			t.Errorf("HexToBytes(%s) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestHexBCC(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"AABBCC", "dd"},
		{"aabbcc", "dd"},
		{"", "00"}, // empty input should return BCC as 0
	}

	for _, test := range tests {
		result, err := convert.HexToBCC(test.input)
		if err != nil {
			t.Errorf("HexBCC(%s) returned an error: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("HexBCC(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestDecToHex(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{255, "FF"},
		{0, "00"},
		{4095, "0FFF"},
	}

	for _, test := range tests {
		result := convert.DecToHex(test.input)
		if result != test.expected {
			t.Errorf("DecToHex(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestHexToDec(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"FF", 255},
		{"0FFF", 4095},
		{"00", 0},
	}

	for _, test := range tests {
		result, err := convert.HexToDec(test.input)
		if err != nil {
			t.Errorf("HexToDec(%s) returned an error: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("HexToDec(%s) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestDecToBin(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "00000000"},
		{1, "00000001"},
		{255, "11111111"},
	}

	for _, test := range tests {
		result := convert.DecToBin(test.input)
		if result != test.expected {
			t.Errorf("DecToBin(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestHexToBin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FF", "11111111"},
		{"0", "00000000"},
	}

	for _, test := range tests {
		result, err := convert.HexToBin(test.input)
		if err != nil {
			t.Errorf("HexToBin(%s) returned an error: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("HexToBin(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
