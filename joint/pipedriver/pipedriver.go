// Package implements an internal mechanism to communicate with an impress terminal.
package pipedriver

import (
	"log"

	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/domain"
	"github.com/codeation/impress/joint/drawsend"
	"github.com/codeation/impress/joint/eventchan"
	"github.com/codeation/impress/joint/eventrecv"
	"github.com/codeation/impress/joint/rpc"
)

type PipeCreator interface {
	NewEventPipe() *rpc.Pipe
	NewStreamPipe() *rpc.Pipe
	NewSyncPipe() *rpc.Pipe
	Close() error
}

type Driver struct {
	driver.Driver
	p         PipeCreator
	eventRecv interface{ Done() }
}

func New(p PipeCreator) *Driver {
	eventPipe := p.NewEventPipe()
	streamPipe := p.NewStreamPipe()
	syncPipe := p.NewSyncPipe()

	eventChan := eventchan.New()
	eventRecv := eventrecv.New(eventChan, eventPipe)
	drawSend := drawsend.New(streamPipe, syncPipe)
	driver := domain.New(drawSend, eventChan, streamPipe)

	return &Driver{
		Driver:    driver,
		p:         p,
		eventRecv: eventRecv,
	}
}

func (d *Driver) Done() {
	d.eventRecv.Done()
	d.Driver.Done()
	if err := d.p.Close(); err != nil {
		log.Printf("pipedriver.Done: p.Close: %v", err)
	}
}
