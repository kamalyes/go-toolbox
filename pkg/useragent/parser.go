/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 15:30:00
 * @FilePath: \go-toolbox\pkg\useragent\parser.go
 * @Description: User-Agent 解析器
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package useragent

import (
	"regexp"
	"strings"
)

// ParsedUserAgent 解析后的 User-Agent 信息
type ParsedUserAgent struct {
	Raw            string `json:"raw"`                // 原始 User-Agent 字符串
	Browser        string `json:"browser"`            // 浏览器名称: Chrome, Firefox, Safari 等
	BrowserVersion string `json:"browser_version"`    // 浏览器版本号
	OS             string `json:"os"`                 // 操作系统: Windows, macOS, Android, iOS 等
	OSVersion      string `json:"os_version"`         // 操作系统版本号
	Device         string `json:"device"`             // 设备名称
	DeviceType     string `json:"device_type"`        // 设备类型: mobile/tablet/desktop/bot
	DeviceVendor   string `json:"device_vendor"`      // 设备厂商: Apple, Samsung, Huawei 等
	IsBot          bool   `json:"is_bot"`             // 是否为爬虫/机器人
	BotName        string `json:"bot_name,omitempty"` // 爬虫名称
	IsMobile       bool   `json:"is_mobile"`          // 是否为移动设备
	IsTablet       bool   `json:"is_tablet"`          // 是否为平板设备
}

// Parse 解析 User-Agent 字符串,返回解析后的结构化信息
// 包括浏览器、操作系统、设备类型等信息
func Parse(ua string) *ParsedUserAgent {
	if ua == "" {
		return &ParsedUserAgent{Raw: ua, DeviceType: DeviceUnknown}
	}

	p := &ParsedUserAgent{Raw: ua}
	lower := strings.ToLower(ua)

	// 优先检测爬虫,如果是爬虫则不再解析其他信息
	p.parseBot(lower)
	if p.IsBot {
		return p
	}

	// 按顺序解析各项信息
	p.parseOS(lower)
	p.parseBrowser(lower)
	p.parseDevice(lower)
	p.determineDeviceType()
	return p
}

// parseBot 检测是否为爬虫/机器人
// 检测常见的搜索引擎爬虫和社交媒体爬虫
func (p *ParsedUserAgent) parseBot(s string) {
	// 常见爬虫关键词映射 (按优先级排序,更具体的放在前面)
	// 注意: 必须先检查具体的爬虫名称,再检查通用关键词
	bots := []struct {
		keyword string
		name    string
	}{
		{"googlebot", Googlebot},                     // Google 搜索爬虫
		{"bingbot", Bingbot},                         // Bing 搜索爬虫
		{"yandexbot", YandexBot},                     // Yandex 搜索爬虫
		{"baiduspider", BotBaidu},                    // 百度搜索爬虫
		{"slurp", BotYahoo},                          // Yahoo 搜索爬虫
		{"twitterbot", Twitterbot},                   // Twitter 爬虫
		{"facebookexternalhit", FacebookExternalHit}, // Facebook 外链爬虫
		{"applebot", Applebot},                       // Apple 搜索爬虫
		{"spider", BotSpider},                        // 通用爬虫
		{"crawler", BotCrawler},                      // 通用爬虫
		{"bot", BotGeneric},                          // 通用爬虫 (必须放在最后)
	}
	for _, b := range bots {
		if strings.Contains(s, b.keyword) {
			p.IsBot, p.BotName, p.DeviceType = true, b.name, DeviceBot
			return
		}
	}
}

