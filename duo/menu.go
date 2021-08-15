package duo

import (
	"github.com/codeation/impress"
)

type menuNode struct {
	driver *driver
	ID     int
}

func (d *driver) NewMenu(label string) impress.Menuer {
	d.onDraw.Lock()
	defer d.onDraw.Unlock()
	d.lastMenuID++
	writeSequence(d.pipeDraw, 'E', d.lastMenuID, 0, label)
	return &menuNode{
		driver: d,
		ID:     d.lastMenuID,
	}
}

func (node *menuNode) NewMenu(label string) impress.Menuer {
	node.driver.onDraw.Lock()
	defer node.driver.onDraw.Unlock()
	node.driver.lastMenuID++
	writeSequence(node.driver.pipeDraw, 'E', node.driver.lastMenuID, node.ID, label)
	return &menuNode{
		driver: node.driver,
		ID:     node.driver.lastMenuID,
	}
}

func (node *menuNode) NewItem(label string, action string) {
	node.driver.onDraw.Lock()
	defer node.driver.onDraw.Unlock()
	node.driver.lastMenuID++
	writeSequence(node.driver.pipeDraw, 'G', node.driver.lastMenuID, node.ID, label, action)
}
