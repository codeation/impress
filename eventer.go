package impress

import (
	"unicode"
)

// Event types
const (
	GeneralEventType  = 10
	KeyboardEventType = 20
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
	UnknownEvent = GeneralEvent{Event: 0}
	DestroyEvent = GeneralEvent{Event: 1}
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
