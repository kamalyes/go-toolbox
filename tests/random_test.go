/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-03 13:51:48
 * @FilePath: \go-toolbox\tests\random_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"encoding/json"
	"math"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/random"
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

	result := random.RandInt(min, max)

	assert.GreaterOrEqual(t, result, min, "Expected result to be greater than or equal to min")
	assert.LessOrEqual(t, result, max, "Expected result to be less than or equal to max")
}

func TestRandFloat(t *testing.T) {
	min := 10.5
	max := 20.5

	result := random.RandFloat(min, max)

	assert.GreaterOrEqual(t, result, min, "Expected result to be greater than or equal to min")
	assert.LessOrEqual(t, result, max, "Expected result to be less than or equal to max")
}

func TestRandString(t *testing.T) {
	str := random.RandString(10, random.CAPITAL|random.LOWERCASE|random.SPECIAL|random.NUMBER)

	assert.Len(t, str, 10, "Expected string length to be 10")
}

func TestRandomStr(t *testing.T) {
	length := 8
	str := random.RandomStr(length)

	assert.Len(t, str, length, "Expected string length to be 8")
}

func TestRandomNum(t *testing.T) {
	length := 6
	num := random.RandomNumber(length)

	assert.Len(t, num, length, "Expected number length to be 6")
}

func TestRandomHex(t *testing.T) {
	bytesLen := 4
	hex := random.RandomHex(bytesLen)

	assert.Len(t, hex, bytesLen*2, "Expected hex length to be twice the bytes length")
}

func TestNewRand(t *testing.T) {
	rd := random.NewRand(1)
	assert.Equal(t, int64(5577006791947779410), rd.Int63())

	rd = random.NewRand()
	for i := 1; i < 1000; i++ {
		assert.Equal(t, true, rd.Intn(i) < i)
		assert.Equal(t, true, rd.Int63n(int64(i)) < int64(i))
		assert.Equal(t, true, random.NewRand().Intn(i) < i)
		assert.Equal(t, true, random.NewRand().Int63n(int64(i)) < int64(i))
	}
}

func TestFRandInt(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, random.FRandInt(1, 2) == 1)
	assert.Equal(t, true, random.FRandInt(-1, 0) == -1)
	assert.Equal(t, true, random.FRandInt(0, 5) >= 0)
	assert.Equal(t, true, random.FRandInt(0, 5) < 5)
	assert.Equal(t, 2, random.FRandInt(2, 2))
	assert.Equal(t, 2, random.FRandInt(3, 2))
}

func TestFRandUint32(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, random.FRandUint32(1, 2) == 1)
	assert.Equal(t, true, random.FRandUint32(0, 5) < 5)
	assert.Equal(t, uint32(2), random.FRandUint32(2, 2))
	assert.Equal(t, uint32(2), random.FRandUint32(3, 2))
}

func TestFastIntn(t *testing.T) {
	t.Parallel()
	for i := 1; i < 10000; i++ {
		assert.Equal(t, true, random.FastRandn(uint32(i)) < uint32(i))
		assert.Equal(t, true, random.FastIntn(i) < i)
	}
	assert.Equal(t, 0, random.FastIntn(-2))
	assert.Equal(t, 0, random.FastIntn(0))
	assert.Equal(t, true, random.FastIntn(math.MaxUint32) < math.MaxUint32)
	assert.Equal(t, true, random.FastIntn(math.MaxInt64) < math.MaxInt64)
}

func TestFRandString(t *testing.T) {
	t.Parallel()
	fns := []func(n int) string{random.FRandString, random.FRandAlphaString, random.FRandHexString, random.FRandDecString}
	ss := []string{random.LETTER_BYTES, random.ALPHA_BYTES, random.HEX_BYTES, random.DEC_BYTES}
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
	assert.Nil(t, random.FRandBytesLetters(10, letters))
	letters = "a"
	assert.Nil(t, random.FRandBytesLetters(10, letters))
	letters = "ab"
	s := convert.B2S(random.FRandBytesLetters(10, letters))
	assert.Equal(t, 10, len(s))
	assert.True(t, strings.Contains(s, "a"))
	assert.True(t, strings.Contains(s, "b"))
	letters = "xxxxxxxxxxxx"
	s = convert.B2S(random.FRandBytesLetters(100, letters))
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
		b := random.FRandBytes(64)
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
		s := random.RandomNumber(64)
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

func TestFRandBytesJSON(t *testing.T) {
	length := 16 // 测试生成的随机字节长度
	// 调用 FRandBytesJSON 函数
	jsonStr, err := random.FRandBytesJSON(length)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 检查 JSON 字符串是否有效
	var result []byte
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	// 检查生成的随机字节长度
	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}
}

// 性能测试函数
func BenchmarkFRandBytesJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := random.FRandBytesJSON(1024) // 测试生成1024字节的随机字节字符串
		if err != nil {
			b.Error(err) // 如果有错误，记录
		}
	}
}

// 定义测试模型
type TestModel struct {
	Name      string         `json:"name"`
	Age       int            `json:"age"`
	Salary    float64        `json:"salary"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	Tags      []string       `json:"tags"`
	Settings  map[string]int `json:"settings"`
}

// TestGenerateRandomModel 测试 GenerateRandomModel 函数
func TestGenerateRandomModel(t *testing.T) {
	// 创建一个 TestModel 的实例
	model := &TestModel{}

	// 调用 GenerateRandomModel
	modelResult, jsonResult, err := random.GenerateRandomModel(model)
	if err != nil || modelResult == nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log("jsonResult", jsonResult)
	// 验证返回的 JSON 字符串是否有效
	var resultMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResult), &resultMap); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	// 验证字段是否被填充
	if resultMap["name"] == "" || resultMap["age"] == nil || resultMap["salary"] == nil || resultMap["is_active"] == nil {
		t.Error("Expected fields to be populated, but they are not")
	}
}

// BenchmarkGenerateRandomModel 性能测试 GenerateRandomModel 函数
func BenchmarkGenerateRandomModel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := &TestModel{}
		_, _, err := random.GenerateRandomModel(model)
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}

func BenchmarkRandBytesParallel(b *testing.B) {
	b.Run("FRandBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandBytes(20)
			}
		})
	})
	b.Run("FRandAlphaBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandAlphaBytes(20)
			}
		})
	})
	b.Run("FRandHexBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandHexBytes(20)
			}
		})
	})
	b.Run("FRandDecBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandDecBytes(20)
			}
		})
	})
}

func BenchmarkRandInt(b *testing.B) {
	b.Run("RandInt", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = random.RandInt(0, i)
		}
	})
	b.Run("FRandInt", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = random.FRandInt(0, i)
		}
	})
	b.Run("FRandUint32", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = random.FRandUint32(0, uint32(i))
		}
	})
	b.Run("FastIntn", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = random.FastIntn(i)
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
				_ = random.FRandInt(0, math.MaxInt32)
			}
		})
	})
	b.Run("FRandUint32", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandUint32(0, math.MaxInt32)
			}
		})
	})
	b.Run("FastIntn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FastIntn(math.MaxInt32)
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
