package rpc

import (
	"encoding/binary"
	"io"
)

type parameterByte struct {
	ptr *byte
}

func (p *parameterByte) set(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, p.ptr)
}

type parameterBytes struct {
	ptr *[]byte
}

func (p *parameterBytes) set(reader io.Reader) error {
	var length int32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return err
	}
	data := make([]byte, length)
	if err := binary.Read(reader, binary.LittleEndian, &data); err != nil {
		return err
	}
	*p.ptr = data
	return nil
}

type parameterBool struct {
	ptr *bool
}

func (p *parameterBool) set(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, p.ptr)
}

type parameterInt struct {
	ptr *int
}

func (p *parameterInt) set(reader io.Reader) error {
	var value int16
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return err
	}
	*p.ptr = int(value)
	return nil
}

type parameterInts struct {
	ptr *[]int
}

func (p parameterInts) set(reader io.Reader) error {
	var length int16
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return err
	}
	*p.ptr = make([]int, length)
	for i := 0; i < int(length); i++ {
		var value int16
		if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
			return err
		}
		(*p.ptr)[i] = int(value)
	}
	return nil
}

type parameterUInt16 struct {
	ptr *uint16
}

func (p *parameterUInt16) set(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, p.ptr)
}

type parameterUInt32 struct {
	ptr *uint32
}

func (p *parameterUInt32) set(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, p.ptr)
}

type parameterString struct {
	ptr *string
}

func (p *parameterString) set(reader io.Reader) error {
	var length int32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return err
	}
	data := make([]byte, length)
	if err := binary.Read(reader, binary.LittleEndian, &data); err != nil {
		return err
	}
	*p.ptr = string(data)
	return nil
}
