/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 10:28:33
 * @FilePath: \go-toolbox\pkg\random\rand.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

import (
	"math"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/json"
)

// Implement Source and Source64 interfaces
type rngSource struct {
	p sync.Pool
}

func (r *rngSource) Int63() (n int64) {
	src := r.p.Get()
	n = src.(rand.Source).Int63()
	r.p.Put(src)
	return
}

// Seed specify seed when using NewRand()
func (r *rngSource) Seed(_ int64) {}

func (r *rngSource) Uint64() (n uint64) {
	src := r.p.Get()
	n = src.(rand.Source64).Uint64()
	r.p.Put(src)
	return
}

// NewRand goroutine-safe rand.Rand, optional seed value
func NewRand(seed ...int64) *rand.Rand {
	n := time.Now().UnixNano()
	if len(seed) > 0 {
		n = seed[0]
	}
	src := &rngSource{
		p: sync.Pool{
			New: func() interface{} {
				return rand.NewSource(n)
			},
		},
	}
	return rand.New(src)
}

var (
	// 设置随机种子
	mathRandSend = rand.New(rand.NewSource(time.Now().Unix()))
	// 大写字母
	matchCapital *[]int
	// 小写字母
	matchLowercase *[]int
	// 特殊符号
	matchSpecial *[]int
	// 数字
	matchNumber *[]int

	once sync.Once

	newRandSend = NewRand()
)

// RandInt
/**
 *  @Description: 随机整数
 *  @param start
 *  @param end
 *  @return v
 */
func RandInt(min, max int) (v int) {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return mathRandSend.Intn(max-min) + min
}

// RandFloat
/**
 *  @Description: 随机小数
 *  @param min
 *  @param max
 *  @return v
 */
func RandFloat(min, max float64) (v float64) {
	return min + mathRandSend.Float64()*(max-min)
}

// initASCII
/**
 *  @Description: 初始化ASCII码列表
 */
func initASCII() {
	once.Do(func() {
		// 大写字母
		c := make([]int, 26)
		for i := 0; i < 26; i++ {
			c[i] = 65 + i
		}
		// 小写字母
		matchCapital = &c
		l := make([]int, 26)
		for i := 0; i < 26; i++ {
			l[i] = 97 + i
		}
		matchLowercase = &l
		// 数字
		n := make([]int, 10)
		for i := 0; i < 10; i++ {
			n[i] = 48 + i
		}
		matchNumber = &n
		// 特殊字符(. @$!%*#_~?&^)
		s := []int{46, 64, 36, 33, 37, 42, 35, 95, 126, 63, 38, 94}
		matchSpecial = &s
	})
}

// RandString
/**
 *  @Description: 随机生成字符串
 *  @param n 字符串长度
 *  @param mode 字符串模式 random.NUMBER|random.LOWERCASE|random.SPECIAL|random.CAPITAL)
 *  @return str 生成的字符串
 */
func RandString(n int, mode RandType) (str string) {
	initASCII()
	var ascii []int
	if mode&CAPITAL >= CAPITAL {
		ascii = append(ascii, *matchCapital...)
	}
	if mode&LOWERCASE >= LOWERCASE {
		ascii = append(ascii, *matchLowercase...)
	}
	if mode&SPECIAL >= SPECIAL {
		ascii = append(ascii, *matchSpecial...)
	}
	if mode&NUMBER >= NUMBER {
		ascii = append(ascii, *matchNumber...)
	}
	if len(ascii) == 0 {
		return
	}
	var build strings.Builder
	for i := 0; i < n; i++ {
		build.WriteString(string(rune(ascii[mathRandSend.Intn(len(ascii))])))
	}
	str = build.String()
	return
}

// RandStringSlice 指定长度的随机字符串
func RandStringSlice(count, len int, mode RandType) (result []string) {
	for i := 0; i < count; i++ {
		result = append(result, RandString(len, mode))
	}
	return result
}

