package duo

import (
	"encoding/binary"
	"io"
	"log"
)

func readBool(conn io.Reader) bool {
	var value bool
	if err := binary.Read(conn, binary.LittleEndian, &value); err != nil {
		log.Fatal(err)
	}
	return value
}

func readChar(conn io.Reader) byte {
	var value byte
	if err := binary.Read(conn, binary.LittleEndian, &value); err != nil {
		log.Fatal(err)
	}
	return value
}

func readInt16(conn io.Reader) int {
	var value int16
	if err := binary.Read(conn, binary.LittleEndian, &value); err != nil {
		log.Fatal(err)
	}
	return int(value)
}

func readUInt32(conn io.Reader) int {
	var value uint32
	if err := binary.Read(conn, binary.LittleEndian, &value); err != nil {
		log.Fatal(err)
	}
	return int(value)
}

func readString(conn io.Reader) string {
	length := readInt16(conn)
	data := make([]byte, length)
	if err := binary.Read(conn, binary.LittleEndian, &data); err != nil {
		log.Fatal(err)
	}
	return string(data)
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
			if err := binary.Write(conn, binary.LittleEndian, int16(len(value))); err != nil {
				log.Fatal(err)
			}
			if err := binary.Write(conn, binary.LittleEndian, []byte(value)); err != nil {
				log.Fatal(err)
			}
		default:
			log.Println("writeSequence,", v)
			log.Fatalf("writeSequence, type %T unexpected", v)
		}
	}
}
