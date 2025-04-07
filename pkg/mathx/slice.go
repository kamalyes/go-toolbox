/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 13:55:18
 * @FilePath: \go-toolbox\pkg\mathx\slice.go
 * @Description: 包含与切片相关的通用函数，例如计算最小值和最大值、差集、并集等。
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package mathx

import (
	"errors"
	"math/rand"
	"reflect"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/types"
	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// MinFunc 类型的实现，用于计算最小值
func MinFunc[T types.Numerical](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// MaxFunc 类型的实现，用于计算最大值
func MaxFunc[T types.Numerical](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// SliceMinMax 计算列表中元素的最小值或最大值。
// 接收一个切片和一个 MinMaxFunc 类型的函数，
// 根据提供的函数决定是计算最小值还是最大值。
// 如果列表为空，则返回错误。
func SliceMinMax[T types.Numerical](list []T, f types.MinMaxFunc[T]) (T, error) {
	if len(list) == 0 {
		var zero T
		return zero, errors.New("列表为空") // 返回错误信息
	}

	result := list[0] // 初始化结果为列表的第一个元素
	for _, v := range list[1:] {
		result = f(result, v) // 使用提供的函数更新结果
	}
	return result, nil // 返回最终结果和 nil 错误
}

// SliceFisherYates 洗牌算法打乱数组
func SliceFisherYates[T comparable](slice []T, maxRetries int) error {
	original := make([]T, len(slice))
	copy(original, slice)

	for retries := 0; retries < maxRetries; retries++ {
		for i := len(slice) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)                   // 生成 0 到 i 之间的随机数
			slice[i], slice[j] = slice[j], slice[i] // 交换
		}

		// 校验洗牌后的切片是否与原始切片相同
		if !SliceEqual(original, slice) {
			return nil // 如果不同，返回 nil 表示成功
		}
	}

	// 如果达到最大重试次数，返回错误
	return errors.New("failed to shuffle slice after max retries")
}

// SliceDiffSetSorted 计算两个已排序数组的差集
func SliceDiffSetSorted[T types.Ordered](arr1, arr2 []T) []T {
	diff := []T{}
	i, j := 0, 0

	// 使用双指针遍历两个已排序的数组
	for i < len(arr1) && j < len(arr2) {
		if arr1[i] < arr2[j] {
			diff = append(diff, arr1[i])
			i++
		} else if arr1[i] > arr2[j] {
			diff = append(diff, arr2[j])
			j++
		} else {
			// 遇到相等的元素，跳过
			i++
			j++
		}
	}

	// 添加 arr1 中剩余的元素
	for i < len(arr1) {
		diff = append(diff, arr1[i])
		i++
	}

	// 添加 arr2 中剩余的元素
	for j < len(arr2) {
		diff = append(diff, arr2[j])
		j++
	}

	// 由于我们只想要差集，应该从结果中去掉在另一个数组中存在的元素
	finalDiff := []T{}
	for _, v := range diff {
		if !SliceContains(arr1, v) || !SliceContains(arr2, v) {
			finalDiff = append(finalDiff, v)
		}
	}

	return finalDiff
}

// SliceUnion 计算两个数组的并集。
// 返回一个新的数组，包含所有元素，不包含重复元素。
func SliceUnion[T comparable](arr1, arr2 []T) []T {
	unionMap := make(map[T]struct{}, len(arr1)+len(arr2)) // 使用映射去重
	for _, element := range arr1 {
		unionMap[element] = struct{}{} // 将 arr1 中的元素加入到 unionMap
	}
	for _, element := range arr2 {
		unionMap[element] = struct{}{} // 将 arr2 中的元素加入到 unionMap
	}

	union := make([]T, 0, len(unionMap)) // 创建并集切片
	for key := range unionMap {
		union = append(union, key) // 将 unionMap 中的键转换为切片
	}
	return union
}

// SliceUniq 集合去重
func SliceUniq[T ~[]E, E comparable](list T) T {
	if len(list) == 0 {
		return list
	}

	ret := make(T, 0, len(list))
	m := make(map[E]struct{}, len(list))
	for _, v := range list {
		if _, exists := m[v]; !exists {
			ret = append(ret, v)
			m[v] = struct{}{}
		}
	}
	return ret
}

// SliceDiff 返回两个集合之间的差异
func SliceDiff[T ~[]E, E comparable](list1 T, list2 T) (ret1 T, ret2 T) {
	m1 := make(map[E]struct{}, len(list1))
	m2 := make(map[E]struct{}, len(list2))

	for _, v := range list1 {
		m1[v] = struct{}{}
	}
	for _, v := range list2 {
		m2[v] = struct{}{}
	}

	// 计算差异
	for _, v := range list1 {
		if _, exists := m2[v]; !exists {
			ret1 = append(ret1, v)
		}
	}
	for _, v := range list2 {
		if _, exists := m1[v]; !exists {
			ret2 = append(ret2, v)
		}
	}

	// 确保返回的切片不是 nil
	if ret1 == nil {
		ret1 = make(T, 0)
	}
	if ret2 == nil {
		ret2 = make(T, 0)
	}
	return ret1, ret2
}

// SliceWithout 返回不包括所有给定值的切片
func SliceWithout[T ~[]E, E comparable](list T, exclude ...E) T {
	if len(list) == 0 {
		return list
	}

	m := make(map[E]struct{}, len(exclude))
	for _, v := range exclude {
		m[v] = struct{}{}
	}

	ret := make(T, 0, len(list))
	for _, v := range list {
		if _, exists := m[v]; !exists {
			ret = append(ret, v)
		}
	}
	return ret
}

// SliceIntersect 返回两个集合的交集
func SliceIntersect[T ~[]E, E comparable](list1 T, list2 T) T {
	m := make(map[E]struct{}, len(list1))
	for _, v := range list1 {
		m[v] = struct{}{}
	}

	ret := make(T, 0, len(list1)) // 预分配内存
	for _, v := range list2 {
		if _, exists := m[v]; exists {
			ret = append(ret, v)
		}
	}
	return ret
}

// SliceContains 检查切片中是否包含某个元素。
// 返回布尔值，表示元素是否存在于切片中。
func SliceContains[T types.Ordered](slice []T, element T) bool {
	length := len(slice)

	switch {
	case length <= 1000:
		// 对于小于1000条数据，直接遍历切片
		return containsLinear(slice, element)
	default:
		// 大数据，使用哈希表
		return containsHash(slice, element)
	}
}

// containsLinear 线性查找
func containsLinear[T types.Ordered](slice []T, element T) bool {
	for _, a := range slice {
		if a == element {
			return true // 找到元素，返回 true
		}
	}
	return false // 未找到元素，返回 false
}

// containsHash 哈希表查找
func containsHash[T types.Ordered](slice []T, element T) bool {
	elementMap := make(map[T]struct{})
	for _, a := range slice {
		elementMap[a] = struct{}{}
	}
	_, found := elementMap[element]
	return found // 返回是否找到该元素
}

// SliceHasDuplicates 检查切片中是否存在重复对象。
// 返回布尔值，表示是否存在重复元素。
func SliceHasDuplicates[T comparable](slice []T) bool {
	const chunkSize = 1000 // 每个 goroutine 处理的块大小
	var wg sync.WaitGroup
	m := make(map[T]struct{})
	mu := sync.Mutex{}

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		wg.Add(1)
		go func(subSlice []T) {
			defer wg.Done()
			localMap := make(map[T]struct{})

			for _, v := range subSlice {
				if _, ok := localMap[v]; ok {
					mu.Lock()
					m[v] = struct{}{}
					mu.Unlock()
					return // 找到重复，提前返回
				}
				localMap[v] = struct{}{}
			}
		}(slice[i:end])
	}

	wg.Wait()

	return len(m) > 0 // 如果 map 非空，表示找到重复元素
}

