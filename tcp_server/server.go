package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/TobiasTheDanish/tcp-chat/tcp_server/internal/ip"
)

type server struct {
	conns []net.Conn
}

func Listen(port string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("0.0.0.0:%s", port))

	if err != nil {
		return err
	}

	ip, err := ip.ExternalIP()
	if err != nil {
		return err
	}

	fmt.Printf("Connect here: %s:%s\n", ip, port)

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		return err
	}

	s := server{conns: make([]net.Conn, 0)}
	channel := make(chan string)
	go s.accept(listener, channel)

	for str := range channel {
		fmt.Printf("%s", str)
		for i := range len(s.conns) {
			s.conns[i].Write([]byte(str))
		}
	}
	return nil
}

func (s *server) accept(listener *net.TCPListener, c chan string) {
	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		s.conns = append(s.conns, conn)
		fmt.Printf("New connection from Local IP: %s\n", conn.LocalAddr().String())
		// Handle new connections in a Goroutine for concurrency
		go HandleConnection(conn, c)
	}
}

func HandleConnection(conn net.Conn, c chan string) {
	defer conn.Close()
	connIp := conn.LocalAddr().String()
	reader := bufio.NewReader(conn)

	conn.Write([]byte("Welcome! What is your username?\n"))
	username, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Printf("Connection to %s closed\n", connIp)
		} else {
			fmt.Printf("ERROR: %s\n", err)
		}
		return
	}
	username = strings.Trim(username, "\r\n \t")

	for {
		// Read from the connection untill a new line is send
		data, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection to %s closed\n", username)
			} else {
				fmt.Printf("ERROR: %s\n", err)
			}
			return
		}

		message := fmt.Sprintf("%s> %s", username, data)
		// Print the data read from the connection to the terminal
		// fmt.Printf(message)
		c <- message
	}
}
