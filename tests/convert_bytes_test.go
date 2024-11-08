/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 10:50:50
 * @FilePath: \go-toolbox\tests\convert_bytes_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/convert"
)

func TestBytesToBCC(t *testing.T) {
	tests := []struct {
		input    []byte
		expected byte
	}{
		{[]byte{0x01, 0x02, 0x03}, 0x00},
		{[]byte{0xFF, 0xFF, 0xFF}, 0xFF},
		{[]byte{0x00}, 0x00},
	}

	for _, test := range tests {
		result := convert.BytesToBCC(test.input)
		if result != test.expected {
			t.Errorf("BytesBCC(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestByteToBinStr(t *testing.T) {
	tests := []struct {
		input    byte
		expected string
	}{
		{0, "00000000"},
		{1, "00000001"},
		{255, "11111111"},
	}

	for _, test := range tests {
		result := convert.ByteToBinStr(test.input)
		if result != test.expected {
			t.Errorf("ByteToBinStr(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestBytesToBinStr(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{0, 1, 2}, "000000000000000100000010"},
		{[]byte{255}, "11111111"},
		{[]byte{}, ""},
	}

	for _, test := range tests {
		result := convert.BytesToBinStr(test.input)
		if result != test.expected {
			t.Errorf("BytesToBinStr(%v) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestBytesToBinStrWithSplit(t *testing.T) {
	tests := []struct {
		input    []byte
		split    string
		expected string
	}{
		{[]byte{0, 1, 2}, " ", "00000000 00000001 00000010"},
		{[]byte{255}, "", "11111111"},
		{[]byte{}, "-", ""},
	}

	for _, test := range tests {
		result := convert.BytesToBinStrWithSplit(test.input, test.split)
		if result != test.expected {
			t.Errorf("BytesToBinStrWithSplit(%v, %s) = %s; want %s", test.input, test.split, result, test.expected)
		}
	}
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// createImage 创建一张简单的图像并返回其 Base64 编码
func createImage() (string, error) {
	// 创建一个 100x100 像素的图像
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// 填充背景为白色
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 画一个红色的矩形
	red := color.RGBA{255, 0, 0, 255}
	draw.Draw(img, image.Rect(10, 10, 90, 90), &image.Uniform{red}, image.Point{}, draw.Over)

	// 创建一个字节缓冲区
	buf := new(bytes.Buffer)

	// 将图像编码为 PNG 格式并写入缓冲区
	err := png.Encode(buf, img)
	if err != nil {
		return "", err
	}

	// 将字节缓冲区转换为 Base64 字符串
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// TestB64ToByte 测试 B64ToByte 函数
func TestB64ToByte(t *testing.T) {
	validB64, err := createImage()
	if err != nil {
		t.Fatalf("Error creating image: %v", err)
	}

	imageBytes, err := convert.B64ToByte(validB64)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(imageBytes) == 0 {
		t.Fatal("Expected non-empty byte slice")
	}
}

func TestSliceByteToString(t *testing.T) {
	b := []byte("hello world")
	s := convert.SliceByteToString(b)

	if s != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", s)
	}
}

func TestStringToSliceByte(t *testing.T) {
	s := "hello world"
	b := convert.StringToSliceByte(s)

	if string(b) != s {
		t.Errorf("expected '%s', got '%s'", s, b)
	}
}
