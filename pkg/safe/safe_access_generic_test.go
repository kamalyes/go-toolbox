/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 23:00:00
 * @FilePath: \go-toolbox\pkg\safe\safe_access_generic_test.go
 * @Description: SafeAccess 泛型和增强功能测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/stretchr/testify/assert"
)

// ==================== 泛型转换测试 ====================

func TestAs_GenericNumerical(t *testing.T) {
	t.Run("As[int]正常转换", func(t *testing.T) {
		s := Safe(42)
		result := As[int](s)
		assert.Equal(t, 42, result)
	})

	t.Run("As[int64]从int转换", func(t *testing.T) {
		s := Safe(100)
		result := As[int64](s)
		assert.Equal(t, int64(100), result)
	})

	t.Run("As[uint]从int转换", func(t *testing.T) {
		s := Safe(50)
		result := As[uint](s)
		assert.Equal(t, uint(50), result)
	})

	t.Run("As[int32]从字符串转换", func(t *testing.T) {
		s := Safe("123")
		result := As[int32](s)
		assert.Equal(t, int32(123), result)
	})

	t.Run("As[int]使用默认值", func(t *testing.T) {
		s := Safe("invalid")
		result := As[int](s, 999)
		assert.Equal(t, 999, result)
	})

	t.Run("As[int]无效值返回零值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := As[int](s)
		assert.Equal(t, 0, result)
	})

	t.Run("As[uint64]大数转换", func(t *testing.T) {
		s := Safe(uint64(18446744073709551615))
		result := As[uint64](s)
		assert.Equal(t, uint64(18446744073709551615), result)
	})
}

func TestAsFloat_GenericFloat(t *testing.T) {
	t.Run("AsFloat[float64]正常转换", func(t *testing.T) {
		s := Safe(3.14)
		result := AsFloat[float64](s, convert.RoundNone)
		assert.Equal(t, 3.14, result)
	})

	t.Run("AsFloat[float32]从int转换", func(t *testing.T) {
		s := Safe(42)
		result := AsFloat[float32](s, convert.RoundNone)
		assert.Equal(t, float32(42.0), result)
	})

	t.Run("AsFloat[float64]从字符串转换", func(t *testing.T) {
		s := Safe("3.14159")
		result := AsFloat[float64](s, convert.RoundNone)
		assert.InDelta(t, 3.14159, result, 0.00001)
	})

	t.Run("AsFloat[float64]四舍五入", func(t *testing.T) {
		s := Safe("3.6")
		result := AsFloat[float64](s, convert.RoundNearest)
		assert.Equal(t, 4.0, result)
	})

	t.Run("AsFloat[float64]向上取整", func(t *testing.T) {
		s := Safe("3.1")
		result := AsFloat[float64](s, convert.RoundUp)
		assert.Equal(t, 4.0, result)
	})

	t.Run("AsFloat[float64]向下取整", func(t *testing.T) {
		s := Safe("3.9")
		result := AsFloat[float64](s, convert.RoundDown)
		assert.Equal(t, 3.0, result)
	})

	t.Run("AsFloat[float32]使用默认值", func(t *testing.T) {
		s := Safe("invalid")
		result := AsFloat[float32](s, convert.RoundNone, 1.23)
		assert.Equal(t, float32(1.23), result)
	})
}

// ==================== 字符串和布尔转换测试 ====================

func TestAsString(t *testing.T) {
	t.Run("AsString从int转换", func(t *testing.T) {
		s := Safe(123)
		result := s.AsString()
		assert.Equal(t, "123", result)
	})

	t.Run("AsString从float转换", func(t *testing.T) {
		s := Safe(3.14)
		result := s.AsString()
		assert.Equal(t, "3.14", result)
	})

	t.Run("AsString从bool转换", func(t *testing.T) {
		s := Safe(true)
		result := s.AsString()
		assert.Equal(t, "true", result)
	})

	t.Run("AsString无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := s.AsString()
		assert.Equal(t, "", result)
	})
}

