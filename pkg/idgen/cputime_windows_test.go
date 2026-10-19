//go:build windows

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-19 19:58:17
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-19 20:09:13
 * @FilePath: \go-toolbox\pkg\idgen\cputime_windows_test.go
 * @Description: 取进程 CPU 时间（Windows 无等价实现，占位）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package idgen

import "time"

// cpuTimeSupported Windows 上没有 getrusage 的等价物
//
// 相关测试只打印观测值、不做断言，所以这里给占位实现并让测试跳过，
// 而不是为一个诊断用途的测试去包一层 GetProcessTimes 的系统调用
//
// [EN] There is no getrusage equivalent on Windows. The test concerned only logs
// observations without asserting, so this stub lets it skip rather than wrapping
// a GetProcessTimes syscall for a diagnostic-only test
const cpuTimeSupported = false

// processCPUTime 占位实现，恒返回零
func processCPUTime() time.Duration { return 0 }
