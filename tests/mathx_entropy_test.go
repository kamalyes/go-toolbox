/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 15:26:58
 * @FilePath: \go-toolbox\tests\mathx_entropy_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"math"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

func TestCalcEntropy(t *testing.T) {
	assert := assert.New(t)

	// 测试用例1：均匀分布
	data1 := map[interface{}]int{
		"a": 2,
		"b": 2,
		"c": 2,
		"d": 2,
	}
	expectedEntropy1 := 2.0 // log2(4)
	assert.InDelta(expectedEntropy1, mathx.CalcEntropy(data1), 1e-9)

	// 测试用例2：非均匀分布
	data2 := map[interface{}]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}
	// 手动计算或使用已知值
	expectedEntropy2 := -(1.0/10.0)*math.Log2(1.0/10.0) - (2.0/10.0)*math.Log2(2.0/10.0) - (3.0/10.0)*math.Log2(3.0/10.0) - (4.0/10.0)*math.Log2(4.0/10.0)
	assert.InDelta(expectedEntropy2, mathx.CalcEntropy(data2), 1e-9)

	// 测试用例3：空分布
	data3 := map[interface{}]int{}
	expectedEntropy3 := 0.0
	assert.Equal(expectedEntropy3, mathx.CalcEntropy(data3))

	// 测试用例4：单个元素分布
	data4 := map[interface{}]int{
		"a": 10,
	}
	expectedEntropy4 := 0.0 // 单个元素没有不确定性
	assert.Equal(expectedEntropy4, mathx.CalcEntropy(data4))
}
