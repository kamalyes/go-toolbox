/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 13:18:15
 * @FilePath: \go-toolbox\pkg\sign\bcrypt_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFromPassword(t *testing.T) {
	t.Run("默认cost", func(t *testing.T) {
		hashed, err := GenerateFromPassword([]byte("hello"))
		require.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})

	t.Run("自定义cost", func(t *testing.T) {
		hashed, err := GenerateFromPassword([]byte("hello"), 4)
		require.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})

	t.Run("cost为0时使用默认值", func(t *testing.T) {
		hashed, err := GenerateFromPassword([]byte("hello"), 0)
		require.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})
}

func TestCompareHashAndPassword(t *testing.T) {
	plain := []byte("test123")
	hashed, err := GenerateFromPassword(plain)
	require.NoError(t, err)

	t.Run("匹配应返回nil", func(t *testing.T) {
		err := CompareHashAndPassword(hashed, plain)
		assert.NoError(t, err)
	})

	t.Run("不匹配应返回错误", func(t *testing.T) {
		err := CompareHashAndPassword(hashed, []byte("wrong"))
		assert.Error(t, err)
	})
}
