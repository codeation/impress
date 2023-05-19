package main

import (
	"image"
	"image/color"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/canvas"
)

func main() {
	app := impress.NewApplication(image.Rect(0, 0, 640, 480), "Hello World Application")
	defer app.Close()

	font := impress.NewFont(15, map[string]string{"family": "Verdana"})
	defer font.Close()

	w := app.NewWindow(image.Rect(0, 0, 640, 480), color.RGBA{240, 240, 240, 0})
	defer w.Drop()

	w.Text("Hello, world!", font, image.Pt(280, 210), color.RGBA{0, 0, 0, 0})
	w.Line(image.Pt(270, 230), image.Pt(380, 230), color.RGBA{255, 0, 0, 0})
	w.Show()
	app.Sync()

	for {
		action := <-app.Chan()
		if action == event.DestroyEvent || action == event.KeyExit {
			break
		}
	}
}
