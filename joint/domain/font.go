package domain

import (
	"image"
	"log"

	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/fontspec"
)

type font struct {
	app        *application
	id         int
	height     int
	lineheight int
	baseline   int
	ascent     int
	descent    int
}

func (app *application) NewFont(height int, attributes map[string]string) driver.Fonter {
	id := app.fontID.Next()
	family, style, variant, weight, stretch := fontspec.Attributes(attributes)
	lineheight, baseline, ascent, descent := app.caller.FontNew(id, height, style, variant, weight, stretch, family)
	return &font{
		app:        app,
		id:         id,
		height:     height,
		lineheight: lineheight,
		baseline:   baseline,
		ascent:     ascent,
		descent:    descent,
	}
}

func (f *font) LineHeight() int {
	return f.lineheight
}

func (f *font) Baseline() int {
	return f.baseline
}

func (f *font) Ascent() int {
	return f.ascent
}

func (f *font) Descent() int {
	return f.descent
}

func (f *font) Close() {
	f.app.caller.FontDrop(f.id)
	f.app.fontID.Back(f.id)
}

func (f *font) Split(text string, edge int, indent int) []string {
	if len(text) > 32767 {
		log.Printf("split text is too large: %d", len(text))
		text = ""
	}
	lengths := f.app.caller.FontSplit(f.id, text, edge, indent)
	return fontspec.SplitByLengths(text, lengths)
}

func (f *font) Size(text string) image.Point {
	x, y := f.app.caller.FontSize(f.id, text)
	return image.Pt(x, y)
}

func getFontID(f driver.Fonter) int {
	for {
		wrappedFont, ok := f.(interface{ Unwrap() driver.Fonter })
		if !ok {
			break
		}
		f = wrappedFont.Unwrap()
	}
	if localFont, ok := f.(*font); ok {
		return localFont.id
	}
	return 0
}
