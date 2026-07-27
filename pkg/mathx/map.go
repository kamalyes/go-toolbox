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
//
// Deprecated: 使用 MapAssign 替代，命名风格统一为 MapXxx
// 迁移示例：
//
//	// 旧：
//	result := ShallowMergeMap(m1, m2, m3)
//	// 新：
//	result := MapAssign(m1, m2, m3)
//
// 注意：MapAssign 额外支持自定义 map 类型（~map[K]V），保留原 map 类型
func ShallowMergeMap[K comparable, V any](maps ...map[K]V) map[K]V {
	return MapAssign[K, V](maps...)
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
	separator = IfNotEmpty(separator, ".")

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
	separator = IfNotEmpty(separator, ".")

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
//
// Deprecated: 使用 MapPickBy 替代，命名风格统一为 MapXxx
// 迁移示例：
//
//	// 旧：
//	result := FilterMap(m, func(k K, v V) bool { return v > 0 })
//	// 新：
//	result := MapPickBy(m, func(k K, v V) bool { return v > 0 })
//
// 注意：MapPickBy 额外支持自定义 map 类型（~map[K]V），保留原 map 类型
func FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	return MapPickBy(m, predicate)
}

// TransformMapValues 转换 map 的所有值
//
// Deprecated: 使用 MapMapValues 替代，回调接收 (value, key) 提供更完整上下文
// 迁移示例：
//
//	// 旧：
//	result := TransformMapValues(m, func(v V) R { return f(v) })
//	// 新：
//	result := MapMapValues(m, func(v V, _ K) R { return f(v) })
//
// 注意：MapMapValues 回调签名为 (value, key)，可忽略 key 实现等价功能
func TransformMapValues[K comparable, V any, R any](m map[K]V, transform func(V) R) map[K]R {
	return MapMapValues(m, func(v V, _ K) R { return transform(v) })
}

// TransformMapKeys 转换 map 的所有键
//
// Deprecated: 使用 MapMapKeys 替代，回调接收 (value, key) 提供更完整上下文
// 迁移示例：
//
//	// 旧：
//	result := TransformMapKeys(m, func(k K) R { return f(k) })
//	// 新：
//	result := MapMapKeys(m, func(_ V, k K) R { return f(k) })
//
// 注意：MapMapKeys 回调签名为 (value, key)，可忽略 value 实现等价功能
func TransformMapKeys[K comparable, V any, R comparable](m map[K]V, transform func(K) R) map[R]V {
	return MapMapKeys(m, func(_ V, k K) R { return transform(k) })
}

// CloneMap 深拷贝 map（浅拷贝键值对，不递归复制值引用的对象）
//
// Deprecated: 使用 MapClone 替代，命名风格统一为 MapXxx
// 迁移示例：
//
//	// 旧：
//	cloned := CloneMap(m)
//	// 新：
//	cloned := MapClone(m)
func CloneMap[K comparable, V any](m map[K]V) map[K]V {
	return MapClone(m)
}

