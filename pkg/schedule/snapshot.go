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

type execStatus int

const (
	Pending execStatus = iota
	Running
	Failure
	Success
	SysTermination
	UserTermination
)

// 定义状态与日志消息的映射
var execStatusLogMessages = map[execStatus]string{
	Failure:         "Task failed: %v",
	Success:         "Task executed successfully",
	Pending:         "Task is pending",
	Running:         "Task is running",
	SysTermination:  "Task terminated by system",
	UserTermination: "Task terminated by user",
}

// ExceedTaskSnapshot 结构体表示任务执行的快照信息
type ExceedTaskSnapshot struct {
	traceId          string   // 跟踪Id
	execFrequency    int      // 执行次数
	failureFrequency int      // 失败次数
	execStatus                // 执行状态
	execLogRecord    []string // 执行链路日志记录
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
	s.execFrequency = frequency
	return s
}

// SetFailureFrequency 设置失败次数并返回指向当前实例的指针
func (s *ExceedTaskSnapshot) SetFailureFrequency(frequency int) *ExceedTaskSnapshot {
	s.failureFrequency = frequency
	return s
}

// SetExecStatus 设置执行状态并返回指向当前实例的指针
func (s *ExceedTaskSnapshot) SetExecStatus(status execStatus) *ExceedTaskSnapshot {
	s.execStatus = status
	return s
}

// AddLogRecord 添加日志记录并返回指向当前实例的指针
func (s *ExceedTaskSnapshot) AddLogRecord(log string) *ExceedTaskSnapshot {
	s.execLogRecord = append(s.execLogRecord, log)
	return s
}

// GetExecFrequency 获取执行次数
func (s *ExceedTaskSnapshot) GetExecFrequency() int {
	return s.execFrequency
}

// GetFailureFrequency 获取失败次数
func (s *ExceedTaskSnapshot) GetFailureFrequency() int {
	return s.failureFrequency
}

// GetExecStatus 获取执行状态
func (s *ExceedTaskSnapshot) GetExecStatus() execStatus {
	return s.execStatus
}

// GetLogRecords 获取日志记录
func (s *ExceedTaskSnapshot) GetLogRecords() []string {
	return s.execLogRecord
}

// GetTraceId 获取跟踪Id
func (s *ExceedTaskSnapshot) GetTraceId() string {
	return s.traceId
}

// SetTraceId 设置跟踪Id
func (s *ExceedTaskSnapshot) SetTraceId(traceId string) *ExceedTaskSnapshot {
	s.traceId = traceId
	return s
}
