package go3d

import (
	"math"
	"sort"

	"github.com/novvoo/go-cairo/pkg/cairo"
)

// Camera 表示相机
type Camera struct {
	Position Vector3
	Target   Vector3
	Up       Vector3
	FOV      float64
	Near     float64
	Far      float64
}

// NewCamera 创建新相机
func NewCamera() *Camera {
	return &Camera{
		Position: NewVector3(0, 0, -5),
		Target:   NewVector3(0, 0, 0),
		Up:       NewVector3(0, 1, 0),
		FOV:      1.0,
		Near:     0.1,
		Far:      100.0,
	}
}

// Light 表示光源
type Light struct {
	Position  Vector3
	Color     [3]float64
	Intensity float64
}

// NewLight 创建新光源
func NewLight(pos Vector3, color [3]float64, intensity float64) *Light {
	return &Light{
		Position:  pos,
		Color:     color,
		Intensity: intensity,
	}
}

// RenderMode 渲染模式
type RenderMode int

const (
	RenderWireframe RenderMode = iota // 线框模式
	RenderFlat                        // 平面着色
	RenderShaded                      // 光照着色
)

// Renderer 3D渲染器
type Renderer struct {
	Surface    cairo.ImageSurface
	Context    cairo.Context
	Width      int
	Height     int
	Camera     *Camera
	Lights     []*Light
	RenderMode RenderMode
	Antialias  bool
}

// NewRenderer 创建新渲染器
func NewRenderer(width, height int) *Renderer {
	surface := cairo.NewImageSurface(cairo.FormatARGB32, width, height)
	context := cairo.NewContext(surface)

	renderer := &Renderer{
		Surface:    surface.(cairo.ImageSurface),
		Context:    context,
		Width:      width,
		Height:     height,
		Camera:     NewCamera(),
		Lights:     make([]*Light, 0),
		RenderMode: RenderWireframe,
		Antialias:  true,
	}

	// 设置合成模式为 SOURCE，确保完全覆盖
	renderer.Context.SetOperator(cairo.OperatorSource)

	// 初始化时清除画布为完全透明的黑色
	renderer.Context.SetSourceRGBA(0, 0, 0, 1.0)
	renderer.Context.Paint()

	// 恢复为正常的 OVER 模式用于后续绘制
	renderer.Context.SetOperator(cairo.OperatorOver)

	return renderer
}

// AddLight 添加光源
func (r *Renderer) AddLight(light *Light) {
	r.Lights = append(r.Lights, light)
}

// SetRenderMode 设置渲染模式
func (r *Renderer) SetRenderMode(mode RenderMode) {
	r.RenderMode = mode
}

// SetAntialias 设置抗锯齿
func (r *Renderer) SetAntialias(enabled bool) {
	r.Antialias = enabled
	if enabled {
		r.Context.SetAntialias(cairo.AntialiasGood)
	} else {
		r.Context.SetAntialias(cairo.AntialiasNone)
	}
}

// Clear 清空画布
func (r *Renderer) Clear(red, green, blue float64) {
	r.Context.SetSourceRGB(red, green, blue)
	r.Context.Paint()
}

// ProjectToScreen 将3D坐标投影到屏幕坐标
func (r *Renderer) ProjectToScreen(v Vector3) (float64, float64, float64) {
	aspect := float64(r.Width) / float64(r.Height)

	// 创建视图矩阵和投影矩阵
	view := LookAt(r.Camera.Position, r.Camera.Target, r.Camera.Up)
	projection := Perspective(r.Camera.FOV, aspect, r.Camera.Near, r.Camera.Far)

	// 先应用视图变换，再应用投影变换
	viewSpace := view.TransformVector(v)
	projected := projection.TransformVector(viewSpace)

	// 转换到屏幕坐标
	x := (projected.X + 1.0) * float64(r.Width) / 2.0
	y := (1.0 - projected.Y) * float64(r.Height) / 2.0

	return x, y, projected.Z
}

