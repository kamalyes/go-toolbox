/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:30:00
 * @FilePath: \go-toolbox\pkg\idgen\shortflake.go
 * @Description: ShortFlake 短 ID 生成器（高性能、紧凑型）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// ShortFlakeGenerator 短 Snowflake 生成器
// 生成 53 位整数（JavaScript Number 安全范围）
// 时间戳(41位) + 机器ID(6位) + 序列号(6位) = 53位
// 最大值: 9007199254740991 (约 9PB，16位数字)
type ShortFlakeGenerator struct {
	epoch   int64         // 自定义纪元（毫秒）
	nodeID  int64         // 节点ID (0-63)
	state   atomic.Uint64 // 无锁状态：高位 epoch 起毫秒数，低 6 位同毫秒序列号
	counter uint64
}

// ShortFlake 位布局
//
// 时间戳(41位) + 节点ID(6位) + 序列号(6位) = 53位，刚好落在 JavaScript 的
// Number 安全整数范围内，前端不需要用字符串承载
//
// 序列号只有 6 位，单实例吞吐上限是 6.4 万每秒，远低于 Snowflake 的 409.6 万
// 这是为了压进 53 位换来的，取舍本身合理，但容量必须让调用方看得见：
// 把常量公开出来，免得压测时才发现卡在一个来路不明的数字上
const (
	// ShortFlakeMaxNodeID 节点 ID 上限（含），6 位
	ShortFlakeMaxNodeID = 63
	// ShortFlakeSequencePerMilli 同一毫秒内可发的 ID 数，6 位
	// 超出后会等到下一毫秒，因此单实例吞吐上限约为 6.4 万每秒
	ShortFlakeSequencePerMilli = 64

	shortFlakeSequenceMask = ShortFlakeSequencePerMilli - 1
	// shortFlakeDefaultEpoch 默认纪元 2022-01-01 00:00:00 UTC（毫秒）
	shortFlakeDefaultEpoch = 1640995200000
)

// ShortFlakeConfig ShortFlake 生成器配置
type ShortFlakeConfig struct {
	Epoch  int64 // 纪元时间戳（毫秒），默认 2022-01-01 00:00:00 UTC
	NodeID int64 // 节点 ID，0 到 ShortFlakeMaxNodeID
}

// DefaultShortFlakeConfig 返回默认配置（纪元 2022-01-01，节点 0）
func DefaultShortFlakeConfig() ShortFlakeConfig {
	return ShortFlakeConfig{Epoch: shortFlakeDefaultEpoch}
}

// Validate 校验配置合法性
//
// 越界的节点 ID 必须报错，不能像原先那样用 nodeID & 0x3F 静默折回
// 折回会让配成 64 的节点和配成 0 的节点变成同一个节点，从此发出完全相同的 ID，
// 负数同理（-1 折成 63）这种错误启动时毫无迹象，只会在数据里留下重复
// 详见 SnowflakeConfig.Validate 的说明
//
// [EN] An out-of-range node ID must be rejected rather than folded back with
// nodeID & 0x3F. Folding makes a node configured as 64 and one configured as 0
// the same node, emitting identical IDs; negatives fold too. See
// SnowflakeConfig.Validate
func (c ShortFlakeConfig) Validate() error {
	if c.Epoch <= 0 {
		return fmt.Errorf("ShortFlakeConfig.Epoch must be > 0, got %d", c.Epoch)
	}
	if c.Epoch > time.Now().UnixMilli() {
		return fmt.Errorf("ShortFlakeConfig.Epoch must not be in the future, got %d", c.Epoch)
	}
	if c.NodeID < 0 || c.NodeID > ShortFlakeMaxNodeID {
		return fmt.Errorf("ShortFlakeConfig.NodeID must be in [0, %d], got %d", ShortFlakeMaxNodeID, c.NodeID)
	}
	return nil
}

// NewShortFlakeGeneratorWithConfig 按配置创建短 Snowflake 生成器
//
// 配置非法时 panic，与包内其他生成器一致
func NewShortFlakeGeneratorWithConfig(cfg ShortFlakeConfig) *ShortFlakeGenerator {
	if cfg.Epoch == 0 {
		cfg.Epoch = shortFlakeDefaultEpoch
	}
	if err := cfg.Validate(); err != nil {
		panic(err)
	}
	return &ShortFlakeGenerator{
		epoch:  cfg.Epoch,
		nodeID: cfg.NodeID,
	}
}

