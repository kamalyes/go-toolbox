/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 18:40:20
 * @FilePath: \go-toolbox\tests\osx_console_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"context"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/go-toolbox/pkg/random"
)

func TestOsxConsole(t *testing.T) {
	console := osx.NewColorConsole(true)
	console.SetLogLevel(osx.INFO)
	console.Success("Hello Success")
	console.Info("Hello Info")
	console.Debug("Hello Debug 1")
	console.SetLogLevel(osx.DEBUG)
	console.Debug("Hello Debug 2")
	console.SetLogLevel(osx.INFO)
	console.Debug("Hello Debug 3")
	console.Warning("Hello Warning")
	console.Error("Hello Error")
	console.SetLogFormat(osx.AdditionalCallerLogFormat)
	console.Success("Hello AdditionalCallerLogFormat Success")
	console.SetLogFormat(osx.AdditionalTimeLogFormat)
	console.ConvertJsonFormat(true)
	console.Success("Hello AdditionalTimeLogFormat Success")
	// 将请求ID存储到上下文中
	ctx := context.Background()
	ctx = context.WithValue(ctx, osx.RequestIDKey, random.RandString(5, random.CAPITAL))

	// 使用上下文打印日志
	console.LogWithContext(ctx, osx.INFO, "This is an info log with context")
}
