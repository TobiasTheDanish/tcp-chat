package main

import (
	"bufio"
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
	tcpAddr, err := client.ResolveAddr(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Connect to the address with tcp
	conn, err := client.Connect(tcpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		// Read from the connection untill a new line is send
		data, err := client.ReadLine(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the data read from the connection to the terminal
		fmt.Print("> ", string(data))

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("ERROR: ", err)
			return
		}

		conn.Write([]byte(text))
	}
}
