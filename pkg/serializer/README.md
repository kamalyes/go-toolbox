# Go-Toolbox Serializer

高性能、类型安全的 Go 泛型序列化器，支持多种序列化格式和压缩算法。

## ✨ 特性

- 🚀 **高性能**: 基于对象池优化，减少内存分配
- 🔒 **类型安全**: 使用 Go 泛型，编译时类型检查
- 🗜️ **智能压缩**: 支持 Gzip/Zlib 压缩，最高可节省 98%+ 空间
- 🔄 **自动回退**: 智能格式检测和兼容性处理
- 📦 **多种格式**: 支持 JSON、GOB 等序列化格式
- ⚡ **Builder 模式**: 链式配置，易于使用
- 🛡️ **并发安全**: 支持高并发场景

## 📊 性能数据

| 数据大小 | 无压缩 | Gzip+GOB | Zlib+GOB | 压缩比 |
|---------|--------|-----------|-----------|--------|
| 小数据   | 1372字符 | 628字符   | 612字符   | 55.4% |
| 中等数据 | 8602字符 | 936字符   | 920字符   | 89.3% |
| 大数据   | 80902字符| 1356字符  | 1340字符  | **98.3%** |

## 🚀 快速开始

### 安装

```bash
go get github.com/kamalyes/go-toolbox/pkg/serializer
```

### 基础使用

```go
package main

import (
    "fmt"
    "github.com/kamalyes/go-toolbox/pkg/serializer"
)

type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    user := User{ID: "123", Name: "张三", Age: 30}
    
    // 创建序列化器
    s := serializer.NewCompact[User]()
    
    // 序列化
    encoded, err := s.EncodeToString(user)
    if err != nil {
        panic(err)
    }
    
    // 反序列化
    decoded, err := s.DecodeFromString(encoded)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("原始: %+v\\n", user)
    fmt.Printf("解码: %+v\\n", decoded)
    fmt.Printf("压缩后大小: %d 字符\\n", len(encoded))
}
```

## 🔧 预设配置

### 快速工厂方法

```go
// 最高压缩率 (Zlib+GOB+Base64)
serializer := serializer.NewZlibCompact[YourType]()

// 平衡性能 (Gzip+GOB+Base64)  
serializer := serializer.NewCompact[YourType]()

// 兼容性优先 (Gzip+JSON+Base64)
serializer := serializer.NewUltraCompact[YourType]()

// 最快速度 (GOB 无压缩)
serializer := serializer.NewFast[YourType]()

// 纯 JSON
serializer := serializer.NewJSON[YourType]()

// 标准 GOB
serializer := serializer.NewGob[YourType]()
```

### 自定义配置

```go
serializer := serializer.New[YourType]().
    WithType(serializer.TypeGob).
    WithCompression(serializer.CompressionGzip).
    WithBase64(true)
```

## 📋 配置选项

### 序列化类型

- `TypeJSON` - JSON 格式（跨语言兼容）
- `TypeGob` - Go 二进制格式（高效）

### 压缩类型

- `CompressionNone` - 无压缩
- `CompressionGzip` - Gzip 压缩
- `CompressionZlib` - Zlib 压缩（通常效果更好）

### Base64 编码

- `WithBase64(true)` - 启用 Base64 编码（字符串安全）
- `WithBase64(false)` - 原始二进制数据

## 🎯 使用场景

### 队列消息序列化

```go
type QueueMessage struct {
    MessageID string                 \`json:"message_id"\`
    Content   string                 \`json:"content"\`
    Metadata  map[string]interface{} \`json:"metadata"\`
}

// 创建专用序列化器
func NewQueueMessageSerializer() *serializer.Serializer[QueueMessage] {
    return serializer.NewZlibCompact[QueueMessage]()
}

// 使用示例
msg := QueueMessage{
    MessageID: "msg-001",
    Content:   "Hello, World!",
    Metadata:  map[string]interface{}{"priority": "high"},
}

s := NewQueueMessageSerializer()
encoded, _ := s.EncodeToString(msg)
decoded, _ := s.DecodeFromString(encoded)
```

