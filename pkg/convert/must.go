/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-27 18:09:55
 * @FilePath: \go-toolbox\pkg\convert\must.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/kamalyes/go-toolbox/pkg/types"
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

// RoundMode 是一个枚举类型，用于指定取整的方式
type RoundMode int

const (
	RoundNone    RoundMode = iota // 不进行四舍五入，保持原值
	RoundNearest                  // 四舍五入到最接近的整数
	RoundDown                     // 向下取整
	RoundUp                       // 向上取整
)

// MustFloatT 将字符串转换为指定类型的浮点数
func MustFloatT[T types.Float](value any, mode RoundMode) (T, error) {
	var f float64

	// 根据输入的类型进行处理
	switch v := value.(type) {
	case string:
		var err error
		f, err = strconv.ParseFloat(v, 64) // 解析为 float64
		if err != nil {
			return 0, fmt.Errorf("无法将字符串转换为浮点数: %v", err)
		}
	case float64:
		f = v // 直接使用 float64 值
	case float32:
		f = float64(v) // 将 float32 转换为 float64
	default:
		return 0, fmt.Errorf("不支持的输入类型: %T", v)
	}

	var result float64

	// 根据取整模式进行处理
	switch mode {
	case RoundNone:
		result = f // 保持原值不变
	case RoundNearest:
		result = math.Round(f) // 使用 math.Round 处理四舍五入
	case RoundUp:
		result = math.Ceil(f) // 向上取整
	case RoundDown:
		result = math.Floor(f) // 向下取整
	default:
		return 0, fmt.Errorf("未知的四舍五入模式")
	}

	return T(result), nil // 将结果转换为 T 类型并返回
}

// MustIntT 将 any 转换为 T 类型
func MustIntT[T types.Numerical](value any, mode *RoundMode) (T, error) {
	const unsupportedConversion = "unsupported conversion"
	// 默认取整模式为向下取整
	if mode == nil {
		defaultMode := RoundDown
		mode = &defaultMode
	}
	var zero T
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
	case float32:
		if *mode == RoundUp {
			return T(math.Ceil(float64(v))), nil
		}
		return T(math.Floor(float64(v))), nil
	case float64:
		if *mode == RoundUp {
			return T(math.Ceil(v)), nil
		}
		return T(math.Floor(v)), nil
	case string:
		// 需要特殊处理下, 坑注意go版本不一致 越界返回的结果不一样
		// GO WIN 1.21.13 input := []string{"9223372036854775807", "9223372036854775806"} actual  : []int64{-9223372036854775808, -9223372036854775808}
		// GO LINUX 1.21.13 input := []string{"9223372036854775807", "9223372036854775806"} actual  : []int64{9223372036854775807, 9223372036854775807}
		var floatValue float64
		err := ParseFloat(v, &floatValue) // 尝试将字符串解析为浮点数
		if err != nil {
			return zero, fmt.Errorf("failed to parse %q: %v", v, err)
		}
		return Float64ToInt[T](floatValue, *mode)
	default:
		return zero, fmt.Errorf("%s: %v (type %T)", unsupportedConversion, value, value)
	}
}

// Float64ToInt 将浮点数转换为整数类型，并进行取整
func Float64ToInt[T types.Numerical](value float64, mode RoundMode) (T, error) {
	var resultFloatValue T

	var convertedValue float64
	if mode == RoundUp {
		convertedValue = math.Ceil(value)
	} else {
		convertedValue = math.Floor(value)
	}

	// 检查转换后的值是否超出 T 的范围
	switch any(convertedValue).(type) {
	case int64:
		if convertedValue < float64(math.MinInt64) || convertedValue > float64(math.MaxInt64) {
			return resultFloatValue, fmt.Errorf("value %f out of range for type %T", convertedValue, resultFloatValue)
		}
	}

	resultFloatValue = T(convertedValue)
	return resultFloatValue, nil
}

// ParseFloat 尝试将字符串解析为指定类型的浮点数
func ParseFloat[T types.Float](v string, value *T) error {
	f, err := strconv.ParseFloat(v, 64) // 先解析为 float64
	if err != nil {
		return fmt.Errorf("failed to parse %q: %v", v, err)
	}
	// 检查是否是 NaN 或无穷大
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Errorf("invalid float value: %q", v)
	}
	*value = T(f) // 转换为目标类型
	return nil
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
		flag, err := MustIntT[int](v, nil)
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

// NumberSliceToStringSlice Number切片转String
func NumberSliceToStringSlice[T types.Numerical](numbers []T) []string {
	if numbers == nil {
		return nil // 处理 nil 切片
	}
	var result []string
	for _, number := range numbers {
		result = append(result, fmt.Sprintf("%v", number)) // 使用 %v 格式化输出
	}
	return result
}

// StringSliceToNumberSlice 将字符串切片转换为数字切片
func StringSliceToNumberSlice[T types.Numerical](input []string, mode *RoundMode) ([]T, error) {
	if input == nil {
		return []T{}, nil // 返回一个空切片而不是 nil
	}
	result := make([]T, 0, len(input))
	for _, str := range input {
		num, err := MustIntT[T](str, mode) // 使用 MustIntT 进行转换
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}

	return result, nil
}

// StringSliceToFloatSlice 将字符串切片转换为浮点数切片
func StringSliceToFloatSlice[T types.Float](input []string, mode RoundMode) ([]T, error) {
	if input == nil {
		return []T{}, nil // 返回一个空切片而不是 nil
	}
	result := make([]T, 0, len(input))
	for _, str := range input {
		num, err := MustFloatT[T](str, mode) // 使用 MustFloatT 进行转换
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}

	return result, nil
}

// StringSliceToInterfaceSlice 将 []string 转换为 []interface{}
func StringSliceToInterfaceSlice(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

var defaultRoundMode = RoundNone

// ToNumberSlice 支持输入是 string 或 []string
func ToNumberSlice[T types.Numerical](input any, decollator string) ([]T, error) {
	var strSlice []string

	switch v := input.(type) {
	case string:
		if v == "" {
			return []T{}, nil // 空字符串直接返回空切片，不继续拆分
		}
		// 自动拆分字符串
		strSlice = strings.Split(v, decollator)
	case []string:
		strSlice = v
	default:
		return nil, fmt.Errorf("unsupported input type %T, want string or []string", input)
	}

	result := make([]T, 0, len(strSlice))
	for _, s := range strSlice {
		s = strings.TrimSpace(s) // 去掉前后空格
		num, err := MustIntT[T](s, &defaultRoundMode)
		if err != nil {
			return nil, fmt.Errorf("转换失败，input=%q, err=%w", s, err)
		}
		result = append(result, num)
	}
	return result, nil
}

// MustToNumberSlice 是不返回错误版本，遇错 panic
func MustToNumberSlice[T types.Numerical](input any, decollator string) []T {
	nums, err := ToNumberSlice[T](input, decollator)
	if err != nil {
		panic(err)
	}
	return nums
}
