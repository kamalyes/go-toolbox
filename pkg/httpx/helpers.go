/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-05 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 11:32:06
 * @FilePath: \go-toolbox\pkg\httpx\helpers.go
 * @Description: HTTP 辅助方法
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"net/http"
	"reflect"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// BuildParams 构建请求参数的辅助方法
// 基础参数通过 base 传入，可选参数通过 opts 传入
func BuildParams(base map[string]string, opts ...func(map[string]string)) map[string]string {
	params := make(map[string]string, len(base))
	for k, v := range base {
		params[k] = v
	}

	for _, opt := range opts {
		if opt != nil {
			opt(params)
		}
	}

	return params
}

// WithParam 条件添加参数
func WithParam(condition bool, key, value string) func(map[string]string) {
	return func(params map[string]string) {
		mathx.IfExec(condition, func() {
			params[key] = value
		})
	}
}

// WithParamNotEmpty 非空时添加参数
func WithParamNotEmpty(key, value string) func(map[string]string) {
	return func(params map[string]string) {
		mathx.IfExec(!validator.IsEmptyValue(reflect.ValueOf(value)), func() {
			params[key] = value
		})
	}
}

// GetUserID 从 HTTP 请求中获取用户ID
// 优先从上下文获取，然后从请求头获取
// contextKey: 上下文中存储用户ID的键
// headerKey: 请求头中存储用户ID的键
func GetUserID(r *http.Request, contextKey interface{}, headerKey string) string {
	// 优先从上下文获取
	if contextKey != nil {
		if userID := r.Context().Value(contextKey); userID != nil {
			if uid, ok := userID.(string); ok {
				return uid
			}
		}
	}

	// 从请求头获取
	if headerKey != "" {
		if userID := r.Header.Get(headerKey); userID != "" {
			return userID
		}
	}

	return ""
}
