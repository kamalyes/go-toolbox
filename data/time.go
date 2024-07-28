/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 11:55:53
 * @FilePath: \go-toolbox\data\time.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package data

import (
	"time"
)

// getTimeLoc
/**
 * @Description: 获取中国时区
 * @return l
 */
func getTimeLoc() (l *time.Location) {
	l, _ = time.LoadLocation("Asia/Shanghai")
	return
}

// NowStr
/**
 * @Description: 获取当时时间字符串（2006-04-02 15:04:05）
 * @return timeStr
 */
func NowStr() (timeStr string) {
	l := getTimeLoc()
	t := time.Now().In(l)
	timeStr = t.Format("2006-04-02 15:04:05")
	return
}
