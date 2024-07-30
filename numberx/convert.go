/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 17:50:03
 * @FilePath: \go-toolbox\numberx\convert.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package numberx

import (
	"fmt"
	"strconv"
)

func ToInt(i interface{}) (num int, err error) {
	switch v := i.(type) {
	case int: // 如果i已经是int类型，直接返回
		return v, nil
	case string: // 如果i是string类型，尝试转换为int
		num, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string '%s' to int: %v", v, err)
		}
		return num, nil
	default:
		return 0, fmt.Errorf("unsupported type %T", i)
	}
}

// ToInt32
/**
 * @Description: 将任意类型转为int32类型
 * @param i
 * @return num
 * @return err
 */
func ToInt32(i interface{}) (num int32, err error) {
	switch v := i.(type) {
	case int32:
		num = v
	case int:
		num = int32(v)
	case float32:
		num = int32(v)
	case float64:
		num = int32(v)
	case string:
		n, e := strconv.ParseInt(v, 10, 32)
		if e != nil {
			err = e
			return
		}
		num = int32(n)
	default:
		err = fmt.Errorf("unable to convert %T to int32", i)
	}
	return
}

// ToInt64
/**
 * @Description: 将任意类型转为int64类型
 * @param i
 * @return num
 * @return err
 */
func ToInt64(i interface{}) (num int64, err error) {
	switch v := i.(type) {
	case int64:
		num = v
	case int:
		num = int64(v)
	case float32:
		num = int64(v)
	case float64:
		num = int64(v)
	case string:
		n, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			err = e
			return
		}
		num = int64(n)
	default:
		err = fmt.Errorf("unable to convert %T to int64", i)
	}
	return
}

// ToFloat32
/**
 * @Description: 将任意类型转为float32类型
 * @param i
 * @return num
 * @return err
 */
func ToFloat32(i interface{}) (num float32, err error) {
	switch v := i.(type) {
	case float32:
		num = v
	case string:
		// string无法直接转换float32，只能先转换为float64，再通过float64转float32
		var num64 float64
		num64, err = strconv.ParseFloat(v, 32)
		num = float32(num64)
	case int:
		num = float32(v)
	case int32:
		num = float32(v)
	case int64:
		num = float32(v)
	case float64:
		// 可能造成精度丢失
		num = float32(v)
	default:
		err = fmt.Errorf("unable to convert %T to float32", i)
	}
	return
}

// ToFloat64
/**
 * @Description: 将任意类型转为float64类型
 * @param i
 * @return num
 * @return err
 */
func ToFloat64(i interface{}) (num float64, err error) {
	switch v := i.(type) {
	case float64:
		num = v
	case string:
		num, err = strconv.ParseFloat(v, 64)
	case int:
		num = float64(v)
	case int32:
		num = float64(v)
	case int64:
		num = float64(v)
	case float32:
		num = float64(v)
	default:
		err = fmt.Errorf("unable to convert %T to float64", i)
	}
	return
}
