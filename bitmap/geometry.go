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

type Point struct {
	X, Y int
}

func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

type Size struct {
	Width, Height int
}

func NewSize(width, height int) Size {
	return Size{
		Width:  width,
		Height: height,
	}
}

type Rect struct {
	Point
	Size
}

func NewRect(x, y, width, height int) Rect {
	return Rect{
		Point: NewPoint(x, y),
		Size:  NewSize(width, height),
	}
}

type Color struct {
	R, G, B int
}

func NewColor(r, g, b int) Color {
	return Color{
		R: r,
		G: g,
		B: b,
	}
}
