/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 01:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-07 17:27:17
 * @FilePath: \go-toolbox\tests\imgix_drawer_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"sync"
	"testing"

	"github.com/disintegration/imaging"
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

	dashOptions := imgix.NewDashOptions(0, 0, 0)
	dashOptionsDashLength := dashOptions.DashLength()
	dashOptionsGapLength := dashOptions.GapLength()
	dashOptionsLineWidth := dashOptions.LineWidth()

	assert.Equal(t, dashOptionsDashLength, imgix.DashStyle(3), "dashOptionsDashLength should return the correct DashOptions")
	assert.Equal(t, dashOptionsGapLength, imgix.DashStyle(6), "dashOptionsGapLength should return the correct DashOptions")
	assert.Equal(t, dashOptionsLineWidth, imgix.DashStyle(2), "dashOptionsLineWidth should return the correct DashOptions")

	dashOptions = imgix.NewDashOptions(5, 7, 6)

	rendererX := imgix.NewGraphicsRenderer(ctx, dashOptions)
	renderXDashOptions := rendererX.GetDashOptions()

	assert.NotNil(t, rendererX, "NewGraphicsRenderer should return a non-nil renderer")
	assert.Equal(t, dashOptions, rendererX.GetDashOptions(), "GetDashOptions() should return the correct DashOptions")
	assert.Equal(t, dashOptions.DashLength(), renderXDashOptions.DashLength(), "GetDashOptions.DashLength should return the correct DashOptions")
	assert.Equal(t, dashOptions.GapLength(), renderXDashOptions.GapLength(), "GetDashOptions.GapLength should return the correct DashOptions")
	assert.Equal(t, dashOptions.LineWidth(), renderXDashOptions.LineWidth(), "GetDashOptions.LineWidth should return the correct DashOptions")
}

func TestUseDefaultDashed(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)
	renderer.DrawLineXYLineWidth(100, 100, 200, 200)
	filePath := "test_default_dashed.png"
	err := saveImgixDrawerImage(ctx, filePath)
	assert.NoError(t, err)
	defer os.Remove(filePath)
}

