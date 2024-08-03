/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 21:35:27
 * @FilePath: \go-toolbox\randx\rand_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package randx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRand(t *testing.T) {
	rd := NewRand(1)
	assert.Equal(t, int64(5577006791947779410), rd.Int63())

	rd = NewRand()
	for i := 1; i < 1000; i++ {
		assert.Equal(t, true, rd.Intn(i) < i)
		assert.Equal(t, true, rd.Int63n(int64(i)) < int64(i))
		assert.Equal(t, true, NewRand().Intn(i) < i)
		assert.Equal(t, true, NewRand().Int63n(int64(i)) < int64(i))
	}
}
