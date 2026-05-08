/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 15:15:16
 * @FilePath: \go-toolbox\pkg\convert\json_slice.go
 * @Description: 切片与JSON字符串互转工具
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package convert

import "encoding/json"

// StringsToJSON 将字符串切片序列化为JSON数组字符串
// 空切片返回空字符串
func StringsToJSON(s []string) string {
	if len(s) == 0 {
		return ""
	}
	data, _ := json.Marshal(s)
	return string(data)
}

// StringsFromJSON 将JSON数组字符串反序列化为字符串切片
// 空字符串返回nil
func StringsFromJSON(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var s []string
	if err := json.Unmarshal([]byte(jsonStr), &s); err != nil {
		return nil, err
	}
	return s, nil
}
