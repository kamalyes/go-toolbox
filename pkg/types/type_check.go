/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-13 13:27:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 13:51:15
 * @FilePath: \go-toolbox\pkg\types\type_check.go
 * @Description: 检查类型兼容性
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package types

import (
	"fmt"
	"reflect"
	"time"
)

var (
	ErrTypeMismatchStrict = "type mismatch: cannot assign %s to %s"
)

// CheckTypeCompatibility 检查类型兼容性
func CheckTypeCompatibility(srcType, dstType reflect.Type) error {
	if srcType == nil || dstType == nil {
		return fmt.Errorf(ErrTypeMismatchStrict, srcType, dstType)
	}

	for srcType.Kind() == reflect.Ptr {
		srcType = srcType.Elem()
	}
	for dstType.Kind() == reflect.Ptr {
		dstType = dstType.Elem()
	}

	if dstType.Kind() == reflect.Interface && dstType.NumMethod() == 0 {
		return nil
	}
	if srcType.Kind() == reflect.Interface && srcType.NumMethod() == 0 {
		return nil
	}

	if srcType == reflect.TypeOf(time.Time{}) && dstType.Kind() == reflect.String {
		return nil
	}

	if srcType.Kind() != dstType.Kind() {
		return fmt.Errorf(ErrTypeMismatchStrict, srcType, dstType)
	}

	switch srcType.Kind() {
	case reflect.Struct:
		dstNum := dstType.NumField()
		for i := 0; i < dstNum; i++ {
			dstField := dstType.Field(i)
			if dstField.PkgPath != "" {
				continue
			}
			srcField, ok := srcType.FieldByName(dstField.Name)
			if !ok {
				continue
			}
			if srcField.PkgPath != "" {
				continue
			}
			if err := CheckTypeCompatibility(srcField.Type, dstField.Type); err != nil {
				return err
			}
		}
		return nil

	case reflect.Slice, reflect.Array:
		return CheckTypeCompatibility(srcType.Elem(), dstType.Elem())

	case reflect.Map:
		if err := CheckTypeCompatibility(srcType.Key(), dstType.Key()); err != nil {
			return err
		}
		return CheckTypeCompatibility(srcType.Elem(), dstType.Elem())

	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return fmt.Errorf("unsupported type in strict mode: %s", srcType.Kind())

	default:
		if srcType != dstType {
			return fmt.Errorf(ErrTypeMismatchStrict, srcType, dstType)
		}
		return nil
	}
}