// SliceRemoveEmpty 移除切片中的空对象。
// 返回一个新切片，包含所有非空元素。
func SliceRemoveEmpty[T any](slice []T) []T {
	result := make([]T, 0, len(slice)) // 创建结果切片
	for _, v := range slice {
		if !validator.IsEmptyValue(reflect.ValueOf(v)) {
			result = append(result, v) // 仅添加非空元素
		}
	}
	return result
}

// SliceRemoveDuplicates 移除切片中的重复值。
// 返回一个新切片，包含所有唯一元素。
func SliceRemoveDuplicates[T comparable](numbers []T) []T {
	m := make(map[T]struct{}, len(numbers))     // 预分配 map 的容量
	uniqueNumbers := make([]T, 0, len(numbers)) // 创建唯一元素切片
	for _, num := range numbers {
		if _, exists := m[num]; !exists {
			m[num] = struct{}{}                        // 添加新元素
			uniqueNumbers = append(uniqueNumbers, num) // 仅添加唯一元素
		}
	}
	return uniqueNumbers
}

// SliceRemoveZero 移除切片中的零值。
// 返回一个新切片，包含所有非零元素。
func SliceRemoveZero(arr []int) []int {
	result := make([]int, 0, len(arr)) // 创建结果切片
	for _, val := range arr {
		if val != 0 {
			result = append(result, val) // 仅添加非零元素
		}
	}
	return result
}

