/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-12 15:27:26
 * @FilePath: \go-toolbox\tests\random_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/stretchr/testify/assert"
)

func TestRandInt(t *testing.T) {
	min := 10
	max := 20

	result := random.RandInt(min, max)

	assert.GreaterOrEqual(t, result, min, "Expected result to be greater than or equal to min")
	assert.LessOrEqual(t, result, max, "Expected result to be less than or equal to max")
}

func TestRandFloat(t *testing.T) {
	min := 10.5
	max := 20.5

	result := random.RandFloat(min, max)

	assert.GreaterOrEqual(t, result, min, "Expected result to be greater than or equal to min")
	assert.LessOrEqual(t, result, max, "Expected result to be less than or equal to max")
}

func TestRandString(t *testing.T) {
	str := random.RandString(10, random.CAPITAL|random.LOWERCASE|random.SPECIAL|random.NUMBER)

	assert.Len(t, str, 10, "Expected string length to be 10")
}

func TestRandNumber(t *testing.T) {
	length := 10
	result := random.RandNumber(length)

	// 使用 assert 检查长度和内容
	assert.Len(t, result, length, "Expected length should be %d", length)

	// 检查结果是否只包含数字
	digitMap := make(map[rune]bool)
	for _, char := range random.DEC_BYTES {
		digitMap[char] = true
	}

	for _, char := range result {
		assert.True(t, digitMap[char], "Result contains non-digit character: %c", char)
	}

	// 测试自定义字节集
	customBytes := "1234567890"
	resultCustom := random.RandNumber(length, customBytes)
	assert.Len(t, resultCustom, length, "Expected length should be %d for custom bytes", length)

	customDigitMap := make(map[rune]bool)
	for _, char := range customBytes {
		customDigitMap[char] = true
	}

	for _, char := range resultCustom {
		assert.True(t, customDigitMap[char], "Result contains character not in custom bytes: %c", char)
	}
}

func TestRandHex(t *testing.T) {
	bytesLen := 5
	result := random.RandHex(bytesLen)

	// 使用 assert 检查长度和内容
	assert.Len(t, result, bytesLen*2, "Expected length should be %d", bytesLen*2)

	// 检查结果是否只包含 hex 字符
	hexMap := make(map[rune]bool)
	for _, char := range random.HEX_BYTES {
		hexMap[char] = true
	}

	for _, char := range result {
		assert.True(t, hexMap[char], "Result contains non-hex character: %c", char)
	}

	// 测试自定义字节集
	customHexBytes := "abcdef"
	resultCustom := random.RandHex(bytesLen, customHexBytes)
	assert.Len(t, resultCustom, bytesLen*2, "Expected length should be %d for custom bytes", bytesLen*2)

	customHexMap := make(map[rune]bool)
	for _, char := range customHexBytes {
		customHexMap[char] = true
	}

	for _, char := range resultCustom {
		assert.True(t, customHexMap[char], "Result contains character not in custom bytes: %c", char)
	}
}

func TestRandNum(t *testing.T) {
	length := 6
	num := random.RandNumber(length)

	assert.Len(t, num, length, "Expected number length to be 6")
}

func TestNewRand(t *testing.T) {
	rd := random.NewRand(1)
	assert.Equal(t, int64(5577006791947779410), rd.Int63())

	rd = random.NewRand()
	for i := 1; i < 1000; i++ {
		assert.Equal(t, true, rd.Intn(i) < i)
		assert.Equal(t, true, rd.Int63n(int64(i)) < int64(i))
		assert.Equal(t, true, random.NewRand().Intn(i) < i)
		assert.Equal(t, true, random.NewRand().Int63n(int64(i)) < int64(i))
	}
}

func TestFRandInt(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, random.FRandInt(1, 2) == 1)
	assert.Equal(t, true, random.FRandInt(-1, 0) == -1)
	assert.Equal(t, true, random.FRandInt(0, 5) >= 0)
	assert.Equal(t, true, random.FRandInt(0, 5) < 5)
	assert.Equal(t, 2, random.FRandInt(2, 2))
	assert.Equal(t, 2, random.FRandInt(3, 2))
}

