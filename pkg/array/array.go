/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 16:53:49
 * @FilePath: \go-toolbox\pkg\array\array.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package array

import (
	"math"
	"reflect"
	"strconv"

	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// InterfaceArrayDiffSet 计算两个任意类型数组的差集，返回一个新数组包含只在一个数组中出现的元素
/**
 *  @Description: 获取两个切片差集
 *  @param a
 *  @param b
 *  @return []string
 */
func InterfaceArrayDiffSet(a, b interface{}) interface{} {
	diffMap := make(map[interface{}]struct{})

	// 将数组 a 中的元素添加到 diffMap 中
	valA := reflect.ValueOf(a)
	for i := 0; i < valA.Len(); i++ {
		diffMap[valA.Index(i).Interface()] = struct{}{}
	}

	// 删除数组 b 中的元素在 diffMap 中出现的元素，不在的元素添加到 diffMap 中
	valB := reflect.ValueOf(b)
	for i := 0; i < valB.Len(); i++ {
		if _, ok := diffMap[valB.Index(i).Interface()]; ok {
			delete(diffMap, valB.Index(i).Interface())
		} else {
			diffMap[valB.Index(i).Interface()] = struct{}{}
		}
	}

	// 将 diffMap 中的元素转换为切片
	diff := make([]interface{}, 0, len(diffMap))
	for key := range diffMap {
		diff = append(diff, key)
	}

	return diff
}

// InterfaceArrayUnion 计算多种类型数组的并集，返回一个新的数组包含所有元素，不包含重复元素
func InterfaceArrayUnion(arr1, arr2 interface{}) interface{} {
	unionMap := make(map[interface{}]bool)
	union := []interface{}{}

	// 将 arr1 中的元素加入到 unionMap 中，去重
	val1 := reflect.ValueOf(arr1)
	for i := 0; i < val1.Len(); i++ {
		element := val1.Index(i).Interface()
		unionMap[element] = true
	}

	// 将 arr2 中的元素加入到 unionMap 中，去重
	val2 := reflect.ValueOf(arr2)
	for i := 0; i < val2.Len(); i++ {
		element := val2.Index(i).Interface()
		unionMap[element] = true
	}

	// 将 map 中的键转换为切片
	for key := range unionMap {
		union = append(union, key)
	}

	return union
}

// IsInterfaceArrayExistElement 检查 Interface 数组中是否包含某个元素
/**
 *  @Description: Interface 数组中是否包含某个元素
 *  @param array
 *  @param element
 *  @return exist
 */
func IsInterfaceArrayExistElement(array []interface{}, element interface{}) (exist bool) {
	for _, a := range array {
		if a == element {
			return true
		}
	}
	return
}

// IsExistRepeatInInterfaceArray 检查 Interface 数组中是否存在重复对象
/**
 *  @Description: Interface 数组中是否存在重复对象
 *  @param array
 *  @return exist
 */
func IsExistRepeatInInterfaceArray(array []interface{}) (exist bool) {
	m := make(map[interface{}]int)
	for _, v := range array {
		_, ok := m[v]
		if ok {
			return true
		} else {
			m[v] = 1
		}
	}
	return false
}

// RemoveEmptyInterfaceInArray 移除空的对象
/**
 *  @Description: 移除空的对象
 *  @param array
 *  @return answer
 */
func RemoveEmptyInterfaceInArray(array []interface{}) []interface{} {
	var result []interface{}
	for _, v := range array {
		// 只移除空字符串和零值，但保留 nil
		if !validator.IsEmptyValue(reflect.ValueOf(v)) || v == nil {
			result = append(result, v)
		}
	}
	return result
}

// Int64ToStringWithDecimals 将 int64 转换为包含小数点后指定位数的字符串
func Int64ToStringWithDecimals(num int64, digit int) string {
	// 计算除数，动态生成指数部分
	divisor := math.Pow10(digit)
	// 将 int64 转换为 float64，然后除以动态生成的除数
	flt := float64(num) / divisor
	// 将 float64 格式化为字符串，保留小数点后指定位数
	str := strconv.FormatFloat(flt, 'f', digit, 64)
	return str
}

// RemoveDuplicatesInInterfaceSlice 移除掉重复值
func RemoveDuplicatesInInterfaceSlice(numbers []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	uniqueNumbers := []interface{}{}

	for _, num := range numbers {
		if !m[num] {
			m[num] = true
			uniqueNumbers = append(uniqueNumbers, num)
		}
	}

	return uniqueNumbers
}

// RemoveZeroInInterfaceSlice 移除掉0的值
func RemoveZeroInInterfaceSlice(arr []interface{}) []interface{} {
	var result []interface{}
	for _, val := range arr {
		if num, ok := val.(int); ok {
			if num != 0 {
				result = append(result, val)
			}
		}
	}
	return result
}
