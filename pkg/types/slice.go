/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-05 15:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 11:26:10
 * @FilePath: \go-toolbox\pkg\types\slice.go
 * @Description: 切片通用操作（泛型实现）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package types

// Contains 检查值是否在切片中（泛型版本，支持任意可比较类型）
func Contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// ContainsAny 检查切片中是否包含任意一个目标值
func ContainsAny[T comparable](slice []T, values ...T) bool {
	for _, item := range slice {
		for _, value := range values {
			if item == value {
				return true
			}
		}
	}
	return false
}

// ContainsAll 检查切片中是否包含所有目标值
func ContainsAll[T comparable](slice []T, values ...T) bool {
	for _, value := range values {
		if !Contains(slice, value) {
			return false
		}
	}
	return true
}

// IndexOf 返回值在切片中的索引，不存在返回 -1
func IndexOf[T comparable](slice []T, value T) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

// Filter 过滤切片，返回满足条件的元素
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// MapTR 映射切片，将每个元素转换为另一种类型
func MapTR[T any, R any](slice []T, mapper func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Unique 去重切片（保持原始顺序）
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Reverse 反转切片（返回新切片，不修改原切片）
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, item := range slice {
		result[len(slice)-1-i] = item
	}
	return result
}

// Chunk 将切片分块
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	chunks := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