func TestFRandUint32(t *testing.T) {
	t.Parallel()
	assert.Equal(t, true, random.FRandUint32(1, 2) == 1)
	assert.Equal(t, true, random.FRandUint32(0, 5) < 5)
	assert.Equal(t, uint32(2), random.FRandUint32(2, 2))
	assert.Equal(t, uint32(2), random.FRandUint32(3, 2))
}

func TestFastIntn(t *testing.T) {
	t.Parallel()
	for i := 1; i < 10000; i++ {
		assert.Equal(t, true, random.FastRandn(uint32(i)) < uint32(i))
		assert.Equal(t, true, random.FastIntn(i) < i)
	}
	assert.Equal(t, 0, random.FastIntn(-2))
	assert.Equal(t, 0, random.FastIntn(0))
	assert.Equal(t, true, random.FastIntn(math.MaxUint32) < math.MaxUint32)
	assert.Equal(t, true, random.FastIntn(math.MaxInt64) < math.MaxInt64)
}

func TestFRandString(t *testing.T) {
	t.Parallel()
	fns := []func(n int) string{random.FRandString, random.FRandAlphaString, random.FRandHexString, random.FRandDecString}
	ss := []string{random.LETTER_BYTES, random.ALPHA_BYTES, random.HEX_BYTES, random.DEC_BYTES}
	for i, fn := range fns {
		a, b := fn(777), fn(777)
		assert.Equal(t, 777, len(a))
		assert.NotEqual(t, a, b)
		assert.Equal(t, "", fn(-1))
		for _, s := range ss[i] {
			assert.True(t, strings.ContainsRune(a, s))
		}
	}
}

// func TestFRandBytesLetters(t *testing.T) {
// 	t.Parallel()
// 	letters := ""
// 	assert.Nil(t, random.FRandBytesLetters(10, letters))
// 	letters = "a"
// 	assert.Nil(t, random.FRandBytesLetters(10, letters))
// 	letters = "ab"
// 	s := convert.B2S(random.FRandBytesLetters(10, letters))
// 	assert.Equal(t, 10, len(s))
// 	assert.True(t, strings.Contains(s, "a"))
// 	assert.True(t, strings.Contains(s, "b"))
// 	letters = "xxxxxxxxxxxx"
// 	s = convert.B2S(random.FRandBytesLetters(100, letters))
// 	assert.Equal(t, 100, len(s))
// 	assert.Equal(t, strings.Repeat("x", 100), s)
// }

var (
	testString = "  Fufu 中　文\u2728->?\n*\U0001F63A   "
	testBytes  = []byte(testString)
)

func TestB2S(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		b := random.FRandBytes(64)
		assert.Equal(t, string(b), convert.B2S(b))
	}

	expected := testString
	actual := convert.B2S([]byte(expected))
	assert.Equal(t, expected, actual)

	assert.Equal(t, true, convert.B2S(nil) == "")
	assert.Equal(t, testString, convert.B2S(testBytes))
}

func TestS2B(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		s := random.RandNumber(64)
		expected := []byte(s)
		actual := convert.S2B(s)
		assert.Equal(t, expected, actual)
		assert.Equal(t, len(expected), len(actual))
	}

	expected := testString
	actual := convert.S2B(expected)
	assert.Equal(t, []byte(expected), actual)

	assert.Equal(t, true, convert.S2B("") == nil)
	assert.Equal(t, testBytes, convert.S2B(testString))
}

func TestFRandBytesJSON(t *testing.T) {
	length := 16 // 测试生成的随机字节长度
	// 调用 FRandBytesJSON 函数
	jsonStr, err := random.FRandBytesJSON(length)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 检查 JSON 字符串是否有效
	var result []byte
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	// 检查生成的随机字节长度
	if len(result) != length {
		t.Errorf("Expected length %d, got %d", length, len(result))
	}
}

