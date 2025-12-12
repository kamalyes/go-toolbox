/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-09 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\convert\yaml.go
 * @Description: YAML/JSON 转换工具
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// ConvertError 用于包装转换过程中产生的错误
type ConvertError struct {
	Op  string // 操作名称
	Err error  // 原始错误
}

// Error 实现 error 接口，返回错误信息
func (e *ConvertError) Error() string {
	return fmt.Sprintf("convert error during %s: %v", e.Op, e.Err)
}

// unmarshalYAML 是一个辅助函数，用于将 YAML 数据反序列化到指定的输出结构体
func unmarshalYAML(yamlData []byte, out interface{}) error {
	if err := yaml.Unmarshal(yamlData, out); err != nil {
		return &ConvertError{"unmarshaling YAML", err}
	}
	return nil
}

// marshalYAML 是一个辅助函数，用于将数据序列化为 YAML 格式
func marshalYAML(data interface{}) ([]byte, error) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return nil, &ConvertError{"marshalling to YAML", err}
	}
	return yamlData, nil
}

// YAMLToJSON 将 YAML 字节数组转换为 JSON 字节数组
func YAMLToJSON(yamlData []byte) ([]byte, error) {
	var data interface{}
	if err := unmarshalYAML(yamlData, &data); err != nil {
		return nil, err
	}

	// 递归转换为 JSON 兼容格式
	data = convertYAMLToJSONCompatible(data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, &ConvertError{"marshalling to JSON", err}
	}

	return jsonData, nil
}

// JSONToYAML 将 JSON 字节数组转换为 YAML 字节数组
func JSONToYAML(jsonData []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, &ConvertError{"unmarshaling JSON", err}
	}

	return marshalYAML(data)
}

// YAMLStringToJSON 将 YAML 字符串转换为 JSON 字符串
func YAMLStringToJSON(yamlStr string) (string, error) {
	jsonData, err := YAMLToJSON([]byte(yamlStr))
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// JSONStringToYAML 将 JSON 字符串转换为 YAML 字符串
func JSONStringToYAML(jsonStr string) (string, error) {
	yamlData, err := JSONToYAML([]byte(jsonStr))
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

// YAMLToInterface 将 YAML 字节数组解析为 interface{}
func YAMLToInterface(yamlData []byte) (interface{}, error) {
	var data interface{}
	if err := unmarshalYAML(yamlData, &data); err != nil {
		return nil, err
	}
	return convertYAMLToJSONCompatible(data), nil
}

// YAMLToMap 将 YAML 字节数组解析为 map[string]interface{}
func YAMLToMap(yamlData []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := unmarshalYAML(yamlData, &data); err != nil {
		return nil, err
	}

	// 确保所有嵌套的 map 键都是字符串
	return convertYAMLMapToStringMap(data), nil
}

// InterfaceToYAML 将 interface{} 转换为 YAML 字节数组
func InterfaceToYAML(data interface{}) ([]byte, error) {
	return marshalYAML(data)
}

// MapToYAML 将 map[string]interface{} 转换为 YAML 字节数组
func MapToYAML(data map[string]interface{}) ([]byte, error) {
	return InterfaceToYAML(data)
}

// convertYAMLToJSONCompatible 递归转换 YAML 数据为 JSON 兼容格式
// YAML 允许任意类型作为 map 的键，但 JSON 只允许字符串
func convertYAMLToJSONCompatible(data interface{}) interface{} {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, value := range v {
			strKey := fmt.Sprintf("%v", key) // 将键转换为字符串
			result[strKey] = convertYAMLToJSONCompatible(value)
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, value := range v {
			result[key] = convertYAMLToJSONCompatible(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = convertYAMLToJSONCompatible(value)
		}
		return result
	default:
		return v // 原样返回其他类型
	}
}

// convertYAMLMapToStringMap 将 map[string]interface{} 中的嵌套 map 键转换为字符串
func convertYAMLMapToStringMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(data))
	for key, value := range data {
		result[key] = convertYAMLToJSONCompatible(value)
	}
	return result
}

// UnmarshalYAML 泛型 YAML 反序列化
func UnmarshalYAML[T any](yamlData []byte) (*T, error) {
	var result T
	if err := unmarshalYAML(yamlData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MarshalYAML 泛型 YAML 序列化
func MarshalYAML[T any](data T) ([]byte, error) {
	return marshalYAML(data)
}

// UnmarshalJSON 泛型 JSON 反序列化（复用现有 json 包功能）
func UnmarshalJSON[T any](jsonData []byte) (*T, error) {
	var result T
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, &ConvertError{"unmarshaling JSON", err}
	}
	return &result, nil
}

// MarshalJSON 泛型 JSON 序列化（复用现有 json 包功能）
func MarshalJSON[T any](data T) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, &ConvertError{"marshalling to JSON", err}
	}
	return jsonData, nil
}
