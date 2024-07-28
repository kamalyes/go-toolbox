/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 10:14:34
 * @FilePath: \go-middleware\array\array.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package array

import "strings"

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
