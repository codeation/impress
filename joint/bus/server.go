package bus

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/codeation/impress/joint/rpc"
)

type ServerPipes pipes

func NewServer() (*ServerPipes, error) {
	randBuffer := make([]byte, 8)
	if _, err := rand.Reader.Read(randBuffer); err != nil {
		return nil, fmt.Errorf("rand.Reader.Read: %w", err)
	}
	suffix := hex.EncodeToString(randBuffer)
	for _, name := range []string{fifoInputPath, fifoStreamPath, fifoOutputPath, fifoEventPath} {
		if err := syscall.Mkfifo(name+suffix, 0600); err != nil {
			return nil, fmt.Errorf("syscall.Mkfifo: %w", err)
		}
	}
	return &ServerPipes{
		suffix: suffix,
	}, nil
}

func (p *ServerPipes) Connect() error {
	var err error
	if p.responseFile, err = os.OpenFile(fifoOutputPath+p.suffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(o): %w", err)
	}
	if p.eventFile, err = os.OpenFile(fifoEventPath+p.suffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(e): %w", err)
	}
	if p.requestFile, err = os.OpenFile(fifoInputPath+p.suffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(i): %w", err)
	}
	if p.streamFile, err = os.OpenFile(fifoStreamPath+p.suffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(s): %w", err)
	}

	streamBuffered := bufio.NewWriterSize(p.streamFile, defaultBufferSize)

	p.StreamPipe = rpc.NewPipe(new(sync.Mutex), streamBuffered, nil)
	p.SyncPipe = rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(p.requestFile), bufio.NewReader(p.responseFile))
	p.EventPipe = rpc.NewPipe(rpc.WithoutMutex(), nil, bufio.NewReader(p.eventFile))

	return nil
}

func (p *ServerPipes) Close() error {
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

	for _, name := range []string{fifoInputPath, fifoStreamPath, fifoOutputPath, fifoEventPath} {
		if _, err := os.Stat(name + p.suffix); err == nil || !errors.Is(err, os.ErrNotExist) {
			if err = os.Remove(name + p.suffix); err != nil {
				return fmt.Errorf("os.Remove: %w", err)
			}
		}
	}
	return nil
}

func (p *ServerPipes) Suffix() string { return p.suffix }
