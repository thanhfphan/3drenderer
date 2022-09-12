package main

import (
	"github.com/thanhfphan/3drenderer/src"
)

func main() {
	app := &src.App{}
	if err := app.Setup(); err != nil {
		panic(err)
	}

	for app.IsRunning() {
		app.ProcessInput()
		app.Update()
		app.Render()
	}

	app.Destroy()
}
