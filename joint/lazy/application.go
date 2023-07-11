// Package implements an internal mechanism to communicate with an impress terminal.
package lazy

import (
	"image"

	"github.com/codeation/impress/driver"
)

type app struct {
	driver.Driver
	rect  image.Rectangle
	title string
}

func New(d driver.Driver) *app {
	return &app{
		Driver: d,
	}
}

func (a *app) Size(rect image.Rectangle) {
	if a.rect == rect {
		return
	}
	a.rect = rect
	a.Driver.Size(rect)
}

func (a *app) Title(title string) {
	if a.title == title {
		return
	}
	a.title = title
	a.Driver.Title(title)
}
