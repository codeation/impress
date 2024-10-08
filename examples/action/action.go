package main

import (
	"image"
	"image/color"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/duo"
)

var (
	white  = color.RGBA{255, 255, 255, 255}
	black  = color.RGBA{0, 0, 0, 255}
	silver = color.RGBA{224, 224, 224, 255}
	red    = color.RGBA{255, 0, 0, 255}

	appRect   = image.Rect(0, 0, 640, 480)
	leftRect  = image.Rect(0, 0, 340, 480)
	rightRect = leftRect.Add(image.Pt(300, 0))
)

type smallWindow struct {
	w    *impress.Window
	rect image.Rectangle
	pos  int
	font *impress.Font
}

func NewSmallWindow(app *impress.Application, rect image.Rectangle,
	background color.Color, font *impress.Font,
) *smallWindow {
	return &smallWindow{
		w:    app.NewWindow(rect, background),
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
	sw.w.Text("Hello, world!", sw.font, image.Pt(105, sw.pos), black)
	if isActive {
		sw.w.Line(image.Pt(105, sw.pos+24), image.Pt(215, sw.pos+24), red)
	}
	sw.w.Show()
}

func (sw *smallWindow) Event(e event.Eventer) {
	switch e {
	case event.KeyUp:
		sw.pos -= 16
		if sw.pos < 0 {
			sw.pos = 0
		}
	case event.KeyDown:
		sw.pos += 16
		if sw.pos > 450 {
			sw.pos = 450
		}
	}
}

func run(app *impress.Application, windows []*smallWindow) {
	activeWindow := windows[0]
	activeWindow.w.Raise()
	for {
		for _, w := range windows {
			w.Redraw(w == activeWindow)
		}
		app.Sync()

		e := <-app.Chan()
		switch {
		case e == event.DestroyEvent || e == event.KeyExit:
			return
		case e == event.KeyLeft:
			activeWindow = windows[0]
			activeWindow.w.Raise()
		case e == event.KeyRight:
			activeWindow = windows[1]
			activeWindow.w.Raise()
		case e.Type() == event.ButtonType:
			clickEvent, ok := e.(event.Button)
			if ok && clickEvent.Action == event.ButtonActionPress && clickEvent.Button == event.ButtonLeft {
				for _, w := range windows {
					if clickEvent.Point.In(w.rect) && w != activeWindow {
						activeWindow = w
						activeWindow.w.Raise()
						break
					}
				}
			}
		default:
			activeWindow.Event(e)
		}
	}
}

func main() {
	app := impress.NewApplication(appRect, "Example")
	defer app.Close()

	font := app.NewFont(15, map[string]string{"family": "Verdana"})
	defer font.Close()

	w1 := NewSmallWindow(app, leftRect, white, font)
	defer w1.Drop()
	w2 := NewSmallWindow(app, rightRect, silver, font)
	defer w2.Drop()

	run(app, []*smallWindow{w1, w2})
}
