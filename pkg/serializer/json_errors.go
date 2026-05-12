/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 13:15:15
 * @FilePath: \go-toolbox\pkg\serializer\json_errors.go
 * @Description: JSON错误定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package serializer

import (
	"errors"
	"fmt"
)

const (
	jsonFieldErrorFormat = "field %s: %w"
	jsonItemErrorFormat  = "item %d: %w"
	jsonKeyErrorFormat   = "key %s: %w"
)

var (
	ErrJSONNilTarget                  = errors.New("serializer json: nil target")
	ErrJSONUnexpectedEndObject        = errors.New("serializer json: unexpected end of object")
	ErrJSONExpectedObject             = errors.New("serializer json: expected object")
	ErrJSONExpectedArray              = errors.New("serializer json: expected array")
	ErrJSONExpectedObjectKeySeparator = errors.New("serializer json: expected ':' after object key")
	ErrJSONInvalidUnknownFieldValue   = errors.New("serializer json: invalid unknown field value")
	ErrJSONExpectedObjectNext         = errors.New("serializer json: expected ',' or '}' after object value")
	ErrJSONExpectedArrayNext          = errors.New("serializer json: expected ',' or ']' after array item")
	ErrJSONMapKeyUnsupported          = errors.New("serializer json: proto-aware map only supports string keys")
)

// NewJSONNilTargetError 创建 JSON 反序列化目标为空错误
func NewJSONNilTargetError() error {
	return ErrJSONNilTarget
}

// IsJSONNilTargetError 判断错误是否为 JSON 反序列化目标为空
func IsJSONNilTargetError(err error) bool {
	return errors.Is(err, ErrJSONNilTarget)
}

// NewJSONUnexpectedEndObjectError 创建 JSON 对象意外结束错误
func NewJSONUnexpectedEndObjectError() error {
	return ErrJSONUnexpectedEndObject
}

// IsJSONUnexpectedEndObjectError 判断错误是否为 JSON 对象意外结束
func IsJSONUnexpectedEndObjectError(err error) bool {
	return errors.Is(err, ErrJSONUnexpectedEndObject)
}

// NewJSONExpectedObjectError 创建期望 JSON 对象错误
func NewJSONExpectedObjectError() error {
	return ErrJSONExpectedObject
}

// IsJSONExpectedObjectError 判断错误是否为期望 JSON 对象错误
func IsJSONExpectedObjectError(err error) bool {
	return errors.Is(err, ErrJSONExpectedObject)
}

// NewJSONExpectedArrayError 创建期望 JSON 数组错误
func NewJSONExpectedArrayError() error {
	return ErrJSONExpectedArray
}

// IsJSONExpectedArrayError 判断错误是否为期望 JSON 数组错误
func IsJSONExpectedArrayError(err error) bool {
	return errors.Is(err, ErrJSONExpectedArray)
}

// NewJSONExpectedObjectKeySeparatorError 创建对象键值分隔符错误
func NewJSONExpectedObjectKeySeparatorError() error {
	return ErrJSONExpectedObjectKeySeparator
}

// IsJSONExpectedObjectKeySeparatorError 判断错误是否为对象键值分隔符错误
func IsJSONExpectedObjectKeySeparatorError(err error) bool {
	return errors.Is(err, ErrJSONExpectedObjectKeySeparator)
}

// NewJSONInvalidUnknownFieldValueError 创建未知字段值非法错误
func NewJSONInvalidUnknownFieldValueError() error {
	return ErrJSONInvalidUnknownFieldValue
}

// IsJSONInvalidUnknownFieldValueError 判断错误是否为未知字段值非法
func IsJSONInvalidUnknownFieldValueError(err error) bool {
	return errors.Is(err, ErrJSONInvalidUnknownFieldValue)
}

// NewJSONExpectedObjectNextError 创建对象值后缺少逗号或结束符错误
func NewJSONExpectedObjectNextError() error {
	return ErrJSONExpectedObjectNext
}

// IsJSONExpectedObjectNextError 判断错误是否为对象值后缺少逗号或结束符
func IsJSONExpectedObjectNextError(err error) bool {
	return errors.Is(err, ErrJSONExpectedObjectNext)
}

// NewJSONExpectedArrayNextError 创建数组元素后缺少逗号或结束符错误
func NewJSONExpectedArrayNextError() error {
	return ErrJSONExpectedArrayNext
}

// IsJSONExpectedArrayNextError 判断错误是否为数组元素后缺少逗号或结束符
func IsJSONExpectedArrayNextError(err error) bool {
	return errors.Is(err, ErrJSONExpectedArrayNext)
}

// NewJSONMapKeyUnsupportedError 创建 proto-aware map 键类型不支持错误
func NewJSONMapKeyUnsupportedError(keyType string) error {
	return fmt.Errorf("%w, got %s", ErrJSONMapKeyUnsupported, keyType)
}

// IsJSONMapKeyUnsupportedError 判断错误是否为 proto-aware map 键类型不支持
func IsJSONMapKeyUnsupportedError(err error) bool {
	return errors.Is(err, ErrJSONMapKeyUnsupported)
}

func isJSONStructScanError(err error) bool {
	return IsJSONUnexpectedEndObjectError(err) ||
		IsJSONExpectedObjectError(err) ||
		IsJSONExpectedObjectKeySeparatorError(err) ||
		IsJSONInvalidUnknownFieldValueError(err) ||
		IsJSONExpectedObjectNextError(err)
}

// NewJSONFieldError 包装字段级 JSON 错误，保留字段名上下文
func NewJSONFieldError(name string, err error) error {
	return fmt.Errorf(jsonFieldErrorFormat, name, err)
}

// NewJSONItemError 包装数组或切片元素级 JSON 错误，保留下标上下文
func NewJSONItemError(index int, err error) error {
	return fmt.Errorf(jsonItemErrorFormat, index, err)
}

// NewJSONKeyError 包装 map 键对应值的 JSON 错误，保留键上下文
func NewJSONKeyError(key string, err error) error {
	return fmt.Errorf(jsonKeyErrorFormat, key, err)
}

// NewJSONArrayTooLongError 创建 JSON 数组长度超过目标数组长度错误
func NewJSONArrayTooLongError(items int, capacity int) error {
	return fmt.Errorf("serializer json: array has %d items, target array has %d", items, capacity)
}
