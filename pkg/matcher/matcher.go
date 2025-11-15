/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-15 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-15 02:00:00
 * @FilePath: \go-toolbox\pkg\matcher\matcher.go
 * @Description: 生产级通用规则匹配引擎 - 高并发、类型安全、高性能
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package matcher

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Rule 规则接口
type Rule[T any] interface {
	// Match 判断是否匹配
	Match(ctx *Context) bool
	// Priority 优先级（数字越大优先级越高）
	Priority() int
	// Result 返回匹配结果
	Result() T
	// ID 规则唯一标识
	ID() string
	// Enabled 是否启用
	Enabled() bool
}

// Context 匹配上下文（并发安全 - 使用 sync.Map 优化性能）
type Context struct {
	data     sync.Map // 高并发优化
	parent   context.Context
	deadline atomic.Int64 // UnixNano
	metadata sync.Map
}

// NewContext 创建上下文
func NewContext() *Context {
	return &Context{
		parent: context.Background(),
	}
}

// NewContextWithParent 创建带父上下文的上下文
func NewContextWithParent(parent context.Context) *Context {
	return &Context{
		parent: parent,
	}
}

// WithTimeout 设置超时
func (c *Context) WithTimeout(timeout time.Duration) *Context {
	c.deadline.Store(time.Now().Add(timeout).UnixNano())
	return c
}

// IsExpired 检查是否超时
func (c *Context) IsExpired() bool {
	dl := c.deadline.Load()
	if dl == 0 {
		return false
	}
	return time.Now().UnixNano() > dl
}

// Done 返回父上下文的 Done channel
func (c *Context) Done() <-chan struct{} {
	return c.parent.Done()
}

// Set 设置上下文数据（并发安全）
func (c *Context) Set(key string, value interface{}) *Context {
	c.data.Store(key, value)
	return c
}

// SetBatch 批量设置（性能优化）
func (c *Context) SetBatch(kvs map[string]interface{}) *Context {
	for k, v := range kvs {
		c.data.Store(k, v)
	}
	return c
}

// Get 获取上下文数据（并发安全）
func (c *Context) Get(key string) (interface{}, bool) {
	return c.data.Load(key)
}

// MustGet 获取数据，不存在则 panic
func (c *Context) MustGet(key string) interface{} {
	if val, ok := c.Get(key); ok {
		return val
	}
	panic(fmt.Sprintf("key %s not found in context", key))
}

// GetString 获取字符串
func (c *Context) GetString(key string) string {
	if val, ok := c.data.Load(key); ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetStringSlice 获取字符串切片
func (c *Context) GetStringSlice(key string) []string {
	if val, ok := c.data.Load(key); ok {
		if slice, ok := val.([]string); ok {
			return slice
		}
	}
	return nil
}

// GetInt 获取整数
func (c *Context) GetInt(key string) int {
	if val, ok := c.data.Load(key); ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case int32:
			return int(v)
		}
	}
	return 0
}

// GetInt64 获取 int64
func (c *Context) GetInt64(key string) int64 {
	if val, ok := c.data.Load(key); ok {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case int32:
			return int64(v)
		}
	}
	return 0
}

