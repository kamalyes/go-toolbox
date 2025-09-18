/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-18 17:55:55
 * @FilePath: \go-toolbox\tests\netx_ip_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/netx"
	"github.com/stretchr/testify/assert"
)

func TestGetLocalInterfaceIPAndExternalIP(t *testing.T) {
	externalIP, internalIP, err := netx.GetLocalInterfaceIPAndExternalIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, externalIP)
	assert.NotEmpty(t, internalIP)
	t.Logf("externalIP %s", externalIP)
	t.Logf("internalIP %s", internalIP)
}

func TestGetLocalInterfaceIPs(t *testing.T) {
	ips, err := netx.GetLocalInterfaceIPs()
	assert.Nil(t, err)
	assert.NotEmpty(t, ips, fmt.Sprintf("Expected at least one global unicast IP, got: %v", ips))
	for _, ip := range ips {
		assert.NotEmpty(t, ip, fmt.Sprintf("Invalid IP address: %s", ip))
	}
}

func TestGetConNetPublicIP(t *testing.T) {
	ip, err := netx.GetConNetPublicIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected public IP, got: %s", ip))
}
