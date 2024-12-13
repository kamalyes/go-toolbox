/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 01:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-13 10:55:11
 * @FilePath: \go-toolbox\tests\imgix_drawer_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/fogleman/gg"
	"github.com/kamalyes/go-toolbox/pkg/imgix"
	"github.com/stretchr/testify/assert"
)

// saveImgixDrawerImage 保存图像到文件
func saveImgixDrawerImage(ctx *gg.Context, filename string) error {
	img := ctx.Image()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// compareImages 比较两个图像的像素
func compareImages(img1, img2 image.Image) bool {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1 != bounds2 {
		return false
	}

	for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
		for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
			if img1.At(x, y) != img2.At(x, y) {
				return false
			}
		}
	}
	return true
}

func TestNewGraphicsRenderer(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	assert.NotNil(t, renderer)
	assert.Equal(t, ctx, renderer.GgCtx)
}

func TestUseDefaultDashed(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	renderer.SetDashStyle(5, 5) // 设置虚线样式
	renderer.UseDefaultDashed() // 使用默认虚线

	renderer.DrawLine(100, 100, 200, 200, 2)
	err := saveImgixDrawerImage(ctx, "test_default_dashed.png")
	assert.NoError(t, err)
	defer os.Remove("test_default_dashed.png")
}

func TestUseSolidLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	renderer.UseSolidLine() // 使用实线
	renderer.DrawLine(100, 100, 200, 200, 2)
	err := saveImgixDrawerImage(ctx, "test_solid_line.png")
	assert.NoError(t, err)
	defer os.Remove("test_solid_line.png")
}

func TestSetDashed(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	renderer.SetDashed(5, 5) // 设置虚线样式
	renderer.DrawLine(100, 100, 200, 200, 2)
	err := saveImgixDrawerImage(ctx, "test_set_dashed.png")
	assert.NoError(t, err)
	defer os.Remove("test_set_dashed.png")
}

func TestDrawCurvedLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	start := &gg.Point{X: 100, Y: 100}
	end := &gg.Point{X: 200, Y: 200}
	control := renderer.CalculateControlPoint(start, end, 50) // 计算控制点

	renderer.DrawCurvedLine(start, end, control)
	err := saveImgixDrawerImage(ctx, "test_curved_line.png")
	assert.NoError(t, err)
	defer os.Remove("test_curved_line.png")
}

func TestDrawPolygon(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	points := []gg.Point{
		{X: 100, Y: 100},
		{X: 200, Y: 100},
		{X: 150, Y: 200},
	}
	renderer.DrawPolygon(points)
	err := saveImgixDrawerImage(ctx, "test_polygon.png")
	assert.NoError(t, err)
	defer os.Remove("test_polygon.png")
}

func TestDrawCenteredMultiLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	startXs := []float64{100}
	endXs := []float64{300}
	startYs := []float64{100}
	endYs := []float64{100}
	textGroups := [][]string{
		{"Hello", "World"},
	}

	renderer.DrawCenteredMultiLine(startXs, endXs, startYs, endYs, textGroups, 0, true)
	err := saveImgixDrawerImage(ctx, "test_centered_multiline.png")
	assert.NoError(t, err)
	defer os.Remove("test_centered_multiline.png")
}

func TestDrawLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	// 绘制线条
	renderer.DrawLine(100, 100, 200, 200, 2)

	// 保存图像
	testImagePath := "test_line.png"
	err := saveImgixDrawerImage(ctx, testImagePath)
	assert.NoError(t, err)
	defer os.Remove(testImagePath) // 测试结束后删除图像

	// 读取并比较图像
	expectedImg := gg.NewContext(800, 600)
	expectedImg.SetColor(color.Black)
	expectedImg.SetLineWidth(2)
	expectedImg.DrawLine(100, 100, 200, 200)
	expectedImg.Stroke()

	// 保存预期图像

	expectedImagePath := "expected_line.png"
	err = saveImgixDrawerImage(expectedImg, expectedImagePath)
	assert.NoError(t, err)
	defer os.Remove(expectedImagePath) // 测试结束后删除预期图像

	// 比较图像
	assert.True(t, compareImages(ctx.Image(), expectedImg.Image()))
}

