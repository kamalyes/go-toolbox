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
	Engine         string `json:"engine,omitempty"`   // 渲染引擎: Blink, WebKit, Gecko 等
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
	p.Engine = engineOf(p.Browser, p.OS)
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

// ============================================================================
// 扩展信息
//
// 为什么引擎是字段、其余三项是方法
//
// 引擎能从已经识别出的浏览器与系统直接推出来，一次 switch，不用再扫字符串，
// 所以放进字段里顺手填上
// 端内浏览器、HTTP 客户端、CPU 架构都需要在 UA 里逐个找关键词，而每个关键词都是
// 一次全串扫描（实测约 20 纳秒一个）端内那张表有二十项，若在 Parse 里一并扫完，
// 光这一项就要多付约四百纳秒，几乎抵掉这个包刚做的全部优化
// 所以这三项做成按需调用的方法：不用的人一分不付，要用的人调一次
// 方法不做缓存，避免给 ParsedUserAgent 引入可变状态——那会让同一个实例在多个
// goroutine 间读取时产生竞争，而 Parse 本来每次都返回新实例
//
// [EN] The engine is a field while the other three are methods. The engine
// follows from the already-detected browser and OS with a single switch and no
// further string scanning, so it is filled in on the spot. In-app browser, HTTP
// client and CPU architecture each need keyword searches through the UA, and
// every keyword is a full scan measured at about 20 nanoseconds; the in-app table
// has twenty entries, so scanning it inside Parse would add some four hundred
// nanoseconds and undo most of this package's recent gains. They are therefore
// computed on demand: callers who do not need them pay nothing. The methods do
// not cache, which would give ParsedUserAgent mutable state and make reads of one
// instance from several goroutines race, while Parse already returns a fresh
// instance each time
// ============================================================================

// engineOf 由浏览器与系统推出渲染引擎
//
// iOS 上所有浏览器都必须用系统的 WebKit（App Store 规则），所以 iOS 上的
// Chrome 与 Firefox 也是 WebKit 而不是 Blink / Gecko
//
// [EN] Derives the rendering engine. Every browser on iOS is required to use the
// system WebKit by App Store rules, so Chrome and Firefox there are WebKit too
func engineOf(browser, os string) string {
	if os == IOS {
		return EngineWebKit
	}
	switch browser {
	case Chrome, Edge, Opera, SamsungBrowser, YandexBrowser:
		return EngineBlink
	case Safari:
		return EngineWebKit
	case Firefox:
		return EngineGecko
	}
	return ""
}

// inAppRules 端内浏览器识别规则
//
// 顺序要紧：更具体的在前企业微信的 UA 同时含 wxwork 与 micromessenger，
// 必须先判企业微信；QQ 内置浏览器含 mqqbrowser，而 QQ 空间等又含 qq，
// 所以具体标识在前、通用标识在后
//
// [EN] Order matters, most specific first: WeCom's UA carries both wxwork and
// micromessenger so WeCom must be checked first, and QQ's browser carries
// mqqbrowser while other QQ surfaces carry only qq
var inAppRules = []struct {
	kw, app string
}{
	{"wxwork", InAppWeCom},
	{"micromessenger", InAppWeChat},
	{"alipayclient", InAppAlipay},
	// 具体应用必须排在通用的 aliapp 之前
	//
	// AliApp(...) 是阿里系通用的 WebView 容器标识，支付宝、钉钉、手机淘宝都带它，
	// 单看它分不出是哪一个钉钉的 UA 形如 DingTalk/7.0.0 AliApp(DingTalk/7.0.0)，
	// 若 aliapp 排在前面就会被误判成支付宝——这条是测试抓出来的
	//
	// [EN] Specific apps must come before the generic aliapp. AliApp(...) is
	// Alibaba's shared WebView container marker, carried by Alipay, DingTalk and
	// Taobao alike, so on its own it cannot tell them apart: DingTalk's UA reads
	// DingTalk/7.0.0 AliApp(DingTalk/7.0.0) and would be taken for Alipay if
	// aliapp came first. The test caught this
	{"dingtalk", InAppDingTalk},
	{"aliapp(tb", InAppTaobao},
	{"aliapp", InAppAlipay},
	{"feishu", InAppFeishu},
	{"lark", InAppFeishu},
	{"mqqbrowser", InAppQQ},
	{"weibo", InAppWeibo},
	{"aweme", InAppDouyin},
	{"bytedancewebview", InAppDouyin},
	{"kuaishou", InAppKuaishou},
	{"xhsdiscover", InAppXiaohongs},
	{"baiduboxapp", InAppBaiduApp},
	{"jdapp", InAppJD},
	{"meituan", InAppMeituan},
	{"line/", InAppLine},
	{"fban", InAppFacebook},
	{"fbav", InAppFacebook},
	{"instagram", InAppInstagram},
	{"twitter", InAppTwitter},
	{"snapchat", InAppSnapchat},
	{"electron", InAppElectron},
	{"qq/", InAppQQ},
}

