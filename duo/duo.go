// Package to connect to GTK driver
package duo

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/codeation/impress"
	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/bus"
	"github.com/codeation/impress/joint/domain"
	"github.com/codeation/impress/joint/drawsend"
	"github.com/codeation/impress/joint/eventchan"
	"github.com/codeation/impress/joint/eventrecv"
	"github.com/codeation/impress/joint/lazy"
)

type doner interface {
	Done()
}

type duo struct {
	driver.Driver
	eventRecv doner
	cmd       *exec.Cmd
	server    *bus.ServerPipes
}

func init() {
	impress.Register(newDuo())
}

func newDuo() *duo {
	d := &duo{}
	if err := d.connect(); err != nil {
		log.Println(err)
		return nil
	}
	eventChan := eventchan.New()
	d.eventRecv = eventrecv.New(eventChan, d.server.EventPipe)
	drawSend := drawsend.New(d.server.StreamPipe, d.server.SyncPipe)
	duoDriver := domain.New(drawSend, eventChan, d.server.StreamPipe)
	d.Driver = lazy.New(duoDriver)
	return d
}

func (d *duo) Done() {
	d.eventRecv.Done()
	d.Driver.Done()
	if err := d.disconnect(); err != nil {
		log.Printf("done: %v", err)
		return
	}
}

func (d *duo) connect() error {
	var err error
	d.server, err = bus.NewServer()
	if err != nil {
		return fmt.Errorf("bus.NewServer: %w", err)
	}

	path := os.Getenv("IMPRESS_TERMINAL_PATH")
	if path == "" {
		path = "./it"
	}
	cmd := exec.Command(path, d.server.Suffix())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start: %w", err)
	}
	d.cmd = cmd

	if err := d.server.Connect(); err != nil {
		return fmt.Errorf("server.Connect: %w", err)
	}

	return nil
}

func (d *duo) disconnect() error {
	if err := d.cmd.Wait(); err != nil {
		return fmt.Errorf("cmd.Wait: %w", err)
	}

	if err := d.server.Close(); err != nil {
		return fmt.Errorf("server.Close: %w", err)
	}

	return nil
}
