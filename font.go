package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// Font represents a font selection
type Font struct {
	driver.Fonter
	Height     int
	LineHeight int
	Baseline   int
	Ascent     int
	Descent    int
	Attributes map[string]string
}

// NewFont return a font selection struct.
// Note than "family" and other attributes are driver specific.
// Open duo/font.go for details.
func (app *Application) NewFont(height int, attributes map[string]string) *Font {
	fonter := app.driver.NewFont(height, attributes)
	return &Font{
		Fonter:     fonter,
		Height:     height,
		LineHeight: fonter.LineHeight(),
		Baseline:   fonter.Baseline(),
		Ascent:     fonter.Ascent(),
		Descent:    fonter.Descent(),
		Attributes: attributes,
	}
}

// Close destroys font selection
func (f *Font) Close() {
	f.Fonter.Close()
	f.Fonter = nil // TODO notice when the font is closed
}

// Split breaks the text into lines that fit in the specified width;
// indent is a width to indent first line
func (f *Font) Split(text string, edge int, indent int) []string {
	return f.Fonter.Split(text, edge, indent)
}

// Size returns the width and height of the drawing area
func (f *Font) Size(text string) image.Point {
	return f.Fonter.Size(text)
}
