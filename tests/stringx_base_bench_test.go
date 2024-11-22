/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 12:55:59
 * @FilePath: \go-toolbox\tests\stringx_base_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

// BenchmarkCoalesce 性能基准测试
func BenchmarkCoalesce(b *testing.B) {
	input := []string{"11", "+", "22", "=", "33"}
	for i := 0; i < b.N; i++ {
		stringx.Coalesce(input...)
	}
}

func BenchmarkToSnakeCase(b *testing.B) {
	input := "ThisIsATestStringForBenchmarking"
	for i := 0; i < b.N; i++ {
		stringx.ConvertCharacterStyle(input, stringx.SnakeCharacterStyle)
	}
}

func BenchmarkToStudlyCase(b *testing.B) {
	input := "this_is_a_test_string_for_benchmarking"
	for i := 0; i < b.N; i++ {
		stringx.ConvertCharacterStyle(input, stringx.StudlyCharacterStyle)
	}
}

func BenchmarkToCamelCase(b *testing.B) {
	input := "this_is_a_test_string_for_benchmarking"
	for i := 0; i < b.N; i++ {
		stringx.ConvertCharacterStyle(input, stringx.CamelCharacterStyle)
	}
}

func BenchmarkToLowerChain(b *testing.B) {
	s := stringx.New("Hello World")
	for i := 0; i < b.N; i++ {
		_ = s.ToLowerChain().Value()
	}
}

func BenchmarkToUpperChain(b *testing.B) {
	s := stringx.New("hello world")
	for i := 0; i < b.N; i++ {
		_ = s.ToUpperChain().Value()
	}
}

func BenchmarkToTitleChain(b *testing.B) {
	s := stringx.New("hello world")
	for i := 0; i < b.N; i++ {
		_ = s.ToTitleChain().Value()
	}
}

func BenchmarkChainedMethods(b *testing.B) {
	s := stringx.New("gO LaNg")
	for i := 0; i < b.N; i++ {
		_ = s.ToLowerChain().ToUpperChain().ToTitleChain().Value()
	}
}
