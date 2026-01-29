/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 10:59:15
 * @FilePath: \go-toolbox\pkg\contextx\getters_test.go
 * @Description: Context 类型安全的 Getter 方法测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package contextx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGetTString 测试泛型获取字符串
func TestGetTString(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试字符串值
	ctx.WithValue(TestKey1, TestValue1)
	result := Get[string](ctx, TestKey1)
	assert.Equal(t, TestValue1, result)

	// 测试空值
	result = Get[string](ctx, TestNonExistentKey)
	assert.Equal(t, "", result)
}

// TestGetTInt 测试泛型获取整数
func TestGetTInt(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试整数值
	ctx.WithValue(TestKey1, TestInt)
	result := Get[int](ctx, TestKey1)
	assert.Equal(t, TestInt, result)

	// 测试字符串转整数
	ctx.WithValue(TestKey2, TestIntStr100)
	result = Get[int](ctx, TestKey2)
	assert.Equal(t, TestInt100, result)

	// 测试空值
	result = Get[int](ctx, TestNonExistentKey)
	assert.Equal(t, 0, result)
}

// TestGetTInt64 测试泛型获取 int64
func TestGetTInt64(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试 int64 值
	ctx.WithValue(TestKey1, TestInt64)
	result := Get[int64](ctx, TestKey1)
	assert.Equal(t, TestInt64, result)

	// 测试字符串转 int64
	ctx.WithValue(TestKey2, TestIntStr999)
	result = Get[int64](ctx, TestKey2)
	assert.Equal(t, int64(TestInt999), result)

	// 测试空值
	result = Get[int64](ctx, TestNonExistentKey)
	assert.Equal(t, int64(0), result)
}

// TestGetTBool 测试泛型获取布尔值
func TestGetTBool(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试布尔值
	ctx.WithValue(TestKey1, true)
	result := Get[bool](ctx, TestKey1)
	assert.True(t, result)

	ctx.WithValue(TestKey2, false)
	result = Get[bool](ctx, TestKey2)
	assert.False(t, result)

	// 测试字符串转布尔
	ctx.WithValue(TestKey3, "true")
	result = Get[bool](ctx, TestKey3)
	assert.True(t, result)

	// 测试空值
	result = Get[bool](ctx, TestNonExistentKey)
	assert.False(t, result)
}

// TestGetTFloat64 测试泛型获取浮点数
func TestGetTFloat64(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试浮点数值
	ctx.WithValue(TestKey1, TestFloat64314)
	result := Get[float64](ctx, TestKey1)
	assert.Equal(t, TestFloat64314, result)

	// 测试字符串转浮点数
	ctx.WithValue(TestKey2, TestFloatStr)
	result = Get[float64](ctx, TestKey2)
	assert.Equal(t, TestFloat271, result)

	// 测试空值
	result = Get[float64](ctx, TestNonExistentKey)
	assert.Equal(t, 0.0, result)
}

// TestGetTStringSlice 测试泛型获取字符串切片
func TestGetTStringSlice(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试字符串切片
	slice := []string{"a", "b", "c"}
	ctx.WithValue(TestKey1, slice)
	result := Get[[]string](ctx, TestKey1)
	assert.Equal(t, slice, result)

	// 测试空值
	result = Get[[]string](ctx, TestNonExistentKey)
	assert.Nil(t, result)
}

// TestGetTDuration 测试泛型获取时间间隔
func TestGetTDuration(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试 Duration 值
	duration := TestTimeout5s
	ctx.WithValue(TestKey1, duration)
	result := Get[time.Duration](ctx, TestKey1)
	assert.Equal(t, duration, result)

	// 测试字符串转 Duration
	ctx.WithValue(TestKey2, "10s")
	result = Get[time.Duration](ctx, TestKey2)
	assert.Equal(t, TestTimeout10s, result)

	// 测试 int64 转 Duration
	ctx.WithValue(TestKey3, int64(1000000000))
	result = Get[time.Duration](ctx, TestKey3)
	assert.Equal(t, time.Second, result)

	// 测试 int 转 Duration
	ctx.WithValue(TestKey4, 2000000000)
	result = Get[time.Duration](ctx, TestKey4)
	assert.Equal(t, TestTimeout2s, result)

	// 测试空值
	result = Get[time.Duration](ctx, TestNonExistentKey)
	assert.Equal(t, time.Duration(0), result)
}

