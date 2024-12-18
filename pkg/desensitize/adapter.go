/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-18 18:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-18 23:05:29
 * @FilePath: \go-toolbox\pkg\desensitize\adapter.go
 * @Description:
 * 该文件实现了数据脱敏的功能，包括注册脱敏器和执行脱敏操作。
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

import (
	"errors"
	"reflect"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// Desensitizer 接口定义
// 该接口用于定义脱敏器的基本行为，包含一个脱敏方法。
type Desensitizer interface {
	Desensitize(value string) string // 输入一个值，返回脱敏后的值
}

// DefaultDesensitizer 默认脱敏适配器
// 该结构体实现了 Desensitizer 接口，用于处理标准的脱敏逻辑。
type DefaultDesensitizer struct {
	desensitizerType DesensitizeType // 脱敏类型
}

// Desensitize 方法实现
// 根据指定的脱敏类型对输入的值进行脱敏处理。
// @param value: 需要脱敏的字符串值。
// @returns
//   - 返回脱敏后的字符串值。
func (e *DefaultDesensitizer) Desensitize(value string) string {
	return Desensitize(value, e.desensitizerType) // 调用外部库进行脱敏
}

var desensitizers map[string]Desensitizer // 存储注册的脱敏器
var desensitizerMu sync.RWMutex           // 读写锁，用于保护脱敏器的并发访问

// RegisterDesensitizer 注册脱敏器（支持现有和自定义）
// @param desensitizerType: 脱敏器的类型标识。
// @param desensitizer: 实现了 Desensitizer 接口的脱敏器实例。
// @returns
//   - 无返回值。
func RegisterDesensitizer(desensitizerType string, desensitizer Desensitizer) {
	syncx.WithLock(&desensitizerMu, func() {
		if desensitizers == nil {
			desensitizers = make(map[string]Desensitizer) // 初始化脱敏器映射
		}
		desensitizers[desensitizerType] = desensitizer // 注册脱敏器
	})
}

// OperateByRule 根据规则进行脱敏操作
// @param desensitizerType: 脱敏器的类型标识。
// @param in: 需要脱敏的输入值，应该是字符串类型。
// @returns
//   - 返回脱敏后的值和可能的错误。
func OperateByRule(desensitizerType string, in interface{}) (interface{}, error) {
	return syncx.WithLockReturn(&desensitizerMu, func() (interface{}, error) {
		operator, ok := desensitizers[desensitizerType] // 查找对应的脱敏器
		if !ok {
			return nil, errors.New("desensitizer not found") // 未找到对应的脱敏器
		}
		return operator.Desensitize(in.(string)), nil // 执行脱敏操作
	})
}

// Desensitization 执行脱敏操作
// @param obj: 需要进行脱敏的对象，应该是结构体或指向结构体的指针。
// @returns
//   - 返回可能的错误。
func Desensitization(obj interface{}) error {
	rt := reflect.TypeOf(obj)  // 获取对象的类型
	rv := reflect.ValueOf(obj) // 获取对象的值

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem() // 获取指针指向的值
		rt = rt.Elem() // 获取指针指向的类型
	}

	if rv.Kind() != reflect.Struct {
		return errors.New("expected a struct or pointer to a struct") // 仅支持结构体或指针
	}

	// 遍历结构体字段
	for idx := 0; idx < rv.NumField(); idx++ {
		fieldValue := rv.Field(idx) // 获取字段的值
		fieldType := rt.Field(idx)  // 获取字段的类型

		desensitizationTag := fieldType.Tag.Get("desensitize") // 获取脱敏标签

		switch fieldValue.Kind() {
		case reflect.Slice, reflect.Array:
			// 处理切片或数组类型
			if fieldType.Type.Elem().Kind() == reflect.Struct {
				// 如果元素是结构体，递归处理
				for i := 0; i < fieldValue.Len(); i++ {
					if err := Desensitization(fieldValue.Index(i).Addr().Interface()); err != nil {
						return err // 处理错误
					}
				}
			} else {
				// 对于基本类型的切片，执行脱敏
				for i := 0; i < fieldValue.Len(); i++ {
					elemValue := fieldValue.Index(i)
					newValue, err := OperateByRule(desensitizationTag, elemValue.Interface()) // 执行脱敏
					if err == nil {
						elemValue.Set(reflect.ValueOf(newValue)) // 设置脱敏后的值
					}
				}
			}
		case reflect.Struct:
			// 处理结构体类型，递归调用
			if err := Desensitization(fieldValue.Addr().Interface()); err != nil {
				return err // 处理错误
			}
		default:
			// 处理其他基本类型
			if desensitizationTag != "" {
				newValue, err := OperateByRule(desensitizationTag, fieldValue.Interface()) // 执行脱敏
				if err == nil {
					fieldValue.Set(reflect.ValueOf(newValue)) // 设置脱敏后的值
				}
			}
		}
	}

	return nil // 返回 nil 表示成功
}

// 初始化时注册现有的脱敏器
func init() {
	// 注册预定义的脱敏器
	RegisterDesensitizer("email", &DefaultDesensitizer{Email})
	RegisterDesensitizer("phoneNumber", &DefaultDesensitizer{PhoneNumber})
	RegisterDesensitizer("name", &DefaultDesensitizer{ChineseName})
	RegisterDesensitizer("identityCard", &DefaultDesensitizer{IDCard})
	RegisterDesensitizer("mobilePhone", &DefaultDesensitizer{MobilePhone})
	RegisterDesensitizer("address", &DefaultDesensitizer{Address})
	RegisterDesensitizer("password", &DefaultDesensitizer{Password})
	RegisterDesensitizer("carLicense", &DefaultDesensitizer{CarLicense})
	RegisterDesensitizer("bankCard", &DefaultDesensitizer{BankCard})
	RegisterDesensitizer("ipv4", &DefaultDesensitizer{IPV4})
	RegisterDesensitizer("ipv6", &DefaultDesensitizer{IPV6})
}
