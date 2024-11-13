/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 15:55:06
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 13:13:57
 * @FilePath: \go-toolbox\tests\mathx_array_bench_test.go
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

// 基准测试 ArrayMinMax
func BenchmarkArrayMinMax(b *testing.B) {
	list := random.RandNumericalLargeSlice[int]()
	minMaxFunc := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	for i := 0; i < b.N; i++ {
		_, err := mathx.ArrayMinMax(list, minMaxFunc)
		if err != nil {
			b.Fatalf("expected no error, got %v", err)
		}
	}
}

// 基准测试 ArrayDiffSet
func BenchmarkArrayDiffSet(b *testing.B) {
	arr1 := random.RandNumericalLargeSlice[int](200)
	arr2 := random.RandNumericalLargeSlice[int](200)

	for i := 0; i < len(arr2); i++ {
		arr2[i] = i + len(arr1)/2 // 使得部分重叠
	}

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayDiffSetSorted(arr1, arr2)
	}
}

// 基准测试 ArrayUnion
func BenchmarkArrayUnion(b *testing.B) {
	arr1 := random.RandNumericalLargeSlice[int](200)
	arr2 := random.RandNumericalLargeSlice[int](200)

	for i := 0; i < len(arr1); i++ {
		arr2[i] = i + len(arr1)/2 // 使得部分重叠
	}

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayUnion(arr1, arr2)
	}
}

// 基准测试 ArrayContains = 300
func BenchmarkArrayContains_300(b *testing.B) {
	intArray := random.RandNumericalLargeSlice[int](300)

	for i := 0; i < len(intArray); i++ {
		intArray[i] = i + len(intArray)/2 // 使得部分重叠
	}

	element := random.RandInt(len(intArray)/2, len(intArray)*2) // 测试查找的元素

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayContains(intArray, element)
	}
}

// 基准测试 ArrayContains = 3000
func BenchmarkArrayContains_3000(b *testing.B) {
	intArray := random.RandNumericalLargeSlice[int](3000)

	for i := 0; i < len(intArray); i++ {
		intArray[i] = i + len(intArray)/2 // 使得部分重叠
	}

	element := random.RandInt(len(intArray)/2, len(intArray)*2) // 测试查找的元素

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayContains(intArray, element)
	}
}

// 基准测试 ArrayContains > 3000
func BenchmarkArrayContains_20000(b *testing.B) {
	intArray := random.RandNumericalLargeSlice[int](20000)

	for i := 0; i < len(intArray); i++ {
		intArray[i] = i + len(intArray)/2 // 使得部分重叠
	}

	element := random.RandInt(len(intArray)/2, len(intArray)*2) // 测试查找的元素

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayContains(intArray, element)
	}
}

// 基准测试 ArrayHasDuplicates
func BenchmarkArrayHasDuplicates(b *testing.B) {
	intArray := random.RandNumericalLargeSlice[int](20000)

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayHasDuplicates(intArray)
	}
}

// 基准测试 ArrayRemoveEmpty
func BenchmarkArrayRemoveEmpty(b *testing.B) {
	defaultSliceSize := 10000
	intArray := make([]interface{}, defaultSliceSize)
	for i := 0; i < defaultSliceSize; i++ {
		if i%10 == 0 {
			intArray[i] = nil // 每10个元素放一个空值
		} else {
			intArray[i] = i
		}
	}

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayRemoveEmpty(intArray)
	}
}

// 基准测试 ArrayRemoveDuplicates
func BenchmarkArrayRemoveDuplicates(b *testing.B) {
	intArray := random.RandNumericalLargeSlice[int]()

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayRemoveDuplicates(intArray)
	}
}

// 基准测试 ArrayRemoveZero
func BenchmarkArrayRemoveZero(b *testing.B) {
	defaultSliceSize := 1000
	arr := make([]int, defaultSliceSize)
	for i := 0; i < defaultSliceSize; i++ {
		arr[i] = i % 10 // 生成一些零值
	}

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayRemoveZero(arr)
	}
}

// 基准测试 ArrayChunk
func BenchmarkArrayChunk(b *testing.B) {
	slice := random.RandNumericalLargeSlice[int]()
	size := 1000 // 每个子切片的大小

	for i := 0; i < b.N; i++ {
		_ = mathx.ArrayChunk(slice, size)
	}
}
