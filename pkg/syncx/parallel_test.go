/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 08:52:55
 * @FilePath: \go-toolbox\pkg\syncx\parallel_test.go
 * @Description: 并发执行工具函数测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 对比 WaitGroup vs Channel 的性能
func BenchmarkWaitGroupVsChannel(b *testing.B) {
	data := make(map[int]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	b.Run("WaitGroup", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for k, v := range data {
				wg.Add(1)
				go func(key, val int) {
					defer wg.Done()
					_ = key + val
				}(k, v)
			}
			wg.Wait()
		}
	})

	b.Run("BufferedChannel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			done := make(chan struct{}, len(data))
			for k, v := range data {
				go func(key, val int) {
					_ = key + val
					done <- struct{}{}
				}(k, v)
			}
			for j := 0; j < len(data); j++ {
				<-done
			}
			close(done)
		}
	})
}

// 小数据集测试
func BenchmarkSmallDataSet(b *testing.B) {
	data := map[int]int{1: 1, 2: 2, 3: 3}

	b.Run("WaitGroup-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for k, v := range data {
				wg.Add(1)
				go func(key, val int) {
					defer wg.Done()
					_ = key + val
				}(k, v)
			}
			wg.Wait()
		}
	})

	b.Run("Channel-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			done := make(chan struct{}, len(data))
			for k, v := range data {
				go func(key, val int) {
					_ = key + val
					done <- struct{}{}
				}(k, v)
			}
			for j := 0; j < len(data); j++ {
				<-done
			}
			close(done)
		}
	})
}

// 大数据集测试
func BenchmarkLargeDataSet(b *testing.B) {
	data := make(map[int]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.Run("WaitGroup-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for k, v := range data {
				wg.Add(1)
				go func(key, val int) {
					defer wg.Done()
					_ = key + val
				}(k, v)
			}
			wg.Wait()
		}
	})

	b.Run("Channel-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			done := make(chan struct{}, len(data))
			for k, v := range data {
				go func(key, val int) {
					_ = key + val
					done <- struct{}{}
				}(k, v)
			}
			for j := 0; j < len(data); j++ {
				<-done
			}
			close(done)
		}
	})
}

func TestParallelForEach(t *testing.T) {
	t.Run("空map", func(t *testing.T) {
		var m map[string]int
		ParallelForEach(m, func(k string, v int) {
			t.Error("不应该执行")
		})
	})

	t.Run("正常执行", func(t *testing.T) {
		m := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}

		var counter atomic.Int32
		ParallelForEach(m, func(k string, v int) {
			counter.Add(1)
			time.Sleep(10 * time.Millisecond) // 模拟耗时操作
		})

		if counter.Load() != 3 {
			t.Errorf("期望执行3次，实际执行%d次", counter.Load())
		}
	})

	t.Run("并发安全", func(t *testing.T) {
		m := make(map[int]int, 100)
		for i := 0; i < 100; i++ {
			m[i] = i
		}

		var mu sync.Mutex
		results := make(map[int]int)

		ParallelForEach(m, func(k int, v int) {
			mu.Lock()
			results[k] = v * 2
			mu.Unlock()
		})

		if len(results) != 100 {
			t.Errorf("期望100个结果，实际%d个", len(results))
		}

		for k, v := range results {
			if v != k*2 {
				t.Errorf("键%d期望值%d，实际值%d", k, k*2, v)
			}
		}
	})

	t.Run("确保所有goroutine完成", func(t *testing.T) {
		m := map[string]int{
			"1": 1,
			"2": 2,
			"3": 3,
		}

		var counter atomic.Int32
		start := time.Now()

		ParallelForEach(m, func(k string, v int) {
			time.Sleep(100 * time.Millisecond)
			counter.Add(1)
		})

		duration := time.Since(start)

		// 由于是并发执行，总时间应该接近100ms而不是300ms
		if duration > 200*time.Millisecond {
			t.Errorf("并发执行耗时过长: %v", duration)
		}

		if counter.Load() != 3 {
			t.Errorf("期望3个任务完成，实际%d个", counter.Load())
		}
	})
}

