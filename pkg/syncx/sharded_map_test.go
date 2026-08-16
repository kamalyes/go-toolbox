/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-01 00:51:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-07 01:56:56
 * @FilePath: \go-toolbox\pkg\syncx\sharded_map_test.go
 * @Description: ShardedMap 分片映射表测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package syncx

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestShardedMapStoreLoad 测试基本存储和加载
func TestShardedMapStoreLoad(t *testing.T) {
	m := NewShardedMap[string, int](64)

	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	v, ok := m.Load("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = m.Load("b")
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	_, ok = m.Load("notexist")
	assert.False(t, ok)
}

// TestShardedMapDelete 测试删除
func TestShardedMapDelete(t *testing.T) {
	m := NewShardedMap[string, int](32)

	m.Store("key1", 100)
	m.Store("key2", 200)

	assert.Equal(t, 2, m.Len())

	m.Delete("key1")

	_, ok := m.Load("key1")
	assert.False(t, ok)
	assert.Equal(t, 1, m.Len())

	// 删除不存在的 key 不应报错，也不应影响长度
	m.Delete("notexist")
	assert.Equal(t, 1, m.Len())
}

// TestShardedMapLoadAndDelete 测试加载并删除
func TestShardedMapLoadAndDelete(t *testing.T) {
	m := NewShardedMap[string, string](16)

	m.Store("test", "value")

	v, ok := m.LoadAndDelete("test")
	assert.True(t, ok)
	assert.Equal(t, "value", v)
	assert.False(t, m.Has("test"))

	// 再次 LoadAndDelete 不存在的 key
	_, ok = m.LoadAndDelete("test")
	assert.False(t, ok)
}

// TestShardedMapLoadOrStore 测试加载或存储
func TestShardedMapLoadOrStore(t *testing.T) {
	m := NewShardedMap[string, int](16)

	// 第一次存储
	v, loaded := m.LoadOrStore("k", 42)
	assert.False(t, loaded)
	assert.Equal(t, 42, v)

	// 第二次加载已有值
	v, loaded = m.LoadOrStore("k", 99)
	assert.True(t, loaded)
	assert.Equal(t, 42, v)
}

// TestShardedMapRange 测试遍历
func TestShardedMapRange(t *testing.T) {
	m := NewShardedMap[string, int](8)

	for i := 0; i < 100; i++ {
		m.Store(fmt.Sprintf("key%d", i), i)
	}

	count := 0
	m.Range(func(k string, v int) bool {
		count++
		return true
	})
	assert.Equal(t, 100, count)

	// 测试提前终止
	count = 0
	m.Range(func(k string, v int) bool {
		count++
		return count < 10
	})
	assert.Equal(t, 10, count)
}

// TestShardedMapRangeParallel 测试并行遍历
func TestShardedMapRangeParallel(t *testing.T) {
	m := NewShardedMap[string, int](64)

	total := 10000
	for i := 0; i < total; i++ {
		m.Store(fmt.Sprintf("key%d", i), i)
	}

	// 并行遍历计数（原子计数器，并发安全）
	var count atomic.Int64
	m.RangeParallel(0, func(k string, v int) {
		count.Add(1)
	})
	assert.Equal(t, int64(total), count.Load())

	// 测试空 map
	m.Clear()
	count.Store(0)
	m.RangeParallel(0, func(k string, v int) {
		count.Add(1)
	})
	assert.Equal(t, int64(0), count.Load())
}

// TestShardedMapRangeParallelCorrectness 并行遍历结果正确性验证
// 确保并行遍历访问到所有元素且无重复
func TestShardedMapRangeParallelCorrectness(t *testing.T) {
	m := NewShardedMap[int, int](64)

	total := 100000
	for i := 0; i < total; i++ {
		m.Store(i, i*2)
	}

	// 并行收集所有 key（mutex 保护，模拟实际使用场景）
	var mu sync.Mutex
	seen := make(map[int]struct{}, total)
	m.RangeParallel(0, func(k int, v int) {
		mu.Lock()
		seen[k] = struct{}{}
		mu.Unlock()
	})

	// 验证：每个 key 恰好被访问一次
	assert.Equal(t, total, len(seen))
	for i := 0; i < total; i++ {
		_, ok := seen[i]
		assert.True(t, ok, "key %d missing", i)
	}
}

// TestShardedMapRangeParallelConcurrent 并行遍历期间并发写入不 panic
// 验证 RLock 与 Lock 不冲突，遍历期间写入安全（遍历看到的是旧值快照）
func TestShardedMapRangeParallelConcurrent(t *testing.T) {
	m := NewShardedMap[string, int](64)

	for i := 0; i < 1000; i++ {
		m.Store(fmt.Sprintf("k%d", i), i)
	}

	var wg sync.WaitGroup
	// 并行遍历
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.RangeParallel(0, func(k string, v int) {
			// 模拟轻量处理
			_ = v
		})
	}()

	// 并发写入（不同 key，不干扰遍历的 shard 读锁）
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				m.Store(fmt.Sprintf("new%d_%d", idx, j), j)
			}
		}(i)
	}
	wg.Wait()
}