// MapClone 浅拷贝 map 的键值对（不递归复制值引用的对象）
// 与 CloneMap 功能等价，命名遵循 MapXxx 风格
func MapClone[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// LayeredMerger 多层级键值对合并器（避免重复传递字段名）
type LayeredMerger[T any, KV any] struct {
	keyFieldName   string
	valueFieldName string
}

// NewLayeredMerger 创建多层级合并器
//
// 参数：
//   - keyFieldName: 键值对结构中 key 字段的名称（如 "Key"）
//   - valueFieldName: 键值对结构中 value 字段的名称（如 "Value"）
//
// 示例：
//
//	merger := NewLayeredMerger[Config, LocalizedText]("Key", "Value")
//	result := merger.Merge(layers, func(c *Config) []LocalizedText { return c.Messages })
func NewLayeredMerger[T any, KV any](keyFieldName, valueFieldName string) *LayeredMerger[T, KV] {
	return &LayeredMerger[T, KV]{
		keyFieldName:   keyFieldName,
		valueFieldName: valueFieldName,
	}
}

// Merge 执行多层级合并
//
// 参数：
//   - layers: 配置层级切片，从低到高优先级（越靠后优先级越高）
//   - fieldGetter: 从配置对象中提取键值对切片的函数
//
// 返回：合并后的键值对切片
func (m *LayeredMerger[T, KV]) Merge(layers []*T, fieldGetter func(*T) []KV) []KV {
	return MergeLayeredKeyValues(layers, fieldGetter, m.keyFieldName, m.valueFieldName)
}

// MergeLayeredKeyValues 多层级键值对合并工具（支持传入任意数量的配置层级）
//
// 参数：
//   - layers: 配置层级切片，从低到高优先级（越靠后优先级越高，后面的会覆盖前面的）
//   - fieldGetter: 从配置对象中提取键值对切片的函数
//   - keyFieldName: 键值对结构中 key 字段的名称（如 "Key"）
//   - valueFieldName: 键值对结构中 value 字段的名称（如 "Value"）
//
// 返回：合并后的键值对切片，保持 key 首次出现的顺序
//
// 特性：
//   - 支持任意层级的链式合并
//   - 自动跳过 nil 层级
//   - 自动跳过空值（值为空字符串的项不会覆盖已有值）
//   - 保持 key 的首次出现顺序
//   - 使用反射自动提取字段，支持任意结构体
//
// 示例：
//
//	type LocalizedText struct {
//	    Key   string
//	    Value string
//	}
//
//	type Config struct {
//	    Messages []LocalizedText
//	}
//
//	result := MergeLayeredKeyValues(
//	    []*Config{hardcodedDefault, ownerConfig, agentConfig},
//	    func(c *Config) []LocalizedText { return c.Messages },
//	    "Key",
//	    "Value",
//	)
func MergeLayeredKeyValues[T any, KV any](
	layers []*T,
	fieldGetter func(*T) []KV,
	keyFieldName string,
	valueFieldName string,
) []KV {
	if len(layers) == 0 {
		return []KV{}
	}

	// 使用 map 存储合并结果，key 为字段 key，value 为字段 value
	valueMap := make(map[string]string)
	// 记录 key 的出现顺序
	keyOrder := make([]string, 0)

	// 按层级顺序合并（从低优先级到高优先级）
	for _, layer := range layers {
		if layer == nil {
			continue
		}

		items := fieldGetter(layer)
		for i := range items {
			item := &items[i]

			// 使用反射提取 key 和 value
			key := extractFieldValue(item, keyFieldName)
			value := extractFieldValue(item, valueFieldName)

			// 跳过空值
			if value == "" {
				continue
			}

			// 记录首次出现的顺序
			if _, exists := valueMap[key]; !exists {
				keyOrder = append(keyOrder, key)
			}

			// 覆盖或新增
			valueMap[key] = value
		}
	}

	// 按顺序构建结果
	result := make([]KV, 0, len(keyOrder))
	for _, key := range keyOrder {
		value := valueMap[key]
		item := buildKeyValueItem[KV](keyFieldName, key, valueFieldName, value)
		result = append(result, item)
	}

	return result
}

// extractFieldValue 通过反射提取结构体字段的字符串值
func extractFieldValue(item interface{}, fieldName string) string {
	val := reflect.ValueOf(item)

	// 处理指针
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}

	// 确保是结构体
	if val.Kind() != reflect.Struct {
		return ""
	}

	// 获取字段值
	fieldVal := val.FieldByName(fieldName)
	if !fieldVal.IsValid() {
		return ""
	}

	// 转换为字符串
	if fieldVal.Kind() == reflect.String {
		return fieldVal.String()
	}

	return ""
}

// buildKeyValueItem 创建 key-value 结构体实例
func buildKeyValueItem[KV any](keyFieldName, keyValue, valueFieldName, valueValue string) KV {
	var item KV
	itemType := reflect.TypeOf(item)
	itemVal := reflect.New(itemType).Elem()

	// 设置 key 字段
	keyField := itemVal.FieldByName(keyFieldName)
	if keyField.IsValid() && keyField.CanSet() && keyField.Kind() == reflect.String {
		keyField.SetString(keyValue)
	}

	// 设置 value 字段
	valueField := itemVal.FieldByName(valueFieldName)
	if valueField.IsValid() && valueField.CanSet() && valueField.Kind() == reflect.String {
		valueField.SetString(valueValue)
	}

	return itemVal.Interface().(KV)
}

