package duo

import (
	"image"
	"log"

	"github.com/codeation/impress/driver"
)

// Values for "style" atribute
var styleValues = map[string]int{
	"":        0,
	"normal":  0,
	"oblique": 1,
	"italic":  2,
}

// Values for "variant" atribute
var variantValues = map[string]int{
	"":           0,
	"normal":     0,
	"small_caps": 1,
	"small caps": 1,
}

// Values for "weight" atribute
var weightValues = map[string]int{
	"thin":       100,
	"ultralight": 200,
	"light":      300,
	"semilight":  350,
	"book":       380,
	"":           400,
	"normal":     400,
	"medium":     500,
	"semibold":   600,
	"bold":       700,
	"ultrabold":  800,
	"heavy":      900,
	"ultraheavy": 1000,
}

// Values for "stretch" atribute
var stretchValues = map[string]int{
	"ultra_condensed": 0,
	"ultra condensed": 0,
	"extra_condensed": 1,
	"extra condensed": 1,
	"condensed":       2,
	"semi_condensed":  3,
	"semi condensed":  3,
	"":                4,
	"normal":          4,
	"semi_expanded":   5,
	"semi expanded":   5,
	"expanded":        6,
	"extra_expanded":  7,
	"extra expanded":  7,
	"ultra_expanded":  8,
	"ultra expanded":  8,
}

type fontface struct {
	driver     *duo
	id         int
	height     int
	baseline   int
	ascent     int
	descent    int
	attributes map[string]string
}

func (d *duo) NewFont(height int, attributes map[string]string) driver.Fonter {
	if d == nil || d.drawPipe == nil {
		log.Fatal("GUI driver not initialized")
	}
	d.lastFontID++
	f := &fontface{
		driver:     d,
		id:         d.lastFontID,
		height:     height,
		attributes: attributes,
	}
	style := f.getValue("style", styleValues)
	variant := f.getValue("variant", variantValues)
	weight := f.getValue("weight", weightValues)
	stretch := f.getValue("stretch", stretchValues)
	d.drawPipe.
		Int16(&f.baseline).
		Int16(&f.ascent).
		Int16(&f.descent).
		Call(
			'N', f.id, f.height, style, variant, weight, stretch, f.attributes["family"])
	return f
}

func (f *fontface) Close() {}

func (f *fontface) Ascent() int   { return f.ascent }
func (f *fontface) Baseline() int { return f.baseline }
func (f *fontface) Descent() int  { return f.descent }

func (f *fontface) Split(text string, edge int) []string {
	if len(text) == 0 {
		return nil
	}
	var lengths []int
	f.driver.drawPipe.
		Int16s(&lengths).
		Call(
			'P', f.id, edge, text)
	output := make([]string, len(lengths))
	pos := 0
	for i := 0; i < len(lengths); i++ {
		output[i] = text[pos : pos+lengths[i]]
		pos += lengths[i]
	}
	return output
}

func (f *fontface) Size(text string) image.Point {
	if len(text) == 0 {
		return image.Pt(0, f.height)
	}
	var width, height int
	f.driver.drawPipe.
		Int16(&width).
		Int16(&height).
		Call(
			'R', f.id, text)
	return image.Pt(width, height)
}

func (f *fontface) getValue(fieldname string, values map[string]int) int {
	source, ok := f.attributes[fieldname]
	if !ok {
		return values[""]
	}
	value, ok := values[source]
	if !ok {
		return values[""]
	}
	return value
}
