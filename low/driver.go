package low

import (
	"github.com/codeation/impress"
)

func init() {
	impress.Register(&binddriver)
}

var binddriver driver

type driver struct {
	app *Application
}

func (d *driver) Init() {
	d.app = NewApplication()
}

func (d *driver) Main() {
	d.app.Main()
}

func (d *driver) Done() {
	d.app.Quit()
}

func (d *driver) Title(title string) {
	d.app.Title(title)
}

func (d *driver) Size(rect impress.Rect) {
	d.app.Size(rect.X, rect.Y, rect.Width, rect.Height)
}

func (d *driver) NewWindow(rect impress.Rect, background impress.Color) impress.Painter {
	p := &paint{
		window:     d.app.NewWindow(),
		rect:       rect,
		background: background,
		snapshot:   NewSnapshot(rect.Size, background),
	}
	p.window.Move(d.app, rect.X, rect.Y)
	return p
}

func (d *driver) NewFont(f *impress.Font) (impress.Fonter, error) {
	return OpenFont(f.Attr["filename"], f.Height)
}

func (d *driver) Event() impress.Eventer {
	return EventDequeue()
}

type paint struct {
	window     *Window
	rect       impress.Rect
	background impress.Color
	snapshot   *Snapshot
}

func (p *paint) Clear() {
	p.snapshot.Fill(p.rect, p.background)
}

func (p *paint) Fill(rect impress.Rect, color impress.Color) {
	p.snapshot.Fill(rect, color)
}

func (p *paint) Line(from impress.Point, to impress.Point, color impress.Color) {
	p.snapshot.Line(from, to, color)
}

func (p *paint) Text(text string, font *impress.Font, from impress.Point, color impress.Color) {
	p.snapshot.Text(text, font, from, color)
}

func (p *paint) Show() {
	p.window.Set(p.snapshot.Picture())
}
