/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 13:15:15
 * @FilePath: \go-toolbox\pkg\serializer\json_test.go
 * @Description: 对象解析和转换工具
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package serializer

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type jsonProtoPayload struct {
	Name   *wrapperspb.StringValue `json:"name"`
	Age    *wrapperspb.Int32Value  `json:"age"`
	Active *wrapperspb.BoolValue   `json:"active"`
}

type jsonProtoLeaf struct {
	Label *wrapperspb.StringValue `json:"label"`
	Score *wrapperspb.Int32Value  `json:"score,omitempty"`
	Note  string                  `json:"note,omitempty"`
}

type jsonProtoComplexPayload struct {
	ID       string                             `json:"id"`
	Primary  jsonProtoLeaf                      `json:"primary"`
	Optional *jsonProtoLeaf                     `json:"optional,omitempty"`
	Items    []jsonProtoLeaf                    `json:"items"`
	PtrItems []*jsonProtoLeaf                   `json:"ptr_items"`
	Matrix   [][]*wrapperspb.Int32Value         `json:"matrix"`
	Labels   map[string]*wrapperspb.StringValue `json:"labels"`
	Buckets  map[string][]*wrapperspb.BoolValue `json:"buckets"`
	Fixed    [2]*wrapperspb.StringValue         `json:"fixed"`
	Ignored  *wrapperspb.StringValue            `json:"-"`
	Empty    *wrapperspb.StringValue            `json:"empty,omitempty"`
}

type jsonGeneratedProtoPayload struct {
	TraceID    string                                       `json:"trace_id"`
	Meta       *structpb.Struct                             `json:"meta"`
	Events     []*structpb.Struct                           `json:"events"`
	Packed     *anypb.Any                                   `json:"packed"`
	CreatedAt  *timestamppb.Timestamp                       `json:"created_at"`
	RetryAfter *durationpb.Duration                         `json:"retry_after"`
	Mask       *fieldmaskpb.FieldMask                       `json:"mask"`
	Empty      *emptypb.Empty                               `json:"empty,omitempty"`
	File       *descriptorpb.FileDescriptorProto            `json:"file"`
	Files      map[string]*descriptorpb.FileDescriptorProto `json:"files"`
	Timeline   []*timestamppb.Timestamp                     `json:"timeline"`
	Values     []*structpb.Value                            `json:"values"`
}

func buildGeneratedProtoPayload(t testing.TB) jsonGeneratedProtoPayload {
	t.Helper()

	packed, err := anypb.New(wrapperspb.String("packed-value"))
	require.NoError(t, err)

	meta, err := structpb.NewStruct(map[string]any{
		"service": "billing",
		"enabled": true,
		"limits": map[string]any{
			"qps":      float64(250),
			"burst":    float64(50),
			"features": []any{"invoice", "refund", "audit"},
		},
	})
	require.NoError(t, err)

	createdEvent, err := structpb.NewStruct(map[string]any{
		"type": "created",
		"payload": map[string]any{
			"order_id": "order-1001",
			"amount":   float64(199.95),
			"items": []any{
				map[string]any{"sku": "book", "qty": float64(2)},
				map[string]any{"sku": "pen", "qty": float64(5)},
			},
		},
	})
	require.NoError(t, err)

	updatedEvent, err := structpb.NewStruct(map[string]any{
		"type": "updated",
		"payload": map[string]any{
			"status": "paid",
			"flags":  []any{true, false, true},
		},
	})
	require.NoError(t, err)

	jsonValue, err := structpb.NewValue(map[string]any{
		"kind":  "mixed",
		"score": float64(98.5),
		"tags":  []any{"fast", "proto", "json"},
	})
	require.NoError(t, err)

	return jsonGeneratedProtoPayload{
		TraceID:    "trace-20260513-0001",
		Meta:       meta,
		Events:     []*structpb.Struct{createdEvent, updatedEvent},
		Packed:     packed,
		CreatedAt:  timestamppb.New(time.Date(2026, 5, 13, 8, 30, 0, 123000000, time.UTC)),
		RetryAfter: durationpb.New(2500 * time.Millisecond),
		Mask:       &fieldmaskpb.FieldMask{Paths: []string{"user.name", "order.total", "items.sku"}},
		Empty:      &emptypb.Empty{},
		File:       buildGeneratedFileDescriptor("billing.proto", "demo.billing", "Order"),
		Files: map[string]*descriptorpb.FileDescriptorProto{
			"audit": buildGeneratedFileDescriptor("audit.proto", "demo.audit", "AuditEvent"),
		},
		Timeline: []*timestamppb.Timestamp{
			timestamppb.New(time.Date(2026, 5, 13, 8, 30, 0, 0, time.UTC)),
			timestamppb.New(time.Date(2026, 5, 13, 8, 30, 5, 0, time.UTC)),
		},
		Values: []*structpb.Value{
			structpb.NewStringValue("alpha"),
			structpb.NewNumberValue(42),
			structpb.NewBoolValue(true),
			jsonValue,
		},
	}
}

