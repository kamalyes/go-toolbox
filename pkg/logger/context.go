/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-01 10:38:52
 * @FilePath: \go-toolbox\pkg\logger\context.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
)

// ContextKeyTraceID 是 traceId 在日志输出的字段名，同时作为 ctx.Value 的 fallback key
// traceId 的首选真相源是追踪 span，由 ContextBridge 提供（见 pkg/logger/otelgrpc 子模块）
// extract 优先经 ContextBridge 取 span traceId；取不到时 fallback 到 ctx.Value("trace_id")，
// 再 fallback 到 RPC incoming 元数据
const ContextKeyTraceID = "trace_id"

type compiledContextKey struct {
	key      string
	keyBytes []byte
	// keyAny 预装箱的 any 形态（ctx.Value 取参）
	//
	// ctx.Value(any) 对非常量 string 实参必须装箱，且因 Value 是接口方法该装箱逃逸到堆——
	// 即每条带 context 的日志、每个 context key 各一次分配。这里在编译期一次性装箱复用，
	// 把该路径的分配降为零（interface 相等性对相等 string 不变，ctx.Value 语义完全一致）
	keyAny any
}

var defaultContextKeys = []string{
	ContextKeyTraceID,
}

var defaultCompiledContextKeys = compileContextKeys(defaultContextKeys)

// ContextBridge 把 logger 与具体的分布式追踪实现、RPC 传输层解耦
// 核心模块因此不依赖 OTel 与 gRPC；需要这两者时空导入 pkg/logger/otelgrpc 子模块即可装配
// 未装配时 trace_id 退化为仅从 ctx.Value(ContextKeyTraceID) 读取，日志照常输出
type ContextBridge struct {
	// TraceIDFromContext 从追踪 span 提取 traceId，ctx 不在 span 内时返回空串
	TraceIDFromContext func(ctx context.Context) string
	// ValueFromIncoming 读取 RPC incoming 元数据中指定 key 的首个非空值（大小写不敏感）
	// 实现须避免深拷贝整个元数据表：该路径是线上 heap profile 中 context 提取的主要分配来源
	ValueFromIncoming func(ctx context.Context, key string) string
	// InjectOutgoing 将 key/value 写入 RPC outgoing 元数据并返回新 ctx
	InjectOutgoing func(ctx context.Context, key, value string) context.Context
}

// contextBridge 持有已装配的桥
// 用原子指针而非逐字段 var：读路径只需一次原子载入，零分配，且与 init 期装配无数据竞争
var contextBridge atomic.Pointer[ContextBridge]

// RegisterContextBridge 装配 ContextBridge，应在 init 阶段调用（空导入子模块即可）
// 重复调用以最后一次为准
func RegisterContextBridge(b ContextBridge) {
	contextBridge.Store(&b)
}

// extractSpanTraceID 从已装配的 ContextBridge 提取 span traceId
// span 是 traceId 的首选真相源；未装配桥、或 ctx 不在 span 内时返回空字符串，
// 调用方 fallback 到 ctx.Value / RPC 元数据
func extractSpanTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	b := contextBridge.Load()
	if b == nil || b.TraceIDFromContext == nil {
		return ""
	}
	return b.TraceIDFromContext(ctx)
}

// DefaultContextKeys 返回默认上下文提取 key
func DefaultContextKeys() []compiledContextKey {
	return compileContextKeys(defaultContextKeys)
}

func compileContextKeys(keys []string) []compiledContextKey {
	if len(keys) == 0 {
		return nil
	}

	compiled := make([]compiledContextKey, 0, len(keys))
	for _, key := range keys {
		compiled = append(compiled, compiledContextKey{
			key:      key,
			keyBytes: []byte(key),
			keyAny:   key, // 预装箱，避免每次 ctx.Value 调用各一次堆分配
		})
	}

	return compiled
}

// valueFromIncomingMetadata 从 RPC incoming 元数据获取指定 key 的首个非空值（大小写不敏感）
// 实际读取委托给已装配的 ContextBridge；未装配时返回空字符串
func valueFromIncomingMetadata(ctx context.Context, key string) string {
	b := contextBridge.Load()
	if b == nil || b.ValueFromIncoming == nil {
		return ""
	}
	return b.ValueFromIncoming(ctx, key)
}

