/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-09 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\mathx\map_test.go
 * @Description: Map 操作工具测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package mathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepMergeMap(t *testing.T) {
	assert := assert.New(t)

	t.Run("简单合并", func(t *testing.T) {
		target := map[string]interface{}{"a": 1, "b": 2}
		source := map[string]interface{}{"b": 3, "c": 4}

		result, err := DeepMergeMap(target, source, nil)
		assert.NoError(err)
		assert.Equal(1, result["a"])
		assert.Equal(3, result["b"])
		assert.Equal(4, result["c"])
	})

	t.Run("嵌套合并", func(t *testing.T) {
		target := map[string]interface{}{
			"user": map[string]interface{}{"name": "Alice", "age": 30},
		}
		source := map[string]interface{}{
			"user": map[string]interface{}{"age": 31, "city": "Beijing"},
		}

		result, err := DeepMergeMap(target, source, nil)
		assert.NoError(err)

		user := result["user"].(map[string]interface{})
		assert.Equal("Alice", user["name"])
		assert.Equal(31, user["age"])
		assert.Equal("Beijing", user["city"])
	})

	t.Run("nil target", func(t *testing.T) {
		source := map[string]interface{}{"a": 1}
		result, err := DeepMergeMap(nil, source, nil)
		assert.NoError(err)
		assert.Equal(1, result["a"])
	})

	t.Run("nil source", func(t *testing.T) {
		target := map[string]interface{}{"a": 1}
		result, err := DeepMergeMap(target, nil, nil)
		assert.NoError(err)
		assert.Equal(1, result["a"])
	})

	t.Run("source with nil value", func(t *testing.T) {
		target := map[string]interface{}{"a": 1}
		source := map[string]interface{}{"b": nil}
		result, err := DeepMergeMap(target, source, nil)
		assert.NoError(err)
		assert.Equal(1, len(result))    // nil值被跳过
		assert.NotContains(result, "b") // b键不应该存在
	})

	t.Run("保持现有策略", func(t *testing.T) {
		target := map[string]interface{}{"a": 1}
		source := map[string]interface{}{"a": 2}
		options := &MapMergeOptions{Strategy: MapMergeStrategyKeepExisting}

		result, err := DeepMergeMap(target, source, options)
		assert.NoError(err)
		assert.Equal(1, result["a"])
	})

	t.Run("冲突报错策略", func(t *testing.T) {
		target := map[string]interface{}{"a": 1}
		source := map[string]interface{}{"a": "conflict"}
		options := &MapMergeOptions{
			Strategy:   MapMergeStrategyError,
			TypeStrict: true,
		}

		_, err := DeepMergeMap(target, source, options)
		assert.Error(err)
	})

	t.Run("自定义冲突处理", func(t *testing.T) {
		target := map[string]interface{}{"a": 1}
		source := map[string]interface{}{"a": 2}
		options := &MapMergeOptions{
			OnConflict: func(key string, target, source interface{}) interface{} {
				return 999 // 自定义值
			},
		}

		result, err := DeepMergeMap(target, source, options)
		assert.NoError(err)
		assert.Equal(999, result["a"])
	})

	t.Run("切片合并-覆盖", func(t *testing.T) {
		target := map[string]interface{}{"tags": []interface{}{"a", "b"}}
		source := map[string]interface{}{"tags": []interface{}{"c"}}
		options := &MapMergeOptions{Strategy: MapMergeStrategyOverwrite}

		result, err := DeepMergeMap(target, source, options)
		assert.NoError(err)
		tags := result["tags"].([]interface{})
		assert.Equal(1, len(tags))
	})

	t.Run("切片合并-保持", func(t *testing.T) {
		target := map[string]interface{}{"tags": []interface{}{"a", "b"}}
		source := map[string]interface{}{"tags": []interface{}{"c"}}
		options := &MapMergeOptions{Strategy: MapMergeStrategyKeepExisting}

		result, err := DeepMergeMap(target, source, options)
		assert.NoError(err)
		tags := result["tags"].([]interface{})
		assert.Equal(2, len(tags))
	})

	t.Run("切片合并-默认（覆盖）", func(t *testing.T) {
		target := map[string]interface{}{"tags": []interface{}{"a"}}
		source := map[string]interface{}{"tags": []interface{}{"b"}}

		result, err := DeepMergeMap(target, source, nil)
		assert.NoError(err)
		tags := result["tags"].([]interface{})
		assert.Equal(1, len(tags)) // 默认策略是覆盖，所以结果是source的值
		assert.Equal("b", tags[0])
	})

	t.Run("超过最大深度", func(t *testing.T) {
		target := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{},
			},
		}
		source := map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": 1,
				},
			},
		}
		options := &MapMergeOptions{MaxDepth: 1}

		_, err := DeepMergeMap(target, source, options)
		assert.Error(err)
		assert.Contains(err.Error(), "exceeded maximum merge depth")
	})
}

func TestShallowMergeMap(t *testing.T) {
	assert := assert.New(t)

	m1 := map[string]int{"a": 1}
	m2 := map[string]int{"b": 2}
	result := ShallowMergeMap(m1, m2)

	assert.Equal(1, result["a"])
	assert.Equal(2, result["b"])
}

func TestConvertMapKeysToString(t *testing.T) {
	assert := assert.New(t)

	data := map[interface{}]interface{}{
		"name": "Alice",
		123:    "number",
	}

	result := ConvertMapKeysToString(data)
	resultMap := result.(map[string]interface{})

	assert.Equal("Alice", resultMap["name"])
	assert.Equal("number", resultMap["123"])
}

func TestFlattenMap(t *testing.T) {
	assert := assert.New(t)

	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{"c": 1},
		},
	}

	result := FlattenMap(data, ".")
	assert.Equal(1, result["a.b.c"])
}

func TestUnflattenMap(t *testing.T) {
	assert := assert.New(t)

	data := map[string]interface{}{"a.b.c": 1}
	result := UnflattenMap(data, ".")

	a := result["a"].(map[string]interface{})
	b := a["b"].(map[string]interface{})

	assert.Equal(1, b["c"])
}

func TestFilterMap(t *testing.T) {
	assert := assert.New(t)

	data := map[string]int{"a": 1, "b": 2, "c": 3}
	result := FilterMap(data, func(k string, v int) bool {
		return v%2 == 0
	})

	assert.Len(result, 1)
	assert.Equal(2, result["b"])
}

func TestTransformMapValues(t *testing.T) {
	assert := assert.New(t)

	data := map[string]int{"a": 1, "b": 2}
	result := TransformMapValues(data, func(v int) string {
		return IF(v%2 == 0, "even", "odd")
	})

	assert.Equal("odd", result["a"])
	assert.Equal("even", result["b"])
}

func TestCloneMap(t *testing.T) {
	assert := assert.New(t)

	original := map[string]int{"a": 1}
	cloned := CloneMap(original)

	cloned["a"] = 2
	assert.Equal(1, original["a"])
}

func TestGetNestedMapValue(t *testing.T) {
	assert := assert.New(t)

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Bob",
		},
	}

	name, ok := GetNestedMapValue[string](data, "user", "name")
	assert.True(ok)
	assert.Equal("Bob", name)
}

func TestSetNestedMapValue(t *testing.T) {
	assert := assert.New(t)

	data := make(map[string]interface{})
	SetNestedMapValue(data, "test", "a", "b", "c")

	a := data["a"].(map[string]interface{})
	b := a["b"].(map[string]interface{})

	assert.Equal("test", b["c"])
}
