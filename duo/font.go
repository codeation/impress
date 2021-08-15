package duo

import (
	"log"

	"github.com/codeation/impress"
)

var tableStyle = map[string]int{
	"":        0,
	"normal":  0,
	"oblique": 1,
	"italic":  2,
}

var tableVariant = map[string]int{
	"":           0,
	"normal":     0,
	"small_caps": 1,
	"small caps": 1,
}

var tableWeight = map[string]int{
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

var tableStretch = map[string]int{
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

type ftfont struct {
	driver *driver
	ID     int
	Font   *impress.Font
}

func (f *ftfont) defaultValue(fieldname string, maps *map[string]int) int {
	source, ok := f.Font.Attr[fieldname]
	if !ok {
		return (*maps)[""]
	}
	value, ok := (*maps)[source]
	if !ok {
		return (*maps)[""]
	}
	return value
}

func (d *driver) NewFont(font *impress.Font) (impress.Fonter, error) {
	if d == nil || d.pipeDraw == nil {
		log.Fatal("GUI driver not initialized")
	}
	d.lastFontID++
	f := &ftfont{
		driver: d,
		ID:     d.lastFontID,
		Font:   font,
	}
	style := f.defaultValue("style", &tableStyle)
	variant := f.defaultValue("variant", &tableVariant)
	weight := f.defaultValue("weight", &tableWeight)
	stretch := f.defaultValue("stretch", &tableStretch)
	f.driver.onDraw.Lock()
	defer f.driver.onDraw.Unlock()
	writeSequence(f.driver.pipeDraw, 'N', f.ID, f.Font.Height, style, variant, weight, stretch,
		f.Font.Attr["family"])
	font.Baseline, _ = readInt16(f.driver.pipeAnswer)
	font.Ascent, _ = readInt16(f.driver.pipeAnswer)
	font.Descent, _ = readInt16(f.driver.pipeAnswer)
	return f, nil
}

func (f *ftfont) Close() {}

func (f *ftfont) Split(text string, edge int) []string {
	if len(text) == 0 {
		return nil
	}
	f.driver.onDraw.Lock()
	defer f.driver.onDraw.Unlock()
	writeSequence(f.driver.pipeDraw, 'P', f.ID, edge, text)
	count, _ := readInt16(f.driver.pipeAnswer)
	pos := 0
	out := make([]string, count)
	for i := 0; i < count; i++ {
		length, _ := readInt16(f.driver.pipeAnswer)
		out[i] = text[pos : pos+length]
		pos += length
	}
	return out
}

func (f *ftfont) Size(text string) impress.Size {
	if len(text) == 0 {
		return impress.NewSize(0, f.Font.Height)
	}
	f.driver.onDraw.Lock()
	defer f.driver.onDraw.Unlock()
	writeSequence(f.driver.pipeDraw, 'R', f.ID, text)
	width, _ := readInt16(f.driver.pipeAnswer)
	height, _ := readInt16(f.driver.pipeAnswer)
	return impress.NewSize(width, height)
}
