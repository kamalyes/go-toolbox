/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-18 16:35:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-18 16:35:27
 * @FilePath: \go-toolbox\pkg\validator\path.go
 * @Description: 路径匹配和验证功能
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import "strings"

// MatchPathInList 检查路径是否匹配列表中的任一模式
//
// 支持的匹配规则：
// 1. 精确匹配：path == pattern
// 2. 前缀匹配：path 以 pattern 开头
//
// 参数：
//   - path: 待检查的路径
//   - patterns: 路径模式列表
//
// 返回：如果匹配任意一个模式则返回 true
//
// 注意：如需更复杂的匹配（如 Glob、正则），请使用 matcher.PathMatcher
//
// 示例：
//
//	patterns := []string{"/api/v1", "/health"}
//	MatchPathInList("/api/v1/users", patterns)  // true (前缀匹配)
//	MatchPathInList("/health", patterns)        // true (精确匹配)
//	MatchPathInList("/admin", patterns)         // false
func MatchPathInList(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if path == pattern || strings.HasPrefix(path, pattern) {
			return true
		}
	}
	return false
}
