//go:build !windows

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-19 19:58:11
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-19 20:09:13
 * @FilePath: \go-toolbox\pkg\idgen\cputime_unix_test.go
 * @Description: 取进程 CPU 时间（类 Unix 平台）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"syscall"
	"time"
)

// cpuTimeSupported 本平台能否取到进程 CPU 时间
const cpuTimeSupported = true

// processCPUTime 返回本进程累计消耗的 CPU 时间（用户态加内核态）
//
// 用 getrusage 而不是 time.Now：要量的是「等待期间有没有在空转」，
// 而空转的特征正是 CPU 时间逼近墙上时间
//
// 拆成按平台编译的两个文件，是因为 syscall.Getrusage 与 syscall.Rusage 在
// Windows 上不存在，放在一起会让整个测试二进制在 windows 上编译不过——
// CI 的平台矩阵里有 windows-latest，这正是它报失败的原因
//
// [EN] Returns the CPU time this process has consumed, user plus system. It uses
// getrusage rather than time.Now because what is being measured is whether the
// wait spins, and spinning shows up as CPU time approaching wall time. The file
// is split per platform because syscall.Getrusage and syscall.Rusage do not
// exist on Windows, and keeping them together stops the whole test binary from
// compiling there — windows-latest is in the CI platform matrix, which is what
// made it report a failure
func processCPUTime() time.Duration {
	var ru syscall.Rusage
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}
