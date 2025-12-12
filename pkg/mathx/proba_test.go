/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\mathx\proba_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTrueOnProba函数用于测试TrueOnProba方法的功能。
func TestTrueOnProba(t *testing.T) {
	const proba = math.Pi / 10
	const total = 100000
	const epsilon = 0.05
	var count int
	p := NewProba()
	for i := 0; i < total; i++ {
		if p.TrueOnProba(proba) {
			count++
		}
	}

	ratio := float64(count) / float64(total)
	assert.InEpsilon(t, proba, ratio, epsilon)
}
