package domain

import (
	"image"
	"image/color"
)

func rectangle(rect image.Rectangle) (int, int, int, int) {
	x0, y0, x1, y1 := rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return x0, y0, x1 - x0, y1 - y0
}

func colors(c color.Color) (uint16, uint16, uint16, uint16) {
	r, g, b, a := c.RGBA()
	return uint16(r), uint16(g), uint16(b), uint16(a)
}
