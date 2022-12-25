package driver

import (
	"image"
	"image/color"

	"github.com/codeation/impress/event"
)

// Driver is a internal interface to a application level functions
type Driver interface {
	Init()
	Done()
	Title(title string)
	Size(rect image.Rectangle)
	NewWindow(rect image.Rectangle, background color.Color) Painter
	NewFont(height int, attributes map[string]string) Fonter
	NewImage(img image.Image) Imager
	NewMenu(label string) Menuer
	Chan() <-chan event.Eventer
	Sync()
}

// Painter is a internal interface to a window functions
type Painter interface {
	Drop()
	Size(rect image.Rectangle)
	Raise()
	Clear()
	Show()
	Fill(rect image.Rectangle, foreground color.Color)
	Line(from image.Point, to image.Point, foreground color.Color)
	Image(rect image.Rectangle, img Imager)
	Text(text string, font Fonter, from image.Point, foreground color.Color)
}

// Fonter is a internal interface to a font functions
type Fonter interface {
	Baseline() int
	Ascent() int
	Descent() int
	Close()
	Split(text string, edge int) []string
	Size(text string) image.Point
}

// Imager is a internal interface to a image functions
type Imager interface {
	Width() int
	Height() int
	Close()
}

// Menuer is a internal interface to a menu node functions
type Menuer interface {
	NewMenu(label string) Menuer
	NewItem(label string, action string)
}