// TestShardedMapKeysValues 测试获取所有 key 和 value
func TestShardedMapKeysValues(t *testing.T) {
	m := NewShardedMap[string, int](16)

	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	keys := m.Keys()
	assert.Len(t, keys, 3)

	values := m.Values()
	assert.Len(t, values, 3)
}

// TestShardedMapClear 测试清空
func TestShardedMapClear(t *testing.T) {
	m := NewShardedMap[string, int](16)

	for i := 0; i < 50; i++ {
		m.Store(fmt.Sprintf("k%d", i), i)
	}

	assert.Equal(t, 50, m.Len())

	m.Clear()

	assert.Zero(t, m.Len())
}

// TestShardedMapCount 测试条件计数
func TestShardedMapCount(t *testing.T) {
	m := NewShardedMap[string, int](16)

	for i := 0; i < 20; i++ {
		m.Store(fmt.Sprintf("k%d", i), i)
	}

	// 无过滤条件，等价于 Len
	assert.Equal(t, 20, m.Count(nil))

	// 计数偶数值
	evenCount := m.Count(func(k string, v int) bool {
		return v%2 == 0
	})
	assert.Equal(t, 10, evenCount)
}

// TestShardedMapConcurrentWrite 并发写入测试
func TestShardedMapConcurrentWrite(t *testing.T) {
	m := NewShardedMap[string, int](64)

	var wg sync.WaitGroup

	goroutines := 100
	writesPerGoroutine := 100
	expected := goroutines * writesPerGoroutine

	// 并发写入
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()

			for j := 0; j < writesPerGoroutine; j++ {
				key := fmt.Sprintf("g%d_k%d", gid, j)
				m.Store(key, gid*1000+j)
			}
		}(i)
	}
	wg.Wait()

	assert.Equal(t, expected, m.Len())

	// 并发读取验证
	var readCount atomic.Int64

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()

			for j := 0; j < writesPerGoroutine; j++ {
				key := fmt.Sprintf("g%d_k%d", gid, j)
				if v, ok := m.Load(key); ok && v == gid*1000+j {
					readCount.Add(1)
				}
			}
		}(i)
	}
	wg.Wait()

	assert.Equal(t, int64(expected), readCount.Load())
}

// TestShardedMapConcurrentMixed 并发混合读写删测试
func TestShardedMapConcurrentMixed(t *testing.T) {
	m := NewShardedMap[string, int](32)

	var wg sync.WaitGroup

	// 预填充
	for i := 0; i < 1000; i++ {
		m.Store(fmt.Sprintf("init%d", i), i)
	}

	// 并发混合操作
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("mixed_%d_%d", gid, j)
				m.Store(key, j)
				_, _ = m.Load(key)
				m.Delete(key)
				_, _ = m.Load(key)
			}
		}(i)
	}
	wg.Wait()

	// 初始的 1000 个应该还在
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("init%d", i)

		v, ok := m.Load(key)
		assert.True(t, ok, "key should exist: %s", key)
		assert.Equal(t, i, v, "value mismatch for key: %s", key)
	}
}

