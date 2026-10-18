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

import (
	"strconv"
	"time"
)

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

// ============================================================================
// 定宽 hex 格式化
//
// ID 的字符串形态几乎都是定宽小写 hex，原先一律走 fmt.Sprintf
// fmt 要做格式串解析与反射取值，对这种固定形状纯属浪费实测（Apple M4）：
//
//	%016x   Sprintf 42.4 ns   手写 12.0 ns
//	%08x    Sprintf 36.4 ns   手写  8.6 ns
//
// 两者分配次数相同（都只分配结果字符串），差的是 CPU
// ============================================================================

// hexDigits 小写 hex 字符表
const hexDigits = "0123456789abcdef"

// putHex16 把 v 写成 16 位小写 hex 填满 dst
//
// uint64 最多就是 16 位 hex，不存在溢出宽度的情况，因此结果与
// fmt.Sprintf("%016x", v) 逐位一致，定宽写入是精确的而不是近似
//
// [EN] Writes v as exactly 16 lowercase hex digits. A uint64 never needs more
// than 16, so this is bit-identical to fmt.Sprintf("%016x", v), not an
// approximation of it
func putHex16(dst *[16]byte, v uint64) {
	for i := 15; i >= 0; i-- {
		dst[i] = hexDigits[v&0xF]
		v >>= 4
	}
}

// hex16 返回 v 的 16 位小写 hex，等价于 fmt.Sprintf("%016x", v)
//
// [EN] Equivalent to fmt.Sprintf("%016x", v)
func hex16(v uint64) string {
	var b [16]byte
	putHex16(&b, v)
	return string(b[:])
}

// hex8 返回 v 的 8 位小写 hex，等价于 fmt.Sprintf("%08x", v)
//
// 入参是 uint32，最多 8 位 hex，同样是精确等价
//
// [EN] Equivalent to fmt.Sprintf("%08x", v). The uint32 input never needs more
// than 8 digits, so this too is exact
func hex8(v uint32) string {
	var b [8]byte
	for i := 7; i >= 0; i-- {
		b[i] = hexDigits[v&0xF]
		v >>= 4
	}
	return string(b[:])
}

// appendHexPadded 以小写 hex 追加 v，不足 width 位时左侧补零
//
// 精确复现 %0Nx 的语义：width 是最小宽度而非截断宽度，值需要更多位时照实输出
// 当 width 小于该类型 hex 的最大位数时必须用这个版本，不能用定宽写入
// 例如 53 位 ID 配 %013x，到约 2056 年就会需要第 14 位
//
// [EN] Appends v as lowercase hex, zero-padded on the left to width. Reproduces
// %0Nx exactly: width is a minimum, not a truncation point, so a value needing
// more digits emits them. Required wherever width is below the type's maximum
// hex length — a 53-bit ID formatted %013x needs a 14th digit around 2056
func appendHexPadded(dst []byte, v uint64, width int) []byte {
	var tmp [16]byte
	h := strconv.AppendUint(tmp[:0], v, 16)
	for i := len(h); i < width; i++ {
		dst = append(dst, '0')
	}
	return append(dst, h...)
}

// sleepToNextMilli 睡到下一毫秒边界
//
// 同毫秒序列号耗尽后必须等到下一毫秒才能继续发号原先两种等法都在空转：
// snowflake 持着互斥锁紧凑轮询 time.Now，shortflake 用 runtime.Gosched 自旋
// Gosched 只把 CPU 让给其他 goroutine，并不交还给操作系统，所以没有别的
// goroutine 可跑时核心照样被占满实测（Apple M4，单核）两者的 CPU 时间都等于
// 墙上时间，也就是一整个核在空转：
//
//	Snowflake   40000 个 ID   墙上 8.68ms   CPU 8.66ms   99.8%
//	ShortFlake     600 个 ID   墙上 8.57ms   CPU 8.58ms  100.1%
//
// 改成按剩余时间精确睡到边界：吞吐由配额决定，本来就不受等法影响，而等待期间
// CPU 可以交给别的工作
//
// [EN] Sleeps to the next millisecond boundary. Once a millisecond's sequence is
// spent, generation has to wait for the next one. Both previous waits spun:
// snowflake polled time.Now in a tight loop while holding its mutex, shortflake
// spun on runtime.Gosched — which only yields to other goroutines, never back to
// the OS, so with nothing else runnable the core still burns. Measured, both
// consumed CPU time equal to wall time, a whole core spinning. Throughput is set
// by the per-millisecond quota either way, so sleeping costs nothing there and
// leaves the CPU for other work
func sleepToNextMilli() {
	rem := time.Millisecond - time.Duration(time.Now().UnixNano()%int64(time.Millisecond))
	time.Sleep(rem)
}
