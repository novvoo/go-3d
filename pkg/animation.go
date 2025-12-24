package go3d

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// AnimationConfig 动画配置
type AnimationConfig struct {
	Width       int     // 宽度
	Height      int     // 高度
	FPS         int     // 帧率
	Duration    float64 // 时长（秒）
	OutputFile  string  // 输出文件名
	TempDir     string  // 临时目录
	Quality     int     // 视频质量 (CRF: 0-51, 越小质量越高)
	CleanupTemp bool    // 是否清理临时文件
	Workers     int     // 并行渲染的工作线程数（默认为1，单线程）
}

// DefaultAnimationConfig 返回默认动画配置
func DefaultAnimationConfig() AnimationConfig {
	return AnimationConfig{
		Width:       1920,
		Height:      1080,
		FPS:         30,
		Duration:    10.0,
		OutputFile:  "animation.mp4",
		TempDir:     "temp_frames",
		Quality:     23,
		CleanupTemp: true,
		Workers:     1, // 默认单线程
	}
}

// FrameRenderer 帧渲染函数类型
type FrameRenderer func(renderer *Renderer, frame int, t float64)

// AnimationGenerator 动画生成器
type AnimationGenerator struct {
	Config   AnimationConfig
	Renderer FrameRenderer
}

// NewAnimationGenerator 创建动画生成器
func NewAnimationGenerator(config AnimationConfig, renderer FrameRenderer) *AnimationGenerator {
	return &AnimationGenerator{
		Config:   config,
		Renderer: renderer,
	}
}

// CheckFFmpeg 检查系统是否安装了 ffmpeg
func CheckFFmpeg() bool {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	return err == nil
}

// GenerateFrames 生成所有帧
func (ag *AnimationGenerator) GenerateFrames() error {
	// 创建临时目录
	if err := os.MkdirAll(ag.Config.TempDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}

	totalFrames := int(float64(ag.Config.FPS) * ag.Config.Duration)
	workers := ag.Config.Workers
	if workers < 1 {
		workers = 1
	}

	fmt.Printf("生成 %d 帧动画 (%dx%d @ %d fps, %d 线程)...\n",
		totalFrames, ag.Config.Width, ag.Config.Height, ag.Config.FPS, workers)

	// 如果只有一个工作线程，使用单线程模式
	if workers == 1 {
		return ag.generateFramesSingleThread(totalFrames)
	}

	// 多线程模式
	return ag.generateFramesMultiThread(totalFrames, workers)
}

// generateFramesSingleThread 单线程生成帧
func (ag *AnimationGenerator) generateFramesSingleThread(totalFrames int) error {
	// 从帧1开始，跳过帧0
	for frame := 1; frame <= totalFrames; frame++ {
		t := float64(frame-1) / float64(totalFrames)

		// 创建渲染器
		renderer := NewRenderer(ag.Config.Width, ag.Config.Height)

		// 调用用户提供的渲染函数
		ag.Renderer(renderer, frame, t)

		// 保存帧，使用frame编号
		framePath := filepath.Join(ag.Config.TempDir, fmt.Sprintf("frame_%04d.png", frame))
		if err := renderer.SaveToPNG(framePath); err != nil {
			renderer.Destroy()
			return fmt.Errorf("保存帧 %d 失败: %w", frame, err)
		}
		renderer.Destroy()

		// 显示进度
		if frame%10 == 0 || frame == totalFrames {
			progress := float64(frame) / float64(totalFrames) * 100
			fmt.Printf("\r  进度: %.1f%% (%d/%d)", progress, frame, totalFrames)
		}
	}
	fmt.Println()
	return nil
}

