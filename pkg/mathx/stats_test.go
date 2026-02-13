/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-24 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-24 10:00:00
 * @FilePath: \go-toolbox\pkg\mathx\stats_test.go
 * @Description: 统计功能测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSortByCount 测试按计数排序
func TestSortByCount(t *testing.T) {
	type Item struct {
		Name  string
		Count uint64
	}

	tests := []struct {
		name     string
		items    []Item
		expected []string // 期望的名称顺序
	}{
		{
			name: "BasicSort",
			items: []Item{
				{"A", 10},
				{"B", 50},
				{"C", 30},
			},
			expected: []string{"B", "C", "A"}, // 降序：50, 30, 10
		},
		{
			name: "SameCount",
			items: []Item{
				{"X", 20},
				{"Y", 20},
				{"Z", 20},
			},
			expected: []string{"X", "Y", "Z"}, // 相同计数保持原顺序
		},
		{
			name:     "EmptySlice",
			items:    []Item{},
			expected: []string{},
		},
		{
			name: "SingleItem",
			items: []Item{
				{"Single", 100},
			},
			expected: []string{"Single"},
		},
		{
			name: "LargeNumbers",
			items: []Item{
				{"Error1", 1000000},
				{"Error2", 999999},
				{"Error3", 1000001},
			},
			expected: []string{"Error3", "Error1", "Error2"}, // 1000001, 1000000, 999999
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制切片避免影响原数据
			items := make([]Item, len(tt.items))
			copy(items, tt.items)

			SortByCount(items, func(item Item) uint64 {
				return item.Count
			})

			// 验证排序后的名称顺序
			actual := make([]string, len(items))
			for i, item := range items {
				actual[i] = item.Name
			}

			assert.Equal(t, tt.expected, actual, "排序结果不符合预期")
		})
	}
}

// TestSortByKey 测试按键升序排序
func TestSortByKey(t *testing.T) {
	type StatusCode struct {
		Code  int
		Count uint64
	}

	tests := []struct {
		name     string
		items    []StatusCode
		expected []int // 期望的状态码顺序
	}{
		{
			name: "IntegerSort",
			items: []StatusCode{
				{500, 5},
				{200, 100},
				{404, 20},
			},
			expected: []int{200, 404, 500}, // 升序
		},
		{
			name: "AlreadySorted",
			items: []StatusCode{
				{100, 1},
				{200, 2},
				{300, 3},
			},
			expected: []int{100, 200, 300},
		},
		{
			name:     "EmptySlice",
			items:    []StatusCode{},
			expected: []int{},
		},
		{
			name: "NegativeNumbers",
			items: []StatusCode{
				{-1, 10},
				{0, 20},
				{1, 30},
			},
			expected: []int{-1, 0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := make([]StatusCode, len(tt.items))
			copy(items, tt.items)

			SortByKey(items, func(s StatusCode) int {
				return s.Code
			})

			actual := make([]int, len(items))
			for i, item := range items {
				actual[i] = item.Code
			}

			assert.Equal(t, tt.expected, actual, "升序排序结果不符合预期")
		})
	}
}

// TestSortByKeyDesc 测试按键降序排序
func TestSortByKeyDesc(t *testing.T) {
	type Score struct {
		Name  string
		Value float64
	}

	tests := []struct {
		name     string
		items    []Score
		expected []string // 期望的名称顺序（按分数降序）
	}{
		{
			name: "FloatSort",
			items: []Score{
				{"Alice", 85.5},
				{"Bob", 92.3},
				{"Charlie", 78.9},
			},
			expected: []string{"Bob", "Alice", "Charlie"}, // 降序：92.3, 85.5, 78.9
		},
		{
			name: "StringSort",
			items: []Score{
				{"Zebra", 0},
				{"Apple", 0},
				{"Mango", 0},
			},
			expected: []string{"Zebra", "Mango", "Apple"}, // 字符串降序
		},
		{
			name:     "EmptySlice",
			items:    []Score{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := make([]Score, len(tt.items))
			copy(items, tt.items)

			// 测试按 Value 降序
			if tt.name == "FloatSort" {
				SortByKeyDesc(items, func(s Score) float64 {
					return s.Value
				})
			} else {
				// 测试按 Name 降序
				SortByKeyDesc(items, func(s Score) string {
					return s.Name
				})
			}

			actual := make([]string, len(items))
			for i, item := range items {
				actual[i] = item.Name
			}

			assert.Equal(t, tt.expected, actual, "降序排序结果不符合预期")
		})
	}
}

