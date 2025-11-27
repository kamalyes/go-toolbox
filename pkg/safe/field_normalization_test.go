/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-27 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-27 23:31:05
 * @FilePath: \go-toolbox\pkg\safe\field_normalization_test.go
 * @Description: 字段名规范化函数100%覆盖率测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// ==================== stringx.NormalizeFieldName 完整覆盖测试 ====================

func TestNormalizeFieldName_AllCases(t *testing.T) {
	t.Run("空字符串", func(t *testing.T) {
		variants := stringx.NormalizeFieldName("")
		assert.Empty(t, variants)
	})

	t.Run("普通驼峰式", func(t *testing.T) {
		variants := stringx.NormalizeFieldName("userName")
		assert.Contains(t, variants, "userName")
		assert.Contains(t, variants, "UserName")  // PascalCase
		assert.Contains(t, variants, "user_name") // snake_case
		assert.Contains(t, variants, "user-name") // kebab-case
	})

	t.Run("下划线格式", func(t *testing.T) {
		variants := stringx.NormalizeFieldName("user_name")
		assert.Contains(t, variants, "user_name")
		assert.Contains(t, variants, "UserName")
		assert.Contains(t, variants, "userName")
	})

	t.Run("连字符格式", func(t *testing.T) {
		variants := stringx.NormalizeFieldName("user-name")
		assert.Contains(t, variants, "user-name")
		assert.Contains(t, variants, "UserName")
		assert.Contains(t, variants, "userName")
	})

	t.Run("PascalCase格式", func(t *testing.T) {
		variants := stringx.NormalizeFieldName("UserName")
		assert.Contains(t, variants, "UserName")
		assert.Contains(t, variants, "userName")
		assert.Contains(t, variants, "user_name")
		assert.Contains(t, variants, "user-name")
	})

	t.Run("单个字母", func(t *testing.T) {
		variants := stringx.NormalizeFieldName("a")
		assert.Contains(t, variants, "a")
		assert.Contains(t, variants, "A")
	})

	t.Run("复杂字段名", func(t *testing.T) {
		variants := stringx.NormalizeFieldName("http_server_port")
		assert.Contains(t, variants, "http_server_port")
		assert.Contains(t, variants, "HttpServerPort")
		assert.Contains(t, variants, "httpServerPort")
		assert.Contains(t, variants, "http-server-port")
	})
}

// ==================== toPascalCase 完整覆盖测试 ====================

func TestToPascalCase_AllCases(t *testing.T) {
	t.Run("下划线分隔", func(t *testing.T) {
		assert.Equal(t, "UserName", stringx.ToPascalCase("user_name"))
		assert.Equal(t, "HttpPort", stringx.ToPascalCase("http_port"))
	})

	t.Run("连字符分隔", func(t *testing.T) {
		assert.Equal(t, "UserName", stringx.ToPascalCase("user-name"))
		assert.Equal(t, "HttpPort", stringx.ToPascalCase("http-port"))
	})

	t.Run("混合分隔符", func(t *testing.T) {
		assert.Equal(t, "UserNameAge", stringx.ToPascalCase("user_name-age"))
	})

	t.Run("已经是PascalCase", func(t *testing.T) {
		assert.Equal(t, "UserName", stringx.ToPascalCase("UserName"))
	})

	t.Run("camelCase转换", func(t *testing.T) {
		assert.Equal(t, "UserName", stringx.ToPascalCase("userName"))
	})

	t.Run("空字符串", func(t *testing.T) {
		assert.Equal(t, "", stringx.ToPascalCase(""))
	})

	t.Run("单字母", func(t *testing.T) {
		assert.Equal(t, "A", stringx.ToPascalCase("a"))
	})

	t.Run("全大写", func(t *testing.T) {
		assert.Equal(t, "HTTP", stringx.ToPascalCase("HTTP"))
	})
}

// ==================== toCamelCase 完整覆盖测试 ====================

