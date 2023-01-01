package duo

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"
)

const (
	fifoStreamPath = "/tmp/it_fifo_stream_"
	fifoInputPath  = "/tmp/it_fifo_input_"
	fifoOutputPath = "/tmp/it_fifo_output_"
	fifoEventPath  = "/tmp/it_fifo_event_"
)

type duo struct {
	cmd          *exec.Cmd
	lastWindowID int
	lastFontID   int
	lastImageID  int
	lastMenuID   int
	events       chan event.Eventer
	fileStream   *os.File
	fileRequest  *os.File
	fileResponse *os.File
	fileEvent    *os.File
	fileSuffix   string
	onExit       bool
	streamPipe   *pipe
	syncPipe     *pipe
	eventPipe    *pipe
}

func init() {
	d := new(duo)
	if err := d.start(); err != nil {
		log.Printf("start duo driver error: %v", err)
		return
	}
	if err := d.connect(); err != nil {
		log.Println("connect duo driver error: %w", err)
		if err := d.cmd.Process.Kill(); err != nil {
			log.Println("kill duo driver error", err)
		}
		return
	}
	impress.Register(d)
}

func (d *duo) start() error {
	randBuffer := make([]byte, 8)
	if _, err := rand.Reader.Read(randBuffer); err != nil {
		return fmt.Errorf("rand.Reader.Read: %w", err)
	}
	d.fileSuffix = hex.EncodeToString(randBuffer)
	for _, name := range []string{fifoInputPath, fifoStreamPath, fifoOutputPath, fifoEventPath} {
		if err := syscall.Mkfifo(name+d.fileSuffix, 0644); err != nil {
			return fmt.Errorf("syscall.Mkfifo: %w", err)
		}
	}
	path := os.Getenv("IMPRESS_TERMINAL_PATH")
	if path == "" {
		path = "./it"
	}
	cmd := exec.Command(path, d.fileSuffix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start: %w", err)
	}
	d.cmd = cmd
	return nil
}

func (d *duo) connect() error {
	var err error
	if d.fileResponse, err = os.OpenFile(fifoOutputPath+d.fileSuffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(p): %w", err)
	}
	if d.fileEvent, err = os.OpenFile(fifoEventPath+d.fileSuffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(e): %w", err)
	}
	if d.fileRequest, err = os.OpenFile(fifoInputPath+d.fileSuffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(q): %w", err)
	}
	if d.fileStream, err = os.OpenFile(fifoStreamPath+d.fileSuffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(s): %w", err)
	}
	d.streamPipe = newPipe(new(sync.Mutex), bufio.NewWriter(d.fileStream), nil)
	d.syncPipe = newPipe(new(sync.Mutex), bufio.NewWriter(d.fileRequest), bufio.NewReader(d.fileResponse))
	d.eventPipe = newPipe(new(dummyMutex), nil, bufio.NewReader(d.fileEvent))
	d.events = make(chan event.Eventer, 1024)
	go d.readEvents()
	return nil
}

func (d *duo) Init() {
	// Version test
	var version string
	d.syncPipe.String(&version).Call(
		'V')
	if version != it_version {
		log.Printf("./it version \"%s\", expected \"%s\"", version, it_version)
	}
}

func (d *duo) Done() {
	if err := d.done(); err != nil {
		log.Printf("done: %v", err)
	}
}

func (d *duo) done() error {
	d.onExit = true
	d.streamPipe.Call('X')
	d.streamPipe.Flush()
	if err := d.fileRequest.Close(); err != nil {
		return fmt.Errorf("fileRequest.Close: %w", err)
	}
	if err := d.fileStream.Close(); err != nil {
		return fmt.Errorf("fileStream.Close: %w", err)
	}
	if err := d.cmd.Wait(); err != nil {
		return fmt.Errorf("cmd.Wait: %w", err)
	}
	if err := d.fileResponse.Close(); err != nil {
		return fmt.Errorf("fileResponse.Close: %w", err)
	}
	if err := d.fileEvent.Close(); err != nil {
		return fmt.Errorf("fileEvent.Close: %w", err)
	}
	for _, name := range []string{fifoInputPath, fifoOutputPath, fifoEventPath} {
		if _, err := os.Stat(name + d.fileSuffix); err == nil || !errors.Is(err, os.ErrNotExist) {
			if err = os.Remove(name + d.fileSuffix); err != nil {
				return fmt.Errorf("os.Remove: %w", err)
			}
		}
	}
	return nil
}

func (d *duo) Sync() {
	d.streamPipe.Flush()
}

func (d *duo) Size(rect image.Rectangle) {
	x, y, width, height := rectangle(rect)
	d.streamPipe.Call(
		'S', x, y, width, height)
}

func (d *duo) Title(title string) {
	d.streamPipe.Call(
		'T', title)
}

func (d *duo) Chan() <-chan event.Eventer {
	return d.events
}

func (d *duo) readEvents() {
	for {
		var command byte
		if err := d.eventPipe.Byte(&command).CallErr(); err != nil {
			if d.onExit {
				close(d.events)
				return
			}
			log.Printf("readEvents: %v", err)
		}
		switch command {
		case 'g':
			var e event.General
			d.eventPipe.
				UInt32(&e.Event).
				Call()
			d.events <- e
		case 'k':
			var e event.Keyboard
			d.eventPipe.
				Rune(&e.Rune).
				Bool(&e.Shift).
				Bool(&e.Control).
				Bool(&e.Alt).
				Bool(&e.Meta).
				String(&e.Name).
				Call()
			d.events <- e
		case 'f':
			var e event.Configure
			d.eventPipe.
				Int16(&e.Size.X).
				Int16(&e.Size.Y).
				Int16(&e.InnerSize.X).
				Int16(&e.InnerSize.Y).
				Call()
			d.events <- e
		case 'b':
			var e event.Button
			d.eventPipe.
				Char(&e.Action).
				Char(&e.Button).
				Int16(&e.Point.X).
				Int16(&e.Point.Y).
				Call()
			d.events <- e
		case 'm':
			var e event.Motion
			d.eventPipe.
				Int16(&e.Point.X).
				Int16(&e.Point.Y).
				Bool(&e.Shift).
				Bool(&e.Control).
				Bool(&e.Alt).
				Bool(&e.Meta).
				Call()
			d.events <- e
		case 'u':
			var e event.Menu
			d.eventPipe.
				String(&e.Action).
				Call()
			d.events <- e
		case 's':
			var e event.Scroll
			d.eventPipe.
				Int16(&e.Direction).
				Int16(&e.DeltaX).
				Int16(&e.DeltaY).
				Call()
			d.events <- e
		default:
			d.events <- event.UnknownEvent
		}
	}
}
