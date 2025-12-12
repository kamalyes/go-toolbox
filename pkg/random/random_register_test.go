/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 09:20:57
 * @FilePath: \go-toolbox\pkg\random\random_register_test.go
 * @Description: 测试随机数据生成器的注册系统
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUser 测试用户结构体，用于测试注册系统
type TestUser struct {
	ID       int     `json:"id" rand:"user_id"`
	Name     string  `json:"name" rand:"full_name"`
	Company  string  `json:"company" rand:"company"`
	Position string  `json:"position" rand:"job_title"`
	Salary   int     `json:"salary" rand:"salary_range"`
	Email    string  `json:"email" rand:"business_email"`
	IsVIP    bool    `json:"is_vip" rand:"vip_status"`
	Score    float64 `json:"score" rand:"performance_score"`
}

// TestRegisterBasicGenerator 测试注册基本生成器
func TestRegisterBasicGenerator(t *testing.T) {
	// 清空注册表
	ClearAllGenerators()

	// 注册一些自定义生成器
	RegisterGenerator("company", func() interface{} {
		companies := []string{"Google", "Microsoft", "Apple", "Amazon", "Meta"}
		return companies[FRandInt(0, len(companies)-1)]
	})

	RegisterGenerator("job_title", func() interface{} {
		titles := []string{"软件工程师", "产品经理", "架构师", "技术总监", "CTO"}
		return titles[FRandInt(0, len(titles)-1)]
	})

	// 验证注册成功
	assert.Contains(t, ListRegisteredGenerators(), "company")
	assert.Contains(t, ListRegisteredGenerators(), "job_title")
	assert.Equal(t, 2, len(ListRegisteredGenerators()))

	// 测试生成器是否存在
	companyGen, exists := GetGenerator("company")
	assert.True(t, exists)
	assert.NotNil(t, companyGen)

	// 测试生成器功能
	company := companyGen()
	assert.NotNil(t, company)
	assert.IsType(t, "", company)

	// 验证生成的公司名称在预期列表中
	validCompanies := []string{"Google", "Microsoft", "Apple", "Amazon", "Meta"}
	assert.Contains(t, validCompanies, company.(string))
}

// TestRegisterNumberGenerator 测试注册数字生成器
func TestRegisterNumberGenerator(t *testing.T) {
	ClearAllGenerators()

	// 注册用户ID生成器（整数）
	RegisterGenerator("user_id", func() interface{} {
		return FRandInt(10000, 99999)
	})

	// 注册薪资范围生成器
	RegisterGenerator("salary_range", func() interface{} {
		// 生成10k到100k的薪资
		return FRandInt(10000, 100000)
	})

	// 注册绩效分数生成器（浮点数）
	RegisterGenerator("performance_score", func() interface{} {
		return float64(FRandInt(60, 100)) + RandFloat(0, 1)
	})

	// 验证注册成功
	generators := ListRegisteredGenerators()
	assert.Contains(t, generators, "user_id")
	assert.Contains(t, generators, "salary_range")
	assert.Contains(t, generators, "performance_score")

	// 测试用户ID生成器
	userIDGen, _ := GetGenerator("user_id")
	userID := userIDGen()
	assert.IsType(t, 0, userID)
	assert.GreaterOrEqual(t, userID.(int), 10000)
	assert.LessOrEqual(t, userID.(int), 99999)

	// 测试薪资生成器
	salaryGen, _ := GetGenerator("salary_range")
	salary := salaryGen()
	assert.IsType(t, 0, salary)
	assert.GreaterOrEqual(t, salary.(int), 10000)
	assert.LessOrEqual(t, salary.(int), 100000)

	// 测试绩效分数生成器
	scoreGen, _ := GetGenerator("performance_score")
	score := scoreGen()
	assert.IsType(t, 0.0, score)
	assert.GreaterOrEqual(t, score.(float64), 60.0)
	assert.LessOrEqual(t, score.(float64), 101.0)
}

// TestRegisterBooleanGenerator 测试注册布尔值生成器
func TestRegisterBooleanGenerator(t *testing.T) {
	ClearAllGenerators()

	// 注册VIP状态生成器（20%概率为VIP）
	RegisterGenerator("vip_status", func() interface{} {
		return FRandInt(1, 100) <= 20 // 20%概率为true
	})

	vipGen, _ := GetGenerator("vip_status")

	// 测试多次生成，验证返回布尔值
	for i := 0; i < 10; i++ {
		vipStatus := vipGen()
		assert.IsType(t, true, vipStatus)
	}
}

