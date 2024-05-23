package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

type Frame struct {
	framer driver.Framer
}

// NewFrame create new inner frame with a specified size
func (app *Application) NewFrame(rect image.Rectangle) *Frame {
	return app.frame.NewFrame(rect)
}

// NewFrame create new child frame with a specified size
func (f *Frame) NewFrame(rect image.Rectangle) *Frame {
	return &Frame{
		framer: f.framer.NewFrame(rect),
	}
}

// Drop deletes frame.
// Note that a dropped frame can no longer be used
func (f *Frame) Drop() {
	f.framer.Drop()
	f.framer = nil // TODO notice when the window is dropped
}

// Size changes frame size and position
func (f *Frame) Size(rect image.Rectangle) {
	f.framer.Size(rect)
}

// Raise brings the frame to the forefront
func (f *Frame) Raise() {
	f.framer.Raise()
}
