/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-16 22:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-16 22:55:55
 * @FilePath: \go-toolbox\safe\mathx_test.go
 * @Description: 安全数学工具函数的测试文件
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package safe

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFastHash(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"", 1}, // 空字符串应返回1
		{"hello", FastHash("hello")},
		{"world", FastHash("world")},
		{"a", FastHash("a")},
		{"1", FastHash("1")},
		{"123456789", FastHash("123456789")},
		{"这是中文", FastHash("这是中文")},
		{"🚀🌟💎", FastHash("🚀🌟💎")},
		{"!@#$%^&*()", FastHash("!@#$%^&*()")},
		{"   ", FastHash("   ")},       // 空格
		{"\n\t\r", FastHash("\n\t\r")}, // 控制字符
		{"a" + string(rune(0)) + "b", FastHash("a" + string(rune(0)) + "b")},                                     // 包含null字符
		{string([]byte{255, 254, 253}), FastHash(string([]byte{255, 254, 253}))},                                 // 高位字节
		{"The quick brown fox jumps over the lazy dog", FastHash("The quick brown fox jumps over the lazy dog")}, // 长字符串
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", FastHash("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
		{"abcdefghijklmnopqrstuvwxyz", FastHash("abcdefghijklmnopqrstuvwxyz")},
		{"0123456789", FastHash("0123456789")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := FastHash(tt.input)
			if tt.input == "" {
				assert.Equal(t, uint64(1), result)
			} else {
				// 检查哈希值的一致性
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestShortHash 测试短哈希生成（默认7位）
func TestShortHash(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(t *testing.T, result string)
	}{
		{
			name:  "空字符串",
			input: "",
			checkFn: func(t *testing.T, result string) {
				// 默认7位 Base36 字符
				assert.Len(t, result, 7)
				assert.Regexp(t, "^[0-9a-z]+$", result)
			},
		},
		{
			name:  "简单字符串",
			input: "hello",
			checkFn: func(t *testing.T, result string) {
				assert.Len(t, result, 7)
				assert.Regexp(t, "^[0-9a-z]+$", result)
			},
		},
		{
			name:  "IP:Port格式",
			input: "192.168.1.100:8080",
			checkFn: func(t *testing.T, result string) {
				assert.Len(t, result, 7)
				assert.Regexp(t, "^[0-9a-z]+$", result)
			},
		},
		{
			name:  "中文字符",
			input: "你好世界",
			checkFn: func(t *testing.T, result string) {
				assert.Len(t, result, 7)
				assert.Regexp(t, "^[0-9a-z]+$", result)
			},
		},
		{
			name:  "长字符串",
			input: "The quick brown fox jumps over the lazy dog",
			checkFn: func(t *testing.T, result string) {
				assert.Len(t, result, 7)
				assert.Regexp(t, "^[0-9a-z]+$", result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShortHash(tt.input)
			tt.checkFn(t, result)

			// 测试一致性：相同输入应该产生相同输出
			result2 := ShortHash(tt.input)
			assert.Equal(t, result, result2, "相同输入应该产生相同的短哈希")
		})
	}
}

// TestShortHashWithLength 测试可配置长度的短哈希
func TestShortHashWithLength(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		want   int // 期望长度
	}{
		{"6位哈希", "test-input", 6, 6},
		{"7位哈希", "test-input", 7, 7},
		{"8位哈希", "test-input", 8, 8},
		{"10位哈希", "test-input", 10, 10},
		{"边界-最小", "test", 1, 1},
		{"边界-最大", "test", 13, 13},
		{"小于最小", "test", 0, 1}, // 自动限制为1
		{"负数", "test", -5, 1},  // 自动限制为1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShortHashWithLength(tt.input, tt.length)
			assert.Len(t, result, tt.want)
			assert.Regexp(t, "^[0-9a-z]+$", result)

			// 测试一致性
			result2 := ShortHashWithLength(tt.input, tt.length)
			assert.Equal(t, result, result2)
		})
	}
}

// TestShortHashUniqueness 测试短哈希的唯一性（7位）
func TestShortHashUniqueness(t *testing.T) {
	inputs := []string{
		"192.168.1.100:8080",
		"192.168.1.100:8081",
		"192.168.1.101:8080",
		"10.0.0.1:8080",
		"localhost:8080",
		"example.com:443",
		"node-1",
		"node-2",
		"pod-abc-123",
		"pod-abc-124",
	}

	hashes := make(map[string]string)
	for _, input := range inputs {
		hash := ShortHash(input)
		if existing, found := hashes[hash]; found {
			t.Errorf("哈希冲突: %s 和 %s 产生了相同的哈希 %s", input, existing, hash)
		}
		hashes[hash] = input
		t.Logf("输入: %-25s -> 哈希: %s", input, hash)
	}
}

