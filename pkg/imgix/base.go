/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-09 12:15:51
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-09 13:58:58
 * @FilePath: \go-toolbox\pkg\imgix\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package imgix

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

// LoadImageFromURL 从给定的 URL 加载图片并返回 image.Image 对象
func LoadImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应体到内存
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 使用 image 包来解码图片
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}

// ImageFormat 定义图像格式类型
type ImageFormat int

// 定义支持的图像格式常量

const (
	JPEG ImageFormat = iota
	JPG
	PNG
	GIF
)

// String 返回图像格式的字符串表示
func (f ImageFormat) String() string {
	switch f {
	case JPEG:
		return "jpeg"
	case PNG:
		return "png"
	case GIF:
		return "gif"
	default:
		return "unknown"
	}
}

// WriterImage 将图像写入输出
func WriterImage(img image.Image, quality int, format ImageFormat, output io.Writer) error {
	// 使用类型来处理不同的图像格式
	switch format {
	case JPEG, JPG:
		if quality < 1 || quality > 100 {
			return errors.New("quality must be between 1 and 100 for JPEG")
		}
		return jpeg.Encode(output, img, &jpeg.Options{Quality: quality}) // 编码为 JPEG 格式
	case PNG:
		return png.Encode(output, img) // 编码为 PNG 格式
	case GIF:
		return gif.Encode(output, img, nil) // 编码为 GIF 格式
	default:
		return fmt.Errorf("unsupported format: %s", format) // 不支持的格式
	}
}

// WriterImageFromBytes 将字节切片转换为图像并写入输出
func WriterImageFromBytes(data []byte, quality int, format ImageFormat, output io.Writer) error {
	// 解码字节切片为图像
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return errors.New("failed to decode image: " + err.Error())
	}

	// 使用 WriterImage 函数将图像写入输出
	return WriterImage(img, quality, format, output)
}

// SaveBufToImageFile 将Buf数据保存到本地文件
func SaveBufToImageFile(buf *bytes.Buffer, filePath string, format ImageFormat) error {
	var img image.Image
	var err error

	// 根据格式解码图像
	switch format {
	case JPEG:
		img, err = jpeg.Decode(buf)
	case PNG:
		img, err = png.Decode(buf)
	case GIF:
		img, err = gif.Decode(buf)
	default:
		return errors.New("unsupported image format")
	}

	if err != nil {
		return err
	}

	// 创建文件
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 根据格式编码图像并写入文件
	switch format {
	case JPEG:
		err = jpeg.Encode(outFile, img, nil)
	case PNG:
		err = png.Encode(outFile, img)
	case GIF:
		err = gif.Encode(outFile, img, nil)
	}

	return err
}
