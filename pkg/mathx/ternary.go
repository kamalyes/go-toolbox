/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-21 19:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-22 08:55:16
 * @FilePath: \go-toolbox\pkg\mathx\ternary.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package mathx

// IF 实现三元运算，使用泛型 T
func IF[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// DoFunc 是一个函数类型，用于执行返回泛型 T 的函数
type DoFunc[T any] func() T

// IfDo 根据条件执行函数并返回结果，支持默认值
func IfDo[T any](condition bool, do DoFunc[T], defaultVal T) T {
	if condition {
		return do() // 执行函数并返回结果
	}
	return defaultVal // 返回默认值
}

// IfDoAF 根据条件执行函数和默认函数并返回结果
func IfDoAF[T any](condition bool, do DoFunc[T], defaultFunc DoFunc[T]) T {
	return IfDo(condition, do, defaultFunc()) // 使用 IfDo 函数
}
