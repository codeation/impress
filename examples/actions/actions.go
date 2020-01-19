package main

import (
	"log"

	"github.com/codeation/impress"
	_ "github.com/codeation/impress/duo"
)

var (
	black  = impress.NewColor(0, 0, 0)
	white  = impress.NewColor(255, 255, 255)
	silver = impress.NewColor(224, 224, 224)
	red    = impress.NewColor(255, 0, 0)
)

// Action loop for any window
func action(act *impress.Action, w *impress.Window, font *impress.Font) {
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
	app := impress.NewApplication()
	defer app.Close()
	app.Title("Example")
	app.Size(impress.NewRect(0, 0, 640, 480))

	font, err := impress.NewFont(`{"family":"Verdana"}`, 15)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()

	// First window
	act1 := app.NewAction()
	w1 := app.NewWindow(impress.NewRect(0, 0, 320, 480), white)
	defer w1.Drop()
	app.Start(func() {
		action(act1, w1, font)
	})

	// Second window
	act2 := app.NewAction()
	w2 := app.NewWindow(impress.NewRect(320, 0, 320, 480), silver)
	defer w2.Drop()
	app.Start(func() {
		action(act2, w2, font)
	})

	// Toggle windows
	act1.Activate()
	app.OnEvent(impress.KeyLeft, act1.Activate)
	app.OnEvent(impress.KeyRight, act2.Activate)

	// Wait for the actions to complete
	app.Wait()
}
