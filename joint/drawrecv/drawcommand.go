package drawrecv

import (
	"github.com/codeation/impress/joint/iface"
	"github.com/codeation/impress/joint/rpc"
)

type DrawCommand struct {
	*drawRecv
}

func NewDrawCommand(calls iface.CallSet, streamPipe, syncPipe *rpc.Pipe) *DrawCommand {
	return &DrawCommand{
		drawRecv: &drawRecv{
			calls:      calls,
			streamPipe: streamPipe,
			syncPipe:   syncPipe,
		},
	}
}

func (d *DrawCommand) StreamCommand() error { return d.drawRecv.streamCommand() }
func (d *DrawCommand) SyncCommand() error   { return d.drawRecv.syncCommand() }
