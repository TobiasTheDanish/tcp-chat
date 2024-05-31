package shared

import (
	"bufio"
	"errors"
	"fmt"
	"reflect"
)

const (
	MAJOR_VERSION byte = 1
	MINOR_VERSION byte = 0
	MAX_DATA_LEN  int  = 65535
)

var (
	InvalidVersion = errors.New("Invalid version.")
	InvalidType    = errors.New("Invalid type.")
)

type PacketHeader struct {
	Version    byte
	DataLength uint16
}

type Packet struct {
	Header PacketHeader
	Data   []byte
}

func (p *Packet) Encode() []byte {
	bytes := make([]byte, p.Header.DataLength+3)

	bytes[0] = p.Header.Version
	bytes[1] = byte(p.Header.DataLength >> 4)
	bytes[2] = byte(p.Header.DataLength & 0x0f)

	for i := range p.Header.DataLength {
		bytes[i+3] = p.Data[i]
	}

	return bytes
}

func (p *Packet) VersionString() string {
	minor := p.Header.Version & 0x0f
	major := p.Header.Version >> 4
	return fmt.Sprintf("%d.%d", major, minor)
}

func PacketFromType(t interface{}) (*Packet, error) {
	rv := reflect.ValueOf(t)

	if rv.Kind() != reflect.Struct {
		return nil, errors.Join(InvalidType, errors.New(fmt.Sprintf("Could not create packet from type: %s", rv.Kind().String())))
	}

	data := getBytesFromStruct(rv)

	return PacketFromData(data)
}

func getBytesFromStruct(v reflect.Value) []byte {
	data := make([]byte, 0)

	for i := range v.NumField() {
		value := v.Field(i)
		kind := value.Kind()

		switch kind {
		case reflect.Struct:
			data = append(data, getBytesFromStruct(value)...)
		case reflect.String:
			b := []byte(value.String())
			fmt.Println("num bytes in string: ", len(b))
			data = append(data, b...)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			data = append(data, getBytesFromInt(value)...)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			data = append(data, getBytesFromUint(value)...)
		default:
			panic(errors.Join(InvalidType, errors.New(fmt.Sprintf("Could not get bytes from field of type: %s", value.Kind().String()))))
		}

	}
	fmt.Println("num bytes in struct: ", len(data))

	return data
}

func getBytesFromInt(v reflect.Value) []byte {
	size := v.Type().Bits()
	fmt.Println("num bytes in int: ", (size / 8))

	switch size {
	case 8:
		return []byte{byte(v.Int())}
	case 16:
		{
			val := int16(v.Int())
			return []byte{byte(val >> 8), byte(val & 0x0f)}
		}
	case 32:
		{
			val := int32(v.Int())
			return []byte{byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val & 0x0f)}
		}
	case 64:
		{
			val := v.Int()
			return []byte{byte(val >> 56), byte(val >> 48), byte(val >> 40), byte(val >> 32), byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val & 0x0f)}
		}
	default:
		panic(fmt.Sprintf("Unreachable size of int in bits %d", size))
	}
}

func getBytesFromUint(v reflect.Value) []byte {
	size := v.Type().Bits()

	fmt.Println("num bytes in uint: ", (size / 8))

	switch size {
	case 8:
		return []byte{byte(v.Uint())}
	case 16:
		{
			val := int16(v.Uint())
			return []byte{byte(val >> 8), byte(val & 0x0f)}
		}
	case 32:
		{
			val := int32(v.Uint())
			return []byte{byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val & 0x0f)}
		}
	case 64:
		{
			val := v.Uint()
			return []byte{byte(val >> 56), byte(val >> 48), byte(val >> 40), byte(val >> 32), byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val & 0x0f)}
		}
	default:
		panic(fmt.Sprintf("Unreachable size of uint in bits %d", size))
	}
}

func PacketFromData(data []byte) (*Packet, error) {
	// store major version in the 4 most significant bits
	// and minor version in the 4 least significant bits
	version := (MAJOR_VERSION << 4) | MINOR_VERSION

	if len(data) > MAX_DATA_LEN {
		return nil, errors.New("Data to long")
	}

	return &Packet{
		Header: PacketHeader{
			Version:    version,
			DataLength: uint16(len(data)),
		},
		Data: data,
	}, nil
}

func ParsePacket(reader *bufio.Reader) (*Packet, error) {
	headerBytes, err := readBytes(3, reader)
	if err != nil {
		return nil, err
	}
	version := headerBytes[0]
	minor := version & 0x0f
	major := version >> 4
	if minor != MINOR_VERSION || major != MAJOR_VERSION {
		return nil, errors.Join(InvalidVersion, errors.New(fmt.Sprintf("Expected: %d.%d, recieved: %d.%d", MAJOR_VERSION, MINOR_VERSION, major, minor)))
	}

	header := PacketHeader{
		Version:    version,
		DataLength: uint16((headerBytes[1] << 4) | headerBytes[2]),
	}

	data, err := readBytes(int(header.DataLength), reader)
	if err != nil {
		return nil, err
	}

	return &Packet{
		Header: header,
		Data:   data,
	}, nil
}

func readBytes(n int, reader *bufio.Reader) ([]byte, error) {
	if n < 0 {
		return nil, errors.New("Cannot read less than 0 bytes")
	}

	res := make([]byte, 0, n)
	if n == 0 {
		return res, nil
	}

	for range n {
		byte, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		res = append(res, byte)
	}

	return res, nil
}
