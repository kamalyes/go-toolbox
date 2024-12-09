/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-09 12:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-09 13:50:55
 * @FilePath: \go-toolbox\tests\imgix_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kamalyes/go-toolbox/pkg/imgix"
)

// mockServer 用于模拟不同的 HTTP 响应
func mockServer(responseBody []byte, statusCode int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode) // 设置响应状态码
		w.Write(responseBody)     // 写入响应体
	})
	return httptest.NewServer(handler) // 返回一个新的测试服务器
}

// TestLoadImageFromURL 测试加载图片的主函数
func TestLoadImageFromURL(t *testing.T) {
	t.Run("Valid PNG Image", testValidPNGImage)
	t.Run("Valid JPEG Image", testValidJPEGImage)
	t.Run("Invalid URL", testInvalidURL)
	t.Run("404 Error", test404Error)
	t.Run("Non-image Content", testNonImageContent)
	t.Run("Network Error", testNetworkError)
}

// testValidPNGImage 测试有效的 PNG 图片
func testValidPNGImage(t *testing.T) {
	pngImage := createTestPNGImage()              // 创建测试用的 PNG 图片
	server := mockServer(pngImage, http.StatusOK) // 启动模拟服务器
	defer server.Close()                          // 确保服务器在测试结束后关闭

	img, err := imgix.LoadImageFromURL(server.URL) // 调用函数加载图片
	assert.NoError(t, err)                         // 期望没有错误
	assert.NotNil(t, img)                          // 期望返回有效的图片
}

// testValidJPEGImage 测试有效的 JPEG 图片
func testValidJPEGImage(t *testing.T) {
	jpegImage := createTestJPEGImage()             // 创建测试用的 JPEG 图片
	server := mockServer(jpegImage, http.StatusOK) // 启动模拟服务器
	defer server.Close()                           // 确保服务器在测试结束后关闭

	img, err := imgix.LoadImageFromURL(server.URL) // 调用函数加载图片
	assert.NoError(t, err)                         // 期望没有错误
	assert.NotNil(t, img)                          // 期望返回有效的图片
}

// testInvalidURL 测试无效的 URL
func testInvalidURL(t *testing.T) {
	_, err := imgix.LoadImageFromURL("http://invalid-url") // 调用函数加载无效的 URL
	assert.Error(t, err)                                   // 期望返回错误
}

// test404Error 测试 404 错误
func test404Error(t *testing.T) {
	server := mockServer(nil, http.StatusNotFound) // 启动模拟服务器，返回 404 状态
	defer server.Close()                           // 确保服务器在测试结束后关闭

	_, err := imgix.LoadImageFromURL(server.URL) // 调用函数加载图片
	assert.Error(t, err)                         // 期望返回错误
}

// testNonImageContent 测试非图片内容
func testNonImageContent(t *testing.T) {
	server := mockServer([]byte("not an image"), http.StatusOK) // 启动模拟服务器，返回非图片内容
	defer server.Close()                                        // 确保服务器在测试结束后关闭

	_, err := imgix.LoadImageFromURL(server.URL) // 调用函数加载图片
	assert.Error(t, err)                         // 期望返回错误
}

// testNetworkError 测试网络错误
func testNetworkError(t *testing.T) {
	_, err := imgix.LoadImageFromURL("http://127.0.0.1:9999") // 假设该端口没有服务在运行
	assert.Error(t, err)                                      // 期望返回错误
}

// createTestPNGImage 创建一个简单的 PNG 图片用于测试
func createTestPNGImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1)) // 创建一个 1x1 像素的 RGBA 图片
	img.Set(0, 0, image.White)                   // 将像素设置为白色
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil { // 编码为 PNG 格式
		panic(err) // 如果编码失败，抛出错误
	}
	return buf.Bytes() // 返回 PNG 图片的字节切片
}

// createTestJPEGImage 创建一个简单的 JPEG 图片用于测试
func createTestJPEGImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1)) // 创建一个 1x1 像素的 RGBA 图片
	img.Set(0, 0, image.White)                   // 将像素设置为白色
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil { // 编码为 JPEG 格式
		panic(err) // 如果编码失败，抛出错误
	}
	return buf.Bytes() // 返回 JPEG 图片的字节切片
}

