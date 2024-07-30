/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 13:35:22
 * @FilePath: \go-toolbox\array\array.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package array

import (
	"math"
	"strconv"
	"strings"
)

// StrArrayDiffSet
/**
 *  @Description: 获取两个切片差集
 *  @param a
 *  @param b
 *  @return []string
 */
func StrArrayDiffSet(a []string, b []string) []string {
	var c []string
	temp := map[string]struct{}{} // map[string]struct{}{}创建了一个key类型为String值类型为空struct的map，Equal -> make(map[string]struct{})
	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{} // 空struct 不占内存空间
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}
	return c
}

// IsStrArrayExistArray
/**
 *  @Description: 字符串数组是否包含字符串
 *  @param array
 *  @param str
 *  @return exist
 */
func IsStrArrayExistArray(array []string, str string) (exist bool) {
	for _, a := range array {
		if a == str {
			return true
		}
	}
	return
}

// IsExistRepeatInArray
/**
 *  @Description: 数组中是否存在重复对象
 *  @param array
 *  @return exist
 */
func IsExistRepeatInArray(array []string) (exist bool) {
	m := make(map[string]int)
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

// RemoveEmptyStrInArray
/**
 *  @Description: 移除
 *  @param array
 *  @return answer
 */
func RemoveEmptyStrInArray(array []string) (answer []string) {
	for _, str := range array {
		if strings.TrimSpace(str) != "" {
			answer = append(answer, str)
		}
	}
	return answer
}

// 将 int64 转换为包含小数点后指定位数的字符串
func Int64ToStringWithDecimals(num int64, digit int) string {
	// 计算除数，动态生成指数部分
	divisor := math.Pow10(digit)
	// 将 int64 转换为 float64，然后除以动态生成的除数
	flt := float64(num) / divisor
	// 将 float64 格式化为字符串，保留小数点后指定位数
	str := strconv.FormatFloat(flt, 'f', digit, 64)
	return str
}

// 比较两个list，返回在slice1中但不在slice2中的元素
func ListDifferenceInt64(slice1 []int64, slice2 []int64) []int64 {
	// 使用 map 存储 slice2 的元素
	slice2Map := make(map[int64]struct{}, len(slice2))
	for _, item := range slice2 {
		slice2Map[item] = struct{}{}
	}

	// 预分配 diff 切片的容量
	diff := make([]int64, 0, len(slice1))

	// 查找在 slice1 中但不在 slice2 中的元素
	for _, item := range slice1 {
		if _, found := slice2Map[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
}

// 移除掉重复值
func RemoveDuplicatesInt32(numbers []int32) []int32 {
	m := make(map[int32]bool)
	uniqueNumbers := []int32{}

	for _, num := range numbers {
		if !m[num] {
			m[num] = true
			uniqueNumbers = append(uniqueNumbers, num)
		}
	}

	return uniqueNumbers
}

// 移除掉0的值
func RemoveZero(arr []int) []int {
	var result []int
	for _, num := range arr {
		if num != 0 {
			result = append(result, num)
		}
	}
	return result
}

// 手机号脱敏
func DesensitizePhoneNum(phoneNum string) string {
	if len(phoneNum) != 11 {
		return phoneNum
	}
	return phoneNum[:3] + "****" + phoneNum[7:]
}
