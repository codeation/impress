package impress

import (
	"image"
	"log"

	"github.com/codeation/impress/event"
)

// Application represents application top level window
type Application struct{}

// NewApplication creates main application window
func NewApplication(rect image.Rectangle, title string) *Application {
	if d == nil {
		log.Fatal("GUI driver is not available")
	}
	d.Init()
	d.Size(rect)
	d.Title(title)
	return &Application{}
}

// Close destroys application resources
func (app *Application) Close() {
	d.Done()
}

// Title sets application window title
func (app *Application) Title(title string) {
	d.Title(title)
}

// Size sets application window size
func (app *Application) Size(rect image.Rectangle) {
	d.Size(rect)
}

// Chan returns event channel
func (app *Application) Chan() <-chan event.Eventer {
	return d.Chan()
}