// TestPercentage 测试百分比计算
func TestPercentage(t *testing.T) {
	tests := []struct {
		name     string
		part     uint64
		total    uint64
		expected float64
	}{
		{"Normal", 50, 100, 50.0},
		{"Zero", 0, 100, 0.0},
		{"ZeroTotal", 50, 0, 0.0},
		{"FullPercent", 100, 100, 100.0},
		{"SmallPercent", 1, 1000, 0.1},
		{"LargeNumbers", 999999, 1000000, 99.9999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Percentage(tt.part, tt.total)
			assert.InDelta(t, tt.expected, result, 0.0001, "百分比计算结果不符合预期")
		})
	}
}

// TestPercentile 测试百分位数计算
func TestPercentile(t *testing.T) {
	tests := []struct {
		name       string
		values     []float64
		percentile float64
		expected   float64
	}{
		{"P50", []float64{1, 2, 3, 4, 5}, 50, 4},                  // Ceil(5 * 0.5) = 3, sorted[3] = 4
		{"P90", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 90, 10}, // Ceil(10 * 0.9) = 9, sorted[9] = 10
		{"P95", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 95, 10},
		{"P99", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 99, 10},
		{"EmptySlice", []float64{}, 50, 0},
		{"SingleValue", []float64{42}, 50, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Percentile(tt.values, tt.percentile)
			assert.Equal(t, tt.expected, result, "百分位数计算结果不符合预期")
		})
	}
}

// TestPercentiles 测试批量百分位数计算
func TestPercentiles(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	percentiles := Percentiles(values, 50, 90, 95, 99)

	assert.Equal(t, 6.0, percentiles[50], "P50 计算错误")  // Ceil(10 * 0.5) = 5, sorted[5] = 6
	assert.Equal(t, 10.0, percentiles[90], "P90 计算错误") // Ceil(10 * 0.9) = 9, sorted[9] = 10
	assert.Equal(t, 10.0, percentiles[95], "P95 计算错误")
	assert.Equal(t, 10.0, percentiles[99], "P99 计算错误")

	// 测试空切片
	emptyPercentiles := Percentiles([]float64{}, 50, 90)
	assert.Equal(t, 0.0, emptyPercentiles[50], "空切片应返回 0")
	assert.Equal(t, 0.0, emptyPercentiles[90], "空切片应返回 0")
}

type Agent struct {
	ID     string
	Weight int32
	Name   string
}

