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

// NewItemFunc adds a item with callback func to menu node
func (m *Menu) NewItemFunc(label string, event MenuEvent, f func()) {
	m.menuer.NewItem(label, event.Action)
	m.app.OnEvent(event, f)
}
