/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 20:09:05
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 19:10:32
 * @FilePath: \go-toolbox\pkg\osx\console.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/moment"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type LogFormat int

const (
	BasicLogFormat LogFormat = iota
	AdditionalCallerLogFormat
	AdditionalTimeLogFormat
)

// Console 接口定义了日志系统的方法
type Console interface {
	// 成功日志
	Success(format string, a ...interface{})
	// 信息日志
	Info(format string, a ...interface{})
	// 调试日志
	Debug(format string, a ...interface{})
	// 警告日志
	Warning(format string, a ...interface{})
	// 错误日志
	Error(format string, a ...interface{})
	// 致命错误日志并退出程序
	Fatalln(format string, a ...interface{})
	// 打印完成标记
	MarkDone()
	// 如果传入错误不为nil，则打印错误日志并退出程序
	Must(err error)
	// 设置日志等级
	SetLogLevel(level LogLevel)
	// 设置日志风格
	SetLogFormat(level LogFormat)
	// 转换Json风格
	ConvertJsonFormat(flag bool)
	// 上下文日志
	LogWithContext(ctx context.Context, level LogLevel, format string, a ...interface{})
}

// colorConsole 结构体实现了Console接口，支持彩色日志输出
type colorConsole struct {
	mu     sync.Mutex
	enable bool // 控制是否启用日志输出
	level  LogLevel
	format LogFormat
	isJson bool
}

type LogEntry struct {
	Timestamp  string         `json:"timestamp"`
	Level      LogLevel       `json:"level"`
	Message    string         `json:"message"`
	CallerInfo *RunTimeCaller `json:"caller_info"`
	RequestId  string         `json:"request_id"`
}

// NewColorConsole 返回一个colorConsole实例，支持彩色日志输出
// enable 参数用于控制是否启用日志输出，默认为 true
// level 参数用于控制日志级别，默认为 INFO
// format 参数用于控制日志格式，默认为 BasicLogFormat
func NewColorConsole(enable ...bool) Console {
	logEnable := true
	if len(enable) > 0 {
		logEnable = enable[0]
	}
	return &colorConsole{
		enable: logEnable,
		level:  INFO,
		format: BasicLogFormat,
	}
}

// SetLogLevel 设置日志级别
func (c *colorConsole) SetLogLevel(level LogLevel) {
	c.level = level
}

// SetLogFormat 设置日志格式
func (c *colorConsole) SetLogFormat(format LogFormat) {
	c.format = format
}

// ConvertJsonFormat 转换Json格式
func (c *colorConsole) ConvertJsonFormat(flag bool) {
	c.isJson = flag
}

// formatString 自定义字符串格式化函数
func formatString(format string, a ...interface{}) (string, error) {
	var result strings.Builder
	var i int
	args := a // 避免在循环中多次访问切片长度

	for i < len(format) {
		idx := strings.IndexRune(format[i:], '%')
		if idx == -1 {
			result.WriteString(format[i:])
			break
		}
		result.WriteString(format[i : i+idx])
		i += idx + 1

		if len(args) == 0 {
			return "", fmt.Errorf("too few arguments in format call")
		}

		str := convert.MustString(args[0])
		args = args[1:]

		result.WriteString(str)
	}

	if len(args) > 0 {
		return "", fmt.Errorf("too many arguments in format call")
	}

	return result.String(), nil
}

