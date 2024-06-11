package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/TobiasTheDanish/tcp-chat/shared"
	"github.com/TobiasTheDanish/tcp-chat/tcp_server"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide port")
		os.Exit(1)
	}
	port := os.Args[1]

	server := tcp_server.Create(handleConnection)

	err := server.Start(port)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	for p := range server.PChan {
		var msg Message
		p.IntoType(&msg)

		fmt.Printf("%s: %s\n", msg.Username, msg.Msg)
		for _, conn := range server.Conns {
			conn.Write(p.Encode())
		}
	}
}

type Message struct {
	Username string
	Msg      string
}

func handleConnection(conn net.Conn, c chan *shared.Packet) {
	defer conn.Close()
	connIp := conn.LocalAddr().String()
	reader := bufio.NewReader(conn)

	welcome := Message{
		Username: "Server",
		Msg:      "Welcome! What is your username?",
	}
	packet, err := shared.PacketFromType(welcome)
	if err != nil {
		fmt.Println(err)
	}

	conn.Write(packet.Encode())
	usernamePacket, err := shared.ParsePacket(reader)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("Connection to %s closed\n", connIp)
		} else {
			fmt.Printf("ERROR: %s\n", err)
		}
		return
	}
	username := strings.Trim(string(usernamePacket.Data), "\r\n \t")

	for {
		// Read from the connection until a new line is send
		data, err := shared.ParsePacket(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection to %s closed\n", username)
			} else {
				fmt.Printf("ERROR: %s\n", err)
			}
			return
		}

		sentMessage := strings.Trim(string(data.Data), "\r\n \t")

		message := Message{
			Username: username,
			Msg:      sentMessage,
		}
		p, err := shared.PacketFromType(message)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		c <- p
	}
}
