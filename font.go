package impress

import (
	"encoding/json"
)

// Font represents a font selection
type Font struct {
	Fonter   Fonter
	Height   int
	Baseline int
	Ascent   int
	Descent  int
	Attr     map[string]string
}

// NewFont return a font selection struct
func NewFont(options string, height int) (*Font, error) {
	var attr map[string]string
	if err := json.Unmarshal([]byte(options), &attr); err != nil {
		return nil, err
	}
	f := &Font{
		Height: height,
		Attr:   attr,
	}
	var err error
	f.Fonter, err = driver.NewFont(f)
	return f, err
}

// Close destroys font selection
func (f *Font) Close() {
	f.Fonter.Close()
}

// Split breaks the text into lines that fit in the specified width
func (f *Font) Split(text string, edge int) []string {
	return f.Fonter.Split(text, edge)
}

// Size returns the width and height of the drawing area
func (f *Font) Size(text string) Size {
	return f.Fonter.Size(text)
}