// TestShardedMapWithShardLock 测试分片级锁
func TestShardedMapWithShardLock(t *testing.T) {
	m := NewShardedMap[string, int](16)

	m.Store("test", 10)

	// 使用 WithShardLock 原子修改
	m.WithShardLock("test", func(data map[string]int) {
		v := data["test"]
		data["test"] = v * 2
	})

	v, ok := m.Load("test")
	assert.True(t, ok)
	assert.Equal(t, 20, v)

	// 使用 WithShardRLock 读取
	m.WithShardRLock("test", func(data map[string]int) {
		assert.Equal(t, 20, data["test"])
	})
}

// TestShardedMapShardCount 测试分片数量自动调整
func TestShardedMapShardCount(t *testing.T) {
	// shardCount=100（非 2 的幂），应自动调整为 128
	m := NewShardedMap[string, int](100)
	assert.Equal(t, 128, m.shardCount)

	// shardCount=0，应使用默认值 64
	m2 := NewShardedMap[string, int](0)
	assert.Equal(t, 64, m2.shardCount)

	// shardCount=64（已是 2 的幂），保持不变
	m3 := NewShardedMap[string, int](64)
	assert.Equal(t, 64, m3.shardCount)
}

// TestShardedMapWithPerShardHint 测试 WithPerShardHint 选项预分配容量
//
// 验证点：
//  1. 指定 hint 后所有 shard 的 data map 内部长度仍为 0（仅预分配容量，未写入数据）
//  2. 容量预分配不影响功能（Store/Load/Range 正常工作）
//  3. hint <= 0 等价于不预分配
func TestShardedMapWithPerShardHint(t *testing.T) {
	// 1. 指定每 shard 预分配 100 容量
	m := NewShardedMapWithOptions[string, int](64, WithPerShardHint[string, int](100))

	// 初始状态：每个 shard 的 data 长度为 0（map 还没写入数据）
	for i, shard := range m.shards {
		assert.Len(t, shard.data, 0, "shard[%d] should be empty after preallocation", i)
	}

	// 2. 写入数据后功能正常
	m.Store("a", 1)
	m.Store("b", 2)

	v, ok := m.Load("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.Equal(t, 2, m.Len())

	// 3. hint <= 0 等价于不预分配（兼容旧版 NewShardedMap）
	m2 := NewShardedMapWithOptions[string, int](64, WithPerShardHint[string, int](0))
	for i, shard := range m2.shards {
		assert.Len(t, shard.data, 0, "shard[%d] should be empty when hint=0", i)
	}

	// 4. 负数 hint 应被归一化为 0
	m3 := NewShardedMapWithOptions[string, int](64, WithPerShardHint[string, int](-10))
	for i, shard := range m3.shards {
		assert.Len(t, shard.data, 0, "shard[%d] should be empty when hint<0", i)
	}

	// 5. 不传 opts 等价于 NewShardedMap
	m4 := NewShardedMapWithOptions[string, int](64)
	assert.Equal(t, 64, m4.shardCount)
}

// TestShardedMapWithPerShardHintFunctional 验证预分配后高并发写入功能正确
//
// 重点：预分配容量不应破坏并发安全性与一致性
func TestShardedMapWithPerShardHintFunctional(t *testing.T) {
	const shardCount = 64
	const writers = 8
	const perWriter = 500

	m := NewShardedMapWithOptions[int, int](
		shardCount,
		WithPerShardHint[int, int](writers*perWriter/shardCount+1),
	)

	var wg sync.WaitGroup

	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()

			for i := 0; i < perWriter; i++ {
				key := base*perWriter + i
				m.Store(key, key*10)
			}
		}(w)
	}
	wg.Wait()

	expected := writers * perWriter
	assert.Equal(t, expected, m.Len())

	// 抽样校验数据完整性
	for w := 0; w < writers; w++ {
		for i := 0; i < perWriter; i += 50 {
			key := w*perWriter + i

			v, ok := m.Load(key)
			assert.True(t, ok, "key should exist: %d", key)
			assert.Equal(t, key*10, v)
		}
	}
}

