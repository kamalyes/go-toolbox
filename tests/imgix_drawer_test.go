/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 01:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-13 13:29:26
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
	"sync"
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

	dashOptions := imgix.NewDashOptions(5, 7)

	rendererX := imgix.NewGraphicsRenderer(ctx, dashOptions)
	renderXDashOptions := rendererX.GetDashOptions()

	assert.NotNil(t, rendererX, "NewGraphicsRenderer should return a non-nil renderer")
	assert.Equal(t, dashOptions, rendererX.GetDashOptions(), "GetDashOptions() should return the correct DashOptions")
	assert.Equal(t, dashOptions.DashLength(), renderXDashOptions.DashLength(), "GetDashOptions.DashLength should return the correct DashOptions")
	assert.Equal(t, dashOptions.GapLength(), renderXDashOptions.GapLength(), "GetDashOptions.GapLength should return the correct DashOptions")
}

func TestImgixGrSaveImage(t *testing.T) {
	// 创建一个新的 gg.Context
	width, height := 100, 100
	ctx := gg.NewContext(width, height)

	// 创建 GraphicsRenderer 实例
	renderer := imgix.NewGraphicsRenderer(ctx)

	// 使用 GraphicsRenderer 绘制一些图形
	renderer.DrawRectangle(&gg.Point{X: 10, Y: 10}, &gg.Point{X: 90, Y: 10}, &gg.Point{X: 90, Y: 90}, &gg.Point{X: 10, Y: 90})

	// 保存图像
	imageName := "test_image"
	imageFormat := imgix.PNG
	quality := 100
	renderer.SaveImage(imageName, quality, imageFormat)

	// 5. 验证保存的图像文件是否存在
	filePath := imageName + ".png"
	_, err := os.Stat(filePath)
	assert.NoError(t, err, "Expected image file to be created")

	// 6. 清理测试生成的文件
	err = os.Remove(filePath)
	assert.NoError(t, err, "Expected to remove test image file")
}

func TestUseDefaultDashed(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)
	renderer.DrawLine(100, 100, 200, 200, 2)
	filePath := "test_default_dashed.png"
	err := saveImgixDrawerImage(ctx, filePath)
	assert.NoError(t, err)
	defer os.Remove(filePath)
}

func TestUseSolidLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	renderer.UseSolidLine() // 使用实线
	renderer.DrawLine(100, 100, 200, 200, 2)
	filePath := "test_solid_line.png"
	err := saveImgixDrawerImage(ctx, filePath)
	assert.NoError(t, err)
	defer os.Remove(filePath)
}

// TestDrawWithStroke 测试 DrawWithStroke 方法
func TestDrawWithStroke(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	drawCalled := false
	drawFunc := func() {
		drawCalled = true
	}

	// 测试不调用 Stroke
	renderer.DrawWithStroke(drawFunc, false)
	assert.True(t, drawCalled, "Expected drawFunc to be called")
}

func TestSetDashed(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	renderer.SetDashed(5, 5) // 设置虚线样式
	renderer.DrawLine(100, 100, 200, 200, 2)
	filePath := "test_set_dashed.png"
	err := saveImgixDrawerImage(renderer.GgCtx, filePath)
	assert.NoError(t, err)
	defer os.Remove(filePath)
}

func TestDrawCurvedLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	start := &gg.Point{X: 100, Y: 100}
	end := &gg.Point{X: 200, Y: 200}
	control := renderer.CalculateControlPoint(start, end, 50) // 计算控制点

	renderer.DrawCurvedLine(start, end, control)
	filePath := "test_curved_line.png"
	err := saveImgixDrawerImage(renderer.GgCtx, filePath)
	assert.NoError(t, err)
	defer os.Remove(filePath)
}

func TestDrawHorizontalLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	// 设置背景色为白色
	ctx.SetColor(color.White)
	ctx.Clear() // 填充背景色

	// 设置线条颜色为黑色
	ctx.SetColor(color.Black)
	renderer := imgix.NewGraphicsRenderer(ctx)

	renderer.DrawHorizontalLine(10, 0, 20)
	filePath := "test_horizontal_line.png"
	err := saveImgixDrawerImage(renderer.GgCtx, filePath)
	assert.NoError(t, err)
	defer os.Remove(filePath)
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
	filePath := "test_polygon.png"
	err := saveImgixDrawerImage(renderer.GgCtx, filePath)
	assert.NoError(t, err)
	defer os.Remove(filePath)
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
	filePath := "test_centered_multiline.png"
	err := saveImgixDrawerImage(renderer.GgCtx, filePath)
	assert.NoError(t, err)
	defer os.Remove(filePath)
}

func TestDrawLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	// 绘制线条
	renderer.DrawLine(100, 100, 200, 200, 2)

	// 保存图像
	testImagePath := "test_line.png"
	err := saveImgixDrawerImage(renderer.GgCtx, testImagePath)
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
	err := saveImgixDrawerImage(renderer.GgCtx, testImagePath)
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
	err := saveImgixDrawerImage(renderer.GgCtx, testImagePath)
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

