/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-12 15:27:26
 * @FilePath: \go-toolbox\pkg\random\rand.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/json"
	"github.com/kamalyes/go-toolbox/pkg/types"
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
	n := time.Now().UnixMicro()
	n += atomic.AddInt64(&randCounter, 1)
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
	// 大写字母
	matchCapital *[]int
	// 小写字母
	matchLowercase *[]int
	// 特殊符号
	matchSpecial *[]int
	// 数字
	matchNumber *[]int

	once sync.Once

	randCounter int64
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
	return NewRand().Intn(max-min) + min
}

// RandFloat
/**
 *  @Description: 随机小数
 *  @param min
 *  @param max
 *  @return v
 */
func RandFloat(min, max float64) (v float64) {
	return min + NewRand().Float64()*(max-min)
}

// initASCII 初始化ASCII码列表
func initASCII() {
	once.Do(func() {
		matchCapital = createASCIIList(65, 90)                                 // 大写字母
		matchLowercase = createASCIIList(97, 122)                              // 小写字母
		matchNumber = createASCIIList(48, 57)                                  // 数字
		matchSpecial = &[]int{46, 64, 36, 33, 37, 42, 35, 95, 126, 63, 38, 94} // 特殊字符 (. @$!%*#_~?&^)
	})
}

// createASCIIList 创建指定范围的ASCII码列表
func createASCIIList(start, end int) *[]int {
	list := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		list[i-start] = i
	}
	return &list
}

// RandString
/**
 *  @Description: 随机生成字符串
 *  @param n 字符串长度
 *  @param mode 字符串模式 random.NUMBER|random.LOWERCASE|random.SPECIAL|random.CAPITAL)
 *  @return str 生成的字符串
 */
func RandString(n int, mode RandType) string {
	var (
		build strings.Builder // 用于高效构建字符串
		ascii []int           // 用于存储所有符合mode的字符的ASCII码集合
		r     = NewRand()     // 创建一个带时间种子的随机数生成器
	)

	// 判断mode中是否包含大写字母字符集，如果包含则追加对应ASCII码
	if mode&CAPITAL != 0 {
		ascii = append(ascii, *matchCapital...)
	}
	// 判断mode中是否包含小写字母字符集，如果包含则追加对应ASCII码
	if mode&LOWERCASE != 0 {
		ascii = append(ascii, *matchLowercase...)
	}
	// 判断mode中是否包含特殊字符集，如果包含则追加对应ASCII码
	if mode&SPECIAL != 0 {
		ascii = append(ascii, *matchSpecial...)
	}
	// 判断mode中是否包含数字字符集，如果包含则追加对应ASCII码
	if mode&NUMBER != 0 {
		ascii = append(ascii, *matchNumber...)
	}

	// 如果没有任何字符集被选中，ascii长度为0，返回空字符串
	if len(ascii) == 0 {
		return ""
	}

	// 循环n次，每次随机选择一个ascii码对应的字符追加到字符串中
	for i := 0; i < n; i++ {
		randomIndex := r.Intn(len(ascii))         // 在ascii切片中随机选一个索引
		build.WriteRune(rune(ascii[randomIndex])) // 将对应ASCII码转换成字符写入builder
	}

	return build.String() // 返回最终生成的随机字符串
}

// RandStringSlice 指定长度的随机字符串
func RandStringSlice(count, len int, mode RandType) (result []string) {
	for i := 0; i < count; i++ {
		result = append(result, RandString(len, mode))
	}
	return result
}

var defaultSliceSize = 1000

// RandNumericalLargeSlice 随机生成大数据整数切片
func RandNumericalLargeSlice[T types.Numerical](largeSize ...int) []T {
	if len(largeSize) > 0 {
		defaultSliceSize = largeSize[0]
	}
	slice := make([]T, defaultSliceSize)
	for i := 0; i < defaultSliceSize; i++ {
		slice[i] = T(i % 100) // 重复一些值以测试去重和重复检查
	}
	return slice
}

// RandNumber 随机一个数字字符串
func RandNumber(length int, customBytes ...string) string {
	var sb strings.Builder
	randBytes := DEC_BYTES
	if len(customBytes) > 0 {
		randBytes = customBytes[0]
	}
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(randBytes[NewRand().Intn(len(randBytes))]))
		}
	}
	return sb.String()
}

// RandHex 随机一个hex字符串
func RandHex(bytesLen int, customBytes ...string) string {
	var sb strings.Builder
	randBytes := HEX_BYTES
	if len(customBytes) > 0 {
		randBytes = customBytes[0]
	}
	for i := 0; i < bytesLen<<1; i++ {
		sb.WriteString(string(randBytes[NewRand().Intn(len(randBytes))]))
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
	return int(NewRand().Int63n(int64(n)))
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

// GenerateRandModel 生成随机模型的 JSON 格式
func GenerateRandModel(model interface{}) (interface{}, string, error) {
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
			if err := setRandValue(field, fieldType); err != nil {
				return err
			}
		}
	}
	return nil
}

