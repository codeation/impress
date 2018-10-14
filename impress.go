package impress

// Driver is the interface to a application level functions
type Driver interface {
	Init()
	Main()
	Done()
	Title(title string)
	Size(rect Rect)
	NewWindow(rect Rect, color Color) Painter
	NewFont(font *Font) (Fonter, error)
	Event() Eventer
}

// Painter is the interface to a window functions
type Painter interface {
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
}

var driver Driver

// Register makes a GUI driver available
func Register(d Driver) {
	driver = d
}
