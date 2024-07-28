/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 21:35:12
 * @FilePath: \go-toolbox\random\rand.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	// CAPITAL 包含大写字母
	CAPITAL = 1
	// LOWERCASE 包含小写字母
	LOWERCASE = 2
	// SPECIAL 包含特殊字符
	SPECIAL = 4
	// NUMBER 包含数字
	NUMBER = 8
	// 自定义NUM区间
	NUM_SOURCE = "0123456789"
	// 自定义HEX区间
	HEX_SOURCE = "ABCDEF0123456789"
	// 自定义字符区间
	STR_SOURCE = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var (
	// 设置随机种子
	r = rand.New(rand.NewSource(time.Now().Unix()))
	// 大写字母
	capital *[]int
	// 小写字母
	lowercase *[]int
	// 特殊符号
	special *[]int
	// 数字
	number *[]int

	once sync.Once
)

// RandInt
/**
 *  @Description: 随机整数
 *  @param start
 *  @param end
 *  @return v
 */
func RandInt(min, max int) (v int) {
	return r.Intn(max-min) + min
}

// RandFloat
/**
 *  @Description: 随机小数
 *  @param min
 *  @param max
 *  @return v
 */
func RandFloat(min, max float64) (v float64) {
	return min + r.Float64()*(max-min)
}

// initASCII
/**
 *  @Description: 初始化ASCII码列表
 */
func initASCII() {
	once.Do(func() {
		fmt.Println("初始化列表")
		// 大写字母
		c := make([]int, 26)
		for i := 0; i < 26; i++ {
			c[i] = 65 + i
		}
		// 小写字母
		capital = &c
		l := make([]int, 26)
		for i := 0; i < 26; i++ {
			l[i] = 97 + i
		}
		lowercase = &l
		// 数字
		n := make([]int, 10)
		for i := 0; i < 10; i++ {
			n[i] = 48 + i
		}
		number = &n
		// 特殊字符(. @$!%*#_~?&^)
		s := []int{46, 64, 36, 33, 37, 42, 35, 95, 126, 63, 38, 94}
		special = &s
	})
}

// RandString
/**
 *  @Description: 随机生成字符串
 *  @param n 字符串长度
 *  @param mode 字符串模式 random.NUMBER|random.LOWERCASE|random.SPECIAL|random.CAPITAL)
 *  @return str 生成的字符串
 */
func RandString(n int, mode int) (str string) {
	initASCII()
	var ascii []int
	if mode&CAPITAL >= CAPITAL {
		ascii = append(ascii, *capital...)
	}
	if mode&LOWERCASE >= LOWERCASE {
		ascii = append(ascii, *lowercase...)
	}
	if mode&SPECIAL >= SPECIAL {
		ascii = append(ascii, *special...)
	}
	if mode&NUMBER >= NUMBER {
		ascii = append(ascii, *number...)
	}
	if len(ascii) == 0 {
		return
	}
	var build strings.Builder
	for i := 0; i < n; i++ {
		build.WriteString(string(rune(ascii[r.Intn(len(ascii))])))
	}
	str = build.String()
	return
}

// RandomStr 随机一个字符串
func RandomStr(length int) string {
	var sb strings.Builder
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(STR_SOURCE[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(STR_SOURCE))]))
		}
	}
	return sb.String()
}

// RandomNum 随机一个数字字符串
func RandomNum(length int) string {
	var sb strings.Builder
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(NUM_SOURCE[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(NUM_SOURCE))]))
		}
	}
	return sb.String()
}

// RandomHex 随机一个hex字符串
func RandomHex(bytesLen int) string {
	var sb strings.Builder
	for i := 0; i < bytesLen<<1; i++ {
		sb.WriteString(string(HEX_SOURCE[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(HEX_SOURCE))]))
	}
	return sb.String()
}
