/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-09 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-09 15:00:00
 * @FilePath: \go-toolbox\pkg\metadata\language_test.go
 * @Description: 语言信息提取测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package metadata

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"小写转换", "zh-cn", "zh-CN"},
		{"大写转换", "EN-US", "en-US"},
		{"下划线替换", "zh_CN", "zh-CN"},
		{"单一语言代码", "en", "en"},
		{"单一大写", "EN", "en"},
		{"复杂格式", "zh_cn", "zh-CN"},
		{"带空格", " en-US ", "en-US"},
		{"空字符串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeLanguage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLanguageExtractor_FromQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com?lang=zh-CN", nil)
	extractor := NewLanguageExtractor("en")

	lang := extractor.Extract(req)
	assert.Equal(t, "zh-CN", lang)
}

func TestLanguageExtractor_FromQueryAlternative(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com?language=ja", nil)
	extractor := NewLanguageExtractor("en")

	lang := extractor.Extract(req)
	assert.Equal(t, "ja", lang)
}

func TestLanguageExtractor_FromHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Language", "fr-FR")

	extractor := NewLanguageExtractor("en")
	lang := extractor.Extract(req)

	assert.Equal(t, "fr-FR", lang)
}

func TestLanguageExtractor_FromCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.AddCookie(&http.Cookie{Name: "language", Value: "de-DE"})

	extractor := NewLanguageExtractor("en")
	lang := extractor.Extract(req)

	assert.Equal(t, "de-DE", lang)
}

func TestLanguageExtractor_FromAcceptLanguage(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	extractor := NewLanguageExtractor("en")
	lang := extractor.Extract(req)

	assert.Equal(t, "zh-CN", lang)
}

func TestLanguageExtractor_DefaultValue(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	extractor := NewLanguageExtractor("zh-CN")
	lang := extractor.Extract(req)

	assert.Equal(t, "zh-CN", lang)
}

func TestLanguageExtractor_Priority(t *testing.T) {
	// Query > Header > Cookie > Accept-Language > Default
	req := httptest.NewRequest("GET", "http://example.com?lang=query-lang", nil)
	req.Header.Set("X-Language", "header-lang")
	req.Header.Set("Accept-Language", "accept-lang")
	req.AddCookie(&http.Cookie{Name: "language", Value: "cookie-lang"})

	extractor := NewLanguageExtractor("default-lang")
	lang := extractor.Extract(req)

	// Query 优先级最高
	assert.Equal(t, "query-lang", lang)
}

func TestLanguageExtractor_PriorityWithoutQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Language", "header-lang")
	req.Header.Set("Accept-Language", "accept-lang")
	req.AddCookie(&http.Cookie{Name: "language", Value: "cookie-lang"})

	extractor := NewLanguageExtractor("default-lang")
	lang := extractor.Extract(req)

	// Header 优先级第二
	assert.Equal(t, "header-lang", lang)
}

func TestLanguageExtractor_PriorityWithoutQueryAndHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Accept-Language", "accept-lang")
	req.AddCookie(&http.Cookie{Name: "language", Value: "cookie-lang"})

	extractor := NewLanguageExtractor("default-lang")
	lang := extractor.Extract(req)

	// Cookie 优先级第三
	assert.Equal(t, "cookie-lang", lang)
}

func TestLanguageExtractor_NormalizeDisabled(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com?lang=zh-cn", nil)

	extractor := NewLanguageExtractor("en")
	extractor.Normalize = false
	lang := extractor.Extract(req)

	// 不标准化，保持原样
	assert.Equal(t, "zh-cn", lang)
}

func TestLanguageExtractor_NormalizeEnabled(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com?lang=zh-cn", nil)

	extractor := NewLanguageExtractor("en")
	extractor.Normalize = true
	lang := extractor.Extract(req)

	// 标准化为 zh-CN
	assert.Equal(t, "zh-CN", lang)
}

func TestLanguageExtractor_CustomKeys(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com?locale=ja-JP", nil)

	extractor := NewLanguageExtractor("en")
	extractor.QueryKeys = []string{"locale"}
	lang := extractor.Extract(req)

	assert.Equal(t, "ja-JP", lang)
}

func TestLanguageExtractor_CustomHeaderKeys(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Locale", "ko-KR")

	extractor := NewLanguageExtractor("en")
	extractor.HeaderKeys = []string{"X-Locale"}
	lang := extractor.Extract(req)

	assert.Equal(t, "ko-KR", lang)
}

func TestLanguageExtractor_DisableAcceptLanguage(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Accept-Language", "zh-CN")

	extractor := NewLanguageExtractor("en")
	extractor.UseAcceptLang = false
	lang := extractor.Extract(req)

	// 禁用 Accept-Language，应返回默认值
	assert.Equal(t, "en", lang)
}

func TestExtractLanguage(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com?lang=fr", nil)
	lang := ExtractLanguage(req)

	assert.Equal(t, "fr", lang)
}

func TestExtractLanguageWithDefault(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	lang := ExtractLanguageWithDefault(req, "ja")

	assert.Equal(t, "ja", lang)
}

func TestLanguageExtractor_ExtractWithContext(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com?lang=zh-CN", nil)
	extractor := NewLanguageExtractor("en")

	const langKey ContextKey = "language"
	newReq, lang := extractor.ExtractWithContext(req, langKey)

	assert.Equal(t, "zh-CN", lang)
	assert.Equal(t, "zh-CN", newReq.Context().Value(langKey))
}

func TestLanguageExtractor_ComplexScenario(t *testing.T) {
	// 模拟一个复杂的真实场景
	req := httptest.NewRequest("GET", "http://example.com/api/resource", nil)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	extractor := NewLanguageExtractor("en")
	lang := extractor.Extract(req)

	assert.Equal(t, "zh-CN", lang)
}

func TestLanguageExtractor_MultipleQueryKeys(t *testing.T) {
	// 测试多个 Query 键名，第一个有值就使用
	req := httptest.NewRequest("GET", "http://example.com?language=es", nil)

	extractor := NewLanguageExtractor("en")
	extractor.QueryKeys = []string{"lang", "language", "locale"}
	lang := extractor.Extract(req)

	assert.Equal(t, "es", lang)
}

func TestLanguageExtractor_MultipleHeaderKeys(t *testing.T) {
	// 测试多个 Header 键名
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Custom-Language", "pt-BR")

	extractor := NewLanguageExtractor("en")
	extractor.QueryKeys = []string{} // 禁用 Query
	extractor.HeaderKeys = []string{"X-Language", "X-Custom-Language"}
	lang := extractor.Extract(req)

	assert.Equal(t, "pt-BR", lang)
}

func TestFromAcceptLanguageSource(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9")

	lang := FromAcceptLanguageSource(req.Context(), req, "Accept-Language")
	assert.Equal(t, "fr-FR", lang)
}

func TestFromAcceptLanguageSource_NilRequest(t *testing.T) {
	lang := FromAcceptLanguageSource(nil, nil, "Accept-Language")
	assert.Equal(t, "", lang)
}

func TestFromAcceptLanguageSource_EmptyHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	lang := FromAcceptLanguageSource(req.Context(), req, "Accept-Language")
	assert.Equal(t, "", lang)
}
