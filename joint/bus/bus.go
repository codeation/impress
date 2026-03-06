// Package implements an internal mechanism to communicate with an impress terminal.
package bus

import (
	"os"
)

const (
	fifoStreamPath = "/tmp/it_fifo_stream_"
	fifoInputPath  = "/tmp/it_fifo_input_"
	fifoOutputPath = "/tmp/it_fifo_output_"
	fifoEventPath  = "/tmp/it_fifo_event_"
)

type files struct {
	streamFile   *os.File
	requestFile  *os.File
	responseFile *os.File
	eventFile    *os.File
}
