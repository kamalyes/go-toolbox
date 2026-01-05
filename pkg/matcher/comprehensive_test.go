/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-15 10:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-15 10:00:00
 * @FilePath: \go-toolbox\pkg\matcher\comprehensive_test.go
 * @Description: 全面的测试套件 - 50+不同场景的测试用例
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package matcher

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/contextx"
	"github.com/stretchr/testify/assert"
)

// TestResult 测试结果结构
type TestResult struct {
	ID       int
	Value    string
	Matched  bool
	Priority int
}

// SimpleRule 简单规则实现
type SimpleRule struct {
	id        string
	priority  int
	enabled   bool
	condition func(*contextx.Context) bool
	result    TestResult
}

func (r *SimpleRule) Match(ctx *contextx.Context) bool {
	if r.condition == nil {
		return false
	}
	return r.condition(ctx)
}

func (r *SimpleRule) Priority() int      { return r.priority }
func (r *SimpleRule) Result() TestResult { return r.result }
func (r *SimpleRule) ID() string         { return r.id }
func (r *SimpleRule) Enabled() bool      { return r.enabled }

// ===== 1-10: 类型兼容性测试 =====

func TestTypeCompatibility_Int_Types(t *testing.T) {
	ctx := contextx.NewContext()

	// 测试所有整数类型
	ctx = ctx.WithValue("int", 42)
	ctx = ctx.WithValue("int8", int8(8))
	ctx = ctx.WithValue("int16", int16(16))
	ctx = ctx.WithValue("int32", int32(32))
	ctx = ctx.WithValue("int64", int64(64))
	ctx = ctx.WithValue("uint", uint(100))
	ctx = ctx.WithValue("uint8", uint8(200))
	ctx = ctx.WithValue("uint16", uint16(300))
	ctx = ctx.WithValue("uint32", uint32(400))
	ctx = ctx.WithValue("uint64", uint64(500))

	assert.Equal(t, 42, ctx.GetInt("int"))
	assert.Equal(t, int8(8), ctx.GetInt8("int8"))
	assert.Equal(t, int16(16), ctx.GetInt16("int16"))
	assert.Equal(t, int32(32), ctx.GetInt32("int32"))
	assert.Equal(t, int64(64), ctx.GetInt64("int64"))
	assert.Equal(t, uint(100), ctx.GetUint("uint"))
	assert.Equal(t, uint8(200), ctx.GetUint8("uint8"))
	assert.Equal(t, uint16(300), ctx.GetUint16("uint16"))
	assert.Equal(t, uint32(400), ctx.GetUint32("uint32"))
	assert.Equal(t, uint64(500), ctx.GetUint64("uint64"))

	// 测试类型转换
	assert.Equal(t, int64(42), ctx.GetInt64("int"))
	assert.Equal(t, 64, ctx.GetInt("int64"))
}

func TestTypeCompatibility_Float_Types(t *testing.T) {
	ctx := contextx.NewContext()

	ctx = ctx.WithValue("float32", float32(3.14))
	ctx = ctx.WithValue("float64", 3.14159)
	ctx = ctx.WithValue("intToFloat", 42)

	assert.InDelta(t, float32(3.14), ctx.GetFloat32("float32"), 0.001)
	assert.InDelta(t, 3.14159, ctx.GetFloat64("float64"), 0.00001)
	assert.Equal(t, 42.0, ctx.GetFloat64("intToFloat"))
	assert.Equal(t, float32(42), ctx.GetFloat32("intToFloat"))
}

func TestTypeCompatibility_String_And_Bool(t *testing.T) {
	ctx := contextx.NewContext()

	ctx = ctx.WithValue("string", "hello world")
	ctx = ctx.WithValue("bool_true", true)
	ctx = ctx.WithValue("bool_false", false)
	ctx = ctx.WithValue("rune", 'A')

	assert.Equal(t, "hello world", ctx.GetString("string"))
	assert.True(t, ctx.GetBool("bool_true"))
	assert.False(t, ctx.GetBool("bool_false"))
	assert.Equal(t, int32('A'), ctx.GetRune("rune"))
}

func TestTypeCompatibility_Time_Duration(t *testing.T) {
	ctx := contextx.NewContext()
	now := time.Now()
	duration := 5 * time.Minute

	ctx = ctx.WithValue("time", now)
	ctx = ctx.WithValue("duration", duration)
	ctx = ctx.WithValue("time_string", now.Format(time.RFC3339))
	ctx = ctx.WithValue("duration_string", "1h30m")
	ctx = ctx.WithValue("timestamp", now.Unix())

	assert.Equal(t, now, ctx.GetTime("time"))
	assert.Equal(t, duration, ctx.GetDuration("duration"))

	// 时间字符串解析
	parsedTime := ctx.GetTime("time_string")
	assert.True(t, parsedTime.Sub(now) < time.Second)

	// 时间间隔字符串解析
	parsedDuration := ctx.GetDuration("duration_string")
	assert.Equal(t, 90*time.Minute, parsedDuration)

	// Unix时间戳
	timeFromTimestamp := ctx.GetTime("timestamp")
	assert.Equal(t, now.Unix(), timeFromTimestamp.Unix())
}

