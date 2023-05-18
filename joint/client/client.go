// Package implements an internal mechanism to communicate with an impress terminal.
package client

import (
	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type client struct {
	callbacks  iface.CallbackSet
	onExit     bool
	eventPipe  *rpc.Pipe
	streamPipe *rpc.Pipe
	syncPipe   *rpc.Pipe
}

func New(callbacks iface.CallbackSet, eventPipe, streamPipe, syncPipe *rpc.Pipe) *client {
	c := &client{
		callbacks:  callbacks,
		eventPipe:  eventPipe,
		streamPipe: streamPipe,
		syncPipe:   syncPipe,
	}
	go c.listen()
	return c
}

func (c *client) Sync() {
	c.streamPipe.Flush()
}
