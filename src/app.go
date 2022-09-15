package src

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	Triangles = []*Triangle{}

	mesh           = &Mesh{Vertices: []*Vec3{}, Faces: []*Face{}}
	CameraPosition = Vec3{X: 0, Y: 0, Z: 0}
	CubeRotation   = Vec3{X: 0, Y: 0, Z: 0}

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

	numMonitor, err := sdl.GetNumVideoDisplays()
	if err != nil {
		return err
	}

	displayAt := 0
	if numMonitor > 0 {
		displayAt = 1
	}

	rect, err := sdl.GetDisplayBounds(displayAt)
	if err != nil {
		return err
	}
	displayMode, err := sdl.GetCurrentDisplayMode(displayAt)
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
		&Face{A: 1, B: 2, C: 3},
		&Face{A: 1, B: 3, C: 4},
		//
		&Face{A: 4, B: 3, C: 6},
		&Face{A: 6, B: 3, C: 5},
		//
		&Face{A: 8, B: 6, C: 5},
		&Face{A: 8, B: 5, C: 7},
		//
		&Face{A: 1, B: 8, C: 7},
		&Face{A: 1, B: 7, C: 2},
		//
		&Face{A: 6, B: 1, C: 4},
		&Face{A: 6, B: 8, C: 1},
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

	CubeRotation.X += 0.01
	CubeRotation.Y += 0.01
	CubeRotation.Z += 0.01

	ProjectedPoints := []*Vec2{}
	for _, item := range mesh.Vertices {
		transformPoint := item.RotateX(CubeRotation.X)
		transformPoint = transformPoint.RotateY(CubeRotation.Y)
		transformPoint = transformPoint.RotateZ(CubeRotation.Z)

		transformPoint.Z += 5

		projectdPoint := &Vec2{
			X: float64(FovFactor) * transformPoint.X / transformPoint.Z,
			Y: float64(FovFactor) * transformPoint.Y / transformPoint.Z,
		}

		projectdPoint.X += float64(a.w_width) / 2
		projectdPoint.Y += float64(a.w_height) / 2
		ProjectedPoints = append(ProjectedPoints, projectdPoint)
	}

	Triangles = []*Triangle{}
	for _, item := range mesh.Faces {
		Triangles = append(Triangles, &Triangle{
			A: ProjectedPoints[item.A-1],
			B: ProjectedPoints[item.B-1],
			C: ProjectedPoints[item.C-1],
		})
	}

}

func (a *App) Render() {
	a.DrawGrid()

	for _, item := range Triangles {
		a.DrawTriangle(item.A.X, item.A.Y, item.B.X, item.B.Y, item.C.X, item.C.Y, 0xFF00FFFF)
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
