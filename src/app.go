package src

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	CubePoints      = make([]*Vec3, N_Points)
	ProjectedPoints = make([]*Vec2, N_Points)
	CameraPosition  = Vec3{x: 0, y: 0, z: -5}
	CubeRotation    = Vec3{x: 0, y: 0, z: 0}

	timePreviousFrame = uint64(0)
)

const (
	FramePerSecond     = 60
	MilisecondPerFrame = FramePerSecond / 1000
	N_Points           = 9 * 9 * 9
	FovFactor          = 640
)

type App struct {
	isRunning         bool
	w_width, w_height int32
	window            *sdl.Window
	renderer          *sdl.Renderer

	colorBuffer        []uint32
	colorBufferTexture *sdl.Texture
}

func (a *App) Setup() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	displayMode, err := sdl.GetCurrentDisplayMode(0)
	if err != nil {
		return err
	}
	a.w_width = displayMode.W
	a.w_height = displayMode.H

	window, err := sdl.CreateWindow("3D Renderer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		a.w_width, a.w_height, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	a.window = window

	renderer, err := sdl.CreateRenderer(a.window, -1, 0)
	if err != nil {
		return err
	}
	a.renderer = renderer

	a.colorBuffer = make([]uint32, a.w_width*a.w_height)

	colorBufferTexture, err := a.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, a.w_width, a.w_height)
	if err != nil {
		return err
	}
	a.colorBufferTexture = colorBufferTexture
	a.isRunning = true

	point_count := 0

	for x := float64(-1); x <= 1; x += 0.25 {
		for y := float64(-1); y <= 1; y += 0.25 {
			for z := float64(-1); z <= 1; z += 0.25 {
				CubePoints[point_count] = &Vec3{
					x: x,
					y: y,
					z: z,
				}
				point_count++
			}
		}
	}

	return nil
}

func (a *App) IsRunning() bool {
	return a.isRunning
}

func (a *App) ProcessInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			fmt.Println("Quit")
			a.isRunning = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				fmt.Println("Quit(ESCAPE)")
				a.isRunning = false
			}
		}
	}
}

func (a *App) Update() {
	timeToWait := MilisecondPerFrame - (sdl.GetTicks64() - timePreviousFrame)
	if timeToWait > 0 && timeToWait < MilisecondPerFrame {
		sdl.Delay(uint32(timeToWait))
	}
	timePreviousFrame = sdl.GetTicks64()

	CubeRotation.x += 0.01
	CubeRotation.y += 0.01
	CubeRotation.z += 0.01

	for i := 0; i < N_Points; i++ {
		point := CubePoints[i]
		transformPoint := point.RotateX(CubeRotation.x)
		transformPoint = transformPoint.RotateY(CubeRotation.y)
		transformPoint = transformPoint.RotateZ(CubeRotation.z)

		transformPoint.z -= CameraPosition.z

		ProjectedPoints[i] = &Vec2{
			x: float64(FovFactor) * transformPoint.x / transformPoint.z,
			y: float64(FovFactor) * transformPoint.y / transformPoint.z,
		}
	}

}

func (a *App) Render() {
	a.DrawGrid()

	for i := 0; i < N_Points; i++ {
		a.DrawRect(
			int32(ProjectedPoints[i].x)+(a.w_width/2),
			int32(ProjectedPoints[i].y)+(a.w_height/2),
			4,
			4,
			0xFF00FFFF,
		)
	}

	a.RenderColorBuffer()
	a.ClearColorBuffer(0xFF000000)

	a.renderer.Present()
}

func (a *App) Destroy() {
	a.renderer.Destroy()
	a.window.Destroy()
	sdl.Quit()
}