// extractKeyValue 从 context 提取单个 key 的值
// 优先级：OTel traceId（仅 trace_id key）> ctx.Value > gRPC metadata
func extractKeyValue(ctx context.Context, key compiledContextKey, spanTraceID string) string {
	if key.key == ContextKeyTraceID && spanTraceID != "" {
		return spanTraceID
	}
	if raw := ctx.Value(key.keyAny); raw != nil {
		if text, ok := raw.(string); ok && text != "" {
			return text
		}
	}
	return valueFromIncomingMetadata(ctx, key.key)
}

func extractContextWithCompiledKeys(ctx context.Context, keys []compiledContextKey) string {
	if ctx == nil {
		return ""
	}

	p := contextPool.Get().(*[]byte)
	buf := (*p)[:0]
	defer func() { putPooledBuf(&contextPool, p, buf) }()

	buf = append(buf, '[')

	var (
		wroteField bool
	)

	// traceId 始终从 OTel span 提取（单一真相源），独立于 keys 配置
	// 即使 keys 为空或不含 trace_id，也输出 traceId，保证全链路日志打通
	spanTraceID := extractSpanTraceID(ctx)
	traceIDWritten := false
	if spanTraceID != "" {
		buf = append(buf, ContextKeyTraceID...)
		buf = append(buf, '=')
		buf = append(buf, spanTraceID...)
		wroteField = true
		traceIDWritten = true
	}

	for _, key := range keys {
		// trace_id 已由 OTel 写入，跳过
		if key.key == ContextKeyTraceID && traceIDWritten {
			continue
		}
		value := extractKeyValue(ctx, key, spanTraceID)
		if value == "" {
			continue
		}
		if wroteField {
			buf = append(buf, ' ')
		}
		buf = append(buf, key.keyBytes...)
		buf = append(buf, '=')
		buf = append(buf, value...)
		wroteField = true
	}

	if !wroteField {
		return ""
	}

	buf = append(buf, ']', ' ')
	return string(buf)
}

// extractContextFieldsWithCompiledKeys 从 context 中提取键值对，返回 map 形式
// 用于 JSON 格式输出，使 traceId 等信息成为 JSON 字段而非 prepend 文本
func extractContextFieldsWithCompiledKeys(ctx context.Context, keys []compiledContextKey) map[string]any {
	if ctx == nil {
		return nil
	}

	fields := make(map[string]any, len(keys)+1)

	// traceId 始终从 OTel span 提取（单一真相源），独立于 keys 配置
	// 即使 keys 为空或不含 trace_id，也输出 traceId，保证全链路日志打通
	spanTraceID := extractSpanTraceID(ctx)
	traceIDWritten := false
	if spanTraceID != "" {
		fields[ContextKeyTraceID] = spanTraceID
		traceIDWritten = true
	}

	for _, key := range keys {
		// trace_id 已由 OTel 写入，跳过
		if key.key == ContextKeyTraceID && traceIDWritten {
			continue
		}
		if value := extractKeyValue(ctx, key, spanTraceID); value != "" {
			fields[key.key] = value
		}
	}

	if len(fields) == 0 {
		return nil
	}
	return fields
}

// appendContextFieldsJSON 将 ctx 字段直写为 JSON 顶层字段（零 map 分配）
// 与 extractContextFields 语义一致（contextExtractor 降级 / OTel traceId 单一真相源），
// 区别在于：跳过已由 KV 或预编码字段（pre）提供的同名 key——两者优先级都高于 ctx，与 map 合并的覆盖语义一致，由调用方在 ctx 之后写入各自的值
func (l *Logger) appendContextFieldsJSON(buf []byte, ctx context.Context, kv []any, pre []logField) []byte {
	if ctx == nil {
		return buf
	}
	if l.contextExtractor != nil {
		info := l.contextExtractor(ctx)
		if info = strings.TrimSpace(info); info != "" {
			if !overridesContextKey(kv, pre, "context") {
				buf = append(buf, ',')
				buf = appendJSONKey(buf, "context")
				buf = appendJSONString(buf, info)
			}
		}
		return buf
	}

	// traceId 始终从 OTel span 提取（单一真相源），独立于 keys 配置
	// 即使 keys 为空或不含 trace_id，也输出 traceId，保证全链路日志打通
	spanTraceID := extractSpanTraceID(ctx)
	traceIDHandled := false
	if spanTraceID != "" {
		if !overridesContextKey(kv, pre, ContextKeyTraceID) {
			buf = append(buf, ',')
			buf = appendJSONKey(buf, ContextKeyTraceID)
			buf = appendJSONString(buf, spanTraceID)
		}
		// 无论是否写出（KV/预编码字段提供同名时跳过），trace_id 均视为已处理，
		// 避免 keys 循环重复写
		traceIDHandled = true
	}

	for _, key := range l.contextKeys {
		// trace_id 已由 OTel 写入或由 KV/预编码字段提供，跳过
		if key.key == ContextKeyTraceID && traceIDHandled {
			continue
		}
		value := extractKeyValue(ctx, key, spanTraceID)
		if value == "" || overridesContextKey(kv, pre, key.key) {
			continue
		}
		buf = append(buf, ',')
		buf = append(buf, '"')
		buf = append(buf, key.keyBytes...)
		buf = append(buf, '"', ':')
		buf = appendJSONString(buf, value)
	}

	return buf
}

