/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\logger\otelgrpc\bridge.go
 * @Description: 为 logger 装配 OpenTelemetry + gRPC 的上下文桥
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

// Package otelgrpc 为 go-toolbox/pkg/logger 装配 OpenTelemetry 与 gRPC 的上下文桥
//
// logger 核心不依赖 OTel 与 gRPC，空导入本包即可让 trace_id 优先取自 OTel span，
// 并支持经 gRPC metadata 跨节点传播：
//
//	import _ "github.com/kamalyes/go-toolbox/pkg/logger/otelgrpc"
package otelgrpc

import (
	"context"

	"github.com/kamalyes/go-toolbox/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func init() {
	logger.RegisterContextBridge(logger.ContextBridge{
		TraceIDFromContext: traceIDFromContext,
		ValueFromIncoming:  valueFromIncoming,
		InjectOutgoing:     injectOutgoing,
	})
}

// traceIDFromContext 从 OTel span 提取 traceId
// OTel span 是 traceId 的真相源：gRPC（otelgrpc StatsHandler）与 HTTP（Tracing 中间件）均由 OTel 创建 span
// ctx 不在 span 内时返回空字符串，由调用方 fallback
func traceIDFromContext(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.HasTraceID() {
		return ""
	}
	return sc.TraceID().String()
}

// valueFromIncoming 从 gRPC incoming metadata 获取指定 key 的首个非空值（大小写不敏感）
// 使用 ValueFromIncomingContext 而非 FromIncomingContext：后者每次调用深拷贝整个 MD
// （map 结构 + 所有 key 的 value slice），是线上 heap profile 中 context 提取路径的
// 主要分配来源（压测观测到该路径持续分配）；前者仅拷贝命中
// key 的单个 value slice，map 本身零分配。查询是只读的，无需持有 MD 副本
func valueFromIncoming(ctx context.Context, key string) string {
	if vals := metadata.ValueFromIncomingContext(ctx, key); len(vals) > 0 && vals[0] != "" {
		return vals[0]
	}
	return ""
}

// injectOutgoing 将 key/value 写入 gRPC outgoing metadata
func injectOutgoing(ctx context.Context, key, value string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy() // 避免修改共享的 metadata
	}
	md.Set(key, value)
	return metadata.NewOutgoingContext(ctx, md)
}
