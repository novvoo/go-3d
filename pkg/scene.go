package go3d

import (
	"math"

	"github.com/novvoo/go-cairo/pkg/cairo"
)

// SceneObject 场景对象接口
type SceneObject interface {
	Render(renderer *Renderer, t float64)
}

// Scene 场景管理器
type Scene struct {
	Objects    []SceneObject
	Lights     []*Light
	Background BackgroundRenderer
}

// NewScene 创建场景
func NewScene() *Scene {
	return &Scene{
		Objects: make([]SceneObject, 0),
		Lights:  make([]*Light, 0),
	}
}

// AddObject 添加场景对象
func (s *Scene) AddObject(obj SceneObject) {
	s.Objects = append(s.Objects, obj)
}

// AddLight 添加光源
func (s *Scene) AddLight(light *Light) {
	s.Lights = append(s.Lights, light)
}

// SetBackground 设置背景渲染器
func (s *Scene) SetBackground(bg BackgroundRenderer) {
	s.Background = bg
}

// Render 渲染整个场景
func (s *Scene) Render(renderer *Renderer, t float64) {
	// 设置光源
	renderer.Lights = s.Lights

	// 渲染背景
	if s.Background != nil {
		s.Background.Render(renderer, t)
	}

	// 渲染所有对象
	for _, obj := range s.Objects {
		obj.Render(renderer, t)
	}
}

// BackgroundRenderer 背景渲染器接口
type BackgroundRenderer interface {
	Render(renderer *Renderer, t float64)
}

// GradientBackground 渐变背景
type GradientBackground struct {
	TopColor    [3]float64
	BottomColor [3]float64
	Steps       int
	Animated    bool // 是否随时间动画
}

// NewGradientBackground 创建渐变背景
func NewGradientBackground(topColor, bottomColor [3]float64) *GradientBackground {
	return &GradientBackground{
		TopColor:    topColor,
		BottomColor: bottomColor,
		Steps:       100,
		Animated:    false,
	}
}

// Render 渲染渐变背景
func (gb *GradientBackground) Render(renderer *Renderer, t float64) {
	renderer.Context.Save()
	defer renderer.Context.Restore()

	topR, topG, topB := gb.TopColor[0], gb.TopColor[1], gb.TopColor[2]
	bottomR, bottomG, bottomB := gb.BottomColor[0], gb.BottomColor[1], gb.BottomColor[2]

	// 如果启用动画，添加色调偏移
	if gb.Animated {
		hueShift := math.Sin(t*math.Pi) * 0.05
		topR += hueShift
		topG += hueShift
		topB += hueShift * 0.5
		bottomR += hueShift * 0.5
		bottomG += hueShift * 0.5
		bottomB += hueShift * 0.3
	}

	// 先用顶部颜色填充整个背景，确保没有空白
	renderer.Context.SetSourceRGB(topR, topG, topB)
	renderer.Context.Rectangle(0, 0, float64(renderer.Width), float64(renderer.Height))
	renderer.Context.Fill()

	// 然后绘制渐变条
	for i := range gb.Steps {
		ratio := float64(i) / float64(gb.Steps)

		r := topR + (bottomR-topR)*ratio
		g := topG + (bottomG-topG)*ratio
		b := topB + (bottomB-topB)*ratio

		y := float64(i) * float64(renderer.Height) / float64(gb.Steps)
		h := float64(renderer.Height) / float64(gb.Steps)

		renderer.Context.SetSourceRGB(r, g, b)
		renderer.Context.Rectangle(0, y, float64(renderer.Width), h)
		renderer.Context.Fill()
	}
}

// SolidBackground 纯色背景
type SolidBackground struct {
	Color [3]float64
}

// NewSolidBackground 创建纯色背景
func NewSolidBackground(color [3]float64) *SolidBackground {
	return &SolidBackground{Color: color}
}

// Render 渲染纯色背景
func (sb *SolidBackground) Render(renderer *Renderer, t float64) {
	renderer.Context.Save()
	renderer.Context.SetSourceRGB(sb.Color[0], sb.Color[1], sb.Color[2])
	renderer.Context.Rectangle(0, 0, float64(renderer.Width), float64(renderer.Height))
	renderer.Context.Fill()
	renderer.Context.Restore()
}

// Label3D 3D 标签
type Label3D struct {
	Position Vector3
	Text     string
	Color    [3]float64
	FontSize float64
	Bold     bool
}

// NewLabel3D 创建 3D 标签
func NewLabel3D(position Vector3, text string, color [3]float64) *Label3D {
	return &Label3D{
		Position: position,
		Text:     text,
		Color:    color,
		FontSize: 20.0,
		Bold:     true,
	}
}

