/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 15:55:59
 * @FilePath: \go-toolbox\tests\stringx_regexp_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/stretchr/testify/assert"
)

var r *stringx.AnyRegs

func init() {
	r = stringx.NewAnyRegs()
}

// TestRegexpData 结构表示测试数据的格式
type TestRegexpData struct {
	name     string
	input    string
	expected bool
}

func runTests(t *testing.T, TestRegexpData []TestRegexpData, matchFunc func(string) bool) {
	// 对每个测试数据进行迭代
	for _, data := range TestRegexpData {
		// 使用测试数据的名称创建子测试
		t.Run(data.name, func(t *testing.T) {
			// 断言匹配函数处理输入数据后的结果与期望值是否一致
			assert.Equal(t, data.expected, matchFunc(data.input))
		})
	}
}

func TestAllstringxFunctions(t *testing.T) {
	t.Run("TestRegIntOrFloat", TestRegIntOrFloat)
	t.Run("TestAnyRegs_RegNumber", TestAnyRegs_RegNumber)
	t.Run("TestAnyRegs_RegLenNNumber", TestAnyRegs_RegLenNNumber)
	t.Run("TestAnyRegs_RegGeNNumber", TestAnyRegs_RegGeNNumber)
	t.Run("TestAnyRegs_RegMNIntervalNumber", TestAnyRegs_RegMNIntervalNumber)
	t.Run("TestAnyRegs_RegStartingWithNonZero", TestAnyRegs_RegStartingWithNonZero)
	t.Run("TestAnyRegs_RegNNovelsOfRealNumber", TestAnyRegs_RegNNovelsOfRealNumber)
	t.Run("TestAnyRegs_RegMNNovelsOfRealNumber", TestAnyRegs_RegMNNovelsOfRealNumber)
	t.Run("TestAnyRegs_RegNanZeroNumber", TestAnyRegs_RegNanZeroNumber)
	t.Run("TestAnyRegs_RegMatchNanZeroNegNumber", TestAnyRegs_RegMatchNanZeroNegNumber)
	t.Run("TestAnyRegs_RegNanZeroNegNumber", TestAnyRegs_RegNLeCharacter)
	t.Run("TestAnyRegs_RegNLeCharacter", TestAnyRegs_RegNLeCharacter)
	t.Run("TestAnyRegs_RegEnCharacter", TestAnyRegs_RegEnCharacter)
	t.Run("TestAnyRegs_RegUpEnCharacter", TestAnyRegs_RegUpEnCharacter)
	t.Run("TestAnyRegs_RegLowerEnCharacter", TestAnyRegs_RegLowerEnCharacter)
	t.Run("TestAnyRegs_RegNumberEnCharacter", TestAnyRegs_RegNumberEnCharacter)
	t.Run("TestAnyRegs_RegNumberEnUnderscores", TestAnyRegs_RegNumberEnUnderscores)
	t.Run("TestAnyRegs_RegPass1", TestAnyRegs_RegPass1)
	t.Run("TestAnyRegs_RegIsContainSpecialCharacter", TestAnyRegs_RegIsContainSpecialCharacter)
	t.Run("TestAnyRegs_RegEmail", TestAnyRegs_RegEmail)
	t.Run("TestAnyRegs_RegChinePhoneNumber", TestAnyRegs_RegChinePhoneNumber)
	t.Run("TestAnyRegs_RegContainChineseCharacter", TestAnyRegs_RegContainChineseCharacter)
	t.Run("TestAnyRegs_MatchDoubleByte", TestAnyRegs_MatchDoubleByte)
	t.Run("TestAnyRegs_MatchEmptyLine", TestAnyRegs_MatchEmptyLine)
	t.Run("TestAnyRegs_MatchIPv4", TestAnyRegs_MatchIPv4)
	t.Run("TestAnyRegs_MatchIPv6", TestAnyRegs_MatchIPv6)
}

func TestRegIntOrFloat(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with integer input": {
			Input:    "123",
			Expected: true,
		},
		"Test with float input": {
			Input:    "12.34",
			Expected: true,
		},
		"Test with non-numeric input": {
			Input:    "abc",
			Expected: false,
		},
		"Test with special characters": {
			Input:    "12@34",
			Expected: false,
		},
		// Add more test cases as needed
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchIntOrFloat(test.Input)
			assert.Equal(t, test.Expected, actual, fmt.Sprintf("Expected: %v, Got: %v for input: %s\n", test.Expected, actual, test.Input))
		})
	}
}