func TestParallelForEachSlice(t *testing.T) {
	t.Run("空slice", func(t *testing.T) {
		var s []int
		ParallelForEachSlice(s, func(i int, v int) {
			t.Error("不应该执行")
		})
	})

	t.Run("正常执行", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}

		var counter atomic.Int32
		ParallelForEachSlice(s, func(i int, v int) {
			counter.Add(1)
		})

		if counter.Load() != 5 {
			t.Errorf("期望执行5次，实际执行%d次", counter.Load())
		}
	})

	t.Run("索引和值正确", func(t *testing.T) {
		s := []string{"a", "b", "c"}

		var mu sync.Mutex
		results := make(map[int]string)

		ParallelForEachSlice(s, func(i int, v string) {
			mu.Lock()
			results[i] = v
			mu.Unlock()
		})

		if len(results) != 3 {
			t.Errorf("期望3个结果，实际%d个", len(results))
		}

		expected := map[int]string{0: "a", 1: "b", 2: "c"}
		for k, v := range expected {
			if results[k] != v {
				t.Errorf("索引%d期望值%s，实际值%s", k, v, results[k])
			}
		}
	})

	t.Run("大量数据", func(t *testing.T) {
		size := 1000
		s := make([]int, size)
		for i := 0; i < size; i++ {
			s[i] = i
		}

		var counter atomic.Int32
		ParallelForEachSlice(s, func(i int, v int) {
			counter.Add(1)
		})

		if counter.Load() != int32(size) {
			t.Errorf("期望执行%d次，实际执行%d次", size, counter.Load())
		}
	})

	t.Run("并发修改安全", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}

		var mu sync.Mutex
		sum := 0

		ParallelForEachSlice(s, func(i int, v int) {
			mu.Lock()
			sum += v
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
		})

		expected := 15 // 1+2+3+4+5
		if sum != expected {
			t.Errorf("期望总和%d，实际%d", expected, sum)
		}
	})
}

// 基准测试
func BenchmarkParallelForEach(b *testing.B) {
	m := make(map[int]int, 100)
	for i := 0; i < 100; i++ {
		m[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelForEach(m, func(k int, v int) {
			_ = k + v
		})
	}
}

func BenchmarkParallelForEachSlice(b *testing.B) {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelForEachSlice(s, func(i int, v int) {
			_ = i + v
		})
	}
}

