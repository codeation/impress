package domain

import (
	"image"
	"image/draw"
	"log"

	"github.com/codeation/impress/driver"
)

type picture struct {
	app  *application
	id   int
	size image.Point
}

func (app *application) NewImage(img image.Image) driver.Imager {
	pix, ok := img.(*image.NRGBA)
	if !ok {
		pix = image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(pix, pix.Bounds(), img, image.Pt(0, 0), draw.Src)
	}
	p := &picture{
		app:  app,
		id:   app.nextImageID(),
		size: img.Bounds().Size(),
	}
	if len(pix.Pix) > 67108863 {
		log.Printf("image size is too large: %d", len(pix.Pix))
		pix = image.NewNRGBA(image.Rect(0, 0, 1, 1))
		p.size = image.Pt(1, 1)
	}
	app.caller.ImageNew(p.id, p.size.X, p.size.Y, pix.Pix)
	return p
}

func (p *picture) Size() image.Point { return p.size }

func (p *picture) Close() {
	p.app.caller.ImageDrop(p.id)
}

func (p *picture) ID() int {
	return p.id
}