// TestGetTTime 测试泛型获取时间值
func TestGetTTime(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试 Time 值
	now := time.Now()
	ctx.WithValue(TestKey1, now)
	result := Get[time.Time](ctx, TestKey1)
	assert.Equal(t, now, result)

	// 测试 RFC3339 字符串转 Time
	timeStr := TestTimeRFC3339
	ctx.WithValue(TestKey2, timeStr)
	expected, _ := time.Parse(time.RFC3339, timeStr)
	result = Get[time.Time](ctx, TestKey2)
	assert.Equal(t, expected, result)

	// 测试 Unix 时间戳转 Time
	timestamp := TestTimestamp
	ctx.WithValue(TestKey3, timestamp)
	result = Get[time.Time](ctx, TestKey3)
	assert.Equal(t, time.Unix(timestamp, 0), result)

	// 测试空值
	result = Get[time.Time](ctx, TestNonExistentKey)
	assert.Equal(t, time.Time{}, result)
}

// TestGetTCustomStruct 测试泛型获取自定义结构体
func TestGetTCustomStruct(t *testing.T) {
	type MyStruct struct {
		Name string
		Age  int
	}

	ctx := NewContext().WithParent(context.Background())

	// 测试自定义结构体
	data := MyStruct{Name: "test", Age: 25}
	ctx.WithValue(TestKey1, data)
	result := Get[MyStruct](ctx, TestKey1)
	assert.Equal(t, data, result)

	// 测试空值
	result = Get[MyStruct](ctx, TestNonExistentKey)
	assert.Equal(t, MyStruct{}, result)
}

// TestGetTPointer 测试泛型获取指针类型
func TestGetTPointer(t *testing.T) {
	type MyStruct struct {
		Value int
	}

	ctx := NewContext().WithParent(context.Background())

	// 测试指针类型
	data := &MyStruct{Value: TestInt99}
	ctx.WithValue(TestKey1, data)
	result := Get[*MyStruct](ctx, TestKey1)
	assert.Equal(t, data, result)
	assert.Equal(t, TestInt99, result.Value)

	// 测试空值
	result = Get[*MyStruct](ctx, TestNonExistentKey)
	assert.Nil(t, result)
}

// TestGetTIntSlice 测试泛型获取整数切片
func TestGetTIntSlice(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试整数切片
	slice := []int{1, 2, 3, 4, 5}
	ctx.WithValue(TestKey1, slice)
	result := Get[[]int](ctx, TestKey1)
	assert.Equal(t, slice, result)

	// 测试空值
	result = Get[[]int](ctx, TestNonExistentKey)
	assert.Nil(t, result)
}

// TestGetTMap 测试泛型获取 map
func TestGetTMap(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试 map
	data := map[string]int{"a": 1, "b": 2}
	ctx.WithValue(TestKey1, data)
	result := Get[map[string]int](ctx, TestKey1)
	assert.Equal(t, data, result)

	// 测试空值
	result = Get[map[string]int](ctx, TestNonExistentKey)
	assert.Nil(t, result)
}

func TestGetString(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试字符串值
	ctx.WithValue(TestKey1, TestValue1)
	assert.Equal(t, TestValue1, Get[string](ctx, TestKey1))

	// 测试空值
	assert.Equal(t, "", Get[string](ctx, TestNonExistentKey))

	// 测试非字符串类型
	ctx.WithValue(TestKey2, TestInt123)
	assert.Equal(t, "", Get[string](ctx, TestKey2))
}

