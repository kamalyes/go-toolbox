/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 17:10:00
 * @FilePath: \go-toolbox\pkg\useragent\parser_test.go
 * @Description: User-Agent 解析器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package useragent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEmptyUA(t *testing.T) {
	ua := Parse("")
	assert.Equal(t, "", ua.Raw)
	assert.Equal(t, DeviceUnknown, ua.DeviceType)
	assert.Empty(t, ua.Browser)
	assert.Empty(t, ua.OS)
}

func TestParseChrome(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	assert.Equal(t, "Chrome", ua.Browser)
	assert.Equal(t, "120", ua.BrowserVersion)
	assert.Equal(t, "Windows", ua.OS)
	assert.Equal(t, "10", ua.OSVersion)
	assert.Equal(t, DeviceDesktop, ua.DeviceType)
	assert.False(t, ua.IsBot)
	assert.False(t, ua.IsMobile)
	assert.False(t, ua.IsTablet)
}

func TestParseFirefox(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0")

	assert.Equal(t, "Firefox", ua.Browser)
	assert.Equal(t, "119", ua.BrowserVersion)
	assert.Equal(t, "Windows", ua.OS)
	assert.Equal(t, "10", ua.OSVersion)
	assert.Equal(t, DeviceDesktop, ua.DeviceType)
}

func TestParseSafari(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15")

	assert.Equal(t, "Safari", ua.Browser)
	assert.Equal(t, "16", ua.BrowserVersion)
	assert.Equal(t, "macOS", ua.OS)
	assert.Equal(t, "10.15", ua.OSVersion)
	assert.Equal(t, "Apple", ua.DeviceVendor)
	assert.Equal(t, DeviceDesktop, ua.DeviceType)
}

func TestParseEdge(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")

	assert.Equal(t, "Edge", ua.Browser)
	assert.Equal(t, "120", ua.BrowserVersion)
	assert.Equal(t, "Windows", ua.OS)
}

func TestParseOpera(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0")

	assert.Equal(t, "Opera", ua.Browser)
	assert.Equal(t, "106", ua.BrowserVersion)
}

func TestParseYandexBrowser(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 YaBrowser/24.1.0.0 Safari/537.36")

	assert.Equal(t, "Yandex Browser", ua.Browser)
	assert.Equal(t, "24", ua.BrowserVersion)
}

func TestParseSamsungBrowser(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 13; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/23.0 Chrome/115.0.0.0 Mobile Safari/537.36")

	assert.Equal(t, "Samsung Browser", ua.Browser)
	assert.Equal(t, "23", ua.BrowserVersion)
}

func TestParseiPhoneSafari(t *testing.T) {
	ua := Parse("Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")

	assert.Equal(t, "Safari", ua.Browser)
	assert.Equal(t, "16", ua.BrowserVersion)
	assert.Equal(t, "iOS", ua.OS)
	assert.Equal(t, "16.6", ua.OSVersion)
	assert.Equal(t, "Apple", ua.DeviceVendor)
	assert.Equal(t, "iphone", ua.Device)
	assert.Equal(t, DeviceMobile, ua.DeviceType)
	assert.True(t, ua.IsMobile)
	assert.False(t, ua.IsTablet)
}

func TestParseiPad(t *testing.T) {
	ua := Parse("Mozilla/5.0 (iPad; CPU OS 15_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Mobile/15E148 Safari/604.1")

	assert.Equal(t, "Safari", ua.Browser)
	assert.Equal(t, "iOS", ua.OS)
	assert.Equal(t, "15.7", ua.OSVersion)
	assert.Equal(t, "Apple", ua.DeviceVendor)
	assert.Equal(t, "ipad", ua.Device)
	assert.Equal(t, DeviceTablet, ua.DeviceType)
	assert.False(t, ua.IsMobile)
	assert.True(t, ua.IsTablet)
}

func TestParseAndroidChrome(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 13; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Mobile Safari/537.36")

	assert.Equal(t, "Chrome", ua.Browser)
	assert.Equal(t, "119", ua.BrowserVersion)
	assert.Equal(t, "Android", ua.OS)
	assert.Equal(t, "13", ua.OSVersion)
	assert.Equal(t, "Samsung", ua.DeviceVendor)
	assert.Equal(t, DeviceMobile, ua.DeviceType)
	assert.True(t, ua.IsMobile)
}

func TestParseAndroidTablet(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 12; SM-T870) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	assert.Equal(t, "Android", ua.OS)
	assert.Equal(t, "12", ua.OSVersion)
	assert.Equal(t, "Samsung", ua.DeviceVendor)
	// 注意: 这个UA没有 "tablet" 或 "mobile" 关键词，但有 android，所以会被判断为 mobile
	assert.True(t, ua.IsMobile)
}

func TestParseAndroidTabletWithKeyword(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 12; SM-T870) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Tablet")

	assert.Equal(t, "Android", ua.OS)
	assert.Equal(t, DeviceTablet, ua.DeviceType)
	assert.False(t, ua.IsMobile)
	assert.True(t, ua.IsTablet)
}

func TestParseGooglebot(t *testing.T) {
	ua := Parse("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Googlebot", ua.BotName)
	assert.Equal(t, DeviceBot, ua.DeviceType)
}

func TestParseBingbot(t *testing.T) {
	ua := Parse("Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Bingbot", ua.BotName)
	assert.Equal(t, DeviceBot, ua.DeviceType)
}

