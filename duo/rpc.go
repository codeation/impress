package duo

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

type pipe struct {
	mutex  mutexer
	writer flushWriter
	reader io.Reader
	queue  []parameter
}

func newPipe(mutex mutexer, writer flushWriter, reader io.Reader) *pipe {
	return &pipe{
		mutex:  mutex,
		writer: writer,
		reader: reader,
	}
}

func (p *pipe) clone() *pipe {
	if len(p.queue) == 0 {
		return &pipe{
			mutex:  p.mutex,
			writer: p.writer,
			reader: p.reader,
			queue:  make([]parameter, 0, 8),
		}
	}
	return p
}

func (p *pipe) add(variable parameter) *pipe {
	output := p.clone()
	output.queue = append(output.queue, variable)
	return output
}

func (p *pipe) Bool(variable *bool) *pipe {
	return p.add(&parameterBool{ptr: variable})
}

func (p *pipe) Byte(variable *byte) *pipe {
	return p.add(&parameterByte{ptr: variable})
}

func (p *pipe) Char(variable *int) *pipe {
	return p.add(&parameterChar{ptr: variable})
}

func (p *pipe) Int16(variable *int) *pipe {
	return p.add(&parameterInt16{ptr: variable})
}

func (p *pipe) Int16s(variable *[]int) *pipe {
	return p.add(&parameterInt16s{ptr: variable})
}

func (p *pipe) Rune(variable *rune) *pipe {
	return p.add(&parameterRune{ptr: variable})
}

func (p *pipe) String(variable *string) *pipe {
	return p.add(&parameterString{ptr: variable})
}

func (p *pipe) UInt32(variable *int) *pipe {
	return p.add(&parameterUInt32{ptr: variable})
}

func (p *pipe) CallErr(values ...interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, v := range values {
		switch value := v.(type) {
		case int32:
			if err := putInt32(p.writer, value); err != nil {
				return err
			}
		case uint32:
			if err := putUInt32(p.writer, value); err != nil {
				return err
			}
		case byte:
			if err := putByte(p.writer, value); err != nil {
				return err
			}
		case int:
			if err := putInt(p.writer, value); err != nil {
				return err
			}
		case string:
			if err := putString(p.writer, value); err != nil {
				return err
			}
		case []byte:
			if err := putBytes(p.writer, value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown type: %T", v)
		}
	}

	if len(values) != 0 && len(p.queue) != 0 {
		if err := p.writer.Flush(); err != nil {
			return err
		}
	}

	for _, v := range p.queue {
		if err := v.set(p.reader); err != nil {
			return err
		}
	}

	return nil
}

func (p *pipe) Call(values ...interface{}) {
	if err := p.CallErr(values...); err != nil {
		log.Fatal(err)
	}
}

func (p *pipe) Flush() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if err := p.writer.Flush(); err != nil {
		log.Fatal(err)
	}
}
