/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-05 20:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-05 20:15:55
 * @FilePath: \go-toolbox\pkg\sqlbuilderx\query.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilderx

import (
	"fmt"
	"strings"
)

// Select 方法用于指定要查询的字段
func (b *Builder) Select(args ...interface{}) *Builder {
	b.methods.field = make([]interface{}, 0) // 初始化字段列表

	if len(args) == 1 {
		fieldArr := make([]string, 0)
		// 处理输入的字段参数
		if field, ok := args[0].(string); ok {
			fieldArr = strings.Split(field, ",") // 将字符串分割为字段数组
		} else if field, ok := args[0].(Raw); ok {
			b.methods.field = append(b.methods.field, field) // 处理 Raw 类型
			return b
		} else if field, ok := args[0].([]string); ok {
			fieldArr = field // 处理字符串数组
		}

		for _, v := range fieldArr {
			b.methods.field = append(b.methods.field, v) // 添加字段到方法中
		}

	} else {
		b.methods.field = append(b.methods.field, args...) // 处理多个字段参数
	}

	return b
}

// Table 方法用于设置查询的表名
func (b *Builder) Table(table interface{}) *Builder {
	b.tmpTableClosureCount, b.tmpTable, b.params["table"], b.TableAlias = b.setTable(table) // 设置表名和别名

	return b
}

// builderHaving 方法用于构建 HAVING 子句
func (b *Builder) builderHaving(sql string) (string, []interface{}) {
	params := make([]interface{}, 0) // 初始化参数列表

	having := strings.Join(b.methods.having, "") // 将 HAVING 条件连接成字符串
	if having != "" {
		having = strings.Trim(having, " ")       // 去掉前后空格
		sql += fmt.Sprintf(" HAVING %s", having) // 添加 HAVING 子句
	}

	if havingParams, ok := b.params["having"]; ok {
		params = append(params, havingParams...) // 添加 HAVING 参数
	}

	return sql, params // 返回构建的 SQL 和参数
}

// builderOrder 方法用于构建 ORDER BY 子句
func (b *Builder) builderOrder(sql string) string {
	if len(b.methods.order) > 0 {
		sql += " ORDER BY " + strings.Join(b.methods.order, ",") // 添加 ORDER BY 条件
	}

	return sql // 返回构建的 SQL
}

// Group 方法用于指定分组字段
func (b *Builder) Group(group ...string) *Builder {
	b.methods.group = append(b.methods.group, group...) // 添加分组字段

	return b
}

// Having 方法用于添加 HAVING 条件
func (b *Builder) Having(args ...interface{}) *Builder {
	var boolean string

	if len(b.methods.having) > 0 {
		boolean = "AND" // 如果已有 HAVING 条件，使用 AND
	}

	b.conditions("having", boolean, args...) // 添加条件

	return b
}

// OrHaving 方法用于添加 OR HAVING 条件
func (b *Builder) OrHaving(args ...interface{}) *Builder {
	var boolean = ""
	if len(b.methods.having) > 0 {
		boolean = "OR" // 如果已有 HAVING 条件，使用 OR
	}

	b.conditions("having", boolean, args...) // 添加条件

	return b
}

// Order 方法用于指定排序字段
func (b *Builder) Order(args ...interface{}) *Builder {
	var (
		field string
		value string
	)
	field, ok := args[0].(string) // 获取字段名
	if !ok {
		return b // 如果字段名不正确，返回
	}

	if len(args) == 1 {
		value = "DESC" // 默认降序
	} else if value, ok = args[1].(string); !ok {
		return b // 如果排序方式不正确，返回
	}

	value = strings.ToUpper(value) // 转为大写

	b.methods.order = append(b.methods.order, b.escapeId(field)+" "+value) // 添加排序条件

	return b
}

// Limit 方法用于指定查询数量
// @Description: 指定查询数量
// @receiver b
// @param int64 offset 起始位置
// @param int64 length 查询数量
// @return *Builder
func (b *Builder) Limit(args ...int64) *Builder {
	switch len(args) {
	case 1:
		b.methods.limit = fmt.Sprintf(" LIMIT %d", args[0]) // 单个参数时指定 LIMIT
	case 2:
		b.methods.limit = fmt.Sprintf(" LIMIT %d,%d", args[0], args[1]) // 两个参数时指定 OFFSET 和 LIMIT
	}

	return b
}

// Page 方法用于指定分页
// param int64 page 页数
// param int64 listRows 每页数量
// return *Builder
func (b *Builder) Page(page int64, listRows int64) *Builder {
	b.methods.limit = fmt.Sprintf(" LIMIT %d,%d", (page-1)*listRows, listRows) // 根据页数和每页数量计算 LIMIT
	return b
}

// ToSql 方法用于生成最终的 SQL 查询语句和参数
func (b *Builder) ToSql() (string, []interface{}) {
	defer b.cleanLastSql() // 清理上一个 SQL

	params := make([]interface{}, 0) // 初始化参数列表

	fieldStr := ""
	if len(b.methods.field) == 0 {
		fieldStr = "*" // 如果没有指定字段，则查询所有字段
	} else {
		fieldStr = b.escapeId(b.methods.field) // 转义字段名
	}

	// 构建基础的 SELECT 语句
	sql := fmt.Sprintf("SELECT %s FROM %s", fieldStr, b.GetTable())

	if tableParams, ok := b.params["table"]; ok {
		params = append(params, tableParams...) // 添加表参数
	}

	if len(b.methods.join) > 0 {
		sql += " " + strings.Join(b.methods.join, " ") // 添加 JOIN 条件

		if joinParams, ok := b.params["join"]; ok {
			params = append(params, joinParams...) // 添加 JOIN 参数
		}
	}

	sql, whereParams := b.builderWhere(sql) // 添加 WHERE 条件
	params = append(params, whereParams...) // 添加 WHERE 参数

	if len(b.methods.group) > 0 {
		sql += " GROUP BY " + b.escapeId(b.methods.group) // 添加 GROUP BY 条件
	}

	sql, havingParams := b.builderHaving(sql) // 添加 HAVING 条件
	params = append(params, havingParams...)  // 添加 HAVING 参数

	sql = b.builderOrder(sql) // 添加 ORDER BY 条件

	if b.methods.limit != "" {
		sql += b.methods.limit // 添加 LIMIT 条件
	}

	return sql, params // 返回最终的 SQL 和参数
}
