/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-01 23:18:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-01 23:52:00
 * @FilePath: \go-toolbox\pkg\matcher\path_cache_test.go
 * @Description: glob 匹配器编译缓存的正确性与并发测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package matcher

import (
	"fmt"
	"sync"
	"testing"
)

// skipPathsLikeConfig 与 go-config 的 logging.SkipPaths 默认值一致
var skipPathsLikeConfig = []string{"/health", "/metrics", "/favicon.ico", "/ping", "/readiness", "/liveness"}

// BenchmarkMatchPathGlob_SkipPathTable 模拟访问日志中间件的真实形态：
// 对一条不在表里的请求路径，把整张模式表逐项匹配一遍
func BenchmarkMatchPathGlob_SkipPathTable(b *testing.B) {
	const path = "/bench.BenchService/Echo"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range skipPathsLikeConfig {
			_ = MatchPathGlob(path, p)
		}
	}
}

// BenchmarkMatchPathGlob_Wildcard 带通配符的模式，必须走正则
func BenchmarkMatchPathGlob_Wildcard(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MatchPathGlob("/v1/audit-logs/42", "/v1/audit-logs*")
	}
}

// TestMatchPathGlobCacheCorrectness 加缓存后匹配结果必须与不走缓存逐次新建完全一致
// 缓存最容易出的错是模式串串味 —— 两个模式共用了同一个编译结果
//
// [EN] Caching must not change any result; the classic cache bug here is two patterns
// sharing one compiled matcher
func TestMatchPathGlobCacheCorrectness(t *testing.T) {
	cases := []struct {
		path, pattern string
	}{
		{"/health", "/health"},
		{"/health", "/metrics"},
		{"/v1/audit-logs", "/v1/audit-logs*"},
		{"/v1/audit-logs/42", "/v1/audit-logs*"},
		{"/v1/audit-logs/42", "/v1/audit-log*"},
		{"/v1/other", "/v1/audit-logs*"},
		{"/a/b/c", "/a/*"},
		{"/a/b/c", "/a/?"},
		{"/a/b", "/a/?"},
		{"/x.json", "/x.?son"},
		{"/x.json", "*.json"},
		{"/anything", "*"},
		{"", ""},
		{"/plus+sign", "/plus+sign"},
		{"/dot.char", "/dot.char"},
		{"/dot.char", "/dotXchar"},
		{"/brack[et]", "/brack[et]"},
		{"/paren(x)", "/paren(x)"},
		{"/a|b", "/a|b"},
		{"/caret^", "/caret^"},
		{"/dollar$", "/dollar$"},
		{"/back\\slash", "/back\\slash"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s_vs_%s", tc.path, tc.pattern), func(t *testing.T) {
			// 参照：绕开缓存，直接构造匹配器
			want := tc.pattern == tc.path
			if !want {
				if pm, err := NewPathMatcher(PathMatchGlob, tc.pattern); err == nil {
					want = pm.Match(tc.path)
				}
			}
			// 连续两次，第一次填缓存第二次命中
			got1 := MatchPathGlob(tc.path, tc.pattern)
			got2 := MatchPathGlob(tc.path, tc.pattern)
			if got1 != want || got2 != want {
				t.Errorf("path=%q pattern=%q: 首次=%v 二次=%v 期望=%v", tc.path, tc.pattern, got1, got2, want)
			}
		})
	}
}

// TestMatchPathGlobCacheConcurrent 缓存项跨 goroutine 共享，-race 下不得有竞争
// [EN] Cache entries are shared across goroutines and must be race-free under -race
func TestMatchPathGlobCacheConcurrent(t *testing.T) {
	patterns := []string{"/a/*", "/b/?", "/c", "*.json", "/d/e/*"}
	paths := []string{"/a/x", "/b/y", "/c", "/f.json", "/d/e/g", "/nomatch"}

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_ = MatchPathGlob(paths[(i+j)%len(paths)], patterns[(i+j)%len(patterns)])
			}
		}(i)
	}
	wg.Wait()
}

// TestMatchPathGlobInvalidPatternCached 编译不出来的模式也要入缓存，不能每次重试
// [EN] A pattern that fails to compile is cached too, rather than retried every call
func TestMatchPathGlobInvalidPatternCached(t *testing.T) {
	// glob 转正则会转义除 * 和 ? 以外的全部字符，所以几乎编不出错；
	// 这里只验证「不 panic、结果稳定」，以及重复调用行为一致
	const weird = "/[unclosed"
	first := MatchPathGlob("/[unclosed", weird)
	second := MatchPathGlob("/[unclosed", weird)
	if first != second {
		t.Errorf("同一模式两次调用结果不一致: %v vs %v", first, second)
	}
}
