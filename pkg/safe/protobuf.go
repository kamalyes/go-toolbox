/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-04 09:59:53
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 10:28:55
 * @FilePath: \go-toolbox\pkg\safe\protobuf.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Ptr 通用指针安全转换
func Ptr[T any, R any](src *T, f func(T) R) *R {
	if src == nil {
		return nil
	}
	r := f(*src)
	return &r
}

// TimeToTimestampPB 安全转换 *time.Time 到 *timestamppb.Timestamp
func TimeToTimestampPB(src *time.Time) *timestamppb.Timestamp {
	if src == nil {
		return nil
	}
	return timestamppb.New(*src)
}

// StringPtr 安全转换 *string 到 string（空字符串）
func StringPtr(src *string) string {
	if src == nil {
		return ""
	}
	return *src
}

// IntPtr 安全转换 *int 到 int（默认0）
func IntPtr(src *int) int {
	if src == nil {
		return 0
	}
	return *src
}

// BoolPtr 安全转换 *bool 到 bool（默认false）
func BoolPtr(src *bool) bool {
	if src == nil {
		return false
	}
	return *src
}

// SlicePtr 安全转换 *[]T 到 []T（空切片）
func SlicePtr[T any](src *[]T) []T {
	if src == nil {
		return make([]T, 0)
	}
	return *src
}

// Float32Ptr 安全转换 *float32 到 float32（默认0）
func Float32Ptr(src *float32) float32 {
	if src == nil {
		return 0
	}
	return *src
}

// Float64Ptr 安全转换 *float64 到 float64（默认0）
func Float64Ptr(src *float64) float64 {
	if src == nil {
		return 0
	}
	return *src
}

// UintPtr 安全转换 *uint 到 uint（默认0）
func UintPtr(src *uint) uint {
	if src == nil {
		return 0
	}
	return *src
}

// Int32Ptr 安全转换 *int32 到 int32（默认0）
func Int32Ptr(src *int32) int32 {
	if src == nil {
		return 0
	}
	return *src
}

// Int64Ptr 安全转换 *int64 到 int64（默认0）
func Int64Ptr(src *int64) int64 {
	if src == nil {
		return 0
	}
	return *src
}

// DurationPtr 安全转换 *time.Duration 到 time.Duration（默认0）
func DurationPtr(src *time.Duration) time.Duration {
	if src == nil {
		return 0
	}
	return *src
}

// StringToPB 安全转换 *string 到 wrapperspb.StringValue
func StringToPB(src *string) *wrapperspb.StringValue {
	if src == nil {
		return nil
	}
	return wrapperspb.String(*src)
}

// BoolToPB 安全转换 *bool 到 wrapperspb.BoolValue
func BoolToPB(src *bool) *wrapperspb.BoolValue {
	if src == nil {
		return nil
	}
	return wrapperspb.Bool(*src)
}

// Int32ToPB 安全转换 *int32 到 wrapperspb.Int32Value
func Int32ToPB(src *int32) *wrapperspb.Int32Value {
	if src == nil {
		return nil
	}
	return wrapperspb.Int32(*src)
}

// Int64ToPB 安全转换 *int64 到 wrapperspb.Int64Value
func Int64ToPB(src *int64) *wrapperspb.Int64Value {
	if src == nil {
		return nil
	}
	return wrapperspb.Int64(*src)
}

// DoubleToPB 安全转换 *float64 到 wrapperspb.DoubleValue
func DoubleToPB(src *float64) *wrapperspb.DoubleValue {
	if src == nil {
		return nil
	}
	return wrapperspb.Double(*src)
}

// BytesPtr 安全转换 *[]byte 到 []byte（空切片）
func BytesPtr(src *[]byte) []byte {
	if src == nil {
		return make([]byte, 0)
	}
	return *src
}

type KV[K comparable, V any] map[K]V

// Ptr 安全转换 *[K]V 到 [K]V（空 ）
func PtrKV[K comparable, V any](src *KV[K, V]) KV[K, V] {
	if src == nil {
		return make(KV[K, V])
	}
	return *src
}

// PtrToTime 安全转换 *timestamppb.Timestamp 到 *time.Time
func PtrToTime(src *timestamppb.Timestamp) *time.Time {
	if src == nil {
		return nil
	}
	t := src.AsTime()
	return &t
}

// PtrToString 安全转换 *wrapperspb.StringValue 到 *string
func PtrToString(src *wrapperspb.StringValue) *string {
	if src == nil {
		return nil
	}
	s := src.GetValue()
	return &s
}

// PtrToBool 安全转换 *wrapperspb.BoolValue 到 *bool
func PtrToBool(src *wrapperspb.BoolValue) *bool {
	if src == nil {
		return nil
	}
	b := src.GetValue()
	return &b
}

// PtrToInt32 安全转换 *wrapperspb.Int32Value 到 *int32
func PtrToInt32(src *wrapperspb.Int32Value) *int32 {
	if src == nil {
		return nil
	}
	i := src.GetValue()
	return &i
}

// PtrToInt64 安全转换 *wrapperspb.Int64Value 到 *int64
func PtrToInt64(src *wrapperspb.Int64Value) *int64 {
	if src == nil {
		return nil
	}
	i := src.GetValue()
	return &i
}

// PtrToDouble 安全转换 *wrapperspb.DoubleValue 到 *float64
func PtrToDouble(src *wrapperspb.DoubleValue) *float64 {
	if src == nil {
		return nil
	}
	f := src.GetValue()
	return &f
}

// PtrToBytes 安全转换 *[]byte 到 *[]byte
func PtrToBytes(src *[]byte) *[]byte {
	if src == nil {
		return nil
	}
	return src
}

// PtrKVToSafe 安全转换 *KV[K, V] 到 *KV[K, V]
func PtrKVToSafe[K comparable, V any](src *KV[K, V]) *KV[K, V] {
	if src == nil {
		return nil
	}
	return src
}
