package impress

// Actor presents the event receiver
type Action struct {
	app    *Application
	events chan Eventer
}

// NewActor creates new event receiver
func (app *Application) NewAction() *Action {
	return &Action{
		app:    app,
		events: make(chan Eventer, 10),
	}
}

func (act *Action) Chan() chan Eventer {
	return act.events
}

// Event gets next event
func (act *Action) Event() Eventer {
	select {
	case e := <-act.events:
		return e
	case <-act.app.destroyed:
		return DoneEvent
	}
}

// Activate enables the actor to receive app events
func (act *Action) Activate() {
	act.app.Activate(act)
}

// Activated returns true when the actor is an active recipient
func (act *Action) Activated() bool {
	return act.app.Activated(act)
}

// Deactivate disables the actor to receive app events
func (act *Action) Deactivate() {
	act.app.Activate(nil)
}
