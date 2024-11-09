/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 22:57:16
 * @FilePath: \go-toolbox\pkg\queue\queue.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package queue

import (
	"context"
)

// Queue 接口定义了队列的基本操作
type Queue interface {
	Enqueue(ctx context.Context, item interface{}) error
	Dequeue(ctx context.Context) (interface{}, error)
	IsEmpty() bool
	Size() int
}