// 对比顺序执行的性能
func BenchmarkSequentialForEach(b *testing.B) {
	m := make(map[int]int, 100)
	for i := 0; i < 100; i++ {
		m[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k, v := range m {
			_ = k + v
		}
	}
}

// TestParallelExecutor_Callbacks 测试回调风格的 API
func TestParallelExecutor_Callbacks(t *testing.T) {
	t.Run("OnSuccess回调", func(t *testing.T) {
		data := map[string]int{"a": 1, "b": 2, "c": 3}
		var successCount int32

		NewParallelExecutor[string, int, int](data).
			OnSuccess(func(key string, val int, result int) {
				atomic.AddInt32(&successCount, 1)
				t.Logf("成功: key=%s, val=%d, result=%d", key, val, result)
			}).
			Execute(func(key string, val int) (int, error) {
				return val * 10, nil
			})

		if successCount != 3 {
			t.Errorf("期望成功回调 3 次, 实际 %d 次", successCount)
		}
	})

	t.Run("OnError回调", func(t *testing.T) {
		data := map[string]int{"a": 1, "b": 2, "c": 3}
		var errorCount int32

		NewParallelExecutor[string, int, int](data).
			OnError(func(key string, val int, err error) {
				atomic.AddInt32(&errorCount, 1)
				t.Logf("错误: key=%s, val=%d, err=%v", key, val, err)
			}).
			Execute(func(key string, val int) (int, error) {
				if val == 2 {
					return 0, fmt.Errorf("值 %d 失败", val)
				}
				return val * 10, nil
			})

		if errorCount != 1 {
			t.Errorf("期望错误回调 1 次, 实际 %d 次", errorCount)
		}
	})

	t.Run("OnComplete回调", func(t *testing.T) {
		data := map[string]int{"a": 1, "b": 2, "c": 3}
		var completed bool

		NewParallelExecutor[string, int, int](data).
			OnComplete(func(results map[string]int, errors map[string]error) {
				completed = true
				t.Logf("完成: results=%d个, errors=%d个", len(results), len(errors))
			}).
			Execute(func(key string, val int) (int, error) {
				return val * 10, nil
			})

		if !completed {
			t.Error("OnComplete 回调未执行")
		}
	})

	t.Run("链式调用所有回调", func(t *testing.T) {
		data := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		var (
			successCount int32
			errorCount   int32
			eachCount    int32
			completed    bool
		)

		NewParallelExecutor[string, int, int](data).
			OnSuccess(func(key string, val int, result int) {
				atomic.AddInt32(&successCount, 1)
			}).
			OnError(func(key string, val int, err error) {
				atomic.AddInt32(&errorCount, 1)
			}).
			OnEachComplete(func(key string) {
				atomic.AddInt32(&eachCount, 1)
			}).
			OnComplete(func(results map[string]int, errors map[string]error) {
				completed = true
				if len(results) != 2 {
					t.Errorf("期望 2 个成功结果, 得到 %d 个", len(results))
				}
				if len(errors) != 2 {
					t.Errorf("期望 2 个错误, 得到 %d 个", len(errors))
				}
			}).
			Execute(func(key string, val int) (int, error) {
				if val%2 == 0 {
					return 0, fmt.Errorf("偶数错误: %d", val)
				}
				return val * 10, nil
			})

		if successCount != 2 {
			t.Errorf("期望成功 2 次, 实际 %d 次", successCount)
		}
		if errorCount != 2 {
			t.Errorf("期望错误 2 次, 实际 %d 次", errorCount)
		}
		if eachCount != 4 {
			t.Errorf("期望每个完成 4 次, 实际 %d 次", eachCount)
		}
		if !completed {
			t.Error("OnComplete 未执行")
		}
	})
}

// TestParallelSliceExecutor_Callbacks 测试 Slice 执行器的回调
func TestParallelSliceExecutor_Callbacks(t *testing.T) {
	t.Run("OnSuccess回调", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		var successCount int32

		NewParallelSliceExecutor[int, int](data).
			OnSuccess(func(idx int, val int, result int) {
				atomic.AddInt32(&successCount, 1)
				t.Logf("成功: idx=%d, val=%d, result=%d", idx, val, result)
			}).
			Execute(func(idx int, val int) (int, error) {
				return val * 10, nil
			})

		if successCount != 5 {
			t.Errorf("期望成功回调 5 次, 实际 %d 次", successCount)
		}
	})

	t.Run("OnComplete获取结果", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}

		NewParallelSliceExecutor[int, int](data).
			OnComplete(func(results []int, errors []error) {
				// 验证结果保持顺序
				expected := []int{10, 20, 30, 40, 50}
				for i, v := range expected {
					if results[i] != v {
						t.Errorf("索引 %d: 期望 %d, 得到 %d", i, v, results[i])
					}
				}
			}).
			Execute(func(idx int, val int) (int, error) {
				time.Sleep(time.Millisecond * time.Duration(6-val)) // 模拟不同执行时间
				return val * 10, nil
			})
	})

	t.Run("混合成功和失败", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		var (
			successCount int32
			errorCount   int32
		)

		NewParallelSliceExecutor[int, int](data).
			OnSuccess(func(idx int, val int, result int) {
				atomic.AddInt32(&successCount, 1)
			}).
			OnError(func(idx int, val int, err error) {
				atomic.AddInt32(&errorCount, 1)
			}).
			OnComplete(func(results []int, errors []error) {
				// 验证偶数索引有错误
				for i, err := range errors {
					if data[i]%2 == 0 {
						if err == nil {
							t.Errorf("索引 %d (偶数) 应该有错误", i)
						}
					} else {
						if err != nil {
							t.Errorf("索引 %d (奇数) 不应该有错误", i)
						}
					}
				}
			}).
			Execute(func(idx int, val int) (int, error) {
				if val%2 == 0 {
					return 0, fmt.Errorf("偶数错误: %d", val)
				}
				return val * 10, nil
			})

		if successCount != 3 {
			t.Errorf("期望成功 3 次, 实际 %d 次", successCount)
		}
		if errorCount != 2 {
			t.Errorf("期望错误 2 次, 实际 %d 次", errorCount)
		}
	})
}

