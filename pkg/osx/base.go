/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 19:31:22
 * @FilePath: \go-toolbox\pkg\osx\base.go
 * @Description: 操作系统相关工具函数
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
	"sync/atomic"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

// WorkerIDConfig WorkerID 配置
type WorkerIDConfig struct {
	maxWorkerID     atomic.Int64 // WorkerID 最大值（默认 1024）
	maxDatacenterID atomic.Int64 // DatacenterID 最大值（默认 32）
	maxSnowflakeID  atomic.Int64 // 雪花算法 ID 最大值（默认 32）
}

// RunTimeCaller 运行时调用栈信息
type RunTimeCaller struct {
	Pc       uintptr // 程序计数器
	File     string  // 文件路径
	Line     int     // 行号
	FuncName string  // 函数名
}

var (
	globalWorkerIDConfig = &WorkerIDConfig{}
	callerPool           = sync.Pool{
		New: func() any {
			return &RunTimeCaller{}
		},
	}
)

func init() {
	// 设置默认值
	globalWorkerIDConfig.maxWorkerID.Store(1024)
	globalWorkerIDConfig.maxDatacenterID.Store(32)
	globalWorkerIDConfig.maxSnowflakeID.Store(32)
}

// SetMaxWorkerID 设置 WorkerID 的最大值（范围限制）
// 默认值: 1024
func SetMaxWorkerID(max int64) {
	globalWorkerIDConfig.maxWorkerID.Store(mathx.IfLeZero(max, 1024))
}

// SetMaxDatacenterID 设置 DatacenterID 的最大值（范围限制）
// 默认值: 32（雪花算法标准）
func SetMaxDatacenterID(max int64) {
	globalWorkerIDConfig.maxDatacenterID.Store(mathx.IfLeZero(max, 32))
}

// SetMaxSnowflakeWorkerID 设置雪花算法 WorkerID 的最大值（范围限制）
// 默认值: 32（雪花算法标准的5位）
func SetMaxSnowflakeWorkerID(max int64) {
	globalWorkerIDConfig.maxSnowflakeID.Store(mathx.IfLeZero(max, 32))
}

// GetMaxWorkerID 获取 WorkerID 的最大值
func GetMaxWorkerID() int64 {
	return globalWorkerIDConfig.maxWorkerID.Load()
}

// GetMaxDatacenterID 获取 DatacenterID 的最大值
func GetMaxDatacenterID() int64 {
	return globalWorkerIDConfig.maxDatacenterID.Load()
}

// GetMaxSnowflakeWorkerID 获取雪花算法 WorkerID 的最大值
func GetMaxSnowflakeWorkerID() int64 {
	return globalWorkerIDConfig.maxSnowflakeID.Load()
}

// GetHostName 获取主机名，如果失败则返回错误
func GetHostName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("无法获取主机名: %v", err)
	}
	return hostname, nil
}

// SafeGetHostName 安全获取主机名，失败时返回随机字符串
func SafeGetHostName() string {
	output, err := GetHostName()
	if err != nil || output == "" {
		// 如果获取主机名失败或返回空字符串，则生成随机字符串
		return stringx.ReplaceSpecialChars(random.FRandAlphaString(8), 'x')
	}
	return stringx.ReplaceSpecialChars(output, 'x')
}

// HashUnixMicroCipherText 生成基于时间戳和主机名的哈希密文
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
// 优先级: K8s Pod Name > K8s Hostname > 环境变量 > 主机名哈希
// 支持环境变量: POD_NAME, HOSTNAME, WORKER_ID, NODE_ID, POD_ORDINAL
// 返回范围: 0 到 (maxWorkerID-1)，默认 0-1023
func GetWorkerId() int64 {
	maxWorkerID := globalWorkerIDConfig.maxWorkerID.Load()

	// 1. 优先从 K8s Pod Name 提取序号（如 myapp-0, myapp-1）
	if podName := os.Getenv("POD_NAME"); podName != "" {
		if id := extractOrdinalFromName(podName); id >= 0 {
			return id % maxWorkerID
		}
	}

	// 2. 从 K8s Hostname 提取序号
	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		if id := extractOrdinalFromName(hostname); id >= 0 {
			return id % maxWorkerID
		}
	}

	// 3. 尝试从环境变量获取
	envVars := []string{"WORKER_ID", "NODE_ID", "POD_ORDINAL"}
	for _, envVar := range envVars {
		if workerID := os.Getenv(envVar); workerID != "" {
			if id, err := mathx.ParseInt64(workerID); err == nil {
				return mathx.Abs(id) % maxWorkerID
			}
		}
	}

	// 4. 从主机名生成唯一的 Worker ID
	hostName := SafeGetHostName()
	hash := sha256.Sum256([]byte(hostName))
	hostNameHash := int64(binary.BigEndian.Uint64(hash[:8]))

	return mathx.Abs(hostNameHash) % maxWorkerID
}

