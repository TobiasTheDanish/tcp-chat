package shared_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"testing"
	"unsafe"

	"github.com/TobiasTheDanish/tcp-chat/shared"
)

func TestPacketFromData(t *testing.T) {
	data := []byte("Hello world")
	dataLen := len(data)

	packet, err := shared.PacketFromData(data)
	if err != nil {
		t.Errorf("Didn't expect error, got: %s\n", err)
	}

	minor := packet.Header.Version & 0x0f
	major := packet.Header.Version >> 4

	if major != shared.MAJOR_VERSION || minor != shared.MINOR_VERSION {
		t.Errorf("Invalid version, expected: %d.%d, got: %d.%d", shared.MAJOR_VERSION, shared.MINOR_VERSION, major, minor)
	}

	if packet.Header.DataLength != uint16(dataLen) {
		t.Errorf("Expected dataLength header to be: %d, got: %d\n", dataLen, packet.Header.DataLength)
	}

	if !bytes.Equal(data, packet.Data) {
		t.Errorf("Expected data to be: \"%s\", got: \"%s\"\n", data, packet.Data)
	}

	tooBigData := make([]byte, shared.MAX_DATA_LEN+10)

	_, err = shared.PacketFromData(tooBigData)
	if err == nil {
		t.Errorf("Expected to error with too big data, but didn't")
	}

	byteArr := []byte{16, 0, 7, 84, 111, 98, 105, 97, 115, 62, 32, 16, 0, 6, 72, 101, 108, 108, 111, 10}
	packet, err = shared.PacketFromData(byteArr)

	if err != nil {
		t.Errorf("Didn't expect error, got: %s\n", err)
	}

	minor = packet.Header.Version & 0x0f
	major = packet.Header.Version >> 4

	if major != shared.MAJOR_VERSION || minor != shared.MINOR_VERSION {
		t.Errorf("Invalid version, expected: %d.%d, got: %d.%d", shared.MAJOR_VERSION, shared.MINOR_VERSION, major, minor)
	}

	if packet.Header.DataLength != uint16(len(byteArr)) {
		t.Errorf("Expected dataLength header to be: %d, got: %d\n", dataLen, packet.Header.DataLength)
	}

	if !bytes.Equal(byteArr, packet.Data) {
		t.Errorf("Expected data to be: \"%s\", got: \"%s\"\n", byteArr, packet.Data)
	}
}

func TestEncode(t *testing.T) {
	packet, _ := shared.PacketFromData([]byte("Hello world"))
	b := packet.Encode()

	if b[0] != packet.Header.Version {
		t.Errorf("Version mismatch between encoded: %d and provided packet: %d", b[0], packet.Header.Version)
	}

	if b[1] != byte(packet.Header.DataLength>>4) {
		t.Errorf("Datalength doesnt match!")
	}

	if b[2] != byte(packet.Header.DataLength&0x0f) {
		t.Errorf("Datalength doesnt match!")
	}

	data := b[3:]
	if strings.Compare(string(packet.Data), string(data)) != 0 {
		t.Errorf("Expected data to be: \"%v\", got: \"%v\"\n", packet.Data, data)
	}
}

