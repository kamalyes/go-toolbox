/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-05 18:56:01
 * @FilePath: \go-toolbox\tests\osx_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/stretchr/testify/assert"
)

func TestAllSysBaseFunctions(t *testing.T) {
	t.Run("TestSafeGetHostName", TestSafeGetHostName)
	t.Run("TestHashUnixMicroCipherText", TestHashUnixMicroCipherText)
}

func TestSafeGetHostName(t *testing.T) {
	actual := osx.SafeGetHostName()
	assert.NotEmpty(t, actual, "HostNames should match")
}

// TestHashUnixMicroCipherText 测试 HashUnixMicroCipherText 函数
func TestHashUnixMicroCipherText(t *testing.T) {
	hash1 := osx.HashUnixMicroCipherText()
	hash2 := osx.HashUnixMicroCipherText()

	// 验证生成的哈希值不为空
	assert.NotEqual(t, hash1, "")
	assert.NotEqual(t, hash2, "")
	assert.Equal(t, len(hash1), 32)
	assert.NotEqual(t, hash1, hash2)
}

func TestStableHashSlot_Complex(t *testing.T) {
	// 多组测试数据，key是测试名，value是map: 输入字符串 -> 预期槽位
	testData := map[string]struct {
		minNum, maxNum int
		inputs         []string
	}{
		"range_0_9": {
			minNum: 0, maxNum: 9,
			inputs: []string{"hello", "world", "golang", "chatgpt", "openai"},
		},
		"range_10_20": {
			minNum: 10, maxNum: 20,
			inputs: []string{"hello", "world", "golang", "chatgpt", "openai"},
		},
		"single_point": {
			minNum: 5, maxNum: 5,
			inputs: []string{"anything", "something", "nothing"},
		},
	}

	for testName, data := range testData {
		t.Run(testName, func(t *testing.T) {
			// 先计算所有输入的预期槽位，确保稳定性
			expectedSlots := make(map[string]int, len(data.inputs))
			for _, input := range data.inputs {
				slot := osx.StableHashSlot(input, data.minNum, data.maxNum)
				expectedSlots[input] = slot
				// 断言单次调用结果一定在范围内
				assert.True(t, slot >= data.minNum && slot <= data.maxNum,
					"slot %d for input %q should be in range [%d,%d]",
					slot, input, data.minNum, data.maxNum)
			}

			// 再次调用，确保稳定性和一致性
			for input, expected := range expectedSlots {
				got := osx.StableHashSlot(input, data.minNum, data.maxNum)
				assert.Equal(t, expected, got, "input %q stable slot mismatch", input)
			}
		})
	}

	// 测试 panic 情况
	t.Run("panic_when_max_less_than_min", func(t *testing.T) {
		assert.Panics(t, func() {
			osx.StableHashSlot("test", 10, 5)
		}, "should panic when maxNum < minNum")
	})
}

func TestGetServerIP(t *testing.T) {
	externalIP, internalIP, err := osx.GetLocalInterfaceIeIp()
	assert.Nil(t, err)
	if externalIP != "" {
		t.Logf("externalIP %s", externalIP)
	}
	if internalIP != "" {
		t.Logf("internalIP %s", internalIP)
	}
}

func TestGetLocalInterfaceIps(t *testing.T) {
	ips, err := osx.GetLocalInterfaceIps()
	assert.Nil(t, err)
	assert.NotEmpty(t, ips, fmt.Sprintf("Expected at least one global unicast IP, got: %v", ips))
	for _, ip := range ips {
		assert.NotEmpty(t, ip, fmt.Sprintf("Invalid IP address: %s", ip))
	}
}

func TestGetClientPublicIP_XForwardedFor(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	testIp := "1.2.3.4"
	req.Header.Set("X-Forwarded-For", testIp)
	ip, err := osx.GetClientPublicIP(req)
	assert.Nil(t, err)
	assert.Equal(t, testIp, ip)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected IP %s, got: %s", testIp, ip))
}

func TestGetClientPublicIP_XRealIp(t *testing.T) {
	testIp := "113.168.80.129"
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-Ip", testIp)
	ip, err := osx.GetClientPublicIP(req)
	assert.Nil(t, err)
	assert.Equal(t, testIp, ip)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected IP %s, got: %s", testIp, ip))
}

func TestGetClientPublicIP_RemoteAddr(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "115.10.11.12:12345"
	spIp := strings.Split(req.RemoteAddr, ":")[0]
	ip, err := osx.GetClientPublicIP(req)
	assert.Nil(t, err)
	assert.Equal(t, spIp, ip)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected IP %s, got: %s", spIp, ip))
}

func TestGetConNetPublicIp(t *testing.T) {
	ip, err := osx.GetConNetPublicIp()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected public IP, got: %s", ip))
}

func TestGetClientPublicIP_NoValidIp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "127.0.0.1") // Localhost IP
	req.Header.Set("X-Real-Ip", "169.254.0.1")     // Link-local IP
	req.RemoteAddr = "192.168.1.1:12345"           // Private IP
	ip, err := osx.GetClientPublicIP(req)
	assert.Nil(t, err)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected public IP, got: %s", ip))
}

func TestGetRuntimeCaller(t *testing.T) {
	// 调用 GetRuntimeCaller，skip=1 表示跳过测试函数本身，获取调用该函数的调用者信息
	caller := osx.GetRuntimeCaller(1)
	defer caller.Release()

	// 断言 File 不为空且不是 unknown_file
	assert.NotEmpty(t, caller.File, "caller.File should not be empty")
	assert.NotEqual(t, "unknown_file", caller.File, "caller.File should not be unknown_file")

	// 断言 Line 大于 0
	assert.Greater(t, caller.Line, 0, "caller.Line should be greater than 0")

	// 断言 FuncName 不为空且不是 unknown_func
	assert.NotEmpty(t, caller.FuncName, "caller.FuncName should not be empty")
	assert.NotEqual(t, "unknown_func", caller.FuncName, "caller.FuncName should not be unknown_func")

	// 断言 Pc 不为 0
	assert.NotZero(t, caller.Pc, "caller.Pc should not be zero")
}

func TestReleaseClearsFields(t *testing.T) {
	caller := osx.GetRuntimeCaller(2)

	// 先断言字段有值
	assert.NotEmpty(t, caller.File)
	assert.NotZero(t, caller.Line)
	assert.NotEmpty(t, caller.FuncName)
	assert.NotZero(t, caller.Pc)

	// 调用 Release
	caller.Release()

	// 断言字段被清空
	assert.Empty(t, caller.File)
	assert.Zero(t, caller.Line)
	assert.Empty(t, caller.FuncName)
	assert.Zero(t, caller.Pc)
}

// TestCommand 测试 Command 函数
func TestCommand(t *testing.T) {
	// 使用 echo 命令进行测试
	output, err := osx.Command("echo", []string{"Hello, World!"}, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedOutput := "Hello, World!\n"
	if string(output) != expectedOutput {
		t.Errorf("expected %q, got %q", expectedOutput, string(output))
	}
}