func buildGeneratedFileDescriptor(name string, pkg string, message string) *descriptorpb.FileDescriptorProto {
	return &descriptorpb.FileDescriptorProto{
		Name:    proto.String(name),
		Package: proto.String(pkg),
		Syntax:  proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String(message),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("id"),
						JsonName: proto.String("id"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name:     proto.String("amount"),
						JsonName: proto.String("amount"),
						Number:   proto.Int32(2),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum(),
					},
				},
			},
		},
	}
}

func TestJSONMarshalProtoMessage(t *testing.T) {
	data, err := JSONMarshal(wrapperspb.String("hello"))
	require.NoError(t, err)
	assert.JSONEq(t, `"hello"`, string(data))
}

func TestJSONUnmarshalProtoMessage(t *testing.T) {
	var value *wrapperspb.StringValue
	err := JSONUnmarshal([]byte(`"hello"`), &value)
	require.NoError(t, err)
	assert.Equal(t, "hello", value.GetValue())
}

func TestJSONUnmarshalLegacyWrappedProtoMessage(t *testing.T) {
	var value *wrapperspb.StringValue
	err := JSONUnmarshal([]byte(`{"Data":"legacy"}`), &value)
	require.NoError(t, err)
	assert.Equal(t, "legacy", value.GetValue())
}

