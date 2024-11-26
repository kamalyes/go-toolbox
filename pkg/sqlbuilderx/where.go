/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-05 20:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-05 20:15:55
 * @FilePath: \go-toolbox\pkg\sqlbuilderx\where.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilderx

import (
	"fmt"
	"strings"
)

// builderWhere 方法构建 WHERE 子句并返回 SQL 语句和参数
func (b *Builder) builderWhere(sql string) (string, []interface{}) {
	params := make([]interface{}, 0)

	where := strings.Join(b.methods.where, "")
	if where != "" {
		where = strings.Trim(where, " ")
		sql += fmt.Sprintf(" WHERE %s", where) // 添加 WHERE 子句
	}

	if whereParams, ok := b.params["where"]; ok {
		params = append(params, whereParams...) // 添加 WHERE 参数
	}

	return sql, params
}

// Where 方法用于添加 WHERE 条件
func (b *Builder) Where(args ...interface{}) *Builder {
	var boolean string

	if len(b.methods.where) > 0 {
		boolean = "AND" // 如果已有条件，则使用 AND
	}

	b.conditions("where", boolean, args...) // 调用 conditions 方法

	return b
}

// OrWhere 方法用于添加 OR WHERE 条件
func (b *Builder) OrWhere(args ...interface{}) *Builder {
	var boolean string

	if len(b.methods.where) > 0 {
		boolean = "OR" // 如果已有条件，则使用 OR
	}

	b.conditions("where", boolean, args...) // 调用 conditions 方法

	return b
}

// WhereExists 方法用于添加 EXISTS 条件
func (b *Builder) WhereExists(where func(*Builder)) *Builder {
	return b.Where("EXISTS", where) // 调用 Where 方法
}

// WhereNotExists 方法用于添加 NOT EXISTS 条件
func (b *Builder) WhereNotExists(where func(*Builder)) *Builder {
	return b.Where("NOT EXISTS", where) // 调用 Where 方法
}

// OrWhereExists 方法用于添加 OR EXISTS 条件
func (b *Builder) OrWhereExists(where func(*Builder)) *Builder {
	return b.OrWhere("EXISTS", where) // 调用 OrWhere 方法
}

// OrWhereNotExists 方法用于添加 OR NOT EXISTS 条件
func (b *Builder) OrWhereNotExists(where func(*Builder)) *Builder {
	return b.OrWhere("NOT EXISTS", where) // 调用 OrWhere 方法
}

// WhereIn 方法用于添加 IN 条件
func (b *Builder) WhereIn(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "IN" // 操作符为 IN

	args = append(args, value...) // 添加值

	return b.Where(args...) // 调用 Where 方法
}

// WhereNotIn 方法用于添加 NOT IN 条件
func (b *Builder) WhereNotIn(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "NOT IN" // 操作符为 NOT IN

	args = append(args, value...) // 添加值

	return b.Where(args...) // 调用 Where 方法
}

// OrWhereIn 方法用于添加 OR IN 条件
func (b *Builder) OrWhereIn(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "IN" // 操作符为 IN

	args = append(args, value...) // 添加值

	return b.OrWhere(args...) // 调用 OrWhere 方法
}

// OrWhereNotIn 方法用于添加 OR NOT IN 条件
func (b *Builder) OrWhereNotIn(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "NOT IN" // 操作符为 NOT IN

	args = append(args, value...) // 添加值

	return b.OrWhere(args...) // 调用 OrWhere 方法
}

// WhereNull 方法用于添加 IS NULL 条件
func (b *Builder) WhereNull(field string) *Builder {
	return b.Where(field, "NULL") // 调用 Where 方法
}

// WhereNotNull 方法用于添加 IS NOT NULL 条件
func (b *Builder) WhereNotNull(field string) *Builder {
	return b.Where(field, "NOT NULL") // 调用 Where 方法
}

// OrWhereNull 方法用于添加 OR IS NULL 条件
func (b *Builder) OrWhereNull(field string) *Builder {
	return b.OrWhere(field, "NULL") // 调用 OrWhere 方法
}

// OrWhereNotNull 方法用于添加 OR IS NOT NULL 条件
func (b *Builder) OrWhereNotNull(field string) *Builder {
	return b.OrWhere(field, "NOT NULL") // 调用 OrWhere 方法
}

// WhereBetween 方法用于添加 BETWEEN 条件
func (b *Builder) WhereBetween(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "BETWEEN" // 操作符为 BETWEEN

	args = append(args, value...) // 添加值

	return b.Where(args...) // 调用 Where 方法
}

// OrWhereBetween 方法用于添加 OR BETWEEN 条件
func (b *Builder) OrWhereBetween(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "BETWEEN" // 操作符为 BETWEEN

	args = append(args, value...) // 添加值

	return b.OrWhere(args...) // 调用 OrWhere 方法
}

// WhereNotBetween 方法用于添加 NOT BETWEEN 条件
func (b *Builder) WhereNotBetween(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "NOT BETWEEN" // 操作符为 NOT BETWEEN

	args = append(args, value...) // 添加值

	return b.Where(args...) // 调用 Where 方法
}

// OrWhereNotBetween 方法用于添加 OR NOT BETWEEN 条件
func (b *Builder) OrWhereNotBetween(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "NOT BETWEEN" // 操作符为 NOT BETWEEN

	args = append(args, value...) // 添加值

	return b.OrWhere(args...) // 调用 OrWhere 方法
}
