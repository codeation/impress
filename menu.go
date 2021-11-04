package impress

// Menu represents any menu node
type Menu struct {
	menuer Menuer
	app    *Application
}

// NewMenu returns a new top-level menu node
func (app *Application) NewMenu(label string) *Menu {
	return &Menu{
		menuer: driver.NewMenu(label),
		app:    app,
	}
}

// NewMenu returns new submenu node
func (m *Menu) NewMenu(label string) *Menu {
	return &Menu{
		menuer: m.menuer.NewMenu(label),
		app:    m.app,
	}
}

// NewItem adds a item to menu node
func (m *Menu) NewItem(label string, event MenuEvent) {
	m.menuer.NewItem(label, event.Action)
}