func TestJSONMarshalProtoStruct(t *testing.T) {
	payload := jsonProtoPayload{
		Name:   wrapperspb.String("test"),
		Age:    wrapperspb.Int32(25),
		Active: wrapperspb.Bool(true),
	}

	data, err := JSONMarshal(&payload)
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"test","age":25,"active":true}`, string(data))
}

func TestJSONUnmarshalProtoStruct(t *testing.T) {
	var payload jsonProtoPayload
	err := JSONUnmarshal([]byte(`{"name":"test","age":25,"active":true}`), &payload)
	require.NoError(t, err)
	assert.Equal(t, "test", payload.Name.GetValue())
	assert.Equal(t, int32(25), payload.Age.GetValue())
	assert.True(t, payload.Active.GetValue())
}

func TestJSONUnmarshalProtoStructRejectsInvalidUnknownField(t *testing.T) {
	var payload jsonProtoPayload
	err := JSONUnmarshal([]byte(`{"name":"test","unknown":[}`), &payload)
	require.Error(t, err)
}

func TestJSONUnmarshalNilTargetError(t *testing.T) {
	err := JSONUnmarshal[jsonProtoPayload]([]byte(`{}`), nil)
	require.Error(t, err)
	assert.True(t, IsJSONNilTargetError(err))
}

func TestNormalizeJSONText(t *testing.T) {
	assert.Equal(t, "{}", NormalizeJSONText(""))
	assert.Equal(t, "{}", NormalizeJSONText(" \n\t "))
	assert.Equal(t, "[]", NormalizeJSONText("", "[]"))
	assert.Equal(t, `{"name":"test"}`, NormalizeJSONText(` {"name":"test"} `))
	assert.Equal(t, `{name:"test"}`, NormalizeJSONText(` {name:"test"} `))
	assert.Equal(t, "{}", NormalizeJSONText("", "not-json"))
}

func TestJSONUnmarshalExpectedObjectError(t *testing.T) {
	var payload jsonProtoPayload
	err := JSONUnmarshal([]byte(`[]`), &payload)
	require.Error(t, err)
	assert.True(t, IsJSONExpectedObjectError(err))
}

func TestJSONErrorClassification(t *testing.T) {
	assert.True(t, IsJSONInvalidUnknownFieldValueError(NewJSONInvalidUnknownFieldValueError()))
	assert.True(t, IsJSONExpectedObjectKeySeparatorError(NewJSONExpectedObjectKeySeparatorError()))
	assert.True(t, IsJSONExpectedObjectNextError(NewJSONExpectedObjectNextError()))
	assert.True(t, IsJSONMapKeyUnsupportedError(NewJSONMapKeyUnsupportedError("int")))

	err := NewJSONFieldError("name", NewJSONExpectedObjectError())
	assert.True(t, IsJSONExpectedObjectError(err))
}

func TestJSONMarshalProtoSlice(t *testing.T) {
	values := []*wrapperspb.StringValue{
		wrapperspb.String("alpha"),
		wrapperspb.String("beta"),
	}

	data, err := JSONMarshal(values)
	require.NoError(t, err)
	assert.JSONEq(t, `["alpha","beta"]`, string(data))
}

func TestJSONUnmarshalProtoSlice(t *testing.T) {
	var values []*wrapperspb.StringValue
	err := JSONUnmarshal([]byte(`["alpha","beta"]`), &values)
	require.NoError(t, err)
	require.Len(t, values, 2)
	assert.Equal(t, "alpha", values[0].GetValue())
	assert.Equal(t, "beta", values[1].GetValue())
}

func TestJSONMarshalOrdinaryValues(t *testing.T) {
	data, err := JSONMarshal(map[string]int{"a": 1, "b": 2})
	require.NoError(t, err)
	assert.JSONEq(t, `{"a":1,"b":2}`, string(data))
}

func TestJSONMarshalComplexProtoPayload(t *testing.T) {
	// 覆盖嵌套 struct、指针、数组、切片、map、nil proto 与 omitempty 的组合场景。
	payload := jsonProtoComplexPayload{
		ID: "root",
		Primary: jsonProtoLeaf{
			Label: wrapperspb.String("main"),
			Score: wrapperspb.Int32(9),
		},
		Optional: &jsonProtoLeaf{
			Label: wrapperspb.String("optional"),
			Score: wrapperspb.Int32(7),
			Note:  "keep",
		},
		Items: []jsonProtoLeaf{
			{Label: wrapperspb.String("one"), Score: wrapperspb.Int32(1)},
			{Label: wrapperspb.String("two"), Score: wrapperspb.Int32(2), Note: "second"},
		},
		PtrItems: []*jsonProtoLeaf{
			{Label: wrapperspb.String("ptr"), Score: wrapperspb.Int32(3)},
			nil,
		},
		Matrix: [][]*wrapperspb.Int32Value{
			{wrapperspb.Int32(1), wrapperspb.Int32(2)},
			{wrapperspb.Int32(3), nil},
		},
		Labels: map[string]*wrapperspb.StringValue{
			"left":  wrapperspb.String("L"),
			"right": wrapperspb.String("R"),
		},
		Buckets: map[string][]*wrapperspb.BoolValue{
			"a": {wrapperspb.Bool(true), wrapperspb.Bool(false)},
			"b": {wrapperspb.Bool(true)},
		},
		Fixed:   [2]*wrapperspb.StringValue{wrapperspb.String("x"), wrapperspb.String("y")},
		Ignored: wrapperspb.String("hidden"),
		Empty:   wrapperspb.String(""),
	}

	data, err := JSONMarshal(&payload)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id":"root",
		"primary":{"label":"main","score":9},
		"optional":{"label":"optional","score":7,"note":"keep"},
		"items":[{"label":"one","score":1},{"label":"two","score":2,"note":"second"}],
		"ptr_items":[{"label":"ptr","score":3},null],
		"matrix":[[1,2],[3,null]],
		"labels":{"left":"L","right":"R"},
		"buckets":{"a":[true,false],"b":[true]},
		"fixed":["x","y"]
	}`, string(data))
	assert.NotContains(t, string(data), "hidden")
	assert.NotContains(t, string(data), "empty")
}

