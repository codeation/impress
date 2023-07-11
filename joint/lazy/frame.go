package lazy

import (
	"image"

	"github.com/codeation/impress/driver"
)

type frame struct {
	driver.Framer
	rect image.Rectangle
}

func (a *app) NewFrame(rect image.Rectangle) driver.Framer {
	return &frame{
		Framer: a.Driver.NewFrame(rect),
		rect:   rect,
	}
}

func (f *frame) Size(rect image.Rectangle) {
	if f.rect == rect {
		return
	}
	f.rect = rect
	f.Framer.Size(rect)
}
