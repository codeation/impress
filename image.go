package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// Image represents a draw-ready image
type Image struct {
	imager driver.Imager
	Width  int
	Height int
}

// NewImage returns a image resources struct
func NewImage(img image.Image) *Image {
	imager := d.NewImage(img)
	return &Image{
		imager: imager,
		Width:  imager.Width(),
		Height: imager.Height(),
	}
}

// Close destroys image resources
func (i *Image) Close() {
	i.imager.Close()
	i.imager = nil // TODO notice when the image is closed
}
