/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 15:55:06
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 13:40:33
 * @FilePath: \go-toolbox\tests\mathx_slice_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/random"
)

// 基准测试 SliceMinMax
func BenchmarkSliceMinMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := random.RandNumericalLargeSlice[int]()
		minMaxFunc := func(a, b int) int {
			if a < b {
				return a
			}
			return b
		}
		_, err := mathx.SliceMinMax(list, minMaxFunc)
		if err != nil {
			b.Fatalf("expected no error, got %v", err)
		}
	}
}

// 基准测试 SliceDiffSet
func BenchmarkSliceDiffSet(b *testing.B) {

	for i := 0; i < b.N; i++ {
		arr1 := random.RandNumericalLargeSlice[int](200)
		arr2 := random.RandNumericalLargeSlice[int](200)

		for i := 0; i < len(arr2); i++ {
			arr2[i] = i + len(arr1)/2 // 使得部分重叠
		}
		_ = mathx.SliceDiffSetSorted(arr1, arr2)
	}
}

// 基准测试 SliceUnion
func BenchmarkSliceUnion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr1 := random.RandNumericalLargeSlice[int](200)
		arr2 := random.RandNumericalLargeSlice[int](200)

		for i := 0; i < len(arr1); i++ {
			arr2[i] = i + len(arr1)/2 // 使得部分重叠
		}

		_ = mathx.SliceUnion(arr1, arr2)
	}
}

// 基准测试 SliceContains = 300
func BenchmarkSliceContains_300(b *testing.B) {

	for i := 0; i < b.N; i++ {
		intSlice := random.RandNumericalLargeSlice[int](300)

		for i := 0; i < len(intSlice); i++ {
			intSlice[i] = i + len(intSlice)/2 // 使得部分重叠
		}

		element := random.RandInt(len(intSlice)/2, len(intSlice)*2) // 测试查找的元素
		_ = mathx.SliceContains(intSlice, element)
	}
}

// 基准测试 SliceContains = 3000
func BenchmarkSliceContains_3000(b *testing.B) {

	for i := 0; i < b.N; i++ {
		intSlice := random.RandNumericalLargeSlice[int](3000)

		for i := 0; i < len(intSlice); i++ {
			intSlice[i] = i + len(intSlice)/2 // 使得部分重叠
		}

		element := random.RandInt(len(intSlice)/2, len(intSlice)*2) // 测试查找的元素
		_ = mathx.SliceContains(intSlice, element)
	}
}

// 基准测试 SliceContains > 3000
func BenchmarkSliceContains_20000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intSlice := random.RandNumericalLargeSlice[int](20000)

		for i := 0; i < len(intSlice); i++ {
			intSlice[i] = i + len(intSlice)/2 // 使得部分重叠
		}

		element := random.RandInt(len(intSlice)/2, len(intSlice)*2) // 测试查找的元素

		_ = mathx.SliceContains(intSlice, element)
	}
}

// 基准测试 SliceHasDuplicates
func BenchmarkSliceHasDuplicates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intSlice := random.RandNumericalLargeSlice[int](20000)
		_ = mathx.SliceHasDuplicates(intSlice)
	}
}

// 基准测试 SliceRemoveEmpty
func BenchmarkSliceRemoveEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		defaultSliceSize := 10000
		intSlice := make([]interface{}, defaultSliceSize)
		for i := 0; i < defaultSliceSize; i++ {
			if i%10 == 0 {
				intSlice[i] = nil // 每10个元素放一个空值
			} else {
				intSlice[i] = i
			}
		}
		_ = mathx.SliceRemoveEmpty(intSlice)
	}
}

// 基准测试 SliceRemoveDuplicates
func BenchmarkSliceRemoveDuplicates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intSlice := random.RandNumericalLargeSlice[int]()
		_ = mathx.SliceRemoveDuplicates(intSlice)
	}
}

// 基准测试 SliceRemoveZero
func BenchmarkSliceRemoveZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		defaultSliceSize := 1000
		arr := make([]int, defaultSliceSize)
		for i := 0; i < defaultSliceSize; i++ {
			arr[i] = i % 10 // 生成一些零值
		}
		_ = mathx.SliceRemoveZero(arr)
	}
}

// 基准测试 SliceChunk
func BenchmarkSliceChunk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := random.RandNumericalLargeSlice[int]()
		size := 1000 // 每个子切片的大小
		_ = mathx.SliceChunk(slice, size)
	}
}

func BenchmarkInsertionSort_100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := random.RandNumericalLargeSlice[int](100)
		mathx.InsertionSort(slice)
	}
}

func BenchmarkQuickSort_100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := random.RandNumericalLargeSlice[int](100)
		mathx.QuickSort(slice, 0, len(slice)-1)
	}
}

func BenchmarkBubbleSort_100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 生成随机数组
		slice := random.RandNumericalLargeSlice[int](100)
		mathx.BubbleSort(slice)
	}
}