func TestJSONUnmarshalComplexProtoPayload(t *testing.T) {
	// 字段顺序、空白、未知字段、嵌套对象和 null 混在一起，验证快速扫描路径不会挑食。
	data := []byte(`{
		"unknown_object":{"nested":[1,{"ok":true}]},
		"labels":{"right":"R","left":"L"},
		"id":"root",
		"matrix":[[1,2],[3,null]],
		"fixed":["x","y"],
		"ptr_items":[{"score":3,"label":"ptr"},null],
		"primary":{"extra":[{"deep":"value"}],"score":9,"label":"main"},
		"buckets":{"a":[true,false],"b":[true]},
		"items":[{"label":"one","score":1},{"note":"second","score":2,"label":"two"}],
		"optional":{"note":"keep","label":"optional","score":7},
		"ignored":"visible-but-skipped"
	}`)

	var payload jsonProtoComplexPayload
	err := JSONUnmarshal(data, &payload)
	require.NoError(t, err)

	assert.Equal(t, "root", payload.ID)
	assert.Equal(t, "main", payload.Primary.Label.GetValue())
	assert.Equal(t, int32(9), payload.Primary.Score.GetValue())
	require.NotNil(t, payload.Optional)
	assert.Equal(t, "optional", payload.Optional.Label.GetValue())
	assert.Equal(t, int32(7), payload.Optional.Score.GetValue())
	assert.Equal(t, "keep", payload.Optional.Note)
	require.Len(t, payload.Items, 2)
	assert.Equal(t, "one", payload.Items[0].Label.GetValue())
	assert.Equal(t, "second", payload.Items[1].Note)
	require.Len(t, payload.PtrItems, 2)
	assert.Equal(t, "ptr", payload.PtrItems[0].Label.GetValue())
	assert.Nil(t, payload.PtrItems[1])
	require.Len(t, payload.Matrix, 2)
	assert.Equal(t, int32(1), payload.Matrix[0][0].GetValue())
	assert.Nil(t, payload.Matrix[1][1])
	assert.Equal(t, "L", payload.Labels["left"].GetValue())
	assert.False(t, payload.Buckets["a"][1].GetValue())
	assert.Equal(t, "x", payload.Fixed[0].GetValue())
	assert.Nil(t, payload.Ignored)
}

func TestJSONUnmarshalComplexProtoPayloadNulls(t *testing.T) {
	// 单独覆盖对象字段为 null 时，指针、slice、map 都应保持 nil/零值。
	data := []byte(`{"optional":null,"items":null,"labels":null,"primary":{"label":"main"},"fixed":["x","y"]}`)

	var payload jsonProtoComplexPayload
	err := JSONUnmarshal(data, &payload)
	require.NoError(t, err)

	assert.Nil(t, payload.Optional)
	assert.Nil(t, payload.Items)
	assert.Nil(t, payload.Labels)
	assert.Equal(t, "main", payload.Primary.Label.GetValue())
	assert.Equal(t, "x", payload.Fixed[0].GetValue())
}

