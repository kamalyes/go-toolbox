/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-26 13:37:12
 * @FilePath: \go-toolbox\pkg\osx\goos_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package osx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMac(t *testing.T) {
	originalGoos := GetGOOS                   // 保存原始 GetGOOS
	defer func() { GetGOOS = originalGoos }() // 测试结束后恢复原始 GetGOOS

	// 测试 macOS
	GetGOOS = func() string { return OSMac }
	assert.True(t, IsMac(), "Expected IsMac() to return true for %s", OSMac)

	// 测试其他操作系统
	GetGOOS = func() string { return OSWindows }
	assert.False(t, IsMac(), "Expected IsMac() to return false for %s", OSWindows)

	GetGOOS = func() string { return OSLinux }
	assert.False(t, IsMac(), "Expected IsMac() to return false for %s", OSLinux)
}

func TestIsWindows(t *testing.T) {
	originalGoos := GetGOOS                   // 保存原始 GetGOOS
	defer func() { GetGOOS = originalGoos }() // 测试结束后恢复原始 GetGOOS

	// 测试 Windows
	GetGOOS = func() string { return OSWindows }
	assert.True(t, IsWindows(), "Expected IsWindows() to return true for %s", IsWindows())

	// 测试其他操作系统
	GetGOOS = func() string { return OSMac }
	assert.False(t, IsWindows(), "Expected IsWindows() to return false for %s", OSMac)

	GetGOOS = func() string { return OSLinux }
	assert.False(t, IsWindows(), "Expected IsWindows() to return false for %s", OSLinux)
}

func TestIsLinux(t *testing.T) {
	originalGoos := GetGOOS                   // 保存原始 GetGOOS
	defer func() { GetGOOS = originalGoos }() // 测试结束后恢复原始 GetGOOS

	// 测试 Linux
	GetGOOS = func() string { return OSLinux }
	assert.True(t, IsLinux(), "Expected IsLinux() to return true for %s", OSLinux)

	// 测试其他操作系统
	GetGOOS = func() string { return OSMac }
	assert.False(t, IsLinux(), "Expected IsLinux() to return false for %s", OSMac)

	GetGOOS = func() string { return OSWindows }
	assert.False(t, IsLinux(), "Expected IsLinux() to return false for %s", OSWindows)
}

func TestIsSupportedOS(t *testing.T) {
	originalGoos := GetGOOS                   // 保存原始 GetGOOS
	defer func() { GetGOOS = originalGoos }() // 测试结束后恢复原始 GetGOOS

	// 测试支持的操作系统
	supportedOS := []string{OSMac, OSWindows, OSLinux}
	for _, os := range supportedOS {
		GetGOOS = func() string { return os }
		assert.True(t, IsSupportedOS(), "Expected IsSupportedOS() to return true for %s", os)
	}

	// 测试不支持的操作系统
	GetGOOS = func() string { return "unsupported_os" }
	assert.False(t, IsSupportedOS(), "Expected IsSupportedOS() to return false for unsupported OS")
}
