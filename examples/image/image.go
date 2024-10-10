package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/duo"
)

func main() {
	rect := image.Rect(0, 0, 640, 480)
	app := impress.NewApplication(rect, "Image Application")
	defer app.Close()

	f, err := os.Open("test_image.png")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	i, err := png.Decode(f)
	if err != nil {
		log.Println(err)
		return
	}
	img := app.NewImage(i)
	defer img.Close()

	w := app.NewWindow(image.Rectangle{Max: rect.Size()}, color.RGBA{255, 255, 255, 255})
	defer w.Drop()

	size := rect.Size()
	for {
		if len(app.Chan()) == 0 {
			w.Size(image.Rectangle{Max: size})
			w.Clear()
			offset := size.Sub(img.Size).Div(2)
			w.Image(image.Rectangle{Min: offset, Max: offset.Add(img.Size)}, img)
			w.Show()
			app.Sync()
		}

		e := <-app.Chan()
		if e == event.DestroyEvent || e == event.KeyExit {
			break
		}
		if ev, ok := e.(event.Configure); ok {
			size = ev.InnerSize
		}
	}
}