// CalculateLighting 计算光照
func (r *Renderer) CalculateLighting(position, normal Vector3, baseColor [3]float64) [3]float64 {
	if len(r.Lights) == 0 {
		return baseColor
	}

	ambient := [3]float64{0.2, 0.2, 0.2}
	diffuse := [3]float64{0, 0, 0}

	for _, light := range r.Lights {
		lightDir := light.Position.Sub(position).Normalize()
		intensity := math.Max(0, normal.Dot(lightDir)) * light.Intensity

		diffuse[0] += light.Color[0] * intensity
		diffuse[1] += light.Color[1] * intensity
		diffuse[2] += light.Color[2] * intensity
	}

	return [3]float64{
		math.Min(1.0, (ambient[0]+diffuse[0])*baseColor[0]),
		math.Min(1.0, (ambient[1]+diffuse[1])*baseColor[1]),
		math.Min(1.0, (ambient[2]+diffuse[2])*baseColor[2]),
	}
}

// triangleWithDepth 带深度信息的三角形
type triangleWithDepth struct {
	tri   Triangle
	depth float64
	color [3]float64
}

// DrawMesh 绘制网格
func (r *Renderer) DrawMesh(mesh *Mesh, color [3]float64) {
	switch r.RenderMode {
	case RenderWireframe:
		r.drawWireframe(mesh, color)
	case RenderFlat:
		r.drawFlat(mesh, color)
	case RenderShaded:
		r.drawShaded(mesh, color)
	}
}

// drawWireframe 绘制线框
func (r *Renderer) drawWireframe(mesh *Mesh, color [3]float64) {
	if len(mesh.Triangles) == 0 {
		return
	}

	r.Context.Save()
	defer r.Context.Restore()

	r.Context.SetSourceRGB(color[0], color[1], color[2])
	r.Context.SetLineWidth(1.5)
	r.Context.SetLineJoin(cairo.LineJoinRound)

	for _, tri := range mesh.Triangles {
		x0, y0, z0 := r.ProjectToScreen(tri.V0)
		x1, y1, z1 := r.ProjectToScreen(tri.V1)
		x2, y2, z2 := r.ProjectToScreen(tri.V2)

		// 简单的视锥剔除
		if z0 < -1 || z0 > 1 || z1 < -1 || z1 > 1 || z2 < -1 || z2 > 1 {
			continue
		}

		r.Context.MoveTo(x0, y0)
		r.Context.LineTo(x1, y1)
		r.Context.LineTo(x2, y2)
		r.Context.ClosePath()
		r.Context.Stroke()
	}
}

// drawFlat 绘制平面着色
func (r *Renderer) drawFlat(mesh *Mesh, color [3]float64) {
	if len(mesh.Triangles) == 0 {
		return
	}

	r.Context.Save()
	defer r.Context.Restore()

	// 预分配切片容量
	triangles := make([]triangleWithDepth, 0, len(mesh.Triangles))

	for _, tri := range mesh.Triangles {
		_, _, z0 := r.ProjectToScreen(tri.V0)
		_, _, z1 := r.ProjectToScreen(tri.V1)
		_, _, z2 := r.ProjectToScreen(tri.V2)

		// 视锥剔除
		if z0 < -1 || z1 < -1 || z2 < -1 {
			continue
		}

		avgDepth := (z0 + z1 + z2) / 3.0

		triangles = append(triangles, triangleWithDepth{
			tri:   tri,
			depth: avgDepth,
			color: color,
		})
	}

	// 从远到近排序
	sort.Slice(triangles, func(i, j int) bool {
		return triangles[i].depth > triangles[j].depth
	})

	// 绘制三角形
	r.Context.SetSourceRGB(color[0], color[1], color[2])
	for _, td := range triangles {
		x0, y0, _ := r.ProjectToScreen(td.tri.V0)
		x1, y1, _ := r.ProjectToScreen(td.tri.V1)
		x2, y2, _ := r.ProjectToScreen(td.tri.V2)

		r.Context.MoveTo(x0, y0)
		r.Context.LineTo(x1, y1)
		r.Context.LineTo(x2, y2)
		r.Context.ClosePath()
		r.Context.Fill()
	}
}

