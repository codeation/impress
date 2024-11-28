package bus

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/codeation/impress/joint/rpc"
)

type ClientPipes pipes

func NewClient(suffix string) (*ClientPipes, error) {
	return &ClientPipes{
		suffix: suffix,
	}, nil
}

func (p *ClientPipes) Connect() error {
	var err error
	if p.responseFile, err = os.OpenFile(fifoOutputPath+p.suffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(o): %w", err)
	}
	if p.eventFile, err = os.OpenFile(fifoEventPath+p.suffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(e): %w", err)
	}
	if p.requestFile, err = os.OpenFile(fifoInputPath+p.suffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(i): %w", err)
	}
	if p.streamFile, err = os.OpenFile(fifoStreamPath+p.suffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(s): %w", err)
	}

	if err = syscall.SetNonblock(int(p.streamFile.Fd()), true); err != nil {
		return fmt.Errorf("syscall.SetNonblck: %w", err)
	}

	p.StreamPipe = rpc.NewPipe(new(sync.Mutex), nil, p.streamFile)
	p.SyncPipe = rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(p.responseFile), p.requestFile)
	p.EventPipe = rpc.NewPipe(rpc.WithoutMutex(), bufio.NewWriter(p.eventFile), nil)

	return nil
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