### 缓存序列化

```go
type CacheData struct {
    Key       string    \`json:"key"\`
    Value     string    \`json:"value"\`
    ExpiresAt time.Time \`json:"expires_at"\`
}

// 缓存专用序列化器（速度优先）
cacheSerializer := serializer.NewFast[CacheData]()

// 大数据缓存序列化器（空间优先）
bigDataSerializer := serializer.NewCompact[CacheData]()
```

## 🧪 性能测试

运行性能测试：

```bash
cd pkg/serializer
go test -bench=. -benchmem
```

压缩效果测试：

```bash
go test -run="TestCompression" -v
```

## 📐 高级功能

### 自定义编解码器

```go
customSerializer := serializer.New[YourType]().
    WithCustomEncoder(func(obj YourType) ([]byte, error) {
        // 自定义编码逻辑
        return customEncode(obj)
    }).
    WithCustomDecoder(func(data []byte) (YourType, error) {
        // 自定义解码逻辑
        return customDecode(data)
    })
```

### 性能统计

```go
stats, err := serializer.GetStats(yourObject)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("JSON 大小: %d 字节\\n", stats.JSONSize)
fmt.Printf("GOB 大小: %d 字节\\n", stats.GobSize)
fmt.Printf("当前大小: %d 字节\\n", stats.CurrentSize)
fmt.Printf("压缩比: %.1f%%\\n", stats.CompressionRatio*100)
fmt.Printf("节省空间: %.1f%%\\n", stats.SpaceSavedPercent)
```

### 性能基准测试

```go
result, err := serializer.Benchmark(yourObject, 1000)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("编码时间: %v\\n", result.EncodeTime)
fmt.Printf("解码时间: %v\\n", result.DecodeTime)
fmt.Printf("数据大小: %d 字节\\n", result.DataSize)
```

## 🔍 错误处理

序列化器提供详细的错误信息：

```go
decoded, err := serializer.DecodeFromString(encoded)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "无法解码数据"):
        // 数据格式错误
        log.Printf("数据格式错误: %v", err)
    case strings.Contains(err.Error(), "压缩失败"):
        // 压缩相关错误
        log.Printf("压缩错误: %v", err)
    default:
        log.Printf("其他错误: %v", err)
    }
}
```

## 🛠️ 最佳实践

### 1. 选择合适的序列化器

```go
// 网络传输（空间敏感）
networkSerializer := serializer.NewZlibCompact[Message]()

// 本地缓存（速度敏感）  
cacheSerializer := serializer.NewFast[CacheItem]()

// 跨语言通信（兼容性敏感）
apiSerializer := serializer.NewJSON[APIResponse]()
```

### 2. 复用序列化器实例

```go
// ✅ 好的做法：复用实例
var messageSerializer = serializer.NewCompact[Message]()

func processMessage(msg Message) {
    encoded, _ := messageSerializer.EncodeToString(msg)
    // ...
}

// ❌ 避免：每次创建新实例
func badProcessMessage(msg Message) {
    serializer := serializer.NewCompact[Message]() // 浪费性能
    encoded, _ := serializer.EncodeToString(msg)
    // ...
}
```

### 3. 处理大数据

```go
// 对于大数据，优先使用压缩
bigDataSerializer := serializer.NewZlibCompact[LargeData]()

// 如果性能敏感，可以考虑并行处理
func processBatch(items []LargeData) {
    results := make([]string, len(items))
    
    // 并行序列化
    var wg sync.WaitGroup
    for i, item := range items {
        wg.Add(1)
        go func(index int, data LargeData) {
            defer wg.Done()
            encoded, _ := bigDataSerializer.EncodeToString(data)
            results[index] = encoded
        }(i, item)
    }
    wg.Wait()
}
```

## 📄 许可证

Copyright (c) 2025 by kamalyes, All Rights Reserved.

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

如有问题，请联系 kamalyes@qq.com