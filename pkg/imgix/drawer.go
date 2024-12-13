/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 09:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-15 09:09:15
 * @FilePath: \go-toolbox\pkg\imgix\drawer.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package imgix

import (
	"bytes"
	"fmt"
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
	bufferPool  *sync.Pool   // bufferPool 指针
	mu          sync.RWMutex // 读写锁
}

type DashStyle float64

type DashOptions struct {
	dashLength DashStyle // 虚线段长度
	gapLength  DashStyle // 虚线间隔长度
	lineWidth  DashStyle // 线条宽度
}

const (
	defaultDashLength = 3
	defaultGapLength  = 6
	defaultLineWidth  = 2
)

// NewDashOptions 创建一个新的 DashOptions 实例
// 如果提供了自定义的段长度和间隔长度，则使用它们；否则使用默认值
func NewDashOptions(dashLength, gapLength, lineWidth DashStyle) DashOptions {
	// 默认值
	if dashLength <= 0 {
		dashLength = defaultDashLength // 默认虚线段长度
	}
	if gapLength <= 0 {
		gapLength = defaultGapLength // 默认虚线间隔长度
	}
	if lineWidth <= 0 {
		lineWidth = defaultLineWidth // 默认线条宽度
	}

	return DashOptions{
		dashLength: dashLength,
		gapLength:  gapLength,
		lineWidth:  lineWidth,
	}
}

// DashLength 返回虚线段长度的 float64 表示
func (d *DashOptions) DashLength() DashStyle {
	return d.dashLength
}

// GapLength 返回虚线间隔长度的 float64 表示
func (d *DashOptions) GapLength() DashStyle {
	return d.gapLength
}

