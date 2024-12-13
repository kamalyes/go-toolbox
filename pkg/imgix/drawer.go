/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 09:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-13 13:05:56
 * @FilePath: \go-toolbox\pkg\imgix\drawer.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package imgix

import (
	"log"
	"math"
	"sync"

	"github.com/fogleman/gg"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// GraphicsRenderer 结构体用于绘制面部特征
type GraphicsRenderer struct {
	GgCtx       *gg.Context
	DashOptions DashOptions
	mu          sync.RWMutex // 读写锁
}

type DashStyle float64

type DashOptions struct {
	dashLength DashStyle // 虚线段长度
	gapLength  DashStyle // 虚线间隔长度
}

// DashLength 返回虚线段长度的 float64 表示
func (d *DashOptions) DashLength() DashStyle {
	return d.dashLength
}

// GapLength 返回虚线间隔长度的 float64 表示
func (d *DashOptions) GapLength() DashStyle {
	return d.gapLength
}

// 获取 DashOptions 的方法
func (g *GraphicsRenderer) GetDashOptions() DashOptions {
	return g.DashOptions
}

// ImageLTRB 表示图像的边界框
type ImageLTRB struct {
	Left   float64 // 左边界
	Top    float64 // 上边界
	Right  float64 // 右边界
	Bottom float64 // 下边界
}

// NewGraphicsRenderer 创建一个新的 GraphicsRenderer 实例
// @param ctx gg.Context 的指针，用于绘制操作
// @param dashOptions
// @return *GraphicsRenderer 返回一个新的 GraphicsRenderer 实例
func NewGraphicsRenderer(ctx *gg.Context, dashOptions ...DashOptions) *GraphicsRenderer {
	defaultDashOptions := DashOptions{
		dashLength: 3,
		gapLength:  6,
	}
	if len(dashOptions) > 0 {
		defaultDashOptions = dashOptions[0]
	}
	renderer := GraphicsRenderer{
		GgCtx:       ctx,
		DashOptions: defaultDashOptions,
	}
	return &renderer
}

// drawWithStroke 是一个通用的绘图函数，接受一个绘图操作的函数作为参数
func (g *GraphicsRenderer) drawWithStroke(drawFunc func(), isStroke bool) {
	syncx.WithLock(&g.mu, func() {
		drawFunc() // 执行绘图操作
		if isStroke {
			g.GgCtx.Stroke() // 在绘图操作完成后统一调用 Stroke
		}
	})
}

// UseDefaultDashed 使用默认虚线
func (g *GraphicsRenderer) UseDefaultDashed() {
	g.drawWithStroke(func() {
		g.GgCtx.SetDash(float64(g.DashOptions.DashLength()), float64(g.DashOptions.GapLength()), 0.0)
	}, false)
}

// UseSolidLine 使用实线
func (g *GraphicsRenderer) UseSolidLine() {
	g.drawWithStroke(func() {
		g.GgCtx.SetDash() // 恢复为实线样式
	}, false)
}

// SetDashed 设置是否使用虚线
// @param dashes
func (g *GraphicsRenderer) SetDashed(dashes ...float64) {
	g.drawWithStroke(func() {
		g.GgCtx.SetDash(dashes...)
	}, false)
}

// DrawLine 绘制线条
// @param startX 起始X坐标
// @param startY 起始Y坐标
// @param endX 结束X坐标
// @param endY 结束Y坐标
// @param lineWidth 线条宽度
func (g *GraphicsRenderer) DrawLine(startX, startY, endX, endY float64, lineWidth float64) {
	g.drawWithStroke(func() {
		// 设置线条宽度
		g.GgCtx.SetLineWidth(lineWidth)
		// 绘制线条
		g.GgCtx.DrawLine(startX, startY, endX, endY)
	}, true)

}

// CalculateFractionPoint 计算任意分数的坐标
// @param startPoint 起始点
// @param endPoint 结束点
// @param fraction 计算的分数（0到1之间）
// @return *gg.Point 返回计算得到的坐标
func (g *GraphicsRenderer) CalculateFractionPoint(startPoint, endPoint *gg.Point, fraction float64) *gg.Point {
	return &gg.Point{
		X: startPoint.X + (endPoint.X-startPoint.X)*fraction,
		Y: startPoint.Y + (endPoint.Y-startPoint.Y)*fraction,
	}
}

// DrawCurvedLine 绘制带有弯曲度的线条
// @param start 起始点
// @param end 结束点
// @param control 控制点
func (g *GraphicsRenderer) DrawCurvedLine(start, end, control *gg.Point) {
	g.drawWithStroke(func() {
		g.GgCtx.MoveTo(start.X, start.Y)
		g.GgCtx.QuadraticTo(control.X, control.Y, end.X, end.Y)
	}, true)
}

// CalculateControlPoint 计算控制点以实现弯曲效果
// @param start 起始点
// @param end 结束点
// @param offset 偏移量
// @return *gg.Point 返回计算得到的控制点
func (g *GraphicsRenderer) CalculateControlPoint(start, end *gg.Point, offset float64) *gg.Point {
	midX := (start.X + end.X) / 2
	midY := (start.Y + end.Y) / 2
	return &gg.Point{X: midX, Y: midY - offset}
}

// DrawRectangle 绘制矩形框
// @param left 矩形左上角的点
// @param top 矩形左上角的点
// @param bottom 矩形右下角的点
// @param right 矩形右下角的点
func (g *GraphicsRenderer) DrawRectangle(left, top, bottom, right *gg.Point) {
	g.drawWithStroke(func() {
		width := right.X - left.X
		height := bottom.Y - top.Y
		g.GgCtx.DrawRectangle(left.X, top.Y, width, height)
	}, true)
}

// DrawPolygon 绘制一个多边形
// @param points 多边形顶点的切片
func (g *GraphicsRenderer) DrawPolygon(points []gg.Point) {
	g.drawWithStroke(func() {
		if len(points) > 0 {
			g.GgCtx.MoveTo(points[0].X, points[0].Y)
			for _, point := range points[1:] {
				g.GgCtx.LineTo(point.X, point.Y)
			}
			g.GgCtx.ClosePath()
		}
	}, true)
}

// DrawVerticalLine 从面部顶部到底部绘制竖线
// @param x 竖线的X坐标
// @param top 竖线的顶部Y坐标
// @param bottom 竖线的底部Y坐标
func (g *GraphicsRenderer) DrawVerticalLine(x float64, top, bottom float64) {
	g.drawWithStroke(func() {
		g.GgCtx.DrawLine(x, top, x, bottom)
	}, true)
}

// DrawHorizontalLine 从面部左侧到右侧绘制横线
// @param y 横线的Y坐标
// @param left 横线的左侧X坐标
// @param right 横线的右侧X坐标
func (g *GraphicsRenderer) DrawHorizontalLine(y float64, left, right float64) {
	g.drawWithStroke(func() {
		g.GgCtx.DrawLine(left, y, right, y)
	}, true)
}

// DrawCircle 绘制一个圆形
// @param centerX 圆心的X坐标
// @param centerY 圆心的Y坐标
// @param radius 圆的半径
func (g *GraphicsRenderer) DrawCircle(centerX, centerY, radius float64) {
	g.drawWithStroke(func() {
		g.GgCtx.DrawCircle(centerX, centerY, radius)
	}, true)
}

// DrawEllipse 绘制一个椭圆
// @param centerX 椭圆中心的X坐标
// @param centerY 椭圆中心的Y坐标
// @param width 椭圆的宽度
// @param height 椭圆的高度
func (g *GraphicsRenderer) DrawEllipse(centerX, centerY, width, height float64) {
	g.drawWithStroke(func() {
		g.GgCtx.DrawEllipse(centerX, centerY, width/2, height/2)
	}, true)
}

// DrawTriangle 绘制一个三角形
// @param x1 第一个顶点的X坐标
// @param y1 第一个顶点的Y坐标
// @param x2 第二个顶点的X坐标
// @param y2 第二个顶点的Y坐标
// @param x3 第三个顶点的X坐标
// @param y3 第三个顶点的Y坐标
func (g *GraphicsRenderer) DrawTriangle(x1, y1, x2, y2, x3, y3 float64) {
	g.drawWithStroke(func() {
		g.GgCtx.MoveTo(x1, y1) // 移动到第一个顶点
		g.GgCtx.LineTo(x2, y2) // 连接到第二个顶点
		g.GgCtx.LineTo(x3, y3) // 连接到第三个顶点
		g.GgCtx.ClosePath()    // 关闭路径形成三角形
	}, true) // 默认调用 Stroke
}

// DrawCenteredMultiLine 在多个指定的坐标区间内绘制多组文本
// @param startXs 文本起始X坐标数组
// @param endXs 文本结束X坐标数组
// @param startYs 文本起始Y坐标数组
// @param endYs 文本结束Y坐标数组
// @param textGroups 文本内容的二维数组
// @param angle 旋转角度
// @param drawDashed 是否绘制虚线
// @param lineSpacing 行间距（可选）
func (g *GraphicsRenderer) DrawCenteredMultiLine(startXs, endXs, startYs, endYs []float64, textGroups [][]string, angle float64, drawDashed bool, lineSpacing ...float64) {
	log.Println("Starting DrawCenteredMultiLine")
	if len(startXs) != len(textGroups) ||
		len(endXs) != len(textGroups) ||
		len(startYs) != len(textGroups) ||
		len(endYs) != len(textGroups) {
		log.Printf("Length mismatch: startXs: %d, endXs: %d, startYs: %d, endYs: %d, textGroups: %d", len(startXs), len(endXs), len(startYs), len(endYs), len(textGroups))
		return
	}

	defaultLineSpacing := 0.0
	if len(lineSpacing) > 0 {
		defaultLineSpacing = lineSpacing[0]
	}

	for i, texts := range textGroups {
		log.Printf("Drawing text group %d", i)
		startX, endX := startXs[i], endXs[i]
		startY, endY := startYs[i], endYs[i]
		midY := (startY + endY) / 2

		if drawDashed {
			g.DrawLine(startX, midY, endX, midY, 0.5)
		}

		totalHeight := 0.0
		lineHeights := make([]float64, len(texts))
		lineSpacings := make([]float64, len(texts))

		for j, text := range texts {
			_, height := g.GgCtx.MeasureString(text)
			lineHeights[j] = height

			if j < len(texts)-1 {
				if defaultLineSpacing > 0 {
					lineSpacings[j] = defaultLineSpacing
				} else {
					lineSpacings[j] = height * 0.5
				}
			}

			totalHeight += height + lineSpacings[j]
		}

		posX := (startX + endX) / 2
		posYStart := (startY + endY - totalHeight) / 2

		for j, text := range texts {
			posY := posYStart + lineHeights[j]/2
			adjustedPosY := posY
			if j == 0 {
				adjustedPosY += 50
				posYStart += 50
			}

			g.GgCtx.Push()
			g.GgCtx.Translate(posX, adjustedPosY)
			g.GgCtx.Rotate(angle * (math.Pi / 180))
			g.GgCtx.DrawStringAnchored(text, 0, 0, 0.5, 0.5)
			g.GgCtx.Pop()

			posYStart += lineHeights[j] + lineSpacings[j]
		}
	}
	log.Println("Finished DrawCenteredMultiLine")
}

// GetLTRB 获取关键点各部分的点横纵坐标最大小值，即：left，top，right，bottom
// @param features 面部特征点的映射
// @return ImageLTRB 返回包含边界信息的结构体
func (g *GraphicsRenderer) GetLTRB(features map[string]gg.Point) ImageLTRB {
	defaultCoordinates := 0.0
	maxFloat := float64(^uint(0) >> 1) // 最大的 float64 值
	minFloat := -maxFloat              // 最小的 float64 值

	left, top, right, bottom := maxFloat, maxFloat, minFloat, minFloat

	for _, coord := range features {
		g.UpdateBounds(coord.X, coord.Y, &left, &top, &right, &bottom)
	}

	// 使用默认值替代未更新的边界值
	if left == maxFloat {
		left = defaultCoordinates
	}
	if top == maxFloat {
		top = defaultCoordinates
	}
	if right == minFloat {
		right = defaultCoordinates
	}
	if bottom == minFloat {
		bottom = defaultCoordinates
	}

	return ImageLTRB{Left: left, Top: top, Right: right, Bottom: bottom}
}

// GetFacialXYByKey 按点的名称获取关键点的横纵坐标
// @param features 面部特征点的映射
// @param keyP 要查找的特征点名称
// @return *gg.Point 返回对应的坐标，如果未找到则返回 nil
func (g *GraphicsRenderer) GetFacialXYByKey(features map[string]gg.Point, keyP string) *gg.Point {
	if coord, exists := features[keyP]; exists {
		return &coord // 直接返回 gg.Point
	}
	return nil // 如果未找到，返回 nil
}

// UpdateBounds 更新边界值
// @param x 当前坐标的X值
// @param y 当前坐标的Y值
// @param left 左边界的指针
// @param top 上边界的指针
// @param right 右边界的指针
// @param bottom 下边界的指针
func (g *GraphicsRenderer) UpdateBounds(x, y float64, left, top, right, bottom *float64) {
	if x < *left {
		*left = x
	}
	if y < *top {
		*top = y
	}
	if x > *right {
		*right = x
	}
	if y > *bottom {
		*bottom = y
	}
}
