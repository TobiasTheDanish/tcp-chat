package tcp_server

import (
	"fmt"
	"net"

	"github.com/TobiasTheDanish/tcp-chat/shared"
	"github.com/TobiasTheDanish/tcp-chat/tcp_server/internal/ip"
)

type ConnectionHandler func(net.Conn, chan *shared.Packet)

type Server struct {
	Conns   []net.Conn
	PChan   chan *shared.Packet
	handler ConnectionHandler
}

func Create(handler ConnectionHandler) Server {
	return Server{
		Conns:   make([]net.Conn, 0),
		PChan:   make(chan *shared.Packet),
		handler: handler,
	}
}

func (s *Server) Start(port string) error {
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
	go s.accept(listener)

	return nil
}

func (s *Server) accept(listener *net.TCPListener) {
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
		go s.handler(conn, s.PChan)
	}
}
