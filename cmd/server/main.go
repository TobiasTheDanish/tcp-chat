package main

import (
	"fmt"
	"os"

	"github.com/TobiasTheDanish/tcp-chat/tcp_server"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide port")
		os.Exit(1)
	}
	port := os.Args[1]

	err := tcp_server.Start(port)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
