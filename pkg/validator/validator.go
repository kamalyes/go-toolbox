/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 19:55:57
 * @FilePath: \go-toolbox\pkg\validator\validator.go
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

var customRules = make(map[string]Rules)

// RegisterRule registers a custom validation rule.
func RegisterRule(key string, rule Rules) error {
	if _, exists := customRules[key]; exists {
		return errors.New(key + "已注册,无法重复注册")
	}
	customRules[key] = rule
	return nil
}

// isEmptyValue checks if a reflect.Value is empty.
func IsEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
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
		return v.IsNil() || IsEmptyValue(v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !IsEmptyValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// HasEmpty checks if any element in the slice is empty.
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

// IsAllEmpty checks if all elements in the slice are empty.
func IsAllEmpty(elems []interface{}) bool {
	for _, elem := range elems {
		if !IsEmptyValue(reflect.ValueOf(elem)) {
			return false
		}
	}
	return true
}

// IsUndefined checks if a string is "undefined" (case insensitive).
func IsUndefined(str string) bool {
	return strings.EqualFold(strings.TrimSpace(str), "undefined")
}

// ContainsChinese checks if a string contains any Chinese characters.
func ContainsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// EmptyToDefault returns defaultStr if str is empty; otherwise, returns str.
func EmptyToDefault(str string, defaultStr string) string {
	if IsEmptyValue(reflect.ValueOf(str)) {
		return defaultStr
	}
	return str
}

// NotEmpty returns the rule string for non-empty validation.
func NotEmpty() string {
	return "notEmpty"
}

// Verify validates the struct fields based on the provided rules.
func Verify(st interface{}, roleMap Rules) error {
	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st)

	if val.Kind() != reflect.Struct {
		return errors.New("expect struct")
	}

	for i := 0; i < val.NumField(); i++ {
		tagVal := typ.Field(i)
		fieldVal := val.Field(i)
		if validations, exists := roleMap[tagVal.Name]; exists {
			for _, v := range validations {
				if err := validateField(tagVal.Name, fieldVal, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// validateField checks a single field against its validation rules.
func validateField(fieldName string, fieldVal reflect.Value, rule string) error {
	switch {
	case rule == NotEmpty():
		if isBlank(fieldVal) {
			return errors.New(fieldName + "值不能为空")
		}
	case strings.HasPrefix(rule, "regexp="):
		if !regexpMatch(strings.TrimPrefix(rule, "regexp="), fieldVal.String()) {
			return errors.New(fieldName + "格式校验不通过")
		}
	default:
		if !compareVerify(fieldVal, rule) {
			return errors.New(fieldName + "长度或值不在合法范围," + rule)
		}
	}
	return nil
}

// compareVerify checks if a value satisfies a comparison rule.
func compareVerify(value reflect.Value, verifyStr string) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		return compare(value.Len(), verifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), verifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), verifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), verifyStr)
	default:
		return false
	}
}

// isBlank checks if a value is blank (zero value).
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

// compare compares a value against a rule.
func compare(value interface{}, verifyStr string) bool {
	parts := strings.Split(verifyStr, "=")
	val := reflect.ValueOf(value)

	if len(parts) != 2 {
		return false
	}

	op, strValue := parts[0], parts[1]
	var compareValue float64
	var err error

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		compareValue, err = strconv.ParseFloat(strconv.FormatInt(val.Int(), 10), 64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		compareValue, err = strconv.ParseFloat(strconv.FormatUint(val.Uint(), 10), 64)
	case reflect.Float32, reflect.Float64:
		compareValue, err = val.Float(), nil
	default:
		return false
	}

	if err != nil {
		return false
	}

	ruleValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return false
	}

	switch op {
	case "lt":
		return compareValue < ruleValue
	case "le":
		return compareValue <= ruleValue
	case "eq":
		return compareValue == ruleValue
	case "ne":
		return compareValue != ruleValue
	case "ge":
		return compareValue >= ruleValue
	case "gt":
		return compareValue > ruleValue
	default:
		return false
	}
}

// regexpMatch checks if a string matches a regex pattern.
func regexpMatch(pattern, matchStr string) bool {
	return regexp.MustCompile(pattern).MatchString(matchStr)
}
