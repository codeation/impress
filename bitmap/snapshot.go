package bitmap

import (
	"image"
	"image/color"
	"image/draw"
)

func rgba(value Color) color.Color {
	return color.RGBA{
		R: uint8(value.R),
		G: uint8(value.G),
		B: uint8(value.B),
		A: 255,
	}
}

type Snapshot struct {
	picture *image.RGBA
}

func NewSnapshot(size Size, background Color) *Snapshot {
	r := image.Rect(0, 0, size.Width, size.Height)
	s := &Snapshot{
		picture: image.NewRGBA(r),
	}
	fill := image.NewUniform(rgba(background))
	draw.Draw(s.picture, r, fill, image.Pt(0, 0), draw.Over)
	return s
}

func (s *Snapshot) Picture() *image.RGBA {
	return s.picture
}

func (s *Snapshot) Fill(rect Rect, foreground Color) {
	r := image.Rect(rect.X, rect.Y, rect.X+rect.Width, rect.Y+rect.Height)
	fill := image.NewUniform(rgba(foreground))
	draw.Draw(s.picture, r, fill, image.Pt(0, 0), draw.Over)
}

func (s *Snapshot) Line(from, to Point, foreground Color) {
	r := image.Rect(min(from.X, to.X), min(from.Y, to.Y), max(from.X, to.X)+1, max(from.Y, to.Y)+1)
	fill := image.NewUniform(rgba(foreground))
	draw.Draw(s.picture, r, fill, image.Pt(0, 0), draw.Over)
}

func (s *Snapshot) Text(text string, font *Font, point Point, foreground Color) (Rect, error) {
	font.context.SetDst(s.picture)
	font.context.SetClip(s.picture.Bounds())
	font.context.SetSrc(image.NewUniform(rgba(foreground)))
	return font.DrawString(text, point)
}
