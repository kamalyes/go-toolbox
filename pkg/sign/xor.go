/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-16 18:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 09:55:15
 * @FilePath: \go-toolbox\pkg\sign\xor.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"encoding/binary"
	"errors"
)

// XORCipher 结构体定义
type XORCipher struct {
	Key byte // 加密密钥
}

// NewXORCipher 创建一个新的 XORCipher 实例
func NewXORCipher(key byte) *XORCipher {
	return &XORCipher{Key: key}
}

// Encrypt 加密函数
func (xc *XORCipher) Encrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("数据不能为空")
	}

	return xorBytes(data, xc.Key), nil
}

// xorBytes 按 8 字节一组做异或
//
// 原实现逐字节循环，而 Go 不会自动向量化，于是 1MB 要 301 微秒（约 3.5 GB/s）
// 异或是逐字节独立的运算，没有进位也没有依赖，所以可以把密钥字节复制成一个
// 64 位字，一次处理八个字节binary.LittleEndian.Uint64 在 amd64 与 arm64 上
// 编译成一条非对齐载入，不分配也不额外检查边界
// 字节序无关紧要：密钥字节在那个 64 位字里是重复的，怎么摆都一样
//
// 顺带把 xc.Key 取进局部变量——它是指针后的结构体字段，编译器无法确定循环中
// 不会被改写，原先每字节都要重读一次
//
// [EN] XORs eight bytes at a time. The original looped byte by byte and Go does
// not auto-vectorise, so 1MB took 301 microseconds, about 3.5 GB/s. XOR is
// independent per byte with no carries and no dependencies, so the key byte can
// be replicated across a 64-bit word and eight bytes handled at once;
// binary.LittleEndian.Uint64 compiles to a single unaligned load on amd64 and
// arm64. Byte order is irrelevant because the key byte repeats throughout the
// word. The key is also hoisted into a local: as a struct field behind a pointer
// the compiler cannot prove it unchanged, so it was re-read every byte
func xorBytes(data []byte, key byte) []byte {
	out := make([]byte, len(data))

	i, n := 0, len(data)
	if n >= 8 {
		// 把密钥字节铺满一个 64 位字
		k := uint64(key)
		k |= k << 8
		k |= k << 16
		k |= k << 32
		for ; i+8 <= n; i += 8 {
			binary.LittleEndian.PutUint64(out[i:], binary.LittleEndian.Uint64(data[i:])^k)
		}
	}
	// 收尾不足八字节的部分
	for ; i < n; i++ {
		out[i] = data[i] ^ key
	}
	return out
}

// Decrypt 解密函数
//
// 异或是自逆运算，解密与加密是同一件事
func (xc *XORCipher) Decrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("加密数据不能为空")
	}
	return xorBytes(data, xc.Key), nil
}
