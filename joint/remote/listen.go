package remote

import (
	"log"

	"github.com/codeation/impress/joint/iface"
)

func (s *server) streamListen() {
	for {
		var command byte
		if err := s.streamPipe.Byte(&command).CallErr(); err != nil {
			log.Printf("readCommands: %v", err)
			return
		}

		switch command {
		case iface.ApplicationSizeCode:
			var x, y, width, height int
			s.streamPipe.
				Int(&x).Int(&y).Int(&width).Int(&height).
				Call()
			s.calls.ApplicationSize(x, y, width, height)

		case iface.ApplicationTitleCode:
			var title string
			s.streamPipe.
				String(&title).
				Call()
			s.calls.ApplicationTitle(title)

		case iface.FrameNewCode:
			var frameID, parentFrameID int
			var x, y, width, height int
			s.streamPipe.
				Int(&frameID).Int(&parentFrameID).
				Int(&x).Int(&y).Int(&width).Int(&height).
				Call()
			s.calls.FrameNew(frameID, parentFrameID, x, y, width, height)

		case iface.FrameDropCode:
			var frameID int
			s.streamPipe.
				Int(&frameID).
				Call()
			s.calls.FrameDrop(frameID)

		case iface.FrameSizeCode:
			var frameID int
			var x, y, width, height int
			s.streamPipe.
				Int(&frameID).
				Int(&x).Int(&y).Int(&width).Int(&height).
				Call()
			s.calls.FrameSize(frameID, x, y, width, height)

		case iface.FrameRaiseCode:
			var frameID int
			s.streamPipe.
				Int(&frameID).
				Call()
			s.calls.FrameRaise(frameID)

		case iface.WindowNewCode:
			var windowID, frameID int
			var x, y, width, height int
			s.streamPipe.
				Int(&windowID).Int(&frameID).
				Int(&x).Int(&y).Int(&width).Int(&height).
				Call()
			s.calls.WindowNew(windowID, frameID, x, y, width, height)

		case iface.WindowDropCode:
			var windowID int
			s.streamPipe.
				Int(&windowID).
				Call()
			s.calls.WindowDrop(windowID)

		case iface.WindowRaiseCode:
			var windowID int
			s.streamPipe.
				Int(&windowID).
				Call()
			s.calls.WindowRaise(windowID)

		case iface.WindowClearCode:
			var windowID int
			s.streamPipe.
				Int(&windowID).
				Call()
			s.calls.WindowClear(windowID)

		case iface.WindowShowCode:
			var windowID int
			s.streamPipe.
				Int(&windowID).
				Call()
			s.calls.WindowShow(windowID)

		case iface.WindowSizeCode:
			var windowID int
			var x, y, width, height int
			s.streamPipe.
				Int(&windowID).
				Int(&x).Int(&y).Int(&width).Int(&height).
				Call()
			s.calls.WindowSize(windowID, x, y, width, height)

		case iface.WindowFillCode:
			var windowID int
			var x, y, width, height int
			var r, g, b uint16
			s.streamPipe.
				Int(&windowID).
				Int(&x).Int(&y).Int(&width).Int(&height).
				UInt16(&r).UInt16(&g).UInt16(&b).
				Call()
			s.calls.WindowFill(windowID, x, y, width, height, r, g, b)

		case iface.WindowLineCode:
			var windowID int
			var x0, y0, x1, y1 int
			var r, g, b uint16
			s.streamPipe.
				Int(&windowID).
				Int(&x0).Int(&y0).Int(&x1).Int(&y1).
				UInt16(&r).UInt16(&g).UInt16(&b).
				Call()
			s.calls.WindowLine(windowID, x0, y0, x1, y1, r, g, b)

		case iface.WindowTextCode:
			var windowID int
			var x, y int
			var r, g, b uint16
			var fontID int
			var height int
			var text string
			s.streamPipe.
				Int(&windowID).
				Int(&x).Int(&y).
				UInt16(&r).UInt16(&g).UInt16(&b).
				Int(&fontID).
				Int(&height).
				String(&text).
				Call()
			s.calls.WindowText(windowID, x, y, r, g, b, fontID, height, text)

		case iface.WindowImageCode:
			var windowID int
			var x, y, width, height int
			var imageID int
			s.streamPipe.
				Int(&windowID).
				Int(&x).Int(&y).Int(&width).Int(&height).
				Int(&imageID).
				Call()
			s.calls.WindowImage(windowID, x, y, width, height, imageID)

		case iface.FontDropCode:
			var fontID int
			s.streamPipe.
				Int(&fontID).
				Call()
			s.calls.FontDrop(fontID)

		case iface.ImageNewCode:
			var imageID int
			var width, height int
			var bitmap []byte
			s.streamPipe.
				Int(&imageID).
				Int(&width).Int(&height).
				Bytes(&bitmap).
				Call()
			s.calls.ImageNew(imageID, width, height, bitmap)

		case iface.ImageDropCode:
			var imageID int
			s.streamPipe.
				Int(&imageID).
				Call()
			s.calls.ImageDrop(imageID)

		case iface.MenuNewCode:
			var menuID int
			var parentMenuID int
			var label string
			s.streamPipe.
				Int(&menuID).
				Int(&parentMenuID).
				String(&label).
				Call()
			s.calls.MenuNew(menuID, parentMenuID, label)

		case iface.MenuItemCode:
			var menuID int
			var parentMenuID int
			var label, action string
			s.streamPipe.
				Int(&menuID).
				Int(&parentMenuID).
				String(&label).String(&action).
				Call()
			s.calls.MenuItem(menuID, parentMenuID, label, action)

		default:
			log.Printf("unknown event: %d", command)
			return
		}
	}
}

func (s *server) syncListen() {
	for {
		var command byte
		if err := s.syncPipe.Byte(&command).CallErr(); err != nil {
			log.Printf("readRequests: %v", err)
			return
		}

		switch command {
		case iface.ApplicationExitCode:
			s.syncPipe.
				Call()
			output := s.calls.ApplicationExit()
			s.syncPipe.Call(output)
			s.syncPipe.Flush()

		case iface.ApplicationVersionCode:
			s.syncPipe.
				Call()
			version := s.calls.ApplicationExit()
			s.syncPipe.Call(version)
			s.syncPipe.Flush()

		case iface.FontNewCode:
			var fontID int
			var height int
			var style, variant, weight, stretch int
			var family string
			s.syncPipe.
				Int(&fontID).
				Int(&height).
				Int(&style).Int(&variant).Int(&weight).Int(&stretch).
				String(&family).
				Call()
			baseline, ascent, descent := s.calls.FontNew(fontID, height, style, variant, weight, stretch, family)
			s.syncPipe.Call(baseline, ascent, descent)
			s.syncPipe.Flush()

		case iface.FontSplitCode:
			var fontID int
			var edge, indent int
			var text string
			s.syncPipe.
				Int(&fontID).
				Int(&edge).Int(&indent).
				String(&text).
				Call()
			lengths := s.calls.FontSplit(fontID, text, edge, indent)
			s.syncPipe.Call(lengths)
			s.syncPipe.Flush()

		case iface.FontSizeCode:
			var fontID int
			var text string
			s.syncPipe.
				Int(&fontID).
				String(&text).
				Call()
			x, y := s.calls.FontSize(fontID, text)
			s.syncPipe.Call(x, y)
			s.syncPipe.Flush()

		default:
			log.Printf("unknown event: %d", command)
			return
		}
	}
}
