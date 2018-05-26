package impress

import (
	"github.com/codeation/impress/bitmap"
	"github.com/codeation/impress/low"
)

type Application struct {
	application *low.Application
	handlers    map[Eventer]func()
}

func NewApplication() *Application {
	return &Application{
		application: low.NewApplication(),
		handlers:    map[Eventer]func(){},
	}
}

func (a *Application) Size(rect bitmap.Rect) {
	a.application.Size(rect.Point.X, rect.Point.Y, rect.Size.Width, rect.Size.Height)
}

func (a *Application) Main() {
	a.application.Main()
}

func (a *Application) Quit() {
	a.application.Quit()
}

func (a *Application) OnEvent(event Eventer, handler func()) {
	a.handlers[event] = handler
}

func (a *Application) Event() Eventer {
	for {
		e := NewEventer(a.application.Event())
		handler, ok := a.handlers[e]
		if !ok {
			return e
		}
		handler()
	}
	return nil
}

type Window struct {
	window      *low.Window
	application *Application
	Rect        bitmap.Rect
	Background  bitmap.Color
	plate       *bitmap.Snapshot
}

func (a *Application) NewWindow(rect bitmap.Rect, background bitmap.Color) *Window {
	w := &Window{
		window:      a.application.NewWindow(),
		application: a,
		Rect:        rect,
		Background:  background,
		plate:       bitmap.NewSnapshot(rect.Size, background),
	}
	w.window.Move(a.application, rect.Point.X, rect.Point.Y)
	return w
}

func (w *Window) Clear() {
	w.plate.Fill(w.Rect, w.Background)
}

func (w *Window) Fill(rect bitmap.Rect, foreground bitmap.Color) {
	w.plate.Fill(rect, foreground)
}

func (w *Window) Line(from, to bitmap.Point, foreground bitmap.Color) {
	w.plate.Line(from, to, foreground)
}

func (w *Window) Text(text string, font *bitmap.Font, point bitmap.Point,
	foreground bitmap.Color) (bitmap.Rect, error) {
	return w.plate.Text(text, font, point, foreground)
}

func (w *Window) Show() {
	w.window.Set(w.plate.Picture())
}
