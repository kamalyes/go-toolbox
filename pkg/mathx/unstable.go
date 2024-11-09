/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 15:19:08
 * @FilePath: \go-toolbox\pkg\mathx\unstable.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"math/rand"
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