func TestParsePacket(t *testing.T) {
	helloWorld := []byte("Hello world")
	packet, _ := shared.PacketFromData(helloWorld)

	reader := bufio.NewReader(bytes.NewReader(packet.Encode()))

	parsed, err := shared.ParsePacket(reader)
	if err != nil {
		t.Errorf("Didn't expect error, got: %s\n", err)
	}

	if parsed.Header.Version != packet.Header.Version {
		t.Errorf("Version mismatch between parsed and provided packet")
	}

	if parsed.Header.DataLength != uint16(packet.Header.DataLength) {
		t.Errorf("Expected dataLength header to be: %d, got: %d\n", packet.Header.DataLength, parsed.Header.DataLength)
	}

	if !bytes.Equal(packet.Data, parsed.Data) {
		t.Errorf("Expected data to be: \"%s\", got: \"%s\"\n", packet.Data, parsed.Data)
	}

	errPacket := *&shared.Packet{
		Header: shared.PacketHeader{
			Version:    0,
			DataLength: uint16(len(helloWorld)),
		},
		Data: helloWorld,
	}

	reader = bufio.NewReader(bytes.NewReader(errPacket.Encode()))
	parsed, err = shared.ParsePacket(reader)
	if err == nil {
		t.Errorf("Expected to error with invalid version, but didn't")
	}

	byteArr := []byte{16, 0, 7, 84, 111, 98, 105, 97, 115, 62, 32, 16, 0, 6, 72, 101, 108, 108, 111, 10}
	packet, err = shared.PacketFromData(byteArr)

	if err != nil {
		t.Errorf("Didn't expect error, got: %s\n", err)
	}

	reader = bufio.NewReader(bytes.NewReader(packet.Encode()))

	parsed, err = shared.ParsePacket(reader)
	if err != nil {
		t.Errorf("Didn't expect error, got: %s\n", err)
	}

	if parsed.Header.Version != packet.Header.Version {
		t.Errorf("Version mismatch between parsed and provided packet")
	}

	if parsed.Header.DataLength != uint16(packet.Header.DataLength) {
		t.Errorf("Expected dataLength header to be: %d, got: %d\n", packet.Header.DataLength, parsed.Header.DataLength)
	}

	if !bytes.Equal(packet.Data, parsed.Data) {
		t.Errorf("Expected data to be: \"%s\", got: \"%s\"\n", packet.Data, parsed.Data)
	}
}

type testStruct struct {
	name string
	age  uint32
}

func TestPacketFromStruct(t *testing.T) {
	test := testStruct{
		name: "Tobias",
		age:  3_000_000,
	}
	packet, err := shared.PacketFromType(test)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
		return
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
		return
	}

	expectedDataLength := uint16(len(test.name) + int(unsafe.Sizeof(test.age)))
	if packet.Header.DataLength != expectedDataLength {
		t.Errorf("Incorrect datalength, expected %d, got %d\n", expectedDataLength, packet.Header.DataLength)
	}

	index := 0
	for i := range len(test.name) {
		if test.name[i] != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, test.name[i], packet.Data[index]))
		}

		index += 1
	}
	for i := range unsafe.Sizeof(test.age) {
		val := byte(test.age >> (8 * (3 - i)))
		if val != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, val, packet.Data[index]))
		}

		index += 1
	}
}

type nestedTestStruct struct {
	data struct {
		fame int16
	}
}

func TestPacketFromNestedStruct(t *testing.T) {
	test := nestedTestStruct{
		data: struct{ fame int16 }{fame: -345},
	}
	packet, err := shared.PacketFromType(test)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
		return
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
		return
	}
	index := 0

	expectedDataLength := uint16(unsafe.Sizeof(test.data.fame))
	if packet.Header.DataLength != expectedDataLength {
		t.Errorf("Incorrect datalength, expected %d, got %d\n", expectedDataLength, packet.Header.DataLength)
	}

	dataFame := int16(packet.Data[index])<<8 | int16(packet.Data[index+1])
	if test.data.fame != dataFame {
		t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, test.data.fame, dataFame))
	}
	for i := range 2 {
		val := byte(test.data.fame >> (8 * (1 - i)))
		if val != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, val, packet.Data[index]))
		}

		index += 1
	}
}

func TestPacketFromByteSlice(t *testing.T) {
	data := []byte{1, 255, 42, 69}

	packet, err := shared.PacketFromType(data)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
		return
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
		return
	}
	index := 0

	expectedDataLength := uint16(len(data))
	if packet.Header.DataLength != expectedDataLength {
		t.Errorf("Incorrect datalength, expected %d, got %d\n", expectedDataLength, packet.Header.DataLength)
	}

	for i := range len(data) {
		if data[i] != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, data[i], packet.Data[index]))
		}
		index += 1
	}
}

func TestPacketFromArray(t *testing.T) {
	truthTable := [4]bool{false, true, true, false}

	packet, err := shared.PacketFromType(truthTable)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
		return
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
		return
	}
	index := 0

	expectedDataLength := uint16(len(truthTable))
	if packet.Header.DataLength != expectedDataLength {
		t.Errorf("Incorrect datalength, expected %d, got %d\n", expectedDataLength, packet.Header.DataLength)
	}

	for i := range len(truthTable) {
		expected := byte(0)
		if truthTable[i] {
			expected = 1
		}
		if expected != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, expected, packet.Data[index]))
		}
		index += 1
	}
}

