/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-26 13:37:12
 * @FilePath: \go-toolbox\tests\osx_goos_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/stretchr/testify/assert"
)

func TestIsMac(t *testing.T) {
	originalGoos := osx.GetGOOS                   // 保存原始 osx.GetGOOS
	defer func() { osx.GetGOOS = originalGoos }() // 测试结束后恢复原始 osx.GetGOOS

	// 测试 macOS
	osx.GetGOOS = func() string { return osx.OSMac }
	assert.True(t, osx.IsMac(), "Expected IsMac() to return true for %s", osx.OSMac)

	// 测试其他操作系统
	osx.GetGOOS = func() string { return osx.OSWindows }
	assert.False(t, osx.IsMac(), "Expected IsMac() to return false for %s", osx.OSWindows)

	osx.GetGOOS = func() string { return osx.OSLinux }
	assert.False(t, osx.IsMac(), "Expected IsMac() to return false for %s", osx.OSLinux)
}

func TestIsWindows(t *testing.T) {
	originalGoos := osx.GetGOOS                   // 保存原始 osx.GetGOOS
	defer func() { osx.GetGOOS = originalGoos }() // 测试结束后恢复原始 osx.GetGOOS

	// 测试 Windows
	osx.GetGOOS = func() string { return osx.OSWindows }
	assert.True(t, osx.IsWindows(), "Expected IsWindows() to return true for %s", osx.IsWindows())

	// 测试其他操作系统
	osx.GetGOOS = func() string { return osx.OSMac }
	assert.False(t, osx.IsWindows(), "Expected IsWindows() to return false for %s", osx.OSMac)

	osx.GetGOOS = func() string { return osx.OSLinux }
	assert.False(t, osx.IsWindows(), "Expected IsWindows() to return false for %s", osx.OSLinux)
}

func TestIsLinux(t *testing.T) {
	originalGoos := osx.GetGOOS                   // 保存原始 osx.GetGOOS
	defer func() { osx.GetGOOS = originalGoos }() // 测试结束后恢复原始 osx.GetGOOS

	// 测试 Linux
	osx.GetGOOS = func() string { return osx.OSLinux }
	assert.True(t, osx.IsLinux(), "Expected IsLinux() to return true for %s", osx.OSLinux)

	// 测试其他操作系统
	osx.GetGOOS = func() string { return osx.OSMac }
	assert.False(t, osx.IsLinux(), "Expected IsLinux() to return false for %s", osx.OSMac)

	osx.GetGOOS = func() string { return osx.OSWindows }
	assert.False(t, osx.IsLinux(), "Expected IsLinux() to return false for %s", osx.OSWindows)
}

func TestIsSupportedOS(t *testing.T) {
	originalGoos := osx.GetGOOS                   // 保存原始 osx.GetGOOS
	defer func() { osx.GetGOOS = originalGoos }() // 测试结束后恢复原始 osx.GetGOOS

	// 测试支持的操作系统
	supportedOS := []string{osx.OSMac, osx.OSWindows, osx.OSLinux}
	for _, os := range supportedOS {
		osx.GetGOOS = func() string { return os }
		assert.True(t, osx.IsSupportedOS(), "Expected IsSupportedOS() to return true for %s", os)
	}

	// 测试不支持的操作系统
	osx.GetGOOS = func() string { return "unsupported_os" }
	assert.False(t, osx.IsSupportedOS(), "Expected IsSupportedOS() to return false for unsupported OS")
}
