/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:15:15
 * @FilePath: \go-toolbox\pkg\mathx\base.go
 * @Description: 计算给定概率分布（map）的熵值
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package mathx

import (
	"math"
)

const epsilon = 1e-6 // 用于避免计算中的概率为0导致的数学错误

// CalcEntropy 计算给定概率分布（map）的熵值
// 熵值表示信息的不确定性或随机性
func CalcEntropy(m map[interface{}]int) float64 {
	// 如果map为空或只有一个元素，熵值为0（这里原函数返回1可能是个误解，因为0个或1个元素的集合没有不确定性）
	if len(m) == 0 {
		return 0
	}

	var entropy float64
	total := 0 // 使用单个变量替代var total int来减少内存分配

	// 计算所有值的总和
	for _, count := range m {
		total += count
	}

	// 遍历map，计算每个元素的概率和对应的熵贡献
	for _, count := range m {
		// 使用float64转换来确保精度，并计算概率
		prob := float64(count) / float64(total)
		// 如果概率非常小，则替换为epsilon以避免log(0)的错误
		if prob < epsilon {
			prob = epsilon
		}
		// 根据熵的公式计算贡献：-p * log2(p)
		entropy -= prob * math.Log2(prob)
	}

	// 熵值应除以log2(元素数量)来归一化，但这里我们直接返回未归一化的熵值，
	// 因为归一化通常取决于具体应用场景。如果确实需要归一化，可以取消下面这行的注释。
	// 注意：原函数中的归一化方式可能是有误的，因为它将熵除以了log2(元素数量)，
	// 这并不符合熵的常规定义。通常，熵是信息不确定性的度量，不需要进一步归一化。
	// return entropy / math.Log2(float64(len(m)))

	// 返回计算得到的熵值
	return entropy
}
