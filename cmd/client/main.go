package main

import (
	"fmt"
	"os"

	"github.com/TobiasTheDanish/tcp-chat/tcp_client"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Please provide host:port to connect to")
		os.Exit(1)
	}

	// Resolve the string address to a TCP address
	tcpClient, err := client.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go tcpClient.Listen()

	for {
		err := tcpClient.ReadWrite(os.Stdin)
		if err != nil {
			fmt.Println("ERROR reading from stdin: ", err)
			return
		}
	}
}