func TestAsBool(t *testing.T) {
	t.Run("AsBool从bool转换", func(t *testing.T) {
		s := Safe(true)
		assert.True(t, s.AsBool())
	})

	t.Run("AsBool从int转换", func(t *testing.T) {
		s := Safe(1)
		assert.True(t, s.AsBool())

		s = Safe(0)
		assert.False(t, s.AsBool())
	})

	t.Run("AsBool从字符串转换", func(t *testing.T) {
		s := Safe("true")
		assert.True(t, s.AsBool())

		s = Safe("false")
		assert.False(t, s.AsBool())

		s = Safe("1")
		assert.True(t, s.AsBool())
	})

	t.Run("AsBool无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.AsBool())
	})
}

func TestAsJSON(t *testing.T) {
	t.Run("AsJSON对象转换", func(t *testing.T) {
		data := map[string]interface{}{
			"name": "test",
			"age":  30,
		}
		s := Safe(data)
		result, err := s.AsJSON(false)
		assert.NoError(t, err)
		assert.Contains(t, result, "test")
		assert.Contains(t, result, "30")
	})

	t.Run("AsJSON缩进格式", func(t *testing.T) {
		data := map[string]interface{}{
			"key": "value",
		}
		s := Safe(data)
		result, err := s.AsJSON(true)
		assert.NoError(t, err)
		assert.Contains(t, result, "\n")
		assert.Contains(t, result, "  ")
	})

	t.Run("AsJSON无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		_, err := s.AsJSON(false)
		assert.Error(t, err)
	})
}

// ==================== 切片操作测试 ====================

func TestAsSlice_GenericNumerical(t *testing.T) {
	t.Run("AsSlice[int]已存在切片", func(t *testing.T) {
		s := Safe([]int{1, 2, 3, 4, 5})
		result, err := AsSlice[int](s)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("AsSlice[int64]从字符串切片转换", func(t *testing.T) {
		s := Safe([]string{"10", "20", "30"})
		result, err := AsSlice[int64](s)
		assert.NoError(t, err)
		assert.Equal(t, []int64{10, 20, 30}, result)
	})

	t.Run("AsSlice[uint]从interface切片转换", func(t *testing.T) {
		s := Safe([]interface{}{1, 2, 3})
		result, err := AsSlice[uint](s)
		assert.NoError(t, err)
		assert.Equal(t, []uint{1, 2, 3}, result)
	})

	t.Run("AsSlice[int]转换错误", func(t *testing.T) {
		s := Safe([]string{"invalid", "data"})
		_, err := AsSlice[int](s)
		assert.Error(t, err)
	})

	t.Run("AsSlice[int]无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		_, err := AsSlice[int](s)
		assert.Error(t, err)
	})
}

func TestAsFloatSlice(t *testing.T) {
	t.Run("AsFloatSlice[float64]已存在切片", func(t *testing.T) {
		s := Safe([]float64{1.1, 2.2, 3.3})
		result, err := AsFloatSlice[float64](s, convert.RoundNone)
		assert.NoError(t, err)
		assert.Equal(t, []float64{1.1, 2.2, 3.3}, result)
	})

	t.Run("AsFloatSlice[float32]从字符串切片转换", func(t *testing.T) {
		s := Safe([]string{"1.5", "2.5", "3.5"})
		result, err := AsFloatSlice[float32](s, convert.RoundNone)
		assert.NoError(t, err)
		assert.Equal(t, []float32{1.5, 2.5, 3.5}, result)
	})

	t.Run("AsFloatSlice[float64]四舍五入", func(t *testing.T) {
		s := Safe([]string{"1.4", "2.6", "3.5"})
		result, err := AsFloatSlice[float64](s, convert.RoundNearest)
		assert.NoError(t, err)
		assert.Equal(t, []float64{1.0, 3.0, 4.0}, result)
	})

	t.Run("AsFloatSlice[float64]从interface切片转换", func(t *testing.T) {
		s := Safe([]interface{}{1, 2.5, 3})
		result, err := AsFloatSlice[float64](s, convert.RoundNone)
		assert.NoError(t, err)
		assert.Equal(t, []float64{1.0, 2.5, 3.0}, result)
	})
}

func TestAsStringSlice(t *testing.T) {
	t.Run("AsStringSlice已存在切片", func(t *testing.T) {
		s := Safe([]string{"a", "b", "c"})
		result := s.AsStringSlice()
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("AsStringSlice从interface切片转换", func(t *testing.T) {
		s := Safe([]interface{}{1, 2.5, true, "test"})
		result := s.AsStringSlice()
		assert.Len(t, result, 4)
		assert.Contains(t, result, "1")
		assert.Contains(t, result, "test")
	})

	t.Run("AsStringSlice无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := s.AsStringSlice()
		assert.Nil(t, result)
	})
}

// ==================== 链式操作增强测试 ====================

func TestMap_Generic(t *testing.T) {
	t.Run("Map泛型转换", func(t *testing.T) {
		s := Safe(10)
		result := Map[int, string](s, func(v int) string {
			return convert.MustString(v * 2)
		})
		assert.True(t, result.IsValid())
		assert.Equal(t, "20", result.String())
	})

	t.Run("Map无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := Map[int, string](s, func(v int) string {
			return "test"
		})
		assert.False(t, result.IsValid())
	})

	t.Run("Map类型不匹配", func(t *testing.T) {
		s := Safe("string")
		result := Map[int, string](s, func(v int) string {
			return convert.MustString(v)
		})
		assert.False(t, result.IsValid())
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("FlatMap正常执行", func(t *testing.T) {
		s := Safe(5)
		result := s.FlatMap(func(v interface{}) *SafeAccess {
			if num, ok := v.(int); ok {
				return Safe(num * 3)
			}
			return &SafeAccess{valid: false}
		})
		assert.True(t, result.IsValid())
		assert.Equal(t, 15, result.Int())
	})

	t.Run("FlatMap无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := s.FlatMap(func(v interface{}) *SafeAccess {
			return Safe(100)
		})
		assert.False(t, result.IsValid())
	})
}

