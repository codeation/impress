package duo

import (
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/codeation/impress"
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
	connDraw     net.Conn
	onExit       bool
	connEvent    net.Conn
}

func init() {
	d := &driver{}
	d.runrpc()
	d.events = make(chan impress.Eventer, 1024)
	impress.Register(d)
}

func (d *driver) runrpc() {
	path := os.Getenv("IMPRESS_TERMINAL_PATH")
	if path == "" {
		path = "./it"
	}
	d.cmd = exec.Command(path)
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr
	if err := d.cmd.Start(); err != nil {
		log.Fatal(err)
	}
}

func (d *driver) Init() {
	// Wait connection to duo driver
	var err error
	for i := 0; i < 100; i++ {
		if d.connDraw, err = net.Dial("tcp", "localhost:1101"); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		log.Fatal(err)
	}
	// Wait connection to event socket
	for i := 0; i < 100; i++ {
		if d.connEvent, err = net.Dial("tcp", "localhost:1102"); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		log.Fatal(err)
	}
	go d.readEvents()
	// Version test
	d.onDraw.Lock()
	writeSequence(d.connDraw, 'V')
	version, _ := readString(d.connDraw)
	d.onDraw.Unlock()
	if version != it_version {
		log.Fatalf("./it version \"%s\", expected \"%s\"", version, it_version)
	}
}

func (d *driver) Done() {
	d.onExit = true
	d.onDraw.Lock()
	writeSequence(d.connDraw, 'X')
	d.onDraw.Unlock()
	if err := d.connEvent.Close(); err != nil {
		log.Fatalf("Close(e) %s", err)
	}
	if err := d.connDraw.Close(); err != nil {
		log.Fatalf("Close(d) %s", err)
	}
	if err := d.cmd.Wait(); err != nil {
		log.Fatalf("Wait %s", err)
	}
}

func (d *driver) Size(rect impress.Rect) {
	d.onDraw.Lock()
	defer d.onDraw.Unlock()
	writeSequence(d.connDraw, 'S', rect.X, rect.Y, rect.Width, rect.Height)
}

func (d *driver) Title(title string) {
	d.onDraw.Lock()
	defer d.onDraw.Unlock()
	writeSequence(d.connDraw, 'T', title)
}

func (d *driver) Chan() <-chan impress.Eventer {
	return d.events
}

// Event

func (d *driver) readEvents() {
	for {
		command, err := readChar(d.connEvent)
		if err != nil {
			if d.onExit {
				break
			}
			log.Fatalf("readEvents %s", err)
		}
		switch command {
		case 'k':
			u, _ := readUInt32(d.connEvent)
			shift, _ := readBool(d.connEvent)
			control, _ := readBool(d.connEvent)
			alt, _ := readBool(d.connEvent)
			meta, _ := readBool(d.connEvent)
			name, _ := readString(d.connEvent)
			d.events <- impress.KeyboardEvent{
				Rune:    rune(u),
				Name:    name,
				Shift:   shift,
				Control: control,
				Alt:     alt,
				Meta:    meta,
			}
		case 'g':
			u, _ := readUInt32(d.connEvent)
			d.events <- impress.GeneralEvent{
				Event: u,
			}
		case 'b':
			btype, _ := readChar(d.connEvent)
			button, _ := readChar(d.connEvent)
			x, _ := readInt16(d.connEvent)
			y, _ := readInt16(d.connEvent)
			d.events <- impress.ButtonEvent{
				Action: int(btype),
				Button: int(button),
				Point:  impress.NewPoint(x, y),
			}
		case 'm':
			x, _ := readInt16(d.connEvent)
			y, _ := readInt16(d.connEvent)
			shift, _ := readBool(d.connEvent)
			control, _ := readBool(d.connEvent)
			alt, _ := readBool(d.connEvent)
			meta, _ := readBool(d.connEvent)
			d.events <- impress.MotionEvent{
				Point:   impress.NewPoint(x, y),
				Shift:   shift,
				Control: control,
				Alt:     alt,
				Meta:    meta,
			}
		case 'u':
			name, _ := readString(d.connEvent)
			d.events <- impress.MenuEvent{
				Action: name,
			}
		default:
			d.events <- impress.UnknownEvent
		}
	}
}
