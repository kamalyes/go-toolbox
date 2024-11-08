/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 01:15:15
 * @FilePath: \go-toolbox\pkg\convert\must.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

// MustString 强制转为字符串
func MustString[T any](v T, timeLayout ...string) string {
	switch v := any(v).(type) {
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
	default:
		return convertToString(v, timeLayout...)
	}
}

// convertToString 将其他类型转换为字符串
func convertToString[T any](v T, timeLayout ...string) string {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case reflect.Struct:
		if val.Type() == reflect.TypeOf(time.Time{}) {
			t := val.Interface().(time.Time) // 这里可以安全地断言为 time.Time
			if len(timeLayout) > 0 {
				return t.Format(timeLayout[0])
			}
			return t.Format(time.RFC3339)
		}
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "" // 返回空字符串以避免 panic
	}
	return string(b)
}

// MustInt 强制转换为整数 (int)
func MustInt[T any](v T) (int, error) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String:
		return parseStringToInt(val.String())
	case reflect.Bool:
		if val.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Invalid:
		return 0, nil
	default:
		return convertToInt(val)
	}
}

// parseStringToInt 解析字符串为整数
func parseStringToInt(s string) (int, error) {
	i, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, errors.New("invalid string to int conversion")
	}
	return i, nil
}

// convertToInt 将其他类型转换为 int
func convertToInt(val reflect.Value) (int, error) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int(val.Float()), nil
	default:
		return 0, errors.New("unsupported type for conversion to int")
	}
}

// MustBool 强制转为 bool
func MustBool[T any](v T) bool {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Bool:
		return val.Bool()
	case reflect.String:
		return stringx.IsTrueString(val.String())
	default:
		flag, err := MustInt(v)
		return err == nil && flag != 0
	}
}

// MustJSONIndent 转 json 返回 []byte
func MustJSONIndent(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// MustJSON 转 json 返回 []byte
func MustJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
