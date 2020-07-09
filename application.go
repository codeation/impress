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
	active    Actor
	events    chan Eventer
	destroyed chan struct{}
	wg        sync.WaitGroup
}

// NewApplication creates main application window
func NewApplication(rect Rect, title string) *Application {
	app := &Application{
		handlers:  map[Eventer]func(){},
		events:    make(chan Eventer, 10),
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

// Done returns a channel that's closed when the application is stopped
func (app *Application) Done() <-chan struct{} {
	return app.destroyed
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
	if app.active == actor {
		if len(regions) > 1 {
			app.active = regions[0].actor
		} else {
			app.active = nil
		}
	}
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

func (app *Application) singleEvent(e Eventer) {
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
					if ev.Point.In(r.rect) && app.active != r.actor {
						app.Activate(r.actor)
						break
					}
				}
				continue
			}
			if a := app.active; a != nil {
				a.Chan() <- e
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

// Activate sets the active event receiver under mutex protection
func (app *Application) Activate(actor Actor) {
	if a := app.active; a != nil {
		a.Chan() <- DeactivatedEvent
	}
	app.active = actor
	if actor != nil {
		actor.Chan() <- ActivatedEvent
	}
}

// Activated returns true when the actor is an active recipient
func (app *Application) Activated(actor Actor) bool {
	return app.active == actor
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
