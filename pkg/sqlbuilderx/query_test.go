/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-05 20:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-05 20:17:55
 * @FilePath: \go-toolbox\pkg\sqlbuilderx\query_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sqlbuilderx

import (
	"reflect"
	"testing"
)

func TestBuilder_Select(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name as n").ToSql()
	if sql == "SELECT `id`,`name` as `n` FROM `user`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Select("id,name n").ToSql()
	if sql == "SELECT `id`,`name` as `n` FROM `user`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Select([]string{"id", "name n"}).ToSql()
	if sql == "SELECT `id`,`name` as `n` FROM `user`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Select_Aggregate(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("max(id) as id_max").ToSql()
	if sql == "SELECT max(`id`) as `id_max` FROM `user`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Select("count(*) c").ToSql()
	if sql == "SELECT count(*) as `c` FROM `user`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Raw(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select(Raw("DISTINCT mobile")).ToSql()
	if sql == "SELECT DISTINCT mobile FROM `user`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Table(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Table("m_users").ToSql()
	if sql == "SELECT * FROM `m_users`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Table_Alias(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Table("m_users u").ToSql()
	if sql == "SELECT * FROM `m_users` as `u`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Table_SubQuery_Multi(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Table(func(m *Builder) {
		m.Table(func(m *Builder) {
			m.Table("m_users").Select("sex", "count(*) as c").Group("sex")
		}).Where("c", ">", 5)
	}).ToSql()
	if sql == "SELECT * FROM (SELECT * FROM (SELECT `sex`,count(*) as `c` FROM `m_users` GROUP BY `sex`) as `tmp2` WHERE `c` > ?) as `tmp1`" &&
		reflect.DeepEqual(params, []interface{}{5}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Table_SubQuery(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Table(func(m *Builder) {
		m.Table("m_users").Select("sex", "count(*) as c").Group("sex")
	}).ToSql()
	if sql == "SELECT * FROM (SELECT `sex`,count(*) as `c` FROM `m_users` GROUP BY `sex`) as `tmp1`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Order(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name").Order("id").ToSql()
	if sql == "SELECT `id`,`name` FROM `user` ORDER BY `id` DESC" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Select("id", "name").Order("id", "asc").ToSql()
	if sql == "SELECT `id`,`name` FROM `user` ORDER BY `id` ASC" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

}

func TestBuilder_Order_Multi(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name").Order("id").Order("age", "asc").ToSql()
	if sql == "SELECT `id`,`name` FROM `user` ORDER BY `id` DESC,`age` ASC" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Group(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("age", "count(*) as c").Group("age").ToSql()
	if sql == "SELECT `age`,count(*) as `c` FROM `user` GROUP BY `age`" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Having(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("age", "count(*) as c").Group("age").Having("c", ">", 10).ToSql()
	if sql == "SELECT `age`,count(*) as `c` FROM `user` GROUP BY `age` HAVING `c` > ?" &&
		reflect.DeepEqual(params, []interface{}{10}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Limit(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name").Limit(10).ToSql()
	if sql == "SELECT `id`,`name` FROM `user` LIMIT 10" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Select("id", "name").Limit(1, 10).ToSql()
	if sql == "SELECT `id`,`name` FROM `user` LIMIT 1,10" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Page(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name").Page(1, 10).ToSql()
	if sql == "SELECT `id`,`name` FROM `user` LIMIT 0,10" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}
