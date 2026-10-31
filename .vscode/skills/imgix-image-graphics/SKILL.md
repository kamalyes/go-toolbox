---
name: imgix-image-graphics
description: 图像处理与二维几何绘图，提供图像加载/裁剪/缩放/叠色/Base64 编解码与多格式写出，以及点运算、分数点、延长线、三角形判定、矩形与虚线渲染器；当需要处理图片或在画布上做坐标计算与绘制时使用
---

# imgix - 图像处理与二维几何

两块能力：**图像 I/O 与变换**，以及**基于 `gg.Point` 的二维几何运算与绘制**

## 先看依赖代价

**这是 go-toolbox 里最重的包。** 引入它会带进来：

- `github.com/fogleman/gg`（画布与绘图）
- `github.com/disintegration/imaging`（缩放重采样）
- `github.com/golang/freetype`（字体栅格化）
- 以及 `golang.org/x/image`

只需要「判断图片格式」或「算两个点的中点」这种小事，不值得引这一串
反过来，真要画图或处理图片时这些都躲不开

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/imgix"
```

图像加载、裁剪、缩放、写出：
```go
img, err := imgix.LoadImageFromFile("in.png")

img = imgix.CropImage(img, &imgix.CropImgOptions{
    MinWidth: 100, MinHeight: 100, MaxWidth: 800, MaxHeight: 800,
})
img = imgix.ResizeImage(img, &imgix.ResizeImgOptions{
    Width: 320, Height: 240, Filter: imaging.Lanczos,
})

