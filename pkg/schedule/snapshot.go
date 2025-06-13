/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-24 10:55:16
 * @FilePath: \go-toolbox\pkg\schedule\snapshot.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

type execStatus int

const (
	Pending execStatus = iota
	Running
	Failure
	Success
	SysTermination
	UserTermination
)

// ExceedTaskSnapshot 结构体表示任务执行的快照信息
type ExceedTaskSnapshot struct {
	traceId          string       // 跟踪Id
	execFrequency    int          // 执行次数
	failureFrequency int          // 失败次数
	execStatus                    // 执行状态
	execLogRecord    []string     // 执行链路日志记录
	mu               sync.RWMutex // 读写锁
}

// NewExceedTaskSnapshot 创建并返回一个新的 ExceedTaskSnapshot 实例
func NewExceedTaskSnapshot() *ExceedTaskSnapshot {
	return &ExceedTaskSnapshot{
		execFrequency:    0,
		failureFrequency: 0,
		execStatus:       Pending,
		execLogRecord:    []string{},
	}
}

// SetExecFrequency 设置执行次数并返回指向当前实例的指针
func (s *ExceedTaskSnapshot) SetExecFrequency(frequency int) *ExceedTaskSnapshot {
	return syncx.WithLockReturnValue(&s.mu, func() *ExceedTaskSnapshot {
		s.execFrequency = frequency
		return s
	})
}

// SetFailureFrequency 设置失败次数并返回指向当前实例的指针
func (s *ExceedTaskSnapshot) SetFailureFrequency(frequency int) *ExceedTaskSnapshot {
	return syncx.WithLockReturnValue(&s.mu, func() *ExceedTaskSnapshot {
		s.failureFrequency = frequency
		return s
	})
}

// IncExecFrequency 对执行频率计数器执行线程安全的自增操作
// 返回当前结构体指针，支持链式调用
func (s *ExceedTaskSnapshot) IncExecFrequency() *ExceedTaskSnapshot {
	return syncx.WithLockReturnValue(&s.mu, func() *ExceedTaskSnapshot {
		s.execFrequency++ // 执行频率加1
		return s
	})
}

// IncFailureFrequency 对失败频率计数器执行线程安全的自增操作
// 返回当前结构体指针，支持链式调用
func (s *ExceedTaskSnapshot) IncFailureFrequency() *ExceedTaskSnapshot {
	return syncx.WithLockReturnValue(&s.mu, func() *ExceedTaskSnapshot {
		s.failureFrequency++ // 失败频率加1
		return s
	})
}

// SetExecStatus 设置执行状态并返回指向当前实例的指针
func (s *ExceedTaskSnapshot) SetExecStatus(status execStatus) *ExceedTaskSnapshot {
	return syncx.WithLockReturnValue(&s.mu, func() *ExceedTaskSnapshot {
		s.execStatus = status
		return s
	})
}

// AddLogRecord 添加日志记录并返回指向当前实例的指针
func (s *ExceedTaskSnapshot) AddLogRecord(log string) *ExceedTaskSnapshot {
	return syncx.WithLockReturnValue(&s.mu, func() *ExceedTaskSnapshot {
		s.execLogRecord = append(s.execLogRecord, log)
		return s
	})
}

// GetExecFrequency 获取执行次数
func (s *ExceedTaskSnapshot) GetExecFrequency() int {
	return syncx.WithRLockReturnValue(&s.mu, func() int {
		return s.execFrequency
	})
}

// GetFailureFrequency 获取失败次数
func (s *ExceedTaskSnapshot) GetFailureFrequency() int {
	return syncx.WithRLockReturnValue(&s.mu, func() int {
		return s.failureFrequency
	})
}

// GetExecStatus 获取执行状态
func (s *ExceedTaskSnapshot) GetExecStatus() execStatus {
	return syncx.WithRLockReturnValue(&s.mu, func() execStatus {
		return s.execStatus
	})
}

// GetLogRecords 获取日志记录
func (s *ExceedTaskSnapshot) GetLogRecords() []string {
	return syncx.WithRLockReturnValue(&s.mu, func() []string {
		return s.execLogRecord
	})
}

// GetTraceId 获取跟踪Id
func (s *ExceedTaskSnapshot) GetTraceId() string {
	return syncx.WithRLockReturnValue(&s.mu, func() string {
		return s.traceId
	})
}

// SetTraceId 设置跟踪Id
func (s *ExceedTaskSnapshot) SetTraceId(traceId string) *ExceedTaskSnapshot {
	return syncx.WithLockReturnValue(&s.mu, func() *ExceedTaskSnapshot {
		s.traceId = traceId
		return s
	})
}
