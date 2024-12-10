/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-09 12:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-10 09:21:55
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

// createImgixTestImage 创建一个简单的图像用于测试
func createImgixTestImage(format imgix.ImageFormat, width, height int, color color.Color) (image.Image, *bytes.Buffer, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color) // 填充指定颜色
		}
	}

	var buf bytes.Buffer
	var err error
	switch format {
	case imgix.JPEG:
		err = jpeg.Encode(&buf, img, nil)
	case imgix.PNG:
		err = png.Encode(&buf, img)
	case imgix.GIF:
		err = gif.Encode(&buf, img, nil)
	default:
		err = nil // 不支持的格式
	}
	return img, &buf, err
}

// createTempFile 创建一个临时文件并返回文件路径
func createTempFile(data []byte) (string, error) {
	tempFile, err := os.CreateTemp("", "test_image_*.jpg")
	if err != nil {
		return "", err
	}
	defer tempFile.Close() // 确保文件在函数结束时关闭

	if _, err := tempFile.Write(data); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
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
	_, buf, err := createImgixTestImage(imgix.PNG, 1, 1, image.White)
	assert.NoError(t, err)

	server := mockServer(buf.Bytes(), http.StatusOK) // 启动模拟服务器
	defer server.Close()                             // 确保服务器在测试结束后关闭

	img, err := imgix.LoadImageFromURL(server.URL) // 调用函数加载图片
	assert.NoError(t, err)                         // 期望没有错误
	assert.NotNil(t, img)                          // 期望返回有效的图片
}

