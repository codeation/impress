package domain

import (
	"image"
	"image/draw"
	"log"

	"github.com/codeation/impress/driver"
)

type picture struct {
	app    *application
	id     int
	width  int
	height int
}

func (app *application) NewImage(img image.Image) driver.Imager {
	pix, ok := img.(*image.NRGBA)
	if !ok {
		pix = image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(pix, pix.Bounds(), img, image.Pt(0, 0), draw.Src)
	}
	p := &picture{
		app:    app,
		id:     app.nextImageID(),
		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
	}
	if len(pix.Pix) > 67108863 {
		log.Printf("image size is too large")
		pix = image.NewNRGBA(image.Rect(0, 0, 1, 1))
		p.width = 1
		p.height = 1
	}
	app.caller.ImageNew(p.id, p.width, p.height, pix.Pix)
	return p
}

func (p *picture) Width() int  { return p.width }
func (p *picture) Height() int { return p.height }

func (p *picture) Close() {
	p.app.caller.ImageDrop(p.id)
}
