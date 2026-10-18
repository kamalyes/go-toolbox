/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:00:00
 * @FilePath: \go-toolbox\pkg\idgen\snowflake.go
 * @Description: Snowflake 分布式ID生成器（高性能版本）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// SnowflakeGenerator Snowflake 分布式ID生成器
type SnowflakeGenerator struct {
	epoch      int64
	workerID   int64
	datacenter int64
	sequence   uint64
	lastTime   int64
	counter    uint64
	mu         sync.Mutex
}

// Snowflake 位布局
//
// 时间戳(41位) + 数据中心(5位) + 工作机器(5位) + 序列号(12位) = 63位
// 这些值决定了容量上限，公开出来便于调用方在部署前核算，而不是读源码里的掩码
const (
	// SnowflakeMaxWorkerID 工作机器 ID 上限（含），5 位
	SnowflakeMaxWorkerID = 31
	// SnowflakeMaxDatacenterID 数据中心 ID 上限（含），5 位
	SnowflakeMaxDatacenterID = 31
	// SnowflakeSequencePerMilli 同一毫秒内可发的 ID 数，12 位
	// 超出后会等到下一毫秒，因此单实例吞吐上限约为 409.6 万每秒
	SnowflakeSequencePerMilli = 4096

	snowflakeSequenceMask = SnowflakeSequencePerMilli - 1
	// snowflakeDefaultEpoch 默认纪元 2022-01-01 00:00:00 UTC（毫秒）
	// 41 位毫秒数可用约 69.7 年，即到 2091 年左右
	snowflakeDefaultEpoch = 1640995200000
)

// SnowflakeConfig Snowflake 生成器配置
type SnowflakeConfig struct {
	Epoch        int64 // 纪元时间戳（毫秒），默认 2022-01-01 00:00:00 UTC
	WorkerID     int64 // 工作机器 ID，0 到 SnowflakeMaxWorkerID
	DatacenterID int64 // 数据中心 ID，0 到 SnowflakeMaxDatacenterID
}

// DefaultSnowflakeConfig 返回默认配置（纪元 2022-01-01，机器与数据中心均为 0）
func DefaultSnowflakeConfig() SnowflakeConfig {
	return SnowflakeConfig{Epoch: snowflakeDefaultEpoch}
}

// Validate 校验配置合法性
//
// 越界的 ID 必须报错，不能像原先那样用 workerID & 0x1F 静默折回
//
// 静默折回的后果是两台机器拿到同一个 ID：配成 32 的机器和配成 0 的机器会被折成
// 同一个 worker，从此发出完全相同的 ID而这恰好是分布式 ID 生成器唯一要保证的事负数同理，-1 会被折成 31
// 这种错误在启动时没有任何迹象，等到数据里出现重复主键才被发现，溯源极难
// 启动即失败远好过静默发重号
//
// [EN] Out-of-range IDs must be rejected rather than silently folded back with
// workerID & 0x1F as before. Folding makes two machines share one ID: a machine
// configured as 32 and one configured as 0 become the same worker and emit
// identical IDs from then on — the one thing a distributed ID generator exists
// to guarantee. Negatives fold too, with -1 becoming 31. Nothing shows at
// startup; it surfaces as duplicate primary keys much later, and tracing it back
// is very hard. Failing at startup is far better than silently issuing duplicates
func (c SnowflakeConfig) Validate() error {
	if c.Epoch <= 0 {
		return fmt.Errorf("SnowflakeConfig.Epoch must be > 0, got %d", c.Epoch)
	}
	if c.Epoch > time.Now().UnixMilli() {
		return fmt.Errorf("SnowflakeConfig.Epoch must not be in the future, got %d", c.Epoch)
	}
	if c.WorkerID < 0 || c.WorkerID > SnowflakeMaxWorkerID {
		return fmt.Errorf("SnowflakeConfig.WorkerID must be in [0, %d], got %d", SnowflakeMaxWorkerID, c.WorkerID)
	}
	if c.DatacenterID < 0 || c.DatacenterID > SnowflakeMaxDatacenterID {
		return fmt.Errorf("SnowflakeConfig.DatacenterID must be in [0, %d], got %d", SnowflakeMaxDatacenterID, c.DatacenterID)
	}
	return nil
}