// 定义测试模型
type TestModel struct {
	Name      string         `json:"name"`
	Age       int            `json:"age"`
	Salary    float64        `json:"salary"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	Tags      []string       `json:"tags"`
	Settings  map[string]int `json:"settings"`
}

func TestRandStringSlice(t *testing.T) {
	count := 5
	length := 10
	mode := random.CAPITAL
	result := random.RandStringSlice(count, length, mode)

	// 验证生成的切片长度
	assert.Equal(t, count, len(result), "生成的切片长度应与请求的 count 相等")
}

// TestGenerateRandModel 测试 GenerateRandModel 函数
func TestGenerateRandModel(t *testing.T) {
	// 创建一个 TestModel 的实例
	model := &TestModel{}

	// 调用 GenerateRandModel
	modelResult, jsonResult, err := random.GenerateRandModel(model)
	if err != nil || modelResult == nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log("jsonResult", jsonResult)
	// 验证返回的 JSON 字符串是否有效
	var resultMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResult), &resultMap); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	// 验证字段是否被填充
	if resultMap["name"] == "" || resultMap["age"] == nil || resultMap["salary"] == nil || resultMap["is_active"] == nil {
		t.Error("Expected fields to be populated, but they are not")
	}
}

// Address 示例嵌套结构体
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
}

// User 示例结构体
type User struct {
	Name       string         `json:"name"`
	Age        *int           `json:"age"` // 指针类型
	Height     float64        `json:"height"`
	IsActive   bool           `json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	Hobbies    []string       `json:"hobbies"`
	Attributes map[string]int `json:"attributes"`
	Address    *Address       `json:"address"` // 指针类型
}

func TestGenerateRandModelComplex(t *testing.T) {
	// 创建一个 User 结构体的指针
	user := &User{}

	// 生成随机模型
	model, jsonOutput, err := random.GenerateRandModel(user)
	if err != nil {
		t.Fatalf("Error generating random model: %v", err)
	}

	// 验证生成的模型不为 nil
	if model == nil {
		t.Fatal("Generated model is nil")
	}

	// 验证 JSON 输出不为空
	if jsonOutput == "" {
		t.Fatal("JSON output is empty")
	}

	// 验证 JSON 格式有效
	var js json.RawMessage
	if err := json.Unmarshal([]byte(jsonOutput), &js); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}

	// 验证指针字段是否被正确填充
	userPtr := model.(*User)
	if userPtr.Age == nil {
		t.Fatal("Age pointer is nil, expected a value")
	}

	// 验证嵌套结构体是否被正确填充
	if userPtr.Address == nil {
		t.Fatal("Address pointer is nil, expected a value")
	}

	// 验证切片和映射是否被正确填充
	if len(userPtr.Hobbies) == 0 {
		t.Fatal("Hobbies slice is empty, expected at least one value")
	}

	if len(userPtr.Attributes) == 0 {
		t.Fatal("Attributes map is empty, expected at least one key-value pair")
	}

	// 可选：打印输出以便调试
	t.Logf("Generated JSON: %s", jsonOutput)
}

func TestRngSource(t *testing.T) {
	rng := random.NewRand()

	// 测试 Seed 方法
	rng.Seed(42) // 这里我们不验证任何状态，因为 Seed 方法是空的

	// 测试 Uint64 方法
	nims := make(map[uint64]struct{})
	for i := 0; i < 100; i++ { // 多次调用以增加不同结果的可能性
		num := rng.Uint64()
		nims[num] = struct{}{}
	}

	// 验证生成的随机数是否在 uint64 的范围内
	for num := range nims {
		assert.LessOrEqual(t, num, ^uint64(0), "Generated number out of uint64 range")
	}

	// 验证至少生成了两个不同的随机数
	if len(nims) < 2 {
		t.Error("Expected at least two different random numbers on multiple calls")
	}
}

const testTimeout = 5 * time.Second

func TestGenerateAvailablePort_DefaultRange(t *testing.T) {
	done := make(chan bool, 1)
	go func() {
		port, err := random.GenerateAvailablePort()
		assert.NoError(t, err, "Failed to generate an available port")

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		assert.NoError(t, err, fmt.Sprintf("Failed to bind to port %d", port))
		listener.Close()

		done <- true
	}()

	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Error("Test timed out waiting for an available port")
	}
}

func TestGenerateAvailablePort_CustomRange(t *testing.T) {
	done := make(chan bool, 1)
	go func() {
		port, err := random.GenerateAvailablePort(2000, 3000)
		assert.NoError(t, err, "Failed to generate an available port within the custom range")
		assert.True(t, port >= 2000 && port <= 3000, fmt.Sprintf("Port %d is not within the range [2000, 3000]", port))

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		assert.NoError(t, err, fmt.Sprintf("Failed to bind to port %d", port))
		listener.Close()

		done <- true
	}()

	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Error("Test timed out waiting for an available port within the custom range")
	}
}

