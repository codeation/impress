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
	active    Actor
}

// NewApplication creates main application window
func NewApplication() *Application {
	app := &Application{
		handlers:  map[Eventer]func(){},
		events:    make(chan Eventer, 10),
		destroyed: make(chan struct{}),
	}
	driver.Init()
	app.Start(app.eventLoop)
	return app
}

// Quit makes invocation of the main loop return
func (app *Application) Quit() {
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

// Redirect drivers events to active window or app itself
func (app *Application) eventLoop() {
	for {
		e := <-driver.Chan()
		if e == DestroyEvent {
			close(app.destroyed)
			break
		}
		if handler, ok := app.handlers[e]; ok {
			handler()
			continue
		}
		if app.active != nil {
			app.active.Chan() <- e
			continue
		}
		app.events <- e
	}
}

// SetActive sets the event receiver
func (app *Application) Activate(act Actor) {
	if app.active != nil {
		app.active.Chan() <- DeactivatedEvent
	}
	app.active = act
	if app.active != nil {
		app.active.Chan() <- ActivatedEvent
	}
}

// Activated returns true when the actor is an active recipient
func (app *Application) Activated(act Actor) bool {
	return app.active == act
}

// Title sets application window title
func (app *Application) Title(s string) {
	driver.Title(s)
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

// Wait waits for the functions to complete
func (app *Application) Wait() {
	app.wg.Wait()
}