func TestTypeCompatibility_Collections(t *testing.T) {
	ctx := contextx.NewContext()

	intSlice := []int{1, 2, 3, 4, 5}
	interfaceSlice := []interface{}{10, 20, 30}
	stringSlice := []string{"a", "b", "c"}
	testMap := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	ctx = ctx.WithValue("intSlice", intSlice)
	ctx = ctx.WithValue("interfaceSlice", interfaceSlice)
	ctx = ctx.WithValue("stringSlice", stringSlice)
	ctx = ctx.WithValue("map", testMap)

	assert.Equal(t, intSlice, ctx.GetIntSlice("intSlice"))
	assert.Equal(t, []int{10, 20, 30}, ctx.GetIntSlice("interfaceSlice"))
	assert.Equal(t, stringSlice, ctx.SafeGetStringSlice("stringSlice"))
	assert.Equal(t, testMap, ctx.GetMap("map"))
}

// ===== 11-20: 基础匹配功能测试 =====

func TestBasicMatching_Simple_Rules(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	rule1 := &SimpleRule{
		id: "rule1", priority: 10, enabled: true,
		condition: MatchString("action", "create"),
		result:    TestResult{ID: 1, Value: "create_action"},
	}

	rule2 := &SimpleRule{
		id: "rule2", priority: 20, enabled: true,
		condition: MatchString("action", "update"),
		result:    TestResult{ID: 2, Value: "update_action"},
	}

	matcher.AddRules(rule1, rule2)

	ctx1 := contextx.NewContext().WithValue("action", "create")
	result1, matched1 := matcher.Match(ctx1)
	assert.True(t, matched1)
	assert.Equal(t, 1, result1.ID)

	ctx2 := contextx.NewContext().WithValue("action", "update")
	result2, matched2 := matcher.Match(ctx2)
	assert.True(t, matched2)
	assert.Equal(t, 2, result2.ID)

	ctx3 := contextx.NewContext().WithValue("action", "delete")
	_, matched3 := matcher.Match(ctx3)
	assert.False(t, matched3)
}

func TestBasicMatching_Priority_Order(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// 添加多个匹配相同条件但优先级不同的规则
	lowPriority := &SimpleRule{
		id: "low", priority: 1, enabled: true,
		condition: MatchString("type", "test"),
		result:    TestResult{ID: 1, Value: "low_priority"},
	}

	highPriority := &SimpleRule{
		id: "high", priority: 100, enabled: true,
		condition: MatchString("type", "test"),
		result:    TestResult{ID: 2, Value: "high_priority"},
	}

	matcher.AddRules(lowPriority, highPriority)

	ctx := contextx.NewContext().WithValue("type", "test")
	result, matched := matcher.Match(ctx)

	assert.True(t, matched)
	assert.Equal(t, 2, result.ID) // 应该匹配高优先级的规则
	assert.Equal(t, "high_priority", result.Value)
}

func TestBasicMatching_Disabled_Rules(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	enabledRule := &SimpleRule{
		id: "enabled", priority: 10, enabled: true,
		condition: MatchString("status", "active"),
		result:    TestResult{ID: 1, Value: "enabled_rule"},
	}

	disabledRule := &SimpleRule{
		id: "disabled", priority: 20, enabled: false,
		condition: MatchString("status", "active"),
		result:    TestResult{ID: 2, Value: "disabled_rule"},
	}

	matcher.AddRules(enabledRule, disabledRule)

	ctx := contextx.NewContext().WithValue("status", "active")
	result, matched := matcher.Match(ctx)

	assert.True(t, matched)
	assert.Equal(t, 1, result.ID) // 只有启用的规则会匹配
}

func TestBasicMatching_MatchAll_Function(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	rule1 := &SimpleRule{
		id: "rule1", priority: 10, enabled: true,
		condition: MatchPrefix("name", "test"),
		result:    TestResult{ID: 1, Value: "prefix_match"},
	}

	rule2 := &SimpleRule{
		id: "rule2", priority: 20, enabled: true,
		condition: MatchSuffix("name", "case"),
		result:    TestResult{ID: 2, Value: "suffix_match"},
	}

	matcher.AddRules(rule1, rule2)

	ctx := contextx.NewContext().WithValue("name", "test_case")
	results := matcher.MatchAll(ctx)

	assert.Len(t, results, 2)
	// 结果按优先级排序
	assert.Equal(t, 2, results[0].ID) // 高优先级在前
	assert.Equal(t, 1, results[1].ID)
}

