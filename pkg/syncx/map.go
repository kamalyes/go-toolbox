/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-15 11:55:06
 * @FilePath: \go-toolbox\pkg\syncx\map.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"net/url"
	"sort"
	"strings"
	"sync"
)

// Map 是一个线程安全的映射，使用泛型 K 和 V。
type Map[K comparable, V comparable] struct { // 确保 V 是可比较的
	mp sync.Map // 使用 sync.Map 来实现线程安全
}

// NewMap 创建一个新的 Map 实例。
func NewMap[K comparable, V comparable]() *Map[K, V] {
	return &Map[K, V]{} // 返回一个新的 Map 实例
}

// CompareAndDelete 比较指定键的值，如果相等则删除该键的值。
func (m *Map[K, V]) CompareAndDelete(key K, value V) bool {
	actual, loaded := m.mp.Load(key)
	if !loaded {
		return false
	}
	if actualValue, ok := actual.(V); ok && actualValue == value {
		m.mp.Delete(key)
		return true
	}
	return false
}

// CompareAndSwap 比较指定键的现有值，如果相等则将其替换为新值。
func (m *Map[K, V]) CompareAndSwap(key K, old, new V) bool {
	actual, loaded := m.mp.Load(key)
	if !loaded {
		return false
	}
	if actualValue, ok := actual.(V); ok && actualValue == old {
		m.mp.Store(key, new)
		return true
	}
	return false
}

// Delete 删除指定键的值。
func (m *Map[K, V]) Delete(key K) {
	m.mp.Delete(key) // 调用 sync.Map 的 Delete 方法
}

// Load 获取指定键的值。
func (m *Map[K, V]) Load(key K) (V, bool) {
	return m.loadValue(m.mp.Load(key))
}

// LoadAndDelete 方法从 Map 中加载并删除指定键的值
func (m *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	return m.loadValue(m.mp.LoadAndDelete(key))
}

// LoadOrStore 获取指定键的值，如果不存在则存储新值。
func (m *Map[K, V]) LoadOrStore(key K, new V) (V, bool) {
	value, ok := m.mp.LoadOrStore(key, new) // 调用 sync.Map 的 LoadOrStore 方法
	return value.(V), ok                    // 返回转换后的值和成功标志
}

// Range 遍历 Map 中的所有键值对。
func (m *Map[K, V]) Range(run func(key K, value V) bool) {
	m.mp.Range(func(k, v any) bool {
		key := k.(K)           // 将 k 转换为 K 类型
		value := v.(V)         // 将 v 转换为 V 类型
		return run(key, value) // 调用用户提供的函数
	})
}

// Size 返回当前 Map 中元素的数量
func (m *Map[K, V]) Size() int {
	count := 0
	m.Range(func(_ K, _ V) bool {
		count++
		return true // 继续遍历
	})
	return count
}

// Clear 清空 Map 中的所有键值对
func (m *Map[K, V]) Clear() {
	m.Range(func(key K, _ V) bool {
		m.Delete(key) // 删除每个键
		return true   // 继续遍历
	})
}

// Keys 返回所有键的切片
func (m *Map[K, V]) Keys() []K {
	var keys []K
	m.Range(func(key K, _ V) bool {
		keys = append(keys, key) // 添加键到切片中
		return true              // 继续遍历
	})
	return keys
}

// Values 返回所有值的切片
func (m *Map[K, V]) Values() []V {
	var values []V
	m.Range(func(_ K, value V) bool {
		values = append(values, value) // 添加值到切片中
		return true                    // 继续遍历
	})
	return values
}

// Store 设置指定键的值。
func (m *Map[K, V]) Store(key K, value V) {
	m.mp.Store(key, value) // 调用 sync.Map 的 Store 方法
}

// Swap 替换指定键的值，并返回之前的值。
func (m *Map[K, V]) Swap(key K, value V) (pre V, ok bool) {
	// 尝试加载当前值
	previous, loaded := m.mp.Load(key)
	if loaded {
		// 如果加载成功，先删除当前键
		m.mp.Delete(key)
		// 然后存储新值
		m.mp.Store(key, value)
		// 返回之前的值
		return m.loadValue(previous, true)
	}
	// 如果键不存在，直接存储新值
	m.mp.Store(key, value)
	return zeroValue[V](), false // 返回零值和 false
}

// loadValue 处理从 sync.Map 加载的值，返回类型 V 和成功标志
func (m *Map[K, V]) loadValue(value any, ok bool) (V, bool) {
	if !ok {
		return zeroValue[V](), false // 如果未找到值，返回零值
	}
	if result, ok := value.(V); ok {
		return result, true // 返回转换后的值和成功标志
	}
	return zeroValue[V](), false // 类型断言失败，返回零值和 false
}

// zeroValue 返回类型 V 的零值
func zeroValue[T any]() T {
	var zero T
	return zero
}

// Equals 函数用于比较两个键值对是否相等，
// 需要用户提供自定义的比较函数。
func (m *Map[K, V]) Equals(key K, value V, cmpFunc func(existing V) bool) bool {
	existing, ok := m.Load(key) // 加载现有值
	if !ok {
		return false // 如果不存在，返回 false
	}
	return cmpFunc(existing) // 使用用户提供的比较函数进行比较
}

// Clone 克隆当前 Map 实例，返回一个新的 Map 实例
func (m *Map[K, V]) Clone() *Map[K, V] {
	newMap := NewMap[K, V]() // 创建一个新的 Map 实例
	m.Range(func(key K, value V) bool {
		newMap.Store(key, value) // 将每个键值对存储到新实例中
		return true              // 继续遍历
	})
	return newMap // 返回新的 Map 实例
}

// CopyMeta 复制 src 中的所有键值对到 dst 中。
// 如果 dst 为 nil，则不进行任何操作。
func CopyMeta(src, dst map[string]string) {
	if dst == nil {
		return // 如果目标 map 为 nil，直接返回
	}

	// 预先分配目标 map 的容量（可选）
	if len(dst) == 0 {
		dst = make(map[string]string, len(src))
	}

	for k, v := range src {
		dst[k] = v // 复制每个键值对
	}
}

// MetaStringToMap 将 meta 字符串转换为键值对的 map
func MetaStringToMap(meta string) map[string]string {
	if meta == "" {
		return nil // 返回 nil 而不是空 map
	}

	// 使用 url.ParseQuery 解析 meta 字符串
	v, err := url.ParseQuery(meta)
	if err != nil {
		return nil // 如果解析失败，返回 nil
	}

	rt := make(map[string]string, len(v)) // 预分配 map 的容量
	for key, values := range v {
		// 只取第一个值，因为 Query 可能返回多个值
		if len(values) > 0 && values[0] != "" {
			rt[key] = values[0] // 将每个键的第一个非空值存入结果 map
		}
	}
	return rt // 返回最终的结果 map
}

// MetaMapToString 将 map 转换为 meta 字符串
func MetaMapToString(meta map[string]string) string {
	if len(meta) == 0 {
		return "" // 如果 map 为空，返回空字符串
	}

	var buf strings.Builder
	keys := make([]string, 0, len(meta))
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys) // 对键进行排序

	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&') // 添加分隔符
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(meta[k]))
	}
	return buf.String()
}
