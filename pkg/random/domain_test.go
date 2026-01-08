/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-05 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 15:16:15
 * @FilePath: \go-toolbox\pkg\random\domain_test.go
 * @Description: 域名随机生成测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package random

import (
	"strings"
	"testing"
)

// TestNewDomainKeywordBuilder 测试默认配置
func TestNewDomainKeywordBuilder(t *testing.T) {
	builder := NewDomainKeywordBuilder("game")
	keywords := builder.Generate()

	if len(keywords) != 3 {
		t.Errorf("Expected 3 keywords, got %d", len(keywords))
	}

	for _, kw := range keywords {
		if !strings.Contains(kw, "game") {
			t.Errorf("Keyword %s should contain 'game'", kw)
		}
	}
}

// TestDomainKeywordBuilder_WithCount 测试自定义数量
func TestDomainKeywordBuilder_WithCount(t *testing.T) {
	tests := []struct {
		name  string
		count int
		want  int
	}{
		{"Count 1", 1, 1},
		{"Count 5", 5, 5},
		{"Count 10", 10, 10},
		{"Invalid Count 0", 0, 3},
		{"Invalid Count -1", -1, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keywords := NewDomainKeywordBuilder("test").WithCount(tt.count).Generate()
			if len(keywords) != tt.want {
				t.Errorf("Expected %d keywords, got %d", tt.want, len(keywords))
			}
		})
	}
}

// TestDomainKeywordBuilder_WithPrefixLength 测试自定义前缀长度
func TestDomainKeywordBuilder_WithPrefixLength(t *testing.T) {
	keywords := NewDomainKeywordBuilder("test").
		WithPrefixLength(5, 5).
		WithCount(10).
		Generate()

	for _, kw := range keywords {
		prefix := strings.Split(kw, "test")[0]
		if len(prefix) != 5 {
			t.Errorf("Expected prefix length 5, got %d in %s", len(prefix), kw)
		}
	}
}

// TestDomainKeywordBuilder_WithSuffixLength 测试自定义后缀长度
func TestDomainKeywordBuilder_WithSuffixLength(t *testing.T) {
	keywords := NewDomainKeywordBuilder("test").
		WithSuffixLength(4, 4).
		WithCount(10).
		Generate()

	for _, kw := range keywords {
		parts := strings.Split(kw, "test")
		if len(parts) < 2 {
			t.Errorf("Keyword %s should have suffix", kw)
			continue
		}
		suffix := parts[1]
		if len(suffix) != 4 {
			t.Errorf("Expected suffix length 4, got %d in %s", len(suffix), kw)
		}
	}
}

// TestDomainKeywordBuilder_FullChain 测试完整链式调用
func TestDomainKeywordBuilder_FullChain(t *testing.T) {
	keywords := NewDomainKeywordBuilder("shop").
		WithCount(5).
		WithPrefixLength(2, 4).
		WithSuffixLength(3, 5).
		Generate()

	if len(keywords) != 5 {
		t.Errorf("Expected 5 keywords, got %d", len(keywords))
	}

	for _, kw := range keywords {
		if !strings.Contains(kw, "shop") {
			t.Errorf("Keyword %s should contain 'shop'", kw)
		}
		prefix := strings.Split(kw, "shop")[0]
		if len(prefix) < 2 || len(prefix) > 4 {
			t.Errorf("Prefix length should be 2-4, got %d in %s", len(prefix), kw)
		}
		parts := strings.Split(kw, "shop")
		if len(parts) >= 2 {
			suffix := parts[1]
			if len(suffix) < 3 || len(suffix) > 5 {
				t.Errorf("Suffix length should be 3-5, got %d in %s", len(suffix), kw)
			}
		}
	}
}

// TestDomainKeywordBuilder_Randomness 测试随机性
func TestDomainKeywordBuilder_Randomness(t *testing.T) {
	set1 := NewDomainKeywordBuilder("test").WithCount(10).Generate()
	set2 := NewDomainKeywordBuilder("test").WithCount(10).Generate()

	identical := 0
	for i := range set1 {
		if set1[i] == set2[i] {
			identical++
		}
	}

	if identical == len(set1) {
		t.Error("Two sets should not be identical")
	}
}

