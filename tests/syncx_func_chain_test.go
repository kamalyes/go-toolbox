/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-08-11 09:27:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-20 13:37:55
 * @FilePath: \go-toolbox\tests\syncx_func_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

// TestFuncChain 测试 FuncChain 的功能
func TestFuncChain(t *testing.T) {
	// 创建 FuncChain 实例
	fc := syncx.NewFuncChain[int]()

	// 测试正常执行的函数
	fc.AddFuncItem(syncx.NewFuncItem(func() (int, error) {
		return 1, nil
	}).WithPriority(1))

	fc.AddFuncItem(syncx.NewFuncItem(func() (int, error) {
		return 2, nil
	}).WithPriority(0))

	// 执行 FuncChain
	fc.Execute()

	// 获取 FuncItems
	funcItems := fc.GetFuncItems()

	// 验证结果
	assert.Equal(t, 2, funcItems[0].GetResult())
	assert.Equal(t, 1, funcItems[1].GetResult())
	assert.NoError(t, funcItems[0].GetError())
	assert.NoError(t, funcItems[1].GetError())
}

// TestFuncChainPanic 测试函数恐慌处理
func TestFuncChainPanic(t *testing.T) {
	fc := syncx.NewFuncChain[int]()

	// 添加一个会引发恐慌的函数
	fc.AddFuncItem(syncx.NewFuncItem(func() (int, error) {
		panic("test panic")
	}))

	// 执行 FuncChain
	fc.Execute()

	// 获取 FuncItems
	funcItems := fc.GetFuncItems()

	// 验证恐慌被处理
	assert.Error(t, funcItems[0].GetError())
	assert.Contains(t, funcItems[0].GetError().Error(), "panic:")
}

// TestFuncChainMultipleScenarios 测试多种场景
func TestFuncChainMultipleScenarios(t *testing.T) {
	fc := syncx.NewFuncChain[int]()

	// 添加正常执行的函数
	fc.AddFuncItem(syncx.NewFuncItem(func() (int, error) {
		return 10, nil
	}).WithPriority(1))

	// 添加超时的函数
	fc.AddFuncItem(syncx.NewFuncItem(func() (int, error) {
		time.Sleep(2 * time.Millisecond) // 模拟长时间运行
		return 0, nil
	}))

	// 添加引发错误的函数
	fc.AddFuncItem(syncx.NewFuncItem(func() (int, error) {
		return 0, errors.New("test error")
	}).WithPriority(0))

	// 执行 FuncChain
	fc.Execute()

	// 获取 FuncItems
	funcItems := fc.GetFuncItems()

	// 验证结果
	assert.Equal(t, 0, funcItems[0].GetResult())
	assert.Error(t, funcItems[1].GetError())
	assert.Equal(t, 10, funcItems[2].GetResult())
}

// TestFuncChainClear 测试清空 FuncChain
func TestFuncChainClear(t *testing.T) {
	fc := syncx.NewFuncChain[int]()

	// 添加函数
	fc.AddFuncItem(syncx.NewFuncItem(func() (int, error) {
		return 1, nil
	}))

	// 清空 FuncChain
	fc.Clear()

	// 获取 FuncItems
	funcItems := fc.GetFuncItems()

	// 验证 FuncChain 是否为空
	assert.Empty(t, funcItems)
}
