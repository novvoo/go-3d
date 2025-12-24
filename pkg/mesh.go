package go3d

import "math"

// Triangle 表示3D三角形
type Triangle struct {
	V0, V1, V2 Vector3
}

// Normal 计算三角形法线
func (t Triangle) Normal() Vector3 {
	edge1 := t.V1.Sub(t.V0)
	edge2 := t.V2.Sub(t.V0)
	normal := edge1.Cross(edge2)
	// 检查退化三角形
	if normal.Length() < 1e-10 {
		return Vector3{0, 1, 0} // 返回默认法线
	}
	return normal.Normalize()
}

// Center 计算三角形中心
func (t Triangle) Center() Vector3 {
	return t.V0.Add(t.V1).Add(t.V2).Scale(1.0 / 3.0)
}

// Mesh 表示3D网格
type Mesh struct {
	Vertices  []Vector3
	Triangles []Triangle
}

// NewMesh 创建新网格
func NewMesh() *Mesh {
	return &Mesh{
		Vertices:  make([]Vector3, 0),
		Triangles: make([]Triangle, 0),
	}
}

// AddVertex 添加顶点
func (m *Mesh) AddVertex(v Vector3) {
	m.Vertices = append(m.Vertices, v)
}

// AddTriangle 添加三角形
func (m *Mesh) AddTriangle(t Triangle) {
	m.Triangles = append(m.Triangles, t)
}

// Transform 变换网格
func (m *Mesh) Transform(matrix Matrix4) *Mesh {
	transformed := NewMesh()
	for _, v := range m.Vertices {
		transformed.AddVertex(matrix.TransformVector(v))
	}
	for _, t := range m.Triangles {
		transformed.AddTriangle(Triangle{
			V0: matrix.TransformVector(t.V0),
			V1: matrix.TransformVector(t.V1),
			V2: matrix.TransformVector(t.V2),
		})
	}
	return transformed
}

// Merge 合并多个网格
func (m *Mesh) Merge(other *Mesh) {
	m.Vertices = append(m.Vertices, other.Vertices...)
	m.Triangles = append(m.Triangles, other.Triangles...)
}

// CreateCube 创建立方体网格
func CreateCube(size float64) *Mesh {
	mesh := NewMesh()
	s := size / 2.0

	// 立方体的8个顶点
	vertices := []Vector3{
		{-s, -s, -s}, {s, -s, -s}, {s, s, -s}, {-s, s, -s},
		{-s, -s, s}, {s, -s, s}, {s, s, s}, {-s, s, s},
	}

	// 12个三角形（6个面，每面2个三角形）
	indices := [][3]int{
		{0, 1, 2}, {0, 2, 3}, // 前面
		{5, 4, 7}, {5, 7, 6}, // 后面
		{4, 0, 3}, {4, 3, 7}, // 左面
		{1, 5, 6}, {1, 6, 2}, // 右面
		{3, 2, 6}, {3, 6, 7}, // 上面
		{4, 5, 1}, {4, 1, 0}, // 下面
	}

	for _, idx := range indices {
		mesh.AddTriangle(Triangle{
			V0: vertices[idx[0]],
			V1: vertices[idx[1]],
			V2: vertices[idx[2]],
		})
	}

	return mesh
}

// CreateSphere 创建球体网格
func CreateSphere(radius float64, segments, rings int) *Mesh {
	mesh := NewMesh()

	for ring := 0; ring <= rings; ring++ {
		theta := float64(ring) * math.Pi / float64(rings)
		sinTheta := math.Sin(theta)
		cosTheta := math.Cos(theta)

		for seg := 0; seg <= segments; seg++ {
			phi := float64(seg) * 2.0 * math.Pi / float64(segments)
			sinPhi := math.Sin(phi)
			cosPhi := math.Cos(phi)

			x := cosPhi * sinTheta
			y := cosTheta
			z := sinPhi * sinTheta

			mesh.AddVertex(NewVector3(x*radius, y*radius, z*radius))
		}
	}

	// 创建三角形
	for ring := 0; ring < rings; ring++ {
		for seg := 0; seg < segments; seg++ {
			first := ring*(segments+1) + seg
			second := first + segments + 1

			mesh.AddTriangle(Triangle{
				V0: mesh.Vertices[first],
				V1: mesh.Vertices[second],
				V2: mesh.Vertices[first+1],
			})

			mesh.AddTriangle(Triangle{
				V0: mesh.Vertices[second],
				V1: mesh.Vertices[second+1],
				V2: mesh.Vertices[first+1],
			})
		}
	}

	return mesh
}

