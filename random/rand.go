/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-15 15:26:39
 * @FilePath: \go-toolbox\random\rand.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/convert"
)

// Implement Source and Source64 interfaces
type rngSource struct {
	p sync.Pool
}

func (r *rngSource) Int63() (n int64) {
	src := r.p.Get()
	n = src.(rand.Source).Int63()
	r.p.Put(src)
	return
}

// Seed specify seed when using NewRand()
func (r *rngSource) Seed(_ int64) {}

func (r *rngSource) Uint64() (n uint64) {
	src := r.p.Get()
	n = src.(rand.Source64).Uint64()
	r.p.Put(src)
	return
}

// NewRand goroutine-safe rand.Rand, optional seed value
func NewRand(seed ...int64) *rand.Rand {
	n := time.Now().UnixNano()
	if len(seed) > 0 {
		n = seed[0]
	}
	src := &rngSource{
		p: sync.Pool{
			New: func() interface{} {
				return rand.NewSource(n)
			},
		},
	}
	return rand.New(src)
}

var (
	// 设置随机种子
	mathRandSend = rand.New(rand.NewSource(time.Now().Unix()))
	// 大写字母
	matchCapital *[]int
	// 小写字母
	matchLowercase *[]int
	// 特殊符号
	matchSpecial *[]int
	// 数字
	matchNumber *[]int

	once sync.Once

	newRandSend = NewRand()
)

// RandInt
/**
 *  @Description: 随机整数
 *  @param start
 *  @param end
 *  @return v
 */
func RandInt(min, max int) (v int) {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return mathRandSend.Intn(max-min) + min
}

// RandFloat
/**
 *  @Description: 随机小数
 *  @param min
 *  @param max
 *  @return v
 */
func RandFloat(min, max float64) (v float64) {
	return min + mathRandSend.Float64()*(max-min)
}

// initASCII
/**
 *  @Description: 初始化ASCII码列表
 */
func initASCII() {
	once.Do(func() {
		fmt.Println("初始化ASCII码列表")
		// 大写字母
		c := make([]int, 26)
		for i := 0; i < 26; i++ {
			c[i] = 65 + i
		}
		// 小写字母
		matchCapital = &c
		l := make([]int, 26)
		for i := 0; i < 26; i++ {
			l[i] = 97 + i
		}
		matchLowercase = &l
		// 数字
		n := make([]int, 10)
		for i := 0; i < 10; i++ {
			n[i] = 48 + i
		}
		matchNumber = &n
		// 特殊字符(. @$!%*#_~?&^)
		s := []int{46, 64, 36, 33, 37, 42, 35, 95, 126, 63, 38, 94}
		matchSpecial = &s
	})
}

// RandString
/**
 *  @Description: 随机生成字符串
 *  @param n 字符串长度
 *  @param mode 字符串模式 random.NUMBER|random.LOWERCASE|random.SPECIAL|random.CAPITAL)
 *  @return str 生成的字符串
 */
func RandString(n int, mode RandType) (str string) {
	initASCII()
	var ascii []int
	if mode&CAPITAL >= CAPITAL {
		ascii = append(ascii, *matchCapital...)
	}
	if mode&LOWERCASE >= LOWERCASE {
		ascii = append(ascii, *matchLowercase...)
	}
	if mode&SPECIAL >= SPECIAL {
		ascii = append(ascii, *matchSpecial...)
	}
	if mode&NUMBER >= NUMBER {
		ascii = append(ascii, *matchNumber...)
	}
	if len(ascii) == 0 {
		return
	}
	var build strings.Builder
	for i := 0; i < n; i++ {
		build.WriteString(string(rune(ascii[mathRandSend.Intn(len(ascii))])))
	}
	str = build.String()
	return
}

// RandomStr 随机一个字符串
func RandomStr(length int, customBytes ...string) string {
	var sb strings.Builder
	randBytes := ALPHA_BYTES
	if len(customBytes) > 0 {
		randBytes = customBytes[0]
	}
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(randBytes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(randBytes))]))
		}
	}
	return sb.String()
}

// RandomNumber 随机一个数字字符串
func RandomNumber(length int, customBytes ...string) string {
	var sb strings.Builder
	randBytes := DEC_BYTES
	if len(customBytes) > 0 {
		randBytes = customBytes[0]
	}
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(randBytes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(randBytes))]))
		}
	}
	return sb.String()
}

// RandomHex 随机一个hex字符串
func RandomHex(bytesLen int, customBytes ...string) string {
	var sb strings.Builder
	randBytes := HEX_BYTES
	if len(customBytes) > 0 {
		randBytes = customBytes[0]
	}
	for i := 0; i < bytesLen<<1; i++ {
		sb.WriteString(string(randBytes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(randBytes))]))
	}
	return sb.String()
}

// FRandInt (>=)min - (<)max
func FRandInt(min, max int) int {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return FastIntn(max-min) + min
}

// FRandUint32 (>=)min - (<)max
func FRandUint32(min, max uint32) uint32 {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return FastRandn(max-min) + min
}

// FastIntn this is similar to rand.Intn, but faster.
// A non-negative pseudo-random number in the half-open interval [0,n).
// Return 0 if n <= 0.
func FastIntn(n int) int {
	if n <= 0 {
		return 0
	}
	if n <= math.MaxUint32 {
		return int(FastRandn(uint32(n)))
	}
	return int(newRandSend.Int63n(int64(n)))
}

// FRandString a random string, which may contain uppercase letters, lowercase letters and numbers.
// Ref: stackoverflow.icza
func FRandString(n int) string {
	return convert.B2S(FRandBytes(n))
}

// FRandHexString 指定长度的随机 hex 字符串
func FRandHexString(n int) string {
	return convert.B2S(FRandHexBytes(n))
}

// FRandAlphaString 指定长度的随机字母字符串
func FRandAlphaString(n int) string {
	return convert.B2S(FRandAlphaBytes(n))
}

// FRandDecString 指定长度的随机数字字符串
func FRandDecString(n int) string {
	return convert.B2S(FRandDecBytes(n))
}

// FRandBytes random bytes, but faster.
func FRandBytes(n int) []byte {
	return FRandBytesLetters(n, LETTER_BYTES)
}

// FRandAlphaBytes generates random alpha bytes.
func FRandAlphaBytes(n int) []byte {
	return FRandBytesLetters(n, ALPHA_BYTES)
}

// FRandHexBytes generates random hexadecimal bytes.
func FRandHexBytes(n int) []byte {
	return FRandBytesLetters(n, HEX_BYTES)
}

// RandDecBytes 指定长度的随机数字切片
func FRandDecBytes(n int) []byte {
	return FRandBytesLetters(n, DEC_BYTES)
}

// FRandBytesLetters 生成指定长度的字符切片
func FRandBytesLetters(n int, letters string) []byte {
	if n < 1 || len(letters) < 2 {
		return nil
	}
	b := make([]byte, n)
	for i, cache, remain := n-1, FastRand(), LETTER_IDX_MAX; i >= 0; {
		if remain == 0 {
			cache, remain = FastRand(), LETTER_IDX_MAX
		}
		if idx := int(cache & LETTER_IDX_MASK); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= LETTER_IDX_BITS
		remain--
	}
	return b
}
