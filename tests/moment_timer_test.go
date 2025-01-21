/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-09 19:15:01
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-21 18:55:55
 * @FilePath: \go-toolbox\tests\moment_timer_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/moment"
	"github.com/stretchr/testify/assert"
)

// TestTimer 测试 Timer 结构体的功能
func TestTimer(t *testing.T) {
	// 创建一个计时器
	timer := moment.NewTimer()
	timer.Run()

	// 等待一段时间，以便计时器可以打印时差
	time.Sleep(250 * time.Millisecond)

	// 结束
	timer.Finish()

	// 获取当前时间与开始时间的差异
	elapsed := timer.GetEndTime().Sub(timer.GetStartTime())

	// 断言时差大于等于 250 毫秒（考虑到可能的延迟）
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(250), "时差应该大于等于 250 毫秒")
}

// TestTimerConcurrent 测试 Timer 结构体的并发功能
func TestTimerConcurrent(t *testing.T) {
	const numTimers = 10  // 并发计时器的数量
	var wg sync.WaitGroup // 用于等待所有 goroutine 完成

	for i := 0; i < numTimers; i++ {
		wg.Add(1) // 增加等待计数
		go func(i int) {
			defer wg.Done() // 完成时减少计数

			// 创建并运行计时器
			timer := moment.NewTimer()
			timer.Run()

			// 等待一段时间
			time.Sleep(time.Duration(100+i*20) * time.Millisecond)

			// 更新结束时间
			timer.Finish()

			// 获取当前时间与开始时间的差异
			elapsed := timer.GetEndTime().Sub(timer.GetStartTime())

			// 断言时差大于等于 100 毫秒（考虑到可能的延迟）
			assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(100), "时差应该大于等于 100 毫秒")
		}(i)
	}

	wg.Wait() // 等待所有 goroutine 完成
}

func TestPauseAndResumeTimer(t *testing.T) {
	timer := moment.NewTimer()
	// 启动计时器
	timer.Run()
	time.Sleep(100 * time.Millisecond) // 等待一定时间

	// 定义暂停和恢复的次数
	const iterations = 5
	var totalPauseDuration time.Duration

	for i := 0; i < iterations; i++ {
		// 暂停计时器
		timer.Pause()
		pauseDuration := 1000 * time.Millisecond
		totalPauseDuration += pauseDuration
		time.Sleep(pauseDuration) // 模拟暂停时间

		// 恢复计时器
		timer.Resume()
		time.Sleep(150 * time.Millisecond) // 等待一段时间
	}

	// 完成计时器
	timer.Finish()

	// 获取持续时间
	expectedDuration := 100*time.Millisecond + (150 * time.Millisecond * time.Duration(iterations))
	actualDuration := timer.GetDuration()

	// 使用 assert 检查持续时间
	assert.GreaterOrEqual(t, actualDuration, expectedDuration, "实际持续时间应大于或等于预期持续时间")
	assert.LessOrEqual(t, actualDuration, expectedDuration+100*time.Millisecond, "实际持续时间应小于或等于预期持续时间加上允许的误差")
	assert.InDelta(t, expectedDuration.Seconds(), actualDuration.Seconds(), 0.1, "实际持续时间应在预期持续时间的范围内")

	// 检查暂停持续时间
	expectedPauseDuration := totalPauseDuration
	actualPauseDuration := timer.GetPauseDuration()

	assert.InDelta(t, expectedPauseDuration.Seconds(), actualPauseDuration.Seconds(), 0.1, "实际暂停持续时间应在预期暂停持续时间的范围内")
}
