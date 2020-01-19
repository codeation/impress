package main

import (
	"log"

	"github.com/codeation/impress"
	_ "github.com/codeation/impress/duo"
)

func main() {
	a := impress.NewApplication()
	a.Title("Hello World Application")
	a.Size(impress.NewRect(0, 0, 640, 480))

	font, err := impress.NewFont(`{"family":"Verdana"}`, 15)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()

	w := a.NewWindow(impress.NewRect(0, 0, 640, 480), impress.NewColor(255, 255, 255))
	w.Text("Hello, world!", font, impress.NewPoint(280, 210), impress.NewColor(0, 0, 0))
	w.Line(impress.NewPoint(270, 230), impress.NewPoint(380, 230), impress.NewColor(255, 0, 0))
	w.Show()

	for a.Event() != impress.DoneEvent {
	}
	a.Quit()
}
