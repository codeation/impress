package duo

import (
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
	onDraw       sync.Mutex
	pipeDraw     *os.File
	pipeAnswer   *os.File
	pipeEvent    *os.File
	pipeSuffix   string
	onExit       bool
}

func init() {
	d := &driver{}
	d.runrpc()
	d.events = make(chan impress.Eventer, 1024)
	impress.Register(d)
}

func (d *driver) runrpc() {
	randBuffer := make([]byte, 8)
	if _, err := rand.Reader.Read(randBuffer); err != nil {
		log.Fatal(err)
	}
	d.pipeSuffix = hex.EncodeToString(randBuffer)
	for _, name := range []string{fifoInputPath, fifoOutputPath, fifoEventPath} {
		if err := syscall.Mkfifo(name+d.pipeSuffix, 0644); err != nil {
			log.Fatal(err)
		}
	}
	path := os.Getenv("IMPRESS_TERMINAL_PATH")
	if path == "" {
		path = "./it"
	}
	d.cmd = exec.Command(path, d.pipeSuffix)
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr
	if err := d.cmd.Start(); err != nil {
		log.Fatal(err)
	}
	var err error
	if d.pipeAnswer, err = os.OpenFile(fifoOutputPath+d.pipeSuffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		log.Fatal(err)
	}
	if d.pipeEvent, err = os.OpenFile(fifoEventPath+d.pipeSuffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		log.Fatal(err)
	}
	if d.pipeDraw, err = os.OpenFile(fifoInputPath+d.pipeSuffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		log.Fatal(err)
	}
	go d.readEvents()
}

func (d *driver) Init() {
	// Version test
	d.onDraw.Lock()
	writeSequence(d.pipeDraw, 'V')
	version, _ := readString(d.pipeAnswer)
	d.onDraw.Unlock()
	if version != it_version {
		log.Fatalf("./it version \"%s\", expected \"%s\"", version, it_version)
	}
}

func (d *driver) Done() {
	d.onExit = true
	d.onDraw.Lock()
	writeSequence(d.pipeDraw, 'X')
	d.onDraw.Unlock()
	if err := d.cmd.Wait(); err != nil {
		log.Fatalf("Wait %s", err)
	}
	if err := d.pipeDraw.Close(); err != nil {
		log.Fatalf("Close(d) %s", err)
	}
	if err := d.pipeAnswer.Close(); err != nil {
		log.Fatalf("Close(a) %s", err)
	}
	if err := d.pipeEvent.Close(); err != nil {
		log.Fatalf("Close(e) %s", err)
	}
	for _, name := range []string{fifoInputPath, fifoOutputPath, fifoEventPath} {
		if _, err := os.Stat(name + d.pipeSuffix); err == nil || !errors.Is(err, os.ErrNotExist) {
			_ = os.Remove(name + d.pipeSuffix)
		}
	}
}

func (d *driver) Size(rect impress.Rect) {
	d.onDraw.Lock()
	defer d.onDraw.Unlock()
	writeSequence(d.pipeDraw, 'S', rect.X, rect.Y, rect.Width, rect.Height)
}

func (d *driver) Title(title string) {
	d.onDraw.Lock()
	defer d.onDraw.Unlock()
	writeSequence(d.pipeDraw, 'T', title)
}

func (d *driver) Chan() <-chan impress.Eventer {
	return d.events
}

// Event

func (d *driver) readEvents() {
	for {
		command, err := readChar(d.pipeEvent)
		if err != nil {
			if d.onExit {
				break
			}
			log.Fatalf("readEvents %s", err)
		}
		switch command {
		case 'k':
			u, _ := readUInt32(d.pipeEvent)
			shift, _ := readBool(d.pipeEvent)
			control, _ := readBool(d.pipeEvent)
			alt, _ := readBool(d.pipeEvent)
			meta, _ := readBool(d.pipeEvent)
			name, _ := readString(d.pipeEvent)
			d.events <- impress.KeyboardEvent{
				Rune:    rune(u),
				Name:    name,
				Shift:   shift,
				Control: control,
				Alt:     alt,
				Meta:    meta,
			}
		case 'g':
			u, _ := readUInt32(d.pipeEvent)
			d.events <- impress.GeneralEvent{
				Event: u,
			}
		case 'f':
			width, _ := readInt16(d.pipeEvent)
			height, _ := readInt16(d.pipeEvent)
			d.events <- impress.ConfigureEvent{
				Size: impress.NewSize(width, height),
			}
		case 'b':
			btype, _ := readChar(d.pipeEvent)
			button, _ := readChar(d.pipeEvent)
			x, _ := readInt16(d.pipeEvent)
			y, _ := readInt16(d.pipeEvent)
			d.events <- impress.ButtonEvent{
				Action: int(btype),
				Button: int(button),
				Point:  impress.NewPoint(x, y),
			}
		case 'm':
			x, _ := readInt16(d.pipeEvent)
			y, _ := readInt16(d.pipeEvent)
			shift, _ := readBool(d.pipeEvent)
			control, _ := readBool(d.pipeEvent)
			alt, _ := readBool(d.pipeEvent)
			meta, _ := readBool(d.pipeEvent)
			d.events <- impress.MotionEvent{
				Point:   impress.NewPoint(x, y),
				Shift:   shift,
				Control: control,
				Alt:     alt,
				Meta:    meta,
			}
		case 'u':
			name, _ := readString(d.pipeEvent)
			d.events <- impress.MenuEvent{
				Action: name,
			}
		default:
			d.events <- impress.UnknownEvent
		}
	}
}