// GetBool 获取布尔值
func (c *Context) GetBool(key string) bool {
	if val, ok := c.data.Load(key); ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// GetFloat64 获取浮点数
func (c *Context) GetFloat64(key string) float64 {
	if val, ok := c.data.Load(key); ok {
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return 0
}

// SetMetadata 设置元数据
func (c *Context) SetMetadata(key, value string) *Context {
	c.metadata.Store(key, value)
	return c
}

// GetMetadata 获取元数据
func (c *Context) GetMetadata(key string) string {
	if val, ok := c.metadata.Load(key); ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// Clone 克隆上下文
func (c *Context) Clone() *Context {
	newCtx := &Context{
		parent: c.parent,
	}
	newCtx.deadline.Store(c.deadline.Load())

	// 复制数据
	c.data.Range(func(k, v interface{}) bool {
		newCtx.data.Store(k, v)
		return true
	})

	// 复制元数据
	c.metadata.Range(func(k, v interface{}) bool {
		newCtx.metadata.Store(k, v)
		return true
	})

	return newCtx
}

// Matcher 规则匹配器（并发安全 - 使用 atomic.Pointer 优化性能）
type Matcher[T any] struct {
	mu          sync.RWMutex              // 仅保护写操作
	rules       atomic.Pointer[[]Rule[T]] // 原子指针避免复制
	sorted      atomic.Bool
	cache       *matchCache[T]
	stats       *MatcherStats
	middlewares atomic.Pointer[[]MatchMiddleware[T]]
}

// MatcherStats 匹配器统计信息
type MatcherStats struct {
	totalMatches   atomic.Int64
	successMatches atomic.Int64
	failedMatches  atomic.Int64
	cacheHits      atomic.Int64
	cacheMisses    atomic.Int64
}

// matchCache 匹配缓存
type matchCache[T any] struct {
	enabled bool
	cache   sync.Map // key: string, value: cacheEntry[T]
	ttl     time.Duration
}

type cacheEntry[T any] struct {
	result    T
	matched   bool
	expiresAt time.Time
}

// MatchMiddleware 匹配中间件
type MatchMiddleware[T any] func(ctx *Context, next func() (T, bool)) (T, bool)

// NewMatcher 创建匹配器
func NewMatcher[T any]() *Matcher[T] {
	m := &Matcher[T]{
		stats: &MatcherStats{},
		cache: &matchCache[T]{
			enabled: false,
			ttl:     5 * time.Minute,
		},
	}
	emptyRules := make([]Rule[T], 0, 16) // 预分配容量
	m.rules.Store(&emptyRules)
	emptyMws := make([]MatchMiddleware[T], 0, 4)
	m.middlewares.Store(&emptyMws)
	return m
}

// EnableCache 启用缓存
func (m *Matcher[T]) EnableCache(ttl time.Duration) *Matcher[T] {
	m.cache.enabled = true
	m.cache.ttl = ttl
	return m
}

// DisableCache 禁用缓存
func (m *Matcher[T]) DisableCache() *Matcher[T] {
	m.cache.enabled = false
	m.cache.cache = sync.Map{}
	return m
}

// Use 添加中间件
func (m *Matcher[T]) Use(middleware MatchMiddleware[T]) *Matcher[T] {
	m.mu.Lock()
	oldMws := m.middlewares.Load()
	newMws := make([]MatchMiddleware[T], len(*oldMws), len(*oldMws)+1)
	copy(newMws, *oldMws)
	newMws = append(newMws, middleware)
	m.middlewares.Store(&newMws)
	m.mu.Unlock()
	return m
}

// AddRule 添加规则
func (m *Matcher[T]) AddRule(rule Rule[T]) *Matcher[T] {
	m.mu.Lock()
	oldRules := m.rules.Load()
	newRules := make([]Rule[T], len(*oldRules), len(*oldRules)+1)
	copy(newRules, *oldRules)
	newRules = append(newRules, rule)
	m.rules.Store(&newRules)
	m.sorted.Store(false)
	m.mu.Unlock()
	return m
}

// AddRules 批量添加规则
func (m *Matcher[T]) AddRules(rules ...Rule[T]) *Matcher[T] {
	m.mu.Lock()
	oldRules := m.rules.Load()
	newRules := make([]Rule[T], len(*oldRules), len(*oldRules)+len(rules))
	copy(newRules, *oldRules)
	newRules = append(newRules, rules...)
	m.rules.Store(&newRules)
	m.sorted.Store(false)
	m.mu.Unlock()
	return m
}

// RemoveRule 移除规则
func (m *Matcher[T]) RemoveRule(id string) *Matcher[T] {
	m.mu.Lock()
	oldRules := m.rules.Load()
	newRules := make([]Rule[T], 0, len(*oldRules))
	for _, rule := range *oldRules {
		if rule.ID() != id {
			newRules = append(newRules, rule)
		}
	}
	m.rules.Store(&newRules)
	m.mu.Unlock()
	return m
}

// ClearRules 清空所有规则
func (m *Matcher[T]) ClearRules() *Matcher[T] {
	m.mu.Lock()
	emptyRules := make([]Rule[T], 0)
	m.rules.Store(&emptyRules)
	m.sorted.Store(false)
	m.mu.Unlock()
	return m
}

// getRules 获取当前规则列表（返回指针，零拷贝）
func (m *Matcher[T]) getRules() *[]Rule[T] {
	return m.rules.Load()
}

// Match 执行匹配（返回第一个匹配的规则）
func (m *Matcher[T]) Match(ctx *Context) (T, bool) {
	m.stats.totalMatches.Add(1)

	// 检查超时
	if ctx.IsExpired() {
		m.stats.failedMatches.Add(1)
		var zero T
		return zero, false
	}

	// 检查缓存
	if m.cache.enabled {
		if cached, ok := m.getCache(ctx); ok {
			m.stats.cacheHits.Add(1)
			if cached.matched {
				m.stats.successMatches.Add(1)
			} else {
				m.stats.failedMatches.Add(1)
			}
			return cached.result, cached.matched
		}
		m.stats.cacheMisses.Add(1)
	}

	// 执行中间件链
	result, matched := m.executeWithMiddlewares(ctx, func() (T, bool) {
		return m.doMatch(ctx)
	})

	// 更新缓存
	if m.cache.enabled {
		m.setCache(ctx, result, matched)
	}

	if matched {
		m.stats.successMatches.Add(1)
	} else {
		m.stats.failedMatches.Add(1)
	}

	return result, matched
}

// doMatch 执行实际匹配（零拷贝优化）
func (m *Matcher[T]) doMatch(ctx *Context) (T, bool) {
	// 确保规则已排序
	m.ensureSorted()

	rules := *m.getRules() // 解引用一次
	for i := range rules {
		rule := &rules[i]
		// 检查是否启用
		if !(*rule).Enabled() {
			continue
		}

		// 检查是否匹配
		if (*rule).Match(ctx) {
			return (*rule).Result(), true
		}
	}

	var zero T
	return zero, false
}

// MatchAll 执行匹配（返回所有匹配的规则）
func (m *Matcher[T]) MatchAll(ctx *Context) []T {
	m.stats.totalMatches.Add(1)

	m.ensureSorted()

	rules := *m.getRules() // 解引用
	results := make([]T, 0)

	for i := range rules {
		rule := &rules[i]
		if !(*rule).Enabled() {
			continue
		}

		if (*rule).Match(ctx) {
			results = append(results, (*rule).Result())
		}
	}

	if len(results) > 0 {
		m.stats.successMatches.Add(1)
	} else {
		m.stats.failedMatches.Add(1)
	}

	return results
}

// executeWithMiddlewares 执行中间件链（零拷贝优化）
func (m *Matcher[T]) executeWithMiddlewares(ctx *Context, final func() (T, bool)) (T, bool) {
	middlewares := *m.getMiddlewares() // 解引用

	if len(middlewares) == 0 {
		return final()
	}

	// 构建中间件链
	var chain func() (T, bool)
	chain = final

	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		next := chain
		chain = func() (T, bool) {
			return middleware(ctx, next)
		}
	}

	return chain()
}

// getMiddlewares 获取中间件列表（返回指针）
func (m *Matcher[T]) getMiddlewares() *[]MatchMiddleware[T] {
	return m.middlewares.Load()
}

// ensureSorted 确保规则已排序
func (m *Matcher[T]) ensureSorted() {
	if m.sorted.Load() {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if m.sorted.Load() {
		return
	}

	// 获取当前规则并复制（排序需要修改）
	oldRules := m.rules.Load()
	newRules := make([]Rule[T], len(*oldRules))
	copy(newRules, *oldRules)

	// 排序规则
	m.sortRules(newRules)

	// 更新指针
	m.rules.Store(&newRules)
	m.sorted.Store(true)
}

// sortRules 排序规则（快速排序）
func (m *Matcher[T]) sortRules(rules []Rule[T]) {
	if len(rules) <= 1 {
		return
	}

	// 快速排序（降序）
	quickSort(rules, 0, len(rules)-1)
}

func quickSort[T any](rules []Rule[T], low, high int) {
	if low < high {
		pi := partition(rules, low, high)
		quickSort(rules, low, pi-1)
		quickSort(rules, pi+1, high)
	}
}

func partition[T any](rules []Rule[T], low, high int) int {
	pivot := rules[high].Priority()
	i := low - 1

	for j := low; j < high; j++ {
		if rules[j].Priority() > pivot {
			i++
			rules[i], rules[j] = rules[j], rules[i]
		}
	}

	rules[i+1], rules[high] = rules[high], rules[i+1]
	return i + 1
}

// getCache 获取缓存
func (m *Matcher[T]) getCache(ctx *Context) (*cacheEntry[T], bool) {
	key := m.getCacheKey(ctx)
	if val, ok := m.cache.cache.Load(key); ok {
		entry := val.(*cacheEntry[T])
		if time.Now().Before(entry.expiresAt) {
			return entry, true
		}
		m.cache.cache.Delete(key)
	}
	return nil, false
}

// setCache 设置缓存
func (m *Matcher[T]) setCache(ctx *Context, result T, matched bool) {
	key := m.getCacheKey(ctx)
	entry := &cacheEntry[T]{
		result:    result,
		matched:   matched,
		expiresAt: time.Now().Add(m.cache.ttl),
	}
	m.cache.cache.Store(key, entry)
}

// getCacheKey 生成缓存键（优化：使用 strings.Builder + sync.Map.Range）
func (m *Matcher[T]) getCacheKey(ctx *Context) string {
	var sb strings.Builder
	sb.Grow(128) // 预分配

	ctx.data.Range(func(k, v interface{}) bool {
		if key, ok := k.(string); ok {
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(fmt.Sprintf("%v", v))
			sb.WriteString(";")
		}
		return true
	})

	return sb.String()
}

// Stats 获取统计信息
func (m *Matcher[T]) Stats() map[string]int64 {
	return map[string]int64{
		"total_matches":   m.stats.totalMatches.Load(),
		"success_matches": m.stats.successMatches.Load(),
		"failed_matches":  m.stats.failedMatches.Load(),
		"cache_hits":      m.stats.cacheHits.Load(),
		"cache_misses":    m.stats.cacheMisses.Load(),
	}
}

// ResetStats 重置统计
func (m *Matcher[T]) ResetStats() {
	m.stats.totalMatches.Store(0)
	m.stats.successMatches.Store(0)
	m.stats.failedMatches.Store(0)
	m.stats.cacheHits.Store(0)
	m.stats.cacheMisses.Store(0)
}

// ChainRule 链式规则构建器
type ChainRule[T any] struct {
	conditions []func(*Context) bool
	priority   int
	result     T
	id         string
	enabled    bool
}

// NewChainRule 创建链式规则
func NewChainRule[T any](result T) *ChainRule[T] {
	return &ChainRule[T]{
		conditions: make([]func(*Context) bool, 0),
		priority:   0,
		result:     result,
		id:         fmt.Sprintf("rule_%d", time.Now().UnixNano()),
		enabled:    true,
	}
}

// When 添加条件
func (r *ChainRule[T]) When(condition func(*Context) bool) *ChainRule[T] {
	r.conditions = append(r.conditions, condition)
	return r
}

// WithPriority 设置优先级
func (r *ChainRule[T]) WithPriority(priority int) *ChainRule[T] {
	r.priority = priority
	return r
}

// WithID 设置ID
func (r *ChainRule[T]) WithID(id string) *ChainRule[T] {
	r.id = id
	return r
}

// WithEnabled 设置是否启用
func (r *ChainRule[T]) WithEnabled(enabled bool) *ChainRule[T] {
	r.enabled = enabled
	return r
}

// Match 实现 Rule 接口
func (r *ChainRule[T]) Match(ctx *Context) bool {
	for _, condition := range r.conditions {
		if !condition(ctx) {
			return false
		}
	}
	return true
}

// Priority 实现 Rule 接口
func (r *ChainRule[T]) Priority() int {
	return r.priority
}

// Result 实现 Rule 接口
func (r *ChainRule[T]) Result() T {
	return r.result
}

// ID 实现 Rule 接口
func (r *ChainRule[T]) ID() string {
	return r.id
}

// Enabled 实现 Rule 接口
func (r *ChainRule[T]) Enabled() bool {
	return r.enabled
}

// ===== 常用条件构建器 =====

// MatchString 字符串精确匹配
func MatchString(key, expected string) func(*Context) bool {
	return func(ctx *Context) bool {
		return ctx.GetString(key) == expected
	}
}

// MatchStringIn 字符串在列表中
func MatchStringIn(key string, list []string) func(*Context) bool {
	return func(ctx *Context) bool {
		val := ctx.GetString(key)
		for _, item := range list {
			if item == val {
				return true
			}
		}
		return false
	}
}

// MatchStringNotIn 字符串不在列表中
func MatchStringNotIn(key string, list []string) func(*Context) bool {
	return func(ctx *Context) bool {
		val := ctx.GetString(key)
		for _, item := range list {
			if item == val {
				return false
			}
		}
		return true
	}
}

// MatchPattern 路径模式匹配
func MatchPattern(key, pattern string) func(*Context) bool {
	return func(ctx *Context) bool {
		val := ctx.GetString(key)
		matched, _ := filepath.Match(pattern, val)
		return matched || pattern == val
	}
}

// MatchPrefix 前缀匹配
func MatchPrefix(key, prefix string) func(*Context) bool {
	return func(ctx *Context) bool {
		return strings.HasPrefix(ctx.GetString(key), prefix)
	}
}

// MatchSuffix 后缀匹配
func MatchSuffix(key, suffix string) func(*Context) bool {
	return func(ctx *Context) bool {
		return strings.HasSuffix(ctx.GetString(key), suffix)
	}
}

// MatchContains 包含匹配
func MatchContains(key, substring string) func(*Context) bool {
	return func(ctx *Context) bool {
		return strings.Contains(ctx.GetString(key), substring)
	}
}

// MatchBool 布尔值匹配
func MatchBool(key string, expected bool) func(*Context) bool {
	return func(ctx *Context) bool {
		return ctx.GetBool(key) == expected
	}
}

// MatchAny 任意条件满足
func MatchAny(conditions ...func(*Context) bool) func(*Context) bool {
	return func(ctx *Context) bool {
		for _, cond := range conditions {
			if cond(ctx) {
				return true
			}
		}
		return false
	}
}

// MatchAll 所有条件满足
func MatchAll(conditions ...func(*Context) bool) func(*Context) bool {
	return func(ctx *Context) bool {
		for _, cond := range conditions {
			if !cond(ctx) {
				return false
			}
		}
		return true
	}
}

// MatchNot 取反
func MatchNot(condition func(*Context) bool) func(*Context) bool {
	return func(ctx *Context) bool {
		return !condition(ctx)
	}
}

// MatchMethodIn HTTP方法匹配
func MatchMethodIn(methods []string) func(*Context) bool {
	if len(methods) == 0 {
		return func(*Context) bool { return true }
	}
	return func(ctx *Context) bool {
		method := ctx.GetString("method")
		for _, m := range methods {
			if strings.EqualFold(m, method) {
				return true
			}
		}
		return false
	}
}

// MatchWildcard 通配符匹配
func MatchWildcard(key, pattern string) func(*Context) bool {
	return func(ctx *Context) bool {
		val := ctx.GetString(key)
		if pattern == "*" {
			return true
		}
		matched, _ := filepath.Match(pattern, val)
		return matched
	}
}
