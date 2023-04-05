package broadcaster

import (
	"fmt"
	"log"
	"net"
)

func startSocket(portNumber int) {
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		log.Fatal("error occured while creating socket: ", err)
	}
	defer listener.Close()

	fmt.Println("Successfully started listening on port:", portNumber)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		// Handle connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	fmt.Print("Connection established:", connection)
}
