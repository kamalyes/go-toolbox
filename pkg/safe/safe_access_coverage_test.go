/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 23:30:00
 * @FilePath: \go-toolbox\pkg\safe\safe_access_coverage_test.go
 * @Description: SafeAccess 完整覆盖率测试 - 50个测试用例
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// ==================== 数值类型转换测试 (覆盖Int64, Int32, Uint, Uint64, Float32, Float64) ====================

func TestInt64_AllCases(t *testing.T) {
	t.Run("Int64正常转换", func(t *testing.T) {
		s := Safe(int64(9223372036854775807))
		assert.Equal(t, int64(9223372036854775807), s.Int64())
	})

	t.Run("Int64从字符串转换", func(t *testing.T) {
		s := Safe("12345678901234")
		assert.Equal(t, int64(12345678901234), s.Int64())
	})

	t.Run("Int64无效值使用默认", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, int64(999), s.Int64(999))
	})

	t.Run("Int64无效SafeAccess使用默认", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, int64(777), s.Int64(777))
	})
}

func TestInt32_AllCases(t *testing.T) {
	t.Run("Int32正常转换", func(t *testing.T) {
		s := Safe(int32(2147483647))
		assert.Equal(t, int32(2147483647), s.Int32())
	})

	t.Run("Int32从字符串转换", func(t *testing.T) {
		s := Safe("12345")
		assert.Equal(t, int32(12345), s.Int32())
	})

	t.Run("Int32无效值使用默认", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, int32(100), s.Int32(100))
	})

	t.Run("Int32无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, int32(0), s.Int32())
	})
}

func TestUint_AllCases(t *testing.T) {
	t.Run("Uint正常转换", func(t *testing.T) {
		s := Safe(uint(42))
		assert.Equal(t, uint(42), s.Uint())
	})

	t.Run("Uint从字符串转换", func(t *testing.T) {
		s := Safe("999")
		assert.Equal(t, uint(999), s.Uint())
	})

	t.Run("Uint无效值使用默认", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, uint(555), s.Uint(555))
	})

	t.Run("Uint无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, uint(0), s.Uint())
	})
}

func TestUint64_AllCases(t *testing.T) {
	t.Run("Uint64正常转换", func(t *testing.T) {
		s := Safe(uint64(18446744073709551615))
		assert.Equal(t, uint64(18446744073709551615), s.Uint64())
	})

	t.Run("Uint64从字符串转换", func(t *testing.T) {
		s := Safe("123456789")
		assert.Equal(t, uint64(123456789), s.Uint64())
	})

	t.Run("Uint64无效值使用默认", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, uint64(888), s.Uint64(888))
	})

	t.Run("Uint64无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, uint64(0), s.Uint64())
	})
}

func TestFloat32_AllCases(t *testing.T) {
	t.Run("Float32正常转换", func(t *testing.T) {
		s := Safe(float32(3.14159))
		result := s.Float32()
		assert.InDelta(t, 3.14159, result, 0.0001)
	})

	t.Run("Float32从字符串转换", func(t *testing.T) {
		s := Safe("2.718")
		result := s.Float32()
		assert.InDelta(t, 2.718, result, 0.001)
	})

	t.Run("Float32无效值使用默认", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, float32(1.5), s.Float32(1.5))
	})

	t.Run("Float32无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, float32(0), s.Float32())
	})
}

func TestFloat64_AllCases(t *testing.T) {
	t.Run("Float64正常转换", func(t *testing.T) {
		s := Safe(float64(3.141592653589793))
		assert.Equal(t, 3.141592653589793, s.Float64())
	})

	t.Run("Float64从字符串转换", func(t *testing.T) {
		s := Safe("2.718281828")
		result := s.Float64()
		assert.InDelta(t, 2.718281828, result, 0.0000001)
	})

	t.Run("Float64无效值使用默认", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, 9.99, s.Float64(9.99))
	})

	t.Run("Float64无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, 0.0, s.Float64())
	})
}

