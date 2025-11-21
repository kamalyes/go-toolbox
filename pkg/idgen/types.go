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
	GeneratorTypeULID       GeneratorType = "ulid"       // ULID
)

// String 转换为字符串
func (t GeneratorType) String() string {
	return string(t)
}
