package impress

import (
	"unicode"
)

// Event types
const (
	GeneralEventType  = 10
	KeyboardEventType = 20
	ButtonEventType   = 30
	MotionEventType   = 40
	MenuEventType     = 50
)

// Eventer is the interface that groups GUI events
type Eventer interface {
	Type() int
}

// GeneralEvent is a general purpose notification
type GeneralEvent struct {
	Event int
}

// Type returns event type
func (e GeneralEvent) Type() int {
	return GeneralEventType
}

// Signal events
var (
	UnknownEvent     = GeneralEvent{Event: 0}
	DestroyEvent     = GeneralEvent{Event: 1}
	DoneEvent        = GeneralEvent{Event: 2}
	ActivatedEvent   = GeneralEvent{Event: 3}
	DeactivatedEvent = GeneralEvent{Event: 4}
)

// KeyboardEvent is a keyboard event
type KeyboardEvent struct {
	Rune    rune
	Shift   bool
	Control bool
	Alt     bool
	Meta    bool
	Name    string
}

// Keyboard events
var (
	KeyLeft      = KeyboardEvent{Name: "Left"}
	KeyRight     = KeyboardEvent{Name: "Right"}
	KeyUp        = KeyboardEvent{Name: "Up"}
	KeyDown      = KeyboardEvent{Name: "Down"}
	KeyEnter     = KeyboardEvent{Rune: 13, Name: "Return"}
	KeyBackSpace = KeyboardEvent{Rune: 8, Name: "BackSpace"}
)

// Type returns event type
func (e KeyboardEvent) Type() int {
	return KeyboardEventType
}

// IsGraphic tests printable rune
func (e KeyboardEvent) IsGraphic() bool {
	return !e.Control && !e.Meta && unicode.IsGraphic(e.Rune)
}

// ButtonEvent is mouse button event
type ButtonEvent struct {
	Action int
	Button int
	Point  Point
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
func (e ButtonEvent) Type() int {
	return ButtonEventType
}

// MotionEvent is mouse motion event
type MotionEvent struct {
	Point   Point
	Shift   bool
	Control bool
	Alt     bool
	Meta    bool
}

// Type returns event type
func (e MotionEvent) Type() int {
	return MotionEventType
}

// MenuEvent is menu action event
type MenuEvent struct {
	Action string
}

// NewMenuEvent returns a menu action event
func NewMenuEvent(short string) MenuEvent {
	return MenuEvent{
		Action: "app." + short,
	}
}

// Type returns event type
func (e MenuEvent) Type() int {
	return MenuEventType
}
