package src

import (
	"math"
	"unsafe"
)

func (a *App) DrawTriangle(x0, y0, x1, y1, x2, y2 float64, color uint32) {
	a.DrawLine(x0, y0, x1, y1, color)
	a.DrawLine(x1, y1, x2, y2, color)
	a.DrawLine(x0, y0, x2, y2, color)
}

func (a *App) DrawLine(x0, y0, x1, y1 float64, color uint32) {
	dx := x1 - x0
	dy := y1 - y0

	steps := math.Abs(dy)
	if math.Abs(dx) > math.Abs(dy) {
		steps = math.Abs(dx)
	}

	Xinc := dx / steps
	Yinc := dy / steps

	X := x0
	Y := y0
	for i := 0; i <= int(steps); i++ {
		a.DrawPixel(int32(math.Round(X)), int32(math.Round(Y)), color)
		X += Xinc // increment in x at each step
		Y += Yinc // increment in y at each step
	}
}

func (a *App) DrawGrid() {
	for y := int32(0); y < a.w_height; y += 5 {
		for x := int32(0); x < a.w_width; x += 5 {
			a.colorBuffer[(a.w_width*y)+x] = 0xFF444444
		}

	}
}

func (a *App) DrawPixel(x, y int32, color uint32) {
	if x >= 0 && x < a.w_width && y >= 0 && y < a.w_height {
		a.colorBuffer[(a.w_width*y)+x] = color
	}
}

func (a *App) DrawRect(x, y, width, height int32, color uint32) {
	for i := int32(0); i < width; i++ {
		for j := int32(0); j < height; j++ {
			a.DrawPixel(x+i, y+j, color)
		}
	}
}

func (a *App) ClearColorBuffer(color uint32) {
	for y := int32(0); y < a.w_height; y++ {
		for x := int32(0); x < a.w_width; x++ {
			a.colorBuffer[(a.w_width*y)+x] = color
		}
	}
}

func (a *App) RenderColorBuffer() {
	size := unsafe.Sizeof(uint32(0))
	a.colorBufferTexture.Update(nil, unsafe.Pointer(&a.colorBuffer[0]), int(a.w_width*int32(size)))
	a.renderer.Copy(a.colorBufferTexture, nil, nil)
}
