/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-10 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-10 00:00:00
 * @FilePath: \go-toolbox\pkg\types\type_check_test.go
 * @Description: type_check.go 的单元测试，覆盖 CheckTypeCompatibility 类型兼容性检查
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package types

import (
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 辅助函数：获取类型的 reflect.Type
func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// MyInt 和 YourInt 是底层类型相同的命名类型，用于测试同 Kind 不同类型的兼容性
type MyInt int
type YourInt int

// TestCheckTypeCompatibilityNil 测试 srcType 或 dstType 为 nil 的情况
func TestCheckTypeCompatibilityNil(t *testing.T) {
	t.Run("srcType 为 nil", func(t *testing.T) {
		err := CheckTypeCompatibility(nil, typeOf[int]())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type mismatch")
	})

	t.Run("dstType 为 nil", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[int](), nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type mismatch")
	})

	t.Run("两者都为 nil", func(t *testing.T) {
		err := CheckTypeCompatibility(nil, nil)
		assert.Error(t, err)
	})
}

// TestCheckTypeCompatibilityPointer 测试指针类型解引用
func TestCheckTypeCompatibilityPointer(t *testing.T) {
	t.Run("双指针兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[*int](), typeOf[*int]())
		assert.NoError(t, err)
	})

	t.Run("指针与非指针兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[*string](), typeOf[string]())
		assert.NoError(t, err)
	})

	t.Run("多层指针", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[**int](), typeOf[*int]())
		assert.NoError(t, err)
	})
}

// TestCheckTypeCompatibilityEmptyInterface 测试空接口
func TestCheckTypeCompatibilityEmptyInterface(t *testing.T) {
	t.Run("目标为空接口总是兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[int](), typeOf[interface{}]())
		assert.NoError(t, err)
	})

	t.Run("源为空接口总是兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[interface{}](), typeOf[int]())
		assert.NoError(t, err)
	})
}

// TestCheckTypeCompatibilityTimeToString 测试 time.Time 到 string 的特殊兼容
func TestCheckTypeCompatibilityTimeToString(t *testing.T) {
	err := CheckTypeCompatibility(reflect.TypeOf(time.Time{}), typeOf[string]())
	assert.NoError(t, err)
}

// TestCheckTypeCompatibilityKindMismatch 测试类型种类不匹配
func TestCheckTypeCompatibilityKindMismatch(t *testing.T) {
	err := CheckTypeCompatibility(typeOf[int](), typeOf[string]())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type mismatch")
}

// TestCheckTypeCompatibilityStruct 测试结构体兼容性
func TestCheckTypeCompatibilityStruct(t *testing.T) {
	type Src struct {
		Name  string
		Age   int
		inner string
	}
	type Dst struct {
		Name  string
		Age   int
		inner string
	}
	type DstMissing struct {
		Name string
	}
	type DstIncompatible struct {
		Name string
		Age  string // 类型不兼容
	}

	t.Run("结构体字段兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[Src](), typeOf[Dst]())
		assert.NoError(t, err)
	})

	t.Run("目标缺少字段也兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[Src](), typeOf[DstMissing]())
		assert.NoError(t, err)
	})

	t.Run("字段类型不兼容返回错误", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[Src](), typeOf[DstIncompatible]())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type mismatch")
	})

	t.Run("目标有源中不存在的字段也兼容", func(t *testing.T) {
		// 目标有 Extra 字段但源中没有，应跳过
		type DstWithExtraField struct {
			Name  string
			Age   int
			Extra string
		}
		err := CheckTypeCompatibility(typeOf[Src](), typeOf[DstWithExtraField]())
		assert.NoError(t, err)
	})
}

// TestCheckTypeCompatibilitySliceAndArray 测试切片和数组
func TestCheckTypeCompatibilitySliceAndArray(t *testing.T) {
	t.Run("切片元素兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[[]int](), typeOf[[]int]())
		assert.NoError(t, err)
	})

	t.Run("切片元素不兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[[]int](), typeOf[[]string]())
		assert.Error(t, err)
	})

	t.Run("数组元素兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[[3]int](), typeOf[[3]int]())
		assert.NoError(t, err)
	})

	t.Run("切片与数组兼容（元素类型相同）", func(t *testing.T) {
		// 注意：切片和数组的 Kind 不同，所以这里应返回错误
		err := CheckTypeCompatibility(typeOf[[]int](), typeOf[[3]int]())
		assert.Error(t, err)
	})
}

// TestCheckTypeCompatibilityMap 测试 map 类型
func TestCheckTypeCompatibilityMap(t *testing.T) {
	t.Run("map 兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[map[string]int](), typeOf[map[string]int]())
		assert.NoError(t, err)
	})

	t.Run("map key 不兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[map[string]int](), typeOf[map[int]int]())
		assert.Error(t, err)
	})

	t.Run("map value 不兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[map[string]int](), typeOf[map[string]string]())
		assert.Error(t, err)
	})
}