// RandomStr 随机一个字符串
func RandomStr(length int, customBytes ...string) string {
	var sb strings.Builder
	randBytes := ALPHA_BYTES
	if len(customBytes) > 0 {
		randBytes = customBytes[0]
	}
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(randBytes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(randBytes))]))
		}
	}
	return sb.String()
}

// RandomNumber 随机一个数字字符串
func RandomNumber(length int, customBytes ...string) string {
	var sb strings.Builder
	randBytes := DEC_BYTES
	if len(customBytes) > 0 {
		randBytes = customBytes[0]
	}
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(randBytes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(randBytes))]))
		}
	}
	return sb.String()
}

// RandomHex 随机一个hex字符串
func RandomHex(bytesLen int, customBytes ...string) string {
	var sb strings.Builder
	randBytes := HEX_BYTES
	if len(customBytes) > 0 {
		randBytes = customBytes[0]
	}
	for i := 0; i < bytesLen<<1; i++ {
		sb.WriteString(string(randBytes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(randBytes))]))
	}
	return sb.String()
}

// FRandInt (>=)min - (<)max
func FRandInt(min, max int) int {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return FastIntn(max-min) + min
}

// FRandUint32 (>=)min - (<)max
func FRandUint32(min, max uint32) uint32 {
	if max == min {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return FastRandn(max-min) + min
}

// FastIntn this is similar to rand.Intn, but faster.
// A non-negative pseudo-random number in the half-open interval [0,n).
// Return 0 if n <= 0.
func FastIntn(n int) int {
	if n <= 0 {
		return 0
	}
	if n <= math.MaxUint32 {
		return int(FastRandn(uint32(n)))
	}
	return int(newRandSend.Int63n(int64(n)))
}

// FRandString a random string, which may contain uppercase letters, lowercase letters and numbers.
// Ref: stackoverflow.icza
func FRandString(n int) string {
	return convert.B2S(FRandBytes(n))
}

// FRandHexString 指定长度的随机 hex 字符串
func FRandHexString(n int) string {
	return convert.B2S(FRandHexBytes(n))
}

// FRandAlphaString 指定长度的随机字母字符串
func FRandAlphaString(n int) string {
	return convert.B2S(FRandAlphaBytes(n))
}

// FRandDecString 指定长度的随机数字字符串
func FRandDecString(n int) string {
	return convert.B2S(FRandDecBytes(n))
}

// FRandBytes random bytes, but faster.
func FRandBytes(n int) []byte {
	return FRandBytesLetters(n, LETTER_BYTES)
}

// FRandAlphaBytes generates random alpha bytes.
func FRandAlphaBytes(n int) []byte {
	return FRandBytesLetters(n, ALPHA_BYTES)
}

// FRandHexBytes generates random hexadecimal bytes.
func FRandHexBytes(n int) []byte {
	return FRandBytesLetters(n, HEX_BYTES)
}

// RandDecBytes 指定长度的随机数字切片
func FRandDecBytes(n int) []byte {
	return FRandBytesLetters(n, DEC_BYTES)
}

// FRandBytesLetters 生成指定长度的字符切片
func FRandBytesLetters(n int, letters string) []byte {
	if n < 1 || len(letters) < 2 {
		return nil
	}
	b := make([]byte, n)
	for i, cache, remain := n-1, FastRand(), LETTER_IDX_MAX; i >= 0; {
		if remain == 0 {
			cache, remain = FastRand(), LETTER_IDX_MAX
		}
		if idx := int(cache & LETTER_IDX_MASK); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= LETTER_IDX_BITS
		remain--
	}
	return b
}

// FRandBytesJSON 生成指定长度的随机字节字符串，并返回 JSON 格式
func FRandBytesJSON(length int) (string, error) {
	// 生成随机字节
	randomBytes := FRandBytes(length)

	// 将字节转换为 JSON 格式
	jsonData, err := json.Marshal(randomBytes)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// 生成随机布尔值
func FRandBool() bool {
	return FRandInt(1, 2) == 1
}

// 生成随机时间
func FRandTime() time.Time {
	return time.Now().Add(time.Duration(FRandInt(1, 1000)) * time.Hour)
}

// GenerateRandomModel 生成随机模型的 JSON 格式
func GenerateRandomModel(model interface{}) (interface{}, string, error) {
	v := reflect.ValueOf(model)

	// 确保传入的是指针类型且非空
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, "", nil
	}

	v = v.Elem() // 获取指针指向的值

	// 填充模型字段的随机值
	if err := populateFields(v); err != nil {
		return nil, "", err
	}

	// 将模型转换为 JSON 格式
	jsonData, err := convert.MustJSONIndent(model)
	if err != nil {
		return nil, "", err
	}
	return model, string(jsonData), nil
}

// populateFields 填充结构体字段的随机值
func populateFields(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		// 仅处理导出字段
		if field.CanSet() {
			// 根据字段类型设置随机值
			if err := setRandomValue(field, fieldType); err != nil {
				return err
			}
		}
	}
	return nil
}

