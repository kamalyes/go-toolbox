/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-11 00:51:18
 * @FilePath: \go-toolbox\pkg\serializer\proto_json_lenient_test.go
 * @Description: LenientProtoJSONUnmarshal 单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestLenientProtoJSONUnmarshal_WrappersInt64_NumberForm(t *testing.T) {
	msg := &wrapperspb.Int64Value{}
	err := LenientProtoJSONUnmarshal([]byte(`1191749077254635521`), msg)
	require.NoError(t, err)
	assert.Equal(t, int64(1191749077254635521), msg.GetValue())
}

func TestLenientProtoJSONUnmarshal_WrappersInt64_StringForm(t *testing.T) {
	msg := &wrapperspb.Int64Value{}
	err := LenientProtoJSONUnmarshal([]byte(`"1191749077254635521"`), msg)
	require.NoError(t, err)
	assert.Equal(t, int64(1191749077254635521), msg.GetValue())
}

func TestLenientProtoJSONUnmarshal_WrappersFloat64_NumberForm(t *testing.T) {
	msg := &wrapperspb.DoubleValue{}
	err := LenientProtoJSONUnmarshal([]byte(`3.14159`), msg)
	require.NoError(t, err)
	assert.InDelta(t, 3.14159, msg.GetValue(), 0.0001)
}

func TestLenientProtoJSONUnmarshal_WrappersFloat64_StringForm(t *testing.T) {
	msg := &wrapperspb.DoubleValue{}
	err := LenientProtoJSONUnmarshal([]byte(`"3.14159"`), msg)
	require.NoError(t, err)
	assert.InDelta(t, 3.14159, msg.GetValue(), 0.0001)
}

func TestLenientProtoJSONUnmarshal_WrappersInt32_StringForm(t *testing.T) {
	msg := &wrapperspb.Int32Value{}
	err := LenientProtoJSONUnmarshal([]byte(`"42"`), msg)
	require.NoError(t, err)
	assert.Equal(t, int32(42), msg.GetValue())
}

func TestLenientProtoJSONUnmarshal_WrappersUint64_StringForm(t *testing.T) {
	msg := &wrapperspb.UInt64Value{}
	err := LenientProtoJSONUnmarshal([]byte(`"18446744073709551615"`), msg)
	require.NoError(t, err)
	assert.Equal(t, uint64(18446744073709551615), msg.GetValue())
}

func TestLenientProtoJSONUnmarshal_WrappersFloat32_StringForm(t *testing.T) {
	msg := &wrapperspb.FloatValue{}
	err := LenientProtoJSONUnmarshal([]byte(`"2.5"`), msg)
	require.NoError(t, err)
	assert.InDelta(t, float32(2.5), msg.GetValue(), 0.001)
}

func TestLenientProtoJSONUnmarshal_WrappersString(t *testing.T) {
	msg := &wrapperspb.StringValue{}
	err := LenientProtoJSONUnmarshal([]byte(`"hello world"`), msg)
	require.NoError(t, err)
	assert.Equal(t, "hello world", msg.GetValue())
}

func TestLenientProtoJSONUnmarshal_WrappersBool(t *testing.T) {
	msg := &wrapperspb.BoolValue{}
	err := LenientProtoJSONUnmarshal([]byte(`false`), msg)
	require.NoError(t, err)
	assert.Equal(t, false, msg.GetValue())
}

func TestLenientProtoJSONUnmarshal_InvalidJSON(t *testing.T) {
	msg := &wrapperspb.Int64Value{}
	err := LenientProtoJSONUnmarshal([]byte(`{invalid`), msg)
	assert.Error(t, err)
}

func TestLenientProtoJSONUnmarshal_NonNumericString(t *testing.T) {
	msg := &wrapperspb.Int64Value{}
	err := LenientProtoJSONUnmarshal([]byte(`"abc"`), msg)
	assert.Error(t, err, "non-numeric string should fail for int64")
}

func TestLenientProtoJSONOptions_Unmarshal_WithDiscardUnknown(t *testing.T) {
	opts := LenientProtoJSONOptions{
		DiscardUnknown: true,
	}
	msg := &wrapperspb.Int64Value{}
	err := opts.Unmarshal([]byte(`42`), msg)
	require.NoError(t, err)
	assert.Equal(t, int64(42), msg.GetValue())
}

func TestLenientProtoJSONOptions_ToProtojsonOptions(t *testing.T) {
	opts := LenientProtoJSONOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	protoOpts := opts.ToProtojsonOptions()
	assert.True(t, protoOpts.DiscardUnknown)
	assert.True(t, protoOpts.AllowPartial)
}

func TestLenientProtoJSONUnmarshal_SuccessPath(t *testing.T) {
	opts := LenientProtoJSONOptions{}
	msg := &wrapperspb.Int64Value{}

	err := opts.Unmarshal([]byte(`42`), msg)
	require.NoError(t, err)
	assert.Equal(t, int64(42), msg.GetValue())
}

func TestLenientProtoJSONUnmarshal_NegativeInt64String(t *testing.T) {
	msg := &wrapperspb.Int64Value{}
	err := LenientProtoJSONUnmarshal([]byte(`"-42"`), msg)
	require.NoError(t, err)
	assert.Equal(t, int64(-42), msg.GetValue())
}

func TestLenientProtoJSONUnmarshal_ZeroValues(t *testing.T) {
	tests := []struct {
		name string
		data string
		msg  *wrapperspb.Int64Value
		want int64
	}{
		{"zero number", `0`, &wrapperspb.Int64Value{}, 0},
		{"zero string", `"0"`, &wrapperspb.Int64Value{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LenientProtoJSONUnmarshal([]byte(tt.data), tt.msg)
			require.NoError(t, err)
			assert.Equal(t, tt.want, tt.msg.GetValue())
		})
	}
}
