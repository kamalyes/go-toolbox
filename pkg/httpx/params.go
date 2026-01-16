/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-16 13:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-16 13:15:00
 * @FilePath: \go-toolbox\pkg\httpx\params.go
 * @Description: 参数构建器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"reflect"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// ParamsBuilder 参数构建器，支持链式调用
type ParamsBuilder struct {
	params map[string]string
}

// NewParams 创建新的参数构建器
func NewParams() *ParamsBuilder {
	return &ParamsBuilder{
		params: make(map[string]string),
	}
}

// NewParamsWithBase 创建带基础参数的构建器
func NewParamsWithBase(base map[string]string) *ParamsBuilder {
	params := make(map[string]string, len(base))
	for k, v := range base {
		params[k] = v
	}
	return &ParamsBuilder{params: params}
}

// Set 设置参数（无条件）
func (b *ParamsBuilder) Set(key, value string) *ParamsBuilder {
	b.params[key] = value
	return b
}

// SetIf 条件设置参数
func (b *ParamsBuilder) SetIf(condition bool, key, value string) *ParamsBuilder {
	mathx.IfExec(condition, func() {
		b.params[key] = value
	})
	return b
}

// SetNotEmpty 非空时设置参数（使用 validator.IsEmptyValue 判断）
func (b *ParamsBuilder) SetNotEmpty(key, value string) *ParamsBuilder {
	mathx.IfExec(!validator.IsEmptyValue(reflect.ValueOf(value)), func() {
		b.params[key] = value
	})
	return b
}

// SetMultiple 批量设置参数
func (b *ParamsBuilder) SetMultiple(params map[string]string) *ParamsBuilder {
	for k, v := range params {
		b.params[k] = v
	}
	return b
}

// Build 构建并返回参数 map
func (b *ParamsBuilder) Build() map[string]string {
	return b.params
}
