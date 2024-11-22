/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 12:57:21
 * @FilePath: \go-toolbox\tests\stringx_base_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/stretchr/testify/assert"
)

func TestLengthChain(t *testing.T) {
	result := stringx.New("Hello, World!").LengthChain()
	assert.Equal(t, 13, result)
}

func TestReverseChain(t *testing.T) {
	result := stringx.New("Hello, World!").ReverseChain().Value()
	assert.Equal(t, "!dlroW ,olleH", result)
}

func TestEqualsChain(t *testing.T) {
	result := stringx.New("hello").EqualsChain("hello")
	assert.True(t, result)
}

func TestEqualsIgnoreCaseChain(t *testing.T) {
	result := stringx.New("HELLO").EqualsIgnoreCaseChain("hello")
	assert.True(t, result)
}

func TestInsertSpacesChain(t *testing.T) {
	result := stringx.New("1234567890").InsertSpacesChain(2).Value()
	assert.Equal(t, "12 34 56 78 90", result)
}

func TestEqualsAnyChain(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := stringx.New("banana").EqualsAnyChain(strList)
	assert.True(t, result)
}

func TestEqualsAnyIgnoreCaseChain(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := stringx.New("OrAnGe").EqualsAnyIgnoreCaseChain(strList)
	assert.True(t, result)
}

func TestEqualsAtChain(t *testing.T) {
	result := stringx.New("hello").EqualsAtChain(1, "e")
	assert.True(t, result)
}

func TestCountChain(t *testing.T) {
	result := stringx.New("banana").CountChain("a")
	assert.Equal(t, 3, result)
}

func TestCompareIgnoreCaseChain(t *testing.T) {
	result := stringx.New("apple").CompareIgnoreCaseChain("BANANA")
	assert.Less(t, result, 0)
}

func TestCoalesceChain(t *testing.T) {
	result := stringx.New("Hello").CoalesceChain(" ", "World", "!")
	assert.Equal(t, "Hello World!", result.Value())
}

func TestConvertCharacterStyleChain(t *testing.T) {
	tests := []struct {
		input    string
		style    stringx.CharacterStyle
		expected string
	}{
		{"HelloWorld", stringx.SnakeCharacterStyle, "hello_world"},
		{"helloWorld", stringx.SnakeCharacterStyle, "hello_world"},
		{"Hello_World", stringx.SnakeCharacterStyle, "hello_world"},
		{" Hello World", stringx.SnakeCharacterStyle, "hello_world"},
		{"Hello World", stringx.SnakeCharacterStyle, "hello_world"},
		{" ", stringx.SnakeCharacterStyle, ""},
		{"", stringx.SnakeCharacterStyle, ""},

		{"hello_world", stringx.StudlyCharacterStyle, "HelloWorld"},
		{"helloWorld", stringx.StudlyCharacterStyle, "HelloWorld"},
		{"hello world", stringx.StudlyCharacterStyle, "HelloWorld"},
		{" Hello World", stringx.StudlyCharacterStyle, "HelloWorld"},
		{"Hello_World", stringx.StudlyCharacterStyle, "HelloWorld"},
		{" ", stringx.StudlyCharacterStyle, ""},
		{"", stringx.StudlyCharacterStyle, ""},

		{"hello_world", stringx.CamelCharacterStyle, "helloWorld"},
		{"HelloWorld", stringx.CamelCharacterStyle, "helloWorld"},
		{"hello world", stringx.CamelCharacterStyle, "helloWorld"},
		{" Hello World", stringx.CamelCharacterStyle, "helloWorld"},
		{"Hello_World", stringx.CamelCharacterStyle, "helloWorld"},
		{" ", stringx.CamelCharacterStyle, ""},
		{"", stringx.CamelCharacterStyle, ""},

		{"HelloWorld", stringx.CharacterStyle(999), "HelloWorld"},
	}

	for _, test := range tests {
		result := stringx.New(test.input).ConvertCharacterStyleChain(test.style).Value()
		assert.Equal(t, test.expected, result, "ConvertCharacterStyle(%q, %v) = %q; want %q", test.input, test.style, result, test.expected)
	}
}

func TestToLowerChain(t *testing.T) {
	s := stringx.New("Hello World")
	result := s.ToLowerChain().Value()
	expected := "hello world"
	assert.Equal(t, expected, result, "ToLowerChain() = %v; want %v", result, expected)
}

func TestToUpperChain(t *testing.T) {
	s := stringx.New("hello world")
	result := s.ToUpperChain().Value()
	expected := "HELLO WORLD"
	assert.Equal(t, expected, result, "ToUpperChain() = %v; want %v", result, expected)
}

func TestToTitleChain(t *testing.T) {
	s := stringx.New("hello world")
	result := s.ToTitleChain().Value()
	expected := "Hello World"
	assert.Equal(t, expected, result, "ToTitleChain() = %v; want %v", result, expected)
}

func TestChainedMethods(t *testing.T) {
	s := stringx.New("gO LaNg")
	result := s.ToLowerChain().ToUpperChain().ToTitleChain().Value()
	expected := "Go Lang"
	assert.Equal(t, expected, result, "Chained methods = %v; want %v", result, expected)
}
