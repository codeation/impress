package main

import (
	"log"

	"github.com/codeation/impress"
	_ "github.com/codeation/impress/duo"
)

func main() {
	app := impress.NewApplication(impress.NewRect(0, 0, 640, 480), "Hello World Application")
	defer app.Close()

	font, err := impress.NewFont(`{"family":"Verdana"}`, 15)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()

	w := app.NewWindow(impress.NewRect(0, 0, 640, 480), impress.NewColor(255, 255, 255))
	w.Text("Hello, world!", font, impress.NewPoint(280, 210), impress.NewColor(0, 0, 0))
	w.Line(impress.NewPoint(270, 230), impress.NewPoint(380, 230), impress.NewColor(255, 0, 0))
	w.Show()

	for {
		event := <-app.Chan()
		if event == impress.DestroyEvent || event == impress.KeyExit {
			break
		}
	}
}