// SliceRemoveValue 移除切片中的指定值。
// 返回一个新切片，包含所有非指定元素。
func SliceRemoveValue[T comparable](arr []T, value T) []T {
	result := make([]T, 0, len(arr)) // 创建结果切片
	for _, val := range arr {
		if val != value {
			result = append(result, val) // 仅添加非指定元素
		}
	}
	return result
}

// SliceChunk 将一个切片分割成多个子切片。
// size 参数指定每个子切片的大小。
// 返回一个包含子切片的切片。
func SliceChunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil // 如果 size <= 0，则返回 nil
	}

	var batches [][]T // 创建子切片切片
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice) // 确保不超出边界
		}
		batches = append(batches, slice[i:end]) // 切片而不复制
	}
	return batches
}

// SliceEqual 比较两个切片是否相等，支持任何类型
func SliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// InsertionSort 对小规模数组使用插入排序
func InsertionSort(arr []int) {
	if len(arr) < 2 {
		return
	}
	low := 0
	high := len(arr) - 1

	for currentIndex := low + 1; currentIndex <= high; currentIndex++ {
		key := arr[currentIndex]
		sortedIndex := currentIndex - 1
		for sortedIndex >= low && arr[sortedIndex] > key {
			arr[sortedIndex+1] = arr[sortedIndex]
			sortedIndex--
		}
		arr[sortedIndex+1] = key
	}
}

// QuickSort 实现快速排序算法
func QuickSort(arr []int, low, high int) {
	if low < high {
		// 获取分区索引
		pi := partition(arr, low, high)

		// 递归排序分区
		QuickSort(arr, low, pi-1)  // 排序基准左侧
		QuickSort(arr, pi+1, high) // 排序基准右侧
	}
}

// partition 进行分区操作
func partition(arr []int, low, high int) int {
	pivot := arr[high]     // 选择最后一个元素作为基准
	sortedIndex := low - 1 // 小于基准的元素的索引

	for currentIndex := low; currentIndex < high; currentIndex++ {
		if arr[currentIndex] < pivot { // 如果当前元素小于基准
			sortedIndex++                                                             // 增加小于基准的元素索引
			arr[sortedIndex], arr[currentIndex] = arr[currentIndex], arr[sortedIndex] // 交换
		}
	}
	// 将基准放到正确的位置
	arr[sortedIndex+1], arr[high] = arr[high], arr[sortedIndex+1]
	return sortedIndex + 1 // 返回基准的索引
}

// BubbleSort 实现冒泡排序算法
func BubbleSort(arr []int) {
	n := len(arr)

	for currentIndex := 0; currentIndex < n-1; currentIndex++ {
		for sortedIndex := 0; sortedIndex < n-currentIndex-1; sortedIndex++ {
			if arr[sortedIndex] > arr[sortedIndex+1] { // 如果当前元素大于下一个元素，则交换
				arr[sortedIndex], arr[sortedIndex+1] = arr[sortedIndex+1], arr[sortedIndex]
			}
		}
	}
}
