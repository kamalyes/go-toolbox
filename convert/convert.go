/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 23:19:16
 * @FilePath: \go-toolbox\convert\convert.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// MustJSONIndent 转 json 返回 []byte
func MustJSONIndent(v interface{}) []byte {
	js, _ := json.MarshalIndent(v, "", "  ")
	return js
}

// MustJSONIndentString 转 json Indent 返回 string
func MustJSONIndentString(v interface{}) string {
	return string(MustJSONIndent(v))
}

// MustJSON 转 json 返回 []byte
func MustJSON(v interface{}) []byte {
	js, _ := json.Marshal(v)
	return js
}

// MustJSONString 转 json 返回 string
func MustJSONString(v interface{}) string {
	return string(MustJSON(v))
}

// MustString 强制转为字符串
func MustString(v interface{}, timeLayout ...string) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case nil:
		return ""
	case bool:
		return strconv.FormatBool(v)
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if len(timeLayout) > 0 {
			return v.Format(timeLayout[0])
		}
		return v.Format(time.RFC3339)
	case fmt.Stringer:
		return v.String()
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// MustInt 强制转换为整数 (int)
func MustInt(v interface{}) (int, error) {
	switch v := v.(type) {
	case string:
		if i, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return i, nil
		}
		return 0, errors.New("invalid string to int conversion")
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return convertToInt(v)
	}
}

// convertToInt 将其他类型转换为 int
func convertToInt(v interface{}) (int, error) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int:
		return int(val.Int()), nil
	case reflect.Int8:
		return int(val.Int()), nil
	case reflect.Int16:
		return int(val.Int()), nil
	case reflect.Int32:
		return int(val.Int()), nil
	case reflect.Int64:
		return int(val.Int()), nil
	case reflect.Uint:
		return int(val.Uint()), nil
	case reflect.Uint8:
		return int(val.Uint()), nil
	case reflect.Uint16:
		return int(val.Uint()), nil
	case reflect.Uint32:
		return int(val.Uint()), nil
	case reflect.Uint64:
		return int(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int(val.Float()), nil
	default:
		return 0, errors.New("unsupported type for conversion to int")
	}
}

// MustBool 强制转为 bool
func MustBool(v interface{}) bool {
	switch v := v.(type) {
	case bool:
		return v
	case string:
		switch v {
		case "1", "t", "T", "true", "TRUE", "True":
			return true
		default:
			return false
		}
	default:
		flag, err := MustInt(v)
		if err != nil {
			return false
		}
		return flag != 0
	}
}

// B64Encode Base64 编码
func B64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// B64Decode Base64 解码
func B64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// B64UrlEncode Base64 URL 安全编码
func B64UrlEncode(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

// B64UrlDecode Base64 URL 安全解码
func B64UrlDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

// BytesToHex 字节数组转为十六进制字符串
func BytesToHex(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

// HexToBytes 将十六进制字符串转换为字节数组
func HexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		return nil, errors.New("hex string must have an even length")
	}
	return hex.DecodeString(hexStr)
}

// BytesBCC 计算字节数组的 BCC
func BytesBCC(data []byte) byte {
	var bcc byte
	for _, b := range data {
		bcc ^= b
	}
	return bcc
}

// HexBCC 计算十六进制字符串的 BCC
func HexBCC(hexStr string) (string, error) {
	bytes, err := HexToBytes(hexStr)
	if err != nil {
		return "", err
	}
	bcc := BytesBCC(bytes)
	return hex.EncodeToString([]byte{bcc}), nil
}

// DecToHex 十进制转为十六进制字符串
func DecToHex(n uint64) string {
	s := strconv.FormatUint(n, 16)
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return strings.ToUpper(s)
}

// HexToDec 十六进制字符串转为十进制
func HexToDec(h string) (uint64, error) {
	return strconv.ParseUint(h, 16, 64)
}

// DecToBin 十进制转为二进制字符串，补齐到8位
func DecToBin(n uint64) string {
	return fmt.Sprintf("%08b", n)
}

// HexToBin 十六进制字符串转为二进制字符串
func HexToBin(h string) (string, error) {
	n, err := HexToDec(h)
	if err != nil {
		return "", err
	}
	return DecToBin(n), nil
}

// ByteToBinStr 将单个字节转为二进制字符串
func ByteToBinStr(b byte) string {
	return fmt.Sprintf("%08b", b)
}

// BytesToBinStr 将字节数组转为二进制字符串
func BytesToBinStr(bs []byte) string {
	var buf bytes.Buffer
	for _, v := range bs {
		buf.WriteString(ByteToBinStr(v))
	}
	return buf.String()
}

// BytesToBinStrWithSplit 将字节数组转为二进制字符串，并添加分隔符
func BytesToBinStrWithSplit(bs []byte, split string) string {
	var buf bytes.Buffer
	for i, v := range bs {
		if i > 0 {
			buf.WriteString(split)
		}
		buf.WriteString(ByteToBinStr(v))
	}
	return buf.String()
}

// Base64ToByte 将 Base64 字符串解码为字节切片
// 参数：imageBase64 - 要解码的 Base64 字符串
// 返回：解码后的字节切片和可能的错误
func Base64ToByte(imageBase64 string) ([]byte, error) {
	// 解码 Base64 字符串
	image, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return nil, err // 返回错误
	}

	return image, nil // 返回解码后的字节
}
