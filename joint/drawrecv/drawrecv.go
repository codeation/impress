// Package implements an internal mechanism to communicate with an impress terminal.
package drawrecv

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

var ErrPipeClosing = errors.New("pipe is closing")

type drawRecv struct {
	calls      iface.CallSet
	streamPipe *rpc.Pipe
	syncPipe   *rpc.Pipe
	wg         sync.WaitGroup
}

func New(calls iface.CallSet, streamPipe, syncPipe *rpc.Pipe) *drawRecv {
	s := &drawRecv{
		calls:      calls,
		streamPipe: streamPipe,
		syncPipe:   syncPipe,
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			if err := s.streamCommand(); err != nil {
				if !errors.Is(err, ErrPipeClosing) {
					log.Printf("stream pipe error: %v", err)
				}
				return
			}
		}
	}()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			if err := s.syncCommand(); err != nil {
				if !errors.Is(err, ErrPipeClosing) {
					log.Printf("sync pipe error: %v", err)
				}
				return
			}
		}
	}()
	return s
}

func (s *drawRecv) Wait() {
	s.wg.Wait()
}

func (s *drawRecv) streamCommand() error {
	var command byte
	if err := s.streamPipe.Get(&command).Err(); err != nil {
		return fmt.Errorf("get (command): %w", err)
	}

	switch command {
	case iface.ApplicationExitCode:
		s.calls.ApplicationExit()
		return ErrPipeClosing

	case iface.ApplicationSizeCode:
		var x, y, width, height int
		if err := s.streamPipe.Get(&x, &y, &width, &height).Err(); err != nil {
			return err
		}
		s.calls.ApplicationSize(x, y, width, height)

	case iface.ApplicationTitleCode:
		var title string
		if err := s.streamPipe.Get(&title).Err(); err != nil {
			return err
		}
		s.calls.ApplicationTitle(title)

	case iface.FrameNewCode:
		var frameID, parentFrameID int
		var x, y, width, height int
		if err := s.streamPipe.Get(&frameID, &parentFrameID, &x, &y, &width, &height).Err(); err != nil {
			return err
		}
		s.calls.FrameNew(frameID, parentFrameID, x, y, width, height)

	case iface.FrameDropCode:
		var frameID int
		s.streamPipe.Get(&frameID)
		s.calls.FrameDrop(frameID)

	case iface.FrameSizeCode:
		var frameID int
		var x, y, width, height int
		if err := s.streamPipe.Get(&frameID, &x, &y, &width, &height).Err(); err != nil {
			return err
		}
		s.calls.FrameSize(frameID, x, y, width, height)

	case iface.FrameRaiseCode:
		var frameID int
		if err := s.streamPipe.Get(&frameID).Err(); err != nil {
			return err
		}
		s.calls.FrameRaise(frameID)

	case iface.WindowNewCode:
		var windowID, frameID int
		var x, y, width, height int
		if err := s.streamPipe.Get(&windowID, &frameID, &x, &y, &width, &height).Err(); err != nil {
			return err
		}
		s.calls.WindowNew(windowID, frameID, x, y, width, height)

	case iface.WindowDropCode:
		var windowID int
		if err := s.streamPipe.Get(&windowID).Err(); err != nil {
			return err
		}
		s.calls.WindowDrop(windowID)

	case iface.WindowRaiseCode:
		var windowID int
		if err := s.streamPipe.Get(&windowID).Err(); err != nil {
			return err
		}
		s.calls.WindowRaise(windowID)

	case iface.WindowClearCode:
		var windowID int
		if err := s.streamPipe.Get(&windowID).Err(); err != nil {
			return err
		}
		s.calls.WindowClear(windowID)

	case iface.WindowShowCode:
		var windowID int
		if err := s.streamPipe.Get(&windowID).Err(); err != nil {
			return err
		}
		s.calls.WindowShow(windowID)

	case iface.WindowSizeCode:
		var windowID int
		var x, y, width, height int
		if err := s.streamPipe.Get(&windowID, &x, &y, &width, &height).Err(); err != nil {
			return err
		}
		s.calls.WindowSize(windowID, x, y, width, height)

	case iface.WindowFillCode:
		var windowID int
		var x, y, width, height int
		var r, g, b, a uint16
		if err := s.streamPipe.Get(&windowID, &x, &y, &width, &height, &r, &g, &b, &a).Err(); err != nil {
			return err
		}
		s.calls.WindowFill(windowID, x, y, width, height, r, g, b, a)

	case iface.WindowLineCode:
		var windowID int
		var x0, y0, x1, y1 int
		var r, g, b, a uint16
		if err := s.streamPipe.Get(&windowID, &x0, &y0, &x1, &y1, &r, &g, &b, &a).Err(); err != nil {
			return err
		}
		s.calls.WindowLine(windowID, x0, y0, x1, y1, r, g, b, a)

	case iface.WindowTextCode:
		var windowID int
		var x, y int
		var r, g, b, a uint16
		var fontID int
		var text string
		if err := s.streamPipe.Get(&windowID, &x, &y, &r, &g, &b, &a, &fontID, &text).Err(); err != nil {
			return err
		}
		s.calls.WindowText(windowID, x, y, r, g, b, a, fontID, text)

	case iface.WindowImageCode:
		var windowID int
		var x, y, width, height int
		var imageID int
		if err := s.streamPipe.Get(&windowID, &x, &y, &width, &height, &imageID).Err(); err != nil {
			return err
		}
		s.calls.WindowImage(windowID, x, y, width, height, imageID)

	case iface.FontDropCode:
		var fontID int
		if err := s.streamPipe.Get(&fontID).Err(); err != nil {
			return err
		}
		s.calls.FontDrop(fontID)

	case iface.ImageNewCode:
		var imageID int
		var width, height int
		var bitmap []byte
		if err := s.streamPipe.Get(&imageID, &width, &height, &bitmap).Err(); err != nil {
			return err
		}
		s.calls.ImageNew(imageID, width, height, bitmap)

	case iface.ImageDropCode:
		var imageID int
		if err := s.streamPipe.Get(&imageID).Err(); err != nil {
			return err
		}
		s.calls.ImageDrop(imageID)

	case iface.MenuNewCode:
		var menuID int
		var parentMenuID int
		var label string
		if err := s.streamPipe.Get(&menuID, &parentMenuID, &label).Err(); err != nil {
			return err
		}
		s.calls.MenuNew(menuID, parentMenuID, label)

	case iface.MenuItemCode:
		var menuID int
		var parentMenuID int
		var label, action string
		if err := s.streamPipe.Get(&menuID, &parentMenuID, &label, &action).Err(); err != nil {
			return err
		}
		s.calls.MenuItem(menuID, parentMenuID, label, action)

	case iface.ClipboardGetCode:
		var typeID int
		if err := s.streamPipe.Get(&typeID).Err(); err != nil {
			return err
		}
		s.calls.ClipboardGet(typeID)

	case iface.ClipboardPutCode:
		var typeID int
		var data []byte
		if err := s.streamPipe.Get(&typeID, data).Err(); err != nil {
			return err
		}
		s.calls.ClipboardPut(typeID, data)

	default:
		return fmt.Errorf("unknown event: %d", command)
	}

	return nil
}

