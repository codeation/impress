package main

import (
	"image/png"
	"log"
	"os"

	"github.com/codeation/impress"
	_ "github.com/codeation/impress/duo"
)

func main() {
	app := impress.NewApplication(impress.NewRect(0, 0, 640, 480), "Image Application")
	defer app.Close()

	f, err := os.Open("test_image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	i, err := impress.NewImage(img)
	if err != nil {
		log.Fatal(err)
	}
	defer i.Close()

	w := app.NewWindow(impress.NewRect(0, 0, 640, 480), impress.NewColor(255, 255, 255))
	w.Image(impress.NewPoint(100, 10), i)
	w.Show()

	app.Wait()
}
