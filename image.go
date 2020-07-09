package impress

import (
	"image"
	"image/draw"
)

// Image represents a draw-ready image
type Image struct {
	Imager        Imager
	Width, Height int
	PixNRGBA      []byte
}

// NewImage returns a image resources struct
func NewImage(img image.Image) (*Image, error) {
	imgNRGBA, ok := img.(*image.NRGBA)
	if !ok {
		imgNRGBA = image.NewNRGBA(image.Rect(0, 0, img.Bounds().Size().X, img.Bounds().Size().Y))
		draw.Draw(imgNRGBA, imgNRGBA.Bounds(), img, image.Pt(0, 0), draw.Src)
	}
	i := &Image{
		Width:    img.Bounds().Size().X,
		Height:   img.Bounds().Size().Y,
		PixNRGBA: imgNRGBA.Pix,
	}
	imager, err := driver.NewImage(i)
	if err != nil {
		return nil, err
	}
	i.Imager = imager
	i.PixNRGBA = nil // remove memory ref
	return i, nil
}

// Close destroys image resources
func (i *Image) Close() {
	i.Imager.Close()
}
