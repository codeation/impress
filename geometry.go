package impress

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

func (p Point) Size(size Size) Rect {
	return Rect{
		Point: p,
		Size:  size,
	}
}

func (p Point) Move(add Point) Point {
	return Point{
		X: p.X + add.X,
		Y: p.Y + add.Y,
	}
}

func (p Point) MoveX(x int) Point {
	return Point{
		X: p.X + x,
		Y: p.Y,
	}
}

func (p Point) MoveY(y int) Point {
	return Point{
		X: p.X,
		Y: p.Y + y,
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

func (r Rect) Move(add Point) Rect {
	return Rect{
		Point: r.Point.Move(add),
		Size:  r.Size,
	}
}

// In returns true when point is inside rect
func (p *Point) In(rect Rect) bool {
	return p.X >= rect.X && p.X < rect.X+rect.Width && p.Y >= rect.Y && p.Y < rect.Y+rect.Height
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
