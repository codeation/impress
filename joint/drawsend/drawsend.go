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
	c.streamPipe.Call(iface.ApplicationSizeCode, x, y, width, height)
}

func (c *drawSend) ApplicationTitle(title string) {
	c.streamPipe.Call(iface.ApplicationTitleCode, title)
}

func (c *drawSend) ApplicationExit() string {
	var output string
	c.syncPipe.
		String(&output).
		Call(iface.ApplicationExitCode)
	return output
}

func (c *drawSend) ApplicationVersion() string {
	var version string
	c.syncPipe.
		String(&version).
		Call(iface.ApplicationVersionCode)
	return version
}

func (c *drawSend) FrameNew(frameID int, parentFrameID int, x, y, width, height int) {
	c.streamPipe.Call(iface.FrameNewCode, frameID, parentFrameID, x, y, width, height)
}

func (c *drawSend) FrameDrop(frameID int) {
	c.streamPipe.Call(iface.FrameDropCode, frameID)
}

func (c *drawSend) FrameSize(frameID int, x, y, width, height int) {
	c.streamPipe.Call(iface.FrameSizeCode, frameID, x, y, width, height)
}

func (c *drawSend) FrameRaise(frameID int) {
	c.streamPipe.Call(iface.FrameRaiseCode, frameID)
}

func (c *drawSend) WindowNew(windowID int, frameID int, x, y, width, height int) {
	c.streamPipe.Call(iface.WindowNewCode, windowID, frameID, x, y, width, height)
}

func (c *drawSend) WindowDrop(windowID int) {
	c.streamPipe.Call(iface.WindowDropCode, windowID)
}

func (c *drawSend) WindowRaise(windowID int) {
	c.streamPipe.Call(iface.WindowRaiseCode, windowID)
}

func (c *drawSend) WindowClear(windowID int) {
	c.streamPipe.Call(iface.WindowClearCode, windowID)
}

func (c *drawSend) WindowShow(windowID int) {
	c.streamPipe.Call(iface.WindowShowCode, windowID)
}

func (c *drawSend) WindowSize(windowID int, x, y, width, height int) {
	c.streamPipe.Call(iface.WindowSizeCode, windowID, x, y, width, height)
}

func (c *drawSend) WindowFill(windowID int, x, y, width, height int, r, g, b uint16) {
	c.streamPipe.Call(iface.WindowFillCode, windowID, x, y, width, height, r, g, b)
}

func (c *drawSend) WindowLine(windowID int, x0, y0, x1, y1 int, r, g, b uint16) {
	c.streamPipe.Call(iface.WindowLineCode, windowID, x0, y0, x1, y1, r, g, b)
}

func (c *drawSend) WindowText(windowID int, x, y int, r, g, b uint16, fontID int, height int, text string) {
	c.streamPipe.Call(iface.WindowTextCode, windowID, x, y, r, g, b, fontID, height, text)
}

func (c *drawSend) WindowImage(windowID int, x, y, width, height int, imageID int) {
	c.streamPipe.Call(iface.WindowImageCode, windowID, x, y, width, height, imageID)
}

func (c *drawSend) FontNew(fontID int, height int, style, variant, weight, stretch int, family string) (int, int, int) {
	var baseline, ascent, descent int
	c.syncPipe.
		Int(&baseline).
		Int(&ascent).
		Int(&descent).
		Call(iface.FontNewCode, fontID, height, style, variant, weight, stretch, family)
	return baseline, ascent, descent
}

func (c *drawSend) FontDrop(fontID int) {}

func (c *drawSend) FontSplit(fontID int, text string, edge, indent int) []int {
	var lengths []int
	c.syncPipe.
		Ints(&lengths).
		Call(iface.FontSplitCode, fontID, edge, indent, text)
	return lengths
}

func (c *drawSend) FontSize(fontID int, text string) (int, int) {
	var x, y int
	c.syncPipe.
		Int(&x).Int(&y).
		Call(iface.FontSizeCode, fontID, text)
	return x, y
}

func (c *drawSend) ImageNew(imageID int, width, height int, bitmap []byte) {
	c.streamPipe.Call(iface.ImageNewCode, imageID, width, height, bitmap)
}

func (c *drawSend) ImageDrop(imageID int) {
	c.streamPipe.Call(iface.ImageDropCode, imageID)
}

func (c *drawSend) MenuNew(menuID int, parentMenuID int, label string) {
	c.streamPipe.Call(iface.MenuNewCode, menuID, parentMenuID, label)
}

func (c *drawSend) MenuItem(menuID int, parentMenuID int, label string, action string) {
	c.streamPipe.Call(iface.MenuItemCode, menuID, parentMenuID, label, action)
}