func TestUseSolidLine(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	renderer.UseSolidLine() // 使用实线
	renderer.DrawLineXYLineWidth(100, 100, 200, 200)
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
	renderer.DrawLineXYLineWidth(100, 100, 200, 200)
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
	control := imgix.CalculateControlPoint(start, end, 50, 0) // 计算控制点

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

	renderer.DrawHorizontalLine(imgix.HorizontalLine{
		Y:      10,
		LeftX:  0,
		RightX: 20,
	})
	renderer.DrawVerticalLine(imgix.VerticalLine{
		X:       10,
		TopY:    0,
		BottomY: 20,
	})
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

func DrawLineXYLineWidth(t *testing.T) {
	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)

	// 绘制线条
	renderer.DrawLineXYLineWidth(100, 100, 200, 200)

	// 保存图像
	testImagePath := "test_line.png"
	err := saveImgixDrawerImage(renderer.GgCtx, testImagePath)
	assert.NoError(t, err)
	defer os.Remove(testImagePath) // 测试结束后删除图像

	// 读取并比较图像
	expectedImg := gg.NewContext(800, 600)
	expectedImg.SetColor(color.Black)
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
	renderer.DrawRectangle(imgix.Rectangle{
		TopLeft:     &gg.Point{X: 100, Y: 100},
		BottomRight: &gg.Point{X: 200, Y: 200},
	}, 0)

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
	renderer.DrawCircle(&gg.Point{X: 400, Y: 300}, 50)

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

// TestCalculateFractionPoint 测试 CalculateFractionPoint 函数
func TestCalculateFractionPoint(t *testing.T) {
	startPoint := &gg.Point{X: 0, Y: 0} // 起始点
	endPoint := &gg.Point{X: 10, Y: 10} // 结束点

	tests := []struct {
		fraction float64
		mode     imgix.CalculateFractionPointMode
		expected *gg.Point
	}{
		{2.0, imgix.Add, &gg.Point{X: 12.0, Y: 12.0}},    // 加法模式
		{2.0, imgix.Subtract, &gg.Point{X: 8.0, Y: 8.0}}, // 减法模式
		{0.5, imgix.Multiply, &gg.Point{X: 5.0, Y: 5.0}}, // 乘法模式
		{2.0, imgix.Divide, &gg.Point{X: 5.0, Y: 5.0}},   // 除法模式
		{0.0, imgix.Divide, &gg.Point{X: 10.0, Y: 10.0}}, // 除法模式，分母为零
	}

	for _, test := range tests {
		point := imgix.CalculateFractionPoint(startPoint, endPoint, test.fraction, test.mode)
		assert.Equal(t, test.expected.X, point.X, "X coordinate mismatch")
		assert.Equal(t, test.expected.Y, point.Y, "Y coordinate mismatch")
	}
}

func TestGetLTRB(t *testing.T) {
	features := map[string]gg.Point{
		"point1": {X: 100, Y: 100},
		"point2": {X: 200, Y: 200},
		"point3": {X: 50, Y: 50},
	}

	ltrb := imgix.GetLTRB(features)

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

func TestUpdateBounds(t *testing.T) {
	left, top, right, bottom := float64(100), float64(100), float64(200), float64(200)
	imgix.UpdateBounds(50, 50, &left, &top, &right, &bottom)

	assert.Equal(t, 50.0, left)
	assert.Equal(t, 50.0, top)
	assert.Equal(t, 200.0, right)
	assert.Equal(t, 200.0, bottom)

	imgix.UpdateBounds(250, 250, &left, &top, &right, &bottom)

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

			renderer.DrawLineXYLineWidth(startX, startY, endX, endY)

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

// TestDrawLine 测试 DrawLine 方法
func TestDrawLine(t *testing.T) {
	startPoint := &gg.Point{X: 10, Y: 10}
	endPoint := &gg.Point{X: 20, Y: 20}

	ctx := gg.NewContext(800, 600)
	renderer := imgix.NewGraphicsRenderer(ctx)
	// 调用要测试的方法
	renderer.DrawLine(startPoint, endPoint)
}

func TestCleanCoordinates(t *testing.T) {
	tests := []struct {
		name     string
		points   map[string]gg.Point
		expected imgix.Coordinates
	}{
		{
			name: "Normal case with multiple points",
			points: map[string]gg.Point{
				"point1": {X: 1, Y: 2},
				"point2": {X: 3, Y: 4},
				"point3": {X: 0, Y: 5},
				"point4": {X: 2, Y: 1},
			},
			expected: imgix.Coordinates{
				Top:    imgix.HorizontalEdge{LeftMost: &gg.Point{X: 0, Y: 5}, RightMost: &gg.Point{X: 3, Y: 5}},
				Bottom: imgix.HorizontalEdge{LeftMost: &gg.Point{X: 0, Y: 1}, RightMost: &gg.Point{X: 3, Y: 1}},
				Left:   imgix.VerticalEdge{Nadir: &gg.Point{X: 0, Y: 1}, Vertex: &gg.Point{X: 0, Y: 5}},
				Right:  imgix.VerticalEdge{Nadir: &gg.Point{X: 3, Y: 1}, Vertex: &gg.Point{X: 3, Y: 5}},
			},
		},
		{
			name: "Single point case",
			points: map[string]gg.Point{
				"point1": {X: 2, Y: 3},
			},
			expected: imgix.Coordinates{
				Top:    imgix.HorizontalEdge{LeftMost: &gg.Point{X: 2, Y: 3}, RightMost: &gg.Point{X: 2, Y: 3}},
				Bottom: imgix.HorizontalEdge{LeftMost: &gg.Point{X: 2, Y: 3}, RightMost: &gg.Point{X: 2, Y: 3}},
				Left:   imgix.VerticalEdge{Nadir: &gg.Point{X: 2, Y: 3}, Vertex: &gg.Point{X: 2, Y: 3}},
				Right:  imgix.VerticalEdge{Nadir: &gg.Point{X: 2, Y: 3}, Vertex: &gg.Point{X: 2, Y: 3}},
			},
		},
		{
			name:     "Empty case",
			points:   map[string]gg.Point{},
			expected: imgix.Coordinates{},
		},
		{
			name: "Overlapping points",
			points: map[string]gg.Point{
				"point1": {X: 1, Y: 1},
				"point2": {X: 1, Y: 1}, // 重叠点
				"point3": {X: 2, Y: 2},
				"point4": {X: 3, Y: 3},
			},
			expected: imgix.Coordinates{
				Top:    imgix.HorizontalEdge{LeftMost: &gg.Point{X: 1, Y: 3}, RightMost: &gg.Point{X: 3, Y: 3}},
				Bottom: imgix.HorizontalEdge{LeftMost: &gg.Point{X: 1, Y: 1}, RightMost: &gg.Point{X: 3, Y: 1}},
				Left:   imgix.VerticalEdge{Nadir: &gg.Point{X: 1, Y: 1}, Vertex: &gg.Point{X: 1, Y: 3}},
				Right:  imgix.VerticalEdge{Nadir: &gg.Point{X: 3, Y: 1}, Vertex: &gg.Point{X: 3, Y: 3}},
			},
		},
		{
			name: "Negative coordinates",
			points: map[string]gg.Point{
				"point1": {X: -1, Y: -2},
				"point2": {X: -3, Y: -4},
				"point3": {X: -5, Y: -1},
				"point4": {X: -2, Y: -3},
			},
			expected: imgix.Coordinates{
				Top:    imgix.HorizontalEdge{LeftMost: &gg.Point{X: -5, Y: -1}, RightMost: &gg.Point{X: -1, Y: -1}},
				Bottom: imgix.HorizontalEdge{LeftMost: &gg.Point{X: -5, Y: -4}, RightMost: &gg.Point{X: -1, Y: -4}},
				Left:   imgix.VerticalEdge{Nadir: &gg.Point{X: -5, Y: -4}, Vertex: &gg.Point{X: -5, Y: -1}},
				Right:  imgix.VerticalEdge{Nadir: &gg.Point{X: -1, Y: -4}, Vertex: &gg.Point{X: -1, Y: -1}},
			},
		},
		{
			name: "Extreme coordinates",
			points: map[string]gg.Point{
				"point1": {X: 1000000, Y: 2000000},
				"point2": {X: -1000000, Y: -2000000},
				"point3": {X: 0, Y: 0},
			},
			expected: imgix.Coordinates{
				Top:    imgix.HorizontalEdge{LeftMost: &gg.Point{X: -1000000, Y: 2000000}, RightMost: &gg.Point{X: 1000000, Y: 2000000}},
				Bottom: imgix.HorizontalEdge{LeftMost: &gg.Point{X: -1000000, Y: -2000000}, RightMost: &gg.Point{X: 1000000, Y: -2000000}},
				Left:   imgix.VerticalEdge{Nadir: &gg.Point{X: -1000000, Y: -2000000}, Vertex: &gg.Point{X: -1000000, Y: 2000000}},
				Right:  imgix.VerticalEdge{Nadir: &gg.Point{X: 1000000, Y: -2000000}, Vertex: &gg.Point{X: 1000000, Y: 2000000}},
			},
		},
		{
			name: "Random distribution of points",
			points: map[string]gg.Point{
				"point1": {X: 5, Y: 7},
				"point2": {X: -3, Y: 2},
				"point3": {X: 1, Y: -1},
				"point4": {X: 4, Y: 6},
				"point5": {X: 0, Y: 0},
				"point6": {X: -2, Y: 3},
			},
			expected: imgix.Coordinates{
				Top:    imgix.HorizontalEdge{LeftMost: &gg.Point{X: -3, Y: 7}, RightMost: &gg.Point{X: 5, Y: 7}},
				Bottom: imgix.HorizontalEdge{LeftMost: &gg.Point{X: -3, Y: -1}, RightMost: &gg.Point{X: 5, Y: -1}},
				Left:   imgix.VerticalEdge{Nadir: &gg.Point{X: -3, Y: -1}, Vertex: &gg.Point{X: -3, Y: 7}},
				Right:  imgix.VerticalEdge{Nadir: &gg.Point{X: 5, Y: -1}, Vertex: &gg.Point{X: 5, Y: 7}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := imgix.CleanCoordinates(tt.points)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCalculatePoint 测试 CalculatePoint 函数
func TestCalculatePoint(t *testing.T) {
	tests := []struct {
		a        *gg.Point
		b        *gg.Point
		mode     imgix.CalculateMode
		axis     imgix.AxisPointMode
		expected *gg.Point
	}{
		{
			a:        &gg.Point{X: 1, Y: 2},
			b:        &gg.Point{X: 3, Y: 4},
			mode:     imgix.CalculateMax,
			axis:     imgix.AxisX,
			expected: &gg.Point{X: 3, Y: 4},
		},
		{
			a:        &gg.Point{X: 1, Y: 2},
			b:        &gg.Point{X: 3, Y: 4},
			mode:     imgix.CalculateMax,
			axis:     imgix.AxisY,
			expected: &gg.Point{X: 3, Y: 4},
		},
		{
			a:        &gg.Point{X: 1, Y: 2},
			b:        &gg.Point{X: 3, Y: 4},
			mode:     imgix.CalculateMin,
			axis:     imgix.AxisX,
			expected: &gg.Point{X: 1, Y: 2},
		},
		{
			a:        &gg.Point{X: 1, Y: 2},
			b:        &gg.Point{X: 3, Y: 4},
			mode:     imgix.CalculateMin,
			axis:     imgix.AxisY,
			expected: &gg.Point{X: 1, Y: 2},
		},
	}

	for _, test := range tests {
		result := imgix.CalculatePoint(test.a, test.b, test.mode, test.axis)
		assert.Equal(t, test.expected, result, "Expected result for input (%v, %v, %v, %v) does not match", test.a, test.b, test.mode, test.axis)
	}
}

// TestCalculateMultiplePoints 测试 CalculateMultiplePoints 函数
func TestCalculateMultiplePoints(t *testing.T) {
	points := []*gg.Point{
		{X: 1, Y: 2},
		{X: 3, Y: 4},
		{X: 5, Y: 1},
		{X: 0, Y: 6},
	}

	// 测试最大值
	expectedMax := &gg.Point{X: 5, Y: 1} // 在X轴上最大值
	resultMax := imgix.CalculateMultiplePoints(points, imgix.CalculateMax, imgix.AxisX)
	assert.Equal(t, expectedMax, resultMax, "Expected max point does not match")

	// 测试最小值
	expectedMin := &gg.Point{X: 0, Y: 6} // 在X轴上最小值
	resultMin := imgix.CalculateMultiplePoints(points, imgix.CalculateMin, imgix.AxisX)
	assert.Equal(t, expectedMin, resultMin, "Expected min point does not match")
}

func TestCanFormTriangle(t *testing.T) {
	tests := []struct {
		points []*gg.Point
		expect bool
	}{
		{[]*gg.Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}}, false}, // 共线
		{[]*gg.Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 0}}, true},  // 可以构成三角形
		{[]*gg.Point{{X: 1, Y: 1}, {X: 2, Y: 2}, {X: 3, Y: 3}}, false}, // 共线
		{[]*gg.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}}, true},  // 可以构成三角形
		{[]*gg.Point{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}}, false}, // 共线
		{[]*gg.Point{{X: 603, Y: 611}, {X: 571.9291914052573, Y: 499.9499975182803}, {X: 603, Y: 462}}, true},
		{[]*gg.Point{{X: 0, Y: 0}, {X: 1, Y: 1}}, false},                             // 不足三个点
		{[]*gg.Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 3, Y: 4}, {X: 5, Y: 6}}, false}, // 大于三个点
	}

	for _, test := range tests {
		area, result, err := imgix.CanFormTriangle(test.points)
		if test.expect {
			assert.NoError(t, err) // 期望没有错误
			assert.NotEqual(t, 0.0, area, "Area should not be zero for a valid triangle")
		}
		assert.Equal(t, test.expect, result, "Points: %+v, Expect: %v, Result: %v\n", test.points, test.expect, result)
	}
}
func TestResizeX(t *testing.T) {
	tests := []struct {
		point    *gg.Point
		resize   float64
		expected *gg.Point
	}{
		{&gg.Point{X: 2, Y: 3}, 2.0, &gg.Point{X: 4, Y: 3}},
		{&gg.Point{X: -1, Y: 5}, 3.0, &gg.Point{X: -3, Y: 5}},
		{&gg.Point{X: 0, Y: 0}, 10.0, &gg.Point{X: 0, Y: 0}},
	}

	for _, test := range tests {
		result := imgix.ResizeX(test.point, test.resize)
		assert.Equal(t, test.expected.X, result.X, "X坐标不匹配")
		assert.Equal(t, test.expected.Y, result.Y, "Y坐标应保持不变")
	}
}