func TestSortByKeyDescUnique(t *testing.T) {
	tests := []struct {
		name     string
		agents   []Agent
		expected []Agent
	}{
		{
			name: "去重保留权重最大的",
			agents: []Agent{
				{ID: "agent1", Weight: 100, Name: "Alice"},
				{ID: "agent2", Weight: 200, Name: "Bob"},
				{ID: "agent1", Weight: 50, Name: "Alice-Duplicate"}, // 重复，权重更小
				{ID: "agent3", Weight: 150, Name: "Charlie"},
				{ID: "agent2", Weight: 180, Name: "Bob-Duplicate"}, // 重复，权重更小
			},
			expected: []Agent{
				{ID: "agent2", Weight: 200, Name: "Bob"},     // 权重最大
				{ID: "agent3", Weight: 150, Name: "Charlie"}, // 第二大
				{ID: "agent1", Weight: 100, Name: "Alice"},   // 第三大（保留权重100的，丢弃50的）
			},
		},
		{
			name: "无重复元素",
			agents: []Agent{
				{ID: "agent1", Weight: 100, Name: "Alice"},
				{ID: "agent2", Weight: 200, Name: "Bob"},
				{ID: "agent3", Weight: 150, Name: "Charlie"},
			},
			expected: []Agent{
				{ID: "agent2", Weight: 200, Name: "Bob"},
				{ID: "agent3", Weight: 150, Name: "Charlie"},
				{ID: "agent1", Weight: 100, Name: "Alice"},
			},
		},
		{
			name: "全部重复",
			agents: []Agent{
				{ID: "agent1", Weight: 100, Name: "Alice-1"},
				{ID: "agent1", Weight: 90, Name: "Alice-2"},
				{ID: "agent1", Weight: 80, Name: "Alice-3"},
			},
			expected: []Agent{
				{ID: "agent1", Weight: 100, Name: "Alice-1"}, // 只保留权重最大的
			},
		},
		{
			name:     "空列表",
			agents:   []Agent{},
			expected: []Agent{},
		},
		{
			name: "单个元素",
			agents: []Agent{
				{ID: "agent1", Weight: 100, Name: "Alice"},
			},
			expected: []Agent{
				{ID: "agent1", Weight: 100, Name: "Alice"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortByKeyDescUnique(
				tt.agents,
				func(a Agent) int32 { return a.Weight },
				func(a Agent) string { return a.ID },
			)

			if len(result) != len(tt.expected) {
				t.Errorf("长度不匹配: got %d, want %d", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i].ID != tt.expected[i].ID ||
					result[i].Weight != tt.expected[i].Weight ||
					result[i].Name != tt.expected[i].Name {
					t.Errorf("索引 %d 不匹配:\ngot  %+v\nwant %+v",
						i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestSortByKeyDescUnique_IntKey(t *testing.T) {
	type Item struct {
		ID    int
		Score float64
	}

	items := []Item{
		{ID: 1, Score: 95.5},
		{ID: 2, Score: 88.0},
		{ID: 1, Score: 92.0}, // 重复，分数更低
		{ID: 3, Score: 99.0},
	}

	result := SortByKeyDescUnique(
		items,
		func(i Item) float64 { return i.Score },
		func(i Item) int { return i.ID },
	)

	expected := []Item{
		{ID: 3, Score: 99.0},
		{ID: 1, Score: 95.5}, // 保留分数更高的
		{ID: 2, Score: 88.0},
	}

	if len(result) != len(expected) {
		t.Fatalf("长度不匹配: got %d, want %d", len(result), len(expected))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("索引 %d 不匹配: got %+v, want %+v", i, result[i], expected[i])
		}
	}
}

// TestMean 测试平均值计算
func TestMean(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"Normal", []float64{1, 2, 3, 4, 5}, 3.0},
		{"SingleValue", []float64{42}, 42.0},
		{"EmptySlice", []float64{}, 0.0},
		{"Negative", []float64{-5, -3, -1}, -3.0},
		{"Mixed", []float64{-10, 0, 10}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mean(tt.values)
			assert.InDelta(t, tt.expected, result, 0.0001, "平均值计算结果不符合预期")
		})
	}
}

// TestStdDev 测试标准差计算
func TestStdDev(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"Normal", []float64{2, 4, 4, 4, 5, 5, 7, 9}, 2.0},
		{"SingleValue", []float64{42}, 0.0},
		{"EmptySlice", []float64{}, 0.0},
		{"SameValues", []float64{5, 5, 5, 5}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StdDev(tt.values)
			assert.InDelta(t, tt.expected, result, 0.0001, "标准差计算结果不符合预期")
		})
	}
}

// TestSummarizeStats 测试统计摘要
func TestSummarizeStats(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	summary := SummarizeStats(values)

	assert.Equal(t, 10, summary.Count, "Count 错误")
	assert.Equal(t, 1.0, summary.Min, "Min 错误")
	assert.Equal(t, 10.0, summary.Max, "Max 错误")
	assert.InDelta(t, 5.5, summary.Mean, 0.1, "Mean 错误")
	assert.Greater(t, summary.StdDev, 0.0, "StdDev 应大于 0")
	assert.Equal(t, 6.0, summary.P50, "P50 错误")  // Ceil(10 * 0.5) = 5, sorted[5] = 6
	assert.Equal(t, 10.0, summary.P90, "P90 错误") // Ceil(10 * 0.9) = 9, sorted[9] = 10
	assert.Equal(t, 10.0, summary.P95, "P95 错误")
	assert.Equal(t, 10.0, summary.P99, "P99 错误")

	// 测试空切片
	emptySummary := SummarizeStats([]float64{})
	assert.Equal(t, 0, emptySummary.Count, "空切片的 Count 应为 0")
}

// BenchmarkSortByCount 基准测试：按计数排序
func BenchmarkSortByCount(b *testing.B) {
	type Item struct {
		Name  string
		Count uint64
	}

	items := make([]Item, 1000)
	for i := range items {
		items[i] = Item{Name: string(rune(i)), Count: uint64(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testItems := make([]Item, len(items))
		copy(testItems, items)
		SortByCount(testItems, func(item Item) uint64 {
			return item.Count
		})
	}
}

// BenchmarkSortByKey 基准测试：按键排序
func BenchmarkSortByKey(b *testing.B) {
	type Item struct {
		ID    int
		Value string
	}

	items := make([]Item, 1000)
	for i := range items {
		items[i] = Item{ID: 1000 - i, Value: "test"}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testItems := make([]Item, len(items))
		copy(testItems, items)
		SortByKey(testItems, func(item Item) int {
			return item.ID
		})
	}
}

// BenchmarkSortByKeyDescUnique 性能测试
func BenchmarkSortByKeyDescUnique(b *testing.B) {
	// 准备测试数据：1000个元素，50%重复
	agents := make([]Agent, 1000)
	for i := 0; i < 1000; i++ {
		agents[i] = Agent{
			ID:     string(rune('A' + (i % 500))), // 50%重复率
			Weight: int32(1000 - i),
			Name:   "Agent",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 每次测试都复制一份数据
		testData := make([]Agent, len(agents))
		copy(testData, agents)

		_ = SortByKeyDescUnique(
			testData,
			func(a Agent) int32 { return a.Weight },
			func(a Agent) string { return a.ID },
		)
	}
}
