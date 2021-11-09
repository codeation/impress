package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// Font represents a font selection
type Font struct {
	fonter     driver.Fonter
	Height     int
	Baseline   int
	Ascent     int
	Descent    int
	Attributes map[string]string
}

// NewFont return a font selection struct.
// Note than "family" and other attributes are driver specific.
// Open duo/font.go for details.
func NewFont(height int, attributes map[string]string) *Font {
	fonter := d.NewFont(height, attributes)
	return &Font{
		fonter:     fonter,
		Height:     height,
		Baseline:   fonter.Baseline(),
		Ascent:     fonter.Ascent(),
		Descent:    fonter.Descent(),
		Attributes: attributes,
	}
}

// Close destroys font selection
func (f *Font) Close() {
	f.fonter.Close()
	f.fonter = nil // TODO notice when the font is closed
}

// Split breaks the text into lines that fit in the specified width
func (f *Font) Split(text string, edge int) []string {
	return f.fonter.Split(text, edge)
}

// Size returns the width and height of the drawing area
func (f *Font) Size(text string) image.Point {
	return f.fonter.Size(text)
}
