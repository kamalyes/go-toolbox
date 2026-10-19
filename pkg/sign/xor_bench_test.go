/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-17 10:06:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 10:06:55
 * @FilePath: \go-toolbox\pkg\sign\xor_bench_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkXORCipherEncrypt(b *testing.B) {
	key := byte(0xAA)
	xor := NewXORCipher(key)
	data := []byte("Hello, World!")

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		var err error
		benchXorSink, err = xor.Encrypt(data)
		if err != nil {
			assert.Fail(b, errEncryption, err)
		}
	}
}

func BenchmarkXORCipherDecrypt(b *testing.B) {
	key := byte(0xAA)
	xor := NewXORCipher(key)
	data := []byte("Hello, World!")
	encrypted, err := xor.Encrypt(data)
	if err != nil {
		b.Fatalf(errEncryption, err)
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		var err error
		benchXorSink, err = xor.Decrypt(encrypted)
		if err != nil {
			assert.Fail(b, errDecryption, err)
		}
	}
}

func BenchmarkXORCipherEncryptLong(b *testing.B) {
	key := byte(0xAA)
	xor := NewXORCipher(key)
	data := make([]byte, 1024*1024) // 1 MB 数据
	for i := range data {
		data[i] = 'A' // 填充数据
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		var err error
		benchXorSink, err = xor.Encrypt(data)
		if err != nil {
			assert.Fail(b, errEncryption, err)
		}
	}
}

func BenchmarkXORCipherDecryptLong(b *testing.B) {
	key := byte(0xAA)
	xor := NewXORCipher(key)
	data := make([]byte, 1024*1024) // 1 MB 数据
	for i := range data {
		data[i] = 'A' // 填充数据
	}
	encrypted, err := xor.Encrypt(data)
	if err != nil {
		b.Fatalf(errEncryption, err)
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		var err error
		benchXorSink, err = xor.Decrypt(encrypted)
		if err != nil {
			assert.Fail(b, errDecryption, err)
		}
	}
}

// benchXorSink 承接基准的返回值
//
// 必须有：原先这些基准把返回值丢给下划线，于是编译器在小输入上把整次分配都
// 消掉，报出 4ns、零分配的假数字，量到的其实不是这段代码
//
// [EN] Required: these benchmarks used to discard the return value, so the
// compiler elided the whole allocation on small inputs and reported a false 4ns
// with no allocations — a number that described nothing
var benchXorSink []byte

// BenchmarkCipherBySize 两种对称变换在各尺寸下的吞吐
//
// 必须带 sink：原先这些基准把返回值丢给下划线，编译器在小输入上把整次分配都
// 消掉，报出 4ns、零分配的假数字
//
// 两者都是逐字节独立的变换，所以都能成批处理：异或没有进位，可以把密钥铺成一个
// 64 位字一次做八个字节；偏移变换只取决于字节自身取值，而字节只有 256 种，
// 所以是一张 256 项的表
//
// [EN] A sink is required: these benchmarks used to discard the return value and
// the compiler elided the whole allocation on small inputs, reporting a false 4ns
// with no allocations. Both transforms are independent per byte and so can be
// batched: XOR has no carries, letting the key be replicated across a 64-bit word
// for eight bytes at a time, while the offset transform depends only on the
// byte's own value, of which there are 256 — a 256-entry table
func BenchmarkCipherBySize(b *testing.B) {
	xor := NewXORCipher(0xAA)
	off := NewProtonOffsetCipherWithPCM(31, 7, 257)

	for _, n := range []int{13, 256, 1024, 4096, 1 << 20} {
		data := make([]byte, n)
		for i := range data {
			data[i] = byte(i)
		}
		b.Run("XOR/"+sizeLabel(n), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(n))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchXorSink, _ = xor.Encrypt(data)
			}
		})
		b.Run("Offset/"+sizeLabel(n), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(n))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchXorSink, _ = off.Encrypt(data)
			}
		})
	}
}

func sizeLabel(n int) string {
	switch {
	case n >= 1<<20:
		return "1MB"
	case n >= 1024:
		return itoaSize(n/1024) + "KB"
	default:
		return itoaSize(n) + "B"
	}
}

func itoaSize(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
