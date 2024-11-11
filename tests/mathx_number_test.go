/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-12 21:51:15
 * @FilePath: \go-toolbox\tests\mathx_NUMBER_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

func TestMathxAtLeast(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, lower, expected int
	}{
		{5, 3, 5},
		{2, 10, 10},
		{-1, -5, -1},
	}
	for _, tt := range tests {
		result := mathx.AtLeast(tt.x, tt.lower)
		assert.Equal(tt.expected, result)
	}
}

func TestMathxAtMost(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, upper, expected int
	}{
		{8, 6, 6},
		{4, 2, 2},
		{-3, -7, -7},
	}
	for _, tt := range tests {
		result := mathx.AtMost(tt.x, tt.upper)
		assert.Equal(tt.expected, result)
	}
}

func TestMathxBetween(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, lower, upper, expected int
	}{
		{4, 1, 10, 4},
		{0, -5, 5, 0},
		{12, 1, 10, 10},
		{-8, -5, -3, -5},
	}
	for _, tt := range tests {
		result := mathx.Between(tt.x, tt.lower, tt.upper)
		assert.Equal(tt.expected, result)
	}
}

func TestAtLeastFloat64(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, lower, expected float64
	}{
		{5.5, 3.3, 5.5},
		{2.2, 10.0, 10.0},
		{-1.1, -5.5, -1.1},
	}
	for _, tt := range tests {
		result := mathx.AtLeast(tt.x, tt.lower)
		assert.Equal(tt.expected, result)
	}
}

func TestBetweenFloat64(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, lower, upper, expected float64
	}{
		{4.4, 1.1, 10.0, 4.4},
		{0.0, -5.5, 5.5, 0.0},
		{12.3, 1.0, 10.0, 10.0},
		{-8.8, -5.5, -3.3, -5.5},
	}
	for _, tt := range tests {
		result := mathx.Between(tt.x, tt.lower, tt.upper)
		assert.Equal(tt.expected, result)
	}
}

func TestZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		expected any
	}{
		{"int", 0},
		{"int8", int8(0)},
		{"int16", int16(0)},
		{"int32", int32(0)},
		{"int64", int64(0)},
		{"uint", uint(0)},
		{"uint8", uint8(0)},
		{"uint16", uint16(0)},
		{"uint32", uint32(0)},
		{"uint64", uint64(0)},
		{"float32", float32(0.0)},
		{"float64", float64(0.0)},
		{"string", ""},
		{"bool", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result any
			switch tt.expected.(type) {
			case int:
				result = mathx.ZeroValue[int]()
			case int8:
				result = mathx.ZeroValue[int8]()
			case int16:
				result = mathx.ZeroValue[int16]()
			case int32:
				result = mathx.ZeroValue[int32]()
			case int64:
				result = mathx.ZeroValue[int64]()
			case uint:
				result = mathx.ZeroValue[uint]()
			case uint8:
				result = mathx.ZeroValue[uint8]()
			case uint16:
				result = mathx.ZeroValue[uint16]()
			case uint32:
				result = mathx.ZeroValue[uint32]()
			case uint64:
				result = mathx.ZeroValue[uint64]()
			case float32:
				result = mathx.ZeroValue[float32]()
			case float64:
				result = mathx.ZeroValue[float64]()
			case string:
				result = mathx.ZeroValue[string]()
			case bool:
				result = mathx.ZeroValue[bool]()
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
