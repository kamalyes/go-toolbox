/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 10:50:01
 * @FilePath: \go-toolbox\pkg\random\business.go
 * @Description: 生成业务相关的随机数据
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

import (
	"fmt"
	"strings"
)

// 邮箱域名列表
var emailDomains = []string{
	"gmail.com", "yahoo.com", "hotmail.com", "outlook.com",
	"qq.com", "163.com", "126.com", "sina.com",
}

// RandomEmail 生成随机邮箱
func RandomEmail() string {
	username := RandString(8, NUMBER|LOWERCASE|CAPITAL)
	domain := emailDomains[RandInt(0, len(emailDomains)-1)]
	return fmt.Sprintf("%s@%s", strings.ToLower(username), domain)
}

// RandomPhone 生成随机手机号（中国大陆）
func RandomPhone() string {
	prefixes := []string{"130", "131", "132", "133", "134", "135", "136", "137", "138", "139",
		"150", "151", "152", "153", "155", "156", "157", "158", "159",
		"180", "181", "182", "183", "184", "185", "186", "187", "188", "189"}

	prefix := prefixes[RandInt(0, len(prefixes)-1)]
	suffix := RandInt(10000000, 99999999)

	return fmt.Sprintf("%s%d", prefix, suffix)
}

// RandomName 生成随机姓名（中文）
func RandomName() string {
	surnames := []string{"王", "李", "张", "刘", "陈", "杨", "黄", "赵", "周", "吴"}
	names := []string{"伟", "芳", "娜", "秀", "英", "敏", "静", "丽", "强", "磊", "军", "勇", "杰", "涛", "超"}

	surname := surnames[RandInt(0, len(surnames)-1)]

	// 60% 双字名，40% 单字名
	if RandInt(0, 99) < 60 {
		name1 := names[RandInt(0, len(names)-1)]
		name2 := names[RandInt(0, len(names)-1)]
		return surname + name1 + name2
	}

	name := names[RandInt(0, len(names)-1)]
	return surname + name
}

// RandomIDCard 生成随机身份证号（仅用于测试）
func RandomIDCard() string {
	// 地区码（随机）
	areaCode := fmt.Sprintf("%06d", RandInt(110000, 659999))

	// 出生日期（1960-2000）
	year := RandInt(1960, 2000)
	month := RandInt(1, 12)
	day := RandInt(1, 28)
	birthDate := fmt.Sprintf("%04d%02d%02d", year, month, day)

	// 顺序码
	sequence := fmt.Sprintf("%03d", RandInt(0, 999))

	// 前17位
	id17 := areaCode + birthDate + sequence

	// 计算校验码
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	sum := 0
	for i, c := range id17 {
		sum += int(c-'0') * weights[i]
	}
	checkCode := checkCodes[sum%11]

	return id17 + checkCode
}

// RandomCompany 生成随机公司名称
func RandomCompany() string {
	// 前缀词汇 - 更丰富的选择
	prefixes := []string{
		"云", "智", "创", "鑫", "盛", "华", "金", "博", "新", "众",
		"联", "汇", "宏", "腾", "飞", "星", "光", "明", "通", "达",
		"远", "诚", "信", "德", "瑞", "优", "美", "天", "地", "海",
		"泰", "安", "和", "平", "正", "中", "方", "圆", "时", "代",
	}

	// 中间词汇 - 行业特征
	middles := []string{
		"科技", "网络", "信息", "数字", "智能", "云端", "数据", "软件",
		"互联", "电子", "通信", "系统", "工程", "服务", "咨询", "传媒",
		"文化", "教育", "医疗", "金融", "商贸", "物流", "环保", "能源",
	}

	// 后缀类型
	types := []string{
		"有限公司", "股份有限公司", "集团有限公司", "科技有限公司",
		"（中国）有限公司", "实业有限公司", "控股有限公司",
	}

	// 随机组合：40%概率使用双前缀，60%概率使用单前缀
	var companyName string
	if RandInt(0, 99) < 40 {
		// 双前缀
		prefix1 := prefixes[RandInt(0, len(prefixes)-1)]
		prefix2 := prefixes[RandInt(0, len(prefixes)-1)]
		// 确保两个前缀不同
		for prefix1 == prefix2 {
			prefix2 = prefixes[RandInt(0, len(prefixes)-1)]
		}
		middle := middles[RandInt(0, len(middles)-1)]
		typeStr := types[RandInt(0, len(types)-1)]
		companyName = fmt.Sprintf("%s%s%s%s", prefix1, prefix2, middle, typeStr)
	} else {
		// 单前缀
		prefix := prefixes[RandInt(0, len(prefixes)-1)]
		middle := middles[RandInt(0, len(middles)-1)]
		typeStr := types[RandInt(0, len(types)-1)]
		companyName = fmt.Sprintf("%s%s%s", prefix, middle, typeStr)
	}

	return companyName
}
