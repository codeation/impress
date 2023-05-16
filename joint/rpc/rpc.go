// Package implements an internal mechanism to communicate with an impress terminal.
package rpc

import (
	"fmt"
	"io"
	"log"
)

type mutexer interface {
	Lock()
	Unlock()
}

type flushWriter interface {
	io.Writer
	Flush() error
}

type parameter interface {
	set(reader io.Reader) error
}

type Pipe struct {
	mutex  mutexer
	writer flushWriter
	reader io.Reader
	queue  []parameter
}

func NewPipe(mutex mutexer, writer flushWriter, reader io.Reader) *Pipe {
	return &Pipe{
		mutex:  mutex,
		writer: writer,
		reader: reader,
	}
}

func (p *Pipe) clone() *Pipe {
	if len(p.queue) == 0 {
		return &Pipe{
			mutex:  p.mutex,
			writer: p.writer,
			reader: p.reader,
			queue:  make([]parameter, 0, 8),
		}
	}
	return p
}

func (p *Pipe) add(variable parameter) *Pipe {
	output := p.clone()
	output.queue = append(output.queue, variable)
	return output
}

func (p *Pipe) Byte(variable *byte) *Pipe {
	return p.add(&parameterByte{ptr: variable})
}

func (p *Pipe) Bytes(variable *[]byte) *Pipe {
	return p.add(&parameterBytes{ptr: variable})
}

func (p *Pipe) Bool(variable *bool) *Pipe {
	return p.add(&parameterBool{ptr: variable})
}

func (p *Pipe) Int(variable *int) *Pipe {
	return p.add(&parameterInt{ptr: variable})
}

func (p *Pipe) Ints(variable *[]int) *Pipe {
	return p.add(&parameterInts{ptr: variable})
}

func (p *Pipe) UInt16(variable *uint16) *Pipe {
	return p.add(&parameterUInt16{ptr: variable})
}

func (p *Pipe) UInt32(variable *uint32) *Pipe {
	return p.add(&parameterUInt32{ptr: variable})
}

func (p *Pipe) String(variable *string) *Pipe {
	return p.add(&parameterString{ptr: variable})
}

func (p *Pipe) CallErr(values ...interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, v := range values {
		switch value := v.(type) {
		case byte:
			if err := putByte(p.writer, value); err != nil {
				return err
			}
		case []byte:
			if err := putBytes(p.writer, value); err != nil {
				return err
			}
		case bool:
			if err := putBool(p.writer, value); err != nil {
				return err
			}
		case int:
			if err := putInt(p.writer, value); err != nil {
				return err
			}
		case []int:
			if err := putInts(p.writer, value); err != nil {
				return err
			}
		case uint16:
			if err := putUInt16(p.writer, value); err != nil {
				return err
			}
		case uint32:
			if err := putUInt32(p.writer, value); err != nil {
				return err
			}
		case string:
			if err := putString(p.writer, value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown type: %T", v)
		}
	}

	if len(values) != 0 && len(p.queue) != 0 {
		if err := p.writer.Flush(); err != nil {
			return fmt.Errorf("flush: %w", err)
		}
	}

	for _, v := range p.queue {
		if err := v.set(p.reader); err != nil {
			return fmt.Errorf("set: %w", err)
		}
	}

	return nil
}

func (p *Pipe) Call(values ...interface{}) {
	if err := p.CallErr(values...); err != nil {
		log.Printf("call: %v", err)
	}
}

func (p *Pipe) Flush() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if err := p.writer.Flush(); err != nil {
		log.Printf("flush: %v", err)
	}
}
