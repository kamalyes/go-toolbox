/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-11 09:15:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-11 10:15:16
 * @FilePath: \go-toolbox\tests\crc_test.go
 * @Description: CRC算法单元测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/crc"
	"github.com/stretchr/testify/assert"
)

var crcData = "123456789"

// 测试不同CRC算法的计算
func TestCRCs(t *testing.T) {
	// 测试数据和预期 CRC 值
	var testCases = []struct {
		name     string      // 测试名称
		data     []byte      // 输入数据
		expected uint64      // 预期的 CRC 值
		factory  crc.Factory // 对应的工厂
	}{
		{"CRC4_ITU Test", []byte(crcData), uint64(0xa), crc.CRC4_ITUFactory},
		{"CRC5_EPC Test", []byte(crcData), uint64(0x1e), crc.CRC5_EPCFactory},
		{"CRC5_ITU Test", []byte(crcData), uint64(0x6), crc.CRC5_ITUFactory},
		{"CRC5_USB Test", []byte(crcData), uint64(0xA), crc.CRC5_USBFactory},
		{"CRC6_ITU Test", []byte(crcData), uint64(0xA), crc.CRC6_ITUFactory},
		{"CRC7_MMC Test", []byte(crcData), uint64(0x63), crc.CRC7_MMCFactory},
		{"CRC8 Test", []byte(crcData), uint64(0xF4), crc.CRC8Factory},
		{"CRC8_ATM Test", []byte(crcData), uint64(0xF4), crc.CRC8_ATMFactory},
		{"CRC8_CDMA2000 Test", []byte(crcData), uint64(0xEA), crc.CRC8_CDMA2000Factory},
		{"CRC8_DALLAS_1WIRE Test", []byte(crcData), uint64(0xA2), crc.CRC8_DALLAS_1WIREFactory},
		{"CRC8_ITU Test", []byte(crcData), uint64(0xA1), crc.CRC8_ITUFactory},
		{"CRC8_ROHC Test", []byte(crcData), uint64(0xD0), crc.CRC8_ROHCFactory},
		{"CRC8_MAXIM Test", []byte(crcData), uint64(0xA1), crc.CRC8_MAXIMFactory},
		{"CRC16_IBM Test", []byte(crcData), uint64(0xbb3d), crc.CRC16_IBMFactory},
		{"CRC16_MAXIM Test", []byte(crcData), uint64(0x44C2), crc.CRC16_MAXIMFactory},
		{"CRC16_USB Test", []byte(crcData), uint64(0xB4C8), crc.CRC16_USBFactory},
		{"CRC16_MODBUS Test", []byte(crcData), uint64(0x4B37), crc.CRC16_MODBUSFactory},
		{"CRC16_ANSI Test", []byte(crcData), uint64(0xA47B), crc.CRC16_ANSIFactory},
		{"CRC16_XMODEM Test", []byte(crcData), uint64(0xA47B), crc.CRC16_XMODEMFactory},
		{"CRC16_CCITT Test", []byte(crcData), uint64(0x29B1), crc.CRC16_CCITTFactory},
		{"CRC16_CCITT_FALSE Test", []byte(crcData), uint64(0x31C3), crc.CRC16_CCITT_FALSEFactory},
		{"CRC16_X25 Test", []byte(crcData), uint64(0x906E), crc.CRC16_X25Factory},
		{"CRC16_DNP Test", []byte(crcData), uint64(0xEA82), crc.CRC16_DNPFactory},
		{"CRC16_CCITT_KERMIT Test", []byte(crcData), uint64(0x2189), crc.CRC16_CCITT_KERMITFactory},
		{"CRC16_GENERIC Test", []byte(crcData), uint64(0xA47B), crc.CRC16_GENERICFactory},
		{"CRC16_CCITT_TRUE Test", []byte(crcData), uint64(0x2189), crc.CRC16_CCITT_TRUEFactory},
		{"CRC24_OPENPGP Test", []byte(crcData), uint64(0x21CF02), crc.CRC24_OPENPGPFactory},
		{"CRC32 Test", []byte(crcData), uint64(0xCBF43926), crc.CRC32Factory},
		{"CRC32_MPEG2 Test", []byte(crcData), uint64(0x376E6E7), crc.CRC32_MPEG2Factory},
		{"CRC32_PKZIP Test", []byte(crcData), uint64(0xFC4F2BE9), crc.CRC32_PKZIPFactory},
		{"CRC32C Test", []byte(crcData), uint64(0xF28417BE), crc.CRC32CFactory},
		{"CRC32_CASTAGNOLI Test", []byte(crcData), uint64(0xE3069283), crc.CRC32_CASTAGNOLIFactory},
		{"CRC32_ADLER32 Test", []byte(crcData), uint64(0x2868AEA8), crc.CRC32_ADLER32Factory},
		{"CRC64_ECMA Test", []byte(crcData), uint64(0x2B9C7EE4E2780C8A), crc.CRC64_ECMAFactory},
		{"CRC64_ISO Test", []byte(crcData), uint64(0x66A2364420E6C605), crc.CRC64_ISOFactory},
		{"CRC64_WEIERSTRASS Test", []byte(crcData), uint64(0x66A2364420E6C605), crc.CRC64_WEIERSTRASSFactory},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calculator, err := tc.factory.Create()
			assert.NoError(t, err, "创建计算器失败: %s", tc.name)

			result := calculator.Compute(tc.data)

			assert.Equal(t, tc.expected, result, "计算结果与预期不符: %s", tc.name)
		})
	}
}

// 测试New函数
func TestNew(t *testing.T) {
	calculator, err := crc.New(crc.CRC8)
	assert.NoError(t, err, "创建计算器时出现错误")
	assert.NotNil(t, calculator, "计算器实例应不为nil")
}

// 测试Compute函数
func TestCompute(t *testing.T) {
	calculator, err := crc.New(crc.CRC8)
	assert.NoError(t, err)

	data := []byte(crcData)
	expectedCRC := uint64(0xF4) // 预期的CRC值
	result := calculator.Compute(data)
	assert.Equal(t, expectedCRC, result, "计算的CRC值与预期不符")
}

// 测试Reset函数
func TestReset(t *testing.T) {
	calculator, err := crc.New(crc.CRC8)
	assert.NoError(t, err)

	// 计算CRC值
	data := []byte(crcData)
	calculator.Compute(data)

	// 重置计算器
	calculator.Reset()
	resultAfterReset := calculator.Compute(data)
	assert.Equal(t, uint64(0xF4), resultAfterReset, "重置后的计算结果应与初始计算结果一致")
}
