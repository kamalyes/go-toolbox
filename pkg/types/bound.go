/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-23 09:11:20
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-24 18:59:29
 * @FilePath: \go-toolbox\pkg\types\bound.go
 * @Description: 边界和范围相关的类型定义
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package types

// Bounds 定义字段的取值范围（泛型）
type Bounds[T Numerical] struct {
	Min   T            // 最小值
	Max   T            // 最大值
	Names map[string]T // 名称映射（如 "jan"->1, "mon"->1）
}

// BoundType 边界类型枚举
type BoundType int

const (
	BoundClosed         BoundType = iota // 闭区间 [min, max]
	BoundOpen                            // 开区间 (min, max)
	BoundLeftOpen                        // 左开右闭 (min, max]
	BoundRightOpen                       // 左闭右开 [min, max)
	BoundLeftUnbounded                   // 左无界 (-∞, max]
	BoundRightUnbounded                  // 右无界 [min, +∞)
	BoundUnbounded                       // 无界 (-∞, +∞)
)

// RangeMode 范围解析模式
type RangeMode int

const (
	RangeModeNormal   RangeMode = iota // 普通模式：精确匹配
	RangeModeWildcard                  // 通配符模式：支持 * 和 ?
	RangeModeStep                      // 步长模式：支持 /step
	RangeModeList                      // 列表模式：支持逗号分隔
)

// BoundError 边界错误类型
type BoundError int

const (
	BoundErrorNone         BoundError = iota // 无错误
	BoundErrorBelowMin                       // 低于最小值
	BoundErrorAboveMax                       // 超过最大值
	BoundErrorInvalidRange                   // 无效范围（min > max）
	BoundErrorZeroStep                       // 步长为零
	BoundErrorNegative                       // 负数错误
)

// RangeValidator 范围验证器类型（函数类型）
type RangeValidator[T Numerical] func(value T, bounds Bounds[T]) BoundError

// RangeParser 范围解析器类型（函数类型）
type RangeParser[T Numerical] func(expr string, bounds Bounds[T]) (T, error)

// RangeTransformer 范围转换器类型（函数类型）
type RangeTransformer[T Numerical, R any] func(value T, bounds Bounds[T]) (R, error)
