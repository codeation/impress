package impress

import (
	"image"
	"image/color"

	"github.com/codeation/impress/driver"
)

// Window represents an inner window.
type Window struct {
	painter driver.Painter
}

// NewWindow creates a new inner window with a specified size and background color.
func (app *Application) NewWindow(rect image.Rectangle, background color.Color) *Window {
	return app.frame.NewWindow(rect, background)
}

// NewWindow creates a new frame window with a specified size and background color.
func (f *Frame) NewWindow(rect image.Rectangle, background color.Color) *Window {
	return &Window{
		painter: f.framer.NewWindow(rect, background),
	}
}

// Drop deletes the window.
// Note that a dropped window can no longer be used.
func (w *Window) Drop() {
	w.painter.Drop()
	w.painter = nil // TODO: Add notice when the window is dropped.
}

// Size changes the window size and position.
func (w *Window) Size(rect image.Rectangle) {
	w.painter.Size(rect)
}

// Clear clears the current window.
func (w *Window) Clear() {
	w.painter.Clear()
}

// Fill draws a rectangle with the specified size and foreground color.
func (w *Window) Fill(rect image.Rectangle, foreground color.Color) {
	w.painter.Fill(rect, foreground)
}

// Line draws a colored line connecting two specified points.
func (w *Window) Line(from image.Point, to image.Point, foreground color.Color) {
	w.painter.Line(from, to, foreground)
}

// Image draws an image into the specified rectangle.
// Specify a different rectangle size to scale the image.
func (w *Window) Image(rect image.Rectangle, img *Image) {
	w.painter.Image(rect, img.imager)
}

// Text draws text at the specified location using the specified font and foreground color.
func (w *Window) Text(text string, font *Font, from image.Point, foreground color.Color) {
	w.painter.Text(text, font.fonter, from, foreground)
}

// Show sends the contents of the window to the screen.
// Note that drawings are not visible until Show is called.
func (w *Window) Show() {
	w.painter.Show()
}

// Raise brings the window to the forefront.
func (w *Window) Raise() {
	w.painter.Raise()
}
