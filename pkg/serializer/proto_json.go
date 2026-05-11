/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 13:20:41
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 15:26:57
 * @FilePath: \apex\go-toolbox\pkg\serializer\proto_json.go
 * @Description: 序列化 protobuf 消息为 JSON 字符串
 */
package serializer

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtoJSONMarshal 序列化 protobuf 消息为 JSON 字符串
// 支持传入 proto.Message 或 interface{} (会自动进行类型断言)
func ProtoJSONMarshal(m interface{}) (string, error) {
	if m == nil {
		return "", nil
	}

	msg, ok := m.(proto.Message)
	if !ok {
		return "", nil
	}

	b, err := protojson.Marshal(msg)
	return string(b), err
}

// ProtoJSONUnmarshal 反序列化 JSON 字符串为 protobuf 消息
// 两个参数均为 interface{}，内部自动识别 proto.Message 和 JSON 字符串，兼容任意顺序
func ProtoJSONUnmarshal(a, b interface{}) error {
	msg, str := extractMessageAndString(a, b)

	if msg == nil || str == "" {
		return nil
	}

	trimmed := strings.TrimSpace(str)
	if trimmed == "" || trimmed == "null" {
		return nil
	}
	return protojson.Unmarshal([]byte(trimmed), msg)
}

// extractMessageAndString 从两个参数中提取 proto.Message 和 JSON 字符串
func extractMessageAndString(a, b interface{}) (proto.Message, string) {
	var msg proto.Message
	var str string

	for _, v := range []interface{}{a, b} {
		m, s := classifyParam(v)
		if m != nil && msg == nil {
			msg = m
		}
		if s != "" && str == "" {
			str = s
		}
	}

	return msg, str
}

// classifyParam 分类参数为 proto.Message 或 JSON 字符串
// 支持直接传入 proto.Message、JSON 字符串、[]byte 或 fmt.Stringer 接口实现
func classifyParam(v interface{}) (proto.Message, string) {
	if v == nil {
		return nil, ""
	}
	switch val := v.(type) {
	case proto.Message:
		return val, ""
	case string:
		return nil, val
	case []byte:
		return nil, string(val)
	case fmt.Stringer:
		return nil, val.String()
	default:
		return nil, ""
	}
}
