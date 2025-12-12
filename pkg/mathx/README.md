# Go 三元运算符库 (mathx/ternary)

一个基于 Go 泛型的强大三元运算符库，提供了丰富的条件判断和值选择功能，支持同步、异步、错误处理等多种场景。

## 特性

- 🚀 **泛型支持** - 完全基于 Go 1.18+ 泛型，类型安全
- 🔄 **同步/异步** - 支持同步和异步执行模式
- 🛡️ **错误处理** - 内置错误处理机制
- 🔗 **链式调用** - 支持优雅的链式 API
- 📦 **零依赖** - 纯标准库实现
- 🎯 **高性能** - 优化的执行路径

## 安装

```bash
go get github.com/kamalyes/go-toolbox
```

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/mathx"

// 基础三元运算
result := mathx.IF(score >= 60, "及格", "不及格")

// 空值检查
name := mathx.IfNotEmpty(user.Name, "匿名用户")

// 安全访问
value := mathx.IfSafeIndex(slice, index, "默认值")
```

## API 文档

### 基础三元运算

#### `IF[T any](condition bool, trueVal, falseVal T) T`

基础的三元运算符，类似于 `condition ? trueVal : falseVal`。

```go
age := 20
status := mathx.IF(age >= 18, "成年人", "未成年人")
// 结果: "成年人"
```

#### `IfNotNil[T any](val *T, defaultVal T) T`

空指针检查，如果指针不为 nil 则返回指针值，否则返回默认值。

```go
var ptr *int = &[]int{42}[0]
result := mathx.IfNotNil(ptr, 0)
// 结果: 42

var nilPtr *int
result = mathx.IfNotNil(nilPtr, 100)
// 结果: 100
```

#### `IfNotEmpty(str string, defaultVal string) string`

字符串空值检查。

```go
username := mathx.IfNotEmpty("", "guest")
// 结果: "guest"

username = mathx.IfNotEmpty("john", "guest")
// 结果: "john"
```

#### `IfNotZero[T comparable](val T, defaultVal T) T`

零值检查，支持任意可比较类型。

```go
count := mathx.IfNotZero(0, 1)
// 结果: 1

count = mathx.IfNotZero(5, 1)
// 结果: 5
```

### 集合操作

#### `IfContains[T comparable](slice []T, target T, trueVal, falseVal T) T`

检查切片是否包含指定元素。

```go
fruits := []string{"apple", "banana", "orange"}
result := mathx.IfContains(fruits, "banana", "找到了", "没找到")
// 结果: "找到了"
```

#### `IfSafeIndex[T any](slice []T, index int, defaultVal T) T`

安全的切片索引访问。

```go
arr := []string{"a", "b", "c"}
result := mathx.IfSafeIndex(arr, 5, "默认值")
// 结果: "默认值"

result = mathx.IfSafeIndex(arr, 1, "默认值")
// 结果: "b"
```

#### `IfSafeKey[K comparable, V any](m map[K]V, key K, defaultVal V) V`

安全的 map 键访问。

```go
config := map[string]string{
    "host": "localhost",
    "port": "8080",
}

host := mathx.IfSafeKey(config, "host", "127.0.0.1")
// 结果: "localhost"

timeout := mathx.IfSafeKey(config, "timeout", "30s")
// 结果: "30s"
```

### 条件组合

#### `IfAny[T any](conditions []bool, trueVal, falseVal T) T`

任意条件满足时返回真值。

```go
conditions := []bool{false, true, false}
result := mathx.IfAny(conditions, "有条件满足", "无条件满足")
// 结果: "有条件满足"
```

#### `IfAll[T any](conditions []bool, trueVal, falseVal T) T`

所有条件都满足时返回真值。

```go
conditions := []bool{true, true, true}
result := mathx.IfAll(conditions, "全部满足", "部分不满足")
// 结果: "全部满足"
```

#### `IfCount[T any](conditions []bool, threshold int, trueVal, falseVal T) T`

满足条件的数量达到阈值时返回真值。

```go
conditions := []bool{true, false, true, true}
result := mathx.IfCount(conditions, 2, "达到阈值", "未达到阈值")
// 结果: "达到阈值"
```

### 函数式操作

#### `IfMap[T, R any](condition bool, val T, mapper func(T) R, defaultVal R) R`

条件映射转换。

```go
text := "hello"
result := mathx.IfMap(true, text, strings.ToUpper, "默认值")
// 结果: "HELLO"
```

#### `IfMapElse[T, R any](condition bool, val T, trueMapper, falseMapper func(T) R) R`

双向映射转换。

```go
text := "Hello"
result := mathx.IfMapElse(true, text, strings.ToUpper, strings.ToLower)
// 结果: "HELLO"
```

#### `IfFilter[T any](useFilter bool, slice []T, predicate func(T) bool) []T`

条件过滤。

```go
numbers := []int{1, 2, 3, 4, 5}
evens := mathx.IfFilter(true, numbers, func(n int) bool { return n%2 == 0 })
// 结果: [2, 4]
```

#### `IfValidate[T, R any](val T, validator func(T) bool, validVal, invalidVal R) R`

验证函数。

```go
email := "user@example.com"
isValid := func(s string) bool { return strings.Contains(s, "@") }
result := mathx.IfValidate(email, isValid, "有效邮箱", "无效邮箱")
// 结果: "有效邮箱"
```

### 类型转换

#### `IfCast[R any](val any, defaultVal R) R`

安全的类型转换。

```go
var value interface{} = "hello"
result := mathx.IfCast[string](value, "默认值")
// 结果: "hello"

