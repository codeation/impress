package duo

import (
	"github.com/codeation/impress/driver"
)

type menuitem struct {
	driver *duo
	ID     int
}

func (d *duo) NewMenu(label string) driver.Menuer {
	d.lastMenuID++
	d.streamPipe.Call(
		'E', d.lastMenuID, 0, label)
	return &menuitem{
		driver: d,
		ID:     d.lastMenuID,
	}
}

func (m *menuitem) NewMenu(label string) driver.Menuer {
	m.driver.lastMenuID++
	m.driver.streamPipe.Call(
		'E', m.driver.lastMenuID, m.ID, label)
	return &menuitem{
		driver: m.driver,
		ID:     m.driver.lastMenuID,
	}
}

func (m *menuitem) NewItem(label string, action string) {
	m.driver.lastMenuID++
	m.driver.streamPipe.Call(
		'G', m.driver.lastMenuID, m.ID, label, action)
}
