/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:15:15
 * @FilePath: \go-toolbox\tests\mathx_int_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

func TestMaxInt(t *testing.T) {
	cases := []struct {
		a      int
		b      int
		expect int
	}{
		{
			a:      0,
			b:      1,
			expect: 1,
		},
		{
			a:      0,
			b:      -1,
			expect: 0,
		},
		{
			a:      1,
			b:      1,
			expect: 1,
		},
	}

	for _, each := range cases {
		each := each
		t.Run(t.Name(), func(t *testing.T) {
			actual := mathx.MaxInt(each.a, each.b)
			assert.Equal(t, each.expect, actual)
		})
	}
}

func TestMinInt(t *testing.T) {
	cases := []struct {
		a      int
		b      int
		expect int
	}{
		{
			a:      0,
			b:      1,
			expect: 0,
		},
		{
			a:      0,
			b:      -1,
			expect: -1,
		},
		{
			a:      1,
			b:      1,
			expect: 1,
		},
	}

	for _, each := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			actual := mathx.MinInt(each.a, each.b)
			assert.Equal(t, each.expect, actual)
		})
	}
}
