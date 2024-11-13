/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 23:20:55
 * @FilePath: \go-toolbox\pkg\mathx\array.go
 * @Description: 包含与数组相关的通用函数，例如计算最小值和最大值、差集、并集等。
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package mathx

import (
	"errors"
	"reflect"

	"github.com/kamalyes/go-toolbox/pkg/types"
	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// ArrayMinMax 计算列表中元素的最小值或最大值。
// 接收一个切片和一个 MinMaxFunc 类型的函数，
// 根据提供的函数决定是计算最小值还是最大值。
// 如果列表为空，则返回错误。
func ArrayMinMax[T any](list []T, f types.MinMaxFunc[T]) (T, error) {
	if len(list) == 0 {
		var zero T
		return zero, errors.New("列表为空") // 返回错误信息
	}

	result := list[0] // 初始化结果为列表的第一个元素
	for _, v := range list[1:] {
		result = f(result, v) // 使用提供的函数更新结果
	}
	return result, nil // 返回最终结果和 nil 错误
}

// ArrayDiffSet 计算两个任意类型数组的差集。
// 返回一个新数组，包含只在一个数组中出现的元素。
func ArrayDiffSet(arr1, arr2 []interface{}) []interface{} {
	// 使用 map 来存储 arr2 的元素，以便快速查找
	set1 := make(map[interface{}]struct{}, len(arr1))
	set2 := make(map[interface{}]struct{}, len(arr2))

	// 将 arr1 的元素存入集合
	for _, item := range arr1 {
		set1[item] = struct{}{}
	}

	// 将 arr2 的元素存入集合
	for _, item := range arr2 {
		set2[item] = struct{}{}
	}

	// 计算只在 arr1 中的元素
	diff := make([]interface{}, 0, len(arr1)+len(arr2)) // 预分配空间

	for item := range set1 {
		if _, found := set2[item]; !found {
			diff = append(diff, item)
		}
	}

	// 计算只在 arr2 中的元素
	for item := range set2 {
		if _, found := set1[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
}

// ArrayUnion 计算两个数组的并集。
// 返回一个新的数组，包含所有元素，不包含重复元素。
func ArrayUnion[T comparable](arr1, arr2 []T) []T {
	unionMap := make(map[T]struct{}) // 使用映射去重
	for _, element := range arr1 {
		unionMap[element] = struct{}{} // 将 arr1 中的元素加入到 unionMap
	}
	for _, element := range arr2 {
		unionMap[element] = struct{}{} // 将 arr2 中的元素加入到 unionMap
	}

	union := make([]T, 0, len(unionMap)) // 创建并集切片
	for key := range unionMap {
		union = append(union, key) // 将 unionMap 中的键转换为切片
	}
	return union
}

// ArrayContains 检查切片中是否包含某个元素。
// 返回布尔值，表示元素是否存在于切片中。
func ArrayContains[T comparable](array []T, element T) bool {
	for _, a := range array {
		if a == element {
			return true // 找到元素，返回 true
		}
	}
	return false // 未找到元素，返回 false
}

// ArrayHasDuplicates 检查切片中是否存在重复对象。
// 返回布尔值，表示是否存在重复元素。
func ArrayHasDuplicates[T comparable](array []T) bool {
	m := make(map[T]struct{}) // 使用映射存储已见元素
	for _, v := range array {
		if _, ok := m[v]; ok {
			return true // 找到重复元素，返回 true
		}
		m[v] = struct{}{} // 添加新元素
	}
	return false // 未找到重复元素，返回 false
}

// ArrayRemoveEmpty 移除切片中的空对象。
// 返回一个新切片，包含所有非空元素。
func ArrayRemoveEmpty[T any](array []T) []T {
	result := make([]T, 0, len(array)) // 创建结果切片
	for _, v := range array {
		if !validator.IsEmptyValue(reflect.ValueOf(v)) {
			result = append(result, v) // 仅添加非空元素
		}
	}
	return result
}

// ArrayRemoveDuplicates 移除切片中的重复值。
// 返回一个新切片，包含所有唯一元素。
func ArrayRemoveDuplicates[T comparable](numbers []T) []T {
	m := make(map[T]struct{})                   // 使用映射去重
	uniqueNumbers := make([]T, 0, len(numbers)) // 创建唯一元素切片
	for _, num := range numbers {
		if _, exists := m[num]; !exists {
			m[num] = struct{}{}                        // 添加新元素
			uniqueNumbers = append(uniqueNumbers, num) // 仅添加唯一元素
		}
	}
	return uniqueNumbers
}

// ArrayRemoveZero 移除切片中的零值。
// 返回一个新切片，包含所有非零元素。
func ArrayRemoveZero(arr []int) []int {
	result := make([]int, 0, len(arr)) // 创建结果切片
	for _, val := range arr {
		if val != 0 {
			result = append(result, val) // 仅添加非零元素
		}
	}
	return result
}

// ArrayChunk 将一个切片分割成多个子切片。
// size 参数指定每个子切片的大小。
// 返回一个包含子切片的切片。
func ArrayChunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil // 如果 size <= 0，则返回 nil
	}

	var batches [][]T // 创建子切片切片
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice) // 确保不超出边界
		}
		batches = append(batches, slice[i:end]) // 切片而不复制
	}
	return batches
}