// clientRules 非浏览器 HTTP 客户端识别规则
var clientRules = []struct {
	kw, client string
}{
	{"curl/", ClientCurl},
	{"wget/", ClientWget},
	{"okhttp", ClientOkHttp},
	{"python-requests", ClientPython},
	{"python-urllib", ClientPython},
	{"postmanruntime", ClientPostman},
	{"insomnia", ClientInsomnia},
	{"java/", ClientJava},
	{"apache-httpclient", ClientApacheHTTP},
	{"go-http-client", ClientGoHTTP},
	{"axios/", ClientAxios},
	{"node-fetch", ClientNodeFetch},
	{"restsharp", ClientRestSharp},
	{"php/", ClientPHP},
	{"ruby", ClientRuby},
}

// InApp 返回承载这个 WebView 的宿主应用名，不是端内浏览器则返回空串
//
// 按需计算，不缓存，原因见本节开头的说明
func (p *ParsedUserAgent) InApp() string {
	if p == nil || p.Raw == "" {
		return ""
	}
	lower := strings.ToLower(p.Raw)
	for _, r := range inAppRules {
		if strings.Contains(lower, r.kw) {
			return r.app
		}
	}
	return ""
}

// IsInApp 是否运行在某个应用内置的 WebView 里
func (p *ParsedUserAgent) IsInApp() bool { return p.InApp() != "" }

// HTTPClient 返回非浏览器 HTTP 客户端的名称，不是则返回空串
//
// 与爬虫分开判断：爬虫是抓取行为，这里是程序化调用
func (p *ParsedUserAgent) HTTPClient() string {
	if p == nil || p.Raw == "" {
		return ""
	}
	lower := strings.ToLower(p.Raw)
	for _, r := range clientRules {
		if strings.Contains(lower, r.kw) {
			return r.client
		}
	}
	return ""
}

// IsHTTPClient 是否为非浏览器的程序化调用
func (p *ParsedUserAgent) IsHTTPClient() bool { return p.HTTPClient() != "" }

// Arch 返回 CPU 架构，无法判断则返回空串
//
// 注意 UA 里的架构信息是不可靠的：Windows 上的 32 位浏览器在 64 位系统上会报
// WOW64，而 Apple Silicon 的 Safari 为了兼容仍然报 Intel Mac OS X
// 这个方法只反映 UA 的字面内容，不做推断
//
// [EN] Architecture as stated in the UA, which is not reliable: a 32-bit browser
// on 64-bit Windows reports WOW64, and Safari on Apple Silicon still reports
// Intel Mac OS X for compatibility. This reflects the literal UA only and infers
// nothing
func (p *ParsedUserAgent) Arch() string {
	if p == nil || p.Raw == "" {
		return ""
	}
	lower := strings.ToLower(p.Raw)
	switch {
	case strings.Contains(lower, "win64") || strings.Contains(lower, "x86_64") ||
		strings.Contains(lower, "x64") || strings.Contains(lower, "wow64") ||
		strings.Contains(lower, "amd64"):
		return ArchX8664
	case strings.Contains(lower, "aarch64") || strings.Contains(lower, "arm64"):
		return ArchARM64
	case strings.Contains(lower, "armv") || strings.Contains(lower, "arm"):
		return ArchARM
	case strings.Contains(lower, "i686") || strings.Contains(lower, "i386") ||
		strings.Contains(lower, "x86"):
		return ArchX86
	}
	return ""
}

// IsDesktop 是否为桌面设备
func (p *ParsedUserAgent) IsDesktop() bool {
	return p != nil && p.DeviceType == DeviceDesktop
}

// IsIOS 是否为 iOS 系统
func (p *ParsedUserAgent) IsIOS() bool { return p != nil && p.OS == IOS }

// IsAndroid 是否为 Android 系统
func (p *ParsedUserAgent) IsAndroid() bool { return p != nil && p.OS == Android }

// IsWindows 是否为 Windows 系统
func (p *ParsedUserAgent) IsWindows() bool { return p != nil && p.OS == Windows }

// IsMacOS 是否为 macOS 系统
func (p *ParsedUserAgent) IsMacOS() bool { return p != nil && p.OS == MacOS }

// IsChromium 是否为 Chromium 系浏览器（引擎为 Blink）
func (p *ParsedUserAgent) IsChromium() bool { return p != nil && p.Engine == EngineBlink }