func TestAnyRegs_RegNumber(t *testing.T) {

	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with pure numeric input": {
			Input:    "123456",
			Expected: true,
		},
		"Test with alphanumeric input": {
			Input:    "abc123",
			Expected: false,
		},
		"Test with special characters": {
			Input:    "12@34",
			Expected: false,
		},
		// Add more test cases as needed
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchNumber(test.Input)
			assert.Equal(t, test.Expected, actual, fmt.Sprintf("Expected: %v, Got: %v for input: %s\n", test.Expected, actual, test.Input))
		})
	}
}

func TestAnyRegs_RegLenNNumber(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Length   int
		Expected bool
	}{
		"Test with valid input of length 5": {
			Input:    "12345",
			Length:   5,
			Expected: true,
		},
		"Test with invalid input of incorrect length": {
			Input:    "4567",
			Length:   5,
			Expected: false,
		},
		"Invalid case - empty string": {
			Input:    "",
			Length:   0,
			Expected: true,
		},
		"Invalid case - non-numeric string": {
			Input:    "abc123",
			Length:   0,
			Expected: false,
		},
		// Add more test cases as needed
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchLenNNumber(test.Input, test.Length)

			assert.Equal(t, test.Expected, actual, fmt.Sprintf("Expected: %v, Got: %v for input: %s\n", test.Expected, actual, test.Input))
		})
	}
}

func TestAnyRegs_RegGeNNumber(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Length   int
		Expected bool
	}{
		"Test with valid input of length 5": {
			Input:    "12345",
			Length:   8,
			Expected: false,
		},
		"Test with invalid input of incorrect length": {
			Input:    "4558567",
			Length:   5,
			Expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchGeNNumber(test.Input, test.Length)

			assert.Equal(t, test.Expected, actual, fmt.Sprintf("Expected: %v, Got: %v for input: %s\n", test.Expected, actual, test.Input))
		})
	}
}

func TestAnyRegs_RegMNIntervalNumber(t *testing.T) {
	tests := []struct {
		Name     string
		Number   string
		Min, Max int
		Expected bool
	}{
		{"Valid case - within interval", "1234567890", 4, 12, true},
		{"Invalid case - too few characters", "123", 4, 12, false},
		{"Invalid case - too many characters", "123456789012345", 4, 12, false},
		{"Valid case - edge case lower bound", "1234", 4, 12, true},
		{"Valid case - edge case upper bound", "123456789012", 4, 12, true},
		{"Edge case - empty string", "", 4, 12, false},
		{"Edge case - interval is just 1 character", "1", 1, 1, true},
		{"Edge case - negative lower bound", "-123", -4, 12, false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if result := r.MatchMNIntervalNumber(test.Number, test.Min, test.Max); result != test.Expected {
				t.Errorf("Test %s failed. Expected %t, got %t", test.Name, test.Expected, result)
			}
		})
	}
}

func TestAnyRegs_RegStartingWithNonZero(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Valid case - string starts with non-zero digit": {
			Input:    "1234",
			Expected: true,
		},
		"Invalid case - string starts with zero": {
			Input:    "0123",
			Expected: false,
		},
		"Invalid case - empty string": {
			Input:    "",
			Expected: false,
		},
		"Invalid case - non-numeric string": {
			Input:    "abc123",
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchStartingWithNonZero(test.Input)
			assert.Equal(t, test.Expected, actual, "Expected input '%s' to return %t", test.Input, test.Expected)
		})
	}
}
func TestAnyRegs_RegNNovelsOfRealNumber(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Novels   int
		Expected bool
	}{
		"Valid case - exact number of novels": {
			Input:    "123",
			Novels:   3,
			Expected: true,
		},
		"Valid case - more novels than required": {
			Input:    "12.345",
			Novels:   3,
			Expected: true,
		},
		"Invalid case - fewer novels than required": {
			Input:    "1.23",
			Novels:   3,
			Expected: false,
		},
		"Invalid case - input contains non-numeric characters": {
			Input:    "12a34",
			Novels:   3,
			Expected: false,
		},
		"Invalid case - empty input": {
			Input:    "",
			Novels:   3,
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchNNovelsOfRealNumber(test.Input, test.Novels)
			assert.Equal(t, test.Expected, actual, "Expected input '%s' to return %t", test.Input, test.Expected)
		})
	}
}

