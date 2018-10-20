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
	eternal      sync.Mutex
	lastWindowID int
	lastFontID   int
	events       chan impress.Eventer
	connDraw     net.Conn
	connEvent    net.Conn
}

func init() {
	d := &driver{}
	d.runrpc()
	d.events = make(chan impress.Eventer)
	impress.Register(d)
}

func (d *driver) runrpc() {
	d.cmd = exec.Command("../impress/terminal/it")
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
	d.eternal.Lock()
}

func (d *driver) Main() {
	d.eternal.Lock()
	d.eternal.Unlock()
}

func (d *driver) Done() {
	d.eternal.Unlock()
	writeSequence(d.connDraw, 'X')
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
	writeSequence(d.connDraw, 'S', rect.X, rect.Y, rect.Width, rect.Height)
}

func (d *driver) Title(title string) {
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
	writeSequence(d.connDraw, 'D', w.id, w.rect.X, w.rect.Y, w.rect.Width, w.rect.Height,
		color.R, color.G, color.B)
	w.Fill(rect, color)
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
	writeSequence(p.driver.connDraw, 'C', p.id)
	p.Fill(p.rect, p.background)
}

func (p *painter) Fill(rect impress.Rect, color impress.Color) {
	writeSequence(p.driver.connDraw, 'F', p.id, rect.X, rect.Y, rect.Width, rect.Height,
		color.R, color.G, color.B)
}

func (p *painter) Line(from impress.Point, to impress.Point, color impress.Color) {
	writeSequence(p.driver.connDraw, 'L', p.id, from.X, from.Y, to.X, to.Y,
		color.R, color.G, color.B)
}

func (p *painter) Text(text string, font *impress.Font, from impress.Point, color impress.Color) {
	f, ok := font.Fonter.(*ftfont)
	if ok {
		f.load()
	}
	writeSequence(p.driver.connDraw, 'U', p.id, from.X, from.Y, color.R, color.G, color.B,
		f.ID, font.Height, text)
}

func (p *painter) Show() {
	writeSequence(p.driver.connDraw, 'W', p.id)
}

// Event

func (d *driver) readEvents() {
	for {
		command := readChar(d.connEvent)
		switch command {
		case 'k':
			u := readUInt32(d.connEvent)
			shift := readBool(d.connEvent)
			control := readBool(d.connEvent)
			alt := readBool(d.connEvent)
			meta := readBool(d.connEvent)
			name := readString(d.connEvent)
			d.events <- impress.KeyboardEvent{
				Rune:    rune(u),
				Name:    name,
				Shift:   shift,
				Control: control,
				Alt:     alt,
				Meta:    meta,
			}
		case 'g':
			u := readUInt32(d.connEvent)
			d.events <- impress.GeneralEvent{
				Event: u,
			}
		}
	}
}
