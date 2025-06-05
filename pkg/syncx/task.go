/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-05 16:25:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-05 16:25:58
 * @FilePath: \go-toolbox\pkg\syncx\task.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"fmt"
	"sync"
)

// 泛型任务函数
type TaskFunc[T any] func() (T, error)

// 任务结果
type TaskResult[T any] struct {
	Result T
	Err    error
}

// 泛型任务管理器，key->任务函数映射
type TaskRunner[T any] struct {
	tasks map[string]TaskFunc[T]
	mu    sync.Mutex // 保护 tasks 添加操作，防止 Add 并发写冲突
}

// 新建任务管理器
func NewTaskRunner[T any]() *TaskRunner[T] {
	return &TaskRunner[T]{
		tasks: make(map[string]TaskFunc[T]),
	}
}

// 链式添加任务，传入唯一标识符key
func (tr *TaskRunner[T]) Add(key string, task TaskFunc[T]) *TaskRunner[T] {
	WithLock(&tr.mu, func() {
		tr.tasks[key] = task
	})
	return tr
}

// 并发执行所有任务，返回 key->结果映射
func (tr *TaskRunner[T]) Run() map[string]TaskResult[T] {
	wg := NewWaitGroup(true)
	results := make(map[string]TaskResult[T])
	var mu sync.Mutex

	for key, task := range tr.tasks {
		key := key
		task := task
		wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("panic: %v", r)
					WithLock(&mu, func() {
						results[key] = TaskResult[T]{Err: err}
					})
				}
			}()

			res, err := task()
			WithLock(&mu, func() {
				results[key] = TaskResult[T]{Result: res, Err: err}
			})
		})
	}

	if err := wg.Wait(); err != nil {
		fmt.Println("Run encountered error:", err)
	}
	return results
}
