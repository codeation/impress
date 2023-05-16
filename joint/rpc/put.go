package rpc

import (
	"encoding/binary"
	"errors"
	"io"
)

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

func putBool(writer io.Writer, value bool) error {
	return binary.Write(writer, binary.LittleEndian, value)
}

func putInt(writer io.Writer, value int) error {
	return binary.Write(writer, binary.LittleEndian, int16(value))
}

func putInts(writer io.Writer, value []int) error {
	if err := binary.Write(writer, binary.LittleEndian, int16(len(value))); err != nil {
		return err
	}
	for _, v := range value {
		if err := binary.Write(writer, binary.LittleEndian, int16(v)); err != nil {
			return err
		}
	}
	return nil
}

func putUInt16(writer io.Writer, value uint16) error {
	return binary.Write(writer, binary.LittleEndian, value)
}

func putUInt32(writer io.Writer, value uint32) error {
	return binary.Write(writer, binary.LittleEndian, value)
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