func TestOrDefault_Generic(t *testing.T) {
	t.Run("OrDefault有效值", func(t *testing.T) {
		s := Safe(42)
		result := OrDefault[int](s, 999)
		assert.Equal(t, 42, result)
	})

	t.Run("OrDefault无效值使用默认", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := OrDefault[int](s, 999)
		assert.Equal(t, 999, result)
	})

	t.Run("OrDefault类型不匹配使用默认", func(t *testing.T) {
		s := Safe("string")
		result := OrDefault[int](s, 999)
		assert.Equal(t, 999, result)
	})
}

func TestMust_Generic(t *testing.T) {
	t.Run("Must正常获取值", func(t *testing.T) {
		s := Safe(100)
		assert.NotPanics(t, func() {
			result := Must[int](s)
			assert.Equal(t, 100, result)
		})
	})

	t.Run("Must无效值panic", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Panics(t, func() {
			Must[int](s)
		})
	})

	t.Run("Must类型不匹配panic", func(t *testing.T) {
		s := Safe("string")
		assert.Panics(t, func() {
			Must[int](s)
		})
	})
}

// ==================== 条件操作测试 ====================

func TestWhen(t *testing.T) {
	t.Run("When条件为真执行", func(t *testing.T) {
		s := Safe(10)
		result := s.When(
			func(v interface{}) bool {
				return v.(int) > 5
			},
			func(v interface{}) interface{} {
				return v.(int) * 2
			},
		)
		assert.Equal(t, 20, result.Int())
	})

	t.Run("When条件为假不执行", func(t *testing.T) {
		s := Safe(3)
		result := s.When(
			func(v interface{}) bool {
				return v.(int) > 5
			},
			func(v interface{}) interface{} {
				return v.(int) * 2
			},
		)
		assert.Equal(t, 3, result.Int())
	})

	t.Run("When无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := s.When(
			func(v interface{}) bool { return true },
			func(v interface{}) interface{} { return 100 },
		)
		assert.False(t, result.IsValid())
	})
}

