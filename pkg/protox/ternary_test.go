/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 15:36:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 15:52:36
 * @FilePath: \go-toolbox\pkg\protox\ternary_test.go
 * @Description: proto 时间三元运算测试（随函数从 mathx 迁入）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package protox

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestIfProtoTimeOr 测试 proto 时间戳转换三元运算
func TestIfProtoTimeOr(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		setup    func() (*timestamppb.Timestamp, time.Duration)
		validate func(*testing.T, time.Time)
	}{
		{
			name: "valid proto timestamp",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				specificTime := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
				return timestamppb.New(specificTime), 0
			},
			validate: func(t *testing.T, result time.Time) {
				expected := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
				assert.Equal(t, expected.Unix(), result.Unix())
			},
		},
		{
			name: "nil proto timestamp with -30 days offset",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				return nil, -30 * 24 * time.Hour
			},
			validate: func(t *testing.T, result time.Time) {
				expected := now.Add(-30 * 24 * time.Hour)
				// 允许1秒的误差
				assert.InDelta(t, expected.Unix(), result.Unix(), 1)
			},
		},
		{
			name: "nil proto timestamp with 0 offset (current time)",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				return nil, 0
			},
			validate: func(t *testing.T, result time.Time) {
				// 应该接近当前时间
				assert.InDelta(t, now.Unix(), result.Unix(), 1)
			},
		},
		{
			name: "nil proto timestamp with 1 hour offset",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				return nil, 1 * time.Hour
			},
			validate: func(t *testing.T, result time.Time) {
				expected := now.Add(1 * time.Hour)
				assert.InDelta(t, expected.Unix(), result.Unix(), 1)
			},
		},
		{
			name: "nil proto timestamp with -30 seconds offset",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				return nil, -30 * time.Second
			},
			validate: func(t *testing.T, result time.Time) {
				expected := now.Add(-30 * time.Second)
				assert.InDelta(t, expected.Unix(), result.Unix(), 1)
			},
		},
		{
			name: "nil proto timestamp with -5 minutes offset",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				return nil, -5 * time.Minute
			},
			validate: func(t *testing.T, result time.Time) {
				expected := now.Add(-5 * time.Minute)
				assert.InDelta(t, expected.Unix(), result.Unix(), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protoTime, duration := tt.setup()
			result := IfProtoTimeOr(protoTime, duration)
			tt.validate(t, result)
		})
	}
}

// TestIfProtoTimeOrPtr 测试 proto 时间戳转指针三元运算
func TestIfProtoTimeOrPtr(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		setup    func() (*timestamppb.Timestamp, time.Duration)
		validate func(*testing.T, *time.Time)
	}{
		{
			name: "valid proto timestamp returns pointer",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				specificTime := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
				return timestamppb.New(specificTime), 0
			},
			validate: func(t *testing.T, result *time.Time) {
				require.NotNil(t, result)
				expected := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
				assert.Equal(t, expected.Unix(), result.Unix())
			},
		},
		{
			name: "nil proto timestamp with -30 days offset returns pointer",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				return nil, -30 * 24 * time.Hour
			},
			validate: func(t *testing.T, result *time.Time) {
				require.NotNil(t, result)
				expected := now.Add(-30 * 24 * time.Hour)
				assert.InDelta(t, expected.Unix(), result.Unix(), 1)
			},
		},
		{
			name: "nil proto timestamp with 0 offset returns current time pointer",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				return nil, 0
			},
			validate: func(t *testing.T, result *time.Time) {
				require.NotNil(t, result)
				assert.InDelta(t, now.Unix(), result.Unix(), 1)
			},
		},
		{
			name: "nil proto timestamp with -5 minutes offset",
			setup: func() (*timestamppb.Timestamp, time.Duration) {
				return nil, -5 * time.Minute
			},
			validate: func(t *testing.T, result *time.Time) {
				require.NotNil(t, result)
				expected := now.Add(-5 * time.Minute)
				assert.InDelta(t, expected.Unix(), result.Unix(), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protoTime, duration := tt.setup()
			result := IfProtoTimeOrPtr(protoTime, duration)
			tt.validate(t, result)
		})
	}
}