func TestAnyRegs_RegMNNovelsOfRealNumber(t *testing.T) {
	tests := map[string]struct {
		Input     string
		MinNovels int
		MaxNovels int
		Expected  bool
	}{
		"Valid case - within the range of min and max novels": {
			Input:     "12.345",
			MinNovels: 2,
			MaxNovels: 6,
			Expected:  true,
		},
		"Invalid case - fewer novels than the minimum required": {
			Input:     "1.2",
			MinNovels: 2,
			MaxNovels: 6,
			Expected:  false,
		},
		"Invalid case - more novels than the maximum allowed": {
			Input:     "12.34567",
			MinNovels: 2,
			MaxNovels: 4,
			Expected:  false,
		},
		"Invalid case - input contains non-numeric characters": {
			Input:     "12.3a",
			MinNovels: 2,
			MaxNovels: 6,
			Expected:  false,
		},
		"Invalid case - empty input": {
			Input:     "",
			MinNovels: 2,
			MaxNovels: 6,
			Expected:  false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchMNNovelsOfRealNumber(test.Input, test.MinNovels, test.MaxNovels)
			assert.Equal(t, test.Expected, actual, "Expected input '%s' to return %t", test.Input, test.Expected)
		})
	}
}

func TestAnyRegs_RegNanZeroNumber(t *testing.T) {
	tests := []TestRegexpData{
		{"MatchNaN", "NaN", false},
		{"MatchInf", "+Inf", false},
		{"MatchPositiveZero", "+0", false},
		{"MatchNegativeZero", "-0", false},
		{"Zero", "0", false},
		{"Integer", "1", true},
	}

	runTests(t, tests, r.MatchNanZeroNumber)
}

func TestAnyRegs_RegMatchNanZeroNegNumber(t *testing.T) {
	tests := []TestRegexpData{
		{"MatchNaN", "NaN", false},
		{"MatchInf", "+Inf", false},
		{"MatchPositiveZero", "+0", false},
		{"MatchNegativeZero", "-0", false},
		{"Zero", "0", false},
		{"Minus", "-1", true},
	}
	runTests(t, tests, r.MatchNanZeroNegNumber)
}

func TestAnyRegs_RegNLeCharacter(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Length   int
		Expected bool
	}{
		"Test with valid input of length 5": {
			Input:    "Hello",
			Length:   5,
			Expected: true,
		},
		"Test with valid input of length 3": {
			Input:    "abc",
			Length:   3,
			Expected: true,
		},
		"Test with valid input of length 7": {
			Input:    "testing",
			Length:   7,
			Expected: true,
		},
		"Test with input containing special characters": {
			Input:    "Hello!",
			Length:   6,
			Expected: true,
		},
		// Add more test cases as needed
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchNLeCharacter(test.Input, test.Length)

			assert.Equal(t, test.Expected, actual, fmt.Sprintf("Expected: %v, Got: %v for input: %s\n", test.Expected, actual, test.Input))
		})
	}
}

func TestAnyRegs_RegEnCharacter(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with valid input of all uppercase letters": {
			Input:    "HELLO",
			Expected: true,
		},
		"Test with valid input of all lowercase letters": {
			Input:    "hello",
			Expected: true,
		},
		"Test with valid input of mixed case letters": {
			Input:    "HeLlO",
			Expected: true,
		},
		"Test with valid input of numbers and special characters": {
			Input:    "123!",
			Expected: false,
		},
		"Test with valid input of special characters": {
			Input:    "@BCD",
			Expected: false,
		},
		"Test with valid input of empty string": {
			Input:    "",
			Expected: false,
		},
		"Test with valid input of mixed Chinese and English characters": {
			Input:    "中abc",
			Expected: false,
		},
		"Test with valid input of Chinese punctuation": {
			Input:    "，。！？",
			Expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchEnCharacter(test.Input)

			assert.Equal(t, test.Expected, actual, fmt.Sprintf("Expected: %v, Got: %v for input: %s\n", test.Expected, actual, test.Input))
		})
	}
}

