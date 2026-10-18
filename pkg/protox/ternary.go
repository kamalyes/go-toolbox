/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 15:21:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 15:52:36
 * @FilePath: \go-toolbox\pkg\protox\ternary.go
 * @Description: proto 时间戳相关的三元运算（从 mathx 分出，避免 mathx 引入 protobuf）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package protox

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// IfProtoTimeOr proto时间戳转换三元运算
// 如果 proto 时间戳不为 nil 且有效，转换为 time.Time；否则返回当前时间加上偏移时间
// 支持任意精度的时间偏移（秒、分钟、小时、天等），正数表示未来，负数表示过去
//
// 示例：
//
//	30天前作为默认值
//	startTime := mathx.IfProtoTimeOr(req.StartTime, -30*24*time.Hour)
//
//	当前时间作为默认值
//	endTime := mathx.IfProtoTimeOr(req.EndTime, 0)
//
//	1小时后作为默认值
//	futureTime := mathx.IfProtoTimeOr(req.FutureTime, 1*time.Hour)
//
//	30秒前作为默认值
//	recentTime := mathx.IfProtoTimeOr(req.RecentTime, -30*time.Second)
func IfProtoTimeOr(protoTime *timestamppb.Timestamp, duration time.Duration) time.Time {
	if protoTime != nil && protoTime.IsValid() {
		return protoTime.AsTime()
	}
	return time.Now().Add(duration)
}

// IfProtoTimeOrPtr proto时间戳转指针三元运算
// 如果 proto 时间戳不为 nil 且有效，转换为 *time.Time；否则返回当前时间加上偏移时间的指针
// 支持任意精度的时间偏移，适用于需要返回指针的场景
//
// 示例：
//
//	startTime := mathx.IfProtoTimeOrPtr(req.StartTime, -30*24*time.Hour)  // 返回 *time.Time
//	endTime := mathx.IfProtoTimeOrPtr(req.EndTime, 0)                      // 返回 *time.Time
//	recentTime := mathx.IfProtoTimeOrPtr(req.RecentTime, -5*time.Minute)  // 5分钟前
func IfProtoTimeOrPtr(protoTime *timestamppb.Timestamp, duration time.Duration) *time.Time {
	t := IfProtoTimeOr(protoTime, duration)
	return &t
}

// IfTimeToProto time.Time 转 proto 时间戳三元运算
// 如果 time 指针不为 nil，转换为 *timestamppb.Timestamp；否则返回当前时间加上偏移时间的 proto 时间戳
// 支持任意精度的时间偏移，适用于反向转换场景
//
// 示例：
//
//	将 time.Time 转为 proto（存在则转换）
//	protoTime := mathx.IfTimeToProto(timePtr, 0)
//
//	将 time.Time 转为 proto（nil时用30天前）
//	protoStartTime := mathx.IfTimeToProto(startTimePtr, -30*24*time.Hour)
//
//	将 time.Time 转为 proto（nil时用1小时后）
//	protoFutureTime := mathx.IfTimeToProto(futureTimePtr, 1*time.Hour)
func IfTimeToProto(t *time.Time, duration time.Duration) *timestamppb.Timestamp {
	if t != nil {
		return timestamppb.New(*t)
	}
	defaultTime := time.Now().Add(duration)
	return timestamppb.New(defaultTime)
}
