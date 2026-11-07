/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\serializer\protojson_codec.go
 * @Description: 可选的 protobuf JSON 编解码器装配
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package serializer

import (
	"reflect"
	"sync/atomic"
)

// ProtoJSONCodec 可选的 protobuf JSON 编解码器
//
// serializer 核心因此不依赖 google.golang.org/protobuf；需要 protobuf 感知的
// JSON 编解码时，空导入 pkg/pbjson 子模块即可装配：
//
//	import _ "github.com/kamalyes/go-toolbox/pkg/pbjson"
//
// 未装配时 JSONMarshal / JSONUnmarshal 对任何类型都退化为标准库 encoding/json 行为
type ProtoJSONCodec struct {
	// IsMessageType 判断 t 是否为 protobuf 消息类型，t 为指针类型
	IsMessageType func(t reflect.Type) bool
	// Marshal 将 protobuf 消息编码为 JSON
	Marshal func(msg any) ([]byte, error)
	// Unmarshal 将 JSON 解码到 protobuf 消息
	Unmarshal func(data []byte, msg any) error
}

// protoJSONCodec 持有已装配的编解码器
// 用原子指针而非逐字段 var：读路径只需一次原子载入，零分配，且与 init 期装配无数据竞争
var protoJSONCodec atomic.Pointer[ProtoJSONCodec]

// RegisterProtoJSONCodec 装配 ProtoJSONCodec，应在 init 阶段调用（空导入子模块即可）
// 装配时清空类型判定缓存，使装配前已缓存的"不含 protobuf 消息"结论重新计算
func RegisterProtoJSONCodec(c ProtoJSONCodec) {
	protoJSONCodec.Store(&c)
	needsProtoJSONCache.Clear()
}

// isProtoMessageType 判断 t 是否为 protobuf 消息类型，未装配编解码器时恒为 false
func isProtoMessageType(t reflect.Type) bool {
	c := protoJSONCodec.Load()
	if c == nil || c.IsMessageType == nil {
		return false
	}
	return c.IsMessageType(t)
}

// marshalProtoJSON 将 protobuf 消息编码为 JSON，未装配编解码器时返回错误
func marshalProtoJSON(msg any) ([]byte, error) {
	c := protoJSONCodec.Load()
	if c == nil || c.Marshal == nil {
		return nil, NewJSONProtoCodecMissingError()
	}
	return c.Marshal(msg)
}
