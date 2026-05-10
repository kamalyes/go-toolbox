/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-09 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-09 00:00:00
 * @FilePath: \go-toolbox\pkg\types\pointer.go
 * @Description: 泛型指针辅助工具
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package types

// Ptr 创建任意类型的指针（泛型版本，比 validator 包中 IntPtr/BoolPtr 等更通用）
// 示例: p := types.Ptr(42) -> *int
func Ptr[T any](v T) *T {
	return &v
}

// Deref 解引用指针，返回底层值；如果指针为 nil 则返回类型零值
// 示例: types.Deref((*int)(nil)) -> 0
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// DerefOrDefault 解引用指针，返回底层值；如果指针为 nil 则返回默认值
// 示例: types.DerefOrDefault((*string)(nil), "hello") -> "hello"
func DerefOrDefault[T any](p *T, defaultValue T) T {
	if p == nil {
		return defaultValue
	}
	return *p
}

// IsNilPtr 检查指针是否为 nil
func IsNilPtr[T any](p *T) bool {
	return p == nil
}

// IsNonNilPtr 检查指针是否非 nil
func IsNonNilPtr[T any](p *T) bool {
	return p != nil
}
