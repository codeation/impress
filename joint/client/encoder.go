package client

import (
	"github.com/codeation/impress/joint/iface"
)

func (c *Client) ApplicationSize(x, y, width, height int) {
	c.streamPipe.Call(iface.ApplicationSizeCode, x, y, width, height)
}

func (c *Client) ApplicationTitle(title string) {
	c.streamPipe.Call(iface.ApplicationTitleCode, title)
}

func (c *Client) ApplicationExit() string {
	c.onExit = true
	var output string
	c.syncPipe.
		String(&output).
		Call(iface.ApplicationExitCode)
	return output
}

func (c *Client) ApplicationVersion() string {
	var version string
	c.syncPipe.
		String(&version).
		Call(iface.ApplicationVersionCode)
	return version
}

func (c *Client) FrameNew(frameID int, parentFrameID int, x, y, width, height int) {
	c.streamPipe.Call(iface.FrameNewCode, frameID, parentFrameID, x, y, width, height)
}

func (c *Client) FrameDrop(frameID int) {
	c.streamPipe.Call(iface.FrameDropCode, frameID)
}

func (c *Client) FrameSize(frameID int, x, y, width, height int) {
	c.streamPipe.Call(iface.FrameSizeCode, frameID, x, y, width, height)
}

func (c *Client) FrameRaise(frameID int) {
	c.streamPipe.Call(iface.FrameRaiseCode, frameID)
}

func (c *Client) WindowNew(windowID int, frameID int, x, y, width, height int) {
	c.streamPipe.Call(iface.WindowNewCode, windowID, frameID, x, y, width, height)
}

func (c *Client) WindowDrop(windowID int) {
	c.streamPipe.Call(iface.WindowDropCode, windowID)
}

func (c *Client) WindowRaise(windowID int) {
	c.streamPipe.Call(iface.WindowRaiseCode, windowID)
}

func (c *Client) WindowClear(windowID int) {
	c.streamPipe.Call(iface.WindowClearCode, windowID)
}

func (c *Client) WindowShow(windowID int) {
	c.streamPipe.Call(iface.WindowShowCode, windowID)
}

func (c *Client) WindowSize(windowID int, x, y, width, height int) {
	c.streamPipe.Call(iface.WindowSizeCode, windowID, x, y, width, height)
}

func (c *Client) WindowFill(windowID int, x, y, width, height int, r, g, b uint16) {
	c.streamPipe.Call(iface.WindowFillCode, windowID, x, y, width, height, r, g, b)
}

func (c *Client) WindowLine(windowID int, x0, y0, x1, y1 int, r, g, b uint16) {
	c.streamPipe.Call(iface.WindowLineCode, windowID, x0, y0, x1, y1, r, g, b)
}

func (c *Client) WindowText(windowID int, x, y int, r, g, b uint16, fontID int, height int, text string) {
	c.streamPipe.Call(iface.WindowTextCode, windowID, x, y, r, g, b, fontID, height, text)
}

func (c *Client) WindowImage(windowID int, x, y, width, height int, imageID int) {
	c.streamPipe.Call(iface.WindowImageCode, windowID, x, y, width, height, imageID)
}

func (c *Client) FontNew(fontID int, height int, style, variant, weight, stretch int, family string) (int, int, int) {
	var baseline, ascent, descent int
	c.syncPipe.
		Int(&baseline).
		Int(&ascent).
		Int(&descent).
		Call(iface.FontNewCode, fontID, height, style, variant, weight, stretch, family)
	return baseline, ascent, descent
}

func (c *Client) FontDrop(fontID int) {}

func (c *Client) FontSplit(fontID int, text string, edge, indent int) []int {
	var lengths []int
	c.syncPipe.
		Ints(&lengths).
		Call(iface.FontSplitCode, fontID, edge, indent, text)
	return lengths
}

func (c *Client) FontSize(fontID int, text string) (int, int) {
	var x, y int
	c.syncPipe.
		Int(&x).Int(&y).
		Call(iface.FontSizeCode, fontID, text)
	return x, y
}

func (c *Client) ImageNew(imageID int, width, height int, bitmap []byte) {
	c.streamPipe.Call(iface.ImageNewCode, imageID, width, height, bitmap)
}

func (c *Client) ImageDrop(imageID int) {
	c.streamPipe.Call(iface.ImageDropCode, imageID)
}

func (c *Client) MenuNew(menuID int, parentMenuID int, label string) {
	c.streamPipe.Call(iface.MenuNewCode, menuID, parentMenuID, label)
}

func (c *Client) MenuItem(menuID int, parentMenuID int, label string, action string) {
	c.streamPipe.Call(iface.MenuItemCode, menuID, parentMenuID, label, action)
}
