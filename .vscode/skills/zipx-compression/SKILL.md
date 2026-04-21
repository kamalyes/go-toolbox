---
name: zipx-compression
description: 压缩解压工具，提供Gzip/Zlib的压缩解压、多级压缩、对象序列化压缩、智能检测解压。当需要对数据/对象进行Gzip或Zlib压缩解压、或自动检测压缩格式解压时使用。
---

# zipx - 压缩解压

提供Gzip/Zlib压缩解压、多级压缩、对象序列化压缩与智能检测解压。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/zipx"
```

Gzip压缩解压：
```go
compressed := zipx.GzipCompress(data)
decompressed := zipx.GzipDecompress(compressed)
```

对象压缩：
```go
compressed, err := zipx.GzipCompressObject[obj](myObj)
obj, err := zipx.GzipDecompressObject[obj](compressed)
```

## 完整API索引

### 函数

#### Gzip

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GzipCompress` | `func(data []byte) ([]byte, error)` | Gzip压缩 |
| `GzipDecompress` | `func(data []byte) ([]byte, error)` | Gzip解压 |
| `MultiGZipCompress` | `func(data []byte, levels ...int) ([]byte, error)` | 多级Gzip压缩 |
| `MultiGZipDecompress` | `func(data []byte) ([]byte, error)` | 多级Gzip解压 |
| `GzipCompressObject[T]` | `func(obj T) ([]byte, error)` | 对象序列化后Gzip压缩 |
| `GzipDecompressObject[T]` | `func(data []byte) (T, error)` | Gzip解压后反序列化对象 |
| `GzipCompressObjectWithInfo[T]` | `func(obj T) ([]byte, Stats, error)` | 带统计信息的对象压缩 |
| `GzipCompressObjectWithSize[T]` | `func(obj T, maxSize int) ([]byte, error)` | 带大小限制的对象压缩 |
| `IsGzipCompressed` | `func(data []byte) bool` | 判断是否Gzip压缩数据 |
| `GzipSmartDecompress` | `func(data []byte) ([]byte, error)` | 智能Gzip解压（自动检测） |
| `GzipSmartDecompressObject[T]` | `func(data []byte) (T, error)` | 智能Gzip对象解压 |

#### Zlib

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ZlibCompress` | `func(data []byte) ([]byte, error)` | Zlib压缩 |
| `ZlibDecompress` | `func(data []byte) ([]byte, error)` | Zlib解压 |
| `MultiZlibCompress` | `func(data []byte, levels ...int) ([]byte, error)` | 多级Zlib压缩 |
| `MultiZlibDecompress` | `func(data []byte) ([]byte, error)` | 多级Zlib解压 |
| `ZlibCompressObject[T]` | `func(obj T) ([]byte, error)` | 对象序列化后Zlib压缩 |
| `ZlibDecompressObject[T]` | `func(data []byte) (T, error)` | Zlib解压后反序列化对象 |
| `ZlibCompressObjectWithInfo[T]` | `func(obj T) ([]byte, Stats, error)` | 带统计信息的对象压缩 |
| `ZlibCompressObjectWithSize[T]` | `func(obj T, maxSize int) ([]byte, error)` | 带大小限制的对象压缩 |
| `IsZlibCompressed` | `func(data []byte) bool` | 判断是否Zlib压缩数据 |
| `ZlibSmartDecompress` | `func(data []byte) ([]byte, error)` | 智能Zlib解压（自动检测） |
| `ZlibSmartDecompressObject[T]` | `func(data []byte) (T, error)` | 智能Zlib对象解压 |

### 常量/变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `GzipPrefix` | []byte | Gzip数据前缀 |
| `GzipPrefixLen` | int | Gzip前缀长度 |
| `ZlibPrefix` | []byte | Zlib数据前缀 |
| `ZlibPrefixLen` | int | Zlib前缀长度 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `CompressResult` | 压缩结果类型 |

## 注意事项

- `IsGzipCompressed` / `IsZlibCompressed` 通过魔数前缀判断
- `MultiGZipCompress` 支持多级压缩，解压时需用 `MultiGZipDecompress`
- `GzipSmartDecompress` 自动检测是否压缩，未压缩则原样返回