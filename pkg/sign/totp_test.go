/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 15:15:24
 * @FilePath: \go-toolbox\pkg\sign\totp_test.go
 * @Description: TOTP 单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTOTPConfig(t *testing.T) {
	config := DefaultTOTPConfig()
	assert.Equal(t, 6, config.Digits)
	assert.Equal(t, 30, config.Period)
	assert.Equal(t, 1, config.Skew)
	assert.Equal(t, "SHA1", config.Algorithm)
}

func TestGenerateTOTPSecret(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	assert.NotEmpty(t, secret)
	assert.True(t, len(secret) > 0)

	secret2 := GenerateTOTPSecret(20)
	assert.NotEqual(t, secret, secret2)

	secretDefault := GenerateTOTPSecret(0)
	assert.NotEmpty(t, secretDefault)

	secretNeg := GenerateTOTPSecret(-1)
	assert.NotEmpty(t, secretNeg)
}

func TestGenerateTOTPSecretLength(t *testing.T) {
	secret10 := GenerateTOTPSecret(10)
	secret20 := GenerateTOTPSecret(20)
	secret32 := GenerateTOTPSecret(32)

	assert.NotEmpty(t, secret10)
	assert.NotEmpty(t, secret20)
	assert.NotEmpty(t, secret32)

	assert.True(t, len(secret20) > len(secret10))
	assert.True(t, len(secret32) > len(secret20))
}

func TestGenerateTOTPURI(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	uri := GenerateTOTPURI(secret, "user@example.com", "TestApp", nil)

	assert.True(t, strings.HasPrefix(uri, "otpauth://totp/"))
	assert.Contains(t, uri, "TestApp:user@example.com")
	assert.Contains(t, uri, "secret="+secret)
	assert.Contains(t, uri, "issuer=TestApp")
	assert.Contains(t, uri, "algorithm=SHA1")
	assert.Contains(t, uri, "digits=6")
	assert.Contains(t, uri, "period=30")
}

func TestGenerateTOTPURIWithCustomConfig(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := &TOTPConfig{
		Digits:    8,
		Period:    60,
		Skew:      2,
		Algorithm: "SHA1",
	}
	uri := GenerateTOTPURI(secret, "user@example.com", "TestApp", config)

	assert.Contains(t, uri, "digits=8")
	assert.Contains(t, uri, "period=60")
}

func TestGenerateTOTPCode(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := DefaultTOTPConfig()

	code, err := GenerateTOTPCode(secret, config)
	assert.NoError(t, err)
	assert.Len(t, code, 6)
	assert.True(t, isNumeric(code))
}

func TestGenerateTOTPCode8Digits(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := &TOTPConfig{
		Digits:    8,
		Period:    30,
		Skew:      1,
		Algorithm: "SHA1",
	}

	code, err := GenerateTOTPCode(secret, config)
	assert.NoError(t, err)
	assert.Len(t, code, 8)
	assert.True(t, isNumeric(code))
}

func TestGenerateTOTPCodeNilConfig(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	code, err := GenerateTOTPCode(secret, nil)
	assert.NoError(t, err)
	assert.Len(t, code, 6)
	assert.True(t, isNumeric(code))
}

func TestGenerateTOTPCodeInvalidSecret(t *testing.T) {
	code, err := GenerateTOTPCode("!!!invalid-base32!!!", nil)
	assert.Error(t, err)
	assert.Empty(t, code)
}

func TestValidateTOTPCode(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := DefaultTOTPConfig()

	code, err := GenerateTOTPCode(secret, config)
	assert.NoError(t, err)

	valid := ValidateTOTPCode(secret, code, config)
	assert.True(t, valid)
}

func TestValidateTOTPCodeNilConfig(t *testing.T) {
	secret := GenerateTOTPSecret(20)

	code, err := GenerateTOTPCode(secret, nil)
	assert.NoError(t, err)

	valid := ValidateTOTPCode(secret, code, nil)
	assert.True(t, valid)
}

func TestValidateTOTPCodeWrongCode(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := DefaultTOTPConfig()

	valid := ValidateTOTPCode(secret, "000000", config)
	assert.False(t, valid)
}

func TestValidateTOTPCodeEmptySecret(t *testing.T) {
	valid := ValidateTOTPCode("", "123456", nil)
	assert.False(t, valid)
}

func TestValidateTOTPCodeEmptyCode(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	valid := ValidateTOTPCode(secret, "", nil)
	assert.False(t, valid)
}

func TestValidateTOTPCodeInvalidSecret(t *testing.T) {
	valid := ValidateTOTPCode("!!!invalid!!!", "123456", nil)
	assert.False(t, valid)
}

func TestValidateTOTPCodeWithSkew(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := &TOTPConfig{
		Digits:    6,
		Period:    30,
		Skew:      2,
		Algorithm: "SHA1",
	}

	code, err := GenerateTOTPCode(secret, config)
	assert.NoError(t, err)

	valid := ValidateTOTPCode(secret, code, config)
	assert.True(t, valid)
}

