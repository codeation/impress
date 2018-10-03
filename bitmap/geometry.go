// Package bitmap implement 2-D image buffer as part of impress library.
// See impress library for more details.
package bitmap

func min(i1, i2 int) int {
	if i1 < i2 {
		return i1
	}
	return i2
}

func max(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y int
}

// NewPoint is same as Point{x, y}.
func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

// A Size is an Width and Height pair.
type Size struct {
	Width, Height int
}

// NewSize is same as Size{width, height}.
func NewSize(width, height int) Size {
	return Size{
		Width:  width,
		Height: height,
	}
}

// A Rect contains the upper left corner coordinates and rectangle size.
type Rect struct {
	Point
	Size
}

// NewRect is same as Rect{Point:Point{x,y}, Size:Size{width, height}}.
func NewRect(x, y, width, height int) Rect {
	return Rect{
		Point: NewPoint(x, y),
		Size:  NewSize(width, height),
	}
}

// Color represents a 24-bit color, having 8 bits for each of red, green, blue.
type Color struct {
	R, G, B int
}

// NewColor is same as Color{r, g, b}
func NewColor(r, g, b int) Color {
	return Color{
		R: r,
		G: g,
		B: b,
	}
}
