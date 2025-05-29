/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-05-29 13:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-29 13:55:55
 * @FilePath: \go-toolbox\tests\mathx_filter_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

func TestMapSliceByKey(t *testing.T) {
	type mapCase[T any, K comparable] struct {
		name    string
		input   []T
		keyFunc func(T) K
		want    map[K]T
	}

	type Person struct {
		ID   int
		Name string
	}

	// 用空接口切片存储不同类型的 mapCase
	cases := []interface{}{
		mapCase[int, int]{
			name:    "int slice, key is element",
			input:   []int{1, 2, 3, 2},
			keyFunc: func(i int) int { return i },
			want:    map[int]int{1: 1, 2: 2, 3: 3},
		},
		mapCase[string, string]{
			name:    "string slice, key is element",
			input:   []string{"a", "b", "a"},
			keyFunc: func(s string) string { return s },
			want:    map[string]string{"a": "a", "b": "b"},
		},
		mapCase[Person, int]{
			name: "struct slice, key is ID",
			input: []Person{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
				{ID: 1, Name: "Ann"},
			},
			keyFunc: func(p Person) int { return p.ID },
			want: map[int]Person{
				1: {ID: 1, Name: "Ann"},
				2: {ID: 2, Name: "Bob"},
			},
		},
		mapCase[*Person, int]{
			name: "pointer slice, key is ID",
			input: []*Person{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			keyFunc: func(p *Person) int { return p.ID },
			want: map[int]*Person{
				1: {ID: 1, Name: "Alice"},
				2: {ID: 2, Name: "Bob"},
			},
		},
	}

	// 统一执行函数，利用类型断言调用泛型函数
	for _, c := range cases {
		switch tc := c.(type) {
		case mapCase[int, int]:
			got := mathx.MapSliceByKey(tc.input, tc.keyFunc)
			assert.Equal(t, tc.want, got, tc.name)
		case mapCase[string, string]:
			got := mathx.MapSliceByKey(tc.input, tc.keyFunc)
			assert.Equal(t, tc.want, got, tc.name)
		case mapCase[Person, int]:
			got := mathx.MapSliceByKey(tc.input, tc.keyFunc)
			assert.Equal(t, tc.want, got, tc.name)
		case mapCase[*Person, int]:
			got := mathx.MapSliceByKey(tc.input, tc.keyFunc)
			assert.Equal(t, tc.want, got, tc.name)
		default:
			t.Errorf("unsupported case type %T", tc)
		}
	}
}

type ABC struct {
	key   string
	count int
}

func TestFilterSliceByFunc(t *testing.T) {
	input := []ABC{
		{key: "123"},
		{key: "456"},
		{key: "123"},
		{key: "789"},
	}

	expected := []ABC{
		{key: "123"},
		{key: "123"},
	}

	result := mathx.FilterSliceByFunc(input, func(d ABC) bool {
		return d.key == "123"
	})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FilterSliceByFunc 结果不正确，期望 %v，实际 %v", expected, result)
	}
}

func TestSliceFilter_BasicFilter(t *testing.T) {
	data := []ABC{
		{key: "a"},
		{key: "b"},
		{key: "a"},
	}

	filtered := mathx.NewSliceFilter(data).
		Condition(func(d ABC) bool {
			return d.key == "a"
		}).
		Result()

	assert.Len(t, filtered, 2)
	for _, item := range filtered {
		assert.Equal(t, "a", item.key)
	}
}

func TestSliceFilter_OnMatchCallback(t *testing.T) {
	data := []ABC{
		{key: "x", count: 1},
		{key: "y", count: 2},
	}

	called := 0
	filtered := mathx.NewSliceFilter(data).
		Condition(func(d ABC) bool {
			return d.key == "x"
		}).
		OnMatch(func(item *ABC) {
			called++
			item.count += 10
		}).
		Result()

	assert.Len(t, filtered, 1)
	assert.Equal(t, "x", filtered[0].key)
	assert.Equal(t, 11, filtered[0].count)
	assert.Equal(t, 1, called)
}

func TestSliceFilter_OnNotMatchCallback(t *testing.T) {
	data := []ABC{
		{key: "x", count: 1},
		{key: "y", count: 2},
	}

	notMatchedKeys := []string{}
	filtered := mathx.NewSliceFilter(data).
		Condition(func(d ABC) bool {
			return d.key == "x"
		}).
		OnNotMatch(func(item *ABC) {
			notMatchedKeys = append(notMatchedKeys, item.key)
		}).
		Result()

	assert.Len(t, filtered, 1)
	assert.Equal(t, "x", filtered[0].key)
	assert.ElementsMatch(t, []string{"y"}, notMatchedKeys)
}