func TestAnyRegs_RegUpEnCharacter(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with valid input of upper case English letters": {
			Input:    "abc",
			Expected: false,
		},
		"Test with valid input of lower case English letters": {
			Input:    "ABC",
			Expected: true,
		},
		"Test with valid input of mixed case English letters": {
			Input:    "AbC",
			Expected: false,
		},
		"Test with valid input of special characters": {
			Input:    "@BCD",
			Expected: false,
		},
		"Test with valid input of empty string": {
			Input:    "",
			Expected: false,
		},
		"Test with valid input of mixed Chinese and English characters": {
			Input:    "中abc",
			Expected: false,
		},
		"Test with valid input of Chinese punctuation": {
			Input:    "，。！？",
			Expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchUpEnCharacter(test.Input)
			assert.Equal(t, test.Expected, actual, fmt.Sprintf("Expected: %v, Got: %v for input: %s\n", test.Expected, actual, test.Input))
		})
	}
}

func TestAnyRegs_RegLowerEnCharacter(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with valid input of upper case English letters": {
			Input:    "ABC",
			Expected: false,
		},
		"Test with valid input of lower case English letters": {
			Input:    "abc",
			Expected: true,
		},
		"Test with valid input of mixed case English letters": {
			Input:    "AbC",
			Expected: false,
		},
		"Test with valid input of special characters": {
			Input:    "@BCD",
			Expected: false,
		},
		"Test with valid input of empty string": {
			Input:    "",
			Expected: false,
		},
		"Test with valid input of mixed Chinese and English characters": {
			Input:    "中abc",
			Expected: false,
		},
		"Test with valid input of Chinese punctuation": {
			Input:    "，。！？",
			Expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchLowerEnCharacter(test.Input)
			assert.Equal(t, test.Expected, actual, fmt.Sprintf("Expected: %v, Got: %v for input: %s\n", test.Expected, actual, test.Input))
		})
	}
}

func TestAnyRegs_RegNumberEnCharacter(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with valid input containing numbers and English letters": {
			Input:    "abc123DEF",
			Expected: true,
		},
		"Test with valid input containing only numbers": {
			Input:    "123456",
			Expected: true,
		},
		"Test with valid input containing only English letters": {
			Input:    "abcDEF",
			Expected: true,
		},
		"Test with invalid input containing special characters": {
			Input:    "abc@123",
			Expected: false,
		},
		"Test with invalid input containing whitespace": {
			Input:    "abc def",
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchNumberEnCharacter(test.Input)

			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("'%s' is a valid string containing numbers and English letters.\n", test.Input))
			} else {
				assert.False(t, actual, fmt.Sprintf("'%s' is an invalid string containing numbers and English letters.\n", test.Input))
			}
		})
	}
}

func TestAnyRegs_RegNumberEnUnderscores(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with valid input containing numbers and English letters": {
			Input:    "abc123_DEF",
			Expected: true,
		},
		"Test with valid input containing only numbers": {
			Input:    "123456",
			Expected: true,
		},
		"Test with valid input containing only English letters": {
			Input:    "abcDEF",
			Expected: true,
		},
		"Test with invalid input containing special characters": {
			Input:    "abc@123",
			Expected: false,
		},
		"Test with invalid input containing whitespace": {
			Input:    "abc def",
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchNumberEnUnderscores(test.Input)

			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("'%s' is a valid string containing numbers and English letters.\n", test.Input))
			} else {
				assert.False(t, actual, fmt.Sprintf("'%s' is an invalid string containing numbers and English letters.\n", test.Input))
			}
		})
	}
}