func TestUnless(t *testing.T) {
	t.Run("Unless条件为假执行", func(t *testing.T) {
		s := Safe(3)
		result := s.Unless(
			func(v interface{}) bool {
				return v.(int) > 5
			},
			func(v interface{}) interface{} {
				return v.(int) * 2
			},
		)
		assert.Equal(t, 6, result.Int())
	})

	t.Run("Unless条件为真不执行", func(t *testing.T) {
		s := Safe(10)
		result := s.Unless(
			func(v interface{}) bool {
				return v.(int) > 5
			},
			func(v interface{}) interface{} {
				return v.(int) * 2
			},
		)
		assert.Equal(t, 10, result.Int())
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("IsEmpty字符串", func(t *testing.T) {
		assert.True(t, Safe("").IsEmpty())
		assert.False(t, Safe("test").IsEmpty())
	})

	t.Run("IsEmpty切片", func(t *testing.T) {
		assert.True(t, Safe([]string{}).IsEmpty())
		assert.False(t, Safe([]string{"a"}).IsEmpty())
	})

	t.Run("IsEmpty map", func(t *testing.T) {
		assert.True(t, Safe(map[string]interface{}{}).IsEmpty())
		assert.False(t, Safe(map[string]interface{}{"key": "value"}).IsEmpty())
	})

	t.Run("IsEmpty无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.True(t, s.IsEmpty())
	})
}

func TestIsNonEmpty(t *testing.T) {
	t.Run("IsNonEmpty", func(t *testing.T) {
		assert.True(t, Safe("test").IsNonEmpty())
		assert.False(t, Safe("").IsNonEmpty())
	})
}

// ==================== 类型检查测试 ====================

func TestIsType_Generic(t *testing.T) {
	t.Run("IsType匹配", func(t *testing.T) {
		s := Safe(42)
		assert.True(t, IsType[int](s))
		assert.False(t, IsType[string](s))
	})

	t.Run("IsType无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, IsType[int](s))
	})
}

func TestIsNumber(t *testing.T) {
	t.Run("IsNumber各种数值类型", func(t *testing.T) {
		assert.True(t, Safe(42).IsNumber())
		assert.True(t, Safe(int64(100)).IsNumber())
		assert.True(t, Safe(uint(50)).IsNumber())
		assert.True(t, Safe(3.14).IsNumber())
		assert.True(t, Safe(float32(2.5)).IsNumber())
		assert.False(t, Safe("123").IsNumber())
		assert.False(t, Safe(true).IsNumber())
	})
}

func TestIsString(t *testing.T) {
	t.Run("IsString", func(t *testing.T) {
		assert.True(t, Safe("test").IsString())
		assert.False(t, Safe(123).IsString())
	})
}

func TestIsBool(t *testing.T) {
	t.Run("IsBool", func(t *testing.T) {
		assert.True(t, Safe(true).IsBool())
		assert.False(t, Safe(1).IsBool())
	})
}

func TestIsSlice(t *testing.T) {
	t.Run("IsSlice", func(t *testing.T) {
		assert.True(t, Safe([]int{1, 2, 3}).IsSlice())
		assert.True(t, Safe([]string{"a", "b"}).IsSlice())
		assert.False(t, Safe(123).IsSlice())
	})
}

func TestIsMap(t *testing.T) {
	t.Run("IsMap", func(t *testing.T) {
		assert.True(t, Safe(map[string]interface{}{}).IsMap())
		assert.False(t, Safe([]int{}).IsMap())
	})
}

// ==================== 集合操作测试 ====================

func TestLen(t *testing.T) {
	t.Run("Len字符串", func(t *testing.T) {
		assert.Equal(t, 5, Safe("hello").Len())
	})

	t.Run("Len切片", func(t *testing.T) {
		assert.Equal(t, 3, Safe([]int{1, 2, 3}).Len())
	})

	t.Run("Len map", func(t *testing.T) {
		assert.Equal(t, 2, Safe(map[string]interface{}{"a": 1, "b": 2}).Len())
	})

	t.Run("Len无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, 0, s.Len())
	})
}

