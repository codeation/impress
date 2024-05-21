// Package implements an internal mechanism to communicate with an impress terminal.
package bus

import (
	"os"

	"github.com/codeation/impress/joint/rpc"
)

const (
	fifoStreamPath = "/tmp/it_fifo_stream_"
	fifoInputPath  = "/tmp/it_fifo_input_"
	fifoOutputPath = "/tmp/it_fifo_output_"
	fifoEventPath  = "/tmp/it_fifo_event_"
)

const defaultBufferSize = 256 * 1024

type pipes struct {
	suffix       string
	streamFile   *os.File
	requestFile  *os.File
	responseFile *os.File
	eventFile    *os.File
	StreamPipe   *rpc.Pipe
	SyncPipe     *rpc.Pipe
	EventPipe    *rpc.Pipe
}