func TestAnyRegs_RegPass1(t *testing.T) {
	tests := map[string]struct {
		Input    string
		MinLen   int
		MaxLen   int
		Expected bool
	}{
		"Test with valid password of minimum length": {
			Input:    "a_abc12",
			MinLen:   4,
			MaxLen:   12,
			Expected: true,
		},
		"Test with valid password of maximum length": {
			Input:    "abcde_12345678",
			MinLen:   6,
			MaxLen:   16,
			Expected: true,
		},
		"Test with password containing invalid characters": {
			Input:    "abc@123",
			MinLen:   4,
			MaxLen:   10,
			Expected: false,
		},
		"Test with too short password": {
			Input:    "abc1",
			MinLen:   5,
			MaxLen:   10,
			Expected: false,
		},
		"Test with too long password": {
			Input:    "abcdefg_1234567890",
			MinLen:   5,
			MaxLen:   15,
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchPass1(test.Input, test.MinLen, test.MaxLen)

			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("'%s' is a valid password within the specified length range.\n", test.Input))
			} else {
				assert.False(t, actual, fmt.Sprintf("'%s' is an invalid password within the specified length range.\n", test.Input))
			}
		})
	}
}

func TestAnyRegs_RegIsContainSpecialCharacter(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with special characters": {
			Input:    "Hello, World!",
			Expected: true,
		},
		"Test with alphanumeric characters": {
			Input:    "GoLang123",
			Expected: false,
		},
		"Test with special characters and symbols": {
			Input:    "Password!@#",
			Expected: true,
		},
		"Test with special characters and underscores": {
			Input:    "Special_Chars%^&*",
			Expected: true,
		},
		"Test with normal string": {
			Input:    "NormalString",
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchIsContainSpecialCharacter(test.Input)

			if test.Expected {
				assert.True(t, actual, "Expected input '%s' to contain special characters", test.Input)
			} else {
				assert.False(t, actual, "Expected input '%s' to not contain special characters", test.Input)
			}
		})
	}
}

func TestAnyRegs_RegChineseCharacter(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with English characters": {
			Input:    "Hello",
			Expected: false,
		},
		"Test with mixed Chinese and English characters": {
			Input:    "你好, World",
			Expected: false,
		},
		"Test with Chinese characters and numbers": {
			Input:    "测试123",
			Expected: false,
		},
		"Test with single Chinese character": {
			Input:    "汉字",
			Expected: true,
		},
		"Test with Chinese characters and special characters": {
			Input:    "你好！",
			Expected: false,
		},
		"Test with only special characters": {
			Input:    "!@#$",
			Expected: false,
		},
		"Test with all numbers": {
			Input:    "12345",
			Expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, _ := stringx.IsChineseCharacter(test.Input)
			if test.Expected {
				assert.True(t, actual, "Expected input '%s' to Chinese characters", test.Input)
			} else {
				assert.False(t, actual, "Expected input '%s' to  Chinese characters", test.Input)
			}
		})
	}
}
func TestAnyRegs_RegEmail(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with a valid email": {
			Input:    "aabc@qq.com",
			Expected: true,
		},
		"Test with an invalid email": {
			Input:    "invalidemail@",
			Expected: false,
		},
		"Test with an email without domain": {
			Input:    "example@",
			Expected: false,
		},
		"Test with an email without username": {
			Input:    "@domain.com",
			Expected: false,
		},
		"Test with multiple domains": {
			Input:    "user@example.co.uk",
			Expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchEmail(test.Input)
			if test.Expected {
				assert.True(t, actual, "Expected input '%s' to be a valid email", test.Input)
			} else {
				assert.False(t, actual, "Expected input '%s' to be an invalid email", test.Input)
			}
		})
	}
}

func TestAnyRegs_RegChinePhoneNumber(t *testing.T) {
	tests := map[string]struct {
		Number   string
		Expected bool
	}{
		"Valid phone number 1": {
			Number:   "13800138000",
			Expected: true,
		},
		"Invalid phone number 1": {
			Number:   "12345678901",
			Expected: false,
		},
		"Valid phone number 2": {
			Number:   "19912345678",
			Expected: true,
		},
		"Invalid phone number 2": {
			Number:   "10000000000",
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchChinesePhoneNumber(test.Number)

			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("%s is a valid phone number.\n", test.Number))
			} else {
				assert.False(t, actual, fmt.Sprintf("%s is an invalid phone number.\n", test.Number))
			}
		})
	}
}

func TestAnyRegs_RegContainChineseCharacter(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with Chinese characters in the string": {
			Input:    "Hello, 世界",
			Expected: true,
		},
		"Test with all Chinese characters in the string": {
			Input:    "你好",
			Expected: true,
		},
		"Test with no Chinese characters in the string": {
			Input:    "Hello, World!",
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchContainChineseCharacter(test.Input)

			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("'%s' contains Chinese characters.\n", test.Input))
			} else {
				assert.False(t, actual, fmt.Sprintf("'%s' does not contain Chinese characters.\n", test.Input))
			}
		})
	}
}

