// Package implements an internal mechanism to communicate with an impress terminal.
package rpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

type flushWriter interface {
	io.Writer
	Flush() error
}

type Pipe struct {
	mutex  sync.Mutex
	writer flushWriter
	reader io.Reader
}

func NewPipe(writer flushWriter, reader io.Reader) *Pipe {
	return &Pipe{
		writer: writer,
		reader: reader,
	}
}

func (p *Pipe) Get(variables ...any) error {
	for _, v := range variables {
		var err error
		switch variable := v.(type) {
		case *byte:
			err = binary.Read(p.reader, binary.LittleEndian, variable)
		case *[]byte:
			err = p.getBytes(variable)
		case *bool:
			err = binary.Read(p.reader, binary.LittleEndian, variable)
		case *int:
			err = p.getInt(variable)
		case *[]int:
			err = p.getInts(variable)
		case *uint16:
			err = binary.Read(p.reader, binary.LittleEndian, variable)
		case *uint32:
			err = binary.Read(p.reader, binary.LittleEndian, variable)
		case *string:
			err = p.getString(variable)
		default:
			err = fmt.Errorf("unknown type: %T", v)
		}
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Fatalf("get: %v", err)
			}
			if !errors.Is(err, os.ErrClosed) && !errors.Is(err, io.EOF) {
				log.Printf("get: %v", err)
			}
			return err
		}
	}
	return nil
}

func (p *Pipe) Put(values ...any) error {
	for _, v := range values {
		var err error
		switch value := v.(type) {
		case byte:
			err = binary.Write(p.writer, binary.LittleEndian, value)
		case []byte:
			err = p.putBytes(value)
		case bool:
			err = binary.Write(p.writer, binary.LittleEndian, value)
		case int:
			err = p.putInt(value)
		case []int:
			err = p.putInts(value)
		case uint16:
			err = binary.Write(p.writer, binary.LittleEndian, value)
		case uint32:
			err = binary.Write(p.writer, binary.LittleEndian, value)
		case string:
			err = p.putString(value)
		default:
			err = fmt.Errorf("unknown type: %T", v)
		}
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Fatalf("put: %v", err)
			}
			log.Printf("put: %v", err)
			return err
		}
	}
	return nil
}

func (p *Pipe) PutTx(values ...any) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.Put(values...)
}

func (p *Pipe) IO(putValues []any, getValues []any) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if err := p.Put(putValues...); err != nil {
		return fmt.Errorf("put: %w", err)
	}
	if err := p.writer.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	if err := p.Get(getValues...); err != nil {
		return fmt.Errorf("get: %w", err)
	}
	return nil
}

func (p *Pipe) Sync() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.writer.Flush()
}
