---
name: json-facade
description: JSON 实现门面，通过 build tag 在标准库/goccy-go-json/jsoniter/sonic 之间切换，并提供向已有 JSON 追加键、递归替换键名、紧凑化等便捷函数；当需要统一 JSON 入口或替换 JSON 实现时使用
---

# json - JSON 实现门面

把 `Marshal`/`Unmarshal`/`MarshalIndent`/`NewDecoder`/`NewEncoder` 转发给可切换的底层实现，
另外提供几个跨实现一致的便捷函数

## 实现由 build tag 选择

| build tag | 实际实现 | 文件 |
|---|---|---|
| 无（默认） | `encoding/json` 标准库 | `json.go` |
| `go_json` | `github.com/goccy/go-json` | `go_json.go` |
| `jsoniter` | `github.com/json-iterator/go`（兼容标准库配置） | `jsoniter.go` |
| `sonic` + `avx` + (linux\|windows\|darwin) + `amd64` | `github.com/bytedance/sonic` | `sonic.go` |

默认文件的约束是
`!jsoniter && !go_json && !(sonic && avx && (linux || windows || darwin) && amd64)`，
也就是「其他三个都没命中时兜底」

```bash
go build                      # 标准库
go build -tags go_json        # goccy
go build -tags jsoniter       # jsoniter
```

## sonic 这条路目前是坏的

实测在 Go 1.25 上，`GOOS=linux GOARCH=amd64 go build -tags='sonic,avx' ./pkg/json/`
会失败，有两个独立的问题：

1. **依赖层先炸**：`github.com/bytedance/sonic/loader@v0.1.1` 引用了
   `_func`、`moduledata`、`compilationUnit` 等 runtime 内部结构，这些在新版本 Go 里已经改名或改形，
   报一串 `undefined:`。这不是本仓库的代码问题，是 sonic v1.11.6 与 Go 1.25 不兼容
2. **即使第 1 条解决了也还会炸**：`sonic.go` 只定义了 5 个转发变量，
   没有 `MarshalWithExtraField` 与 `Compact`——另外三个实现文件都各自提供了这两个函数
   所以 sonic 约束生效时，这两个符号会 undefined

注意 **arm64 机器上加 `-tags sonic` 不会报错，但也完全没生效**：
sonic 的约束要求 `amd64`，不满足时 `sonic.go` 被排除，兜底的 `json.go` 反而被编进来
所以「本地加了标签能编过」不能证明 sonic 可用，必须指定 `GOARCH=amd64` 才算测到

要用 sonic 得先升 sonic 版本，并把缺的两个函数补进 `sonic.go`

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/json"
```

转发变量，用法与标准库完全一致：
```go
b, err := json.Marshal(v)
err = json.Unmarshal(b, &v)
b, err = json.MarshalIndent(v, "", "  ")
dec := json.NewDecoder(r)
enc := json.NewEncoder(w)
```

向已有 JSON 字符串追加键值对：
```go
pairs := json.NewKeyValuePairs().
    Add("trace_id", traceID).
    Add("ts", time.Now().Unix())
out, err := json.AppendKeysToJSONMarshal(`{"name":"x"}`, pairs)
// {"name":"x","trace_id":"...","ts":...}
```

递归替换键名中的子串（map 与 slice 都会递归进去）：
```go
replaced, err := json.ReplaceKeys(data, "_", "-")   // user_name -> user-name
```

## 完整API索引

### 转发变量（随 build tag 变化）

| 导出名称 | 等价于 | 说明 |
|---|---|---|
| `Marshal` | `json.Marshal` | 序列化 |
| `Unmarshal` | `json.Unmarshal` | 反序列化 |
| `MarshalIndent` | `json.MarshalIndent` | 带缩进序列化 |
| `NewDecoder` | `json.NewDecoder` | 创建流式解码器 |
| `NewEncoder` | `json.NewEncoder` | 创建流式编码器 |

它们是 **变量而不是函数**，所以可以在测试里临时替换，但也意味着不能用作泛型实参

### 便捷函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `AppendKeysToJSONMarshal` | `func(originalJSON string, pairs *KeyValuePairs) ([]byte, error)` | 把键值对并入已有 JSON；解析失败返回空切片与错误 |
| `ReplaceKeys` | `func(data interface{}, oldStr, newStr string) (interface{}, error)` | 递归替换 map 键名中的子串，slice 内元素一并递归 |
| `MarshalWithExtraField` | `func(v any, extraKey string, extraValue any) ([]byte, error)` | 序列化并附加一个额外字段 |
| `Compact` | `func(data []byte) string` | 去掉空白的紧凑化输出 |

### KeyValuePairs

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewKeyValuePairs` | `func() *KeyValuePairs` | 创建键值对集合 |
| `Add` | `func(key string, value interface{}) *KeyValuePairs` | 添加一对，返回自身支持链式 |

## 注意事项

- `AppendKeysToJSONMarshal` 内部是「反序列化成 `map[string]interface{}` → 合并 → 再序列化」，
  因此**不保留原 JSON 的键顺序**，数字也会走 `float64`；需要保序或精确数字请自己拼接
- `KeyValuePairs` 底层是 map，`Add` 同一个键会覆盖；map 也意味着追加的那几个键之间无序
- `ReplaceKeys` 替换的是**键名**，不动值；它对 `map[string]interface{}` 与 `[]interface{}`
  递归，其他类型原样返回
- 切换 build tag 会改变序列化的细节行为（HTML 转义、浮点格式、错误文案等）。
  依赖这些细节的测试在不同 tag 下可能结果不同，断言请对着解析后的结构而不是字节串
- 默认（无标签）路径的测试是绿的；换 tag 后请自己重跑一遍 `go test -tags <tag> ./pkg/json/`
