package impress

// Application represents application top level window
type Application struct {
	handlers map[Eventer]func()
}

// NewApplication creates main application window
func NewApplication() *Application {
	a := &Application{
		handlers: map[Eventer]func(){},
	}
	driver.Init()
	return a
}

// Main runs the main GUI loop until Quit is called
func (a *Application) Main() {
	driver.Main()
}

// Quit makes invocation of the main loop return
func (a *Application) Quit() {
	driver.Done()
}

// OnEvent connects function call back to an event
func (a *Application) OnEvent(event Eventer, handler func()) {
	a.handlers[event] = handler
}

// Event returns next application event
func (a *Application) Event() Eventer {
	for {
		e := driver.Event()
		handler, ok := a.handlers[e]
		if !ok {
			return e
		}
		handler()
	}
}

// Title sets application window title
func (a *Application) Title(s string) {
	driver.Title(s)
}

// Size sets application window size
func (a *Application) Size(rect Rect) {
	driver.Size(rect)
}
