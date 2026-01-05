/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 10:00:58
 * @FilePath: \go-toolbox\pkg\contextx\utils_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package contextx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMergeContext 测试合并多个上下文
func TestMergeContext(t *testing.T) {
	// 创建上下文并设置一些值
	ctx1 := NewContext().
		WithValue(TestKey1, TestValue1).
		WithValue(TestKey2, TestValue2)

	ctx2 := NewContext().
		WithValue(TestKey2, TestNewValue2). // 这个值会覆盖 ctx1 中的值
		WithValue(TestKey3, TestValue3)

	ctx3 := NewContext().
		WithValue(TestKey4, TestValue4)

	// 合并上下文
	merged := MergeContext(ctx1, ctx2, ctx3)

	// 断言合并后的值
	assert.Equal(t, TestValue1, merged.Value(TestKey1), "期望值为 'value1'")
	assert.Equal(t, TestNewValue2, merged.Value(TestKey2), "期望值为 'newValue2'，应覆盖之前的值")
	assert.Equal(t, TestValue3, merged.Value(TestKey3), "期望值为 'value3'")
	assert.Equal(t, TestValue4, merged.Value(TestKey4), "期望值为 'value4'")
	assert.Nil(t, merged.Value("key5"), "期望值为 nil，因为 key5 不存在")
}

// TestMergeContextEmpty 测试合并空上下文
func TestMergeContextEmpty(t *testing.T) {
	merged := MergeContext()

	assert.NotNil(t, merged, "期望合并后的上下文不为 nil")
	assert.Equal(t, context.Background(), merged.Context, "期望合并后的上下文为背景上下文")
}
