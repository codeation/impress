package impress

import (
	"image"

	"github.com/codeation/impress/clipboard"
	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/event"
)

// Application represents the top-level window of the application.
type Application struct {
	driver driver.Driver
	frame  *Frame
}

// MakeApplication creates the top application window for the specified GUI driver.
func MakeApplication(d driver.Driver, rect image.Rectangle, title string) *Application {
	app := &Application{
		driver: d,
	}
	app.driver.Init()
	app.driver.Size(rect)
	app.driver.Title(title)
	app.frame = &Frame{framer: app.driver.NewFrame(image.Rect(0, 0, rect.Dx(), rect.Dy()))}
	return app
}

// Close destroys application resources.
func (app *Application) Close() {
	app.frame.Drop()
	app.driver.Done()
}

// Sync flushes the graphics content to the screen driver.
func (app *Application) Sync() {
	app.driver.Sync()
}

// Title sets the application window title.
func (app *Application) Title(title string) {
	app.driver.Title(title)
}

// Size sets the application window size.
func (app *Application) Size(rect image.Rectangle) {
	app.driver.Size(rect)
}

// ClipboardGet requests an event with the clipboard content.
func (app *Application) ClipboardGet(typeID int) {
	app.driver.ClipboardGet(typeID)
}

// ClipboardPut sets the content to the OS clipboard.
func (app *Application) ClipboardPut(c clipboard.Clipboarder) {
	app.driver.ClipboardPut(c)
}

// Chan returns the event channel.
func (app *Application) Chan() <-chan event.Eventer {
	return app.driver.Chan()
}