// CreateCylinder 创建圆柱体网格
func CreateCylinder(radius, height float64, segments int) *Mesh {
	mesh := NewMesh()
	halfHeight := height / 2.0

	// 顶部和底部圆心
	topCenter := NewVector3(0, halfHeight, 0)
	bottomCenter := NewVector3(0, -halfHeight, 0)

	// 创建顶部和底部的顶点
	for i := 0; i <= segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(segments)
		x := math.Cos(angle) * radius
		z := math.Sin(angle) * radius

		mesh.AddVertex(NewVector3(x, halfHeight, z))
		mesh.AddVertex(NewVector3(x, -halfHeight, z))
	}

	// 侧面三角形
	for i := 0; i < segments; i++ {
		topIdx := i * 2
		bottomIdx := i*2 + 1
		nextTopIdx := (i + 1) * 2 % ((segments + 1) * 2)
		nextBottomIdx := ((i+1)*2 + 1) % ((segments + 1) * 2)

		mesh.AddTriangle(Triangle{
			V0: mesh.Vertices[topIdx],
			V1: mesh.Vertices[bottomIdx],
			V2: mesh.Vertices[nextTopIdx],
		})

		mesh.AddTriangle(Triangle{
			V0: mesh.Vertices[bottomIdx],
			V1: mesh.Vertices[nextBottomIdx],
			V2: mesh.Vertices[nextTopIdx],
		})
	}

	// 顶部和底部的三角形
	for i := 0; i < segments; i++ {
		topIdx := i * 2
		nextTopIdx := (i + 1) * 2 % ((segments + 1) * 2)
		bottomIdx := i*2 + 1
		nextBottomIdx := ((i+1)*2 + 1) % ((segments + 1) * 2)

		mesh.AddTriangle(Triangle{
			V0: topCenter,
			V1: mesh.Vertices[topIdx],
			V2: mesh.Vertices[nextTopIdx],
		})

		mesh.AddTriangle(Triangle{
			V0: bottomCenter,
			V1: mesh.Vertices[nextBottomIdx],
			V2: mesh.Vertices[bottomIdx],
		})
	}

	return mesh
}

// CreatePlane 创建平面网格
func CreatePlane(width, height float64, subdivisions int) *Mesh {
	mesh := NewMesh()
	w := width / 2.0
	h := height / 2.0
	step := 1.0 / float64(subdivisions)

	// 创建顶点
	for i := 0; i <= subdivisions; i++ {
		for j := 0; j <= subdivisions; j++ {
			x := -w + float64(j)*width*step
			z := -h + float64(i)*height*step
			mesh.AddVertex(NewVector3(x, 0, z))
		}
	}

	// 创建三角形
	for i := 0; i < subdivisions; i++ {
		for j := 0; j < subdivisions; j++ {
			idx := i*(subdivisions+1) + j
			mesh.AddTriangle(Triangle{
				V0: mesh.Vertices[idx],
				V1: mesh.Vertices[idx+subdivisions+1],
				V2: mesh.Vertices[idx+1],
			})
			mesh.AddTriangle(Triangle{
				V0: mesh.Vertices[idx+1],
				V1: mesh.Vertices[idx+subdivisions+1],
				V2: mesh.Vertices[idx+subdivisions+2],
			})
		}
	}

	return mesh
}

// CreateCone 创建圆锥体网格
func CreateCone(radius, height float64, segments int) *Mesh {
	mesh := NewMesh()
	halfHeight := height / 2.0

	// 顶点
	apex := NewVector3(0, halfHeight, 0)
	bottomCenter := NewVector3(0, -halfHeight, 0)

	// 底部圆的顶点
	for i := 0; i <= segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(segments)
		x := math.Cos(angle) * radius
		z := math.Sin(angle) * radius
		mesh.AddVertex(NewVector3(x, -halfHeight, z))
	}

	// 侧面三角形（从顶点到底部圆）
	for i := 0; i < segments; i++ {
		mesh.AddTriangle(Triangle{
			V0: apex,
			V1: mesh.Vertices[i],
			V2: mesh.Vertices[(i+1)%(segments+1)],
		})
	}

	// 底面三角形
	for i := 0; i < segments; i++ {
		mesh.AddTriangle(Triangle{
			V0: bottomCenter,
			V1: mesh.Vertices[(i+1)%(segments+1)],
			V2: mesh.Vertices[i],
		})
	}

	return mesh
}

// CreateTorus 创建圆环网格
func CreateTorus(majorRadius, minorRadius float64, majorSegments, minorSegments int) *Mesh {
	mesh := NewMesh()

	for i := 0; i <= majorSegments; i++ {
		theta := float64(i) * 2.0 * math.Pi / float64(majorSegments)
		cosTheta := math.Cos(theta)
		sinTheta := math.Sin(theta)

		for j := 0; j <= minorSegments; j++ {
			phi := float64(j) * 2.0 * math.Pi / float64(minorSegments)
			cosPhi := math.Cos(phi)
			sinPhi := math.Sin(phi)

			x := (majorRadius + minorRadius*cosPhi) * cosTheta
			y := minorRadius * sinPhi
			z := (majorRadius + minorRadius*cosPhi) * sinTheta

			mesh.AddVertex(NewVector3(x, y, z))
		}
	}

	// 创建三角形
	for i := 0; i < majorSegments; i++ {
		for j := 0; j < minorSegments; j++ {
			first := i*(minorSegments+1) + j
			second := first + minorSegments + 1

			mesh.AddTriangle(Triangle{
				V0: mesh.Vertices[first],
				V1: mesh.Vertices[second],
				V2: mesh.Vertices[first+1],
			})

			mesh.AddTriangle(Triangle{
				V0: mesh.Vertices[second],
				V1: mesh.Vertices[second+1],
				V2: mesh.Vertices[first+1],
			})
		}
	}

	return mesh
}
