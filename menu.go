package impress

import (
	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/event"
)

// Menu represents any menu node.
type Menu struct {
	menuer driver.Menuer
}

// NewMenu returns a new top-level menu node with the specified label.
func (app *Application) NewMenu(label string) *Menu {
	return &Menu{
		menuer: app.driver.NewMenu(label),
	}
}

// NewMenu returns a new submenu node with the specified label.
func (m *Menu) NewMenu(label string) *Menu {
	return &Menu{
		menuer: m.menuer.NewMenu(label),
	}
}

// NewItem adds an item to the menu node with the specified label and event.
func (m *Menu) NewItem(label string, event event.Menu) {
	m.menuer.NewItem(label, event.Action)
}
