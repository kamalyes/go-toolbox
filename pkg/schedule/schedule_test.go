/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 15:56:27
 * @FilePath: \go-toolbox\pkg\schedule\schedule_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSchedule_InitialState(t *testing.T) {
	s := NewSchedule()
	assert.NotNil(t, s)
	assert.NotNil(t, s.jobs)
	assert.NotNil(t, s.taskChan)
	assert.Equal(t, 20, s.workerCount)
	assert.False(t, s.IsBroken())
	assert.Equal(t, 100, cap(s.taskChan))
}

func TestGetJobRulesAndCopy(t *testing.T) {
	s := NewSchedule()

	// 初始应为空
	assert.Empty(t, s.GetJobRules())
	assert.Empty(t, s.GetJobRulesCopy())

	// 添加模拟任务
	job1 := &JobRule{id: "job1"}
	job2 := &JobRule{id: "job2"}

	s.jobs["job1"] = job1
	s.jobs["job2"] = job2

	// 获取任务规则
	rules := s.GetJobRules()
	assert.Len(t, rules, 2)
	assert.Contains(t, rules, job1)
	assert.Contains(t, rules, job2)

	// 获取副本
	rulesCopy := s.GetJobRulesCopy()
	assert.Len(t, rulesCopy, 2)
	assert.Contains(t, rulesCopy, job1)
	assert.Contains(t, rulesCopy, job2)

	// 确认副本是新切片（地址不同）
	assert.NotSame(t, &rules, &rulesCopy)
}

func TestIsBrokenAndWorkerCount(t *testing.T) {
	s := NewSchedule()

	assert.False(t, s.IsBroken())
	assert.Equal(t, 20, s.GetWorkerCount())

	// 模拟修改状态
	s.mu.Lock()
	s.isBroken = true
	s.workerCount = 42
	s.mu.Unlock()

	assert.True(t, s.IsBroken())
	assert.Equal(t, 42, s.GetWorkerCount())
}

func TestSetTaskChanCapacity(t *testing.T) {
	s := NewSchedule()

	oldChan := s.taskChan
	oldCap := cap(oldChan)
	assert.Equal(t, 100, oldCap)

	newCapacity := 50
	s2 := s.SetTaskChanCapacity(newCapacity)
	assert.Equal(t, s, s2)

	newChan := s.taskChan
	assert.NotEqual(t, oldChan, newChan)
	assert.Equal(t, newCapacity, cap(newChan))
}

func TestSetWorkerCount(t *testing.T) {
	s := NewSchedule()
	newCount := 10
	s2 := s.SetWorkerCount(newCount)
	assert.Equal(t, s, s2)
	assert.Equal(t, newCount, s.GetWorkerCount())
}

func TestValidate(t *testing.T) {
	s := NewSchedule()

	// 0任务时应返回错误
	err := s.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrAtLeastOneRuleMustBeDefined)

	// 添加任务后应无错误
	s.jobs["job1"] = &JobRule{id: "job1"}
	err = s.Validate()
	assert.NoError(t, err)
}

func TestConcurrentAccess(t *testing.T) {
	s := NewSchedule()

	const goroutines = 50
	const iterations = 100

	var wg sync.WaitGroup

	// 并发设置工作池大小
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(count int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.SetWorkerCount(count)
				_ = s.GetWorkerCount()
			}
		}(i)
	}

	// 并发设置任务通道容量
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(capacity int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.SetTaskChanCapacity(capacity + 1)
			}
		}(i)
	}

	// 并发读取IsBroken状态
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = s.IsBroken()
			}
		}()
	}

	wg.Wait()
}