func TestJSONRoundTripGeneratedProtoPayload(t *testing.T) {
	// 使用真实 generated PB 类型覆盖 protojson 的特殊编码形态：Any、Struct、Timestamp、Duration、FieldMask、DescriptorProto。
	payload := buildGeneratedProtoPayload(t)

	data, err := JSONMarshal(&payload)
	require.NoError(t, err)
	assert.JSONEq(t, string(data), string(data))
	assert.Contains(t, string(data), "type.googleapis.com/google.protobuf.StringValue")
	assert.Contains(t, string(data), "2026-05-13T08:30:00.123Z")
	assert.Contains(t, string(data), "2.500s")
	assert.Contains(t, string(data), "user.name,order.total,items.sku")

	var restored jsonGeneratedProtoPayload
	err = JSONUnmarshal(data, &restored)
	require.NoError(t, err)
	assertGeneratedProtoPayload(t, &payload, &restored)
}

func TestJSONUnmarshalGeneratedProtoPayloadFromRealProtoJSON(t *testing.T) {
	// 这里直接写 protojson 形态，模拟外部服务返回的真实 protobuf JSON。
	data := []byte(`{
		"trace_id":"trace-inline",
		"meta":{"service":"payment","enabled":true,"limits":{"qps":100,"features":["pay","audit"]}},
		"events":[{"type":"received","payload":{"status":"ok","attempts":2}}],
		"packed":{"@type":"type.googleapis.com/google.protobuf.StringValue","value":"inline-packed"},
		"created_at":"2026-05-13T09:00:01.456Z",
		"retry_after":"3.250s",
		"mask":"account.id,account.status",
		"empty":{},
		"file":{"name":"payment.proto","package":"demo.payment","messageType":[{"name":"Payment","field":[{"name":"id","number":1,"label":"LABEL_OPTIONAL","type":"TYPE_STRING","jsonName":"id"}]}],"syntax":"proto3"},
		"files":{"audit":{"name":"audit.proto","package":"demo.audit","messageType":[{"name":"Audit","field":[{"name":"amount","number":1,"label":"LABEL_OPTIONAL","type":"TYPE_DOUBLE","jsonName":"amount"}]}],"syntax":"proto3"}},
		"timeline":["2026-05-13T09:00:01Z","2026-05-13T09:00:02Z"],
		"values":["alpha",42,true,{"nested":{"ok":true}}],
		"ignored_unknown":{"deep":[{"still":"valid"}]}
	}`)

	var payload jsonGeneratedProtoPayload
	err := JSONUnmarshal(data, &payload)
	require.NoError(t, err)

	assert.Equal(t, "trace-inline", payload.TraceID)
	assert.Equal(t, "payment", payload.Meta.Fields["service"].GetStringValue())
	assert.Equal(t, "received", payload.Events[0].Fields["type"].GetStringValue())
	assert.Equal(t, time.Date(2026, 5, 13, 9, 0, 1, 456000000, time.UTC), payload.CreatedAt.AsTime())
	assert.Equal(t, 3250*time.Millisecond, payload.RetryAfter.AsDuration())
	assert.Equal(t, []string{"account.id", "account.status"}, payload.Mask.Paths)
	assert.NotNil(t, payload.Empty)
	assert.Equal(t, "payment.proto", payload.File.GetName())
	assert.Equal(t, "Audit", payload.Files["audit"].MessageType[0].GetName())
	assert.Equal(t, "alpha", payload.Values[0].GetStringValue())
	assert.Equal(t, float64(42), payload.Values[1].GetNumberValue())

	var unpacked wrapperspb.StringValue
	err = payload.Packed.UnmarshalTo(&unpacked)
	require.NoError(t, err)
	assert.Equal(t, "inline-packed", unpacked.GetValue())
}