// ==================== String和Duration边界测试 ====================

func TestString_EdgeCases(t *testing.T) {
	t.Run("String从各种类型转换", func(t *testing.T) {
		assert.Equal(t, "123", Safe(123).String())
		assert.Equal(t, "3.14", Safe(3.14).String())
		assert.Equal(t, "true", Safe(true).String())
	})

	t.Run("String空值处理", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, "", s.String())
		assert.Equal(t, "default", s.String("default"))
	})

	t.Run("String从nil转换", func(t *testing.T) {
		s := Safe(nil)
		assert.Equal(t, "", s.String())
	})
}

func TestStringOr_AllCases(t *testing.T) {
	t.Run("StringOr有效值", func(t *testing.T) {
		s := Safe("hello")
		assert.Equal(t, "hello", s.StringOr("default"))
	})

	t.Run("StringOr空字符串使用默认", func(t *testing.T) {
		s := Safe("")
		assert.Equal(t, "default", s.StringOr("default"))
	})

	t.Run("StringOr无效值使用默认", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, "fallback", s.StringOr("fallback"))
	})

	t.Run("StringOr非字符串类型", func(t *testing.T) {
		s := Safe(123)
		assert.Equal(t, "123", s.StringOr("default"))
	})

	t.Run("StringOr nil值", func(t *testing.T) {
		s := Safe(nil)
		assert.Equal(t, "fallback", s.StringOr("fallback"))
	})
}

func TestDuration_AllCases(t *testing.T) {
	t.Run("Duration正常时长", func(t *testing.T) {
		d := 5 * time.Second
		s := Safe(d)
		assert.Equal(t, d, s.Duration())
	})

	t.Run("Duration从字符串解析", func(t *testing.T) {
		s := Safe("10s")
		assert.Equal(t, 10*time.Second, s.Duration())
	})

	t.Run("Duration从整数转换", func(t *testing.T) {
		s := Safe(1000000000) // 1秒的纳秒数
		assert.Equal(t, time.Second, s.Duration())
	})

	t.Run("Duration无效字符串使用默认", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, 3*time.Second, s.Duration(3*time.Second))
	})

	t.Run("Duration无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, time.Duration(0), s.Duration())
	})

	t.Run("Duration复杂时长", func(t *testing.T) {
		s := Safe("1h30m45s")
		expected := 1*time.Hour + 30*time.Minute + 45*time.Second
		assert.Equal(t, expected, s.Duration())
	})
}

// ==================== Value和OrElse测试 ====================

func TestValue_AllCases(t *testing.T) {
	t.Run("Value有效值", func(t *testing.T) {
		s := Safe("test")
		val := s.Value()
		assert.NotNil(t, val)
		assert.Equal(t, "test", val)
	})

	t.Run("Value无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		val := s.Value()
		assert.Nil(t, val)
	})

	t.Run("Value复杂对象", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		s := Safe(data)
		val := s.Value()
		assert.NotNil(t, val)
		assert.Equal(t, data, val)
	})
}

func TestOrElse_AllCases(t *testing.T) {
	t.Run("OrElse有效值不使用默认", func(t *testing.T) {
		s := Safe(42)
		result := s.OrElse(100)
		assert.Equal(t, 42, result.Int())
	})

	t.Run("OrElse无效值使用默认", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := s.OrElse("default")
		assert.Equal(t, "default", result.String())
	})

	t.Run("OrElse链式调用", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := s.OrElse(999).OrElse(111)
		assert.Equal(t, 999, result.Int())
	})
}

