package duo

import (
	"encoding/binary"
	"io"
	"log"
)

func readBool(conn io.Reader) (bool, error) {
	var value bool
	err := binary.Read(conn, binary.LittleEndian, &value)
	return value, err
}

func readChar(conn io.Reader) (byte, error) {
	var value byte
	err := binary.Read(conn, binary.LittleEndian, &value)
	return value, err
}

func readInt16(conn io.Reader) (int, error) {
	var value int16
	err := binary.Read(conn, binary.LittleEndian, &value)
	return int(value), err
}

func readUInt32(conn io.Reader) (int, error) {
	var value uint32
	err := binary.Read(conn, binary.LittleEndian, &value)
	return int(value), err
}

func readString(conn io.Reader) (string, error) {
	length, err := readInt16(conn)
	if err != nil {
		return "", err
	}
	data := make([]byte, length)
	err = binary.Read(conn, binary.LittleEndian, &data)
	return string(data), err
}

func writeSequence(conn io.Writer, values ...interface{}) {
	for i, v := range values {
		switch value := v.(type) {
		case int32:
			if i != 0 {
				log.Println("writeSequence,", v)
				log.Fatalf("writeSequence, type %T unexpected", v)
			}
			if err := binary.Write(conn, binary.LittleEndian, byte(value)); err != nil {
				log.Fatal(err)
			}
		case byte:
			if err := binary.Write(conn, binary.LittleEndian, value); err != nil {
				log.Fatal(err)
			}
		case int:
			if err := binary.Write(conn, binary.LittleEndian, int16(value)); err != nil {
				log.Fatal(err)
			}
		case string:
			if len(value) > 32767 {
				log.Fatalf("writeSequence, string too big (%d)", len(value))
			}
			if err := binary.Write(conn, binary.LittleEndian, int32(len(value))); err != nil {
				log.Fatal(err)
			}
			if err := binary.Write(conn, binary.LittleEndian, []byte(value)); err != nil {
				log.Fatal(err)
			}
		case []byte:
			if len(value) > 8388607 {
				log.Fatalf("writeSequence, string too big (%d)", len(value))
			}
			if err := binary.Write(conn, binary.LittleEndian, int32(len(value))); err != nil {
				log.Fatal(err)
			}
			if err := binary.Write(conn, binary.LittleEndian, value); err != nil {
				log.Fatal(err)
			}
		default:
			log.Println("writeSequence,", v)
			log.Fatalf("writeSequence, type %T unexpected", v)
		}
	}
}
