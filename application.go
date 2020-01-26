package impress

import (
	"sync"
)

type region struct {
	actor Actor
	rect  Rect
}

// Application represents application top level window
type Application struct {
	handlers  map[Eventer]func()
	regions   []*region
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

// Close makes invocation of the main loop return
func (app *Application) Close() {
	driver.Done()
}

// OnEvent connects function call back to an event
func (app *Application) OnEvent(event Eventer, handler func()) {
	app.handlers[event] = handler
}

// AddActor adds actor for rect
func (app *Application) AddActor(actor Actor, rect Rect) {
	app.regions = append(app.regions, &region{actor, rect})
}

// RemoveActor removes actor from actor list
func (app *Application) RemoveActor(actor Actor) {
	regions := []*region{}
	for _, r := range app.regions {
		if r.actor == actor {
			continue
		}
		regions = append(regions, r)
	}
	app.regions = regions
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

// Redirect drivers events to active Actor or Application itself
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
			if ev, ok := e.(ButtonEvent); ok &&
				ev.Action == ButtonActionPress && ev.Button == ButtonLeft {
				for _, r := range app.regions {
					if ev.Point.In(r.rect) && !app.Activated(r.actor) {
						app.Activate(r.actor)
						break
					}
				}
			}
			if app.active != nil {
				app.active.Chan() <- e
				continue
			}
			app.events <- e
		case <-app.destroyed:
			return
		}
	}
}

// Activate sets the active event receiver
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

// Wait waits for main application loop to complete
func (app *Application) Wait() {
	app.wg.Wait()
}
