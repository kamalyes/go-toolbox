/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-11 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 15:11:00
 * @FilePath: \go-toolbox\pkg\schedule\snapshot_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExceedTaskSnapshot_BasicOperations(t *testing.T) {
	snap := NewExceedTaskSnapshot()
	assert.NotNil(t, snap)

	// 初始状态检查
	assert.Equal(t, 0, snap.GetExecFrequency())
	assert.Equal(t, 0, snap.GetFailureFrequency())
	assert.Equal(t, Pending, snap.GetExecStatus())
	assert.Empty(t, snap.GetLogRecords())
	assert.Empty(t, snap.GetTraceId())

	// 设置并获取 ExecFrequency
	snap.SetExecFrequency(10)
	assert.Equal(t, 10, snap.GetExecFrequency())

	// 设置并获取 FailureFrequency
	snap.SetFailureFrequency(5)
	assert.Equal(t, 5, snap.GetFailureFrequency())

	// 设置并获取 ExecStatus
	snap.SetExecStatus(Running)
	assert.Equal(t, Running, snap.GetExecStatus())

	// 添加日志记录
	snap.AddLogRecord("log1")
	snap.AddLogRecord("log2")
	logs := snap.GetLogRecords()
	assert.Len(t, logs, 2)
	assert.Equal(t, "log1", logs[0])
	assert.Equal(t, "log2", logs[1])

	// 设置并获取 TraceId
	snap.SetTraceId("trace-123")
	assert.Equal(t, "trace-123", snap.GetTraceId())
}

func TestExceedTaskSnapshot_ConcurrentAccess(t *testing.T) {
	snap := NewExceedTaskSnapshot()
	const goroutines = 100
	const increments = 1000

	var wg sync.WaitGroup

	// 并发增加 ExecFrequency
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				snap.IncExecFrequency()
			}
		}()
	}

	// 并发增加 FailureFrequency
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				snap.IncFailureFrequency()
			}
		}()
	}

	// 并发添加日志
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				snap.AddLogRecord(
					// 简单日志内容
					"goroutine " + string(rune(id)) + " log " + string(rune(j)),
				)
			}
		}(i)
	}

	// 并发设置状态和traceId
	wg.Add(goroutines)
	statuses := []execStatus{Pending, Running, Failure, Success, SysTermination, UserTermination}
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			snap.SetExecStatus(statuses[i%len(statuses)])
			snap.SetTraceId("trace-" + string(rune(i)))
		}(i)
	}

	wg.Wait()

	// 验证结果
	assert.Equal(t, goroutines*increments, snap.GetExecFrequency(), "ExecFrequency should match increments")
	assert.Equal(t, goroutines*increments, snap.GetFailureFrequency(), "FailureFrequency should match increments")

	// 日志数量检查 (至少 goroutines * 10)
	logs := snap.GetLogRecords()
	assert.GreaterOrEqual(t, len(logs), goroutines*10)

	// 状态断言：必须是有效枚举值 取最后一次设置的值，无法确定具体值，只检查不为空
	status := snap.GetExecStatus()
	validStatuses := map[execStatus]bool{
		Pending: true, Running: true, Failure: true, Success: true, SysTermination: true, UserTermination: true,
	}
	assert.True(t, validStatuses[status], "ExecStatus should be a valid enum value")

	// traceId 不为空
	assert.NotEmpty(t, snap.GetTraceId())
}

func TestExceedTaskSnapshot_ChainCalls(t *testing.T) {
	snap := NewExceedTaskSnapshot()

	// 链式调用
	snap.SetExecFrequency(1).
		SetFailureFrequency(2).
		SetExecStatus(Failure).
		AddLogRecord("log1").
		SetTraceId("trace-xyz").
		AddLogRecord("log2")

	assert.Equal(t, 1, snap.GetExecFrequency())
	assert.Equal(t, 2, snap.GetFailureFrequency())
	assert.Equal(t, Failure, snap.GetExecStatus())
	assert.Equal(t, "trace-xyz", snap.GetTraceId())

	logs := snap.GetLogRecords()
	assert.Len(t, logs, 2)
	assert.Equal(t, "log1", logs[0])
	assert.Equal(t, "log2", logs[1])
}
