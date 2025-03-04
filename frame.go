package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// Frame represents a GUI frame.
type Frame struct {
	framer driver.Framer
}

// NewFrame creates a new inner frame with the specified size.
func (app *Application) NewFrame(rect image.Rectangle) *Frame {
	return app.frame.NewFrame(rect)
}

// NewFrame creates a new child frame with the specified size.
func (f *Frame) NewFrame(rect image.Rectangle) *Frame {
	return &Frame{
		framer: f.framer.NewFrame(rect),
	}
}

// Drop deletes the frame.
// Note that a dropped frame can no longer be used.
func (f *Frame) Drop() {
	f.framer.Drop()
	f.framer = nil // TODO: Add notice when the frame is dropped.
}

// Size changes the frame size and position.
func (f *Frame) Size(rect image.Rectangle) {
	f.framer.Size(rect)
}

// Raise brings the frame to the forefront.
func (f *Frame) Raise() {
	f.framer.Raise()
}