// 示例: 真实场景 - 批量发送消息
func ExampleParallelExecutor() {
	clients := map[string]string{
		"user1": "client1@example.com",
		"user2": "client2@example.com",
		"user3": "client3@example.com",
	}

	NewParallelExecutor[string, string, bool](clients).
		OnSuccess(func(userID string, email string, sent bool) {
			fmt.Printf("✓ 消息已发送到 %s (%s)\n", userID, email)
		}).
		OnError(func(userID string, email string, err error) {
			fmt.Printf("✗ 发送失败 %s: %v\n", userID, err)
		}).
		OnComplete(func(results map[string]bool, errors map[string]error) {
			fmt.Printf("\n总计: 成功 %d, 失败 %d\n", len(results), len(errors))
		}).
		Execute(func(userID string, email string) (bool, error) {
			// 模拟发送消息
			if userID == "user2" {
				return false, fmt.Errorf("网络超时")
			}
			return true, nil
		})

	// Unordered output:
	// ✓ 消息已发送到 user1 (client1@example.com)
	// ✗ 发送失败 user2: 网络超时
	// ✓ 消息已发送到 user3 (client3@example.com)
	//
	// 总计: 成功 2, 失败 1
}

// 示例: Slice 场景 - 批量处理数据
func ExampleParallelSliceExecutor() {
	data := []int{10, 20, 30, 40, 50}

	NewParallelSliceExecutor[int, int](data).
		OnComplete(func(results []int, errors []error) {
			// 在 OnComplete 中按顺序输出,保证输出顺序
			for i, result := range results {
				fmt.Printf("索引 %d: %d -> %d\n", i, data[i], result)
			}
			fmt.Println("所有任务完成!")
		}).
		Execute(func(idx int, val int) (int, error) {
			return val * 2, nil
		})

	// Output:
	// 索引 0: 10 -> 20
	// 索引 1: 20 -> 40
	// 索引 2: 30 -> 60
	// 索引 3: 40 -> 80
	// 索引 4: 50 -> 100
	// 所有任务完成!
}

