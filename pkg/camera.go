package go3d

import "math"

// CameraPath 相机路径接口
type CameraPath interface {
	GetPosition(t float64) Vector3
	GetTarget(t float64) Vector3
	GetFOV(t float64) float64
}

// CameraKeyframe 相机关键帧
type CameraKeyframe struct {
	Time     float64 // 时间点 (0-1)
	Position Vector3
	Target   Vector3
	FOV      float64
}

// InterpolatedCameraPath 插值相机路径
type InterpolatedCameraPath struct {
	Keyframes      []CameraKeyframe
	SmoothFunction func(float64) float64 // 平滑函数
}

// NewInterpolatedCameraPath 创建插值相机路径
func NewInterpolatedCameraPath(keyframes []CameraKeyframe) *InterpolatedCameraPath {
	return &InterpolatedCameraPath{
		Keyframes:      keyframes,
		SmoothFunction: Smoothstep, // 默认使用 smoothstep
	}
}

// GetPosition 获取指定时间的相机位置
func (cp *InterpolatedCameraPath) GetPosition(t float64) Vector3 {
	return cp.interpolateVector(t, func(kf CameraKeyframe) Vector3 { return kf.Position })
}

// GetTarget 获取指定时间的相机目标
func (cp *InterpolatedCameraPath) GetTarget(t float64) Vector3 {
	return cp.interpolateVector(t, func(kf CameraKeyframe) Vector3 { return kf.Target })
}

// GetFOV 获取指定时间的 FOV
func (cp *InterpolatedCameraPath) GetFOV(t float64) float64 {
	return cp.interpolateFloat(t, func(kf CameraKeyframe) float64 { return kf.FOV })
}

// interpolateVector 插值向量
func (cp *InterpolatedCameraPath) interpolateVector(t float64, getter func(CameraKeyframe) Vector3) Vector3 {
	if len(cp.Keyframes) == 0 {
		return NewVector3(0, 0, 0)
	}
	if len(cp.Keyframes) == 1 {
		return getter(cp.Keyframes[0])
	}

	// 找到相邻的两个关键帧
	var kf1, kf2 CameraKeyframe
	for i := 0; i < len(cp.Keyframes)-1; i++ {
		if t >= cp.Keyframes[i].Time && t <= cp.Keyframes[i+1].Time {
			kf1 = cp.Keyframes[i]
			kf2 = cp.Keyframes[i+1]
			break
		}
	}

	// 如果超出范围，使用首尾关键帧
	if t < cp.Keyframes[0].Time {
		return getter(cp.Keyframes[0])
	}
	if t > cp.Keyframes[len(cp.Keyframes)-1].Time {
		return getter(cp.Keyframes[len(cp.Keyframes)-1])
	}

	// 计算局部插值参数
	localT := (t - kf1.Time) / (kf2.Time - kf1.Time)
	if cp.SmoothFunction != nil {
		localT = cp.SmoothFunction(localT)
	}

	// 线性插值
	v1 := getter(kf1)
	v2 := getter(kf2)
	return v1.Scale(1 - localT).Add(v2.Scale(localT))
}

// interpolateFloat 插值浮点数
func (cp *InterpolatedCameraPath) interpolateFloat(t float64, getter func(CameraKeyframe) float64) float64 {
	if len(cp.Keyframes) == 0 {
		return 0
	}
	if len(cp.Keyframes) == 1 {
		return getter(cp.Keyframes[0])
	}

	var kf1, kf2 CameraKeyframe
	for i := 0; i < len(cp.Keyframes)-1; i++ {
		if t >= cp.Keyframes[i].Time && t <= cp.Keyframes[i+1].Time {
			kf1 = cp.Keyframes[i]
			kf2 = cp.Keyframes[i+1]
			break
		}
	}

	if t < cp.Keyframes[0].Time {
		return getter(cp.Keyframes[0])
	}
	if t > cp.Keyframes[len(cp.Keyframes)-1].Time {
		return getter(cp.Keyframes[len(cp.Keyframes)-1])
	}

	localT := (t - kf1.Time) / (kf2.Time - kf1.Time)
	if cp.SmoothFunction != nil {
		localT = cp.SmoothFunction(localT)
	}

	v1 := getter(kf1)
	v2 := getter(kf2)
	return v1*(1-localT) + v2*localT
}

// OrbitCameraPath 环绕相机路径
type OrbitCameraPath struct {
	Center       Vector3 // 环绕中心
	Radius       float64 // 环绕半径
	Height       float64 // 相机高度
	Speed        float64 // 环绕速度（圈数）
	FOV          float64 // 视场角
	HeightOffset func(float64) float64
	RadiusOffset func(float64) float64
}

// NewOrbitCameraPath 创建环绕相机路径
func NewOrbitCameraPath(center Vector3, radius, height, speed, fov float64) *OrbitCameraPath {
	return &OrbitCameraPath{
		Center: center,
		Radius: radius,
		Height: height,
		Speed:  speed,
		FOV:    fov,
	}
}

// GetPosition 获取相机位置
func (ocp *OrbitCameraPath) GetPosition(t float64) Vector3 {
	angle := t * ocp.Speed * 2 * math.Pi

	radius := ocp.Radius
	if ocp.RadiusOffset != nil {
		radius += ocp.RadiusOffset(t)
	}

	height := ocp.Height
	if ocp.HeightOffset != nil {
		height += ocp.HeightOffset(t)
	}

	return NewVector3(
		ocp.Center.X+radius*math.Cos(angle),
		ocp.Center.Y+height,
		ocp.Center.Z+radius*math.Sin(angle),
	)
}

// GetTarget 获取相机目标
func (ocp *OrbitCameraPath) GetTarget(t float64) Vector3 {
	return ocp.Center
}

// GetFOV 获取 FOV
func (ocp *OrbitCameraPath) GetFOV(t float64) float64 {
	return ocp.FOV
}

// Smoothstep 平滑插值函数
func Smoothstep(t float64) float64 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t * t * (3 - 2*t)
}

// Smootherstep 更平滑的插值函数
func Smootherstep(t float64) float64 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t * t * t * (t*(t*6-15) + 10)
}

// EaseInOut 缓入缓出函数
func EaseInOut(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - math.Pow(-2*t+2, 2)/2
}

// ApplyCameraPath 应用相机路径到渲染器
func ApplyCameraPath(renderer *Renderer, path CameraPath, t float64) {
	renderer.Camera.Position = path.GetPosition(t)
	renderer.Camera.Target = path.GetTarget(t)
	renderer.Camera.FOV = path.GetFOV(t)
}
