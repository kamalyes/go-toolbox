/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-11 13:20:41
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 15:27:24
 * @FilePath: \go-toolbox\pkg\stringx\quote_test.go
 * @Description: 引用 JSON 字符串
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteJSONBytes(t *testing.T) {
	assert.Equal(t, []byte(`"name"`), QuoteJSONBytes("name"))
	assert.Equal(t, []byte(`"a\\b"`), QuoteJSONBytes(`a\b`))
}