func TestResizeY(t *testing.T) {
	tests := []struct {
		point    *gg.Point
		resize   float64
		expected *gg.Point
	}{
		{&gg.Point{X: 2, Y: 3}, 0.5, &gg.Point{X: 2, Y: 1.5}},
		{&gg.Point{X: -1, Y: 5}, 2.0, &gg.Point{X: -1, Y: 10}},
		{&gg.Point{X: 0, Y: 0}, 10.0, &gg.Point{X: 0, Y: 0}},
	}

	for _, test := range tests {
		result := imgix.ResizeY(test.point, test.resize)
		assert.Equal(t, test.expected.X, result.X, "X坐标应保持不变")
		assert.Equal(t, test.expected.Y, result.Y, "Y坐标不匹配")
	}
}

// TestResizePoint 测试 ResizePoint 函数
func TestResizePoint(t *testing.T) {
	tests := []struct {
		point    *gg.Point
		resize   float64
		expected *gg.Point
	}{
		{
			point:    &gg.Point{X: 1, Y: 2},
			resize:   0.5,
			expected: &gg.Point{X: 0.5, Y: 1},
		},
		{
			point:    &gg.Point{X: -1, Y: -2},
			resize:   2,
			expected: &gg.Point{X: -2, Y: -4},
		},
		{
			point:    &gg.Point{X: 0, Y: 0},
			resize:   1,
			expected: &gg.Point{X: 0, Y: 0}, // 不缩放
		},
		{
			point:    &gg.Point{X: 10, Y: 20},
			resize:   0,
			expected: &gg.Point{X: 0, Y: 0}, // 缩放到原点
		},
	}

	for _, test := range tests {
		scaledPoint := imgix.ResizePoint(test.point, test.resize, test.resize)
		assert.Equal(t, test.expected.X, scaledPoint.X, "X coordinate mismatch")
		assert.Equal(t, test.expected.Y, scaledPoint.Y, "Y coordinate mismatch")
	}
}

