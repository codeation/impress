// Package implements an internal mechanism to communicate with an impress terminal.
package drawsend

import (
	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type drawSend struct {
	streamPipe *rpc.Pipe
	syncPipe   *rpc.Pipe
}

func New(streamPipe, syncPipe *rpc.Pipe) *drawSend {
	c := &drawSend{
		streamPipe: streamPipe,
		syncPipe:   syncPipe,
	}
	return c
}

func (c *drawSend) ApplicationSize(x, y, width, height int) {
	c.streamPipe.
		Lock().
		Put(iface.ApplicationSizeCode, x, y, width, height).
		Unlock()
}

func (c *drawSend) ApplicationTitle(title string) {
	c.streamPipe.
		Lock().
		Put(iface.ApplicationTitleCode, title).
		Unlock()
}

func (c *drawSend) ApplicationExit() {
	c.streamPipe.
		Lock().
		Put(iface.ApplicationExitCode).
		Flush().
		Unlock()
	c.syncPipe.
		Lock().
		Put(iface.ApplicationExitCode).
		Flush().
		Unlock()
}

func (c *drawSend) ApplicationVersion() string {
	var version string
	c.syncPipe.
		Lock().
		Put(iface.ApplicationVersionCode).
		Flush().
		Get(&version).
		Unlock()
	return version
}

func (c *drawSend) FrameNew(frameID int, parentFrameID int, x, y, width, height int) {
	c.streamPipe.
		Lock().
		Put(iface.FrameNewCode, frameID, parentFrameID, x, y, width, height).
		Unlock()
}

func (c *drawSend) FrameDrop(frameID int) {
	c.streamPipe.
		Lock().
		Put(iface.FrameDropCode, frameID).
		Unlock()
}

func (c *drawSend) FrameSize(frameID int, x, y, width, height int) {
	c.streamPipe.
		Lock().
		Put(iface.FrameSizeCode, frameID, x, y, width, height).
		Unlock()
}

func (c *drawSend) FrameRaise(frameID int) {
	c.streamPipe.
		Lock().
		Put(iface.FrameRaiseCode, frameID).
		Unlock()
}

func (c *drawSend) WindowNew(windowID int, frameID int, x, y, width, height int) {
	c.streamPipe.
		Lock().
		Put(iface.WindowNewCode, windowID, frameID, x, y, width, height).
		Unlock()
}

func (c *drawSend) WindowDrop(windowID int) {
	c.streamPipe.
		Lock().
		Put(iface.WindowDropCode, windowID).
		Unlock()
}

func (c *drawSend) WindowRaise(windowID int) {
	c.streamPipe.
		Lock().
		Put(iface.WindowRaiseCode, windowID).
		Unlock()
}

func (c *drawSend) WindowClear(windowID int) {
	c.streamPipe.
		Lock().
		Put(iface.WindowClearCode, windowID).
		Unlock()
}

func (c *drawSend) WindowShow(windowID int) {
	c.streamPipe.
		Lock().
		Put(iface.WindowShowCode, windowID).
		Unlock()
}

func (c *drawSend) WindowSize(windowID int, x, y, width, height int) {
	c.streamPipe.
		Lock().
		Put(iface.WindowSizeCode, windowID, x, y, width, height).
		Unlock()
}

func (c *drawSend) WindowFill(windowID int, x, y, width, height int, r, g, b, a uint16) {
	c.streamPipe.
		Lock().
		Put(iface.WindowFillCode, windowID, x, y, width, height, r, g, b, a).
		Unlock()
}

func (c *drawSend) WindowLine(windowID int, x0, y0, x1, y1 int, r, g, b, a uint16) {
	c.streamPipe.
		Lock().
		Put(iface.WindowLineCode, windowID, x0, y0, x1, y1, r, g, b, a).
		Unlock()
}

func (c *drawSend) WindowText(windowID int, x, y int, r, g, b, a uint16, fontID int, text string) {
	c.streamPipe.
		Lock().
		Put(iface.WindowTextCode, windowID, x, y, r, g, b, a, fontID, text).
		Unlock()
}

func (c *drawSend) WindowImage(windowID int, x, y, width, height int, imageID int) {
	c.streamPipe.
		Lock().
		Put(iface.WindowImageCode, windowID, x, y, width, height, imageID).
		Unlock()
}

func (c *drawSend) FontNew(fontID int, height int, style, variant, weight, stretch int, family string) (int, int, int, int) {
	var lineheight, baseline, ascent, descent int
	c.syncPipe.
		Lock().
		Put(iface.FontNewCode, fontID, height, style, variant, weight, stretch, family).
		Flush().
		Get(&lineheight, &baseline, &ascent, &descent).
		Unlock()
	return lineheight, baseline, ascent, descent
}

func (c *drawSend) FontDrop(fontID int) {
	c.streamPipe.
		Lock().
		Put(iface.FontDropCode, fontID).
		Unlock()
}

func (c *drawSend) FontSplit(fontID int, text string, edge, indent int) []int {
	var lengths []int
	c.syncPipe.
		Lock().
		Put(iface.FontSplitCode, fontID, edge, indent, text).
		Flush().
		Get(&lengths).
		Unlock()
	return lengths
}

func (c *drawSend) FontSize(fontID int, text string) (int, int) {
	var x, y int
	c.syncPipe.
		Lock().
		Put(iface.FontSizeCode, fontID, text).
		Flush().
		Get(&x, &y).
		Unlock()
	return x, y
}

func (c *drawSend) ImageNew(imageID int, width, height int, bitmap []byte) {
	c.streamPipe.
		Lock().
		Put(iface.ImageNewCode, imageID, width, height, bitmap).
		Flush().
		Unlock()
}

func (c *drawSend) ImageDrop(imageID int) {
	c.streamPipe.
		Lock().
		Put(iface.ImageDropCode, imageID).
		Unlock()
}

func (c *drawSend) MenuNew(menuID int, parentMenuID int, label string) {
	c.streamPipe.
		Lock().
		Put(iface.MenuNewCode, menuID, parentMenuID, label).
		Unlock()
}

func (c *drawSend) MenuItem(menuID int, parentMenuID int, label string, action string) {
	c.streamPipe.
		Lock().
		Put(iface.MenuItemCode, menuID, parentMenuID, label, action).
		Unlock()
}

func (c *drawSend) ClipboardGet(typeID int) {
	c.streamPipe.
		Lock().
		Put(iface.ClipboardGetCode, typeID).
		Unlock()
}

func (c *drawSend) ClipboardPut(typeID int, data []byte) {
	c.streamPipe.
		Lock().
		Put(iface.ClipboardPutCode, typeID, data).
		Flush().
		Unlock()
}
