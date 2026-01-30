/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-30 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-30 00:00:00
 * @FilePath: \go-toolbox\pkg\types\convertible.go
 * @Description: 可转换类型约束定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package types

// Convertible 是一个约束，包含所有可以通过 convert.MustConvertTo 转换的类型
// 包括：基础类型（string, bool）、所有数字类型、字节切片、字典和切片
type Convertible interface {
	~string | ~bool |
		Numerical |
		~[]byte |
		~map[string]any |
		~[]any
}
