# SafeAccess 泛型和增强功能使用指南

## 概述

`SafeAccess` 现在支持强大的泛型操作和类型转换，集成了 `convert` 包的转换能力，提供了类似 JavaScript 可选链的安全访问体验。

## 核心特性

### 1. 泛型数值转换

使用泛型函数进行类型安全的数值转换：

```go
// 基础用法
s := safe.Safe(42)

// 转换为不同的数值类型
intVal := safe.As[int](s)           // 42
int64Val := safe.As[int64](s)       // int64(42)
uintVal := safe.As[uint](s)         // uint(42)

// 从字符串转换
s = safe.Safe("123")
result := safe.As[int32](s)         // int32(123)

// 使用默认值
s = safe.Safe("invalid")
result = safe.As[int](s, 999)       // 999（转换失败时使用默认值）
```

### 2. 泛型浮点数转换

支持多种取整模式的浮点数转换：

```go
s := safe.Safe("3.14159")

// 不取整
f1 := safe.AsFloat[float64](s, convert.RoundNone)     // 3.14159

// 四舍五入
f2 := safe.AsFloat[float64](s, convert.RoundNearest)  // 3.0

// 向上取整
f3 := safe.AsFloat[float64](s, convert.RoundUp)       // 4.0

// 向下取整
f4 := safe.AsFloat[float64](s, convert.RoundDown)     // 3.0
```

### 3. 智能类型转换

#### AsString - 智能字符串转换
```go
safe.Safe(123).AsString()           // "123"
safe.Safe(3.14).AsString()          // "3.14"
safe.Safe(true).AsString()          // "true"

// 时间格式化
t := time.Now()
safe.Safe(t).AsString("2006-01-02") // "2025-11-21"
```

#### AsBool - 智能布尔转换
```go
safe.Safe(true).AsBool()            // true
safe.Safe(1).AsBool()               // true
safe.Safe(0).AsBool()               // false
safe.Safe("true").AsBool()          // true
safe.Safe("1").AsBool()             // true
```

#### AsJSON - JSON 转换
```go
data := map[string]interface{}{
    "name": "test",
    "age":  30,
}
s := safe.Safe(data)

// 紧凑格式
json, _ := s.AsJSON(false)
// {"name":"test","age":30}

// 缩进格式
json, _ := s.AsJSON(true)
// {
//   "name": "test",
//   "age": 30
// }
```

### 4. 泛型切片转换

#### AsSlice - 数值切片转换
```go
// 字符串切片转数值切片
s := safe.Safe([]string{"10", "20", "30"})
nums, _ := safe.AsSlice[int](s)
// []int{10, 20, 30}

// interface{} 切片转换
s = safe.Safe([]interface{}{1, 2, 3})
result, _ := safe.AsSlice[int64](s)
// []int64{1, 2, 3}
```

#### AsFloatSlice - 浮点数切片转换
```go
s := safe.Safe([]string{"1.5", "2.5", "3.5"})
floats, _ := safe.AsFloatSlice[float64](s, convert.RoundNone)
// []float64{1.5, 2.5, 3.5}

// 带取整
s = safe.Safe([]string{"1.4", "2.6", "3.5"})
rounded, _ := safe.AsFloatSlice[float64](s, convert.RoundNearest)
// []float64{1.0, 3.0, 4.0}
```

#### AsStringSlice - 字符串切片转换
```go
s := safe.Safe([]interface{}{1, 2.5, true, "test"})
strs := s.AsStringSlice()
// []string{"1", "2.5", "true", "test"}
```

### 5. 链式操作增强

#### Map - 泛型映射
```go
s := safe.Safe(10)
result := safe.Map[int, string](s, func(v int) string {
    return fmt.Sprintf("值是: %d", v * 2)
})
// SafeAccess("值是: 20")
```

#### FlatMap - 扁平化映射
```go
s := safe.Safe(5)
result := s.FlatMap(func(v interface{}) *safe.SafeAccess {
    if num, ok := v.(int); ok {
        return safe.Safe(num * 3)
    }
    return &safe.SafeAccess{valid: false}
})
// SafeAccess(15)
```

#### OrDefault - 默认值
```go
s := safe.Safe(42)
result := safe.OrDefault[int](s, 999)  // 42

s = &safe.SafeAccess{valid: false}
result = safe.OrDefault[int](s, 999)   // 999
```

#### Must - 强制获取值
```go
s := safe.Safe(100)
result := safe.Must[int](s)  // 100

// 无效值会 panic
s = &safe.SafeAccess{valid: false}
result = safe.Must[int](s)  // panic!
```

### 6. 条件操作

#### When - 条件执行
```go
s := safe.Safe(10)
result := s.When(
    func(v interface{}) bool {
        return v.(int) > 5
    },
    func(v interface{}) interface{} {
        return v.(int) * 2
    },
)
// SafeAccess(20)
```

