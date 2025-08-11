# Go-ToolbOX

> Go-toolbox 的特点在于日常工作需求和扩展开发，封装了通用的工具类

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/kamalyes/go-toolbox)
[![license](https://img.shields.io/github/license/kamalyes/go-toolbox)]()
[![download](https://img.shields.io/github/downloads/kamalyes/go-toolbox/total)]()
[![release](https://img.shields.io/github/v/release/kamalyes/go-toolbox)]()
[![commit](https://img.shields.io/github/last-commit/kamalyes/go-toolbox)]()
[![issues](https://img.shields.io/github/issues/kamalyes/go-toolbox)]()
[![pull](https://img.shields.io/github/issues-pr/kamalyes/go-toolbox)]()
[![fork](https://img.shields.io/github/forks/kamalyes/go-toolbox)]()
[![star](https://img.shields.io/github/stars/kamalyes/go-toolbox)]()
[![go](https://img.shields.io/github/go-mod/go-version/kamalyes/go-toolbox)]()
[![size](https://img.shields.io/github/repo-size/kamalyes/go-toolbox)]()
[![contributors](https://img.shields.io/github/contributors/kamalyes/go-toolbox)]()
[![codecov](https://codecov.io/gh/kamalyes/go-toolbox/branch/master/graph/badge.svg)](https://codecov.io/gh/kamalyes/go-toolbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamalyes/go-toolbox)](https://goreportcard.com/report/github.com/kamalyes/go-toolbox)
[![Go Reference](https://pkg.go.dev/badge/github.com/kamalyes/go-toolbox?status.svg)](https://pkg.go.dev/github.com/kamalyes/go-toolbox?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/kamalyes/go-toolbox/-/badge.svg)](https://sourcegraph.com/github.com/kamalyes/go-toolbox?badge)

### Go-toolbox 的主要特性

- **转换 (convert)**: 数据类型之间的转换，例如将字符串转换为整数，或将日期格式从一种形式转换为另一种形式

- **脱敏 (desensitize)**: 去除或模糊敏感信息，以防止数据泄露，例如，去除个人身份信息 (PII) 或对数据进行加密

- **CRC (crc)**: 循环冗余校验，用于检测数据传输中的错误

- **错误处理 (errorx)**: 提供增强的错误处理功能，简化错误管理

- **HTTP 扩展 (httpx)**: 提供 HTTP 请求和响应的辅助工具

- **图像处理 (imgix)**: 用于图像处理和操作的工具

- **JSON 处理 (json)**: 轻量级的数据交换格式的处理

- **位置服务 (location)**: 与 IP 区域等相关的信息

- **数学扩展 (mathx)**: 数字计算的扩展功能

- **时间处理 (moment)**: 解析、验证、操作和显示日期和时间，简化日期和时间的处理过程

- **操作系统接口 (osx)**: 与操作系统交互的编程接口

- **队列 (queue)**: 提供队列数据结构的实现

- **随机数 (random)**: 随机数生成器，适用于多种应用

- **重试机制 (retry)**: 在操作失败时重试该操作的过程，常用于网络请求和数据库操作，以提高系统的可靠性

- **调度 (schedule)**: 任务调度工具，支持定时任务的执行

- **签名 (sign)**: 数据的完整性和来源验证，用于验证数据的完整性和来源，例如单词签名和消息签名

- **SQL 构建器 (sqlbuilderx)**: 用于构建 SQL 查询的工具

- **字符串处理 (stringx)**: 字符串的扩展处理功能，提供字符串的扩展功能，如格式化、拆分、连接等

- **同步工具 (syncx)**: 提供并发编程中的同步工具

- **类型 (types)**: 提供各种类型的定义和操作

- **单位转换 (units)**: 单位之间的转换工具

- **用户代理 (useragent)**: 处理用户代理字符串的工具

- **UUID (uuid)**: 生成通用唯一识别码 (UUID)

- **验证器 (validator)**: 用于验证数据有效性的工具，例如，表单验证、数据格式验证等，以确保输入数据符合预期的格式和规则

- **压缩工具 (zipx)**: 与数据压缩和解压缩相关的工具

## 开始使用

建议需要 [Go](https://go.dev/) 版本 [1.20](https://go.dev/doc/devel/release#go1.20.0) 

### 获取

使用 [Go 的模块支持](https://go.dev/wiki/Modules#how-to-use-modules)，当您在代码中添加导入时，`go [build|run|test]` 将自动获取所需的依赖项：

```go
import "github.com/kamalyes/go-toolbox"
```

或者，使用 `go get` 命令：

```sh
go get -u github.com/kamalyes/go-toolbox
```