func TestAnyRegs_MatchDoubleByte(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with double-byte characters in the string": {
			Input:    "Hello, 世界",
			Expected: true,
		},
		"Test with all double-byte characters in the string": {
			Input:    "你好",
			Expected: true,
		},
		"Test with no double-byte characters in the string": {
			Input:    "Hello, World!",
			Expected: false,
		},
		"Test with Japanese Hiragana characters in the string": {
			Input:    "こんにちは",
			Expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchDoubleByte(test.Input)

			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("'%s' contains double-byte characters.\n", test.Input))
			} else {
				assert.False(t, actual, fmt.Sprintf("'%s' does not contain double-byte characters.\n", test.Input))
			}
		})
	}
}

func TestAnyRegs_MatchEmptyLine(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Empty line with spaces": {
			Input:    "   \n",
			Expected: true,
		},
		"Empty line with tabs": {
			Input:    "\t\t\n",
			Expected: true,
		},
		"Empty line with mixed spaces and tabs": {
			Input:    "  \t \t\n",
			Expected: true,
		},
		"Non-empty new line": {
			Input:    "empty newline\n",
			Expected: true,
		},
		"Multiple empty lines": {
			Input:    "\n\n\n",
			Expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchEmptyLine(test.Input)

			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("The string contains empty lines:\n'%s'\n", test.Input))
			} else {
				assert.False(t, actual, fmt.Sprintf("The string does not contain empty lines:\n'%s'\n", test.Input))
			}
		})
	}
}

func TestAnyRegs_MatchIPv4(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with a valid IPv4 address": {
			Input:    "192.168.1.1",
			Expected: true,
		},
		"Test with another valid IPv4 address": {
			Input:    "255.255.255.255",
			Expected: true,
		},
		"Test with one more valid IPv4 address": {
			Input:    "0.0.0.0",
			Expected: true,
		},
		"Test with an invalid IPv4 address": {
			Input:    "256.256.256.256",
			Expected: false,
		},
		"Test with another invalid IPv4 address": {
			Input:    "192.168.1.256",
			Expected: false,
		},
		"Test with an incomplete IPv4 address": {
			Input:    "192.168.1",
			Expected: false,
		},
		"Test with a non-numeric IPv4 address": {
			Input:    "abc.def.ghi.jkl",
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchIPv4(test.Input)
			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("%s is a valid IPv4 address.\n", test.Input))
			} else {
				assert.False(t, actual, fmt.Sprintf("%s is an invalid IPv4 address.\n", test.Input))
			}
		})
	}
}

func TestAnyRegs_MatchIPv6(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected bool
	}{
		"Test with a valid long IPv6 address": {
			Input:    "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			Expected: true,
		},
		"Test with a valid shortened IPv6 address": {
			Input:    "2001:db8:85a3::8a2e:370:7334",
			Expected: true,
		},
		"Test with a valid loopback IPv6 address": {
			Input:    "::1",
			Expected: true,
		},
		"Test with another valid IPv6 address": {
			Input:    "fe80::1ff:fe23:4567:890a",
			Expected: true,
		},
		"Test with an invalid IPv4 address": {
			Input:    "192.168.1.1",
			Expected: false,
		},
		"Test with an invalid character in IPv6 address": {
			Input:    "2001:db8:85a3:0:0:8a2e:370g:7334",
			Expected: false,
		},
		"Test with double '::' in IPv6 address": {
			Input:    "2001::85a3::8a2e:370:7334",
			Expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := r.MatchIPv6(test.Input)

			if test.Expected {
				assert.True(t, actual, fmt.Sprintf("%s is a valid IPv6 address.\n", test.Input))
			} else {
				assert.False(t, actual, fmt.Sprintf("%s is an invalid IPv6 address.\n", test.Input))
			}
		})
	}
}

// TestMatchPass2Valid 测试 MatchPass2 的有效情况
func TestMatchPass2Valid(t *testing.T) {
	input := "Valid1@Password"
	result := r.MatchPass2(input)
	assert.True(t, result, "Expected true for valid password")
}

