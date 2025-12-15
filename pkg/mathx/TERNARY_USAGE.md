# 三元运算函数使用指南 (Ternary Usage Guide)

## 📚 目录

- [基础三元运算](#基础三元运算)
- [条件执行（副作用）](#条件执行副作用)
- [延迟执行](#延迟执行)
- [错误处理](#错误处理)
- [多条件判断](#多条件判断)
- [空值检查](#空值检查)
- [集合操作](#集合操作)
- [类型转换与验证](#类型转换与验证)
- [链式构建器](#链式构建器)
- [特殊场景](#特殊场景)

---

## 基础三元运算

### 场景 1：简单条件返回值

**推荐使用：`IF`**

```go
// ❌ 传统写法
var status string
if user.IsActive {
    status = "在线"
} else {
    status = "离线"
}

// ✅ 推荐写法
status := mathx.IF(user.IsActive, "在线", "离线")
```

**适用场景：**

- 简单的条件判断返回值
- 等同于三元运算符 `condition ? trueVal : falseVal`
- 最常用的基础函数

---

## 条件执行（副作用）

### 场景 2：只需要执行操作，不需要返回值

**推荐使用：`IfExec`**

```go
// ❌ 传统写法
if err != nil {
    logger.Error("操作失败: %v", err)
}

// ✅ 推荐写法
mathx.IfExec(err != nil, func() {
    logger.Error("操作失败: %v", err)
})
```

**适用场景：**

- 日志记录
- 通知发送
- 统计计数

---

### 场景 3：根据条件执行不同操作

**推荐使用：`IfExecElse`**

```go
// ❌ 传统写法
if err == nil {
    logger.Info("成功")
} else {
    logger.Error("失败: %v", err)
}

// ✅ 推荐写法（完整版本）
mathx.IfExecElse(err == nil,
    func() { logger.Info("成功") },
    func() { logger.Error("失败: %v", err) },
)

// ✅ 推荐写法（省略 false 分支，可选参数）
mathx.IfExecElse(needLog,
    func() { logger.Info("处理完成") },
)
```

**可变参数说明：**

- `onFalse` 参数可选，可以省略或传 `nil`
- 省略时，条件为 `false` 时不执行任何操作

**适用场景：**

- 需要明确的成功/失败处理
- 双分支副作用操作
- 只需要 true 分支的场景（省略 false 分支）

---

### 场景 4：需要传递结果和错误给回调

**推荐使用：`IfCall`**

```go
// ✅ 完整版本（两个回调）
mathx.IfCall(err != nil, result, err,
    func(r T, e error) { onSuccess(r) },
    func(r T, e error) { onError(e) },
)

// ✅ 只需要 true 分支
mathx.IfCall(success, data, nil,
    func(r T, e error) { log.Info("成功: %v", r) },
)

// ✅ 只需要 false 分支（第一个回调传 nil）
mathx.IfCall(err != nil, nil, err,
    nil,
    func(r T, e error) { log.Error("错误: %v", e) },
)

// ✅ 不提供回调（只做条件判断）
mathx.IfCall(condition, value, err)
```

**可变参数说明：**

- `callbacks` 可以传 0-2 个回调函数
- `callbacks[0]` 是 true 分支，`callbacks[1]` 是 false 分支
- 可以省略任意回调，或传 `nil`

**适用场景：**

- 异步回调场景
- 需要同时传递结果和错误

---

## 延迟执行

### 场景 5：避免不必要的计算（惰性求值）

**推荐使用：`IfDoAF`**

```go
// ❌ 传统写法（两个函数都会执行）
var result string
if condition {
    result = expensiveComputation()  // 昂贵计算
} else {
    result = cheapDefault()
}

// ✅ 推荐写法（只执行需要的函数）
result := mathx.IfDoAF(condition,
    func() string { return expensiveComputation() },  // 仅 condition=true 时执行
    func() string { return cheapDefault() },          // 仅 condition=false 时执行
)
```

**适用场景：**

- 计算成本高的操作
- 数据库查询
- 网络请求
- 避免提前求值

**⚠️ 已废弃：** `IfLazy` → 请使用 `IfDoAF`

---

### 场景 6：条件执行单个延迟函数

**推荐使用：`IfDo`**

```go
// ✅ 单函数延迟执行
result := mathx.IfDo(needCompute,
    func() int { return heavyCalculation() },
    0, // 默认值
)
```

**适用场景：**

- 条件满足时才执行计算
- 有明确的默认值

---

## 错误处理

### 场景 7：带错误返回的延迟执行

**推荐使用：`IfDoWithError`**

```go
// ❌ 传统写法
var result string
var err error
if shouldProcess {
    result, err = processData()
} else {
    result = ""
    err = nil
}

// ✅ 推荐写法
result, err := mathx.IfDoWithError(shouldProcess,
    func() (string, error) {
        return processData()
    },
    "", // 默认值
)
```

**适用场景：**

- 可能返回错误的操作
- 数据库操作
- 文件读写

---

### 场景 8：错误时返回默认值（不关心错误）

**推荐使用：`IfDoWithErrorDefault`**

```go
// ✅ 忽略错误，返回默认值
result := mathx.IfDoWithErrorDefault(condition,
    func() (int, error) { return parseValue() },
    0, // 默认值
)
```

**适用场景：**

- 允许失败的操作
- 降级处理

---

### 场景 9：简化错误检查

**推荐使用：`ReturnIfErr`**

```go
// ✅ 简化错误返回
return mathx.ReturnIfErr(result, err)
// 等同于：
// if err != nil {
//     return zero, err
// }
// return result, nil
```

**适用场景：**

- 函数末尾错误检查
- 减少样板代码

---

## 多条件判断

### 场景 10：多个条件，返回第一个满足的值

**推荐使用：`IfElse`**

```go
// ❌ 传统写法
var status string
if score >= 90 {
    status = "优秀"
} else if score >= 60 {
    status = "及格"
} else if score >= 0 {
    status = "不及格"
} else {
    status = "无效"
}

// ✅ 推荐写法
status := mathx.IfElse(
    []bool{score >= 90, score >= 60, score >= 0},
    []string{"优秀", "及格", "不及格"},
    "无效", // 默认值
)
```

**适用场景：**

- 多级条件判断
- 类似 switch-case 逻辑
- 评分、等级划分

**⚠️ 已废弃：** `IfDefault` → 请使用 `IfElse`

---

### 场景 11：结构化多条件判断

**推荐使用：`IfChain`**

```go
// ✅ 结构化条件
result := mathx.IfChain([]mathx.ConditionValue[string]{
    {Cond: x > 0, Value: "正数"},
    {Cond: x == 0, Value: "零"},
    {Cond: x < 0, Value: "负数"},
}, "未知")
```

**适用场景：**

- 条件与值需要配对
- 代码可读性优先

---

### 场景 12：开关式映射

**推荐使用：`IfSwitch`**

```go
// ❌ 传统写法
var message string
switch statusCode {
case 200:
    message = "成功"
case 404:
    message = "未找到"
case 500:
    message = "服务器错误"
default:
    message = "未知状态"
}

// ✅ 推荐写法
message := mathx.IfSwitch(statusCode, map[int]string{
    200: "成功",
    404: "未找到",
    500: "服务器错误",
}, "未知状态")
```

**适用场景：**

- 状态码映射
- 枚举值转换

---

## 空值检查

### 场景 13：指针空值检查

**推荐使用：`IfNotNil`**

```go
// ❌ 传统写法
var value int
if ptr != nil {
    value = *ptr
} else {
    value = defaultValue
}

// ✅ 推荐写法
value := mathx.IfNotNil(ptr, defaultValue)
```

**适用场景：**

- 指针类型安全访问
- 避免 nil panic

---

### 场景 14：字符串空值检查

**推荐使用：`IfNotEmpty`**

```go
// ✅ 字符串默认值
name := mathx.IfNotEmpty(user.Name, "匿名用户")
```

**适用场景：**

- 字符串默认值设置
- 配置项回退

---

### 场景 15：零值检查

**推荐使用：`IfNotZero`**

```go
// ✅ 零值检查（支持任意可比较类型）
timeout := mathx.IfNotZero(config.Timeout, 30)
```

**适用场景：**

- 数值配置默认值
- 任意可比较类型零值检查

---

### 场景 16：错误或空值组合检查

**推荐使用：`IfErrOrNil`**

```go
// ✅ 组合检查
message := mathx.IfErrOrNil(val, err, "失败", "成功")
```

**适用场景：**

- 同时检查错误和值
- 简化双重判断

---

## 集合操作

### 场景 17：检查切片是否包含元素

**推荐使用：`IfContains`**

```go
// ❌ 传统写法
var message string
found := false
for _, p := range permissions {
    if p == "admin" {
        found = true
        break
    }
}
if found {
    message = "有权限"
} else {
    message = "无权限"
}

// ✅ 推荐写法
message := mathx.IfContains(permissions, "admin", "有权限", "无权限")
```

**适用场景：**

- 权限检查
- 标签匹配

---

### 场景 18：任意条件满足

**推荐使用：`IfAny`**

```go
// ✅ 任一满足
canAccess := mathx.IfAny(
    []bool{isAdmin, isOwner, hasPermission},
    true,
    false,
)
```

**适用场景：**

- OR 逻辑
- 多权限检查

---

### 场景 19：所有条件满足

**推荐使用：`IfAll`**

```go
// ✅ 全部满足
isValid := mathx.IfAll(
    []bool{hasName, hasEmail, agreeTerms},
    true,
    false,
)
```

**适用场景：**

- AND 逻辑
- 表单验证

---

### 场景 20：条件计数

**推荐使用：`IfCount`**

```go
// ✅ 满足数量判断
level := mathx.IfCount(
    []bool{hasA, hasB, hasC},
    2, // 阈值
    "高级",
    "普通",
)
```

**适用场景：**

- 等级划分
- 积分系统

---

### 场景 21：切片过滤

**推荐使用：`IfFilter`**

```go
// ✅ 条件过滤
filtered := mathx.IfFilter(needFilter, users, func(u User) bool {
    return u.IsActive
})
```

**适用场景：**

- 动态过滤
- 搜索功能

---

### 场景 22：安全索引访问

**推荐使用：`IfSafeIndex`**

```go
// ❌ 传统写法
var item string
if index >= 0 && index < len(slice) {
    item = slice[index]
} else {
    item = defaultItem
}

// ✅ 推荐写法
item := mathx.IfSafeIndex(slice, index, defaultItem)
```

**适用场景：**

- 避免 panic
- 数组边界安全

---

### 场景 23：安全字典访问

**推荐使用：`IfSafeKey`**

```go
// ✅ 防空值访问
value := mathx.IfSafeKey(cache, key, defaultValue)
```

**适用场景：**

- map 安全读取
- 缓存降级

---

### 场景 24：切片长度判断

**推荐使用：`IfEmptySlice` / `IfLenGt` / `IfLenEq`**

```go
// ✅ 空切片检查
message := mathx.IfEmptySlice(items, "列表为空", "有数据")

// ✅ 长度大于
tip := mathx.IfLenGt(results, 0, "查询成功", "无结果")

// ✅ 长度等于
status := mathx.IfLenEq(queue, 1, "单任务", "多任务")
```

**适用场景：**

- 列表状态判断
- 批量操作提示

---

## 类型转换与验证

### 场景 25：类型断言

**推荐使用：`IfCast`**

```go
// ✅ 安全类型转换
str := mathx.IfCast[string](value, "默认值")
```

**适用场景：**

- interface{} 类型转换
- 类型安全降级

---

### 场景 26：验证函数

**推荐使用：`IfValidate`**

```go
// ✅ 验证逻辑
message := mathx.IfValidate(email,
    func(s string) bool { return strings.Contains(s, "@") },
    "邮箱有效",
    "邮箱无效",
)
```

**适用场景：**

- 输入验证
- 数据校验

---

### 场景 27：解析尝试

**推荐使用：`IfTryParse`**

```go
// ✅ 解析失败返回默认值
num := mathx.IfTryParse("123",
    func(s string) (int, error) { return strconv.Atoi(s) },
    0, // 默认值
)
```

**适用场景：**

- 字符串解析
- 容错处理

---

### 场景 28：数值区间判断

**推荐使用：`IfBetween`**

```go
// ✅ 区间检查
level := mathx.IfBetween(score, 60, 100, "及格", "不及格")
```

**适用场景：**

- 分数判断
- 范围验证

**⚠️ 已删除：** `IfInRange` (实现不正确，请使用 `IfBetween`)

---

## 链式构建器

### 场景 29：复杂条件链（副作用）

**推荐使用：`When().Then().Else().Do()`**

```go
// ✅ 链式副作用
mathx.When(err != nil).
    Then(func() { log.Error("失败") }).
    Else(func() { log.Info("成功") }).
    Do()
```

**适用场景：**

- 日志链式调用
- 清晰的条件分支

---

### 场景 30：链式返回值

**推荐使用：`WhenValue().ThenReturn().ElseReturn().Get()`**

```go
// ✅ 链式返回值
result := mathx.WhenValue[int](score >= 60).
    ThenReturn(100).
    ElseReturn(0).
    Get()
```

**适用场景：**

- 可读性优先
- 复杂条件判断

---

### 场景 31：多级条件提前返回

**推荐使用：`IFChainFor().When().ThenReturn()`**

```go
// ❌ 传统写法（深度嵌套）
func validateUser(name, email string, age int) error {
    if name == "" {
        return errors.New("名称为空")
    }
    if age < 0 {
        return errors.New("年龄无效")
    }
    if email == "" {
        return errors.New("邮箱为空")
    }
    return nil
}

// ✅ 推荐写法（链式调用，更清晰）
err := mathx.IFChainFor[error]().
    When(name == "").ThenReturn(errors.New("名称为空")).
    When(age < 0).ThenReturn(errors.New("年龄无效")).
    When(email == "").ThenReturn(errors.New("邮箱为空")).
    ExecuteOr(nil)
```

**适用场景：**

- 参数验证
- 提前退出逻辑
- 避免深度嵌套

---

### 场景 32：错误链式构建

**推荐使用：`IFErrorChain()`**

```go
// ✅ 错误处理链
err := mathx.IFErrorChain().
    When(user == nil).ThenReturn(ErrUserNotFound).
    When(!user.IsActive).ThenReturn(ErrUserInactive).
    ExecuteOr(nil)
```

**适用场景：**

- 错误验证链
- 业务规则检查

---

## 特殊场景

### 场景 33：映射转换

**推荐使用：`IfMap` / `IfMapElse`**

```go
// ✅ 条件映射（有默认值）
result := mathx.IfMap(hasData, rawData,
    func(d Data) string { return d.Format() },
    "无数据",
)

// ✅ 双向映射（完整版本）
output := mathx.IfMapElse(isJSON, data,
    func(d Data) string { return d.ToJSON() },
    func(d Data) string { return d.ToXML() },
)

// ✅ 省略 false 分支（返回零值）
output := mathx.IfMapElse(needFormat, data,
    func(d Data) string { return d.Format() },
)
// 当 needFormat=false 时，返回 ""（string 零值）
```

**可变参数说明：**

- `IfMapElse` 的 `falseMapper` 参数可选
- 省略时，条件为 `false` 时返回类型零值
- 适合只关心 true 分支转换的场景

**适用场景：**

- 数据转换
- 格式化输出
- 条件性的类型转换

---

### 场景 34：管道式处理

**推荐使用：`IfPipeline`**

```go
// ✅ 管道链
result := mathx.IfPipeline(shouldProcess, "hello", []func(string) string{
    strings.ToUpper,
    func(s string) string { return ">>> " + s },
    func(s string) string { return s + "!" },
}, "default")
// 结果: ">>> HELLO!"
```

**适用场景：**

- 数据处理流水线
- 多步转换

---

### 场景 35：带缓存的计算

**推荐使用：`IfMemoized`**

```go
// ✅ 缓存计算结果
cache := make(map[string]int)
result := mathx.IfMemoized(needCompute, "key1", cache,
    func() int { return expensiveCalculation() },
    0,
)
```

**适用场景：**

- 计算缓存
- 避免重复计算

---

### 场景 36：多值匹配

**推荐使用：`IfMulti`**

```go
// ✅ 多值匹配
isSpecial := mathx.IfMulti(code,
    []int{200, 201, 204},
    true,
    false,
)
```

**适用场景：**

- 状态码匹配
- 多值比较

---

### 场景 37：格式化字符串选择

**推荐使用：`IfStrFmt`**

```go
// ✅ 条件格式化
format, args := mathx.IfStrFmt(err != nil,
    "失败: %v", []any{err},
    "成功: %s", []any{result},
)
logger.Info(format, args...)
```

**适用场景：**

- 日志格式化
- 动态消息生成

---

### 场景 38：计数比较

**推荐使用：`IfCountGt`**

```go
// ✅ 计数阈值判断
message := mathx.IfCountGt(userCount, 1000, "热门", "普通")
```

**适用场景：**

- 数值阈值判断
- 统计展示

---

### 场景 39：异步执行

**推荐使用：`IfDoAsync` / `IfDoAsyncWithTimeout`**

```go
// ✅ 异步执行（提供默认值）
ch := mathx.IfDoAsync(needFetch,
    func() Data { return fetchData() },
    defaultData,
)
result := <-ch

// ✅ 异步执行（不提供默认值，返回零值）
ch := mathx.IfDoAsync(needFetch,
    func() Data { return fetchData() },
)
result := <-ch // 条件为 false 时返回 Data 类型零值

// ✅ 带超时（提供默认值）
ch := mathx.IfDoAsyncWithTimeout(needFetch,
    func() Data { return fetchData() },
    5000,       // 超时时间（毫秒）
    defaultData, // 默认值（可选）
)
result := <-ch

// ✅ 带超时（不提供默认值）
ch := mathx.IfDoAsyncWithTimeout(needFetch,
    func() Data { return fetchData() },
    5000, // 超时时间（毫秒）
)
result := <-ch // 条件为 false 或超时时返回零值
```

**可变参数说明：**

- `IfDoAsync` 的 `defaultVal` 参数可选
- `IfDoAsyncWithTimeout` 的 `defaultVal` 参数可选
- 参数顺序：`condition, do, timeoutMs, [defaultVal]`
- 省略 `defaultVal` 时，条件为 `false` 时返回类型零值
- ⚠️ **注意**：超时时返回零值，而不是 `defaultVal`

**适用场景：**

- 异步任务
- 超时控制
- 非关键性数据获取（允许返回零值）

---

### 场景 40：JSON 序列化

**推荐使用：`MarshalJSONOrDefault`**

```go
// ✅ 安全 JSON 序列化
jsonStr := mathx.MarshalJSONOrDefault(data, "{}")
```

**适用场景：**

- JSON 字段确保非空
- MySQL JSON 列

---

## ⚠️ 已废弃函数清单

| 已废弃函数 | 替代函数 | 原因 |
|-----------|---------|------|
| `IfElseFn` | `IF` | 完全相同，无需额外包装 |
| `IfV` | `IfExecElse` | 功能重复，语义不清晰 |
| `IfLazy` | `IfDoAF` | 功能重复，IfDoAF 语义更清晰 |
| `IfDefault` | `IfElse` | 功能完全相同 |
| `IfInRange` | `IfBetween` | 实现不正确（仅检查边界） |
| `IfIn` | `IfContains` | 功能重复（已删除 IfIn，保留 IfContains） |

---

## 📖 最佳实践

### ✅ DO（推荐）

```go
// 1. 简单判断用 IF
status := mathx.IF(isOnline, "在线", "离线")

// 2. 惰性求值用 IfDoAF
data := mathx.IfDoAF(needLoad,
    func() Data { return loadFromDB() },
    func() Data { return getCached() },
)

// 3. 副作用用 IfExec/IfExecElse
mathx.IfExec(debug, func() { log.Debug("调试信息") })

// 4. 链式构建提高可读性
err := mathx.IFChainFor[error]().
    When(invalid).ThenReturn(ErrInvalid).
    When(expired).ThenReturn(ErrExpired).
    ExecuteOr(nil)
```

### ❌ DON'T（不推荐）

```go
// ❌ 不要滥用，简单逻辑用传统 if
// 过度使用会降低可读性
result := mathx.IF(mathx.IF(a > b, true, false), 1, 0)

// ❌ 不要在性能敏感场景使用函数形式
// 函数调用有开销
for i := 0; i < 1000000; i++ {
    _ = mathx.IfDoAF(condition, expensiveFn, defaultFn) // 避免
}

// ❌ 不要忽略错误
// 错误处理还是要明确
_, _ = mathx.IfDoWithError(true, riskyOp, defaultVal) // 不好
```

---

## 🎯 选择决策树

```
需要返回值？
├─ 是
│  ├─ 简单条件 → IF
│  ├─ 需要延迟执行 → IfDo / IfDoAF
│  ├─ 多条件 → IfElse / IfChain
│  ├─ 带错误 → IfDoWithError / ReturnIfErr
│  └─ 链式调用 → WhenValue / IFChainFor
│
└─ 否（副作用）
   ├─ 单分支 → IfExec
   ├─ 双分支 → IfExecElse
   ├─ 带回调 → IfCall
   └─ 链式调用 → When().Then().Else()
```

---

## 📞 联系方式

- 作者：kamalyes
- 邮箱：<501893067@qq.com>
- 版本：v2.0 (2025-12-15 重构版)

---

**Happy Coding! 🚀**
