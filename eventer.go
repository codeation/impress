package impress

import (
	"github.com/codeation/impress/low"
)

const (
	GeneralEventType  = 10
	KeyboardEventType = 20
)

type Eventer interface {
	Type() int
}

type GeneralEvent struct {
	Event int
}

func (e GeneralEvent) Type() int {
	return GeneralEventType
}

var (
	UnknownEvent = GeneralEvent{Event: 0}
	DestroyEvent = GeneralEvent{Event: 1}
)

type KeyboardEvent struct {
	Rune    rune
	Name    string
	Shift   bool
	Control bool
	Alt     bool
	Meta    bool
}

var (
	KeyLeft  = KeyboardEvent{Name: "Left"}
	KeyRight = KeyboardEvent{Name: "Right"}
	KeyUp    = KeyboardEvent{Name: "Up"}
	KeyDown  = KeyboardEvent{Name: "Down"}
	KeyEnter = KeyboardEvent{Rune: 13, Name: "Return"}
	KeySave  = KeyboardEvent{Rune: 115, Name: "s", Meta: true}
)

func (e KeyboardEvent) Type() int {
	return KeyboardEventType
}

func NewEventer(e interface{}) Eventer {
	switch e {
	case low.DestroyEevent:
		return DestroyEvent
	}
	switch event := e.(type) {
	case low.KeyboardEvent:
		{
			return KeyboardEvent{
				Rune:    low.KeyRune(event),
				Name:    low.KeyName(event),
				Shift:   low.KeyModifier(event, low.ShiftKeyModifier),
				Control: low.KeyModifier(event, low.ControlKeyModifier),
				Alt:     low.KeyModifier(event, low.AltKeyModifier),
				Meta:    low.KeyModifier(event, low.MetaKeyModifier),
			}
		}
	default:
		return UnknownEvent
	}
}
