package impress

// Driver is the interface to a application level functions
type Driver interface {
	Init()
	Done()
	Title(title string)
	Size(rect Rect)
	NewWindow(rect Rect, color Color) Painter
	NewFont(font *Font) (Fonter, error)
	Event() Eventer
}

// Painter is the interface to a window functions
type Painter interface {
	Drop()
	Size(rect Rect)
	Clear()
	Show()
	Fill(rect Rect, color Color)
	Line(from Point, to Point, color Color)
	Text(text string, font *Font, from Point, color Color)
}

// Fonter is the interface to a font functions
type Fonter interface {
	Close()
	Split(text string, edge int) []string
	Size(text string) Size
}

var driver Driver

// Register makes a GUI driver available
func Register(d Driver) {
	driver = d
}
