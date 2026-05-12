/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-13 13:19:57
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 14:42:55
 * @FilePath: \go-toolbox\pkg\stringx\conv.go
 * @Description: 字符串转换工具
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"strconv"
	"time"
	"unsafe"
)

// FastItoa 整数转字符串
func FastItoa(val int) string {
	if val >= 0 && val < 100 {
		return intStrings0To99[val]
	}

	if val >= 100 && val < 10000 {
		return intStrings100To9999[val-100]
	}

	return strconv.Itoa(val)
}

// FastFloat 浮点数转字符串
func FastFloat(val float64, prec int) string {
	return strconv.FormatFloat(val, 'f', prec, 64)
}

// FastAppendInt 将整数追加到缓冲区
func FastAppendInt(buf []byte, val int) []byte {
	if val == 0 {
		return append(buf, '0')
	}

	if val < 0 {
		buf = append(buf, '-')
		val = -val
	}

	if val < 10 {
		return append(buf, byte('0'+val))
	}

	if val < 100 {
		return append(buf, byte('0'+val/10), byte('0'+val%10))
	}

	if val < 1000 {
		return append(buf, byte('0'+val/100), byte('0'+(val/10)%10), byte('0'+val%10))
	}

	return strconv.AppendInt(buf, int64(val), 10)
}

var (
	intStrings0To99     [100]string
	intStrings100To9999 [9900]string
)

// Init 初始化整数字符串缓存
func init() {
	for i := 0; i < 10; i++ {
		buf := [1]byte{byte('0' + i)}
		intStrings0To99[i] = unsafe.String(&buf[0], 1)
	}

	for i := 10; i < 100; i++ {
		buf := [2]byte{
			byte('0' + i/10),
			byte('0' + i%10),
		}
		intStrings0To99[i] = unsafe.String(&buf[0], 2)
	}

	for i := 0; i < 9900; i++ {
		val := i + 100
		if val < 1000 {
			buf := [3]byte{
				byte('0' + val/100),
				byte('0' + (val/10)%10),
				byte('0' + val%10),
			}
			intStrings100To9999[i] = unsafe.String(&buf[0], 3)
		} else {
			buf := [4]byte{
				byte('0' + val/1000),
				byte('0' + (val/100)%10),
				byte('0' + (val/10)%10),
				byte('0' + val%10),
			}
			intStrings100To9999[i] = unsafe.String(&buf[0], 4)
		}
	}
}

// FastFormatTime 格式化时间到缓冲区
func FastFormatTime(buf []byte, t time.Time) []byte {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	buf = FastAppendInt(buf, year)
	buf = append(buf, '/')

	buf = FastAppendInt(buf, int(month))
	buf = append(buf, '/')

	buf = FastAppendInt(buf, day)
	buf = append(buf, ' ')

	buf = FastAppendInt(buf, hour)
	buf = append(buf, ':')

	if min < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, min)
	buf = append(buf, ':')

	if sec < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, sec)
	buf = append(buf, ' ')

	return buf
}

// FastFormatTimeISO 格式化时间到缓冲区，ISO格式
func FastFormatTimeISO(buf []byte, t time.Time) []byte {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	buf = FastAppendInt(buf, year)
	buf = append(buf, '-')

	if month < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, int(month))
	buf = append(buf, '-')

	if day < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, day)
	buf = append(buf, ' ')

	if hour < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, hour)
	buf = append(buf, ':')

	if min < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, min)
	buf = append(buf, ':')

	if sec < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, sec)

	return buf
}

// FastFormatTimeCompact 格式化时间到缓冲区，紧凑格式
func FastFormatTimeCompact(buf []byte, t time.Time) []byte {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	buf = FastAppendInt(buf, year)

	if month < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, int(month))

	if day < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, day)

	if hour < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, hour)

	if min < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, min)

	if sec < 10 {
		buf = append(buf, '0')
	}
	buf = FastAppendInt(buf, sec)

	return buf
}
