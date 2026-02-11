// Package implements an internal mechanism to communicate with an impress terminal.
package drawsend

import (
	"log"

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
	must(c.streamPipe.PutTx(iface.ApplicationSizeCode, x, y, width, height))
}

func (c *drawSend) ApplicationTitle(title string) {
	must(c.streamPipe.PutTx(iface.ApplicationTitleCode, title))
}

func (c *drawSend) ApplicationExit() {
	must(c.streamPipe.PutTx(iface.ApplicationExitCode))
	must(c.streamPipe.Sync())
	must(c.syncPipe.PutTx(iface.ApplicationExitCode))
	must(c.syncPipe.Sync())
}

func (c *drawSend) ApplicationVersion() string {
	var version string
	must(c.syncPipe.IO(
		[]any{iface.ApplicationVersionCode},
		[]any{&version},
	))
	return version
}

func (c *drawSend) FrameNew(frameID int, parentFrameID int, x, y, width, height int) {
	must(c.streamPipe.PutTx(iface.FrameNewCode, frameID, parentFrameID, x, y, width, height))
}

func (c *drawSend) FrameDrop(frameID int) {
	must(c.streamPipe.PutTx(iface.FrameDropCode, frameID))
}

func (c *drawSend) FrameSize(frameID int, x, y, width, height int) {
	must(c.streamPipe.PutTx(iface.FrameSizeCode, frameID, x, y, width, height))
}

func (c *drawSend) FrameRaise(frameID int) {
	must(c.streamPipe.PutTx(iface.FrameRaiseCode, frameID))
}

func (c *drawSend) WindowNew(windowID int, frameID int, x, y, width, height int) {
	must(c.streamPipe.PutTx(iface.WindowNewCode, windowID, frameID, x, y, width, height))
}

func (c *drawSend) WindowDrop(windowID int) {
	must(c.streamPipe.PutTx(iface.WindowDropCode, windowID))
}

func (c *drawSend) WindowRaise(windowID int) {
	must(c.streamPipe.PutTx(iface.WindowRaiseCode, windowID))
}

func (c *drawSend) WindowClear(windowID int) {
	must(c.streamPipe.PutTx(iface.WindowClearCode, windowID))
}

func (c *drawSend) WindowShow(windowID int) {
	must(c.streamPipe.PutTx(iface.WindowShowCode, windowID))
}

func (c *drawSend) WindowSize(windowID int, x, y, width, height int) {
	must(c.streamPipe.PutTx(iface.WindowSizeCode, windowID, x, y, width, height))
}

func (c *drawSend) WindowFill(windowID int, x, y, width, height int, r, g, b, a uint16) {
	must(c.streamPipe.PutTx(iface.WindowFillCode, windowID, x, y, width, height, r, g, b, a))
}

func (c *drawSend) WindowLine(windowID int, x0, y0, x1, y1 int, r, g, b, a uint16) {
	must(c.streamPipe.PutTx(iface.WindowLineCode, windowID, x0, y0, x1, y1, r, g, b, a))
}

func (c *drawSend) WindowText(windowID int, x, y int, r, g, b, a uint16, fontID int, text string) {
	must(c.streamPipe.PutTx(iface.WindowTextCode, windowID, x, y, r, g, b, a, fontID, text))
}

func (c *drawSend) WindowImage(windowID int, x, y, width, height int, imageID int) {
	must(c.streamPipe.PutTx(iface.WindowImageCode, windowID, x, y, width, height, imageID))
}

func (c *drawSend) FontNew(fontID int, height int, style, variant, weight, stretch int, family string) {
	must(c.streamPipe.PutTx(iface.FontNewCode, fontID, height, style, variant, weight, stretch, family))
}

func (c *drawSend) FontDrop(fontID int) {
	must(c.streamPipe.PutTx(iface.FontDropCode, fontID))
}

func (c *drawSend) FontMetricNew(fontID int, height int, style, variant, weight, stretch int, family string) (int, int, int, int) {
	var lineheight, baseline, ascent, descent int
	must(c.syncPipe.IO(
		[]any{iface.FontNewCode, fontID, height, style, variant, weight, stretch, family},
		[]any{&lineheight, &baseline, &ascent, &descent},
	))
	return lineheight, baseline, ascent, descent
}

func (c *drawSend) FontMetricDrop(fontID int) {
	must(c.syncPipe.IO([]any{iface.FontDropCode, fontID}, nil))
}

func (c *drawSend) FontMetricSplit(fontID int, text string, edge, indent int) []int {
	var lengths []int
	must(c.syncPipe.IO(
		[]any{iface.FontSplitCode, fontID, edge, indent, text},
		[]any{&lengths},
	))
	return lengths
}

func (c *drawSend) FontMetricSize(fontID int, text string) (int, int) {
	var x, y int
	must(c.syncPipe.IO(
		[]any{iface.FontSizeCode, fontID, text},
		[]any{&x, &y},
	))
	return x, y
}

func (c *drawSend) ImageNew(imageID int, width, height int, bitmap []byte) {
	must(c.streamPipe.PutTx(iface.ImageNewCode, imageID, width, height, bitmap))
}

func (c *drawSend) ImageDrop(imageID int) {
	must(c.streamPipe.PutTx(iface.ImageDropCode, imageID))
}

func (c *drawSend) MenuNew(menuID int, parentMenuID int, label string) {
	must(c.streamPipe.PutTx(iface.MenuNewCode, menuID, parentMenuID, label))
}

func (c *drawSend) MenuItem(menuID int, parentMenuID int, label string, action string) {
	must(c.streamPipe.PutTx(iface.MenuItemCode, menuID, parentMenuID, label, action))
}

func (c *drawSend) ClipboardGet(typeID int) {
	must(c.streamPipe.PutTx(iface.ClipboardGetCode, typeID))
}

func (c *drawSend) ClipboardPut(typeID int, data []byte) {
	must(c.streamPipe.PutTx(iface.ClipboardPutCode, typeID, data))
}

func must(err error) {
	if err != nil {
		log.Fatalf("drawsend.must: %v", err)
	}
}
