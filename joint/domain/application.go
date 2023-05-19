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

type flusher interface {
	Flush()
}

type application struct {
	caller       iface.CallSet
	chaner       chaner
	flusher      flusher
	lastFrameID  int
	lastWindowID int
	lastFontID   int
	lastImageID  int
	lastMenuID   int
	mutex        sync.Mutex
}

func New(caller iface.CallSet, ch chaner, flusher flusher) *application {
	return &application{
		caller:  caller,
		chaner:  ch,
		flusher: flusher,
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
	a.flusher.Flush()
}