func assertGeneratedProtoPayload(t *testing.T, expected *jsonGeneratedProtoPayload, actual *jsonGeneratedProtoPayload) {
	t.Helper()

	assert.Equal(t, expected.TraceID, actual.TraceID)
	assert.True(t, proto.Equal(expected.Meta, actual.Meta))
	require.Len(t, actual.Events, len(expected.Events))
	for i := range expected.Events {
		assert.True(t, proto.Equal(expected.Events[i], actual.Events[i]))
	}
	assert.Equal(t, expected.CreatedAt.AsTime(), actual.CreatedAt.AsTime())
	assert.Equal(t, expected.RetryAfter.AsDuration(), actual.RetryAfter.AsDuration())
	assert.Equal(t, expected.Mask.Paths, actual.Mask.Paths)
	assert.Nil(t, actual.Empty)
	assert.True(t, proto.Equal(expected.File, actual.File))
	assert.True(t, proto.Equal(expected.Files["audit"], actual.Files["audit"]))
	require.Len(t, actual.Timeline, len(expected.Timeline))
	for i := range expected.Timeline {
		assert.Equal(t, expected.Timeline[i].AsTime(), actual.Timeline[i].AsTime())
	}
	require.Len(t, actual.Values, len(expected.Values))
	for i := range expected.Values {
		assert.True(t, proto.Equal(expected.Values[i], actual.Values[i]))
	}

	var unpacked wrapperspb.StringValue
	err := actual.Packed.UnmarshalTo(&unpacked)
	require.NoError(t, err)
	assert.Equal(t, "packed-value", unpacked.GetValue())
}

func traditionalMarshalGeneratedProtoPayload(payload *jsonGeneratedProtoPayload) ([]byte, error) {
	result := make(map[string]json.RawMessage, 11)
	result["trace_id"] = mustMarshalRaw(payload.TraceID)
	if err := putProtoRaw(result, "meta", payload.Meta); err != nil {
		return nil, err
	}
	if err := putProtoSliceRaw(result, "events", payload.Events); err != nil {
		return nil, err
	}
	if err := putProtoRaw(result, "packed", payload.Packed); err != nil {
		return nil, err
	}
	if err := putProtoRaw(result, "created_at", payload.CreatedAt); err != nil {
		return nil, err
	}
	if err := putProtoRaw(result, "retry_after", payload.RetryAfter); err != nil {
		return nil, err
	}
	if err := putProtoRaw(result, "mask", payload.Mask); err != nil {
		return nil, err
	}
	if err := putProtoRaw(result, "file", payload.File); err != nil {
		return nil, err
	}
	if err := putProtoMapRaw(result, "files", payload.Files); err != nil {
		return nil, err
	}
	if err := putProtoSliceRaw(result, "timeline", payload.Timeline); err != nil {
		return nil, err
	}
	if err := putProtoSliceRaw(result, "values", payload.Values); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

func traditionalUnmarshalGeneratedProtoPayload(data []byte, payload *jsonGeneratedProtoPayload) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if value := raw["trace_id"]; len(value) > 0 {
		if err := json.Unmarshal(value, &payload.TraceID); err != nil {
			return err
		}
	}
	if err := getProtoRaw(raw, "meta", &payload.Meta); err != nil {
		return err
	}
	if err := getProtoSliceRaw(raw, "events", &payload.Events); err != nil {
		return err
	}
	if err := getProtoRaw(raw, "packed", &payload.Packed); err != nil {
		return err
	}
	if err := getProtoRaw(raw, "created_at", &payload.CreatedAt); err != nil {
		return err
	}
	if err := getProtoRaw(raw, "retry_after", &payload.RetryAfter); err != nil {
		return err
	}
	if err := getProtoRaw(raw, "mask", &payload.Mask); err != nil {
		return err
	}
	if err := getProtoRaw(raw, "empty", &payload.Empty); err != nil {
		return err
	}
	if err := getProtoRaw(raw, "file", &payload.File); err != nil {
		return err
	}
	if err := getProtoMapRaw(raw, "files", &payload.Files); err != nil {
		return err
	}
	if err := getProtoSliceRaw(raw, "timeline", &payload.Timeline); err != nil {
		return err
	}
	return getProtoSliceRaw(raw, "values", &payload.Values)
}

func mustMarshalRaw(value any) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}

func putProtoRaw[M proto.Message](result map[string]json.RawMessage, key string, value M) error {
	if reflect.ValueOf(value).IsNil() {
		result[key] = json.RawMessage("null")
		return nil
	}
	data, err := protojson.Marshal(value)
	if err != nil {
		return err
	}
	result[key] = data
	return nil
}