// NewShortFlakeGenerator 创建短 Snowflake 生成器
//
// nodeID 取值 0 到 ShortFlakeMaxNodeID，越界会 panic
// 原因见 ShortFlakeConfig.Validate
func NewShortFlakeGenerator(nodeID int64) *ShortFlakeGenerator {
	return NewShortFlakeGeneratorWithConfig(ShortFlakeConfig{
		Epoch:  shortFlakeDefaultEpoch,
		NodeID: nodeID,
	})
}

// GenerateTraceID 生成跟踪ID（hex 编码，13字符）
// 格式: hex(shortflake_id)
// 示例: "018f5a3c7d2e4"
// 与 RequestID 的区别: hex 编码更紧凑，适合 HTTP Header 传递
func (g *ShortFlakeGenerator) GenerateTraceID() string {
	id := g.generate()
	// 13 位低于 53 位 ID 的 hex 上限（最多 14 位），必须用补零版而非定宽截断
	var buf [16]byte
	return string(appendHexPadded(buf[:0], uint64(id), 13))
}

// GenerateSpanID 生成跨度ID（hex 低8字符）
// 格式: hex(shortflake_id & 0xFFFFFFFF)
// 示例: "7d2e4b10"
// 与 TraceID 的区别: 更短，仅保留低位部分
func (g *ShortFlakeGenerator) GenerateSpanID() string {
	id := g.generate()
	return hex8(uint32(id))
}

