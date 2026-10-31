---
name: uuid-generation
description: UUID 与确定性 ID 生成，提供无连字符 UUID v4、按字段内容生成的确定性 MD5 ID、MD5 十六进制摘要与子串重复；当需要随机唯一 ID 或按业务字段生成稳定去重键时使用
---

# uuid - UUID 与确定性 ID

四个函数：随机 UUID、按字段内容生成的确定性 ID、MD5 摘要、子串重复

## 先分清两个语义完全不同的函数

| 函数 | 随机还是确定 | 相同输入会怎样 |
|---|---|---|
| `UUID()` | **随机**（UUID v4） | 每次都不同 |
| `UniqueID(fields...)` | **确定**（MD5） | **每次都相同** |

`UniqueID` 的名字容易让人以为它也是随机的，但它是把所有字段拼成字符串后取 MD5，
**同样的入参永远得到同样的结果**
这正是它的用途：做去重键、幂等键、缓存键——而不是做「新建一条记录的 ID」

只有在 `len(fields) == 0`，或者所有字段拼接并 `TrimSpace` 后是空串时，
它才退化成调用 `UUID()` 返回一个随机值。**这个退化是个坑**：
传了字段但全是空串，你会拿到一个随机 ID 而不是稳定 ID，调用处通常察觉不到

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/uuid"
```

```go
id := uuid.UUID()                              // 32 位十六进制，无连字符
key := uuid.UniqueID(tenantID, "order", orderNo) // 确定性：同样入参同样结果
sum := uuid.Md5("hello")                       // 32 位十六进制 MD5
pad := uuid.SameSubStr("ab", 3)                // "ababab"
```

## 完整API索引

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `UUID` | `func() string` | google/uuid 的 v4，去掉连字符后的 32 位十六进制 |
| `UniqueID` | `func(fields ...interface{}) string` | 字段经 `convert.MustString` 拼接后取 MD5；空入参或拼接后为空串时回退到 `UUID()` |
| `Md5` | `func(src string) string` | MD5 的 32 位十六进制摘要 |
| `SameSubStr` | `func(subStr string, repeat int) string` | 把 `subStr` 重复 `repeat` 次 |

## 注意事项

- **`Md5` 和 `UniqueID` 不要用于任何安全目的。** MD5 早已不具备抗碰撞性，
  可以用来做标识、去重、缓存键，但不能用来校验完整性、签名、存口令
  需要安全摘要请用 `pkg/sign` 里的 SHA-256 及以上
- `UniqueID` 的字段拼接**没有分隔符**：`UniqueID("a","bc")` 与 `UniqueID("ab","c")`
  拼出来都是 `"abc"`，MD5 相同。字段边界有歧义时，自己加分隔符再传
- `UniqueID` 依赖 `convert.MustString` 对各类型的字符串化规则，
  入参类型变了（比如 `int` 改 `int64`）结果可能跟着变，当作持久化键时要留意
- `SameSubStr` 的行为与标准库 `strings.Repeat(subStr, repeat)` 一致，
  新代码直接用标准库即可；`repeat <= 0` 返回空串
- `UUID()` 目前是 `strings.ReplaceAll(id.String(), "-", "")`，
  即先生成 36 字符带连字符的串再去掉连字符，多付一次分配
  实测（M 系列，`-count=3`）现写法 254–268 ns、96 B、3 allocs，
  改成 `hex.EncodeToString(id[:])` 是 199–229 ns、48 B、2 allocs，输出完全一样的 32 字符
  绝对值差距不大（大头是 `uuid.New()` 里的 crypto/rand），
  但在每请求都生成 ID 的路径上这是白送的一次分配
