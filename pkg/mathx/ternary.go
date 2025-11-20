/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-21 19:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 11:19:04
 * @FilePath: \go-toolbox\pkg\mathx\ternary.go
 * @Description: 包提供了一组基于 Go 泛型实现的三元运算及条件执行函数，支持同步、异步、带错误处理等多种场景
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"encoding/json"
	"time"
)

// IF 实现三元运算，使用泛型 T
// 根据布尔条件 condition 返回 trueVal 或 falseVal
// 作用类似于三元表达式 condition ? trueVal : falseVal
func IF[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// DoFunc 是一个泛型函数类型，表示无参返回 T 类型值的函数
// 用于延迟执行并返回结果
type DoFunc[T any] func() T

// IfDo 根据条件 condition 决定是否执行函数 do 并返回结果
// 如果 condition 为 true，执行 do() 并返回其结果；否则返回默认值 defaultVal
func IfDo[T any](condition bool, do DoFunc[T], defaultVal T) T {
	if condition {
		return do() // 执行函数并返回结果
	}
	return defaultVal // 返回默认值
}

// IfDoAF 根据条件 condition 决定执行 do 或 defaultFunc 函数并返回结果
// 如果 condition 为 true，执行 do()；否则执行 defaultFunc()
func IfDoAF[T any](condition bool, do DoFunc[T], defaultFunc DoFunc[T]) T {
	return IfDo(condition, do, defaultFunc()) // 复用 IfDo 函数
}

// DoFuncWithError 是一个泛型函数类型，表示无参返回 (T, error) 的函数
// 用于执行可能返回错误的延迟操作
type DoFuncWithError[T any] func() (T, error)

// IfDoWithError 根据条件 condition 执行带错误返回的函数 do
// 如果 condition 为 true，执行 do() 并返回结果与错误；否则返回默认值 defaultVal 和 nil 错误
func IfDoWithError[T any](condition bool, do DoFuncWithError[T], defaultVal T) (T, error) {
	if condition {
		return do()
	}
	return defaultVal, nil
}

// IfDoAsync 支持异步执行延迟函数 do，返回结果的通道
// 根据条件 condition 决定是否执行 do()，否则返回默认值 defaultVal
// 结果通过带缓冲的通道返回，避免阻塞
func IfDoAsync[T any](condition bool, do DoFunc[T], defaultVal T) <-chan T {
	ch := make(chan T, 1)
	go func() {
		if condition {
			ch <- do()
		} else {
			ch <- defaultVal
		}
		close(ch)
	}()
	return ch
}

// IfDoAsyncWithTimeout 异步执行延迟函数 do，支持超时控制
// condition 为 true 时执行 do()，否则返回默认值 defaultVal
// timeoutMs 指定超时时间（毫秒），超时则返回类型 T 的零值
// 返回一个通道，异步获取结果
func IfDoAsyncWithTimeout[T any](condition bool, do DoFunc[T], defaultVal T, timeoutMs int) <-chan T {
	ch := make(chan T, 1)
	go func() {
		if condition {
			ch <- do()
		} else {
			ch <- defaultVal
		}
		close(ch)
	}()

	out := make(chan T, 1)
	go func() {
		select {
		case v := <-ch:
			out <- v
		case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
			var zero T
			out <- zero
		}
		close(out)
	}()
	return out
}

// IfElse 实现多条件链式判断，安全泛型版本
// conds 和 values 长度必须相等，依次判断 conds 中的条件
// 返回第一个为 true 的对应 values 元素；如果没有满足条件，返回 defaultVal
// 适合实现类似 if ... else if ... else 的多分支逻辑
func IfElse[T any](conds []bool, values []T, defaultVal T) T {
	if len(conds) != len(values) {
		panic("IfElse: 条件和值长度必须相等")
	}
	for i, cond := range conds {
		if cond {
			return values[i]
		}
	}
	return defaultVal
}

// ConditionValue 是一个泛型结构体，封装一个条件和对应的值
// 用于链式多条件判断，方便组合使用
type ConditionValue[T any] struct {
	Cond  bool // 条件
	Value T    // 条件为 true 时对应的返回值
}

// IfChain 根据传入的 ConditionValue 切片，依次判断每个条件
// 返回第一个条件为 true 的对应值；如果所有条件均为 false，返回默认值 defaultVal
// 实现类似 if-else if-else 的多条件链判断，且支持任意泛型类型
//
// 示例：
//
//	pairs := []ConditionValue[int]{
//	    {Cond: x > 0, Value: 1},
//	    {Cond: x == 0, Value: 0},
//	    {Cond: x < 0, Value: -1},
//	}
//	result := IfChain(pairs, 999) // 若都不满足，返回 999
func IfChain[T any](pairs []ConditionValue[T], defaultVal T) T {
	for _, pair := range pairs {
		if pair.Cond {
			return pair.Value
		}
	}
	return defaultVal
}

// ResultWithError 泛型结构体，封装结果和错误
type ResultWithError[T any] struct {
	Result T
	Err    error
}

// IfDoWithErrorAsync 支持异步执行带错误返回的函数 do
// 根据条件 condition 决定是否执行 do()，否则返回默认值 defaultVal 和 nil 错误
// 返回一个通道，通道元素是 ResultWithError[T] 类型
func IfDoWithErrorAsync[T any](condition bool, do DoFuncWithError[T], defaultVal T) <-chan ResultWithError[T] {
	ch := make(chan ResultWithError[T], 1)
	go func() {
		if condition {
			r, e := do()
			ch <- ResultWithError[T]{r, e}
		} else {
			ch <- ResultWithError[T]{defaultVal, nil}
		}
		close(ch)
	}()
	return ch
}

// ReturnIfErr 简化错误检查和返回
// 如果 err 不为 nil，返回 T 类型的零值和 err；否则返回 val 和 nil
func ReturnIfErr[T any](val T, err error) (T, error) {
	if err != nil {
		var zero T
		return zero, err
	}
	return val, nil
}

// IfDoWithErrorDefault 根据条件 condition 执行带错误返回的函数 do
// 如果 condition 为 false 或 do 返回错误，则返回 defaultVal
func IfDoWithErrorDefault[T any](condition bool, do DoFuncWithError[T], defaultVal T) T {
	if !condition {
		return defaultVal
	}
	val, err := do()
	if err != nil {
		return defaultVal
	}
	return val
}

// IfCall 根据布尔条件 condition，选择性调用对应的回调函数
//
// Params
//   - condition: 判断条件，true 时调用 onTrue，false 时调用 onFalse
//   - result: 泛型参数，传递给回调函数的结果值
//   - err: 错误信息，传递给回调函数
//   - onTrue: 条件为 true 时调用的回调函数，接收 (result, err)
//   - onFalse: 条件为 false 时调用的回调函数，接收 (result, err)
//
// 函数逻辑：
//  1. 根据 condition 选择要调用的回调函数 cb（onTrue 或 onFalse）
//  2. 如果 cb 不为 nil，则调用 cb(result, err)
//  3. 如果 cb 为 nil，则跳过调用，避免空指针异常
//
// 作用：简化根据条件调用不同回调的代码，避免重复写 if-else 和 nil 判断，提高代码简洁性和安全性
func IfCall[T any](condition bool, result T, err error, onTrue func(T, error), onFalse func(T, error)) {
	if condition && onTrue != nil {
		onTrue(result, err)
		return
	}
	if !condition && onFalse != nil {
		onFalse(result, err)
	}
}

// IfExec 根据条件执行副作用操作（无返回值）
// 适用于只需要执行代码块，不需要返回值的场景
//
// Params
//   - condition: 判断条件，true 时执行 action
//   - action: 条件为 true 时执行的函数
//
// 示例：
//
//	mathx.IfExec(user != nil, func() {
//	    log.Printf("User: %s", user.Name)
//	})
func IfExec(condition bool, action func()) {
	if condition && action != nil {
		action()
	}
}

// IfExecElse 根据条件执行不同的副作用操作
// 类似三元运算符，但用于执行代码块而非返回值
//
// Params
//   - condition: 判断条件
//   - onTrue: 条件为 true 时执行的函数
//   - onFalse: 条件为 false 时执行的函数
//
// 示例：
//
//	mathx.IfExecElse(err == nil,
//	    func() { log.Info("Success") },
//	    func() { log.Error("Failed: " + err.Error()) },
//	)
func IfExecElse(condition bool, onTrue func(), onFalse func()) {
	if condition {
		if onTrue != nil {
			onTrue()
		}
	} else {
		if onFalse != nil {
			onFalse()
		}
	}
}

// MarshalJSONOrDefault 将任意值序列化为 JSON 字符串，失败或空值时返回默认值
// 适用于需要确保 JSON 字段不为空的场景（如 MySQL JSON 列）
//
// 参数：
//   - value: 待序列化的值（可以是 map、struct、slice 等任意类型）
//   - defaultVal: 序列化失败或值为空时返回的默认值（通常为 "{}" 或 "[]"）
//
// 返回：
//   - JSON 字符串或默认值
//
// 示例：
//
//	map[string]string 序列化
//	extra := map[string]string{"key": "value"}
//	json := mathx.MarshalJSONOrDefault(extra, "{}") // 返回 {"key":"value"}
//
//	空 map 返回默认值
//	empty := map[string]string{}
//	json := mathx.MarshalJSONOrDefault(empty, "{}") // 返回 {}
//
//	nil 值返回默认值
//	json := mathx.MarshalJSONOrDefault(nil, "{}") // 返回 {}
//
//	序列化失败返回默认值
//	invalid := make(chan int) // channel 不能序列化
//	json := mathx.MarshalJSONOrDefault(invalid, "{}") // 返回 {}
func MarshalJSONOrDefault(value any, defaultVal string) string {
	// 处理 nil 值
	if value == nil {
		return defaultVal
	}

	// 处理 map 类型的空值检测
	switch v := value.(type) {
	case map[string]string:
		if len(v) == 0 {
			return defaultVal
		}
	case map[string]any:
		if len(v) == 0 {
			return defaultVal
		}
	}

	// 尝试序列化
	bytes, err := json.Marshal(value)
	if err != nil {
		return defaultVal
	}

	result := string(bytes)

	// 处理 JSON null 的情况
	if result == "null" {
		return defaultVal
	}

	return result
}
