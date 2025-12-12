/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-22 10:07:57
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:30:15
 * @FilePath: \go-toolbox\pkg\stringx\replace_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceChain(t *testing.T) {
	result := New("Hello World").ReplaceChain("World", "Golang", 1).Value()
	assert.Equal(t, "Hello Golang", result)
}

func TestReplaceAllChain(t *testing.T) {
	result := New("Hello World World").ReplaceAllChain("World", "Golang").Value()
	assert.Equal(t, "Hello Golang Golang", result)
}

func TestReplaceWithMatcherChain(t *testing.T) {
	result := New("Hello 123").ReplaceWithMatcherChain(`\d+`, func(s string) string {
		return "456"
	}).Value()
	assert.Equal(t, "Hello 456", result)
}

func TestHideChain(t *testing.T) {
	result := New("SensitiveData").HideChain(3, 8).Value()
	assert.Equal(t, "Sen*****eData", result)
}

func TestReplaceSpecialCharsChain(t *testing.T) {
	result := New("Hello, World!").ReplaceSpecialCharsChain('*').Value()
	assert.Equal(t, "Hello**World*", result)
}
