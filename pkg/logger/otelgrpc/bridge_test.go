/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\logger\otelgrpc\bridge_test.go
 * @Description: 装配 OTel + gRPC 桥后的 logger 上下文行为
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package otelgrpc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	_ "github.com/kamalyes/go-toolbox/pkg/logger/otelgrpc"

	. "github.com/kamalyes/go-toolbox/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// TestInfoContextMetadataFallback 验证 InfoContext（Format 路径）从 gRPC metadata 提取字段
// 覆盖 valueFromIncomingMetadata 直写路径：字段存在、缺失不写、大小写不敏感、JSON 结构合法
func TestInfoContextMetadataFallback(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger().
		WithOutput(&buf).
		WithFormat(FormatJSON).
		WithShowCaller(false).
		WithContextKeys("app_id", "namespace", "uid")

	md := metadata.Pairs("APP_ID", "10001", "namespace", "live")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	logger.InfoContext(ctx, "metadata fallback")

	out := strings.TrimSpace(buf.String())
	assert.Contains(t, out, `"app_id":"10001"`, "metadata 大小写不敏感：APP_ID 应命中 app_id")
	assert.Contains(t, out, `"namespace":"live"`)
	assert.NotContains(t, out, `"uid"`, "metadata 缺失的 key 不应输出字段")
	assert.NotContains(t, out, `"trace_id"`, "无 OTel/ctx.Value/metadata 来源时 trace_id 不输出")

	var entry map[string]any
	assert.NoError(t, json.Unmarshal([]byte(out), &entry), "输出必须是合法 JSON")
	assert.Equal(t, "10001", entry["app_id"])
	assert.Equal(t, "live", entry["namespace"])
}

// TestExtractTraceFromIncomingPriority 验证 trace 提取优先级：OTel span > metadata > ctx.Value
// （ExtractTraceFromIncoming 已改为 ValueFromIncomingContext 单 value 拷贝，语义需保持等价）
func TestExtractTraceFromIncomingPriority(t *testing.T) {
	// metadata fallback + 大小写不敏感
	md := metadata.Pairs("TRACE_ID", "md-trace")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	assert.Equal(t, "md-trace", ExtractTraceFromIncoming(ctx))

	// OTel span 优先于 metadata
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01},
		SpanID:     trace.SpanID{0x02},
		TraceFlags: trace.FlagsSampled,
	})
	otelCtx := trace.ContextWithSpanContext(ctx, sc)
	assert.Equal(t, sc.TraceID().String(), ExtractTraceFromIncoming(otelCtx))

	// OTel span 优先于 ctx.Value
	valueCtx := context.WithValue(context.Background(), ContextKeyTraceID, "value-trace")
	assert.Equal(t, sc.TraceID().String(), ExtractTraceFromIncoming(trace.ContextWithSpanContext(valueCtx, sc)))

	// 空输入返回空
	assert.Empty(t, ExtractTraceFromIncoming(context.Background()))
}

// TestInfoContextOTelTraceIDJSON 验证 OTel span 的 trace_id 经直写路径输出且 JSON 合法
func TestInfoContextOTelTraceIDJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger().
		WithOutput(&buf).
		WithFormat(FormatJSON).
		WithShowCaller(false)

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x0a, 0x0b},
		SpanID:     trace.SpanID{0x0c},
		TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)
	logger.InfoContext(ctx, "otel trace")

	out := strings.TrimSpace(buf.String())
	assert.Equal(t, 1, strings.Count(out, `"trace_id"`))

	var entry map[string]any
	assert.NoError(t, json.Unmarshal([]byte(out), &entry))
	assert.Equal(t, sc.TraceID().String(), entry["trace_id"])
}
