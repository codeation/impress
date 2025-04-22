package drawrecv

import (
	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type SyncCommand struct {
	*drawRecv
}

func NewSyncCommand(calls iface.CallSet, streamPipe, syncPipe *rpc.Pipe) *SyncCommand {
	return &SyncCommand{
		drawRecv: &drawRecv{
			calls:      calls,
			streamPipe: streamPipe,
			syncPipe:   syncPipe,
		},
	}
}

func (d *SyncCommand) SyncCommand() error { return d.syncCommand() }
