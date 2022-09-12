package src

import "unsafe"

func (a *App) DrawGrid() {
	for y := int32(0); y < a.w_height; y += 10 {
		for x := int32(0); x < a.w_width; x += 10 {
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
