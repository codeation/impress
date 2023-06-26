package main

import (
	"image"
	"image/color"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/duo"
)

var (
	red  = color.RGBA{255, 0, 0, 255}
	blue = color.RGBA{0, 0, 255, 255}
)

func main() {
	app := impress.NewApplication(image.Rect(0, 0, 740, 480), "Frames demo")
	defer app.Close()

	leftFrame := app.NewFrame(image.Rect(0, 0, 320, 480))
	defer leftFrame.Drop()
	leftDialog := leftFrame.NewWindow(image.Rect(100, 100, 300, 200), red)
	defer leftDialog.Drop()

	rightFrame := app.NewFrame(image.Rect(320, 0, 640, 480))
	defer rightFrame.Drop()
	rightDialog := rightFrame.NewWindow(image.Rect(100, 100, 300, 200), blue)
	defer rightDialog.Drop()

	leftDialog.Show()
	rightDialog.Show()
	app.Sync()

	for {
		action := <-app.Chan()
		if action == event.DestroyEvent || action == event.KeyExit {
			break
		}
	}
}