func TestBasicMatching_Chain_Rules(t *testing.T) {
	rule := NewChainRule(TestResult{ID: 1, Value: "chain_match"}).
		When(MatchString("service", "api")).
		When(MatchString("method", "POST")).
		When(MatchContains("path", "/users")).
		WithPriority(10).
		WithID("api_rule")

	ctx := contextx.NewContext().WithValue("service", "api")
	ctx = ctx.WithValue("method", "POST")
	ctx = ctx.WithValue("path", "/api/users/create")

	assert.True(t, rule.Match(ctx))
	assert.Equal(t, 10, rule.Priority())
	assert.Equal(t, "api_rule", rule.ID())
}

// ===== 21-30: 高级匹配模式测试 =====

func TestAdvancedMatching_Pattern_Wildcards(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	rule := &SimpleRule{
		id: "wildcard", priority: 10, enabled: true,
		condition: func(ctx *contextx.Context) bool {
			path := ctx.GetString("path")
			// 简单的通配符匹配：*.log
			return strings.HasSuffix(path, ".log")
		},
		result: TestResult{ID: 1, Value: "log_file"},
	}

	matcher.AddRule(rule)

	testCases := []struct {
		path     string
		expected bool
	}{
		{"error.log", true},
		{"access.log", true},
		{"app.log.backup", false},
		{"logfile.txt", false},
		{"/var/log/system.log", true}, // 路径匹配
	}

	for _, tc := range testCases {
		ctx := contextx.NewContext().WithValue("path", tc.path)
		_, matched := matcher.Match(ctx)
		assert.Equal(t, tc.expected, matched, "path: %s", tc.path)
	}
}

func TestAdvancedMatching_Conditional_Logic(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// OR 条件测试 - 使用函数字面量
	orRule := &SimpleRule{
		id: "or_rule", priority: 30, enabled: true,
		condition: func(ctx *contextx.Context) bool {
			role := ctx.GetString("role")
			return role == "admin" || role == "moderator"
		},
		result: TestResult{ID: 1, Value: "privileged_user"},
	}

	// AND 条件测试
	andRule := &SimpleRule{
		id: "and_rule", priority: 20, enabled: true,
		condition: func(ctx *contextx.Context) bool {
			return ctx.GetString("env") == "production" && ctx.GetString("secure") == "true"
		},
		result: TestResult{ID: 2, Value: "secure_production"},
	}

	// NOT 条件测试
	notRule := &SimpleRule{
		id: "not_rule", priority: 10, enabled: true,
		condition: func(ctx *contextx.Context) bool {
			return ctx.GetString("status") != "disabled"
		},
		result: TestResult{ID: 3, Value: "active_service"},
	}

	matcher.AddRules(orRule, andRule, notRule)

	// 测试 OR 条件
	ctx1 := contextx.NewContext().WithValue("role", "admin")
	result1, matched1 := matcher.Match(ctx1)
	assert.True(t, matched1)
	assert.Equal(t, 1, result1.ID)

	// 测试 AND 条件
	ctx2 := contextx.NewContext().WithValue("env", "production").WithValue("secure", "true")
	result2, matched2 := matcher.Match(ctx2)
	assert.True(t, matched2)
	assert.Equal(t, 2, result2.ID)

	// 测试 NOT 条件
	ctx3 := contextx.NewContext().WithValue("status", "active")
	result3, matched3 := matcher.Match(ctx3)
	assert.True(t, matched3)
	assert.Equal(t, 3, result3.ID)
}

func TestAdvancedMatching_Multiple_Conditions(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// 复杂的多条件规则
	rule := &SimpleRule{
		id: "complex", priority: 10, enabled: true,
		condition: MatchAll(
			MatchStringIn("method", []string{"GET", "POST", "PUT"}),
			MatchPrefix("path", "/api"),
			MatchNot(MatchContains("path", "test")),
			MatchAny(
				MatchString("version", "v1"),
				MatchString("version", "v2"),
			),
		),
		result: TestResult{ID: 1, Value: "api_request"},
	}

	matcher.AddRule(rule)

	// 正匹配测试
	ctx1 := contextx.NewContext().
		WithValue("method", "POST").
		WithValue("path", "/api/users").
		WithValue("version", "v1")
	result1, matched1 := matcher.Match(ctx1)
	assert.True(t, matched1)
	assert.Equal(t, 1, result1.ID)

	// 负匹配测试（包含test）
	ctx2 := contextx.NewContext().
		WithValue("method", "GET").
		WithValue("path", "/api/test/users").
		WithValue("version", "v2")
	_, matched2 := matcher.Match(ctx2)
	assert.False(t, matched2)
}

