// Package implements an internal mechanism to communicate with an impress terminal.
package remote

import (
	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type server struct {
	calls      iface.CallSet
	streamPipe *rpc.Pipe
	syncPipe   *rpc.Pipe
}

func New(calls iface.CallSet, streamPipe, syncPipe *rpc.Pipe) *server {
	s := &server{
		calls:      calls,
		streamPipe: streamPipe,
		syncPipe:   syncPipe,
	}
	go s.streamListen()
	go s.syncListen()
	return s
}