// generateFramesMultiThread 多线程生成帧
func (ag *AnimationGenerator) generateFramesMultiThread(totalFrames, workers int) error {
	// 创建任务通道和错误通道
	jobs := make(chan int, totalFrames)
	errors := make(chan error, workers)

	// 用于进度显示的通道
	progress := make(chan int, totalFrames)

	var wg sync.WaitGroup

	// 启动工作协程
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for frame := range jobs {
				t := float64(frame-1) / float64(totalFrames)

				// 创建渲染器
				renderer := NewRenderer(ag.Config.Width, ag.Config.Height)

				// 调用用户提供的渲染函数
				ag.Renderer(renderer, frame, t)

				// 保存帧，使用frame编号
				framePath := filepath.Join(ag.Config.TempDir, fmt.Sprintf("frame_%04d.png", frame))
				if err := renderer.SaveToPNG(framePath); err != nil {
					renderer.Destroy()
					errors <- fmt.Errorf("保存帧 %d 失败: %w", frame, err)
					return
				}
				renderer.Destroy()

				// 报告进度
				progress <- 1
			}
		}(w)
	}

	// 启动进度显示协程
	done := make(chan bool)
	go func() {
		completed := 0
		for range progress {
			completed++
			if completed%10 == 0 || completed == totalFrames {
				percent := float64(completed) / float64(totalFrames) * 100
				fmt.Printf("\r  进度: %.1f%% (%d/%d)", percent, completed, totalFrames)
			}
			if completed == totalFrames {
				break
			}
		}
		fmt.Println()
		done <- true
	}()

	// 分发任务，从帧1开始
	for frame := 1; frame <= totalFrames; frame++ {
		jobs <- frame
	}
	close(jobs)

	// 等待所有工作完成
	wg.Wait()
	close(progress)
	close(errors)

	// 等待进度显示完成
	<-done

	// 检查是否有错误
	if len(errors) > 0 {
		return <-errors
	}

	return nil
}

// ComposeVideo 使用 ffmpeg 合成视频
func (ag *AnimationGenerator) ComposeVideo() error {
	fmt.Println("\n使用 ffmpeg 合成视频...")

	cmd := exec.Command("ffmpeg",
		"-y",
		"-framerate", fmt.Sprintf("%d", ag.Config.FPS),
		"-start_number", "1", // 从帧1开始
		"-i", filepath.Join(ag.Config.TempDir, "frame_%04d.png"),
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-crf", fmt.Sprintf("%d", ag.Config.Quality),
		ag.Config.OutputFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg 错误: %w\n输出: %s", err, string(output))
	}

	fmt.Printf("\n✓ 动画已生成: %s\n", ag.Config.OutputFile)
	fmt.Printf("  分辨率: %dx%d\n", ag.Config.Width, ag.Config.Height)
	fmt.Printf("  帧率: %d fps\n", ag.Config.FPS)
	fmt.Printf("  时长: %.1f 秒\n", ag.Config.Duration)

	return nil
}

// Generate 生成完整动画（帧 + 视频）
func (ag *AnimationGenerator) Generate() error {
	// 生成帧
	if err := ag.GenerateFrames(); err != nil {
		return err
	}

	// 合成视频
	if err := ag.ComposeVideo(); err != nil {
		return err
	}

	// 清理临时文件
	if ag.Config.CleanupTemp {
		if err := os.RemoveAll(ag.Config.TempDir); err != nil {
			fmt.Printf("警告: 清理临时文件失败: %v\n", err)
		}
	}

	return nil
}

// GenerateFramesOnly 仅生成帧序列（不合成视频）
func (ag *AnimationGenerator) GenerateFramesOnly(outputDir string) error {
	// 临时修改配置
	originalTempDir := ag.Config.TempDir
	ag.Config.TempDir = outputDir
	ag.Config.CleanupTemp = false

	defer func() {
		ag.Config.TempDir = originalTempDir
	}()

	if err := ag.GenerateFrames(); err != nil {
		return err
	}

	fmt.Printf("\n✓ 序列帧已生成到目录: %s\n", outputDir)
	fmt.Println("\n要生成视频，请安装 ffmpeg 并运行:")
	fmt.Printf("  ffmpeg -framerate %d -i %s/frame_%%04d.png -c:v libx264 -pix_fmt yuv420p animation.mp4\n",
		ag.Config.FPS, outputDir)

	return nil
}
