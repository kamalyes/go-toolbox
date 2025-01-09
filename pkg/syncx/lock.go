/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:05:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-08 15:15:59
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

// RLocker 是一个接口，定义了读锁和解锁的方法
type RLocker interface {
	RLock()
	RUnlock()
}

// WithLock 是一个通用的函数，用于在给定的锁上执行操作
func WithLock(lock Locker, operation func()) {
	lock.Lock()         // 获取锁
	defer lock.Unlock() // 确保在操作完成后释放锁
	operation()         // 执行操作
}

// WithLockReturn 是一个支持返回值的函数，用于在给定的锁上执行操作
func WithLockReturn[T any](lock Locker, operation func() (T, error)) (T, error) {
	lock.Lock()         // 获取锁
	defer lock.Unlock() // 确保在操作完成后释放锁
	return operation()  // 执行操作并返回结果
}

// WithRLock 是一个用于在给定的读锁上执行操作的函数
func WithRLock(lock RLocker, operation func()) {
	lock.RLock()         // 获取读锁
	defer lock.RUnlock() // 确保在操作完成后释放读锁
	operation()          // 执行操作
}

// WithRLockReturn 是一个支持返回值的函数，用于在给定的读锁上执行操作
func WithRLockReturn[T any](lock RLocker, operation func() (T, error)) (T, error) {
	lock.RLock()         // 获取读锁
	defer lock.RUnlock() // 确保在操作完成后释放读锁
	return operation()   // 执行操作并返回结果
}
