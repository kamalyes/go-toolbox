/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 13:00:16
 * @FilePath: \go-toolbox\chinese\chinese.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package chinese

var Chinese = new(chinese)

type chinese struct{}

// Len
/**
 *  @Description: 获取中文字符串长度
 *  @receiver c
 *  @param str
 *  @return int
 */
func (c chinese) Len(str string) int {
	rt := []rune(str)
	return len(rt)
}

// Cut
/**
 *  @Description: 截取中文字符串
 *  @receiver c
 *  @param str
 *  @param start
 *  @param end
 *  @return string
 */
func (c chinese) Cut(str string, start int, end int) string {
	rt := []rune(str)
	return string(rt[start:end])
}
