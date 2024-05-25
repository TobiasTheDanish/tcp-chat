package client

import (
	"bufio"
	"net"
)

func ResolveAddr(addr string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp4", addr)
}

func Connect(addr *net.TCPAddr) (*net.TCPConn, error) {
	return net.DialTCP("tcp", nil, addr)
}

func ConnectAndWrite(addr *net.TCPAddr, bytes []byte) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte("Hello TCP Server\n"))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ReadLine(conn *net.TCPConn) (string, error) {
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return data, nil
}
