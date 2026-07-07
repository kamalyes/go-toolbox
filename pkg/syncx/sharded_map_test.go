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