func TestResizePoints(t *testing.T) {
	tests := []struct {
		points   []*gg.Point
		resize   float64
		expected []*gg.Point
	}{
		{
			points:   []*gg.Point{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}},
			resize:   0.5,
			expected: []*gg.Point{{X: 0.5, Y: 1}, {X: 1.5, Y: 2}, {X: 2.5, Y: 3}},
		},
		{
			points:   []*gg.Point{{X: -1, Y: -2}, {X: -3, Y: -4}, {X: -5, Y: -6}},
			resize:   2,
			expected: []*gg.Point{{X: -2, Y: -4}, {X: -6, Y: -8}, {X: -10, Y: -12}},
		},
		{
			points:   []*gg.Point{{X: 0, Y: 0}, {X: 1, Y: 1}},
			resize:   1,
			expected: []*gg.Point{{X: 0, Y: 0}, {X: 1, Y: 1}}, // 不缩放
		},
		{
			points:   []*gg.Point{{X: 10, Y: 20}, {X: 30, Y: 40}},
			resize:   0,
			expected: []*gg.Point{{X: 0, Y: 0}, {X: 0, Y: 0}}, // 缩放到原点
		},
	}

	for _, test := range tests {
		scaledPoints := imgix.ResizePoints(test.points, test.resize, test.resize)
		assert.Equal(t, len(test.expected), len(scaledPoints)) // 确保长度相等

		for i, point := range scaledPoints {
			assert.Equal(t, test.expected[i].X, point.X, "X coordinate mismatch")
			assert.Equal(t, test.expected[i].Y, point.Y, "Y coordinate mismatch")
		}
	}
}

