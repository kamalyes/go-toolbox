---
name: osx-system
description: 系统工具包，提供OS检测、文件操作、哈希计算、环境变量、WorkerID配置、主机信息、运行时调用、资源监控。当需要判断操作系统、操作文件目录、获取系统信息、或监控系统资源时使用。
---

# osx - 系统工具

提供操作系统检测、文件操作、环境变量、哈希计算、主机信息、WorkerID配置与资源监控。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/osx"
```

OS检测：
```go
if osx.IsMac() { /* macOS特定逻辑 */ }
if osx.IsWindows() { /* Windows特定逻辑 */ }
```

文件操作：
```go
osx.MkdirIfNotExist("/path/to/dir")
osx.Copy("/src/file", "/dst/file")
```

## 完整API索引

### 函数

#### OS检测

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `IsMac` | `func() bool` | 判断是否macOS |
| `IsWindows` | `func() bool` | 判断是否Windows |
| `IsLinux` | `func() bool` | 判断是否Linux |
| `IsSupportedOS` | `func() bool` | 判断是否受支持的OS |

#### 文件操作

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `MkdirIfNotExist` | `func(path string) error` | 目录不存在则创建 |
| `Copy` | `func(src, dst string) error` | 复制文件 |
| `MkdirTemp` | `func(dir, prefix string) (string, error)` | 创建临时目录 |
| `JoinPaths` | `func(paths ...string) string` | 拼接路径 |
| `JoinURL` | `func(parts ...string) string` | 拼接URL |
| `ParseUrlPath` | `func(url string) string` | 解析URL路径 |
| `DirHasContent` | `func(path string) bool` | 目录是否有内容 |
| `GetDirFiles` | `func(path string) ([]string, error)` | 获取目录文件列表 |
| `FindFiles` | `func(root, pattern string) ([]string, error)` | 查找文件 |
| `FindFilesRecursive` | `func(root, pattern string) ([]string, error)` | 递归查找文件 |

#### 文件检查与操作

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `CheckImageExists` | `func(path string) bool` | 检查图片是否存在 |
| `SaveImage` | `func(path string, data []byte) error` | 保存图片 |
| `WriteContentToFile` | `func(path, content string) error` | 写入内容到文件 |
| `CreateIfNotExist` | `func(path string) error` | 文件不存在则创建 |
| `RemoveIfExist` | `func(path string) error` | 文件存在则删除 |
| `FileExists` | `func(path string) bool` | 判断文件是否存在 |
| `FileNameWithoutExt` | `func(path string) string` | 获取无扩展名的文件名 |
| `ComputeHashes` | `func(filePath string) (map[string]string, error)` | 计算文件哈希 |

#### 环境变量

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Getenv[T]` | `func(key string, defaultVal T) T` | 获取环境变量并转为泛型类型 |

#### WorkerID配置

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `SetMaxWorkerID` | `func(id int64)` | 设置最大WorkerID |
| `SetMaxDatacenterID` | `func(id int64)` | 设置最大数据中心ID |
| `SetMaxSnowflakeWorkerID` | `func(id int64)` | 设置最大Snowflake WorkerID |
| `GetMaxWorkerID` | `func() int64` | 获取最大WorkerID |
| `GetMaxDatacenterID` | `func() int64` | 获取最大数据中心ID |
| `GetMaxSnowflakeWorkerID` | `func() int64` | 获取最大Snowflake WorkerID |

#### 主机信息

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GetHostName` | `func() string` | 获取主机名 |
| `SafeGetHostName` | `func() string` | 安全获取主机名 |
| `HashUnixMicroCipherText` | `func() string` | 微秒级哈希密文 |
| `GetWorkerId` | `func() int64` | 获取WorkerId |
| `GetDatacenterId` | `func() int64` | 获取数据中心ID |
| `GetWorkerIdForSnowflake` | `func() int64` | 获取Snowflake WorkerId |

#### 哈希与运行时

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `StableHashSlot` | `func(s string, min, max int) int` | 稳定哈希槽位 |
| `GetRuntimeCaller` | `func(skip int) RunTimeCaller` | 获取运行时调用信息 |
| `Command` | `func(bin string, argv []string, baseDir string) error` | 执行命令 |

#### 资源监控

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewMonitor` | `func(threshold Threshold) *Monitor` | 创建监控器 |
| `NewAdvancedMonitor` | `func() *Monitor` | 创建高级监控器 |
| `GetMemoryStats` | `func() MemoryStats` | 获取内存统计 |
| `GetCurrentUsage` | `func() float64` | 获取当前内存使用率 |
| `GetCurrentSnapshot` | `func() Snapshot` | 获取当前快照 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `WorkerIDConfig` | WorkerID配置类型 |
| `RunTimeCaller` | 运行时调用信息类型 |
| `ThresholdLevel` | 阈值级别类型 |
| `MetricType` | 指标类型 |
| `Threshold` | 阈值类型 |
| `Snapshot` | 快照类型 |
| `GrowthRate` | 增长率类型 |
| `MonitorStats` | 监控统计类型 |
| `Monitor` | 监控器类型 |

### 常量/变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `OSMac` | string | macOS标识 |
| `OSWindows` | string | Windows标识 |
| `OSLinux` | string | Linux标识 |

## 注意事项

- `IsSupportedOS` 仅支持 macOS/Windows/Linux
- `ComputeHashes` 返回多种哈希算法（MD5/SHA1/SHA256等）的结果
- `NewMonitor` 和 `NewAdvancedMonitor` 用于内存资源监控