func TestValidateTOTPCodeDifferentSecrets(t *testing.T) {
	secret1 := GenerateTOTPSecret(20)
	secret2 := GenerateTOTPSecret(20)

	code1, err := GenerateTOTPCode(secret1, nil)
	assert.NoError(t, err)

	valid := ValidateTOTPCode(secret2, code1, nil)
	assert.False(t, valid)
}

func TestGenerateBackupCodes(t *testing.T) {
	codes := GenerateBackupCodes(10)
	assert.Len(t, codes, 10)

	for _, code := range codes {
		assert.True(t, len(code) > 0)
		assert.True(t, len(code) <= 8)
		assert.True(t, isHex(code))
	}
}

func TestGenerateBackupCodesZero(t *testing.T) {
	codes := GenerateBackupCodes(0)
	assert.Len(t, codes, 0)
}

func TestGenerateBackupCodesUniqueness(t *testing.T) {
	codes := GenerateBackupCodes(100)
	seen := make(map[string]bool)
	for _, code := range codes {
		assert.False(t, seen[code], "Backup codes should be unique, but found duplicate: %s", code)
		seen[code] = true
	}
}

func TestGenerateBackupCodesSingle(t *testing.T) {
	codes := GenerateBackupCodes(1)
	assert.Len(t, codes, 1)
	assert.NotEmpty(t, codes[0])
}

func TestGenerateAndValidateFlow(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := DefaultTOTPConfig()

	code, err := GenerateTOTPCode(secret, config)
	assert.NoError(t, err)
	assert.Len(t, code, 6)

	valid := ValidateTOTPCode(secret, code, config)
	assert.True(t, valid)

	wrongCode := "999999"
	invalidValid := ValidateTOTPCode(secret, wrongCode, config)
	assert.False(t, invalidValid)
}

func TestGenerateTOTPCodeConsistency(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := DefaultTOTPConfig()

	code1, err1 := GenerateTOTPCode(secret, config)
	code2, err2 := GenerateTOTPCode(secret, config)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, code1, code2)
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func TestValidateTOTPCodeWithSpaces(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := DefaultTOTPConfig()

	code, err := GenerateTOTPCode(secret, config)
	assert.NoError(t, err)

	// 模拟 Google Authenticator 显示格式 "123 456"
	spacedCode := code[:3] + " " + code[3:]
	valid := ValidateTOTPCode(secret, spacedCode, config)
	assert.True(t, valid, "应该支持带空格的验证码格式")
}

func TestValidateTOTPCodeOnlySpaces(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	valid := ValidateTOTPCode(secret, "   ", nil)
	assert.False(t, valid)
}

func TestValidateTOTPCodePeriodZero(t *testing.T) {
	secret := GenerateTOTPSecret(20)
	config := &TOTPConfig{
		Digits:    6,
		Period:    0,
		Skew:      1,
		Algorithm: "SHA1",
	}

	// 不应 panic，应回退到默认 Period=30
	code, err := GenerateTOTPCode(secret, config)
	assert.NoError(t, err)
	assert.Len(t, code, 6)

	valid := ValidateTOTPCode(secret, code, config)
	assert.True(t, valid)
}

func TestConsumeBackupCode(t *testing.T) {
	codesJSON := `["ABCD1234","EFGH5678","IJKL9012"]`

	t.Run("消耗存在的码", func(t *testing.T) {
		ok, remaining := ConsumeBackupCode(codesJSON, "EFGH5678")
		assert.True(t, ok)
		assert.NotContains(t, remaining, "EFGH5678")
		assert.Contains(t, remaining, "ABCD1234")
		assert.Contains(t, remaining, "IJKL9012")
	})

	t.Run("消耗不存在的码", func(t *testing.T) {
		ok, remaining := ConsumeBackupCode(codesJSON, "ZZZZ0000")
		assert.False(t, ok)
		assert.Equal(t, codesJSON, remaining)
	})

	t.Run("大小写不敏感", func(t *testing.T) {
		ok, remaining := ConsumeBackupCode(codesJSON, "abcd1234")
		assert.True(t, ok)
		assert.NotContains(t, remaining, "ABCD1234")
	})

	t.Run("空码返回false", func(t *testing.T) {
		ok, _ := ConsumeBackupCode(codesJSON, "")
		assert.False(t, ok)
	})

	t.Run("空JSON返回false", func(t *testing.T) {
		ok, _ := ConsumeBackupCode("", "ABCD1234")
		assert.False(t, ok)
	})

	t.Run("非法JSON返回false", func(t *testing.T) {
		ok, _ := ConsumeBackupCode("not-json", "ABCD1234")
		assert.False(t, ok)
	})

	t.Run("消耗最后一个码返回空数组", func(t *testing.T) {
		ok, remaining := ConsumeBackupCode(`["ONLY1234"]`, "ONLY1234")
		assert.True(t, ok)
		assert.Equal(t, "", remaining)
	})
}
