package action

import (
	"github.com/codeation/impress"
)

// Action presents a simple event receiver
type Action struct {
	app    *impress.Application
	events chan impress.Eventer
}

// NewAction launches new event receiver func
func NewAction(app *impress.Application, rect impress.Rect, loop func(act *Action)) *Action {
	act := &Action{
		app:    app,
		events: make(chan impress.Eventer, 10),
	}
	app.AddActor(act, rect)
	app.Start(func() {
		loop(act)
		app.RemoveActor(act)
	})
	if app.Activated(nil) {
		app.Activate(act)
	}
	return act
}

// Chan returns receiver events channel
func (act *Action) Chan() chan impress.Eventer {
	return act.events
}

// Event gets next event when the Action is an active recipient
func (act *Action) Event() impress.Eventer {
	select {
	case e := <-act.events:
		return e
	case <-act.app.Done():
		return impress.DoneEvent
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
