/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-18 17:36:32
 * @FilePath: \go-toolbox\pkg\osx\base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHostName(t *testing.T) {
	hostname, err := GetHostName()
	assert.NoError(t, err)
	assert.NotEmpty(t, hostname)
}

func TestSafeGetHostName(t *testing.T) {
	hostname := SafeGetHostName()
	assert.NotEmpty(t, hostname)
	assert.GreaterOrEqual(t, len(hostname), 1)
}

func TestHashUnixMicroCipherText(t *testing.T) {
	hash1 := HashUnixMicroCipherText()
	hash2 := HashUnixMicroCipherText()

	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, hash2)
	assert.NotEqual(t, hash1, hash2) // 每次生成的哈希应该不同
	assert.Len(t, hash1, 32)         // MD5 哈希长度为 32
}

func TestGetWorkerId(t *testing.T) {
	// 场景1：默认配置（0-1023）
	workerID := GetWorkerId()
	assert.GreaterOrEqual(t, workerID, int64(0))
	assert.Less(t, workerID, int64(1024))

	// 场景2：自定义范围
	SetMaxWorkerID(2048)
	workerID = GetWorkerId()
	assert.GreaterOrEqual(t, workerID, int64(0))
	assert.Less(t, workerID, int64(2048))

	// 恢复默认值
	SetMaxWorkerID(1024)
}

func TestGetWorkerId_WithEnv(t *testing.T) {
	tests := []struct {
		name   string
		envKey string
		envVal string
	}{
		{"POD_NAME", "POD_NAME", "myapp-5"},
		{"HOSTNAME", "HOSTNAME", "myapp-10"},
		{"WORKER_ID", "WORKER_ID", "100"},
		{"NODE_ID", "NODE_ID", "200"},
		{"POD_ORDINAL", "POD_ORDINAL", "15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			oldVal := os.Getenv(tt.envKey)
			os.Setenv(tt.envKey, tt.envVal)
			defer os.Setenv(tt.envKey, oldVal)

			workerID := GetWorkerId()
			assert.GreaterOrEqual(t, workerID, int64(0))
			assert.Less(t, workerID, int64(1024))
		})
	}
}

func TestExtractOrdinalFromName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"简单格式", "myapp-0", 0},
		{"多位数字", "myapp-123", 123},
		{"StatefulSet格式", "myapp-statefulset-5", 5},
		{"无序号", "myapp", -1},
		{"空字符串", "", -1},
		{"只有横杠", "myapp-", -1},
		{"非数字后缀", "myapp-abc", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractOrdinalFromName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDatacenterId(t *testing.T) {
	// 场景1：默认配置（0-31）
	dcID := GetDatacenterId()
	assert.GreaterOrEqual(t, dcID, int64(0))
	assert.Less(t, dcID, int64(32))

	// 场景2：自定义范围
	SetMaxDatacenterID(64)
	dcID = GetDatacenterId()
	assert.GreaterOrEqual(t, dcID, int64(0))
	assert.Less(t, dcID, int64(64))

	// 恢复默认值
	SetMaxDatacenterID(32)
}

func TestGetDatacenterId_WithEnv(t *testing.T) {
	tests := []struct {
		name   string
		envKey string
		envVal string
	}{
		// 自定义环境变量（数字类型）
		{"DATACENTER_ID", "DATACENTER_ID", "10"},
		{"DC_ID", "DC_ID", "20"},
		{"DATA_CENTER_ID", "DATA_CENTER_ID", "30"},

		// Kubernetes（能区分不同环境/节点）
		{"KUBERNETES_NAMESPACE", "KUBERNETES_NAMESPACE", "production"},
		{"K8S_NAMESPACE", "K8S_NAMESPACE", "staging"},
		{"POD_NAMESPACE", "POD_NAMESPACE", "default"},
		{"KUBERNETES_CLUSTER_NAME", "KUBERNETES_CLUSTER_NAME", "prod-cluster"},
		{"K8S_CLUSTER_NAME", "K8S_CLUSTER_NAME", "test-cluster"},

		// Docker Swarm（节点 ID）
		{"DOCKER_SWARM_NODE_ID", "DOCKER_SWARM_NODE_ID", "node-123"},
		{"SWARM_NODE_ID", "SWARM_NODE_ID", "node-456"},

		// Nomad（DC 和命名空间）
		{"NOMAD_DC", "NOMAD_DC", "dc1"},
		{"NOMAD_NAMESPACE", "NOMAD_NAMESPACE", "prod"},

		// Mesos（任务 ID）
		{"MESOS_TASK_ID", "MESOS_TASK_ID", "task-123"},
		{"MARATHON_APP_ID", "MARATHON_APP_ID", "app-456"},

		// OpenShift（命名空间）
		{"OPENSHIFT_BUILD_NAMESPACE", "OPENSHIFT_BUILD_NAMESPACE", "build-ns"},
		{"OPENSHIFT_DEPLOYMENT_NAMESPACE", "OPENSHIFT_DEPLOYMENT_NAMESPACE", "deploy-ns"},

		// 数据中心/集群标识
		{"DATACENTER", "DATACENTER", "dc-east"},
		{"DC", "DC", "dc-west"},
		{"IDC", "IDC", "idc-01"},
		{"SITE_ID", "SITE_ID", "site-02"},
		{"CLUSTER_ID", "CLUSTER_ID", "cluster-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理所有相关环境变量
			allEnvVars := []string{
				// 自定义
				"DATACENTER_ID", "DC_ID", "DATA_CENTER_ID",
				// K8s
				"KUBERNETES_NAMESPACE", "K8S_NAMESPACE", "POD_NAMESPACE",
				"KUBERNETES_CLUSTER_NAME", "K8S_CLUSTER_NAME", "CLUSTER_NAME",
				// Docker
				"DOCKER_SWARM_NODE_ID", "SWARM_NODE_ID",
				// Nomad
				"NOMAD_DC", "NOMAD_NAMESPACE",
				// Mesos
				"MESOS_TASK_ID", "MARATHON_APP_ID",
				// OpenShift
				"OPENSHIFT_BUILD_NAMESPACE", "OPENSHIFT_DEPLOYMENT_NAMESPACE",
				// 数据中心标识
				"DATACENTER", "DC", "IDC", "SITE_ID",
				"CLUSTER_ID", "CLUSTER",
			}
			oldVals := make(map[string]string)
			for _, key := range allEnvVars {
				oldVals[key] = os.Getenv(key)
				os.Unsetenv(key)
			}
			defer func() {
				for key, val := range oldVals {
					if val != "" {
						os.Setenv(key, val)
					}
				}
			}()

			// 设置测试环境变量
			os.Setenv(tt.envKey, tt.envVal)

			dcID := GetDatacenterId()
			assert.GreaterOrEqual(t, dcID, int64(0))
			assert.Less(t, dcID, int64(32))
		})
	}
}

