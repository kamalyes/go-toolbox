/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\logger\otelgrpc\bridge_benchmark_test.go
 * @Description: metadata 回退路径基准
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package otelgrpc_test

import (
	"context"
	"io"
	"testing"

	_ "github.com/kamalyes/go-toolbox/pkg/logger/otelgrpc"

	. "github.com/kamalyes/go-toolbox/pkg/logger"
	"google.golang.org/grpc/metadata"
)

func benchMetadataCtx() context.Context {
	md := metadata.Pairs(
		"app_id", "10001",
		"namespace", "live",
		"uid", "123456",
		"authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.test-token",
		"user-agent", "grpc-go/1.66.0",
		"content-type", "application/grpc",
	)
	return metadata.NewIncomingContext(context.Background(), md)
}

// BenchmarkJSON_InfoContext_MetadataFallback 测量经 ContextBridge 从 gRPC incoming
// metadata 提取字段的直写路径
// 原先在 pkg/logger 内部用 benchLogger 构造，随 metadata 路径一并移入本模块后改用导出构建器
func BenchmarkJSON_InfoContext_MetadataFallback(b *testing.B) {
	l := NewLogger().
		WithOutput(io.Discard).
		WithFormat(FormatJSON).
		WithShowCaller(true).
		WithContextKeys("app_id", "namespace", "uid")
	ctx := benchMetadataCtx()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.InfoContext(ctx, "处理请求 method=%s path=%s", "GET", "/api/v1/users")
	}
}