// 获取 LineWidth 的方法
func (g *DashOptions) LineWidth() DashStyle {
	return g.lineWidth
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

// Width 计算边界框的宽度
func (ltrb ImageLTRB) Width() float64 {
	return ltrb.Right - ltrb.Left
}

// Height 计算边界框的高度
func (ltrb ImageLTRB) Height() float64 {
	return ltrb.Bottom - ltrb.Top
}

// Center 返回边界框的中心点
func (ltrb ImageLTRB) Center() (float64, float64) {
	centerX := ltrb.Left + ltrb.Width()/2
	centerY := ltrb.Top + ltrb.Height()/2
	return centerX, centerY
}

// Contains 检查一个点是否在边界框内
func (ltrb ImageLTRB) Contains(x, y float64) bool {
	return x >= ltrb.Left && x <= ltrb.Right && y >= ltrb.Top && y <= ltrb.Bottom
}

// String 返回边界框的字符串表示
func (ltrb ImageLTRB) String() string {
	return fmt.Sprintf("ImageLTRB(Left: %f, Top: %f, Right: %f, Bottom: %f)", ltrb.Left, ltrb.Top, ltrb.Right, ltrb.Bottom)
}

// NewGraphicsRenderer 创建一个新的 GraphicsRenderer 实例
// @param ctx gg.Context 的指针，用于绘制操作
// @param dashOptions
// @return *GraphicsRenderer 返回一个新的 GraphicsRenderer 实例
func NewGraphicsRenderer(ctx *gg.Context, dashOptions ...DashOptions) *GraphicsRenderer {
	defaultDashOptions := NewDashOptions(defaultDashLength, defaultGapLength, defaultLineWidth)
	if len(dashOptions) > 0 {
		defaultDashOptions = dashOptions[0]
	}
	// 初始化 bufferPool
	bufferPool := &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	renderer := GraphicsRenderer{
		GgCtx:       ctx,
		DashOptions: defaultDashOptions,
		bufferPool:  bufferPool,
	}
	return &renderer
}

// DrawWithStroke 是一个通用的绘图函数，接受一个绘图操作的函数作为参数
func (g *GraphicsRenderer) DrawWithStroke(drawFunc func(), isStroke bool) {
	syncx.WithLock(&g.mu, func() {
		drawFunc() // 执行绘图操作
		if isStroke {
			g.GgCtx.Stroke() // 在绘图操作完成后统一调用 Stroke
		}
	})
}

// UseDefaultDashed 使用默认虚线
func (g *GraphicsRenderer) UseDefaultDashed() {
	g.DrawWithStroke(func() {
		g.GgCtx.SetDash(float64(g.DashOptions.DashLength()), float64(g.DashOptions.GapLength()), 0.0)
	}, false)
}

// UseSolidLine 使用实线
func (g *GraphicsRenderer) UseSolidLine() {
	g.DrawWithStroke(func() {
		g.GgCtx.SetDash() // 恢复为实线样式
	}, false)
}

// SetDashed 设置是否使用虚线
// @param dashes
func (g *GraphicsRenderer) SetDashed(dashes ...float64) {
	g.DrawWithStroke(func() {
		g.GgCtx.SetDash(dashes...)
	}, false)
}

// SaveImage 存储图片
// @param name 文件名名称
// @param quality 图片质量
// @param imageFormat 图片类型 (目前仅支持PNG/JPG/JPEG)
func (g *GraphicsRenderer) SaveImage(name string, quality int, imageFormat ImageFormat) {
	// 创建图像缓冲区 从池中获取一个 bytes.Buffer
	imgBuffer := g.bufferPool.Get().(*bytes.Buffer)
	defer g.bufferPool.Put(imgBuffer) // 使用完后放回池中

	// 清空缓冲区
	imgBuffer.Reset()
	WriterImage(g.GgCtx.Image(), quality, imageFormat, imgBuffer)
	// 保存生成的图像，使用策略名称
	imageFileName := fmt.Sprintf("%s.%s", name, imageFormat.String())
	SaveBufToImageFile(imgBuffer, imageFileName, imageFormat)
}

// DrawLineXYLineWidth 绘制线条
// @param startX 起始X坐标
// @param startY 起始Y坐标
// @param endX 结束X坐标
// @param endY 结束Y坐标
// @param lineWidth 线条宽度
func (g *GraphicsRenderer) DrawLineXYLineWidth(startX, startY, endX, endY float64) {
	g.DrawWithStroke(func() {
		g.GgCtx.DrawLine(startX, startY, endX, endY)
	}, true)

}

// DrawCurvedLine 绘制带有弯曲度的线条
// @param start 起始点
// @param end 结束点
// @param control 控制点
func (g *GraphicsRenderer) DrawCurvedLine(start, end, control *gg.Point) {
	g.DrawWithStroke(func() {
		g.GgCtx.MoveTo(start.X, start.Y)
		g.GgCtx.QuadraticTo(control.X, control.Y, end.X, end.Y)
	}, true)
}

// DrawRectangle 绘制矩形框
// @param left 矩形左上角的点
// @param top 矩形左上角的点
// @param bottom 矩形右下角的点
// @param right 矩形右下角的点
func (g *GraphicsRenderer) DrawRectangle(left, top, bottom, right *gg.Point) {
	g.DrawWithStroke(func() {
		width := right.X - left.X
		height := bottom.Y - top.Y
		g.GgCtx.DrawRectangle(left.X, top.Y, width, height)
	}, true)
}

// DrawPolygon 绘制一个多边形
// @param points 多边形顶点的切片
func (g *GraphicsRenderer) DrawPolygon(points []gg.Point) {
	g.DrawWithStroke(func() {
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
	g.DrawWithStroke(func() {
		g.GgCtx.DrawLine(x, top, x, bottom)
	}, true)
}

// DrawHorizontalLine 从面部左侧到右侧绘制横线
// @param y 横线的Y坐标
// @param left 横线的左侧X坐标
// @param right 横线的右侧X坐标
func (g *GraphicsRenderer) DrawHorizontalLine(y float64, left, right float64) {
	g.DrawWithStroke(func() {
		g.GgCtx.DrawLine(left, y, right, y)
	}, true)
}

// DrawLine 绘制从 startPoint 到 endPoint 的线条。
// startPoint 和 endPoint: 线条的起始和结束点。
func (g *GraphicsRenderer) DrawLine(startPoint, endPoint *gg.Point) {
	g.DrawWithStroke(func() {
		g.GgCtx.DrawLine(startPoint.X, startPoint.Y, endPoint.X, endPoint.Y)
	}, true)
}

// DrawCircle 绘制一个圆形
// @param centerX 圆心的X坐标
// @param centerY 圆心的Y坐标
// @param radius 圆的半径
func (g *GraphicsRenderer) DrawCircle(centerX, centerY, radius float64) {
	g.DrawWithStroke(func() {
		g.GgCtx.DrawCircle(centerX, centerY, radius)
	}, true)
}

// DrawArc 绘制一条弧线
// @param startPoint 起始点的X坐标/起始点的Y坐标
// @param endPoint 结束点的X坐标/结束点的Y坐标
// @param radius 圆的半径
// @param angleStart 起始角度（以度为单位）
// @param angleExtent 角度范围（以度为单位）
func (g *GraphicsRenderer) DrawArc(startPoint, endPoint *gg.Point, radius float64, angleStart, angleExtent float64) {
	g.DrawWithStroke(func() {
		g.GgCtx.DrawArc(startPoint.X, startPoint.Y, radius, angleStart*(math.Pi/180), angleExtent*(math.Pi/180))
		g.GgCtx.LineTo(endPoint.X, endPoint.Y) // 连接到结束点
	}, true)
}

// DrawEllipse 绘制一个椭圆
// @param centerX 椭圆中心的X坐标
// @param centerY 椭圆中心的Y坐标
// @param width 椭圆的宽度
// @param height 椭圆的高度
func (g *GraphicsRenderer) DrawEllipse(centerX, centerY, width, height float64) {
	g.DrawWithStroke(func() {
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
	g.DrawWithStroke(func() {
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
			g.DrawLineXYLineWidth(startX, midY, endX, midY)
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

// VerticalEdge 结构体表示一个垂直边，包括顶点和底点
type VerticalEdge struct {
	Nadir  *gg.Point // 底点，表示边的下端点
	Vertex *gg.Point // 顶点，表示边的上端点
}

// HorizontalEdge 结构体表示一个水平边，包括最左边和最右边的点
type HorizontalEdge struct {
	LeftMost  *gg.Point // 最左边的点，表示边的左端点
	RightMost *gg.Point // 最右边的点，表示边的右端点
}

// Coordinates 结构体表示边界框的坐标，包含四条边
type Coordinates struct {
	Top    HorizontalEdge // 上边，包含最左边和最右边的点
	Bottom HorizontalEdge // 下边，包含最左边和最右边的点
	Left   VerticalEdge   // 左边，包含顶点和底点
	Right  VerticalEdge   // 右边，包含顶点和底点
}

// CleanCoordinates 函数接受一个 map[string]Point，返回清洗后的边界框坐标
func CleanCoordinates(pointsMap map[string]gg.Point) Coordinates {
	if len(pointsMap) == 0 {
		return Coordinates{}
	}

	// 初始化最小值和最大值
	minX := math.Inf(1)
	maxX := math.Inf(-1)
	minY := math.Inf(1)
	maxY := math.Inf(-1)

	// 遍历所有点，找到最小和最大坐标
	for _, p := range pointsMap {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// 创建边界框的四条边
	top := HorizontalEdge{
		LeftMost:  &gg.Point{X: minX, Y: maxY},
		RightMost: &gg.Point{X: maxX, Y: maxY},
	}
	bottom := HorizontalEdge{
		LeftMost:  &gg.Point{X: minX, Y: minY},
		RightMost: &gg.Point{X: maxX, Y: minY},
	}
	left := VerticalEdge{
		Nadir:  &gg.Point{X: minX, Y: minY},
		Vertex: &gg.Point{X: minX, Y: maxY},
	}
	right := VerticalEdge{
		Nadir:  &gg.Point{X: maxX, Y: minY},
		Vertex: &gg.Point{X: maxX, Y: maxY},
	}

	// 返回构建的 Coordinates
	return Coordinates{
		Top:    top,
		Bottom: bottom,
		Left:   left,
		Right:  right,
	}
}

// GetLTRB 获取关键点各部分的点横纵坐标最大小值，即：left，top，right，bottom
// @param features 面部特征点的映射
// @return ImageLTRB 返回包含边界信息的结构体
func GetLTRB(features map[string]gg.Point) ImageLTRB {
	defaultCoordinates := 0.0
	maxFloat := float64(^uint(0) >> 1) // 最大的 float64 值
	minFloat := -maxFloat              // 最小的 float64 值

	left, top, right, bottom := maxFloat, maxFloat, minFloat, minFloat

	for _, coord := range features {
		UpdateBounds(coord.X, coord.Y, &left, &top, &right, &bottom)
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

// UpdateBounds 更新边界值
// @param x 当前坐标的X值
// @param y 当前坐标的Y值
// @param left 左边界的指针
// @param top 上边界的指针
// @param right 右边界的指针
// @param bottom 下边界的指针
func UpdateBounds(x, y float64, left, top, right, bottom *float64) {
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

// CalculateControlPoint 计算控制点以实现弯曲效果
// @param start 起始点
// @param end 结束点
// @param offset 偏移量
// @return *gg.Point 返回计算得到的控制点
func CalculateControlPoint(start, end *gg.Point, offset float64) *gg.Point {
	midX := (start.X + end.X) / 2
	midY := (start.Y + end.Y) / 2
	return &gg.Point{X: midX, Y: midY - offset}
}

// CalculateFractionPointMode 类型定义，用于表示计算模式
type CalculateFractionPointMode int

const (
	Add CalculateFractionPointMode = iota
	Subtract
	Multiply
	Divide
)

// CalculateFractionPoint 计算任意分数的坐标
// @param startPoint 起始点
// @param endPoint 结束点
// @param fraction 计算的分数
// @param mode 计算模式（加、减、乘、除）
// @return *gg.Point 返回计算得到的坐标
func CalculateFractionPoint(startPoint, endPoint *gg.Point, fraction float64, mode CalculateFractionPointMode) *gg.Point {
	var x, y float64

	// 计算差值
	deltaX := endPoint.X - startPoint.X
	deltaY := endPoint.Y - startPoint.Y

	switch mode {
	case Add:
		x = deltaX + fraction
		y = deltaY + fraction
	case Subtract:
		x = deltaX - fraction
		y = deltaY - fraction
	case Multiply:
		x = deltaX * fraction
		y = deltaY * fraction
	case Divide:
		if fraction != 0 {
			x = deltaX / fraction
			y = deltaY / fraction
		} else {
			x = deltaX // 避免除以零
			y = deltaY
		}
	default:
		x = startPoint.X
		y = startPoint.Y
	}

	return &gg.Point{
		X: x,
		Y: y,
	}
}

// 定义轴的枚举类型
type AxisPointMode int

const (
	AxisX AxisPointMode = iota // X轴
	AxisY                      // Y轴
)

// 定义最大值和最小值的模式
type PointMaxMin int

const (
	PointMax PointMaxMin = iota // 最大值模式
	PointMin                    // 最小值模式
)

// CalculatePoint 返回两个点中在指定维度上较大的或较小的点
func CalculatePoint(a, b *gg.Point, mode PointMaxMin, axis AxisPointMode) *gg.Point {
	switch axis {
	case AxisY: // 如果比较Y轴
		if (mode == PointMax && a.Y > b.Y) || (mode == PointMin && a.Y < b.Y) {
			return a // 返回a点
		}
		return b // 返回b点
	case AxisX: // 默认比较X轴
		fallthrough // 允许使用fallthrough以简化代码
	default:
		if (mode == PointMax && a.X > b.X) || (mode == PointMin && a.X < b.X) {
			return a // 返回a点
		}
		return b // 返回b点
	}
}

// CalculateMultiplePoints 计算多个点中的最大或最小点
func CalculateMultiplePoints(points []*gg.Point, mode PointMaxMin, axis AxisPointMode) *gg.Point {
	if len(points) == 0 {
		return nil // 如果没有点，返回nil
	}

	// 初始化结果为第一个点
	result := points[0]
	for _, point := range points[1:] {
		result = CalculatePoint(result, point, mode, axis) // 使用CalculatePoint逐个比较
	}
	return result // 返回最终结果
}

// 判断三个点是否能构成三角形
func CanFormTriangle(points []*gg.Point) (float64, bool, error) {
	maxPointLen := 3
	if len(points) < maxPointLen || len(points) > maxPointLen {
		return 0.0, false, fmt.Errorf("points length number must be %d", maxPointLen)
	}

	// 计算三角形的面积
	// 公式的推导：
	// 给定三角形的三个顶点坐标 A(x1, y1)、B(x2, y2)、C(x3, y3)
	// 三角形的面积可以通过以下公式计算：
	// Area = 0.5 * | x1(y2 - y3) + x2(y3 - y1) + x3(y1 - y2) |
	// 其中，points[0] = A, points[1] = B, points[2] = C
	area := 0.5 * (points[0].X*(points[1].Y-points[2].Y) +
		points[1].X*(points[2].Y-points[0].Y) +
		points[2].X*(points[0].Y-points[1].Y))

	return area, area != 0, nil // 如果面积不为零，则可以构成三角形
}

// ResizeX 只缩放 x 坐标
func ResizeX(point *gg.Point, scaleFactor float64) *gg.Point {
	return &gg.Point{
		X: point.X * scaleFactor,
		Y: point.Y, // y 坐标保持不变
	}
}

// ResizeY 只缩放 y 坐标
func ResizeY(point *gg.Point, scaleFactor float64) *gg.Point {
	return &gg.Point{
		Y: point.Y * scaleFactor,
		X: point.X, // x 坐标保持不变
	}
}

// ResizePoint 单个坐标轴的等比例缩放
func ResizePoint(point *gg.Point, resizeX, resizeY float64) *gg.Point {
	return &gg.Point{
		X: ResizeX(point, resizeX).X,
		Y: ResizeY(point, resizeY).Y,
	}
}

// ResizeUpTriangle 根据固定点放大三角形的两个顶点
// @params vertexA: 固定点的坐标
// @params vertexB: 需要放大的第一个点的坐标
// @params vertexC: 需要放大的第二个点的坐标
// @params resizeXy: 放大因子，值大于1表示放大
// @return: 新的 vertexB 和 vertexC 的坐标
func ResizeUpTriangle(vertexA, vertexB, vertexC *gg.Point, resizeXy float64) (*gg.Point, *gg.Point) {
	// 计算 vertexA 到 vertexB 和 vertexA 到 vertexC 的向量
	vectorAB := &gg.Point{X: vertexB.X - vertexA.X, Y: vertexB.Y - vertexA.Y}
	vectorAC := &gg.Point{X: vertexC.X - vertexA.X, Y: vertexC.Y - vertexA.Y}

	// 缩放向量
	resizedVectorAB := ResizePoint(vectorAB, resizeXy, resizeXy)
	resizedVectorAC := ResizePoint(vectorAC, resizeXy, resizeXy)

	// 计算放大后的 vertexB 和 vertexC 的新坐标
	newVertexB := &gg.Point{X: vertexA.X + resizedVectorAB.X, Y: vertexA.Y + resizedVectorAB.Y}
	newVertexC := &gg.Point{X: vertexA.X + resizedVectorAC.X, Y: vertexA.Y + resizedVectorAC.Y}

	return newVertexB, newVertexC
}

// ResizeDownTriangle 根据固定点缩小三角形的两个顶点
// @params vertexA: 固定点的坐标
// @params vertexB: 需要缩小的第一个点的坐标
// @params vertexC: 需要缩小的第二个点的坐标
// @params resizeXy: 缩小因子，值小于1表示缩小
// @return: 新的 vertexB 和 vertexC 的坐标
func ResizeDownTriangle(vertexA, vertexB, vertexC *gg.Point, resizeXy float64) (*gg.Point, *gg.Point) {
	// 确保缩小因子在0到1之间
	if resizeXy >= 1 || resizeXy <= 0 {
		return vertexB, vertexC // 如果因子不在范围内，返回原始坐标
	}

	// 计算 vertexA 到 vertexB 和 vertexA 到 vertexC 的向量
	vectorAB := &gg.Point{X: vertexB.X - vertexA.X, Y: vertexB.Y - vertexA.Y}
	vectorAC := &gg.Point{X: vertexC.X - vertexA.X, Y: vertexC.Y - vertexA.Y}

	// 缩放向量
	resizedVectorAB := ResizePoint(vectorAB, resizeXy, resizeXy)
	resizedVectorAC := ResizePoint(vectorAC, resizeXy, resizeXy)

	// 计算缩小后的 vertexB 和 vertexC 的新坐标
	newVertexB := &gg.Point{X: vertexA.X + resizedVectorAB.X, Y: vertexA.Y + resizedVectorAB.Y}
	newVertexC := &gg.Point{X: vertexA.X + resizedVectorAC.X, Y: vertexA.Y + resizedVectorAC.Y}

	return newVertexB, newVertexC
}

// ResizePoints 批量缩放多个点，支持不同的 x 和 y 缩放因子
// @param points: 需要缩放的点的切片
// @param scaleX: x 轴的缩放因子
// @param scaleY: y 轴的缩放因子
// @return []*gg.Point: 返回缩放后的点的切片
func ResizePoints(points []*gg.Point, scaleX, scaleY float64) []*gg.Point {
	resizedPoints := make([]*gg.Point, len(points)) // 创建一个新的切片用于存储缩放后的点
	for i, point := range points {
		resizedPoints[i] = ResizePoint(point, scaleX, scaleY) // 使用缩放函数逐个缩放点
	}
	return resizedPoints // 返回缩放后的点的切片
}