func TestToCamelCase_AllCases(t *testing.T) {
	t.Run("下划线分隔", func(t *testing.T) {
		assert.Equal(t, "userName", stringx.ToCamelCase("user_name"))
	})

	t.Run("连字符分隔", func(t *testing.T) {
		assert.Equal(t, "userName", stringx.ToCamelCase("user-name"))
	})

	t.Run("PascalCase转换", func(t *testing.T) {
		assert.Equal(t, "userName", stringx.ToCamelCase("UserName"))
	})

	t.Run("已经是camelCase", func(t *testing.T) {
		assert.Equal(t, "userName", stringx.ToCamelCase("userName"))
	})

	t.Run("空字符串", func(t *testing.T) {
		assert.Equal(t, "", stringx.ToCamelCase(""))
	})

	t.Run("单字母", func(t *testing.T) {
		assert.Equal(t, "a", stringx.ToCamelCase("A"))
	})
}

// ==================== toSnakeCase 完整覆盖测试 ====================

func TestToSnakeCase_AllCases(t *testing.T) {
	t.Run("PascalCase转换", func(t *testing.T) {
		assert.Equal(t, "user_name", stringx.ToSnakeCase("UserName"))
	})

	t.Run("camelCase转换", func(t *testing.T) {
		assert.Equal(t, "user_name", stringx.ToSnakeCase("userName"))
	})

	t.Run("连字符转下划线", func(t *testing.T) {
		assert.Equal(t, "user_name", stringx.ToSnakeCase("user-name"))
	})

	t.Run("已经是snake_case", func(t *testing.T) {
		assert.Equal(t, "user_name", stringx.ToSnakeCase("user_name"))
	})

	t.Run("连续大写字母", func(t *testing.T) {
		// 连续大写会被视为一个单词，转换后首字母小写
		assert.Equal(t, "httpserver", stringx.ToSnakeCase("HTTPServer"))
	})

	t.Run("首字母大写", func(t *testing.T) {
		assert.Equal(t, "name", stringx.ToSnakeCase("Name"))
	})

	t.Run("空字符串", func(t *testing.T) {
		assert.Equal(t, "", stringx.ToSnakeCase(""))
	})

	t.Run("单字母", func(t *testing.T) {
		assert.Equal(t, "a", stringx.ToSnakeCase("a"))
		assert.Equal(t, "a", stringx.ToSnakeCase("A"))
	})

	t.Run("混合情况", func(t *testing.T) {
		// 连续大写后跟小写会在第一个大写前加下划线
		assert.Equal(t, "my_varname", stringx.ToSnakeCase("myVARName"))
	})
}

// ==================== toKebabCase 完整覆盖测试 ====================

func TestToKebabCase_AllCases(t *testing.T) {
	t.Run("PascalCase转换", func(t *testing.T) {
		assert.Equal(t, "user-name", stringx.ToKebabCase("UserName"))
	})

	t.Run("camelCase转换", func(t *testing.T) {
		assert.Equal(t, "user-name", stringx.ToKebabCase("userName"))
	})

	t.Run("下划线转连字符", func(t *testing.T) {
		assert.Equal(t, "user-name", stringx.ToKebabCase("user_name"))
	})

	t.Run("已经是kebab-case", func(t *testing.T) {
		assert.Equal(t, "user-name", stringx.ToKebabCase("user-name"))
	})
}

// ==================== Field方法多种命名风格测试 ====================

func TestField_NamingConventions(t *testing.T) {
	// 测试结构体字段
	type Config struct {
		ServerPort int
		DbHost     string
		UiPath     string
	}

	config := Config{
		ServerPort: 8080,
		DbHost:     "localhost",
		UiPath:     "/ui",
	}

	s := Safe(config)

	t.Run("PascalCase访问", func(t *testing.T) {
		assert.Equal(t, 8080, s.Field("ServerPort").Int())
		assert.Equal(t, "localhost", s.Field("DbHost").String())
		assert.Equal(t, "/ui", s.Field("UiPath").String())
	})

	t.Run("camelCase访问", func(t *testing.T) {
		assert.Equal(t, 8080, s.Field("serverPort").Int())
		assert.Equal(t, "localhost", s.Field("dbHost").String())
		assert.Equal(t, "/ui", s.Field("uiPath").String())
	})

	t.Run("snake_case访问", func(t *testing.T) {
		assert.Equal(t, 8080, s.Field("server_port").Int())
		assert.Equal(t, "localhost", s.Field("db_host").String())
		assert.Equal(t, "/ui", s.Field("ui_path").String())
	})

	t.Run("kebab-case访问", func(t *testing.T) {
		assert.Equal(t, 8080, s.Field("server-port").Int())
		assert.Equal(t, "localhost", s.Field("db-host").String())
		assert.Equal(t, "/ui", s.Field("ui-path").String())
	})

	// 测试map字段
	t.Run("map多种命名访问", func(t *testing.T) {
		data := map[string]interface{}{
			"serverPort": 9090,
			"db_host":    "127.0.0.1",
			"ui-path":    "/admin",
		}
		sm := Safe(data)

		// camelCase原始键
		assert.Equal(t, 9090, sm.Field("serverPort").Int())
		assert.Equal(t, 9090, sm.Field("server_port").Int())
		assert.Equal(t, 9090, sm.Field("server-port").Int())
		assert.Equal(t, 9090, sm.Field("ServerPort").Int())

		// snake_case原始键
		assert.Equal(t, "127.0.0.1", sm.Field("db_host").String())
		assert.Equal(t, "127.0.0.1", sm.Field("dbHost").String())
		assert.Equal(t, "127.0.0.1", sm.Field("DbHost").String())

		// kebab-case原始键
		assert.Equal(t, "/admin", sm.Field("ui-path").String())
		assert.Equal(t, "/admin", sm.Field("uiPath").String())
		assert.Equal(t, "/admin", sm.Field("UiPath").String())
	})
}

