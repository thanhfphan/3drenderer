package src

import (
	"fmt"
	"sort"

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
	isRunning bool
	isDebug   bool

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
	if numMonitor > 1 {
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
	a.isDebug = false

	vertices, faces, err := a.LoadOBJFile("./assets/jug/jug.obj")
	if err != nil {
		return err
	}

	mesh.Vertices = vertices
	mesh.Faces = faces

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
			} else if t.Keysym.Sym == sdl.K_d {
				if t.State == sdl.RELEASED {
					fmt.Println("Pressed d")
					a.isDebug = !a.isDebug
				}
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

	// CubeRotation.X += 0.01
	// CubeRotation.Y += 0.01
	CubeRotation.Z += 0.01

	Triangles = []*Triangle{}

	for _, item := range mesh.Faces {
		// ***** Transform Vertices *****
		ta := mesh.Vertices[item.A-1].Rotate(CubeRotation)
		ta = ta.RotateX(3.14 / 2)
		ta.Z += 2
		tb := mesh.Vertices[item.B-1].Rotate(CubeRotation)
		tb = tb.RotateX(3.14 / 2)
		tb.Z += 2
		tc := mesh.Vertices[item.C-1].Rotate(CubeRotation)
		tc = tc.RotateX(3.14 / 2)
		tc.Z += 2

		if a.isDebug {
			ac := Vec3Sub(tc, ta)
			ab := Vec3Sub(tb, ta)
			n := Vec3Cross(ac, ab)
			n.Normalize()

			vc := Vec3Sub(CameraPosition, ta)
			if Vec3Dot(vc, n) < 0 {
				continue
			}
		}

		//  ***** Projection *******
		// A
		projectA := &Vec3{
			X: float64(FovFactor) * ta.X / ta.Z,
			Y: float64(FovFactor) * ta.Y / ta.Z,
			Z: ta.Z,
		}
		projectA.X += float64(a.w_width) / 2
		projectA.Y += float64(a.w_height)/2 + 200
		// B
		projectB := &Vec3{
			X: float64(FovFactor) * tb.X / tb.Z,
			Y: float64(FovFactor) * tb.Y / tb.Z,
			Z: tb.Z,
		}
		projectB.X += float64(a.w_width) / 2
		projectB.Y += float64(a.w_height)/2 + 200
		// C
		projectC := &Vec3{
			X: float64(FovFactor) * tc.X / tc.Z,
			Y: float64(FovFactor) * tc.Y / tc.Z,
			Z: tc.Z,
		}
		projectC.X += float64(a.w_width) / 2
		projectC.Y += float64(a.w_height)/2 + 200

		Triangles = append(Triangles, &Triangle{
			A:        projectA,
			B:        projectB,
			C:        projectC,
			AvgDepth: (projectA.Z + projectB.Z + projectC.Z) / 3,
		})
	}
}

func (a *App) Render() {
	a.DrawGrid()

	sort.Slice(Triangles, func(i, j int) bool {
		return Triangles[i].AvgDepth > Triangles[j].AvgDepth
	})

	for _, item := range Triangles {
		a.DrawTriangle(item.A.X, item.A.Y, item.B.X, item.B.Y, item.C.X, item.C.Y, 0xFF00FFFF)
		a.FillTriangle(*item.A, *item.B, *item.C, 0xFF808080)
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
