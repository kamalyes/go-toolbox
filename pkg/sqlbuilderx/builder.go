/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-05 20:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-06 10:03:59
 * @FilePath: \go-toolbox\pkg\sqlbuilderx\builder.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilderx

import (
	"fmt"
)

// Raw 定义一个类型，表示原始 SQL 字符串
type Raw string

// methods 结构体用于存储 SQL 查询的各个部分
type methods struct {
	field        []interface{}          // 字段列表
	where        []string               // where 条件
	order        []string               // 排序条件
	limit        string                 // 限制条件
	group        []string               // 分组条件
	having       []string               // having 条件
	join         []string               // join 条件
	duplicateKey map[string]interface{} // 重复键
}

// Builder 结构体用于构建 SQL 查询
type Builder struct {
	TableName            string                   // 表名
	tmpTable             string                   // 临时表名
	TableAlias           string                   // 表别名
	tmpTableClosureCount uint8                    // 临时表闭合计数
	params               map[string][]interface{} // 参数列表

	// 链式操作方法列表
	methods methods // 存储 SQL 查询的各个部分
}

// NewBuilder 创建一个新的 Builder 实例
func NewBuilder(tableName string) *Builder {
	obj := &Builder{TableName: tableName} // 初始化 Builder
	obj.initialize()                      // 调用初始化方法
	return obj
}

// GetField 获取字段列表
func (b *Builder) GetField() []interface{} {
	return b.methods.field
}

// GetWhere 获取 where 条件
func (b *Builder) GetWhere() []string {
	return b.methods.where
}

// TmpTable 获取临时表名
func (b *Builder) TmpTable() string {
	return b.tmpTable
}

// GetTable 获取当前表的名称
func (b *Builder) GetTable() string {
	if b.tmpTable != "" {
		table := b.tmpTable
		if b.tmpTableClosureCount == 0 {
			table = fmt.Sprintf("`%s`", table) // 使用反引号包裹表名

			if b.TableAlias != "" {
				table = fmt.Sprintf("%s as `%s`", table, b.TableAlias) // 添加表别名
			}
		}
		return table
	} else {
		if b.TableName == "" {
			return "" // 如果表名为空，返回空字符串
		}

		return fmt.Sprintf("`%s`", b.TableName) // 使用反引号包裹表名
	}
}

// GetOrder 获取排序条件
func (b *Builder) GetOrder() []string {
	return b.methods.order
}

// GetLimit 获取限制条件
func (b *Builder) GetLimit() string {
	return b.methods.limit
}

// GetGroup 获取分组条件
func (b *Builder) GetGroup() []string {
	return b.methods.group
}

// GetHaving 获取 having 条件
func (b *Builder) GetHaving() []string {
	return b.methods.having
}

// GetJoin 获取 join 条件
func (b *Builder) GetJoin() []string {
	return b.methods.join
}

// Clone 克隆当前 Builder 实例，返回一个新的 Builder
func (b *Builder) Clone() *Builder {

	// 克隆 params
	paramsClone := make(map[string][]interface{})
	for k, v := range b.params {
		// 深拷贝每个切片
		vClone := make([]interface{}, len(v))
		copy(vClone, v)
		paramsClone[k] = vClone
	}

	// 使用简单的切片拷贝
	methodsClone := methods{
		field:  append([]interface{}{}, b.methods.field...), // 克隆字段列表
		where:  append([]string{}, b.methods.where...),      // 克隆 where 条件
		order:  append([]string{}, b.methods.order...),      // 克隆排序条件
		limit:  b.methods.limit,                             // 直接复制限制条件
		group:  append([]string{}, b.methods.group...),      // 克隆分组条件
		having: append([]string{}, b.methods.having...),     // 克隆 having 条件
		join:   append([]string{}, b.methods.join...),       // 克隆 join 条件
	}

	// 如果 duplicateKey 被设置，则复制它
	if len(b.methods.duplicateKey) > 0 {
		methodsClone.duplicateKey = make(map[string]interface{})
		for k, v := range b.methods.duplicateKey {
			methodsClone.duplicateKey[k] = v
		}
	}

	return &Builder{
		TableName:            b.TableName,
		tmpTable:             b.tmpTable,
		tmpTableClosureCount: b.tmpTableClosureCount,
		params:               paramsClone,  // 克隆参数
		methods:              methodsClone, // 使用克隆后的 methods
	}
}
