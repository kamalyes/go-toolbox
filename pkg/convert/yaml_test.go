/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-09 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\convert\yaml_test.go
 * @Description: YAML/JSON 转换工具
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLToJSON(t *testing.T) {
	yamlInput := `
name: John Doe
age: 30
skills:
  - Go
  - Python
`
	expectedJSON := `{"age":30,"name":"John Doe","skills":["Go","Python"]}`

	jsonOutput, err := YAMLToJSON([]byte(yamlInput))
	assert.NoError(t, err, "YAML 转换为 JSON 时出错")
	assert.JSONEq(t, expectedJSON, string(jsonOutput), "转换后的 JSON 不匹配")
}

func TestJSONToYAML(t *testing.T) {
	jsonInput := `{"name":"John Doe","age":30,"skills":["Go","Python"]}`

	yamlOutput, err := JSONToYAML([]byte(jsonInput))
	assert.NoError(t, err, "JSON 转换为 YAML 时出错")

	// YAML字段顺序不固定，检查关键内容
	yamlStr := string(yamlOutput)
	assert.Contains(t, yamlStr, "name: John Doe", "YAML应包含name字段")
	assert.Contains(t, yamlStr, "age: 30", "YAML应包含age字段")
	assert.Contains(t, yamlStr, "- Go", "YAML应包含Go技能")
	assert.Contains(t, yamlStr, "- Python", "YAML应包含Python技能")
}

func TestYAMLStringToJSON(t *testing.T) {
	yamlInput := `
name: John Doe
age: 30
skills:
  - Go
  - Python
`
	expectedJSON := `{"age":30,"name":"John Doe","skills":["Go","Python"]}`

	jsonOutput, err := YAMLStringToJSON(yamlInput)
	assert.NoError(t, err, "YAML 字符串转换为 JSON 时出错")
	assert.JSONEq(t, expectedJSON, jsonOutput, "转换后的 JSON 不匹配")
}

func TestJSONStringToYAML(t *testing.T) {
	jsonInput := `{"name":"John Doe","age":30,"skills":["Go","Python"]}`

	yamlOutput, err := JSONStringToYAML(jsonInput)
	assert.NoError(t, err, "JSON 字符串转换为 YAML 时出错")

	// YAML字段顺序不固定，检查关键内容
	assert.Contains(t, yamlOutput, "name: John Doe", "YAML应包含name字段")
	assert.Contains(t, yamlOutput, "age: 30", "YAML应包含age字段")
	assert.Contains(t, yamlOutput, "- Go", "YAML应包含Go技能")
	assert.Contains(t, yamlOutput, "- Python", "YAML应包含Python技能")
}

func TestYAMLToInterface(t *testing.T) {
	yamlInput := `
name: John Doe
age: 30
skills:
  - Go
  - Python
`

	result, err := YAMLToInterface([]byte(yamlInput))
	assert.NoError(t, err, "YAML 转换为 interface{} 时出错")

	// 进行类型断言并验证结果
	data, ok := result.(map[string]interface{})
	assert.True(t, ok, "转换结果应为 map[string]interface{}")
	assert.Equal(t, "John Doe", data["name"], "姓名不匹配")
	assert.Equal(t, 30, data["age"], "年龄不匹配")
	assert.ElementsMatch(t, []interface{}{"Go", "Python"}, data["skills"], "技能不匹配")
}

func TestYAMLToMap(t *testing.T) {
	yamlInput := `
name: John Doe
age: 30
skills:
  - Go
  - Python
`

	result, err := YAMLToMap([]byte(yamlInput))
	assert.NoError(t, err, "YAML 转换为 map 时出错")
	assert.Equal(t, "John Doe", result["name"], "姓名不匹配")
	assert.Equal(t, 30, result["age"], "年龄不匹配")
	assert.ElementsMatch(t, []interface{}{"Go", "Python"}, result["skills"], "技能不匹配")
}

// 测试 InterfaceToYAML
func TestInterfaceToYAML(t *testing.T) {
	data := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}

	yamlOutput, err := InterfaceToYAML(data)
	assert.NoError(t, err, "interface{} 转换为 YAML 时出错")
	assert.Contains(t, string(yamlOutput), "name: Alice")
	assert.Contains(t, string(yamlOutput), "age: 25")
}

// 测试 MapToYAML
func TestMapToYAML(t *testing.T) {
	data := map[string]interface{}{
		"service": "api",
		"port":    8080,
	}

	yamlOutput, err := MapToYAML(data)
	assert.NoError(t, err, "map 转换为 YAML 时出错")
	assert.Contains(t, string(yamlOutput), "service: api")
	assert.Contains(t, string(yamlOutput), "port: 8080")
}

