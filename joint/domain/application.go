// Package implements an internal mechanism to communicate with an impress terminal.
package domain

import (
	"image"
	"sync"

	"github.com/codeation/impress/event"
	"github.com/codeation/impress/joint/iface"
)

type chaner interface {
	Chan() <-chan event.Eventer
}

type callSetSyncer interface {
	iface.CallSet
	Sync()
}

type application struct {
	caller       callSetSyncer
	chaner       chaner
	lastFrameID  int
	lastWindowID int
	lastFontID   int
	lastImageID  int
	lastMenuID   int
	mutex        sync.Mutex
}

func New(caller callSetSyncer, ch chaner) *application {
	return &application{
		caller: caller,
		chaner: ch,
	}
}

func (a *application) nextFrameID() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.lastFrameID++
	return a.lastFrameID
}

func (a *application) nextWindowID() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.lastWindowID++
	return a.lastWindowID
}

func (a *application) nextFontID() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.lastFontID++
	return a.lastFontID
}

func (a *application) nextImageID() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.lastImageID++
	return a.lastImageID
}

func (a *application) nextMenuID() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.lastMenuID++
	return a.lastMenuID
}

func (a *application) Init() {}

func (a *application) Done() {
	a.caller.ApplicationExit()
}

func (a *application) Title(title string) {
	a.caller.ApplicationTitle(title)
}

func (a *application) Size(rect image.Rectangle) {
	x, y, width, height := rectangle(rect)
	a.caller.ApplicationSize(x, y, width, height)
}

func (a *application) Chan() <-chan event.Eventer {
	return a.chaner.Chan()
}

func (a *application) Sync() {
	a.caller.Sync()
}
