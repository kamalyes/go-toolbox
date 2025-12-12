/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\convert\radix_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		result, err := HexToBytes(test.input)
		if test.expected == nil {
			assert.Error(t, err, "HexToBytes(%s) should return an error", test.input)
		} else {
			assert.NoError(t, err, "HexToBytes(%s) returned an error: %v", test.input, err)
			assert.Equal(t, test.expected, result, "HexToBytes(%s) = %v; want %v", test.input, result, test.expected)
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
		result, err := HexToBCC(test.input)
		assert.NoError(t, err, "HexBCC(%s) returned an error: %v", test.input, err)
		assert.Equal(t, test.expected, result, "HexBCC(%s) = %s; want %s", test.input, result, test.expected)
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
		result := DecToHex(test.input)
		assert.Equal(t, test.expected, result, "DecToHex(%d) = %s; want %s", test.input, result, test.expected)
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
		result, err := HexToDec(test.input)
		assert.NoError(t, err, "HexToDec(%s) returned an error: %v", test.input, err)
		assert.Equal(t, test.expected, result, "HexToDec(%s) = %d; want %d", test.input, result, test.expected)
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
		result := DecToBin(test.input)
		assert.Equal(t, test.expected, result, "DecToBin(%d) = %s; want %s", test.input, result, test.expected)
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
		result, err := HexToBin(test.input)
		assert.NoError(t, err, "HexToBin(%s) returned an error: %v", test.input, err)
		assert.Equal(t, test.expected, result, "HexToBin(%s) = %s; want %s", test.input, result, test.expected)
	}
}
