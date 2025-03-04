package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// Font represents a font selection.
type Font struct {
	fonter     driver.Fonter
	Height     int
	LineHeight int
	Baseline   int
	Ascent     int
	Descent    int
	Attributes map[string]string
}

// NewFont returns a Font struct representing a font selection.
// Note that "family" and other attributes are driver-specific.
// Open duo/font.go for details.
func (app *Application) NewFont(height int, attributes map[string]string) *Font {
	fonter := app.driver.NewFont(height, attributes)
	return &Font{
		fonter:     fonter,
		Height:     height,
		LineHeight: fonter.LineHeight(),
		Baseline:   fonter.Baseline(),
		Ascent:     fonter.Ascent(),
		Descent:    fonter.Descent(),
		Attributes: attributes,
	}
}

// Close destroys the font selection.
// Note that a closed font can no longer be used.
func (f *Font) Close() {
	f.fonter.Close()
	f.fonter = nil // TODO: Add notice when the font is closed.
}

// Split breaks the text into lines that fit within the specified width.
// The indent parameter specifies the width to indent the first line.
func (f *Font) Split(text string, edge int, indent int) []string {
	return f.fonter.Split(text, edge, indent)
}

// Size returns the width and height of the drawing area for the given text.
func (f *Font) Size(text string) image.Point {
	return f.fonter.Size(text)
}