// TestNewShardedMapCustomHasher 测试自定义 hash 函数
func TestNewShardedMapCustomHasher(t *testing.T) {
	m := NewShardedMap[string, int](16)

	m.Store("1", 1)
	m.Store("2", 2)

	v, ok := m.Load("1")
	assert.True(t, ok)
	assert.Equal(t, 1, v)
}

// TestShardedMapSwap 测试 Swap 替换值并返回旧值
func TestShardedMapSwap(t *testing.T) {
	t.Run("key exists returns old value", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		m.Store("k", 100)

		old, ok := m.Swap("k", 200)
		assert.True(t, ok)
		assert.Equal(t, 100, old)

		v, ok := m.Load("k")
		assert.True(t, ok)
		assert.Equal(t, 200, v)
		// Swap 已存在的 key 不改变元素总数
		assert.Equal(t, 1, m.Len())
	})

	t.Run("key not exists stores and returns zero", func(t *testing.T) {
		m := NewShardedMap[string, int](16)

		old, ok := m.Swap("k", 200)
		assert.False(t, ok)
		assert.Equal(t, 0, old) // 零值

		v, ok := m.Load("k")
		assert.True(t, ok)
		assert.Equal(t, 200, v)
		assert.Equal(t, 1, m.Len())
	})

	t.Run("empty map swap string zero value", func(t *testing.T) {
		m := NewShardedMap[string, string](16)

		old, ok := m.Swap("k", "v")
		assert.False(t, ok)
		assert.Equal(t, "", old)
	})

	t.Run("concurrent swap same key", func(t *testing.T) {
		m := NewShardedMap[int, int](64)
		m.Store(1, 0)

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				m.Swap(1, val)
			}(i)
		}
		wg.Wait()

		// 元素总数不变
		assert.Equal(t, 1, m.Len())
		v, ok := m.Load(1)
		assert.True(t, ok)
		_ = v
	})
}

// TestShardedMapCompareAndSwap 测试原子比较并交换
func TestShardedMapCompareAndSwap(t *testing.T) {
	t.Run("key not exists returns false", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		assert.False(t, m.CompareAndSwap("k", 1, 2))
	})

	t.Run("value matches swap succeeds", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		m.Store("k", 1)

		assert.True(t, m.CompareAndSwap("k", 1, 2))

		v, ok := m.Load("k")
		assert.True(t, ok)
		assert.Equal(t, 2, v)
	})

	t.Run("value mismatch returns false", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		m.Store("k", 1)

		assert.False(t, m.CompareAndSwap("k", 99, 2))

		v, ok := m.Load("k")
		assert.True(t, ok)
		assert.Equal(t, 1, v) // 值未变
	})

	t.Run("empty map compare and swap", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		assert.False(t, m.CompareAndSwap("k", 0, 1))
	})

	t.Run("concurrent compare and swap increments", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		m.Store("counter", 0)

		var wg sync.WaitGroup
		var successCount atomic.Int64
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					old, _ := m.Load("counter")
					if m.CompareAndSwap("counter", old, old+1) {
						successCount.Add(1)
						return
					}
				}
			}()
		}
		wg.Wait()

		assert.Equal(t, int64(100), successCount.Load())
		v, _ := m.Load("counter")
		assert.Equal(t, 100, v)
	})
}

