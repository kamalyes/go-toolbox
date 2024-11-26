/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-05 20:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-06 11:55:55
 * @FilePath: \go-toolbox\pkg\sqlbuilderx\utils.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilderx

import (
	"fmt"
	"reflect"
	"strings"
)

// setTable 设置表名，支持字符串和闭包
func (b *Builder) setTable(table interface{}) (tmpTableClosureCount uint8, tmpTable string, param []interface{}, tableAlias string) {
	switch table.(type) {
	case string:
		// 如果是字符串，获取别名
		tmpTable, tableAlias = b.getAlias(table.(string))
	case func(*Builder):
		// 如果是闭包，创建新的 Builder
		bw := NewBuilder("")
		bw.tmpTableClosureCount = b.tmpTableClosureCount
		bw.tmpTableClosureCount++
		tmpTableClosureCount = bw.tmpTableClosureCount
		table.(func(*Builder))(bw) // 执行闭包
		tmpTable, param = bw.ToSql()
		tableAlias = fmt.Sprintf("tmp%d", tmpTableClosureCount)
		tmpTable = fmt.Sprintf("(%s) as `tmp%d`", tmpTable, tmpTableClosureCount)
	case func() *Builder:
		// 如果是返回 Builder 的闭包
		tmpTableClosureCount = b.tmpTableClosureCount + 1
		tmpTable, param = table.(func() *Builder)().ToSql()
		tableAlias = fmt.Sprintf("tmp%d", tmpTableClosureCount)
		tmpTable = fmt.Sprintf("(%s) as `tmp%d`", tmpTable, tmpTableClosureCount)
	}

	return tmpTableClosureCount, tmpTable, param, tableAlias
}

// placeholders 生成占位符
func (b *Builder) placeholders(n int) string {
	var s strings.Builder
	for i := 0; i < n-1; i++ {
		s.WriteString("?,") // 添加占位符
	}
	if n > 0 {
		s.WriteString("?") // 最后一个占位符不加逗号
	}
	return s.String()
}

// escapeId 转义字段名
func (b *Builder) escapeId(field interface{}) (fieldStr string) {
	if field, ok := field.(Raw); ok {
		fieldStr += fmt.Sprintf("%s", field)
		return
	}

	comma := ""
	var fieldArr []string
	if field, ok := field.(string); ok {
		fieldArr = strings.Split(field, ",") // 支持多个字段
	}

	if field, ok := field.([]string); ok {
		fieldArr = field
	}

	if len(fieldArr) > 0 {
		for k, v := range fieldArr {
			if k > 0 {
				comma = ","
			}
			fieldStr += b.strEscapeId(v, comma) // 转义字段
		}
		return
	}

	if field, ok := field.([]interface{}); ok {
		for k, v := range field {
			if k > 0 {
				comma = ","
			}
			switch v.(type) {
			case string:
				fieldStr += b.strEscapeId(v.(string), comma) // 转义字符串字段
			case Raw:
				fieldStr += fmt.Sprintf("%s%s", comma, v) // 原始字段
			}
		}
		return
	}

	return
}

// getAlias 获取字段和别名
func (b *Builder) getAlias(field string) (string, string) {
	var alias string

	containsAs := strings.Contains(field, " as ")

	if containsAs || strings.Contains(field, " ") {
		var fieldArr []string
		if containsAs {
			fieldArr = strings.Split(field, " as ")
		} else {
			fieldArr = strings.Split(field, " ")
		}

		field = strings.Trim(fieldArr[0], " ")

		for i := 1; i < len(fieldArr); i++ {
			if fieldArr[i] == "" {
				continue
			}
			alias = strings.Trim(fieldArr[i], " ")
			break
		}
	}

	return field, alias
}

