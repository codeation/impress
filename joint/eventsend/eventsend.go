// Package implements an internal mechanism to communicate with an impress terminal.
package eventsend

import (
	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type eventSend struct {
	eventPipe *rpc.Pipe
}

func New(eventPipe *rpc.Pipe) *eventSend {
	return &eventSend{
		eventPipe: eventPipe,
	}
}

func (cb *eventSend) EventGeneral(eventID int) {
	cb.eventPipe.
		Lock().
		Put(iface.EventGeneralCode, uint32(eventID)).
		Flush().
		Unlock()
}

func (cb *eventSend) EventKeyboard(r rune, shift, control, alt, meta bool, name string) {
	cb.eventPipe.
		Lock().
		Put(iface.EventKeyboardCode, uint32(r), shift, control, alt, meta, name).
		Flush().
		Unlock()
}

func (cb *eventSend) EventConfigure(width, height, innerWidth, innerHeight int) {
	cb.eventPipe.
		Lock().
		Put(iface.EventConfigureCode, width, height, innerWidth, innerHeight).
		Flush().
		Unlock()
}

func (cb *eventSend) EventButton(action, button int, x, y int) {
	cb.eventPipe.
		Lock().
		Put(iface.EventButtonCode, byte(action), byte(button), x, y).
		Flush().
		Unlock()
}

func (cb *eventSend) EventMotion(x, y int, shift, control, alt, meta bool) {
	cb.eventPipe.
		Lock().
		Put(iface.EventMotionCode, x, y, shift, control, alt, meta).
		Flush().
		Unlock()
}

func (cb *eventSend) EventMenu(action string) {
	cb.eventPipe.
		Lock().
		Put(iface.EventMenuCode, action).
		Flush().
		Unlock()
}

func (cb *eventSend) EventScroll(direction int, deltaX, deltaY int) {
	cb.eventPipe.
		Lock().
		Put(iface.EventScrollCode, direction, deltaX, deltaY).
		Flush().
		Unlock()
}

func (cb *eventSend) EventClipboard(typeID int, data []byte) {
	cb.eventPipe.
		Lock().
		Put(iface.EventClipboard, typeID, data).
		Flush().
		Unlock()
}