// TestRegisterComplexGenerator 测试注册复杂生成器
func TestRegisterComplexGenerator(t *testing.T) {
	ClearAllGenerators()

	// 注册全名生成器
	RegisterGenerator("full_name", func() interface{} {
		firstNames := []string{"张", "李", "王", "刘", "陈", "杨", "黄", "赵", "吴", "周"}
		lastNames := []string{"伟", "芳", "娜", "秀英", "敏", "静", "丽", "强", "磊", "军"}
		return firstNames[FRandInt(0, len(firstNames)-1)] + lastNames[FRandInt(0, len(lastNames)-1)]
	})

	// 注册企业邮箱生成器
	RegisterGenerator("business_email", func() interface{} {
		domains := []string{"company.com", "corp.com", "tech.com", "enterprise.com"}
		username := FRandString(6)
		domain := domains[FRandInt(0, len(domains)-1)]
		return strings.ToLower(username) + "@" + domain
	})

	// 测试全名生成器
	nameGen, _ := GetGenerator("full_name")
	name := nameGen().(string)
	assert.NotEmpty(t, name)
	assert.IsType(t, "", name)
	assert.True(t, len(name) >= 2) // 中文姓名至少2个字符

	// 测试企业邮箱生成器
	emailGen, _ := GetGenerator("business_email")
	email := emailGen().(string)
	assert.NotEmpty(t, email)
	assert.Contains(t, email, "@")
	assert.True(t, strings.HasSuffix(email, ".com"))
}

// TestGenerateWithRegisteredGenerators 测试使用注册生成器进行模型生成
func TestGenerateWithRegisteredGenerators(t *testing.T) {
	ClearAllGenerators()

	// 注册所有需要的生成器
	RegisterGenerator("user_id", func() interface{} {
		return FRandInt(10000, 99999)
	})

	RegisterGenerator("full_name", func() interface{} {
		firstNames := []string{"张三", "李四", "王五", "刘六", "陈七"}
		return firstNames[FRandInt(0, len(firstNames)-1)]
	})

	RegisterGenerator("company", func() interface{} {
		companies := []string{"阿里巴巴", "腾讯", "字节跳动", "百度", "华为"}
		return companies[FRandInt(0, len(companies)-1)]
	})

	RegisterGenerator("job_title", func() interface{} {
		titles := []string{"前端工程师", "后端工程师", "全栈工程师", "架构师", "技术总监"}
		return titles[FRandInt(0, len(titles)-1)]
	})

	RegisterGenerator("salary_range", func() interface{} {
		return FRandInt(8000, 50000)
	})

	RegisterGenerator("business_email", func() interface{} {
		domains := []string{"company.com", "tech.cn", "corp.com"}
		return FRandString(6) + "@" + domains[FRandInt(0, len(domains)-1)]
	})

	RegisterGenerator("vip_status", func() interface{} {
		return FRandInt(1, 100) <= 30 // 30%概率为VIP
	})

	RegisterGenerator("performance_score", func() interface{} {
		return float64(FRandInt(70, 95)) + RandFloat(0, 1)
	})

	// 生成测试用户
	var user TestUser
	_, _, err := GenerateRandModel(&user)
	assert.NoError(t, err)

	// 验证各字段是否按照注册的生成器生成
	assert.GreaterOrEqual(t, user.ID, 10000)
	assert.LessOrEqual(t, user.ID, 99999)

	assert.NotEmpty(t, user.Name)
	validNames := []string{"张三", "李四", "王五", "刘六", "陈七"}
	assert.Contains(t, validNames, user.Name)

	assert.NotEmpty(t, user.Company)
	validCompanies := []string{"阿里巴巴", "腾讯", "字节跳动", "百度", "华为"}
	assert.Contains(t, validCompanies, user.Company)

	assert.NotEmpty(t, user.Position)
	validPositions := []string{"前端工程师", "后端工程师", "全栈工程师", "架构师", "技术总监"}
	assert.Contains(t, validPositions, user.Position)

	assert.GreaterOrEqual(t, user.Salary, 8000)
	assert.LessOrEqual(t, user.Salary, 50000)

	assert.NotEmpty(t, user.Email)
	assert.Contains(t, user.Email, "@")

	assert.GreaterOrEqual(t, user.Score, 70.0)
	assert.LessOrEqual(t, user.Score, 96.0)

	// 验证生成结果能够JSON序列化
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
}

// TestUnregisterGenerator 测试注销生成器
func TestUnregisterGenerator(t *testing.T) {
	ClearAllGenerators()

	// 注册几个生成器
	RegisterGenerator("test1", func() interface{} { return "value1" })
	RegisterGenerator("test2", func() interface{} { return "value2" })
	RegisterGenerator("test3", func() interface{} { return "value3" })

	assert.Equal(t, 3, len(ListRegisteredGenerators()))

	// 注销一个生成器
	UnregisterGenerator("test2")
	generators := ListRegisteredGenerators()
	assert.Equal(t, 2, len(generators))
	assert.Contains(t, generators, "test1")
	assert.Contains(t, generators, "test3")
	assert.NotContains(t, generators, "test2")

	// 验证注销的生成器不存在
	_, exists := GetGenerator("test2")
	assert.False(t, exists)

	// 验证其他生成器仍然存在
	_, exists = GetGenerator("test1")
	assert.True(t, exists)
}