// ==================== Field访问nil map和nil值测试 ====================

func TestField_NilCases(t *testing.T) {
	t.Run("Field访问nil SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false, value: nil}
		result := s.Field("anything")
		assert.False(t, result.IsValid())
	})

	t.Run("Field访问值为nil的SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: true, value: nil}
		result := s.Field("anything")
		assert.False(t, result.IsValid())
	})

	t.Run("Field访问map中值为nil的键", func(t *testing.T) {
		data := map[string]interface{}{
			"key": nil,
		}
		s := Safe(data)
		result := s.Field("key")
		assert.False(t, result.IsValid())
	})

	t.Run("Field访问nil指针结构体", func(t *testing.T) {
		type TestStruct struct {
			Name string
		}
		var ptr *TestStruct = nil
		s := Safe(ptr)
		result := s.Field("Name")
		assert.False(t, result.IsValid())
	})

	t.Run("Field访问结构体中nil指针字段", func(t *testing.T) {
		type Inner struct {
			Value string
		}
		type Outer struct {
			Inner *Inner
		}
		outer := Outer{Inner: nil}
		s := Safe(outer)
		result := s.Field("Inner")
		assert.False(t, result.IsValid())
	})

	t.Run("Field访问有效指针字段", func(t *testing.T) {
		type Inner struct {
			Value string
		}
		type Outer struct {
			Inner *Inner
		}
		outer := Outer{Inner: &Inner{Value: "test"}}
		s := Safe(outer)
		result := s.Field("Inner").Field("Value")
		assert.True(t, result.IsValid())
		assert.Equal(t, "test", result.String())
	})
}

// ==================== Duration全覆盖测试 ====================

func TestDuration_PointerTypes(t *testing.T) {
	t.Run("Duration指针类型", func(t *testing.T) {
		d := 5 * time.Second
		s := Safe(&d)
		assert.Equal(t, 5*time.Second, s.Duration())
	})

	t.Run("Duration nil指针", func(t *testing.T) {
		var d *time.Duration = nil
		s := Safe(d)
		assert.Equal(t, 10*time.Second, s.Duration(10*time.Second))
	})

	t.Run("Duration复杂解析", func(t *testing.T) {
		s := Safe("2h30m")
		expected := 2*time.Hour + 30*time.Minute
		assert.Equal(t, expected, s.Duration())
	})

	t.Run("Duration解析失败", func(t *testing.T) {
		s := Safe("invalid-duration")
		assert.Equal(t, 7*time.Second, s.Duration(7*time.Second))
	})
}

// ==================== splitFieldPath 完整测试 ====================

func TestSplitFieldPath_AllCases(t *testing.T) {
	t.Run("splitFieldPath空字符串", func(t *testing.T) {
		result := splitFieldPath("")
		assert.Empty(t, result)
	})

	t.Run("splitFieldPath单个字段", func(t *testing.T) {
		result := splitFieldPath("field")
		assert.Equal(t, []string{"field"}, result)
	})

	t.Run("splitFieldPath多级路径", func(t *testing.T) {
		result := splitFieldPath("level1.level2.level3")
		assert.Equal(t, []string{"level1", "level2", "level3"}, result)
	})
}