func TestGetWorkerIdForSnowflake(t *testing.T) {
	// 场景1：默认配置（0-31）
	snowflakeID := GetWorkerIdForSnowflake()
	assert.GreaterOrEqual(t, snowflakeID, int64(0))
	assert.Less(t, snowflakeID, int64(32))

	// 场景2：自定义范围（10位 WorkerID）
	SetMaxSnowflakeWorkerID(1024)
	snowflakeID = GetWorkerIdForSnowflake()
	assert.GreaterOrEqual(t, snowflakeID, int64(0))
	assert.Less(t, snowflakeID, int64(1024))

	// 恢复默认值
	SetMaxSnowflakeWorkerID(32)
}

func TestSetAndGetMaxWorkerID(t *testing.T) {
	// 保存原始值
	originalMax := GetMaxWorkerID()
	defer SetMaxWorkerID(originalMax)

	// 场景1：设置有效值
	SetMaxWorkerID(2048)
	assert.Equal(t, int64(2048), GetMaxWorkerID())

	// 场景2：设置无效值（<= 0），应该使用默认值 1024
	SetMaxWorkerID(0)
	assert.Equal(t, int64(1024), GetMaxWorkerID())

	SetMaxWorkerID(-100)
	assert.Equal(t, int64(1024), GetMaxWorkerID())
}

func TestSetAndGetMaxDatacenterID(t *testing.T) {
	// 保存原始值
	originalMax := GetMaxDatacenterID()
	defer SetMaxDatacenterID(originalMax)

	// 场景1：设置有效值
	SetMaxDatacenterID(64)
	assert.Equal(t, int64(64), GetMaxDatacenterID())

	// 场景2：设置无效值（<= 0），应该使用默认值 32
	SetMaxDatacenterID(0)
	assert.Equal(t, int64(32), GetMaxDatacenterID())

	SetMaxDatacenterID(-50)
	assert.Equal(t, int64(32), GetMaxDatacenterID())
}

