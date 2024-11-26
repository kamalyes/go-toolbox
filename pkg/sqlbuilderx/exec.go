/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-05 20:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-06 11:06:59
 * @FilePath: \go-toolbox\pkg\sqlbuilderx\exec.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilderx

import (
	"fmt"
	"strings"
)

// Delete 方法构建 DELETE SQL 语句
func (b *Builder) Delete() (string, []interface{}) {
	defer b.cleanLastSql() // 清理上一个 SQL

	params := make([]interface{}, 0) // 用于存储参数

	// 构建基本的 DELETE 语句
	sql := fmt.Sprintf("DELETE FROM %s", b.GetTable())
	sql, whereParams := b.builderWhere(sql) // 添加 WHERE 条件

	// 如果有 ORDER BY 条件，添加到 SQL 语句中
	if len(b.methods.order) > 0 {
		sql += " ORDER BY " + strings.Join(b.methods.order, ",")
	}

	// 如果有 LIMIT 条件，添加到 SQL 语句中
	if b.methods.limit != "" {
		sql += b.methods.limit
	}

	// 将 WHERE 参数添加到参数列表中
	params = append(params, whereParams...)

	return sql, params // 返回构建的 SQL 语句和参数
}

// DuplicateKey 设置重复键
func (b *Builder) DuplicateKey(duplicateKey map[string]interface{}) *Builder {
	b.methods.duplicateKey = duplicateKey
	return b
}

// Insert 方法构建 INSERT SQL 语句
func (b *Builder) Insert(args ...interface{}) (string, []interface{}) {
	return b.insertReplace("INSERT", args...)
}

// Replace 方法构建 REPLACE SQL 语句
func (b *Builder) Replace(args ...interface{}) (string, []interface{}) {
	return b.insertReplace("REPLACE", args...)
}

// insertReplace 方法用于构建 INSERT 或 REPLACE SQL 语句
func (b *Builder) insertReplace(mode string, args ...interface{}) (string, []interface{}) {
	params := make([]interface{}, 0) // 用于存储参数
	sql := ""
	defer b.cleanLastSql() // 清理上一个 SQL

	var field []string
	var values [][]interface{}
	if len(args) == 2 {
		field, ok := args[0].([]string) // 获取字段名
		if query, ok1 := args[1].(func(*Builder)); ok && ok1 {
			bw := NewBuilder("")      // 创建新的 Builder 实例
			query(bw)                 // 执行传入的查询函数
			sql, params := bw.ToSql() // 获取 SQL 和参数
			sql = fmt.Sprintf("%s INTO %s (%s) %s", mode, b.GetTable(), b.escapeId(field), sql)

			return sql, params // 返回构建的 SQL 语句和参数
		}
	}

	// 处理传入的参数
	for k, arg := range args {
		isContinue := false
		if arg, ok := arg.(map[string]interface{}); ok {
			if k == 0 {
				for f, _ := range arg {
					field = append(field, f) // 获取字段名
				}
			}

			value := make([]interface{}, 0)

			for _, v := range field {
				if val, ok := arg[v]; ok {
					value = append(value, val) // 获取字段对应的值
				} else {
					isContinue = true // 如果字段缺失，继续下一项
				}
			}

			if isContinue {
				continue // 如果字段缺失，跳过当前项
			}

			values = append(values, value) // 添加值到值列表
		}
	}

	// 构建 INSERT 或 REPLACE 语句
	sql = fmt.Sprintf("%s INTO %s (%s) VALUES", mode, b.GetTable(), b.escapeId(field))

	comma := ""
	for k, value := range values {
		if k > 0 {
			comma = ","
		}
		sql += fmt.Sprintf("%s(%s)", comma, strings.Trim(strings.Repeat("?,", len(value)), ","))
		params = append(params, value...) // 添加值到参数列表
	}

	// 处理重复键更新
	if b.methods.duplicateKey != nil {
		duplicateKey := ""
		for k, value := range b.methods.duplicateKey {
			duplicateKey += fmt.Sprintf("%s=?,", k) // 添加重复键更新字段
			params = append(params, value)          // 添加值到参数列表
		}
		duplicateKey = strings.Trim(duplicateKey, ",")                        // 去掉最后的逗号
		sql = fmt.Sprintf("%s ON DUPLICATE KEY UPDATE %s", sql, duplicateKey) // 添加 ON DUPLICATE KEY UPDATE
	}

	return sql, params // 返回构建的 SQL 语句和参数
}

// Update 方法构建 UPDATE SQL 语句
func (b *Builder) Update(data map[string]interface{}) (string, []interface{}) {
	defer b.cleanLastSql() // 清理上一个 SQL

	params := make([]interface{}, 0) // 用于存储参数
	setVal := ""
	for k, v := range data {
		setVal += b.escapeId(k) + "=?,"
		params = append(params, v) // 添加值到参数列表
	}

	setVal = strings.Trim(setVal, ",") // 去掉最后的逗号

	// 构建 UPDATE 语句
	sql := fmt.Sprintf("UPDATE %s SET %s", b.GetTable(), setVal)
	sql, whereParams := b.builderWhere(sql) // 添加 WHERE 条件
	params = append(params, whereParams...) // 添加 WHERE 参数

	return sql, params // 返回构建的 SQL 语句和参数
}
