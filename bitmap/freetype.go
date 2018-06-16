package bitmap

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"io/ioutil"
)

type Font struct {
	context *freetype.Context
	face    font.Face
	rect    Rect
}

func OpenFont(name string, fontsize int) (*Font, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, errors.Wrap(err, "NewFont")
	}
	return NewFont(data, fontsize)
}

func NewFont(data []byte, fontsize int) (*Font, error) {
	fnt, err := freetype.ParseFont(data)
	if err != nil {
		return nil, errors.Wrap(err, "NewFont")
	}

	opts := &truetype.Options{
		Size:              float64(fontsize),
		DPI:               72,
		Hinting:           font.HintingFull,
		GlyphCacheEntries: 64,
	}
	face := truetype.NewFace(fnt, opts)

	bounds := fnt.Bounds(fixed.I(fontsize))
	rect := NewRect(bounds.Min.X.Floor(), bounds.Min.Y.Floor(),
		-bounds.Min.X.Floor(), bounds.Max.Y.Round()-bounds.Min.Y.Floor())

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fnt)
	c.SetFontSize(float64(fontsize))
	c.SetHinting(font.HintingFull)

	return &Font{
		context: c,
		face:    face,
		rect:    rect,
	}, nil
}

func (f *Font) Close() {
	f.face.Close()
}

func (f *Font) DrawString(text string, point Point) (Rect, error) {
	pt, err := f.context.DrawString(text, freetype.Pt(point.X, point.Y))
	rect := f.rect
	rect.X += point.X
	rect.Y += point.Y
	rect.Width += pt.X.Ceil() - point.X
	return rect, err
}

func (f *Font) Ascent() int {
	return f.face.Metrics().Ascent.Round()
}

func (f *Font) Descent() int {
	return f.face.Metrics().Descent.Round()
}

func (f *Font) Height() int {
	return f.face.Metrics().Height.Round()
}

func (f *Font) Size(text string) Size {
	width := fixed.I(0)
	var prevrune rune
	for i, r := range text {
		adv, ok := f.face.GlyphAdvance(r)
		if ok {
			width += adv
		}
		if i > 0 {
			width += f.face.Kern(prevrune, r)
		}
		prevrune = r
	}
	return Size{
		Width:  width.Ceil(),
		Height: f.face.Metrics().Height.Round(),
	}
}

func (f *Font) Split(text string, edge int) []string {
	out := make([]string, 0)
	edgeI := fixed.I(edge)
	for len(text) > 0 {
		width := fixed.I(0)
		lastspace := 0
		spacewidth := fixed.I(0)
		current := 0
		var prevrune rune
		for i, r := range text {
			if r == ' ' {
				lastspace = i
				spacewidth = width
			}
			adv, ok := f.face.GlyphAdvance(r)
			if ok {
				width += adv
			}
			if i > 0 {
				width += f.face.Kern(prevrune, r)
			}
			prevrune = r
			current = i
			if width > edgeI {
				break
			}
		}
		if width <= edgeI {
			out = append(out, text)
			break
		} else if spacewidth > edgeI*2/3 {
			out = append(out, text[:lastspace])
			text = text[lastspace+1:]
		} else {
			out = append(out, text[:current])
			text = text[current:]
		}
	}
	return out
}
