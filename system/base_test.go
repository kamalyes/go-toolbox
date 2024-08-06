/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-06 15:57:30
 * @FilePath: \go-toolbox\system\base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeGetHostName(t *testing.T) {
	actual := SafeGetHostName()
	assert.NotEmpty(t, actual, "HostNames should match")
}

func TestServerIP(t *testing.T) {
	externalIP, internalIP := ServerIP()
	assert.NotEmpty(t, externalIP, "External IP should not be empty, got %s", externalIP)
	assert.NotEmpty(t, internalIP, "Internal IP should not be empty, got %s", internalIP)
}
