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
	"sort"
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

func TestMapClone(t *testing.T) {
	assert := assert.New(t)

	// 正常 map
	original := map[string]int{"a": 1, "b": 2}
	cloned := MapClone(original)
	assert.Equal(original, cloned)

	// 修改克隆不影响原 map
	cloned["a"] = 999
	assert.Equal(1, original["a"])

	// 新增键
	cloned["c"] = 3
	_, exists := original["c"]
	assert.False(exists)

	// nil map
	var nilMap map[string]int
	assert.Nil(MapClone(nilMap))

	// 空 map
	emptyClone := MapClone(map[string]int{})
	assert.NotNil(emptyClone)
	assert.Len(emptyClone, 0)
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

// ============================================================================
// 键/值提取
// ============================================================================

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := MapKeys(m)
	sort.Strings(keys) // map 遍历顺序不保证
	assert.Equal(t, []string{"a", "b", "c"}, keys)

	// 多个 map
	m2 := map[string]int{"d": 4}
	keys = MapKeys(m, m2)
	sort.Strings(keys)
	assert.Equal(t, []string{"a", "b", "c", "d"}, keys)

	// 空 map
	assert.Equal(t, []string{}, MapKeys(map[string]int{}))
}

func TestMapUniqKeys(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4} // b 重复
	keys := MapUniqKeys(m1, m2)
	sort.Strings(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	values := MapValues(m)
	sort.Ints(values)
	assert.Equal(t, []int{1, 2, 3}, values)
}

func TestMapUniqValues(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"c": 1, "d": 3} // 1 重复
	values := MapUniqValues(m1, m2)
	sort.Ints(values)
	assert.Equal(t, []int{1, 2, 3}, values)
}

// ============================================================================
// 键检查 / 值获取
// ============================================================================

func TestMapHasKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	assert.True(t, MapHasKey(m, "a"))
	assert.False(t, MapHasKey(m, "c"))
}

func TestMapValueOr(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	assert.Equal(t, 1, MapValueOr(m, "a", 0))
	assert.Equal(t, 99, MapValueOr(m, "c", 99))
}

// ============================================================================
// 选择 / 移除（PickBy / OmitBy）
// ============================================================================

func TestMapPickBy(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := MapPickBy(m, func(key string, value int) bool {
		return value%2 == 0
	})
	assert.Equal(t, map[string]int{"b": 2, "d": 4}, result)
}

func TestMapOmitBy(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := MapOmitBy(m, func(key string, value int) bool {
		return value%2 == 0
	})
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, result)
}

func TestMapPickByKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := MapPickByKeys(m, []string{"a", "c", "x"}) // x 不存在
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, result)
}

func TestMapOmitByKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := MapOmitByKeys(m, []string{"a", "c", "x"})
	assert.Equal(t, map[string]int{"b": 2, "d": 4}, result)
}

func TestMapPickByValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := MapPickByValues(m, []int{2, 4})
	assert.Equal(t, map[string]int{"b": 2, "d": 4}, result)
}

func TestMapOmitByValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := MapOmitByValues(m, []int{2, 4})
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, result)
}

// ============================================================================
// map <-> 键值对切片互转
// ============================================================================

func TestMapEntriesAndFromEntries(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// map -> entries
	entries := MapEntries(m)
	assert.Len(t, entries, 3)

	// entries -> map
	m2 := MapFromEntries(entries)
	assert.Equal(t, m, m2)
}

func TestMapFromEntriesOverwrite(t *testing.T) {
	// 同 key 后者覆盖前者
	entries := []MapEntry[string, int]{
		{"a", 1},
		{"a", 2},
		{"b", 3},
	}
	result := MapFromEntries(entries)
	assert.Equal(t, 2, result["a"]) // 后者覆盖
	assert.Equal(t, 3, result["b"])
}

// ============================================================================
// 反转 / 合并
// ============================================================================

func TestMapInvert(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	inverted := MapInvert(m)
	assert.Equal(t, "a", inverted[1])
	assert.Equal(t, "b", inverted[2])
	assert.Equal(t, "c", inverted[3])
}

