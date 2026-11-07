/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\pbjson\core_proto_test.go
 * @Description: 核心包在 protobuf 类型上的行为测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package pbjson_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/types"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// TestMustStringWithProtobufTimestamp 测试 protobuf Timestamp 转换
func TestMustStringWithProtobufTimestamp(t *testing.T) {
	// 创建一个 protobuf Timestamp (2024-01-01 12:00:00 UTC)
	pbTime := timestamppb.New(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))

	// 默认格式 (RFC3339)
	result := convert.MustString(pbTime)
	assert.Equal(t, "2024-01-01T12:00:00Z", result)

	// 自定义格式
	customLayout := "2006-01-02 15:04:05"
	result2 := convert.MustString(pbTime, customLayout)
	assert.Equal(t, "2024-01-01 12:00:00", result2)

	// 测试 nil protobuf Timestamp
	var nilPbTime *timestamppb.Timestamp
	result3 := convert.MustString(nilPbTime)
	assert.Equal(t, "", result3) // nil 转换为 ""，而不是 "null"
}

// TestEnsureStructDefaultsProtoPointer 验证 types.EnsureStructDefaults 会初始化 nil protobuf 指针字段
// 原为 pkg/types 的子测试，因携带 protobuf 依赖而移入本模块
func TestEnsureStructDefaultsProtoPointer(t *testing.T) {
	type S struct {
		PB *wrapperspb.StringValue
	}
	s := &S{}
	v := reflect.ValueOf(s).Elem()
	types.EnsureStructDefaults(v)
	assert.NotNil(t, s.PB)
}
