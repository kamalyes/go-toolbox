/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 19:31:22
 * @FilePath: \go-toolbox\pkg\osx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

// GetHostName 获取主机名，如果失败则返回错误
func GetHostName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("无法获取主机名: %v", err)
	}
	return hostname, nil
}

// 获取主机名函数
func SafeGetHostName() string {
	output, err := GetHostName()
	if err != nil || output == "" {
		// 如果获取主机名失败或返回空字符串，则生成随机字符串
		return stringx.ReplaceSpecialChars(random.FRandAlphaString(8), 'x')
	}
	return stringx.ReplaceSpecialChars(output, 'x')
}

// HashUnixMicroCipherText
func HashUnixMicroCipherText() string {
	var (
		nowUnixMicro = time.Now().UnixMicro()
		hostName     = SafeGetHostName()
		randStr      = random.RandString(10, 4)
		plainText    = fmt.Sprintf("%s%s%d", hostName, randStr, nowUnixMicro)
		cipherText   = stringx.CalculateMD5Hash(plainText)
	)
	return cipherText
}

// GetWorkerId 获取唯一的 Worker ID (智能增强版)
// 优先级: 环境变量 WORKER_ID > 主机名哈希
// 支持环境变量: WORKER_ID, NODE_ID, POD_ORDINAL
// 返回范围: 0-1023
func GetWorkerId() int64 {
	// 1. 尝试从环境变量获取
	envVars := []string{"WORKER_ID", "NODE_ID", "POD_ORDINAL"}
	for _, envVar := range envVars {
		if workerID := os.Getenv(envVar); workerID != "" {
			if id, err := mathx.ParseInt64(workerID); err == nil {
				// 确保在 0-1023 范围内
				return mathx.Abs(id) % 1024
			}
		}
	}

	// 2. 从主机名生成唯一的 Worker ID
	hostName := SafeGetHostName()
	hash := sha256.Sum256([]byte(hostName))
	hostNameHash := int64(binary.BigEndian.Uint64(hash[:8]))

	// 确保在 0-1023 范围内
	workerId := mathx.Abs(hostNameHash) % 1024
	return workerId
}

// GetDatacenterId 获取数据中心ID (新增)
// 优先级: 环境变量 DATACENTER_ID > 区域环境变量 > 默认值1
// 返回范围: 0-31 (雪花算法标准)
func GetDatacenterId() int64 {
	// 1. 尝试从环境变量获取
	envVars := []string{"DATACENTER_ID", "DC_ID", "ZONE_ID", "REGION_ID"}
	for _, envVar := range envVars {
		if dcID := os.Getenv(envVar); dcID != "" {
			if id, err := mathx.ParseInt64(dcID); err == nil {
				// 确保在 0-31 范围内 (雪花算法的datacenter位数限制)
				return mathx.Abs(id) % 32
			}
		}
	}

	// 2. 从Kubernetes相关环境变量推导
	if namespace := os.Getenv("KUBERNETES_NAMESPACE"); namespace != "" {
		hash := sha256.Sum256([]byte(namespace))
		return int64(binary.BigEndian.Uint32(hash[:4])) % 32
	}

	// 3. 默认返回 1
	return 1
}

// GetWorkerIdForSnowflake 专门为雪花算法获取WorkerId (新增)
// 返回范围: 0-31 (雪花算法标准的5位worker id)
func GetWorkerIdForSnowflake() int64 {
	workerId := GetWorkerId()
	// 雪花算法的worker id只有5位，范围0-31
	return workerId % 32
}

// StableHashSlot 根据输入字符串 s 和范围 [minNum, maxNum]，返回一个稳定且范围内的整数
// 使用加密哈希 sha256，抗碰撞更强
// 如果 maxNum < minNum，会 panic
// 如果 maxNum == minNum，直接返回 minNum
func StableHashSlot(s string, minNum, maxNum int) int {
	if maxNum < minNum {
		panic("maxNum must be >= minNum")
	}
	if maxNum == minNum {
		return minNum
	}

	var (
		hashBytes = sha256.Sum256([]byte(s))
		// 取哈希结果的前8字节转为uint64
		hashVal   = binary.BigEndian.Uint64(hashBytes[:8])
		rangeSize = maxNum - minNum + 1
		result    = int(hashVal%uint64(rangeSize)) + minNum
	)
	return mathx.IF(result > 0, result, 1)
}

// RunTimeCaller 结构体用于存储调用栈信息
type RunTimeCaller struct {
	Pc       uintptr // 程序计数器
	File     string  // 文件名
	Line     int     // 行号
	FuncName string  // 函数名
}

var callerPool = sync.Pool{
	New: func() interface{} {
		return &RunTimeCaller{}
	},
}

// GetRuntimeCaller 获取调用栈信息，调用者使用完需调用 Release() 归还对象
func GetRuntimeCaller(skip int) *RunTimeCaller {
	caller := callerPool.Get().(*RunTimeCaller)
	caller.init(skip)
	return caller
}

// init 初始化 RunTimeCaller 内容
func (c *RunTimeCaller) init(skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		c.Pc = 0
		c.File = "unknown_file"
		c.Line = 0
		c.FuncName = "unknown_func"
		return
	}

	c.Pc = pc
	c.File = file
	c.Line = line

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		c.FuncName = "unknown_func"
		return
	}

	fullName := fn.Name()
	// 只保留函数名（去掉包路径）
	if lastSlash := strings.LastIndex(fullName, "/"); lastSlash != -1 {
		fullName = fullName[lastSlash+1:]
	}
	if lastDot := strings.LastIndex(fullName, "."); lastDot != -1 {
		fullName = fullName[lastDot+1:]
	}
	c.FuncName = fullName
}

// Release 放回对象池，调用者必须调用
func (c *RunTimeCaller) Release() {
	// 清理字段，防止内存泄漏
	c.Pc = 0
	c.File = ""
	c.Line = 0
	c.FuncName = ""
	callerPool.Put(c)
}

// String 返回调用栈的格式化字符串，包含函数名、文件名和行号
func (c *RunTimeCaller) String() string {
	return fmt.Sprintf("FuncName:%s, File:%s:%d", c.FuncName, c.File, c.Line)
}

// Command 执行系统命令
func Command(bin string, argv []string, baseDir string) ([]byte, error) {
	cmd := exec.Command(bin, argv...)

	if baseDir != "" {
		cmd.Dir = baseDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stdout.Bytes(), fmt.Errorf("command failed: %s: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}