func TestMapInvertDuplicate(t *testing.T) {
	// 重复值，后者覆盖前者
	m := map[string]int{"a": 1, "b": 1}
	inverted := MapInvert(m)
	assert.Len(t, inverted, 1)
	assert.True(t, inverted[1] == "a" || inverted[1] == "b")
}

func TestMapAssign(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4} // b 覆盖
	m3 := map[string]int{"d": 5}

	result := MapAssign(m1, m2, m3)
	assert.Equal(t, 1, result["a"])
	assert.Equal(t, 3, result["b"]) // m2 覆盖 m1
	assert.Equal(t, 4, result["c"])
	assert.Equal(t, 5, result["d"])
}

// ============================================================================
// map 转切片
// ============================================================================

func TestMapToSlice(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := MapToSlice(m, func(key string, value int) string {
		return key + ":" + string(rune('0'+value))
	})
	assert.Len(t, result, 3)
	// map 遍历顺序不保证，用 contains 检查
	assert.Contains(t, result, "a:1")
	assert.Contains(t, result, "b:2")
	assert.Contains(t, result, "c:3")
}

func TestMapFilterToSlice(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := MapFilterToSlice(m, func(key string, value int) (string, bool) {
		if value%2 == 0 {
			return key, true
		}
		return "", false
	})
	assert.Len(t, result, 2)
	assert.Contains(t, result, "b")
	assert.Contains(t, result, "d")
}

func TestMapFilterKeys(t *testing.T) {
	m := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	keys := MapFilterKeys(m, func(key string, value int) bool {
		return len(key) > 1
	})
	sort.Strings(keys)
	assert.Equal(t, []string{"bb", "ccc"}, keys)
}

func TestMapFilterValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	values := MapFilterValues(m, func(key string, value int) bool {
		return value > 2
	})
	sort.Ints(values)
	assert.Equal(t, []int{3, 4}, values)
}

// ============================================================================
// map 键/值/键值对转换
// ============================================================================

func TestMapMapEntries(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := MapMapEntries(m, func(key string, value int) (int, string) {
		return value, key
	})
	assert.Equal(t, "a", result[1])
	assert.Equal(t, "b", result[2])
}

func TestMapMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := MapMapKeys(m, func(value int, key string) string {
		return key + key
	})
	assert.Equal(t, 1, result["aa"])
	assert.Equal(t, 2, result["bb"])
}

func TestMapMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := MapMapValues(m, func(value int, key string) string {
		return key + ":" + string(rune('0'+value))
	})
	assert.Equal(t, "a:1", result["a"])
	assert.Equal(t, "b:2", result["b"])
}

// ============================================================================
// 边界情况：空 map / nil map
// ============================================================================

func TestMapExtEmptyMap(t *testing.T) {
	empty := map[string]int{}

	assert.Equal(t, []string{}, MapKeys(empty))
	assert.Equal(t, []int{}, MapValues(empty))
	assert.False(t, MapHasKey(empty, "x"))
	assert.Equal(t, 99, MapValueOr(empty, "x", 99))
	assert.Equal(t, map[string]int{}, MapPickBy(empty, func(k string, v int) bool { return true }))
	assert.Equal(t, map[string]int{}, MapOmitBy(empty, func(k string, v int) bool { return true }))
	assert.Equal(t, []MapEntry[string, int]{}, MapEntries(empty))
	assert.Equal(t, map[string]int{}, MapFromEntries([]MapEntry[string, int]{}))
}

func TestMapExtNilMap(t *testing.T) {
	var nilMap map[string]int

	// nil map 应安全处理（不 panic）
	assert.Equal(t, []string{}, MapKeys(nilMap))
	assert.Equal(t, []int{}, MapValues(nilMap))
	assert.False(t, MapHasKey(nilMap, "x"))
	assert.Equal(t, 99, MapValueOr(nilMap, "x", 99))
	assert.Equal(t, map[string]int{}, MapPickBy(nilMap, func(k string, v int) bool { return true }))
	assert.Equal(t, []MapEntry[string, int]{}, MapEntries(nilMap))
}