func TestPacketFromUintPtr(t *testing.T) {
	data := uintptr(3200879)

	packet, err := shared.PacketFromType(data)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
		return
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
		return
	}
	index := 0

	ptrSize := unsafe.Sizeof(data)
	expectedDataLength := uint16(ptrSize)
	if packet.Header.DataLength != expectedDataLength {
		t.Errorf("Incorrect datalength, expected %d, got %d\n", expectedDataLength, packet.Header.DataLength)
	}

	for i := range ptrSize {
		val := byte(data >> (8 * ((ptrSize - 1) - i)))
		if val != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, val, packet.Data[index]))
		}

		index += 1
	}
}

func TestPacketFromPointer(t *testing.T) {
	data := testStruct{
		name: "Tobias",
		age:  256,
	}

	packet, err := shared.PacketFromType(&data)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
		return
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
		return
	}

	expectedDataLength := uint16(len(data.name) + int(unsafe.Sizeof(data.age)))
	if packet.Header.DataLength != expectedDataLength {
		t.Errorf("Incorrect datalength, expected %d, got %d\n", expectedDataLength, packet.Header.DataLength)
	}

	index := 0
	for i := range len(data.name) {
		if data.name[i] != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, data.name[i], packet.Data[index]))
		}

		index += 1
	}
	for i := range unsafe.Sizeof(data.age) {
		val := byte(data.age >> (8 * (3 - i)))
		if val != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, val, packet.Data[index]))
		}

		index += 1
	}
}

func TestPacketFromFloat32(t *testing.T) {
	data := float32(67000.123)

	packet, err := shared.PacketFromType(data)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
		return
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
		return
	}
	index := 0

	ptrSize := unsafe.Sizeof(data)
	expectedDataLength := uint16(ptrSize)
	if packet.Header.DataLength != expectedDataLength {
		t.Errorf("Incorrect datalength, expected %d, got %d\n", expectedDataLength, packet.Header.DataLength)
	}
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, float32(data))
	if err != nil {
		t.Errorf("binary.Write failed with error: %s", err)
		return
	}
	expected := buf.Bytes()

	actual := uint32(0)
	for i := range ptrSize {
		shiftVal := int32((8 * ((ptrSize - 1) - i)))
		val := expected[i]
		if val != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, val, packet.Data[index]))
		}
		actual = (uint32(val) << shiftVal) | actual

		index += 1
	}

	actualFloat := math.Float32frombits(actual)
	if data != actualFloat {
		t.Error(fmt.Sprintf("Error creating packet. Float value malformed. Expected %f, got %f\n", data, actualFloat))
	}
}

func TestPacketFromFloat64(t *testing.T) {
	data := float64(67000.123)

	packet, err := shared.PacketFromType(data)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
		return
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
		return
	}
	index := 0

	ptrSize := unsafe.Sizeof(data)
	expectedDataLength := uint16(ptrSize)
	if packet.Header.DataLength != expectedDataLength {
		t.Errorf("Incorrect datalength, expected %d, got %d\n", expectedDataLength, packet.Header.DataLength)
	}
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, data)
	if err != nil {
		t.Errorf("binary.Write failed with error: %s", err)
		return
	}
	expected := buf.Bytes()

	actual := uint64(0)
	for i := range ptrSize {
		shiftVal := int32((8 * ((ptrSize - 1) - i)))
		val := expected[i]
		if val != packet.Data[index] {
			t.Error(fmt.Sprintf("Error creating packet. At byte pos %d, expected %d, got %d", index, val, packet.Data[index]))
		}
		actual = (uint64(val) << shiftVal) | actual

		index += 1
	}

	actualFloat := math.Float64frombits(actual)
	if data != actualFloat {
		t.Error(fmt.Sprintf("Error creating packet. Float value malformed. Expected %f, got %f\n", data, actualFloat))
	}
}

func TestPacketFromComplex128(t *testing.T) {

}
