/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-18 11:15:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-18 11:37:55
 * @FilePath: \go-toolbox\pkg\logx\empty.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logx

import "context"

// Log 接口定义
type Log interface {
	Debug(ctx context.Context, msg string)
	Trace(ctx context.Context, msg string)
	Notice(ctx context.Context, msg string)
	Warning(ctx context.Context, msg string)
	Error(ctx context.Context, msg string)
	Fatal(ctx context.Context, msg string)
}

// EmptyLog 是 Log 接口的空实现
type EmptyLog struct{}

// NewEmptyLog 创建一个新的 EmptyLog 实例
func NewEmptyLog() *EmptyLog {
	return &EmptyLog{}
}

// 下面的所有方法都没有实现任何功能
func (e *EmptyLog) Debug(ctx context.Context, msg string)   {}
func (e *EmptyLog) Trace(ctx context.Context, msg string)   {}
func (e *EmptyLog) Notice(ctx context.Context, msg string)  {}
func (e *EmptyLog) Warning(ctx context.Context, msg string) {}
func (e *EmptyLog) Error(ctx context.Context, msg string)   {}
func (e *EmptyLog) Fatal(ctx context.Context, msg string)   {}

// NoLog 是一个全局的空日志实例
var NoLog Log = NewEmptyLog()
