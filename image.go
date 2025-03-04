package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// Image represents a drawable image.
type Image struct {
	imager driver.Imager
	Size   image.Point
}

// NewImage returns an Image struct containing the image resources.
func (app *Application) NewImage(img image.Image) *Image {
	imager := app.driver.NewImage(img)
	return &Image{
		imager: imager,
		Size:   imager.Size(),
	}
}

// Close destroys the image resources.
// Note that a closed image can no longer be used.
func (i *Image) Close() {
	i.imager.Close()
	i.imager = nil // TODO: Add notice when the image is closed.
}
