/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 12:57:21
 * @FilePath: \go-toolbox\pkg\stringx\base_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLengthChain(t *testing.T) {
	result := New("Hello, World!").LengthChain()
	assert.Equal(t, 13, result)
}

func TestReverseChain(t *testing.T) {
	result := New("Hello, World!").ReverseChain().Value()
	assert.Equal(t, "!dlroW ,olleH", result)
}

func TestEqualsChain(t *testing.T) {
	result := New("hello").EqualsChain("hello")
	assert.True(t, result)
}

func TestEqualsIgnoreCaseChain(t *testing.T) {
	result := New("HELLO").EqualsIgnoreCaseChain("hello")
	assert.True(t, result)
}

func TestInsertSpacesChain(t *testing.T) {
	result := New("1234567890").InsertSpacesChain(2).Value()
	assert.Equal(t, "12 34 56 78 90", result)
}

func TestEqualsAnyChain(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := New("banana").EqualsAnyChain(strList)
	assert.True(t, result)
}

func TestEqualsAnyIgnoreCaseChain(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := New("OrAnGe").EqualsAnyIgnoreCaseChain(strList)
	assert.True(t, result)
}

func TestEqualsAtChain(t *testing.T) {
	result := New("hello").EqualsAtChain(1, "e")
	assert.True(t, result)
}

func TestCountChain(t *testing.T) {
	result := New("banana").CountChain("a")
	assert.Equal(t, 3, result)
}

func TestCompareIgnoreCaseChain(t *testing.T) {
	result := New("apple").CompareIgnoreCaseChain("BANANA")
	assert.Less(t, result, 0)
}

func TestCoalesceChain(t *testing.T) {
	result := New("Hello").CoalesceChain(" ", "World", "!")
	assert.Equal(t, "Hello World!", result.Value())
}

func TestConvertCharacterStyleChain(t *testing.T) {
	tests := []struct {
		input    string
		style    CharacterStyle
		expected string
	}{
		{"HelloWorld", SnakeCharacterStyle, "hello_world"},
		{"helloWorld", SnakeCharacterStyle, "hello_world"},
		{"Hello_World", SnakeCharacterStyle, "hello_world"},
		{" Hello World", SnakeCharacterStyle, "hello_world"},
		{"Hello World", SnakeCharacterStyle, "hello_world"},
		{" ", SnakeCharacterStyle, ""},
		{"", SnakeCharacterStyle, ""},

		{"hello_world", StudlyCharacterStyle, "HelloWorld"},
		{"helloWorld", StudlyCharacterStyle, "HelloWorld"},
		{"hello world", StudlyCharacterStyle, "HelloWorld"},
		{" Hello World", StudlyCharacterStyle, "HelloWorld"},
		{"Hello_World", StudlyCharacterStyle, "HelloWorld"},
		{" ", StudlyCharacterStyle, ""},
		{"", StudlyCharacterStyle, ""},

		{"hello_world", CamelCharacterStyle, "helloWorld"},
		{"HelloWorld", CamelCharacterStyle, "helloWorld"},
		{"hello world", CamelCharacterStyle, "helloWorld"},
		{" Hello World", CamelCharacterStyle, "helloWorld"},
		{"Hello_World", CamelCharacterStyle, "helloWorld"},
		{" ", CamelCharacterStyle, ""},
		{"", CamelCharacterStyle, ""},

		{"HelloWorld", CharacterStyle(999), "HelloWorld"},
	}

	for _, test := range tests {
		result := New(test.input).ConvertCharacterStyleChain(test.style).Value()
		assert.Equal(t, test.expected, result, "ConvertCharacterStyle(%q, %v) = %q; want %q", test.input, test.style, result, test.expected)
	}
}

func TestToLowerChain(t *testing.T) {
	s := New("Hello World")
	result := s.ToLowerChain().Value()
	expected := "hello world"
	assert.Equal(t, expected, result, "ToLowerChain() = %v; want %v", result, expected)
}

func TestToUpperChain(t *testing.T) {
	s := New("hello world")
	result := s.ToUpperChain().Value()
	expected := "HELLO WORLD"
	assert.Equal(t, expected, result, "ToUpperChain() = %v; want %v", result, expected)
}

func TestToTitleChain(t *testing.T) {
	s := New("hello world")
	result := s.ToTitleChain().Value()
	expected := "Hello World"
	assert.Equal(t, expected, result, "ToTitleChain() = %v; want %v", result, expected)
}

func TestChainedMethods(t *testing.T) {
	s := New("gO LaNg")
	result := s.ToLowerChain().ToUpperChain().ToTitleChain().Value()
	expected := "Go Lang"
	assert.Equal(t, expected, result, "Chained methods = %v; want %v", result, expected)
}
