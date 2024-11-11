/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 15:39:53
 * @FilePath: \go-toolbox\tests\osx_console_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
)

// BenchmarkInfo 测试 colorConsole 的 Info 方法在高并发下的性能
func BenchmarkInfo(b *testing.B) {
	console := osx.NewColorConsole(true) // 创建一个 colorConsole 实例
	var wg sync.WaitGroup

	// 定义一个匿名函数作为 goroutine 的工作负载
	worker := func() {
		defer wg.Done()
		for i := 0; i < b.N/10; i++ { // b.N 是基准测试框架提供的循环次数
			console.Info("This is an info message %d", i)
		}
	}

	// 启动多个 goroutine 来模拟高并发
	numWorkers := 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// 等待所有 goroutine 完成
	wg.Wait()
}
