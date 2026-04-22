---
name: errorx-error-handling
description: 错误处理体系，提供错误分类、错误类型注册、错误包装、预定义错误常量。当需要构建结构化错误、按类型分类错误、或使用预定义错误码时使用。
---

# errorx - 错误处理体系

提供错误类型注册、错误分类与包装、预定义错误常量，构建结构化错误处理链。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/errorx"
```

创建结构化错误：

```go
err := errorx.NewBaseError("order not found", errorx.ErrorTypeBusiness)
wrapped := errorx.WrapError("process order", err)
```

注册与分类：

```go
errorx.RegisterError(1001, "user not found")
classified := errorx.ClassifyError(err)
```

## 完整API索引

### 函数

| 导出名称         | 签名                                                          | 说明                     |
| ---------------- | ------------------------------------------------------------- | ------------------------ |
| `WrapError`      | `func(message string, err ...error) error`                    | 包装错误并附加消息       |
| `WrapTypedError` | `func(errType ErrorType, message string, err ...error) error` | 包装错误并附加类型和消息 |
| `NewBaseError`   | `func(msg string, errTypes ...ErrorType) *BaseError`          | 创建带类型的基础错误     |
| `RegisterError`  | `func(errType ErrorType, msg string)`                         | 注册错误码与描述         |
| `NewError`       | `func(errType ErrorType, args ...interface{}) error`          | 按类型创建已注册错误     |
| `ClassifyError`  | `func(err error) ErrorType`                                   | 按类型分类错误           |
| `PrintErrorMap`  | `func()`                                                      | 打印错误映射             |
| `GetErrorMap`    | `func() ErrorMapType`                                         | 获取错误映射             |
| `ResetErrorMap`  | `func()`                                                      | 重置错误映射             |
| `New`            | `func(message string) error`                                  | 快捷创建错误             |
| `Newf`           | `func(format string, args ...interface{}) error`              | 格式化创建错误           |

### 类型

| 导出名称       | 说明                                        |
| -------------- | ------------------------------------------- |
| `BaseError`    | 结构化错误基础类型，含 `Msg` 和 `Type` 字段 |
| `ErrorType`    | 错误分类枚举类型                            |
| `ErrorMapType` | 错误映射类型                                |

### BaseError 方法

| 导出名称            | 签名               | 说明         |
| ------------------- | ------------------ | ------------ |
| `BaseError.Error`   | `func() string`    | 返回错误消息 |
| `BaseError.GetType` | `func() ErrorType` | 返回错误类型 |

## 注意事项

- `ClassifyError` 依赖已注册的错误类型，未注册的错误返回默认分类
- `WrapError` 会保留原始错误的完整堆栈信息
- `NewError` 需先 `RegisterError` 注册类型，否则返回空消息