func TestKeys(t *testing.T) {
	t.Run("Keys正常获取", func(t *testing.T) {
		m := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		s := Safe(m)
		keys := s.Keys()
		assert.Len(t, keys, 3)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
		assert.Contains(t, keys, "key3")
	})

	t.Run("Keys非map返回nil", func(t *testing.T) {
		s := Safe([]int{1, 2, 3})
		assert.Nil(t, s.Keys())
	})
}

func TestValues(t *testing.T) {
	t.Run("Values正常获取", func(t *testing.T) {
		m := map[string]interface{}{
			"key1": 100,
			"key2": "test",
		}
		s := Safe(m)
		values := s.Values()
		assert.Len(t, values, 2)
	})

	t.Run("Values非map返回nil", func(t *testing.T) {
		s := Safe(123)
		assert.Nil(t, s.Values())
	})
}

func TestContains(t *testing.T) {
	t.Run("Contains在map中", func(t *testing.T) {
		m := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		s := Safe(m)
		assert.True(t, s.Contains("key1"))
		assert.False(t, s.Contains("key3"))
	})

	t.Run("Contains在切片中", func(t *testing.T) {
		s := Safe([]int{1, 2, 3, 4, 5})
		assert.True(t, s.Contains(3))
		assert.False(t, s.Contains(10))
	})

	t.Run("Contains在字符串切片中", func(t *testing.T) {
		s := Safe([]string{"apple", "banana", "orange"})
		assert.True(t, s.Contains("banana"))
		assert.False(t, s.Contains("grape"))
	})

	t.Run("Contains无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.Contains("anything"))
	})
}

// ==================== 综合场景测试 ====================

func TestComplexScenario_ChainOperations(t *testing.T) {
	t.Run("链式调用组合", func(t *testing.T) {
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"age":   "25",
				"score": []string{"85", "90", "95"},
			},
		}

		s := Safe(data)

		// 测试链式访问
		age := As[int](s.Field("user").Field("age"))
		assert.Equal(t, 25, age)

		// 测试切片转换
		scores, err := AsSlice[int](s.Field("user").Field("score"))
		assert.NoError(t, err)
		assert.Equal(t, []int{85, 90, 95}, scores)
	})

	t.Run("条件转换组合", func(t *testing.T) {
		s := Safe(100)
		result := s.
			When(func(v interface{}) bool {
				return v.(int) > 50
			}, func(v interface{}) interface{} {
				return v.(int) / 2
			}).
			Map(func(v interface{}) interface{} {
				return v.(int) + 10
			})

		assert.Equal(t, 60, result.Int())
	})
}

func TestComplexScenario_DataTransform(t *testing.T) {
	t.Run("数据转换场景", func(t *testing.T) {
		// 模拟配置数据
		config := map[string]interface{}{
			"server": map[string]interface{}{
				"port":    "8080",
				"timeout": "30",
				"enabled": "true",
			},
			"features": []interface{}{"feature1", "feature2", "feature3"},
		}

		s := Safe(config)

		// 获取服务器配置
		port := As[int](s.Field("server").Field("port"))
		assert.Equal(t, 8080, port)

		timeout := As[int](s.Field("server").Field("timeout"))
		assert.Equal(t, 30, timeout)

		enabled := s.Field("server").Field("enabled").AsBool()
		assert.True(t, enabled)

		// 获取特性列表
		features := s.Field("features").AsStringSlice()
		assert.Len(t, features, 3)
		assert.Contains(t, features, "feature1")
	})
}

func TestComplexScenario_ErrorHandling(t *testing.T) {
	t.Run("错误处理场景", func(t *testing.T) {
		data := map[string]interface{}{
			"valid":   "123",
			"invalid": "not-a-number",
		}

		s := Safe(data)

		// 有效数据转换
		validNum := As[int](s.Field("valid"))
		assert.Equal(t, 123, validNum)

		// 无效数据使用默认值
		invalidNum := As[int](s.Field("invalid"), 999)
		assert.Equal(t, 999, invalidNum)

		// 不存在的字段
		missing := As[int](s.Field("missing"), 777)
		assert.Equal(t, 777, missing)
	})
}
