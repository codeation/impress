package impress

// Driver is a internal interface to a application level functions
type Driver interface {
	Init()
	Done()
	Title(title string)
	Size(rect Rect)
	NewWindow(rect Rect, color Color) Painter
	NewFont(font *Font) (Fonter, error)
	NewImage(img *Image) (Imager, error)
	NewMenu(label string) Menuer
	Chan() <-chan Eventer
}

// Painter is a internal interface to a window functions
type Painter interface {
	Drop()
	Size(rect Rect)
	Raise()
	Clear()
	Show()
	Fill(rect Rect, color Color)
	Line(from Point, to Point, color Color)
	Image(from Point, img *Image)
	Text(text string, font *Font, from Point, color Color)
}

// Fonter is a internal interface to a font functions
type Fonter interface {
	Close()
	Split(text string, edge int) []string
	Size(text string) Size
}

// Imager is a internal interface to a image functions
type Imager interface {
	Close()
}

// Menuer is a internal interface to a menu node functions
type Menuer interface {
	NewMenu(label string) Menuer
	NewItem(label string, action string)
}

var driver Driver

// Register makes a GUI driver available. Internal
func Register(d Driver) {
	driver = d
}
