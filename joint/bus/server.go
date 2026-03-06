package bus

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"syscall"

	"github.com/codeation/impress/joint/rpc"
)

const streamBufferSize = 64 * 1024

type Runner interface {
	Run(suffix string) error
	Wait() error
}

type ServerPipes struct {
	f      files
	runner Runner
}

func NewServer(runner Runner) (*ServerPipes, error) {
	var randBuffer [8]byte
	if _, err := rand.Reader.Read(randBuffer[:]); err != nil {
		return nil, fmt.Errorf("rand.Reader.Read: %w", err)
	}
	suffix := hex.EncodeToString(randBuffer[:])

	for _, name := range []string{fifoInputPath, fifoStreamPath, fifoOutputPath, fifoEventPath} {
		if err := syscall.Mkfifo(name+suffix, 0600); err != nil {
			return nil, fmt.Errorf("syscall.Mkfifo: %w", err)
		}
		defer os.Remove(name + suffix)
	}

	if err := runner.Run(suffix); err != nil {
		return nil, fmt.Errorf("runner.Run: %w", err)
	}

	// opening sequence is the same as for client
	responseFile, err := os.OpenFile(fifoOutputPath+suffix, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(o): %w", err)
	}
	eventFile, err := os.OpenFile(fifoEventPath+suffix, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(e): %w", err)
	}
	requestFile, err := os.OpenFile(fifoInputPath+suffix, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(i): %w", err)
	}
	streamFile, err := os.OpenFile(fifoStreamPath+suffix, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(s): %w", err)
	}

	return &ServerPipes{
		f: files{
			streamFile:   streamFile,
			requestFile:  requestFile,
			responseFile: responseFile,
			eventFile:    eventFile,
		},
		runner: runner,
	}, nil
}

func (p *ServerPipes) NewEventPipe() *rpc.Pipe {
	return rpc.NewPipe(nil, bufio.NewReader(p.f.eventFile))
}

func (p *ServerPipes) NewStreamPipe() *rpc.Pipe {
	streamBuffered := bufio.NewWriterSize(p.f.streamFile, streamBufferSize)
	return rpc.NewPipe(streamBuffered, nil)
}

func (p *ServerPipes) NewSyncPipe() *rpc.Pipe {
	return rpc.NewPipe(bufio.NewWriter(p.f.requestFile), bufio.NewReader(p.f.responseFile))
}

func (p *ServerPipes) Close() error {
	if err := p.runner.Wait(); err != nil {
		return fmt.Errorf("runner.Wait: %w", err)
	}

	if err := p.f.requestFile.Close(); err != nil {
		return fmt.Errorf("requestFile.Close(): %w", err)
	}
	if err := p.f.streamFile.Close(); err != nil {
		return fmt.Errorf("streamFile.Close: %w", err)
	}

	if err := p.f.responseFile.Close(); err != nil {
		return fmt.Errorf("responseFile.Close: %w", err)
	}
	if err := p.f.eventFile.Close(); err != nil {
		return fmt.Errorf("eventFile.Close: %w", err)
	}

	return nil
}
