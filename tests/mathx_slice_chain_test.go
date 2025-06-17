/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-17 13:06:30
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-17 13:55:21
 * @FilePath: \go-toolbox\tests\mathx_slice_chain_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

// 通用断言函数
func assertSliceChainEqual[T comparable](t *testing.T, sc *mathx.SliceChain[T], want []T) {
	t.Helper()
	assert.Equal(t, want, sc.Data())
}

func TestSliceChain_Append(t *testing.T) {
	sc := mathx.FromSlice([]int{1, 2}).Append(3, 4)
	assertSliceChainEqual(t, sc, []int{1, 2, 3, 4})

	// 追加空元素不改变切片
	sc.Append()
	assertSliceChainEqual(t, sc, []int{1, 2, 3, 4})

	// 空切片追加元素
	scEmpty := mathx.FromSlice([]int{}).Append(10)
	assertSliceChainEqual(t, scEmpty, []int{10})
}

func TestSliceChain_Uniq(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want []int
	}{
		{"normal", []int{1, 2, 2, 3, 1, 4}, []int{1, 2, 3, 4}},
		{"empty", []int{}, []int{}},
		{"single", []int{5}, []int{5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := mathx.FromSlice(tt.data)
			sc.Uniq()
			assertSliceChainEqual(t, sc, tt.want)
		})
	}
}

func TestSliceChain_RemoveValue(t *testing.T) {
	tests := []struct {
		name   string
		data   []string
		remove string
		want   []string
	}{
		{"remove exists", []string{"a", "b", "a", "c"}, "a", []string{"b", "c"}},
		{"remove not exists", []string{"b", "c"}, "x", []string{"b", "c"}},
		{"empty slice", []string{}, "a", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := mathx.FromSlice(tt.data)
			sc.RemoveValue(tt.remove)
			assertSliceChainEqual(t, sc, tt.want)
		})
	}
}

func TestSliceChain_RemoveEmpty(t *testing.T) {
	tests := []struct {
		name string
		data []string
		want []string
	}{
		{"mixed empty", []string{"", "a", "", "b"}, []string{"a", "b"}},
		{"all empty", []string{"", "", ""}, []string{}},
		{"empty slice", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := mathx.FromSlice(tt.data)
			sc.RemoveEmpty()
			assertSliceChainEqual(t, sc, tt.want)
		})
	}
}

func TestSliceChain_Filter(t *testing.T) {
	sc := mathx.FromSlice([]int{1, 2, 3, 4, 5})
	sc.Filter(func(x int) bool { return x%2 == 1 })
	assertSliceChainEqual(t, sc, []int{1, 3, 5})

	sc.Filter(func(x int) bool { return false })
	assertSliceChainEqual(t, sc, []int{})

	scEmpty := mathx.FromSlice([]int{})
	scEmpty.Filter(func(x int) bool { return true })
	assertSliceChainEqual(t, scEmpty, []int{})
}

func TestSliceChain_Sort(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want []int
	}{
		{"normal", []int{5, 3, 4, 1, 2}, []int{1, 2, 3, 4, 5}},
		{"empty", []int{}, []int{}},
		{"single", []int{42}, []int{42}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := mathx.FromSlice(tt.data)
			sc.Sort(func(a, b int) bool { return a < b })
			assertSliceChainEqual(t, sc, tt.want)
		})
	}
}

func TestSliceChain_Data(t *testing.T) {
	original := []string{"x", "y", "z"}
	sc := mathx.FromSlice(original)
	data := sc.Data()
	assertSliceChainEqual(t, sc, original)

	// 修改返回切片会影响内部数据，因为是同一引用
	data[0] = "changed"
	assertSliceChainEqual(t, sc, data)
}

func TestSliceChain_ChainUsage(t *testing.T) {
	sc := mathx.FromSlice([]int{5, 3, 3, 2, 1, 1, 4})

	// 链式调用模拟
	ops := []func(*mathx.SliceChain[int]) *mathx.SliceChain[int]{
		func(sc *mathx.SliceChain[int]) *mathx.SliceChain[int] { return sc.Uniq() },
		func(sc *mathx.SliceChain[int]) *mathx.SliceChain[int] {
			return sc.Sort(func(a, b int) bool { return a < b })
		},
		func(sc *mathx.SliceChain[int]) *mathx.SliceChain[int] {
			return sc.Filter(func(x int) bool { return x > 2 })
		},
		func(sc *mathx.SliceChain[int]) *mathx.SliceChain[int] { return sc.Append(6, 7) },
		func(sc *mathx.SliceChain[int]) *mathx.SliceChain[int] { return sc.RemoveValue(3) },
	}

	for _, op := range ops {
		sc = op(sc)
	}

	assertSliceChainEqual(t, sc, []int{4, 5, 6, 7})
}

// 并发安全测试
func TestSliceChain_ConcurrentSafety(t *testing.T) {
	sc := mathx.FromSlice([]int{})

	var wg sync.WaitGroup
	concurrency := 50

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sc.Append(id)
			sc.RemoveValue(id - 1)                        // 尝试移除不存在的元素，测试稳定性
			sc.Filter(func(x int) bool { return x >= 0 }) // 过滤所有元素（不过滤）
			sc.Uniq()
		}(i)
	}

	wg.Wait()

	// 最终切片元素个数不超过并发数，且无重复
	data := sc.Data()
	seen := make(map[int]struct{})
	for _, v := range data {
		if _, ok := seen[v]; ok {
			t.Errorf("duplicate element detected: %v", v)
		}
		seen[v] = struct{}{}
	}
	assert.LessOrEqual(t, len(data), concurrency)
}

// 性能基准测试

func BenchmarkSliceChain_Append(b *testing.B) {
	sc := mathx.FromSlice([]int{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Append(i)
	}
}

func BenchmarkSliceChain_RemoveValue(b *testing.B) {
	// 预填充大量数据
	data := make([]int, 10000)
	for i := range data {
		data[i] = i % 100 // 0~99重复
	}
	sc := mathx.FromSlice(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.RemoveValue(i % 100)
	}
}

func BenchmarkSliceChain_Uniq(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i % 1000 // 0~999重复
	}
	sc := mathx.FromSlice(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Uniq()
	}
}

func BenchmarkSliceChain_Filter(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	sc := mathx.FromSlice(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Filter(func(x int) bool { return x%2 == 0 })
	}
}

func BenchmarkSliceChain_Sort(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = 10000 - i
	}
	sc := mathx.FromSlice(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Sort(func(a, b int) bool { return a < b })
	}
}