// drawShaded 绘制光照着色
func (r *Renderer) drawShaded(mesh *Mesh, color [3]float64) {
	if len(mesh.Triangles) == 0 {
		return
	}

	r.Context.Save()
	defer r.Context.Restore()

	// 预分配切片容量
	triangles := make([]triangleWithDepth, 0, len(mesh.Triangles))

	for _, tri := range mesh.Triangles {
		_, _, z0 := r.ProjectToScreen(tri.V0)
		_, _, z1 := r.ProjectToScreen(tri.V1)
		_, _, z2 := r.ProjectToScreen(tri.V2)

		// 视锥剔除
		if z0 < -1 || z1 < -1 || z2 < -1 {
			continue
		}

		avgDepth := (z0 + z1 + z2) / 3.0

		// 计算法线
		normal := tri.Normal()

		// 背面剔除
		viewDir := r.Camera.Position.Sub(tri.Center()).Normalize()
		if normal.Dot(viewDir) < 0 {
			continue
		}

		// 计算三角形中心
		center := tri.Center()

		// 计算光照颜色
		litColor := r.CalculateLighting(center, normal, color)

		triangles = append(triangles, triangleWithDepth{
			tri:   tri,
			depth: avgDepth,
			color: litColor,
		})
	}

	// 从远到近排序
	sort.Slice(triangles, func(i, j int) bool {
		return triangles[i].depth > triangles[j].depth
	})

	// 绘制三角形
	for _, td := range triangles {
		x0, y0, _ := r.ProjectToScreen(td.tri.V0)
		x1, y1, _ := r.ProjectToScreen(td.tri.V1)
		x2, y2, _ := r.ProjectToScreen(td.tri.V2)

		r.Context.MoveTo(x0, y0)
		r.Context.LineTo(x1, y1)
		r.Context.LineTo(x2, y2)
		r.Context.ClosePath()

		r.Context.SetSourceRGB(td.color[0], td.color[1], td.color[2])
		r.Context.Fill()
	}
}

// DrawMeshWithGradient 使用渐变绘制网格
func (r *Renderer) DrawMeshWithGradient(mesh *Mesh, color1, color2 [3]float64) {
	if len(mesh.Triangles) == 0 {
		return
	}

	r.Context.Save()
	defer r.Context.Restore()

	// 预分配切片容量
	triangles := make([]triangleWithDepth, 0, len(mesh.Triangles))

	for _, tri := range mesh.Triangles {
		_, _, z0 := r.ProjectToScreen(tri.V0)
		_, _, z1 := r.ProjectToScreen(tri.V1)
		_, _, z2 := r.ProjectToScreen(tri.V2)

		// 视锥剔除
		if z0 < -1 || z1 < -1 || z2 < -1 {
			continue
		}

		avgDepth := (z0 + z1 + z2) / 3.0

		// 根据深度计算渐变颜色
		t := (avgDepth + 1.0) / 2.0 // 归一化到 0-1
		color := [3]float64{
			color1[0]*(1-t) + color2[0]*t,
			color1[1]*(1-t) + color2[1]*t,
			color1[2]*(1-t) + color2[2]*t,
		}

		triangles = append(triangles, triangleWithDepth{
			tri:   tri,
			depth: avgDepth,
			color: color,
		})
	}

	// 从远到近排序
	sort.Slice(triangles, func(i, j int) bool {
		return triangles[i].depth > triangles[j].depth
	})

	// 绘制三角形
	for _, td := range triangles {
		x0, y0, _ := r.ProjectToScreen(td.tri.V0)
		x1, y1, _ := r.ProjectToScreen(td.tri.V1)
		x2, y2, _ := r.ProjectToScreen(td.tri.V2)

		r.Context.MoveTo(x0, y0)
		r.Context.LineTo(x1, y1)
		r.Context.LineTo(x2, y2)
		r.Context.ClosePath()

		r.Context.SetSourceRGB(td.color[0], td.color[1], td.color[2])
		r.Context.Fill()
	}
}

// SaveToPNG 保存为PNG文件
func (r *Renderer) SaveToPNG(filename string) error {
	r.Surface.WriteToPNG(filename)
	return nil
}

// Destroy 释放资源
func (r *Renderer) Destroy() {
	r.Context.Destroy()
	r.Surface.Destroy()
}
