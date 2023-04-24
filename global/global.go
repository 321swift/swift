package global

import (
	"fmt"
	"log"
	"net"
)

const (
	BroadcastPort = 51413
)

var BackendServerPort = 0

func GetAvailablePort() int {
	var serverPort = 5050
	for {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
		if err == nil {
			listener.Close()
			return serverPort
		}
		// Port is already in use, so try the next one
		log.Printf("Port %d already in use, trying next port\n", serverPort)
		serverPort++

	}
}
