/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-30 20:21:57
 * @FilePath: \go-toolbox\pkg\httpx\utils.go
 * @Description: HTTP 相关工具
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package httpx

import "net/url"

func ParseUrlPath(urlString string) (path string) {
	var (
		err error
		u   *url.URL
	)
	if u, err = url.Parse(urlString); err != nil {
		return path
	}
	return u.Path
}
