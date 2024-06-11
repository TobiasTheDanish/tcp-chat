package main

import (
	"fmt"
	"os"

	"github.com/TobiasTheDanish/tcp-chat/shared"
	"github.com/TobiasTheDanish/tcp-chat/tcp_client"
)

type Message struct {
	Username string
	Msg      string
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Please provide host:port to connect to")
		os.Exit(1)
	}

	// Resolve the string address to a TCP address
	tcpClient, err := tcp_client.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go tcpClient.Listen(messageHandler)

	for {
		err := tcpClient.ReadSend(os.Stdin)
		if err != nil {
			fmt.Println("ERROR reading from stdin: ", err)
			return
		}
	}
}

func messageHandler(p *shared.Packet) {
	var msg Message
	p.IntoType(&msg)

	fmt.Printf("%s: %s\n", msg.Username, msg.Msg)
}