func putProtoSliceRaw[M proto.Message](result map[string]json.RawMessage, key string, values []M) error {
	items := make([]json.RawMessage, 0, len(values))
	for _, value := range values {
		data, err := protojson.Marshal(value)
		if err != nil {
			return err
		}
		items = append(items, data)
	}
	result[key] = mustMarshalRaw(items)
	return nil
}

func putProtoMapRaw[M proto.Message](result map[string]json.RawMessage, key string, values map[string]M) error {
	items := make(map[string]json.RawMessage, len(values))
	for name, value := range values {
		data, err := protojson.Marshal(value)
		if err != nil {
			return err
		}
		items[name] = data
	}
	result[key] = mustMarshalRaw(items)
	return nil
}

func getProtoRaw[M proto.Message](raw map[string]json.RawMessage, key string, target *M) error {
	data := raw[key]
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	msg := newProtoMessage[M]()
	if err := protojson.Unmarshal(data, msg); err != nil {
		return err
	}
	*target = msg
	return nil
}

func getProtoSliceRaw[M proto.Message](raw map[string]json.RawMessage, key string, target *[]M) error {
	data := raw[key]
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	result := make([]M, 0, len(items))
	for _, item := range items {
		msg := newProtoMessage[M]()
		if err := protojson.Unmarshal(item, msg); err != nil {
			return err
		}
		result = append(result, msg)
	}
	*target = result
	return nil
}

func getProtoMapRaw[M proto.Message](raw map[string]json.RawMessage, key string, target *map[string]M) error {
	data := raw[key]
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var items map[string]json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	result := make(map[string]M, len(items))
	for name, item := range items {
		msg := newProtoMessage[M]()
		if err := protojson.Unmarshal(item, msg); err != nil {
			return err
		}
		result[name] = msg
	}
	*target = result
	return nil
}

func newProtoMessage[M proto.Message]() M {
	var zero M
	return reflect.New(reflect.TypeOf(zero).Elem()).Interface().(M)
}

func BenchmarkJSONMarshalProtoStruct(b *testing.B) {
	payload := jsonProtoPayload{
		Name:   wrapperspb.String("benchmark"),
		Age:    wrapperspb.Int32(30),
		Active: wrapperspb.Bool(true),
	}
	_, _ = JSONMarshal(payload)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := JSONMarshal(payload); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalGeneratedProtoPayload(b *testing.B) {
	payload := buildGeneratedProtoPayload(b)
	_, _ = JSONMarshal(&payload)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := JSONMarshal(&payload); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalGeneratedProtoPayloadTraditional(b *testing.B) {
	payload := buildGeneratedProtoPayload(b)
	_, _ = traditionalMarshalGeneratedProtoPayload(&payload)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := traditionalMarshalGeneratedProtoPayload(&payload); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalGeneratedProtoPayload(b *testing.B) {
	payload := buildGeneratedProtoPayload(b)
	data, err := JSONMarshal(&payload)
	require.NoError(b, err)
	var warmup jsonGeneratedProtoPayload
	_ = JSONUnmarshal(data, &warmup)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var restored jsonGeneratedProtoPayload
		if err := JSONUnmarshal(data, &restored); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalGeneratedProtoPayloadTraditional(b *testing.B) {
	payload := buildGeneratedProtoPayload(b)
	data, err := traditionalMarshalGeneratedProtoPayload(&payload)
	require.NoError(b, err)
	var warmup jsonGeneratedProtoPayload
	_ = traditionalUnmarshalGeneratedProtoPayload(data, &warmup)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var restored jsonGeneratedProtoPayload
		if err := traditionalUnmarshalGeneratedProtoPayload(data, &restored); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalProtoStruct(b *testing.B) {
	data := []byte(`{"name":"benchmark","age":30,"active":true}`)
	var warmup jsonProtoPayload
	_ = JSONUnmarshal(data, &warmup)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var payload jsonProtoPayload
		if err := JSONUnmarshal(data, &payload); err != nil {
			b.Fatal(err)
		}
	}
}
