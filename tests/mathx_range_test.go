/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:15:56
 * @FilePath: \go-toolbox\tests\mathx_range_test.go
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
