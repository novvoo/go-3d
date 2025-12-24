//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"math"

	go3d "github.com/novvoo/go-3d/pkg"
)

func main() {
	// 检查 ffmpeg 是否存在
	if !go3d.CheckFFmpeg() {
		fmt.Println("未检测到 ffmpeg，将生成 PNG 序列帧")
		fmt.Println("提示: 安装 ffmpeg 可以直接生成 MP4 视频")
		generateFrames()
		return
	}

	fmt.Println("检测到 ffmpeg，将直接生成 MP4 视频")
	generateMP4Animation()
}

// generateMP4Animation 生成 MP4 动画
func generateMP4Animation() {
	// 配置动画参数
	config := go3d.DefaultAnimationConfig()
	config.Duration = 10.0
	config.FPS = 30
	config.Workers = 5 // 使用 5 个线程并行渲染

	// 创建动画生成器
	generator := go3d.NewAnimationGenerator(config, renderFrame)

	// 生成动画
	if err := generator.Generate(); err != nil {
		fmt.Printf("生成动画失败: %v\n", err)
	}
}

// generateFrames 仅生成 PNG 序列帧
func generateFrames() {
	config := go3d.DefaultAnimationConfig()
	config.Duration = 10.0
	config.FPS = 30
	config.Workers = 1 // 改为单线程，避免多线程竞态问题

	generator := go3d.NewAnimationGenerator(config, renderFrame)

	if err := generator.GenerateFramesOnly("animation_frames"); err != nil {
		fmt.Printf("生成帧序列失败: %v\n", err)
	}
}

// renderFrame 帧渲染函数
func renderFrame(renderer *go3d.Renderer, frame int, t float64) {
	// 设置动态相机
	setupDynamicCamera(renderer, t)

	// 添加光源
	light1 := go3d.NewLight(
		go3d.NewVector3(-5, 8, -5),
		[3]float64{1.0, 0.9, 0.8},
		0.8,
	)
	light2 := go3d.NewLight(
		go3d.NewVector3(5, 5, 5),
		[3]float64{0.6, 0.7, 1.0},
		0.6,
	)
	renderer.AddLight(light1)
	renderer.AddLight(light2)
	renderer.SetRenderMode(go3d.RenderShaded)

	// 创建场景
	scene := go3d.NewScene()

	// 设置渐变背景 - MUI 风格：使用深色主题配色
	background := go3d.NewGradientBackground(
		[3]float64{0.08, 0.09, 0.12}, // MUI Dark Background (Grey 900)
		[3]float64{0.15, 0.16, 0.20}, // MUI Dark Paper (Grey 800)
	)
	background.Animated = true
	scene.SetBackground(background)

	// 添加太阳系
	solarSystem := go3d.CreateDefaultSolarSystem()
	scene.AddObject(solarSystem)

	// 添加坐标系统
	coordSystem := go3d.NewCoordinateSystem(5.0)
	scene.AddObject(coordSystem)

	// 添加标题标签 - MUI 风格：使用 Primary 和 Secondary 颜色
	// 临时注释掉以测试是否是标签导致矩形
	/*
		titleLabel := go3d.NewLabel3D(
			go3d.NewVector3(0, 6.0+math.Sin(t*2*math.Pi)*0.3, 0),
			"太阳系",
			[3]float64{0.56, 0.93, 0.56}, // MUI Light Green A200
		)
		titleLabel.FontSize = 28
		scene.AddObject(titleLabel)

		subtitleLabel := go3d.NewLabel3D(
			go3d.NewVector3(0, 5.5+math.Sin(t*2*math.Pi)*0.3, 0),
			"Solar System",
			[3]float64{0.74, 0.76, 0.78}, // MUI Grey 400 (Secondary Text)
		)
		subtitleLabel.FontSize = 20
		scene.AddObject(subtitleLabel)
	*/

	// 渲染场景 - 使用加速的时间让行星运动更明显
	// t 是 0-1 的归一化时间，乘以一个系数让行星运动更快
	animationTime := t * 3.0 // 3倍速度，让行星在动画期间完成更多轨道运动
	scene.Render(renderer, animationTime)
}

// setupDynamicCamera 设置动态相机视角 - 同时绕X、Y、Z三个轴旋转
func setupDynamicCamera(renderer *go3d.Renderer, t float64) {
	// 基础半径
	baseRadius := 20.0

	// 三个轴独立的旋转角度，使用不同的频率让运动更丰富
	// Y轴旋转（水平环绕）：完整旋转一圈
	angleY := t * 2 * math.Pi

	// X轴旋转（垂直环绕）：上下大幅度旋转
	angleX := t * 1.5 * math.Pi // 旋转270度

	// Z轴旋转（前后环绕）：前后方向旋转
	angleZ := t * 1.0 * math.Pi // 旋转180度

	// 使用欧拉角计算相机位置
	// 从初始位置 (0, 0, baseRadius) 开始，依次应用三个轴的旋转

	// 初始位置：相机在Z轴正方向
	x := 0.0
	y := 0.0
	z := baseRadius

	// 应用X轴旋转（绕X轴旋转会改变Y和Z）
	cosX := math.Cos(angleX)
	sinX := math.Sin(angleX)
	newY := y*cosX - z*sinX
	newZ := y*sinX + z*cosX
	y = newY
	z = newZ

	// 应用Y轴旋转（绕Y轴旋转会改变X和Z）
	cosY := math.Cos(angleY)
	sinY := math.Sin(angleY)
	newX := x*cosY + z*sinY
	newZ = -x*sinY + z*cosY
	x = newX
	z = newZ

	// 应用Z轴旋转（绕Z轴旋转会改变X和Y）
	cosZ := math.Cos(angleZ)
	sinZ := math.Sin(angleZ)
	newX = x*cosZ - y*sinZ
	newY = x*sinZ + y*cosZ
	x = newX
	y = newY

	cameraPos := go3d.NewVector3(x, y, z)

	// 相机目标：始终看向太阳系中心
	targetPos := go3d.NewVector3(0, 0, 0)

	// 计算相机的上方向向量，让它随着相机旋转
	// 初始上方向是 (0, 1, 0)
	upX := 0.0
	upY := 1.0
	upZ := 0.0

	// 应用相同的旋转变换到上方向向量
	// X轴旋转
	newUpY := upY*cosX - upZ*sinX
	newUpZ := upY*sinX + upZ*cosX
	upY = newUpY
	upZ = newUpZ

	// Y轴旋转
	newUpX := upX*cosY + upZ*sinY
	newUpZ = -upX*sinY + upZ*cosY
	upX = newUpX
	upZ = newUpZ

	// Z轴旋转
	newUpX = upX*cosZ - upY*sinZ
	newUpY = upX*sinZ + upY*cosZ
	upX = newUpX
	upY = newUpY

	upVector := go3d.NewVector3(upX, upY, upZ)

	// 视场角
	fov := 0.75

	renderer.Camera.Position = cameraPos
	renderer.Camera.Target = targetPos
	renderer.Camera.Up = upVector
	renderer.Camera.FOV = fov
}
