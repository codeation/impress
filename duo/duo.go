// Package to connect to GTK driver
package next

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/codeation/impress"
	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/domain"
	"github.com/codeation/impress/joint/drawsend"
	"github.com/codeation/impress/joint/eventchan"
	"github.com/codeation/impress/joint/eventrecv"
	"github.com/codeation/impress/joint/rpc"
)

const (
	fifoStreamPath = "/tmp/it_fifo_stream_"
	fifoInputPath  = "/tmp/it_fifo_input_"
	fifoOutputPath = "/tmp/it_fifo_output_"
	fifoEventPath  = "/tmp/it_fifo_event_"
)

type doner interface {
	Done()
}

type flusher interface {
	Sync() error
}

type duo struct {
	driver.Driver
	eventRecv    doner
	cmd          *exec.Cmd
	fileSuffix   string
	streamFile   *os.File
	requestFile  *os.File
	responseFile *os.File
	eventFile    *os.File
	streamPipe   *rpc.Pipe
	syncPipe     *rpc.Pipe
	eventPipe    *rpc.Pipe
}

func init() {
	impress.Register(newDuo())
}

func newDuo() *duo {
	d := &duo{}
	if err := d.connect(); err != nil {
		log.Println(err)
		return nil
	}
	eventChan := eventchan.New()
	d.eventRecv = eventrecv.New(eventChan, d.eventPipe)
	drawSend := drawsend.New(d.streamPipe, d.syncPipe)
	d.Driver = domain.New(drawSend, eventChan, d.streamPipe)
	return d
}

func (d *duo) Done() {
	d.eventRecv.Done()
	d.Driver.Done()
	if err := d.disconnect(); err != nil {
		log.Println(err)
		return
	}
}

func (d *duo) connect() error {
	randBuffer := make([]byte, 8)
	if _, err := rand.Reader.Read(randBuffer); err != nil {
		return fmt.Errorf("rand.Reader.Read: %w", err)
	}
	d.fileSuffix = hex.EncodeToString(randBuffer)
	for _, name := range []string{fifoInputPath, fifoStreamPath, fifoOutputPath, fifoEventPath} {
		if err := syscall.Mkfifo(name+d.fileSuffix, 0644); err != nil {
			return fmt.Errorf("syscall.Mkfifo: %w", err)
		}
	}

	path := os.Getenv("IMPRESS_TERMINAL_PATH")
	if path == "" {
		path = "./it"
	}
	cmd := exec.Command(path, d.fileSuffix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start: %w", err)
	}
	d.cmd = cmd

	var err error
	if d.responseFile, err = os.OpenFile(fifoOutputPath+d.fileSuffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(o): %w", err)
	}
	if d.eventFile, err = os.OpenFile(fifoEventPath+d.fileSuffix, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(e): %w", err)
	}
	if d.requestFile, err = os.OpenFile(fifoInputPath+d.fileSuffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(i): %w", err)
	}
	if d.streamFile, err = os.OpenFile(fifoStreamPath+d.fileSuffix, os.O_WRONLY, os.ModeNamedPipe); err != nil {
		return fmt.Errorf("os.OpenFile(s): %w", err)
	}

	streamBuffered := bufio.NewWriter(d.streamFile)

	d.streamPipe = rpc.NewPipe(new(sync.Mutex), streamBuffered, nil)
	d.syncPipe = rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(d.requestFile), bufio.NewReader(d.responseFile))
	d.eventPipe = rpc.NewPipe(rpc.WithoutMutex(), nil, bufio.NewReader(d.eventFile))

	return nil
}

func (d *duo) disconnect() error {
	if err := d.requestFile.Close(); err != nil {
		return fmt.Errorf("fileRequest.Close: %w", err)
	}
	if err := d.streamFile.Close(); err != nil {
		return fmt.Errorf("fileStream.Close: %w", err)
	}

	if err := d.cmd.Wait(); err != nil {
		return fmt.Errorf("cmd.Wait: %w", err)
	}

	if err := d.responseFile.Close(); err != nil {
		return fmt.Errorf("fileResponse.Close: %w", err)
	}
	if err := d.eventFile.Close(); err != nil {
		return fmt.Errorf("fileEvent.Close: %w", err)
	}

	for _, name := range []string{fifoInputPath, fifoOutputPath, fifoEventPath} {
		if _, err := os.Stat(name + d.fileSuffix); err == nil || !errors.Is(err, os.ErrNotExist) {
			if err = os.Remove(name + d.fileSuffix); err != nil {
				return fmt.Errorf("os.Remove: %w", err)
			}
		}
	}
	return nil
}
