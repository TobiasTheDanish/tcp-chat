package tcp_client

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/TobiasTheDanish/tcp-chat/shared"
)

type Client struct {
	addr *net.TCPAddr
	conn *net.TCPConn
}

func Connect(addr string) (*Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return &Client{addr: tcpAddr, conn: tcpConn}, nil
}

func (c *Client) ReadPacket() (*shared.Packet, error) {
	data, err := shared.ParsePacket(bufio.NewReader(c.conn))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ReadSend(r io.Reader) error {
	reader := bufio.NewReader(r)
	// fmt.Print("> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	return c.SendBytes([]byte(text))
}

func (c *Client) SendPacket(p *shared.Packet) error {
	_, err := c.conn.Write(p.Encode())
	return err
}

func (c *Client) SendBytes(bytes []byte) error {
	p, err := shared.PacketFromData(bytes)
	if err != nil {
		return err
	}

	return c.SendPacket(p)
}

func (c *Client) Listen() {
	for {
		// Read from the connection untill a new line is send
		packet, err := c.ReadPacket()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the data read from the connection to the terminal
		fmt.Print(string(packet.Data))
	}
}
