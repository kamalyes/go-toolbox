/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 10:50:01
 * @FilePath: \go-toolbox\pkg\mathx\stats.go
 * @Description: 统计功能
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"cmp"
	"fmt"
	"math"
	"sort"
)

// Percentile 计算百分位数（支持50, 90, 95, 99）
func Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	index := int(math.Ceil(float64(len(sorted)) * p / 100.0))
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// Percentiles 批量计算多个百分位数
func Percentiles(values []float64, percentiles ...float64) map[float64]float64 {
	result := make(map[float64]float64, len(percentiles))

	if len(values) == 0 {
		for _, p := range percentiles {
			result[p] = 0
		}
		return result
	}

	// 只排序一次
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	for _, p := range percentiles {
		index := int(math.Ceil(float64(len(sorted)) * p / 100.0))
		if index >= len(sorted) {
			index = len(sorted) - 1
		}
		result[p] = sorted[index]
	}

	return result
}

// Percentage 计算百分比
func Percentage(part, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

// FormatPercentage 格式化百分比
func FormatPercentage(part, total uint64, precision int) string {
	return fmt.Sprintf("%.*f%%", precision, Percentage(part, total))
}

// Mean 计算平均值
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// StdDev 计算标准差
func StdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := Mean(values)
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}

	return math.Sqrt(sumSquares / float64(len(values)))
}

// SortByCount 按计数排序统计数据（降序）
func SortByCount[T any](items []T, getCount func(T) uint64) {
	sort.Slice(items, func(i, j int) bool {
		return getCount(items[i]) > getCount(items[j])
	})
}

// SortByKey 按键排序统计数据（升序）- 支持任意可比较类型
func SortByKey[T any, K cmp.Ordered](items []T, getKey func(T) K) {
	sort.Slice(items, func(i, j int) bool {
		return cmp.Compare(getKey(items[i]), getKey(items[j])) < 0
	})
}

// SortByKeyDesc 按键排序统计数据（降序）- 支持任意可比较类型
func SortByKeyDesc[T any, K cmp.Ordered](items []T, getKey func(T) K) {
	sort.Slice(items, func(i, j int) bool {
		return cmp.Compare(getKey(items[i]), getKey(items[j])) > 0
	})
}

// SortByKeyDescUnique 按键排序统计数据（降序）并去重
// 去重规则：保留每个唯一标识的第一个元素（即排序后权重最大的）
// getKey: 用于排序的键提取函数
// getID: 用于去重的唯一标识提取函数
func SortByKeyDescUnique[T any, K cmp.Ordered, ID comparable](items []T, getKey func(T) K, getID func(T) ID) []T {
	// 先按键降序排序
	sort.Slice(items, func(i, j int) bool {
		return cmp.Compare(getKey(items[i]), getKey(items[j])) > 0
	})

	// 去重：保留每个 ID 的第一个元素（权重最大的）
	seen := make(map[ID]bool, len(items))
	result := make([]T, 0, len(items))
	for _, item := range items {
		id := getID(item)
		if !seen[id] {
			seen[id] = true
			result = append(result, item)
		}
	}

	return result
}

// StatsSummary 统计摘要
type StatsSummary struct {
	Count  int     `json:"count"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
	P50    float64 `json:"p50"`
	P90    float64 `json:"p90"`
	P95    float64 `json:"p95"`
	P99    float64 `json:"p99"`
}

// SummarizeStats 生成统计摘要
func SummarizeStats(values []float64) StatsSummary {
	if len(values) == 0 {
		return StatsSummary{}
	}

	percentiles := Percentiles(values, 50, 90, 95, 99)

	return StatsSummary{
		Count:  len(values),
		Min:    Min(values...),
		Max:    Max(values...),
		Mean:   Mean(values),
		StdDev: StdDev(values),
		P50:    percentiles[50],
		P90:    percentiles[90],
		P95:    percentiles[95],
		P99:    percentiles[99],
	}
}
