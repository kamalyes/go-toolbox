/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-23 23:50:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-23 23:50:00
 * @FilePath: \go-toolbox\pkg\breaker\errors.go
 * @Description: 错误定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package breaker

import "errors"

var (
	// ErrOpen 熔断器打开错误
	ErrOpen = errors.New("circuit breaker is open")

	// ErrRateLimitExceeded 限流错误
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)
