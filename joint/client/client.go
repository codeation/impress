// Package implements an internal mechanism to communicate with an impress terminal.
package client

import (
	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type Client struct {
	callbacks  iface.CallbackSet
	onExit     bool
	eventPipe  *rpc.Pipe
	streamPipe *rpc.Pipe
	syncPipe   *rpc.Pipe
}

func NewClient(callbacks iface.CallbackSet, eventPipe, streamPipe, syncPipe *rpc.Pipe) *Client {
	c := &Client{
		callbacks:  callbacks,
		eventPipe:  eventPipe,
		streamPipe: streamPipe,
		syncPipe:   syncPipe,
	}
	go c.listen()
	return c
}

func (c *Client) Sync() {
	c.streamPipe.Flush()
}
