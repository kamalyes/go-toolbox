/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-17 16:03:19
 * @FilePath: \go-toolbox\convert\convert_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"testing"
	"time"
)

func TestAllConvertFunctions(t *testing.T) {
	t.Run("TestMustString", TestMustString)
	t.Run("TestMustInt", TestMustInt)
	t.Run("TestMustBool", TestMustBool)
	t.Run("TestB64Encode", TestB64Encode)
	t.Run("TestB64Decode", TestB64Decode)
	t.Run("TestHexToBytes", TestHexToBytes)
	t.Run("TestBytesBCC", TestBytesBCC)
	t.Run("TestHexBCC", TestHexBCC)
	t.Run("TestDecToHex", TestDecToHex)
	t.Run("TestHexToDec", TestHexToDec)
	t.Run("TestDecToBin", TestDecToBin)
	t.Run("TestByteToBinStr", TestByteToBinStr)
	t.Run("TestBytesToBinStr", TestBytesToBinStr)
	t.Run("TestBytesToBinStrWithSplit", TestBytesToBinStrWithSplit)
}

func TestMustString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"hello", "hello"},
		{[]byte("world"), "world"},
		{nil, ""},
		{true, "true"},
		{42, "42"},
		{3.14, "3.14"},
		{time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), "2024-01-01T12:00:00Z"},
	}

	for _, test := range tests {
		result := MustString(test.input)
		if result != test.expected {
			t.Errorf("MustString(%v) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestMustInt(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int
	}{
		{"123", 123},
		{123, 123},
		{nil, 0},
		{true, 1},
		{false, 0},
		{3.14, 3},
		{convertToInt, 0},
	}

	for _, test := range tests {
		result, _ := MustInt(test.input)
		if result != test.expected {
			t.Errorf("MustInt(%v) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestMustBool(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{"1", true},
		{"true", true},
		{"false", false},
		{0, false},
		{1, true},
		{nil, false},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		result := MustBool(test.input)
		if result != test.expected {
			t.Errorf("MustBool(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

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
		result := B64Encode(test.input)
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
		result, err := B64Decode(test.input)
		if err != nil {
			t.Errorf("B64Decode(%s) returned an error: %v", test.input, err)
			continue
		}
		if !equalBytes(result, test.expected) {
			t.Errorf("B64Decode(%s) = %v; want %v", test.input, result, test.expected)
		}
	}
}

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

func TestBytesBCC(t *testing.T) {
	tests := []struct {
		input    []byte
		expected byte
	}{
		{[]byte{0x01, 0x02, 0x03}, 0x00},
		{[]byte{0xFF, 0xFF, 0xFF}, 0xFF},
		{[]byte{0x00}, 0x00},
	}

	for _, test := range tests {
		result := BytesBCC(test.input)
		if result != test.expected {
			t.Errorf("BytesBCC(%v) = %v; want %v", test.input, result, test.expected)
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
		result, err := HexBCC(test.input)
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
		result := DecToHex(test.input)
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
		result, err := HexToDec(test.input)
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
		result := DecToBin(test.input)
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
		result, err := HexToBin(test.input)
		if err != nil {
			t.Errorf("HexToBin(%s) returned an error: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("HexToBin(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestByteToBinStr(t *testing.T) {
	tests := []struct {
		input    byte
		expected string
	}{
		{0, "00000000"},
		{1, "00000001"},
		{255, "11111111"},
	}

	for _, test := range tests {
		result := ByteToBinStr(test.input)
		if result != test.expected {
			t.Errorf("ByteToBinStr(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestBytesToBinStr(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{0, 1, 2}, "000000000000000100000010"},
		{[]byte{255}, "11111111"},
		{[]byte{}, ""},
	}

	for _, test := range tests {
		result := BytesToBinStr(test.input)
		if result != test.expected {
			t.Errorf("BytesToBinStr(%v) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestBytesToBinStrWithSplit(t *testing.T) {
	tests := []struct {
		input    []byte
		split    string
		expected string
	}{
		{[]byte{0, 1, 2}, " ", "00000000 00000001 00000010"},
		{[]byte{255}, "", "11111111"},
		{[]byte{}, "-", ""},
	}

	for _, test := range tests {
		result := BytesToBinStrWithSplit(test.input, test.split)
		if result != test.expected {
			t.Errorf("BytesToBinStrWithSplit(%v, %s) = %s; want %s", test.input, test.split, result, test.expected)
		}
	}
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
