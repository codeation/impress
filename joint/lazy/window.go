package lazy

import (
	"image"
	"image/color"

	"github.com/codeation/impress/driver"
)

type fillPaint struct {
	rect       image.Rectangle
	foreground color.Color
}

type linePaint struct {
	from, to   image.Point
	foreground color.Color
}

type imagePaint struct {
	rect image.Rectangle
	img  driver.Imager
}

type textPaint struct {
	text       string
	font       driver.Fonter
	from       image.Point
	foreground color.Color
}

type window struct {
	driver.Painter
	history  []interface{}
	reAdjust bool
	position int
	rect     image.Rectangle
}

func (f *frame) NewWindow(rect image.Rectangle, background color.Color) driver.Painter {
	return &window{
		Painter: f.Framer.NewWindow(rect, background),
		rect:    rect,
	}
}

func (w *window) Size(rect image.Rectangle) {
	if w.rect == rect {
		return
	}
	w.reAdjust = true
	w.rect = rect
	w.Painter.Size(rect)
}

func (w *window) Clear() {
	w.position = 0
}

func (w *window) Show() {
	if !w.reAdjust && len(w.history) == w.position {
		return
	}
	w.reAdjust = false
	w.history = w.history[:w.position]
	w.Painter.Clear()
	for _, h := range w.history {
		switch p := h.(type) {
		case *fillPaint:
			w.Painter.Fill(p.rect, p.foreground)
		case *linePaint:
			w.Painter.Line(p.from, p.to, p.foreground)
		case *imagePaint:
			w.Painter.Image(p.rect, p.img)
		case *textPaint:
			w.Painter.Text(p.text, p.font, p.from, p.foreground)
		}
	}
	w.Painter.Show()
}

func (w *window) Fill(rect image.Rectangle, foreground color.Color) {
	if len(w.history) > w.position {
		p, ok := w.history[w.position].(*fillPaint)
		if ok && p.rect == rect && p.foreground == foreground {
			w.position++
			return
		}
		w.history = w.history[:w.position]
	}
	w.reAdjust = true
	w.history = append(w.history, &fillPaint{
		rect:       rect,
		foreground: foreground,
	})
	w.position++
}

func (w *window) Line(from image.Point, to image.Point, foreground color.Color) {
	if len(w.history) > w.position {
		p, ok := w.history[w.position].(*linePaint)
		if ok && p.from == from && p.to == to && p.foreground == foreground {
			w.position++
			return
		}
		w.history = w.history[:w.position]
	}
	w.reAdjust = true
	w.history = append(w.history, &linePaint{
		from:       from,
		to:         to,
		foreground: foreground,
	})
	w.position++
}

func (w *window) Image(rect image.Rectangle, img driver.Imager) {
	if len(w.history) > w.position {
		p, ok := w.history[w.position].(*imagePaint)
		if ok && p.rect == rect && p.img == img {
			w.position++
			return
		}
		w.history = w.history[:w.position]
	}
	w.reAdjust = true
	w.history = append(w.history, &imagePaint{
		rect: rect,
		img:  img,
	})
	w.position++
}

func (w *window) Text(text string, font driver.Fonter, from image.Point, foreground color.Color) {
	if len(w.history) > w.position {
		p, ok := w.history[w.position].(*textPaint)
		if ok && p.text == text && p.font == font && p.from == from && p.foreground == foreground {
			w.position++
			return
		}
		w.history = w.history[:w.position]
	}
	w.reAdjust = true
	w.history = append(w.history, &textPaint{
		text:       text,
		font:       font,
		from:       from,
		foreground: foreground,
	})
	w.position++
}