// overridesContextKey 报告 KV 参数或预编码字段是否已提供指定 key
func overridesContextKey(kv []any, pre []logField, key string) bool {
	return kvContainsKey(kv, key) || fieldsContainKey(pre, key)
}

// kvContainsKey 报告键值对参数中是否已提供指定 key（KV 优先判重用）
// string key 直接比较零分配；非 string key 经 fmt.Sprint 转换后比较（罕见路径允许分配）
func kvContainsKey(kv []any, key string) bool {
	for i := 0; i < len(kv); i += 2 {
		if k, ok := kv[i].(string); ok {
			if k == key {
				return true
			}
		} else if fmt.Sprint(kv[i]) == key {
			return true
		}
	}
	return false
}

// WithContextKeys 配置 Logger 在记录 Context 日志时提取哪些 key
func (l *Logger) WithContextKeys(keys ...string) *Logger {
	l.contextKeys = compileContextKeys(keys)
	l.contextExtractor = nil
	return l
}

// ExtractTraceID 从 ctx 提取 trace_id
// 优先级：OTel span > ctx.Value(ContextKeyTraceID) > 空
// 直接复用 context.go 的 extractSpanTraceID，保证与日志输出使用同一真相源
func ExtractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	// 1. 追踪 span（首选真相源）- 复用 context.go 内部函数
	if id := extractSpanTraceID(ctx); id != "" {
		return id
	}
	// 2. ctx.Value fallback（ContextKeyTraceID = "trace_id"）
	if v := ctx.Value(ContextKeyTraceID); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// ContextWithTraceID 创建携带 trace_id 的 context
// 如果 ctx 已有 trace_id 则不覆盖（尊重上游传入的值）
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}
	if ExtractTraceID(ctx) != "" {
		return ctx // 已有则不覆盖
	}
	return context.WithValue(ctx, ContextKeyTraceID, traceID)
}

// InjectTraceToOutgoing 将 trace_id 注入 RPC outgoing 元数据
// 用于客户端调用远程节点时传播 trace，使追踪未启用时仍可传播
// 实际写入委托给已装配的 ContextBridge；未装配时原样返回 ctx
func InjectTraceToOutgoing(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}
	b := contextBridge.Load()
	if b == nil || b.InjectOutgoing == nil {
		return ctx
	}
	return b.InjectOutgoing(ctx, ContextKeyTraceID, traceID)
}

// ExtractTraceFromIncoming 从 gRPC incoming metadata 提取 trace_id
// 优先级：OTel span > metadata > ctx.Value
// 用于 gRPC 服务端接收到远程节点请求时恢复 trace
func ExtractTraceFromIncoming(ctx context.Context) string {
	// 1. 追踪 span - 复用 context.go 内部函数
	if id := extractSpanTraceID(ctx); id != "" {
		return id
	}
	// 2. RPC 元数据 fallback（经 ContextBridge，实现侧避免全量元数据深拷贝）
	if id := valueFromIncomingMetadata(ctx, ContextKeyTraceID); id != "" {
		return id
	}
	// 3. ctx.Value fallback
	if v := ctx.Value(ContextKeyTraceID); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// RestoreTraceFromIncoming 从 gRPC incoming 恢复 trace_id 到 ctx
// 如果提取到 trace_id 但 ctx 中没有，则注入到 ctx.Value
// 返回增强后的 ctx，供下游 logger 自动输出 trace_id
func RestoreTraceFromIncoming(ctx context.Context) context.Context {
	traceID := ExtractTraceFromIncoming(ctx)
	if traceID == "" {
		return ctx
	}
	// 只在 ctx 中缺失时注入，不覆盖已有值
	return ContextWithTraceID(ctx, traceID)
}