// TestParallelSliceExecutor_WithConcurrency 测试并发控制功能
func TestParallelSliceExecutor_WithConcurrency(t *testing.T) {
	t.Run("限制并发数为1", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		var (
			mu             sync.Mutex
			currentRunning int32
			maxRunning     int32
		)

		NewParallelSliceExecutor[int, int](data).
			WithConcurrency(1).
			Execute(func(idx int, val int) (int, error) {
				current := atomic.AddInt32(&currentRunning, 1)

				// 记录最大并发数
				mu.Lock()
				if current > maxRunning {
					maxRunning = current
				}
				mu.Unlock()

				time.Sleep(10 * time.Millisecond) // 模拟耗时操作
				atomic.AddInt32(&currentRunning, -1)
				return val * 2, nil
			})

		if maxRunning != 1 {
			t.Errorf("期望最大并发数为 1，实际为 %d", maxRunning)
		}
	})

	t.Run("限制并发数为3", func(t *testing.T) {
		data := make([]int, 10)
		for i := range data {
			data[i] = i + 1
		}

		var (
			mu             sync.Mutex
			currentRunning int32
			maxRunning     int32
		)

		NewParallelSliceExecutor[int, int](data).
			WithConcurrency(3).
			Execute(func(idx int, val int) (int, error) {
				current := atomic.AddInt32(&currentRunning, 1)

				mu.Lock()
				if current > maxRunning {
					maxRunning = current
				}
				mu.Unlock()

				time.Sleep(20 * time.Millisecond)
				atomic.AddInt32(&currentRunning, -1)
				return val * 2, nil
			})

		if maxRunning > 3 {
			t.Errorf("期望最大并发数不超过 3，实际为 %d", maxRunning)
		}
		if maxRunning < 2 {
			t.Errorf("期望最大并发数至少为 2，实际为 %d", maxRunning)
		}
	})

	t.Run("不限制并发数", func(t *testing.T) {
		data := make([]int, 20)
		for i := range data {
			data[i] = i + 1
		}

		var (
			mu             sync.Mutex
			currentRunning int32
			maxRunning     int32
		)

		NewParallelSliceExecutor[int, int](data).
			Execute(func(idx int, val int) (int, error) {
				current := atomic.AddInt32(&currentRunning, 1)

				mu.Lock()
				if current > maxRunning {
					maxRunning = current
				}
				mu.Unlock()

				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&currentRunning, -1)
				return val * 2, nil
			})

		// 不限制时，并发数应该接近数据量
		if maxRunning < 10 {
			t.Errorf("期望最大并发数至少为 10，实际为 %d", maxRunning)
		}
	})

	t.Run("并发控制下结果正确性", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		var results []int
		var mu sync.Mutex

		NewParallelSliceExecutor[int, int](data).
			WithConcurrency(3).
			OnSuccess(func(idx int, val int, result int) {
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}).
			Execute(func(idx int, val int) (int, error) {
				time.Sleep(5 * time.Millisecond)
				return val * 2, nil
			})

		if len(results) != 10 {
			t.Errorf("期望 10 个结果，实际 %d 个", len(results))
		}

		// 验证所有结果都存在
		resultMap := make(map[int]bool)
		for _, r := range results {
			resultMap[r] = true
		}
		for _, v := range data {
			expected := v * 2
			if !resultMap[expected] {
				t.Errorf("缺少结果 %d", expected)
			}
		}
	})

	t.Run("并发控制下错误处理", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		var (
			successCount int32
			errorCount   int32
		)

		NewParallelSliceExecutor[int, int](data).
			WithConcurrency(2).
			OnSuccess(func(idx int, val int, result int) {
				atomic.AddInt32(&successCount, 1)
			}).
			OnError(func(idx int, val int, err error) {
				atomic.AddInt32(&errorCount, 1)
			}).
			Execute(func(idx int, val int) (int, error) {
				time.Sleep(10 * time.Millisecond)
				if val%2 == 0 {
					return 0, fmt.Errorf("偶数错误: %d", val)
				}
				return val * 2, nil
			})

		if successCount != 3 {
			t.Errorf("期望成功 3 次，实际 %d 次", successCount)
		}
		if errorCount != 2 {
			t.Errorf("期望错误 2 次，实际 %d 次", errorCount)
		}
	})

	t.Run("链式调用WithConcurrency", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		var completed bool

		NewParallelSliceExecutor[int, int](data).
			WithConcurrency(2).
			OnSuccess(func(idx int, val int, result int) {
				// 成功回调
			}).
			OnError(func(idx int, val int, err error) {
				// 错误回调
			}).
			OnComplete(func(results []int, errors []error) {
				completed = true
			}).
			Execute(func(idx int, val int) (int, error) {
				return val * 2, nil
			})

		if !completed {
			t.Error("OnComplete 未执行")
		}
	})
}

// BenchmarkParallelSliceExecutor_Concurrency 测试不同并发数的性能
func BenchmarkParallelSliceExecutor_Concurrency(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = i + 1
	}

	b.Run("无限制", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			NewParallelSliceExecutor[int, int](data).
				Execute(func(idx int, val int) (int, error) {
					time.Sleep(time.Microsecond)
					return val * 2, nil
				})
		}
	})

	b.Run("并发数10", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			NewParallelSliceExecutor[int, int](data).
				WithConcurrency(10).
				Execute(func(idx int, val int) (int, error) {
					time.Sleep(time.Microsecond)
					return val * 2, nil
				})
		}
	})

	b.Run("并发数50", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			NewParallelSliceExecutor[int, int](data).
				WithConcurrency(50).
				Execute(func(idx int, val int) (int, error) {
					time.Sleep(time.Microsecond)
					return val * 2, nil
				})
		}
	})
}

