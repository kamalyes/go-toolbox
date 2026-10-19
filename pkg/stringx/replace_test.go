/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:30:55
 * @FilePath: \go-toolbox\pkg\stringx\replace_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	result := Replace("hello, world", "hello", "hi", 1)
	assert.Equal(t, "hi, world", result)
}

func TestReplaceAll(t *testing.T) {
	result := ReplaceAll("hello, hello, world", "hello", "hi")
	assert.Equal(t, "hi, hi, world", result)
}

func TestReplaceWithIndex(t *testing.T) {
	tests := []struct {
		input       string
		startIndex  int
		endIndex    int
		replacedStr string
		expected    string
	}{
		{"hello world", 6, 11, "*", "hello *****"},
		{"hello world", 0, 5, "*", "***** world"},
		{"hello world", 0, 0, "*", "hello world"}, // 替换长度为0
		{"hello world", 5, 5, "*", "hello world"}, // 替换长度为0
		{"hello world", 5, 6, "*", "hello*world"}, // 替换单个字符
		{"hello", 1, 4, "X", "hXXXo"},             // 替换多个字符
		{"", 0, 1, "*", ""},                       // 空字符串
		{"test", -1, 2, "#", "##st"},              // startIndex < 0
		{"test", 1, 10, "#", "t###"},              // endIndex 超出范围
		{"test", 3, 1, "#", "test"},               // startIndex > endIndex
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("input: %s, start: %d, end: %d", test.input, test.startIndex, test.endIndex), func(t *testing.T) {
			result := ReplaceWithIndex(test.input, test.startIndex, test.endIndex, test.replacedStr)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestPad(t *testing.T) {
	resultDefault := Pad("hello", 10)
	assert.Equal(t, "hell*****o", resultDefault)

	result := Pad("hello", 10, &Paddler{Position: Middle})
	assert.Equal(t, "hell*****o", result)

	result = Pad("world", 8, &Paddler{Position: Left})
	assert.Equal(t, "***world", result)

	result = Pad("ok", 5, &Paddler{Position: Right})
	assert.Equal(t, "ok***", result)
}

func TestReplaceWithMatcher(t *testing.T) {
	result := ReplaceWithMatcher("hello 123 world 456", `\d+`, func(s string) string {
		return "xxx"
	})
	assert.Equal(t, "hello xxx world xxx", result)
}

func TestHide(t *testing.T) {
	result := Hide("password12345", 8, 10)
	assert.Equal(t, "password**345", result)
}

func TestReplaceSpecialChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "HelloXXWorldX"},
		{"Go is fun.", "GoXisXfunX"},
		{"Special #chars#", "SpecialXXcharsX"},
		{"1234-5678", "1234X5678"},
		{"NoSpecialChars", "NoSpecialChars"},
		{"", ""},
		{"!@#$%^&*()", "XXXXXXXXXX"},
	}

	for _, test := range tests {
		output := ReplaceSpecialChars(test.input, 'X')
		if output != test.expected {
			t.Errorf("ReplaceSpecialChars(%q) = %q; expected %q", test.input, output, test.expected)
		}
	}
}

// TestReplaceWithMatcherCache 正则缓存必须返回可用的正则，且并发安全
//
// 原实现每次调用都 regexp.MustCompile，而编译比匹配贵得多：实测整次调用
// 1201ns / 2454B / 24 次分配，其中编译本身占 545ns / 2248B / 17 次
// 加缓存后降到 550ns / 177B / 7 次，缓存查找只要 8.1ns
//
// [EN] The regexp cache must return a working regexp and be safe for concurrent
// use. The original compiled on every call, and compiling costs far more than
// matching: the whole call measured 1201ns with 24 allocations, of which
// compilation alone was 545ns and 17. With the cache it is 550ns and 7, and a
// cache lookup costs 8.1ns
func TestReplaceWithMatcherCache(t *testing.T) {
	upper := func(s string) string { return strings.ToUpper(s) }

	// 同一模式反复调用，结果必须一致
	for i := 0; i < 50; i++ {
		assert.Equal(t, "ABC-123", ReplaceWithMatcher("abc-123", `[a-z]+`, upper), "第 %d 次", i)
	}
	// 不同模式互不干扰
	assert.Equal(t, "abc-XXX", ReplaceWithMatcher("abc-123", `\d+`, func(string) string { return "XXX" }))
	assert.Equal(t, "ABC-123", ReplaceWithMatcher("abc-123", `[a-z]+`, upper))

	// 缓存命中的是同一个编译结果
	assert.Same(t, cachedRegex(`\d+`), cachedRegex(`\d+`), "同一模式应取到同一个 *Regexp")
	assert.NotSame(t, cachedRegex(`\d+`), cachedRegex(`[a-z]+`), "不同模式不应取到同一个")

	// 非法模式仍然 panic，与 regexp.MustCompile 的既有行为一致
	assert.Panics(t, func() { ReplaceWithMatcher("x", `([`, upper) }, "非法模式应 panic")

	// 并发使用同一模式与不同模式
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				if ReplaceWithMatcher("abc-123", `[a-z]+`, upper) != "ABC-123" {
					t.Error("并发下结果不一致")
					return
				}
				_ = cachedRegex(`pat` + string(rune('a'+n%26)))
			}
		}(i)
	}
	wg.Wait()
}

var benchReplaced string

func BenchmarkReplaceWithMatcher(b *testing.B) {
	s := "phone 13800138000 and 13900139000, mail a@b.com"
	mask := func(m string) string { return m[:3] + "****" + m[7:] }

	b.Run("缓存命中", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchReplaced = ReplaceWithMatcher(s, `\d{11}`, mask)
		}
	})
	b.Run("对照/每次编译", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchReplaced = regexp.MustCompile(`\d{11}`).ReplaceAllStringFunc(s, mask)
		}
	})
	b.Run("仅缓存查找", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = cachedRegex(`\d{11}`)
		}
	})
}
