/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-16 22:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-16 22:55:55
 * @FilePath: \go-toolbox\pkg\safe\mathx.go
 * @Description: 安全的数学工具函数集合，提供各种数学算法的安全实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package safe

import (
	"errors"
	"math"
	"math/big"
)

// 定义常用的数学常数
const (
	MaxSafeInteger = 1<<53 - 1      // JavaScript安全整数最大值
	MinSafeInteger = -(1<<53 - 1)   // JavaScript安全整数最小值
	GoldenRatio    = 1.618033988749 // 黄金比例
	EulerNumber    = 2.718281828459 // 欧拉数
)

// FastHash 快速哈希函数 - 安全版本
// 使用FNV-1a算法，避免哈希冲突和安全问题
func FastHash(s string) uint64 {
	if len(s) == 0 {
		return 1 // 返回非0值避免哈希冲突
	}

	// 使用更安全的字符串访问方式
	var h uint64 = 14695981039346656037 // FNV offset basis
	b := []byte(s)                      // 转换为byte切片，更安全
	for _, v := range b {
		h ^= uint64(v)
		h *= 1099511628211 // FNV prime
	}
	return h
}

// NextPowerOfTwo 获取下一个2的幂 - 安全版本
// 防止整数溢出和无效输入
func NextPowerOfTwo(n int) int {
	if n <= 1 {
		return 2
	}
	if n > MaxSafeInteger>>1 {
		return MaxSafeInteger // 防止溢出
	}
	
	// 如果n已经是2的幂，返回下一个
	if n&(n-1) == 0 {
		return n << 1
	}

	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n + 1
}

// SafeAdd 安全整数加法，防止溢出
func SafeAdd(a, b int64) (int64, error) {
	if b > 0 && a > math.MaxInt64-b {
		return 0, errors.New("integer overflow in addition")
	}
	if b < 0 && a < math.MinInt64-b {
		return 0, errors.New("integer underflow in addition")
	}
	return a + b, nil
}

// SafeSubtract 安全整数减法，防止溢出
func SafeSubtract(a, b int64) (int64, error) {
	if b > 0 && a < math.MinInt64+b {
		return 0, errors.New("integer underflow in subtraction")
	}
	if b < 0 && a > math.MaxInt64+b {
		return 0, errors.New("integer overflow in subtraction")
	}
	return a - b, nil
}

// SafeMultiply 安全整数乘法，防止溢出
func SafeMultiply(a, b int64) (int64, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}

	// 检查溢出
	if a > 0 && b > 0 && a > math.MaxInt64/b {
		return 0, errors.New("integer overflow in multiplication")
	}
	if a > 0 && b < 0 && b < math.MinInt64/a {
		return 0, errors.New("integer underflow in multiplication")
	}
	if a < 0 && b > 0 && a < math.MinInt64/b {
		return 0, errors.New("integer underflow in multiplication")
	}
	if a < 0 && b < 0 && a < math.MaxInt64/b {
		return 0, errors.New("integer overflow in multiplication")
	}

	return a * b, nil
}

// SafeDivide 安全整数除法，防止除零和溢出
func SafeDivide(a, b int64) (int64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	if a == math.MinInt64 && b == -1 {
		return 0, errors.New("integer overflow in division")
	}
	return a / b, nil
}

// SafeModulo 安全取模运算，防止除零
func SafeModulo(a, b int64) (int64, error) {
	if b == 0 {
		return 0, errors.New("modulo by zero")
	}
	return a % b, nil
}

// SafePower 安全幂运算，使用快速幂算法防止溢出
func SafePower(base, exp int64) (int64, error) {
	if exp < 0 {
		return 0, errors.New("negative exponent not supported for integer power")
	}
	if exp == 0 {
		return 1, nil
	}
	if base == 0 {
		return 0, nil
	}
	if base == 1 {
		return 1, nil
	}
	if base == -1 {
		if exp%2 == 0 {
			return 1, nil
		}
		return -1, nil
	}

	var result int64 = 1
	var currentBase = base

	for exp > 0 {
		if exp%2 == 1 {
			newResult, err := SafeMultiply(result, currentBase)
			if err != nil {
				return 0, err
			}
			result = newResult
		}

		if exp > 1 {
			newBase, err := SafeMultiply(currentBase, currentBase)
			if err != nil {
				return 0, err
			}
			currentBase = newBase
		}
		exp /= 2
	}

	return result, nil
}

