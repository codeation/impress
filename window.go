package impress

import (
	"image"
	"image/color"

	"github.com/codeation/impress/driver"
)

// Window represents inner window
type Window struct {
	painter driver.Painter
}

// NewWindow creates new inner window with a specified size and background color
func (app *Application) NewWindow(rect image.Rectangle, background color.Color) *Window {
	return app.frame.NewWindow(rect, background)
}

// NewWindow creates new frame window with a specified size and background color
func (f *Frame) NewWindow(rect image.Rectangle, background color.Color) *Window {
	return &Window{
		painter: f.framer.NewWindow(rect, background),
	}
}

// Drop deletes window.
// Note that a dropped window can no longer be used
func (w *Window) Drop() {
	w.painter.Drop()
	w.painter = nil // TODO notice when the window is dropped
}

// Size changes window size and position
func (w *Window) Size(rect image.Rectangle) {
	w.painter.Size(rect)
}

// Clear clears current window
func (w *Window) Clear() {
	w.painter.Clear()
}

// Fill draws a rectangle with specified size and foreground color
func (w *Window) Fill(rect image.Rectangle, foreground color.Color) {
	w.painter.Fill(rect, foreground)
}

// Line draws a color line connecting two specified points
func (w *Window) Line(from image.Point, to image.Point, foreground color.Color) {
	w.painter.Line(from, to, foreground)
}

// Image draws a image into specified rectangle.
// Specify a different rectangle size to scale the image
func (w *Window) Image(rect image.Rectangle, img *Image) {
	w.painter.Image(rect, img.Imager)
}

// Text draws a text at specified location using a specified font and foreground color
func (w *Window) Text(text string, font *Font, from image.Point, foreground color.Color) {
	w.painter.Text(text, font.Fonter, from, foreground)
}

// Show sends the contents of the window to the screen.
// Note that a drawings are not visible until Show
func (w *Window) Show() {
	w.painter.Show()
}

// Raise brings the window to the forefront
func (w *Window) Raise() {
	w.painter.Raise()
}
