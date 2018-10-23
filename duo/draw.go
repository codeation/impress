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
	events       chan impress.Eventer
	onDraw       sync.Mutex
	connDraw     net.Conn
	onExit       bool
	connEvent    net.Conn
}

func init() {
	d := &driver{}
	d.runrpc()
	d.events = make(chan impress.Eventer)
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
		log.Fatal(err)
	}
	if err := d.connDraw.Close(); err != nil {
		log.Fatal(err)
	}
	if err := d.cmd.Wait(); err != nil {
		log.Fatal(err)
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

func (d *driver) NewWindow(rect impress.Rect, color impress.Color) impress.Painter {
	d.lastWindowID++
	w := &painter{
		driver:     d,
		id:         d.lastWindowID,
		rect:       rect,
		background: color,
	}
	d.onDraw.Lock()
	defer d.onDraw.Unlock()
	writeSequence(d.connDraw, 'D', w.id, w.rect.X, w.rect.Y, w.rect.Width, w.rect.Height,
		color.R, color.G, color.B)
	writeSequence(d.connDraw, 'F', w.id, w.rect.X, w.rect.Y, w.rect.Width, w.rect.Height,
		w.background.R, w.background.G, w.background.B)
	return w
}

func (d *driver) Event() impress.Eventer {
	return <-d.events
}

// Paint

type painter struct {
	driver     *driver
	id         int
	rect       impress.Rect
	background impress.Color
}

func (p *painter) Clear() {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'C', p.id)
	writeSequence(p.driver.connDraw, 'F', p.id, p.rect.X, p.rect.Y, p.rect.Width, p.rect.Height,
		p.background.R, p.background.G, p.background.B)
}

func (p *painter) Fill(rect impress.Rect, color impress.Color) {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'F', p.id, rect.X, rect.Y, rect.Width, rect.Height,
		color.R, color.G, color.B)
}

func (p *painter) Line(from impress.Point, to impress.Point, color impress.Color) {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'L', p.id, from.X, from.Y, to.X, to.Y,
		color.R, color.G, color.B)
}

func (p *painter) Text(text string, font *impress.Font, from impress.Point, color impress.Color) {
	f, ok := font.Fonter.(*ftfont)
	if ok {
		f.load()
	}
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'U', p.id, from.X, from.Y, color.R, color.G, color.B,
		f.ID, font.Height, text)
}

func (p *painter) Show() {
	p.driver.onDraw.Lock()
	defer p.driver.onDraw.Unlock()
	writeSequence(p.driver.connDraw, 'W', p.id)
}

// Event

func (d *driver) readEvents() {
	for !d.onExit {
		command, err := readChar(d.connEvent)
		if d.onExit {
			break
		}
		if err != nil {
			log.Println(err)
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
		}
	}
}
