package tcp_server

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/TobiasTheDanish/tcp-chat/tcp_server"
)

func TestPacketFromData(t *testing.T) {
	data := []byte("Hello world")
	dataLen := len(data)

	packet, err := tcp_server.PacketFromData(data)
	if err != nil {
		t.Errorf("Didn't expect error, got: %s\n", err)
	}

	minor := packet.Header.Version & 0x0f
	major := packet.Header.Version >> 4

	if major != tcp_server.MAJOR_VERSION || minor != tcp_server.MINOR_VERSION {
		t.Errorf("Invalid version, expected: %d.%d, got: %d.%d", tcp_server.MAJOR_VERSION, tcp_server.MINOR_VERSION, major, minor)
	}

	if packet.Header.DataLength != uint16(dataLen) {
		t.Errorf("Expected dataLength header to be: %d, got: %d\n", dataLen, packet.Header.DataLength)
	}

	if !bytes.Equal(data, packet.Data) {
		t.Errorf("Expected data to be: \"%s\", got: \"%s\"\n", data, packet.Data)
	}

	tooBigData := make([]byte, tcp_server.MAX_DATA_LEN+10)

	_, err = tcp_server.PacketFromData(tooBigData)
	if err == nil {
		t.Errorf("Expected to error with too big data, but didn't")
	}
}

func TestEncode(t *testing.T) {
	packet, _ := tcp_server.PacketFromData([]byte("Hello world"))
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
	packet, _ := tcp_server.PacketFromData(helloWorld)

	reader := bufio.NewReader(bytes.NewReader(packet.Encode()))

	parsed, err := tcp_server.ParsePacket(reader)
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

	errPacket := *&tcp_server.Packet{
		Header: tcp_server.PacketHeader{
			Version:    0,
			DataLength: uint16(len(helloWorld)),
		},
		Data: helloWorld,
	}

	reader = bufio.NewReader(bytes.NewReader(errPacket.Encode()))
	parsed, err = tcp_server.ParsePacket(reader)
	if err == nil {
		t.Errorf("Expected to error with too big data, but didn't")
	}
}
