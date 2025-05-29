/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-05-29 13:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-29 13:55:21
 * @FilePath: \go-toolbox\pkg\mathx\filter.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package mathx

// MapSliceByKey 将一个切片 slice 转换为一个 map，
// 键由 keyFunc 从切片元素中提取，值为对应的切片元素。
//
// 泛型参数说明：
//   T - 切片元素的类型，支持任意类型。
//   K - map 的键类型，必须是可比较类型（comparable），
//       这是 Go map 键的要求。
//
// 参数说明：
//   slice - 输入的切片，类型为 []T。
//   keyFunc - 函数，接收一个 T 类型元素，返回对应的键 K。
//
// 返回值说明：
//   返回一个 map[K]T，键是 keyFunc 生成的键，值是对应的切片元素。
//   如果 slice 中有重复的键，后出现的元素会覆盖前面的。
//
// 例子：
//   users := []User{{ID:1, Name:"A"}, {ID:2, Name:"B"}}
//   userMap := MapSliceByKey(users, func(u User) int { return u.ID })
//   userMap[1] == User{ID:1, Name:"A"}
//
func MapSliceByKey[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T {
	// 预分配 map 容量为切片长度，避免多次扩容
	result := make(map[K]T, len(slice))

	// 遍历切片元素
	for _, item := range slice {
		// 计算当前元素对应的键
		key := keyFunc(item)
		// 将元素存入 map，以 key 为索引
		result[key] = item
	}

	// 返回构造好的 map
	return result
}

// FilterSliceByFunc 是一个通用的切片过滤函数。
// 它接收一个切片 slice 和一个判断函数 predicate，
// 返回一个新的切片，包含所有满足 predicate 条件的元素。
//
// 泛型参数说明：
//   T 表示切片元素的类型，可以是任意类型。
//
// 参数说明：
//   slice - 输入的切片，类型为 []T。
//   predicate - 一个函数，接收一个 T 类型元素，返回 bool，表示该元素是否满足条件。
//
// 返回值说明：
//   返回一个新的切片，包含所有满足 predicate 条件的元素，顺序与原切片一致。
//
// 示例：
//   results := []ABC{...}
//   frontFaces := FilterSliceByFunc(results, func(d ABC) bool {
//       return d.key == "123"
//   })
func FilterSliceByFunc[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	// 遍历切片中的每个元素
	for _, item := range slice {
		// 判断当前元素是否满足条件
		if predicate(item) {
			// 满足条件则追加到结果切片
			result = append(result, item)
		}
	}
	// 返回筛选后的切片
	return result
}

// Predicate 是判断条件的函数类型，输入元素，返回是否满足条件
type Predicate[T any] func(item T) bool

// FilterCallback 是过滤后回调函数类型，输入元素指针
type FilterCallback[T any] func(item *T)

// combineStrategy 定义策略接口，封装不同的合并判断逻辑
type combineStrategy[T any] interface {
	// Match 根据策略判断该元素是否满足所有谓词条件
	Match(item *T, preds []Predicate[T]) bool
}

// andStrategy 实现“全部为真”策略
type andStrategy[T any] struct{}

// Match 遇到false立即返回false，遍历完都没false返回true
func (s andStrategy[T]) Match(item *T, preds []Predicate[T]) bool {
	// 遇到 false 就短路返回 false
	return matchGeneric(item, preds,
		func(b bool) bool { return !b }, // b == false 时短路
		false,                           // 短路时返回 false
		true,                            // 遍历完没短路返回 true
	)
}

// orStrategy 实现“任意为真”策略
type orStrategy[T any] struct{}

// Match 遇到true立即返回true，遍历完都没true返回false
func (s orStrategy[T]) Match(item *T, preds []Predicate[T]) bool {
	// 遇到 true 就短路返回 true
	return matchGeneric(item, preds,
		func(b bool) bool { return b }, // b == true 时短路
		true,                           // 短路时返回 true
		false,                          // 遍历完没短路返回 false
	)
}

func matchGeneric[T any](
	item *T,
	preds []Predicate[T],
	shortCircuitCond func(bool) bool, // 如果返回true表示满足短路条件，立即返回对应结果
	shortCircuitResult bool, // 达到短路条件时返回的结果
	defaultResult bool, // 遍历完毕未触发短路时返回的结果
) bool {
	for _, pred := range preds {
		var res bool
		if pred == nil {
			res = true // nil谓词视为true
		} else {
			res = pred(*item)
		}
		if shortCircuitCond(res) {
			return shortCircuitResult
		}
	}
	return defaultResult
}

// customStrategy 支持自定义合并函数的策略
type customStrategy[T any] struct {
	combine func([]bool) bool // 用户自定义的合并逻辑
}

// Match 先计算每个谓词结果，再调用自定义合并函数决定返回值
func (s customStrategy[T]) Match(item *T, preds []Predicate[T]) bool {
	results := make([]bool, len(preds))
	for i, pred := range preds {
		if pred == nil {
			results[i] = true // nil谓词视为true
		} else {
			results[i] = pred(*item)
		}
	}
	return s.combine(results)
}

// SliceFilter 泛型切片过滤器，支持多谓词和策略模式
type SliceFilter[T any] struct {
	slice      []T                 // 原始切片
	predicates []Predicate[T]      // 过滤条件谓词集合
	strategy   combineStrategy[T]  // 合并策略接口
	onMatch    []FilterCallback[T] // 匹配成功回调
	onNotMatch []FilterCallback[T] // 不匹配回调
	result     []T                 // 过滤结果缓存
	filtered   bool                // 是否已过滤过，避免重复计算
}

// NewSliceFilter 创建新过滤器实例，默认使用and策略
func NewSliceFilter[T any](slice []T) *SliceFilter[T] {
	return &SliceFilter[T]{
		slice:    slice,
		strategy: andStrategy[T]{}, // 默认全部条件为真
	}
}

// UseAnd 设置策略为“全部为真”
func (f *SliceFilter[T]) UseAnd() *SliceFilter[T] {
	f.strategy = andStrategy[T]{}
	return f
}

// UseOr 设置策略为“任意为真”
func (f *SliceFilter[T]) UseOr() *SliceFilter[T] {
	f.strategy = orStrategy[T]{}
	return f
}

// Condition 添加一个或多个过滤条件，支持链式调用
func (f *SliceFilter[T]) Condition(preds ...Predicate[T]) *SliceFilter[T] {
	f.predicates = append(f.predicates, preds...) // 将新的条件追加到条件列表
	return f
}

// OnMatch 添加满足筛选条件时的回调函数，支持多次调用添加多个回调。
// 回调函数接受元素指针，可以修改元素。
func (f *SliceFilter[T]) OnMatch(cb FilterCallback[T]) *SliceFilter[T] {
	f.onMatch = append(f.onMatch, cb)
	return f
}

// OnNotMatch 添加不满足筛选条件时的回调函数，支持多次调用添加多个回调。
func (f *SliceFilter[T]) OnNotMatch(cb FilterCallback[T]) *SliceFilter[T] {
	f.onNotMatch = append(f.onNotMatch, cb)
	return f
}

// UseCustom 设置自定义合并策略
func (f *SliceFilter[T]) UseCustom(combine func([]bool) bool) *SliceFilter[T] {
	f.strategy = customStrategy[T]{combine}
	return f
}

// Result 执行过滤操作，返回符合条件的元素切片
func (f *SliceFilter[T]) Result() []T {
	if f.filtered {
		// 已过滤过，直接返回缓存结果，避免重复计算
		return f.result
	}
	f.filtered = true

	var res []T

	// 辅助函数，执行回调函数数组
	callCallbacks := func(cbs []FilterCallback[T], item *T) {
		for _, cb := range cbs {
			if cb != nil {
				cb(item)
			}
		}
	}

	for i := range f.slice {
		item := &f.slice[i]

		if len(f.predicates) == 0 {
			// 无谓词，默认全部元素都匹配
			callCallbacks(f.onMatch, item)
			res = append(res, *item)
			continue
		}

		// 使用策略判断该元素是否满足谓词条件
		match := f.strategy.Match(item, f.predicates)

		if match {
			callCallbacks(f.onMatch, item)
			res = append(res, *item)
		} else {
			callCallbacks(f.onNotMatch, item)
		}
	}

	f.result = res
	return res
}

type FindResult[T any, V any] struct {
	item  *T
	val   V
	found bool
	// 允许写回 map 的引用
	dataMap map[string]V
	key     string
	stopped bool // 是否停止链式调用
}

func (r *FindResult[T, V]) Item() *T {
	return r.item
}

// 查找函数，返回 FindResult
func FindUpdate[T any, V any](
	item *T,
	dataMap map[string]V,
	getKey func(*T, ...any) string,
	keyArgs ...any,
) *FindResult[T, V] {
	key := getKey(item, keyArgs...)
	val, found := dataMap[key]
	return &FindResult[T, V]{item, val, found, dataMap, key, false}
}

// 如果找到，执行回调，支持修改 item、val
func (r *FindResult[T, V]) IfFound(f func(*T, *V)) *FindResult[T, V] {
	if r.stopped {
		return r
	}
	if r.found {
		f(r.item, &r.val)
		// 写回 map
		r.dataMap[r.key] = r.val
	}
	return r
}

// 如果没找到，执行回调
func (r *FindResult[T, V]) OrElse(f func(*T)) *FindResult[T, V] {
	if r.stopped {
		return r
	}
	if !r.found {
		f(r.item)
	}
	return r
}

// 统一处理找到和没找到
func (r *FindResult[T, V]) Then(onFound func(*T, *V), onNotFound func(*T)) *FindResult[T, V] {
	if r.stopped {
		return r
	}
	if r.found {
		onFound(r.item, &r.val)
		r.dataMap[r.key] = r.val
	} else {
		onNotFound(r.item)
	}
	return r
}

// 支持链式调用中断
func (r *FindResult[T, V]) Stop() *FindResult[T, V] {
	r.stopped = true
	return r
}

// 支持条件执行，类似 if
func (r *FindResult[T, V]) When(cond bool, f func(*FindResult[T, V])) *FindResult[T, V] {
	if r.stopped {
		return r
	}
	if cond {
		f(r)
	}
	return r
}

// 支持无条件执行，方便写一些操作
func (r *FindResult[T, V]) Do(f func(*FindResult[T, V])) *FindResult[T, V] {
	if r.stopped {
		return r
	}
	f(r)
	return r
}