func TestGenerateAvailablePort_InvalidRange(t *testing.T) {
	_, err := random.GenerateAvailablePort(65536, 1024)
	assert.Error(t, err, "Expected an error for an invalid port range")
}

func TestRandNumericalStep(t *testing.T) {
	// 指定步长为2，整数类型
	got := random.RandNumerical(2, 10, 2)
	want := []int{2, 4, 6, 8, 10}
	assert.Equal(t, want, got)

	// 指定步长为3，整数类型，end未整除步长
	got = random.RandNumerical(1, 10, 3)
	want = []int{1, 4, 7, 10}
	assert.Equal(t, want, got)

	// 步长大于区间长度，结果只有一个元素
	got = random.RandNumerical(1, 3, 5)
	want = []int{1}
	assert.Equal(t, want, got)

	// 浮点数类型，指定步长1.5
	gotF := random.RandNumerical(0.0, 6.0, 1.5)
	wantF := []float64{0.0, 1.5, 3.0, 4.5, 6.0}
	assert.Equal(t, wantF, gotF)

	// 步长为0，返回空切片
	got = random.RandNumerical(1, 10, 0)
	assert.Empty(t, got)

	// 负步长，返回空切片（根据你函数逻辑）
	got = random.RandNumerical(1, 10, -1)
	assert.Empty(t, got)
}

func TestRandNumericalInt(t *testing.T) {
	got := random.RandNumerical(3, 7)
	want := []int{3, 4, 5, 6, 7}
	assert.Equal(t, want, got)

	got = random.RandNumerical(5, 3)
	assert.Empty(t, got) // end < start，返回空切片

	got = random.RandNumerical(0, 0)
	want = []int{0}
	assert.Equal(t, want, got)
}

func TestRandNumericalUint8(t *testing.T) {
	got := random.RandNumerical[uint8](1, 5)
	want := []uint8{1, 2, 3, 4, 5}
	assert.Equal(t, want, got)
}

func TestRandNumericalFloat64(t *testing.T) {
	got := random.RandNumerical(1.0, 2.0, 0.3)
	want := []float64{1.0, 1.3, 1.6, 1.9}
	assert.InDeltaSlice(t, want, got, 1e-9) // 浮点数允许误差

	got = random.RandNumerical(2.0, 1.0, 0.1)
	assert.Empty(t, got) // end < start，空切片

	got = random.RandNumerical[float64](0, 1, 0)
	assert.Empty(t, got) // 步长0，空切片
}

func TestRandNumericalFloat32(t *testing.T) {
	got := random.RandNumerical[float32](0, 1, 0.25)
	want := []float32{0, 0.25, 0.5, 0.75, 1}
	assert.InDeltaSlice(t, want, got, 1e-6)
}

func TestRandNumericalWithRandomStepInt(t *testing.T) {
	start, end := 1, 20
	minStep, maxStep := 1, 3

	res := random.RandNumericalWithRandomStep[int](start, end, minStep, maxStep)
	assert.NotEmpty(t, res, "结果切片不应该为空")

	for i, val := range res {
		assert.GreaterOrEqual(t, val, start, "元素 %d 小于 start", i)
		assert.LessOrEqual(t, val, end, "元素 %d 大于 end", i)
		if i > 0 {
			step := val - res[i-1]
			assert.GreaterOrEqual(t, step, minStep, "步长 %d 小于最小步长", step)
			assert.LessOrEqual(t, step, maxStep, "步长 %d 大于最大步长", step)
		}
	}
}

func TestRandNumericalWithRandomStepFloat64(t *testing.T) {
	start, end := 0.0, 2.0
	minStep, maxStep := 0.1, 0.5

	res := random.RandNumericalWithRandomStep[float64](start, end, minStep, maxStep)
	assert.NotEmpty(t, res, "结果切片不应该为空")

	const epsilon = 1e-9
	for i, val := range res {
		assert.True(t, val >= start-epsilon, "元素 %d 小于 start", i)
		assert.True(t, val <= end+epsilon, "元素 %d 大于 end", i)
		if i > 0 {
			step := val - res[i-1]
			assert.True(t, step >= minStep-epsilon, "步长 %v 小于最小步长", step)
			assert.True(t, step <= maxStep+epsilon, "步长 %v 大于最大步长", step)
		}
	}
}
