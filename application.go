package impress

import (
	"image"
	"log"

	"github.com/codeation/impress/clipboard"
	"github.com/codeation/impress/event"
)

// Application represents application top level window
type Application struct {
	frame *Frame
}

// NewApplication creates main application window
func NewApplication(rect image.Rectangle, title string) *Application {
	if d == nil {
		log.Printf("GUI driver is not available")
		return nil
	}
	d.Init()
	d.Size(rect)
	d.Title(title)
	return &Application{
		frame: &Frame{framer: d.NewFrame(image.Rect(0, 0, rect.Dx(), rect.Dy()))},
	}
}

// Close destroys application resources
func (app *Application) Close() {
	app.frame.Drop()
	d.Done()
}

// Sync flushes graphics content to screen driver
func (app *Application) Sync() {
	d.Sync()
}

// Title sets application window title
func (app *Application) Title(title string) {
	d.Title(title)
}

// Size sets application window size
func (app *Application) Size(rect image.Rectangle) {
	d.Size(rect)
}

// ClipboardGet requests event with clipboard content
func (app *Application) ClipboardGet(typeID int) {
	d.ClipboardGet(typeID)
}

// ClipboardPut set content to OS clipboard
func (app *Application) ClipboardPut(c clipboard.Clipboarder) {
	d.ClipboardPut(c)
}

// Chan returns event channel
func (app *Application) Chan() <-chan event.Eventer {
	return d.Chan()
}