func TestMap_AllCases(t *testing.T) {
	t.Run("Map有效转换", func(t *testing.T) {
		s := Safe(10)
		result := s.Map(func(v interface{}) interface{} {
			return v.(int) * 2
		})
		assert.Equal(t, 20, result.Int())
	})

	t.Run("Map无效值不执行", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		executed := false
		result := s.Map(func(v interface{}) interface{} {
			executed = true
			return v
		})
		assert.False(t, executed)
		assert.False(t, result.IsValid())
	})

	t.Run("Map类型转换", func(t *testing.T) {
		s := Safe("123")
		result := s.Map(func(v interface{}) interface{} {
			return len(v.(string))
		})
		assert.Equal(t, 3, result.Int())
	})
}

// ==================== SafeGetXXX函数测试 ====================

func TestSafeGetString_AllCases(t *testing.T) {
	t.Run("SafeGetString正常获取", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		assert.Equal(t, "value", SafeGetString(data, "key"))
	})

	t.Run("SafeGetString键不存在", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		assert.Equal(t, "", SafeGetString(data, "missing"))
	})

	t.Run("SafeGetString非字符串转换", func(t *testing.T) {
		data := map[string]interface{}{"key": 123}
		assert.Equal(t, "123", SafeGetString(data, "key"))
	})
}

func TestSafeGetBool_AllCases(t *testing.T) {
	t.Run("SafeGetBool正常获取", func(t *testing.T) {
		data := map[string]interface{}{"key": true}
		assert.True(t, SafeGetBool(data, "key"))
	})

	t.Run("SafeGetBool键不存在", func(t *testing.T) {
		data := map[string]interface{}{"key": true}
		assert.False(t, SafeGetBool(data, "missing"))
	})

	t.Run("SafeGetBool从字符串转换", func(t *testing.T) {
		data := map[string]interface{}{"key": "true"}
		assert.True(t, SafeGetBool(data, "key"))
	})

	t.Run("SafeGetBool从整数转换", func(t *testing.T) {
		data := map[string]interface{}{"key": 1}
		assert.True(t, SafeGetBool(data, "key"))
	})
}

func TestSafeGetStringSlice_AllCases(t *testing.T) {
	t.Run("SafeGetStringSlice正常获取", func(t *testing.T) {
		data := map[string]interface{}{"key": []string{"a", "b", "c"}}
		result := SafeGetStringSlice(data, "key")
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("SafeGetStringSlice键不存在", func(t *testing.T) {
		data := map[string]interface{}{"key": []string{"a"}}
		result := SafeGetStringSlice(data, "missing")
		assert.Nil(t, result)
	})

	t.Run("SafeGetStringSlice从interface切片转换", func(t *testing.T) {
		data := map[string]interface{}{"key": []interface{}{"x", "y", "z"}}
		result := SafeGetStringSlice(data, "key")
		assert.Len(t, result, 3)
		assert.Contains(t, result, "x")
	})

	t.Run("SafeGetStringSlice非切片类型", func(t *testing.T) {
		data := map[string]interface{}{"key": "not-a-slice"}
		result := SafeGetStringSlice(data, "key")
		assert.Nil(t, result)
	})
}

// ==================== GetIntValue和GetInt64Value测试 ====================

func TestGetIntValue_AllCases(t *testing.T) {
	t.Run("GetIntValue从int", func(t *testing.T) {
		s := Safe(42)
		assert.Equal(t, 42, s.GetIntValue(0))
	})

	t.Run("GetIntValue从字符串", func(t *testing.T) {
		s := Safe("123")
		assert.Equal(t, 123, s.GetIntValue(0))
	})

	t.Run("GetIntValue从float", func(t *testing.T) {
		s := Safe(3.14)
		assert.Equal(t, 3, s.GetIntValue(0))
	})

	t.Run("GetIntValue无效值", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, 999, s.GetIntValue(999))
	})

	t.Run("GetIntValue从bool", func(t *testing.T) {
		s := Safe(true)
		// bool无法直接转int，应该返回默认值
		assert.Equal(t, 999, s.GetIntValue(999))
	})
}

