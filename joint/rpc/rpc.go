// Package implements an internal mechanism to communicate with an impress terminal.
package rpc

import (
	"encoding/binary"
	"fmt"
	"io"
)

type mutexer interface {
	Lock()
	Unlock()
}

type flushWriter interface {
	io.Writer
	Flush() error
}

type Pipe struct {
	mutex  mutexer
	writer flushWriter
	reader io.Reader
	err    error
}

func NewPipe(mutex mutexer, writer flushWriter, reader io.Reader) *Pipe {
	return &Pipe{
		mutex:  mutex,
		writer: writer,
		reader: reader,
	}
}

func (p *Pipe) Lock() *Pipe {
	p.err = nil
	p.mutex.Lock()
	return p
}

func (p *Pipe) Unlock() *Pipe {
	p.mutex.Unlock()
	return p
}

func (p *Pipe) Flush() *Pipe {
	if p.err != nil {
		return p
	}
	p.err = p.writer.Flush()
	return p
}

func (p *Pipe) Err() error {
	return p.err
}

func (p *Pipe) Get(variables ...interface{}) *Pipe {
	for _, v := range variables {
		if p.err != nil {
			return p
		}
		switch variable := v.(type) {
		case *byte:
			p.err = binary.Read(p.reader, binary.LittleEndian, variable)
		case *[]byte:
			p.err = p.getBytes(variable)
		case *bool:
			p.err = binary.Read(p.reader, binary.LittleEndian, variable)
		case *int:
			p.err = p.getInt(variable)
		case *[]int:
			p.err = p.getInts(variable)
		case *uint16:
			p.err = binary.Read(p.reader, binary.LittleEndian, variable)
		case *uint32:
			p.err = binary.Read(p.reader, binary.LittleEndian, variable)
		case *string:
			p.err = p.getString(variable)
		default:
			p.err = fmt.Errorf("unknown type: %T", v)
		}
	}
	return p
}

func (p *Pipe) Put(values ...interface{}) *Pipe {
	for _, v := range values {
		if p.err != nil {
			return p
		}
		switch value := v.(type) {
		case byte:
			p.err = binary.Write(p.writer, binary.LittleEndian, value)
		case []byte:
			p.err = p.putBytes(value)
		case bool:
			p.err = binary.Write(p.writer, binary.LittleEndian, value)
		case int:
			p.err = p.putInt(value)
		case []int:
			p.err = p.putInts(value)
		case uint16:
			p.err = binary.Write(p.writer, binary.LittleEndian, value)
		case uint32:
			p.err = binary.Write(p.writer, binary.LittleEndian, value)
		case string:
			p.err = p.putString(value)
		default:
			p.err = fmt.Errorf("unknown type: %T", v)
		}
	}
	return p
}