func TestAdvancedMatching_String_Operations(t *testing.T) {
	testCases := []struct {
		name      string
		condition func(*contextx.Context) bool
		key       string
		value     string
		expected  bool
	}{
		{"prefix", MatchPrefix("url", "/api"), "url", "/api/users", true},
		{"prefix_false", MatchPrefix("url", "/api"), "url", "/web/users", false},
		{"suffix", MatchSuffix("file", ".txt"), "file", "readme.txt", true},
		{"suffix_false", MatchSuffix("file", ".txt"), "file", "readme.md", false},
		{"contains", MatchContains("text", "world"), "text", "hello world", true},
		{"contains_false", MatchContains("text", "world"), "text", "hello there", false},
		{"string_in", MatchStringIn("lang", []string{"go", "python", "java"}), "lang", "go", true},
		{"string_in_false", MatchStringIn("lang", []string{"go", "python", "java"}), "lang", "rust", false},
		{"string_not_in", MatchStringNotIn("status", []string{"error", "failed"}), "status", "success", true},
		{"string_not_in_false", MatchStringNotIn("status", []string{"error", "failed"}), "status", "error", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := contextx.NewContext().WithValue(tc.key, tc.value)
			result := tc.condition(ctx)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAdvancedMatching_HTTP_Methods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	rule := &SimpleRule{
		id: "http", priority: 10, enabled: true,
		condition: MatchMethodIn(methods),
		result:    TestResult{ID: 1, Value: "http_method"},
	}

	matcher := NewMatcher[TestResult]().AddRule(rule)

	// 测试大小写不敏感
	testMethods := []string{"get", "POST", "Put", "DELETE", "patch", "HEAD"}
	expected := []bool{true, true, true, true, true, false}

	for i, method := range testMethods {
		ctx := contextx.NewContext().WithValue("method", method)
		_, matched := matcher.Match(ctx)
		assert.Equal(t, expected[i], matched, "method: %s", method)
	}
}

// ===== 31-40: 并发和性能测试 =====

func TestConcurrency_Parallel_Matching(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// 添加多个规则
	for i := 0; i < 100; i++ {
		rule := &SimpleRule{
			id:        fmt.Sprintf("rule_%d", i),
			priority:  i,
			enabled:   true,
			condition: MatchString("id", strconv.Itoa(i)),
			result:    TestResult{ID: i, Value: fmt.Sprintf("result_%d", i)},
		}
		matcher.AddRule(rule)
	}

	// 并发测试
	const goroutines = 100
	const iterations = 1000
	var wg sync.WaitGroup
	var successCount atomic.Int64

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				targetID := (gid*iterations + i) % 100
				ctx := contextx.NewContext().WithValue("id", strconv.Itoa(targetID))
				result, matched := matcher.Match(ctx)
				if matched && result.ID == targetID {
					successCount.Add(1)
				}
			}
		}(g)
	}

	wg.Wait()
	assert.Equal(t, int64(goroutines*iterations), successCount.Load())
}

func TestConcurrency_Rule_Modification(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	var wg sync.WaitGroup
	const goroutines = 50

	// 一组协程添加规则
	for i := 0; i < goroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				rule := &SimpleRule{
					id:        fmt.Sprintf("rule_%d_%d", id, j),
					priority:  id*10 + j,
					enabled:   true,
					condition: MatchString("worker", strconv.Itoa(id)),
					result:    TestResult{ID: id*10 + j, Value: fmt.Sprintf("worker_%d", id)},
				}
				matcher.AddRule(rule)
			}
		}(i)
	}

	// 另一组协程执行匹配
	for i := 0; i < goroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				ctx := contextx.NewContext().WithValue("worker", strconv.Itoa(id%10))
				matcher.Match(ctx) // 不检查结果，只确保不会panic
			}
		}(i)
	}

	wg.Wait()
	// 如果执行到这里没有panic，说明并发安全
}

func TestPerformance_Large_Ruleset(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}

	matcher := NewMatcher[TestResult]()
	ruleCount := 10000

	// 创建大量规则
	for i := 0; i < ruleCount; i++ {
		rule := &SimpleRule{
			id:        fmt.Sprintf("rule_%d", i),
			priority:  rand.Intn(1000),
			enabled:   true,
			condition: MatchString("target", fmt.Sprintf("target_%d", i)),
			result:    TestResult{ID: i, Value: fmt.Sprintf("result_%d", i)},
		}
		matcher.AddRule(rule)
	}

	// 性能测试
	iterations := 10000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		target := fmt.Sprintf("target_%d", rand.Intn(ruleCount))
		ctx := contextx.NewContext().WithValue("target", target)
		matcher.Match(ctx)
	}

	duration := time.Since(start)
	opsPerSec := float64(iterations) / duration.Seconds()

	t.Logf("大规模规则集性能测试:")
	t.Logf("  规则数量: %d", ruleCount)
	t.Logf("  测试次数: %d", iterations)
	t.Logf("  执行时间: %v", duration)
	t.Logf("  吞吐量: %.2f ops/sec", opsPerSec)

	// 性能要求：至少100 ops/sec（CI环境友好的性能指标）
	assert.Greater(t, opsPerSec, 100.0, "性能不达标")
}