func TestGetInt64Value_AllCases(t *testing.T) {
	t.Run("GetInt64Value从int64", func(t *testing.T) {
		s := Safe(int64(9223372036854775807))
		assert.Equal(t, int64(9223372036854775807), s.GetInt64Value(0))
	})

	t.Run("GetInt64Value从字符串", func(t *testing.T) {
		s := Safe("12345678901234")
		assert.Equal(t, int64(12345678901234), s.GetInt64Value(0))
	})

	t.Run("GetInt64Value从float", func(t *testing.T) {
		s := Safe(3.99)
		assert.Equal(t, int64(3), s.GetInt64Value(0))
	})

	t.Run("GetInt64Value无效值", func(t *testing.T) {
		s := Safe("invalid")
		assert.Equal(t, int64(888), s.GetInt64Value(888))
	})

	t.Run("GetInt64Value从int", func(t *testing.T) {
		s := Safe(42)
		assert.Equal(t, int64(42), s.GetInt64Value(0))
	})
}

// ==================== AsFloat边界测试 ====================

func TestAsFloat_EdgeCases(t *testing.T) {
	t.Run("AsFloat无默认值", func(t *testing.T) {
		s := Safe(3.14)
		result := AsFloat[float64](s, convert.RoundNone)
		assert.Equal(t, 3.14, result)
	})

	t.Run("AsFloat无效值返回零值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := AsFloat[float64](s, convert.RoundNone)
		assert.Equal(t, 0.0, result)
	})

	t.Run("AsFloat类型转换失败使用默认", func(t *testing.T) {
		s := Safe("invalid")
		result := AsFloat[float64](s, convert.RoundNone, 9.99)
		assert.Equal(t, 9.99, result)
	})
}

// ==================== 类型检查边界测试 ====================

func TestIsNumber_EdgeCases(t *testing.T) {
	t.Run("IsNumber无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.IsNumber())
	})

	t.Run("IsNumber各种整数类型", func(t *testing.T) {
		assert.True(t, Safe(int8(1)).IsNumber())
		assert.True(t, Safe(int16(1)).IsNumber())
		assert.True(t, Safe(uint8(1)).IsNumber())
		assert.True(t, Safe(uint16(1)).IsNumber())
		assert.True(t, Safe(uint32(1)).IsNumber())
	})
}

func TestIsString_EdgeCases(t *testing.T) {
	t.Run("IsString无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.IsString())
	})

	t.Run("IsString空字符串", func(t *testing.T) {
		assert.True(t, Safe("").IsString())
	})
}

func TestIsBool_EdgeCases(t *testing.T) {
	t.Run("IsBool无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.IsBool())
	})

	t.Run("IsBool true和false", func(t *testing.T) {
		assert.True(t, Safe(true).IsBool())
		assert.True(t, Safe(false).IsBool())
	})
}

func TestIsSlice_EdgeCases(t *testing.T) {
	t.Run("IsSlice无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.IsSlice())
	})

	t.Run("IsSlice各种切片类型", func(t *testing.T) {
		assert.True(t, Safe([]int{}).IsSlice())
		assert.True(t, Safe([]string{}).IsSlice())
		assert.True(t, Safe([]interface{}{}).IsSlice())
	})
}

func TestIsMap_EdgeCases(t *testing.T) {
	t.Run("IsMap无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.IsMap())
	})

	t.Run("IsMap各种map类型", func(t *testing.T) {
		assert.True(t, Safe(map[string]int{}).IsMap())
		assert.True(t, Safe(map[string]interface{}{}).IsMap())
		assert.True(t, Safe(map[int]string{}).IsMap())
	})
}

// ==================== Len边界测试 ====================

func TestLen_EdgeCases(t *testing.T) {
	t.Run("Len非容器类型", func(t *testing.T) {
		s := Safe(123)
		assert.Equal(t, 0, s.Len())
	})

	t.Run("Len nil值", func(t *testing.T) {
		s := Safe(nil)
		assert.Equal(t, 0, s.Len())
	})

	t.Run("Len数组类型", func(t *testing.T) {
		arr := [3]int{1, 2, 3}
		s := Safe(arr)
		assert.Equal(t, 3, s.Len())
	})
}