func TestDrawRectangle(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	// 绘制矩形
	left := &gg.Point{X: 100, Y: 100}
	top := &gg.Point{X: 100, Y: 100}
	bottom := &gg.Point{X: 200, Y: 200}
	right := &gg.Point{X: 200, Y: 200}

	renderer.DrawRectangle(left, top, bottom, right)

	// 保存图像
	testImagePath := "test_rectangle.png"
	err := saveImgixDrawerImage(ctx, testImagePath)
	assert.NoError(t, err)
	defer os.Remove(testImagePath) // 测试结束后删除图像

	// 读取并比较图像
	expectedImg := gg.NewContext(800, 600)
	expectedImg.SetColor(color.Black)
	expectedImg.DrawRectangle(100, 100, 100, 100) // 宽度和高度都是100
	expectedImg.Stroke()

	// 保存预期图像
	expectedImagePath := "expected_rectangle.png"
	err = saveImgixDrawerImage(expectedImg, expectedImagePath)
	assert.NoError(t, err)
	defer os.Remove(expectedImagePath) // 测试结束后删除预期图像

	// 比较图像
	assert.True(t, compareImages(ctx.Image(), expectedImg.Image()))
}

func TestDrawCircle(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	// 绘制圆形
	renderer.DrawCircle(400, 300, 50)

	// 保存图像
	testImagePath := "test_circle.png"
	err := saveImgixDrawerImage(ctx, testImagePath)
	assert.NoError(t, err)
	defer os.Remove(testImagePath) // 测试结束后删除图像

	// 读取并比较图像
	expectedImg := gg.NewContext(800, 600)
	expectedImg.SetColor(color.Black)
	expectedImg.DrawCircle(400, 300, 50)
	expectedImg.Stroke()

	// 保存预期图像
	expectedImagePath := "expected_circle.png"
	err = saveImgixDrawerImage(expectedImg, expectedImagePath)
	assert.NoError(t, err)
	defer os.Remove(expectedImagePath) // 测试结束后删除预期图像

	// 比较图像
	assert.True(t, compareImages(ctx.Image(), expectedImg.Image()))
}

func TestCalculateFractionPoint(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	start := &gg.Point{X: 0, Y: 0}
	end := &gg.Point{X: 100, Y: 100}
	fraction := 0.5

	point := renderer.CalculateFractionPoint(start, end, fraction)

	assert.Equal(t, 50.0, point.X)
	assert.Equal(t, 50.0, point.Y)
}

func TestGetLTRB(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	features := map[string]gg.Point{
		"point1": {X: 100, Y: 100},
		"point2": {X: 200, Y: 200},
		"point3": {X: 50, Y: 50},
	}

	ltrb := renderer.GetLTRB(features)

	assert.Equal(t, 50.0, ltrb.Left)
	assert.Equal(t, 50.0, ltrb.Top)
	assert.Equal(t, 200.0, ltrb.Right)
	assert.Equal(t, 200.0, ltrb.Bottom)
}

func TestGetFacialXYByKey(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	features := map[string]gg.Point{
		"point1": {X: 100, Y: 100},
		"point2": {X: 200, Y: 200},
	}

	point := renderer.GetFacialXYByKey(features, "point1")
	assert.NotNil(t, point)
	assert.Equal(t, 100.0, point.X)
	assert.Equal(t, 100.0, point.Y)

	point = renderer.GetFacialXYByKey(features, "point3") // point3 不存在
	assert.Nil(t, point)
}

func TestUpdateBounds(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	left, top, right, bottom := float64(100), float64(100), float64(200), float64(200)
	renderer.UpdateBounds(50, 50, &left, &top, &right, &bottom)

	assert.Equal(t, 50.0, left)
	assert.Equal(t, 50.0, top)
	assert.Equal(t, 200.0, right)
	assert.Equal(t, 200.0, bottom)

	renderer.UpdateBounds(250, 250, &left, &top, &right, &bottom)

	assert.Equal(t, 50.0, left)
	assert.Equal(t, 50.0, top)
	assert.Equal(t, 250.0, right)  // right 更新为 250
	assert.Equal(t, 250.0, bottom) // bottom 更新为 250
}