// 示例：使用并发控制处理大量任务
func ExampleParallelSliceExecutor_WithConcurrency() {
	// 模拟100个任务，但只允许10个并发
	tasks := make([]int, 100)
	for i := range tasks {
		tasks[i] = i + 1
	}

	var completed int32

	NewParallelSliceExecutor[int, string](tasks).
		WithConcurrency(10). // 限制最多10个并发
		OnSuccess(func(idx int, val int, result string) {
			count := atomic.AddInt32(&completed, 1)
			if count%10 == 0 {
				fmt.Printf("已完成 %d 个任务\n", count)
			}
		}).
		OnComplete(func(results []string, errors []error) {
			fmt.Printf("全部完成！总计 %d 个任务\n", len(results))
		}).
		Execute(func(idx int, val int) (string, error) {
			// 模拟耗时操作
			time.Sleep(time.Millisecond)
			return fmt.Sprintf("任务 %d 完成", val), nil
		})

	// Output:
	// 已完成 10 个任务
	// 已完成 20 个任务
	// 已完成 30 个任务
	// 已完成 40 个任务
	// 已完成 50 个任务
	// 已完成 60 个任务
	// 已完成 70 个任务
	// 已完成 80 个任务
	// 已完成 90 个任务
	// 已完成 100 个任务
	// 全部完成！总计 100 个任务
}

// TestParallelSliceExecutor_OnPanic_BasicPanicHandling 测试基本的 panic 捕获
func TestParallelSliceExecutor_OnPanic_BasicPanicHandling(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	panicIndex := 2 // 第3个元素会 panic

	var (
		mu           sync.Mutex
		panicCalls   []int
		panicValues  []any
		successCalls []int
	)

	NewParallelSliceExecutor[int, string](data).
		OnSuccess(func(idx int, val int, result string) {
			mu.Lock()
			defer mu.Unlock()
			successCalls = append(successCalls, idx)
		}).
		OnPanic(func(idx int, val int, panicVal any) {
			mu.Lock()
			defer mu.Unlock()
			panicCalls = append(panicCalls, idx)
			panicValues = append(panicValues, panicVal)
		}).
		Execute(func(idx int, val int) (string, error) {
			if idx == panicIndex {
				panic("intentional panic")
			}
			return "success", nil
		})

	// 验证 panic 被捕获
	assert.Len(t, panicCalls, 1, "期望捕获 1 次 panic")
	assert.Equal(t, panicIndex, panicCalls[0], "panic 索引应该匹配")
	assert.Len(t, panicValues, 1, "期望有 1 个 panic 值")
	assert.Equal(t, "intentional panic", panicValues[0], "panic 值应该匹配")

	// 验证其他任务正常执行
	assert.Len(t, successCalls, len(data)-1, "其他任务应该正常执行")
}

// TestParallelSliceExecutor_OnPanic_MultiplePanics 测试多个任务同时 panic
func TestParallelSliceExecutor_OnPanic_MultiplePanics(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	panicIndices := map[int]bool{2: true, 5: true, 8: true}

	var (
		mu           sync.Mutex
		panicCount   int
		successCount int
	)

	NewParallelSliceExecutor[int, string](data).
		OnSuccess(func(idx int, val int, result string) {
			mu.Lock()
			defer mu.Unlock()
			successCount++
		}).
		OnPanic(func(idx int, val int, panicVal any) {
			mu.Lock()
			defer mu.Unlock()
			panicCount++
			assert.True(t, panicIndices[idx], "panic 索引应该在预期范围内: %d", idx)
		}).
		Execute(func(idx int, val int) (string, error) {
			if panicIndices[idx] {
				panic("intentional panic")
			}
			return "success", nil
		})

	// 验证所有 panic 都被捕获
	assert.Equal(t, len(panicIndices), panicCount, "应该捕获所有 panic")

	// 验证其他任务正常执行
	expectedSuccess := len(data) - len(panicIndices)
	assert.Equal(t, expectedSuccess, successCount, "其他任务应该正常执行")
}