// GenerateRequestID 生成请求ID（纯数字+计数器后缀）
// 格式: shortflake数字-递增计数器
// 示例: "3425234523452-1"
// 与 TraceID 的区别: 纯数字+计数器，方便日志检索和排序
func (g *ShortFlakeGenerator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	id := g.generate()

	var sb strings.Builder
	sb.Grow(24)
	sb.WriteString(strconv.FormatInt(id, 10))
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成关联ID（UUID 风格的 hex 格式）
// 格式: 8-4-4-4-4 hex（由 shortflake ID 拆分 + 计数器填充）
// 示例: "018f5a3c-7d2e-4b10-a1c3-e5f67890"
// 与 TraceID 的区别: UUID 格式，含连字符，适合跨系统传递
func (g *ShortFlakeGenerator) GenerateCorrelationID() string {
	id1 := g.generate()
	id2 := g.generate()

	// 写进栈上数组后就地改那两位
	// 原先是两次 Sprintf 加两次字符串拼接，四个中间字符串各分配一次，
	// 全都只为被 copy 一次就丢掉
	var h1, h2 [16]byte
	putHex16(&h1, uint64(id1))
	putHex16(&h2, uint64(id2))
	h1[12] = '4'
	h2[0] = '8'

	var buf [36]byte
	copy(buf[0:8], h1[0:8])
	buf[8] = '-'
	copy(buf[9:13], h1[8:12])
	buf[13] = '-'
	copy(buf[14:18], h1[12:16])
	buf[18] = '-'
	copy(buf[19:23], h2[0:4])
	buf[23] = '-'
	copy(buf[24:36], h2[4:16])

	return string(buf[:])
}

// Generate 生成 53 位整数 ID
func (g *ShortFlakeGenerator) Generate() int64 {
	return g.generate()
}

// generate 内部生成方法（无锁 CAS 实现）
// state 布局：高位为 epoch 起的毫秒数，低 6 位为同毫秒序列号
// 同毫秒序列耗尽时 CAS 自旋等待下一毫秒，不持有任何锁，其余调用方不受阻
func (g *ShortFlakeGenerator) generate() int64 {
	for {
		old := g.state.Load()
		last := old >> 6
		now := uint64(time.Now().UnixMilli() - g.epoch)

		// 时钟回拨时钳到已记录的毫秒数上
		//
		// 不钳位的话 last != now 成立，序列号从 0 重新开始，于是把那一毫秒已经
		// 发出去的号重新发一遍——ID 生成器唯一必须保证的事就此失效
		// NTP 校正、虚拟机挂起恢复、手工改系统时间都会造成回拨
		// 注意这种重复查不出于「本轮生成的 ID 是否唯一」，重复的对象在回拨之前，
		// 所以 TestShortFlakeClockRegression 直接断言时间戳不回头
		//
		// [EN] Clamps to the recorded millisecond when the clock steps back.
		// Without it last != now holds, the sequence restarts at 0, and the
		// numbers already issued in that millisecond are issued again — the one
		// thing an ID generator must guarantee. NTP correction, VM
		// suspend/resume and manual clock changes all cause this
		if now < last {
			now = last
		}

		seq := uint64(0)
		if last == now {
			// 同一毫秒内，序列号递增
			seq = (old & shortFlakeSequenceMask) + 1
			if seq > shortFlakeSequenceMask {
				// 本毫秒的序列号用尽（ShortFlakeSequencePerMilli 个），等到下一毫秒
				// 这里不持任何锁，睡眠期间其余调用方不受影响
				// 原先用 runtime.Gosched 自旋，而 Gosched 只让给其他 goroutine、
				// 不交还操作系统，实测仍是 100% 占满一个核空转
				sleepToNextMilli()
				continue
			}
		}

		newState := (now << 6) | seq
		if g.state.CompareAndSwap(old, newState) {
			// 时间戳(41位) + 节点ID(6位) + 序列号(6位)
			return (int64(now) << 12) | (g.nodeID << 6) | int64(seq)
		}
		// CAS 失败：并发竞争，重试
	}
}

// ShortFlakeBase62Generator Base62 编码的短 ID 生成器
// 将 53 位数字编码为 9-10 字符的字符串
type ShortFlakeBase62Generator struct {
	*ShortFlakeGenerator
}

// NewShortFlakeBase62Generator 创建 Base62 编码生成器
func NewShortFlakeBase62Generator(nodeID int64) *ShortFlakeBase62Generator {
	return &ShortFlakeBase62Generator{
		ShortFlakeGenerator: NewShortFlakeGenerator(nodeID),
	}
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateTraceID 生成 Base62 编码的 TraceID（9-10字符）
// 格式: base62(shortflake_id)
// 示例: "aB3xK9mPqR"
// 与 RequestID 的区别: Base62 编码更短，URL 安全
func (g *ShortFlakeBase62Generator) GenerateTraceID() string {
	return g.encodeBase62(g.generate())
}

// GenerateSpanID 生成 Base62 编码的 SpanID（6字符）
// 格式: base62(shortflake_id & 0x3FFFFFFFFF)，取低36位
// 示例: "K9mPxR"
// 与 TraceID 的区别: 更短，仅保留低位部分
func (g *ShortFlakeBase62Generator) GenerateSpanID() string {
	id := g.generate()
	return g.encodeBase62(id & 0x3FFFFFFFFF)
}

// GenerateRequestID 生成 Base62 编码的 RequestID（前缀+计数器）
// 格式: base62前6字符-递增计数器
// 示例: "aB3xK9-1"
// 与 TraceID 的区别: 带计数器后缀，可排序
func (g *ShortFlakeBase62Generator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	id := g.generate()
	prefix := g.encodeBase62(id)

	var sb strings.Builder
	sb.Grow(16)
	if len(prefix) >= 6 {
		sb.WriteString(prefix[:6])
	} else {
		sb.WriteString(prefix)
	}
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成 Base62 编码的 CorrelationID（双段拼接）
// 格式: base62(shortflake_id)-base62(shortflake_id)，18-22字符
// 示例: "aB3xK9mPqR-xY7wN4qL2v"
// 与 TraceID 的区别: 双段拼接，增强唯一性，适合跨系统关联
func (g *ShortFlakeBase62Generator) GenerateCorrelationID() string {
	id1 := g.generate()
	id2 := g.generate()

	var sb strings.Builder
	sb.Grow(22)
	sb.WriteString(g.encodeBase62(id1))
	sb.WriteByte('-')
	sb.WriteString(g.encodeBase62(id2))

	return sb.String()
}

// encodeBase62 将整数编码为 Base62 字符串（9-10字符）
func (g *ShortFlakeBase62Generator) encodeBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	var result [11]byte
	pos := len(result) - 1

	for num > 0 {
		result[pos] = base62Chars[num%62]
		num /= 62
		pos--
	}

	return string(result[pos+1:])
}
