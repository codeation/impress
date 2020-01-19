package impress

// Action presents the event receiver
type Action struct {
	app    *Application
	events chan Eventer
}

// NewAction creates new event receiver
func (app *Application) NewAction() *Action {
	return &Action{
		app:    app,
		events: make(chan Eventer, 10),
	}
}

// Chan returns receiver events channel
func (act *Action) Chan() chan Eventer {
	return act.events
}

// Event gets next event when the Action is an active recepient
func (act *Action) Event() Eventer {
	select {
	case e := <-act.events:
		return e
	case <-act.app.destroyed:
		return DoneEvent
	}
}

// Activate enables the Action to receive app events
func (act *Action) Activate() {
	act.app.Activate(act)
}

// Activated returns true when the Action is an active recipient
func (act *Action) Activated() bool {
	return act.app.Activated(act)
}

// Deactivate disables the Action to receive app events
func (act *Action) Deactivate() {
	act.app.Activate(nil)
}

// ExitApplication exit Application of event receiver
func (act *Action) ExitApplication() {
	close(act.app.destroyed)
}
