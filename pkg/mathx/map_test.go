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

// TestMergeLayeredKeyValues 测试多层级 key-value 合并
func TestMergeLayeredKeyValues(t *testing.T) {
	assert := assert.New(t)

	// 定义测试用的结构体
	type LocalizedText struct {
		Key   string
		Value string
	}

	type Config struct {
		Messages []LocalizedText
	}

	// 创建合并器（只需创建一次）
	merger := NewLayeredMerger[Config, LocalizedText]("Key", "Value")

	t.Run("基本三层合并", func(t *testing.T) {
		// 第一层：硬编码默认值
		hardcoded := &Config{
			Messages: []LocalizedText{
				{Key: "en", Value: "Hello"},
				{Key: "zh", Value: "你好"},
			},
		}

		// 第二层：owner 配置
		owner := &Config{
			Messages: []LocalizedText{
				{Key: "en", Value: "Hi"}, // 覆盖英文
				{Key: "fr", Value: "Bonjour"},
			},
		}

		// 第三层：agent 配置
		agent := &Config{
			Messages: []LocalizedText{
				{Key: "en", Value: "Hey"}, // 再次覆盖英文
				{Key: "es", Value: "Hola"},
			},
		}

		result := merger.Merge(
			[]*Config{hardcoded, owner, agent},
			func(c *Config) []LocalizedText { return c.Messages },
		)

		// 验证结果
		assert.Equal(4, len(result))
		assert.Equal("en", result[0].Key)
		assert.Equal("Hey", result[0].Value) // agent 优先级最高
		assert.Equal("zh", result[1].Key)
		assert.Equal("你好", result[1].Value) // 只在 hardcoded 中
		assert.Equal("fr", result[2].Key)
		assert.Equal("Bonjour", result[2].Value) // 只在 owner 中
		assert.Equal("es", result[3].Key)
		assert.Equal("Hola", result[3].Value) // 只在 agent 中
	})

	t.Run("跳过nil层级", func(t *testing.T) {
		layer1 := &Config{
			Messages: []LocalizedText{{Key: "en", Value: "Hello"}},
		}
		layer3 := &Config{
			Messages: []LocalizedText{{Key: "zh", Value: "你好"}},
		}

		result := merger.Merge(
			[]*Config{layer1, nil, layer3},
			func(c *Config) []LocalizedText { return c.Messages },
		)

		assert.Equal(2, len(result))
		assert.Equal("en", result[0].Key)
		assert.Equal("zh", result[1].Key)
	})

	t.Run("跳过空值", func(t *testing.T) {
		layer1 := &Config{
			Messages: []LocalizedText{
				{Key: "en", Value: "Hello"},
				{Key: "zh", Value: "你好"},
			},
		}
		layer2 := &Config{
			Messages: []LocalizedText{
				{Key: "en", Value: ""}, // 空值，不应覆盖
				{Key: "fr", Value: "Bonjour"},
			},
		}

		result := merger.Merge(
			[]*Config{layer1, layer2},
			func(c *Config) []LocalizedText { return c.Messages },
		)

		assert.Equal(3, len(result))
		assert.Equal("Hello", result[0].Value) // 空值未覆盖
		assert.Equal("你好", result[1].Value)
		assert.Equal("Bonjour", result[2].Value)
	})

	t.Run("保持key首次出现顺序", func(t *testing.T) {
		layer1 := &Config{
			Messages: []LocalizedText{
				{Key: "a", Value: "1"},
				{Key: "b", Value: "2"},
				{Key: "c", Value: "3"},
			},
		}
		layer2 := &Config{
			Messages: []LocalizedText{
				{Key: "d", Value: "4"},
				{Key: "b", Value: "22"}, // 覆盖，但顺序保持在原位
			},
		}

		result := merger.Merge(
			[]*Config{layer1, layer2},
			func(c *Config) []LocalizedText { return c.Messages },
		)

		assert.Equal(4, len(result))
		assert.Equal("a", result[0].Key)
		assert.Equal("b", result[1].Key)
		assert.Equal("22", result[1].Value) // 值被覆盖
		assert.Equal("c", result[2].Key)
		assert.Equal("d", result[3].Key) // 新 key 在最后
	})

	t.Run("空layers数组", func(t *testing.T) {
		result := merger.Merge(
			[]*Config{},
			func(c *Config) []LocalizedText { return c.Messages },
		)

		assert.Equal(0, len(result))
	})

	t.Run("全nil layers", func(t *testing.T) {
		result := merger.Merge(
			[]*Config{nil, nil, nil},
			func(c *Config) []LocalizedText { return c.Messages },
		)

		assert.Equal(0, len(result))
	})

	t.Run("支持更多层级", func(t *testing.T) {
		layers := make([]*Config, 5)
		for i := 0; i < 5; i++ {
			layers[i] = &Config{
				Messages: []LocalizedText{
					{Key: "key", Value: string(rune('A' + i))}, // A, B, C, D, E
				},
			}
		}

		result := merger.Merge(
			layers,
			func(c *Config) []LocalizedText { return c.Messages },
		)

		assert.Equal(1, len(result))
		assert.Equal("E", result[0].Value) // 最后一层优先级最高
	})

	t.Run("多语言场景真实测试", func(t *testing.T) {
		// 模拟实际的多语言配置场景
		systemDefault := &Config{
			Messages: []LocalizedText{
				{Key: "en", Value: "Welcome"},
				{Key: "zh", Value: "欢迎"},
				{Key: "es", Value: "Bienvenido"},
			},
		}

		companyConfig := &Config{
			Messages: []LocalizedText{
				{Key: "en", Value: "Welcome to our company"},
				{Key: "vi", Value: "Chào mừng"}, // 新增越南语
			},
		}

		agentCustom := &Config{
			Messages: []LocalizedText{
				{Key: "en", Value: "Hi, I'm your agent!"}, // 个性化英文
			},
		}

		result := merger.Merge(
			[]*Config{systemDefault, companyConfig, agentCustom},
			func(c *Config) []LocalizedText { return c.Messages },
		)

		// 验证合并结果
		assert.Equal(4, len(result))

		msgMap := make(map[string]string)
		for _, msg := range result {
			msgMap[msg.Key] = msg.Value
		}

		assert.Equal("Hi, I'm your agent!", msgMap["en"]) // agent 覆盖
		assert.Equal("欢迎", msgMap["zh"])                  // 系统默认
		assert.Equal("Bienvenido", msgMap["es"])          // 系统默认
		assert.Equal("Chào mừng", msgMap["vi"])           // 公司配置
	})

	t.Run("直接调用MergeLayeredKeyValues兼容性测试", func(t *testing.T) {
		// 验证直接调用底层函数也能正常工作
		layer1 := &Config{
			Messages: []LocalizedText{{Key: "en", Value: "Hello"}},
		}

		result := MergeLayeredKeyValues(
			[]*Config{layer1},
			func(c *Config) []LocalizedText { return c.Messages },
			"Key",
			"Value",
		)

		assert.Equal(1, len(result))
		assert.Equal("Hello", result[0].Value)
	})
}