// strEscapeId 转义字段名和别名
func (b *Builder) strEscapeId(field string, comma string) string {
	var alias, table string

	field, alias = b.getAlias(field)

	if alias != "" {
		alias = " as `" + alias + "`"
	}

	if strings.Contains(field, ".") {
		fieldArr := strings.Split(field, ".")
		table = strings.Trim(fieldArr[0], " ")
		field = strings.Trim(fieldArr[1], " ")
	}

	if table != "" {
		table = "`" + table + "`."
	}

	leftBracketIndex := strings.Index(field, "(")
	rightBracketIndex := strings.Index(field, ")")

	if leftBracketIndex >= 0 && rightBracketIndex >= 0 {
		param := strings.Trim(field[leftBracketIndex+1:rightBracketIndex], " `")
		field = field[:leftBracketIndex]

		if param != "" && param != "*" {
			param = fmt.Sprintf("`%s`", param)
		}
		field = fmt.Sprintf("%s(%s)", field, param)
	} else {
		field = fmt.Sprintf("`%s`", field)
	}

	return fmt.Sprintf("%s%s%s%s", comma, table, field, alias)
}

// convertInterfaceSlice 将接口切片转换为接口数组
func (b *Builder) convertInterfaceSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	vLen := v.Len()
	ret := make([]interface{}, vLen)
	for i := 0; i < vLen; i++ {
		ret[i] = v.Index(i).Interface() // 将每个元素转换为接口
	}
	return ret
}

// conditions 添加条件
func (b *Builder) conditions(mode string, boolean string, args ...interface{}) *Builder {
	var conditions string

	argsLen := len(args)
	if argsLen == 1 {
		if query, ok := args[0].(func(*Builder)); ok {
			bw := NewBuilder("")
			query(bw)
			conditions = fmt.Sprintf(" %s (%s)", boolean, strings.Join(bw.methods.where, ""))
			b.params[mode] = append(b.params[mode], bw.params[mode]...) // 添加参数
		} else if condition, ok := args[0].(Raw); ok {
			conditions = fmt.Sprintf(" %s %s", boolean, condition) // 原始条件
		}
	} else if argsLen > 1 {
		field := ""
		operator := ""

		var value interface{}
		switch argsLen {
		case 2:
			field = args[0].(string)
			operator = "="
			value = args[1]
		case 3:
			field = args[0].(string)
			operator = args[1].(string)
			value = args[2]
		default:
			field = args[0].(string)
			operator = args[1].(string)
			value = args[2:]
		}

		valueKind := reflect.TypeOf(value).Kind()

		if strings.Contains(operator, "BETWEEN") {
			args = b.convertInterfaceSlice(value)

			b.params[mode] = append(b.params[mode], args[:2]...)
			conditions = fmt.Sprintf(" %s %s %s ? AND ?", boolean, b.escapeId(field), operator)
		} else {
			switch valueKind {
			case reflect.Array:
			case reflect.Slice:
				vi := b.convertInterfaceSlice(value)
				conditions = fmt.Sprintf(" %s %s %s (%s)", boolean, b.escapeId(field), operator, b.placeholders(len(vi)))
				b.params[mode] = append(b.params[mode], vi...) // 添加切片参数
			case reflect.Func:
				if query, ok := value.(func(*Builder)); ok {
					bw := NewBuilder("")
					query(bw)
					bwSql, bwParams := bw.ToSql()
					if field == "EXISTS" || field == "NOT EXISTS" {
						operator = field
						conditions = fmt.Sprintf(" %s %s (%s)", boolean, operator, bwSql)
					} else {
						conditions = fmt.Sprintf(" %s %s %s (%s)", boolean, b.escapeId(field), operator, bwSql)
					}

					b.params[mode] = append(b.params[mode], bwParams...) // 添加子查询参数
				}
			default:
				if value == "NULL" || value == "NOT NULL" {
					conditions = fmt.Sprintf(" %s %s IS %s", boolean, b.escapeId(field), value)
				} else {
					conditions = fmt.Sprintf(" %s %s %s ?", boolean, b.escapeId(field), operator)
					b.params[mode] = append(b.params[mode], value) // 添加其他参数
				}
			}
		}
	}

	switch mode {
	case "where":
		b.methods.where = append(b.methods.where, conditions) // 添加 WHERE 条件
	case "having":
		b.methods.having = append(b.methods.having, conditions) // 添加 HAVING 条件
	}

	return b
}

// cleanLastSql 清理上一个 SQL 状态
func (b *Builder) cleanLastSql() {
	b.tmpTable = ""
	b.tmpTableClosureCount = 0
	b.methods = methods{}                        // 重置方法
	b.params = make(map[string][]interface{}, 0) // 重置参数
}
