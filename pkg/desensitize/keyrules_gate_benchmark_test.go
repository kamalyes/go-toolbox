/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-03 18:33:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-03 19:11:00
 * @FilePath: \go-toolbox\pkg\desensitize\keyrules_gate_benchmark_test.go
 * @Description: ActionForKey 首末字节否决的常驻基准 - 收益与代价两边都量
 *
 * 闸门不是纯收益：被否决的键名省掉整趟扫描，命中精确表的键名则白走一遍闸门
 * 所以两组键名都留一档，去掉闸门的实现（actionForKeyNoGate）当对照
 * 只量前一组就是在挑对自己有利的数字
 *
 * [EN] Permanent benchmarks for ActionForKey's first-last rejection, measuring both sides
 * The gate is not a pure win: a rejected key skips the whole scan while a key that does hit
 * the exact table walks the gate for nothing, so both key sets are kept with the gateless
 * implementation as the control
 * Measuring only the first set would be choosing the flattering number
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package desensitize

import "testing"

var gwKeys = []string{
	"method", "path", "status", "bytes", "duration_ms", "ip", "user_agent",
	"trace_id", "request_id", "user_id", "tenant_id", "domain", "platform_id", "region_id",
}

var aSink KeyAction

func BenchmarkGatewayKeys(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, k := range gwKeys {
			aSink = ActionForKey(k)
		}
	}
}

func BenchmarkGatewayKeysNoGate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, k := range gwKeys {
			aSink = actionForKeyNoGate(k)
		}
	}
}

func BenchmarkGatewayKeysParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, k := range gwKeys {
				aSink = ActionForKey(k)
			}
		}
	})
}

func BenchmarkGatewayKeysNoGateParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, k := range gwKeys {
				aSink = actionForKeyNoGate(k)
			}
		}
	})
}

// 命中精确表的键名：闸门放行，这一组是闸门的代价侧
// [EN] Keys that hit the exact table pass the gate; this set is the gate's cost side
func BenchmarkSensitiveKeys(b *testing.B) {
	keys := []string{"password", "token", "secret", "authorization", "phone", "email"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, k := range keys {
			aSink = ActionForKey(k)
		}
	}
}

func BenchmarkSensitiveKeysNoGate(b *testing.B) {
	keys := []string{"password", "token", "secret", "authorization", "phone", "email"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, k := range keys {
			aSink = actionForKeyNoGate(k)
		}
	}
}

// keyShape 的无分支实现与逐字节参照实现对照
// 常驻是因为它的内层循环对每个字段每条日志都要走一遍
//
// [EN] The branchless keyShape against the byte-wise reference, kept permanently because
// its inner loop runs for every field of every line
func BenchmarkKeyShape(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, k := range gwKeys {
			shapeSink, _, _, _ = keyShape(k)
		}
	}
}

func BenchmarkKeyShapeRef(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, k := range gwKeys {
			shapeSink, _, _, _ = keyShapeRef(k)
		}
	}
}

var shapeSink byte
