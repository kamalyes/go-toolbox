/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-17 13:55:21
 * @FilePath: \go-toolbox\pkg\mathx\slice.go
 * @Description: 包含与切片相关的通用函数，例如计算最小值和最大值、差集、并集等
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package mathx

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/kamalyes/go-toolbox/pkg/types"
	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// SliceMinMax 计算列表中元素的最小值或最大值
// 接收一个切片和一个 MinMaxFunc 类型的函数，
// 根据提供的函数决定是计算最小值还是最大值
// 如果列表为空，则返回错误
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

// SliceUnion 计算两个数组的并集
// 返回一个新的数组，包含所有元素，不包含重复元素
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

// SliceContains 检查切片中是否包含某个元素
// 返回布尔值，表示元素是否存在于切片中
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

// SliceContainsComparable 检查切片中是否包含某个元素（支持comparable类型）
// 适用于枚举、字符串等comparable类型，比SliceContains支持更广泛的类型
func SliceContainsComparable[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// SliceHasDuplicates 检查切片中是否存在重复对象
// 返回布尔值，表示是否存在重复元素
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

// SliceRemoveEmpty 移除切片中的空对象
// 返回一个新切片，包含所有非空元素
func SliceRemoveEmpty[T any](slice []T) []T {
	result := make([]T, 0, len(slice)) // 创建结果切片
	for _, v := range slice {
		if !validator.IsEmptyValue(reflect.ValueOf(v)) {
			result = append(result, v) // 仅添加非空元素
		}
	}
	return result
}

// SliceRemoveDuplicates 移除切片中的重复值
// 返回一个新切片，包含所有唯一元素
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

// SliceRemove 根据给定的条件函数移除切片中的元素
// 返回一个新切片，包含所有满足条件的元素
func SliceRemove[T any](arr []T, condition func(T) bool) []T {
	result := make([]T, 0, len(arr)) // 创建结果切片
	for _, val := range arr {
		if condition(val) {
			result = append(result, val) // 仅添加满足条件的元素
		}
	}
	return result
}

// SliceRemoveZero 移除切片中的零值
// 返回一个新切片，包含所有非零元素
func SliceRemoveZero[T comparable](arr []T) []T {
	return SliceRemove(arr, func(val T) bool {
		return val != *new(T) // 仅保留非零元素，使用零值判断
	})
}

// SliceRemoveValue 移除切片中的指定值
// 返回一个新切片，包含所有非指定元素
func SliceRemoveValue[T comparable](arr []T, value T) []T {
	return SliceRemove(arr, func(val T) bool {
		return val != value // 仅保留非指定元素
	})
}

// SliceChunk 将一个切片分割成多个子切片
// size 参数指定每个子切片的大小
// 返回一个包含子切片的切片
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

// RepeatField 返回一个长度为 count 的切片，
// 切片中的每个元素都是传入的 field 值的副本
func RepeatField[T any](field T, count int) []T {
	s := make([]T, count) // 创建一个长度为 count 的切片，元素类型为 T
	for i := range s {
		s[i] = field // 将 field 赋值给切片的每个元素
	}
	return s // 返回填充好的切片
}

// SliceChain 是一个支持链式调用的泛型切片操作结构体，带并发安全锁保护
// 通过链式调用可以方便地对切片进行追加、去重、排序、过滤等操作
type SliceChain[T comparable] struct {
	mu   sync.RWMutex // 读写锁
	data []T
}

// FromSlice 根据普通切片创建一个新的 SliceChain 实例
// Params:
//   - slice: 需要转换成链式操作结构的普通切片，类型为 []T
//
// Returns:
//   - 返回一个指向新创建的 SliceChain 实例的指针，内部包含了传入的切片数据副本
func FromSlice[T comparable](slice []T) *SliceChain[T] {
	sc := &SliceChain[T]{}
	if len(slice) > 0 {
		sc.data = append(sc.data, slice...)
	}
	return sc
}

// Append 追加元素到当前切片，支持链式调用
// Params:
//   - elements: 可变参数，表示需要追加到当前切片中的元素，类型为 T
//
// Returns:
//   - 返回当前 SliceChain 实例指针，方便链式调用
//
// Examples:
//
//	sc.Append(1, 2, 3)
func (sc *SliceChain[T]) Append(elements ...T) *SliceChain[T] {
	return syncx.WithLockReturnValue(&sc.mu, func() *SliceChain[T] {
		sc.data = append(sc.data, elements...)
		return sc
	})
}

// Uniq 去重，移除切片中重复的元素，保持元素顺序，支持链式调用
// Returns:
//   - 返回当前 SliceChain 实例指针，去重后的数据保存在内部切片中
func (sc *SliceChain[T]) Uniq() *SliceChain[T] {
	return syncx.WithLockReturnValue(&sc.mu, func() *SliceChain[T] {
		sc.data = SliceRemoveDuplicates(sc.data)
		return sc
	})
}

// RemoveValue 移除切片中所有等于指定值的元素，支持链式调用
// Params：
//   - value: 需要从切片中移除的元素值，类型为 T
//
// Returns:
//   - 返回当前 SliceChain 实例指针，移除指定值后的数据保存在内部切片中
func (sc *SliceChain[T]) RemoveValue(value T) *SliceChain[T] {
	return syncx.WithLockReturnValue(&sc.mu, func() *SliceChain[T] {
		n := 0
		for _, v := range sc.data {
			if v != value {
				sc.data[n] = v
				n++
			}
		}
		sc.data = sc.data[:n]
		return sc
	})
}

// RemoveEmpty 移除“空值”元素，空值定义为元素等于类型零值，支持链式调用
// 适用于数字、字符串、指针等类型的零值判断
// Returns:
//   - 返回当前 SliceChain 实例指针，移除零值元素后的数据保存在内部切片中
func (sc *SliceChain[T]) RemoveEmpty() *SliceChain[T] {
	return syncx.WithLockReturnValue(&sc.mu, func() *SliceChain[T] {
		result := sc.data[:0]
		for _, v := range sc.data {
			if !validator.IsCEmpty(v) {
				result = append(result, v)
			}
		}
		sc.data = result
		return sc
	})
}

// Filter 根据传入的过滤函数 f，保留满足条件的元素，支持链式调用
// Params：
//   - f: 过滤函数，接收一个元素 T，返回 bool，返回 true 表示保留该元素，false 表示过滤掉
//
// Returns:
//   - 返回当前 SliceChain 实例指针，过滤后的数据保存在内部切片中
//
// Examples:
//
//	sc.Filter(func(x int) bool { return x%2 == 0 })
func (sc *SliceChain[T]) Filter(f func(T) bool) *SliceChain[T] {
	return syncx.WithLockReturnValue(&sc.mu, func() *SliceChain[T] {
		result := sc.data[:0] // 利用切片复用内存，避免额外分配
		for _, v := range sc.data {
			if f(v) {
				result = append(result, v)
			}
		}
		sc.data = result
		return sc
	})
}

// Sort 使用传入的 less 函数对切片进行排序，支持链式调用
// Params:
//   - less: 比较函数，接收两个元素 a, b，返回 bool，返回 true 表示 a < b
//
// Returns:
//   - 返回当前 SliceChain 实例指针，排序后的数据保存在内部切片中
//
// Examples:
//
//	sc.Sort(func(a, b int) bool { return a < b })
func (sc *SliceChain[T]) Sort(less func(a, b T) bool) *SliceChain[T] {
	return syncx.WithLockReturnValue(&sc.mu, func() *SliceChain[T] {
		sort.Slice(sc.data, func(i, j int) bool {
			return less(sc.data[i], sc.data[j])
		})
		return sc
	})
}

// Data 返回当前链式操作后的切片数据，方便与普通切片交互
// Returns:
//   - 返回当前内部切片的引用，类型为 []T
//
// 注意：返回的是内部切片的引用，修改返回值会影响 SliceChain 内部数据
func (sc *SliceChain[T]) Data() []T {
	return syncx.WithRLockReturnValue(&sc.mu, func() []T {
		return IF(len(sc.data) > 0, sc.data, []T{}) // 返回非nil空切片，避免断言错误
	})
}

// String 实现 fmt.Stringer 接口，方便打印 SliceChain 内容
// Returns:
//   - 返回当前切片数据的字符串表示
func (sc *SliceChain[T]) String() string {
	return syncx.WithRLockReturnValue(&sc.mu, func() string {
		return fmt.Sprintf("%v", sc.data)
	})
}
