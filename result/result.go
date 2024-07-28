/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 09:36:36
 * @FilePath: \go-middleware\result\result.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (

	// FAIL 失败默认code返回1
	FAIL = "1"

	// SUCCESS 成功默认code返回0
	SUCCESS = "0"
)

// Response 统一 json 结构体
type Response struct {

	/** 状态码 */
	Code string `json:"code"`

	/** 内容体 */
	Content interface{} `json:"content"`

	/** 消息 */
	Message string `json:"message"`
}

// Result gin 统一返回
func Result(code string, content interface{}, message string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		content,
		message,
	})
}

// Ok 成功
func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "成功", c)
}

// OkMsg 带message消息的成功
func OkMsg(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

// OkData 带数据的成功
func OkData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "成功", c)
}

// OkDataMsg 带数据和返回消息的成功
func OkDataMsg(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

// Fail 失败
func Fail(c *gin.Context) {
	Result(FAIL, map[string]interface{}{}, "失败", c)
}

// FailMsg 带message消息的失败
func FailMsg(message string, c *gin.Context) {
	Result(FAIL, map[string]interface{}{}, message, c)
}

// FailDataMsg 带数据和返回消息的失败
func FailDataMsg(data interface{}, message string, c *gin.Context) {
	Result(FAIL, data, message, c)
}
