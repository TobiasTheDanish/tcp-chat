package shared

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"reflect"
)

const (
	MAJOR_VERSION byte = 1
	MINOR_VERSION byte = 0
	MAX_DATA_LEN  int  = 65535
)

var (
	InvalidVersion  = errors.New("Invalid version.")
	InvalidType     = errors.New("Invalid type.")
	UnsupportedType = errors.New("Unsupported type.")
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
	if t == nil {
		return PacketFromData([]byte{})
	}

	rv := reflect.ValueOf(t)

	data, err := getBytesFromValue(rv)
	if err != nil {
		return nil, err
	}

	return PacketFromData(data)
}

func getBytesFromValue(v reflect.Value) ([]byte, error) {
	var (
		err  error
		data []byte
	)
	kind := v.Kind()
	switch kind {
	case reflect.Struct:
		data, err = getBytesFromStruct(v)
	case reflect.Slice, reflect.Array:
		data, err = getBytesFromSliceOrArray(v)
	case reflect.Pointer, reflect.Interface:
		data, err = getBytesFromValue(v.Elem())
	case reflect.String:
		data = []byte(v.String())
	case reflect.Bool:
		data = getBytesFromBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		data = getBytesFromInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		data = getBytesFromUint(v)
	case reflect.Float32, reflect.Float64:
		data = getBytesFromFloat(v)
	case reflect.Complex64, reflect.Complex128:
		return nil, errors.Join(UnsupportedType, errors.New(fmt.Sprintf("Type %s is not currently supported", v.Kind().String())))
	case reflect.Invalid:
		return []byte{}, nil
	default:
		return data, errors.Join(InvalidType, errors.New(fmt.Sprintf("Could not get bytes from field of type: %s", v.Kind().String())))
	}

	return data, err
}

func getBytesFromSliceOrArray(v reflect.Value) ([]byte, error) {
	data := make([]byte, 0)

	for i := range v.Len() {
		b, err := getBytesFromValue(v.Index(i))
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}

	return data, nil
}

func getBytesFromStruct(v reflect.Value) ([]byte, error) {
	data := make([]byte, 0)

	for i := range v.NumField() {
		value := v.Field(i)

		b, err := getBytesFromValue(value)
		if err != nil {
			return nil, err
		}

		data = append(data, b...)
	}

	return data, nil
}

func getBytesFromBool(v reflect.Value) []byte {
	val := 0
	if v.Bool() {
		val = 1
	}

	return []byte{byte(val)}
}

func getBytesFromInt(v reflect.Value) []byte {
	size := v.Type().Bits()

	switch size {
	case 8:
		return []byte{byte(v.Int())}
	case 16:
		{
			val := int16(v.Int())
			return []byte{byte(val >> 8), byte(val)}
		}
	case 32:
		{
			val := int32(v.Int())
			return []byte{byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)}
		}
	case 64:
		{
			val := v.Int()
			return []byte{byte(val >> 56), byte(val >> 48), byte(val >> 40), byte(val >> 32), byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)}
		}
	default:
		panic(fmt.Sprintf("Unreachable size of int in bits %d", size))
	}
}

func getBytesFromUint(v reflect.Value) []byte {
	size := v.Type().Bits()

	switch size {
	case 8:
		return []byte{byte(v.Uint())}
	case 16:
		{
			val := uint16(v.Uint())
			return []byte{byte(val >> 8), byte(val)}
		}
	case 32:
		{
			val := uint32(v.Uint())
			return []byte{byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)}
		}
	case 64:
		{
			val := v.Uint()
			return []byte{byte(val >> 56), byte(val >> 48), byte(val >> 40), byte(val >> 32), byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)}
		}
	default:
		panic(fmt.Sprintf("Unreachable size of uint in bits %d", size))
	}
}

func getBytesFromFloat(v reflect.Value) []byte {
	size := v.Type().Bits()

	switch size {
	case 32:
		{
			val := math.Float32bits(float32(v.Float()))
			return []byte{byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)}
		}
	case 64:
		{
			val := math.Float64bits(v.Float())
			return []byte{byte(val >> 56), byte(val >> 48), byte(val >> 40), byte(val >> 32), byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)}
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
