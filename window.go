package impress

// Window represents inner window
type Window struct {
	paint      Painter
	background Color
}

// NewWindow creates new inner window with a specified size and background color
func (a *Application) NewWindow(rect Rect, color Color) *Window {
	return &Window{
		paint:      driver.NewWindow(rect, color),
		background: NewColor(255, 255, 255),
	}
}

// Clear clears current window
func (w *Window) Clear() {
	w.paint.Clear()
}

// Fill draws a rectangle with specified size and foreground color
func (w *Window) Fill(rect Rect, color Color) {
	w.paint.Fill(rect, color)
}

// Line draws a color line connecting two specified points
func (w *Window) Line(from Point, to Point, color Color) {
	w.paint.Line(from, to, color)
}

// Text draws a text at specified location using a specified font and foreground color
func (w *Window) Text(text string, font *Font, from Point, color Color) {
	w.paint.Text(text, font, from, color)
}

// Show sends the contents of the window to the screen
func (w *Window) Show() {
	w.paint.Show()
}