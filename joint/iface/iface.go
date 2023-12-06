// Package implements an internal mechanism to communicate with an impress terminal.
package iface

type CallSet interface {
	ApplicationExit() string
	ApplicationVersion() string
	ApplicationSize(x, y, width, height int)
	ApplicationTitle(title string)
	FrameNew(frameID int, parentFrameID int, x, y, width, height int)
	FrameDrop(frameID int)
	FrameSize(frameID int, x, y, width, height int)
	FrameRaise(frameID int)
	WindowNew(windowID int, frameID int, x, y, width, height int)
	WindowDrop(windowID int)
	WindowRaise(windowID int)
	WindowClear(windowID int)
	WindowShow(windowID int)
	WindowSize(windowID int, x, y, width, height int)
	WindowFill(windowID int, x, y, width, height int, r, g, b, a uint16)
	WindowLine(windowID int, x0, y0, x1, y1 int, r, g, b, a uint16)
	WindowText(windowID int, x, y int, r, g, b, a uint16, fontID int, text string)
	WindowImage(windowID int, x, y, width, height int, imageID int)
	FontNew(fontID int, height int, style, variant, weight, stretch int, family string) (int, int, int, int)
	FontDrop(fontID int)
	FontSplit(fontID int, text string, edge, indent int) []int
	FontSize(fontID int, text string) (int, int)
	ImageNew(imageID int, width, height int, bitmap []byte)
	ImageDrop(imageID int)
	MenuNew(menuID int, parentMenuID int, label string)
	MenuItem(menuID int, parentMenuID int, label string, action string)
	ClipboardGet(typeID int)
	ClipboardPut(typeID int, data []byte)
}

type CallbackSet interface {
	EventGeneral(eventID int)
	EventKeyboard(r rune, shift, control, alt, meta bool, name string)
	EventConfigure(width, height, innerWidth, innerHeight int)
	EventButton(action, button int, x, y int)
	EventMotion(x, y int, shift, control, alt, meta bool)
	EventMenu(action string)
	EventScroll(direction int, deltaX, deltaY int)
	EventClipboard(typeID int, data []byte)
}