// ==================== Keys和Values边界测试 ====================

func TestKeys_EdgeCases(t *testing.T) {
	t.Run("Keys空map", func(t *testing.T) {
		s := Safe(map[string]interface{}{})
		keys := s.Keys()
		assert.NotNil(t, keys)
		assert.Len(t, keys, 0)
	})

	t.Run("Keys无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Nil(t, s.Keys())
	})
}

func TestValues_EdgeCases(t *testing.T) {
	t.Run("Values空map", func(t *testing.T) {
		s := Safe(map[string]interface{}{})
		values := s.Values()
		assert.NotNil(t, values)
		assert.Len(t, values, 0)
	})

	t.Run("Values无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Nil(t, s.Values())
	})
}

// ==================== Contains边界测试 ====================

func TestContains_EdgeCases(t *testing.T) {
	t.Run("Contains空切片", func(t *testing.T) {
		s := Safe([]int{})
		assert.False(t, s.Contains(1))
	})

	t.Run("Contains空map", func(t *testing.T) {
		s := Safe(map[string]interface{}{})
		assert.False(t, s.Contains("key"))
	})

	t.Run("Contains非容器类型", func(t *testing.T) {
		s := Safe(123)
		assert.False(t, s.Contains(123))
	})

	t.Run("Contains interface切片", func(t *testing.T) {
		s := Safe([]interface{}{1, "test", true})
		assert.True(t, s.Contains("test"))
		assert.False(t, s.Contains("missing"))
	})
}

// ==================== IsEmpty边界测试 ====================

func TestIsEmpty_AllTypes(t *testing.T) {
	t.Run("IsEmpty nil", func(t *testing.T) {
		s := Safe(nil)
		assert.True(t, s.IsEmpty())
	})

	t.Run("IsEmpty数组", func(t *testing.T) {
		arr := [3]int{1, 2, 3}
		s := Safe(arr)
		// 数组类型会被认为是非容器，返回false
		assert.False(t, s.IsEmpty())
	})

	t.Run("IsEmpty非容器类型", func(t *testing.T) {
		s := Safe(123)
		assert.False(t, s.IsEmpty())
	})
}

// ==================== Unless边界测试 ====================

func TestUnless_EdgeCases(t *testing.T) {
	t.Run("Unless无效值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := s.Unless(
			func(v interface{}) bool { return false },
			func(v interface{}) interface{} { return 100 },
		)
		assert.False(t, result.IsValid())
	})

	t.Run("Unless条件为false执行", func(t *testing.T) {
		s := Safe(10)
		result := s.Unless(
			func(v interface{}) bool { return v.(int) < 5 },
			func(v interface{}) interface{} { return v.(int) * 10 },
		)
		// 条件为false(10 < 5 = false)，所以执行转换函数
		assert.Equal(t, 100, result.Int())
	})
}

// ==================== Field边界测试 ====================

func TestField_EdgeCases(t *testing.T) {
	t.Run("Field访问私有字段", func(t *testing.T) {
		type TestStruct struct {
			Public  string
			private string
		}
		s := Safe(TestStruct{Public: "visible", private: "hidden"})
		assert.True(t, s.Field("Public").IsValid())
		assert.False(t, s.Field("private").IsValid())
	})

	t.Run("Field访问嵌套nil指针", func(t *testing.T) {
		type Inner struct {
			Value string
		}
		type Outer struct {
			Inner *Inner
		}
		s := Safe(Outer{Inner: nil})
		assert.False(t, s.Field("Inner").IsValid())
	})

	t.Run("Field访问非结构体", func(t *testing.T) {
		s := Safe("not a struct")
		assert.False(t, s.Field("anything").IsValid())
	})

	t.Run("Field访问数组", func(t *testing.T) {
		s := Safe([]int{1, 2, 3})
		assert.False(t, s.Field("0").IsValid())
	})
}

