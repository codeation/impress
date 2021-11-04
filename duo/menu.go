package duo

import (
	"github.com/codeation/impress"
)

type menuNode struct {
	driver *driver
	ID     int
}

func (d *driver) NewMenu(label string) impress.Menuer {
	d.lastMenuID++
	d.drawPipe.Call(
		'E', d.lastMenuID, 0, label)
	return &menuNode{
		driver: d,
		ID:     d.lastMenuID,
	}
}

func (node *menuNode) NewMenu(label string) impress.Menuer {
	node.driver.lastMenuID++
	node.driver.drawPipe.Call(
		'E', node.driver.lastMenuID, node.ID, label)
	return &menuNode{
		driver: node.driver,
		ID:     node.driver.lastMenuID,
	}
}

func (node *menuNode) NewItem(label string, action string) {
	node.driver.lastMenuID++
	node.driver.drawPipe.Call(
		'G', node.driver.lastMenuID, node.ID, label, action)
}