#### Unless - 条件排除
```go
s := safe.Safe(3)
result := s.Unless(
    func(v interface{}) bool {
        return v.(int) > 5
    },
    func(v interface{}) interface{} {
        return v.(int) * 2
    },
)
// SafeAccess(6)
```

### 7. 类型检查

```go
s := safe.Safe(42)

s.IsNumber()              // true
s.IsString()              // false
safe.IsType[int](s)       // true
safe.IsType[string](s)    // false

s = safe.Safe([]int{1, 2, 3})
s.IsSlice()               // true

s = safe.Safe(map[string]interface{}{})
s.IsMap()                 // true
```

### 8. 集合操作

#### 长度获取
```go
safe.Safe("hello").Len()                    // 5
safe.Safe([]int{1, 2, 3}).Len()             // 3
safe.Safe(map[string]int{"a": 1}).Len()     // 1
```

#### Keys 和 Values
```go
m := map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
}
s := safe.Safe(m)

keys := s.Keys()      // []string{"key1", "key2"}
values := s.Values()  // []interface{}{"value1", "value2"}
```

#### Contains - 包含检查
```go
// Map 检查
m := map[string]interface{}{"key1": "value1"}
safe.Safe(m).Contains("key1")  // true

// 切片检查
slice := []int{1, 2, 3, 4, 5}
safe.Safe(slice).Contains(3)   // true
```

### 9. 空值检查

```go
safe.Safe("").IsEmpty()              // true
safe.Safe("test").IsEmpty()          // false
safe.Safe([]string{}).IsEmpty()      // true
safe.Safe("test").IsNonEmpty()       // true
```

## 完整使用示例

### 配置解析场景
```go
config := map[string]interface{}{
    "server": map[string]interface{}{
        "port":    "8080",
        "timeout": "30",
        "enabled": "true",
    },
    "database": map[string]interface{}{
        "hosts": []string{"host1", "host2"},
        "port":  "5432",
    },
}

s := safe.Safe(config)

// 获取服务器配置
serverPort := safe.As[int](s.Field("server").Field("port"))
// 8080

timeout := safe.As[int](s.Field("server").Field("timeout"))
// 30

enabled := s.Field("server").Field("enabled").AsBool()
// true

// 获取数据库主机列表
dbHosts := s.Field("database").Field("hosts").AsStringSlice()
// []string{"host1", "host2"}
```

### 链式转换场景
```go
data := map[string]interface{}{
    "user": map[string]interface{}{
        "age":   "25",
        "score": []string{"85", "90", "95"},
    },
}

s := safe.Safe(data)

// 链式访问并转换
age := safe.As[int](s.Field("user").Field("age"))
// 25

// 切片转换
scores, _ := safe.AsSlice[int](s.Field("user").Field("score"))
// []int{85, 90, 95}
```

### 条件处理场景
```go
s := safe.Safe(100)

result := s.
    When(func(v interface{}) bool {
        return v.(int) > 50
    }, func(v interface{}) interface{} {
        return v.(int) / 2
    }).
    Map(func(v interface{}) interface{} {
        return v.(int) + 10
    })

finalValue := result.Int()
// 60 (100 > 50 ? 100/2 : 100) + 10 = 60
```

### 错误处理场景
```go
data := map[string]interface{}{
    "valid":   "123",
    "invalid": "not-a-number",
    "missing": nil,
}

s := safe.Safe(data)

// 有效数据
validNum := safe.As[int](s.Field("valid"))
// 123

// 无效数据使用默认值
invalidNum := safe.As[int](s.Field("invalid"), 999)
// 999

// 不存在的字段
missingNum := safe.As[int](s.Field("missing"), 777)
// 777
```

## 最佳实践

1. **使用泛型函数进行类型转换**
   - 优先使用 `As[T]` 和 `AsFloat[T]` 而不是具体类型方法
   - 提供合理的默认值处理转换失败

2. **链式调用时注意空值处理**
   - 使用 `IsValid()` 检查中间结果
   - 利用 `OrElse()` 提供备用值

3. **切片转换时处理错误**
   - 切片转换可能失败，务必检查错误
   - 使用合适的取整模式处理浮点数

4. **合理使用条件操作**
   - `When` 和 `Unless` 适合简单的条件转换
   - 复杂逻辑建议先获取值再处理

5. **类型检查优先于类型断言**
   - 使用 `IsType[T]` 进行安全的类型检查
   - 使用 `Must[T]` 时确保值一定存在

## 性能考虑

- 泛型函数在编译时进行类型特化，性能接近直接类型转换
- 链式调用会创建中间对象，避免过深的链式调用
- 大量数据转换时考虑批量处理而非逐个转换

## 总结

通过泛型和增强功能，`SafeAccess` 提供了：

✅ 类型安全的泛型转换  
✅ 强大的类型推断和转换  
✅ 灵活的链式操作  
✅ 丰富的条件处理  
✅ 完善的集合操作  
✅ 优雅的错误处理  

享受类型安全和便捷开发的双重体验！
