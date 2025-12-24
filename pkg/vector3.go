package go3d

import "math"

// Vector3 表示3D空间中的向量
type Vector3 struct {
	X, Y, Z float64
}

// NewVector3 创建新的3D向量
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

// Add 向量加法
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Sub 向量减法
func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// Scale 向量缩放
func (v Vector3) Scale(s float64) Vector3 {
	return Vector3{v.X * s, v.Y * s, v.Z * s}
}

// Dot 点积
func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross 叉积
func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X,
	}
}

// Length 向量长度
func (v Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize 归一化向量
func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if length < 1e-10 { // 使用更安全的阈值
		return Vector3{0, 0, 0}
	}
	return v.Scale(1.0 / length)
}