// TestResizePointOneselfX 测试 ResizePointOneselfX 函数
func TestResizePointOneselfX(t *testing.T) {
	point := &gg.Point{X: 10, Y: 20}

	tests := []struct {
		scaleFactor float64
		operation   imgix.CalculateFractionPointMode
		expectedX   float64
	}{
		{2.0, imgix.Add, 30.0},       // 10 + (10 * 2) = 30
		{0.5, imgix.Subtract, 5.0},   // 10 - (10 * 0.5) = 5
		{2.0, imgix.Multiply, 200.0}, // 10 * (10 * 2) = 200
		{2.0, imgix.Divide, 0.5},     // 10 / (10 * 2) = 0.5
		{1.0, imgix.Subtract, 0.0},   // 10 - (10 * 1) = 0
	}

	for _, test := range tests {
		result := imgix.ResizePointOneselfX(point, test.scaleFactor, test.operation)
		assert.Equal(t, test.expectedX, result.X, "Expected X value did not match")
	}
}

// TestResizePointOneselfY 测试 ResizePointOneselfY 函数
func TestResizePointOneselfY(t *testing.T) {
	point := &gg.Point{X: 10, Y: 20}

	tests := []struct {
		scaleFactor float64
		operation   imgix.CalculateFractionPointMode
		expectedY   float64
	}{
		{2.0, imgix.Add, 60.0},       // 20 + (20 * 2) = 60
		{0.5, imgix.Subtract, 10.0},  // 20 - (20 * 0.5) = 10
		{2.0, imgix.Multiply, 800.0}, // 20 * (20 * 2) = 800
		{2.0, imgix.Divide, 0.5},     // 20 / (20 * 2) = 0.5
		{1.0, imgix.Subtract, 0.0},   // 20 - (20 * 1) = 0
	}

	for _, test := range tests {
		result := imgix.ResizePointOneselfY(point, test.scaleFactor, test.operation)
		assert.Equal(t, test.expectedY, result.Y, "Expected Y value did not match")
	}
}

