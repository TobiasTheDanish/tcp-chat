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

	server, err := tcp_server.Start(port, HandleConnection)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	for p := range server.PChan {
		fmt.Print(string(p.Data))
		for _, conn := range server.Conns {
			conn.Write(p.Encode())
		}
	}
}

func HandleConnection(conn net.Conn, c chan *shared.Packet) {
	defer conn.Close()
	connIp := conn.LocalAddr().String()
	reader := bufio.NewReader(conn)

	packet, err := shared.PacketFromData([]byte("Welcome! What is your username?\n"))
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

		sentMessage := string(data.Data)

		message := fmt.Sprintf("%s > %s", username, sentMessage)
		p, err := shared.PacketFromData([]byte(message))
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		c <- p
	}
}