result = mathx.IfCast[int](value, 0)
// 结果: 0
```

#### `IfBetween[T int | int64 | float32 | float64](val, min, max T, trueVal, falseVal T) T`

数值区间检查。

```go
score := 85
grade := mathx.IfBetween(score, 80, 100, 90, 60)
// 结果: 90 (因为 85 在 80-100 区间内)
```

### 高级功能

#### `IfSwitch[K comparable, V any](key K, cases map[K]V, defaultVal V) V`

开关式选择。

```go
status := "success"
cases := map[string]string{
    "success": "操作成功",
    "error":   "操作失败",
    "pending": "操作进行中",
}
message := mathx.IfSwitch(status, cases, "未知状态")
// 结果: "操作成功"
```

#### `IfTryParse[T, R any](input T, parser func(T) (R, error), defaultVal R) R`

尝试解析操作。

```go
parser := func(s string) (int, error) { return strconv.Atoi(s) }
result := mathx.IfTryParse("123", parser, 0)
// 结果: 123

result = mathx.IfTryParse("abc", parser, 0)
// 结果: 0
```

### 异步操作

#### `IfDoAsync[T any](condition bool, do DoFunc[T], defaultVal T) <-chan T`

异步执行函数。

```go
ch := mathx.IfDoAsync(true, func() string {
    time.Sleep(100 * time.Millisecond)
    return "异步结果"
}, "默认值")

result := <-ch
// 结果: "异步结果"
```

#### `IfDoAsyncWithTimeout[T any](condition bool, do DoFunc[T], defaultVal T, timeoutMs int) <-chan T`

带超时的异步执行。

```go
ch := mathx.IfDoAsyncWithTimeout(true, func() string {
    time.Sleep(200 * time.Millisecond)
    return "结果"
}, "默认值", 100) // 100ms 超时

result := <-ch
// 结果: 零值 (超时)
```

### 错误处理

#### `IfDoWithError[T any](condition bool, do DoFuncWithError[T], defaultVal T) (T, error)`

带错误处理的函数执行。

```go
result, err := mathx.IfDoWithError(true, func() (int, error) {
    return strconv.Atoi("123")
}, 0)
// result: 123, err: nil
```

#### `ReturnIfErr[T any](val T, err error) (T, error)`

错误检查简化。

```go
value, err := someFunction()
return mathx.ReturnIfErr(value, err)
```

### 链式调用

#### 执行链

```go
mathx.When(err != nil).
    Then(func() { log.Error("操作失败") }).
    Else(func() { log.Info("操作成功") }).
    Do()
```

#### 值链

```go
result := mathx.WhenValue(score >= 90).
    ThenReturn("优秀").
    ElseReturn("良好").
    Get()
```

### 实用功能

#### `IfPipeline[T any](condition bool, input T, funcs []func(T) T, defaultVal T) T`

管道式处理。

```go
funcs := []func(string) string{
    strings.ToUpper,
    func(s string) string { return s + "!" },
    func(s string) string { return ">>> " + s },
}

result := mathx.IfPipeline(true, "hello", funcs, "默认值")
// 结果: ">>> HELLO!"
```

#### `IfLazy[T any](condition bool, trueFn, falseFn func() T) T`

惰性求值。

```go
result := mathx.IfLazy(condition, 
    func() string { return expensiveComputation() },
    func() string { return "快速默认值" })
```

#### `IfMemoized[T any](condition bool, key string, cache map[string]T, computeFn func() T, defaultVal T) T`

带缓存的计算。

```go
cache := make(map[string]string)
result := mathx.IfMemoized(true, "key1", cache, 
    func() string { return expensiveComputation() }, 
    "默认值")
```

## 使用场景

### 1. 配置处理

```go
config := map[string]string{
    "env": "production",
}

dbHost := mathx.IfSwitch(
    mathx.IfSafeKey(config, "env", "development"),
    map[string]string{
        "production":  "prod-db.example.com",
        "staging":     "stage-db.example.com",
        "development": "localhost",
    },
    "localhost",
)
```

### 2. 用户输入验证

```go
username := mathx.IfNotEmpty(
    strings.TrimSpace(input.Username), 
    "anonymous",
)

age := mathx.IfBetween(input.Age, 0, 150, input.Age, 0)
```

### 3. API 响应处理

```go
status := mathx.IfValidate(user, 
    func(u User) bool { return u.IsActive }, 
    "active", 
    "inactive",
)

response := mathx.IfMap(len(results) > 0, results,
    func(r []Item) ApiResponse { 
        return ApiResponse{Data: r, Success: true} 
    },
    ApiResponse{Error: "No data found"},
)
```

### 4. 错误处理

```go
mathx.When(err != nil).
    Then(func() {
        log.Error("操作失败", "error", err)
        metrics.IncrementCounter("errors")
    }).
    Else(func() {
        log.Info("操作成功")
        metrics.IncrementCounter("success")
    }).
    Do()
```

## 性能考虑

- 所有函数都经过性能优化
- 惰性求值避免不必要的计算
- 内联友好的实现
- 零内存分配（除异步操作外）

## 最佳实践

1. **优先使用简单的 IF 函数**：对于简单的条件判断，使用基础的 `IF` 函数。

2. **合理使用异步操作**：只在真正需要并发的场景使用异步版本。

3. **错误处理**：对于可能出错的操作，优先使用带错误处理的版本。

4. **链式调用**：对于复杂的条件逻辑，使用链式调用提高可读性。

5. **缓存计算**：对于昂贵的计算操作，考虑使用 `IfMemoized`。

## 许可证

MIT License - 详见 [LICENSE](../../LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

### v1.0.0

- 初始版本发布
- 基础三元运算符功能
- 泛型支持
- 异步操作
- 错误处理
- 链式调用API
