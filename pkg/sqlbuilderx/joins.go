/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-05 20:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-05 20:15:55
 * @FilePath: \go-toolbox\pkg\sqlbuilderx\joins.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilderx

import "fmt"

// Joins 方法用于添加连接（JOIN）到查询中
func (b *Builder) Joins(table interface{}, condition string, joinType string, params ...interface{}) *Builder {
	b.initialize() // 初始化构建器

	// 设置表名，获取是否为闭包、表名、参数和表别名
	isClosure, table, param, tableAlias := b.setTable(table)

	if isClosure == 0 {
		// 如果不是闭包，则格式化表名
		table = fmt.Sprintf("`%s`", table)

		// 如果有表别名，则添加别名
		if tableAlias != "" {
			table = fmt.Sprintf("%s as `%s`", table, tableAlias)
		}
	}

	// 将连接参数添加到构建器的参数中
	b.params["join"] = append(b.params["join"], param...)
	b.params["join"] = append(b.params["join"], params...)

	// 将连接信息添加到方法列表中
	b.methods.join = append(b.methods.join, fmt.Sprintf("%s JOIN %s %s", joinType, table, condition))
	return b
}

// LefJoin 方法用于添加左连接（LEFT JOIN）到查询中
func (b *Builder) LefJoin(table interface{}, condition string, params ...interface{}) *Builder {
	return b.Joins(table, condition, "LEFT", params...)
}

// RightJoin 方法用于添加右连接（RIGHT JOIN）到查询中
func (b *Builder) RightJoin(table interface{}, condition string, params ...interface{}) *Builder {
	return b.Joins(table, condition, "RIGHT", params...)
}

// Join 方法用于添加内连接（INNER JOIN）到查询中
func (b *Builder) Join(table interface{}, condition string, params ...interface{}) *Builder {
	return b.Joins(table, condition, "INNER", params...)
}
