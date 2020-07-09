package main

import (
	"log"

	"github.com/codeation/impress"
	"github.com/codeation/impress/action"

	_ "github.com/codeation/impress/duo"
)

var (
	black  = impress.NewColor(0, 0, 0)
	white  = impress.NewColor(255, 255, 255)
	silver = impress.NewColor(224, 224, 224)
	red    = impress.NewColor(255, 0, 0)

	leftRect  = impress.NewRect(0, 0, 320, 480)
	rightRect = impress.NewRect(320, 0, 320, 480)
)

// loop is endless func for any window
func loop(act *action.Action, w *impress.Window, font *impress.Font) {
	var pos int
	for {
		// Draw line
		w.Clear()
		w.Text("Hello, world!", font, impress.NewPoint(105, pos), black)
		if act.Activated() {
			w.Line(impress.NewPoint(105, pos+30), impress.NewPoint(215, pos+30), red)
		}
		w.Show()

		// Move line position
		switch act.Event() {
		case impress.KeyUp:
			pos -= 16
			if pos < 0 {
				pos = 0
			}
		case impress.KeyDown:
			pos += 16
			if pos > 450 {
				pos = 450
			}
		case impress.DoneEvent:
			return
		}
	}
}

func main() {
	app := impress.NewApplication(impress.NewRect(0, 0, 640, 480), "Example")
	defer app.Close()

	font, err := impress.NewFont(`{"family":"Verdana"}`, 15)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()

	// Left window and actor
	w1 := app.NewWindow(leftRect, white)
	defer w1.Drop()
	act1 := action.NewAction(app, leftRect, func(act *action.Action) { loop(act, w1, font) })
	app.OnEvent(impress.KeyLeft, act1.Activate)

	// Right window and actor
	w2 := app.NewWindow(rightRect, silver)
	defer w2.Drop()
	act2 := action.NewAction(app, rightRect, func(act *action.Action) { loop(act, w2, font) })
	app.OnEvent(impress.KeyRight, act2.Activate)

	app.Wait()
}
