package duo

import (
	"image"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func rectangle(rect image.Rectangle) (x, y, width, height int) {
	x = min(rect.Min.X, rect.Max.X)
	y = min(rect.Min.Y, rect.Max.Y)
	width = abs(rect.Dx())
	height = abs(rect.Dy())
	return
}
