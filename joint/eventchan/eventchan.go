// Package implements an internal mechanism to communicate with an impress terminal.
package eventchan

import (
	"image"

	"github.com/codeation/impress/event"
)

type eventChan struct {
	events chan event.Eventer
}

func New() *eventChan {
	return &eventChan{
		events: make(chan event.Eventer, 32),
	}
}

func (e *eventChan) Chan() <-chan event.Eventer {
	return e.events
}

func (e *eventChan) EventGeneral(eventID int) {
	e.events <- event.General{
		Event: eventID,
	}
}

func (e *eventChan) EventKeyboard(r rune, shift, control, alt, meta bool, name string) {
	e.events <- event.Keyboard{
		Rune:    r,
		Shift:   shift,
		Control: control,
		Alt:     alt,
		Meta:    meta,
		Name:    name,
	}
}

func (e *eventChan) EventConfigure(width, height, innerWidth, innerHeight int) {
	e.events <- event.Configure{
		Size:      image.Pt(width, height),
		InnerSize: image.Pt(innerWidth, innerHeight),
	}
}

func (e *eventChan) EventButton(action, button int, x, y int) {
	e.events <- event.Button{
		Action: action,
		Button: button,
		Point:  image.Pt(x, y),
	}
}

func (e *eventChan) EventMotion(x, y int, shift, control, alt, meta bool) {
	e.events <- event.Motion{
		Point:   image.Pt(x, y),
		Shift:   shift,
		Control: control,
		Alt:     alt,
		Meta:    meta,
	}
}

func (e *eventChan) EventMenu(action string) {
	e.events <- event.Menu{
		Action: action,
	}
}

func (e *eventChan) EventScroll(direction int, deltaX, deltaY int) {
	e.events <- event.Scroll{
		Direction: direction,
		DeltaX:    deltaX,
		DeltaY:    deltaY,
	}
}
