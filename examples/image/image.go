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

var (
	background = color.RGBA{255, 255, 255, 0}
)

func main() {
	app := impress.NewApplication(image.Rect(0, 0, 640, 480), "Image Application")
	defer app.Close()

	f, err := os.Open("test_image.png")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Println(err)
		return
	}
	i := impress.NewImage(img)
	defer i.Close()

	w := app.NewWindow(image.Rect(0, 0, 640, 480), background)
	w.Image(image.Rect(100, 10, 100+i.Size.X, 10+i.Size.Y), i)
	w.Show()
	app.Sync()

	for {
		action := <-app.Chan()
		if action == event.DestroyEvent || action == event.KeyExit {
			break
		}
	}
}
