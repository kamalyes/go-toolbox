/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:05:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-09 11:31:32
 * @FilePath: \go-toolbox\pkg\syncx\lock.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import "errors"

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

// WithLockReturnValue 是一个支持返回值的函数，用于在给定的锁上执行操作，不返回错误
func WithLockReturnValue[T any](lock Locker, operation func() T) T {
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

// WithRLockReturnValue 是一个支持返回值的函数，用于在给定的读锁上执行操作，不返回错误
func WithRLockReturnValue[T any](lock RLocker, operation func() T) T {
	lock.RLock()         // 获取读锁
	defer lock.RUnlock() // 确保在操作完成后释放读锁
	return operation()   // 执行操作并返回结果
}

var ErrLockNotAcquired = errors.New("lock not acquired")

// TryLocker 支持 TryLock 的锁接口
type TryLocker interface {
	Locker
	TryLock() bool
}

// TryRLocker 支持 TryRLock 的读锁接口
type TryRLocker interface {
	RLocker
	TryRLock() bool
}

// WithTryLock 在支持 TryLock 的锁上尝试执行操作，成功获取锁才执行
func WithTryLock(lock TryLocker, operation func()) error {
	if !lock.TryLock() {
		return ErrLockNotAcquired
	}
	defer lock.Unlock()
	operation()
	return nil
}

// WithTryLockReturn 在支持 TryLock 的锁上尝试执行操作，成功获取锁才执行，支持返回值和错误
func WithTryLockReturn[T any](lock TryLocker, operation func() (T, error)) (T, error) {
	var zero T
	if !lock.TryLock() {
		return zero, ErrLockNotAcquired
	}
	defer lock.Unlock()
	return operation()
}

// WithTryLockReturnValue 在支持 TryLock 的锁上尝试执行操作，成功获取锁才执行，支持返回值，不返回错误
func WithTryLockReturnValue[T any](lock TryLocker, operation func() T) (T, error) {
	var zero T
	if !lock.TryLock() {
		return zero, ErrLockNotAcquired
	}
	defer lock.Unlock()
	return operation(), nil
}

// WithTryRLock 在支持 TryRLock 的读锁上尝试执行操作，成功获取读锁才执行
func WithTryRLock(lock TryRLocker, operation func()) error {
	if !lock.TryRLock() {
		return ErrLockNotAcquired
	}
	defer lock.RUnlock()
	operation()
	return nil
}

// WithTryRLockReturn 在支持 TryRLock 的读锁上尝试执行操作，成功获取读锁才执行，支持返回值和错误
func WithTryRLockReturn[T any](lock TryRLocker, operation func() (T, error)) (T, error) {
	var zero T
	if !lock.TryRLock() {
		return zero, ErrLockNotAcquired
	}
	defer lock.RUnlock()
	return operation()
}

// WithTryRLockReturnValue 在支持 TryRLock 的读锁上尝试执行操作，成功获取读锁才执行，支持返回值，不返回错误
func WithTryRLockReturnValue[T any](lock TryRLocker, operation func() T) (T, error) {
	var zero T
	if !lock.TryRLock() {
		return zero, ErrLockNotAcquired
	}
	defer lock.RUnlock()
	return operation(), nil
}
