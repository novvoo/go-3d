package go3d

import "math"

// Orbit 轨道
type Orbit struct {
	Radius    float64
	Color     [3]float64
	Thickness float64
	Segments  int
}

// NewOrbit 创建轨道
func NewOrbit(radius float64, color [3]float64) *Orbit {
	return &Orbit{
		Radius:    radius,
		Color:     color,
		Thickness: 0.01,
		Segments:  64,
	}
}

// Render 渲染轨道
func (o *Orbit) Render(renderer *Renderer, t float64) {
	orbit := CreateTorus(o.Radius, o.Thickness, o.Segments, 4)
	transform := Identity()
	// 不需要旋转，轨道默认就在XY平面上
	transformedOrbit := orbit.Transform(transform)
	renderer.DrawMesh(transformedOrbit, o.Color)
}

// Star 星星
type Star struct {
	Position   Vector3
	Radius     float64
	Color      [3]float64
	Brightness float64
	Twinkle    bool
	Phase      float64 // 闪烁相位
}

// NewStar 创建星星
func NewStar(position Vector3, radius float64, color [3]float64) *Star {
	return &Star{
		Position:   position,
		Radius:     radius,
		Color:      color,
		Brightness: 1.0,
		Twinkle:    false,
		Phase:      0,
	}
}

// SetTwinkle 设置闪烁
func (s *Star) SetTwinkle(phase float64) *Star {
	s.Twinkle = true
	s.Phase = phase
	return s
}

// Render 渲染星星
func (s *Star) Render(renderer *Renderer, t float64) {
	brightness := s.Brightness
	if s.Twinkle {
		brightness = s.Brightness * (0.7 + 0.3*math.Sin(t*4*math.Pi+s.Phase))
	}

	color := [3]float64{
		s.Color[0] * brightness,
		s.Color[1] * brightness,
		s.Color[2] * brightness,
	}

	star := CreateSphere(s.Radius, 6, 6)
	transform := Identity()
	transform = transform.Multiply(Translation(s.Position.X, s.Position.Y, s.Position.Z))
	transformedStar := star.Transform(transform)
	renderer.DrawMesh(transformedStar, color)
}

// StarField 星空场
type StarField struct {
	Stars []Star
}

// NewStarField 创建星空场
func NewStarField(numStars int, distance float64) *StarField {
	sf := &StarField{
		Stars: make([]Star, numStars),
	}

	for i := range numStars {
		// 伪随机位置
		angle1 := float64(i) * 2.4
		angle2 := float64(i) * 1.7
		dist := distance + float64(i%10)*2.0

		x := dist * math.Cos(angle1) * math.Sin(angle2)
		y := dist * math.Sin(angle1) * 0.5
		z := dist * math.Cos(angle1) * math.Cos(angle2)

		// 不同颜色的星星
		var starColor [3]float64
		colorType := i % 3
		switch colorType {
		case 0:
			starColor = [3]float64{1.0, 1.0, 1.0} // 白色
		case 1:
			starColor = [3]float64{0.8, 0.9, 1.0} // 淡蓝色
		case 2:
			starColor = [3]float64{1.0, 0.95, 0.7} // 淡黄色
		}

		sf.Stars[i] = *NewStar(NewVector3(x, y, z), 0.05, starColor).SetTwinkle(float64(i))
	}

	return sf
}

// Render 渲染星空场
func (sf *StarField) Render(renderer *Renderer, t float64) {
	for i := range sf.Stars {
		sf.Stars[i].Render(renderer, t)
	}
}
