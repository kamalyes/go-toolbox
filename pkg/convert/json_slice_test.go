/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 15:15:16
 * @FilePath: \go-toolbox\pkg\convert\json_slice_test.go
 * @Description: 切片与JSON字符串互转测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringsToJSON(t *testing.T) {
	t.Run("nil slice", func(t *testing.T) {
		result := StringsToJSON(nil)
		assert.Empty(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := StringsToJSON([]string{})
		assert.Empty(t, result)
	})

	t.Run("single element", func(t *testing.T) {
		result := StringsToJSON([]string{"CN"})
		assert.Equal(t, `["CN"]`, result)
	})

	t.Run("multiple elements", func(t *testing.T) {
		result := StringsToJSON([]string{"CN", "US", "JP"})
		assert.Equal(t, `["CN","US","JP"]`, result)
	})

	t.Run("elements with special characters", func(t *testing.T) {
		result := StringsToJSON([]string{"hello world", "a\"b", "c\\d"})
		assert.Contains(t, result, `"hello world"`)
		assert.Contains(t, result, `"a\"b"`)
	})

	t.Run("empty strings in slice", func(t *testing.T) {
		result := StringsToJSON([]string{"", "a", ""})
		assert.Equal(t, `["","a",""]`, result)
	})
}

func TestStringsFromJSON(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		result, err := StringsFromJSON("")
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("valid JSON array", func(t *testing.T) {
		result, err := StringsFromJSON(`["CN","US","JP"]`)
		require.NoError(t, err)
		assert.Equal(t, []string{"CN", "US", "JP"}, result)
	})

	t.Run("single element array", func(t *testing.T) {
		result, err := StringsFromJSON(`["CN"]`)
		require.NoError(t, err)
		assert.Equal(t, []string{"CN"}, result)
	})

	t.Run("empty JSON array", func(t *testing.T) {
		result, err := StringsFromJSON(`[]`)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := StringsFromJSON(`not json`)
		assert.Error(t, err)
	})

	t.Run("JSON object instead of array", func(t *testing.T) {
		_, err := StringsFromJSON(`{"key":"value"}`)
		assert.Error(t, err)
	})

	t.Run("null JSON", func(t *testing.T) {
		result, err := StringsFromJSON(`null`)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestStringsToFromJSONRoundTrip(t *testing.T) {
	t.Run("round trip preserves data", func(t *testing.T) {
		original := []string{"CN", "US", "JP", "KR"}
		jsonStr := StringsToJSON(original)
		restored, err := StringsFromJSON(jsonStr)
		require.NoError(t, err)
		assert.Equal(t, original, restored)
	})

	t.Run("round trip empty slice", func(t *testing.T) {
		jsonStr := StringsToJSON([]string{})
		assert.Empty(t, jsonStr)
		restored, err := StringsFromJSON(jsonStr)
		require.NoError(t, err)
		assert.Nil(t, restored)
	})

	t.Run("round trip with empty strings", func(t *testing.T) {
		original := []string{"", "a", ""}
		jsonStr := StringsToJSON(original)
		restored, err := StringsFromJSON(jsonStr)
		require.NoError(t, err)
		assert.Equal(t, original, restored)
	})
}
