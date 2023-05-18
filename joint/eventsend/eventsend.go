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
	cb.eventPipe.Call(iface.EventGeneralCode, uint32(eventID))
	cb.eventPipe.Flush()
}

func (cb *eventSend) EventKeyboard(r rune, shift, control, alt, meta bool, name string) {
	cb.eventPipe.Call(iface.EventKeyboardCode, uint32(r), shift, control, alt, meta, name)
	cb.eventPipe.Flush()
}

func (cb *eventSend) EventConfigure(width, height, innerWidth, innerHeight int) {
	cb.eventPipe.Call(iface.EventConfigureCode, width, height, innerWidth, innerHeight)
	cb.eventPipe.Flush()
}

func (cb *eventSend) EventButton(action, button int, x, y int) {
	cb.eventPipe.Call(iface.EventButtonCode, byte(action), byte(button), x, y)
	cb.eventPipe.Flush()
}

func (cb *eventSend) EventMotion(x, y int, shift, control, alt, meta bool) {
	cb.eventPipe.Call(iface.EventMotionCode, x, y, shift, control, alt, meta)
	cb.eventPipe.Flush()
}

func (cb *eventSend) EventMenu(action string) {
	cb.eventPipe.Call(iface.EventMenuCode, action)
	cb.eventPipe.Flush()
}

func (cb *eventSend) EventScroll(direction int, deltaX, deltaY int) {
	cb.eventPipe.Call(iface.EventScrollCode, direction, deltaX, deltaY)
	cb.eventPipe.Flush()
}
