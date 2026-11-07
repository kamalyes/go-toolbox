/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 00:00:00
 * @FilePath: \go-toolbox\pkg\pbjson\codec.go
 * @Description: 为 serializer 装配 protobuf JSON 编解码器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

// Package pbjson 提供 protobuf 消息的 JSON 编解码，并为 pkg/serializer 装配
// protobuf 感知能力
//
// go-toolbox 核心不携带任何三方依赖，protobuf 相关代码因此独立成 module
// 导入本包即可让 serializer.JSONMarshal / JSONUnmarshal 正确处理 protobuf 消息；
// 若只需要装配而不直接调用本包函数，使用空导入：
//
//	import _ "github.com/kamalyes/go-toolbox/pkg/pbjson"
package pbjson

import (
	"fmt"
	"reflect"

	"github.com/kamalyes/go-toolbox/pkg/serializer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// protoMessageType proto.Message 的反射类型
var protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

func init() {
	serializer.RegisterProtoJSONCodec(serializer.ProtoJSONCodec{
		IsMessageType: isMessageType,
		Marshal:       marshalMessage,
		Unmarshal:     unmarshalMessage,
	})
}

// isMessageType 判断类型是否实现 proto.Message 接口
func isMessageType(t reflect.Type) bool {
	return t != nil && t.Implements(protoMessageType)
}

// marshalMessage 将 protobuf 消息编码为 JSON
func marshalMessage(msg any) ([]byte, error) {
	m, ok := msg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("pbjson: %T is not a proto.Message", msg)
	}
	return protojson.Marshal(m)
}

// unmarshalMessage 将 JSON 解码到 protobuf 消息
func unmarshalMessage(data []byte, msg any) error {
	m, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("pbjson: %T is not a proto.Message", msg)
	}
	return protojson.Unmarshal(data, m)
}
