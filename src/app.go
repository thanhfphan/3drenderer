package src

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	Triangles = []*Triangle{}

	mesh           = &Mesh{Vertices: []*Vec3{}, Faces: []*Face{}}
	CameraPosition = &Vec3{X: 0, Y: 0, Z: 0}
	CubeRotation   = &Vec3{X: 0, Y: 0, Z: 0}

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

	Triangles = []*Triangle{}
	for _, item := range mesh.Faces {
		// ***** Transform Vertices *****
		ta := mesh.Vertices[item.A-1].Rotate(CubeRotation)
		ta.Z += 5
		tb := mesh.Vertices[item.B-1].Rotate(CubeRotation)
		tb.Z += 5
		tc := mesh.Vertices[item.C-1].Rotate(CubeRotation)
		tc.Z += 5

		//  ***** Projection *******
		// A
		projectA := &Vec3{
			X: float64(FovFactor) * ta.X / ta.Z,
			Y: float64(FovFactor) * ta.Y / ta.Z,
		}
		projectA.X += float64(a.w_width) / 2
		projectA.Y += float64(a.w_height) / 2
		// B
		projectB := &Vec3{
			X: float64(FovFactor) * tb.X / tb.Z,
			Y: float64(FovFactor) * tb.Y / tb.Z,
		}
		projectB.X += float64(a.w_width) / 2
		projectB.Y += float64(a.w_height) / 2
		// C
		projectC := &Vec3{
			X: float64(FovFactor) * tc.X / tc.Z,
			Y: float64(FovFactor) * tc.Y / tc.Z,
		}
		projectC.X += float64(a.w_width) / 2
		projectC.Y += float64(a.w_height) / 2

		Triangles = append(Triangles, &Triangle{
			A: projectA,
			B: projectB,
			C: projectC,
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
