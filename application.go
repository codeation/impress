package impress

// Application represents application top level window
type Application struct{}

// NewApplication creates main application window
func NewApplication(rect Rect, title string) *Application {
	app := &Application{}
	driver.Init()
	driver.Size(rect)
	driver.Title(title)
	return app
}

// Close makes invocation of the main loop return
func (app *Application) Close() {
	driver.Done()
}

// Title sets application window title
func (app *Application) Title(title string) {
	driver.Title(title)
}

// Size sets application window size
func (app *Application) Size(rect Rect) {
	driver.Size(rect)
}

func (app *Application) Chan() <-chan Eventer {
	return driver.Chan()
}
