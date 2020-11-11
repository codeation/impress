package main

import (
	"fmt"
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

	menu1 := app.NewMenu("File")
	menu1.NewItem("Open", impress.NewMenuEvent("open"))
	menu1.NewItem("Save", impress.NewMenuEvent("save"))
	menu1.NewItem("Exit", impress.NewMenuEvent("exit"))

	menu2 := app.NewMenu("Edit")
	menu2.NewItemFunc("Paste", impress.NewMenuEvent("check"), func() { fmt.Println("Paste event") })

	menu3 := app.NewMenu("Help")
	menu3.NewItem("Impress Help", impress.NewMenuEvent("help"))
	menuM := menu3.NewMenu("Open man Page")
	menuM.NewItem("Index", impress.NewMenuEvent("index"))
	menuM.NewItem("Search", impress.NewMenuEvent("search"))

	app.Wait()
}
