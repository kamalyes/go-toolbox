/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 13:20:41
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 15:27:31
 * @FilePath: \go-toolbox\pkg\stringx\quote.go
 * @Description: 引用 JSON 字符串
 */
package stringx

import "strconv"

// QuoteJSONBytes 将字符串按 JSON 字符串规则转义，并返回带双引号的字节切片
func QuoteJSONBytes(str string) []byte {
	return strconv.AppendQuote(nil, str)
}