func TestGetInt(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试整数值
	ctx.WithValue(TestKey1, TestInt)
	assert.Equal(t, TestInt, Get[int](ctx, TestKey1))

	// 测试空值
	assert.Equal(t, 0, Get[int](ctx, TestNonExistentKey))

	// 测试字符串转整数
	ctx.WithValue(TestKey2, TestIntStr100)
	assert.Equal(t, TestInt100, Get[int](ctx, TestKey2))
}

func TestGetInt64(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试 int64 值
	ctx.WithValue(TestKey1, TestInt64)
	assert.Equal(t, TestInt64, ctx.GetInt64(TestKey1))

	// 测试空值
	assert.Equal(t, int64(0), ctx.GetInt64(TestNonExistentKey))

	// 测试字符串转 int64
	ctx.WithValue(TestKey2, TestIntStr999)
	assert.Equal(t, int64(TestInt999), ctx.GetInt64(TestKey2))
}

func TestGetBool(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试布尔值
	ctx.WithValue(TestKey1, true)
	assert.True(t, ctx.GetBool(TestKey1))

	ctx.WithValue(TestKey2, false)
	assert.False(t, ctx.GetBool(TestKey2))

	// 测试空值
	assert.False(t, ctx.GetBool(TestNonExistentKey))

	// 测试字符串转布尔
	ctx.WithValue(TestKey3, "true")
	assert.True(t, ctx.GetBool(TestKey3))
}

func TestGetFloat64(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试浮点数值
	ctx.WithValue(TestKey1, TestFloat64314)
	assert.Equal(t, TestFloat64314, Get[float64](ctx, TestKey1))

	// 测试空值
	assert.Equal(t, 0.0, Get[float64](ctx, TestNonExistentKey))

	// 测试字符串转浮点数
	ctx.WithValue(TestKey2, TestFloatStr)
	assert.Equal(t, TestFloat271, Get[float64](ctx, TestKey2))
}

func TestGetStringSlice(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试字符串切片
	slice := []string{"a", "b", "c"}
	ctx.WithValue(TestKey1, slice)
	assert.Equal(t, slice, ctx.GetStringSlice(TestKey1))

	// 测试空值
	assert.Nil(t, ctx.GetStringSlice(TestNonExistentKey))

	// 测试从 []interface{} 转换
	interfaceSlice := []interface{}{"x", "y", "z"}
	ctx.WithValue(TestKey2, interfaceSlice)
	result := ctx.GetStringSlice(TestKey2)
	assert.Equal(t, []string{"x", "y", "z"}, result)
}

func TestGetIntSlice(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试整数切片
	intSlice := []int{10, 20, 30}
	ctx.WithValue(TestKey1, intSlice)
	assert.Equal(t, intSlice, ctx.GetIntSlice(TestKey1))

	// 测试从 []interface{} 转换
	interfaceSlice := []interface{}{10, 20, 30}
	ctx.WithValue(TestKey2, interfaceSlice)
	result := ctx.GetIntSlice(TestKey2)
	assert.Equal(t, []int{10, 20, 30}, result)

	// 测试空值
	assert.Nil(t, ctx.GetIntSlice(TestNonExistentKey))
}

func TestSafeGetStringSlice(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	stringSlice := []string{"a", "b", "c"}
	ctx.WithValue(TestKey1, stringSlice)
	assert.Equal(t, stringSlice, ctx.SafeGetStringSlice(TestKey1))

	// 测试从 []interface{} 转换
	interfaceSlice := []interface{}{"test1", "test2"}
	ctx.WithValue(TestKey2, interfaceSlice)
	result := ctx.SafeGetStringSlice(TestKey2)
	assert.Equal(t, []string{"test1", "test2"}, result)
}

