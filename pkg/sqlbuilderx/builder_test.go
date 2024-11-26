/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-05 20:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-06 09:59:55
 * @FilePath: \go-toolbox\pkg\sqlbuilderx\builder_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilderx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	// 测试创建 Builder 实例
	builder := NewBuilder("users")

	// 验证初始化的属性
	assert.Equal(t, "users", builder.TableName, "Table name should be 'users'")
	assert.Empty(t, builder.tmpTable, "Temporary table should be empty")
	assert.Equal(t, uint8(0), builder.tmpTableClosureCount, "Temporary table closure count should be 0")
	assert.NotNil(t, builder.params, "Params should be initialized")
}

func TestBuilder_SetAndGetMethods(t *testing.T) {
	// 创建一个新的 Builder 实例
	builder := NewBuilder("users")

	// 设置字段
	builder.methods.field = []interface{}{"id", "name"}
	assert.Equal(t, []interface{}{"id", "name"}, builder.GetField(), "Should return the correct field list")

	// 设置 where 条件
	builder.methods.where = []string{"age > 30"}
	assert.Equal(t, []string{"age > 30"}, builder.GetWhere(), "Should return the correct where conditions")

	// 设置排序条件
	builder.methods.order = []string{"name ASC"}
	assert.Equal(t, []string{"name ASC"}, builder.GetOrder(), "Should return the correct order conditions")

	// 设置限制条件
	builder.methods.limit = "10"
	assert.Equal(t, "10", builder.GetLimit(), "Should return the correct limit")

	// 设置分组条件
	builder.methods.group = []string{"country"}
	assert.Equal(t, []string{"country"}, builder.GetGroup(), "Should return the correct group conditions")

	// 设置 having 条件
	builder.methods.having = []string{"COUNT(*) > 1"}
	assert.Equal(t, []string{"COUNT(*) > 1"}, builder.GetHaving(), "Should return the correct having conditions")

	// 设置 join 条件
	builder.methods.join = []string{"JOIN orders ON users.id = orders.user_id"}
	assert.Equal(t, []string{"JOIN orders ON users.id = orders.user_id"}, builder.GetJoin(), "Should return the correct join conditions")
}

func TestBuilder_Clone(t *testing.T) {
	// 创建一个新的 Builder 实例并设置一些值
	builder := NewBuilder("users")
	builder.methods.field = []interface{}{"id", "name"}
	builder.methods.where = []string{"age > 30"}
	builder.methods.order = []string{"name ASC"}
	builder.methods.limit = "10"
	builder.methods.group = []string{"country"}
	builder.methods.having = []string{"COUNT(*) > 1"}
	builder.methods.join = []string{"JOIN orders ON users.id = orders.user_id"}
	builder.params = map[string][]interface{}{
		"age": {30},
	}

	// 克隆 Builder
	clonedBuilder := builder.Clone()

	// 验证克隆后的 Builder 是否与原 Builder 相同
	assert.Equal(t, builder.TableName, clonedBuilder.TableName, "Table names should match")
	assert.Equal(t, builder.methods, clonedBuilder.methods, "Methods should match")

	// 验证克隆后的 params 是否与原 Builder 相同
	assert.Equal(t, builder.params, clonedBuilder.params, "Params should match")

	// 验证克隆后的 Builder 是否是独立的
	clonedBuilder.methods.field[0] = "changed"
	assert.NotEqual(t, builder.methods.field[0], clonedBuilder.methods.field[0], "Cloned builder should be independent")

	clonedBuilder.params["age"][0] = 40
	assert.NotEqual(t, builder.params["age"][0], clonedBuilder.params["age"][0], "Cloned builder params should be independent")
}

func TestBuilder_EmptyTableName(t *testing.T) {
	// 创建一个新的 Builder 实例，表名为空
	builder := NewBuilder("")

	// 验证获取表名
	assert.Empty(t, builder.GetTable(), "Table name should be empty")
}

func TestBuilder_ComplexClone(t *testing.T) {
	// 创建一个新的 Builder 实例并设置复杂的值
	builder := NewBuilder("users")
	builder.methods.duplicateKey = map[string]interface{}{"unique_key": "value"}
	builder.params = map[string][]interface{}{
		"age": {30, 40},
	}

	// 克隆 Builder
	clonedBuilder := builder.Clone()

	// 验证克隆后的重复键
	assert.Equal(t, builder.methods.duplicateKey, clonedBuilder.methods.duplicateKey, "Duplicate keys should match")
	assert.Equal(t, builder.params, clonedBuilder.params, "Params should match")

	// 验证独立性
	clonedBuilder.methods.duplicateKey["unique_key"] = "changed"
	assert.NotEqual(t, builder.methods.duplicateKey["unique_key"], clonedBuilder.methods.duplicateKey["unique_key"], "Cloned builder should be independent")
}
