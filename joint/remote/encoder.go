package remote

import (
	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type callback struct {
	eventPipe *rpc.Pipe
}

func NewCallBacks(eventPipe *rpc.Pipe) *callback {
	return &callback{
		eventPipe: eventPipe,
	}
}

func (cb *callback) EventGeneral(eventID int) {
	cb.eventPipe.Call(iface.EventGeneralCode, uint32(eventID))
	cb.eventPipe.Flush()
}

func (cb *callback) EventKeyboard(r rune, shift, control, alt, meta bool, name string) {
	cb.eventPipe.Call(iface.EventKeyboardCode, uint32(r), shift, control, alt, meta, name)
	cb.eventPipe.Flush()
}

func (cb *callback) EventConfigure(width, height, innerWidth, innerHeight int) {
	cb.eventPipe.Call(iface.EventConfigureCode, width, height, innerWidth, innerHeight)
	cb.eventPipe.Flush()
}

func (cb *callback) EventButton(action, button int, x, y int) {
	cb.eventPipe.Call(iface.EventButtonCode, byte(action), byte(button), x, y)
	cb.eventPipe.Flush()
}

func (cb *callback) EventMotion(x, y int, shift, control, alt, meta bool) {
	cb.eventPipe.Call(iface.EventMotionCode, x, y, shift, control, alt, meta)
	cb.eventPipe.Flush()
}

func (cb *callback) EventMenu(action string) {
	cb.eventPipe.Call(iface.EventMenuCode, action)
	cb.eventPipe.Flush()
}

func (cb *callback) EventScroll(direction int, deltaX, deltaY int) {
	cb.eventPipe.Call(iface.EventScrollCode, direction, deltaX, deltaY)
	cb.eventPipe.Flush()
}