// NewSnowflakeGeneratorWithConfig 按配置创建 Snowflake 生成器
//
// 配置非法时 panic，与 NewNumericIDGeneratorWithConfig 保持一致：ID 生成器的
// 配置属于启动期参数，带着错误配置继续跑只会污染数据
func NewSnowflakeGeneratorWithConfig(cfg SnowflakeConfig) *SnowflakeGenerator {
	if cfg.Epoch == 0 {
		cfg.Epoch = snowflakeDefaultEpoch
	}
	if err := cfg.Validate(); err != nil {
		panic(err)
	}
	return &SnowflakeGenerator{
		epoch:      cfg.Epoch,
		workerID:   cfg.WorkerID,
		datacenter: cfg.DatacenterID,
		lastTime:   0,
		sequence:   0,
	}
}

// NewSnowflakeGenerator 创建 Snowflake 生成器
//
// workerID 与 datacenter 的取值范围见 SnowflakeMaxWorkerID / SnowflakeMaxDatacenterID
// 越界会 panic，原因见 SnowflakeConfig.Validate
func NewSnowflakeGenerator(workerID, datacenter int64) *SnowflakeGenerator {
	return NewSnowflakeGeneratorWithConfig(SnowflakeConfig{
		Epoch:        snowflakeDefaultEpoch,
		WorkerID:     workerID,
		DatacenterID: datacenter,
	})
}

// GenerateTraceID 生成跟踪ID（16进制格式，带时间排序特征）
// 格式: hex(snowflake_id)，16字符
// 示例: "018f5a3c7d2e4b10"
// 与 RequestID 的区别: 使用 hex 编码，更适合作为追踪标识
func (g *SnowflakeGenerator) GenerateTraceID() string {
	id := g.Generate()
	return hex16(uint64(id))
}

// GenerateSpanID 生成跨度ID（hex截取后8字符）
// 格式: hex(snowflake_id) 后8字符
// 示例: "7d2e4b10"
// 与 TraceID 的区别: 更短，仅保留低位部分，同 Trace 内唯一
func (g *SnowflakeGenerator) GenerateSpanID() string {
	id := g.Generate()
	return hex8(uint32(id))
}

// GenerateRequestID 生成请求ID（纯数字+计数器后缀）
// 格式: snowflake数字-递增计数器
// 示例: "1732184000123456789-1"
// 与 TraceID 的区别: 纯数字+计数器，可排序，方便日志检索
func (g *SnowflakeGenerator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	id := g.Generate()

	var sb strings.Builder
	sb.Grow(28)
	sb.WriteString(strconv.FormatInt(id, 10))
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成关联ID（UUID 风格的 hex 格式）
// 格式: 8-4-4-4-4 hex（由 snowflake ID 拆分填充）
// 示例: "018f5a3c-7d2e-4b10-a1c3-e5f678901234"
// 与 TraceID 的区别: UUID 格式，含连字符，适合跨系统传递和关联
func (g *SnowflakeGenerator) GenerateCorrelationID() string {
	id := g.Generate()

	// 直接把 hex 写进栈上数组，不再先 Sprintf 出中间字符串再 copy 出来
	// 那两个中间字符串只为被 copy 一次就丢掉，是纯粹的分配浪费
	var h [16]byte
	putHex16(&h, uint64(id))

	var buf [36]byte
	copy(buf[0:8], h[0:8])
	buf[8] = '-'
	copy(buf[9:13], h[8:12])
	buf[13] = '-'
	copy(buf[14:18], "4b10")
	buf[18] = '-'
	copy(buf[19:23], h[12:16])
	buf[23] = '-'

	counter := atomic.AddUint64(&g.counter, 1)
	// %012x 取的是 counter 的低 12 位 hex，即 16 位 hex 的后 12 位
	var c [16]byte
	putHex16(&c, counter)
	copy(buf[24:36], c[4:16])

	return string(buf[:])
}

// Generate 生成 Snowflake ID
func (g *SnowflakeGenerator) Generate() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < g.lastTime {
		now = g.lastTime
	}

	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & snowflakeSequenceMask
		if g.sequence == 0 {
			// 本毫秒的序列号用尽（SnowflakeSequencePerMilli 个），等下一毫秒
			// 等待期间继续持锁：序列耗尽时其余调用方同样无号可发，让它们排队是对的，问题只在于原先是空转轮询，现在精确睡到边界
			for now <= g.lastTime {
				sleepToNextMilli()
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	timestamp := (now - g.epoch) << 22
	worker := g.workerID << 17
	dc := g.datacenter << 12

	return timestamp | worker | dc | int64(g.sequence)
}
