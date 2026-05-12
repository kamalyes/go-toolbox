/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-25 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\helpers.go
 * @Description: 验证器通用辅助函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"regexp"
	"sync"
)

// 正则缓存 - 避免重复编译
var (
	regexCache   = make(map[string]*regexp.Regexp)
	regexCacheMu sync.RWMutex
)

// GetCompiledRegex 获取编译的正则（带缓存）- 公共函数供其他模块使用
func GetCompiledRegex(pattern string) (*regexp.Regexp, error) {
	regexCacheMu.RLock()
	re, exists := regexCache[pattern]
	regexCacheMu.RUnlock()

	if exists {
		return re, nil
	}

	regexCacheMu.Lock()
	defer regexCacheMu.Unlock()

	// 双重检查
	if re, exists := regexCache[pattern]; exists {
		return re, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	regexCache[pattern] = re
	return re, nil
}

// ClearRegexCache 清空正则缓存（用于测试或内存释放）
func ClearRegexCache() {
	regexCacheMu.Lock()
	defer regexCacheMu.Unlock()
	regexCache = make(map[string]*regexp.Regexp)
}