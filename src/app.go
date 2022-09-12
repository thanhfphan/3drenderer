package src

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	ProjectedPoints = []*Vec2{}
	mesh            = &Mesh{Vertices: []*Vec3{}, Faces: []*Triangle{}}
	CameraPosition  = Vec3{x: 0, y: 0, z: -5}
	CubeRotation    = Vec3{x: 0, y: 0, z: 0}

	timePreviousFrame = uint64(0)
)

const (
	FramePerSecond      = 60
	MilisecondsPerFrame = 1000 / FramePerSecond
	N_Points            = 9 * 9 * 9
	FovFactor           = 640
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

	rect, err := sdl.GetDisplayBounds(1)
	if err != nil {
		return err
	}
	displayMode, err := sdl.GetCurrentDisplayMode(1)
	if err != nil {
		return err
	}
	a.w_width = displayMode.W
	a.w_height = displayMode.H

	window, err := sdl.CreateWindow("3D Renderer", rect.X, rect.Y,
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

	mesh.Vertices = append(mesh.Vertices,
		&Vec3{-1, -1, -1}, // 1
		&Vec3{-1, 1, -1},  // 2
		&Vec3{1, 1, -1},   // 3
		&Vec3{1, -1, -1},  // 4
		&Vec3{1, 1, 1},    // 5
		&Vec3{1, -1, 1},   // 6
		&Vec3{-1, 1, 1},   // 7
		&Vec3{-1, -1, 1},  // 8
	)

	mesh.Faces = append(mesh.Faces,
		&Triangle{A: 1, B: 2, C: 3},
		&Triangle{A: 1, B: 3, C: 4},
		//
		&Triangle{A: 4, B: 3, C: 6},
		&Triangle{A: 6, B: 3, C: 5},
		//
		&Triangle{A: 8, B: 6, C: 5},
		&Triangle{A: 8, B: 5, C: 7},
		//
		&Triangle{A: 1, B: 8, C: 7},
		&Triangle{A: 1, B: 7, C: 2},
		//
		&Triangle{A: 6, B: 1, C: 4},
		&Triangle{A: 6, B: 8, C: 1},
	)

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
	timeToWait := MilisecondsPerFrame - (sdl.GetTicks64() - timePreviousFrame)
	if timeToWait > 0 && timeToWait <= MilisecondsPerFrame {
		sdl.Delay(uint32(timeToWait))
	}
	timePreviousFrame = sdl.GetTicks64()

	CubeRotation.x += 0.01
	CubeRotation.y += 0.01
	CubeRotation.z += 0.01

	ProjectedPoints = []*Vec2{}
	for _, item := range mesh.Vertices {
		transformPoint := item.RotateX(CubeRotation.x)
		transformPoint = transformPoint.RotateY(CubeRotation.y)
		transformPoint = transformPoint.RotateZ(CubeRotation.z)

		transformPoint.z -= CameraPosition.z

		projectdPoint := &Vec2{
			x: float64(FovFactor) * transformPoint.x / transformPoint.z,
			y: float64(FovFactor) * transformPoint.y / transformPoint.z,
		}

		projectdPoint.x += float64(a.w_width) / 2
		projectdPoint.y += float64(a.w_height) / 2
		ProjectedPoints = append(ProjectedPoints, projectdPoint)
	}

}

func (a *App) Render() {
	a.DrawGrid()

	for _, item := range ProjectedPoints {
		a.DrawRect(
			int32(item.x),
			int32(item.y),
			4,
			4,
			0xFF00FFFF,
		)
	}

	for _, item := range mesh.Faces {
		pA := ProjectedPoints[item.A-1]
		pB := ProjectedPoints[item.B-1]
		pC := ProjectedPoints[item.C-1]
		a.DrawTriangle(pA.x, pA.y, pB.x, pB.y, pC.x, pC.y, 0xFF00FFFF)
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