func TestGetMap(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试 map[string]interface{}
	testMap := map[string]interface{}{"key1": "value1", "key2": 123}
	ctx.WithValue(TestKey1, testMap)
	assert.Equal(t, testMap, ctx.GetMap(TestKey1))

	// 测试从 map[interface{}]interface{} 转换
	interfaceMap := map[interface{}]interface{}{"key1": "value1", "key2": 456}
	ctx.WithValue(TestKey2, interfaceMap)
	result := ctx.GetMap(TestKey2)
	assert.Equal(t, map[string]interface{}{"key1": "value1", "key2": 456}, result)

	// 测试空值
	assert.Nil(t, ctx.GetMap(TestNonExistentKey))
}

func TestGetInt8(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, int8(8))
	assert.Equal(t, int8(8), ctx.GetInt8(TestKey1))

	ctx.WithValue(TestKey2, "8")
	assert.Equal(t, int8(8), ctx.GetInt8(TestKey2))

	assert.Equal(t, int8(0), ctx.GetInt8(TestNonExistentKey))
}

func TestGetInt16(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, int16(16))
	assert.Equal(t, int16(16), ctx.GetInt16(TestKey1))

	ctx.WithValue(TestKey2, "16")
	assert.Equal(t, int16(16), ctx.GetInt16(TestKey2))

	assert.Equal(t, int16(0), ctx.GetInt16(TestNonExistentKey))
}

func TestGetInt32(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, int32(32))
	assert.Equal(t, int32(32), ctx.GetInt32(TestKey1))

	ctx.WithValue(TestKey2, "32")
	assert.Equal(t, int32(32), ctx.GetInt32(TestKey2))

	assert.Equal(t, int32(0), ctx.GetInt32(TestNonExistentKey))
}

func TestGetRune(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, 'A')
	assert.Equal(t, 'A', ctx.GetRune(TestKey1))

	ctx.WithValue(TestKey2, int32(65)) // 'A' 的 ASCII 码
	assert.Equal(t, rune(65), ctx.GetRune(TestKey2))

	assert.Equal(t, rune(0), ctx.GetRune(TestNonExistentKey))
}

func TestGetUint(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, uint(100))
	assert.Equal(t, uint(100), ctx.GetUint(TestKey1))

	ctx.WithValue(TestKey2, "100")
	assert.Equal(t, uint(100), ctx.GetUint(TestKey2))

	assert.Equal(t, uint(0), ctx.GetUint(TestNonExistentKey))
}

func TestGetUint8(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, uint8(200))
	assert.Equal(t, uint8(200), ctx.GetUint8(TestKey1))

	ctx.WithValue(TestKey2, "200")
	assert.Equal(t, uint8(200), ctx.GetUint8(TestKey2))

	assert.Equal(t, uint8(0), ctx.GetUint8(TestNonExistentKey))
}

func TestGetUint16(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, uint16(300))
	assert.Equal(t, uint16(300), ctx.GetUint16(TestKey1))

	ctx.WithValue(TestKey2, "300")
	assert.Equal(t, uint16(300), ctx.GetUint16(TestKey2))

	assert.Equal(t, uint16(0), ctx.GetUint16(TestNonExistentKey))
}

func TestGetUint32(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, uint32(400))
	assert.Equal(t, uint32(400), ctx.GetUint32(TestKey1))

	ctx.WithValue(TestKey2, "400")
	assert.Equal(t, uint32(400), ctx.GetUint32(TestKey2))

	assert.Equal(t, uint32(0), ctx.GetUint32(TestNonExistentKey))
}

func TestGetUint64(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, uint64(500))
	assert.Equal(t, uint64(500), ctx.GetUint64(TestKey1))

	ctx.WithValue(TestKey2, "500")
	assert.Equal(t, uint64(500), ctx.GetUint64(TestKey2))

	assert.Equal(t, uint64(0), ctx.GetUint64(TestNonExistentKey))
}

func TestGetFloat32(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	ctx.WithValue(TestKey1, float32(3.14))
	assert.InDelta(t, float32(3.14), ctx.GetFloat32(TestKey1), 0.001)

	ctx.WithValue(TestKey2, "3.14")
	assert.InDelta(t, float32(3.14), ctx.GetFloat32(TestKey2), 0.001)

	assert.Equal(t, float32(0), ctx.GetFloat32(TestNonExistentKey))
}

