/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 19:31:55
 * @FilePath: \go-toolbox\pkg\idgen\factory.go
 * @Description: ID 生成器工厂
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

// NewIDGenerator 创建 ID 生成器
func NewIDGenerator(generatorType interface{}) IDGenerator {
	var typeStr string

	switch v := generatorType.(type) {
	case GeneratorType:
		typeStr = v.String()
	case string:
		typeStr = v
	default:
		return NewDefaultIDGenerator()
	}

	switch typeStr {
	case "uuid":
		return NewUUIDGenerator()
	case "nanoid":
		return NewNanoIDGenerator()
	case "snowflake":
		return NewSnowflakeGenerator(1, 1)
	case "shortflake", "short":
		return NewShortFlakeGenerator(1)
	case "ulid":
		return NewULIDGenerator()
	case "default", "hex", "logger", "":
		return NewDefaultIDGenerator()
	default:
		return NewDefaultIDGenerator()
	}
}

// NewIDGeneratorFromString 从字符串创建 ID 生成器（已废弃，使用 NewIDGenerator 代替）
// Deprecated: 使用 NewIDGenerator 代替
func NewIDGeneratorFromString(generatorType string) IDGenerator {
	return NewIDGenerator(generatorType)
}
