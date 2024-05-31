package tcp_server

import (
	"fmt"
	"net"

	"github.com/TobiasTheDanish/tcp-chat/shared"
	"github.com/TobiasTheDanish/tcp-chat/tcp_server/internal/ip"
)

type ConnectionHandler func(net.Conn, chan *shared.Packet)

type Server struct {
	Conns []net.Conn
	PChan chan *shared.Packet
}

func Start(port string, handler ConnectionHandler) (*Server, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("0.0.0.0:%s", port))

	if err != nil {
		return nil, err
	}

	ip, err := ip.ExternalIP()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Connect here: %s:%s\n", ip, port)

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		return nil, err
	}

	s := &Server{
		Conns: make([]net.Conn, 0),
		PChan: make(chan *shared.Packet),
	}
	go s.accept(listener, handler)

	return s, nil
}

func (s *Server) accept(listener *net.TCPListener, handler ConnectionHandler) {
	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		s.Conns = append(s.Conns, conn)
		fmt.Printf("New connection from Local IP: %s\n", conn.LocalAddr().String())
		// Handle new connections in a Goroutine for concurrency
		go handler(conn, s.PChan)
	}
}
