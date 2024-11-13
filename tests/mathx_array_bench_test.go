/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 15:55:06
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 23:20:55
 * @FilePath: \go-toolbox\tests\mathx_array_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
)

// BenchmarkArrayChunk 基准测试 ArrayChunk 函数
func BenchmarkArrayChunk(b *testing.B) {
	input := make([]int, 10000) // 创建一个包含 10000 个元素的切片
	for i := 0; i < 10000; i++ {
		input[i] = i
	}

	b.ResetTimer() // 重置计时器以排除设置时间
	for i := 0; i < b.N; i++ {
		mathx.ArrayChunk(input, 100) // 测试分块大小为 100
	}
}
