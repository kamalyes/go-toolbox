/*
* @Author: kamalyes 501893067@qq.com
* @Date: 2025-06-09 17:15:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-11 10:15:26
 * @FilePath: \go-toolbox\pkg\crc\crc.go
* @Description: CRC算法核心实现
*
* Copyright (c) 2025 by kamalyes, All Rights Reserved.
*/
package crc

import (
	"fmt"
	"math/bits"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// Config 表示CRC算法的配置结构体
type Config struct {
	Width  uint   // 校验位宽 (8/16/24/32/64)
	Poly   uint64 // 生成多项式 (需注意位宽对齐)
	Init   uint64 // 初始值
	RefIn  bool   // 输入数据字节位反转
	RefOut bool   // 输出结果位反转
	XorOut uint64 // 最终异或值
}

// Calculator 定义CRC计算器接口
type Calculator interface {
	Compute(data []byte) uint64 // 计算数据的CRC校验值
	Reset()                     // 重置计算器状态
}

// calculatorImpl 是CRC计算器的核心实现
type calculatorImpl struct {
	config Config     // CRC算法配置
	table  []uint64   // 查表法预计算表
	mask   uint64     // 位掩码
	crc    uint64     // 当前CRC值
	mu     sync.Mutex // 添加锁
	once   sync.Once  // 确保表只初始化一次
}

// New 创建CRC计算器实例
// 参数: cfg - CRC算法配置
// 返回: 计算器实例
func New(cfg Config) (Calculator, error) {
	if cfg.Width == 0 || cfg.Width > 64 {
		return nil, fmt.Errorf("invalid width: %d, must be between 1 and 64", cfg.Width)
	}
	if cfg.Poly == 0 {
		return nil, fmt.Errorf("invalid polynomial: must be non-zero")
	}
	c := &calculatorImpl{config: cfg}
	c.mask = (1 << cfg.Width) - 1
	c.Reset() // 初始化CRC值
	return c, nil
}

// generateTable 生成CRC查表法预计算表 (256项)
// 使用sync.Once确保线程安全
func (c *calculatorImpl) generateTable() {
	c.once.Do(func() {
		c.table = make([]uint64, 256)
		for i := range c.table {
			crc := uint64(i)
			if c.config.Width > 8 {
				crc <<= (c.config.Width - 8)
			}

			for j := 0; j < 8; j++ {
				if crc&(1<<(c.config.Width-1)) != 0 {
					crc = (crc << 1) ^ c.config.Poly
				} else {
					crc <<= 1
				}
			}
			c.table[i] = crc & c.mask // 保存CRC值
		}
	})
}

// Compute 计算数据的CRC校验值
func (c *calculatorImpl) Compute(data []byte) uint64 {
	return syncx.WithLockReturnValue(&c.mu, func() uint64 {
		c.generateTable() // 初始化查表 (首次调用时)
		length := len(data)

		for i := 0; i < length; i++ {
			b := data[i]
			// 输入位反转处理
			if c.config.RefIn {
				b = bits.Reverse8(b)
			}

			// 核心查表计算
			if c.config.Width > 8 {
				idx := byte(c.crc>>(c.config.Width-8)) ^ b
				c.crc = (c.crc << 8) ^ c.table[idx]
			} else {
				c.crc = c.table[byte(c.crc)^b]
			}
			c.crc &= c.mask // 应用位掩码
		}

		// 输出位反转处理
		result := c.crc
		if c.config.RefOut {
			result = reverseOutput(result, c.config.Width)
		}

		// 重置计算器状态
		c.Reset()
		return (result ^ c.config.XorOut) & c.mask // 最终异或处理
	})
}

// Reset 重置计算器状态
func (c *calculatorImpl) Reset() {
	c.crc = c.config.Init & c.mask // 使用初始值重置CRC
}

// reverseOutput 自定义输出位反转
// 参数: val - 待反转值, width - 位宽
// 返回: 反转后的值
func reverseOutput(val uint64, width uint) uint64 {
	switch width {
	case 8:
		return uint64(bits.Reverse8(uint8(val)))
	case 16:
		return uint64(bits.Reverse16(uint16(val)))
	case 32:
		return uint64(bits.Reverse32(uint32(val)))
	case 64:
		return bits.Reverse64(val)
	default:
		return reverseCustom(val, width) // 自定义位宽反转
	}
}

// reverseCustom 自定义位宽反转
func reverseCustom(val uint64, width uint) uint64 {
	rev := uint64(0)
	for i := uint(0); i < width; i++ {
		if val&(1<<i) != 0 {
			rev |= 1 << (width - 1 - i)
		}
	}
	return rev
}
