package domain

import (
	"image"
	"image/color"
	"log"

	"github.com/codeation/impress/driver"
)

type window struct {
	app        *application
	id         int
	rect       image.Rectangle
	background color.Color
}

func (f *frame) NewWindow(rect image.Rectangle, background color.Color) driver.Painter {
	id := f.app.nextWindowID()
	x, y, width, height := rectangle(rect)
	r, g, b, a := colors(background)
	f.app.caller.WindowNew(id, f.id, x, y, width, height)
	f.app.caller.WindowFill(id, 0, 0, width, height, r, g, b, a)
	return &window{
		app:        f.app,
		id:         id,
		rect:       rect,
		background: background,
	}
}

func (w *window) Drop() {
	w.app.caller.WindowDrop(w.id)
}

func (w *window) Size(rect image.Rectangle) {
	w.rect = rect
	x, y, width, height := rectangle(w.rect)
	w.app.caller.WindowSize(w.id, x, y, width, height)
}

func (w *window) Raise() {
	w.app.caller.WindowRaise(w.id)
}

func (w *window) Clear() {
	_, _, width, height := rectangle(w.rect)
	r, g, b, a := colors(w.background)
	w.app.caller.WindowClear(w.id)
	w.app.caller.WindowFill(w.id, 0, 0, width, height, r, g, b, a)
}

func (w *window) Show() {
	w.app.caller.WindowShow(w.id)
}

func (w *window) Fill(rect image.Rectangle, foreground color.Color) {
	x, y, width, height := rectangle(rect)
	r, g, b, a := colors(foreground)
	w.app.caller.WindowFill(w.id, x, y, width, height, r, g, b, a)
}

func (w *window) Line(from image.Point, to image.Point, foreground color.Color) {
	r, g, b, a := colors(foreground)
	w.app.caller.WindowLine(w.id, from.X, from.Y, to.X, to.Y, r, g, b, a)
}

func (w *window) Text(text string, fonter driver.Fonter, from image.Point, foreground color.Color) {
	f, ok := fonter.(interface{ ID() int })
	if !ok {
		log.Println("wrong fonter type")
		return
	}
	r, g, b, a := colors(foreground)
	w.app.caller.WindowText(w.id, from.X, from.Y, r, g, b, a, f.ID(), text)
}

func (w *window) Image(rect image.Rectangle, img driver.Imager) {
	p, ok := img.(interface{ ID() int })
	if !ok {
		log.Println("wrong imager type")
		return
	}
	x, y, width, height := rectangle(rect)
	w.app.caller.WindowImage(w.id, x, y, width, height, p.ID())
}