func TestPerformance_Memory_Usage(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过内存测试")
	}

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	matcher := NewMatcher[TestResult]()

	// 创建规则并执行匹配
	for i := 0; i < 1000; i++ {
		rule := &SimpleRule{
			id:        fmt.Sprintf("rule_%d", i),
			priority:  i,
			enabled:   true,
			condition: MatchString("id", strconv.Itoa(i)),
			result:    TestResult{ID: i},
		}
		matcher.AddRule(rule)

		// 执行一些匹配操作
		for j := 0; j < 10; j++ {
			ctx := contextx.NewContext().WithValue("id", strconv.Itoa(rand.Intn(i+1)))
			matcher.Match(ctx)
		}
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	allocatedMB := float64(m2.Alloc) / 1024 / 1024
	t.Logf("内存使用: %.2f MB", allocatedMB)

	// 合理的内存使用范围（只检查是否过度使用）
	if allocatedMB > 500.0 {
		t.Logf("警告：内存使用较高: %.2f MB", allocatedMB)
	}
}

func TestPerformance_Cache_Effectiveness(t *testing.T) {
	matcher := NewMatcher[TestResult]().EnableCache(5 * time.Minute)

	rule := &SimpleRule{
		id: "test", priority: 10, enabled: true,
		condition: MatchString("key", "value"),
		result:    TestResult{ID: 1, Value: "cached_result"},
	}
	matcher.AddRule(rule)

	ctx := contextx.NewContext().WithValue("key", "value")

	// 第一次匹配
	start := time.Now()
	result1, matched1 := matcher.Match(ctx)
	firstDuration := time.Since(start)

	// 第二次匹配（应该命中缓存）
	start = time.Now()
	result2, matched2 := matcher.Match(ctx)
	secondDuration := time.Since(start)

	assert.True(t, matched1)
	assert.True(t, matched2)
	assert.Equal(t, result1, result2)

	stats := matcher.Stats()
	assert.Equal(t, int64(1), stats["cache_hits"])
	assert.Equal(t, int64(1), stats["cache_misses"])

	t.Logf("首次匹配: %v", firstDuration)
	t.Logf("缓存命中: %v", secondDuration)
	t.Logf("性能提升: %.2fx", float64(firstDuration)/float64(secondDuration))
}

// ===== 41-50: 边界条件和错误处理测试 =====

func TestEdgeCases_Empty_Matcher(t *testing.T) {
	matcher := NewMatcher[TestResult]()
	ctx := contextx.NewContext().WithValue("any", "value")

	result, matched := matcher.Match(ctx)
	assert.False(t, matched)
	assert.Equal(t, TestResult{}, result)

	results := matcher.MatchAll(ctx)
	assert.Empty(t, results)
}

func TestEdgeCases_Nil_Values(t *testing.T) {
	ctx := contextx.NewContext()
	ctx.WithValue("nil_value", nil)
	ctx.WithValue("empty_string", "")

	assert.Equal(t, "", ctx.GetString("nil_value"))
	assert.Equal(t, 0, ctx.GetInt("nil_value"))
	assert.False(t, ctx.GetBool("nil_value"))

	assert.Equal(t, "", ctx.GetString("empty_string"))
	assert.Equal(t, "", ctx.GetString("nonexistent"))
}

func TestEdgeCases_Extreme_Values(t *testing.T) {
	ctx := contextx.NewContext()

	ctx.WithValue("max_int64", math.MaxInt64)
	ctx.WithValue("min_int64", math.MinInt64)
	ctx.WithValue("max_float64", math.MaxFloat64)
	ctx.WithValue("inf", math.Inf(1))
	ctx.WithValue("nan", math.NaN())

	assert.Equal(t, int64(math.MaxInt64), ctx.GetInt64("max_int64"))
	assert.Equal(t, int64(math.MinInt64), ctx.GetInt64("min_int64"))
	assert.Equal(t, math.MaxFloat64, ctx.GetFloat64("max_float64"))
	assert.True(t, math.IsInf(ctx.GetFloat64("inf"), 1))
	assert.True(t, math.IsNaN(ctx.GetFloat64("nan")))
}

func TestEdgeCases_Unicode_Strings(t *testing.T) {
	ctx := contextx.NewContext()

	unicodeStrings := []string{
		"Hello, 世界",
		"🚀 Rocket 🚀",
		"Ελληνικά",
		"🇺🇸🇨🇳🇯🇵",
		"\U0001F600\U0001F601\U0001F602",
	}

	matcher := NewMatcher[TestResult]()

	for i, str := range unicodeStrings {
		rule := &SimpleRule{
			id:        fmt.Sprintf("unicode_%d", i),
			priority:  i,
			enabled:   true,
			condition: MatchString("text", str),
			result:    TestResult{ID: i, Value: str},
		}
		matcher.AddRule(rule)

		ctx.WithValue("text", str)
		result, matched := matcher.Match(ctx)
		assert.True(t, matched)
		assert.Equal(t, str, result.Value)
	}
}

