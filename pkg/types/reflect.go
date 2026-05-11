/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-09 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 13:17:27
 * @FilePath: \apex\go-toolbox\pkg\types\reflect.go
 * @Description: 反射工具函数
 */
package types

import (
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"
)

// ProtoMessageType protobuf 消息类型
var ProtoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

// ExtractJSONKey 从结构体字段标签中提取 JSON 键名
func ExtractJSONKey(fieldType reflect.StructField) string {
	jsonTag := fieldType.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return ""
	}
	if idx := strings.Index(jsonTag, ","); idx >= 0 {
		jsonTag = jsonTag[:idx]
	}
	if jsonTag == "" {
		return fieldType.Name
	}
	return jsonTag
}

// EnsureStructDefaults 确保结构体字段默认值
func EnsureStructDefaults(v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			fieldType := t.Field(i)
			if fieldType.Type.Implements(ProtoMessageType) {
				field.Set(reflect.New(fieldType.Type.Elem()))
			} else if fieldType.Type.Elem().Kind() == reflect.Struct {
				field.Set(reflect.New(fieldType.Type.Elem()))
			}
		}
	}
}

// NewProtoMessage 创建 protobuf 消息实例
func NewProtoMessage[T proto.Message]() T {
	var zero T
	t := reflect.TypeOf(zero).Elem()
	return reflect.New(t).Interface().(T)
}
