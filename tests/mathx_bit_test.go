/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 15:26:58
 * @FilePath: \go-toolbox\tests\mathx_bit_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"math/big"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

func TestGetBit64(t *testing.T) {
	tests := []struct {
		min, max, step uint
		expected       uint64
	}{

		{0, 3, 1, 0b1111},              // bits 0-3 set
		{1, 4, 1, 0b11110},             // bits 1-4 set
		{0, 10, 2, 0b10101010101},      // bits 0,2,4,6,8,10 set
		{5, 5, 1, 0b100000},            // only bit 5 set
		{0, 0, 1, 0b1},                 // only bit 0 set
		{1, 1, 1, 0x2},                 // only bit 1 set
		{1, 5, 2, 0x2a},                // bits 1,3,5 set (101010)
		{1, 4, 2, 0xa},                 // bits 1,3 set (1010)
		{1, 1, 0, 0},                   // step=0 返回0
		{1, 65, 1, 0},                  // max>64 返回0
		{10, 5, 1, 0},                  // min > max 返回0
		{0, 63, 1, ^uint64(0)},         // 全64位都置1
		{0, 63, 2, 0x5555555555555555}, // 偶数位全1
		{1, 63, 2, 0xAAAAAAAAAAAAAAAA}, // 奇数位全1
		{63, 63, 1, 1 << 63},           // 最高位单独设置
		{62, 63, 1, (3 << 62)},         // 最高两位设置
		{64, 70, 1, 0},                 // 超出范围，返回0（位移无效）
	}

	for _, tt := range tests {
		got := mathx.GetBit64(tt.min, tt.max, tt.step)
		assert.Equalf(t, tt.expected, got, "GetBit64(%d, %d, %d) got wrong result", tt.min, tt.max, tt.step)
	}
}

func TestBit64ToArray(t *testing.T) {
	tests := []struct {
		bit      uint64
		expected []uint
	}{
		{0b1011, []uint{0, 1, 3}},
		{0, []uint{}},
		{^uint64(0), func() []uint {
			arr := make([]uint, 64)
			for i := 0; i < 64; i++ {
				arr[i] = uint(i)
			}
			return arr
		}()},
		{1 << 63, []uint{63}},
	}

	for _, tt := range tests {
		got := mathx.Bit64ToArray(tt.bit)
		assert.Equal(t, tt.expected, got, "Bit64ToArray(%b) got wrong result", tt.bit)
	}
}

func TestGetBitBig(t *testing.T) {
	cases := []struct {
		min, max, step uint
		exp            *big.Int
	}{
		{1, 10, 0, big.NewInt(0)},
		{10, 5, 1, big.NewInt(0)},
		{0, 5, 1, func() *big.Int {
			b := big.NewInt(0)
			for i := uint(0); i <= 5; i++ {
				b.SetBit(b, int(i), 1)
			}
			return b
		}()},
		{0, 6, 2, func() *big.Int {
			b := big.NewInt(0)
			for i := uint(0); i <= 6; i += 2 {
				b.SetBit(b, int(i), 1)
			}
			return b
		}()},
		{1, 100, 10, func() *big.Int {
			b := big.NewInt(0)
			for i := uint(1); i <= 100; i += 10 {
				b.SetBit(b, int(i), 1)
			}
			return b
		}()},
	}

	for _, c := range cases {
		got := mathx.GetBitBig(c.min, c.max, c.step)
		assert.Equalf(t, 0, got.Cmp(c.exp), "GetBitBig(%d,%d,%d) failed", c.min, c.max, c.step)
	}
}

func TestBitToArrayBig(t *testing.T) {
	cases := []struct {
		input    *big.Int
		expected []uint
	}{
		{big.NewInt(0), []uint{}},
		{func() *big.Int { b := big.NewInt(0); b.SetBit(b, 0, 1); return b }(), []uint{0}},
		{func() *big.Int {
			b := big.NewInt(0)
			for _, p := range []uint{0, 2, 5, 63, 100} {
				b.SetBit(b, int(p), 1)
			}
			return b
		}(), []uint{0, 2, 5, 63, 100}},
		{func() *big.Int {
			b := big.NewInt(0)
			for i := 0; i <= 10; i++ {
				b.SetBit(b, i, 1)
			}
			return b
		}(), func() []uint {
			a := make([]uint, 11)
			for i := range a {
				a[i] = uint(i)
			}
			return a
		}()},
		{func() *big.Int {
			b := big.NewInt(0)
			for i := 0; i < 64; i++ {
				b.SetBit(b, i, 1)
			}
			return b
		}(), func() []uint {
			a := make([]uint, 64)
			for i := range a {
				a[i] = uint(i)
			}
			return a
		}()},
	}

	for _, c := range cases {
		got := mathx.BitToArrayBig(c.input)
		assert.Equalf(t, c.expected, got, "BitToArrayBig(%b) failed", c.input)
	}
}

func TestBit64GenerateAndParse(t *testing.T) {
	tests := []struct {
		min, max, step uint
		expected       []uint
	}{
		{0, 3, 1, []uint{0, 1, 2, 3}},
		{1, 4, 1, []uint{1, 2, 3, 4}},
		{0, 10, 2, []uint{0, 2, 4, 6, 8, 10}},
		{5, 5, 1, []uint{5}},
		{1, 5, 2, []uint{1, 3, 5}},
		{0, 63, 2, func() []uint {
			a := make([]uint, 32)
			for i := 0; i < 32; i++ {
				a[i] = uint(i * 2)
			}
			return a
		}()},
	}

	for _, tt := range tests {
		mask := mathx.GetBit64(tt.min, tt.max, tt.step)
		got := mathx.Bit64ToArray(mask)
		assert.Equalf(t, tt.expected, got, "Generate and parse failed for min=%d max=%d step=%d", tt.min, tt.max, tt.step)
	}
}

func TestBitBigGenerateAndParse(t *testing.T) {
	tests := []struct {
		min, max, step uint
		expected       []uint
	}{
		{0, 5, 1, []uint{0, 1, 2, 3, 4, 5}},
		{0, 6, 2, []uint{0, 2, 4, 6}},
		{1, 10, 3, []uint{1, 4, 7, 10}},
		{0, 100, 10, []uint{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}},
	}

	for _, tt := range tests {
		mask := mathx.GetBitBig(tt.min, tt.max, tt.step)
		got := mathx.BitToArrayBig(mask)
		assert.Equalf(t, tt.expected, got, "Generate and parse big.Int failed for min=%d max=%d step=%d", tt.min, tt.max, tt.step)
	}
}
