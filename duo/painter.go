package duo

import (
	"image"
	"image/color"

	"github.com/codeation/impress/driver"
)

type canvas struct {
	driver     *duo
	layout     *layout
	id         int
	rect       image.Rectangle
	background color.Color
}

func (l *layout) NewWindow(rect image.Rectangle, background color.Color) driver.Painter {
	l.driver.lastWindowID++
	c := &canvas{
		driver:     l.driver,
		layout:     l,
		id:         l.driver.lastWindowID,
		rect:       rect,
		background: background,
	}
	x, y, width, height := rectangle(c.rect)
	r, g, b, _ := c.background.RGBA()
	c.driver.streamPipe.Call(
		'D', c.id, l.id, x, y, width, height,
		'F', c.id, 0, 0, width, height, r, g, b)
	return c
}

func (c *canvas) Drop() {
	c.driver.streamPipe.Call(
		'O', c.id)
}

func (c *canvas) Raise() {
	c.driver.streamPipe.Call(
		'A', c.id)
}

func (c *canvas) Size(rect image.Rectangle) {
	c.rect = rect
	x, y, width, height := rectangle(c.rect)
	c.driver.streamPipe.Call(
		'Z', c.id, x, y, width, height)
}

func (c *canvas) Clear() {
	_, _, width, height := rectangle(c.rect)
	r, g, b, _ := c.background.RGBA()
	c.driver.streamPipe.Call(
		'C', c.id,
		'F', c.id, 0, 0, width, height, r, g, b)
}

func (c *canvas) Fill(rect image.Rectangle, foreground color.Color) {
	r, g, b, _ := foreground.RGBA()
	x, y, width, height := rectangle(rect)
	c.driver.streamPipe.Call(
		'F', c.id, x, y, width, height, r, g, b)
}

func (c *canvas) Line(from image.Point, to image.Point, foreground color.Color) {
	r, g, b, _ := foreground.RGBA()
	c.driver.streamPipe.Call(
		'L', c.id, from.X, from.Y, to.X, to.Y, r, g, b)
}

func (c *canvas) Image(rect image.Rectangle, img driver.Imager) {
	b := img.(*bitmap)
	c.driver.streamPipe.Call(
		'I', c.id, rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(), b.id)
}

func (c *canvas) Text(text string, fonter driver.Fonter, from image.Point, foreground color.Color) {
	f := fonter.(*fontface)
	r, g, b, _ := foreground.RGBA()
	c.driver.streamPipe.Call(
		'U', c.id, from.X, from.Y, r, g, b, f.id, f.height, text)
}

func (c *canvas) Show() {
	c.driver.streamPipe.Call(
		'W', c.id)
}