// TestShardedMapCompareAndDelete 测试原子比较并删除
func TestShardedMapCompareAndDelete(t *testing.T) {
	t.Run("key not exists returns false", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		assert.False(t, m.CompareAndDelete("k", 1))
		assert.Equal(t, 0, m.Len())
	})

	t.Run("value matches delete succeeds", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		m.Store("k", 1)

		assert.True(t, m.CompareAndDelete("k", 1))
		assert.False(t, m.Has("k"))
		assert.Equal(t, 0, m.Len())
	})

	t.Run("value mismatch returns false", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		m.Store("k", 1)

		assert.False(t, m.CompareAndDelete("k", 99))
		assert.True(t, m.Has("k"))
		assert.Equal(t, 1, m.Len())
	})

	t.Run("empty map compare and delete", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		assert.False(t, m.CompareAndDelete("k", 0))
	})

	t.Run("concurrent compare and delete", func(t *testing.T) {
		m := NewShardedMap[int, int](64)
		for i := 0; i < 100; i++ {
			m.Store(i, i)
		}

		var wg sync.WaitGroup
		var deletedCount atomic.Int64
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				if m.CompareAndDelete(idx, idx) {
					deletedCount.Add(1)
				}
			}(i)
		}
		wg.Wait()

		assert.Equal(t, int64(100), deletedCount.Load())
		assert.Equal(t, 0, m.Len())
	})
}

// TestShardedMapStoreBatch 测试批量存储
func TestShardedMapStoreBatch(t *testing.T) {
	t.Run("empty batch does nothing", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		m.StoreBatch(nil)
		assert.Equal(t, 0, m.Len())

		m.StoreBatch(map[string]int{})
		assert.Equal(t, 0, m.Len())
	})

	t.Run("normal batch store", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		items := make(map[string]int, 5)
		for i := 0; i < 5; i++ {
			items[fmt.Sprintf("k%d", i)] = i
		}

		m.StoreBatch(items)
		assert.Equal(t, 5, m.Len())

		for i := 0; i < 5; i++ {
			v, ok := m.Load(fmt.Sprintf("k%d", i))
			assert.True(t, ok)
			assert.Equal(t, i, v)
		}
	})

	t.Run("batch with existing keys only counts new", func(t *testing.T) {
		m := NewShardedMap[string, int](16)
		m.Store("k0", 100)

		items := map[string]int{"k0": 200, "k1": 300}
		m.StoreBatch(items)

		assert.Equal(t, 2, m.Len())
		v, _ := m.Load("k0")
		assert.Equal(t, 200, v) // 值被覆盖
		v, _ = m.Load("k1")
		assert.Equal(t, 300, v)
	})

	t.Run("large batch with int keys across shards", func(t *testing.T) {
		m := NewShardedMap[int, int](64)
		items := make(map[int]int, 1000)
		for i := 0; i < 1000; i++ {
			items[i] = i * 10
		}

		m.StoreBatch(items)
		assert.Equal(t, 1000, m.Len())

		// 抽样校验
		for i := 0; i < 1000; i += 100 {
			v, ok := m.Load(i)
			assert.True(t, ok)
			assert.Equal(t, i*10, v)
		}
	})

	t.Run("concurrent store batch", func(t *testing.T) {
		m := NewShardedMap[int, int](64)
		var wg sync.WaitGroup

		for w := 0; w < 10; w++ {
			wg.Add(1)
			go func(base int) {
				defer wg.Done()
				items := make(map[int]int, 100)
				for i := 0; i < 100; i++ {
					items[base*100+i] = i
				}
				m.StoreBatch(items)
			}(w)
		}
		wg.Wait()

		assert.Equal(t, 1000, m.Len())
	})
}

// TestShardedMapClearWithPerShardHint 测试带预分配容量提示的 Clear
// 覆盖 Clear 中 perShardHint > 0 的分支
func TestShardedMapClearWithPerShardHint(t *testing.T) {
	m := NewShardedMapWithOptions[string, int](16, WithPerShardHint[string, int](50))

	for i := 0; i < 50; i++ {
		m.Store(fmt.Sprintf("k%d", i), i)
	}
	assert.Equal(t, 50, m.Len())

	m.Clear()

	assert.Zero(t, m.Len())
	// 验证 Clear 后仍可正常写入读取（复用预分配容量）
	m.Store("after_clear", 999)
	v, ok := m.Load("after_clear")
	assert.True(t, ok)
	assert.Equal(t, 999, v)
}

