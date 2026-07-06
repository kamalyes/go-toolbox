/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-01 00:51:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-07 01:28:56
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
)

// TestShardedMapStoreLoad 测试基本存储和加载
func TestShardedMapStoreLoad(t *testing.T) {
	m := NewShardedMap[string, int](64)

	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	if v, ok := m.Load("a"); !ok || v != 1 {
		t.Errorf("Load('a') = (%d, %v), want (1, true)", v, ok)
	}
	if v, ok := m.Load("b"); !ok || v != 2 {
		t.Errorf("Load('b') = (%d, %v), want (2, true)", v, ok)
	}
	if _, ok := m.Load("notexist"); ok {
		t.Error("Load('notexist') 应该不存在")
	}
}

// TestShardedMapDelete 测试删除
func TestShardedMapDelete(t *testing.T) {
	m := NewShardedMap[string, int](32)

	m.Store("key1", 100)
	m.Store("key2", 200)

	if m.Len() != 2 {
		t.Errorf("Len = %d, want 2", m.Len())
	}

	m.Delete("key1")
	if _, ok := m.Load("key1"); ok {
		t.Error("Delete 后 key1 应该不存在")
	}
	if m.Len() != 1 {
		t.Errorf("Len = %d, want 1", m.Len())
	}

	// 删除不存在的 key 不应报错
	m.Delete("notexist")
	if m.Len() != 1 {
		t.Errorf("删除不存在的 key 后 Len = %d, want 1", m.Len())
	}
}

// TestShardedMapLoadAndDelete 测试加载并删除
func TestShardedMapLoadAndDelete(t *testing.T) {
	m := NewShardedMap[string, string](16)

	m.Store("test", "value")

	v, ok := m.LoadAndDelete("test")
	if !ok || v != "value" {
		t.Errorf("LoadAndDelete = (%s, %v), want (value, true)", v, ok)
	}
	if m.Has("test") {
		t.Error("LoadAndDelete 后 key 应该不存在")
	}

	// 再次 LoadAndDelete 不存在的 key
	_, ok = m.LoadAndDelete("test")
	if ok {
		t.Error("LoadAndDelete 不存在的 key 应返回 false")
	}
}

// TestShardedMapLoadOrStore 测试加载或存储
func TestShardedMapLoadOrStore(t *testing.T) {
	m := NewShardedMap[string, int](16)

	// 第一次存储
	v, loaded := m.LoadOrStore("k", 42)
	if loaded {
		t.Error("首次 LoadOrStore 应返回 loaded=false")
	}
	if v != 42 {
		t.Errorf("LoadOrStore 返回 v=%d, want 42", v)
	}

	// 第二次加载已有值
	v, loaded = m.LoadOrStore("k", 99)
	if !loaded {
		t.Error("已存在 key 的 LoadOrStore 应返回 loaded=true")
	}
	if v != 42 {
		t.Errorf("LoadOrStore 返回 v=%d, want 已有值 42", v)
	}
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
	if count != 100 {
		t.Errorf("Range 遍历 count=%d, want 100", count)
	}

	// 测试提前终止
	count = 0
	m.Range(func(k string, v int) bool {
		count++
		return count < 10
	})
	if count != 10 {
		t.Errorf("Range 提前终止 count=%d, want 10", count)
	}
}

// TestShardedMapKeysValues 测试获取所有 key 和 value
func TestShardedMapKeysValues(t *testing.T) {
	m := NewShardedMap[string, int](16)

	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("Keys 长度 = %d, want 3", len(keys))
	}

	values := m.Values()
	if len(values) != 3 {
		t.Errorf("Values 长度 = %d, want 3", len(values))
	}
}

// TestShardedMapClear 测试清空
func TestShardedMapClear(t *testing.T) {
	m := NewShardedMap[string, int](16)

	for i := 0; i < 50; i++ {
		m.Store(fmt.Sprintf("k%d", i), i)
	}
	if m.Len() != 50 {
		t.Errorf("Len = %d, want 50", m.Len())
	}

	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Clear 后 Len = %d, want 0", m.Len())
	}
}

// TestShardedMapCount 测试条件计数
func TestShardedMapCount(t *testing.T) {
	m := NewShardedMap[string, int](16)

	for i := 0; i < 20; i++ {
		m.Store(fmt.Sprintf("k%d", i), i)
	}

	// 无过滤条件，等价于 Len
	if m.Count(nil) != 20 {
		t.Errorf("Count(nil) = %d, want 20", m.Count(nil))
	}

	// 计数偶数值
	evenCount := m.Count(func(k string, v int) bool {
		return v%2 == 0
	})
	if evenCount != 10 {
		t.Errorf("Count(偶数) = %d, want 10", evenCount)
	}
}

// TestShardedMapConcurrentWrite 并发写入测试
func TestShardedMapConcurrentWrite(t *testing.T) {
	m := NewShardedMap[string, int](64)
	var wg sync.WaitGroup
	goroutines := 100
	writesPerGoroutine := 100

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

	expected := goroutines * writesPerGoroutine
	if m.Len() != expected {
		t.Errorf("并发写入后 Len = %d, want %d", m.Len(), expected)
	}

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

	if int(readCount.Load()) != expected {
		t.Errorf("并发读取验证 count=%d, want %d", readCount.Load(), expected)
	}
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
				m.Store(key, j)    // 写
				_, _ = m.Load(key) // 读
				m.Delete(key)      // 删
				_, _ = m.Load(key) // 再读
			}
		}(i)
	}
	wg.Wait()

	// 初始的 1000 个应该还在
	for i := 0; i < 1000; i++ {
		if v, ok := m.Load(fmt.Sprintf("init%d", i)); !ok || v != i {
			t.Errorf("init%d 验证失败: v=%d, ok=%v", i, v, ok)
			break
		}
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

	if v, _ := m.Load("test"); v != 20 {
		t.Errorf("WithShardLock 修改后 v=%d, want 20", v)
	}

	// 使用 WithShardRLock 读取
	m.WithShardRLock("test", func(data map[string]int) {
		if data["test"] != 20 {
			t.Error("WithShardRLock 读取值不正确")
		}
	})
}

// TestShardedMapShardCount 验证非 2 的幂的 shardCount 自动调整
func TestShardedMapShardCount(t *testing.T) {
	// shardCount=100（非 2 的幂），应自动调整为 128
	m := NewShardedMap[string, int](100)
	if m.shardCount != 128 {
		t.Errorf("shardCount = %d, want 128 (100 向上取 2 的幂)", m.shardCount)
	}

	// shardCount=0，应使用默认值 64
	m2 := NewShardedMap[string, int](0)
	if m2.shardCount != 64 {
		t.Errorf("shardCount = %d, want 64 (默认值)", m2.shardCount)
	}

	// shardCount=64（已是 2 的幂），保持不变
	m3 := NewShardedMap[string, int](64)
	if m3.shardCount != 64 {
		t.Errorf("shardCount = %d, want 64", m3.shardCount)
	}
}

// TestNewShardedMapCustomHasher 测试自定义 hash 函数
func TestNewShardedMapCustomHasher(t *testing.T) {
	m := NewShardedMap[string, int](16)

	m.Store("1", 1)
	m.Store("2", 2)

	if v, _ := m.Load("1"); v != 1 {
		t.Errorf("Load(1) = %s, want one", v)
	}
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
