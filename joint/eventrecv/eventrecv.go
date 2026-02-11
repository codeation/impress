// Package implements an internal mechanism to communicate with an impress terminal.
package eventrecv

import (
	"log"
	"sync/atomic"

	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type eventRecv struct {
	callbacks iface.CallbackSet
	onExit    atomic.Bool
	eventPipe *rpc.Pipe
}

func New(callbacks iface.CallbackSet, eventPipe *rpc.Pipe) *eventRecv {
	c := &eventRecv{
		callbacks: callbacks,
		eventPipe: eventPipe,
	}
	go c.listen()
	return c
}

func (c *eventRecv) Done() {
	c.onExit.Store(true)
}

func (c *eventRecv) listen() {
	for {
		var command byte
		if err := c.eventPipe.Get(&command); err != nil {
			if c.onExit.Load() {
				return
			}
			log.Fatalf("readEvents (client): %v", err)
		}

		switch command {
		case iface.EventGeneralCode:
			var eventID uint32
			must(c.eventPipe.Get(&eventID))
			c.callbacks.EventGeneral(int(eventID))

		case iface.EventKeyboardCode:
			var r uint32
			var shift, control, alt, meta bool
			var name string
			must(c.eventPipe.Get(&r, &shift, &control, &alt, &meta, &name))
			c.callbacks.EventKeyboard(rune(r), shift, control, alt, meta, name)

		case iface.EventConfigureCode:
			var width, height, innerWidth, innerHeight int
			must(c.eventPipe.Get(&width, &height, &innerWidth, &innerHeight))
			c.callbacks.EventConfigure(width, height, innerWidth, innerHeight)

		case iface.EventButtonCode:
			var action, button byte
			var x, y int
			must(c.eventPipe.Get(&action, &button, &x, &y))
			c.callbacks.EventButton(int(action), int(button), x, y)

		case iface.EventMotionCode:
			var x, y int
			var shift, control, alt, meta bool
			must(c.eventPipe.Get(&x, &y, &shift, &control, &alt, &meta))
			c.callbacks.EventMotion(x, y, shift, control, alt, meta)

		case iface.EventMenuCode:
			var action string
			must(c.eventPipe.Get(&action))
			c.callbacks.EventMenu(action)

		case iface.EventScrollCode:
			var direction int
			var deltaX, deltaY int
			must(c.eventPipe.Get(&direction, &deltaX, &deltaY))
			c.callbacks.EventScroll(direction, deltaX, deltaY)

		case iface.EventClipboard:
			var typeID int
			var data []byte
			must(c.eventPipe.Get(&typeID, &data))
			c.callbacks.EventClipboard(typeID, data)

		default:
			log.Fatalf("unknown event (client): %d", command)
		}
	}
}

func must(err error) {
	if err != nil {
		log.Fatalf("eventrecv.must: %v", err)
	}
}