// ==================== As和AsFloat未覆盖分支 ====================

func TestAs_UncoveredBranches(t *testing.T) {
	t.Run("As无效SafeAccess无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := As[int](s)
		assert.Equal(t, 0, result)
	})

	t.Run("As转换错误无默认值", func(t *testing.T) {
		s := Safe("not-a-number")
		result := As[int](s)
		assert.Equal(t, 0, result)
	})
}

func TestAsFloat_UncoveredBranches(t *testing.T) {
	t.Run("AsFloat无效SafeAccess有默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result := AsFloat[float64](s, convert.RoundNone, 3.14)
		assert.Equal(t, 3.14, result)
	})

	t.Run("AsFloat转换错误无默认值", func(t *testing.T) {
		s := Safe("not-a-number")
		result := AsFloat[float64](s, convert.RoundNone)
		assert.Equal(t, 0.0, result)
	})
}

// ==================== AsSlice和AsFloatSlice错误分支 ====================

func TestAsSlice_ErrorBranches(t *testing.T) {
	t.Run("AsSlice无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result, err := AsSlice[int](s)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("AsSlice不可转换类型", func(t *testing.T) {
		s := Safe("not-a-slice")
		result, err := AsSlice[int](s)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("AsSlice interface切片转换错误", func(t *testing.T) {
		s := Safe([]interface{}{"not", "numbers"})
		result, err := AsSlice[int](s)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("AsSlice字符串切片转换", func(t *testing.T) {
		s := Safe([]string{"1", "2", "3"})
		result, err := AsSlice[int](s)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func TestAsFloatSlice_ErrorBranches(t *testing.T) {
	t.Run("AsFloatSlice无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		result, err := AsFloatSlice[float64](s, convert.RoundNone)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("AsFloatSlice不可转换类型", func(t *testing.T) {
		s := Safe("not-a-slice")
		result, err := AsFloatSlice[float64](s, convert.RoundNone)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("AsFloatSlice interface切片转换错误", func(t *testing.T) {
		s := Safe([]interface{}{"not", "numbers"})
		result, err := AsFloatSlice[float64](s, convert.RoundNone)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("AsFloatSlice字符串切片转换", func(t *testing.T) {
		s := Safe([]string{"1.1", "2.2", "3.3"})
		result, err := AsFloatSlice[float64](s, convert.RoundNone)
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.InDelta(t, 1.1, result[0], 0.01)
	})
}

// ==================== IsEmpty数组和默认分支 ====================

func TestIsEmpty_ArrayAndDefault(t *testing.T) {
	t.Run("IsEmpty字符串切片", func(t *testing.T) {
		s := Safe([]string{})
		assert.True(t, s.IsEmpty())
	})

	t.Run("IsEmpty非空字符串切片", func(t *testing.T) {
		s := Safe([]string{"item"})
		assert.False(t, s.IsEmpty())
	})

	t.Run("IsEmpty interface切片", func(t *testing.T) {
		s := Safe([]interface{}{})
		assert.True(t, s.IsEmpty())
	})

	t.Run("IsEmpty非空interface切片", func(t *testing.T) {
		s := Safe([]interface{}{1, 2})
		assert.False(t, s.IsEmpty())
	})

	t.Run("IsEmpty map", func(t *testing.T) {
		s := Safe(map[string]interface{}{})
		assert.True(t, s.IsEmpty())
	})

	t.Run("IsEmpty非空map", func(t *testing.T) {
		s := Safe(map[string]interface{}{"key": "value"})
		assert.False(t, s.IsEmpty())
	})

	t.Run("IsEmpty空字符串", func(t *testing.T) {
		s := Safe("")
		assert.True(t, s.IsEmpty())
	})

	t.Run("IsEmpty非空字符串", func(t *testing.T) {
		s := Safe("hello")
		assert.False(t, s.IsEmpty())
	})

	t.Run("IsEmpty其他类型", func(t *testing.T) {
		s := Safe(123)
		assert.False(t, s.IsEmpty())
	})
}

// ==================== Len各种类型测试 ====================

func TestLen_AllTypes(t *testing.T) {
	t.Run("Len字符串", func(t *testing.T) {
		s := Safe("hello")
		assert.Equal(t, 5, s.Len())
	})

	t.Run("Len字符串切片", func(t *testing.T) {
		s := Safe([]string{"a", "b", "c"})
		assert.Equal(t, 3, s.Len())
	})

	t.Run("Len interface切片", func(t *testing.T) {
		s := Safe([]interface{}{1, 2, 3, 4})
		assert.Equal(t, 4, s.Len())
	})

	t.Run("Len map", func(t *testing.T) {
		s := Safe(map[string]interface{}{"a": 1, "b": 2})
		assert.Equal(t, 2, s.Len())
	})

	t.Run("Len其他切片通过反射", func(t *testing.T) {
		s := Safe([]int{1, 2, 3, 4, 5})
		assert.Equal(t, 5, s.Len())
	})

	t.Run("Len数组通过反射", func(t *testing.T) {
		arr := [5]int{1, 2, 3, 4, 5}
		s := Safe(arr)
		assert.Equal(t, 5, s.Len())
	})

	t.Run("Len字符串通过反射", func(t *testing.T) {
		s := Safe("测试")
		// 中文字符串的字节长度是6 (UTF-8编码，每个中文字符3字节)
		assert.Equal(t, 6, s.Len())
	})

	t.Run("Len非容器类型", func(t *testing.T) {
		s := Safe(123)
		assert.Equal(t, 0, s.Len())
	})

	t.Run("Len无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, 0, s.Len())
	})
}

// ==================== Contains map非字符串键测试 ====================

func TestContains_MapNonStringKey(t *testing.T) {
	t.Run("Contains map使用非字符串target", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		s := Safe(data)
		// target不是字符串类型
		assert.False(t, s.Contains(123))
	})

	t.Run("Contains map存在的键", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		s := Safe(data)
		assert.True(t, s.Contains("key"))
	})

	t.Run("Contains无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.Contains("anything"))
	})
}

// ==================== At路径访问空路径无默认值 ====================

func TestAt_EmptyPathNoDefault(t *testing.T) {
	t.Run("At空路径无默认值", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		s := Safe(data)
		result := s.At("")
		assert.False(t, result.IsValid())
	})

	t.Run("At路径中途中断", func(t *testing.T) {
		data := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": "value",
			},
		}
		s := Safe(data)
		// level3不存在，没有默认值
		result := s.At("level1.level3")
		assert.False(t, result.IsValid())
	})
}

// ==================== Bool未覆盖分支 ====================

func TestBool_UncoveredBranches(t *testing.T) {
	t.Run("Bool无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.False(t, s.Bool())
	})

	t.Run("Bool转换成功", func(t *testing.T) {
		s := Safe("yes")
		// convert.MustBool会把"yes"转换为true
		assert.True(t, s.Bool())
	})
}

// ==================== Int各类型未覆盖分支 ====================

func TestInt_UncoveredBranches(t *testing.T) {
	t.Run("Int无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, 0, s.Int())
	})
}

