package duo

import (
	"log"

	"github.com/codeation/impress"
)

type bitmap struct {
	driver *driver
	ID     int
	Image  *impress.Image
}

func (d *driver) NewImage(img *impress.Image) (impress.Imager, error) {
	if d == nil || d.pipeDraw == nil {
		log.Fatal("GUI driver not initialized")
	}
	d.lastImageID++
	b := &bitmap{
		driver: d,
		ID:     d.lastImageID,
		Image:  img,
	}
	b.driver.onDraw.Lock()
	defer b.driver.onDraw.Unlock()
	writeSequence(b.driver.pipeDraw, 'B', b.ID, img.Width, img.Height, img.PixNRGBA)
	return b, nil
}

func (b *bitmap) Close() {
	b.driver.onDraw.Lock()
	defer b.driver.onDraw.Unlock()
	writeSequence(b.driver.pipeDraw, 'M', b.ID)
}