// setRandValue 根据字段类型设置随机值
func setRandValue(field reflect.Value, fieldType reflect.StructField) error {
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
		return setRandSlice(field, fieldType) // 处理切片类型
	case reflect.Map:
		return setRandMap(field, fieldType) // 处理映射类型
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

// setRandSlice 随机生成切片并设置到字段
func setRandSlice(field reflect.Value, fieldType reflect.StructField) error {
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

// setRandMap 随机生成映射并设置到字段
func setRandMap(field reflect.Value, fieldType reflect.StructField) error {
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

// GenerateAvailablePort 返回一个随机的可用端口号
// 如果未提供ports参数或提供的参数长度不为2，则使用默认端口范围1024到65535
func GenerateAvailablePort(ports ...int) (int, error) {
	// 设置默认端口范围
	minPort, maxPort := 1024, 65535

	// 检查是否提供了有效的ports参数
	if len(ports) == 2 {
		minPort = ports[0]
		maxPort = ports[1]
		// 验证提供的端口范围是否有效
		if minPort > maxPort {
			return 0, fmt.Errorf("minimum port cannot be greater than maximum port")
		}
	}

	for {
		port := RandInt(minPort, maxPort)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			// 成功监听端口，关闭监听器并返回端口号
			listener.Close()
			return port, nil
		}
		// 如果端口已被使用或发生其他错误，则尝试下一个端口
	}
}

// RandNumerical 泛型函数，生成从 start 到 end（包含）的切片
func RandNumerical[T types.Numerical](start, end T, step ...T) []T {
	// 如果结束值小于开始值，返回 nil 表示无效输入
	if end < start {
		return nil
	}

	// 默认步长为 1，如果传入了步长参数，则使用传入的步长
	var stepVal T = 1
	if len(step) > 0 {
		stepVal = step[0]
	}

	// 步长必须大于 0，否则返回 nil 表示无效步长
	if stepVal <= 0 {
		return nil
	}

	// 根据 start 的具体类型区分处理浮点数和整数
	switch any(start).(type) {
	// 如果是浮点数类型（float32 或 float64）
	case float32, float64:
		// 计算序列长度 size，向下取整后加 1，确保包含起点
		size := int((end-start)/stepVal) + 1
		if size <= 0 {
			return nil
		}

		// 创建结果切片，容量为 size
		result := make([]T, size)

		// 将步长和起点转换为 float64 方便计算
		stepF := float64(stepVal)
		startF := float64(start)

		// 依次计算每个元素的值：start + i * step
		for i := 0; i < size; i++ {
			result[i] = T(startF + float64(i)*stepF) // 直接转换为泛型浮点数类型 T
		}

		return result

	// 其他类型视为整数类型处理
	default:
		// 将整数类型参数转换为 float64，方便做除法和乘法计算
		startF := float64(start)
		endF := float64(end)
		stepF := float64(stepVal)

		// 计算序列长度 size，向下取整后加 1，确保包含起点
		size := int((endF-startF)/stepF) + 1
		if size <= 0 {
			return nil
		}

		// 创建结果切片，容量为 size
		result := make([]T, size)

		// 依次计算每个元素的值（浮点数计算）
		for i := 0; i < size; i++ {
			valF := startF + float64(i)*stepF

			// 利用 Float64ToInt 函数将浮点数安全转换回整数泛型 T
			// RoundNone 表示不做额外取整，直接转换
			val, err := convert.Float64ToInt[T](valF, convert.RoundNone)
			if err != nil {
				// 转换失败时 panic，提示错误信息
				panic(fmt.Sprintf("RandNumerical: 转换失败 %v", err))
			}
			result[i] = val
		}

		return result
	}
}

// RandNumericalWithRandomStep 泛型函数，生成从 start 到 end（包含或不超过）
// 每一步随机步长在 [minStep, maxStep] 范围内（浮点数或整数）
func RandNumericalWithRandomStep[T types.Numerical](start, end, minStep, maxStep T) []T {
	if end < start {
		return nil
	}
	if minStep <= 0 || maxStep < minStep {
		return nil
	}

	var result []T
	current := start

	switch any(start).(type) {
	case float32, float64:
		for current <= end {
			result = append(result, current)
			// 生成随机浮点步长
			step := RandFloat(float64(minStep), float64(maxStep))
			current += T(step)
		}
	default:
		for current <= end {
			result = append(result, current)
			// 生成随机整数步长
			step := RandInt(int(minStep), int(maxStep))
			current += T(step)
		}
	}

	return result
}

func init() {
	// 初始化ASCII码列表，只执行一次
	initASCII()
}