// setRandomValue 根据字段类型设置随机值
func setRandomValue(field reflect.Value, fieldType reflect.StructField) error {
	switch fieldType.Type.Kind() {
	case reflect.String:
		field.SetString(FRandString(10)) // 随机生成10个字符的字符串
	case reflect.Int:
		field.SetInt(int64(FRandInt(18, 65))) // 随机生成18到65之间的整数
	case reflect.Float64:
		field.SetFloat(math.Round(float64(FRandInt(1, 100)) / 1.5)) // 随机生成1.0到100.0之间的浮点数
	case reflect.Bool:
		field.SetBool(FRandBool()) // 随机生成布尔值
	case reflect.Struct:
		if fieldType.Type == reflect.TypeOf(time.Time{}) {
			field.Set(reflect.ValueOf(FRandTime())) // 随机生成时间
		} else {
			// 递归填充嵌套结构体字段
			if err := populateFields(field); err != nil {
				return err
			}
		}
	case reflect.Slice:
		return setRandomSlice(field, fieldType) // 处理切片类型
	case reflect.Map:
		return setRandomMap(field, fieldType) // 处理映射类型
	case reflect.Ptr: // 处理指针类型
		if field.IsNil() {
			field.Set(reflect.New(fieldType.Type.Elem())) // 确保指针被分配
		}
		// 仅当指针指向结构体时才递归填充
		if fieldType.Type.Elem().Kind() == reflect.Struct {
			return populateFields(field.Elem())
		}
	default:
		// 你可以添加更多类型的处理逻辑
	}
	return nil
}

// setRandomSlice 随机生成切片并设置到字段
func setRandomSlice(field reflect.Value, fieldType reflect.StructField) error {
	if fieldType.Type.Elem().Kind() == reflect.String {
		length := FRandInt(1, 5) // 随机长度
		slice := reflect.MakeSlice(fieldType.Type, length, length)
		for j := 0; j < length; j++ {
			slice.Index(j).SetString(FRandString(5)) // 随机生成5个字符的字符串
		}
		field.Set(slice) // 设置生成的切片
	} else if fieldType.Type.Elem().Kind() == reflect.Struct {
		length := FRandInt(1, 5) // 随机长度
		slice := reflect.MakeSlice(fieldType.Type, length, length)
		for j := 0; j < length; j++ {
			// 递归填充每个嵌套结构体
			if err := populateFields(slice.Index(j)); err != nil {
				return err
			}
		}
		field.Set(slice) // 设置生成的切片
	}
	return nil
}

// setRandomMap 随机生成映射并设置到字段
func setRandomMap(field reflect.Value, fieldType reflect.StructField) error {
	if fieldType.Type.Key().Kind() == reflect.String && fieldType.Type.Elem().Kind() == reflect.Int {
		m := reflect.MakeMap(fieldType.Type) // 创建映射
		length := FRandInt(1, 5)             // 随机长度
		for j := 0; j < length; j++ {
			key := FRandString(5)                                       // 随机生成字符串作为键
			value := FRandInt(1, 100)                                   // 随机生成整数作为值
			m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value)) // 设置映射的键值对
		}
		field.Set(m) // 设置生成的映射
	}
	return nil
}
