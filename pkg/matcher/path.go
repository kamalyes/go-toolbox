/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-05 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-05 00:00:00
 * @FilePath: \go-toolbox\pkg\matcher\path.go
 * @Description: 路径匹配增强功能 - 为 HTTP 路由、文件路径等提供高性能匹配
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package matcher

import (
	"regexp"
	"strings"
	"sync"
)

// PathMatcherType 路径匹配器类型
type PathMatcherType int

const (
	// PathMatchExact 精确匹配
	PathMatchExact PathMatcherType = iota
	// PathMatchPrefix 前缀匹配
	PathMatchPrefix
	// PathMatchSuffix 后缀匹配
	PathMatchSuffix
	// PathMatchGlob Glob 模式匹配 (如 /api/*/users)
	PathMatchGlob
	// PathMatchRegex 正则表达式匹配
	PathMatchRegex
	// PathMatchContains 包含匹配
	PathMatchContains
)

// PathMatcher 路径匹配器
type PathMatcher struct {
	matchType PathMatcherType
	pattern   string
	regex     *regexp.Regexp
}

// NewPathMatcher 创建路径匹配器
func NewPathMatcher(matchType PathMatcherType, pattern string) (*PathMatcher, error) {
	pm := &PathMatcher{
		matchType: matchType,
		pattern:   pattern,
	}

	// 如果是正则表达式或 Glob 模式，预编译
	switch matchType {
	case PathMatchRegex:
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		pm.regex = regex
	case PathMatchGlob:
		// 将 glob 模式转换为正则：* 匹配任意字符（包括 /），? 匹配单个字符
		var sb strings.Builder
		sb.WriteString("^")
		for _, r := range pattern {
			switch r {
			case '*':
				sb.WriteString(".*")
			case '?':
				sb.WriteString(".")
			default:
				sb.WriteString(regexp.QuoteMeta(string(r)))
			}
		}
		sb.WriteString("$")
		regex, err := regexp.Compile(sb.String())
		if err != nil {
			return nil, err
		}
		pm.regex = regex
	}

	return pm, nil
}

// Match 执行路径匹配
func (pm *PathMatcher) Match(path string) bool {
	switch pm.matchType {
	case PathMatchExact:
		return path == pm.pattern
	case PathMatchPrefix:
		return strings.HasPrefix(path, pm.pattern)
	case PathMatchSuffix:
		return strings.HasSuffix(path, pm.pattern)
	case PathMatchGlob:
		if pm.regex != nil {
			return pm.regex.MatchString(path)
		}
		return false
	case PathMatchRegex:
		if pm.regex != nil {
			return pm.regex.MatchString(path)
		}
		return false
	case PathMatchContains:
		return strings.Contains(path, pm.pattern)
	default:
		return false
	}
}

// globMatcherCache 按 glob 模式缓存已编译的匹配器
//
// MatchPathGlob 原先每次调用都走 NewPathMatcher，也就是每次都把 glob 转成正则、
// 再 regexp.Compile 一遍 而调用方的典型形态是在请求路径上按一张配置好的模式表
// 逐项匹配，于是同一批模式每个请求都要重新编译一轮
//
// 实测 go-rpc-gateway 的访问日志中间件：默认 6 条 SkipPaths 让单个请求产生
// 375 次分配，其中 80% 落在 regexp.compile 里 同仓的限流与白名单中间件
// 经 MatchPathWithMethod 走同一条路，同样每请求重编译
//
// 以模式串为键是安全的，不会无界增长：模式来自配置，是有界集合；
// 被匹配的 path 来自请求（可被外部控制），但它是匹配的输入而不是缓存的键
//
// PathMatcher 构造后不可变，唯一的方法 Match 只读字段，且 *regexp.Regexp 本身
// 并发安全，所以缓存项可以跨 goroutine 共享
//
// [EN] Compiled matchers cached by glob pattern
// MatchPathGlob used to call NewPathMatcher every time, converting the glob to a regex
// and compiling it again on each call, while the typical caller matches a request path
// against a configured table of patterns, recompiling that whole batch per request
// Measured in go-rpc-gateway's access-log middleware, the six default SkipPaths produced
// 375 allocations per request with 80% of them inside regexp.compile, and the rate-limit
// and whitelist middlewares reach the same code through MatchPathWithMethod
// Keying by pattern cannot grow without bound: patterns come from configuration and form
// a bounded set, while the matched path is externally controlled but is the input rather
// than the key
// A PathMatcher is immutable once built and its only method reads fields, and
// *regexp.Regexp is itself safe for concurrent use, so entries are shared across goroutines
var globMatcherCache sync.Map // map[string]*PathMatcher，编译失败时存 nil