func TestGetDuration(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试 Duration 值
	duration := TestTimeout5s
	ctx.WithValue(TestKey1, duration)
	assert.Equal(t, duration, ctx.GetDuration(TestKey1))

	// 测试空值
	assert.Equal(t, time.Duration(0), ctx.GetDuration(TestNonExistentKey))

	// 测试字符串转 Duration
	ctx.WithValue(TestKey2, "10s")
	assert.Equal(t, TestTimeout10s, ctx.GetDuration(TestKey2))

	// 测试 int64 转 Duration
	ctx.WithValue(TestKey3, int64(1000000000))
	assert.Equal(t, time.Second, ctx.GetDuration(TestKey3))

	// 测试 int 转 Duration
	ctx.WithValue(TestKey4, 2000000000)
	assert.Equal(t, TestTimeout2s, ctx.GetDuration(TestKey4))
}

func TestGetTime(t *testing.T) {
	ctx := NewContext().WithParent(context.Background())

	// 测试 Time 值
	now := time.Now()
	ctx.WithValue(TestKey1, now)
	assert.Equal(t, now, ctx.GetTime(TestKey1))

	// 测试空值
	assert.Equal(t, time.Time{}, ctx.GetTime(TestNonExistentKey))

	// 测试 RFC3339 字符串转 Time
	timeStr := TestTimeRFC3339
	ctx.WithValue(TestKey2, timeStr)
	expected, _ := time.Parse(time.RFC3339, timeStr)
	assert.Equal(t, expected, ctx.GetTime(TestKey2))

	// 测试 Unix 时间戳转 Time
	timestamp := TestTimestamp
	ctx.WithValue(TestKey3, timestamp)
	assert.Equal(t, time.Unix(timestamp, 0), ctx.GetTime(TestKey3))
}