func TestSliceFilter_MultipleCallbacks(t *testing.T) {
	data := []ABC{
		{key: "a"},
		{key: "b"},
	}

	matchCalls := 0
	notMatchCalls := 0

	filtered := mathx.NewSliceFilter(data).
		Condition(func(d ABC) bool {
			return d.key == "a"
		}).
		OnMatch(func(item *ABC) {
			matchCalls++
		}).
		OnMatch(func(item *ABC) {
			item.key = "modified"
		}).
		OnNotMatch(func(item *ABC) {
			notMatchCalls++
		}).
		Result()

	assert.Len(t, filtered, 1)
	assert.Equal(t, "modified", filtered[0].key)
	assert.Equal(t, 1, matchCalls)
	assert.Equal(t, 1, notMatchCalls)
}

func TestSliceFilter_NoPredicate(t *testing.T) {
	data := []ABC{
		{key: "x"},
		{key: "y"},
	}

	filtered := mathx.NewSliceFilter(data).Result()

	// 无筛选条件，应该返回全部
	assert.Len(t, filtered, 2)
}

func TestSliceFilter_ResultCalledMultipleTimes(t *testing.T) {
	data := []ABC{
		{key: "a"},
		{key: "b"},
	}

	f := mathx.NewSliceFilter(data).
		Condition(func(d ABC) bool {
			return d.key == "a"
		})

	res1 := f.Result()
	res2 := f.Result()

	assert.Equal(t, res1, res2)
}

func TestSliceFilter_EmptySlice(t *testing.T) {
	data := []ABC{}

	filtered := mathx.NewSliceFilter(data).
		Condition(func(d ABC) bool {
			return d.key == "a"
		}).
		Result()

	assert.Empty(t, filtered)
}

func TestSliceFilter_NilCallbacks(t *testing.T) {
	data := []ABC{
		{key: "a"},
		{key: "b"},
	}

	// 测试传入nil回调不会panic
	filtered := mathx.NewSliceFilter(data).
		Condition(func(d ABC) bool {
			return d.key == "a"
		}).
		OnMatch(nil).
		OnNotMatch(nil).
		Result()

	assert.Len(t, filtered, 1)
	assert.Equal(t, "a", filtered[0].key)
}

func TestSliceFilter_MultipleConditions(t *testing.T) {
	data := []ABC{
		{key: "a", count: 10},
		{key: "b", count: 5},
		{key: "c", count: 10},
		{key: "d", count: 20},
		{key: "e", count: 1},
	}

	filtered := mathx.NewSliceFilter(data).
		Condition(func(d ABC) bool {
			return d.key == "a"
		}, func(d ABC) bool {
			return d.count >= 20
		}).
		Condition(func(d ABC) bool {
			return d.count == 1
		}).
		UseOr().
		Result()

	expected := []ABC{
		{key: "a", count: 10},
		{key: "d", count: 20},
		{key: "e", count: 1},
	}

	assert.Equal(t, expected, filtered)
}

// 生成测试数据：构造大量 ABC 结构体，key 为 "key"+数字字符串，count 递增
func generateTestData(n int) []ABC {
	data := make([]ABC, n)
	for i := 0; i < n; i++ {
		data[i] = ABC{
			key:   "key" + strconv.Itoa(i%100), // 100 个不同 key，重复使用
			count: i,
		}
	}
	return data
}

// 基准测试：测试 SliceFilter 在大数据量下的性能
func BenchmarkSliceFilter_Result(b *testing.B) {
	data := generateTestData(100000) // 10 万条数据

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mathx.NewSliceFilter(data).
			Condition(func(d ABC) bool {
				// 筛选 key == "key50" 且 count >= 50000
				return d.key == "key50"
			}).
			Condition(func(d ABC) bool {
				return d.count >= 50000
			}).
			Result()
	}
}

// 基准测试：测试 FilterSliceByFunc 函数的性能对比
func BenchmarkFilterSliceByFunc(b *testing.B) {
	data := generateTestData(100000) // 10 万条数据

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mathx.FilterSliceByFunc(data, func(d ABC) bool {
			return d.key == "key50" && d.count >= 50000
		})
	}
}
