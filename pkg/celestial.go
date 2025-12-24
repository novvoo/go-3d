package go3d

import "math"

// Planet 行星
type Planet struct {
	Name          string
	NameCN        string
	Radius        float64
	OrbitRadius   float64
	OrbitSpeed    float64
	RotationSpeed float64
	Color         [3]float64
	UseGradient   bool
	GradientColor [3]float64
	HasMoon       bool
	HasRings      bool
	RingColors    [][3]float64
}

// NewPlanet 创建行星
func NewPlanet(name, nameCN string, radius, orbitRadius, orbitSpeed, rotationSpeed float64, color [3]float64) *Planet {
	return &Planet{
		Name:          name,
		NameCN:        nameCN,
		Radius:        radius,
		OrbitRadius:   orbitRadius,
		OrbitSpeed:    orbitSpeed,
		RotationSpeed: rotationSpeed,
		Color:         color,
		UseGradient:   false,
	}
}

// SetGradient 设置渐变色
func (p *Planet) SetGradient(color1, color2 [3]float64) *Planet {
	p.UseGradient = true
	p.Color = color1
	p.GradientColor = color2
	return p
}

// AddMoon 添加月球
func (p *Planet) AddMoon() *Planet {
	p.HasMoon = true
	return p
}

// AddRings 添加光环
func (p *Planet) AddRings(colors [][3]float64) *Planet {
	p.HasRings = true
	p.RingColors = colors
	return p
}

// GetPosition 获取行星在指定时间的位置
func (p *Planet) GetPosition(t float64) Vector3 {
	angle := t * p.OrbitSpeed * math.Pi
	x := p.OrbitRadius * math.Cos(angle)
	z := p.OrbitRadius * math.Sin(angle)
	y := 0.0 // 行星在XZ平面上运动，与轨道圈一致
	return NewVector3(x, y, z)
}

// Render 渲染行星
func (p *Planet) Render(renderer *Renderer, t float64) {
	pos := p.GetPosition(t)

	// 创建行星球体
	planetMesh := CreateSphere(p.Radius, 16, 16)

	// 应用变换
	transform := Identity()
	transform = transform.Multiply(Translation(pos.X, pos.Y, pos.Z))
	transform = transform.Multiply(RotationY(t * p.RotationSpeed * math.Pi))

	transformedPlanet := planetMesh.Transform(transform)

	// 渲染行星
	if p.UseGradient {
		renderer.DrawMeshWithGradient(transformedPlanet, p.Color, p.GradientColor)
	} else {
		renderer.DrawMesh(transformedPlanet, p.Color)
	}

	// 渲染标签
	labelPos := NewVector3(pos.X, pos.Y+p.Radius+0.3, pos.Z)
	label := NewLabel3D(labelPos, p.NameCN, [3]float64{1, 1, 1})
	label.Render(renderer, t)

	// 渲染月球
	if p.HasMoon {
		p.renderMoon(renderer, pos, t)
	}

	// 渲染光环
	if p.HasRings {
		p.renderRings(renderer, pos, t)
	}
}

// renderMoon 渲染月球
func (p *Planet) renderMoon(renderer *Renderer, planetPos Vector3, t float64) {
	moonOrbitRadius := p.Radius * 2
	moonAngle := t * 8.0 * math.Pi

	moonX := planetPos.X + moonOrbitRadius*math.Cos(moonAngle)
	moonZ := planetPos.Z + moonOrbitRadius*math.Sin(moonAngle)
	moonY := planetPos.Y + math.Sin(moonAngle)*0.1

	moon := CreateSphere(p.Radius*0.3, 10, 10)
	transform := Identity()
	transform = transform.Multiply(Translation(moonX, moonY, moonZ))
	transformedMoon := moon.Transform(transform)
	renderer.DrawMesh(transformedMoon, [3]float64{0.95, 0.95, 0.95})
}

// renderRings 渲染光环
func (p *Planet) renderRings(renderer *Renderer, planetPos Vector3, t float64) {
	if len(p.RingColors) == 0 {
		return
	}

	baseRadius := p.Radius * 1.4
	for i, color := range p.RingColors {
		radius := baseRadius + float64(i)*p.Radius*0.3
		ring := CreateTorus(radius, 0.02, 48, 6)

		transform := Identity()
		transform = transform.Multiply(Translation(planetPos.X, planetPos.Y, planetPos.Z))
		transform = transform.Multiply(RotationX(math.Pi/2 + 0.3))
		transform = transform.Multiply(RotationY(t * math.Pi))

		transformedRing := ring.Transform(transform)
		renderer.DrawMesh(transformedRing, color)
	}
}