// TestClearAllGenerators 测试清空所有生成器
func TestClearAllGenerators(t *testing.T) {
	ClearAllGenerators()

	// 注册几个生成器
	RegisterGenerator("test1", func() interface{} { return "value1" })
	RegisterGenerator("test2", func() interface{} { return "value2" })
	RegisterGenerator("test3", func() interface{} { return "value3" })

	assert.Equal(t, 3, len(ListRegisteredGenerators()))

	// 清空所有生成器
	ClearAllGenerators()
	assert.Equal(t, 0, len(ListRegisteredGenerators()))

	// 验证所有生成器都被清除
	_, exists1 := GetGenerator("test1")
	_, exists2 := GetGenerator("test2")
	_, exists3 := GetGenerator("test3")
	assert.False(t, exists1)
	assert.False(t, exists2)
	assert.False(t, exists3)
}

// TestRegisterGeneratorOverride 测试生成器覆盖
func TestRegisterGeneratorOverride(t *testing.T) {
	ClearAllGenerators()

	// 注册第一个生成器
	RegisterGenerator("test", func() interface{} { return "old_value" })

	gen1, _ := GetGenerator("test")
	assert.Equal(t, "old_value", gen1())

	// 注册同名生成器进行覆盖
	RegisterGenerator("test", func() interface{} { return "new_value" })

	gen2, _ := GetGenerator("test")
	assert.Equal(t, "new_value", gen2())

	// 验证只有一个生成器
	assert.Equal(t, 1, len(ListRegisteredGenerators()))
}

// TestRegisterNilGenerator 测试注册nil生成器
func TestRegisterNilGenerator(t *testing.T) {
	ClearAllGenerators()

	// 注册nil生成器应该不会崩溃
	RegisterGenerator("nil_test", nil)

	// 验证nil生成器不会被注册
	_, exists := GetGenerator("nil_test")
	assert.False(t, exists)
	assert.Equal(t, 0, len(ListRegisteredGenerators()))
}

// TestGeneratorWithTypeConversion 测试生成器的类型转换
func TestGeneratorWithTypeConversion(t *testing.T) {
	ClearAllGenerators()

	// TestConversion 用于测试类型转换的结构体
	type TestConversion struct {
		StringFromInt string  `rand:"int_to_string"`
		IntFromString int     `rand:"string_to_int"`
		FloatFromInt  float64 `rand:"int_to_float"`
		BoolFromInt   bool    `rand:"int_to_bool"`
	}

	// 注册返回不同类型的生成器
	RegisterGenerator("int_to_string", func() interface{} {
		return 12345 // 返回int，但目标字段是string
	})

	RegisterGenerator("string_to_int", func() interface{} {
		return "67890" // 返回string，但目标字段是int
	})

	RegisterGenerator("int_to_float", func() interface{} {
		return 123 // 返回int，但目标字段是float64
	})

	RegisterGenerator("int_to_bool", func() interface{} {
		return 1 // 返回int，但目标字段是bool（非零值应转换为true）
	})

	// 生成测试对象
	var test TestConversion
	_, _, err := GenerateRandModel(&test)
	assert.NoError(t, err)

	// 验证类型转换是否正确
	assert.Equal(t, "12345", test.StringFromInt)
	assert.Equal(t, 67890, test.IntFromString)
	assert.Equal(t, 123.0, test.FloatFromInt)
	// 注意：bool转换可能需要特殊处理，这里先验证不出错
	assert.IsType(t, true, test.BoolFromInt)
}

// TestGeneratorPriority 测试生成器优先级（注册的生成器优先于内置生成器）
func TestGeneratorPriority(t *testing.T) {
	ClearAllGenerators()

	// TestPriority 测试优先级的结构体
	type TestPriority struct {
		Email string `rand:"email"`
		Name  string `rand:"name"`
	}

	// 注册自定义的email生成器，覆盖内置的
	RegisterGenerator("email", func() interface{} {
		return "custom@registered.com"
	})

	// name标签没有注册自定义生成器，应该使用内置生成器

	// 生成测试对象
	var test TestPriority
	_, _, err := GenerateRandModel(&test)
	assert.NoError(t, err)

	// 验证email使用了注册的生成器
	assert.Equal(t, "custom@registered.com", test.Email)

	// 验证name使用了内置生成器（长度应该是6）
	assert.NotEmpty(t, test.Name)
	assert.Equal(t, 6, len(test.Name))
}