// printLog 是一个通用的日志打印函数，增加了时间和函数名称前缀
func (c *colorConsole) printLog(level LogLevel, colorFunc func(string, ...interface{}) string, format string, a ...interface{}) {
	if !c.enable || level < c.level {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	// 使用自定义的格式化函数格式化日志消息
	message, err := formatString(format, a...)
	if err != nil {
		// 如果格式化错误，则直接打印错误（这里简单处理为打印到标准错误，并退出程序）
		fmt.Fprintf(os.Stderr, "Error formatting log message: %v\n", err)
		os.Exit(1)
	}

	// 根据日志格式添加额外信息
	var prefixedMessage string
	var callerInfo *RunTimeCaller
	_, _, currentTime := moment.GetCurrentTimeInfo()
	currentDataTime := currentTime.Format(time.DateTime)
	switch c.format {
	case AdditionalCallerLogFormat:
		callerInfo = GetCallerInfo(2)
		prefixedMessage = fmt.Sprintf("[%s] [%s:%d] %s: %s", currentDataTime, callerInfo.File, callerInfo.Line, callerInfo.FuncName, message)
	case AdditionalTimeLogFormat:
		callerInfo = GetCallerInfo(2)
		prefixedMessage = fmt.Sprintf("[%s] %s", currentDataTime, message)
	default:
		prefixedMessage = message
	}

	// 处理颜色函数（如果提供了颜色函数）
	var coloredMessage string
	if colorFunc != nil {
		coloredMessage = colorFunc(prefixedMessage)
	} else {
		coloredMessage = prefixedMessage // 如果没有提供颜色函数，则直接使用原始消息
	}

	if c.isJson {
		// JSON 格式处理（保持不变）
		logEntry := LogEntry{
			Timestamp:  currentDataTime,
			Level:      level,
			Message:    message,
			CallerInfo: callerInfo,
		}
		jsonMessage, err := json.Marshal(logEntry)
		if err != nil {
			// 如果JSON序列化错误，则直接打印错误（这里简单处理为打印到标准错误，并退出程序）
			fmt.Fprintf(os.Stderr, "Error serializing JSON log message: %v\n", err)
			os.Exit(1)
		}
		coloredMessage = string(jsonMessage)
	}

	os.Stdout.Write([]byte(coloredMessage + "\n"))
	os.Stdout.Sync() // 确保缓冲区被刷新
}

// 实现 Console 接口的方法

// Info 打印信息日志（无颜色）
func (c *colorConsole) Info(format string, a ...interface{}) {
	c.printLog(INFO, func(s string, args ...interface{}) string { return s }, format, a...)
}

// Debug 打印调试日志（青色）
func (c *colorConsole) Debug(format string, a ...interface{}) {
	c.printLog(DEBUG, color.LightCyan.Sprintf, format, a...)
}

// Success 打印成功日志（绿色）
func (c *colorConsole) Success(format string, a ...interface{}) {
	c.printLog(INFO, color.LightGreen.Sprintf, format, a...)
}

// Warning 打印警告日志（黄色）
func (c *colorConsole) Warning(format string, a ...interface{}) {
	c.printLog(WARNING, color.LightYellow.Sprintf, format, a...)
}

// Error 打印错误日志（红色）
func (c *colorConsole) Error(format string, a ...interface{}) {
	c.printLog(ERROR, color.LightRed.Sprintf, format, a...)
}

// Fatalln 打印致命错误日志并退出程序（红色）
func (c *colorConsole) Fatalln(format string, a ...interface{}) {
	c.printLog(FATAL, color.LightRed.Sprintf, format, a...)
	os.Exit(1)
}

// MarkDone 打印完成标记（绿色）
func (c *colorConsole) MarkDone() {
	c.Success("Done.")
}

// Must 如果传入错误不为nil，则打印错误日志并退出程序（红色）
func (c *colorConsole) Must(err error) {
	if err != nil {
		c.Fatalln("%+v", err)
	}
}

// ContextKey 是一个类型，用于创建上下文中唯一的键
type ContextKey string

// 定义一个键来从上下文中提取信息（例如，请求ID）
const RequestIDKey ContextKey = "request_id"

// LogWithContext 方法接受一个上下文，并根据上下文和日志级别打印日志
func (c *colorConsole) LogWithContext(ctx context.Context, level LogLevel, format string, a ...interface{}) {
	// 从上下文中提取请求ID（或其他信息）
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok || strings.TrimSpace(requestID) == "" {
		requestID = "<no-request-id>" // 如果没有找到或为空，则使用特殊占位符
	}

	// 构造包含请求ID的日志消息
	formattedMessage := fmt.Sprintf("[%s] %s", requestID, format)

	// 使用现有的printLog方法打印日志
	c.printLog(level, nil, formattedMessage, a...) // 注意这里我们不需要再封装颜色函数，因为printLog内部会根据level处理
}