// TestShortHashWithLengthUniqueness 测试不同长度的唯一性
func TestShortHashWithLengthUniqueness(t *testing.T) {
	inputs := []string{
		"192.168.1.100:8080",
		"192.168.1.100:8081",
		"192.168.1.101:8080",
		"10.0.0.1:8080",
		"localhost:8080",
	}

	lengths := []int{6, 7, 8}
	for _, length := range lengths {
		t.Run(fmt.Sprintf("长度%d", length), func(t *testing.T) {
			hashes := make(map[string]string)
			for _, input := range inputs {
				hash := ShortHashWithLength(input, length)
				if existing, found := hashes[hash]; found {
					t.Errorf("哈希冲突: %s 和 %s 产生了相同的哈希 %s", input, existing, hash)
				}
				hashes[hash] = input
			}
		})
	}
}

// TestShortHashConsistency 测试短哈希的一致性
func TestShortHashConsistency(t *testing.T) {
	input := "192.168.1.100:8080"

	// 多次调用应该返回相同结果
	results := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		results[i] = ShortHash(input)
	}

	first := results[0]
	for i, result := range results {
		assert.Equal(t, first, result, "第 %d 次调用返回了不同的结果", i)
	}
}

// TestShortHashCollisionResistance 测试短哈希的抗冲突能力（7位）
func TestShortHashCollisionResistance(t *testing.T) {
	const testCount = 10000
	hashes := make(map[string]string, testCount)
	collisions := 0

	for i := 0; i < testCount; i++ {
		input := fmt.Sprintf("test-input-%d", i)
		hash := ShortHash(input)

		if existing, found := hashes[hash]; found {
			collisions++
			t.Logf("发现冲突 #%d: %s 和 %s 产生了相同的哈希 %s", collisions, input, existing, hash)
		}
		hashes[hash] = input
	}

	// 7位 Base36 约 36^7 = 78,364,164,096 种可能
	// 10000 个输入，期望冲突 < 1 个
	assert.LessOrEqual(t, collisions, 1, "在 %d 个输入中冲突应该 <= 1", testCount)
	t.Logf("测试完成: %d 个输入，%d 个唯一哈希，%d 个冲突", testCount, len(hashes), collisions)
}

// TestShortHashWithLengthCollisionResistance 测试不同长度的抗冲突能力
func TestShortHashWithLengthCollisionResistance(t *testing.T) {
	const testCount = 10000

	tests := []struct {
		length           int
		maxCollisions    int
		expectedCapacity int64
	}{
		{6, 5, 2176782336},    // 36^6 ≈ 21亿
		{7, 1, 78364164096},   // 36^7 ≈ 783亿
		{8, 0, 2821109907456}, // 36^8 ≈ 2.8万亿
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("长度%d", tt.length), func(t *testing.T) {
			hashes := make(map[string]string, testCount)
			collisions := 0

			for i := 0; i < testCount; i++ {
				input := fmt.Sprintf("test-input-%d", i)
				hash := ShortHashWithLength(input, tt.length)

				if existing, found := hashes[hash]; found {
					collisions++
					if collisions <= 3 { // 只记录前3个冲突
						t.Logf("冲突 #%d: %s 和 %s -> %s", collisions, input, existing, hash)
					}
				}
				hashes[hash] = input
			}

			assert.LessOrEqual(t, collisions, tt.maxCollisions,
				"长度 %d 在 %d 个输入中冲突应该 <= %d", tt.length, testCount, tt.maxCollisions)
			t.Logf("长度 %d: %d 输入，%d 唯一哈希，%d 冲突（容量: %d）",
				tt.length, testCount, len(hashes), collisions, tt.expectedCapacity)
		})
	}
}

// BenchmarkShortHash 短哈希性能测试（默认7位）
func BenchmarkShortHash(b *testing.B) {
	input := "192.168.1.100:8080"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ShortHash(input)
	}
}

// BenchmarkShortHashWithLength 不同长度的性能测试
func BenchmarkShortHashWithLength(b *testing.B) {
	input := "192.168.1.100:8080"
	lengths := []int{6, 7, 8, 10, 13}

	for _, length := range lengths {
		b.Run(fmt.Sprintf("长度%d", length), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ShortHashWithLength(input, length)
			}
		})
	}
}

// BenchmarkShortHashWithVariedInputs 不同输入的性能测试
func BenchmarkShortHashWithVariedInputs(b *testing.B) {
	inputs := []string{
		"short",
		"192.168.1.100:8080",
		"pod-name-abc-123-xyz",
		"这是一个中文字符串测试",
		"The quick brown fox jumps over the lazy dog",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		_ = ShortHash(input)
	}
}

func TestNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{-1, 2},
		{0, 2},
		{1, 2},
		{2, 4},
		{3, 4},
		{4, 8},
		{5, 8},
		{7, 8},
		{8, 16},
		{9, 16},
		{15, 16},
		{16, 32},
		{17, 32},
		{31, 32},
		{32, 64},
		{33, 64},
		{63, 64},
		{64, 128},
		{127, 128},
		{128, 256},
		{255, 256},
		{256, 512},
		{511, 512},
		{512, 1024},
		{1023, 1024},
		{1024, 2048},
		{2047, 2048},
		{4095, 4096},
		{65535, 65536},
		{131071, 131072},
		{262143, 262144},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := NextPowerOfTwo(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeAdd(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int64
		expected  int64
		expectErr bool
	}{
		{"正常加法", 10, 20, 30, false},
		{"零加法1", 0, 10, 10, false},
		{"零加法2", 10, 0, 10, false},
		{"零加法3", 0, 0, 0, false},
		{"负数加法1", -10, 5, -5, false},
		{"负数加法2", -10, -5, -15, false},
		{"负数加法3", 10, -5, 5, false},
		{"大数加法", 1000000, 2000000, 3000000, false},
		{"最大正数", math.MaxInt64 - 1, 1, math.MaxInt64, false},
		{"最小负数", math.MinInt64 + 1, -1, math.MinInt64, false},
		{"溢出测试1", math.MaxInt64, 1, 0, true},
		{"溢出测试2", math.MaxInt64, 100, 0, true},
		{"溢出测试3", math.MaxInt64 - 10, 20, 0, true},
		{"下溢测试1", math.MinInt64, -1, 0, true},
		{"下溢测试2", math.MinInt64, -100, 0, true},
		{"下溢测试3", math.MinInt64 + 10, -20, 0, true},
		{"边界测试1", math.MaxInt64 / 2, math.MaxInt64 / 2, math.MaxInt64 - 1, false},
		{"边界测试2", math.MinInt64 / 2, math.MinInt64 / 2, math.MinInt64, false},
		{"1", 1, 1, 2, false},
		{"-1", -1, -1, -2, false},
		{"混合符号1", 100, -50, 50, false},
		{"混合符号2", -100, 50, -50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeAdd(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafeSubtract(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int64
		expected  int64
		expectErr bool
	}{
		{"正常减法", 20, 10, 10, false},
		{"零减法", 10, 0, 10, false},
		{"负数减法", -5, -10, 5, false},
		{"下溢测试", math.MinInt64, 1, 0, true},
		{"溢出测试", math.MaxInt64, -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeSubtract(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafeMultiply(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int64
		expected  int64
		expectErr bool
	}{
		{"正常乘法", 10, 20, 200, false},
		{"零乘法1", 0, 10, 0, false},
		{"零乘法2", 10, 0, 0, false},
		{"零乘法3", 0, 0, 0, false},
		{"零乘法4", 0, -10, 0, false},
		{"一乘法1", 1, 10, 10, false},
		{"一乘法2", 10, 1, 10, false},
		{"一乘法3", 1, 1, 1, false},
		{"负一乘法1", -1, 10, -10, false},
		{"负一乘法2", 10, -1, -10, false},
		{"负一乘法3", -1, -1, 1, false},
		{"负数乘法1", -10, 5, -50, false},
		{"负数乘法2", -10, -5, 50, false},
		{"负数乘法3", 10, -5, -50, false},
		{"大数乘法", 1000, 2000, 2000000, false},
		{"平方根测试", 46340, 46340, 2147395600, false}, // sqrt(MaxInt64) 约等于 3037000499
		{"边界测试1", 3037000499, 3, 9111001497, false},
		{"溢出测试1", math.MaxInt64, 2, 0, true},
		{"溢出测试2", math.MaxInt64, -2, 0, true},
		{"溢出测试3", 2, math.MaxInt64, 0, true},
		{"溢出测试4", -2, math.MaxInt64, 0, true},
		{"溢出测试5", math.MaxInt64/2 + 1, 2, 0, true},
		{"下溢测试1", math.MinInt64, 2, 0, true},
		{"下溢测试2", math.MinInt64/2 - 1, 2, 0, true},
		{"下溢测试3", -1000000, 10000000, -10000000000000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeMultiply(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int64
		expected  int64
		expectErr bool
	}{
		{"正常除法", 20, 10, 2, false},
		{"零被除", 0, 10, 0, false},
		{"除零错误", 10, 0, 0, true},
		{"溢出测试", math.MinInt64, -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeDivide(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafeModulo(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int64
		expected  int64
		expectErr bool
	}{
		{"正常取模", 23, 10, 3, false},
		{"零被模", 0, 10, 0, false},
		{"模零错误", 10, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeModulo(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafePower(t *testing.T) {
	tests := []struct {
		name      string
		base, exp int64
		expected  int64
		expectErr bool
	}{
		{"正常幂运算", 2, 3, 8, false},
		{"零的幂", 0, 5, 0, false},
		{"任何数的0次幂", 5, 0, 1, false},
		{"1的任何次幂", 1, 100, 1, false},
		{"-1的偶数次幂", -1, 4, 1, false},
		{"-1的奇数次幂", -1, 3, -1, false},
		{"负指数", 2, -1, 0, true},
		{"大数幂运算溢出", 100, 20, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafePower(tt.base, tt.exp)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafeSqrt(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		expected  float64
		expectErr bool
		tolerance float64
	}{
		{"正数平方根", 16.0, 4.0, false, 1e-10},
		{"零的平方根", 0.0, 0.0, false, 1e-10},
		{"小数平方根", 2.25, 1.5, false, 1e-10},
		{"负数平方根", -1.0, 0.0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeSqrt(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, tt.tolerance)
			}
		})
	}
}

func TestSafeLog(t *testing.T) {
	tests := []struct {
		name      string
		n, base   float64
		expected  float64
		expectErr bool
		tolerance float64
	}{
		{"正常对数", 8.0, 2.0, 3.0, false, 1e-10},
		{"自然对数", math.E, math.E, 1.0, false, 1e-10},
		{"非正数", -1.0, 2.0, 0.0, true, 0},
		{"零", 0.0, 2.0, 0.0, true, 0},
		{"无效底数", 2.0, 0.0, 0.0, true, 0},
		{"底数为1", 2.0, 1.0, 0.0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeLog(tt.n, tt.base)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, tt.tolerance)
			}
		})
	}
}

func TestSafeGCD(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int64
		expected int64
	}{
		{"正数GCD1", 48, 18, 6},
		{"正数GCD2", 12, 8, 4},
		{"正数GCD3", 15, 25, 5},
		{"正数GCD4", 21, 14, 7},
		{"正数GCD5", 35, 49, 7},
		{"正数GCD6", 56, 42, 14},
		{"正数GCD7", 72, 108, 36},
		{"正数GCD8", 100, 75, 25},
		{"正数GCD9", 144, 96, 48},
		{"正数GCD10", 200, 150, 50},
		{"一个为零1", 10, 0, 10},
		{"一个为零2", 0, 15, 15},
		{"两个都为零", 0, 0, 0},
		{"负数GCD1", -48, 18, 6},
		{"负数GCD2", 48, -18, 6},
		{"负数GCD3", -48, -18, 6},
		{"负数GCD4", -12, 8, 4},
		{"负数GCD5", 12, -8, 4},
		{"负数GCD6", -12, -8, 4},
		{"互质数1", 17, 13, 1},
		{"互质数2", 23, 29, 1},
		{"互质数3", 31, 37, 1},
		{"互质数4", 41, 43, 1},
		{"相同数1", 10, 10, 10},
		{"相同数2", 25, 25, 25},
		{"相同数3", 100, 100, 100},
		{"倍数关系1", 10, 5, 5},
		{"倍数关系2", 20, 4, 4},
		{"倍数关系3", 30, 6, 6},
		{"大数GCD1", 1071, 462, 21},
		{"大数GCD2", 2016, 1512, 504},
		{"大数GCD3", 12345, 6789, 3},
		{"大数GCD4", 98765, 54321, 1},
		{"连续数字1", 6, 8, 2},
		{"连续数字2", 9, 12, 3},
		{"连续数字3", 15, 18, 3},
		{"连续数字4", 20, 24, 4},
		{"幂数关系1", 16, 8, 8},    // 2^4, 2^3
		{"幂数关系2", 32, 16, 16},  // 2^5, 2^4
		{"幂数关系3", 64, 32, 32},  // 2^6, 2^5
		{"幂数关系4", 27, 9, 9},    // 3^3, 3^2
		{"幂数关系5", 81, 27, 27},  // 3^4, 3^3
		{"斐波那契数列1", 21, 13, 1}, // F(8), F(7)
		{"斐波那契数列2", 34, 21, 1}, // F(9), F(8)
		{"斐波那契数列3", 55, 34, 1}, // F(10), F(9)
		{"素数组合1", 11, 7, 1},
		{"素数组合2", 13, 11, 1},
		{"素数组合3", 17, 13, 1},
		{"素数组合4", 19, 17, 1},
		{"素数组合5", 23, 19, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeGCD(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeLCM(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int64
		expected  int64
		expectErr bool
	}{
		{"正数LCM", 12, 8, 24, false},
		{"一个为零", 10, 0, 0, false},
		{"互质数", 7, 11, 77, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeLCM(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		input    int64
		expected bool
	}{
		{-5, false},
		{-1, false},
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{10, false},
		{11, true},
		{12, false},
		{13, true},
		{14, false},
		{15, false},
		{16, false},
		{17, true},
		{18, false},
		{19, true},
		{20, false},
		{21, false},
		{22, false},
		{23, true},
		{24, false},
		{25, false},
		{29, true},
		{31, true},
		{37, true},
		{41, true},
		{43, true},
		{47, true},
		{49, false}, // 7*7
		{51, false}, // 3*17
		{53, true},
		{59, true},
		{61, true},
		{67, true},
		{71, true},
		{73, true},
		{79, true},
		{83, true},
		{89, true},
		{91, false}, // 7*13
		{97, true},
		{100, false},
		{101, true},
		{103, true},
		{107, true},
		{109, true},
		{113, true},
		{121, false}, // 11*11
		{127, true},
		{131, true},
		{137, true},
		{139, true},
		{149, true},
		{151, true},
		{157, true},
		{163, true},
		{167, true},
		{169, false}, // 13*13
		{173, true},
		{179, true},
		{181, true},
		{191, true},
		{193, true},
		{197, true},
		{199, true},
		{211, true},
		{223, true},
		{227, true},
		{229, true},
		{233, true},
		{239, true},
		{241, true},
		{251, true},
		{257, true},
		{263, true},
		{269, true},
		{271, true},
		{277, true},
		{281, true},
		{283, true},
		{289, false}, // 17*17
		{293, true},
		{307, true},
		{311, true},
		{313, true},
		{317, true},
		{331, true},
		{337, true},
		{347, true},
		{349, true},
		{353, true},
		{359, true},
		{367, true},
		{373, true},
		{379, true},
		{383, true},
		{389, true},
		{397, true},
		{401, true},
		{409, true},
		{419, true},
		{421, true},
		{431, true},
		{433, true},
		{439, true},
		{443, true},
		{449, true},
		{457, true},
		{461, true},
		{463, true},
		{467, true},
		{479, true},
		{487, true},
		{491, true},
		{499, true},
		{503, true},
		{509, true},
		{521, true},
		{523, true},
		{541, true},
		{547, true},
		{557, true},
		{563, true},
		{569, true},
		{571, true},
		{577, true},
		{587, true},
		{593, true},
		{599, true},
		{601, true},
		{607, true},
		{613, true},
		{617, true},
		{619, true},
		{631, true},
		{641, true},
		{643, true},
		{647, true},
		{653, true},
		{659, true},
		{661, true},
		{673, true},
		{677, true},
		{683, true},
		{691, true},
		{701, true},
		{709, true},
		{719, true},
		{727, true},
		{733, true},
		{739, true},
		{743, true},
		{751, true},
		{757, true},
		{761, true},
		{769, true},
		{773, true},
		{787, true},
		{797, true},
		{809, true},
		{811, true},
		{821, true},
		{823, true},
		{827, true},
		{829, true},
		{839, true},
		{853, true},
		{857, true},
		{859, true},
		{863, true},
		{877, true},
		{881, true},
		{883, true},
		{887, true},
		{907, true},
		{911, true},
		{919, true},
		{929, true},
		{937, true},
		{941, true},
		{947, true},
		{953, true},
		{967, true},
		{971, true},
		{977, true},
		{983, true},
		{991, true},
		{997, true},
		{1000, false},
		{1001, false}, // 7*11*13
		{1009, true},
		{1013, true},
		{1019, true},
		{1021, true},
		{1031, true},
		{1033, true},
		{1039, true},
		{1049, true},
		{1051, true},
		{1061, true},
		{1063, true},
		{1069, true},
		{1087, true},
		{1091, true},
		{1093, true},
		{1097, true},
		{1103, true},
		{1109, true},
		{1117, true},
		{1123, true},
		{1129, true},
		{1151, true},
		{1153, true},
		{1163, true},
		{1171, true},
		{1181, true},
		{1187, true},
		{1193, true},
		{1201, true},
		{1213, true},
		{1217, true},
		{1223, true},
		{1229, true},
		{1231, true},
		{1237, true},
		{1249, true},
		{1259, true},
		{1277, true},
		{1279, true},
		{1283, true},
		{1289, true},
		{1291, true},
		{1297, true},
		{1301, true},
		{1303, true},
		{1307, true},
		{1319, true},
		{1321, true},
		{1327, true},
		{10007, true},      // 大素数
		{10009, true},      // 大素数
		{100003, true},     // 更大素数
		{982451653, true},  // 超大素数
		{982451654, false}, // 非素数
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := IsPrime(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		expected  int64
		expectErr bool
	}{
		{"F(-10)", -10, 0, true},
		{"F(-1)", -1, 0, true},
		{"F(0)", 0, 0, false},
		{"F(1)", 1, 1, false},
		{"F(2)", 2, 1, false},
		{"F(3)", 3, 2, false},
		{"F(4)", 4, 3, false},
		{"F(5)", 5, 5, false},
		{"F(6)", 6, 8, false},
		{"F(7)", 7, 13, false},
		{"F(8)", 8, 21, false},
		{"F(9)", 9, 34, false},
		{"F(10)", 10, 55, false},
		{"F(11)", 11, 89, false},
		{"F(12)", 12, 144, false},
		{"F(13)", 13, 233, false},
		{"F(14)", 14, 377, false},
		{"F(15)", 15, 610, false},
		{"F(16)", 16, 987, false},
		{"F(17)", 17, 1597, false},
		{"F(18)", 18, 2584, false},
		{"F(19)", 19, 4181, false},
		{"F(20)", 20, 6765, false},
		{"F(25)", 25, 75025, false},
		{"F(30)", 30, 832040, false},
		{"F(35)", 35, 9227465, false},
		{"F(40)", 40, 102334155, false},
		{"F(45)", 45, 1134903170, false},
		{"F(50)", 50, 12586269025, false},
		{"F(55)", 55, 139583862445, false},
		{"F(60)", 60, 1548008755920, false},
		{"F(70)", 70, 190392490709135, false},
		{"F(80)", 80, 23416728348467685, false},
		{"F(90)", 90, 2880067194370816120, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Fibonacci(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		expected  string // 使用字符串比较big.Int
		expectErr bool
	}{
		{"0!", 0, "1", false},
		{"1!", 1, "1", false},
		{"5!", 5, "120", false},
		{"10!", 10, "3628800", false},
		{"负数", -1, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Factorial(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestSafeAverage(t *testing.T) {
	tests := []struct {
		name      string
		input     []int64
		expected  float64
		expectErr bool
		tolerance float64
	}{
		{"空切片", []int64{}, 0.0, true, 0},
		{"单个元素", []int64{42}, 42.0, false, 1e-10},
		{"两个元素", []int64{10, 20}, 15.0, false, 1e-10},
		{"三个元素", []int64{1, 2, 3}, 2.0, false, 1e-10},
		{"四个元素", []int64{10, 20, 30, 40}, 25.0, false, 1e-10},
		{"五个元素", []int64{1, 2, 3, 4, 5}, 3.0, false, 1e-10},
		{"六个元素", []int64{2, 4, 6, 8, 10, 12}, 7.0, false, 1e-10},
		{"七个元素", []int64{1, 3, 5, 7, 9, 11, 13}, 7.0, false, 1e-10},
		{"八个元素", []int64{10, 15, 20, 25, 30, 35, 40, 45}, 27.5, false, 1e-10},
		{"九个元素", []int64{100, 200, 300, 400, 500, 600, 700, 800, 900}, 500.0, false, 1e-10},
		{"十个元素", []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5.5, false, 1e-10},
		{"负数平坁01", []int64{-1, -2, -3}, -2.0, false, 1e-10},
		{"负数平坁02", []int64{-10, -20, -30}, -20.0, false, 1e-10},
		{"负数平坁03", []int64{-5, -10, -15, -20}, -12.5, false, 1e-10},
		{"混合符号平坁01", []int64{-10, 0, 10}, 0.0, false, 1e-10},
		{"混合符号平坁02", []int64{-5, -3, 0, 3, 5}, 0.0, false, 1e-10},
		{"混合符号平坁03", []int64{-100, 50, 25, 75}, 12.5, false, 1e-10},
		{"混合符号平坁04", []int64{-20, -10, 0, 10, 20}, 0.0, false, 1e-10},
		{"全零", []int64{0, 0, 0, 0}, 0.0, false, 1e-10},
		{"包含零", []int64{0, 1, 2, 3}, 1.5, false, 1e-10},
		{"相同数字", []int64{5, 5, 5, 5}, 5.0, false, 1e-10},
		{"大数平均", []int64{1000000, 2000000, 3000000}, 2000000.0, false, 1e-10},
		{"小数平均", []int64{1, 1, 1, 1, 1}, 1.0, false, 1e-10},
		{"不规则数列", []int64{13, 7, 23, 2, 18}, 12.6, false, 1e-10},
		{"斐波那契数列", []int64{1, 1, 2, 3, 5, 8, 13}, 33.0 / 7.0, false, 1e-10},
		{"平方数", []int64{1, 4, 9, 16, 25}, 11.0, false, 1e-10},
		{"立方数", []int64{1, 8, 27, 64, 125}, 45.0, false, 1e-10},
		{"素数", []int64{2, 3, 5, 7, 11, 13, 17, 19, 23}, float64(2+3+5+7+11+13+17+19+23) / 9.0, false, 1e-10},
		{"交替正负", []int64{-1, 1, -2, 2, -3, 3}, 0.0, false, 1e-10},
		{"递增数列", []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, 55.0, false, 1e-10},
		{"递减数列", []int64{100, 90, 80, 70, 60, 50, 40, 30, 20, 10}, 55.0, false, 1e-10},
		{"最大值边界", []int64{math.MaxInt64 - 10, 0, 10}, float64(math.MaxInt64-10) / 3.0, false, 1e-10},
		{"最小值边界", []int64{math.MinInt64 + 10, 0, -10}, float64(math.MinInt64+10) / 3.0, false, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeAverage(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, tt.tolerance)
			}
		})
	}
}

func TestSafeMax(t *testing.T) {
	tests := []struct {
		name      string
		input     []int64
		expected  int64
		expectErr bool
	}{
		{"正常最大值", []int64{1, 5, 3, 9, 2}, 9, false},
		{"空切片", []int64{}, 0, true},
		{"单个元素", []int64{42}, 42, false},
		{"负数", []int64{-5, -1, -10}, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeMax(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafeMin(t *testing.T) {
	tests := []struct {
		name      string
		input     []int64
		expected  int64
		expectErr bool
	}{
		{"正常最小值", []int64{1, 5, 3, 9, 2}, 1, false},
		{"空切片", []int64{}, 0, true},
		{"单个元素", []int64{42}, 42, false},
		{"负数", []int64{-5, -1, -10}, -10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeMin(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafeClamp(t *testing.T) {
	tests := []struct {
		name            string
		value, min, max int64
		expected        int64
		expectErr       bool
	}{
		{"范围内", 5, 1, 10, 5, false},
		{"小于最小值", -5, 1, 10, 1, false},
		{"大于最大值", 15, 1, 10, 10, false},
		{"无效范围", 5, 10, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeClamp(tt.value, tt.min, tt.max)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSafeAbs(t *testing.T) {
	tests := []struct {
		name      string
		input     int64
		expected  int64
		expectErr bool
	}{
		{"正数", 42, 42, false},
		{"负数", -42, 42, false},
		{"零", 0, 0, false},
		{"最小整数", math.MinInt64, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeAbs(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// 性能测试
func BenchmarkFastHash(b *testing.B) {
	testString := "这是一个测试字符串用于性能测试"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FastHash(testString)
	}
}

func BenchmarkNextPowerOfTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextPowerOfTwo(1000)
	}
}

func BenchmarkSafeAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SafeAdd(12345, 67890)
	}
}

func BenchmarkSafeMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SafeMultiply(12345, 67890)
	}
}

func BenchmarkSafeGCD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SafeGCD(123456, 789012)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(982451653)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(30)
	}
}

// SafeAccess新API的测试用例

func TestSafeAccessAtMethods(t *testing.T) {
	// 创建测试数据结构
	type DatabaseConfig struct {
		Host     string        `json:"host"`
		Port     int           `json:"port"`
		Enabled  bool          `json:"enabled"`
		Timeout  time.Duration `json:"timeout"`
		Username *string       `json:"username"`
	}

	type ServerConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	type Config struct {
		Database *DatabaseConfig `json:"database"`
		Server   *ServerConfig   `json:"server"`
	}

	username := "admin"
	config := &Config{
		Database: &DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Enabled:  true,
			Timeout:  time.Second * 30,
			Username: &username,
		},
		Server: &ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}

	safe := Safe(config)

	t.Run("StringAt正常访问", func(t *testing.T) {
		result := safe.StringAt("Database.Host", "unknown")
		assert.Equal(t, "localhost", result)
	})

	t.Run("StringAt不存在字段使用默认值", func(t *testing.T) {
		result := safe.StringAt("Database.NonExistent", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("IntAt正常访问", func(t *testing.T) {
		result := safe.IntAt("Database.Port", 3306)
		assert.Equal(t, 5432, result)
	})

	t.Run("IntAt嵌套路径访问", func(t *testing.T) {
		result := safe.IntAt("Server.Port", 8000)
		assert.Equal(t, 8080, result)
	})

	t.Run("BoolAt正常访问", func(t *testing.T) {
		result := safe.BoolAt("Database.Enabled", false)
		assert.Equal(t, true, result)
	})

	t.Run("BoolAt不存在字段使用默认值", func(t *testing.T) {
		result := safe.BoolAt("Database.NonExistent", false)
		assert.Equal(t, false, result)
	})

	t.Run("DurationAt正常访问", func(t *testing.T) {
		result := safe.DurationAt("Database.Timeout", time.Second*10)
		assert.Equal(t, time.Second*30, result)
	})

	t.Run("ValueAt获取原始值", func(t *testing.T) {
		result := safe.ValueAt("Database.Port")
		assert.Equal(t, 5432, result)
	})

	t.Run("StringOrAt空值处理", func(t *testing.T) {
		// 测试不存在的字段
		result := safe.StringOrAt("Database.EmptyField", "guest")
		assert.Equal(t, "guest", result)
	})

	t.Run("At方法链式调用", func(t *testing.T) {
		result := safe.At("Database.Host").String("unknown")
		assert.Equal(t, "localhost", result)
	})

	t.Run("At方法验证有效性", func(t *testing.T) {
		validAccess := safe.At("Database.Host")
		assert.True(t, validAccess.IsValid())

		invalidAccess := safe.At("NonExistent.Field")
		assert.False(t, invalidAccess.IsValid())
	})
}

func TestSafeAccessAtMethodsWithMap(t *testing.T) {
	// 测试Map数据
	mapData := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "MyApp",
			"version": "1.0.0",
			"port":    9000,
			"debug":   true,
			"config": map[string]interface{}{
				"timeout": "30s",
				"retries": 3,
			},
		},
		"database": map[string]interface{}{
			"host": "db.example.com",
			"port": 5432,
		},
	}

	safe := Safe(mapData)

	t.Run("Map数据StringAt访问", func(t *testing.T) {
		result := safe.StringAt("app.name", "Unknown")
		assert.Equal(t, "MyApp", result)
	})

	t.Run("Map数据IntAt访问", func(t *testing.T) {
		result := safe.IntAt("app.port", 8080)
		assert.Equal(t, 9000, result)
	})

	t.Run("Map数据BoolAt访问", func(t *testing.T) {
		result := safe.BoolAt("app.debug", false)
		assert.Equal(t, true, result)
	})

	t.Run("Map数据深层嵌套访问", func(t *testing.T) {
		result := safe.StringAt("app.config.timeout", "10s")
		assert.Equal(t, "30s", result)

		intResult := safe.IntAt("app.config.retries", 1)
		assert.Equal(t, 3, intResult)
	})

	t.Run("Map数据不同路径访问", func(t *testing.T) {
		dbHost := safe.StringAt("database.host", "localhost")
		assert.Equal(t, "db.example.com", dbHost)

		dbPort := safe.IntAt("database.port", 3306)
		assert.Equal(t, 5432, dbPort)
	})
}

func TestSafeAccessAtMethodsEdgeCases(t *testing.T) {
	t.Run("空路径处理", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		safe := Safe(data)

		result := safe.StringAt("", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("nil数据处理", func(t *testing.T) {
		safe := Safe(nil)

		result := safe.StringAt("any.path", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("单层路径访问", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		safe := Safe(data)

		result := safe.StringAt("key", "default")
		assert.Equal(t, "value", result)
	})

	t.Run("路径中包含点号的字段名", func(t *testing.T) {
		data := map[string]interface{}{
			"normal": map[string]interface{}{
				"field": "value",
			},
		}
		safe := Safe(data)

		result := safe.StringAt("normal.field", "default")
		assert.Equal(t, "value", result)
	})

	t.Run("多层级深度访问", func(t *testing.T) {
		data := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": map[string]interface{}{
						"value": "deep-value",
					},
				},
			},
		}
		safe := Safe(data)

		result := safe.StringAt("level1.level2.level3.value", "default")
		assert.Equal(t, "deep-value", result)
	})
}

// 性能测试

func BenchmarkSafeAccessStringAt(b *testing.B) {
	data := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "TestApp",
		},
	}
	safe := Safe(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		safe.StringAt("app.name", "default")
	}
}

func BenchmarkSafeAccessTraditional(b *testing.B) {
	data := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "TestApp",
		},
	}
	safe := Safe(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		safe.Field("app").Field("name").String("default")
	}
}
