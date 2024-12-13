/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:05:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-13 13:05:15
 * @FilePath: \go-toolbox\pkg\syncx\lock.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

// Locker 是一个接口，定义了锁定和解锁的方法
type Locker interface {
	Lock()
	Unlock()
}

// WithLock 是一个通用的函数，用于在给定的锁上执行操作
func WithLock(lock Locker, operation func()) {
	lock.Lock()         // 获取锁
	defer lock.Unlock() // 确保在操作完成后释放锁
	operation()         // 执行操作
}