// TestOffsetPointX 测试 OffsetPointX 函数
func TestOffsetPointX(t *testing.T) {
	point := &gg.Point{X: 10, Y: 20}

	tests := []struct {
		offset    float64
		operation imgix.CalculateFractionPointMode
		expectedX float64
	}{
		{5.0, imgix.Add, 15.0},      // 10 + 5 = 15
		{3.0, imgix.Subtract, 7.0},  // 10 - 3 = 7
		{2.0, imgix.Multiply, 20.0}, // 10 * 2 = 20
		{2.0, imgix.Divide, 5.0},    // 10 / 2 = 5
		{0.0, imgix.Subtract, 10.0}, // 10 - 0 = 10
	}

	for _, test := range tests {
		result := imgix.OffsetPointX(point, test.offset, test.operation)
		assert.Equal(t, test.expectedX, result.X, "Expected X value did not match")
	}
}

// TestOffsetPointY 测试 OffsetPointY 函数
func TestOffsetPointY(t *testing.T) {
	point := &gg.Point{X: 10, Y: 20}

	tests := []struct {
		offset    float64
		operation imgix.CalculateFractionPointMode
		expectedY float64
	}{
		{5.0, imgix.Add, 25.0},      // 20 + 5 = 25
		{3.0, imgix.Subtract, 17.0}, // 20 - 3 = 17
		{2.0, imgix.Multiply, 40.0}, // 20 * 2 = 40
		{2.0, imgix.Divide, 10.0},   // 20 / 2 = 10
		{0.0, imgix.Subtract, 20.0}, // 20 - 0 = 20
	}

	for _, test := range tests {
		result := imgix.OffsetPointY(point, test.offset, test.operation)
		assert.Equal(t, test.expectedY, result.Y, "Expected Y value did not match")
	}
}

// TestResizeUpTriangle 测试 ResizeUpTriangle 函数
func TestResizeUpTriangle(t *testing.T) {
	tests := []struct {
		vertexA   *gg.Point
		vertexB   *gg.Point
		vertexC   *gg.Point
		resize    float64
		expectedB *gg.Point
		expectedC *gg.Point
	}{
		{
			vertexA:   &gg.Point{X: 0, Y: 0},
			vertexB:   &gg.Point{X: 2, Y: 0},
			vertexC:   &gg.Point{X: 1, Y: 2},
			resize:    2.0,
			expectedB: &gg.Point{X: 4, Y: 0},
			expectedC: &gg.Point{X: 2, Y: 4},
		},
	}

	for _, test := range tests {
		newVertexB, newVertexC := imgix.ResizeUpTriangle(test.vertexA, test.vertexB, test.vertexC, test.resize)

		assert.Equal(t, test.expectedB.X, newVertexB.X, "vertexB X坐标不匹配")
		assert.Equal(t, test.expectedB.Y, newVertexB.Y, "vertexB Y坐标不匹配")
		assert.Equal(t, test.expectedC.X, newVertexC.X, "vertexC X坐标不匹配")
		assert.Equal(t, test.expectedC.Y, newVertexC.Y, "vertexC Y坐标不匹配")
	}
}

