/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-05-27 18:51:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-27 19:05:28
 * @FilePath: \go-toolbox\pkg\convert\reflect.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"fmt"
	"reflect"
	"time"
)

type TransformFieldsOptions struct {
	StrictTypeCheck bool   // 是否严格类型匹配，类型不符时panic
	TimeFormat      string // 时间转换格式
}

func TransformFields(dst any, src any, opts *TransformFieldsOptions) {
	if opts == nil {
		opts = &TransformFieldsOptions{}
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.DateTime
	}
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		panic("dst must be a non-nil pointer")
	}
	dstVal = dstVal.Elem()
	srcVal := reflect.ValueOf(src)
	transformValue(dstVal, srcVal, opts)
}

var (
	timeType    = reflect.TypeOf(time.Time{})
	timePtrType = reflect.TypeOf(&time.Time{})
)

func transformValue(dst, src reflect.Value, opts *TransformFieldsOptions) {
	if !src.IsValid() {
		return
	}

	switch dst.Kind() {
	case reflect.Struct:
		if src.Kind() == reflect.Struct {
			for i := 0; i < dst.NumField(); i++ {
				dstField := dst.Field(i)
				dstFieldType := dst.Type().Field(i)
				if !dstField.CanSet() || dstFieldType.PkgPath != "" {
					continue
				}
				srcField := src.FieldByName(dstFieldType.Name)
				if !srcField.IsValid() {
					continue
				}
				transformValue(dstField, srcField, opts)
			}
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to %s", src.Type(), dst.Type()))
		}
	case reflect.String:
		if src.Type() == timeType {
			t := src.Interface().(time.Time)
			dst.SetString(t.Format(opts.TimeFormat))
			return
		}
		if src.Type() == timePtrType {
			if src.IsNil() {
				dst.SetString("")
				return
			}
			t := src.Elem().Interface().(time.Time)
			dst.SetString(t.Format(opts.TimeFormat))
			return
		}
		if src.Kind() == reflect.String {
			dst.SetString(src.String())
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to string", src.Type()))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if src.Kind() >= reflect.Int && src.Kind() <= reflect.Int64 {
			dst.SetInt(src.Int())
			return
		}
		if src.Kind() >= reflect.Uint && src.Kind() <= reflect.Uint64 {
			dst.SetInt(int64(src.Uint()))
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to %s", src.Type(), dst.Type()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if src.Kind() >= reflect.Uint && src.Kind() <= reflect.Uint64 {
			dst.SetUint(src.Uint())
			return
		}
		if src.Kind() >= reflect.Int && src.Kind() <= reflect.Int64 {
			v := src.Int()
			if v < 0 {
				if opts.StrictTypeCheck {
					panic(fmt.Sprintf("negative int %d cannot assign to %s", v, dst.Type()))
				}
				return
			}
			dst.SetUint(uint64(v))
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to %s", src.Type(), dst.Type()))
		}
	case reflect.Float32, reflect.Float64:
		if src.Kind() == reflect.Float32 || src.Kind() == reflect.Float64 {
			dst.SetFloat(src.Float())
			return
		}
		if src.Kind() >= reflect.Int && src.Kind() <= reflect.Int64 {
			dst.SetFloat(float64(src.Int()))
			return
		}
		if src.Kind() >= reflect.Uint && src.Kind() <= reflect.Uint64 {
			dst.SetFloat(float64(src.Uint()))
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to %s", src.Type(), dst.Type()))
		}
	case reflect.Bool:
		if src.Kind() == reflect.Bool {
			dst.SetBool(src.Bool())
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to bool", src.Type()))
		}
	case reflect.Ptr:
		if src.Kind() == reflect.Ptr {
			if src.IsNil() {
				dst.Set(reflect.Zero(dst.Type()))
				return
			}
			if dst.IsNil() {
				dst.Set(reflect.New(dst.Type().Elem()))
			}
			transformValue(dst.Elem(), src.Elem(), opts)
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to %s", src.Type(), dst.Type()))
		}
	case reflect.Slice:
		if src.Kind() == reflect.Slice {
			newSlice := reflect.MakeSlice(dst.Type(), src.Len(), src.Cap())
			for i := 0; i < src.Len(); i++ {
				transformValue(newSlice.Index(i), src.Index(i), opts)
			}
			dst.Set(newSlice)
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to slice %s", src.Type(), dst.Type()))
		}
	case reflect.Map:
		if src.Kind() == reflect.Map {
			if src.IsNil() {
				dst.Set(reflect.Zero(dst.Type()))
				return
			}
			newMap := reflect.MakeMapWithSize(dst.Type(), src.Len())
			for _, key := range src.MapKeys() {
				valSrc := src.MapIndex(key)
				valDst := reflect.New(dst.Type().Elem()).Elem()
				transformValue(valDst, valSrc, opts)
				newMap.SetMapIndex(key, valDst)
			}
			dst.Set(newMap)
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("type mismatch: cannot assign %s to map %s", src.Type(), dst.Type()))
		}
	default:
		if src.Type().AssignableTo(dst.Type()) {
			dst.Set(src)
			return
		}
		if opts.StrictTypeCheck {
			panic(fmt.Sprintf("unsupported or mismatched type: %s -> %s", src.Type(), dst.Type()))
		}
	}
}
