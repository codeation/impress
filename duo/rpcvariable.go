package duo

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
)

type parameterBool struct {
	ptr *bool
}

func (p *parameterBool) set(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, p.ptr)
}

type parameterByte struct {
	ptr *byte
}

func (p *parameterByte) set(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, p.ptr)
}

type parameterChar struct {
	ptr *int
}

func (p *parameterChar) set(reader io.Reader) error {
	var value byte
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return err
	}
	*p.ptr = int(value)
	return nil
}

type parameterInt16 struct {
	ptr *int
}

func (p *parameterInt16) set(reader io.Reader) error {
	var value int16
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return err
	}
	*p.ptr = int(value)
	return nil
}

type parameterInt16s struct {
	ptr *[]int
}

func (p parameterInt16s) set(reader io.Reader) error {
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

type parameterRune struct {
	ptr *rune
}

func (p *parameterRune) set(reader io.Reader) error {
	var value uint32
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return err
	}
	*p.ptr = rune(value)
	return nil
}

type parameterString struct {
	ptr *string
}

func (p *parameterString) set(reader io.Reader) error {
	var length int16
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		log.Fatal(err)
	}
	data := make([]byte, length)
	if err := binary.Read(reader, binary.LittleEndian, &data); err != nil {
		log.Fatal(err)
	}
	*p.ptr = string(data)
	return nil
}

type parameterUInt32 struct {
	ptr *int
}

func (p *parameterUInt32) set(reader io.Reader) error {
	var value uint32
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return err
	}
	*p.ptr = int(value)
	return nil
}

func putByte(writer io.Writer, value byte) error {
	return binary.Write(writer, binary.LittleEndian, value)
}

func putBytes(writer io.Writer, value []byte) error {
	if len(value) > 67108863 {
		return errors.New("slice too large")
	}
	if err := binary.Write(writer, binary.LittleEndian, int32(len(value))); err != nil {
		return err
	}
	return binary.Write(writer, binary.LittleEndian, value)
}

func putInt(writer io.Writer, value int) error {
	return binary.Write(writer, binary.LittleEndian, int16(value))
}

func putUInt32(writer io.Writer, value uint32) error {
	return binary.Write(writer, binary.LittleEndian, uint16(value))
}

func putInt32(writer io.Writer, value int32) error {
	return binary.Write(writer, binary.LittleEndian, byte(value))
}

func putString(writer io.Writer, value string) error {
	if len(value) > 32767 {
		return errors.New("string too large")
	}
	if err := binary.Write(writer, binary.LittleEndian, int32(len(value))); err != nil {
		return err
	}
	return binary.Write(writer, binary.LittleEndian, []byte(value))
}

type dummyMutex struct{}

func (*dummyMutex) Lock()   {}
func (*dummyMutex) Unlock() {}