// extractOrdinalFromName 从名称中提取序号
// 支持格式: myapp-0, myapp-1, myapp-statefulset-0 等
// 返回 -1 表示未找到序号
func extractOrdinalFromName(name string) int64 {
	// 从后往前查找最后一个 '-' 后的数字
	lastDash := strings.LastIndex(name, "-")
	if lastDash == -1 || lastDash == len(name)-1 {
		return -1
	}

	ordinalStr := name[lastDash+1:]
	if id, err := mathx.ParseInt64(ordinalStr); err == nil {
		return id
	}

	return -1
}

// GetDatacenterId 获取数据中心ID
// 优先级: 自定义环境变量 > 容器编排平台 > 数据中心标识 > 默认值1
// 注意：不使用 REGION/ZONE 等区域变量，因为同一区域内的多台机器值相同，无法区分节点
// 返回范围: 0 到 (maxDatacenterID-1)，默认 0-31（雪花算法标准）
func GetDatacenterId() int64 {
	maxDatacenterID := globalWorkerIDConfig.maxDatacenterID.Load()

	// 1. 优先使用自定义环境变量（数字类型，直接使用）
	customEnvVars := []string{"DATACENTER_ID", "DC_ID", "DATA_CENTER_ID"}
	for _, envVar := range customEnvVars {
		if dcID := os.Getenv(envVar); dcID != "" {
			if id, err := mathx.ParseInt64(dcID); err == nil {
				return mathx.Abs(id) % maxDatacenterID
			}
		}
	}

	// 2. 从容器编排平台环境变量推导（这些能区分不同节点/环境）
	orchestrationEnvVars := []string{
		// Kubernetes - 命名空间和集群名可以区分不同环境
		"KUBERNETES_NAMESPACE", "K8S_NAMESPACE", "POD_NAMESPACE",
		"KUBERNETES_CLUSTER_NAME", "K8S_CLUSTER_NAME",

		// Docker Swarm - 节点 ID 可以区分不同节点
		"DOCKER_SWARM_NODE_ID", "SWARM_NODE_ID",

		// Nomad - DC 和命名空间可以区分
		"NOMAD_DC", "NOMAD_NAMESPACE",

		// Mesos - 任务 ID 可以区分
		"MESOS_TASK_ID", "MARATHON_APP_ID",

		// OpenShift - 命名空间可以区分
		"OPENSHIFT_BUILD_NAMESPACE", "OPENSHIFT_DEPLOYMENT_NAMESPACE",
	}
	for _, envVar := range orchestrationEnvVars {
		if value := os.Getenv(envVar); value != "" {
			hash := sha256.Sum256([]byte(value))
			return int64(binary.BigEndian.Uint32(hash[:4])) % maxDatacenterID
		}
	}

	// 3. 从明确的数据中心/集群标识推导
	datacenterEnvVars := []string{
		// 数据中心/机房标识
		"DATACENTER", "DC", "IDC", "SITE_ID",

		// 集群标识
		"CLUSTER_ID",
	}
	for _, envVar := range datacenterEnvVars {
		if value := os.Getenv(envVar); value != "" {
			hash := sha256.Sum256([]byte(value))
			return int64(binary.BigEndian.Uint32(hash[:4])) % maxDatacenterID
		}
	}

	// 4. 默认返回 1
	return 1
}

// GetWorkerIdForSnowflake 获取适合雪花算法的 WorkerID（0-31，5位）
func GetWorkerIdForSnowflake() int64 {
	maxSnowflakeID := globalWorkerIDConfig.maxSnowflakeID.Load()
	return GetWorkerId() % maxSnowflakeID
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
	return mathx.IfGtZero(result, 1)
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
	cmd.Dir = mathx.IF(baseDir != "", baseDir, "")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stdout.Bytes(), fmt.Errorf("command failed: %s: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}
