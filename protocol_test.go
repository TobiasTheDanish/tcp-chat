package shared_test

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

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
	data struct {
		fame int16
	}
}

func TestPacketFromType(t *testing.T) {
	test := testStruct{
		name: "Tobias",
		age:  14,
		data: struct {
			fame int16
		}{
			fame: -16,
		},
	}
	packet, err := shared.PacketFromType(test)
	if err != nil {
		t.Errorf("Did not expect error, but got: %s", err)
	}

	if packet == nil {
		t.Error("Expected packet but got nil")
	}

	fmt.Println("DataLength: ", packet.Header.DataLength)
}
