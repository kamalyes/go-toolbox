# User-Agent 使用指南

User-Agent 解析器和生成器，支持浏览器识别、设备检测、爬虫识别以及随机 UA 生成。

## 快速开始

### 基础用法

```go
import "github.com/kamalyes/go-toolbox/pkg/useragent"

// 解析 User-Agent
ua := useragent.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")

fmt.Println("浏览器:", ua.Browser)          // Chrome
fmt.Println("操作系统:", ua.OS)             // Windows
fmt.Println("设备类型:", ua.DeviceType)     // desktop
```

## 主要功能

### 1. User-Agent 解析

```go
parsed := useragent.Parse(userAgentString)

// 浏览器信息
parsed.Browser           // Chrome, Firefox, Safari, Edge 等
parsed.BrowserVersion    // 浏览器主版本号

// 操作系统信息
parsed.OS                // Windows, macOS, iOS, Android, Linux 等
parsed.OSVersion         // 操作系统版本号

// 设备信息
parsed.Device            // 设备名称
parsed.DeviceType        // mobile, tablet, desktop, bot
parsed.DeviceVendor      // Apple, Samsung, Huawei 等
parsed.IsMobile          // 是否为移动设备
parsed.IsTablet          // 是否为平板设备

// 爬虫检测
parsed.IsBot             // 是否为爬虫
parsed.BotName           // 爬虫名称: Googlebot, Bingbot 等
```

### 2. 设备类型检测

支持识别的浏览器：
- Chrome, Firefox, Safari, Edge, Opera
- Samsung Browser, Yandex Browser
- 移动端浏览器

支持识别的操作系统：
- Windows (XP, Vista, 7, 8, 8.1, 10)
- macOS
- iOS (iPhone, iPad)
- Android
- Linux, ChromeOS, FreeBSD
- HarmonyOS, Windows Phone

支持识别的爬虫：
- Googlebot, Bingbot, YandexBot
- Baiduspider, Yahoo Slurp
- Twitterbot, FacebookExternalHit
- Applebot

### 3. 随机 UA 生成器

```go
// 创建生成器
gen := useragent.New()

// 生成随机 UA
gen.GenerateRand()
fullUA := gen.GetFullValue()

// 获取信息
browser := gen.GetName()           // 浏览器名称
version := gen.GetFullVersion()    // 完整版本号
os := gen.GetOS()                  // 操作系统
osVersion := gen.GetFullOSVersion() // OS 完整版本
```

### 4. 稳定 UA 生成

```go
gen := useragent.New()

// 生成特定设备类型的稳定 UA
desktopUA := gen.GenerateStabilize(useragent.DeviceTypeDesktop)
mobileUA := gen.GenerateStabilize(useragent.DeviceTypeMobile)
tabletUA := gen.GenerateStabilize(useragent.DeviceTypeTablet)
```

### 5. 设置生成类型

```go
gen := useragent.New()

// 设置生成类型
gen.SetRgType(useragent.RgTypePopular)    // 流行浏览器
gen.SetRgType(useragent.RgTypeAll)        // 所有浏览器
gen.SetRgType(useragent.RgTypeUnpopular)  // 非主流浏览器

gen.GenerateRand()
```

### 6. 操作系统判断

```go
gen := useragent.New()
gen.GenerateRand()

// 判断操作系统类型
if gen.IsWindows() {
    // Windows 系统
}
if gen.IsMacOS() {
    // macOS 系统
}
if gen.IsAndroid() {
    // Android 系统
}
if gen.IsIOS() {
    // iOS 系统
}
```

## 应用场景

### 爬虫识别

```go
func DetectBot(userAgent string) bool {
    parsed := useragent.Parse(userAgent)
    if parsed.IsBot {
        log.Printf("检测到爬虫: %s", parsed.BotName)
        return true
    }
    return false
}
```

### 设备适配

```go
func GetTemplate(userAgent string) string {
    parsed := useragent.Parse(userAgent)
    
    switch parsed.DeviceType {
    case useragent.DeviceMobile:
        return "mobile.html"
    case useragent.DeviceTablet:
        return "tablet.html"
    default:
        return "desktop.html"
    }
}
```

### 浏览器兼容性检测

```go
func CheckBrowserSupport(userAgent string) bool {
    parsed := useragent.Parse(userAgent)
    
    // 检查是否为现代浏览器
    if parsed.Browser == "Chrome" || parsed.Browser == "Firefox" {
        version := convert.MustInt(parsed.BrowserVersion)
        return version >= 90
    }
    return false
}
```

### 爬虫模拟 (测试)

```go
func CreateCrawler() *http.Client {
    gen := useragent.New()
    ua := gen.GenerateStabilize(useragent.DeviceTypeDesktop)
    
    client := &http.Client{
        Transport: &http.Transport{
            // 配置
        },
    }
    
    // 在请求中使用生成的 UA
    req.Header.Set("User-Agent", ua)
    return client
}
```

## 常量定义

### 浏览器常量
```go
Chrome, Firefox, Safari, Edge, Opera
SamsungBrowser, YandexBrowser
Googlebot, Bingbot, YandexBot, 等
```

### 操作系统常量
```go
Windows, MacOS, Android, IOS
Linux, FreeBSD, ChromeOS
OpenHarmony, WindowsPhone
```

### 设备类型常量
```go
DeviceMobile   = "mobile"
DeviceTablet   = "tablet"
DeviceDesktop  = "desktop"
DeviceBot      = "bot"
```

### 设备厂商常量
```go
VendorApple, VendorSamsung, VendorHuawei
VendorHonor, VendorXiaomi, VendorOPPO, VendorVivo
```

## 测试覆盖率

- **覆盖率**: 99%+
- **测试场景**: 覆盖主流浏览器、操作系统、设备、爬虫
- **并发安全**: 所有公共 API 都是并发安全的