// TestResizeDownTriangle 测试 ResizeDownTriangle 函数
func TestResizeDownTriangle(t *testing.T) {
	tests := []struct {
		vertexA   *gg.Point
		vertexB   *gg.Point
		vertexC   *gg.Point
		resize    float64
		expectedB *gg.Point
		expectedC *gg.Point
	}{
		{
			vertexA:   &gg.Point{X: 0, Y: 0},
			vertexB:   &gg.Point{X: 4, Y: 0},
			vertexC:   &gg.Point{X: 0, Y: 4},
			resize:    0.5,
			expectedB: &gg.Point{X: 2, Y: 0},
			expectedC: &gg.Point{X: 0, Y: 2},
		},
		{
			vertexA:   &gg.Point{X: 1, Y: 1},
			vertexB:   &gg.Point{X: 3, Y: 1},
			vertexC:   &gg.Point{X: 2, Y: 3},
			resize:    0.75,
			expectedB: &gg.Point{X: 2.5, Y: 1},
			expectedC: &gg.Point{X: 1.75, Y: 2.5},
		},
		{
			vertexA:   &gg.Point{X: 0, Y: 0},
			vertexB:   &gg.Point{X: 4, Y: 0},
			vertexC:   &gg.Point{X: 0, Y: 4},
			resize:    0.1,
			expectedB: &gg.Point{X: 0.4, Y: 0},
			expectedC: &gg.Point{X: 0, Y: 0.4},
		},
		{
			vertexA:   &gg.Point{X: 0, Y: 0},
			vertexB:   &gg.Point{X: 4, Y: 0},
			vertexC:   &gg.Point{X: 0, Y: 4},
			resize:    1.5, // 超出范围，返回原始坐标
			expectedB: &gg.Point{X: 4, Y: 0},
			expectedC: &gg.Point{X: 0, Y: 4},
		},
	}

	for _, test := range tests {
		newVertexB, newVertexC := imgix.ResizeDownTriangle(test.vertexA, test.vertexB, test.vertexC, test.resize)

		assert.Equal(t, test.expectedB.X, newVertexB.X, "vertexB X坐标不匹配")
		assert.Equal(t, test.expectedB.Y, newVertexB.Y, "vertexB Y坐标不匹配")
		assert.Equal(t, test.expectedC.X, newVertexC.X, "vertexC X坐标不匹配")
		assert.Equal(t, test.expectedC.Y, newVertexC.Y, "vertexC Y坐标不匹配")
	}
}

// TestExtendLine 测试 ExtendLine 函数
func TestExtendLine(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		p1     gg.Point
		p2     gg.Point
		length float64
		expect gg.Point
	}{
		{gg.Point{X: 200, Y: 100}, gg.Point{X: 300, Y: 50}, 50, gg.Point{X: 244.72135954999578, Y: 77.63932022500211}},
		{gg.Point{X: 0, Y: 0}, gg.Point{X: 3, Y: 4}, 5, gg.Point{X: 3.0, Y: 4.0}}, // 3-4-5 三角形
		{gg.Point{X: 1, Y: 1}, gg.Point{X: 4, Y: 5}, 0, gg.Point{X: 1.0, Y: 1.0}}, // 不延长
	}

	// 执行测试
	for _, test := range tests {
		result := imgix.ExtendLine(&test.p1, &test.p2, test.length)
		if result.X != test.expect.X || result.Y != test.expect.Y {
			t.Errorf("ExtendLine(%v, %v, %v) = %v; want %v", test.p1, test.p2, test.length, result, test.expect)
		}
	}
}

// 测试 ResizeImage 函数
func TestResizeImage(t *testing.T) {
	// 创建一个简单的测试图片（100x100 的红色方块）
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			testImg.Set(x, y, color.RGBA{255, 0, 0, 255}) // 红色
		}
	}

	// 定义缩放选项
	resizeOptions := &imgix.ResizeImgOptions{
		Width:  50,
		Height: 50,
		Filter: imaging.Lanczos,
	}

	// 调用 ResizeImage 函数
	resizedImg := imgix.ResizeImage(testImg, resizeOptions)

	// 断言结果
	assert.NotNil(t, resizedImg, "Resized image should not be nil")
	assert.Equal(t, 50, resizedImg.Bounds().Dx(), "Resized image width should be 50")
	assert.Equal(t, 50, resizedImg.Bounds().Dy(), "Resized image height should be 50")

	// 可选：将结果保存到文件以便手动检查
	outFile, err := os.Create("test_resize_mage.png")
	assert.NoError(t, err, "Error creating output file")
	defer outFile.Close()

	err = jpeg.Encode(outFile, resizedImg, nil)
	assert.NoError(t, err, "Error encoding image to JPEG")
}

