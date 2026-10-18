/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 15:21:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 15:52:36
 * @FilePath: \go-toolbox\pkg\protox\safe.go
 * @Description: proto 空安全转换（从 safe 迁入：safe 的其余能力不需要 protobuf）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package protox

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// TimeToTimestampPB 安全转换 *time.Time 到 *timestamppb.Timestamp
func TimeToTimestampPB(src *time.Time) *timestamppb.Timestamp {
	if src == nil {
		return nil
	}
	return timestamppb.New(*src)
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
