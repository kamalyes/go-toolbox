/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-12 22:26:27
 * @FilePath: \go-toolbox\pkg\convert\must.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
	default:
		// 对于未知类型，使用 %v 格式化为默认字符串表示
		return fmt.Sprintf("%v", val)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "" // 返回空字符串以避免 panic
	}
	return string(b)
}

// Number 是一个接口，限制了可以作为类型参数的数值类型
type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// MustIntT 将 any 转换为 T 类型
func MustIntT[T Number](value any) (T, error) {
	const exceedsIntRange = "value exceeds int range"
	const unsupportedConversion = "unsupported conversion"
	switch v := value.(type) {
	case int:
		return T(v), nil
	case int8:
		return T(v), nil
	case int16:
		return T(v), nil
	case int32:
		return T(v), nil
	case int64:
		return T(v), nil
	case uint:
		return T(v), nil
	case uint8:
		return T(v), nil
	case uint16:
		return T(v), nil
	case uint32:
		return T(v), nil
	case uint64:
		return T(v), nil
	default:
		var zero T
		return zero, fmt.Errorf("%s: %v (type %T)", unsupportedConversion, value, value)
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
		flag, err := MustIntT[int](v)
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
