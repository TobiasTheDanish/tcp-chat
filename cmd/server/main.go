package main

import (
	"fmt"
	"net"
	"os"

	"github.com/TobiasTheDanish/tcp-chat/tcp_server"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide port")
		os.Exit(1)
	}

	// Resolve the string address to a TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("0.0.0.0:%s", os.Args[1]))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server.DisplayIp()

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("New connection from Local IP: %s\n", conn.LocalAddr().String())
		// Handle new connections in a Goroutine for concurrency
		go server.HandleConnection(conn)
	}
}
