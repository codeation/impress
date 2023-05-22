package rpc

import (
	"encoding/binary"
	"errors"
)

func (p *Pipe) putBytes(value []byte) error {
	if len(value) > 67108863 {
		return errors.New("slice too large")
	}
	if err := binary.Write(p.writer, binary.LittleEndian, int32(len(value))); err != nil {
		return err
	}
	return binary.Write(p.writer, binary.LittleEndian, value)
}

func (p *Pipe) putInt(value int) error {
	return binary.Write(p.writer, binary.LittleEndian, int16(value))
}

func (p *Pipe) putInts(value []int) error {
	if err := binary.Write(p.writer, binary.LittleEndian, int16(len(value))); err != nil {
		return err
	}
	for _, v := range value {
		if err := binary.Write(p.writer, binary.LittleEndian, int16(v)); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipe) putString(value string) error {
	if len(value) > 32767 {
		return errors.New("string too large")
	}
	if err := binary.Write(p.writer, binary.LittleEndian, int32(len(value))); err != nil {
		return err
	}
	return binary.Write(p.writer, binary.LittleEndian, []byte(value))
}