// SafeSqrt 安全平方根计算，使用牛顿法
func SafeSqrt(n float64) (float64, error) {
	if n < 0 {
		return 0, errors.New("square root of negative number")
	}
	if n == 0 {
		return 0, nil
	}

	// 使用 math.Sqrt 进行安全计算
	result := math.Sqrt(n)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, errors.New("invalid square root result")
	}

	return result, nil
}

// SafeLog 安全对数计算
func SafeLog(n, base float64) (float64, error) {
	if n <= 0 {
		return 0, errors.New("logarithm of non-positive number")
	}
	if base <= 0 || base == 1 {
		return 0, errors.New("invalid logarithm base")
	}

	result := math.Log(n) / math.Log(base)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, errors.New("invalid logarithm result")
	}

	return result, nil
}

// SafeGCD 安全最大公约数计算，使用欧几里得算法
func SafeGCD(a, b int64) int64 {
	// 处理负数
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}

	// 欧几里得算法
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// SafeLCM 安全最小公倍数计算
func SafeLCM(a, b int64) (int64, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}

	gcd := SafeGCD(a, b)
	if gcd == 0 {
		return 0, errors.New("invalid GCD result")
	}

	// 防止溢出
	absA := a
	if a < 0 {
		absA = -a
	}
	absB := b
	if b < 0 {
		absB = -b
	}

	result, err := SafeMultiply(absA/gcd, absB)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// IsPrime 安全素数检测，使用Miller-Rabin算法
func IsPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	// 对于较小的数，使用简单试除法
	if n < 100 {
		for i := int64(3); i*i <= n; i += 2 {
			if n%i == 0 {
				return false
			}
		}
		return true
	}

	// 对于较大的数，使用big.Int的ProbablyPrime方法
	bigN := big.NewInt(n)
	return bigN.ProbablyPrime(20) // 20次Miller-Rabin测试
}

// Fibonacci 安全斐波那契数列计算，防止溢出
func Fibonacci(n int) (int64, error) {
	if n < 0 {
		return 0, errors.New("negative Fibonacci index not supported")
	}
	if n == 0 {
		return 0, nil
	}
	if n == 1 {
		return 1, nil
	}

	var a, b int64 = 0, 1
	for i := 2; i <= n; i++ {
		next, err := SafeAdd(a, b)
		if err != nil {
			return 0, err
		}
		a, b = b, next
	}

	return b, nil
}

// Factorial 安全阶乘计算，防止溢出
func Factorial(n int) (*big.Int, error) {
	if n < 0 {
		return nil, errors.New("negative factorial not defined")
	}

	result := big.NewInt(1)
	for i := 2; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}

	return result, nil
}

// SafeAverage 安全平均值计算，防止溢出
func SafeAverage(numbers []int64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("cannot calculate average of empty slice")
	}

	var sum int64 = 0
	for _, num := range numbers {
		newSum, err := SafeAdd(sum, num)
		if err != nil {
			// 如果溢出，使用浮点数计算
			floatSum := float64(sum)
			for i := indexOf(numbers, num); i < len(numbers); i++ {
				floatSum += float64(numbers[i])
			}
			return floatSum / float64(len(numbers)), nil
		}
		sum = newSum
	}

	return float64(sum) / float64(len(numbers)), nil
}

// indexOf 辅助函数：查找元素在切片中的索引
func indexOf(slice []int64, item int64) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// SafeMax 安全获取最大值
func SafeMax(numbers []int64) (int64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("cannot find max of empty slice")
	}

	max := numbers[0]
	for _, num := range numbers[1:] {
		if num > max {
			max = num
		}
	}

	return max, nil
}

// SafeMin 安全获取最小值
func SafeMin(numbers []int64) (int64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("cannot find min of empty slice")
	}

	min := numbers[0]
	for _, num := range numbers[1:] {
		if num < min {
			min = num
		}
	}

	return min, nil
}

// SafeClamp 安全值范围限制
func SafeClamp(value, min, max int64) (int64, error) {
	if min > max {
		return 0, errors.New("min cannot be greater than max")
	}

	if value < min {
		return min, nil
	}
	if value > max {
		return max, nil
	}

	return value, nil
}

// SafeAbs 安全绝对值计算，防止溢出
func SafeAbs(n int64) (int64, error) {
	if n == math.MinInt64 {
		return 0, errors.New("absolute value of MinInt64 causes overflow")
	}

	if n < 0 {
		return -n, nil
	}

	return n, nil
}