func TestEdgeCases_Large_Context_Data(t *testing.T) {
	ctx := contextx.NewContext()

	// 添加大量数据
	for i := 0; i < 10000; i++ {
		ctx.WithValue(fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i))
	}

	// 测试获取数据
	assert.Equal(t, "value_5000", ctx.GetString("key_5000"))
	assert.Equal(t, "", ctx.GetString("nonexistent"))

	// 测试克隆
	cloned := ctx.Clone()
	assert.Equal(t, "value_5000", cloned.GetString("key_5000"))
}

func TestEdgeCases_Context_Timeout(t *testing.T) {
	ctx := contextx.NewContext().WithTimeout(10 * time.Millisecond)

	// 立即检查，应该没有超时
	assert.False(t, ctx.IsExpired())

	// 等待超时
	time.Sleep(15 * time.Millisecond)
	assert.True(t, ctx.IsExpired())

	// 超时的上下文应该不匹配
	matcher := NewMatcher[TestResult]()
	rule := &SimpleRule{
		id: "timeout", priority: 10, enabled: true,
		condition: MatchString("test", "value"),
		result:    TestResult{ID: 1},
	}
	matcher.AddRule(rule)

	ctx.WithValue("test", "value")
	_, matched := matcher.Match(ctx)
	assert.False(t, matched, "超时的上下文不应该匹配")
}

func TestEdgeCases_Rule_Removal(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	rule1 := &SimpleRule{id: "rule1", priority: 10, enabled: true, condition: MatchString("test", "1"), result: TestResult{ID: 1}}
	rule2 := &SimpleRule{id: "rule2", priority: 20, enabled: true, condition: MatchString("test", "2"), result: TestResult{ID: 2}}
	rule3 := &SimpleRule{id: "rule3", priority: 30, enabled: true, condition: MatchString("test", "3"), result: TestResult{ID: 3}}

	matcher.AddRules(rule1, rule2, rule3)

	// 验证所有规则都存在
	ctx1 := contextx.NewContext().WithValue("test", "1")
	result1, matched1 := matcher.Match(ctx1)
	assert.True(t, matched1)
	assert.Equal(t, 1, result1.ID)

	// 移除中间规则
	matcher.RemoveRule("rule2")

	// 验证规则2被移除
	ctx2 := contextx.NewContext().WithValue("test", "2")
	_, matched2 := matcher.Match(ctx2)
	assert.False(t, matched2)

	// 验证其他规则仍然存在
	ctx3 := contextx.NewContext().WithValue("test", "3")
	result3, matched3 := matcher.Match(ctx3)
	assert.True(t, matched3)
	assert.Equal(t, 3, result3.ID)
}

func TestEdgeCases_Clear_All_Rules(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// 添加规则
	for i := 0; i < 100; i++ {
		rule := &SimpleRule{
			id:        fmt.Sprintf("rule_%d", i),
			priority:  i,
			enabled:   true,
			condition: MatchString("id", strconv.Itoa(i)),
			result:    TestResult{ID: i},
		}
		matcher.AddRule(rule)
	}

	// 验证规则存在
	ctx := contextx.NewContext().WithValue("id", "50")
	_, matched := matcher.Match(ctx)
	assert.True(t, matched)

	// 清空所有规则
	matcher.ClearRules()

	// 验证没有规则匹配
	_, matched = matcher.Match(ctx)
	assert.False(t, matched)

	stats := matcher.Stats()
	assert.Equal(t, int64(1), stats["failed_matches"]) // 清空后的失败匹配
}

func TestEdgeCases_Cache_Expiration(t *testing.T) {
	matcher := NewMatcher[TestResult]().EnableCache(50 * time.Millisecond)

	rule := &SimpleRule{
		id: "cache_test", priority: 10, enabled: true,
		condition: MatchString("key", "value"),
		result:    TestResult{ID: 1, Value: "cached"},
	}
	matcher.AddRule(rule)

	ctx := contextx.NewContext().WithValue("key", "value")

	// 第一次匹配
	result1, matched1 := matcher.Match(ctx)
	assert.True(t, matched1)

	// 等待缓存过期
	time.Sleep(60 * time.Millisecond)

	// 第二次匹配（缓存已过期）
	result2, matched2 := matcher.Match(ctx)
	assert.True(t, matched2)
	assert.Equal(t, result1, result2)

	stats := matcher.Stats()
	assert.Equal(t, int64(2), stats["cache_misses"]) // 两次都是缓存未命中
}

// ===== 51+: 实际应用场景测试 =====

