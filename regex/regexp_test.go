/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 16:29:12
 * @FilePath: \go-toolbox\regex\regexp_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package regex

import (
	"fmt"
	"testing"
)

var r *AnyRegs

func init() {
	r = NewAnyRegs()
}

func TestRegIntOrFloat(t *testing.T) {
	//var a = "1234你好"
	var a = "1234.45"
	fmt.Println(r.MatchIntOrFloat(a))
}

func TestAnyRegs_RegNumber(t *testing.T) {
	//var a = "1234你好"
	//var a = "1234.45"
	var a = "1234.4$5"
	fmt.Println(r.MatchIntOrFloat(a))
}

func TestAnyRegs_RegLenNNumber(t *testing.T) {
	//var a = "1234.45"
	var a = "123445"
	fmt.Println(r.MatchLenNNumber(a, 6))
}

func TestAnyRegs_RegGeNNumber(t *testing.T) {
	var a = "123446"
	fmt.Println(r.MatchLenNNumber(a, 6))
}

func TestAnyRegs_RegMNIntervalNumber(t *testing.T) {
	var a = "12312211121212"
	fmt.Println(r.MatchMNIntervalNumber(a, 4, 12))
}

func TestAnyRegs_RegStartingWithNonZero(t *testing.T) {
	var a = "00123445"
	fmt.Println(r.MatchStartingWithNonZero(a))
}

func TestAnyRegs_RegNNovelsOfRealNumber(t *testing.T) {
	//var a = "12.00"
	var a = "12"
	fmt.Println(r.MatchNNovelsOfRealNumber(a, 3))
}

func TestAnyRegs_RegMNNovelsOfRealNumber(t *testing.T) {
	var a = "12.1222"
	//var a = "12"
	fmt.Println(r.MatchMNNovelsOfRealNumber(a, 2, 6))
}

func TestAnyRegs_RegNanZeroNumber(t *testing.T) {
	var a = "-1"
	//var a = "12"
	fmt.Println(r.MatchNanZeroNumber(a))
}

func TestAnyRegs_RegNanZeroNegNumber(t *testing.T) {
	var a = "-1"
	//var a = "12"
	fmt.Println(r.MatchNanZeroNegNumber(a))
}

func TestAnyRegs_RegNLeCharacter(t *testing.T) {
	//var a = "-1)1"
	var a = "小泽玛"
	fmt.Println(r.MatchNLeCharacter(a, 3))
}

func TestAnyRegs_RegEnCharacter(t *testing.T) {
	//var a = "-1)1"
	var a = "abc0"
	fmt.Println(r.MatchEnCharacter(a))
}

func TestAnyRegs_RegUpEnCharacter(t *testing.T) {
	var a = "ABC"
	fmt.Println(r.MatchUpEnCharacter(a))
}

func TestAnyRegs_RegLowerEnCharacter(t *testing.T) {
	var a = "abc"
	fmt.Println(r.MatchLowerEnCharacter(a))
}

func TestAnyRegs_RegNumberEnCharacter(t *testing.T) {
	var a = "abc12?"
	fmt.Println(r.MatchNumberEnCharacter(a))
}

func TestAnyRegs_RegNumberEnUnderscores(t *testing.T) {
	var a = "abc12_"
	//var a = "_abc12"
	fmt.Println(r.MatchNumberEnUnderscores(a))
}

func TestAnyRegs_RegPass1(t *testing.T) {
	//var a = "abc12_121232131231?"
	var a = "a_abc12"
	fmt.Println(r.MatchPass1(a, 4, 12))
}

func TestAnyRegs_RegIsContainSpecialCharacter(t *testing.T) {
	// 测试字符串
	testStrings := []string{
		"Hello, World!",
		"GoLang123",
		"Password!@#",
		"Special_Chars%^&*",
		"NormalString",
	}

	// 验证每个字符串是否包含特殊字符
	for _, str := range testStrings {
		if r.MatchIsContainSpecialCharacter(str) {
			fmt.Printf("The string \"%s\" contains special characters.\n", str)
		} else {
			fmt.Printf("The string \"%s\" does not contain special characters.\n", str)
		}
	}
}