// TestDomainKeywordBuilder_AutoSwapMinMax 测试自动交换最小/最大值
func TestDomainKeywordBuilder_AutoSwapMinMax(t *testing.T) {
	keywords := NewDomainKeywordBuilder("test").
		WithPrefixLength(5, 2).
		WithSuffixLength(6, 3).
		WithCount(5).
		Generate()

	for _, kw := range keywords {
		prefix := strings.Split(kw, "test")[0]
		if len(prefix) < 2 || len(prefix) > 5 {
			t.Errorf("Prefix length should be 2-5 after swap, got %d in %s", len(prefix), kw)
		}

		parts := strings.Split(kw, "test")
		if len(parts) >= 2 {
			suffix := parts[1]
			if len(suffix) < 3 || len(suffix) > 6 {
				t.Errorf("Suffix length should be 3-6 after swap, got %d in %s", len(suffix), kw)
			}
		}
	}
}

// TestJoinDomainsWithTLDs 测试域名与TLD拼接
func TestJoinDomainsWithTLDs(t *testing.T) {
	domains := []string{"game", "shop", "test"}
	tlds := []string{"com", "net"}

	result := JoinDomainsWithTLDs(domains, tlds)

	if len(result) != 6 {
		t.Errorf("Expected 6 domains, got %d", len(result))
	}

	expected := []string{
		"game.com", "game.net",
		"shop.com", "shop.net",
		"test.com", "test.net",
	}

	for _, exp := range expected {
		found := false
		for _, res := range result {
			if res == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected domain %s not found in result", exp)
		}
	}
}

// TestJoinDomainsWithTLDs_EmptyInputs 测试空输入
func TestJoinDomainsWithTLDs_EmptyInputs(t *testing.T) {
	tests := []struct {
		name    string
		domains []string
		tlds    []string
		want    int
	}{
		{"Empty domains", []string{}, []string{"com", "net"}, 0},
		{"Empty TLDs", []string{"game", "shop"}, []string{}, 0},
		{"Both empty", []string{}, []string{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinDomainsWithTLDs(tt.domains, tt.tlds)
			if len(result) != tt.want {
				t.Errorf("Expected %d results, got %d", tt.want, len(result))
			}
		})
	}
}

// TestJoinDomainsWithTLDs_SingleTLD 测试单个TLD
func TestJoinDomainsWithTLDs_SingleTLD(t *testing.T) {
	domains := []string{"example", "demo"}
	tlds := []string{"io"}

	result := JoinDomainsWithTLDs(domains, tlds)

	expected := []string{"example.io", "demo.io"}
	if len(result) != len(expected) {
		t.Errorf("Expected %d domains, got %d", len(expected), len(result))
	}

	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("Expected %s, got %s", exp, result[i])
		}
	}
}

// BenchmarkDomainKeywordBuilder_Default 基准测试：默认配置
func BenchmarkDomainKeywordBuilder_Default(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDomainKeywordBuilder("test").Generate()
	}
}

// BenchmarkDomainKeywordBuilder_LargeCount 基准测试：大量生成
func BenchmarkDomainKeywordBuilder_LargeCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDomainKeywordBuilder("test").WithCount(100).Generate()
	}
}

// BenchmarkDomainKeywordBuilder_CustomLengths 基准测试：自定义长度
func BenchmarkDomainKeywordBuilder_CustomLengths(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDomainKeywordBuilder("test").
			WithCount(10).
			WithPrefixLength(2, 5).
			WithSuffixLength(3, 6).
			Generate()
	}
}

// BenchmarkJoinDomainsWithTLDs 基准测试：TLD拼接
func BenchmarkJoinDomainsWithTLDs(b *testing.B) {
	domains := []string{"game", "shop", "test", "demo", "app"}
	tlds := []string{"com", "net", "org", "io", "co"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JoinDomainsWithTLDs(domains, tlds)
	}
}
