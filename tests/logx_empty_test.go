/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-18 11:15:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-18 11:37:55
 * @FilePath: \go-toolbox\tests\empty_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"context"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/logx"
	"github.com/stretchr/testify/assert"
)

func TestEmptyLog(t *testing.T) {
	// 创建一个 EmptyLog 实例
	logger := logx.NewEmptyLog()

	// 创建上下文
	ctx := context.Background()

	// 调用所有日志方法，验证没有崩溃或抛出错误
	assert.NotPanics(t, func() { logger.Debug(ctx, "Debug message") }, "Expected no panic on Debug")
	assert.NotPanics(t, func() { logger.Trace(ctx, "Trace message") }, "Expected no panic on Trace")
	assert.NotPanics(t, func() { logger.Notice(ctx, "Notice message") }, "Expected no panic on Notice")
	assert.NotPanics(t, func() { logger.Warning(ctx, "Warning message") }, "Expected no panic on Warning")
	assert.NotPanics(t, func() { logger.Error(ctx, "Error message") }, "Expected no panic on Error")
	assert.NotPanics(t, func() { logger.Fatal(ctx, "Fatal message") }, "Expected no panic on Fatal")
}
