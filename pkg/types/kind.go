/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 13:08:33
 * @FilePath: \go-toolbox\pkg\types\kind.go
 * @Description: 类型判断辅助函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package types

// ToFloat64OK 尝试将值转换为 float64 类型，返回转换结果和是否成功
func ToFloat64OK(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}
