package main

import (
	"image"
	"image/color"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/duo"
)

func main() {
	app := impress.NewApplication(image.Rect(0, 0, 480, 240), "Hello World Application")
	defer app.Close()

	font := impress.NewFont(15, map[string]string{"family": "Verdana"})
	defer font.Close()

	w := app.NewWindow(image.Rect(0, 0, 480, 240), color.RGBA{255, 255, 255, 255})
	defer w.Drop()

	w.Text("Hello, world!", font, image.Pt(200, 100), color.RGBA{0, 0, 0, 255})
	w.Line(image.Pt(200, 120), image.Pt(300, 120), color.RGBA{255, 0, 0, 255})
	w.Show()
	app.Sync()

	for {
		action := <-app.Chan()
		if action == event.DestroyEvent || action == event.KeyExit {
			break
		}
	}
}
