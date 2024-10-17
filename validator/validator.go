/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 16:56:00
 * @FilePath: \go-toolbox\validator\validator.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Rules map[string][]string

type RulesMap map[string]Rules

var CustomizeMap = make(map[string]Rules)

// RegisterRule 注册自定义规则方案建议在路由初始化层即注册
func RegisterRule(key string, rule Rules) (err error) {
	if CustomizeMap[key] != nil {
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CustomizeMap[key] = rule
		return nil
	}
}

// IsEmptyValue 是否空对象
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0.0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return true
		}
		return IsEmptyValue(v.Elem())
	default:
		return false
	}
}

// HasEmpty 判断提供的接口数组是否包含空值，并返回空值数量
func HasEmpty(elems []interface{}) (bool, int) {
	if len(elems) == 0 {
		return true, 0
	}

	emptyCount := 0

	for _, elem := range elems {
		if IsEmptyValue(reflect.ValueOf(elem)) {
			emptyCount++
		}
	}

	return emptyCount > 0, emptyCount
}

// IsAllEmpty 判断提供的接口数组是否全是空值
func IsAllEmpty(elems []interface{}) bool {
	if len(elems) == 0 {
		return true
	}

	for _, elem := range elems {
		if !IsEmptyValue(reflect.ValueOf(elem)) {
			return false
		}
	}

	return true
}

// IsUndefined 检查字符串是否等于 "undefined"（不区分大小写，忽略前后空格）
func IsUndefined(str string) bool {
	trimmedStr := strings.TrimSpace(str)
	return strings.EqualFold(trimmedStr, "undefined")
}

// 校验字符串是否包含中文字符
func ContainsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// EmptyToDefault 如果字符串是""，则返回指定默认字符串，否则返回字符串本身。
func EmptyToDefault(str string, defaultStr string) string {
	if IsEmptyValue(reflect.ValueOf(str)) {
		return defaultStr
	} else {
		return str
	}
}

// NotEmpty 非空 不能为其对应类型的0值
func NotEmpty() string {
	return "notEmpty"
}

// RegexpMatch 正则校验 校验输入项是否满足正则表达式
func RegexpMatch(rule string) string {
	return "regexp=" + rule
}

// Lt 小于入参(<) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Lt(mark string) string {
	return "lt=" + mark
}

// Le 小于等于入参(<=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Le(mark string) string {
	return "le=" + mark
}

// Eq 等于入参(==) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Eq(mark string) string {
	return "eq=" + mark
}

// Ne 不等于入参(!=)  如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Ne(mark string) string {
	return "ne=" + mark
}

// Ge 大于等于入参(>=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Ge(mark string) string {
	return "ge=" + mark
}

// Gt 大于入参(>) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Gt(mark string) string {
	return "gt=" + mark
}

// Verify 校验方法
func Verify(st interface{}, roleMap Rules) (err error) {
	compareMap := map[string]bool{
		"lt": true,
		"le": true,
		"eq": true,
		"ne": true,
		"ge": true,
		"gt": true,
	}

	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st) // 获取reflect.Type类型

	kd := val.Kind() // 获取到st对应的类别
	if kd != reflect.Struct {
		return errors.New("expect struct")
	}
	num := val.NumField()
	// 遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		val := val.Field(i)
		if len(roleMap[tagVal.Name]) > 0 {
			for _, v := range roleMap[tagVal.Name] {
				switch {
				case v == "notEmpty":
					if isBlank(val) {
						return errors.New(tagVal.Name + "值不能为空")
					}
				case strings.Split(v, "=")[0] == "regexp":
					if !regexpMatch(strings.Split(v, "=")[1], val.String()) {
						return errors.New(tagVal.Name + "格式校验不通过")
					}
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(val, v) {
						return errors.New(tagVal.Name + "长度或值不在合法范围," + v)
					}
				}
			}
		}
	}
	return nil
}

// 长度和数字的校验方法 根据类型自动校验
func compareVerify(value reflect.Value, VerifyStr string) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}
}

// 非空校验
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// 比较函数
func compare(value interface{}, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < VInt
		case VerifyStrArr[0] == "le":
			return val.Int() <= VInt
		case VerifyStrArr[0] == "eq":
			return val.Int() == VInt
		case VerifyStrArr[0] == "ne":
			return val.Int() != VInt
		case VerifyStrArr[0] == "ge":
			return val.Int() >= VInt
		case VerifyStrArr[0] == "gt":
			return val.Int() > VInt
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(VInt)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(VInt)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(VInt)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(VInt)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(VInt)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(VInt)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	default:
		return false
	}
}

// 正则匹配规则
func regexpMatch(rule, matchStr string) bool {
	return regexp.MustCompile(rule).MatchString(matchStr)
}
