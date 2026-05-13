/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 22:55:58
 * @FilePath: \go-toolbox\pkg\serializer\json_test.go
 * @Description: 对象解析和转换工具
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package safe

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type taggedStringStruct struct {
	ID          string `json:"id" gorm:"column:id;type:varchar"`
	Detail      string `json:"detail" gorm:"column:detail;type:json"`
	JSONBDetail string `json:"jsonb_detail" gorm:"column:jsonb_detail;type:jsonb"`
	Plain       string `json:"plain" gorm:"column:plain;type:text"`
}

func tagTypeIsJSON(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "json" || value == "jsonb"
}

func normalizeTaggedJSON(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "{}"
	}
	return value
}

func TestExtractStructuredTagValue(t *testing.T) {
	assert.Equal(t, "detail", ExtractStructuredTagValue("column:detail;type:json", "column"))
	assert.Equal(t, "json", ExtractStructuredTagValue("column:detail;type:json", "type"))
	assert.Equal(t, "jsonb", ExtractStructuredTagValue("column:detail;TYPE: jsonb ", "type"))
	assert.Empty(t, ExtractStructuredTagValue("column:detail", "size"))
}

func TestStringFieldAliasesByTagType(t *testing.T) {
	aliases := StringFieldAliasesByTagType[taggedStringStruct]("gorm", tagTypeIsJSON)

	for _, name := range []string{"Detail", "detail", "JSONBDetail", "jsonb_detail"} {
		_, ok := aliases[name]
		assert.True(t, ok, "alias %s should exist", name)
	}
	for _, name := range []string{"ID", "id", "Plain", "plain"} {
		_, ok := aliases[name]
		assert.False(t, ok, "alias %s should not exist", name)
	}
}

func TestNormalizeStringFieldsByTagType(t *testing.T) {
	config := &taggedStringStruct{
		ID:          "",
		Detail:      " \n\t ",
		JSONBDetail: ` {"ok":true} `,
		Plain:       "",
	}

	NormalizeStringFieldsByTagType(config, "gorm", tagTypeIsJSON, normalizeTaggedJSON)

	assert.Equal(t, "", config.ID)
	assert.Equal(t, "{}", config.Detail)
	assert.Equal(t, `{"ok":true}`, config.JSONBDetail)
	assert.Equal(t, "", config.Plain)
}

func TestNormalizeStringFieldMapByTagType(t *testing.T) {
	jsonbValue := " "
	fields := map[string]interface{}{
		"detail":       "",
		"jsonb_detail": &jsonbValue,
		"plain":        "",
	}

	NormalizeStringFieldMapByTagType[taggedStringStruct](fields, "gorm", tagTypeIsJSON, normalizeTaggedJSON)

	assert.Equal(t, "{}", fields["detail"])
	assert.Equal(t, "{}", fields["jsonb_detail"])
	assert.Equal(t, "", fields["plain"])
}
