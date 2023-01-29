package duo

import (
	"image"

	"github.com/codeation/impress/driver"
)

type layout struct {
	driver *duo
	id     int
	rect   image.Rectangle
}

func (d *duo) NewFrame(rect image.Rectangle) driver.Framer {
	d.lastFrameID++
	l := &layout{
		driver: d,
		id:     d.lastFrameID,
		rect:   rect,
	}
	x, y, width, height := rectangle(l.rect)
	l.driver.streamPipe.Call(
		'Y', l.id, int(0), x, y, width, height)
	return l
}

func (l *layout) NewFrame(rect image.Rectangle) driver.Framer {
	l.driver.lastFrameID++
	child := &layout{
		driver: l.driver,
		id:     l.driver.lastFrameID,
		rect:   rect,
	}
	x, y, width, height := rectangle(child.rect)
	l.driver.streamPipe.Call(
		'Y', child.id, l.id, x, y, width, height)
	return child
}

func (l *layout) Drop() {
	l.driver.streamPipe.Call(
		'Q', l.id)
}

func (l *layout) Raise() {
	l.driver.streamPipe.Call(
		'J', l.id)
}

func (l *layout) Size(rect image.Rectangle) {
	l.rect = rect
	x, y, width, height := rectangle(l.rect)
	l.driver.streamPipe.Call(
		'H', l.id, x, y, width, height)
}
