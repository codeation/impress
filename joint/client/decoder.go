package client

import (
	"log"

	"github.com/codeation/impress/joint/iface"
)

func (c *Client) listen() {
	for {
		var command byte
		if err := c.eventPipe.Byte(&command).CallErr(); err != nil {
			if c.onExit {
				return
			}
			log.Printf("readEvents (client): %v", err)
			return
		}

		switch command {
		case iface.EventGeneralCode:
			var eventID uint32
			c.eventPipe.
				UInt32(&eventID).
				Call()
			c.callbacks.EventGeneral(int(eventID))

		case iface.EventKeyboardCode:
			var r uint32
			var shift, control, alt, meta bool
			var name string
			c.eventPipe.
				UInt32(&r).
				Bool(&shift).Bool(&control).Bool(&alt).Bool(&meta).
				String(&name).
				Call()
			c.callbacks.EventKeyboard(rune(r), shift, control, alt, meta, name)

		case iface.EventConfigureCode:
			var width, height, innerWidth, innerHeight int
			c.eventPipe.
				Int(&width).Int(&height).
				Int(&innerWidth).Int(&innerHeight).
				Call()
			c.callbacks.EventConfigure(width, height, innerWidth, innerHeight)

		case iface.EventButtonCode:
			var action, button byte
			var x, y int
			c.eventPipe.
				Byte(&action).Byte(&button).
				Int(&x).Int(&y).
				Call()
			c.callbacks.EventButton(int(action), int(button), x, y)

		case iface.EventMotionCode:
			var x, y int
			var shift, control, alt, meta bool
			c.eventPipe.
				Int(&x).Int(&y).
				Bool(&shift).Bool(&control).Bool(&alt).Bool(&meta).
				Call()
			c.callbacks.EventMotion(x, y, shift, control, alt, meta)

		case iface.EventMenuCode:
			var action string
			c.eventPipe.
				String(&action).
				Call()
			c.callbacks.EventMenu(action)

		case iface.EventScrollCode:
			var direction int
			var deltaX, deltaY int
			c.eventPipe.
				Int(&direction).
				Int(&deltaX).Int(&deltaY).
				Call()
			c.callbacks.EventScroll(direction, deltaX, deltaY)

		default:
			log.Printf("unknown event (client): %d", command)
			return
		}
	}
}
