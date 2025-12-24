package go3d

import "math"

// Matrix4 表示4x4变换矩阵
type Matrix4 [16]float64

// Identity 返回单位矩阵
func Identity() Matrix4 {
	return Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Multiply 矩阵乘法
func (m Matrix4) Multiply(other Matrix4) Matrix4 {
	var result Matrix4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				result[i*4+j] += m[i*4+k] * other[k*4+j]
			}
		}
	}
	return result
}

// TransformVector 变换向量
func (m Matrix4) TransformVector(v Vector3) Vector3 {
	x := m[0]*v.X + m[1]*v.Y + m[2]*v.Z + m[3]
	y := m[4]*v.X + m[5]*v.Y + m[6]*v.Z + m[7]
	z := m[8]*v.X + m[9]*v.Y + m[10]*v.Z + m[11]
	w := m[12]*v.X + m[13]*v.Y + m[14]*v.Z + m[15]

	// 检查 w 是否接近零，避免除零错误
	if math.Abs(w) > 1e-10 && math.Abs(w-1.0) > 1e-10 {
		return Vector3{x / w, y / w, z / w}
	}
	return Vector3{x, y, z}
}

// Translation 创建平移矩阵
func Translation(x, y, z float64) Matrix4 {
	return Matrix4{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
}

// Scale 创建缩放矩阵
func Scale(x, y, z float64) Matrix4 {
	return Matrix4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

// RotationX 创建绕X轴旋转矩阵
func RotationX(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	}
}

// RotationY 创建绕Y轴旋转矩阵
func RotationY(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		c, 0, s, 0,
		0, 1, 0, 0,
		-s, 0, c, 0,
		0, 0, 0, 1,
	}
}

// RotationZ 创建绕Z轴旋转矩阵
func RotationZ(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Perspective 创建透视投影矩阵
func Perspective(fov, aspect, near, far float64) Matrix4 {
	// 防止除零和无效参数
	if math.Abs(fov) < 1e-10 || math.Abs(aspect) < 1e-10 || math.Abs(far-near) < 1e-10 {
		return Identity()
	}

	f := 1.0 / math.Tan(fov/2.0)
	return Matrix4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / (near - far), (2 * far * near) / (near - far),
		0, 0, -1, 0,
	}
}

// LookAt 创建视图矩阵
func LookAt(eye, target, up Vector3) Matrix4 {
	// 计算相机坐标系
	forward := target.Sub(eye).Normalize()
	if forward.Length() < 1e-10 {
		forward = Vector3{0, 0, 1}
	}

	right := forward.Cross(up).Normalize()
	if right.Length() < 1e-10 {
		// up 和 forward 平行，选择另一个 up 向量
		if math.Abs(forward.Y) < 0.9 {
			right = forward.Cross(Vector3{0, 1, 0}).Normalize()
		} else {
			right = forward.Cross(Vector3{1, 0, 0}).Normalize()
		}
	}

	newUp := right.Cross(forward).Normalize()

	// 构建视图矩阵
	return Matrix4{
		right.X, right.Y, right.Z, -right.Dot(eye),
		newUp.X, newUp.Y, newUp.Z, -newUp.Dot(eye),
		-forward.X, -forward.Y, -forward.Z, forward.Dot(eye),
		0, 0, 0, 1,
	}
}
