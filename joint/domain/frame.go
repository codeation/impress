package domain

import (
	"image"

	"github.com/codeation/impress/driver"
)

type frame struct {
	app *application
	id  int
}

func (app *application) NewFrame(rect image.Rectangle) driver.Framer {
	id := app.nextFrameID()
	x, y, width, height := rectangle(rect)
	app.caller.FrameNew(id, 0, x, y, width, height)
	return &frame{
		app: app,
		id:  id,
	}
}

func (f *frame) Drop() {
	f.app.caller.FrameDrop(f.id)
}

func (f *frame) Size(rect image.Rectangle) {
	x, y, width, height := rectangle(rect)
	f.app.caller.FrameSize(f.id, x, y, width, height)
}

func (f *frame) Raise() {
	f.app.caller.FrameDrop(f.id)
}

func (f *frame) NewFrame(rect image.Rectangle) driver.Framer {
	id := f.app.nextFrameID()
	x, y, width, height := rectangle(rect)
	f.app.caller.FrameNew(id, f.id, x, y, width, height)
	return &frame{
		app: f.app,
		id:  id,
	}
}
