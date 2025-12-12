/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-09 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-09 01:00:00
 * @FilePath: \go-toolbox\pkg\mathx\map.go
 * @Description: Map 操作工具函数 - 深度合并、转换、扁平化等
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package mathx

import (
	"fmt"
	"reflect"
	"strconv"
)

// MapMergeStrategy Map合并策略
type MapMergeStrategy int

const (
	// MapMergeStrategyOverwrite 覆盖策略：源覆盖目标
	MapMergeStrategyOverwrite MapMergeStrategy = iota
	// MapMergeStrategyKeepExisting 保持现有：保留目标值，忽略源
	MapMergeStrategyKeepExisting
	// MapMergeStrategyError 冲突报错：发现冲突时返回错误
	MapMergeStrategyError
)

// MapMergeOptions Map合并选项
type MapMergeOptions struct {
	Strategy     MapMergeStrategy                                         // 合并策略
	MaxDepth     int                                                      // 最大递归深度，0表示不限制
	currentDepth int                                                      // 当前递归深度（内部使用）
	TypeStrict   bool                                                     // 是否严格类型检查
	OnConflict   func(key string, target, source interface{}) interface{} // 冲突处理回调
}

// DeepMergeMap 深度合并两个 map[string]interface{}
// target: 目标map（会被修改）
// source: 源map
// options: 合并选项，nil则使用默认选项
// 返回合并后的 map 和可能的错误
func DeepMergeMap(target, source map[string]interface{}, options *MapMergeOptions) (map[string]interface{}, error) {
	if target == nil {
		target = make(map[string]interface{})
	}
	if source == nil {
		return target, nil
	}
	if options == nil {
		options = &MapMergeOptions{
			Strategy: MapMergeStrategyOverwrite,
			MaxDepth: 100,
		}
	}

	// 检查递归深度
	if options.MaxDepth > 0 && options.currentDepth >= options.MaxDepth {
		return nil, fmt.Errorf("exceeded maximum merge depth of %d", options.MaxDepth)
	}

	for key, srcValue := range source {
		if srcValue == nil {
			continue // 跳过 nil 值
		}

		targetValue, exists := target[key]

		// 如果目标中不存在该键，直接设置
		if !exists {
			target[key] = srcValue
			continue
		}

		// 处理冲突
		merged, err := mergeMapValues(key, targetValue, srcValue, options)
		if err != nil {
			return nil, err
		}
		target[key] = merged
	}

	return target, nil
}

// mergeMapValues 合并两个值
func mergeMapValues(key string, targetValue, sourceValue interface{}, options *MapMergeOptions) (interface{}, error) {
	// 使用自定义冲突处理器
	if options.OnConflict != nil {
		return options.OnConflict(key, targetValue, sourceValue), nil
	}

	// 类型检查
	targetType := reflect.TypeOf(targetValue)
	sourceType := reflect.TypeOf(sourceValue)

	if options.TypeStrict && targetType != sourceType {
		if options.Strategy == MapMergeStrategyError {
			return nil, fmt.Errorf("type mismatch for key '%s': target is %v, source is %v", key, targetType, sourceType)
		}
	}

	// 如果两个都是 map[string]interface{}，递归合并
	targetMap, targetIsMap := targetValue.(map[string]interface{})
	sourceMap, sourceIsMap := sourceValue.(map[string]interface{})

	if targetIsMap && sourceIsMap {
		// 创建新的选项，增加递归深度
		newOptions := *options
		newOptions.currentDepth++
		return DeepMergeMap(targetMap, sourceMap, &newOptions)
	}

	// 如果两个都是切片，根据策略处理
	targetSlice, targetIsSlice := interfaceToSlice(targetValue)
	sourceSlice, sourceIsSlice := interfaceToSlice(sourceValue)

	if targetIsSlice && sourceIsSlice {
		return mergeMapSlices(targetSlice, sourceSlice, options)
	}

	// 其他情况根据策略处理
	switch options.Strategy {
	case MapMergeStrategyOverwrite:
		return sourceValue, nil
	case MapMergeStrategyKeepExisting:
		return targetValue, nil
	case MapMergeStrategyError:
		return nil, fmt.Errorf("conflict for key '%s': target=%v, source=%v", key, targetValue, sourceValue)
	default:
		return sourceValue, nil
	}
}

// interfaceToSlice 将 interface{} 转换为 []interface{}
func interfaceToSlice(v interface{}) ([]interface{}, bool) {
	if v == nil {
		return nil, false
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, false
	}

	result := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result, true
}

// mergeMapSlices 合并两个切片
func mergeMapSlices(target, source []interface{}, options *MapMergeOptions) (interface{}, error) {
	switch options.Strategy {
	case MapMergeStrategyOverwrite:
		return source, nil // 完全覆盖
	case MapMergeStrategyKeepExisting:
		return target, nil // 保留原有
	default:
		// 合并（默认行为）
		merged := append([]interface{}{}, target...)
		merged = append(merged, source...)
		return merged, nil
	}
}

