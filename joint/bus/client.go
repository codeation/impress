package bus

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"github.com/codeation/impress/joint/rpc"
)

type ClientPipes files

func NewClient(suffix string) (*ClientPipes, error) {
	responseFile, err := os.OpenFile(fifoOutputPath+suffix, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(o): %w", err)
	}
	eventFile, err := os.OpenFile(fifoEventPath+suffix, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(e): %w", err)
	}
	requestFile, err := os.OpenFile(fifoInputPath+suffix, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(i): %w", err)
	}
	streamFile, err := os.OpenFile(fifoStreamPath+suffix, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(s): %w", err)
	}

	if err := syscall.SetNonblock(int(streamFile.Fd()), true); err != nil {
		return nil, fmt.Errorf("syscall.SetNonblck: %w", err)
	}

	return &ClientPipes{
		streamFile:   streamFile,
		requestFile:  requestFile,
		responseFile: responseFile,
		eventFile:    eventFile,
	}, nil
}

func (p *ClientPipes) NewEventPipe() *rpc.Pipe {
	return rpc.NewPipe(bufio.NewWriter(p.eventFile), nil)
}

func (p *ClientPipes) NewStreamPipe() *rpc.Pipe {
	return rpc.NewPipe(nil, p.streamFile)
}

func (p *ClientPipes) NewSyncPipe() *rpc.Pipe {
	return rpc.NewPipe(bufio.NewWriter(p.responseFile), p.requestFile)
}

func (p *ClientPipes) Close() error {
	if err := p.requestFile.Close(); err != nil {
		return fmt.Errorf("requestFile.Close(): %w", err)
	}
	if err := p.streamFile.Close(); err != nil {
		return fmt.Errorf("streamFile.Close: %w", err)
	}

	if err := p.responseFile.Close(); err != nil {
		return fmt.Errorf("responseFile.Close: %w", err)
	}
	if err := p.eventFile.Close(); err != nil {
		return fmt.Errorf("eventFile.Close: %w", err)
	}

	return nil
}

func (p *ClientPipes) StreamFile() *os.File  { return p.streamFile }
func (p *ClientPipes) RequestFile() *os.File { return p.requestFile }
