package go3d

// SolarSystem 太阳系
type SolarSystem struct {
	Sun     *CelestialBody
	Planets []*Planet
	Orbits  []*Orbit
	Stars   *StarField
}

// NewSolarSystem 创建太阳系
func NewSolarSystem() *SolarSystem {
	ss := &SolarSystem{
		Planets: make([]*Planet, 0),
		Orbits:  make([]*Orbit, 0),
	}

	// 创建太阳 - MUI 风格：使用 Yellow/Orange 渐变
	ss.Sun = NewCelestialBody("Sun", "太阳", 0.8, [3]float64{1.0, 0.92, 0.23})
	ss.Sun.SetGradient([3]float64{1.0, 0.92, 0.23}, [3]float64{1.0, 0.6, 0.0}) // Yellow 400 -> Orange 500
	ss.Sun.RotationSpeed = 1.0

	// 创建星空背景
	ss.Stars = NewStarField(50, 20.0)

	return ss
}

// AddPlanet 添加行星
func (ss *SolarSystem) AddPlanet(planet *Planet) {
	ss.Planets = append(ss.Planets, planet)
}

// AddOrbit 添加轨道
func (ss *SolarSystem) AddOrbit(orbit *Orbit) {
	ss.Orbits = append(ss.Orbits, orbit)
}

// CreateDefaultSolarSystem 创建默认太阳系（8大行星）
func CreateDefaultSolarSystem() *SolarSystem {
	ss := NewSolarSystem()

	// 定义行星数据
	planetsData := []struct {
		name, nameCN                                   string
		radius, orbitRadius, orbitSpeed, rotationSpeed float64
		color                                          [3]float64
		useGradient                                    bool
		gradientColor                                  [3]float64
		hasMoon, hasRings                              bool
		ringColors                                     [][3]float64
	}{
		// MUI 风格配色 - 使用 Material Design 色板
		{"Mercury", "水星", 0.15, 2.0, 4.0, 10.0, [3]float64{0.62, 0.62, 0.62}, true, [3]float64{0.38, 0.38, 0.38}, false, false, nil}, // Grey 500 -> 700
		{"Venus", "金星", 0.25, 3.0, 3.0, 8.0, [3]float64{1.0, 0.92, 0.23}, true, [3]float64{1.0, 0.76, 0.03}, false, false, nil},      // Amber 400 -> 600
		{"Earth", "地球", 0.28, 4.0, 2.0, 12.0, [3]float64{0.25, 0.59, 0.95}, true, [3]float64{0.13, 0.59, 0.95}, true, false, nil},    // Blue 500 -> 600
		{"Mars", "火星", 0.20, 5.0, 1.5, 11.0, [3]float64{1.0, 0.34, 0.13}, true, [3]float64{0.96, 0.26, 0.21}, false, false, nil},     // Deep Orange 500 -> 600
		{"Jupiter", "木星", 0.60, 7.0, 0.8, 15.0, [3]float64{1.0, 0.6, 0.0}, true, [3]float64{0.96, 0.49, 0.0}, false, false, nil},     // Orange 500 -> 700
		{"Saturn", "土星", 0.50, 9.0, 0.6, 13.0, [3]float64{1.0, 0.92, 0.23}, true, [3]float64{0.98, 0.84, 0.0}, false, true, [][3]float64{ // Yellow 400 -> 600
			{1.0, 0.95, 0.4}, {0.98, 0.89, 0.2}, {0.96, 0.84, 0.1},
		}},
		{"Uranus", "天王星", 0.35, 11.0, 0.4, 9.0, [3]float64{0.0, 0.74, 0.83}, true, [3]float64{0.0, 0.59, 0.65}, false, false, nil},    // Cyan 500 -> 700
		{"Neptune", "海王星", 0.33, 12.5, 0.3, 8.0, [3]float64{0.25, 0.32, 0.71}, true, [3]float64{0.16, 0.25, 0.63}, false, false, nil}, // Indigo 600 -> 800
	}

	// 添加行星和轨道
	for _, pd := range planetsData {
		planet := NewPlanet(pd.name, pd.nameCN, pd.radius, pd.orbitRadius, pd.orbitSpeed, pd.rotationSpeed, pd.color)

		if pd.useGradient {
			planet.SetGradient(pd.color, pd.gradientColor)
		}
		if pd.hasMoon {
			planet.AddMoon()
		}
		if pd.hasRings && len(pd.ringColors) > 0 {
			planet.AddRings(pd.ringColors)
		}

		ss.AddPlanet(planet)
		ss.AddOrbit(NewOrbit(pd.orbitRadius, [3]float64{0.26, 0.27, 0.29})) // MUI Grey 800 (更柔和的轨道线)
	}

	return ss
}

// Render 渲染太阳系
func (ss *SolarSystem) Render(renderer *Renderer, t float64) {
	// 渲染星空
	if ss.Stars != nil {
		ss.Stars.Render(renderer, t)
	}

	// 渲染太阳
	if ss.Sun != nil {
		ss.Sun.Render(renderer, t)
	}

	// 渲染轨道
	for _, orbit := range ss.Orbits {
		orbit.Render(renderer, t)
	}

	// 渲染行星
	for _, planet := range ss.Planets {
		planet.Render(renderer, t)
	}
}

// CelestialBody 天体（太阳、恒星等）
type CelestialBody struct {
	Name          string
	NameCN        string
	Radius        float64
	Color         [3]float64
	UseGradient   bool
	GradientColor [3]float64
	RotationSpeed float64
	Position      Vector3
}

// NewCelestialBody 创建天体
func NewCelestialBody(name, nameCN string, radius float64, color [3]float64) *CelestialBody {
	return &CelestialBody{
		Name:          name,
		NameCN:        nameCN,
		Radius:        radius,
		Color:         color,
		UseGradient:   false,
		RotationSpeed: 0,
		Position:      NewVector3(0, 0, 0),
	}
}

// SetGradient 设置渐变色
func (cb *CelestialBody) SetGradient(color1, color2 [3]float64) *CelestialBody {
	cb.UseGradient = true
	cb.Color = color1
	cb.GradientColor = color2
	return cb
}

// Render 渲染天体
func (cb *CelestialBody) Render(renderer *Renderer, t float64) {
	body := CreateSphere(cb.Radius, 20, 20)

	transform := Identity()
	transform = transform.Multiply(Translation(cb.Position.X, cb.Position.Y, cb.Position.Z))
	if cb.RotationSpeed != 0 {
		transform = transform.Multiply(RotationY(t * cb.RotationSpeed * 3.14159))
	}

	transformedBody := body.Transform(transform)

	if cb.UseGradient {
		renderer.DrawMeshWithGradient(transformedBody, cb.Color, cb.GradientColor)
	} else {
		renderer.DrawMesh(transformedBody, cb.Color)
	}

	// 渲染标签
	labelPos := NewVector3(cb.Position.X, cb.Position.Y+cb.Radius+0.5, cb.Position.Z)
	label := NewLabel3D(labelPos, cb.NameCN, [3]float64{1, 1, 1})
	label.Render(renderer, t)
}