// testValidJPEGImage 测试有效的 JPEG 图片
func testValidJPEGImage(t *testing.T) {
	_, buf, err := createImgixTestImage(imgix.JPEG, 1, 1, image.White)
	assert.NoError(t, err)

	server := mockServer(buf.Bytes(), http.StatusOK) // 启动模拟服务器
	defer server.Close()                             // 确保服务器在测试结束后关闭

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

func TestWriterImage(t *testing.T) {
	// 创建一个测试用的图像
	img, _, err := createImgixTestImage(imgix.JPEG, 100, 100, color.RGBA{255, 0, 0, 255})
	assert.NoError(t, err)

	// 测试 JPEG
	var buf bytes.Buffer
	err = imgix.WriterImage(img, 80, imgix.JPEG, &buf)
	assert.NoError(t, err)          // 期望没有错误
	assert.Greater(t, buf.Len(), 0) // 期望输出不为空

	// 测试 PNG
	buf.Reset() // 清空缓冲区
	img, _, err = createImgixTestImage(imgix.PNG, 100, 100, color.RGBA{255, 0, 0, 255})
	assert.NoError(t, err)
	err = imgix.WriterImage(img, 0, imgix.PNG, &buf)
	assert.NoError(t, err)          // 期望没有错误
	assert.Greater(t, buf.Len(), 0) // 期望输出不为空

	// 测试 GIF
	buf.Reset() // 清空缓冲区
	img, _, err = createImgixTestImage(imgix.GIF, 100, 100, color.RGBA{255, 0, 0, 255})
	assert.NoError(t, err)
	err = imgix.WriterImage(img, 0, imgix.GIF, &buf)
	assert.NoError(t, err)          // 期望没有错误
	assert.Greater(t, buf.Len(), 0) // 期望输出不为空
}

func TestWriterImageFromBytes(t *testing.T) {
	// 创建一个测试用的 JPEG 图像
	_, jpegImageBuf, err := createImgixTestImage(imgix.JPEG, 1, 1, image.White)
	assert.NoError(t, err)

	// 测试 WriterImageFromBytes 函数
	var buf bytes.Buffer
	err = imgix.WriterImageFromBytes(jpegImageBuf.Bytes(), 80, imgix.PNG, &buf)
	assert.NoError(t, err)          // 期望没有错误
	assert.Greater(t, buf.Len(), 0) // 期望输出不为空

	// 测试从无效的字节切片加载图像
	err = imgix.WriterImageFromBytes([]byte("not an image"), 80, imgix.PNG, &buf)
	assert.Error(t, err) // 期望返回错误
}

func TestSaveBufToImageFile(t *testing.T) {
	formats := []imgix.ImageFormat{imgix.JPEG, imgix.PNG, imgix.GIF}

	for _, format := range formats {
		t.Run(format.String(), func(t *testing.T) {
			_, buf, err := createImgixTestImage(format, 100, 100, image.White)
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

	_, buf, err := createImgixTestImage(imgix.JPEG, 100, 100, image.White)

	// 测试使用不支持的格式
	err = imgix.SaveBufToImageFile(buf, "test_image.invalid", imgix.ImageFormat(100))
	assert.Error(t, err) // 期望返回错误

	// 测试文件创建失败的情况
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

// TestLoadImageFromFile 测试从文件加载图像的功能
func TestLoadImageFromFile(t *testing.T) {
	// 创建一个测试用的 JPEG 图像文件
	_, buf, err := createImgixTestImage(imgix.JPEG, 10, 10, color.RGBA{255, 0, 0, 255})
	assert.NoError(t, err)

	filePath, err := createTempFile(buf.Bytes())
	assert.NoError(t, err)    // 期望没有错误
	defer os.Remove(filePath) // 确保测试结束后删除文件

	// 测试加载有效的图像文件
	loadedImg, err := imgix.LoadImageFromFile(filePath)
	assert.NoError(t, err)      // 期望没有错误
	assert.NotNil(t, loadedImg) // 期望返回有效的图像
}

// TestLoadImageFromFileInvalidPath 测试无效文件路径
func TestLoadImageFromFileInvalidPath(t *testing.T) {
	_, err := imgix.LoadImageFromFile("invalid/path/to/image.jpg") // 调用函数加载无效路径
	assert.Error(t, err)                                           // 期望返回错误
}

// TestLoadImageFromFileUnsupportedFormat 测试不支持的图像格式
func TestLoadImageFromFileUnsupportedFormat(t *testing.T) {
	// 创建一个临时文件并写入无效数据
	tempFile, err := os.CreateTemp("", "test_invalid_image_*.txt")
	assert.NoError(t, err)           // 期望没有错误
	defer os.Remove(tempFile.Name()) // 确保测试结束后删除文件

	_, err = tempFile.WriteString("this is not an image") // 写入非图像数据
	assert.NoError(t, err)                                // 期望没有错误
	tempFile.Close()                                      // 关闭文件

	// 测试加载无效的图像文件
	_, err = imgix.LoadImageFromFile(tempFile.Name())
	assert.Error(t, err) // 期望返回错误
}

func TestEncodeImageToBase64(t *testing.T) {
	// 创建一个测试用的图像
	img, _, err := createImgixTestImage(imgix.JPEG, 100, 100, color.RGBA{255, 0, 0, 255})
	assert.NoError(t, err) // 期望没有错误

	// 测试有效的图像编码为 Base64
	base64Str, err := imgix.EncodeImageToBase64(img, imgix.JPEG, 80)
	assert.NoError(t, err)                                   // 期望没有错误
	assert.NotEmpty(t, base64Str)                            // 期望返回的字符串不为空
	assert.Contains(t, base64Str, "data:image/jpeg;base64,") // 确保包含正确的前缀
}

func TestLoadImageFromFileAndEncodeToBase64(t *testing.T) {
	// 创建一个测试用的 JPEG 图像文件
	_, buf, err := createImgixTestImage(imgix.JPEG, 10, 10, color.RGBA{255, 0, 0, 255})
	assert.NoError(t, err)

	filePath, err := createTempFile(buf.Bytes())
	assert.NoError(t, err)    // 期望没有错误
	defer os.Remove(filePath) // 确保测试结束后删除文件

	// 测试加载有效的图像文件并编码为 Base64
	base64Str, err := imgix.LoadImageFromFileAndEncodeToBase64(filePath, imgix.JPEG, 80)
	assert.NoError(t, err)                                   // 期望没有错误
	assert.NotEmpty(t, base64Str)                            // 期望返回的字符串不为空
	assert.Contains(t, base64Str, "data:image/jpeg;base64,") // 确保包含正确的前缀
}

func TestLoadImageFromFileAndEncodeToBase64InvalidPath(t *testing.T) {
	// 测试加载无效文件路径
	base64Str, err := imgix.LoadImageFromFileAndEncodeToBase64("invalid/path/to/image.jpg", imgix.JPEG, 80)
	assert.Error(t, err)       // 期望返回错误
	assert.Empty(t, base64Str) // 期望返回的字符串为空
}

func TestLoadImageFromFileAndEncodeToBase64UnsupportedFormat(t *testing.T) {
	// 创建一个临时文件并写入无效数据
	tempFile, err := os.CreateTemp("", "test_invalid_image_*.txt")
	assert.NoError(t, err)           // 期望没有错误
	defer os.Remove(tempFile.Name()) // 确保测试结束后删除文件

	_, err = tempFile.WriteString("this is not an image") // 写入非图像数据
	assert.NoError(t, err)                                // 期望没有错误
	tempFile.Close()                                      // 关闭文件

	// 测试加载无效的图像文件并编码为 Base64
	base64Str, err := imgix.LoadImageFromFileAndEncodeToBase64(tempFile.Name(), imgix.JPEG, 80)
	assert.Error(t, err)       // 期望返回错误
	assert.Empty(t, base64Str) // 期望返回的字符串为空
}