// TestMatchPass2Invalid 测试 MatchPass2 的无效情况
func TestMatchPass2Invalid(t *testing.T) {
	input := "short"
	result := r.MatchPass2(input)
	assert.False(t, result, "Expected false for invalid password")
}

func TestParseWeek(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Weekday
		err      bool
	}{
		{"M", time.Monday, false},
		{"T", time.Tuesday, false},
		{"W", time.Wednesday, false},
		{"R", time.Thursday, false},
		{"F", time.Friday, false},
		{"S", time.Saturday, false},
		{"U", time.Sunday, false},
		{"MON", time.Monday, false},
		{"TUE", time.Tuesday, false},
		{"WED", time.Wednesday, false},
		{"THU", time.Thursday, false},
		{"FRI", time.Friday, false},
		{"SAT", time.Saturday, false},
		{"SUN", time.Sunday, false},
		{"Monday", time.Monday, false},
		{"Tuesday", time.Tuesday, false},
		{"Wednesday", time.Wednesday, false},
		{"Thursday", time.Thursday, false},
		{"Friday", time.Friday, false},
		{"Saturday", time.Saturday, false},
		{"Sunday", time.Sunday, false},
		{"1", time.Monday, false},
		{"2", time.Tuesday, false},
		{"3", time.Wednesday, false},
		{"4", time.Thursday, false},
		{"5", time.Friday, false},
		{"6", time.Saturday, false},
		{"0", time.Sunday, false},
		{"7", 0, true}, // 超出范围
		{"Invalid", 0, true},
		{"  Tue  ", time.Tuesday, false}, // 测试空格
	}

	for _, test := range tests {
		result, err := stringx.ParseWeek(test.input)
		if test.err {
			assert.Error(t, err, "ParseWeek Expected an error for input: %q", test.input)
		} else {
			assert.NoError(t, err, "ParseWeek Did not expect an error for input: %q", test.input)
			assert.Equal(t, test.expected, result, "ParseWeek Expected %v for input %q", test.expected, test.input)
		}
	}
}

func TestIsValidMonth(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Month
		err      bool
	}{
		{"JAN", time.January, false},
		{"FEB", time.February, false},
		{"MAR", time.March, false},
		{"APR", time.April, false},
		{"MAY", time.May, false},
		{"JUN", time.June, false},
		{"JUL", time.July, false},
		{"AUG", time.August, false},
		{"SEP", time.September, false},
		{"OCT", time.October, false},
		{"NOV", time.November, false},
		{"DEC", time.December, false},
		{"0", 0, true},
		{"1", time.January, false},
		{"2", time.February, false},
		{"3", time.March, false},
		{"4", time.April, false},
		{"5", time.May, false},
		{"6", time.June, false},
		{"7", time.July, false},
		{"8", time.August, false},
		{"9", time.September, false},
		{"10", time.October, false},
		{"11", time.November, false},
		{"12", time.December, false},
		{"Invalid", 0, true},
		{"January", time.January, false},
		{"February", time.February, false},
		{"  JUL  ", time.July, false}, // 测试空格
		{"J", time.January, false},    // 单字母缩写
		{"F", time.February, false},   // 单字母缩写
		{"M", time.March, false},      // 单字母缩写
		{"A", time.April, false},      // 单字母缩写
		{"Y", time.May, false},        // 单字母缩写
		{"N", time.June, false},       // 单字母缩写
		{"L", time.July, false},       // 单字母缩写
		{"G", time.August, false},     // 单字母缩写
		{"S", time.September, false},  // 单字母缩写
		{"T", time.October, false},    // 单字母缩写
		{"V", time.November, false},   // 单字母缩写
		{"C", time.December, false},   // 单字母缩写
		{"Invalid", 0, true},
	}

	for _, test := range tests {
		result, err := stringx.ParseMonth(test.input)
		if test.err {
			assert.Error(t, err, "ParseMonth Expected an error for input: %q", test.input)
		} else {
			assert.NoError(t, err, "ParseMonth Did not expect an error for input: %q", test.input)
			assert.Equal(t, test.expected, result, "ParseMonth Expected %v for input %q", test.expected, test.input)
		}
	}
}
