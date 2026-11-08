/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-08 00:00:00
 * @FilePath: \go-toolbox\pkg\argus\protoenum\protoenum.go
 * @Description: protobuf 枚举值合法性的声明式校验规则，注册进 argus 的包级规则表
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

// Package protoenum 给 argus 补一条 protobuf 枚举校验规则
//
// 为什么是独立子模块而不是直接写在 pkg/argus 里：argus 在 go-toolbox 的根模块
// 里，而根模块不依赖 google.golang.org/protobuf。把这条规则写进去会把 protobuf
// 拖进整个 go-toolbox 的依赖图，影响每一个使用者 —— 而 pkg/protox、pkg/pbmo、
// pkg/pbjson 已经确立了相反的做法：需要 protobuf 的东西各自成为子模块。本包
// 照这个先例走，于是 argus 保持 protobuf-free，而需要这条规则的人显式引入它
//
// 用法：空导入即可，注册在 init 里完成
//
//	import _ "github.com/kamalyes/go-toolbox/pkg/argus/protoenum"
//
// 然后在字段上写 `validate:"required,protoenum"`
//
// [EN] Package protoenum adds a protobuf enum validation rule to argus.
//
// Why a separate submodule rather than code inside pkg/argus: argus lives in
// go-toolbox's root module, and that module does not depend on
// google.golang.org/protobuf. Putting this rule there would pull protobuf into every
// consumer of go-toolbox, while pkg/protox, pkg/pbmo and pkg/pbjson have already
// established the opposite practice of carving protobuf-dependent code into submodules.
// This package follows that precedent, so argus stays protobuf-free and whoever needs the
// rule imports it explicitly.
//
// Usage is a blank import, since registration happens in init:
//
//	import _ "github.com/kamalyes/go-toolbox/pkg/argus/protoenum"
//
// then declare `validate:"required,protoenum"` on a field
package protoenum

import (
	"reflect"

	"github.com/kamalyes/go-toolbox/pkg/argus/rule"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Rule 枚举值合法性规则的标签名
//
// 配合 required 使用：required 拒掉 UNSPECIFIED = 0，本规则拒掉不在枚举里的数字
// 两者职责不同，所以通常一起写 `validate:"required,protoenum"`
//
// [EN] Rule is the tag name of the enum validity rule.
//
// It pairs with required, which rejects the UNSPECIFIED zero value while this rule rejects
// a number the enum does not declare. The two do different jobs, so they are usually
// written together as `validate:"required,protoenum"`
const Rule = "protoenum"

// init 把规则注册进 argus 的包级规则表
//
// 为什么注册到 rule.BuiltinRules 而不是 validator.RegisterValidation：执行校验的
// 实例往往不在调用方手里。go-rpc-gateway 的 struct tag 校验拦截器用的是一个
// sync.Once 的包私有单例，没有导出入口，所以实例级注册够不着它。BuiltinRules 是
// 导出的包级 map，而规则是在建校验计划时查表（rule/plan.go）而不是 New() 时快照，
// 所以 init 注册对那个单例同样生效
//
// 漏引入会怎样：argus 对未知规则名是 fail-closed 的，带 protoenum 标签的字段会
// 一律校验失败，而不是静默放行。所以忘了空导入的后果是接口立刻全挂、一眼就能发现，不是悄悄失去校验 —— 这个方向是刻意选的
//
// [EN] init registers the rule into argus's package-level rule table.
//
// It goes through rule.BuiltinRules rather than validator.RegisterValidation because the
// instance doing the validating is often not the caller's: go-rpc-gateway's struct tag
// interceptor uses a sync.Once package-private singleton with no exported accessor, so
// instance-level registration cannot reach it. BuiltinRules is an exported package map and
// rules are looked up while building the validation plan rather than snapshotted at New(),
// so registering here does reach that singleton.
//
// Omitting the import is loud rather than silent: argus fails closed on an unknown rule
// name, so a protoenum field rejects everything instead of passing everything. That
// direction is deliberate
func init() {
	rule.BuiltinRules[Rule] = validate
}

// validate 字段值必须是该枚举声明过的某个值
//
// 走 protoreflect 的描述符判定，不是写死一张允许值清单。两者的差别在新增枚举值
// 那天：清单（无论是 Go 里的 switch 还是 tag 里的 oneof=A B C）会悄悄漏掉新值，
// 表现为协议里加了值、服务端却拒绝它，而且没有任何地方会提醒漏改了；描述符判定自动跟着协议走
//
// 顺带说明为什么不能用 argus 自带的 oneof：它把字段值转成字符串再比，而 protobuf
// 枚举的 String() 返回的是名字，所以 `oneof=1 2 3` 这种数字候选值会把合法请求全拒掉
// 改写成名字候选值能工作，但那是把生成代码的符号名硬编码进协议注释，重命名枚举值时同样会悄悄失效
//
// 零值的处置交给 required：UNSPECIFIED 是枚举声明过的值，描述符判定会放它过去。
// 这是对的 —— 「允不允许不传」是每个字段自己的事，不该由本规则替它决定
//
// 非枚举字段返回 false：把 protoenum 标在一个不是枚举的字段上是写错了标签，
// 而放过去会让那个字段此后完全不校验，且没有任何提示
//
// [EN] validate requires a field's value to be one the enum declares.
//
// It decides from the protoreflect descriptor rather than a hard-coded list of allowed
// values. The difference shows up the day a value is added: a list, whether a Go switch or
// a oneof tag, silently omits it, so the contract has the value while the server rejects it
// and nothing points at the omission; a descriptor check follows the contract automatically.
//
// argus's own oneof is unsuitable because it stringifies the field first and a protobuf
// enum's String returns its name, so numeric candidates reject every valid request. Name
// candidates do work but hard-code generated symbol names into the contract and break just
// as quietly on a rename.
//
// The zero value is left to required, since UNSPECIFIED is a declared value and whether it
// is acceptable is each field's own business, not this rule's.
//
// A non-enum field yields false: tagging a non-enum field with protoenum is a mistake in the
// tag, and passing it would leave that field unvalidated from then on with no hint
func validate(field reflect.Value, _ string, _ bool) bool {
	if !field.IsValid() || !field.CanInterface() {
		return false
	}
	e, ok := field.Interface().(protoreflect.Enum)
	if !ok {
		return false
	}
	return e.Descriptor().Values().ByNumber(e.Number()) != nil
}