// Render 渲染标签
func (l *Label3D) Render(renderer *Renderer, t float64) {
	x, y, z := renderer.ProjectToScreen(l.Position)

	// 只绘制在视野内的标签
	if z > -1 && z < 1 {
		renderer.Context.Save()
		defer renderer.Context.Restore()

		// 根据深度调整大小，但保持完全不透明
		depth := (z + 1) / 2
		fontSize := l.FontSize * (1.0 - depth*0.3)

		// 创建 Pango 布局用于文字渲染
		layout := renderer.Context.PangoCairoCreateLayout()
		defer func() {
			// 确保布局资源被释放
			if pangoLayout, ok := layout.(*cairo.PangoCairoLayout); ok {
				pangoLayout.Destroy()
			}
		}()

		if pangoLayout, ok := layout.(*cairo.PangoCairoLayout); ok {
			fontDesc := cairo.NewPangoFontDescription()

			fontDesc.SetFamily("sans-serif")
			if l.Bold {
				fontDesc.SetWeight(700)
			}
			fontDesc.SetSize(fontSize)

			pangoLayout.SetFontDescription(fontDesc)
			pangoLayout.SetText(l.Text)

			extents := pangoLayout.GetPixelExtents()
			textWidth := float64(extents.Width)
			textHeight := float64(extents.Height)

			textX := x - textWidth/2
			textY := y - textHeight

			// 使用完全不透明的颜色，alpha = 1.0
			renderer.Context.SetSourceRGBA(l.Color[0], l.Color[1], l.Color[2], 1.0)
			renderer.Context.MoveTo(textX, textY)
			renderer.Context.PangoCairoShowText(layout)
		}
	}
}

// CoordinateSystem 坐标系统
type CoordinateSystem struct {
	Length     float64
	Thickness  float64
	ShowLabels bool
}

// NewCoordinateSystem 创建坐标系统
func NewCoordinateSystem(length float64) *CoordinateSystem {
	return &CoordinateSystem{
		Length:     length,
		Thickness:  0.03,
		ShowLabels: true,
	}
}

// Render 渲染坐标系统
func (cs *CoordinateSystem) Render(renderer *Renderer, t float64) {
	// X轴 - 红色
	cs.drawAxis(renderer, NewVector3(0, 0, 0), NewVector3(cs.Length, 0, 0),
		[3]float64{1.0, 0.3, 0.3}, "X")

	// Y轴 - 绿色
	cs.drawAxis(renderer, NewVector3(0, 0, 0), NewVector3(0, cs.Length, 0),
		[3]float64{0.3, 1.0, 0.3}, "Y")

	// Z轴 - 蓝色
	cs.drawAxis(renderer, NewVector3(0, 0, 0), NewVector3(0, 0, cs.Length),
		[3]float64{0.3, 0.3, 1.0}, "Z")
}

// drawAxis 绘制单个坐标轴
func (cs *CoordinateSystem) drawAxis(renderer *Renderer, start, end Vector3, color [3]float64, label string) {
	direction := end.Sub(start)
	length := direction.Length()

	// 绘制轴线（圆柱体）
	cylinder := CreateCylinder(cs.Thickness, length, 8)

	up := NewVector3(0, 1, 0)
	axis := up.Cross(direction.Normalize())
	angle := math.Acos(up.Dot(direction.Normalize()))

	transform := Identity()
	transform = transform.Multiply(Translation(
		(start.X+end.X)/2,
		(start.Y+end.Y)/2,
		(start.Z+end.Z)/2,
	))

	if axis.Length() > 0.001 {
		transform = transform.Multiply(RotationFromAxisAngle(axis.Normalize(), angle))
	}

	transformedCylinder := cylinder.Transform(transform)
	renderer.DrawMesh(transformedCylinder, color)

	// 绘制箭头
	cone := CreateCone(cs.Thickness*4, cs.Length*0.05, 8)
	coneTransform := Identity()
	coneTransform = coneTransform.Multiply(Translation(end.X, end.Y, end.Z))
	if axis.Length() > 0.001 {
		coneTransform = coneTransform.Multiply(RotationFromAxisAngle(axis.Normalize(), angle))
	}
	transformedCone := cone.Transform(coneTransform)
	renderer.DrawMesh(transformedCone, color)

	// 绘制标签
	if cs.ShowLabels {
		labelPos := end.Scale(1.15)
		labelObj := NewLabel3D(labelPos, label, color)
		labelObj.Render(renderer, 0)
	}
}

// RotationFromAxisAngle 从轴角创建旋转矩阵
func RotationFromAxisAngle(axis Vector3, angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	t := 1 - c

	x, y, z := axis.X, axis.Y, axis.Z

	return Matrix4{
		t*x*x + c, t*x*y - s*z, t*x*z + s*y, 0,
		t*x*y + s*z, t*y*y + c, t*y*z - s*x, 0,
		t*x*z - s*y, t*y*z + s*x, t*z*z + c, 0,
		0, 0, 0, 1,
	}
}
