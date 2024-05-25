package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/TobiasTheDanish/tcp-chat/tcp_server/internal/ip"
)

func DisplayIp() {
	ip, err := ip.ExternalIP()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ip)
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	connIp := conn.LocalAddr().String()
	reader := bufio.NewReader(conn)

	conn.Write([]byte("What is your username?\n"))
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
		// Write back the same message to the client
		conn.Write([]byte(fmt.Sprintf("Hello %s\n", username)))

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

		// Print the data read from the connection to the terminal
		fmt.Printf("%s> %s", username, data)
	}
}