// ShallowMergeMap 浅合并多个map（不递归）
// 使用泛型，支持任意可比较的键类型
func ShallowMergeMap[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// ConvertMapKeysToString 递归地将 map 的所有键转换为字符串
// 支持嵌套的 map 和 slice
func ConvertMapKeysToString(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case map[interface{}]interface{}:
		return convertInterfaceMapToStringMap(v)
	case map[string]interface{}:
		return convertStringMapRecursive(v)
	case []interface{}:
		return convertSliceRecursive(v)
	default:
		return data
	}
}

// convertInterfaceMapToStringMap 将 map[interface{}]interface{} 转换为 map[string]interface{}
func convertInterfaceMapToStringMap(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		strKey := interfaceToString(k)
		result[strKey] = ConvertMapKeysToString(v)
	}
	return result
}

// convertStringMapRecursive 递归处理 map[string]interface{} 的值
func convertStringMapRecursive(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = ConvertMapKeysToString(v)
	}
	return result
}

// convertSliceRecursive 递归转换切片中的元素
func convertSliceRecursive(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[i] = ConvertMapKeysToString(v)
	}
	return result
}

// interfaceToString 将 interface{} 转换为字符串
func interfaceToString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(val).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(val).Uint(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(val).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// GetNestedMapValue 从嵌套的 map 中获取值，支持路径访问
// 例如: GetNestedMapValue[string](data, "user", "profile", "name")
func GetNestedMapValue[T any](m map[string]interface{}, keys ...string) (T, bool) {
	var zero T
	if len(keys) == 0 {
		return zero, false
	}

	current := interface{}(m)
	for i, key := range keys {
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return zero, false
		}

		value, exists := currentMap[key]
		if !exists {
			return zero, false
		}

		// 如果是最后一个键，尝试类型转换
		if i == len(keys)-1 {
			result, ok := value.(T)
			return result, ok
		}

		current = value
	}

	return zero, false
}

// SetNestedMapValue 在嵌套的 map 中设置值，如果路径不存在则创建
// 例如: SetNestedMapValue(data, "John", "user", "profile", "name")
func SetNestedMapValue(m map[string]interface{}, value interface{}, keys ...string) {
	if len(keys) == 0 {
		return
	}

	current := m
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		next, exists := current[key]
		if !exists {
			next = make(map[string]interface{})
			current[key] = next
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			// 如果存在但不是 map，替换为 map
			nextMap = make(map[string]interface{})
			current[key] = nextMap
		}
		current = nextMap
	}

	// 设置最后一个键的值
	current[keys[len(keys)-1]] = value
}

// FlattenMap 扁平化嵌套的 map，使用点号分隔键
// 例如: {"a": {"b": {"c": 1}}} => {"a.b.c": 1}
func FlattenMap(m map[string]interface{}, separator string) map[string]interface{} {
	if separator == "" {
		separator = "."
	}

	result := make(map[string]interface{})
	flattenMapRecursive(m, "", separator, result)
	return result
}

// flattenMapRecursive 递归扁平化
func flattenMapRecursive(m map[string]interface{}, prefix, separator string, result map[string]interface{}) {
	for key, value := range m {
		newKey := key
		if prefix != "" {
			newKey = prefix + separator + key
		}

		if nestedMap, ok := value.(map[string]interface{}); ok {
			flattenMapRecursive(nestedMap, newKey, separator, result)
		} else {
			result[newKey] = value
		}
	}
}

// UnflattenMap 将扁平化的 map 还原为嵌套结构
// 例如: {"a.b.c": 1} => {"a": {"b": {"c": 1}}}
func UnflattenMap(m map[string]interface{}, separator string) map[string]interface{} {
	if separator == "" {
		separator = "."
	}

	result := make(map[string]interface{})
	for key, value := range m {
		keys := splitMapKey(key, separator)
		SetNestedMapValue(result, value, keys...)
	}
	return result
}

// splitMapKey 分割键字符串
func splitMapKey(key, separator string) []string {
	if separator == "" {
		return []string{key}
	}
	var result []string
	current := ""
	sepLen := len(separator)
	for i := 0; i < len(key); i++ {
		if i+sepLen <= len(key) && key[i:i+sepLen] == separator {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			i += sepLen - 1
		} else {
			current += string(key[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// FilterMap 过滤 map，保留满足条件的键值对
func FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// TransformMapValues 转换 map 的所有值
func TransformMapValues[K comparable, V any, R any](m map[K]V, transform func(V) R) map[K]R {
	result := make(map[K]R, len(m))
	for k, v := range m {
		result[k] = transform(v)
	}
	return result
}

// TransformMapKeys 转换 map 的所有键
func TransformMapKeys[K comparable, V any, R comparable](m map[K]V, transform func(K) R) map[R]V {
	result := make(map[R]V, len(m))
	for k, v := range m {
		result[transform(k)] = v
	}
	return result
}

// CloneMap 深拷贝 map
func CloneMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
