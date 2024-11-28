// Package implements an internal mechanism to communicate with an impress terminal.
package drawwait

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/codeation/impress/joint/iface"
)

var ErrPipeClosing = errors.New("pipe is closing")

const defaultBufferSize = 64 * 1024

type decoder interface {
	Decode(reader io.Reader) (bool, error)
}

type streamCommand struct {
	calls  iface.CallSet
	reader *bufio.Reader
	queue  []decoder
	next   func()
}

func NewStreamCommand(calls iface.CallSet, reader io.Reader) *streamCommand {
	s := &streamCommand{
		calls:  calls,
		reader: bufio.NewReaderSize(reader, defaultBufferSize),
	}
	return s
}

func (s *streamCommand) StreamCommand() error {
	for {
		if err := s.streamCommand(); err != nil {
			return err
		}
		if s.reader.Buffered() == 0 {
			break
		}
	}
	return nil
}

func (s *streamCommand) streamCommand() error {
	if s.next == nil {
		return s.decodeCommand()
	}
	if len(s.queue) > 0 {
		ok, err := s.queue[0].Decode(s.reader)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		if ok {
			s.queue = s.queue[1:]
		}
	}
	if len(s.queue) == 0 {
		fn := s.next
		s.next = nil
		fn()
	}
	return nil
}

func (s *streamCommand) decodeCommand() error {
	var command byte
	if err := binary.Read(s.reader, binary.LittleEndian, &command); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	switch command {
	case iface.ApplicationExitCode:
		s.calls.ApplicationExit()
		return ErrPipeClosing

	case iface.ApplicationSizeCode:
		var x, y, width, height int
		s.queue = newParameters(&x, &y, &width, &height)
		s.next = func() {
			s.calls.ApplicationSize(x, y, width, height)
		}

	case iface.ApplicationTitleCode:
		var title string
		s.queue = newParameters(&title)
		s.next = func() {
			s.calls.ApplicationTitle(title)
		}

	case iface.FrameNewCode:
		var frameID, parentFrameID int
		var x, y, width, height int
		s.queue = newParameters(&frameID, &parentFrameID, &x, &y, &width, &height)
		s.next = func() {
			s.calls.FrameNew(frameID, parentFrameID, x, y, width, height)
		}

	case iface.FrameDropCode:
		var frameID int
		s.queue = newParameters(&frameID)
		s.next = func() {
			s.calls.FrameDrop(frameID)
		}

	case iface.FrameSizeCode:
		var frameID int
		var x, y, width, height int
		s.queue = newParameters(&frameID, &x, &y, &width, &height)
		s.next = func() {
			s.calls.FrameSize(frameID, x, y, width, height)
		}

	case iface.FrameRaiseCode:
		var frameID int
		s.queue = newParameters(&frameID)
		s.next = func() {
			s.calls.FrameRaise(frameID)
		}

	case iface.WindowNewCode:
		var windowID, frameID int
		var x, y, width, height int
		s.queue = newParameters(&windowID, &frameID, &x, &y, &width, &height)
		s.next = func() {
			s.calls.WindowNew(windowID, frameID, x, y, width, height)
		}

	case iface.WindowDropCode:
		var windowID int
		s.queue = newParameters(&windowID)
		s.next = func() {
			s.calls.WindowDrop(windowID)
		}

	case iface.WindowRaiseCode:
		var windowID int
		s.queue = newParameters(&windowID)
		s.next = func() {
			s.calls.WindowRaise(windowID)
		}

	case iface.WindowClearCode:
		var windowID int
		s.queue = newParameters(&windowID)
		s.next = func() {
			s.calls.WindowClear(windowID)
		}

	case iface.WindowShowCode:
		var windowID int
		s.queue = newParameters(&windowID)
		s.next = func() {
			s.calls.WindowShow(windowID)
		}

	case iface.WindowSizeCode:
		var windowID int
		var x, y, width, height int
		s.queue = newParameters(&windowID, &x, &y, &width, &height)
		s.next = func() {
			s.calls.WindowSize(windowID, x, y, width, height)
		}

	case iface.WindowFillCode:
		var windowID int
		var x, y, width, height int
		var r, g, b, a uint16
		s.queue = newParameters(&windowID, &x, &y, &width, &height, &r, &g, &b, &a)
		s.next = func() {
			s.calls.WindowFill(windowID, x, y, width, height, r, g, b, a)
		}

	case iface.WindowLineCode:
		var windowID int
		var x0, y0, x1, y1 int
		var r, g, b, a uint16
		s.queue = newParameters(&windowID, &x0, &y0, &x1, &y1, &r, &g, &b, &a)
		s.next = func() {
			s.calls.WindowLine(windowID, x0, y0, x1, y1, r, g, b, a)
		}

	case iface.WindowTextCode:
		var windowID int
		var x, y int
		var r, g, b, a uint16
		var fontID int
		var text string
		s.queue = newParameters(&windowID, &x, &y, &r, &g, &b, &a, &fontID, &text)
		s.next = func() {
			s.calls.WindowText(windowID, x, y, r, g, b, a, fontID, text)
		}

	case iface.WindowImageCode:
		var windowID int
		var x, y, width, height int
		var imageID int
		s.queue = newParameters(&windowID, &x, &y, &width, &height, &imageID)
		s.next = func() {
			s.calls.WindowImage(windowID, x, y, width, height, imageID)
		}

	case iface.FontDropCode:
		var fontID int
		s.queue = newParameters(&fontID)
		s.next = func() {
			s.calls.FontDrop(fontID)
		}

	case iface.ImageNewCode:
		var imageID int
		var width, height int
		var bitmap []byte
		s.queue = newParameters(&imageID, &width, &height, &bitmap)
		s.next = func() {
			s.calls.ImageNew(imageID, width, height, bitmap)
		}

	case iface.ImageDropCode:
		var imageID int
		s.queue = newParameters(&imageID)
		s.next = func() {
			s.calls.ImageDrop(imageID)
		}

	case iface.MenuNewCode:
		var menuID int
		var parentMenuID int
		var label string
		s.queue = newParameters(&menuID, &parentMenuID, &label)
		s.next = func() {
			s.calls.MenuNew(menuID, parentMenuID, label)
		}

	case iface.MenuItemCode:
		var menuID int
		var parentMenuID int
		var label, action string
		s.queue = newParameters(&menuID, &parentMenuID, &label, &action)
		s.next = func() {
			s.calls.MenuItem(menuID, parentMenuID, label, action)
		}

	case iface.ClipboardGetCode:
		var typeID int
		s.queue = newParameters(&typeID)
		s.next = func() {
			s.calls.ClipboardGet(typeID)
		}

	case iface.ClipboardPutCode:
		var typeID int
		var data []byte
		s.queue = newParameters(&typeID, &data)
		s.next = func() {
			s.calls.ClipboardPut(typeID, data)
		}

	default:
		return fmt.Errorf("unknown event: %d", command)
	}
	return nil
}
