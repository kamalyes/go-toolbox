/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-08 00:00:00
 * @FilePath: \go-toolbox\pkg\argus\protoenum\protoenum_test.go
 * @Description: protoenum 规则的测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package protoenum

import (
	"reflect"
	"testing"

	validator "github.com/kamalyes/go-toolbox/pkg/argus"
	"github.com/kamalyes/go-toolbox/pkg/argus/rule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/descriptorpb"
)

// TestRuleIsRegistered 规则必须出现在 argus 的包级规则表里
//
// 这是本包存在的全部意义：注册没发生，带 protoenum 标签的字段会因为「未知规则」
// 一律校验失败
//
// [EN] TestRuleIsRegistered: the rule must appear in argus's package-level table.
//
// That is this package's entire purpose: without the registration, a protoenum-tagged
// field fails validation outright as an unknown rule
func TestRuleIsRegistered(t *testing.T) {
	_, ok := rule.BuiltinRules[Rule]
	assert.True(t, ok, "protoenum 未注册进 rule.BuiltinRules | the rule is not registered")
}

// TestValidateAcceptsDeclaredValues 声明过的枚举值必须通过
// [EN] TestValidateAcceptsDeclaredValues: a declared enum value must pass
func TestValidateAcceptsDeclaredValues(t *testing.T) {
	for _, v := range []descriptorpb.FieldDescriptorProto_Type{
		descriptorpb.FieldDescriptorProto_TYPE_STRING,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_BOOL,
	} {
		assert.True(t, validate(reflect.ValueOf(v), "", false),
			"声明过的枚举值被拒绝 | a declared value was rejected: %v", v)
	}
}

// TestValidateRejectsUndeclaredNumber 枚举外的数字必须被拒
//
// 这正是写死允许值清单抓不住的那一类：protobuf 的枚举字段在 Go 里就是一个
// int32，赋一个没声明过的数字编译完全合法
//
// [EN] TestValidateRejectsUndeclaredNumber: a number outside the enum must be rejected.
//
// This is what a hard-coded list of allowed values cannot catch: a protobuf enum field is
// an int32 in Go and assigning an undeclared number compiles perfectly well
func TestValidateRejectsUndeclaredNumber(t *testing.T) {
	assert.False(t, validate(reflect.ValueOf(descriptorpb.FieldDescriptorProto_Type(9999)), "", false),
		"枚举外的数字被放过 | an undeclared number passed")
}

// TestValidateRejectsNonEnum 把规则标在非枚举字段上必须被拒
//
// 标错了标签应当立刻失败，而不是让那个字段此后完全不校验
//
// [EN] TestValidateRejectsNonEnum: the rule on a non-enum field must be rejected, since a
// mistagged field should fail at once rather than silently go unvalidated from then on
func TestValidateRejectsNonEnum(t *testing.T) {
	assert.False(t, validate(reflect.ValueOf("string"), "", false), "字符串被当成了枚举 | a string passed as an enum")
	assert.False(t, validate(reflect.ValueOf(42), "", false), "整数被当成了枚举 | an int passed as an enum")
	assert.False(t, validate(reflect.Value{}, "", false), "零值 reflect.Value 被放过 | an invalid reflect.Value passed")
}

// TestZeroValueIsLeftToRequired 零值由 required 决定，本规则放它过去
//
// UNSPECIFIED 是枚举声明过的值，所以描述符判定本就该通过。「允不允许不传」是
// 每个字段自己的事：有的字段可以缺省，本规则替它决定就错了
//
// [EN] TestZeroValueIsLeftToRequired: the zero value is required's business and this rule
// passes it.
//
// UNSPECIFIED is a declared value, so a descriptor check should pass it. Whether absence is
// acceptable is each field's own business, and deciding that here would be wrong
func TestZeroValueIsLeftToRequired(t *testing.T) {
	zero := descriptorpb.FieldDescriptorProto_Type(0)
	if zero.Descriptor().Values().ByNumber(0) == nil {
		t.Skip("该枚举没有声明 0 值，换一个枚举才能验证这件事 | this enum declares no zero value")
	}
	assert.True(t, validate(reflect.ValueOf(zero), "", false),
		"零值被本规则拒绝了，那是 required 的职责 | the zero value was rejected here instead of by required")
}

// TestEndToEndThroughValidator 经由 argus 的校验器走通一遍
//
// 前面几条都是直接调规则函数，只能证明规则本身对。这一条证明的是注册真的接上了：
// 校验器在建计划时查 BuiltinRules，所以 init 里的注册必须对它生效 —— 这正是本包
// 不走 validator.RegisterValidation 的原因（那只对单个实例生效，而执行校验的实例
// 常常不在调用方手里）
//
// [EN] TestEndToEndThroughValidator goes through argus's validator.
//
// The earlier cases call the rule function directly and only prove the rule itself. This one
// proves the registration connects: the validator looks BuiltinRules up while building its
// plan, so an init-time registration must reach it — which is exactly why this package does
// not use validator.RegisterValidation, as that only affects one instance while the instance
// doing the validating is often not the caller's
func TestEndToEndThroughValidator(t *testing.T) {
	type request struct {
		Kind descriptorpb.FieldDescriptorProto_Type `validate:"required,protoenum"`
	}

	v := validator.New(validator.WithRequiredStructEnabled())

	require.NoError(t, v.Struct(&request{Kind: descriptorpb.FieldDescriptorProto_TYPE_STRING}),
		"合法枚举值被拒绝 | a valid enum value was rejected")
	assert.Error(t, v.Struct(&request{Kind: descriptorpb.FieldDescriptorProto_Type(9999)}),
		"枚举外的数字通过了校验器 | an undeclared number passed the validator")
	assert.Error(t, v.Struct(&request{}),
		"零值通过了 required | the zero value passed required")
}
