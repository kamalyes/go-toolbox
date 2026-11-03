/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-03 18:31:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-03 19:06:00
 * @FilePath: \go-toolbox\pkg\desensitize\keyrules_gate_test.go
 * @Description: ActionForKey 首末字节 O(1) 否决的差分测试
 *
 * 为什么必须逐项差分而不是抽查几个用例
 *
 * 这个闸门的失败形态是静默的：误否决一个键名，它的脱敏动作就从"擦除"变成
 * "原样输出"，于是一个密码明文落进日志，而且不报任何错
 * 所以这里拿"去掉闸门的原实现"当预言机，对规则表全量、规则表的大小写与
 * 分隔符变体、以及随机串逐个比对两条路径的结果
 *
 * [EN] A differential test rather than a handful of spot cases
 * This gate fails silently: a wrongly rejected key has its action downgraded from wiping to
 * passing through, so a password lands in the log in plaintext and nothing is reported
 * The oracle here is the original implementation without the gate, compared case by case
 * across the whole rule table, its case and separator variants, and random strings
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package desensitize

import (
	"math/rand"
	"strings"
	"testing"
)

// actionForKeyNoGate 是加闸门之前的 ActionForKey，逐行照抄，作为差分预言机
// [EN] ActionForKey as it was before the gate, copied line for line as the oracle
func actionForKeyNoGate(key string) KeyAction {
	first, last, n, needsNorm := keyShape(key)
	if n == 0 {
		return KeepAsIs
	}

	if hasKeyShape(first, last, n) {
		if !needsNorm {
			if action, ok := exactKeyRules[key]; ok {
				return action
			}
		} else if n <= maxNormalizedKeyLen {
			var buf [maxNormalizedKeyLen]byte
			if m, ok := normalizeKeyInto(buf[:], key); ok {
				if action, hit := exactKeyRules[string(buf[:m])]; hit {
					return action
				}
			}
		}
	}

	return matchSuffix(key)
}

// titleCase 首字母大写，不用已废弃的 strings.Title
// [EN] Capitalizes the first letter without the deprecated strings.Title
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// gateCorpus 汇总所有要比对的键名
func gateCorpus() []string {
	var base []string
	for k := range exactKeyRules {
		base = append(base, k)
	}
	for _, r := range suffixKeyRules {
		base = append(base, r.suffix, "user"+r.suffix, "x_"+r.suffix)
	}
	base = append(base, visibleKeySuffixes...)
	base = append(base,
		"method", "path", "status", "bytes", "duration_ms", "ip", "user_agent",
		"trace_id", "request_id", "user_id", "tenant_id", "domain", "platform_id", "region_id")

	out := make([]string, 0, len(base)*8)
	for _, k := range base {
		out = append(out, k, strings.ToUpper(k), titleCase(k))
		// 首末带分隔符：归一化会去掉它们，闸门必须放行而不是否决
		out = append(out, "_"+k, k+"_", "-"+k, k+"-", "."+k, k+".", "_"+k+"_")
		// 中间插分隔符
		if len(k) > 2 {
			mid := len(k) / 2
			out = append(out, k[:mid]+"_"+k[mid:], k[:mid]+"-"+k[mid:], k[:mid]+"."+k[mid:])
		}
	}

	// 边界形态
	out = append(out, "", "_", "-", ".", "___", "_-.", " ", "  ")
	for c := byte('a'); c <= 'z'; c++ {
		out = append(out, string(c), strings.ToUpper(string(c)))
	}
	for c := byte('0'); c <= '9'; c++ {
		out = append(out, string(c), "k"+string(c), string(c)+"k")
	}
	// 非 ASCII：keyShape 原样保留这些字节，闸门也必须一致
	out = append(out, "密码", "口令token", "token口令", "héllo", "a\xffb", "\xff", "\x80token")

	return out
}

// TestActionForKeyGateMatchesNoGate 闸门不得改变任何键名的判定结果
// [EN] The gate must not change the verdict for any key
func TestActionForKeyGateMatchesNoGate(t *testing.T) {
	for _, k := range gateCorpus() {
		if got, want := ActionForKey(k), actionForKeyNoGate(k); got != want {
			t.Errorf("键名 %q 判定不一致：带闸门=%v 预言机=%v", k, got, want)
		}
	}
}

// TestActionForKeyGateMatchesNoGateRandom 随机串同样不得分歧
//
// 固定种子，失败可复现 字母表刻意混入分隔符、数字、大写与非 ASCII，
// 覆盖闸门里每个条件分支
//
// [EN] Random strings must agree too, with a fixed seed so a failure reproduces
// The alphabet deliberately mixes separators, digits, uppercase and non-ASCII to reach
// every branch of the gate
func TestActionForKeyGateMatchesNoGateRandom(t *testing.T) {
	alphabet := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-.\xff\x80")
	rng := rand.New(rand.NewSource(20261103))
	buf := make([]byte, 0, 24)
	for i := 0; i < 400000; i++ {
		buf = buf[:0]
		n := rng.Intn(17)
		for j := 0; j < n; j++ {
			buf = append(buf, alphabet[rng.Intn(len(alphabet))])
		}
		k := string(buf)
		if got, want := ActionForKey(k), actionForKeyNoGate(k); got != want {
			t.Fatalf("随机键名 %q 判定不一致：带闸门=%v 预言机=%v", k, got, want)
		}
	}
}

// TestActionForKeyGateAfterRuleMutation 运行时增删规则后闸门仍要跟上
//
// 闸门与 hasKeyShape 共用位图正是为了这一点：若另建一张表而忘了同步，
// 新注册的字段名会被闸门挡在门外，注册了却不生效
// 反向也要对——UnregisterSensitiveKeys 刻意不清位图，残留的位只会多一次
// 本不必要的查表，不会给出错误答案
//
// [EN] The gate must follow runtime rule changes, which is why it shares hasKeyShape's
// bitmap: a separate table left unsynchronized would have the gate reject a freshly added
// rule, so adding it would have no effect
func TestActionForKeyGateAfterRuleMutation(t *testing.T) {
	const key = "zzqqvv"
	if got := ActionForKey(key); got != KeepAsIs {
		t.Fatalf("前置条件不成立：%q 本应无规则，实得 %v", key, got)
	}

	RegisterSensitiveKeys(key)
	if got, want := ActionForKey(key), actionForKeyNoGate(key); got != want {
		t.Errorf("新增规则后判定不一致：带闸门=%v 预言机=%v", got, want)
	}
	if got := ActionForKey(key); got != WipeValue {
		t.Errorf("新增的规则被闸门挡住了：实得 %v", got)
	}

	UnregisterSensitiveKeys(key)
	if got, want := ActionForKey(key), actionForKeyNoGate(key); got != want {
		t.Errorf("注销后判定不一致：带闸门=%v 预言机=%v", got, want)
	}
	if got := ActionForKey(key); got != KeepAsIs {
		t.Errorf("注销后仍被判为敏感：实得 %v", got)
	}
}
