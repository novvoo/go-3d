# Go-3D 太阳系动画生成器

一个使用 Go 语言编写的 3D 太阳系动画渲染库，支持生成 PNG 序列帧或直接输出 MP4 视频。

## 功能特性

- 🌍 完整的太阳系模拟（太阳、八大行星及其轨道）
- 🎥 动态相机系统，支持多种视角切换
- 💡 多光源渲染系统
- 🎨 Material UI 风格的渐变背景
- 📝 3D 文字标签支持
- 🎬 自动检测 ffmpeg 并生成 MP4 视频
- 📸 支持导出 PNG 序列帧

## 演示效果

![演示动画](animation.mp4)

> 注：如果视频无法播放，请[点击这里下载查看](animation.mp4)

## 系统要求

- Go 1.24.4 或更高版本
- ffmpeg（可选，用于生成 MP4 视频）

## FFmpeg 安装

### Windows

**方法 1: 使用 Chocolatey（推荐）**
```cmd
choco install ffmpeg
```

**方法 2: 使用 Scoop**
```cmd
scoop install ffmpeg
```

**方法 3: 手动安装**
1. 访问 [FFmpeg 官网](https://ffmpeg.org/download.html)
2. 下载 Windows 构建版本（推荐 [gyan.dev](https://www.gyan.dev/ffmpeg/builds/)）
3. 解压到目录（如 `C:\ffmpeg`）
4. 将 `C:\ffmpeg\bin` 添加到系统环境变量 PATH
5. 打开新的命令行窗口，运行 `ffmpeg -version` 验证安装

### macOS

**方法 1: 使用 Homebrew（推荐）**
```bash
brew install ffmpeg
```

**方法 2: 使用 MacPorts**
```bash
sudo port install ffmpeg
```

### Linux

**Ubuntu/Debian**
```bash
sudo apt update
sudo apt install ffmpeg
```

**Fedora**
```bash
sudo dnf install ffmpeg
```

**Arch Linux**
```bash
sudo pacman -S ffmpeg
```

**CentOS/RHEL**
```bash
# 启用 EPEL 和 RPM Fusion 仓库
sudo yum install epel-release
sudo yum localinstall --nogpgcheck https://download1.rpmfusion.org/free/el/rpmfusion-free-release-7.noarch.rpm
sudo yum install ffmpeg
```

### 验证安装

安装完成后，在终端运行以下命令验证：
```bash
ffmpeg -version
```

如果显示版本信息，说明安装成功。

## 安装

```bash
go get github.com/novvoo/go-3d
```

## 快速开始

### 生成 MP4 动画（需要 ffmpeg）

```go
package main

import (
    "fmt"
    go3d "github.com/novvoo/go-3d/pkg"
)

func main() {
    // 配置动画参数
    config := go3d.DefaultAnimationConfig()
    config.Duration = 10.0  // 10秒动画
    config.FPS = 30         // 30帧/秒
    
    // 创建动画生成器
    generator := go3d.NewAnimationGenerator(config, renderFrame)
    
    // 生成动画
    if err := generator.Generate(); err != nil {
        fmt.Printf("生成动画失败: %v\n", err)
    }
}

func renderFrame(renderer *go3d.Renderer, frame int, t float64) {
    // 设置相机
    renderer.Camera.Position = go3d.NewVector3(15, 10, 15)
    renderer.Camera.Target = go3d.NewVector3(0, 0, 0)
    
    // 创建场景
    scene := go3d.NewScene()
    
    // 添加太阳系
    solarSystem := go3d.CreateDefaultSolarSystem()
    scene.AddObject(solarSystem)
    
    // 渲染
    scene.Render(renderer, t)
}
```

### 仅生成 PNG 序列帧（不需要 ffmpeg）

```go
config := go3d.DefaultAnimationConfig()
generator := go3d.NewAnimationGenerator(config, renderFrame)

// 生成帧序列到指定目录
if err := generator.GenerateFramesOnly("output_frames"); err != nil {
    fmt.Printf("生成帧序列失败: %v\n", err)
}
```

## 运行示例

```bash
# 克隆仓库
git clone https://github.com/novvoo/go-3d.git
cd go-3d

# 运行示例（自动检测 ffmpeg）
go run example/animation.go
```

程序会自动检测系统中是否安装了 ffmpeg：
- ✅ 如果已安装：直接生成 `animation.mp4` 视频文件
- ⚠️ 如果未安装：生成 PNG 序列帧到 `animation_frames` 目录

## 项目结构

```
go-3d/
├── pkg/                    # 核心库代码
│   ├── animation.go       # 动画生成器
│   ├── camera.go          # 相机系统
│   ├── celestial.go       # 天体对象
│   ├── matrix4.go         # 4x4 矩阵运算
│   ├── mesh.go            # 网格渲染
│   ├── orbit.go           # 轨道系统
│   ├── renderer.go        # 渲染器
│   ├── scene.go           # 场景管理
│   ├── solarsystem.go     # 太阳系配置
│   └── vector3.go         # 3D 向量运算
├── example/               # 示例代码
│   └── animation.go       # 太阳系动画示例
└── README.md
```

## API 文档

### 动画配置

```go
type AnimationConfig struct {
    Width      int     // 画面宽度（默认：1920）
    Height     int     // 画面高度（默认：1080）
    FPS        int     // 帧率（默认：30）
    Duration   float64 // 动画时长（秒，默认：5.0）
    OutputFile string  // 输出文件名（默认：animation.mp4）
}
```

### 相机控制

```go
camera := renderer.Camera
camera.Position = go3d.NewVector3(x, y, z)  // 相机位置
camera.Target = go3d.NewVector3(x, y, z)    // 观察目标
camera.FOV = 0.8                             // 视场角
```

### 光源系统

```go
light := go3d.NewLight(
    go3d.NewVector3(x, y, z),      // 光源位置
    [3]float64{r, g, b},            // 光源颜色（0-1）
    intensity,                       // 光照强度
)
renderer.AddLight(light)
```

## 常见问题

### Q: 为什么生成的是 PNG 序列帧而不是视频？
A: 这说明系统中未安装 ffmpeg。请参考上面的 [FFmpeg 安装](#ffmpeg-安装) 部分进行安装。

### Q: 如何手动将 PNG 序列帧转换为视频？
A: 使用以下 ffmpeg 命令：
```bash
ffmpeg -framerate 30 -i frame_%04d.png -c:v libx264 -pix_fmt yuv420p output.mp4
```

### Q: 如何调整视频质量？
A: 修改 `AnimationConfig` 中的 `Width` 和 `Height` 参数，或在生成视频时调整 ffmpeg 参数。

### Q: 支持哪些渲染模式？
A: 目前支持：
- `RenderWireframe` - 线框模式
- `RenderShaded` - 着色模式（默认）
- `RenderTextured` - 纹理模式（开发中）

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 相关链接

- [FFmpeg 官网](https://ffmpeg.org/)
- [Go Cairo 绑定](https://github.com/novvoo/go-cairo)