func TestParseYandexBot(t *testing.T) {
	ua := Parse("Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "YandexBot", ua.BotName)
}

func TestParseBaiduSpider(t *testing.T) {
	ua := Parse("Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Baidu", ua.BotName)
}

func TestParseYahooSlurp(t *testing.T) {
	ua := Parse("Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Yahoo", ua.BotName)
}

func TestParseTwitterbot(t *testing.T) {
	ua := Parse("Twitterbot/1.0")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Twitterbot", ua.BotName)
}

func TestParseFacebookExternalHit(t *testing.T) {
	ua := Parse("facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "FacebookExternalHit", ua.BotName)
}

func TestParseApplebot(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) AppleBot/0.1")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Applebot", ua.BotName)
}

func TestParseGenericSpider(t *testing.T) {
	ua := Parse("Mozilla/5.0 (compatible; Spider/1.0)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Spider", ua.BotName)
}

func TestParseGenericCrawler(t *testing.T) {
	ua := Parse("Mozilla/5.0 (compatible; MyCrawler/1.0)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Crawler", ua.BotName)
}

func TestParseGenericBot(t *testing.T) {
	ua := Parse("Mozilla/5.0 (compatible; MyBot/1.0)")

	assert.True(t, ua.IsBot)
	assert.Equal(t, "Bot", ua.BotName)
}

func TestParseWindows81(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	assert.Equal(t, "Windows", ua.OS)
	assert.Equal(t, "8.1", ua.OSVersion)
}

func TestParseWindows8(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36")

	assert.Equal(t, "Windows", ua.OS)
	assert.Equal(t, "8", ua.OSVersion)
}

func TestParseWindows7(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36")

	assert.Equal(t, "Windows", ua.OS)
	assert.Equal(t, "7", ua.OSVersion)
}

func TestParseWindowsVista(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36")

	assert.Equal(t, "Windows", ua.OS)
	assert.Equal(t, "Vista", ua.OSVersion)
}

func TestParseWindowsXP(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36")

	assert.Equal(t, "Windows", ua.OS)
	assert.Equal(t, "XP", ua.OSVersion)
}

func TestParseHarmonyOS(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; HarmonyOS 3.0; NOH-AN00) AppleWebKit/537.36")

	assert.Equal(t, "OpenHarmony", ua.OS)
	assert.Equal(t, "3", ua.OSVersion)
}

func TestParseLinux(t *testing.T) {
	ua := Parse("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	assert.Equal(t, "Linux", ua.OS)
	assert.Empty(t, ua.OSVersion)
}

func TestParseWindowsPhone(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows Phone 10.0; Android 6.0.1) AppleWebKit/537.36")

	// 注意: 这个 UA 包含 Android 关键词，会优先匹配 Android
	assert.Equal(t, "Android", ua.OS)
	assert.Equal(t, "6.0", ua.OSVersion)
}

func TestParseChromeOS(t *testing.T) {
	ua := Parse("Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36")

	assert.Equal(t, "ChromeOS", ua.OS)
}

func TestParseFreeBSD(t *testing.T) {
	ua := Parse("Mozilla/5.0 (X11; FreeBSD amd64) AppleWebKit/537.36")

	assert.Equal(t, "FreeBSD", ua.OS)
}

func TestParseHuaweiDevice(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 12; ELS-AN00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Mobile Safari/537.36")

	assert.Equal(t, "Android", ua.OS)
	assert.True(t, ua.IsMobile)
	// Huawei device detection needs "huawei" keyword
}

func TestParseHuaweiWithKeyword(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 12; HUAWEI ELS-AN00) AppleWebKit/537.36")

	assert.Equal(t, "Huawei", ua.DeviceVendor)
}

func TestParseHonorDevice(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 12; HONOR X9a) AppleWebKit/537.36")

	assert.Equal(t, "Honor", ua.DeviceVendor)
}

func TestParseXiaomiDevice(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 13; Xiaomi 13 Pro) AppleWebKit/537.36")

	assert.Equal(t, "Xiaomi", ua.DeviceVendor)
}

func TestParseOPPODevice(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 13; OPPO Find X5 Pro) AppleWebKit/537.36")

	assert.Equal(t, "OPPO", ua.DeviceVendor)
}

func TestParseVivoDevice(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 13; vivo X90 Pro) AppleWebKit/537.36")

	assert.Equal(t, "Vivo", ua.DeviceVendor)
}

func TestParseMacDevice(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	assert.Equal(t, "macOS", ua.OS)
	assert.Equal(t, "Apple", ua.DeviceVendor)
	assert.Equal(t, "macintosh", ua.Device)
}

func TestParseNoVersionInBrowser(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome Safari/537.36")

	// Chrome keyword exists but no version number
	assert.Equal(t, "Chrome", ua.Browser)
	assert.Empty(t, ua.BrowserVersion)
}

func TestParseAndroidWithoutMobileKeyword(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 10; SM-T510) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	assert.Equal(t, "Android", ua.OS)
	// Android 关键词会触发 IsMobile
	assert.True(t, ua.IsMobile)
	assert.Equal(t, DeviceMobile, ua.DeviceType)
}

func TestParseComplexUA(t *testing.T) {
	ua := Parse("Mozilla/5.0 (Linux; Android 13; SM-S918B Build/TP1A.220624.014; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/119.0.6045.193 Mobile Safari/537.36")

	assert.Equal(t, "Chrome", ua.Browser)
	assert.Equal(t, "119", ua.BrowserVersion)
	assert.Equal(t, "Android", ua.OS)
	assert.Equal(t, "13", ua.OSVersion)
	assert.Equal(t, "Samsung", ua.DeviceVendor)
	assert.True(t, ua.IsMobile)
}
