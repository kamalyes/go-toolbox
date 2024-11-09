/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:55:15
 * @FilePath: \go-toolbox\pkg\sign\hmac.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

const salt string = "*$salt@*"

// HmacSha256Base64 计算hmac
func HmacSha256Base64(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, err := io.WriteString(h, message)
	if err != nil {
		log.Fatalln("计算签名错误：" + err.Error())
		return ""
	}
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sign
}

// HmacSha256Hex 字符串计算sha256之后转hex
func HmacSha256Hex(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, err := io.WriteString(h, message)
	if err != nil {
		log.Fatalln("计算签名错误：" + err.Error())
		return ""
	}
	sign := hex.EncodeToString(h.Sum(nil))
	return sign
}

// SHA256 Sha 算签名
func SHA256(text string) string {
	hash := sha256.New()
	text = salt + text + salt
	hash.Write([]byte(text))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