func (s *drawRecv) syncCommand() error {
	var command byte
	if err := s.syncPipe.Get(&command).Err(); err != nil {
		return fmt.Errorf("get (command): %w", err)
	}

	switch command {
	case iface.ApplicationExitCode:
		s.calls.ApplicationExit()
		return ErrPipeClosing

	case iface.ApplicationVersionCode:
		version := s.calls.ApplicationVersion()
		if err := s.syncPipe.Put(version).Flush().Err(); err != nil {
			return err
		}

	case iface.FontNewCode:
		var fontID int
		var height int
		var style, variant, weight, stretch int
		var family string
		if err := s.syncPipe.Get(&fontID, &height, &style, &variant, &weight, &stretch, &family).Err(); err != nil {
			return err
		}
		lineheight, baseline, ascent, descent := s.calls.FontNew(fontID, height, style, variant, weight, stretch, family)
		if err := s.syncPipe.Put(lineheight, baseline, ascent, descent).Flush().Err(); err != nil {
			return err
		}

	case iface.FontSplitCode:
		var fontID int
		var edge, indent int
		var text string
		if err := s.syncPipe.Get(&fontID, &edge, &indent, &text).Err(); err != nil {
			return err
		}
		lengths := s.calls.FontSplit(fontID, text, edge, indent)
		if err := s.syncPipe.Put(lengths).Flush().Err(); err != nil {
			return err
		}

	case iface.FontSizeCode:
		var fontID int
		var text string
		if err := s.syncPipe.Get(&fontID, &text).Err(); err != nil {
			return err
		}
		x, y := s.calls.FontSize(fontID, text)
		if err := s.syncPipe.Put(x, y).Flush().Err(); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown event: %d", command)
	}

	return nil
}
