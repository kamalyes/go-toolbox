/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-02 21:33:02
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-02 21:33:02
 * @FilePath: \go-toolbox\pkg\osx\id_bench_test.go
 * @Description: 唯一 ID 生成的性能与格式回归 | Performance and format regression for unique id generation
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package osx

import (
	"regexp"
	"sync"
	"testing"
)

var hexID32 = regexp.MustCompile(`^[0-9a-f]{32}$`)

// 格式是契约的一部分：调用方把它当 OpenTelemetry trace ID 用，宽度必须是 32 位小写 hex
// [EN] The format is part of the contract: callers use it as an OpenTelemetry trace id, so it
// has to stay thirty-two lowercase hex characters
func TestHashUnixMicroCipherTextFormat(t *testing.T) {
	for i := 0; i < 1000; i++ {
		got := HashUnixMicroCipherText()
		if !hexID32.MatchString(got) {
			t.Fatalf("第 %d 次拿到 %q，不是 32 位小写 hex", i, got)
		}
	}
}

// 唯一性：同一进程内不能重复，并发下也不能
// [EN] Uniqueness within one process, including under concurrency
func TestHashUnixMicroCipherTextUnique(t *testing.T) {
	const perG, goroutines = 2000, 8

	var mu sync.Mutex
	seen := make(map[string]struct{}, perG*goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			local := make([]string, 0, perG)
			for i := 0; i < perG; i++ {
				local = append(local, HashUnixMicroCipherText())
			}
			mu.Lock()
			defer mu.Unlock()
			for _, s := range local {
				if _, dup := seen[s]; dup {
					t.Errorf("重复的 ID：%s", s)
					return
				}
				seen[s] = struct{}{}
			}
		}()
	}
	wg.Wait()

	if len(seen) != perG*goroutines {
		t.Errorf("期望 %d 个不重复的 ID，实际 %d 个", perG*goroutines, len(seen))
	}
}

func BenchmarkHashUnixMicroCipherText(b *testing.B) {
	b.ReportAllocs()
	var sink string
	for i := 0; i < b.N; i++ {
		sink = HashUnixMicroCipherText()
	}
	_ = sink
}

// 网关每个请求调两次（trace id 与 request id 各一次），所以这个口径更贴近真实成本
// [EN] The gateway calls it twice per request, once for the trace id and once for the request
// id, so this shape is closer to the real cost
func BenchmarkHashUnixMicroCipherText_TwicePerRequest(b *testing.B) {
	b.ReportAllocs()
	var a, c string
	for i := 0; i < b.N; i++ {
		a = HashUnixMicroCipherText()
		c = HashUnixMicroCipherText()
	}
	_, _ = a, c
}

func BenchmarkHashUnixMicroCipherText_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var sink string
		for pb.Next() {
			sink = HashUnixMicroCipherText()
		}
		_ = sink
	})
}
