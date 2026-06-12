// Package event defines various types of GUI events
package event

import (
	"image"
	"unicode"

	"github.com/codeation/impress/clipboard"
)

// Eventer is the interface that groups GUI events.
type Eventer any

// General represents a general-purpose notification event.
type General struct {
	Event int // Event identifies the specific general event
}

// Predefined general events.
var (
	UnknownEvent = General{Event: 0} // UnknownEvent represents an unknown event
	DestroyEvent = General{Event: 1} // DestroyEvent represents a destroy event
)

// Configure represents a window configuration event.
type Configure struct {
	Size      image.Point // Size represents the size of the window
	InnerSize image.Point // InnerSize represents the size of the inner part of the window
}

// Keyboard represents a keyboard event.
type Keyboard struct {
	Rune    rune   // Rune represents the character input from the keyboard
	Shift   bool   // Shift indicates if the Shift key is pressed
	Control bool   // Control indicates if the Control key is pressed
	Alt     bool   // Alt indicates if the Alt key is pressed
	Meta    bool   // Meta indicates if the Meta key is pressed
	Name    string // Name represents the name of the key
}

// Predefined keyboard events.
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

// IsGraphic tests if the rune is a printable character.
func (e Keyboard) IsGraphic() bool {
	return !e.Control && !e.Meta && unicode.IsGraphic(e.Rune)
}

// Button represents a mouse button event.
type Button struct {
	Action int         // Action represents the button action type
	Button int         // Button represents the mouse button number
	Point  image.Point // Point represents the location of the mouse pointer
}

// Button action types.
const (
	ButtonActionPress   = 4 // ButtonActionPress represents a button press action
	ButtonActionDouble  = 5 // ButtonActionDouble represents a double-click action
	ButtonActionTriple  = 6 // ButtonActionTriple represents a triple-click action
	ButtonActionRelease = 7 // ButtonActionRelease represents a button release action
)

// Button numbers.
const (
	ButtonLeft   = 1 // ButtonLeft represents the left mouse button
	ButtonMiddle = 2 // ButtonMiddle represents the middle mouse button
	ButtonRight  = 3 // ButtonRight represents the right mouse button
)

// Motion represents a mouse motion event.
type Motion struct {
	Point   image.Point // Point represents the location of the mouse pointer
	Shift   bool        // Shift indicates if the Shift key is pressed during motion
	Control bool        // Control indicates if the Control key is pressed during motion
	Alt     bool        // Alt indicates if the Alt key is pressed during motion
	Meta    bool        // Meta indicates if the Meta key is pressed during motion
}

// Menu represents a menu action event.
type Menu struct {
	Action string // Action represents the menu action command
}

// NewMenu returns a new menu action event.
func NewMenu(short string) Menu {
	return Menu{
		Action: "app." + short, // Prefix the action with "app."
	}
}

// Scroll represents a scroll event.
type Scroll struct {
	Direction int // Direction represents the direction of the scroll
	DeltaX    int // DeltaX represents the horizontal scroll delta
	DeltaY    int // DeltaY represents the vertical scroll delta
}

// Scroll direction types.
const (
	ScrollUp     = 0 // ScrollUp represents scrolling up
	ScrollDown   = 1 // ScrollDown represents scrolling down
	ScrollLeft   = 2 // ScrollLeft represents scrolling left
	ScrollRight  = 3 // ScrollRight represents scrolling right
	ScrollSmooth = 4 // ScrollSmooth represents smooth scrolling
)

// Clipboard represents an event with clipboard data.
type Clipboard struct {
	Data clipboard.Clipboarder // Data represents the clipboard content
}
