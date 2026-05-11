/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:00:00
 * @FilePath: \go-toolbox\pkg\idgen\types.go
 * @Description: ID 生成器类型定义
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

// IDType ID 语义类型
// 不同语义的 ID 有不同的格式要求和长度约束
type IDType string

const (
	IDTypeTraceID       IDType = "trace_id"       // 全链路追踪 ID，长格式、时间排序、全局唯一
	IDTypeSpanID        IDType = "span_id"        // 单次操作跨度 ID，短格式、同 Trace 内唯一
	IDTypeRequestID     IDType = "request_id"     // 请求唯一标识，带计数器、可排序
	IDTypeCorrelationID IDType = "correlation_id" // 跨系统关联 ID，UUID 格式、跨服务传递
)

// IDSpec ID 规格配置
// 每种 IDType 对应不同的格式规格，生成器根据规格产出不同格式的 ID
type IDSpec struct {
	TraceLen       int  // TraceID 长度
	SpanLen        int  // SpanID 长度（截取自 TraceID 或独立生成）
	RequestCounter bool // RequestID 是否附加计数器后缀
	CorrelationFmt bool // CorrelationID 是否使用 UUID 格式（含连字符）
}

// DefaultSpec 默认规格
var DefaultSpec = IDSpec{
	TraceLen:       32,
	SpanLen:        16,
	RequestCounter: true,
	CorrelationFmt: true,
}

// SpecForGenerator 各生成器的规格映射
var SpecForGenerator = map[GeneratorType]IDSpec{
	GeneratorTypeDefault:    {TraceLen: 32, SpanLen: 16, RequestCounter: true, CorrelationFmt: true},
	GeneratorTypeUUID:       {TraceLen: 36, SpanLen: 16, RequestCounter: true, CorrelationFmt: true},
	GeneratorTypeNanoID:     {TraceLen: 21, SpanLen: 16, RequestCounter: true, CorrelationFmt: false},
	GeneratorTypeSnowflake:  {TraceLen: 0, SpanLen: 0, RequestCounter: true, CorrelationFmt: true},
	GeneratorTypeShortFlake: {TraceLen: 0, SpanLen: 0, RequestCounter: true, CorrelationFmt: true},
	GeneratorTypeShortID:    {TraceLen: 10, SpanLen: 8, RequestCounter: true, CorrelationFmt: false},
	GeneratorTypeNumeric:    {TraceLen: 8, SpanLen: 8, RequestCounter: true, CorrelationFmt: false},
	GeneratorTypeULID:       {TraceLen: 26, SpanLen: 16, RequestCounter: true, CorrelationFmt: false},
}

// IDGenerator ID生成器接口
type IDGenerator interface {
	GenerateTraceID() string
	GenerateSpanID() string
	GenerateRequestID() string
	GenerateCorrelationID() string
}

// GeneratorType 生成器类型
type GeneratorType string

const (
	GeneratorTypeDefault    GeneratorType = "default"    // 默认 Hex 生成器
	GeneratorTypeUUID       GeneratorType = "uuid"       // UUID v4
	GeneratorTypeNanoID     GeneratorType = "nanoid"     // NanoID
	GeneratorTypeSnowflake  GeneratorType = "snowflake"  // Snowflake
	GeneratorTypeShortFlake GeneratorType = "shortflake" // ShortFlake (短ID)
	GeneratorTypeShortID    GeneratorType = "shortid"    // ShortID (8~10位Base62)
	GeneratorTypeNumeric    GeneratorType = "numeric"    // Numeric (8位纯数字)
	GeneratorTypeULID       GeneratorType = "ulid"       // ULID
)

// String 转换为字符串
func (t GeneratorType) String() string {
	return string(t)
}

// Spec 获取生成器对应的 ID 规格
func (t GeneratorType) Spec() IDSpec {
	if spec, ok := SpecForGenerator[t]; ok {
		return spec
	}
	return DefaultSpec
}
