/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:15:15
 * @FilePath: \go-toolbox\pkg\mathx\proba.go
 * @Description:
 * 此文件定义了一个Proba结构体和相关的方法，用于根据给定的概率判断是否为真。
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"math/rand"
	"sync"
	"time"
)

// Proba结构体用于根据给定的概率判断事件是否为真。
// 它内部使用了一个线程不安全的随机数生成器，并通过互斥锁来保证并发安全。
type Proba struct {
	// r 是一个随机数生成器，使用当前时间的纳秒数作为种子来初始化，以保证每次运行时的随机性。
	// 注意：rand.New(...) 返回的随机数生成器本身不是线程安全的。
	r *rand.Rand
	// lock 是一个互斥锁，用于在并发访问随机数生成器时保证线程安全。
	lock sync.Mutex
}

// NewProba函数用于创建一个新的Proba实例。
// 它返回一个初始化了随机数生成器和互斥锁的Proba结构体指针。
func NewProba() *Proba {
	// 使用当前时间的纳秒数作为种子来创建一个新的随机数源，并基于此源创建一个随机数生成器。
	return &Proba{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TrueOnProba方法用于根据给定的概率判断事件是否为真。
// 它接收一个浮点数作为概率参数（范围应在0到1之间），并返回一个布尔值。
// 如果生成的随机数小于给定的概率，则返回true；否则返回false。
// 此方法通过加锁来保证在并发访问时的线程安全。
func (p *Proba) TrueOnProba(proba float64) bool {
	// 加锁以保证在并发访问随机数生成器时的线程安全。
	p.lock.Lock()
	defer p.lock.Unlock() // 确保在函数返回前解锁，以防止死锁。

	// 使用随机数生成器生成一个0到1之间的浮点数，并与给定的概率进行比较。
	// 如果生成的随机数小于给定的概率，则返回true；否则返回false。
	truth := p.r.Float64() < proba
	return truth
}
