/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 20:35:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 21:10:10
 * @FilePath: \go-toolbox\pkg\stringx\width_test.go
 * @Description: 字符显示宽度计算测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuneWidth(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		expected int
	}{
		// ASCII 字符
		{"ASCII space", ' ', 1},
		{"ASCII letter", 'A', 1},
		{"ASCII digit", '1', 1},
		{"ASCII symbol", '@', 1},

		// 控制字符
		{"Control char null", '\x00', 0},
		{"Control char tab", '\t', 0},
		{"Control char newline", '\n', 0},
		{"Control char DEL", '\x7F', 0},

		// 中文字符
		{"Chinese common", '中', 2},
		{"Chinese common2", '国', 2},
		{"Chinese rare", '䶮', 2},

		// 日文字符
		{"Hiragana", 'あ', 2},
		{"Katakana", 'ア', 2},
		{"Japanese kanji", '日', 2},

		// 韩文字符
		{"Hangul", '한', 2},
		{"Hangul2", '글', 2},

		// 杂项技术符号
		{"Clock emoji ⏰", '⏰', 2},
		{"Timer ⏱", '⏱', 2},
		{"Alarm ⏲", '⏲', 2},

		// 杂项符号
		{"Star ★", '★', 2},
		{"Check ✓", '✓', 2},
		{"Cross ✗", '✗', 2},
		{"Heart ♥", '♥', 2},
		{"Snowman ☃", '☃', 2},

		// Emoji 表情
		{"Smile emoji", '😀', 2},
		{"Heart emoji", '❤', 2},
		{"Fire emoji", '🔥', 2},
		{"Rocket emoji", '🚀', 2},
		{"Star emoji ⭐", '⭐', 2},
		{"Party emoji 🎉", '🎉', 2},
		{"Money emoji 💰", '💰', 2},

		// 全角字符
		{"Fullwidth A", 'Ａ', 2},
		{"Fullwidth 1", '１', 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuneWidth(tt.r)
			assert.Equal(t, tt.expected, got, "RuneWidth(%q) should return %d", tt.r, tt.expected)
		})
	}
}

func TestDisplayWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// 纯 ASCII
		{"Empty string", "", 0},
		{"ASCII only", "Hello", 5},
		{"ASCII with space", "Hello World", 11},

		// 纯中文
		{"Chinese only", "你好世界", 8},
		{"Chinese sentence", "中文测试", 8},

		// 中英文混合
		{"Mixed CN EN", "你好 World", 10}, // 你(2)+好(2)+空格(1)+W(1)+o(1)+r(1)+l(1)+d(1) = 10
		{"Mixed CN EN 2", "Hello 世界", 10},

		// 带表情
		{"With emoji", "😀😁😂", 6},
		{"Mixed with emoji", "Hello 🌍", 8},
		{"CN with emoji", "你好 🎉", 7},

		// 复杂混合
		{"Complex mix", "⏰ 结束时间", 11},
		{"Complex mix 2", "✅ 状态", 7},
		{"Complex mix 3", "🎉 活动名称", 11},
		{"Complex mix 4", "👥 参与人数", 11},
		{"Complex mix 5", "🔥 热度", 7},

		// 特殊符号
		{"Special symbols", "★☆♥♦", 8},
		{"Math symbols", "≈ ≤ ≥", 5}, // ≈(1)+空格(1)+≤(1)+空格(1)+≥(1)

		// 全角字符
		{"Fullwidth", "ＡＢＣ１２３", 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DisplayWidth(tt.input)
			assert.Equal(t, tt.expected, got, "DisplayWidth(%q) should return %d", tt.input, tt.expected)
		})
	}
}

func TestDisplayWidthChain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"ASCII", "Hello", 5},
		{"Chinese", "你好", 4},
		{"Mixed", "Hello 世界", 10},
		{"With emoji", "🎉 Party", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.input).DisplayWidthChain()
			assert.Equal(t, tt.expected, got, "New(%q).DisplayWidthChain() should return %d", tt.input, tt.expected)
		})
	}
}

// BenchmarkRuneWidth 基准测试单个字符宽度计算
func BenchmarkRuneWidth(b *testing.B) {
	chars := []rune{'A', '中', '😀', '⏰', '★'}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range chars {
			RuneWidth(r)
		}
	}
}

// BenchmarkDisplayWidth 基准测试字符串宽度计算
func BenchmarkDisplayWidth(b *testing.B) {
	tests := []string{
		"Hello World",
		"你好世界",
		"Hello 世界 🌍",
		"⏰ 结束时间 🎉 活动名称",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range tests {
			DisplayWidth(s)
		}
	}
}
