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

// mustCompile 预编译正则，空 pattern 返回 nil
func mustCompile(pat string) *regexp.Regexp {
	if pat == "" {
		return nil
	}
	return regexp.MustCompile(pat)
}

// osRules 操作系统识别规则（预编译正则，避免循环内重复编译）
var osRules = []struct {
	kw, os string
	re     *regexp.Regexp
}{
	{"iphone", IOS, mustCompile(`os (\d+[._]\d+)`)},
	{"ipad", IOS, mustCompile(`os (\d+[._]\d+)`)},
	{"mac os x", MacOS, mustCompile(`mac os x (\d+[._]\d+)`)},
	{"android", Android, mustCompile(`android[ /](\d+(?:\.\d+)?)`)},
	{"harmonyos", OpenHarmony, mustCompile(`harmonyos[ /]?(\d+)`)},
	{"linux", Linux, nil},
	{"windows phone", WindowsPhone, mustCompile(`windows phone (?:os )?(\d+)`)},
	{"cros", ChromeOS, nil},
	{"freebsd", FreeBSD, nil},
}

// browserRules 浏览器识别规则（预编译正则，避免循环内重复编译）
var browserRules = []struct {
	kw, name string
	re       *regexp.Regexp
}{
	{"edg", Edge, mustCompile(`edg[e]?/(\d+)`)},
	{"opr", Opera, mustCompile(`opr/(\d+)`)},
	{"yabrowser", YandexBrowser, mustCompile(`yabrowser/(\d+)`)},
	{"samsungbrowser", SamsungBrowser, mustCompile(`samsungbrowser/(\d+)`)},
	{"chrome", Chrome, mustCompile(`chrome/(\d+)`)},
	{"firefox", Firefox, mustCompile(`firefox/(\d+)`)},
	{"safari", Safari, mustCompile(`version/(\d+)`)},
}

// 以下三张表原先写在各自的函数体内，每次 Parse 都要重建一遍
//
// osRules 与 browserRules 当初已经提到包级（注释写着「避免循环内重复编译」），
// 这三张被漏掉了其中两张还是 map 字面量：map 是可变引用类型，编译器无法把它
// 放进静态数据，每次调用都要分配桶并逐个插入键
//
// 更要紧的是 map 换成了有序切片：Go 遍历 map 的顺序是随机的，而这里的循环是
// 「谁先匹配上就用谁」一个 UA 若同时含 "sm-" 与 "samsung"（三星设备的 UA 常
// 同时出现型号与品牌），或 iPad 在桌面模式下同时含 "ipad" 与 "macintosh"，
// 那么 Device 字段会在多次调用之间随机变化——不报错，只是同一个 UA 有时得到
// 这个值有时得到那个值顺序按源码里原本的书写次序固定下来，那个次序本身读起来
// 就是有意安排的（具体型号在前、通用品牌在后）
//
// [EN] These three tables used to sit inside their functions and were rebuilt on
// every Parse. osRules and browserRules had already been hoisted — their comment
// says it is to avoid recompiling inside the loop — and these were missed. Two of
// them were map literals, and a map is a mutable reference type the compiler
// cannot place in static data, so every call allocated buckets and inserted each
// key. More importantly the maps became ordered slices: Go iterates a map in
// random order while these loops take the first match, so a UA containing both
// "sm-" and "samsung", or an iPad in desktop mode containing both "ipad" and
// "macintosh", would set the Device field to a value that varied between calls —
// no error, just the same input giving different answers. The order is fixed to
// the one the source already used, which reads as deliberate: specific models
// first, generic brands after

// botRules 爬虫识别规则（按优先级排序，更具体的放在前面，通用关键词最后）
var botRules = []struct {
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

// windowsVersions Windows 版本映射（按版本号降序，必须先匹配更具体的）
var windowsVersions = []struct {
	kw, version string
}{
	{"windows nt 10.0", "10"},
	{"windows nt 6.3", "8.1"},
	{"windows nt 6.2", "8"},
	{"windows nt 6.1", "7"},
	{"windows nt 6.0", "Vista"},
	{"windows nt 5", "XP"},
}

// vendorRules 设备厂商识别规则（具体型号前缀在前，通用品牌名在后）
var vendorRules = []struct {
	kw, vendor string
}{
	{"iphone", VendorApple},
	{"ipad", VendorApple},
	{"macintosh", VendorApple},
	{"sm-", VendorSamsung},
	{"samsung", VendorSamsung},
	{"huawei", VendorHuawei},
	{"honor", VendorHonor},
	{"xiaomi", VendorXiaomi},
	{"oppo", VendorOPPO},
	{"vivo", VendorVivo},
}

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
	for _, b := range botRules {
		if strings.Contains(s, b.keyword) {
			p.IsBot, p.BotName, p.DeviceType = true, b.name, DeviceBot
			return
		}
	}
}

// parseOS 解析操作系统信息
// 识别 Windows, macOS, iOS, Android, Linux 等主流操作系统
func (p *ParsedUserAgent) parseOS(s string) {
	for _, w := range windowsVersions {
		if strings.Contains(s, w.kw) {
			p.OS, p.OSVersion = Windows, w.version
			return
		}
	}

	// 其他操作系统规则: 关键词, 系统名称, 版本提取正则
	// 注意: iOS 设备的 UA 中包含 "Mac OS X", 所以必须先检测 iPhone/iPad
	rules := osRules

	for _, r := range rules {
		if strings.Contains(s, r.kw) {
			p.OS = r.os
			// 提取版本号
			if r.re != nil {
				if m := r.re.FindStringSubmatch(s); len(m) > 1 {
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
	for _, b := range browserRules {
		if strings.Contains(s, b.kw) {
			p.Browser = b.name
			// 提取版本号(主版本号)
			if b.re != nil {
				if m := b.re.FindStringSubmatch(s); len(m) > 1 {
					p.BrowserVersion = m[1]
				}
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

	for _, v := range vendorRules {
		if strings.Contains(s, v.kw) {
			p.DeviceVendor = v.vendor
			p.Device = v.kw
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