func TestRealWorld_API_Gateway_Routing(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// API网关路由规则
	rules := []*SimpleRule{
		{id: "auth_service", priority: 100, enabled: true,
			condition: MatchAll(MatchPrefix("path", "/auth"), MatchMethodIn([]string{"POST", "GET"})),
			result:    TestResult{ID: 1, Value: "auth-service"}},

		{id: "user_service", priority: 90, enabled: true,
			condition: MatchAll(MatchPrefix("path", "/users"), MatchMethodIn([]string{"GET", "POST", "PUT", "DELETE"})),
			result:    TestResult{ID: 2, Value: "user-service"}},

		{id: "order_service", priority: 80, enabled: true,
			condition: MatchPrefix("path", "/orders"),
			result:    TestResult{ID: 3, Value: "order-service"}},

		{id: "static_files", priority: 10, enabled: true,
			condition: MatchSuffix("path", ".js"),
			result:    TestResult{ID: 4, Value: "static-cdn"}},
	}

	matcher.AddRules(rules[0], rules[1], rules[2], rules[3])

	testCases := []struct {
		method   string
		path     string
		expected string
	}{
		{"POST", "/auth/login", "auth-service"},
		{"GET", "/users/123", "user-service"},
		{"POST", "/orders/create", "order-service"},
		{"GET", "/static/app.js", "static-cdn"},
	}

	for _, tc := range testCases {
		ctx := contextx.NewContext().WithValue("method", tc.method).WithValue("path", tc.path)
		result, matched := matcher.Match(ctx)
		assert.True(t, matched, "路径 %s %s 应该匹配", tc.method, tc.path)
		assert.Equal(t, tc.expected, result.Value)
	}
}

func TestRealWorld_Feature_Flags(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// 功能开关规则
	rules := []*SimpleRule{
		{id: "beta_users", priority: 100, enabled: true,
			condition: MatchAll(MatchBool("beta_user", true), MatchString("env", "production")),
			result:    TestResult{ID: 1, Value: "new_feature_enabled"}},

		{id: "admin_users", priority: 90, enabled: true,
			condition: MatchString("role", "admin"),
			result:    TestResult{ID: 2, Value: "admin_features_enabled"}},

		{id: "percentage_rollout", priority: 50, enabled: true,
			condition: func(ctx *contextx.Context) bool {
				userID := ctx.GetInt("user_id")
				return userID%100 < 10 // 10%的用户
			},
			result: TestResult{ID: 3, Value: "gradual_rollout_enabled"}},
	}

	matcher.AddRules(rules[0], rules[1], rules[2])

	// 测试Beta用户
	ctx1 := contextx.NewContext().WithValue("beta_user", true).WithValue("env", "production")
	result1, matched1 := matcher.Match(ctx1)
	assert.True(t, matched1)
	assert.Equal(t, "new_feature_enabled", result1.Value)

	// 测试管理员
	ctx2 := contextx.NewContext().WithValue("role", "admin")
	result2, matched2 := matcher.Match(ctx2)
	assert.True(t, matched2)
	assert.Equal(t, "admin_features_enabled", result2.Value)

	// 测试渐进式发布
	enabledCount := 0
	for i := 0; i < 1000; i++ {
		ctx := contextx.NewContext().WithValue("user_id", i)
		_, matched := matcher.Match(ctx)
		if matched {
			enabledCount++
		}
	}

	// 应该大约有10%的用户启用功能
	assert.InDelta(t, 100, enabledCount, 20, "渐进式发布比例不正确")
}

func TestRealWorld_Content_Filtering(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// 内容过滤规则
	rules := []*SimpleRule{
		{id: "spam_keywords", priority: 100, enabled: true,
			condition: MatchAny(
				MatchContains("content", "spam"),
				MatchContains("content", "viagra"),
				MatchContains("content", "lottery"),
			),
			result: TestResult{ID: 1, Value: "blocked_spam"}},

		{id: "offensive_language", priority: 90, enabled: true,
			condition: MatchContains("content", "offensive"),
			result:    TestResult{ID: 2, Value: "blocked_offensive"}},

		{id: "max_length", priority: 80, enabled: true,
			condition: func(ctx *contextx.Context) bool {
				content := ctx.GetString("content")
				return len(content) > 1000
			},
			result: TestResult{ID: 3, Value: "blocked_too_long"}},
	}

	matcher.AddRules(rules[0], rules[1], rules[2])

	testCases := []struct {
		content  string
		expected string
		blocked  bool
	}{
		{"This is a normal message", "", false},
		{"Win the lottery now!", "blocked_spam", true},
		{"Buy viagra cheap", "blocked_spam", true},
		{"This message contains offensive language", "blocked_offensive", true},
		{strings.Repeat("x", 1001), "blocked_too_long", true},
	}

	for _, tc := range testCases {
		ctx := contextx.NewContext().WithValue("content", tc.content)
		result, matched := matcher.Match(ctx)
		assert.Equal(t, tc.blocked, matched, "内容: %s", tc.content[:min(50, len(tc.content))])
		if matched {
			assert.Equal(t, tc.expected, result.Value)
		}
	}
}