func TestInt64_UncoveredBranches(t *testing.T) {
	t.Run("Int64无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, int64(0), s.Int64())
	})
}

func TestInt32_UncoveredBranches(t *testing.T) {
	t.Run("Int32无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, int32(0), s.Int32())
	})
}

func TestUint_UncoveredBranches(t *testing.T) {
	t.Run("Uint无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, uint(0), s.Uint())
	})
}

func TestUint64_UncoveredBranches(t *testing.T) {
	t.Run("Uint64无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, uint64(0), s.Uint64())
	})
}

func TestFloat32_UncoveredBranches(t *testing.T) {
	t.Run("Float32无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, float32(0), s.Float32())
	})
}

func TestFloat64_UncoveredBranches(t *testing.T) {
	t.Run("Float64无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, 0.0, s.Float64())
	})
}

func TestString_UncoveredBranches(t *testing.T) {
	t.Run("String无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, "", s.String())
	})
}

func TestDuration_UncoveredBranches(t *testing.T) {
	t.Run("Duration无效且无默认值", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, time.Duration(0), s.Duration())
	})
}

func TestGetIntValue_UncoveredBranches(t *testing.T) {
	t.Run("GetIntValue无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, 999, s.GetIntValue(999))
	})
}

func TestGetInt64Value_UncoveredBranches(t *testing.T) {
	t.Run("GetInt64Value无效SafeAccess", func(t *testing.T) {
		s := &SafeAccess{valid: false}
		assert.Equal(t, int64(888), s.GetInt64Value(888))
	})
}