func TestAnyRegs_RegChineseCharacter(t *testing.T) {
	// 测试字符串
	testStrings := []string{
		"你好世界",      // 纯汉字
		"Hello",     // 非汉字
		"你好, World", // 混合
		"测试123",     // 混合
		"汉字",        // 纯汉字
	}

	// 验证每个字符串是否只包含汉字
	for _, str := range testStrings {
		if r.MatchChineseCharacter(str) {
			fmt.Printf("The string \"%s\" contains only Chinese characters.\n", str)
		} else {
			fmt.Printf("The string \"%s\" does not contain only Chinese characters.\n", str)
		}
	}
}

func TestAnyRegs_RegEmail(t *testing.T) {
	var a = "aabc@qq.com"
	fmt.Println(r.MatchEmail(a))
}

func TestAnyRegs_RegChinePhoneNumber(t *testing.T) {
	testNumbers := []string{
		"13800138000", // valid
		"12345678901", // invalid
		"19912345678", // valid
		"10000000000", // invalid
	}

	for _, number := range testNumbers {
		if r.MatchChinesePhoneNumber(number) {
			fmt.Printf("%s is a valid phone number.\n", number)
		} else {
			fmt.Printf("%s is an invalid phone number.\n", number)
		}
	}
}

func TestAnyRegs_RegContainChineseCharacter(t *testing.T) {
	testStrings := []string{
		"Hello, 世界",     // contains Chinese characters
		"你好",            // all Chinese characters
		"Hello, World!", // no Chinese characters
		"Go语言",          // contains Chinese characters
	}

	for _, str := range testStrings {
		if r.MatchContainChineseCharacter(str) {
			fmt.Printf("'%s' contains Chinese characters.\n", str)
		} else {
			fmt.Printf("'%s' does not contain Chinese characters.\n", str)
		}
	}
}

func TestAnyRegs_MatchDoubleByte(t *testing.T) {
	testStrings := []string{
		"Hello, 世界",     // contains double-byte characters
		"你好",            // all double-byte characters
		"Hello, World!", // no double-byte characters
		"Go语言",          // contains double-byte characters
		"こんにちは",         // Japanese Hiragana (double-byte characters)
	}

	for _, str := range testStrings {
		if r.MatchDoubleByte(str) {
			fmt.Printf("'%s' contains double-byte characters.\n", str)
		} else {
			fmt.Printf("'%s' does not contain double-byte characters.\n", str)
		}
	}
}

func TestAnyRegs_MatchEmptyLine(t *testing.T) {
	testStrings := []string{
		"Hello, World!\n\nThis is a test.\n\n", // contains empty lines
		"No empty lines here.",                 // no empty lines
		"\n\n",                                 // all empty lines
		"   \nThis line has spaces.\n",         // contains an empty line with spaces
	}

	for _, str := range testStrings {
		if r.MatchEmptyLine(str) {
			fmt.Printf("The string contains empty lines:\n'%s'\n", str)
		} else {
			fmt.Printf("The string does not contain empty lines:\n'%s'\n", str)
		}
	}
}

func TestAnyRegs_MatchIPv4(t *testing.T) {
	testIPs := []string{
		"192.168.1.1",     // valid
		"255.255.255.255", // valid
		"0.0.0.0",         // valid
		"256.256.256.256", // invalid
		"192.168.1.256",   // invalid
		"192.168.1",       // invalid
		"abc.def.ghi.jkl", // invalid
	}

	for _, ip := range testIPs {
		if r.MatchIPv4(ip) {
			fmt.Printf("%s is a valid IPv4 address.\n", ip)
		} else {
			fmt.Printf("%s is an invalid IPv4 address.\n", ip)
		}
	}
}

func TestAnyRegs_MatchIPv6(t *testing.T) {
	testIPs := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334", // valid
		"2001:db8:85a3::8a2e:370:7334",            // valid
		"::1",                                     // valid
		"fe80::1ff:fe23:4567:890a",                // valid
		"192.168.1.1",                             // invalid (IPv4)
		"2001:db8:85a3:0:0:8a2e:370g:7334",        // invalid (contains 'g')
		"2001::85a3::8a2e:370:7334",               // invalid (double '::')
	}

	for _, ip := range testIPs {
		if r.MatchIPv6(ip) {
			fmt.Printf("%s is a valid IPv6 address.\n", ip)
		} else {
			fmt.Printf("%s is an invalid IPv6 address.\n", ip)
		}
	}
}