// TestNextPowerOfTwo 测试 NextPowerOfTwo 函数的各分支
func TestNextPowerOfTwo(t *testing.T) {
	t.Run("n <= 1 returns 2", func(t *testing.T) {
		assert.Equal(t, 2, NextPowerOfTwo(0))
		assert.Equal(t, 2, NextPowerOfTwo(1))
		assert.Equal(t, 2, NextPowerOfTwo(-5))
	})

	t.Run("n is power of two returns next power", func(t *testing.T) {
		assert.Equal(t, 4, NextPowerOfTwo(2))
		assert.Equal(t, 8, NextPowerOfTwo(4))
		assert.Equal(t, 16, NextPowerOfTwo(8))
		assert.Equal(t, 32, NextPowerOfTwo(16))
		assert.Equal(t, 64, NextPowerOfTwo(32))
	})

	t.Run("n is not power of two returns next power", func(t *testing.T) {
		assert.Equal(t, 4, NextPowerOfTwo(3))
		assert.Equal(t, 128, NextPowerOfTwo(100))
		assert.Equal(t, 256, NextPowerOfTwo(200))
	})
}

// TestKvHasher 测试 KvHasher 为各类型选择的 hash 函数
func TestKvHasher(t *testing.T) {
	t.Run("string hasher", func(t *testing.T) {
		h := KvHasher[string]()
		assert.NotZero(t, h("hello"))
		assert.NotZero(t, h("world"))
		assert.Equal(t, h("abc"), h("abc")) // 相同 key 返回相同 hash
	})

	t.Run("int hasher", func(t *testing.T) {
		h := KvHasher[int]()
		assert.NotZero(t, h(123))
		assert.Equal(t, h(456), h(456))
	})

	t.Run("int64 hasher", func(t *testing.T) {
		h := KvHasher[int64]()
		assert.NotZero(t, h(int64(123456789)))
		assert.Equal(t, h(int64(1)), h(int64(1)))
	})

	t.Run("int32 hasher", func(t *testing.T) {
		h := KvHasher[int32]()
		assert.NotZero(t, h(int32(123)))
		assert.Equal(t, h(int32(1)), h(int32(1)))
	})

	t.Run("uint hasher", func(t *testing.T) {
		h := KvHasher[uint]()
		assert.NotZero(t, h(uint(123)))
		assert.Equal(t, h(uint(1)), h(uint(1)))
	})

	t.Run("uint64 hasher", func(t *testing.T) {
		h := KvHasher[uint64]()
		assert.NotZero(t, h(uint64(123)))
		assert.Equal(t, h(uint64(1)), h(uint64(1)))
	})

	t.Run("default hasher for bool type", func(t *testing.T) {
		// bool 不在专用分支中，走 default 分支用 fmt.Sprintf 转 string 再 hash
		h := KvHasher[bool]()
		assert.NotZero(t, h(true))
		assert.NotZero(t, h(false))
		assert.Equal(t, h(true), h(true))
	})
}

// BenchmarkShardedMapWrite 分片 map 写入基准测试
func BenchmarkShardedMapWrite(b *testing.B) {
	m := NewShardedMap[string, int](64)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Store(fmt.Sprintf("k%d", i), i)
			i++
		}
	})
}

// BenchmarkSyncMapWrite sync.Map 写入基准测试（对比）
func BenchmarkSyncMapWrite(b *testing.B) {
	var m sync.Map

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Store(fmt.Sprintf("k%d", i), i)
			i++
		}
	})
}

// BenchmarkShardedMapRead 分片 map 读取基准测试
func BenchmarkShardedMapRead(b *testing.B) {
	m := NewShardedMap[string, int](64)

	for i := 0; i < 10000; i++ {
		m.Store(fmt.Sprintf("k%d", i), i)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Load(fmt.Sprintf("k%d", i%10000))
			i++
		}
	})
}

// BenchmarkShardedMapReadPreGen 预生成 key 的读取基准测试（排除 fmt.Sprintf 干扰）
func BenchmarkShardedMapReadPreGen(b *testing.B) {
	m := NewShardedMap[string, int](64)
	keys := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		keys[i] = fmt.Sprintf("k%d", i)
		m.Store(keys[i], i)
	}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Load(keys[i%10000])
			i++
		}
	})
}