// TestCropImage 测试 CropImage 函数
func TestCropImage(t *testing.T) {
	// 创建并保存测试图像
	filePath := "test_crop_image.png"
	err := createTestImage(filePath)
	if err != nil {
		t.Fatalf("创建测试图像失败: %v", err)
	}
	defer os.Remove(filePath) // 测试完成后删除文件

	// 读取生成的图像
	testImg, err := imaging.Open(filePath)
	if err != nil {
		t.Fatalf("打开测试图像失败: %v", err)
	}

	// 定义裁剪选项
	cropOptions := &imgix.CropImgOptions{
		MinWidth:  10,
		MinHeight: 10,
		MaxWidth:  50,
		MaxHeight: 50,
	}

	// 调用 CropImage 函数
	croppedImg := imgix.CropImage(testImg, cropOptions)

	// 断言裁剪后的图像尺寸
	assert.Equal(t, 40, croppedImg.Bounds().Dx(), "裁剪后的宽度应为 40")
	assert.Equal(t, 40, croppedImg.Bounds().Dy(), "裁剪后的高度应为 40")
}

func TestAdjustValues(t *testing.T) {
	tests := []struct {
		start     float64
		end       float64
		target    float64
		wantStart float64
		wantEnd   float64
	}{
		// 测试用例 1: 正常情况
		{2.0, 3.0, 5.0, 0.0, 5.0}, // 输入: start=2.0, end=3.0, target=5.0
		// 计算: diff = 3.0 - 2.0 = 1.0
		// 补充差值: 5.0 - 1.0 = 4.0
		// 增量: 4.0 / 2 = 2.0
		// 更新: start = 2.0 - 2.0 = 0.0, end = 3.0 + 2.0 = 5.0

		// 测试用例 2: 起始和结束相等
		{3.0, 3.0, 5.0, 0.5, 5.5}, // 输入: start=3.0, end=3.0, target=5.0
		// 计算: diff = 3.0 - 3.0 = 0.0
		// 补充差值: 5.0 - 0.0 = 5.0
		// 增量: 5.0 / 2 = 2.5
		// 更新: start = 3.0 - 2.5 = 0.5, end = 3.0 + 2.5 = 5.5

		// 测试用例 3: 差值大于目标值
		{2.0, 7.0, 5.0, 2.0, 7.0}, // 输入: start=2.0, end=7.0, target=5.0
		// 计算: diff = 7.0 - 2.0 = 5.0
		// 补充差值: 5.0 - 5.0 = 0.0
		// 不需要调整: start = 2.0, end = 7.0

		// 测试用例 4: 应该颠倒
		{8.0, 2.0, 5.0, 2.0, 8.0}, // 输入: start=8.0, end=2.0, target=5.0
		// 颠倒: start = 2.0, end = 8.0
		// 计算: diff = 8.0 - 2.0 = 6.0
		// 补充差值: 5.0 - 6.0 = -1.0
		// 不需要调整: start = 2.0, end = 8.0

		// 测试用例 5: 颠倒并调整
		{8.0, 2.0, 9.0, 0.5, 9.5}, // 输入: start=8.0, end=2.0, target=9.0
		// 颠倒: start = 2.0, end = 8.0
		// 计算: diff = 8.0 - 2.0 = 6.0
		// 补充差值: 9.0 - 6.0 = 3.0
		// 增量: 3.0 / 2 = 1.5
		// 更新: start = 2.0 - 1.5 = 0.5, end = 8.0 + 1.5 = 9.5
	}

	for _, test := range tests {
		start, end := imgix.AdjustValues(test.start, test.end, test.target)
		assert.Equal(t, test.wantStart, start, "起始值不匹配")
		assert.Equal(t, test.wantEnd, end, "结束值不匹配")
	}
}
