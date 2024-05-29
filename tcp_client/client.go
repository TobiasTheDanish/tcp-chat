package tcp_client

import (
	"bufio"
	"fmt"
	"io"
	"net"
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

func (c *Client) Write(bytes []byte) error {
	_, err := c.conn.Write(bytes)
	return err
}

func (c *Client) ReadLine() (string, error) {
	data, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return data, nil
}

func (c *Client) ReadWrite(r io.Reader) error {
	reader := bufio.NewReader(r)
	fmt.Print("> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	c.Write([]byte(text))

	return nil
}

func (c *Client) Listen() {
	for {
		// Read from the connection untill a new line is send
		data, err := c.ReadLine()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the data read from the connection to the terminal
		fmt.Print("> ", string(data))
	}
}