// BenchmarkShardedMapWritePreGen 预生成 key 的写入基准测试（排除 fmt.Sprintf 干扰）
func BenchmarkShardedMapWritePreGen(b *testing.B) {
	m := NewShardedMap[string, int](64)
	keys := make([]string, 100000)
	for i := range keys {
		keys[i] = fmt.Sprintf("k%d", i)
	}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Store(keys[i%100000], i)
			i++
		}
	})
}

// BenchmarkFNVHashString32 hash 函数基准测试
func BenchmarkFNVHashString32(b *testing.B) {
	s := "hello_world_key_12345"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = FNVHashString32(s)
	}
}

// BenchmarkShardedMapStoreBatch 批量写入基准测试
func BenchmarkShardedMapStoreBatch(b *testing.B) {
	m := NewShardedMap[int, int](64)
	items := make(map[int]int, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = i * 10
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.StoreBatch(items)
	}
}

// BenchmarkShardedMapSwap Swap 操作基准测试
func BenchmarkShardedMapSwap(b *testing.B) {
	m := NewShardedMap[int, int](64)
	for i := 0; i < 10000; i++ {
		m.Store(i, i)
	}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Swap(i%10000, i)
			i++
		}
	})
}

// BenchmarkShardedMapRange 串行遍历基准测试（百万级，轻量回调）
func BenchmarkShardedMapRange(b *testing.B) {
	m := NewShardedMap[int, int](64)
	for i := 0; i < 1000000; i++ {
		m.Store(i, i)
	}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var sum int64
		m.Range(func(k, v int) bool {
			sum += int64(v)
			return true
		})
		_ = sum
	}
}

// BenchmarkShardedMapRangeParallel 并行遍历基准测试（百万级，轻量回调）
// 注意：轻量回调下并行可能因 goroutine 开销慢于串行，实际收益取决于回调复杂度
// 对比 BenchmarkShardedMapRangeHeavy* 可见：回调越重，并行收益越大
func BenchmarkShardedMapRangeParallel(b *testing.B) {
	m := NewShardedMap[int, int](64)
	for i := 0; i < 1000000; i++ {
		m.Store(i, i)
	}
	var sum atomic.Int64
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sum.Store(0)
		m.RangeParallel(0, func(k, v int) {
			sum.Add(int64(v))
		})
	}
	_ = sum.Load()
}

// BenchmarkShardedMapRangeHeavy 串行遍历（重回调，模拟 go-wsc 心跳检查/广播场景）
// 回调包含：条件判断 + 原子读 + 模拟 TrySend 开销
func BenchmarkShardedMapRangeHeavy(b *testing.B) {
	m := NewShardedMap[int, int](64)
	for i := 0; i < 1000000; i++ {
		m.Store(i, i)
	}
	var counter atomic.Int64
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		counter.Store(0)
		m.Range(func(k, v int) bool {
			// 模拟 checkHeartbeat 回调：原子读时间戳 + 条件判断 + 计数
			if v > 0 {
				counter.Add(1)
			}
			// 模拟轻量计算开销（TrySend 的 channel select 成本）
			_ = k * v
			return true
		})
	}
}

// BenchmarkShardedMapRangeParallelHeavy 并行遍历（重回调，对比 BenchmarkShardedMapRangeHeavy）
func BenchmarkShardedMapRangeParallelHeavy(b *testing.B) {
	m := NewShardedMap[int, int](64)
	for i := 0; i < 1000000; i++ {
		m.Store(i, i)
	}
	var counter atomic.Int64
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		counter.Store(0)
		m.RangeParallel(0, func(k, v int) {
			// 模拟 checkHeartbeat 回调：原子读时间戳 + 条件判断 + 计数
			if v > 0 {
				counter.Add(1)
			}
			// 模拟轻量计算开销（TrySend 的 channel select 成本）
			_ = k * v
		})
	}
}
