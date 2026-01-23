/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 10:50:01
 * @FilePath: \go-toolbox\pkg\random\business_test.go
 * @Description: 业务随机数据生成测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomEmail(t *testing.T) {
	// 测试生成100个邮箱
	for i := 0; i < 100; i++ {
		email := RandomEmail()

		// 验证邮箱格式
		assert.NotEmpty(t, email, "邮箱不能为空")
		assert.Contains(t, email, "@", "邮箱必须包含@符号")

		// 验证邮箱格式是否合法
		emailRegex := regexp.MustCompile(`^[a-z0-9._-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
		assert.True(t, emailRegex.MatchString(email), "邮箱格式不合法: %s", email)

		// 验证域名是否在预定义列表中
		parts := strings.Split(email, "@")
		assert.Len(t, parts, 2, "邮箱格式错误")

		domain := parts[1]
		validDomain := false
		for _, d := range emailDomains {
			if d == domain {
				validDomain = true
				break
			}
		}
		assert.True(t, validDomain, "邮箱域名不在预定义列表中: %s", domain)
	}
}

func TestRandomPhone(t *testing.T) {
	// 测试生成100个手机号
	for i := 0; i < 100; i++ {
		phone := RandomPhone()

		// 验证手机号长度
		assert.Len(t, phone, 11, "手机号长度必须为11位")

		// 验证手机号格式（全是数字）
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		assert.True(t, phoneRegex.MatchString(phone), "手机号格式不合法: %s", phone)

		// 验证前缀
		prefix := phone[:3]
		validPrefixes := []string{
			"130", "131", "132", "133", "134", "135", "136", "137", "138", "139",
			"150", "151", "152", "153", "155", "156", "157", "158", "159",
			"180", "181", "182", "183", "184", "185", "186", "187", "188", "189",
		}

		validPrefix := false
		for _, p := range validPrefixes {
			if p == prefix {
				validPrefix = true
				break
			}
		}
		assert.True(t, validPrefix, "手机号前缀不合法: %s", prefix)
	}
}

func TestRandomName(t *testing.T) {
	singleCount := 0
	doubleCount := 0

	// 测试生成1000个姓名，统计单字名和双字名的比例
	for i := 0; i < 1000; i++ {
		name := RandomName()

		// 验证姓名不为空
		assert.NotEmpty(t, name, "姓名不能为空")

		// 验证姓名长度（2-3个字符）
		runeCount := len([]rune(name))
		assert.True(t, runeCount == 2 || runeCount == 3, "姓名长度应为2或3个字符，实际为: %d, 姓名: %s", runeCount, name)

		if runeCount == 2 {
			singleCount++
		} else {
			doubleCount++
		}
	}

	// 验证比例（60%双字名，40%单字名，允许10%误差）
	doubleRatio := float64(doubleCount) / 1000.0
	assert.InDelta(t, 0.6, doubleRatio, 0.1, "双字名比例应接近60%%，实际为: %.1f%%", doubleRatio*100)
}

func TestRandomIDCard(t *testing.T) {
	// 测试生成100个身份证号
	for i := 0; i < 100; i++ {
		idCard := RandomIDCard()

		// 验证身份证号长度
		assert.Len(t, idCard, 18, "身份证号长度必须为18位")

		// 验证前17位是数字
		for j := 0; j < 17; j++ {
			assert.True(t, idCard[j] >= '0' && idCard[j] <= '9', "前17位必须是数字")
		}

		// 验证最后一位是数字或X
		lastChar := idCard[17]
		assert.True(t, (lastChar >= '0' && lastChar <= '9') || lastChar == 'X', "最后一位必须是数字或X")

		// 验证地区码范围
		areaCode := idCard[:6]
		assert.Regexp(t, `^[1-6]\d{5}$`, areaCode, "地区码格式不正确")

		// 验证出生日期格式
		birthDate := idCard[6:14]
		year := birthDate[:4]
		month := birthDate[4:6]
		day := birthDate[6:8]

		assert.Regexp(t, `^(19[6-9]\d|20[0-9]\d)$`, year, "年份应在1960-2000之间")
		assert.Regexp(t, `^(0[1-9]|1[0-2])$`, month, "月份应在01-12之间")
		assert.Regexp(t, `^(0[1-9]|[12]\d)$`, day, "日期应在01-28之间")

		// 验证校验码算法
		weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		checkCodes := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

		sum := 0
		for k := 0; k < 17; k++ {
			sum += int(idCard[k]-'0') * weights[k]
		}
		expectedCheck := checkCodes[sum%11]
		actualCheck := string(idCard[17])

		assert.Equal(t, expectedCheck, actualCheck, "身份证校验码不正确，身份证号: %s", idCard)
	}
}

func TestRandomCompany(t *testing.T) {
	singlePrefixCount := 0
	doublePrefixCount := 0

	// 测试生成1000个公司名称
	for i := 0; i < 1000; i++ {
		company := RandomCompany()

		// 验证公司名称不为空
		assert.NotEmpty(t, company, "公司名称不能为空")

		// 验证公司名称包含"有限公司"
		assert.True(t, strings.Contains(company, "有限公司") ||
			strings.Contains(company, "股份有限公司") ||
			strings.Contains(company, "集团有限公司"),
			"公司名称应包含公司类型后缀: %s", company)

		// 统计单前缀和双前缀的比例
		// 移除所有已知后缀
		name := company
		for _, suffix := range []string{"（中国）有限公司", "股份有限公司", "集团有限公司", "科技有限公司", "实业有限公司", "控股有限公司", "有限公司"} {
			if strings.HasSuffix(name, suffix) {
				name = strings.TrimSuffix(name, suffix)
				break
			}
		}

		// 移除所有已知中间词（行业特征）
		middles := []string{"科技", "网络", "信息", "数字", "智能", "云端", "数据", "软件",
			"互联", "电子", "通信", "系统", "工程", "服务", "咨询", "传媒",
			"文化", "教育", "医疗", "金融", "商贸", "物流", "环保", "能源"}
		for _, mid := range middles {
			if strings.Contains(name, mid) {
				name = strings.Replace(name, mid, "", 1)
				break
			}
		}

		// 判断前缀数量：移除行业词后，单前缀1个字=单前缀，2个字及以上=双前缀
		runeCount := len([]rune(name))
		if runeCount >= 2 {
			doublePrefixCount++
		} else if runeCount >= 1 {
			singlePrefixCount++
		}
	}

	// 验证比例（40%双前缀，60%单前缀，允许10%误差）
	doubleRatio := float64(doublePrefixCount) / 1000.0
	assert.InDelta(t, 0.4, doubleRatio, 0.1, "双前缀比例应接近40%%，实际为: %.1f%%", doubleRatio*100)
}

func TestRandomCompanyUniqueness(t *testing.T) {
	// 测试生成的公司名称有一定的唯一性
	companies := make(map[string]bool)
	duplicates := 0

	for i := 0; i < 500; i++ {
		company := RandomCompany()
		if companies[company] {
			duplicates++
		}
		companies[company] = true
	}

	// 重复率应该很低（小于5%）
	duplicateRatio := float64(duplicates) / 500.0
	assert.Less(t, duplicateRatio, 0.05, "公司名称重复率过高: %.1f%%", duplicateRatio*100)
}

// Benchmark tests
func BenchmarkRandomEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomEmail()
	}
}

func BenchmarkRandomPhone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomPhone()
	}
}

func BenchmarkRandomName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomName()
	}
}

func BenchmarkRandomIDCard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomIDCard()
	}
}

func BenchmarkRandomCompany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomCompany()
	}
}