func TestWriterImage(t *testing.T) {
	// 创建一个测试用的图像
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255}) // 填充红色
		}
	}

	// 创建一个缓冲区用于保存压缩后的图像
	var buf bytes.Buffer

	// 测试 JPEG
	err := imgix.WriterImage(img, 80, imgix.JPEG, &buf)
	assert.NoError(t, err)          // 期望没有错误
	assert.Greater(t, buf.Len(), 0) // 期望输出不为空

	// 测试 PNG
	buf.Reset() // 清空缓冲区
	err = imgix.WriterImage(img, 0, imgix.PNG, &buf)
	assert.NoError(t, err)          // 期望没有错误
	assert.Greater(t, buf.Len(), 0) // 期望输出不为空

	// 测试 GIF
	buf.Reset() // 清空缓冲区
	err = imgix.WriterImage(img, 0, imgix.GIF, &buf)
	assert.NoError(t, err)          // 期望没有错误
	assert.Greater(t, buf.Len(), 0) // 期望输出不为空
}

func TestWriterImageFromBytes(t *testing.T) {
	// 创建一个测试用的 JPEG 图像
	jpegImage := createTestJPEGImage()

	// 测试 WriterImageFromBytes 函数
	var buf bytes.Buffer
	err := imgix.WriterImageFromBytes(jpegImage, 80, imgix.PNG, &buf)
	assert.NoError(t, err)          // 期望没有错误
	assert.Greater(t, buf.Len(), 0) // 期望输出不为空

	// 测试从无效的字节切片加载图像
	err = imgix.WriterImageFromBytes([]byte("not an image"), 80, imgix.PNG, &buf)
	assert.Error(t, err) // 期望返回错误
}

func createTestImageBuf(format imgix.ImageFormat) (*bytes.Buffer, error) {
	// 创建一个测试用的图像
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{0, 255, 0, 255}) // 填充绿色
		}
	}

	// 将图像编码到缓冲区
	buf := new(bytes.Buffer)
	var err error
	switch format {
	case imgix.JPEG:
		err = jpeg.Encode(buf, img, nil)
	case imgix.PNG:
		err = png.Encode(buf, img)
	case imgix.GIF:
		err = gif.Encode(buf, img, nil)
	}

	if err != nil {
		return nil, err
	}
	return buf, nil
}

func TestSaveBufToImageFile(t *testing.T) {
	formats := []imgix.ImageFormat{imgix.JPEG, imgix.PNG, imgix.GIF}

	for _, format := range formats {
		t.Run(format.String(), func(t *testing.T) {
			buf, err := createTestImageBuf(format)
			assert.NoError(t, err) // 期望没有错误

			// 测试保存有效的图像数据
			filePath := "test_image." + format.String()
			err = imgix.SaveBufToImageFile(buf, filePath, format)
			assert.NoError(t, err) // 期望没有错误

			// 检查文件是否存在
			_, err = os.Stat(filePath)
			assert.NoError(t, err) // 期望文件存在

			// 清理生成的文件
			os.Remove(filePath)
		})
	}

	// 测试使用无效的图像数据
	invalidBuf := bytes.NewBufferString("invalid data")
	err := imgix.SaveBufToImageFile(invalidBuf, "test_invalid_image.png", imgix.PNG)
	assert.Error(t, err) // 期望返回错误

	buf, err := createTestImageBuf(imgix.JPEG)

	// 测试使用不支持的格式
	err = imgix.SaveBufToImageFile(buf, "test_image.invalid", imgix.ImageFormat(100))
	assert.Error(t, err) // 期望返回错误

	// 测试文件创建失败的情况
	// 这里我们可以模拟文件创建失败，可以通过创建一个只读的目录来实现
	os.Mkdir("readonly_dir", 0755)
	defer os.Remove("readonly_dir")

	// 将目录权限改为只读
	os.Chmod("readonly_dir", 0400) // 只读权限

	err = imgix.SaveBufToImageFile(buf, "readonly_dir/test_image.png", imgix.PNG)
	assert.Error(t, err) // 期望返回错误

	// 恢复目录权限
	os.Chmod("readonly_dir", 0755)
	os.Remove("readonly_dir")
}
