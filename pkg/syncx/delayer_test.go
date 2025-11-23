/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 19:03:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 19:33:15
 * @FilePath: \go-toolbox\pkg\syncx\delayer_test.go
 * @Description: 泛型延迟器测试文件
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestDelayerBasic 测试泛型基本功能
func TestDelayerBasic(t *testing.T) {
	// 创建一个返回字符串的泛型延迟器
	delayer := NewDelayer[string]().
		WithDelay(50 * time.Millisecond).
		WithTimes(3).
		WithTaskFunc(func(ctx *ExecutionContext) (string, error) {
			return fmt.Sprintf("Task %d completed", ctx.Index), nil
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	results := delayer.GetResults()
	assert.Equal(t, 3, len(results), "Expected 3 results")

	for i, result := range results {
		expected := fmt.Sprintf("Task %d completed", i)
		assert.Equal(t, expected, result, "Result mismatch")
	}

	delayer.Close()
}

// TestDelayerWithNumbers 测试泛型数字类型
func TestDelayerWithNumbers(t *testing.T) {
	delayer := NewDelayer[int]().
		WithDelay(30 * time.Millisecond).
		WithTimes(5).
		WithTaskFunc(func(ctx *ExecutionContext) (int, error) {
			return ctx.Index * 10, nil
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	results := delayer.GetResults()
	t.Logf("Actual results: %v", results)
	assert.Equal(t, 5, len(results), "Expected 5 results")

	// 由于并发执行可能导致顺序不同，我们检查结果是否正确
	resultMap := make(map[int]bool)
	for _, result := range results {
		resultMap[result] = true
	}

	expectedResults := []int{0, 10, 20, 30, 40}
	for _, expected := range expectedResults {
		assert.True(t, resultMap[expected], fmt.Sprintf("Expected result %d not found", expected))
	}

	delayer.Close()
}

// TestDelayerWithStruct 测试泛型结构体类型
func TestDelayerWithStruct(t *testing.T) {
	type TaskResult struct {
		ID    int
		Name  string
		Value float64
	}

	delayer := NewDelayer[TaskResult]().
		WithDelay(25 * time.Millisecond).
		WithTimes(3).
		WithTaskFunc(func(ctx *ExecutionContext) (TaskResult, error) {
			return TaskResult{
				ID:    ctx.Index,
				Name:  fmt.Sprintf("Task_%d", ctx.Index),
				Value: float64(ctx.Index) * 1.5,
			}, nil
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	results := delayer.GetResults()
	assert.Equal(t, 3, len(results), "Expected 3 results")

	for i, result := range results {
		assert.Equal(t, i, result.ID)
		assert.Equal(t, fmt.Sprintf("Task_%d", i), result.Name)
		assert.Equal(t, float64(i)*1.5, result.Value)
	}

	delayer.Close()
}

// TestDelayerWithCallback 测试泛型回调
func TestDelayerWithCallback(t *testing.T) {
	callbackResults := make([]string, 0)
	var callbackMutex sync.Mutex

	delayer := NewDelayer[string]().
		WithDelay(40 * time.Millisecond).
		WithTimes(3).
		WithTaskFunc(func(ctx *ExecutionContext) (string, error) {
			return strconv.Itoa(ctx.Index), nil
		}).
		WithOnSuccess(func(ctx *ExecutionContext, result string) {
			callbackMutex.Lock()
			defer callbackMutex.Unlock()
			callbackResults = append(callbackResults, fmt.Sprintf("Success: %s", result))
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	callbackMutex.Lock()
	assert.Equal(t, 3, len(callbackResults), "Expected 3 callback results")
	callbackMutex.Unlock()

	delayer.Close()
}

// TestDelayerResultChannel 测试结果通道
func TestDelayerResultChannel(t *testing.T) {
	delayer := NewDelayer[int]().
		WithDelay(20 * time.Millisecond).
		WithTimes(5).
		WithTaskFunc(func(ctx *ExecutionContext) (int, error) {
			return ctx.Index + 100, nil
		})

	// 从通道读取结果
	var channelResults []int
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for result := range delayer.GetResultChannel() {
			channelResults = append(channelResults, result)
		}
	}()

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()
	delayer.Close()

	// 等待通道读取完成
	wg.Wait()

	assert.Equal(t, 5, len(channelResults), "Expected 5 channel results")

	// 验证结果（注意：由于并发，顺序可能不一样）
	expectedSum := 100 + 101 + 102 + 103 + 104 // 500
	actualSum := 0
	for _, result := range channelResults {
		actualSum += result
	}

	assert.Equal(t, expectedSum, actualSum, "Sum of results should match")
}

// TestDelayerWithError 测试泛型错误处理
func TestDelayerWithError(t *testing.T) {
	delayer := NewDelayer[string]().
		WithDelay(30 * time.Millisecond).
		WithTimes(5).
		WithTaskFunc(func(ctx *ExecutionContext) (string, error) {
			if ctx.Index == 2 {
				return "", fmt.Errorf("simulated error at index %d", ctx.Index)
			}
			return fmt.Sprintf("Success_%d", ctx.Index), nil
		}).
		WithOnErrorContext(func(ctx *ExecutionContext) bool {
			t.Logf("Error occurred at index %d: %v", ctx.Index, ctx.Error)
			return true // 继续执行
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	results := delayer.GetResults()
	assert.Equal(t, 4, len(results), "Expected 4 successful results") // 5 - 1 error

	stats := delayer.GetStats()
	// 注意：新的统计字段名称
	totalExecutions := stats.SuccessCount + stats.ErrorCount
	assert.Equal(t, int64(5), totalExecutions, "Expected 5 total executions")
	assert.Equal(t, int64(4), stats.SuccessCount, "Expected 4 successful executions")
	assert.Equal(t, int64(1), stats.ErrorCount, "Expected 1 failed execution")

	delayer.Close()
}

// TestDelayerPerformance 性能测试
func TestDelayerPerformance(t *testing.T) {
	const taskCount = 1000

	start := time.Now()

	delayer := NewDelayer[int]().
		WithDelay(1 * time.Millisecond).
		WithTimes(taskCount).
		WithConcurrent(true).
		WithMaxConcurrency(50).
		WithDisableCallbacks(true). // 禁用回调以提升性能
		WithTaskFunc(func(ctx *ExecutionContext) (int, error) {
			// 简单的计算任务
			return ctx.Index * ctx.Index, nil
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	duration := time.Since(start)
	results := delayer.GetResults()

	assert.Equal(t, taskCount, len(results), "Expected all tasks to complete")
	t.Logf("Completed %d tasks in %v (%.2f tasks/sec)",
		taskCount, duration, float64(taskCount)/duration.Seconds())

	delayer.Close()
}

// TestDelayerBasic 测试泛型基本功能
func TestDelayerString(t *testing.T) {
	// 创建一个返回字符串的泛型延迟器
	delayer := NewDelayer[string]().
		WithDelay(50 * time.Millisecond).
		WithTimes(3).
		WithTaskFunc(func(ctx *ExecutionContext) (string, error) {
			return fmt.Sprintf("Task %d completed", ctx.Index), nil
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	results := delayer.GetResults()
	assert.Equal(t, 3, len(results), "Expected 3 results")

	for i, result := range results {
		expected := fmt.Sprintf("Task %d completed", i)
		assert.Equal(t, expected, result, "Result mismatch")
	}

	delayer.Close()
}

// TestDelayerWithNumbers 测试泛型数字类型
func TestDelayerInt(t *testing.T) {
	delayer := NewDelayer[int]().
		WithDelay(30 * time.Millisecond).
		WithTimes(5).
		WithTaskFunc(func(ctx *ExecutionContext) (int, error) {
			return ctx.Index * 10, nil
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	results := delayer.GetResults()
	assert.Equal(t, 5, len(results), "Expected 5 results")

	for i, result := range results {
		expected := i * 10
		assert.Equal(t, expected, result, "Number result mismatch")
	}

	delayer.Close()
}

// TaskResult 自定义结构体类型
type TaskResult struct {
	ID      int
	Message string
	Success bool
}

// TestDelayerWithStruct 测试泛型结构体类型
func TestDelayerStruct(t *testing.T) {
	delayer := NewDelayer[TaskResult]().
		WithDelay(25 * time.Millisecond).
		WithTimes(3).
		WithTaskFunc(func(ctx *ExecutionContext) (TaskResult, error) {
			return TaskResult{
				ID:      ctx.Index,
				Message: fmt.Sprintf("Task %d executed", ctx.Index),
				Success: true,
			}, nil
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	results := delayer.GetResults()
	assert.Equal(t, 3, len(results), "Expected 3 results")

	for i, result := range results {
		assert.Equal(t, i, result.ID, "ID mismatch")
		assert.Equal(t, fmt.Sprintf("Task %d executed", i), result.Message, "Message mismatch")
		assert.True(t, result.Success, "Success should be true")
	}

	delayer.Close()
}

// TestDelayerWithCallback 测试泛型回调功能
func TestDelayerCallback(t *testing.T) {
	var callbackResults []string
	var mu sync.Mutex

	delayer := NewDelayer[string]().
		WithDelay(40 * time.Millisecond).
		WithTimes(3).
		WithTaskFunc(func(ctx *ExecutionContext) (string, error) {
			return fmt.Sprintf("Result %d", ctx.Index), nil
		}).
		WithOnSuccess(func(ctx *ExecutionContext, result string) {
			mu.Lock()
			defer mu.Unlock()
			callbackResults = append(callbackResults, result)
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	mu.Lock()
	assert.Equal(t, 3, len(callbackResults), "Expected 3 callback results")
	mu.Unlock()

	delayer.Close()
}

// TestDelayerConcurrent 测试泛型并发执行
func TestDelayerConcurrent(t *testing.T) {
	delayer := NewDelayer[int]().
		WithDelay(20 * time.Millisecond).
		WithTimes(5).
		WithConcurrent(true).
		WithMaxConcurrency(3).
		WithTaskFunc(func(ctx *ExecutionContext) (int, error) {
			return ctx.Index + 100, nil
		})

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()

	results := delayer.GetResults()
	assert.Equal(t, 5, len(results), "Expected 5 results")

	// 验证所有结果都存在（并发执行时顺序可能不同）
	resultMap := make(map[int]bool)
	for _, result := range results {
		resultMap[result] = true
	}

	for i := 0; i < 5; i++ {
		expected := i + 100
		assert.True(t, resultMap[expected], fmt.Sprintf("Expected result %d not found", expected))
	}

	delayer.Close()
}

// TestDelayerResultChannel 测试结果通道
func TestDelayerChannel(t *testing.T) {
	delayer := NewDelayer[string]().
		WithDelay(30 * time.Millisecond).
		WithTimes(5).
		WithTaskFunc(func(ctx *ExecutionContext) (string, error) {
			return "channel-" + strconv.Itoa(ctx.Index), nil
		})

	// 启动一个goroutine来收集通道结果
	var channelResults []string
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for result := range delayer.GetResultChannel() {
			channelResults = append(channelResults, result)
		}
	}()

	err := delayer.Execute()
	assert.NoError(t, err, "Execute should not return error")

	delayer.WaitForCompletion()
	delayer.Close() // 关闭通道

	wg.Wait() // 等待结果收集完成

	assert.Equal(t, 5, len(channelResults), "Expected 5 channel results")

	// 验证结果内容
	resultMap := make(map[string]bool)
	for _, result := range channelResults {
		resultMap[result] = true
	}

	for i := 0; i < 5; i++ {
		expected := "channel-" + strconv.Itoa(i)
		assert.True(t, resultMap[expected], fmt.Sprintf("Expected channel result %s not found", expected))
	}
}

// BenchmarkDelayerPerformance 性能测试
func BenchmarkDelayerPerformance(b *testing.B) {
	taskCount := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delayer := NewDelayer[int]().
			WithDelay(1 * time.Millisecond).
			WithTimes(taskCount).
			WithConcurrent(true).
			WithMaxConcurrency(50).
			WithDisableCallbacks(true). // 禁用回调以提升性能
			WithTaskFunc(func(ctx *ExecutionContext) (int, error) {
				return ctx.Index * 2, nil
			})

		delayer.Execute()
		delayer.WaitForCompletion()
		delayer.Close()
	}
}

// BenchmarkHighConcurrencyAtomic 高并发原子操作基准测试
func BenchmarkHighConcurrencyAtomic(b *testing.B) {
	taskCount := 10000
	concurrency := 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delayer := NewDelayer[int]().
			WithDelay(0). // 无延迟，纯测试并发性能
			WithTimes(taskCount).
			WithConcurrent(true).
			WithMaxConcurrency(concurrency).
			WithDisableCallbacks(true).
			WithTaskFunc(func(ctx *ExecutionContext) (int, error) {
				// 简单的计算任务
				return ctx.Index * ctx.Index, nil
			})

		delayer.Execute()
		delayer.WaitForCompletion()
		delayer.Close()
	}
}

// BenchmarkChannelOperationsAtomic 原子操作通道基准测试
func BenchmarkChannelOperationsAtomic(b *testing.B) {
	delayer := NewDelayer[int]().
		WithDelay(0).
		WithTimes(1).
		WithTaskFunc(func(ctx *ExecutionContext) (int, error) {
			return ctx.Index, nil
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 测试高频的通道发送操作
		delayer.safeChannelSend(i)
	}

	delayer.Close()
}