err = imgix.WriterImage(img, 85, imgix.JPEG, w)        // 写到 io.Writer
b64, err := imgix.EncodeImageToBase64(img, imgix.PNG, 100)
```

点运算与分数点：
```go
mid  := imgix.CalculateFractionPoint(a, b, 0.5, imgix.Add)   // 线段上的分数点
peak := imgix.CalculateMultiplePoints(pts, imgix.CalculateMax, imgix.AxisY)
ext  := imgix.ExtendLine(p1, p2, 20)                         // 沿 p1→p2 延长 20
```

虚线渲染：
```go
r := imgix.NewGraphicsRenderer(ctx, imgix.NewDashOptions(6, 3, 1.5))
```

## 完整API索引

### 图像加载与输出

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `LoadImageFromFile` | `func(filePath string) (image.Image, error)` | 从文件加载 |
| `LoadImageFromURL` | `func(url string) (image.Image, error)` | **从网络加载**，会发起 HTTP 请求 |
| `LoadImageFromFileAndEncodeToBase64` | `func(filePath string, format ImageFormat, quality int) (string, error)` | 加载并直接转 Base64 |
| `EncodeImageToBase64` | `func(img image.Image, format ImageFormat, quality int) (string, error)` | 编码为 Base64 |
| `WriterImage` | `func(img image.Image, quality int, format ImageFormat, output io.Writer) error` | 按格式写出 |
| `WriterImageFromBytes` | `func(data []byte, quality int, format ImageFormat, output io.Writer) error` | 从字节写出 |
| `SaveBufToImageFile` | `func(buf *bytes.Buffer, filePath string, format ImageFormat) error` | 缓冲写入文件 |

### 图像变换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `CropImage` | `func(img image.Image, opts *CropImgOptions) image.Image` | 按最小/最大宽高裁剪 |
| `ResizeImage` | `func(img image.Image, opts *ResizeImgOptions) image.Image` | 缩放，重采样算法由 `Filter` 指定 |
| `AddOverlay` | `func(originalImg image.Image, overlayColor color.RGBA) image.Image` | 叠加一层纯色 |

`CropImgOptions`：`MinWidth`、`MinHeight`、`MaxWidth`、`MaxHeight`
`ResizeImgOptions`：`Width`、`Height`、`Filter`（`imaging.ResampleFilter`）

### 点与线运算（入参多为 `*gg.Point`）

| 导出名称 | 说明 |
|---|---|
| `CalculatePoint` | 两点按 `CalculateMode` + `AxisPointMode` 运算 |
| `CalculateMultiplePoints` | 多点按模式归约（最大/最小/平均） |
| `CalculateFractionPoint` | 线段上的分数点 |
| `CalculateControlPoint` | 带偏移的贝塞尔控制点 |
| `ExtendLine` | 沿方向延长指定长度 |
| `ModifyPoint` / `OffsetPointX` / `OffsetPointY` | 按值修改 / 偏移单轴 |
| `ResizePoint` / `ResizePoints` / `ResizeX` / `ResizeY` | 缩放点或点集 |
| `ResizePointOneselfX` / `ResizePointOneselfY` | 按自身比例缩放单轴 |
| `ResizeUpTriangle` / `ResizeDownTriangle` | 三角形顶点缩放 |
| `CanFormTriangle` | 能否构成三角形，返回面积与判定 |
| `AdjustValues` | 调整起止值使其包含目标 |
| `UpdateBounds` | 用点更新 LTRB 边界（指针原地改） |
| `GetPointOperation` | 取变参里的运算模式，缺省有默认值 |

### 枚举

| 类型 | 取值 |
|---|---|
| `ImageFormat` | `JPEG`、`PNG`、`GIF`（有 `String()`） |
| `CalculateMode` | `CalculateMax`、`CalculateMin`、`CalculateAvg` |
| `AxisPointMode` | `AxisX`、`AxisY` |
| `CalculateFractionPointMode` | `Add`、`Subtract`、`Multiply`、`Divide` |

### 几何类型与绘制

| 导出名称 | 说明 |
|---|---|
| `Rectangle` / `NewRectangle` | 矩形，可带 `RectangleExtensions` |
| `RectangleExtensions` / `NewRectangleExtensions` | 矩形扩展行为 |
| `ImageLTRB` / `GetLTRB` | 左上右下边界；`GetLTRB` 从特征点 map 求包围盒 |
| `Coordinates` / `CleanCoordinates` | 坐标集合；`CleanCoordinates` 清洗特征点 map |
| `HorizontalLine` / `VerticalLine` / `LineSegment` | 水平线 / 垂直线 / 线段 |
| `HorizontalEdge` / `VerticalEdge` | 水平边 / 垂直边 |
| `DelayLine` | 延迟线 |
| `DashStyle` / `DashOptions` / `NewDashOptions` | 虚线样式（线长、间隔、线宽） |
| `GraphicsRenderer` / `NewGraphicsRenderer` | 绘图渲染器，包装 `*gg.Context` |
| `MatchesAnyPattern` | `func(input string, patterns []string) bool`，多模式匹配 |

### 点运算策略

`PointOperationStrategy` 接口 + `PointAddStrategy`、`PointSubtractStrategy`、
`PointMultiplyStrategy`、`PointDivideStrategy` 四个空结构体实现

## 注意事项

- **`LoadImageFromURL` 会发起网络请求**，且不接受 `context.Context`，
  因此无法随调用方取消；URL 来自用户输入时要自己做 SSRF 防护与大小限制
- **`quality` 只对 JPEG 有意义**，PNG/GIF 是无损的，传什么都不改变输出体积
- `ImageFormat` 只有 `JPEG`/`PNG`/`GIF` 三种，**没有 WebP / BMP / TIFF**，
  传了超出范围的值要看具体函数如何兜底，不要假定会报错
- `CropImage` 与 `ResizeImage` 的入参是**指针**（`*CropImgOptions`），传 nil 要小心
- `UpdateBounds` 通过 `*float64` 指针原地修改边界，是少见的输出参数写法，
  调用前要先把四个变量初始化成合理的初值（否则第一次比较就错）
- `gg.Point` 的坐标是 `float64`，几何函数不做取整；
  最终落到像素时的取整策略由调用方决定，不同取整会让相邻图元出现一像素缝
- 这些几何函数**绝大多数返回新的 `*gg.Point`**，不修改入参；
  例外是上面点明的 `UpdateBounds`