// TestCheckTypeCompatibilityUnsupportedTypes 测试不支持的类型
func TestCheckTypeCompatibilityUnsupportedTypes(t *testing.T) {
	t.Run("Func 类型不支持", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[func()](), typeOf[func()]())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported type")
	})

	t.Run("Chan 类型不支持", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[chan int](), typeOf[chan int]())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported type")
	})

	t.Run("UnsafePointer 不支持", func(t *testing.T) {
		// 使用 unsafe.Pointer 类型，其 Kind 为 reflect.UnsafePointer
		var p unsafe.Pointer
		unsafePtrType := reflect.TypeOf(p)
		err := CheckTypeCompatibility(unsafePtrType, unsafePtrType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported type")
	})
}

// TestCheckTypeCompatibilityDefault 测试默认分支
func TestCheckTypeCompatibilityDefault(t *testing.T) {
	t.Run("相同基础类型兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[int](), typeOf[int]())
		assert.NoError(t, err)
	})

	t.Run("不同基础类型不兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[int](), typeOf[int64]())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type mismatch")
	})

	t.Run("string 类型兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[string](), typeOf[string]())
		assert.NoError(t, err)
	})

	t.Run("bool 类型兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[bool](), typeOf[bool]())
		assert.NoError(t, err)
	})

	t.Run("同 Kind 不同命名类型不兼容", func(t *testing.T) {
		// MyInt 和 YourInt 的底层类型都是 int（Kind 相同），但类型不同
		err := CheckTypeCompatibility(typeOf[MyInt](), typeOf[YourInt]())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type mismatch")
	})

	t.Run("相同命名类型兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[MyInt](), typeOf[MyInt]())
		assert.NoError(t, err)
	})
}

// TestCheckTypeCompatibilityStructWithUnexported 测试结构体中未导出字段的处理
func TestCheckTypeCompatibilityStructWithUnexported(t *testing.T) {
	type WithUnexported struct {
		Exported   string
		unexported string
	}
	type DstWithUnexported struct {
		Exported   string
		unexported int // 类型不同但未导出应跳过
	}

	t.Run("未导出字段类型不匹配也应跳过", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[WithUnexported](), typeOf[DstWithUnexported]())
		assert.NoError(t, err)
	})
}

// TestCheckTypeCompatibilityNestedStruct 测试嵌套结构体
func TestCheckTypeCompatibilityNestedStruct(t *testing.T) {
	type Inner struct {
		Value int
	}
	type OuterSrc struct {
		Inner Inner
		Name  string
	}
	type OuterDst struct {
		Inner Inner
		Name  string
	}
	type OuterDstBad struct {
		Inner struct {
			Value string // 不兼容
		}
		Name string
	}

	t.Run("嵌套结构体兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[OuterSrc](), typeOf[OuterDst]())
		assert.NoError(t, err)
	})

	t.Run("嵌套结构体不兼容", func(t *testing.T) {
		err := CheckTypeCompatibility(typeOf[OuterSrc](), typeOf[OuterDstBad]())
		assert.Error(t, err)
	})
}

// TestCheckTypeCompatibilitySrcFieldUnexported 使用 reflect.StructOf 构造源结构体
// 其字段名大写但 PkgPath 非空，以覆盖 srcField.PkgPath != "" 分支
func TestCheckTypeCompatibilitySrcFieldUnexported(t *testing.T) {
	// 使用 reflect.StructOf 构造一个字段名大写但 PkgPath 非空的结构体类型
	srcType := reflect.StructOf([]reflect.StructField{
		{
			Name:    "Value",
			PkgPath: "fakepkg", // 设置 PkgPath 使字段被视为未导出
			Type:    reflect.TypeOf(int(0)),
		},
	})
	// 目标结构体有正常导出字段
	dstType := reflect.StructOf([]reflect.StructField{
		{
			Name: "Value",
			Type: reflect.TypeOf(int(0)),
		},
	})

	// 源字段 PkgPath 非空时应跳过该字段的兼容性检查
	err := CheckTypeCompatibility(srcType, dstType)
	assert.NoError(t, err)
}

// TestCheckTypeCompatibilityDstFieldUnexported 测试目标结构体字段未导出时的跳过逻辑
func TestCheckTypeCompatibilityDstFieldUnexported(t *testing.T) {
	type Src struct {
		hidden int // 未导出
	}
	type Dst struct {
		hidden string // 未导出，类型不同
	}

	// 目标字段未导出应跳过
	err := CheckTypeCompatibility(typeOf[Src](), typeOf[Dst]())
	assert.NoError(t, err)
}
