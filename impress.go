package impress

import (
	"image"

	"github.com/codeation/impress/driver"
)

// NewApplication creates main application window for default GUI driver.
// Consider using MakeApplication func instead
var NewApplication func(rect image.Rectangle, title string) *Application

// NewImage returns a image resources struct for default GUI driver.
//
// Deprecated: Use *Application.NewImage func instead.
var NewImage func(img image.Image) *Image

// NewFont return a font selection struct for default GUI driver.
//
// Deprecated: Use *Application.NewFont func instead.
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
