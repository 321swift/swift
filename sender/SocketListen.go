package sender

import (
	"fmt"
	"net"
)

func StartServer(portNumber int) {
	// Listen for incoming connections on port 12345
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		fmt.Println(err)
		fmt.Println("from sender @ listener")
		return
	}
	defer listener.Close()

	fmt.Println("Waiting for incoming connections...")

	// Accept an incoming connection
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Connection established!")
}
