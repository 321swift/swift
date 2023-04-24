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

func GetAvailablePort(desiredPort int) int {
	for {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", desiredPort))
		if err == nil {
			listener.Close()
			return desiredPort
		}
		// Port is already in use, so try the next one
		log.Printf("Port %d already in use, trying next port\n", desiredPort)
		desiredPort++

	}
}

type Message struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}
