package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// Image represents a draw-ready image
type Image struct {
	driver.Imager
	Size image.Point
}

// NewImage returns a image resources struct
func (app *Application) NewImage(img image.Image) *Image {
	imager := app.driver.NewImage(img)
	return &Image{
		Imager: imager,
		Size:   imager.Size(),
	}
}

// Close destroys image resources
func (i *Image) Close() {
	i.Imager.Close()
	i.Imager = nil // TODO notice when the image is closed
}
