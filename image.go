package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// Image represents a draw-ready image
type Image struct {
	imager driver.Imager
	Size   image.Point
}

// NewImage returns a image resources struct
func (app *Application) NewImage(img image.Image) *Image {
	imager := app.driver.NewImage(img)
	return &Image{
		imager: imager,
		Size:   imager.Size(),
	}
}

// Close destroys image resources
func (i *Image) Close() {
	i.imager.Close()
	i.imager = nil // TODO notice when the image is closed
}