// ==================== At方法边界测试 ====================

func TestAt_EdgeCases(t *testing.T) {
	t.Run("At空路径返回默认值", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		s := Safe(data)
		result := s.At("", "default-value")
		// 空路径应该返回默认值
		assert.True(t, result.IsValid())
		assert.Equal(t, "default-value", result.String())
	})

	t.Run("At深层嵌套不存在", func(t *testing.T) {
		data := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": "value",
				},
			},
		}
		s := Safe(data)
		result := s.At("level1.level2.missing.level4")
		assert.False(t, result.IsValid())
	})

	t.Run("At中途遇到非map", func(t *testing.T) {
		data := map[string]interface{}{
			"level1": "string-value",
		}
		s := Safe(data)
		result := s.At("level1.level2")
		assert.False(t, result.IsValid())
	})
}

// ==================== AsSlice边界测试 ====================

func TestAsSlice_EdgeCases(t *testing.T) {
	t.Run("AsSlice空切片", func(t *testing.T) {
		s := Safe([]int{})
		result, err := AsSlice[int](s)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("AsSlice混合数值类型转换", func(t *testing.T) {
		s := Safe([]interface{}{1, "2", 3.0, "4"})
		result, err := AsSlice[int](s)
		assert.NoError(t, err)
		assert.Len(t, result, 4)
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})
}

// ==================== AsFloatSlice边界测试 ====================

func TestAsFloatSlice_EdgeCases(t *testing.T) {
	t.Run("AsFloatSlice空切片", func(t *testing.T) {
		s := Safe([]float64{})
		result, err := AsFloatSlice[float64](s, convert.RoundNone)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("AsFloatSlice混合类型", func(t *testing.T) {
		s := Safe([]interface{}{1, 2.5, "3.7"})
		result, err := AsFloatSlice[float64](s, convert.RoundNone)
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.InDelta(t, 3.7, result[2], 0.01)
	})
}

// ==================== AsStringSlice边界测试 ====================

func TestAsStringSlice_EdgeCases(t *testing.T) {
	t.Run("AsStringSlice空切片", func(t *testing.T) {
		s := Safe([]string{})
		result := s.AsStringSlice()
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("AsStringSlice非切片类型", func(t *testing.T) {
		s := Safe("not a slice")
		result := s.AsStringSlice()
		assert.Nil(t, result)
	})
}

// ==================== Bool边界测试 ====================

func TestBool_EdgeCases(t *testing.T) {
	t.Run("Bool从各种真值", func(t *testing.T) {
		assert.True(t, Safe(1).Bool())
		assert.True(t, Safe("true").Bool())
		assert.True(t, Safe("1").Bool())
		assert.True(t, Safe(true).Bool())
	})

	t.Run("Bool从各种假值", func(t *testing.T) {
		assert.False(t, Safe(0).Bool())
		assert.False(t, Safe("false").Bool())
		assert.False(t, Safe("0").Bool())
		assert.False(t, Safe(false).Bool())
	})

	t.Run("Bool无效值使用默认", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.True(t, s.Bool(true))
	})

	t.Run("Bool假值但有默认", func(t *testing.T) {
		s := Safe(false)
		assert.True(t, s.Bool(true))
	})
}

// ==================== Int边界测试 ====================

func TestInt_EdgeCases(t *testing.T) {
	t.Run("Int从负数", func(t *testing.T) {
		s := Safe(-999)
		assert.Equal(t, -999, s.Int())
	})

	t.Run("Int从大数", func(t *testing.T) {
		s := Safe(2147483647)
		assert.Equal(t, 2147483647, s.Int())
	})

	t.Run("Int无效转换使用默认", func(t *testing.T) {
		s := Safe("not-a-number")
		assert.Equal(t, 777, s.Int(777))
	})
}
