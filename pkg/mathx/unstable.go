/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-06 09:15:55
 * @FilePath: \go-toolbox\pkg\mathx\unstable.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"math/rand"
	"reflect"
	"sync"
	"time"
)

// Unstable 结构体用于基于给定的偏差值生成围绕均值附近的随机值。
type Unstable struct {
	deviation float64     // 偏差值
	r         *rand.Rand  // 随机数生成器
	lock      *sync.Mutex // 互斥锁，用于并发安全
}

// NewUnstable 创建一个新的 Unstable 实例。
func NewUnstable(deviation float64) Unstable {
	// 确保偏差值在合理范围内
	if deviation < 0 {
		deviation = 0
	}
	if deviation > 1 {
		deviation = 1
	}
	return Unstable{
		deviation: deviation,
		r:         rand.New(rand.NewSource(time.Now().UnixNano())), // 使用当前时间的纳秒数作为随机数种子
		lock:      new(sync.Mutex),                                 // 初始化互斥锁
	}
}

// AroundDuration 根据给定的基础时长和偏差值返回一个随机的时长。
func (u Unstable) AroundDuration(base time.Duration) time.Duration {
	u.lock.Lock() // 加锁以确保并发安全
	// 根据公式计算随机值，公式为：(1 + deviation - 2*deviation*随机数) * 基础值
	val := time.Duration((1 + u.deviation - 2*u.deviation*u.r.Float64()) * float64(base))
	u.lock.Unlock() // 解锁
	return val
}

// AroundInt 根据给定的基础整数值和偏差值返回一个随机的 int64 值。
func (u Unstable) AroundInt(base int64) int64 {
	u.lock.Lock() // 加锁以确保并发安全
	// 根据公式计算随机值，公式为：(1 + deviation - 2*deviation*随机数) * 基础值
	val := int64((1 + u.deviation - 2*u.deviation*u.r.Float64()) * float64(base))
	u.lock.Unlock() // 解锁
	return val
}

// Clone 深拷贝任意类型的值
func Clone(value interface{}, seen map[uintptr]interface{}) interface{} {
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)

	// 检查是否已经克隆过
	if v.Kind() == reflect.Ptr {
		ptrAddr := v.Pointer()
		if existing, found := seen[ptrAddr]; found {
			return existing
		}
	}

	// 创建一个新的映射以存储克隆的对象
	if seen == nil {
		seen = make(map[uintptr]interface{})
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		newValue := reflect.New(v.Elem().Type())
		seen[v.Pointer()] = newValue.Interface()
		newValue.Elem().Set(reflect.ValueOf(Clone(v.Elem().Interface(), seen)))
		return newValue.Interface()

	case reflect.Slice:
		clone := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			clone.Index(i).Set(reflect.ValueOf(Clone(v.Index(i).Interface(), seen)))
		}
		return clone.Interface()

	case reflect.Map:
		clone := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			clone.SetMapIndex(reflect.ValueOf(Clone(key.Interface(), seen)), reflect.ValueOf(Clone(value.Interface(), seen)))
		}
		return clone.Interface()

	case reflect.Struct:
		clone := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			clone.Field(i).Set(reflect.ValueOf(Clone(field.Interface(), seen)))
		}
		return clone.Interface()

	default:
		return value
	}
}