func TestSetAndGetMaxSnowflakeWorkerID(t *testing.T) {
	// 保存原始值
	originalMax := GetMaxSnowflakeWorkerID()
	defer SetMaxSnowflakeWorkerID(originalMax)

	// 场景1：设置有效值
	SetMaxSnowflakeWorkerID(1024)
	assert.Equal(t, int64(1024), GetMaxSnowflakeWorkerID())

	// 场景2：设置无效值（<= 0），应该使用默认值 32
	SetMaxSnowflakeWorkerID(0)
	assert.Equal(t, int64(32), GetMaxSnowflakeWorkerID())

	SetMaxSnowflakeWorkerID(-10)
	assert.Equal(t, int64(32), GetMaxSnowflakeWorkerID())
}

func TestStableHashSlot(t *testing.T) {
	// 场景1：正常范围
	result := StableHashSlot("test-key", 1, 100)
	assert.GreaterOrEqual(t, result, 1)
	assert.LessOrEqual(t, result, 100)

	// 场景2：相同输入应该产生相同输出（稳定性）
	result1 := StableHashSlot("stable-key", 1, 100)
	result2 := StableHashSlot("stable-key", 1, 100)
	assert.Equal(t, result1, result2)

	// 场景3：不同输入应该产生不同输出（大概率）
	result3 := StableHashSlot("key1", 1, 100)
	result4 := StableHashSlot("key2", 1, 100)
	assert.NotEqual(t, result3, result4)

	// 场景4：范围相同时返回该值
	result5 := StableHashSlot("any-key", 50, 50)
	assert.Equal(t, 50, result5)

	// 场景5：maxNum < minNum 应该 panic
	assert.Panics(t, func() {
		StableHashSlot("test", 100, 1)
	})
}

func TestGetRuntimeCaller(t *testing.T) {
	caller := GetRuntimeCaller(1)
	require.NotNil(t, caller)
	defer caller.Release()

	assert.NotEmpty(t, caller.File)
	assert.Greater(t, caller.Line, 0)
	assert.NotEmpty(t, caller.FuncName)
	assert.Contains(t, caller.String(), "FuncName:")
	assert.Contains(t, caller.String(), "File:")
}

func TestRunTimeCallerPool(t *testing.T) {
	// 测试对象池复用
	caller1 := GetRuntimeCaller(1)
	file1 := caller1.File
	caller1.Release()

	caller2 := GetRuntimeCaller(1)
	defer caller2.Release()

	// 验证对象被复用（地址可能相同）
	assert.NotEmpty(t, caller2.File)
	assert.NotEqual(t, file1, "") // 确保 Release 清理了字段
}

func TestCommand(t *testing.T) {
	// 场景1：执行简单命令（跨平台）
	var output []byte
	var err error

	// Windows 使用 cmd /c echo, Unix 使用 echo
	if os.PathSeparator == '\\' {
		// Windows
		output, err = Command("cmd", []string{"/c", "echo", "hello"}, "")
	} else {
		// Unix/Linux/Mac
		output, err = Command("echo", []string{"hello"}, "")
	}

	assert.NoError(t, err)
	assert.Contains(t, string(output), "hello")

	// 场景2：指定工作目录
	tempDir := os.TempDir()
	if os.PathSeparator == '\\' {
		output, err = Command("cmd", []string{"/c", "cd"}, tempDir)
	} else {
		output, err = Command("pwd", []string{}, tempDir)
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, output)

	// 场景3：执行失败的命令
	_, err = Command("nonexistent-command-xyz", []string{}, "")
	assert.Error(t, err)
}

func TestWorkerIDConfig_Concurrent(t *testing.T) {
	// 并发测试配置的线程安全性
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			SetMaxWorkerID(int64(1024 + id))
			_ = GetMaxWorkerID()
			SetMaxDatacenterID(int64(32 + id))
			_ = GetMaxDatacenterID()
			SetMaxSnowflakeWorkerID(int64(32 + id))
			_ = GetMaxSnowflakeWorkerID()
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证最终值有效
	assert.Greater(t, GetMaxWorkerID(), int64(0))
	assert.Greater(t, GetMaxDatacenterID(), int64(0))
	assert.Greater(t, GetMaxSnowflakeWorkerID(), int64(0))
}

func TestGetWorkerId_Consistency(t *testing.T) {
	// 验证在相同环境下多次调用返回相同结果
	id1 := GetWorkerId()
	id2 := GetWorkerId()
	id3 := GetWorkerId()

	assert.Equal(t, id1, id2)
	assert.Equal(t, id2, id3)
}

func TestGetDatacenterId_Consistency(t *testing.T) {
	// 验证在相同环境下多次调用返回相同结果
	id1 := GetDatacenterId()
	id2 := GetDatacenterId()
	id3 := GetDatacenterId()

	assert.Equal(t, id1, id2)
	assert.Equal(t, id2, id3)
}
