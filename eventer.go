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
	Rune rune
	Name string
}

var (
	KeyLeft  = KeyboardEvent{Rune: 0, Name: "Left"}
	KeyRight = KeyboardEvent{Rune: 0, Name: "Right"}
	KeyUp    = KeyboardEvent{Rune: 0, Name: "Up"}
	KeyDown  = KeyboardEvent{Rune: 0, Name: "Down"}
	KeyEnter = KeyboardEvent{Rune: 13, Name: "Return"}
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
				Rune: low.KeyRune(event),
				Name: low.KeyName(event),
			}
		}
	default:
		return UnknownEvent
	}
}