// parseOS 解析操作系统信息
// 识别 Windows, macOS, iOS, Android, Linux 等主流操作系统
func (p *ParsedUserAgent) parseOS(s string) {
	// Windows 版本映射表
	winVer := map[string]string{
		"windows nt 10.0": "10",
		"windows nt 6.3":  "8.1",
		"windows nt 6.2":  "8",
		"windows nt 6.1":  "7",
		"windows nt 6.0":  "Vista",
		"windows nt 5":    "XP",
	}
	for k, v := range winVer {
		if strings.Contains(s, k) {
			p.OS, p.OSVersion = Windows, v
			return
		}
	}

	// 其他操作系统规则: 关键词, 系统名称, 版本提取正则
	// 注意: iOS 设备的 UA 中包含 "Mac OS X", 所以必须先检测 iPhone/iPad
	rules := []struct {
		kw, os, pat string
	}{
		{"iphone", IOS, `os (\d+[._]\d+)`},                             // iOS (iPhone) - 必须在 macOS 之前
		{"ipad", IOS, `os (\d+[._]\d+)`},                               // iOS (iPad) - 必须在 macOS 之前
		{"mac os x", MacOS, `mac os x (\d+[._]\d+)`},                   // macOS
		{"android", Android, `android[ /](\d+(?:\.\d+)?)`},             // Android
		{"harmonyos", OpenHarmony, `harmonyos[ /]?(\d+)`},              // HarmonyOS
		{"linux", Linux, ""},                                           // Linux
		{"windows phone", WindowsPhone, `windows phone (?:os )?(\d+)`}, // Windows Phone
		{"cros", ChromeOS, ""},                                         // ChromeOS
		{"freebsd", FreeBSD, ""},                                       // FreeBSD
	}

	for _, r := range rules {
		if strings.Contains(s, r.kw) {
			p.OS = r.os
			// 提取版本号
			if r.pat != "" {
				if m := regexp.MustCompile(r.pat).FindStringSubmatch(s); len(m) > 1 {
					p.OSVersion = strings.ReplaceAll(m[1], "_", ".")
				}
			}
			return
		}
	}
}

// parseBrowser 解析浏览器信息
// 识别主流浏览器: Chrome, Firefox, Safari, Edge, Opera 等
func (p *ParsedUserAgent) parseBrowser(s string) {
	// 浏览器规则: 关键词, 浏览器名称, 版本提取正则
	browsers := []struct {
		kw, name, pat string
	}{
		{"edg", Edge, `edg[e]?/(\d+)`},                             // Microsoft Edge
		{"opr", Opera, `opr/(\d+)`},                                // Opera
		{"yabrowser", YandexBrowser, `yabrowser/(\d+)`},            // Yandex Browser
		{"samsungbrowser", SamsungBrowser, `samsungbrowser/(\d+)`}, // Samsung Browser
		{"chrome", Chrome, `chrome/(\d+)`},                         // Google Chrome
		{"firefox", Firefox, `firefox/(\d+)`},                      // Mozilla Firefox
		{"safari", Safari, `version/(\d+)`},                        // Apple Safari
	}

	for _, b := range browsers {
		if strings.Contains(s, b.kw) {
			p.Browser = b.name
			// 提取版本号(主版本号)
			if m := regexp.MustCompile(b.pat).FindStringSubmatch(s); len(m) > 1 {
				p.BrowserVersion = m[1]
			}
			return
		}
	}
}

// parseDevice 解析设备信息
// 识别移动设备、平板设备及设备厂商
func (p *ParsedUserAgent) parseDevice(s string) {
	// 检测移动设备和平板
	p.IsMobile = strings.Contains(s, "mobile") || strings.Contains(s, "android")
	p.IsTablet = strings.Contains(s, "tablet") || strings.Contains(s, "ipad")

	// 平板不算移动设备
	if p.IsTablet {
		p.IsMobile = false
	}

	// 设备厂商映射
	vendors := map[string]string{
		"iphone":    VendorApple,   // iPhone
		"ipad":      VendorApple,   // iPad
		"macintosh": VendorApple,   // Mac
		"sm-":       VendorSamsung, // 三星 (SM- 开头的设备型号)
		"samsung":   VendorSamsung, // 三星
		"huawei":    VendorHuawei,  // 华为
		"honor":     VendorHonor,   // 荣耀
		"xiaomi":    VendorXiaomi,  // 小米
		"oppo":      VendorOPPO,    // OPPO
		"vivo":      VendorVivo,    // Vivo
	}

	for k, v := range vendors {
		if strings.Contains(s, k) {
			p.DeviceVendor = v
			p.Device = k
			return
		}
	}
}

// determineDeviceType 根据设备特征确定最终的设备类型
// 优先级: bot > tablet > mobile > desktop
func (p *ParsedUserAgent) determineDeviceType() {
	if p.IsBot {
		p.DeviceType = DeviceBot
	} else if p.IsTablet {
		p.DeviceType = DeviceTablet
	} else if p.IsMobile {
		p.DeviceType = DeviceMobile
	} else {
		p.DeviceType = DeviceDesktop
	}
}