func TestImageLTRB(t *testing.T) {
	// 创建一个 ImageLTRB 实例
	ltrb := imgix.ImageLTRB{
		Left:   10,
		Top:    20,
		Right:  50,
		Bottom: 80,
	}

	// 测试 Width 方法
	expectedWidth := ltrb.Right - ltrb.Left // 40.0
	assert.Equal(t, expectedWidth, ltrb.Width(), "Width() should be equal")

	// 测试 Height 方法
	expectedHeight := ltrb.Bottom - ltrb.Top // 60.0
	assert.Equal(t, expectedHeight, ltrb.Height(), "Height() should be equal")

	// 测试 Center 方法
	expectedCenterX := ltrb.Left + (expectedWidth / 2.0) // 30.0
	expectedCenterY := ltrb.Top + (expectedHeight / 2.0) // 50.0
	gotX, gotY := ltrb.Center()
	assert.Equal(t, expectedCenterX, gotX, "Center X should be equal")
	assert.Equal(t, expectedCenterY, gotY, "Center Y should be equal")

	// 测试 Contains 方法
	pointInsideX, pointInsideY := 30.0, 30.0
	assert.True(t, ltrb.Contains(pointInsideX, pointInsideY), "Contains should return true for point inside")

	pointOutsideX, pointOutsideY := 5.0, 5.0
	assert.False(t, ltrb.Contains(pointOutsideX, pointOutsideY), "Contains should return false for point outside")
	assert.NotEmpty(t, ltrb.String(), "String() should return an empty string when all bounds are zero")
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

func TestConcurrentDrawLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	const numGoroutines = 5
	var wg sync.WaitGroup
	lines := make(chan []float64, numGoroutines)

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			startX := float64(index * 2)
			startY := float64(index * 2)
			endX := startX + 50
			endY := startY + 50

			renderer.DrawLine(startX, startY, endX, endY, 2)

			lines <- []float64{startX, startY, endX, endY}
		}(i)
	}

	wg.Wait()
	close(lines)

	for line := range lines {
		startX, startY, endX, endY := line[0], line[1], line[2], line[3]
		if startX < 0 || startY < 0 || endX > 800 || endY > 600 {
			t.Errorf("Line out of bounds: start(%f, %f), end(%f, %f)", startX, startY, endX, endY)
		}
	}

	err := saveImgixDrawerImage(renderer.GgCtx, "test_concurrent_draw_line.png")
	assert.NoError(t, err)
	defer os.Remove("test_concurrent_draw_line.png")
}

func TestConcurrentGraphicsRenderer(t *testing.T) {
	var wg sync.WaitGroup
	numTests := 10            // 定义可用测试的数量
	numConcurrentTests := 100 // 并发运行的测试数量

	for i := 0; i < numConcurrentTests; i++ {
		wg.Add(1)
		go func(testNum int) {
			defer wg.Done()
			ctx := gg.NewContext(800, 600)
			renderer := imgix.NewGraphicsRenderer(ctx)

			// 使用取模来选择测试
			testIndex := testNum % numTests
			switch testIndex {
			case 0:
				assert.NotNil(t, renderer)
				assert.Equal(t, ctx, renderer.GgCtx)
			case 1:
				renderer.DrawLine(100, 100, 200, 200, 2)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_default_dashed.png")
				assert.NoError(t, err)
				defer os.Remove("test_default_dashed.png")
			case 2:
				renderer.UseSolidLine()
				renderer.DrawLine(100, 100, 200, 200, 2)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_solid_line.png")
				assert.NoError(t, err)
				defer os.Remove("test_solid_line.png")
			case 3:
				renderer.SetDashed(5, 5)
				renderer.DrawLine(100, 100, 200, 200, 2)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_set_dashed.png")
				assert.NoError(t, err)
				defer os.Remove("test_set_dashed.png")
			case 4:
				start := &gg.Point{X: 100, Y: 100}
				end := &gg.Point{X: 200, Y: 200}
				control := renderer.CalculateControlPoint(start, end, 50)
				renderer.DrawCurvedLine(start, end, control)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_curved_line.png")
				assert.NoError(t, err)
				defer os.Remove("test_curved_line.png")
			case 5:
				points := []gg.Point{
					{X: 100, Y: 100},
					{X: 200, Y: 100},
					{X: 150, Y: 200},
				}
				renderer.DrawPolygon(points)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_polygon.png")
				assert.NoError(t, err)
				defer os.Remove("test_polygon.png")
			case 6:
				startXs := []float64{100}
				endXs := []float64{300}
				startYs := []float64{100}
				endYs := []float64{100}
				textGroups := [][]string{
					{"Hello", "World"},
				}
				renderer.DrawCenteredMultiLine(startXs, endXs, startYs, endYs, textGroups, 0, true)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_centered_multiline.png")
				assert.NoError(t, err)
				defer os.Remove("test_centered_multiline.png")
			case 7:
				renderer.DrawLine(100, 100, 200, 200, 2)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_line.png")
				assert.NoError(t, err)
				defer os.Remove("test_line.png")
			case 8:
				left := &gg.Point{X: 100, Y: 100}
				top := &gg.Point{X: 100, Y: 100}
				bottom := &gg.Point{X: 200, Y: 200}
				right := &gg.Point{X: 200, Y: 200}
				renderer.DrawRectangle(left, top, bottom, right)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_rectangle.png")
				assert.NoError(t, err)
				defer os.Remove("test_rectangle.png")
			case 9:
				renderer.DrawCircle(400, 300, 50)
				err := saveImgixDrawerImage(renderer.GgCtx, "test_circle.png")
				assert.NoError(t, err)
				defer os.Remove("test_circle.png")
			}
		}(i)
	}

	wg.Wait() // 等待所有测试完成
}