// TestGetValue 测试从标准 context.Context 获取值
func TestGetValue(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		result := GetValue[string](nil, "key")
		assert.Equal(t, "", result)
	})

	t.Run("string value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", "value")
		result := GetValue[string](ctx, "key")
		assert.Equal(t, "value", result)
	})

	t.Run("int value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "count", 123)
		result := GetValue[int](ctx, "count")
		assert.Equal(t, 123, result)
	})

	t.Run("bool value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "enabled", true)
		result := GetValue[bool](ctx, "enabled")
		assert.True(t, result)
	})

	t.Run("float64 value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "price", 99.99)
		result := GetValue[float64](ctx, "price")
		assert.Equal(t, 99.99, result)
	})

	t.Run("custom struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		user := User{Name: "Alice", Age: 30}
		ctx := context.WithValue(context.Background(), "user", user)
		result := GetValue[User](ctx, "user")
		assert.Equal(t, user, result)
	})

	t.Run("pointer type", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}
		config := &Config{Host: "localhost", Port: 8080}
		ctx := context.WithValue(context.Background(), "config", config)
		result := GetValue[*Config](ctx, "config")
		assert.Equal(t, config, result)
		assert.Equal(t, "localhost", result.Host)
		assert.Equal(t, 8080, result.Port)
	})

	t.Run("slice type", func(t *testing.T) {
		slice := []string{"a", "b", "c"}
		ctx := context.WithValue(context.Background(), "items", slice)
		result := GetValue[[]string](ctx, "items")
		assert.Equal(t, slice, result)
	})

	t.Run("map type", func(t *testing.T) {
		data := map[string]int{"x": 1, "y": 2}
		ctx := context.WithValue(context.Background(), "data", data)
		result := GetValue[map[string]int](ctx, "data")
		assert.Equal(t, data, result)
	})

	t.Run("time.Duration", func(t *testing.T) {
		duration := 5 * time.Second
		ctx := context.WithValue(context.Background(), "timeout", duration)
		result := GetValue[time.Duration](ctx, "timeout")
		assert.Equal(t, duration, result)
	})

	t.Run("time.Time", func(t *testing.T) {
		now := time.Now()
		ctx := context.WithValue(context.Background(), "timestamp", now)
		result := GetValue[time.Time](ctx, "timestamp")
		assert.Equal(t, now, result)
	})

	t.Run("key not found", func(t *testing.T) {
		ctx := context.Background()
		result := GetValue[string](ctx, "nonexistent")
		assert.Equal(t, "", result)
	})

	t.Run("nil value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", nil)
		result := GetValue[string](ctx, "key")
		assert.Equal(t, "", result)
	})

	t.Run("type mismatch returns zero value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", "string value")
		result := GetValue[int](ctx, "key")
		assert.Equal(t, 0, result)
	})

	t.Run("zero value for different types", func(t *testing.T) {
		ctx := context.Background()

		assert.Equal(t, "", GetValue[string](ctx, "key"))
		assert.Equal(t, 0, GetValue[int](ctx, "key"))
		assert.Equal(t, int64(0), GetValue[int64](ctx, "key"))
		assert.Equal(t, false, GetValue[bool](ctx, "key"))
		assert.Equal(t, 0.0, GetValue[float64](ctx, "key"))
		assert.Nil(t, GetValue[[]string](ctx, "key"))
		assert.Nil(t, GetValue[map[string]int](ctx, "key"))
		assert.Equal(t, time.Duration(0), GetValue[time.Duration](ctx, "key"))
		assert.Equal(t, time.Time{}, GetValue[time.Time](ctx, "key"))
	})

	t.Run("nested context values", func(t *testing.T) {
		ctx1 := context.WithValue(context.Background(), "key1", "value1")
		ctx2 := context.WithValue(ctx1, "key2", "value2")

		assert.Equal(t, "value1", GetValue[string](ctx2, "key1"))
		assert.Equal(t, "value2", GetValue[string](ctx2, "key2"))
	})

	t.Run("interface{} type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "data", "any value")
		result := GetValue[interface{}](ctx, "data")
		assert.Equal(t, "any value", result)
	})

	t.Run("any type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "data", 12345)
		result := GetValue[any](ctx, "data")
		assert.Equal(t, 12345, result)
	})
}

// TestGetValueEdgeCases 测试边界情况
func TestGetValueEdgeCases(t *testing.T) {
	t.Run("empty string key", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "", "value")
		result := GetValue[string](ctx, "")
		assert.Equal(t, "value", result)
	})

	t.Run("pointer value", func(t *testing.T) {
		str := "test"
		ctx := context.WithValue(context.Background(), "key", &str)
		result := GetValue[*string](ctx, "key")
		assert.NotNil(t, result)
		assert.Equal(t, "test", *result)
	})

	t.Run("nil pointer value", func(t *testing.T) {
		var ptr *string
		ctx := context.WithValue(context.Background(), "key", ptr)
		result := GetValue[*string](ctx, "key")
		assert.Nil(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		slice := []string{}
		ctx := context.WithValue(context.Background(), "key", slice)
		result := GetValue[[]string](ctx, "key")
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]int{}
		ctx := context.WithValue(context.Background(), "key", m)
		result := GetValue[map[string]int](ctx, "key")
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})
}

// BenchmarkGetValue 性能测试
func BenchmarkGetValue(b *testing.B) {
	ctx := context.WithValue(context.Background(), "key", "value")

	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = GetValue[string](ctx, "key")
		}
	})

	b.Run("int", func(b *testing.B) {
		ctx := context.WithValue(context.Background(), "count", 123)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GetValue[int](ctx, "count")
		}
	})

	b.Run("struct", func(b *testing.B) {
		type User struct {
			Name string
			Age  int
		}
		user := User{Name: "Alice", Age: 30}
		ctx := context.WithValue(context.Background(), "user", user)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GetValue[User](ctx, "user")
		}
	})

	b.Run("missing key", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = GetValue[string](ctx, "nonexistent")
		}
	})
}
