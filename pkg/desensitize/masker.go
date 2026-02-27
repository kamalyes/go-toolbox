/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-19 10:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-19 10:00:00
 * @FilePath: \go-toolbox\pkg\desensitize\masker.go
 * @Description: 数据脱敏器，用于日志和 API 响应的敏感数据脱敏
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

import (
	"encoding/json"
	"regexp"
	"strings"
	"sync"
)

// MaskerConfig 脱敏器配置
type MaskerConfig struct {
	SensitiveKeys []string // 敏感字段列表（如 password, token, secret）
	SensitiveMask string   // 掩码字符（默认 "***"）
	MaxBodySize   int      // 最大处理数据大小（字节）
}

// DataMasker 数据脱敏器（单例模式）
type DataMasker struct {
	config     *MaskerConfig
	regexCache map[string]*regexp.Regexp
	regexMutex sync.RWMutex
}

// NewMasker 创建数据脱敏器（支持可选配置）
func NewMasker(configs ...*MaskerConfig) *DataMasker {
	var config *MaskerConfig
	if len(configs) > 0 && configs[0] != nil {
		config = configs[0]
	} else {
		config = DefaultMaskerConfig()
	}
	return &DataMasker{
		config:     config,
		regexCache: make(map[string]*regexp.Regexp),
	}
}

func DefaultMaskerConfig() *MaskerConfig {
	return &MaskerConfig{
		SensitiveKeys: []string{
			"password", "passwd", "pwd",
			"token", "accesstoken", "access_token",
			"secret", "secretkey", "secret_key",
			"apikey", "api_key",
			"authorization",
			"cookie", "session",
			"credit_card", "creditcard",
			"ssn", "id_card", "idcard",
		},
		SensitiveMask: "***",
		MaxBodySize:   10240, // 10KB
	}
}

// Mask 脱敏数据（主入口）
func (dm *DataMasker) Mask(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// 截断超长数据
	maxSize := dm.getMaxBodySize()
	if len(data) > maxSize {
		data = data[:maxSize]
	}

	// 快速判断：是否可能是 JSON（以 { 或 [ 开头）
	if len(data) > 0 && (data[0] == '{' || data[0] == '[') {
		if masked := dm.maskJSON(data); masked != "" {
			return masked
		}
	}

	// 文本脱敏
	return dm.maskText(data)
}

// maskJSON 脱敏 JSON 数据
func (dm *DataMasker) maskJSON(data []byte) string {
	var jsonData map[string]any
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return ""
	}

	dm.maskJSONFields(jsonData)

	masked, err := json.Marshal(jsonData)
	if err != nil {
		return ""
	}
	return string(masked)
}

// maskJSONFields 递归脱敏 JSON 字段
func (dm *DataMasker) maskJSONFields(data any) {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			if dm.isSensitive(key) {
				v[key] = dm.getMask()
			} else {
				dm.maskJSONFields(value)
			}
		}
	case []any:
		for _, item := range v {
			dm.maskJSONFields(item)
		}
	}
}

// getRegex 获取或编译正则表达式（带缓存）
func (dm *DataMasker) getRegex(key string) *regexp.Regexp {
	// 快速路径：读锁检查缓存
	dm.regexMutex.RLock()
	re, exists := dm.regexCache[key]
	dm.regexMutex.RUnlock()

	if exists {
		return re
	}

	// 慢速路径：写锁编译并缓存
	dm.regexMutex.Lock()
	defer dm.regexMutex.Unlock()

	// Double-check（避免重复编译）
	if re, exists := dm.regexCache[key]; exists {
		return re
	}

	// 编译正则表达式（使用 QuoteMeta 防止注入）
	pattern := `(?i)"?` + regexp.QuoteMeta(key) + `"?\s*[:=]\s*"?[^"&,}\s]+`
	re = regexp.MustCompile(pattern)
	dm.regexCache[key] = re
	return re
}

// maskText 脱敏文本数据（使用正则缓存）
func (dm *DataMasker) maskText(data []byte) string {
	result := string(data)
	mask := dm.getMask()
	lowerResult := strings.ToLower(result)

	// 只对包含敏感字段的数据进行正则替换
	for _, key := range dm.config.SensitiveKeys {
		if !strings.Contains(lowerResult, key) {
			continue
		}
		re := dm.getRegex(key)
		result = re.ReplaceAllString(result, key+"="+mask)
	}

	return result
}

// isSensitive 检查是否为敏感字段
func (dm *DataMasker) isSensitive(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, sensitive := range dm.config.SensitiveKeys {
		if strings.Contains(lowerKey, sensitive) {
			return true
		}
	}
	return false
}

// getMask 获取掩码
func (dm *DataMasker) getMask() string {
	if dm.config.SensitiveMask != "" {
		return dm.config.SensitiveMask
	}
	return DefaultMaskerConfig().SensitiveMask
}

// getMaxBodySize 获取最大 body 大小
func (dm *DataMasker) getMaxBodySize() int {
	if dm.config.MaxBodySize > 0 {
		return dm.config.MaxBodySize
	}
	return DefaultMaskerConfig().MaxBodySize
}

// WithSensitiveKeys 设置敏感字段列表（链式调用）
func (dm *DataMasker) WithSensitiveKeys(keys ...string) *DataMasker {
	dm.config.SensitiveKeys = keys
	return dm
}

// AddSensitiveKeys 添加敏感字段（链式调用）
func (dm *DataMasker) AddSensitiveKeys(keys ...string) *DataMasker {
	dm.config.SensitiveKeys = append(dm.config.SensitiveKeys, keys...)
	return dm
}

// WithMask 设置掩码字符（链式调用）
func (dm *DataMasker) WithMask(mask string) *DataMasker {
	dm.config.SensitiveMask = mask
	return dm
}

// WithMaxBodySize 设置最大处理数据大小（链式调用）
func (dm *DataMasker) WithMaxBodySize(size int) *DataMasker {
	dm.config.MaxBodySize = size
	return dm
}

// MaskString 脱敏字符串（链式调用）
func (dm *DataMasker) MaskString(data string) string {
	return dm.Mask([]byte(data))
}

// MaskBytes 脱敏字节数组（链式调用）
func (dm *DataMasker) MaskBytes(data []byte) string {
	return dm.Mask(data)
}
