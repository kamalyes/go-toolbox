/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-28 00:00:00
 * @FilePath: \go-toolbox\pkg\sign\totp.go
 * @Description: TOTP（基于时间的一次性密码）实现
 * 基于RFC 6238算法，支持Google Authenticator等验证器应用
 * 用于双因素认证（2FA）场景
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"
)

// TOTPConfig TOTP配置参数
type TOTPConfig struct {
	Digits    int    // 验证码位数，默认6
	Period    int    // 时间步长（秒），默认30
	Skew      int    // 允许的时间窗口偏移量，默认1（前后各1个窗口）
	Algorithm string // 哈希算法，默认SHA1
}

// DefaultTOTPConfig 返回默认TOTP配置
func DefaultTOTPConfig() *TOTPConfig {
	return &TOTPConfig{
		Digits:    6,
		Period:    30,
		Skew:      1,
		Algorithm: "SHA1",
	}
}

// GenerateTOTPSecret 生成TOTP密钥（Base32编码的随机字节）
// secretLength: 密钥字节长度，推荐20
func GenerateTOTPSecret(secretLength int) string {
	if secretLength <= 0 {
		secretLength = 20
	}
	secret := make([]byte, secretLength)
	if _, err := rand.Read(secret); err != nil {
		secret = []byte(fmt.Sprintf("%d.%d", time.Now().UnixNano(), secretLength))
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)
}

// GenerateTOTPURI 构建TOTP URI（otpauth://totp/...）
// 供Google Authenticator等验证器应用扫描
func GenerateTOTPURI(secret, account, issuer string, config *TOTPConfig) string {
	if config == nil {
		config = DefaultTOTPConfig()
	}
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=%s&digits=%d&period=%d",
		issuer, account, secret, issuer, config.Algorithm, config.Digits, config.Period)
}

// ValidateTOTPCode 验证TOTP验证码
// 基于RFC 6238算法，允许前后Skew个时间窗口
func ValidateTOTPCode(secret, code string, config *TOTPConfig) bool {
	if secret == "" || code == "" {
		return false
	}

	if config == nil {
		config = DefaultTOTPConfig()
	}

	secretUpper := strings.ToUpper(strings.TrimSpace(secret))
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secretUpper)
	if err != nil {
		return false
	}

	now := time.Now().Unix()
	timeStep := int64(config.Period)

	for offset := int64(-int64(config.Skew)); offset <= int64(config.Skew); offset++ {
		counter := (now + offset*timeStep) / timeStep
		expected := generateTOTPCode(key, counter, config.Digits)
		if expected == code {
			return true
		}
	}
	return false
}

// GenerateTOTPCode 根据密钥和当前时间生成TOTP验证码
func GenerateTOTPCode(secret string, config *TOTPConfig) (string, error) {
	if config == nil {
		config = DefaultTOTPConfig()
	}

	secretUpper := strings.ToUpper(strings.TrimSpace(secret))
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secretUpper)
	if err != nil {
		return "", fmt.Errorf("invalid base32 secret: %w", err)
	}

	now := time.Now().Unix()
	counter := now / int64(config.Period)

	return generateTOTPCode(key, counter, config.Digits), nil
}

// generateTOTPCode 根据密钥和时间计数器生成TOTP验证码
func generateTOTPCode(key []byte, counter int64, digits int) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	truncated := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

	code := truncated % uint32(math.Pow10(digits))
	return fmt.Sprintf("%0*d", digits, code)
}

// GenerateBackupCodes 生成指定数量的恢复码
// 每个恢复码为8位随机十六进制字符串
func GenerateBackupCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			b = []byte(fmt.Sprintf("%04d%04d", i, time.Now().Nanosecond()))
		}
		codes[i] = fmt.Sprintf("%X", b)
		if len(codes[i]) > 8 {
			codes[i] = codes[i][:8]
		}
	}
	return codes
}
