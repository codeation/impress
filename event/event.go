package event

import (
	"image"
	"unicode"
)

// Event types
const (
	_ int = iota
	GeneralType
	KeyboardType
	ButtonType
	MotionType
	MenuType
	ConfigureType
	ScrollType
)

// Eventer is the interface that groups GUI events
type Eventer interface {
	Type() int
}

// General is a general purpose notification
type General struct {
	Event int
}

// Type returns event type
func (e General) Type() int {
	return GeneralType
}

// Signal events
var (
	UnknownEvent = General{Event: 0}
	DestroyEvent = General{Event: 1}
)

// Configure event
type Configure struct {
	Size      image.Point
	InnerSize image.Point
}

// Type returns event type
func (e Configure) Type() int {
	return ConfigureType
}

// Keyboard is a keyboard event
type Keyboard struct {
	Rune    rune
	Shift   bool
	Control bool
	Alt     bool
	Meta    bool
	Name    string
}

// Keyboard events
var (
	KeyLeft      = Keyboard{Name: "Left"}
	KeyRight     = Keyboard{Name: "Right"}
	KeyUp        = Keyboard{Name: "Up"}
	KeyDown      = Keyboard{Name: "Down"}
	KeyBackSpace = Keyboard{Rune: 8, Name: "BackSpace"}
	KeyTab       = Keyboard{Rune: 9, Name: "Tab"}
	KeyEnter     = Keyboard{Rune: 13, Name: "Return"}
	KeyEscape    = Keyboard{Rune: 27, Name: "Escape"}
	KeySpace     = Keyboard{Rune: 32, Name: "space"}
	KeyDelete    = Keyboard{Rune: 127, Name: "Delete"}
	KeyPageUp    = Keyboard{Name: "Page_Up"}
	KeyPageDown  = Keyboard{Name: "Page_Down"}
	KeyHome      = Keyboard{Name: "Home"}
	KeyEnd       = Keyboard{Name: "End"}
)

// Type returns event type
func (e Keyboard) Type() int {
	return KeyboardType
}

// IsGraphic tests printable rune
func (e Keyboard) IsGraphic() bool {
	return !e.Control && !e.Meta && unicode.IsGraphic(e.Rune)
}

// Button is mouse button event
type Button struct {
	Action int
	Button int
	Point  image.Point
}

// Button action type
const (
	ButtonActionPress   = 4
	ButtonActionDouble  = 5
	ButtonActionTriple  = 6
	ButtonActionRelease = 7
)

// Button number
const (
	ButtonLeft   = 1
	ButtonMiddle = 2
	ButtonRight  = 3
)

// Type returns event type
func (e Button) Type() int {
	return ButtonType
}

// Motion is mouse motion event
type Motion struct {
	Point   image.Point
	Shift   bool
	Control bool
	Alt     bool
	Meta    bool
}

// Type returns event type
func (e Motion) Type() int {
	return MotionType
}

// Menu is menu action event
type Menu struct {
	Action string
}

// NewMenu returns a menu action event
func NewMenu(short string) Menu {
	return Menu{
		Action: "app." + short,
	}
}

// Type returns event type
func (e Menu) Type() int {
	return MenuType
}

// Scrool is scrool event
type Scroll struct {
	Direction int
	DeltaX    int
	DeltaY    int
}

// Scroll direction type
const (
	ScrollUp     = 0
	ScrollDown   = 1
	ScrollLeft   = 2
	ScrollRight  = 3
	ScrollSmooth = 4
)

// Type returns event type
func (e Scroll) Type() int {
	return ScrollType
}
