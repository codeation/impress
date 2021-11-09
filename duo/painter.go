package duo

import (
	"image"
	"image/color"

	"github.com/codeation/impress/driver"
)

type canvas struct {
	driver     *duo
	id         int
	rect       image.Rectangle
	background color.Color
}

func (d *duo) NewWindow(rect image.Rectangle, background color.Color) driver.Painter {
	d.lastWindowID++
	c := &canvas{
		driver:     d,
		id:         d.lastWindowID,
		rect:       rect,
		background: background,
	}
	x, y, width, height := rectangle(c.rect)
	r, g, b, _ := c.background.RGBA()
	c.driver.drawPipe.Call(
		'D', c.id, x, y, width, height, r, g, b,
		'F', c.id, 0, 0, width, height, r, g, b)
	return c
}

func (c *canvas) Drop() {
	c.driver.drawPipe.Call(
		'O', c.id)
}

func (c *canvas) Raise() {
	c.driver.drawPipe.Call(
		'A', c.id)
}

func (c *canvas) Size(rect image.Rectangle) {
	c.rect = rect
	x, y, width, height := rectangle(c.rect)
	r, g, b, _ := c.background.RGBA()
	c.driver.drawPipe.Call(
		'Z', c.id, x, y, width, height, r, g, b)
}

func (c *canvas) Clear() {
	_, _, width, height := rectangle(c.rect)
	r, g, b, _ := c.background.RGBA()
	c.driver.drawPipe.Call(
		'C', c.id,
		'F', c.id, 0, 0, width, height, r, g, b)
}

func (c *canvas) Fill(rect image.Rectangle, foreground color.Color) {
	r, g, b, _ := foreground.RGBA()
	x, y, width, height := rectangle(rect)
	c.driver.drawPipe.Call(
		'F', c.id, x, y, width, height, r, g, b)
}

func (c *canvas) Line(from image.Point, to image.Point, foreground color.Color) {
	r, g, b, _ := foreground.RGBA()
	c.driver.drawPipe.Call(
		'L', c.id, from.X, from.Y, to.X, to.Y, r, g, b)
}

func (p *canvas) Image(from image.Point, img driver.Imager) {
	b := img.(*bitmap)
	p.driver.drawPipe.Call(
		'I', p.id, from.X, from.Y, b.id)
}

func (c *canvas) Text(text string, fonter driver.Fonter, from image.Point, foreground color.Color) {
	f := fonter.(*fontface)
	r, g, b, _ := foreground.RGBA()
	c.driver.drawPipe.Call(
		'U', c.id, from.X, from.Y, r, g, b, f.id, f.height, text)
}

func (c *canvas) Show() {
	var zero int
	c.driver.drawPipe.
		Int16(&zero).Call(
		'W', c.id)
}
