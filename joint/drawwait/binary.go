package drawwait

import (
	"encoding/binary"
	"fmt"
	"io"
)

func newParameters(variables ...any) []decoder {
	output := make([]decoder, 0, len(variables))
	for _, v := range variables {
		switch variable := v.(type) {
		case *byte:
			output = append(output, newAny(variable))
		case *[]byte:
			output = append(output, newBytes(variable))
		case *bool:
			output = append(output, newAny(variable))
		case *int:
			output = append(output, newInt(variable))
		case *[]int:
			output = append(output, newInts(variable))
		case *uint16:
			output = append(output, newAny(variable))
		case *uint32:
			output = append(output, newAny(variable))
		case *string:
			output = append(output, newString(variable))
		default:
			fmt.Printf("unknown type: %T", v)
		}
	}
	return output
}

type binaryBuffer struct {
	buffer []byte
	offset int
}

func newBuffer(size int) *binaryBuffer {
	return &binaryBuffer{
		buffer: make([]byte, size),
	}
}

func (b *binaryBuffer) Decode(reader io.Reader) (bool, error) {
	length, err := reader.Read(b.buffer[b.offset:])
	if err != nil {
		return false, fmt.Errorf("reader.Read: %w", err)
	}
	b.offset += length
	return b.offset >= len(b.buffer), nil
}

type binaryAny[T any] struct {
	value     *T
	anyBuffer *binaryBuffer
}

func newAny[T any](value *T) *binaryAny[T] {
	return &binaryAny[T]{
		value:     value,
		anyBuffer: newBuffer(binary.Size(*value)),
	}
}

func (b *binaryAny[T]) Decode(reader io.Reader) (bool, error) {
	ok, err := b.anyBuffer.Decode(reader)
	if ok && err == nil {
		if _, err := binary.Decode(b.anyBuffer.buffer, binary.LittleEndian, b.value); err != nil {
			return false, fmt.Errorf("binary.Decode: %w", err)
		}
	}
	return ok, err
}

type binaryInt struct {
	value      *int
	valueInt16 *binaryAny[int16]
}

func newInt(value *int) *binaryInt {
	return &binaryInt{
		value:      value,
		valueInt16: newAny(new(int16)),
	}
}

func (b *binaryInt) Decode(reader io.Reader) (bool, error) {
	ok, err := b.valueInt16.Decode(reader)
	if ok && err == nil {
		*b.value = int(*b.valueInt16.value)
	}
	return ok, err
}

type binaryString struct {
	value       *string
	valueInt32  *binaryAny[int32]
	valueBuffer *binaryBuffer
}

func newString(value *string) *binaryString {
	return &binaryString{
		value:      value,
		valueInt32: newAny(new(int32)),
	}
}

func (b *binaryString) Decode(reader io.Reader) (bool, error) {
	if b.valueBuffer == nil {
		ok, err := b.valueInt32.Decode(reader)
		if ok && err == nil {
			b.valueBuffer = newBuffer(int(*b.valueInt32.value))
			return false, nil
		}
		return ok, err
	}
	ok, err := b.valueBuffer.Decode(reader)
	if ok && err == nil {
		*b.value = string(b.valueBuffer.buffer)
	}
	return ok, err
}

type binaryInts struct {
	value       *[]int
	valueInt16  *binaryAny[int16]
	valueBuffer *binaryBuffer
}

func newInts(value *[]int) *binaryInts {
	return &binaryInts{
		value:      value,
		valueInt16: newAny(new(int16)),
	}
}

func (b *binaryInts) Decode(reader io.Reader) (bool, error) {
	if b.valueBuffer == nil {
		ok, err := b.valueInt16.Decode(reader)
		if ok && err == nil {
			b.valueBuffer = newBuffer(int(*b.valueInt16.value) * 2)
			return false, nil
		}
		return ok, err
	}
	ok, err := b.valueBuffer.Decode(reader)
	if ok && err == nil {
		valuesInt16 := make([]int16, int(*b.valueInt16.value))
		if _, err := binary.Decode(b.valueBuffer.buffer, binary.LittleEndian, &valuesInt16); err != nil {
			return false, fmt.Errorf("binary.Decode: %w", err)
		}
		values := make([]int, 0, int(*b.valueInt16.value))
		for i := range values {
			values[i] = int(valuesInt16[i])
		}
		b.value = &values
	}
	return ok, err
}

type binaryBytes struct {
	value       *[]byte
	valueInt32  *binaryAny[int32]
	valueBuffer *binaryBuffer
}

func newBytes(value *[]byte) *binaryBytes {
	return &binaryBytes{
		value:      value,
		valueInt32: newAny(new(int32)),
	}
}

func (b *binaryBytes) Decode(reader io.Reader) (bool, error) {
	if b.valueBuffer == nil {
		ok, err := b.valueInt32.Decode(reader)
		if ok && err == nil {
			b.valueBuffer = newBuffer(int(*b.valueInt32.value))
			return false, nil
		}
		return ok, err
	}
	ok, err := b.valueBuffer.Decode(reader)
	if ok && err == nil {
		*b.value = b.valueBuffer.buffer
	}
	return ok, nil
}
