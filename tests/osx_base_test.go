/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-18 17:36:32
 * @FilePath: \go-toolbox\tests\osx_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"sync"
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

// BenchmarkGetHostName 测试 GetHostName 的性能
func BenchmarkGetHostName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		osx.GetHostName()
	}
}

// BenchmarkSafeGetHostName 测试 SafeGetHostName 的性能
func BenchmarkSafeGetHostName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		osx.SafeGetHostName()
	}
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

func TestGetWorkerId(t *testing.T) {
	workerId := osx.GetWorkerId()
	fmt.Printf("TestGetWorkerId workerId %#v", workerId)
	assert.NotEmpty(t, workerId)
	assert.Less(t, workerId, int64(1024)) // Worker ID 范围为 0-1023
}

// 测试多次调用以确保一致性
func TestGetWorkerIdConsistency(t *testing.T) {
	const numCalls = 1000
	var wg sync.WaitGroup
	results := make([]int64, numCalls)
	errors := make([]error, numCalls)

	// 使用 WaitGroup 来等待所有 goroutine 完成
	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = osx.GetWorkerId()
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 检查是否有错误
	for _, err := range errors {
		assert.NoError(t, err)
	}

	// 检查所有结果是否一致
	for i := 1; i < numCalls; i++ {
		assert.Equal(t, results[0], results[i], "WorkerInfo 应该在多次调用中保持一致")
	}
}

// BenchmarkGetWorkerId 测试 GetWorkerId 的性能
func BenchmarkGetWorkerId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		osx.GetWorkerId()
	}
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
	var cmd string
	var args []string
	var expectedOutput string

	// 根据操作系统选择合适的命令
	if osx.IsWindows() {
		// Windows 环境使用 cmd /c echo（不使用引号）
		cmd = "cmd"
		args = []string{"/c", "echo", "Hello, World!"}
		expectedOutput = "\"Hello, World!\"\r\n" // Windows echo 会添加引号
	} else {
		// Unix-like 环境使用 echo
		cmd = "echo"
		args = []string{"Hello, World!"}
		expectedOutput = "Hello, World!\n"
	}

	output, err := osx.Command(cmd, args, "")
	assert.NoError(t, err, "Command execution should not return error")
	assert.Equal(t, expectedOutput, string(output), "Command output should match expected")
}
