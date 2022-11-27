package duo

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
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
	fileDraw     *os.File
	fileAnswer   *os.File
	fileEvent    *os.File
	fileSuffix   string
	onExit       bool
	drawPipe     *pipe
	eventPipe    *pipe
}

func init() {
	impress.Register(newDuo())
}

func newDuo() *duo {
	d := new(duo)
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
	d.events = make(chan event.Eventer, 1024)
	go d.readEvents()
	return d
}

func (d *duo) Init() {
	// Version test
	var version string
	d.drawPipe.String(&version).Call(
		'V')
	if version != it_version {
		log.Fatalf("./it version \"%s\", expected \"%s\"", version, it_version)
	}
}

func (d *duo) Done() {
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

func (d *duo) Size(rect image.Rectangle) {
	x, y, width, height := rectangle(rect)
	d.drawPipe.Call(
		'S', x, y, width, height)
}

func (d *duo) Title(title string) {
	d.drawPipe.Call(
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
			log.Fatal(err)
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
