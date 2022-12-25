package duo

import (
	"image"
	"image/draw"
	"log"

	"github.com/codeation/impress/driver"
)

type bitmap struct {
	driver *duo
	id     int
	width  int
	height int
}

func (d *duo) NewImage(img image.Image) driver.Imager {
	if d == nil || d.streamPipe == nil {
		log.Fatal("GUI driver not initialized")
	}
	pix, ok := img.(*image.NRGBA)
	if !ok {
		pix = image.NewNRGBA(image.Rect(0, 0, img.Bounds().Size().X, img.Bounds().Size().Y))
		draw.Draw(pix, pix.Bounds(), img, image.Pt(0, 0), draw.Src)
	}
	d.lastImageID++
	b := &bitmap{
		driver: d,
		id:     d.lastImageID,
		width:  img.Bounds().Size().X,
		height: img.Bounds().Size().Y,
	}
	b.driver.streamPipe.Call(
		'B', b.id, b.width, b.height, pix.Pix)
	return b
}

func (b *bitmap) Close() {
	b.driver.streamPipe.Call(
		'M', b.id)
}

func (b *bitmap) Width() int  { return b.width }
func (b *bitmap) Height() int { return b.height }
