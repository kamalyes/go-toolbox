/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-03 19:05:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-03 19:13:00
 * @FilePath: \go-toolbox\pkg\desensitize\keyshape_test.go
 * @Description: keyShape 无分支实现与原逐字节实现的差分测试
 *
 * 为什么要单独给 keyShape 一个参照实现
 *
 * keyrules_gate_test.go 的预言机 actionForKeyNoGate 自己就调 keyShape，
 * keyShape 一改，预言机跟着一起错，那组差分就测不出 keyShape 的回归了
 * 所以这里把改之前的逐字节实现原样留一份当参照，四个返回值逐个比
 *
 * [EN] Why keyShape needs its own reference implementation
 * The oracle in keyrules_gate_test.go, actionForKeyNoGate, calls keyShape itself, so a
 * change to keyShape moves the oracle with it and that suite can no longer catch a
 * keyShape regression
 * The pre-change byte-wise implementation is kept here verbatim as the reference, with all
 * four return values compared
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package desensitize

import (
	"math/rand"
	"testing"
)

// keyShapeRef 改为无分支之前的逐字节实现，逐行照抄
// [EN] The byte-wise implementation as it was before going branchless, copied line for line
func keyShapeRef(key string) (first, last byte, n int, needsNorm bool) {
	for i := 0; i < len(key); i++ {
		c := key[i]
		if c == '_' || c == '-' || c == '.' {
			needsNorm = true
			continue
		}
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
			needsNorm = true
		}
		if n == 0 {
			first = c
		}
		last = c
		n++
	}
	return first, last, n, needsNorm
}

func checkShape(t *testing.T, key string) {
	t.Helper()
	f1, l1, n1, m1 := keyShape(key)
	f2, l2, n2, m2 := keyShapeRef(key)
	if f1 != f2 || l1 != l2 || n1 != n2 || m1 != m2 {
		t.Fatalf("键名 %q 形状不一致\n  无分支: 首=%q 末=%q 长=%d 需归一化=%v\n  参照  : 首=%q 末=%q 长=%d 需归一化=%v",
			key, f1, l1, n1, m1, f2, l2, n2, m2)
	}
}

// TestKeyShapeMatchesRef 规则表全量与各种变体下形状必须一致
// [EN] The shape must agree across the whole rule table and its variants
func TestKeyShapeMatchesRef(t *testing.T) {
	for _, k := range gateCorpus() {
		checkShape(t, k)
	}
}

// TestKeyShapeMatchesRefRandom 随机串同样不得分歧
//
// 字母表混入分隔符、数字、大写与非 ASCII：非 ASCII 字节既不是分隔符也不是大写，
// 必须原样计入长度并能当首末字符，这正是类别表最容易写错的一格
//
// [EN] Random strings must agree too
// The alphabet mixes separators, digits, uppercase and non-ASCII, the last of which is
// neither a separator nor uppercase and must count towards the length and be eligible as
// the first or last character, the cell of the class table easiest to get wrong
func TestKeyShapeMatchesRefRandom(t *testing.T) {
	alphabet := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-.\xff\x80\x7f")
	rng := rand.New(rand.NewSource(20261103))
	buf := make([]byte, 0, 32)
	for i := 0; i < 500000; i++ {
		buf = buf[:0]
		n := rng.Intn(25)
		for j := 0; j < n; j++ {
			buf = append(buf, alphabet[rng.Intn(len(alphabet))])
		}
		checkShape(t, string(buf))
	}
}

// TestKeyShapeSeparatorOnly 全分隔符的键名返回零值首末与零长度
// 这是无分支版本里唯一需要提前返回的分支，单独钉住
//
// [EN] An all-separator key returns zero first and last characters and zero length, the
// one early return in the branchless version, pinned on its own
func TestKeyShapeSeparatorOnly(t *testing.T) {
	for _, k := range []string{"_", "-", ".", "__", "_-.", "---", "....", "_-_-_"} {
		checkShape(t, k)
		if f, l, n, m := keyShape(k); f != 0 || l != 0 || n != 0 || !m {
			t.Errorf("键名 %q 期望 (0, 0, 0, true)，实得 (%q, %q, %d, %v)", k, f, l, n, m)
		}
	}
}