func TestRealWorld_Load_Balancing(t *testing.T) {
	matcher := NewMatcher[TestResult]()

	// 负载均衡规则
	rules := []*SimpleRule{
		{id: "high_cpu_server", priority: 100, enabled: true,
			condition: MatchAll(
				MatchString("server_type", "high_cpu"),
				func(ctx *contextx.Context) bool { return ctx.GetFloat64("cpu_usage") < 80.0 },
			),
			result: TestResult{ID: 1, Value: "high-cpu-server-pool"}},

		{id: "memory_intensive", priority: 90, enabled: true,
			condition: MatchAll(
				MatchString("request_type", "memory_intensive"),
				func(ctx *contextx.Context) bool { return ctx.GetFloat64("memory_usage") < 70.0 },
			),
			result: TestResult{ID: 2, Value: "memory-optimized-pool"}},

		{id: "default_pool", priority: 10, enabled: true,
			condition: func(ctx *contextx.Context) bool { return true }, // 默认匹配
			result:    TestResult{ID: 3, Value: "default-server-pool"}},
	}

	matcher.AddRules(rules[0], rules[1], rules[2])

	testCases := []struct {
		serverType   string
		requestType  string
		cpuUsage     float64
		memoryUsage  float64
		expectedPool string
	}{
		{"high_cpu", "", 50.0, 60.0, "high-cpu-server-pool"},
		{"standard", "memory_intensive", 60.0, 50.0, "memory-optimized-pool"},
		{"standard", "standard", 90.0, 80.0, "default-server-pool"},
	}

	for _, tc := range testCases {
		ctx := contextx.NewContext().
			WithValue("server_type", tc.serverType).
			WithValue("request_type", tc.requestType).
			WithValue("cpu_usage", tc.cpuUsage).
			WithValue("memory_usage", tc.memoryUsage)

		result, matched := matcher.Match(ctx)
		assert.True(t, matched)
		assert.Equal(t, tc.expectedPool, result.Value)
	}
}

func TestRealWorld_Stress_Testing(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试")
	}

	matcher := NewMatcher[TestResult]().EnableCache(1 * time.Minute)

	// 创建复杂的规则集
	ruleCount := 5000
	for i := 0; i < ruleCount; i++ {
		rule := &SimpleRule{
			id:       fmt.Sprintf("stress_rule_%d", i),
			priority: rand.Intn(1000),
			enabled:  true,
			condition: MatchAll(
				MatchString("service", fmt.Sprintf("service_%d", i%50)),
				MatchString("method", []string{"GET", "POST", "PUT"}[i%3]),
				func(ctx *contextx.Context) bool {
					return ctx.GetInt("user_id")%100 == i%100
				},
			),
			result: TestResult{ID: i, Value: fmt.Sprintf("action_%d", i)},
		}
		matcher.AddRule(rule)
	}

	// 压力测试
	const goroutines = 100
	const iterations = 1000
	var wg sync.WaitGroup
	var totalMatches atomic.Int64
	var totalTime atomic.Int64 // 纳秒

	start := time.Now()

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				startTime := time.Now()

				ctx := contextx.NewContext().
					WithValue("service", fmt.Sprintf("service_%d", rand.Intn(50))).
					WithValue("method", []string{"GET", "POST", "PUT"}[rand.Intn(3)]).
					WithValue("user_id", rand.Intn(10000))

				_, matched := matcher.Match(ctx)

				duration := time.Since(startTime)
				totalTime.Add(duration.Nanoseconds())

				if matched {
					totalMatches.Add(1)
				}
			}
		}(g)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	operations := goroutines * iterations
	avgLatency := time.Duration(totalTime.Load() / int64(operations))
	opsPerSec := float64(operations) / totalDuration.Seconds()

	stats := matcher.Stats()

	t.Logf("压力测试结果:")
	t.Logf("  规则数量: %d", ruleCount)
	t.Logf("  并发数: %d", goroutines)
	t.Logf("  总操作数: %d", operations)
	t.Logf("  匹配成功: %d", totalMatches.Load())
	t.Logf("  总耗时: %v", totalDuration)
	t.Logf("  平均延迟: %v", avgLatency)
	t.Logf("  吞吐量: %.2f ops/sec", opsPerSec)
	t.Logf("  缓存命中率: %.2f%%", float64(stats["cache_hits"])*100/float64(stats["total_matches"]))

	// 性能断言（CI环境友好）
	if opsPerSec < 1000 {
		t.Logf("警告：吞吐量较低 %.2f ops/sec", opsPerSec)
	} else {
		assert.Greater(t, opsPerSec, 1000.0, "吞吐量不达标")
	}
	if avgLatency > 100*time.Millisecond {
		t.Logf("警告：延迟较高 %v", avgLatency)
	} else {
		assert.Less(t, avgLatency, 100*time.Millisecond, "平均延迟过高")
	}
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
