// Package to connect to GTK driver
package duodriver

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/bus"
	"github.com/codeation/impress/joint/domain"
	"github.com/codeation/impress/joint/drawsend"
	"github.com/codeation/impress/joint/eventchan"
	"github.com/codeation/impress/joint/eventrecv"
	"github.com/codeation/impress/joint/lazy"
)

type duo struct {
	driver.Driver
	eventRecv interface{ Done() }
	cmd       *exec.Cmd
	server    *bus.ServerPipes
}

// New returns a new GTK driver
func New() (*duo, error) {
	d := &duo{}
	if err := d.connect(); err != nil {
		return nil, err
	}
	eventChan := eventchan.New()
	d.eventRecv = eventrecv.New(eventChan, d.server.EventPipe)
	drawSend := drawsend.New(d.server.StreamPipe, d.server.SyncPipe)
	duoDriver := domain.New(drawSend, eventChan, d.server.StreamPipe)
	d.Driver = lazy.New(duoDriver)
	return d, nil
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

	path, err := itPath()
	if errors.Is(err, errITNotFound) {
		log.Print("impress terminal not found")
		log.Print("set IMPRESS_TERMINAL_PATH environment variable")
		return fmt.Errorf("itPath: %w", err)
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
