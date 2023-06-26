// Package implements an internal mechanism to communicate with an impress terminal.
package drawrecv

import (
	"log"
	"sync"

	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type drawRecv struct {
	calls      iface.CallSet
	streamPipe *rpc.Pipe
	syncPipe   *rpc.Pipe
	wg         sync.WaitGroup
	onExit     bool
}

func New(calls iface.CallSet, streamPipe, syncPipe *rpc.Pipe) *drawRecv {
	s := &drawRecv{
		calls:      calls,
		streamPipe: streamPipe,
		syncPipe:   syncPipe,
	}
	go s.streamListen()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.syncListen()
	}()
	return s
}

func (s *drawRecv) Wait() {
	s.wg.Wait()
}

func (s *drawRecv) streamListen() {
	for {
		var command byte
		if err := s.streamPipe.Get(&command).Err(); err != nil {
			if s.onExit {
				return
			}
			log.Printf("readCommands: %v", err)
			return
		}

		switch command {
		case iface.ApplicationSizeCode:
			var x, y, width, height int
			s.streamPipe.Get(&x, &y, &width, &height)
			s.calls.ApplicationSize(x, y, width, height)

		case iface.ApplicationTitleCode:
			var title string
			s.streamPipe.Get(&title)
			s.calls.ApplicationTitle(title)

		case iface.FrameNewCode:
			var frameID, parentFrameID int
			var x, y, width, height int
			s.streamPipe.Get(&frameID, &parentFrameID, &x, &y, &width, &height)
			s.calls.FrameNew(frameID, parentFrameID, x, y, width, height)

		case iface.FrameDropCode:
			var frameID int
			s.streamPipe.Get(&frameID)
			s.calls.FrameDrop(frameID)

		case iface.FrameSizeCode:
			var frameID int
			var x, y, width, height int
			s.streamPipe.Get(&frameID, &x, &y, &width, &height)
			s.calls.FrameSize(frameID, x, y, width, height)

		case iface.FrameRaiseCode:
			var frameID int
			s.streamPipe.Get(&frameID)
			s.calls.FrameRaise(frameID)

		case iface.WindowNewCode:
			var windowID, frameID int
			var x, y, width, height int
			s.streamPipe.Get(&windowID, &frameID, &x, &y, &width, &height)
			s.calls.WindowNew(windowID, frameID, x, y, width, height)

		case iface.WindowDropCode:
			var windowID int
			s.streamPipe.Get(&windowID)
			s.calls.WindowDrop(windowID)

		case iface.WindowRaiseCode:
			var windowID int
			s.streamPipe.Get(&windowID)
			s.calls.WindowRaise(windowID)

		case iface.WindowClearCode:
			var windowID int
			s.streamPipe.Get(&windowID)
			s.calls.WindowClear(windowID)

		case iface.WindowShowCode:
			var windowID int
			s.streamPipe.Get(&windowID)
			s.calls.WindowShow(windowID)

		case iface.WindowSizeCode:
			var windowID int
			var x, y, width, height int
			s.streamPipe.Get(&windowID, &x, &y, &width, &height)
			s.calls.WindowSize(windowID, x, y, width, height)

		case iface.WindowFillCode:
			var windowID int
			var x, y, width, height int
			var r, g, b, a uint16
			s.streamPipe.Get(&windowID, &x, &y, &width, &height, &r, &g, &b, &a)
			s.calls.WindowFill(windowID, x, y, width, height, r, g, b, a)

		case iface.WindowLineCode:
			var windowID int
			var x0, y0, x1, y1 int
			var r, g, b, a uint16
			s.streamPipe.Get(&windowID, &x0, &y0, &x1, &y1, &r, &g, &b, &a)
			s.calls.WindowLine(windowID, x0, y0, x1, y1, r, g, b, a)

		case iface.WindowTextCode:
			var windowID int
			var x, y int
			var r, g, b, a uint16
			var fontID int
			var height int
			var text string
			s.streamPipe.Get(&windowID, &x, &y, &r, &g, &b, &a, &fontID, &height, &text)
			s.calls.WindowText(windowID, x, y, r, g, b, a, fontID, height, text)

		case iface.WindowImageCode:
			var windowID int
			var x, y, width, height int
			var imageID int
			s.streamPipe.Get(&windowID, &x, &y, &width, &height, &imageID)
			s.calls.WindowImage(windowID, x, y, width, height, imageID)

		case iface.FontDropCode:
			var fontID int
			s.streamPipe.Get(&fontID)
			s.calls.FontDrop(fontID)

		case iface.ImageNewCode:
			var imageID int
			var width, height int
			var bitmap []byte
			s.streamPipe.Get(&imageID, &width, &height, &bitmap)
			s.calls.ImageNew(imageID, width, height, bitmap)

		case iface.ImageDropCode:
			var imageID int
			s.streamPipe.Get(&imageID)
			s.calls.ImageDrop(imageID)

		case iface.MenuNewCode:
			var menuID int
			var parentMenuID int
			var label string
			s.streamPipe.Get(&menuID, &parentMenuID, &label)
			s.calls.MenuNew(menuID, parentMenuID, label)

		case iface.MenuItemCode:
			var menuID int
			var parentMenuID int
			var label, action string
			s.streamPipe.Get(&menuID, &parentMenuID, &label, &action)
			s.calls.MenuItem(menuID, parentMenuID, label, action)

		default:
			log.Printf("unknown event: %d", command)
			return
		}
	}
}

func (s *drawRecv) syncListen() {
	for {
		var command byte
		if err := s.syncPipe.Get(&command).Err(); err != nil {
			log.Printf("readRequests: %v", err)
			return
		}

		switch command {
		case iface.ApplicationExitCode:
			output := s.calls.ApplicationExit()
			s.syncPipe.Put(output).Flush()
			s.onExit = true
			return

		case iface.ApplicationVersionCode:
			version := s.calls.ApplicationVersion()
			s.syncPipe.Put(version).Flush()

		case iface.FontNewCode:
			var fontID int
			var height int
			var style, variant, weight, stretch int
			var family string
			s.syncPipe.Get(&fontID, &height, &style, &variant, &weight, &stretch, &family)
			lineheight, baseline, ascent, descent := s.calls.FontNew(fontID, height, style, variant, weight, stretch, family)
			s.syncPipe.Put(lineheight, baseline, ascent, descent).Flush()

		case iface.FontSplitCode:
			var fontID int
			var edge, indent int
			var text string
			s.syncPipe.Get(&fontID, &edge, &indent, &text)
			lengths := s.calls.FontSplit(fontID, text, edge, indent)
			s.syncPipe.Put(lengths).Flush()

		case iface.FontSizeCode:
			var fontID int
			var text string
			s.syncPipe.Get(&fontID, &text)
			x, y := s.calls.FontSize(fontID, text)
			s.syncPipe.Put(x, y).Flush()

		default:
			log.Printf("unknown event: %d", command)
			return
		}
	}
}
