package duo

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/codeation/impress"
)

const (
	fifoInputPath  = "/tmp/it_fifo_input_"
	fifoOutputPath = "/tmp/it_fifo_output_"
	fifoEventPath  = "/tmp/it_fifo_event_"
)

// Driver

type driver struct {
	cmd          *exec.Cmd
	lastWindowID int
	lastFontID   int
	lastImageID  int
	lastMenuID   int
	events       chan impress.Eventer
	fileDraw     *os.File
	fileAnswer   *os.File
	fileEvent    *os.File
	fileSuffix   string
	onExit       bool
	drawPipe     *pipe
	eventPipe    *pipe
}

func init() {
	d := &driver{}
	d.runrpc()
	impress.Register(d)
}

func (d *driver) runrpc() {
	randBuffer := make([]byte, 8)
	if _, err := rand.Reader.Read(randBuffer); err != nil {
		log.Fatal(err)
	}
	d.fileSuffix = hex.EncodeToString(randBuffer)
	for _, name := range []string{fifoInputPath, fifoOutputPath, fifoEventPath} {
		if err := syscall.Mkfifo(name+d.fileSuffix, 0644); err != nil {
			log.Fatal(err)
		}
	}
	path := os.Getenv("IMPRESS_TERMINAL_PATH")
	if path == "" {
		path = "./it"
	}
	d.cmd = exec.Command(path, d.fileSuffix)
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr
	if err := d.cmd.Start(); err != nil {
		log.Fatal(err)
	}
	var err error
	if d.fileAnswer, err = os.OpenFile(fifoOutputPath+d.fileSuffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		log.Fatal(err)
	}
	if d.fileEvent, err = os.OpenFile(fifoEventPath+d.fileSuffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		log.Fatal(err)
	}
	if d.fileDraw, err = os.OpenFile(fifoInputPath+d.fileSuffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		log.Fatal(err)
	}
	d.drawPipe = newPipe(new(sync.Mutex), bufio.NewWriter(d.fileDraw), bufio.NewReader(d.fileAnswer))
	d.eventPipe = newPipe(new(dummyMutex), nil, bufio.NewReader(d.fileEvent))
	d.events = make(chan impress.Eventer, 1024)
	go d.readEvents()
}

func (d *driver) Init() {
	// Version test
	var version string
	d.drawPipe.String(&version).Call(
		'V')
	if version != it_version {
		log.Fatalf("./it version \"%s\", expected \"%s\"", version, it_version)
	}
}

func (d *driver) Done() {
	d.onExit = true
	d.drawPipe.Call('X')
	d.drawPipe.Flush()
	if err := d.fileDraw.Close(); err != nil {
		log.Fatalf("Close(d) %s", err)
	}
	if err := d.cmd.Wait(); err != nil {
		log.Fatalf("Wait %s", err)
	}
	if err := d.fileAnswer.Close(); err != nil {
		log.Fatalf("Close(a) %s", err)
	}
	if err := d.fileEvent.Close(); err != nil {
		log.Fatalf("Close(e) %s", err)
	}
	for _, name := range []string{fifoInputPath, fifoOutputPath, fifoEventPath} {
		if _, err := os.Stat(name + d.fileSuffix); err == nil || !errors.Is(err, os.ErrNotExist) {
			_ = os.Remove(name + d.fileSuffix)
		}
	}
}

func (d *driver) Size(rect impress.Rect) {
	d.drawPipe.Call(
		'S', rect.X, rect.Y, rect.Width, rect.Height)
}

func (d *driver) Title(title string) {
	d.drawPipe.Call(
		'T', title)
}

func (d *driver) Chan() <-chan impress.Eventer {
	return d.events
}

// Event

func (d *driver) readEvents() {
	for {
		var command byte
		if err := d.eventPipe.Byte(&command).CallErr(); err != nil {
			if d.onExit {
				close(d.events)
				return
			}
			log.Fatal(err)
		}
		switch command {
		case 'g':
			var e impress.GeneralEvent
			d.eventPipe.
				UInt32(&e.Event).
				Call()
			d.events <- e
		case 'k':
			var e impress.KeyboardEvent
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
			var e impress.ConfigureEvent
			d.eventPipe.
				Int16(&e.Size.Width).
				Int16(&e.Size.Height).
				Call()
			d.events <- e
		case 'b':
			var e impress.ButtonEvent
			d.eventPipe.
				Char(&e.Action).
				Char(&e.Button).
				Int16(&e.Point.X).
				Int16(&e.Point.Y).
				Call()
			d.events <- e
		case 'm':
			var e impress.MotionEvent
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
			var e impress.MenuEvent
			d.eventPipe.
				String(&e.Action).
				Call()
			d.events <- e
		default:
			d.events <- impress.UnknownEvent
		}
	}
}