// 测试 UnmarshalYAML 泛型
func TestUnmarshalYAML(t *testing.T) {
	type Config struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
	}

	yamlInput := []byte(`
name: test-service
port: 9000
`)

	result, err := UnmarshalYAML[Config](yamlInput)
	assert.NoError(t, err, "YAML 反序列化失败")
	assert.Equal(t, "test-service", result.Name)
	assert.Equal(t, 9000, result.Port)
}

// 测试 MarshalYAML 泛型
func TestMarshalYAML(t *testing.T) {
	type User struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	}

	user := User{
		Name:  "Bob",
		Email: "bob@example.com",
	}

	yamlOutput, err := MarshalYAML(user)
	assert.NoError(t, err, "YAML 序列化失败")
	assert.Contains(t, string(yamlOutput), "name: Bob")
	assert.Contains(t, string(yamlOutput), "email: bob@example.com")
}

// 测试 UnmarshalJSON 泛型
func TestUnmarshalJSON(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonInput := []byte(`{"name":"Charlie","age":35}`)

	result, err := UnmarshalJSON[Person](jsonInput)
	assert.NoError(t, err, "JSON 反序列化失败")
	assert.Equal(t, "Charlie", result.Name)
	assert.Equal(t, 35, result.Age)
}

// 测试 MarshalJSON 泛型
func TestMarshalJSON(t *testing.T) {
	type Product struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	product := Product{
		ID:   1,
		Name: "Book",
	}

	jsonOutput, err := MarshalJSON(product)
	assert.NoError(t, err, "JSON 序列化失败")
	assert.Contains(t, string(jsonOutput), `"id":1`)
	assert.Contains(t, string(jsonOutput), `"name":"Book"`)
}

// 测试嵌套 map 的键转换
func TestConvertYAMLToJSONCompatibleNestedMap(t *testing.T) {
	data := map[interface{}]interface{}{
		"user": map[interface{}]interface{}{
			"name": "Test",
			123:    "number_key",
		},
	}

	result := convertYAMLToJSONCompatible(data)
	resultMap := result.(map[string]interface{})
	user := resultMap["user"].(map[string]interface{})

	assert.Equal(t, "Test", user["name"])
	assert.Equal(t, "number_key", user["123"])
}

// 测试 slice 中的 map 键转换
func TestConvertYAMLToJSONCompatibleSliceWithMaps(t *testing.T) {
	data := []interface{}{
		map[interface{}]interface{}{
			"id": 1,
		},
		map[interface{}]interface{}{
			456: "number_key",
		},
	}

	result := convertYAMLToJSONCompatible(data)
	resultSlice := result.([]interface{})

	first := resultSlice[0].(map[string]interface{})
	second := resultSlice[1].(map[string]interface{})

	assert.Equal(t, 1, first["id"])
	assert.Equal(t, "number_key", second["456"])
}

// 测试 map[string]interface{} 递归转换
func TestConvertYAMLToJSONCompatibleStringMap(t *testing.T) {
	data := map[string]interface{}{
		"nested": map[string]interface{}{
			"key": "value",
		},
	}

	result := convertYAMLToJSONCompatible(data)
	resultMap := result.(map[string]interface{})
	nested := resultMap["nested"].(map[string]interface{})

	assert.Equal(t, "value", nested["key"])
}

// 测试错误处理
func TestYAMLToJSONInvalidYAML(t *testing.T) {
	invalidYAML := []byte(`invalid: [yaml`)

	_, err := YAMLToJSON(invalidYAML)
	assert.Error(t, err, "应该返回错误")
}

func TestJSONToYAMLInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{invalid json}`)

	_, err := JSONToYAML(invalidJSON)
	assert.Error(t, err, "应该返回错误")
}

// 测试往返转换
func TestYAMLJSONRoundTrip(t *testing.T) {
	originalYAML := []byte(`
name: Test
value: 123
nested:
  key: value
`)

	// YAML -> JSON -> YAML
	jsonData, err := YAMLToJSON(originalYAML)
	assert.NoError(t, err)

	yamlData, err := JSONToYAML(jsonData)
	assert.NoError(t, err)

	// 验证数据完整性
	result, err := YAMLToMap(yamlData)
	assert.NoError(t, err)

	assert.Equal(t, "Test", result["name"])
	assert.Equal(t, 123, result["value"])

	nested := result["nested"].(map[string]interface{})
	assert.Equal(t, "value", nested["key"])
}
