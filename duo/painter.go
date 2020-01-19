package duo

import (
	"github.com/codeation/impress"
)

// Paint

type painter struct {
	driver     *driver
	id         int
	rect       impress.Rect
	background impress.Color
}

func (p *painter) Drop() {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'O', p.id)
}

func (p *painter) Size(rect impress.Rect) {
	p.rect = rect
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'Z', p.id, p.rect.X, p.rect.Y, p.rect.Width, p.rect.Height,
		p.background.R, p.background.G, p.background.B)
}

func (p *painter) Clear() {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'C', p.id)
	writeSequence(p.driver.connDraw, 'F', p.id, 0, 0, p.rect.Width, p.rect.Height,
		p.background.R, p.background.G, p.background.B)
}

func (p *painter) Fill(rect impress.Rect, color impress.Color) {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'F', p.id, rect.X, rect.Y, rect.Width, rect.Height,
		color.R, color.G, color.B)
}

func (p *painter) Line(from impress.Point, to impress.Point, color impress.Color) {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'L', p.id, from.X, from.Y, to.X, to.Y,
		color.R, color.G, color.B)
}

func (p *painter) Text(text string, font *impress.Font, from impress.Point, color impress.Color) {
	f := font.Fonter.(*ftfont)
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'U', p.id, from.X, from.Y, color.R, color.G, color.B,
		f.ID, font.Height, text)
}

func (p *painter) Show() {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'W', p.id)
}