// MatchPathGlob Glob 模式匹配路径（支持 * 和 ? 通配符，* 可匹配 /）
// 示例: MatchPathGlob("/api/v1/resource", "/api/*") => true
// 同一个模式只编译一次，之后命中缓存
//
// [EN] A given pattern is compiled once and served from cache afterwards
func MatchPathGlob(path, pattern string) bool {
	if pattern == path {
		return true
	}
	cached, ok := globMatcherCache.Load(pattern)
	if !ok {
		// 编译失败也入缓存（存 nil）：否则一条写错的配置会让每个请求都重试编译
		// [EN] A failed compile is cached as nil too, or one malformed config entry would
		// retry the compile on every request
		var value any
		if pm, err := NewPathMatcher(PathMatchGlob, pattern); err == nil {
			value = pm
		}
		cached, _ = globMatcherCache.LoadOrStore(pattern, value)
	}
	pm, _ := cached.(*PathMatcher)
	if pm == nil {
		return false
	}
	return pm.Match(path)
}

// MatchPathWithMethod 匹配路径和 HTTP 方法
// pathPattern 支持 glob 通配符，allowedMethods 为空时匹配所有方法
func MatchPathWithMethod(path, method, pathPattern string, allowedMethods []string) bool {
	// 检查路径是否匹配
	if !MatchPathGlob(path, pathPattern) {
		return false
	}

	// 如果没有指定方法限制，则允许所有方法
	if len(allowedMethods) == 0 {
		return true
	}

	// 检查方法是否在允许列表中
	return MatchMethod(allowedMethods, method)
}

// MatchMethod 检查 HTTP 方法是否匹配
// 如果 methods 为空，则匹配所有方法
// 方法名比较不区分大小写
func MatchMethod(methods []string, method string) bool {
	if len(methods) == 0 {
		return true // 没有指定方法，匹配所有
	}
	for _, m := range methods {
		if strings.EqualFold(m, method) {
			return true
		}
	}
	return false
}

// NormalizePath 标准化路径（移除重复的斜杠，确保正确的开头和结尾）
// 示例: NormalizePath("//api///users/") => "/api/users"
func NormalizePath(path string) string {
	if path == "" {
		return "/"
	}
	path = strings.ReplaceAll(path, "//", "/")
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	return path
}

// ExtractPathSegments 提取路径段
// 示例: ExtractPathSegments("/api/v1/users") => ["api", "v1", "users"]
func ExtractPathSegments(path string) []string {
	path = NormalizePath(path)
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

// PathMatcherBuilder 路径匹配器构建器
type PathMatcherBuilder struct {
	matchers []*PathMatcher
}

// NewPathMatcherBuilder 创建路径匹配器构建器
func NewPathMatcherBuilder() *PathMatcherBuilder {
	return &PathMatcherBuilder{
		matchers: make([]*PathMatcher, 0),
	}
}

// AddExact 添加精确匹配
func (b *PathMatcherBuilder) AddExact(pattern string) *PathMatcherBuilder {
	pm, _ := NewPathMatcher(PathMatchExact, pattern)
	if pm != nil {
		b.matchers = append(b.matchers, pm)
	}
	return b
}

// AddPrefix 添加前缀匹配
func (b *PathMatcherBuilder) AddPrefix(pattern string) *PathMatcherBuilder {
	pm, _ := NewPathMatcher(PathMatchPrefix, pattern)
	if pm != nil {
		b.matchers = append(b.matchers, pm)
	}
	return b
}

// AddSuffix 添加后缀匹配
func (b *PathMatcherBuilder) AddSuffix(pattern string) *PathMatcherBuilder {
	pm, _ := NewPathMatcher(PathMatchSuffix, pattern)
	if pm != nil {
		b.matchers = append(b.matchers, pm)
	}
	return b
}

// AddGlob 添加 Glob 匹配
func (b *PathMatcherBuilder) AddGlob(pattern string) *PathMatcherBuilder {
	pm, _ := NewPathMatcher(PathMatchGlob, pattern)
	if pm != nil {
		b.matchers = append(b.matchers, pm)
	}
	return b
}

// AddRegex 添加正则匹配
func (b *PathMatcherBuilder) AddRegex(pattern string) *PathMatcherBuilder {
	pm, _ := NewPathMatcher(PathMatchRegex, pattern)
	if pm != nil {
		b.matchers = append(b.matchers, pm)
	}
	return b
}

// AddContains 添加包含匹配
func (b *PathMatcherBuilder) AddContains(pattern string) *PathMatcherBuilder {
	pm, _ := NewPathMatcher(PathMatchContains, pattern)
	if pm != nil {
		b.matchers = append(b.matchers, pm)
	}
	return b
}

// MatchAny 匹配任意一个模式
func (b *PathMatcherBuilder) MatchAny(path string) bool {
	for _, matcher := range b.matchers {
		if matcher.Match(path) {
			return true
		}
	}
	return false
}

// MatchAll 匹配所有模式
func (b *PathMatcherBuilder) MatchAll(path string) bool {
	if len(b.matchers) == 0 {
		return false
	}

	for _, matcher := range b.matchers {
		if !matcher.Match(path) {
			return false
		}
	}
	return true
}

// Build 构建完成
func (b *PathMatcherBuilder) Build() []*PathMatcher {
	return b.matchers
}
