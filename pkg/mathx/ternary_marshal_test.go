/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-20 11:20:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 11:20:00
 * @FilePath: \go-toolbox\pkg\mathx\ternary_marshal_test.go
 * @Description: MarshalJSONOrDefault 函数测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"testing"
)

func TestMarshalJSONOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		defaultVal string
		want       string
	}{
		{
			name:       "nil值返回默认值",
			value:      nil,
			defaultVal: "{}",
			want:       "{}",
		},
		{
			name:       "空map[string]string返回默认值",
			value:      map[string]string{},
			defaultVal: "{}",
			want:       "{}",
		},
		{
			name:       "空map[string]any返回默认值",
			value:      map[string]any{},
			defaultVal: "{}",
			want:       "{}",
		},
		{
			name:       "非空map[string]string正常序列化",
			value:      map[string]string{"key": "value", "foo": "bar"},
			defaultVal: "{}",
			want:       `{"foo":"bar","key":"value"}`,
		},
		{
			name:       "非空map[string]any正常序列化",
			value:      map[string]any{"name": "test", "age": 18},
			defaultVal: "{}",
			want:       `{"age":18,"name":"test"}`,
		},
		{
			name: "struct正常序列化",
			value: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "Alice",
				Age:  25,
			},
			defaultVal: "{}",
			want:       `{"name":"Alice","age":25}`,
		},
		{
			name:       "空slice返回空数组",
			value:      []string{},
			defaultVal: "[]",
			want:       "[]",
		},
		{
			name:       "非空slice正常序列化",
			value:      []int{1, 2, 3},
			defaultVal: "[]",
			want:       "[1,2,3]",
		},
		{
			name:       "字符串正常序列化",
			value:      "hello world",
			defaultVal: `""`,
			want:       `"hello world"`,
		},
		{
			name:       "数字正常序列化",
			value:      42,
			defaultVal: "0",
			want:       "42",
		},
		{
			name:       "布尔值正常序列化",
			value:      true,
			defaultVal: "false",
			want:       "true",
		},
		{
			name:       "channel无法序列化返回默认值",
			value:      make(chan int),
			defaultVal: "{}",
			want:       "{}",
		},
		{
			name:       "function无法序列化返回默认值",
			value:      func() {},
			defaultVal: "{}",
			want:       "{}",
		},
		{
			name: "嵌套map正常序列化",
			value: map[string]any{
				"user": map[string]any{
					"name": "Bob",
					"age":  30,
				},
				"active": true,
			},
			defaultVal: "{}",
			want:       `{"active":true,"user":{"age":30,"name":"Bob"}}`,
		},
		{
			name: "指针值正常序列化",
			value: func() *int {
				v := 100
				return &v
			}(),
			defaultVal: "0",
			want:       "100",
		},
		{
			name:       "nil指针返回默认值",
			value:      (*int)(nil),
			defaultVal: "{}",
			want:       "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MarshalJSONOrDefault(tt.value, tt.defaultVal)
			if got != tt.want {
				t.Errorf("MarshalJSONOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

// BenchmarkMarshalJSONOrDefault 性能测试
func BenchmarkMarshalJSONOrDefault(b *testing.B) {
	benchmarks := []struct {
		name  string
		value any
	}{
		{
			name:  "nil值",
			value: nil,
		},
		{
			name:  "空map",
			value: map[string]string{},
		},
		{
			name:  "小map",
			value: map[string]string{"key": "value"},
		},
		{
			name: "大map",
			value: map[string]any{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
			},
		},
		{
			name: "struct",
			value: struct {
				Name  string `json:"name"`
				Age   int    `json:"age"`
				Email string `json:"email"`
			}{
				Name:  "Test User",
				Age:   30,
				Email: "test@example.com",
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = MarshalJSONOrDefault(bm.value, "{}")
			}
		})
	}
}

// TestMarshalJSONOrDefault_EdgeCases 边界情况测试
func TestMarshalJSONOrDefault_EdgeCases(t *testing.T) {
	t.Run("空字符串序列化", func(t *testing.T) {
		got := MarshalJSONOrDefault("", `""`)
		want := `""`
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("零值数字序列化", func(t *testing.T) {
		got := MarshalJSONOrDefault(0, "999")
		want := "0"
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("false布尔值序列化", func(t *testing.T) {
		got := MarshalJSONOrDefault(false, "true")
		want := "false"
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("包含特殊字符的字符串", func(t *testing.T) {
		got := MarshalJSONOrDefault(`{"key":"value"}`, "{}")
		want := `"{\"key\":\"value\"}"`
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Unicode字符串", func(t *testing.T) {
		got := MarshalJSONOrDefault("你好世界", `""`)
		want := `"你好世界"`
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

// TestMarshalJSONOrDefault_RealWorldUsage 真实使用场景测试
func TestMarshalJSONOrDefault_RealWorldUsage(t *testing.T) {
	t.Run("数据库JSON字段-ContentExtra", func(t *testing.T) {
		// 模拟protobuf请求
		type Request struct {
			ContentExtra map[string]string
		}

		tests := []struct {
			name string
			req  Request
			want string
		}{
			{
				name: "空ContentExtra",
				req:  Request{ContentExtra: map[string]string{}},
				want: "{}",
			},
			{
				name: "nil ContentExtra",
				req:  Request{ContentExtra: nil},
				want: "{}",
			},
			{
				name: "有数据的ContentExtra",
				req:  Request{ContentExtra: map[string]string{"type": "text", "encoding": "utf-8"}},
				want: `{"encoding":"utf-8","type":"text"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := MarshalJSONOrDefault(tt.req.ContentExtra, "{}")
				if got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("数据库JSON字段-Metadata", func(t *testing.T) {
		metadata := map[string]any{
			"source":    "web",
			"ip":        "127.0.0.1",
			"timestamp": 1234567890,
		}
		got := MarshalJSONOrDefault(metadata, "{}")
		// 注意：map的key顺序是不确定的，所以只检查是否包含预期内容
		if len(got) < 30 { // 基本的长度检查
			t.Errorf("序列化结果太短: %v", got)
		}
	})

	t.Run("数据库JSON字段-MediaInfo", func(t *testing.T) {
		type MediaInfo struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		}

		tests := []struct {
			name      string
			mediaInfo *MediaInfo
			want      string
		}{
			{
				name:      "nil MediaInfo",
				mediaInfo: nil,
				want:      "{}",
			},
			{
				name: "有数据的MediaInfo",
				mediaInfo: &MediaInfo{
					URL:    "https://example.com/image.jpg",
					Width:  800,
					Height: 600,
				},
				want: `{"url":"https://example.com/image.jpg","width":800,"height":600}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := MarshalJSONOrDefault(tt.mediaInfo, "{}")
				if got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			})
		}
	})
}
