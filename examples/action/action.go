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

	appRect   = impress.NewRect(0, 0, 640, 480)
	leftRect  = impress.NewRect(0, 0, 340, 480)
	rightRect = impress.NewRect(300, 0, 340, 480)
)

type smallWindow struct {
	w    *impress.Window
	rect impress.Rect
	pos  int
	font *impress.Font
}

func NewSmallWindow(app *impress.Application, rect impress.Rect, color impress.Color, font *impress.Font) *smallWindow {
	return &smallWindow{
		w:    app.NewWindow(rect, color),
		rect: rect,
		pos:  0,
		font: font,
	}
}

func (sw *smallWindow) Drop() {
	sw.w.Drop()
}

func (sw *smallWindow) Redraw(isActive bool) {
	sw.w.Clear()
	sw.w.Text("Hello, world!", sw.font, impress.NewPoint(105, sw.pos), black)
	if isActive {
		sw.w.Line(impress.NewPoint(105, sw.pos+30), impress.NewPoint(215, sw.pos+30), red)
		sw.w.Raise()
	}
	sw.w.Show()
}

func (sw *smallWindow) Event(event impress.Eventer) {
	switch event {
	case impress.KeyUp:
		sw.pos -= 16
		if sw.pos < 0 {
			sw.pos = 0
		}
	case impress.KeyDown:
		sw.pos += 16
		if sw.pos > 450 {
			sw.pos = 450
		}
	}
}

func action(app *impress.Application, windows []*smallWindow) {
	activeWindow := windows[0]
	for {
		for _, w := range windows {
			w.Redraw(w == activeWindow)
		}

		event := <-app.Chan()
		if event == impress.DestroyEvent || event == impress.KeyExit {
			return
		}
		switch {
		case event.Type() == impress.ButtonEventType:
			clickEvent, ok := event.(impress.ButtonEvent)
			if ok && clickEvent.Action == impress.ButtonActionPress && clickEvent.Button == impress.ButtonLeft {
				for _, w := range windows {
					if clickEvent.Point.In(w.rect) && w != activeWindow {
						activeWindow = w
						break
					}
				}
			}
		case event == impress.KeyLeft:
			activeWindow = windows[0]
		case event == impress.KeyRight:
			activeWindow = windows[1]
		default:
			activeWindow.Event(event)
		}
	}
}

func main() {
	app := impress.NewApplication(appRect, "Example")
	defer app.Close()

	font, err := impress.NewFont(`{"family":"Verdana"}`, 15)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()

	w1 := NewSmallWindow(app, leftRect, white, font)
	defer w1.Drop()
	w2 := NewSmallWindow(app, rightRect, silver, font)
	defer w2.Drop()

	action(app, []*smallWindow{w1, w2})
}
