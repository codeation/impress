// Package implements an internal mechanism to communicate with an impress terminal.
package eventrecv

import (
	"log"

	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type eventRecv struct {
	callbacks iface.CallbackSet
	onExit    bool
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
	c.onExit = true
}

func (c *eventRecv) listen() {
	for {
		var command byte
		if err := c.eventPipe.Get(&command).Err(); err != nil {
			if c.onExit {
				return
			}
			log.Printf("readEvents (client): %v", err)
			return
		}

		switch command {
		case iface.EventGeneralCode:
			var eventID uint32
			c.eventPipe.Get(&eventID)
			c.callbacks.EventGeneral(int(eventID))

		case iface.EventKeyboardCode:
			var r uint32
			var shift, control, alt, meta bool
			var name string
			c.eventPipe.Get(&r, &shift, &control, &alt, &meta, &name)
			c.callbacks.EventKeyboard(rune(r), shift, control, alt, meta, name)

		case iface.EventConfigureCode:
			var width, height, innerWidth, innerHeight int
			c.eventPipe.Get(&width, &height, &innerWidth, &innerHeight)
			c.callbacks.EventConfigure(width, height, innerWidth, innerHeight)

		case iface.EventButtonCode:
			var action, button byte
			var x, y int
			c.eventPipe.Get(&action, &button, &x, &y)
			c.callbacks.EventButton(int(action), int(button), x, y)

		case iface.EventMotionCode:
			var x, y int
			var shift, control, alt, meta bool
			c.eventPipe.Get(&x, &y, &shift, &control, &alt, &meta)
			c.callbacks.EventMotion(x, y, shift, control, alt, meta)

		case iface.EventMenuCode:
			var action string
			c.eventPipe.Get(&action)
			c.callbacks.EventMenu(action)

		case iface.EventScrollCode:
			var direction int
			var deltaX, deltaY int
			c.eventPipe.Get(&direction, &deltaX, &deltaY)
			c.callbacks.EventScroll(direction, deltaX, deltaY)

		default:
			log.Printf("unknown event (client): %d", command)
			return
		}
	}
}