// TestParallelSliceExecutor_OnPanic_WithError 测试 panic 和 error 同时存在
func TestParallelSliceExecutor_OnPanic_WithError(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	var (
		mu           sync.Mutex
		panicCount   int
		errorCount   int
		successCount int
	)

	NewParallelSliceExecutor[int, string](data).
		OnSuccess(func(idx int, val int, result string) {
			mu.Lock()
			defer mu.Unlock()
			successCount++
		}).
		OnError(func(idx int, val int, err error) {
			mu.Lock()
			defer mu.Unlock()
			errorCount++
		}).
		OnPanic(func(idx int, val int, panicVal any) {
			mu.Lock()
			defer mu.Unlock()
			panicCount++
		}).
		Execute(func(idx int, val int) (string, error) {
			switch idx {
			case 1:
				panic("panic at index 1")
			case 3:
				return "", &testError{msg: "error at index 3"}
			default:
				return "success", nil
			}
		})

	// 验证统计
	assert.Equal(t, 1, panicCount, "应该有 1 次 panic")
	assert.Equal(t, 1, errorCount, "应该有 1 次 error")
	assert.Equal(t, 3, successCount, "应该有 3 次成功")
}

// TestParallelSliceExecutor_OnPanic_NoPanicCallback 测试没有设置 OnPanic 时不会崩溃
func TestParallelSliceExecutor_OnPanic_NoPanicCallback(t *testing.T) {
	data := []int{1, 2, 3}

	var (
		mu           sync.Mutex
		successCount int
	)

	// 不设置 OnPanic，验证不会导致程序崩溃
	NewParallelSliceExecutor[int, string](data).
		OnSuccess(func(idx int, val int, result string) {
			mu.Lock()
			defer mu.Unlock()
			successCount++
		}).
		Execute(func(idx int, val int) (string, error) {
			if idx == 1 {
				panic("panic without callback")
			}
			return "success", nil
		})

	// 验证其他任务仍然执行
	assert.Equal(t, 2, successCount, "其他任务应该仍然执行")
}

// TestParallelSliceExecutor_OnPanic_WithConcurrency 测试带并发限制的 panic 处理
func TestParallelSliceExecutor_OnPanic_WithConcurrency(t *testing.T) {
	data := make([]int, 20)
	for i := range data {
		data[i] = i
	}

	var (
		mu           sync.Mutex
		panicCount   int
		successCount int
	)

	NewParallelSliceExecutor[int, string](data).
		WithConcurrency(5). // 限制并发数为 5
		OnSuccess(func(idx int, val int, result string) {
			mu.Lock()
			defer mu.Unlock()
			successCount++
		}).
		OnPanic(func(idx int, val int, panicVal any) {
			mu.Lock()
			defer mu.Unlock()
			panicCount++
		}).
		Execute(func(idx int, val int) (string, error) {
			// 每隔 5 个元素 panic 一次
			if idx%5 == 0 && idx > 0 {
				panic("periodic panic")
			}
			return "success", nil
		})

	// 验证统计（索引 5, 10, 15 会 panic）
	expectedPanics := 3
	assert.Equal(t, expectedPanics, panicCount, "应该捕获预期数量的 panic")

	expectedSuccess := len(data) - expectedPanics
	assert.Equal(t, expectedSuccess, successCount, "其他任务应该正常执行")
}

// TestParallelSliceExecutor_OnPanic_PanicValue 测试不同类型的 panic 值
func TestParallelSliceExecutor_OnPanic_PanicValue(t *testing.T) {
	data := []int{1, 2, 3, 4}

	var (
		mu          sync.Mutex
		panicValues []any
	)

	NewParallelSliceExecutor[int, string](data).
		OnPanic(func(idx int, val int, panicVal any) {
			mu.Lock()
			defer mu.Unlock()
			panicValues = append(panicValues, panicVal)
		}).
		Execute(func(idx int, val int) (string, error) {
			switch idx {
			case 0:
				panic("string panic")
			case 1:
				panic(42)
			case 2:
				panic(&testError{msg: "error panic"})
			default:
				return "success", nil
			}
		})

	// 验证捕获了 3 种不同类型的 panic
	assert.Len(t, panicValues, 3, "应该捕获 3 次 panic")

	// 验证 panic 值类型
	hasString := false
	hasInt := false
	hasError := false

	for _, pv := range panicValues {
		switch pv.(type) {
		case string:
			hasString = true
		case int:
			hasInt = true
		case *testError:
			hasError = true
		}
	}

	assert.True(t, hasString, "应该捕获 string 类型的 panic")
	assert.True(t, hasInt, "应该捕获 int 类型的 panic")
	assert.True(t, hasError, "应该捕获 error 类型的 panic")
}

// testError 测试用的错误类型
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
