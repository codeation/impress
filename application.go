package impress

import (
	"sync"
)

// Application represents application top level window
type Application struct {
	handlers  map[Eventer]func()
	events    chan Eventer
	destroyed chan struct{}
	wg        sync.WaitGroup
}

// NewApplication creates main application window
func NewApplication(rect Rect, title string) *Application {
	app := &Application{
		handlers:  map[Eventer]func(){},
		events:    make(chan Eventer, 64),
		destroyed: make(chan struct{}),
	}
	driver.Init()
	driver.Size(rect)
	driver.Title(title)
	app.Start(app.eventLoop)
	return app
}

// Close makes invocation of the main loop return
func (app *Application) Close() {
	driver.Done()
}

// OnEvent connects function call back to an event
func (app *Application) OnEvent(event Eventer, handler func()) {
	app.handlers[event] = handler
}

// Event returns next application event
func (app *Application) Event() Eventer {
	for {
		select {
		case e := <-app.events:
			return e
		case <-app.destroyed:
			return DoneEvent
		}
	}
}

// Redirect drivers events to Application events chan
func (app *Application) eventLoop() {
	for {
		select {
		case e := <-driver.Chan():
			if e == DestroyEvent {
				close(app.destroyed)
				return
			}
			if handler, ok := app.handlers[e]; ok {
				handler()
				continue
			}
			select {
			case app.events <- e:
			default:
			}
		case <-app.destroyed:
			return
		}
	}
}

// Title sets application window title
func (app *Application) Title(title string) {
	driver.Title(title)
}

// Size sets application window size
func (app *Application) Size(rect Rect) {
	driver.Size(rect)
}

// Start starts the specified func but does not wait for it to complete
func (app *Application) Start(f func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		f()
	}()
}

// Wait waits for main application loop to complete
func (app *Application) Wait() {
	app.wg.Wait()
}
