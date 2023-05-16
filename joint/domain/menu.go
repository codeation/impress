package domain

import (
	"github.com/codeation/impress/driver"
)

type menu struct {
	app *application
	id  int
}

func (app *application) NewMenu(label string) driver.Menuer {
	id := app.nextMenuID()
	app.caller.MenuNew(id, 0, label)
	return &menu{
		app: app,
		id:  id,
	}
}

func (m *menu) NewMenu(label string) driver.Menuer {
	id := m.app.nextMenuID()
	m.app.caller.MenuNew(id, m.id, label)
	return &menu{
		app: m.app,
		id:  id,
	}
}

func (m *menu) NewItem(label string, action string) {
	id := m.app.nextMenuID()
	m.app.caller.MenuItem(id, m.id, label, action)
}