// MapKeys 提取一个或多个 map 的所有键，返回切片
// 支持多个 map，键来自所有 map 的并集（不去重，见 MapUniqKeys 去重版本）
func MapKeys[K comparable, V any](in ...map[K]V) []K {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	result := make([]K, 0, size)
	for i := range in {
		for k := range in[i] {
			result = append(result, k)
		}
	}
	return result
}

// MapUniqKeys 提取一个或多个 map 的所有去重键
func MapUniqKeys[K comparable, V any](in ...map[K]V) []K {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	seen := make(map[K]struct{}, size)
	result := make([]K, 0, size)
	for i := range in {
		for k := range in[i] {
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
			result = append(result, k)
		}
	}
	return result
}

// MapValues 提取一个或多个 map 的所有值，返回切片
func MapValues[K comparable, V any](in ...map[K]V) []V {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	result := make([]V, 0, size)
	for i := range in {
		for _, v := range in[i] {
			result = append(result, v)
		}
	}
	return result
}

// MapUniqValues 提取一个或多个 map 的所有去重值（V 必须可比较）
func MapUniqValues[K, V comparable](in ...map[K]V) []V {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	seen := make(map[V]struct{}, size)
	result := make([]V, 0, size)
	for i := range in {
		for _, v := range in[i] {
			if _, exists := seen[v]; exists {
				continue
			}
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// ============================================================================
// 键检查 / 值获取
// ============================================================================

// MapHasKey 检查键是否存在
func MapHasKey[K comparable, V any](in map[K]V, key K) bool {
	_, ok := in[key]
	return ok
}

// MapValueOr 返回键对应的值或 fallback
func MapValueOr[K comparable, V any](in map[K]V, key K, fallback V) V {
	if v, ok := in[key]; ok {
		return v
	}
	return fallback
}

// ============================================================================
// 选择 / 移除（PickBy / OmitBy）
// ============================================================================

// MapPickBy 按 predicate 过滤，保留满足条件的键值对，返回同类型 map
func MapPickBy[K comparable, V any, M ~map[K]V](in M, predicate func(key K, value V) bool) M {
	r := M{}
	for k, v := range in {
		if predicate(k, v) {
			r[k] = v
		}
	}
	return r
}

// MapOmitBy 按 predicate 过滤，移除满足条件的键值对（与 PickBy 相反）
func MapOmitBy[K comparable, V any, M ~map[K]V](in M, predicate func(key K, value V) bool) M {
	r := M{}
	for k, v := range in {
		if !predicate(k, v) {
			r[k] = v
		}
	}
	return r
}

// MapPickByKeys 按键列表过滤，仅保留 keys 中存在的键
func MapPickByKeys[K comparable, V any, M ~map[K]V](in M, keys []K) M {
	r := M{}
	for i := range keys {
		if v, ok := in[keys[i]]; ok {
			r[keys[i]] = v
		}
	}
	return r
}

// MapOmitByKeys 按键列表过滤，移除 keys 中存在的键
func MapOmitByKeys[K comparable, V any, M ~map[K]V](in M, keys []K) M {
	r := M{}
	for k, v := range in {
		r[k] = v
	}
	for i := range keys {
		delete(r, keys[i])
	}
	return r
}

// MapPickByValues 按值列表过滤，仅保留 values 中存在的值（V 必须可比较）
func MapPickByValues[K, V comparable, M ~map[K]V](in M, values []V) M {
	set := make(map[V]struct{}, len(values))
	for i := range values {
		set[values[i]] = struct{}{}
	}
	r := M{}
	for k, v := range in {
		if _, ok := set[v]; ok {
			r[k] = v
		}
	}
	return r
}

// MapOmitByValues 按值列表过滤，移除 values 中存在的值（V 必须可比较）
func MapOmitByValues[K, V comparable, M ~map[K]V](in M, values []V) M {
	set := make(map[V]struct{}, len(values))
	for i := range values {
		set[values[i]] = struct{}{}
	}
	r := M{}
	for k, v := range in {
		if _, ok := set[v]; !ok {
			r[k] = v
		}
	}
	return r
}

// ============================================================================
// map <-> 键值对切片互转
// ============================================================================

// MapEntry 键值对
type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// MapEntries 将 map 转为键值对切片
func MapEntries[K comparable, V any](in map[K]V) []MapEntry[K, V] {
	entries := make([]MapEntry[K, V], 0, len(in))
	for k, v := range in {
		entries = append(entries, MapEntry[K, V]{Key: k, Value: v})
	}
	return entries
}

// MapFromEntries 将键值对切片转为 map（同 key 后者覆盖前者）
func MapFromEntries[K comparable, V any](entries []MapEntry[K, V]) map[K]V {
	out := make(map[K]V, len(entries))
	for i := range entries {
		out[entries[i].Key] = entries[i].Value
	}
	return out
}

// ============================================================================
// 反转 / 合并
// ============================================================================

// MapInvert 反转 map 的键和值（V 必须可比较）
// 若有重复值，后者覆盖前者
func MapInvert[K, V comparable](in map[K]V) map[V]K {
	out := make(map[V]K, len(in))
	for k, v := range in {
		out[v] = k
	}
	return out
}

// MapAssign 从左到右合并多个 map，后者覆盖前者
// 与 ShallowMergeMap 区别：ShallowMergeMap 接收可变参数 map[K]V，
// MapAssign 接收自定义 map 类型 M（~map[K]V）保持类型
func MapAssign[K comparable, V any, M ~map[K]V](maps ...M) M {
	count := 0
	for i := range maps {
		count += len(maps[i])
	}
	out := M{}
	for i := range maps {
		for k, v := range maps[i] {
			out[k] = v
		}
	}
	return out
}

// ============================================================================
// map 转切片
// ============================================================================

// MapToSlice 将 map 转为切片，iteratee 接收 (key, value) 返回结果元素
func MapToSlice[K comparable, V any, R any](in map[K]V, iteratee func(key K, value V) R) []R {
	result := make([]R, 0, len(in))
	for k, v := range in {
		result = append(result, iteratee(k, v))
	}
	return result
}

// MapFilterToSlice 过滤+转换 map 到切片
// iteratee 返回 (结果, 是否保留)，仅保留的元素进入切片
func MapFilterToSlice[K comparable, V any, R any](in map[K]V, iteratee func(key K, value V) (R, bool)) []R {
	result := make([]R, 0, len(in))
	for k, v := range in {
		if r, ok := iteratee(k, v); ok {
			result = append(result, r)
		}
	}
	return result
}

// MapFilterKeys 按 predicate 过滤，返回满足条件的键切片（Filter + Keys 组合）
func MapFilterKeys[K comparable, V any](in map[K]V, predicate func(key K, value V) bool) []K {
	result := make([]K, 0, len(in))
	for k, v := range in {
		if predicate(k, v) {
			result = append(result, k)
		}
	}
	return result
}

// MapFilterValues 按 predicate 过滤，返回满足条件的值切片（Filter + Values 组合）
func MapFilterValues[K comparable, V any](in map[K]V, predicate func(key K, value V) bool) []V {
	result := make([]V, 0, len(in))
	for k, v := range in {
		if predicate(k, v) {
			result = append(result, v)
		}
	}
	return result
}

// ============================================================================
// map 键/值/键值对转换（MapEntries 风格的全量转换）
// ============================================================================

// MapMapEntries 同时转换 key 和 value，返回新类型的 map
func MapMapEntries[K1 comparable, V1 any, K2 comparable, V2 any](in map[K1]V1, iteratee func(key K1, value V1) (K2, V2)) map[K2]V2 {
	result := make(map[K2]V2, len(in))
	for k1 := range in {
		k2, v2 := iteratee(k1, in[k1])
		result[k2] = v2
	}
	return result
}

// MapMapKeys 转换所有键，返回新 map（值类型不变）
func MapMapKeys[K comparable, V any, R comparable](in map[K]V, iteratee func(value V, key K) R) map[R]V {
	result := make(map[R]V, len(in))
	for k, v := range in {
		result[iteratee(v, k)] = v
	}
	return result
}

// MapMapValues 转换所有值，返回新 map（键类型不变）
// 与 TransformMapValues 区别：TransformMapValues 只接收 value，本方法接收 (value, key)
func MapMapValues[K comparable, V any, R any](in map[K]V, iteratee func(value V, key K) R) map[K]R {
	result := make(map[K]R, len(in))
	for k, v := range in {
		result[k] = iteratee(v, k)
	}
	return result
}
