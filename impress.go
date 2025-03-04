package impress

import (
	"image"
	"log"

	"github.com/codeation/impress/driver"
)

// NewApplication creates the main application window for the default GUI driver.
// Consider using the MakeApplication function instead.
var NewApplication = func(rect image.Rectangle, title string) *Application {
	log.Fatalf("GUI driver must be registered. Add GTK driver to use by default:\nimport _ \"github.com/codeation/impress/duo\"")
	return nil
}

// NewImage returns an Image struct containing the image resources for the default GUI driver.
//
// Deprecated: Use the *Application.NewImage function instead.
var NewImage func(img image.Image) *Image

// NewFont returns a Font struct for the default GUI driver.
//
// Deprecated: Use the *Application.NewFont function instead.
var NewFont func(height int, attributes map[string]string) *Font

// Register is an internal function that makes the GUI driver available.
func Register(d driver.Driver) {
	NewApplication = func(rect image.Rectangle, title string) *Application {
		app := MakeApplication(d, rect, title)
		NewImage = app.NewImage
		NewFont = app.NewFont
		return app
	}
}
