// Package implements an internal mechanism to communicate with an impress terminal.
package domain

import (
	"image"
	"log"

	"github.com/codeation/impress/clipboard"
	"github.com/codeation/impress/event"
	"github.com/codeation/impress/joint/idcycle"
	"github.com/codeation/impress/joint/iface"
)

type chaner interface {
	Chan() <-chan event.Eventer
}

type syncer interface {
	Sync() error
}

type application struct {
	caller   iface.CallSet
	chaner   chaner
	syncer   syncer
	frameID  *idcycle.ID
	windowID *idcycle.ID
	fontID   *idcycle.ID
	imageID  *idcycle.ID
	menuID   *idcycle.ID
}

func New(caller iface.CallSet, ch chaner, syncer syncer) *application {
	return &application{
		caller:   caller,
		chaner:   ch,
		syncer:   syncer,
		frameID:  idcycle.New(),
		windowID: idcycle.New(),
		fontID:   idcycle.New(),
		imageID:  idcycle.New(),
		menuID:   idcycle.New(),
	}
}

func (a *application) Init() {
	driverVersion := a.caller.ApplicationVersion()
	if driverVersion != apiVersion {
		log.Fatalf("Unexpected driver API version: %s", driverVersion)
	}
}

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
	go a.syncer.Sync() // nolint:errcheck
}

func (a *application) ClipboardGet(typeID int) {
	a.caller.ClipboardGet(typeID)
}

func (a *application) ClipboardPut(c clipboard.Clipboarder) {
	a.caller.ClipboardPut(c.Type(), c.Data())
}
