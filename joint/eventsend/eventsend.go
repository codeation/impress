// Package implements an internal mechanism to communicate with an impress terminal.
package eventsend

import (
	"log"

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
	must(cb.eventPipe.PutTx(iface.EventGeneralCode, uint32(eventID)))
	must(cb.eventPipe.Sync())
}

func (cb *eventSend) EventKeyboard(r rune, shift, control, alt, meta bool, name string) {
	must(cb.eventPipe.PutTx(iface.EventKeyboardCode, uint32(r), shift, control, alt, meta, name))
	must(cb.eventPipe.Sync())
}

func (cb *eventSend) EventConfigure(width, height, innerWidth, innerHeight int) {
	must(cb.eventPipe.PutTx(iface.EventConfigureCode, width, height, innerWidth, innerHeight))
	must(cb.eventPipe.Sync())
}

func (cb *eventSend) EventButton(action, button int, x, y int) {
	must(cb.eventPipe.PutTx(iface.EventButtonCode, byte(action), byte(button), x, y))
	must(cb.eventPipe.Sync())
}

func (cb *eventSend) EventMotion(x, y int, shift, control, alt, meta bool) {
	must(cb.eventPipe.PutTx(iface.EventMotionCode, x, y, shift, control, alt, meta))
	must(cb.eventPipe.Sync())
}

func (cb *eventSend) EventMenu(action string) {
	must(cb.eventPipe.PutTx(iface.EventMenuCode, action))
	must(cb.eventPipe.Sync())
}

func (cb *eventSend) EventScroll(direction int, deltaX, deltaY int) {
	must(cb.eventPipe.PutTx(iface.EventScrollCode, direction, deltaX, deltaY))
	must(cb.eventPipe.Sync())
}

func (cb *eventSend) EventClipboard(typeID int, data []byte) {
	must(cb.eventPipe.PutTx(iface.EventClipboard, typeID, data))
	must(cb.eventPipe.Sync())
}

func must(err error) {
	if err != nil {
		log.Fatalf("eventsend.must: %v", err)
	}
}