// TestIfTimeToProto 测试 time.Time 转 proto 时间戳三元运算
func TestIfTimeToProto(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		setup    func() (*time.Time, time.Duration)
		validate func(*testing.T, *timestamppb.Timestamp)
	}{
		{
			name: "valid time.Time pointer converts to proto",
			setup: func() (*time.Time, time.Duration) {
				specificTime := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
				return &specificTime, 0
			},
			validate: func(t *testing.T, result *timestamppb.Timestamp) {
				require.NotNil(t, result)
				require.True(t, result.IsValid())
				expected := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
				assert.Equal(t, expected.Unix(), result.AsTime().Unix())
			},
		},
		{
			name: "nil time.Time pointer with 0 offset",
			setup: func() (*time.Time, time.Duration) {
				return nil, 0
			},
			validate: func(t *testing.T, result *timestamppb.Timestamp) {
				require.NotNil(t, result)
				require.True(t, result.IsValid())
				assert.InDelta(t, now.Unix(), result.AsTime().Unix(), 1)
			},
		},
		{
			name: "nil time.Time pointer with -30 days offset",
			setup: func() (*time.Time, time.Duration) {
				return nil, -30 * 24 * time.Hour
			},
			validate: func(t *testing.T, result *timestamppb.Timestamp) {
				require.NotNil(t, result)
				require.True(t, result.IsValid())
				expected := now.Add(-30 * 24 * time.Hour)
				assert.InDelta(t, expected.Unix(), result.AsTime().Unix(), 1)
			},
		},
		{
			name: "nil time.Time pointer with 1 hour offset",
			setup: func() (*time.Time, time.Duration) {
				return nil, 1 * time.Hour
			},
			validate: func(t *testing.T, result *timestamppb.Timestamp) {
				require.NotNil(t, result)
				require.True(t, result.IsValid())
				expected := now.Add(1 * time.Hour)
				assert.InDelta(t, expected.Unix(), result.AsTime().Unix(), 1)
			},
		},
		{
			name: "nil time.Time pointer with -30 seconds offset",
			setup: func() (*time.Time, time.Duration) {
				return nil, -30 * time.Second
			},
			validate: func(t *testing.T, result *timestamppb.Timestamp) {
				require.NotNil(t, result)
				require.True(t, result.IsValid())
				expected := now.Add(-30 * time.Second)
				assert.InDelta(t, expected.Unix(), result.AsTime().Unix(), 1)
			},
		},
		{
			name: "nil time.Time pointer with -5 minutes offset",
			setup: func() (*time.Time, time.Duration) {
				return nil, -5 * time.Minute
			},
			validate: func(t *testing.T, result *timestamppb.Timestamp) {
				require.NotNil(t, result)
				require.True(t, result.IsValid())
				expected := now.Add(-5 * time.Minute)
				assert.InDelta(t, expected.Unix(), result.AsTime().Unix(), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timePtr, duration := tt.setup()
			result := IfTimeToProto(timePtr, duration)
			tt.validate(t, result)
		})
	}
}

// BenchmarkIfProtoTimeOr 基准测试 IfProtoTimeOr
func BenchmarkIfProtoTimeOr(b *testing.B) {
	protoTime := timestamppb.New(time.Now())
	for i := 0; i < b.N; i++ {
		IfProtoTimeOr(protoTime, -30*24*time.Hour)
	}
}

// BenchmarkIfProtoTimeOrPtr 基准测试 IfProtoTimeOrPtr
func BenchmarkIfProtoTimeOrPtr(b *testing.B) {
	protoTime := timestamppb.New(time.Now())
	for i := 0; i < b.N; i++ {
		IfProtoTimeOrPtr(protoTime, -30*24*time.Hour)
	}
}

// BenchmarkIfTimeToProto 基准测试 IfTimeToProto
func BenchmarkIfTimeToProto(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		IfTimeToProto(&now, -30*24*time.Hour)
	}
}

// TestTimeConversionRoundTrip 测试时间转换的往返一致性
func TestTimeConversionRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "specific time round trip",
			time: time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC),
		},
		{
			name: "current time round trip",
			time: time.Now(),
		},
		{
			name: "past time round trip",
			time: time.Now().Add(-7 * 24 * time.Hour),
		},
		{
			name: "future time round trip",
			time: time.Now().Add(3 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// time.Time -> proto -> time.Time
			protoTime := timestamppb.New(tt.time)
			result1 := IfProtoTimeOr(protoTime, 0)
			assert.Equal(t, tt.time.Unix(), result1.Unix())

			// time.Time -> proto (via IfTimeToProto) -> time.Time
			timePtr := tt.time
			protoTime2 := IfTimeToProto(&timePtr, 0)
			result2 := IfProtoTimeOr(protoTime2, 0)
			assert.Equal(t, tt.time.Unix(), result2.Unix())
		})
	}
}
