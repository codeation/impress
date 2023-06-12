package main

import (
	"image"
	"image/color"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/duo"
)

var (
	background = color.RGBA{255, 255, 255, 0}
	red        = color.RGBA{255, 0, 0, 0}
	blue       = color.RGBA{0, 0, 255, 0}
)

func main() {
	app := impress.NewApplication(image.Rect(0, 0, 640, 480), "Try to resize window")
	defer app.Close()

	windowRect := image.Rect(0, 0, 640, 480)
	w := app.NewWindow(windowRect, background)
	defer w.Drop()
	readyRect := image.Rectangle{}

	for {
		if windowRect != readyRect {
			readyRect = windowRect

			w.Size(windowRect)
			w.Clear()

			w.Line(image.Pt(windowRect.Min.X, windowRect.Min.Y), image.Pt(windowRect.Max.X, windowRect.Min.Y), blue)     // up
			w.Line(image.Pt(windowRect.Min.X, windowRect.Min.Y), image.Pt(windowRect.Min.X, windowRect.Max.Y), blue)     // left
			w.Line(image.Pt(windowRect.Max.X-1, windowRect.Min.Y), image.Pt(windowRect.Max.X-1, windowRect.Max.Y), blue) // right
			w.Line(image.Pt(windowRect.Min.X, windowRect.Max.Y-1), image.Pt(windowRect.Max.X, windowRect.Max.Y-1), blue) // down

			w.Line(image.Pt(windowRect.Min.X, windowRect.Min.Y), image.Pt(windowRect.Max.X, windowRect.Max.Y), blue)
			w.Line(image.Pt(windowRect.Max.X, windowRect.Min.Y), image.Pt(windowRect.Min.X, windowRect.Max.Y), blue)

			w.Fill(image.Rect(100, 100, 200, 200), red)

			w.Line(image.Pt(100, 100-2), image.Pt(200, 100-2), red) // up
			w.Line(image.Pt(100-2, 100), image.Pt(100-2, 200), red) // left
			w.Line(image.Pt(199+2, 100), image.Pt(199+2, 200), red) // right
			w.Line(image.Pt(100, 199+2), image.Pt(200, 199+2), red) // down

			w.Show()
			app.Sync()
		}

		action := <-app.Chan()
		if action == event.DestroyEvent || action == event.KeyExit {
			break
		}
		if resizeEvent, ok := action.(event.Configure); ok {
			windowRect = image.Rect(0, 0, resizeEvent.InnerSize.X, resizeEvent.InnerSize.Y)
		}
	}
}
