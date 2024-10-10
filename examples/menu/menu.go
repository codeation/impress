package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/duo"
)

var (
	background = color.RGBA{255, 255, 255, 255}
	foreground = color.RGBA{0, 0, 0, 255}
	underline  = color.RGBA{255, 0, 0, 255}
)

func main() {
	app := impress.NewApplication(image.Rect(0, 0, 640, 480), "Menu Demo Application")
	defer app.Close()

	font := app.NewFont(15, map[string]string{"family": "Verdana"})
	defer font.Close()

	menu1 := app.NewMenu("File")
	menu1.NewItem("Open", event.NewMenu("open"))
	menu1.NewItem("Save", event.NewMenu("save"))
	menu1.NewItem("Exit", event.NewMenu("exit"))

	menu2 := app.NewMenu("Edit")
	menu2.NewItem("Paste", event.NewMenu("paste"))

	menu3 := app.NewMenu("Help")
	menu3.NewItem("Impress Help", event.NewMenu("help"))
	menuM := menu3.NewMenu("Open man Page")
	menuM.NewItem("Index", event.NewMenu("index"))
	menuM.NewItem("Search", event.NewMenu("search"))

	w := app.NewWindow(image.Rect(0, 0, 640, 480), background)
	defer w.Drop()

	w.Text("Hello, world!", font, image.Pt(280, 210), foreground)
	w.Line(image.Pt(270, 230), image.Pt(380, 230), underline)
	w.Show()
	app.Sync()

	for {
		e := <-app.Chan()
		if e == event.DestroyEvent || e == event.KeyExit {
			break
		}
		if e.Type() == event.MenuType {
			if ev, ok := e.(event.Menu); ok {
				if ev.Action == "app.exit" {
					break
				}
				fmt.Println("Menu", ev.Action)
			}
		}
	}
}
