package receiver

import (
	"fmt"
	"net"
)

func ConnectToSocket(address string) {
	_, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Error connecting to %s: %s\n", address, err.Error())
		return
	}
	// defer conn.Close()
	fmt.Printf("Connected to %s\n", address)

	// Use the connection here...
}
