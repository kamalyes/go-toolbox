/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 00:38:34
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 01:28:32
 * @FilePath: \go-toolbox\conv\byte_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexAllFunctions(t *testing.T) {
	t.Run("TestBytesToHex", TestBytesToHex)
	t.Run("TestHexToBytes", TestHexToBytes)
	t.Run("TestHexBCC", TestHexBCC)
	t.Run("TestBytesBCC", TestBytesBCC)
	t.Run("TestConversions", TestConversions)
}
func TestBytesToHex(t *testing.T) {
	testCases := map[string]struct {
		expectedHexStr string
		hexBytes       []byte
	}{
		"Test 1": {hexBytes: []byte{0xAA, 0xBB, 0xCC}, expectedHexStr: "AABBCC"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := BytesToHex(tc.hexBytes)
			assert.Equal(t, tc.expectedHexStr, result, "Expected hexadecimal string does not match")
		})
	}
}

func TestHexToBytes(t *testing.T) {
	testCases := map[string]struct {
		hexStr           string
		expectedHexBytes []byte
	}{
		"Test 2": {hexStr: "AABBCC", expectedHexBytes: []byte{0xAA, 0xBB, 0xCC}},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := HexToBytes(tc.hexStr)
			assert.Equal(t, tc.expectedHexBytes, result, "Expected hexadecimal string does not match")
		})
	}
}

func TestHexBCC(t *testing.T) {
	testCases := map[string]struct {
		hexStr      string
		expectedBCC string
	}{
		"Test 3": {hexStr: "112233", expectedBCC: "00"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := HexBCC(tc.hexStr)
			assert.Equal(t, tc.expectedBCC, result, "Expected BCC value does not match")
		})
	}
}

func TestBytesBCC(t *testing.T) {
	testCases := []struct {
		name        string
		inputBytes  []byte
		expectedBCC byte
	}{
		{"Test 4", []byte{0xAA, 0xBB, 0xCC}, byte(0xdd)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := BytesBCC(tc.inputBytes)
			assert.Equal(t, tc.expectedBCC, result, "Unexpected BCC value")
		})
	}
}

func TestConversions(t *testing.T) {
	assert.Equal(t, "2A", DecToHex(42), "Decimal to Hexadecimal conversion failed")
	assert.Equal(t, uint64(42), HexToDec("2A"), "Hexadecimal to Decimal conversion failed")
	assert.Equal(t, "00101010", DecToBin(42), "Decimal to Binary conversion failed")
	assert.Equal(t, "00101010", HexToBin("2A"), "Hexadecimal to Binary conversion failed")
	assert.Equal(t, "00101010", ByteToBinStr(byte(42)), "Byte to Binary conversion failed")
	assert.Equal(t, "00100001001000100011001101000100", BytesToBinStr([]byte{0x21, 0x22, 0x33, 0x44}), "Bytes to Binary string conversion failed")
	assert.Equal(t, "00100001-00100010-00110011-01000100", BytesToBinStrWithSplit([]byte{0x21, 0x22, 0x33, 0x44}, "-"), "Bytes to Binary string with split failed")
}
