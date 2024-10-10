package impress

import (
	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/event"
)

// Menu represents any menu node
type Menu struct {
	menuer driver.Menuer
}

// NewMenu returns a new top-level menu node
func (app *Application) NewMenu(label string) *Menu {
	return &Menu{
		menuer: app.driver.NewMenu(label),
	}
}

// NewMenu returns new submenu node
func (m *Menu) NewMenu(label string) *Menu {
	return &Menu{
		menuer: m.menuer.NewMenu(label),
	}
}

// NewItem adds a item to menu node
func (m *Menu) NewItem(label string, event event.Menu) {
	m.menuer.NewItem(label, event.Action)
}
