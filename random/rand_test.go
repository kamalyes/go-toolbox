/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 09:15:09
 * @FilePath: \go-toolbox\random\rand_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package random

import (
	"math"
	"math/rand"
	"strings"
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/convert"
	"github.com/stretchr/testify/assert"
)

func TestAllRandomFunctions(t *testing.T) {
	t.Run("TestRandInt", TestRandInt)
	t.Run("TestRandFloat", TestRandFloat)
	t.Run("TestRandString", TestRandString)
	t.Run("TestRandomStr", TestRandomStr)
	t.Run("TestRandomNum", TestRandomNum)
	t.Run("TestRandomHex", TestRandomHex)
}

func TestRandInt(t *testing.T) {
	min := 10
	max := 20

	result := RandInt(min, max)

	assert.GreaterOrEqual(t, result, min, "Expected result to be greater than or equal to min")
	assert.LessOrEqual(t, result, max, "Expected result to be less than or equal to max")
}

func TestRandFloat(t *testing.T) {
	min := 10.5
	max := 20.5

	result := RandFloat(min, max)

	assert.GreaterOrEqual(t, result, min, "Expected result to be greater than or equal to min")
	assert.LessOrEqual(t, result, max, "Expected result to be less than or equal to max")
}

func TestRandString(t *testing.T) {
	str := RandString(10, CAPITAL|LOWERCASE|SPECIAL|NUMBER)

	assert.Len(t, str, 10, "Expected string length to be 10")
}

func TestRandomStr(t *testing.T) {
	length := 8
	str := RandomStr(length)

	assert.Len(t, str, length, "Expected string length to be 8")
}

func TestRandomNum(t *testing.T) {
	length := 6
	num := RandomNumber(length)

	assert.Len(t, num, length, "Expected number length to be 6")
}

func TestRandomHex(t *testing.T) {
	bytesLen := 4
	hex := RandomHex(bytesLen)

	assert.Len(t, hex, bytesLen*2, "Expected hex length to be twice the bytes length")
}

func TestNewRand(t *testing.T) {
	rd := NewRand(1)
	assert.Equal(t, int64(5577006791947779410), rd.Int63())

	rd = NewRand()
	for i := 1; i < 1000; i++ {
		assert.Equal(t, true, rd.Intn(i) < i)
		assert.Equal(t, true, rd.Int63n(int64(i)) < int64(i))
		assert.Equal(t, true, NewRand().Intn(i) < i)
		assert.Equal(t, true, NewRand().Int63n(int64(i)) < int64(i))
	}
}

func TestFRandInt(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, FRandInt(1, 2) == 1)
	assert.Equal(t, true, FRandInt(-1, 0) == -1)
	assert.Equal(t, true, FRandInt(0, 5) >= 0)
	assert.Equal(t, true, FRandInt(0, 5) < 5)
	assert.Equal(t, 2, FRandInt(2, 2))
	assert.Equal(t, 2, FRandInt(3, 2))
}

func TestFRandUint32(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, FRandUint32(1, 2) == 1)
	assert.Equal(t, true, FRandUint32(0, 5) < 5)
	assert.Equal(t, uint32(2), FRandUint32(2, 2))
	assert.Equal(t, uint32(2), FRandUint32(3, 2))
}

func TestFastIntn(t *testing.T) {
	t.Parallel()
	for i := 1; i < 10000; i++ {
		assert.Equal(t, true, FastRandn(uint32(i)) < uint32(i))
		assert.Equal(t, true, FastIntn(i) < i)
	}
	assert.Equal(t, 0, FastIntn(-2))
	assert.Equal(t, 0, FastIntn(0))
	assert.Equal(t, true, FastIntn(math.MaxUint32) < math.MaxUint32)
	assert.Equal(t, true, FastIntn(math.MaxInt64) < math.MaxInt64)
}

func TestFRandString(t *testing.T) {
	t.Parallel()
	fns := []func(n int) string{FRandString, FRandAlphaString, FRandHexString, FRandDecString}
	ss := []string{LETTER_BYTES, ALPHA_BYTES, HEX_BYTES, DEC_BYTES}
	for i, fn := range fns {
		a, b := fn(777), fn(777)
		assert.Equal(t, 777, len(a))
		assert.NotEqual(t, a, b)
		assert.Equal(t, "", fn(-1))
		for _, s := range ss[i] {
			assert.True(t, strings.ContainsRune(a, s))
		}
	}
}

func TestFRandBytesLetters(t *testing.T) {
	t.Parallel()
	letters := ""
	assert.Nil(t, FRandBytesLetters(10, letters))
	letters = "a"
	assert.Nil(t, FRandBytesLetters(10, letters))
	letters = "ab"
	s := convert.B2S(FRandBytesLetters(10, letters))
	assert.Equal(t, 10, len(s))
	assert.True(t, strings.Contains(s, "a"))
	assert.True(t, strings.Contains(s, "b"))
	letters = "xxxxxxxxxxxx"
	s = convert.B2S(FRandBytesLetters(100, letters))
	assert.Equal(t, 100, len(s))
	assert.Equal(t, strings.Repeat("x", 100), s)
}

var (
	testString = "  Fufu 中　文\u2728->?\n*\U0001F63A   "
	testBytes  = []byte(testString)
)

func TestB2S(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		b := FRandBytes(64)
		assert.Equal(t, string(b), convert.B2S(b))
	}

	expected := testString
	actual := convert.B2S([]byte(expected))
	assert.Equal(t, expected, actual)

	assert.Equal(t, true, convert.B2S(nil) == "")
	assert.Equal(t, testString, convert.B2S(testBytes))
}

func TestS2B(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		s := RandomNumber(64)
		expected := []byte(s)
		actual := convert.S2B(s)
		assert.Equal(t, expected, actual)
		assert.Equal(t, len(expected), len(actual))
	}

	expected := testString
	actual := convert.S2B(expected)
	assert.Equal(t, []byte(expected), actual)

	assert.Equal(t, true, convert.S2B("") == nil)
	assert.Equal(t, testBytes, convert.S2B(testString))
}

func BenchmarkRandBytesParallel(b *testing.B) {
	b.Run("FRandBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandBytes(20)
			}
		})
	})
	b.Run("FRandAlphaBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandAlphaBytes(20)
			}
		})
	})
	b.Run("FRandHexBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandHexBytes(20)
			}
		})
	})
	b.Run("FRandDecBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandDecBytes(20)
			}
		})
	})
}

func BenchmarkRandInt(b *testing.B) {
	b.Run("RandInt", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = RandInt(0, i)
		}
	})
	b.Run("FRandInt", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = FRandInt(0, i)
		}
	})
	b.Run("FRandUint32", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = FRandUint32(0, uint32(i))
		}
	})
	b.Run("FastIntn", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = FastIntn(i)
		}
	})
	b.Run("Rand.Intn", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = mathRandSend.Intn(i)
		}
	})
	b.Run("std.rand.Intn", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = rand.Intn(i)
		}
	})
}

func BenchmarkRandInt32Parallel(b *testing.B) {
	b.Run("FRandInt", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandInt(0, math.MaxInt32)
			}
		})
	})
	b.Run("FRandUint32", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandUint32(0, math.MaxInt32)
			}
		})
	})
	b.Run("FastIntn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FastIntn(math.MaxInt32)
			}
		})
	})
	b.Run("Rand.Intn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = newRandSend.Intn(math.MaxInt32)
			}
		})
	})
	var mu sync.Mutex
	b.Run("std.rand.Intn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mu.Lock()
				_ = rand.Intn(math.MaxInt32)
				mu.Unlock()
			}
		})
	})
}
