package rpc

import (
	"encoding/binary"
)

func (p *Pipe) getBytes(variable *[]byte) error {
	var length int32
	if err := binary.Read(p.reader, binary.LittleEndian, &length); err != nil {
		return err
	}
	data := make([]byte, length)
	if err := binary.Read(p.reader, binary.LittleEndian, &data); err != nil {
		return err
	}
	*variable = data
	return nil
}

func (p *Pipe) getInt(variable *int) error {
	var value int16
	if err := binary.Read(p.reader, binary.LittleEndian, &value); err != nil {
		return err
	}
	*variable = int(value)
	return nil
}

func (p *Pipe) getInts(variable *[]int) error {
	var length int16
	if err := binary.Read(p.reader, binary.LittleEndian, &length); err != nil {
		return err
	}
	*variable = make([]int, length)
	for i := 0; i < int(length); i++ {
		var value int16
		if err := binary.Read(p.reader, binary.LittleEndian, &value); err != nil {
			return err
		}
		(*variable)[i] = int(value)
	}
	return nil
}

func (p *Pipe) getString(variable *string) error {
	var length int32
	if err := binary.Read(p.reader, binary.LittleEndian, &length); err != nil {
		return err
	}
	data := make([]byte, length)
	if err := binary.Read(p.reader, binary.LittleEndian, &data); err != nil {
		return err
	}
	*variable = string(data)
	return nil
}
