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
		sw.w.Line(image.Pt(105, sw.pos+30), image.Pt(215, sw.pos+30), red)
		sw.w.Raise()
	}
	sw.w.Show()
}

func (sw *smallWindow) Event(action event.Eventer) {
	switch action {
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

func action(app *impress.Application, windows []*smallWindow) {
	activeWindow := windows[0]
	for {
		for _, w := range windows {
			w.Redraw(w == activeWindow)
		}
		app.Sync()

		action := <-app.Chan()
		switch {
		case action == event.DestroyEvent || action == event.KeyExit:
			return
		case action == event.KeyLeft:
			activeWindow = windows[0]
		case action == event.KeyRight:
			activeWindow = windows[1]
		case action.Type() == event.ButtonType:
			clickEvent, ok := action.(event.Button)
			if ok && clickEvent.Action == event.ButtonActionPress && clickEvent.Button == event.ButtonLeft {
				for _, w := range windows {
					if clickEvent.Point.In(w.rect) && w != activeWindow {
						activeWindow = w
						break
					}
				}
			}
		default:
			activeWindow.Event(action)
		}
	}
}

func main() {
	app := impress.NewApplication(appRect, "Example")
	defer app.Close()

	font := impress.NewFont(15, map[string]string{"family": "Verdana"})
	defer font.Close()

	w1 := NewSmallWindow(app, leftRect, white, font)
	defer w1.Drop()
	w2 := NewSmallWindow(app, rightRect, silver, font)
	defer w2.Drop()

	action(app, []*smallWindow{w1, w2})
}
